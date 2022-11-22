package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/service/util"
	"github.com/DerAndereAndi/eebus-go/ship"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/gorilla/websocket"
)

type shipDelegate interface {
	requestUserTrustForService(service *ServiceDetails)
	shipIDUpdateForService(service *ServiceDetails)

	// new spine connection established, inform SPINE
	addRemoteDeviceConnection(ski string, writeI spine.WriteMessageI) spine.ReadMessageI

	// remove an existing connection from SPINE
	removeRemoteDeviceConnection(ski string)
}

// A shipConnection handles websocket connections.
type shipConnection struct {
	// The ship connection mode of this connection
	role shipRole

	// The remote service
	remoteService *ServiceDetails

	// The local service
	localService *ServiceDetails

	// The actual websocket connection
	conn *websocket.Conn

	// Where to pass incoming SPINE messages to
	readHandler spine.ReadMessageI

	// The ship write channel for outgoing SHIP messages
	shipWriteChannel chan []byte

	// The connection was closed
	closeChannel chan struct{}

	// The channel for handling local trust state update for the remote service
	shipTrustChannel chan bool

	// The current SHIP state
	smeState shipMessageExchangeState

	// contains the error message if smeState is in state smeHandshakeError
	handshakeError error

	// handles timeouts for the current smeState
	handshakeTimer        *time.Timer
	handshakeTimerRunning bool

	// stores if the connection should be trusted right away
	connectionIsTrusted bool

	unregisterChannel  chan<- *shipConnection
	connectionDelegate shipDelegate

	// internal handling of closed connections
	isConnectionClosed bool

	mux          sync.Mutex
	shutdownOnce sync.Once
}

func newConnectionHandler(unregisterChannel chan<- *shipConnection, connectionDelegate shipDelegate, role shipRole, localService, remoteService *ServiceDetails, conn *websocket.Conn) *shipConnection {
	return &shipConnection{
		unregisterChannel:  unregisterChannel,
		connectionDelegate: connectionDelegate,
		role:               role,
		localService:       localService,
		remoteService:      remoteService,
		conn:               conn,
	}
}

func (c *shipConnection) startup() {
	c.shipWriteChannel = make(chan []byte, 1) // Send outgoing ship messages
	c.shipTrustChannel = make(chan bool, 1)   // Listen to trust state update
	c.closeChannel = make(chan struct{}, 1)   // Listen to close events

	c.handshakeTimer = time.NewTimer(time.Hour * 1)
	if !c.handshakeTimer.Stop() {
		<-c.handshakeTimer.C
	}

	go c.readShipPump()
	go c.writeShipPump()

	// if the user trusted this connection e.g. via the UI or if we already have a stored SHIP ID for this SKI
	c.connectionIsTrusted = c.remoteService.userTrust || len(c.remoteService.ShipID) > 0
	c.smeState = cmiStateInitStart
	c.handleShipState(false, nil)
}

// shutdown the connection and all internals
// may only invoked after startup() is invoked!
func (c *shipConnection) shutdown(safeShutdown bool) {
	c.shutdownOnce.Do(func() {
		if c.isConnectionClosed {
			return
		}

		smeState := c.smeState

		if smeState == smeComplete {
			c.connectionDelegate.removeRemoteDeviceConnection(c.remoteService.SKI)
		}

		c.unregisterChannel <- c

		c.mux.Lock()

		if !util.IsChannelClosed(c.shipWriteChannel) {
			close(c.shipWriteChannel)
			c.shipWriteChannel = nil
		}

		if !util.IsChannelClosed(c.shipTrustChannel) {
			close(c.shipTrustChannel)
			c.shipTrustChannel = nil
		}

		if !util.IsChannelClosed(c.closeChannel) {
			close(c.closeChannel)
			c.closeChannel = nil
		}

		c.mux.Unlock()

		if c.conn != nil {
			smeState = c.smeState
			if smeState == smeComplete && safeShutdown {
				// close the SHIP connection according to the SHIP protocol
				c.shipClose()
			}

			c.conn.Close()
		}

		c.isConnectionClosed = true
	})
}

// WriteMessageI interface implementation
func (c *shipConnection) WriteMessage(message []byte) {
	if err := c.sendSpineData(message); err != nil {
		logging.Log.Error(c.remoteService.SKI, "Error sending spine message: ", err)
		return
	}

}

// writePump pumps messages from the SPINE and SHIP writeChannels to the websocket connection
func (c *shipConnection) writeShipPump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.shutdown(false)
	}()

	for {
		select {
		case <-c.closeChannel:
			return

		case message, ok := <-c.shipWriteChannel:
			if c.isConnClosed() {
				return
			}

			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				logging.Log.Debug(c.remoteService.SKI, "Ship write channel closed")
				// The write channel has been closed
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
				logging.Log.Error(c.remoteService.SKI, "Error writing to websocket: ", err)
				return
			}
		case <-ticker.C:
			if c.isConnClosed() {
				return
			}
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logging.Log.Error(c.remoteService.SKI, "Error writing to websocket: ", err)
				return
			}
		}
	}
}

