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

package diameter

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/fiorix/go-diameter/v4/diam/sm"
	"github.com/golang/glog"
)

// KeyAndAnswer wraps the information to be returned to an AnswerHandler.
// When an answer is received, the handler should return a parsed answer
// (struct form) and the key that maps it to a request
type KeyAndAnswer struct{ Answer, Key interface{} }

// AnswerHandler is called when an answer is received. The handler is responsible
// for parsing a raw message into a usable message and then returning it. The return
// value is used to untrack the request and send the answer to the given channel
// Input: message received
// Output: struct containing the related request key and parsed answer
type AnswerHandler func(message *diam.Message) KeyAndAnswer

// Client is a wrapper around a sm.Client that handles connection management,
// request tracking, and configuration. Using this, the application should not
// know anything about the underlying diameter connection
type Client struct {
	mux            *sm.StateMachine
	smClient       *sm.Client
	connMan        *ConnectionManager
	requestTracker *RequestTracker
	cfg            *DiameterClientConfig
	originStateID  uint32
}

// OriginRealm returns client's config Realm
func (c *Client) OriginRealm() string {
	if c != nil && c.cfg != nil && len(c.cfg.Realm) > 0 {
		return c.cfg.Realm
	}
	return "magma"
}

// OriginHost returns client's config Host
func (c *Client) OriginHost() string {
	if c != nil && c.cfg != nil && len(c.cfg.Host) > 0 {
		return c.cfg.Host
	}
	return "magma"
}

// OriginStateID returns client's Origin-State-ID
func (c *Client) OriginStateID() uint32 {
	if c != nil {
		return c.originStateID
	}
	return 0
}

// ServiceContextId returns client's config ServiceContextId
func (c *Client) ServiceContextId() string {
	if c != nil && c.cfg != nil && len(c.cfg.ServiceContextId) > 0 {
		return c.cfg.ServiceContextId
	}
	return ServiceContextIDDefault
}

// NewClient creates a new client based on the config passed.
// Input: clientCfg containing relevant diameter settings
func NewClient(clientCfg *DiameterClientConfig) *Client {
	originStateID := uint32(time.Now().Unix())
	mux := sm.New(&sm.Settings{
		OriginHost:       datatype.DiameterIdentity(clientCfg.Host),
		OriginRealm:      datatype.DiameterIdentity(clientCfg.Realm),
		VendorID:         datatype.Unsigned32(Vendor3GPP),
		ProductName:      datatype.UTF8String(clientCfg.ProductName),
		OriginStateID:    datatype.Unsigned32(originStateID),
		FirmwareRevision: 1,
		HostIPAddress:    datatype.Address(net.ParseIP("127.0.0.1")),
	})

	appIdAvp := diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(clientCfg.AppID))

	var authAppIdAvps []*diam.AVP
	if clientCfg.AuthAppID != 0 {
		authAppIdAvps = []*diam.AVP{
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(clientCfg.AuthAppID))}
	} else {
		authAppIdAvps = []*diam.AVP{appIdAvp}
	}

	vendorSpecificApplicationIDs := getVendorSpecificApplicationIDAVPs(clientCfg, appIdAvp)

	// Add the standard 3gpp vendor ID
	vendorSpecificApplicationIDs = append(vendorSpecificApplicationIDs,
		diam.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				appIdAvp,
				diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(Vendor3GPP)),
			},
		}))

	cli := &sm.Client{
		Dict:               dict.Default,
		Handler:            mux,
		MaxRetransmits:     clientCfg.Retransmits,
		RetransmitInterval: time.Second,
		EnableWatchdog:     clientCfg.WatchdogInterval > 0,
		WatchdogInterval:   time.Second * time.Duration(clientCfg.WatchdogInterval),
		SupportedVendorID: []*diam.AVP{
			diam.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, datatype.Unsigned32(Vendor3GPP)),
		},
		AuthApplicationID:           authAppIdAvps,
		VendorSpecificApplicationID: vendorSpecificApplicationIDs,
	}
	client := &Client{
		mux:            mux,
		smClient:       cli,
		connMan:        NewConnectionManager(),
		requestTracker: NewRequestTracker(),
		cfg:            clientCfg,
		originStateID:  originStateID,
	}
	go client.handleErrors(mux.ErrorReports())
	return client
}

const (
	connectionRecoveryBackoff    = time.Millisecond * 100
	connectionRecoveryMaxRetries = 10
)

