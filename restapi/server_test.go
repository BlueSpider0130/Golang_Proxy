package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gitlab.com/devskiller-tasks/messaging-app-golang/smsproxy"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const shutdownDuration = 5 * time.Second
const port = 80

func TestSmsAcceptedToSend(t *testing.T) {
	// given
	server := NewServer(port)
	server.BindEndpoints()
	defer server.Stop(shutdownDuration)
	content, err := json.Marshal(SendSmsRequest{
		PhoneNumber: "123",
		Content:     "Some content",
	})
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPost, "/sms", bytes.NewReader(content))
	require.NoError(t, err)
	recorder := httptest.NewRecorder()

	// when
	server.serveMux.ServeHTTP(recorder, request)

	// then
	require.EqualValues(t, http.StatusAccepted, recorder.Code)
	messageID, err := readMessageID(recorder.Body)
	require.NoError(t, err)

	// and given
	request, err = http.NewRequest(http.MethodGet, fmt.Sprintf("/sms/%s", messageID.String()), bytes.NewReader(content))
	require.NoError(t, err)
	recorder = httptest.NewRecorder()

	// when
	server.serveMux.ServeHTTP(recorder, request)

	// then
	require.EqualValues(t, http.StatusOK, recorder.Code)
	var statusHttpResponse SmsStatusResponse
	decoder := json.NewDecoder(recorder.Body)
	err = decoder.Decode(&statusHttpResponse)
	require.NoError(t, err)
	require.EqualValues(t, smsproxy.Accepted, statusHttpResponse.Status)
}

func TestURINotFound(t *testing.T) {
	tests := []struct {
		url string
	}{
		{"/das"},
		{"/sms/asdads/adsdasads"},
		{"/sms/"},
	}

	server := NewServer(port)
	server.BindEndpoints()
	defer server.Stop(shutdownDuration)

	for _, test := range tests {
		t.Run(fmt.Sprintf("testing url : %s", test.url), func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, test.url, nil)
			require.NoError(t, err)
			result := httptest.NewRecorder()
			server.serveMux.ServeHTTP(result, request)
			response := result.Result()
			require.Equal(t, http.StatusNotFound, response.StatusCode)
		})
	}
}

func TestErrorSendingSms(t *testing.T) {
	// given
	server := NewServer(port)
	server.BindEndpoints()
	server.smsProxyService = newErrorSmsProxy("some error")
	defer server.Stop(shutdownDuration)

	content, err := json.Marshal(SendSmsRequest{
		PhoneNumber: "123",
		Content:     "Some content",
	})
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPost, "/sms", bytes.NewReader(content))
	require.NoError(t, err)
	recorder := httptest.NewRecorder()

	// when
	server.serveMux.ServeHTTP(recorder, request)

	// then
	require.EqualValues(t, http.StatusInternalServerError, recorder.Code)

	responseError, err := readErrorContent(recorder.Body)
	require.NoError(t, err)
	require.Equal(t, "some error", responseError)
}

func TestBadRequestSendingWrongDataAsSMS(t *testing.T) {
	// given
	server := NewServer(port)
	server.BindEndpoints()
	defer server.Stop(shutdownDuration)

	content, err := json.Marshal(struct {
		Field string
	}{Field: "something"})
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPost, "/sms", bytes.NewReader(content))
	require.NoError(t, err)
	recorder := httptest.NewRecorder()

	// when
	server.serveMux.ServeHTTP(recorder, request)

	// then
	require.EqualValues(t, http.StatusBadRequest, recorder.Code)
}

func readMessageID(body io.Reader) (uuid.UUID, error) {
	var httpResponse SmsSendResponse
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&httpResponse)
	if err != nil {
		return uuid.UUID{}, err
	}
	messageID, err := uuid.Parse(httpResponse.ID)
	if err != nil {
		return uuid.UUID{}, err
	}
	return messageID, nil
}

func readErrorContent(body io.Reader) (string, error) {
	var x struct{ Error string }
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&x)
	if err != nil {
		return "", err
	}
	return x.Error, nil
}
