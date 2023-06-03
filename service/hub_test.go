package service

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/ship"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestHubSuite(t *testing.T) {
	suite.Run(t, new(HubSuite))
}

type testStruct struct {
	counter   int
	timeRange connectionInitiationDelayTimeRange
}

type HubSuite struct {
	suite.Suite

	serviceProvider *MockServiceProvider
	mdnsService     *MockMdnsService

	tests []testStruct
}

func (s *HubSuite) SetupSuite() {
	s.tests = []testStruct{
		{0, connectionInitiationDelayTimeRanges[0]},
		{1, connectionInitiationDelayTimeRanges[1]},
		{2, connectionInitiationDelayTimeRanges[2]},
		{3, connectionInitiationDelayTimeRanges[2]},
		{4, connectionInitiationDelayTimeRanges[2]},
		{5, connectionInitiationDelayTimeRanges[2]},
		{6, connectionInitiationDelayTimeRanges[2]},
		{7, connectionInitiationDelayTimeRanges[2]},
		{8, connectionInitiationDelayTimeRanges[2]},
		{9, connectionInitiationDelayTimeRanges[2]},
		{10, connectionInitiationDelayTimeRanges[2]},
	}

	ctrl := gomock.NewController(s.T())

	s.serviceProvider = NewMockServiceProvider(ctrl)
	s.serviceProvider.EXPECT().RemoteSKIConnected(gomock.Any()).AnyTimes()
	s.serviceProvider.EXPECT().ServiceShipIDUpdate(gomock.Any(), gomock.Any()).AnyTimes()
	s.serviceProvider.EXPECT().ServicePairingDetailUpdate(gomock.Any(), gomock.Any()).AnyTimes()

	s.mdnsService = NewMockMdnsService(ctrl)
	s.mdnsService.EXPECT().SetupMdnsService().AnyTimes()
	s.mdnsService.EXPECT().AnnounceMdnsEntry().AnyTimes()
	s.mdnsService.EXPECT().UnannounceMdnsEntry().AnyTimes()
	s.mdnsService.EXPECT().UnregisterMdnsSearch(gomock.Any()).AnyTimes()
}

func (s *HubSuite) Test_NewConnectionsHub() {
	ski := "12af9e"
	localService := NewServiceDetails(ski)
	configuration := &Configuration{
		interfaces: []string{"en0"},
	}

	hub := newConnectionsHub(s.serviceProvider, s.mdnsService, nil, configuration, localService)
	assert.NotNil(s.T(), hub)

	hub.start()
}

func (s *HubSuite) Test_IsRemoteSKIPaired() {
	sut := connectionsHub{
		connections:              make(map[string]*ship.ShipConnection),
		connectionAttemptCounter: make(map[string]int),
		remoteServices:           make(map[string]*ServiceDetails),
		serviceProvider:          s.serviceProvider,
	}
	ski := "test"

	paired := sut.IsRemoteServiceForSKIPaired(ski)
	assert.Equal(s.T(), false, paired)

	// mark it as connected, so mDNS is not triggered
	sut.connections[ski] = &ship.ShipConnection{}
	sut.RegisterRemoteSKI(ski, true)

	paired = sut.IsRemoteServiceForSKIPaired(ski)
	assert.Equal(s.T(), true, paired)

	// remove the connection, so the test doesn't try to close it
	delete(sut.connections, ski)
	sut.RegisterRemoteSKI(ski, false)
	paired = sut.IsRemoteServiceForSKIPaired(ski)
	assert.Equal(s.T(), false, paired)
}

func (s *HubSuite) Test_CheckRestartMdnsSearch() {
	sut := connectionsHub{
		connections:              make(map[string]*ship.ShipConnection),
		connectionAttemptCounter: make(map[string]int),
	}
	sut.checkRestartMdnsSearch()
	// Nothing to verify yet
}

func (s *HubSuite) Test_ReportServiceShipID() {
	sut := connectionsHub{
		serviceProvider: s.serviceProvider,
	}
	sut.ReportServiceShipID("", "")
	// Nothing to verify yet
}