// handleErrors logs errors received during transmission & tries to recover errored out connections
func (client *Client) handleErrors(ec <-chan *diam.ErrorReport) {
	cm := client.connMan
	if cm == nil {
		glog.Error("<nil> Connection Manager")
	}

	for err := range ec {
		if err != nil && err.Conn != nil && client != nil {
			dc := err.Conn
			connStr := dc.LocalAddr().String() + "->" + dc.RemoteAddr().String()
			glog.Errorf("diameter connection %s error: %v", connStr, err)
			if cm == nil {
				continue
			}
			conn := cm.Find(dc)
			if conn != nil {
				// first, try to close the existing connection if it hasn't been reestablished yet
				conn.destroyConnection(dc)
				// recover connection in a dedicated routine, it can take long time
				// getDiamConnection(0 will just return success if the connection was already recovered
				go func(conn *Connection) {
					backoff := connectionRecoveryBackoff
					for retry := 0; retry < connectionRecoveryMaxRetries; retry++ {
						_, _, retryErr := conn.getDiamConnection()
						if retryErr == nil {
							glog.Infof("diameter connection %s is successfully recovered", connStr)
							return
						}
						glog.Errorf(
							"failed to recover diameter connection %s; attempt #%d: %v",
							connStr, retry, retryErr)

						time.Sleep(backoff)
						backoff *= 2
					}
				}(conn)
			} else {
				glog.Errorf("cannot find connection for %s", connStr)
			}
		} else {
			glog.Error(err)
		}
	}
}

// BeginConnection attempts to begin a new connection with the server
func (client *Client) BeginConnection(server *DiameterServerConfig) error {
	if client.connMan == nil {
		err := fmt.Errorf("No connection manager to initiate connection with")
		glog.Error(err)
		return err
	}
	_, err := client.connMan.GetConnection(client.smClient, server)
	if err != nil {
		glog.Error(err)
	}
	return err
}

func (client *Client) Retries() uint {
	if client != nil && client.cfg != nil {
		return client.cfg.RetryCount
	}
	return 0
}

// EnableConnectionCreation enables the connection manager to create new connections
func (client *Client) EnableConnectionCreation() {
	if client.connMan == nil {
		glog.Errorf("No connection manager to enable connection creation with")
		return
	}
	client.connMan.Enable()
}

// DisableConnectionCreation closes all existing connections and disables the
// connection manager to create new connections for the period of time specified
func (client *Client) DisableConnectionCreation(period time.Duration) {
	if client.connMan == nil {
		glog.Errorf("No connection manager to disable connection creation with")
		return
	}
	client.connMan.DisableFor(period)
}

// SendRequest sends a diameter request message to the given server and sends
// back the answer on the given channel. A key is required to identify the
// corresponding answer. Additionally, SendRequest will add the OriginHost/Realm
// AVPs to the message because they are mandatory for all requests
// Input: server - cfg containing info on what server to send to
// 				done - channel to send the answer to when received
//				message - request to send
//				key - something to uniquely identify the request
// Output: error if message sending failed, nil otherwise
func (client *Client) SendRequest(
	server *DiameterServerConfig,
	done chan interface{},
	message *diam.Message,
	key interface{},
) error {
	client.requestTracker.RegisterRequest(key, done)
	conn, err := client.connMan.GetConnection(client.smClient, server)
	if err == nil {
		m := client.AddOriginAVPsToMessage(message)
		err = conn.SendRequestToServer(m, client.cfg.RetryCount, server)
		if err != nil {
			client.requestTracker.DeregisterRequest(key)
		}
	}
	return err
}

// AddOriginAVPsToMessage adds the host/realm to the message
func (client *Client) AddOriginAVPsToMessage(message *diam.Message) *diam.Message {
	message.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity(client.cfg.Host))
	message.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity(client.cfg.Realm))
	// add Origin-State-ID
	if client.originStateID != 0 {
		originAVP, err := message.FindAVP(avp.OriginStateID, 0)
		if err != nil {
			message.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(client.originStateID))
		} else if originAVP != nil {
			// apply new originStateID
			originAVP.Data = datatype.Unsigned32(client.originStateID)
		}
	}
	return message
}

// IgnoreAnswer untracks a request if the application, say, times out
// Input: key identifying request
func (client *Client) IgnoreAnswer(key interface{}) {
	client.requestTracker.DeregisterRequest(key)
}

// RegisterAnswerHandlerForAppID registers a function to be called when an answer message
// matching the given command is received. The AnswerHandler is responsible for
// parsing the diameter message into something usable and extracting a key to identify
// the corresponding request
// Input: command - the diameter code for the command (like diam.CreditControl)
// 				handler - the function to call when a message is received
func (client *Client) RegisterAnswerHandlerForAppID(command uint32, appID uint32, handler AnswerHandler) {
	index := diam.CommandIndex{AppID: appID, Code: command, Request: false}
	muxHandler := diam.HandlerFunc(func(c diam.Conn, m *diam.Message) {
		answerKey := handler(m)
		if answerKey.Key == nil {
			return
		}
		doneChan := client.requestTracker.DeregisterRequest(answerKey.Key)
		doneChan <- answerKey.Answer
	})
	client.mux.HandleIdx(index, muxHandler)
}

