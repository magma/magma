/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
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

func MiddlewareSink(t *testing.T) modules.Middleware {
	return func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
		require.Fail(t, "Should never be called (ofpanalytics module should not call next())")
		return nil, nil
	}
}

func CreateAuthRequest() *radius.Request {

	packet := radius.New(radius.CodeAccessRequest, SecretCode)
	req := &radius.Request{}
	rfc2865.NASIPAddress_Add(packet, net.ParseIP(NASIPAddress))
	rfc2865.NASIdentifier_AddString(packet, NASIdentifier)
	rfc2865.CalledStationID_AddString(packet, CalledStationID)
	rfc2865.CallingStationID_AddString(packet, CallingStationID)
	req.Packet = packet

	return req
}

func TestV2(t *testing.T) {
	cases := []struct {
		authCode      string
		radiusResCode radius.Code
		name          string
	}{
		{
			authCode:      acceptCode,
			radiusResCode: radius.CodeAccessAccept,
			name:          "accept",
		},
		{
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

			req := CreateAuthRequest()
			res, err := Handle(mCtx, &modules.RequestContext{
				RequestID:      0,
				Logger:         logger,
				SessionStorage: nil,
			}, req, MiddlewareSink(t))
			require.NoError(t, err)

			require.Equal(t, tc.radiusResCode, res.Code)

		})
	}
}
