package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/DerAndereAndi/eebus-go/service/util"
	"github.com/DerAndereAndi/eebus-go/ship"
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/gorilla/websocket"
)

type ShipRole string

const (
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second // SHIP 4.2: ping interval + pong timeout
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 60 * time.Second // SHIP 4.2: ping interval

	// SHIP 9.2: Set maximum fragment length to 1024 bytes
	maxMessageSize = 1024

	ShipRoleServer ShipRole = "server"
	ShipRoleClient ShipRole = "client"
)

type ConnectionHandlerDelegate interface {
	requestUserTrustForService(service *ServiceDetails)
	shipIDUpdateForService(service *ServiceDetails)

	// new spine connection established, inform SPINE
	addRemoteDeviceConnection(ski, deviceCode string, deviceType model.DeviceTypeType, readC <-chan []byte, writeC chan<- []byte)

	// remove an existing connection from SPINE
	removeRemoteDeviceConnection(ski string)
}

// A ConnectionHandler handles websocket connections.
type ConnectionHandler struct {
	// The ship connection mode of this connection
	role ShipRole

	// The remote service
	remoteService *ServiceDetails

	// The local service
	localService *ServiceDetails

	// The actual websocket connection
	conn *websocket.Conn

	// internal connection id, so we can identify it uniquely instead of using the SKI which can be used in multiple connections
	connID uint64

	// The read channel for incoming messages
	readChannel chan []byte

	// The write channel for outgoing messages
	writeChannel chan []byte

	// The ship read channel for incoming messages
	shipReadChannel chan []byte

	// The ship write channel for outgoing messages
	shipWriteChannel chan []byte

	// The channel for handling local trust state update for the remote service
	shipTrustChannel chan bool

	// The current SHIP state
	smeState shipMessageExchangeState

	unregisterChannel  chan<- *ConnectionHandler
	connectionDelegate ConnectionHandlerDelegate
}

func newConnectionHandler(unregisterChannel chan<- *ConnectionHandler, connectionDelegate ConnectionHandlerDelegate, role ShipRole, localService, remoteService *ServiceDetails, conn *websocket.Conn) *ConnectionHandler {
	return &ConnectionHandler{
		unregisterChannel:  unregisterChannel,
		connectionDelegate: connectionDelegate,
		role:               role,
		localService:       localService,
		remoteService:      remoteService,
		conn:               conn,
	}
}

// Connection handler when the service initiates a connection to a remote service
func (c *ConnectionHandler) handleConnection() {
	if len(c.remoteService.SKI) == 0 {
		fmt.Println("SKI is not set")
		c.conn.Close()
		return
	}

	c.startup()
}

func (c *ConnectionHandler) startup() {
	c.readChannel = make(chan []byte, 1)      // Listen to incoming websocket messages
	c.writeChannel = make(chan []byte, 1)     // Send outgoing websocket messages
	c.shipReadChannel = make(chan []byte, 1)  // Listen to incoming ship messages
	c.shipWriteChannel = make(chan []byte, 1) // Send outgoing ship messages
	c.shipTrustChannel = make(chan bool, 1)   // Listen to trust state update

	go c.readPump()
	go c.writePump()

	go func() {
		if err := c.shipHandshake(c.remoteService.userTrust || len(c.remoteService.ShipID) > 0); err != nil {
			fmt.Println("SHIP handshake error: ", err)
			c.shutdown(false)
			return
		}

		// Report to SPINE local device about this remote device connection
		c.connectionDelegate.addRemoteDeviceConnection(c.remoteService.SKI, c.remoteService.ShipID, c.remoteService.deviceType, c.readChannel, c.writeChannel)
		c.shipMessageHandler()
	}()
}

// shutdown the connection and all internals
// may only invoked after startup() is invoked!
func (c *ConnectionHandler) shutdown(safeShutdown bool) {
	if c.smeState == smeComplete {
		c.connectionDelegate.removeRemoteDeviceConnection(c.remoteService.SKI)
	}

	c.unregisterChannel <- c

	if !isChannelClosed(c.readChannel) {
		close(c.readChannel)
	}

	if !isChannelClosed(c.writeChannel) {
		close(c.writeChannel)
	}

	if !isChannelClosed(c.shipReadChannel) {
		close(c.shipReadChannel)
	}

	if !isChannelClosed(c.shipWriteChannel) {
		close(c.shipWriteChannel)
	}

	if c.conn != nil {
		if c.smeState == smeComplete && safeShutdown {
			// close the SHIP connection according to the SHIP protocol
			c.shipClose()
		}

		c.conn.Close()
	}
}

