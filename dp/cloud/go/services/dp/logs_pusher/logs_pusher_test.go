package logs_pusher

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubFluentdServer struct {
	expectedPayload string
	response        string
	t               *testing.T
}

func (s *stubFluentdServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p, _ := ioutil.ReadAll(r.Body)
	assert.JSONEq(s.t, s.expectedPayload, string(p))
	assert.Equal(s.t, "/dp", r.URL.Path)
}

func TestLogsPusher(t *testing.T) {
	testServer := httptest.NewServer(&stubFluentdServer{
		expectedPayload: `{"cbsd_serial_number":"cbsdId1234", "event_timestamp":12345, "fcc_id":"someFccId", "log_from":"SAS", "log_message":"some log message", "log_name":"someLogName", "log_to":"DP", "network_id":"someNetwork"}`,
		t:               t,
	})
	defer testServer.Close()
	log := &DPLog{
		EventTimestamp:   12345,
		LogFrom:          "SAS",
		LogTo:            "DP",
		LogName:          "someLogName",
		LogMessage:       "some log message",
		CbsdSerialNumber: "cbsdId1234",
		NetworkId:        "someNetwork",
		FccId:            "someFccId",
	}
	err := PushDPLog(context.Background(), log, fmt.Sprintf("%s/%s", testServer.URL, "dp"))
	assert.NoError(t, err)
}
