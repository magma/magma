/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fbc/cwf/radius/config"
	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/monitoring"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	"fbc/lib/go/machine"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/expresswifi"
	"fbc/lib/go/radius/rfc2865"
	"fbc/lib/go/radius/rfc2866"
	"fbc/lib/go/radius/rfc2869"
	"fbc/lib/go/radius/ruckus"

	"github.com/mitchellh/mapstructure"
	"go.opencensus.io/tag"

	"github.com/donovanhide/eventsource"
	"go.uber.org/zap"
)

type (
	// SSEListener listens to Radius udp packets
	SSEListener struct {
		Listener
		Config           config.ListenerConfig
		Logger           *zap.Logger
		Extra            SSEListenerExtraConfig
		ShutdownSignal   chan bool
		ShutdownComplete chan bool
		EventSource      *eventsource.Stream
		Counters         monitoring.ListenerCounters
		ready            chan bool
	}

	// SSEEvent a struct depicting a single SSE Event coming from AAA
	SSEEvent struct {
		Code       int                      `json:"code"`
		Identifier byte                     `json:"identifier"`
		AVPs       map[string][]interface{} `json:"avps"`
		ProxyState string                   `json:"proxy_state"`
	}

	// SSEListenerExtraConfig extra config for SSE listener
	SSEListenerExtraConfig struct {
		Cookie         string `json:"cookie"`
		EventStreamURL string `json:"eventStreamURL"`
		ResponseURL    string `json:"responseURL"`
		SSEMacAddress  string `json:"sseMACAddress"`
	}

	wwwReq struct {
		Data string `json:"data,omitempty"`
	}
)

// NewSSEListener ...
func NewSSEListener() *SSEListener {
	return &SSEListener{
		ready: make(chan bool),
	}
}

// Init override
func (l *SSEListener) Init(
	server *Server,
	serverConfig config.ServerConfig,
	listenerConfig config.ListenerConfig,
	ctrs monitoring.ListenerCounters,
) error {
	// Parse configuration
	var cfg SSEListenerExtraConfig
	err := mapstructure.Decode(listenerConfig.Extra, &cfg)
	if err != nil {
		return err
	}

	if cfg.EventStreamURL == "" {
		return errors.New("You must configure 'eventStreamURL' in 'extra' section")
	}

	if cfg.ResponseURL == "" {
		return errors.New("You must configure 'responseURL' in 'extra' section")
	}

	l.Server = server
	l.Logger = server.logger
	l.Config = listenerConfig
	l.Counters = ctrs
	l.Extra = cfg
	l.ShutdownSignal = make(chan bool)
	l.ShutdownComplete = make(chan bool)
	return nil
}

// ListenAndServe override
func (l *SSEListener) ListenAndServe() error {
	// Build request
	url, err := url.Parse(l.Extra.EventStreamURL)
	if err != nil {
		l.ready <- false
		return err
	}

	// Get MAC Address to represent this machine
	if l.Extra.SSEMacAddress == "" {
		l.Extra.SSEMacAddress = machine.GetMachineMACAddressID()
	}

	l.Logger.Info(
		"Subscribing for CoA requests",
		zap.String("sse_mac_address", l.Extra.SSEMacAddress),
		zap.String("url", l.Extra.EventStreamURL),
	)

	req := http.Request{
		Method:     "GET",
		URL:        url,
		Proto:      "HTTP/2",
		ProtoMajor: 2,
		Header: http.Header{
			"sse-client-mac-address": []string{l.Extra.SSEMacAddress},
			"radius-packet-encoding": []string{"utf8/json"},
		},
	}

	if l.Extra.Cookie != "" {
		req.Header.Set("cookie", l.Extra.Cookie)
	}

	// Initiate stream
	es, err := eventsource.SubscribeWithRequest("", &req)
	if err != nil {
		l.Logger.Error("failed to start listenning on event source", zap.Error(err))
		l.ready <- false
		return err
	}
	l.EventSource = es

	// Start listener
	go func(l *SSEListener) {
		for {
			select {
			case <-l.ShutdownSignal:
				l.EventSource.Close()
				l.ShutdownComplete <- true
				return
			case evt := <-l.EventSource.Events:
				go l.handleCoaRequest(evt)
				break
			case err = <-l.EventSource.Errors:
				l.Logger.Error(err.Error())
				break
			}
		}
	}(l)

	// Signal ready
	l.ready <- true

	return nil
}