func isChannelClosed[T any](ch <-chan T) bool {
	select {
	case <-ch:
		return false
	default:
		return true
	}
}

// writePump pumps messages from the writeChannel to the websocket connection
func (c *ConnectionHandler) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.writeChannel:
			if !ok {
				// The write channel is closed
				return
			}
			if err := c.sendSpineData(message); err != nil {
				return
			}

		case message, ok := <-c.shipWriteChannel:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The write channel has been closed
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages from the websocket connection into the read message channel
func (c *ConnectionHandler) readPump() {
	defer func() {
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		message, err := c.readWebsocketMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("Error reading message: ", err)
			}

			c.shutdown(false)
			return
		}

		// Check if this is a SHIP SME or SPINE message
		isShipMessage := false
		if c.smeState != smeComplete {
			isShipMessage = true
		} else {
			isShipMessage = bytes.Contains([]byte("datagram:"), message)
		}

		if isShipMessage {
			c.shipReadChannel <- message
		} else {
			_, jsonData := c.parseMessage(message, true)

			// Get the datagram from the message
			data := ship.ShipData{}
			if err := json.Unmarshal(jsonData, &data); err != nil {
				fmt.Println("Error unmarshalling message: ", err)
				continue
			}

			if data.Data.Payload == nil {
				fmt.Println("Received no valid payload")
				continue
			}
			c.readChannel <- []byte(data.Data.Payload)
		}
	}
}

// handles incoming ship specific messages outside of the handshake process
func (c *ConnectionHandler) shipMessageHandler() {
	for {
		select {
		case msg := <-c.shipReadChannel:
			// TODO: implement this
			// This should only be a close/abort message, right?
			fmt.Println(string(msg))
		}
	}
}

// read a message from the websocket connection
func (c *ConnectionHandler) readWebsocketMessage() ([]byte, error) {
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
func (c *ConnectionHandler) writeWebsocketMessage(message []byte) error {
	if c.conn == nil {
		return errors.New("Connection is not initialized")
	}
	c.shipWriteChannel <- message
	return nil
}

func (c *ConnectionHandler) sendSpineData(data []byte) error {
	payload := json.RawMessage(data)

	// Create the message
	shipMessage := ship.ShipData{
		Data: ship.DataType{
			Header: ship.HeaderType{
				ProtocolId: ship.ShipProtocolId,
			},
			Payload: payload,
		},
	}

	// Send the message
	return c.sendModel(ship.MsgTypeData, &shipMessage)
}

// send a json message for a provided model to the websocket connection
func (c *ConnectionHandler) sendModel(typ byte, model interface{}) error {
	msg, err := json.Marshal(model)
	if err != nil {
		return err
	}

	eebusMsg, err := util.JsonIntoEEBUSJson(msg)
	if err != nil {
		return err
	}

	fmt.Println("Send: ", string(eebusMsg))

	// Wrap the message into a binary message with the ship header
	shipMsg := []byte{typ}
	shipMsg = append(shipMsg, eebusMsg...)

	err = c.writeWebsocketMessage(shipMsg)
	if err != nil {
		return err
	}

	return nil
}

// enable jsonFormat if the return message is expected to be encoded in
// the eebus json format
// return the SHIP message type, the SHIP message and an error
func (c *ConnectionHandler) parseMessage(msg []byte, jsonFormat bool) (byte, []byte) {
	// Extract the SHIP header byte
	shipHeaderByte := msg[0]
	// remove the SHIP header byte from the message
	msg = msg[1:]

	if len(msg) > 1 {
		fmt.Println("Recv: ", string(msg))
	}

	if jsonFormat {
		return shipHeaderByte, util.JsonFromEEBUSJson(msg)
	}

	return shipHeaderByte, msg
}
