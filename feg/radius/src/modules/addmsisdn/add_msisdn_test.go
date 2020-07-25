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

package addmsisdn

import (
	"context"
	"errors"
	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/session"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"
)

const (
	ExpressWiFiVendor uint32 = 40981
)

func TestMsisdnAdded(t *testing.T) {
	// Arrange
	var sessionID = "sessionID"
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "failed to get logger")
	ctx, err := Init(logger, modules.ModuleConfig{})
	require.NoError(t, err, "failed to init")

	var outputMsisdn string
	var outputVendor uint32
	var msisdn = "+1234567890"

	sessionStorage := session.NewSessionStorage(session.NewMultiSessionMemoryStorage(), sessionID)
	sessionStorage.Set(session.State{MACAddress: "fa:ce:b0:0c:12:34", MSISDN: msisdn})

	// Act
	_, err = Handle(
		ctx,
		&modules.RequestContext{
			RequestID:      0,
			Logger:         logger,
			SessionStorage: sessionStorage,
		},
		createRadiusRequest("called", "calling"),
		func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
			attrs, ok := r.Attributes[rfc2865.VendorSpecific_Type]
			if !ok {
				return nil, errors.New("cannot find VendorSpecific attribute in response")
			}

			for _, vsattr := range attrs {
				vendor, attr, err := radius.VendorSpecific(vsattr)
				if err != nil {
					return nil, err
				}
				outputVendor = vendor
				if attr[0] == 7 {
					outputMsisdn = string(attr[2:])
				}
			}
			return nil, nil
		},
	)

	// Act and Assert
	require.Nil(t, err)
	require.Equal(t, msisdn, outputMsisdn)
	require.Equal(t, ExpressWiFiVendor, outputVendor)
}

func TestMissingSessionState(t *testing.T) {
	// Arrange
	sessionID := "sessionID"
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "failed to get logger")
	Init(logger, modules.ModuleConfig{})
	sessionStorage := session.NewSessionStorage(session.NewMultiSessionMemoryStorage(), sessionID)

	// Act
	_, err = Handle(
		nil,
		&modules.RequestContext{
			RequestID:      0,
			Logger:         logger,
			SessionStorage: sessionStorage,
		},
		createRadiusRequest("called", "calling"),
		func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
			require.Fail(t, "Should never have been called, expected to fail on missing state")
			return nil, nil
		},
	)

	// Act and Assert
	require.NotNil(t, err)
	require.Equal(t, "session sessionID no found in storage", err.Error())
}

func TestInvalidVendorSpecificAttr(t *testing.T) {
	// Arrange
	var sessionID = strings.Repeat("a", 300)
	logger, err := zap.NewDevelopment()
	require.NoError(t, err, "failed to get logger")
	Init(logger, modules.ModuleConfig{})
	sessionStorage := session.NewSessionStorage(session.NewMultiSessionMemoryStorage(), sessionID)
	sessionStorage.Set(session.State{MACAddress: "fa:ce:b0:0c:12:34", MSISDN: sessionID})

	// Act
	_, err = Handle(
		nil,
		&modules.RequestContext{
			RequestID:      0,
			Logger:         logger,
			SessionStorage: sessionStorage,
		},
		createRadiusRequest("called", "calling"),
		func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
			require.Fail(t, "Should never have been called, expected to fail on missing state")
			return nil, nil
		},
	)

	// Act and Assert
	require.NotNil(t, err)
	require.Equal(t, "Failed encoding MSISDN attribute: value too long", err.Error())
}

func createRadiusRequest(calledStationID string, callingStationID string) *radius.Request {
	packet := radius.New(radius.CodeAccessRequest, []byte{0x01, 0x02, 0x03, 0x4, 0x05, 0x06})
	packet.Attributes[rfc2865.CallingStationID_Type] = []radius.Attribute{radius.Attribute(callingStationID)}
	packet.Attributes[rfc2865.CalledStationID_Type] = []radius.Attribute{radius.Attribute(calledStationID)}
	req := &radius.Request{}
	req = req.WithContext(context.Background())
	req.Packet = packet
	return req
}
