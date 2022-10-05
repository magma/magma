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

package ofpanalytics

import (
	"encoding/json"
	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var (
	AccessToken      = "123"
	SecretCode       = []byte{0x01, 0x02, 0x03, 0x4, 0x05, 0x06}
	NASIPAddress     = "1.1.1.1"
	NASIdentifier    = "1C:B9:C4:3C:C1:80"
	CalledStationID  = "1C-B9-C4-3C-C1-80:Globe Express Wi-Fi by Facebook"
	CallingStationID = "1C-B9-C4-3C-C1-81"
)

type RequestOptions struct {
	NASIpAddress     net.IP
	NASIdentifier    string
	CalledStationID  string
	CallingStationID string
}

func isZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func MiddlewareSink(t *testing.T) modules.Middleware {
	return func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
		require.Fail(t, "Should never be called (ofpanalytics module should not call next())")
		return nil, nil
	}
}

func DefaultAuthRequest() RequestOptions {
	return RequestOptions{
		NASIpAddress:     net.ParseIP(NASIPAddress),
		NASIdentifier:    NASIdentifier,
		CalledStationID:  CalledStationID,
		CallingStationID: CallingStationID,
	}
}

func CreateAuthRequest(options RequestOptions) *radius.Request {
	packet := radius.New(radius.CodeAccessRequest, SecretCode)
	req := &radius.Request{}
	if !isZeroOfUnderlyingType(options.NASIpAddress) {
		rfc2865.NASIPAddress_Add(packet, net.ParseIP(NASIPAddress))
	}
	if !isZeroOfUnderlyingType(options.NASIdentifier) {
		rfc2865.NASIdentifier_AddString(packet, NASIdentifier)
	}
	if !isZeroOfUnderlyingType(options.CalledStationID) {
		rfc2865.CalledStationID_AddString(packet, CalledStationID)
	}
	if !isZeroOfUnderlyingType(options.CallingStationID) {
		rfc2865.CallingStationID_AddString(packet, CallingStationID)
	}
	req.Packet = packet

	return req
}

func TestV2(t *testing.T) {
	cases := []struct {
		req           *radius.Request
		authCode      string
		radiusResCode radius.Code
		name          string
	}{
		{
			req:           CreateAuthRequest(DefaultAuthRequest()),
			authCode:      acceptCode,
			radiusResCode: radius.CodeAccessAccept,
			name:          "accept",
		},
		{
			req: CreateAuthRequest(RequestOptions{
				CalledStationID:  CalledStationID,
				CallingStationID: CallingStationID,
			}),
			authCode:      acceptCode,
			radiusResCode: radius.CodeAccessAccept,
			name:          "accept missing fields",
		},
		{
			req:           CreateAuthRequest(DefaultAuthRequest()),
			authCode:      rejectCode,
			radiusResCode: radius.CodeAccessReject,
			name:          "reject",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test server
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Returning radius auth response
				jsonPacket, _ := json.Marshal(map[string][]string{
					"config:Auth-Type": []string{tc.authCode},
				})
				w.Write(jsonPacket)
			}))
			defer ts.Close()

			logger, _ := zap.NewDevelopment()
			mCtx, err := Init(logger, modules.ModuleConfig{
				"URI":         ts.URL,
				"AccessToken": AccessToken,
				"DryRun":      true,
			})
			require.Nil(t, err)

			res, err := Handle(mCtx, &modules.RequestContext{
				RequestID:      0,
				Logger:         logger,
				SessionStorage: nil,
			}, tc.req, MiddlewareSink(t))
			require.NoError(t, err)

			require.Equal(t, tc.radiusResCode, res.Code)

		})
	}
}