// RegisterAnswerHandler registers a function to be called when an answer message
// matching the given command is received. The AnswerHandler is responsible for
// parsing the diameter message into something usable and extracting a key to identify
// the corresponding request
// Input: command - the diameter code for the command (like diam.CreditControl)
// 				handler - the function to call when a message is received
func (client *Client) RegisterAnswerHandler(command uint32, handler AnswerHandler) {
	client.RegisterAnswerHandlerForAppID(command, client.cfg.AppID, handler)
}

// RegisterRequestHandlerForAppID registers a function to be called when a request message
// matching the command is received. The RequestHandler is responsible for parsing
// the diameter message, taking any actions required, and sending a response back
// through the responder argument of the handler. This responder is given so that
// the response can happen asynchonously in another go routine.
// Input: command - the diameter code for the command (like diam.CreditControl)
//				handler - the function to call when a message is received
func (client *Client) RegisterRequestHandlerForAppID(command uint32, appID uint32, handler diam.HandlerFunc) {
	client.mux.HandleIdx(diam.CommandIndex{AppID: appID, Code: command, Request: true}, handler)
}

// RegisterHandler registers diameter handler to be used for given command and app
func (client *Client) RegisterHandler(command uint32, appID uint32, request bool, handler diam.Handler) {
	client.mux.HandleIdx(diam.CommandIndex{AppID: appID, Code: command, Request: request}, handler)
}

// GenSessionIDOpt generates rfc6733 compliant session ID:
//     <DiameterIdentity>;<high 32 bits>;<low 32 bits>[;<optional value>]
func GenSessionIDOpt(identity, protocol, opt string) string {
	if len(identity) == 0 {
		identity = "magma"
	}
	nano := time.Now().UnixNano()
	ts := uint(nano<<32) ^ uint(nano)
	if len(protocol) != 0 {
		return fmt.Sprintf("%s-%s;%d;%d;%s", identity, protocol, ts, rand.Uint32(), opt)
	}
	return fmt.Sprintf("%s;%d;%d;%s", identity, ts, rand.Uint32(), opt)
}

// GenSessionIDOpt generates rfc6733 compliant session ID:
//     <DiameterIdentity>;<high 32 bits>;<low 32 bits>;IMSI<imsi value>
func GenSessionIdImsi(identity, protocol, imsi string) string {
	if len(identity) == 0 {
		identity = "magma"
	}
	nano := time.Now().UnixNano()
	ts := uint(nano<<32) ^ uint(nano)
	if len(protocol) != 0 {
		return fmt.Sprintf("%s-%s;%d;%d;IMSI%s", identity, protocol, ts, rand.Uint32(), imsi)
	}
	return fmt.Sprintf("%s;%d;%d;IMSI%s", identity, ts, rand.Uint32(), imsi)
}

// GenSessionID generates rfc6733 compliant session ID:
//     <DiameterIdentity>;<high 32 bits>;<low 32 bits>[;<optional value>]
// Where <optional value> is base 16 uint32 random number
func GenSessionID(identity, protocol string) string {
	return GenSessionIDOpt(identity, protocol, strconv.FormatUint(uint64(rand.Uint32()), 16))
}

// GenSessionIDOpt generates rfc6733 compliant session ID:
//     <DiameterIdentity>;<high 32 bits>;<low 32 bits>[;<optional value>]
// Where <DiameterIdentity> is client.Host|ProductName-protocol
func (client *DiameterClientConfig) GenSessionIDOpt(protocol, opt string) string {
	if client != nil {
		if len(client.Host) != 0 {
			return GenSessionIDOpt(client.Host, protocol, opt)
		} else {
			return GenSessionIDOpt(client.ProductName, protocol, opt)
		}
	}
	return GenSessionID("", protocol)
}

// GenSessionID generates rfc6733 compliant session ID:
//     <DiameterIdentity>;<high 32 bits>;<low 32 bits>[;<optional value>]
// Where <DiameterIdentity> is client.Host|ProductName-protocol
//     and <optional value> is base 16 uint32 random number
func (client *DiameterClientConfig) GenSessionID(protocol string) string {
	return client.GenSessionIDOpt(protocol, strconv.FormatUint(uint64(rand.Uint32()), 16))
}

