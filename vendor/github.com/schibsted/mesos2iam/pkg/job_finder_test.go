package pkg

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type DiscoverJobIdFromContainerTestSuite struct {
	suite.Suite
	pidFinder *CommandPidFinder
}

func (suite *DiscoverJobIdFromContainerTestSuite) SetupTest() {
	suite.pidFinder = NewPidFinder()
}

func (suite *DiscoverJobIdFromContainerTestSuite) TestExtractPortFromFuserCommand() {
	cm := fmt.Sprintf("35812/tcp:           28320")

	pid, err := suite.pidFinder.extractPidFromFuserCommand(cm)

	if err != nil {
		suite.T().Error(err)
	}
	assert.Equal(suite.T(), int32(28320), pid)
}
func (suite *DiscoverJobIdFromContainerTestSuite) TestExtractPortFromFuserCommandWithSeveralLines() {
	cm := fmt.Sprintf("35812/tcp:           28320\n       ")

	pid, err := suite.pidFinder.extractPidFromFuserCommand(cm)

	if err != nil {
		suite.T().Error(err)

	}
	assert.Equal(suite.T(), int32(28320), pid)
}
func (suite *DiscoverJobIdFromContainerTestSuite) TestExtractPortFromFuserCommandFailsIfComandOutputIsInvalid() {
	pid, err := suite.pidFinder.extractPidFromFuserCommand("8080/tcp")

	assert.Equal(suite.T(), int32(0), pid)

	if err == nil {
		suite.T().Error("Expected error didn't happen")
	}
}

func TestDiscoverJobIdFromContainerTestSuite(t *testing.T) {
	suite.Run(t, new(DiscoverJobIdFromContainerTestSuite))
}

func TestFindJobIdFromRequestWhenHostMode(t *testing.T) {
	req, err := http.NewRequest("GET", "/v2/credentials", nil)
	req.RemoteAddr = "52.52.52.52:10000"

	if err != nil {
		t.Fatal(err)
	}

	mockedPidFinder := getPidFinderMock()
	mockedRepository := getRepositoryMock()

	finder := ContainerJobFinder{
		mockedRepository,
		mockedPidFinder,
		"52.52.52.52",
		"TARDIS_SCHID=",
	}

	jobId, err := finder.FindJobIdFromRequest(req)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "4ea13548-caa8-48dc-af69-58a651d9fa3b", jobId)
	mockedPidFinder.AssertExpectations(t)
	mockedRepository.AssertExpectations(t)
}

func TestFindJobIdFromRequestWhenBridgeMode(t *testing.T) {
	req, err := http.NewRequest("GET", "/v2/credentials", nil)
	req.RemoteAddr = "172.17.0.2:10000"

	if err != nil {
		t.Fatal(err)
	}

	mockedPidFinder := getPidFinderMock()

	container := &docker.Container{
		Config: &docker.Config{
			Env: []string{"TARDIS_SCHID=4ea13548-caa8-48dc-af69-58a651d9fa3b"},
		},
		NetworkSettings: &docker.NetworkSettings{
			IPAddress: "172.17.0.2",
		},
	}
	mockedRepository := &MockedIpRepository{}
	mockedRepository.On("FindContainerUsingIp", "172.17.0.2").Return(container, nil)

	finder := ContainerJobFinder{
		repository: mockedRepository,
		pidFinder:  mockedPidFinder,
		idPrefix:   "TARDIS_SCHID=",
	}

	jobId, err := finder.FindJobIdFromRequest(req)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "4ea13548-caa8-48dc-af69-58a651d9fa3b", jobId)
	mockedRepository.AssertExpectations(t)
}

func getRepositoryMock() *MockedCommandRepository {
	container := &docker.Container{
		Config: &docker.Config{
			Env: []string{"TARDIS_SCHID=4ea13548-caa8-48dc-af69-58a651d9fa3b"},
		},
	}
	mockedRepository := &MockedCommandRepository{}
	mockedRepository.On("FindContainerUsingCommandPID", int32(800)).Return(container, nil)
	return mockedRepository
}

func getPidFinderMock() *MockedPidFinder {
	mockedPidFinder := &MockedPidFinder{}
	mockedPidFinder.On("GetCommandPidByPort", "10000").Return(800, nil)
	return mockedPidFinder
}

type MockedCommandRepository struct {
	mock.Mock
}

func (m *MockedCommandRepository) FindContainerUsingCommandPID(pid int32) (*docker.Container, error) {
	args := m.Called(pid)
	container := args.Get(0).(*docker.Container)
	return container, args.Error(1)
}
func (m *MockedCommandRepository) FindContainerUsingIp(ip string) (*docker.Container, error) {
	return nil, errors.New("Not implemented")
}

type MockedIpRepository struct {
	mock.Mock
}

func (m *MockedIpRepository) FindContainerUsingCommandPID(pid int32) (*docker.Container, error) {
	return nil, errors.New("Not implemented")
}
func (m *MockedIpRepository) FindContainerUsingIp(ip string) (*docker.Container, error) {
	args := m.Called(ip)
	container := args.Get(0).(*docker.Container)
	return container, args.Error(1)
}

type MockedPidFinder struct {
	mock.Mock
}

func (m *MockedPidFinder) GetCommandPidByPort(port string) (int32, error) {
	args := m.Called(port)
	return int32(args.Int(0)), args.Error(1)
}