// Ready ...
func (l *SSEListener) Ready() chan bool {
	return l.ready
}

// Shutdown ...
func (l *SSEListener) Shutdown(ctx context.Context) error {
	l.ShutdownSignal <- true
	<-l.ShutdownComplete
	return nil
}

// Handlers
func (l *SSEListener) handleCoaRequest(evt eventsource.Event) {
	// Get SSE Event
	var sseEvent SSEEvent
	unmarshalCounter := monitoring.NewOperation(
		"sse_unmarshal",
		tag.Upsert(monitoring.ListenerTag, "sse"),
	)
	if err := json.Unmarshal([]byte(evt.Data()), &sseEvent); err != nil {
		l.Logger.Error("Failed to unmarshal request", zap.String("payload", evt.Data()))
		unmarshalCounter.Failure("unmarshal")
		return
	}

	// Verify and extract
	if radius.Code(sseEvent.Code) != radius.CodeCoARequest &&
		radius.Code(sseEvent.Code) != radius.CodeDisconnectRequest {
		l.Logger.Error(
			"Got invalid CoA radius code (expected: 43 CoA-Request or 40 Disconnect-Request)",
			zap.Int("value", sseEvent.Code),
		)
		unmarshalCounter.Failure("invalid_radius_code")
		return
	}
	unmarshalCounter.Success()

	requestCounter := l.Counters.StartRequest(radius.Code(sseEvent.Code))

	// Convert to RADIUS request
	correlationField := zap.Uint32("correlation", rand.Uint32())
	c := &modules.RequestContext{
		RequestID: correlationField.Integer,
		Logger:    l.Logger.With(correlationField),

		// Note: When using SSE Listener, we cannot reconstruct RADIUS session id
		//       (calling + called station IDs) from CoA Event.
		SessionID:      "",
		SessionStorage: nil,
	}

	var r *radius.Request = &radius.Request{
		Packet: &radius.Packet{
			Code:       radius.Code(sseEvent.Code),
			Identifier: sseEvent.Identifier,
			Attributes: radius.Attributes{},
			Secret:     []byte(l.Server.config.Secret),
		},
	}
	apply(r.Packet, sseEvent.AVPs)

	// Handle the request
	l.Logger.Debug("handling CoA request")
	radiusResponse, err := l.HandleRequest(c, r)
	if err != nil {
		l.Logger.Error("failed to handle SSE event", zap.Error(err))
		requestCounter.Failure("handling")
		return
	}
	l.Logger.Debug("CoA request handled successfully")

	if radiusResponse == nil {
		l.Logger.Warn("Got null response from handler. Dropping packet")
		requestCounter.Failure("nil_response")
		return
	}

	// Send CoA Response back to AAA
	responseJSON, err := json.Marshal(&wwwReq{
		Data: base64.StdEncoding.EncodeToString(radiusResponse.Raw),
	})
	if err != nil {
		l.Logger.Error("failed to serialize response body", zap.Error(err))
		requestCounter.Failure("marshal")
		return
	}

	httpRequest, err := http.NewRequest(
		"POST",
		l.Extra.ResponseURL,
		bytes.NewReader(responseJSON),
	)
	if err != nil {
		l.Logger.Error("failed to serialize response body", zap.Error(err))
		requestCounter.Failure("create_request")
		return
	}

	httpRequest.Header.Set("radius-packet-encoding", "application/json")
	client := &http.Client{}
	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		l.Logger.Error("failed sending CoA Response to AAA", zap.Error(err))
		requestCounter.Failure("send_response")
		return
	}

	requestCounter.GotResponse(radiusResponse.Code)
	defer httpResponse.Body.Close()
}