// GenSessionIdImsi generates rfc6733 compliant session ID:
//     <DiameterIdentity>;<high 32 bits>;<low 32 bits>;IMSI<imsi>]
// Where <DiameterIdentity> is client.Host|ProductName-protocol
//     and <optional value> is base 16 uint32 random number
func (client *DiameterClientConfig) GenSessionIdImsi(protocol, imsi string) string {
	if len(imsi) == 0 {
		return client.GenSessionID(protocol)
	}
	return GenSessionIdImsi("", protocol, imsi)
}

// DecodeSessionID extracts and returns session ID if available,
// or original diam SessionId (diamSid) string otherwise
// Input: OriginHost;rand1#;rand2#;IMSIxyz
// Returns: IMSIxyz-rand#
// rand# = rand1# + rand2#, where + means concatenation
func DecodeSessionID(diamSid string) string {
	split := strings.Split(diamSid, ";")
	n1 := len(split) - 1
	if n1 >= 3 && strings.HasPrefix(split[n1], "IMSI") {
		return split[n1] + "-" + split[n1-2] + split[n1-1]
	}
	return diamSid // not magma encoded SID, return as is
}

// ExtractImsiFromSessionID extracts and returns IMSI (without 'IMSI' prefix) from diameter session ID if available,
// or original diam SessionId (diamSid) string otherwise
// Input: OriginHost;[rand1#;rand2#;]IMSIxyz
// Returns: xyz
func ExtractImsiFromSessionID(diamSid string) (string, error) {
	split := strings.Split(diamSid, ";")
	n1 := len(split) - 1
	if n1 > 0 {
		imsi := strings.TrimSpace(split[n1])
		if strings.HasPrefix(imsi, "IMSI") {
			imsi = imsi[4:]
			if len(imsi) < 10 || len(imsi) > 16 {
				return diamSid, fmt.Errorf("Invalid length of IMSI: %s, SessionID: %s", imsi, diamSid)
			}
			for p, l := range imsi {
				if l < '0' || l > '9' {
					return diamSid, fmt.Errorf("Invalid char '%c' in IMSI[%d]: %s, SessionID: %s", l, p, imsi, diamSid)
				}
			}
			return imsi, nil
		}
	}
	return diamSid, fmt.Errorf("Non Magma SessionID: %s", diamSid)
}

// EncodeSessionID encodes SessionID in rfc6733 compliant form:
// <DiameterIdentity>;<high 32 bits>;<low 32 bits>[;<optional value>]
// OriginHost/Realm;rand#;rand#;IMSIxyz
func EncodeSessionID(diamIdentity, sid string) string {
	split := strings.Split(sid, "-")
	if len(split) > 1 && strings.HasPrefix(split[0], "IMSI") {
		rndPart := split[1]
		r2l := len(rndPart) / 2
		return fmt.Sprintf("%s;%s;%s;%s", diamIdentity, rndPart[:r2l], rndPart[r2l:], split[0])
	}
	return sid // not magma generated SID, return as is

}

// ParseDiamSessionID parses given session ID in the form of: // OriginHost;req#;rand#;IMSIxyz_BearerId
// and returns OriginHost, Request Number, Rand, IMSI (without prefix) and bearrerId if present
func ParseDiamSessionID(sessionID string) (host, rnd1, rnd2, imsi, bearrerId string) {
	parts := strings.Split(sessionID, ";")
	l := len(parts)
	switch {
	case l >= 4:
		t := strings.Split(strings.TrimPrefix(parts[3], "IMSI"), "_")
		imsi = t[0]
		if len(t) > 1 {
			bearrerId = t[1]
		}
		fallthrough
	case l == 3:
		rnd2 = parts[2]
		fallthrough
	case l == 2:
		rnd1 = parts[1]
		fallthrough
	case l == 1:
		host = parts[0]
	}
	return
}

func getVendorSpecificApplicationIDAVPs(clientCfg *DiameterClientConfig,
	appIdAvp *diam.AVP) []*diam.AVP {

	var vendorSpecificApplicationIDs []*diam.AVP

	if clientCfg.SupportedVendorIDs != "" {
		// Split the vendor specific application ID string in tokens
		strIds := strings.Split(clientCfg.SupportedVendorIDs, ",")
		// Iterate over each string and append them to the AVP
		for index := range strIds {
			u32, err := strconv.ParseUint(strIds[index], 10, 32)
			if err != nil {
				break
			}
			vendorSpecificApplicationIDs = append(vendorSpecificApplicationIDs,
				diam.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
					AVP: []*diam.AVP{
						appIdAvp,
						diam.NewAVP(avp.VendorID, avp.Mbit, 0,
							datatype.Unsigned32(u32)),
					},
				}))
		}
	}
	return vendorSpecificApplicationIDs

}
