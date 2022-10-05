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

package lbcanary

import (
	"context"
	"fbc/cwf/radius/config"
	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/session"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"
	"math"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"
)

const dummyListener = "dummyListener"

func TestLBCanaryFailsWhenUsingCanaryReservedNameLive(t *testing.T) {

	// Arrange
	serverConfig := config.ServerConfig{
		LoadBalance: config.LoadBalanceConfig{
			Canaries: []config.Canary{
				{
					Name:                "c1",
					TrafficSlicePercent: 5,
				},
				{
					Name:                config.LiveTier,
					TrafficSlicePercent: 25,
				},
				{
					Name:                "c3",
					TrafficSlicePercent: 30,
				},
			},
		},
	}

	// Act
	err := Init(&serverConfig)

	// Assert
	require.NotNil(t, err)
}

func TestLBCanaryFailsWhenCanariesOverallocated(t *testing.T) {

	// Arrange
	serverConfig := config.ServerConfig{
		LoadBalance: config.LoadBalanceConfig{
			Canaries: []config.Canary{
				{
					Name:                "c1",
					TrafficSlicePercent: 5,
				},
				{
					Name:                "c2",
					TrafficSlicePercent: 25,
				},
				{
					Name:                "c3",
					TrafficSlicePercent: 30,
				},
			},
		},
	}

	// Act
	err := Init(&serverConfig)

	// Assert
	require.NotNil(t, err)
}

func TestLBCanaryAllocatesTiersProportionally(t *testing.T) {
	// Arrange
	const invocations = 100000
	var sessionID = "sessionID"
	c1Count, c2Count, liveTierCount := 0, 0, 0
	c1Proportion, c2Proportion, liveTierProportion := 5, 25, 70
	serverConfig := config.ServerConfig{
		LoadBalance: config.LoadBalanceConfig{
			Canaries: []config.Canary{
				{
					Name:                "c1",
					TrafficSlicePercent: c1Proportion,
				},
				{
					Name:                "c2",
					TrafficSlicePercent: c2Proportion,
				},
			},
		},
	}

	// Act
	Init(&serverConfig)

	for i := 0; i < invocations; i++ {
		logger, _ := zap.NewDevelopment()

		state := session.State{}
		sessionStorage := session.NewSessionStorage(session.NewMultiSessionMemoryStorage(), sessionID)
		sessionStorage.Set(state)

		Process(
			&modules.RequestContext{
				RequestID:      0,
				Logger:         logger,
				SessionStorage: sessionStorage,
			},
			dummyListener,
			createRadiusRequest("called", "calling"),
		)

		s, _ := sessionStorage.Get()

		switch s.Tier {
		case "c1":
			c1Count++
		case "c2":
			c2Count++
		case config.LiveTier:
			liveTierCount++
		}
	}

	// Assert
	c1ActualProportion := 100.0 * float64(c1Count) / invocations
	c2ActualProportion := 100.0 * float64(c2Count) / invocations
	liveTierActualProportion := 100 * float64(liveTierCount) / invocations
	c1Deviation := math.Abs(float64(c1Proportion) - c1ActualProportion)
	c2Deviation := math.Abs(float64(c2Proportion) - c2ActualProportion)
	liveTierDeviation := math.Abs(float64(liveTierProportion) - liveTierActualProportion)
	require.True(t, c1Deviation < 3.0)
	require.True(t, c2Deviation < 3.0)
	require.True(t, liveTierDeviation < 3.0)
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