func (s *HubSuite) Test_DisconnectSKI() {
	sut := connectionsHub{
		connections: make(map[string]*ship.ShipConnection),
	}
	ski := "test"
	sut.DisconnectSKI(ski, "none")
}

func (s *HubSuite) Test_RegisterConnection() {
	ski := "12af9e"
	localService := NewServiceDetails(ski)

	sut := connectionsHub{
		connections:  make(map[string]*ship.ShipConnection),
		mdns:         s.mdnsService,
		localService: localService,
	}

	ski = "test"
	con := &ship.ShipConnection{
		RemoteSKI: ski,
	}
	sut.registerConnection(con)
	assert.Equal(s.T(), 1, len(sut.connections))
	con = sut.connectionForSKI(ski)
	assert.NotNil(s.T(), con)
}

func (s *HubSuite) Test_IncreaseConnectionAttemptCounter() {

	// we just need a dummy for this test
	sut := connectionsHub{
		connectionAttemptCounter: make(map[string]int),
	}
	ski := "test"

	for _, test := range s.tests {
		sut.increaseConnectionAttemptCounter(ski)

		sut.muxConAttempt.Lock()
		counter, exists := sut.connectionAttemptCounter[ski]
		timeRange := connectionInitiationDelayTimeRanges[counter]
		sut.muxConAttempt.Unlock()

		assert.Equal(s.T(), true, exists)
		assert.Equal(s.T(), test.timeRange.min, timeRange.min)
		assert.Equal(s.T(), test.timeRange.max, timeRange.max)
	}
}

func (s *HubSuite) Test_RemoveConnectionAttemptCounter() {
	// we just need a dummy for this test
	sut := connectionsHub{
		connectionAttemptCounter: make(map[string]int),
	}
	ski := "test"

	sut.increaseConnectionAttemptCounter(ski)
	_, exists := sut.connectionAttemptCounter[ski]
	assert.Equal(s.T(), true, exists)

	sut.removeConnectionAttemptCounter(ski)
	_, exists = sut.connectionAttemptCounter[ski]
	assert.Equal(s.T(), false, exists)
}

func (s *HubSuite) Test_GetCurrentConnectionAttemptCounter() {
	// we just need a dummy for this test
	sut := connectionsHub{
		connectionAttemptCounter: make(map[string]int),
	}
	ski := "test"

	sut.increaseConnectionAttemptCounter(ski)
	_, exists := sut.connectionAttemptCounter[ski]
	assert.Equal(s.T(), exists, true)
	sut.increaseConnectionAttemptCounter(ski)

	value, exists := sut.getCurrentConnectionAttemptCounter(ski)
	assert.Equal(s.T(), 1, value)
	assert.Equal(s.T(), true, exists)
}

func (s *HubSuite) Test_GetConnectionInitiationDelayTime() {
	// we just need a dummy for this test
	ski := "12af9e"
	localService := NewServiceDetails(ski)
	sut := connectionsHub{
		localService:             localService,
		connectionAttemptCounter: make(map[string]int),
	}

	counter, duration := sut.getConnectionInitiationDelayTime(ski)
	assert.Equal(s.T(), 0, counter)
	assert.LessOrEqual(s.T(), float64(s.tests[counter].timeRange.min), float64(duration/time.Second))
	assert.GreaterOrEqual(s.T(), float64(s.tests[counter].timeRange.max), float64(duration/time.Second))
}

func (s *HubSuite) Test_ConnectionAttemptRunning() {
	// we just need a dummy for this test
	ski := "test"
	sut := connectionsHub{
		connectionAttemptRunning: make(map[string]bool),
	}

	sut.setConnectionAttemptRunning(ski, true)
	status := sut.isConnectionAttemptRunning(ski)
	assert.Equal(s.T(), true, status)
	sut.setConnectionAttemptRunning(ski, false)
	status = sut.isConnectionAttemptRunning(ski)
	assert.Equal(s.T(), false, status)
}
