package http_test

import (
	"encoding/json"
	"github.com/aws/amazon-ecs-agent/agent/credentials"
	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	http_pkg "github.com/schibsted/mesos2iam/http"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSecurityRequestHandler(t *testing.T) {
	jobId := "4ea13548-caa8-48dc-af69-58a651d9fa3b"
	req, err := http.NewRequest("GET", "/v2/credentials", nil)
	if err != nil {
		t.Fatal(err)
	}

	mockedJobFinder := &MockedJobFinder{}
	mockedJobFinder.On("FindJobIdFromRequest", req).Return(jobId, nil)

	netClient := getMockNetClient("/credentials/" + jobId)

	securityRequestHandler := http_pkg.NewSecurityRequestHandler(mockedJobFinder, netClient, "http://fakeSmaugUrl", "TARDIS_SCHID")
	writer := httptest.NewRecorder()

	securityRequestHandler.ServeHTTP(writer, req)

	writer.Flush()
	body, _ := ioutil.ReadAll(writer.Body)
	assert.Equal(t, 200, writer.Code)
	assert.Equal(t, "{\"RoleArn\":\"roleArn\",\"AccessKeyId\":\"AccessKey\",\"SecretAccessKey\":\"Secret\",\"Token\":\"Token\",\"Expiration\":\"Expiration Date\"}", string(body))
}

func TestSecurityRequestHandlerInvalidJobId(t *testing.T) {
	req, err := http.NewRequest("GET", "/v2/credentials", nil)
	if err != nil {
		t.Fatal(err)
	}

	mockedJobFinder := &MockedJobFinder{}
	mockedJobFinder.On("FindJobIdFromRequest", req).Return("invalidJobid", nil)

	netClient := getMockNetClient("/credentials/")
	securityRequestHandler := http_pkg.NewSecurityRequestHandler(mockedJobFinder, netClient, "http://fakeSmaugUrl", "TARDIS_SCHID")
	writer := httptest.NewRecorder()

	securityRequestHandler.ServeHTTP(writer, req)

	writer.Flush()
	body, _ := ioutil.ReadAll(writer.Body)
	assert.Equal(t, 400, writer.Code)
	assert.Equal(t, "Invalid JobId in http request: invalidJobid", string(body))
}

type MockedJobFinder struct {
	mock.Mock
}

func (m *MockedJobFinder) FindJobIdFromRequest(r *http.Request) (string, error) {
	args := m.Called(r)
	return args.String(0), args.Error(1)
}

func getMockNetClient(requestUrl string) *http.Client {
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}
	netClient.Transport = newMockTransport(requestUrl)
	return netClient
}

type mockTransport struct {
	mockCredentials *credentials.IAMRoleCredentials
	requestUrl      string
}

func newMockTransport(requestUrl string) http.RoundTripper {
	creds := &credentials.IAMRoleCredentials{
		"id",
		"roleArn",
		"AccessKey",
		"Secret",
		"Token",
		"Expiration Date",
	}
	return &mockTransport{
		creds,
		requestUrl,
	}
}

// Implement http.RoundTripper
func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	expectedUrl := req.URL.EscapedPath()
	foundUrl := t.requestUrl
	if expectedUrl != foundUrl {
		return nil, errors.Errorf("Url does not match. Found: %s. Expected: %s", foundUrl, expectedUrl)
	}

	// Create mocked http.Response
	response := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusOK,
	}
	response.Header.Set("Content-Type", "application/json")

	responseBody, _ := json.Marshal(t.mockCredentials)

	response.Body = ioutil.NopCloser(strings.NewReader(string(responseBody[:])))
	return response, nil
}