// readShipPump pumps messages from the websocket connection into the read message channel
func (c *shipConnection) readShipPump() {
	defer func() {
		c.shutdown(false)
	}()

	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		select {
		case <-c.closeChannel:
			return
		case <-c.handshakeTimer.C:
			if c.isConnClosed() {
				return
			}

			if !c.handshakeTimerRunning {
				continue
			}

			c.handleShipState(true, nil)
		case trust := <-c.shipTrustChannel:
			if trust {
				c.connectionIsTrusted = true
				c.handshakeHello_Trust()
			}
		default:
			if c.isConnClosed() {
				return
			}

			message, err := c.readWebsocketMessage()
			if err != nil {
				if c.isConnClosed() {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						logging.Log.Error(c.remoteService.SKI, "Error reading message: ", err)
					}

					return
				}

				logging.Log.Error(c.remoteService.SKI, "Websocket read error: ", err)
				c.shutdown(false)
				return
			}

			// Check if this is a SHIP SME or SPINE message
			if !bytes.Contains(message, []byte("datagram")) {
				c.handleShipState(false, message)
				continue
			}

			_, jsonData := c.parseMessage(message, true)

			// Get the datagram from the message
			data := ship.ShipData{}
			if err := json.Unmarshal(jsonData, &data); err != nil {
				logging.Log.Error(c.remoteService.SKI, "Error unmarshalling message: ", err)
				continue
			}

			if data.Data.Payload == nil {
				logging.Log.Error(c.remoteService.SKI, "Received no valid payload")
				continue
			}

			if c.readHandler == nil {
				return
			}
			_, _ = c.readHandler.ReadMessage([]byte(data.Data.Payload))
		}
	}
}

// read a message from the websocket connection
func (c *shipConnection) readWebsocketMessage() ([]byte, error) {
	if c.conn == nil {
		return nil, errors.New("Connection is not initialized")
	}

	msgType, b, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	if msgType != websocket.BinaryMessage {
		return nil, errors.New("Message is not a binary message")
	}

	if len(b) < 2 {
		return nil, fmt.Errorf("Invalid ship message length")
	}

	return b, nil
}

// write a message to the websocket connection
func (c *shipConnection) writeWebsocketMessage(message []byte) error {
	if c.conn == nil {
		return errors.New("Connection is not initialized")
	}

	c.shipWriteChannel <- message
	return nil
}

const payloadPlaceholder = `{"place":"holder"}`

func (c *shipConnection) transformSpineDataIntoShipJson(data []byte) ([]byte, error) {
	spineMsg, err := util.JsonIntoEEBUSJson(data)
	if err != nil {
		return nil, err
	}

	payload := json.RawMessage([]byte(spineMsg))

	// Workaround for the fact that SHIP payload is a json.RawMessage
	// which would also be transformed into an array element but it shouldn't
	// hence patching the payload into the message later after the SHIP
	// and SPINE model are transformed independently

	// Create the message
	shipMessage := ship.ShipData{
		Data: ship.DataType{
			Header: ship.HeaderType{
				ProtocolId: ship.ShipProtocolId,
			},
			Payload: json.RawMessage([]byte(payloadPlaceholder)),
		},
	}

	msg, err := json.Marshal(shipMessage)
	if err != nil {
		return nil, err
	}

	eebusMsg, err := util.JsonIntoEEBUSJson(msg)
	if err != nil {
		return nil, err
	}

	eebusMsg = strings.ReplaceAll(eebusMsg, `[`+payloadPlaceholder+`]`, string(payload))

	return []byte(eebusMsg), nil
}

func (c *shipConnection) sendSpineData(data []byte) error {
	eebusMsg, err := c.transformSpineDataIntoShipJson(data)
	if err != nil {
		return err
	}

	logging.Log.Trace("Send:", c.remoteService.SKI, string(eebusMsg))

	// Wrap the message into a binary message with the ship header
	shipMsg := []byte{ship.MsgTypeData}
	shipMsg = append(shipMsg, eebusMsg...)

	err = c.writeWebsocketMessage(shipMsg)
	if err != nil {
		logging.Log.Error("Error sending message: ", err)
		return err
	}

	return nil
}

// send a json message for a provided model to the websocket connection
func (c *shipConnection) sendShipModel(typ byte, model interface{}) error {
	shipMsg, err := c.shipMessage(typ, model)
	if err != nil {
		return err
	}

	err = c.writeWebsocketMessage(shipMsg)
	if err != nil {
		return err
	}

	return nil
}

func (c *shipConnection) shipMessage(typ byte, model interface{}) ([]byte, error) {
	msg, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}

	eebusMsg, err := util.JsonIntoEEBUSJson(msg)
	if err != nil {
		return nil, err
	}

	logging.Log.Trace("Send:", c.remoteService.SKI, string(eebusMsg))

	// Wrap the message into a binary message with the ship header
	shipMsg := []byte{typ}
	shipMsg = append(shipMsg, eebusMsg...)

	return shipMsg, nil
}

// enable jsonFormat if the return message is expected to be encoded in
// the eebus json format
// return the SHIP message type, the SHIP message and an error
func (c *shipConnection) parseMessage(msg []byte, jsonFormat bool) (byte, []byte) {
	// Extract the SHIP header byte
	shipHeaderByte := msg[0]
	// remove the SHIP header byte from the message
	msg = msg[1:]

	if len(msg) > 1 && c.smeState == smeComplete {
		logging.Log.Trace("Recv:", c.remoteService.SKI, string(msg))
	}

	if jsonFormat {
		return shipHeaderByte, util.JsonFromEEBUSJson(msg)
	}

	return shipHeaderByte, msg
}

func (c *shipConnection) isConnClosed() bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.isConnectionClosed
}