// TODO: This piece of code is used to convert JSON formatted CoA into
//       RADIUS packet. If we serialize this in the AAA server, and
//       send the serialized RADIUS packet, we can just pass it on here.
func getFirstAttrValue(attrs map[string][]interface{}, attr string) string {
	values, ok := attrs[attr]
	if !ok {
		return ""
	}

	if len(values) == 0 {
		return ""
	}

	strVal, ok := values[0].(string)
	if !ok {
		return ""
	}

	return strVal
}

func apply(p *radius.Packet, attributes map[string][]interface{}) {
	// Get parameters
	cpToken := getFirstAttrValue(attributes, "XWF-Captive-Portal-Token")
	if cpToken != "" {
		expresswifi.XWFCaptivePortalToken_Set(p, []byte(cpToken))
	}

	acctSessionID := getFirstAttrValue(attributes, "Acct-Session-Id")
	if acctSessionID != "" {
		rfc2866.AcctSessionID_Set(p, []byte(acctSessionID))
	}

	callingStationID := getFirstAttrValue(attributes, "Calling-Station-Id")
	if callingStationID != "" {
		rfc2865.CallingStationID_Set(p, []byte(callingStationID))
	}

	username := getFirstAttrValue(attributes, "User-Name")
	if username != "" {
		rfc2865.UserName_Set(p, []byte(username))
	}

	proxyState := getFirstAttrValue(attributes, "Proxy-State")
	if proxyState != "" {
		rfc2865.ProxyState_Set(p, []byte(proxyState))
	}

	acctInterimInterval := getFirstAttrValue(attributes, "Acct-Interim-Interval")
	if acctInterimInterval != "" {
		val, err := strconv.Atoi(acctInterimInterval)
		if err != nil {
			rfc2869.AcctInterimInterval_Set(p, rfc2869.AcctInterimInterval(val))
		}
	}

	nasIdentifier := getFirstAttrValue(attributes, "NAS-Identifier")
	if nasIdentifier != "" {
		rfc2865.NASIdentifier_Set(p, []byte(nasIdentifier))
	}

	nasIPAddress := getFirstAttrValue(attributes, "NAS-IP-Address")
	if nasIPAddress != "" {
		rfc2865.NASIPAddress_Set(p, []byte(nasIPAddress))
	}

	cptoken := getFirstAttrValue(attributes, "Ruckus-CP-Token")
	if cptoken != "" {
		ruckus.RuckusCPToken_Set(p, []byte(cptoken))
	}

	tcs, ok := attributes["XWF-Authorize-Traffic-Classes"]
	if ok {
		parsedTCs := []expresswifi.XWFAuthorizeTrafficClasses{}
		for _, tcVal := range tcs {
			tc, ok := tcVal.(map[string]interface{})
			if ok {
				parsedTCs = append(parsedTCs, expresswifi.XWFAuthorizeTrafficClasses{
					XWFAuthorizeClassName: tc["XWF-Authorize-Class-Name"].(string),
					XWFAuthorizeBytesLeft: uint64(tc["XWF-Authorize-Bytes-Left"].(float64)),
				})
			}
		}

		expresswifi.XWFAuthorizeTrafficClasses_Set(
			p,
			parsedTCs,
		)
	}

	tcs, ok = attributes["Ruckus-TC-Attr-Ids-With-Quota"]
	if ok {
		parsedTCs := []ruckus.RuckusTCAttrIdsWithQuota{}
		for _, tcVal := range tcs {
			tc, ok := tcVal.(map[string]interface{})
			if ok {
				parsedTCs = append(parsedTCs, ruckus.RuckusTCAttrIdsWithQuota{
					RuckusTCNameQuota: tc["Ruckus-TC-Name-Quota"].(string),
					RuckusTCQuota:     uint64(tc["Ruckus-TC-Quota"].(float64)),
				})
			}
		}

		ruckus.RuckusTCAttrIdsWithQuota_Set(
			p,
			parsedTCs,
		)
	}
}
