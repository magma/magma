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

package lballocate

import (
	"context"
	"fbc/cwf/radius/config"
	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/session"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"
	"fmt"
	"math"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"
)

const upstreamRadiusHost1 = "192.168.100.101"
const upstreamRadiusHost2 = "192.168.100.102"
const upstreamRadiusHost3 = "192.168.100.103"
const listener1Name = "listener1"
const listener2Name = "listener2"
const listener3Name = "listener3"
const serviceTier1Name = "serviceTier1"
const serviceTier2Name = "serviceTier2"
const serviceTier3Name = "serviceTier3"
const canary1Name = "c1"
const canary2Name = "c2"
const canary3Name = "c3"

func TestLBAllocateFromLiveTier(t *testing.T) {
	serverConfig := config.ServerConfig{
		LoadBalance: config.LoadBalanceConfig{
			ServiceTiers: []config.ServiceTier{
				{
					Name:          serviceTier1Name,
					UpstreamHosts: []string{upstreamRadiusHost1},
				},
				{
					Name:          serviceTier2Name,
					UpstreamHosts: []string{upstreamRadiusHost2},
				},
				{
					Name:          serviceTier3Name,
					UpstreamHosts: []string{upstreamRadiusHost3},
				},
			},
			LiveTier: config.TierRouting{
				Routes: []config.ListenerRoute{
					{
						Listener:    listener1Name,
						ServiceTier: serviceTier1Name,
					},
					{
						Listener:    listener2Name,
						ServiceTier: serviceTier2Name,
					},
					{
						Listener:    listener3Name,
						ServiceTier: serviceTier3Name,
					},
				},
			},
		},
	}

	doTestLBAllocateSimple(t, &serverConfig, &session.State{Tier: config.LiveTier}, listener2Name, upstreamRadiusHost2)
}

func TestLBAllocateFromCanary(t *testing.T) {
	serverConfig := config.ServerConfig{
		LoadBalance: config.LoadBalanceConfig{
			ServiceTiers: []config.ServiceTier{
				{
					Name:          serviceTier1Name,
					UpstreamHosts: []string{upstreamRadiusHost1},
				},
				{
					Name:          serviceTier2Name,
					UpstreamHosts: []string{upstreamRadiusHost2},
				},
				{
					Name:          serviceTier3Name,
					UpstreamHosts: []string{upstreamRadiusHost3},
				},
			},
			Canaries: []config.Canary{
				{
					Name: canary1Name,
				},
				{
					Name: canary2Name,
					Routing: config.TierRouting{
						Routes: []config.ListenerRoute{
							{
								Listener:    listener1Name,
								ServiceTier: serviceTier1Name,
							},
							{
								Listener:    listener2Name,
								ServiceTier: serviceTier2Name,
							},
							{
								Listener:    listener3Name,
								ServiceTier: serviceTier3Name,
							},
						},
					},
				},
				{
					Name: canary3Name,
				},
			},
		},
	}

	state := session.State{Tier: canary2Name}
	doTestLBAllocateSimple(t, &serverConfig, &state, listener2Name, upstreamRadiusHost2)
}

func TestLBAllocateFromDefaultTier(t *testing.T) {
	serverConfig := config.ServerConfig{
		LoadBalance: config.LoadBalanceConfig{
			ServiceTiers: []config.ServiceTier{
				{
					Name:          serviceTier1Name,
					UpstreamHosts: []string{upstreamRadiusHost1},
				},
				{
					Name:          serviceTier2Name,
					UpstreamHosts: []string{upstreamRadiusHost2},
				},
				{
					Name:          serviceTier3Name,
					UpstreamHosts: []string{upstreamRadiusHost3},
				},
			},
			LiveTier: config.TierRouting{
				Routes: []config.ListenerRoute{
					{
						Listener:    listener1Name,
						ServiceTier: serviceTier1Name,
					},
					{
						Listener:    listener2Name,
						ServiceTier: serviceTier2Name,
					},
					{
						Listener:    listener3Name,
						ServiceTier: serviceTier3Name,
					},
				},
			},
			DefaultTier: config.LiveTier,
		},
	}

	doTestLBAllocateSimple(t, &serverConfig, &session.State{}, listener2Name, upstreamRadiusHost2)
}

func doTestLBAllocateSimple(t *testing.T, serverConfig *config.ServerConfig, state *session.State, listenerName, expectedHost string) {
	// Arrange
	var sessionID = "sessionID"
	Init(serverConfig)

	logger, _ := zap.NewDevelopment()

	sessionStorage := session.NewSessionStorage(session.NewMultiSessionMemoryStorage(), sessionID)
	sessionStorage.Set(*state)

	// Act
	err := Process(
		&modules.RequestContext{
			RequestID:      0,
			Logger:         logger,
			SessionStorage: sessionStorage,
		},
		listenerName,
		createRadiusRequest(),
	)

	// Assert
	require.Nil(t, err)
	s, err := sessionStorage.Get()
	require.Nil(t, err)
	require.Equal(t, expectedHost, s.UpstreamHost)
}

func TestLBAllocateWithNoTierSpecifiedAndNoDefaultTierFails(t *testing.T) {
	serverConfig := config.ServerConfig{
		LoadBalance: config.LoadBalanceConfig{
			ServiceTiers: []config.ServiceTier{
				{
					Name:          serviceTier1Name,
					UpstreamHosts: []string{upstreamRadiusHost1},
				},
				{
					Name:          serviceTier2Name,
					UpstreamHosts: []string{upstreamRadiusHost2},
				},
				{
					Name:          serviceTier3Name,
					UpstreamHosts: []string{upstreamRadiusHost3},
				},
			},
			LiveTier: config.TierRouting{
				Routes: []config.ListenerRoute{
					{
						Listener:    listener1Name,
						ServiceTier: serviceTier1Name,
					},
					{
						Listener:    listener2Name,
						ServiceTier: serviceTier2Name,
					},
					{
						Listener:    listener3Name,
						ServiceTier: serviceTier3Name,
					},
				},
			},
		},
	}

	doTestLBAllocateFails(t, &serverConfig, &session.State{}, listener2Name)
}

func TestLBAllocateFromLiveTierWithNoUpstreamHostsFails(t *testing.T) {
	serverConfig := config.ServerConfig{
		LoadBalance: config.LoadBalanceConfig{
			ServiceTiers: []config.ServiceTier{
				{
					Name:          serviceTier1Name,
					UpstreamHosts: []string{},
				},
			},
			LiveTier: config.TierRouting{
				Routes: []config.ListenerRoute{
					{
						Listener:    listener1Name,
						ServiceTier: serviceTier1Name,
					},
				},
			},
		},
	}

	doTestLBAllocateFails(t, &serverConfig, &session.State{Tier: config.LiveTier}, listener1Name)
}

func TestLBAllocateFromLiveTierWithNoMatchingServiceTierFails(t *testing.T) {
	serverConfig := config.ServerConfig{
		LoadBalance: config.LoadBalanceConfig{
			ServiceTiers: []config.ServiceTier{
				{
					Name:          serviceTier2Name,
					UpstreamHosts: []string{upstreamRadiusHost2},
				},
			},
			LiveTier: config.TierRouting{
				Routes: []config.ListenerRoute{
					{
						Listener:    listener1Name,
						ServiceTier: serviceTier1Name,
					},
				},
			},
		},
	}

	doTestLBAllocateFails(t, &serverConfig, &session.State{Tier: config.LiveTier}, listener1Name)
}

func TestLBAllocateFromLiveTierWithNoMatchingListenerRouteFails(t *testing.T) {
	serverConfig := config.ServerConfig{
		LoadBalance: config.LoadBalanceConfig{
			ServiceTiers: []config.ServiceTier{
				{
					Name:          serviceTier2Name,
					UpstreamHosts: []string{upstreamRadiusHost2},
				},
			},
			LiveTier: config.TierRouting{
				Routes: []config.ListenerRoute{
					{
						Listener:    listener2Name,
						ServiceTier: serviceTier2Name,
					},
				},
			},
		},
	}

	doTestLBAllocateFails(t, &serverConfig, &session.State{Tier: config.LiveTier}, listener1Name)
}

func TestLBAllocateFromCanaryWithNoUpstreamHostsFails(t *testing.T) {
	serverConfig := config.ServerConfig{
		LoadBalance: config.LoadBalanceConfig{
			ServiceTiers: []config.ServiceTier{
				{
					Name:          serviceTier1Name,
					UpstreamHosts: []string{},
				},
			},
			Canaries: []config.Canary{
				{
					Name: canary2Name,
					Routing: config.TierRouting{
						Routes: []config.ListenerRoute{
							{
								Listener:    listener1Name,
								ServiceTier: serviceTier1Name,
							},
						},
					},
				},
			},
		},
	}

	state := session.State{Tier: canary2Name}
	doTestLBAllocateFails(t, &serverConfig, &state, listener1Name)
}

func TestLBAllocateFromCanaryWithNoMatchingServiceTierFails(t *testing.T) {
	serverConfig := config.ServerConfig{
		LoadBalance: config.LoadBalanceConfig{
			ServiceTiers: []config.ServiceTier{
				{
					Name:          serviceTier2Name,
					UpstreamHosts: []string{upstreamRadiusHost2},
				},
			},
			Canaries: []config.Canary{
				{
					Name: canary2Name,
					Routing: config.TierRouting{
						Routes: []config.ListenerRoute{
							{
								Listener:    listener1Name,
								ServiceTier: serviceTier1Name,
							},
						},
					},
				},
			},
		},
	}

	state := session.State{Tier: canary2Name}
	doTestLBAllocateFails(t, &serverConfig, &state, listener1Name)
}

func TestLBAllocateFromCanaryWithNoMatchingListenerRouteFails(t *testing.T) {
	serverConfig := config.ServerConfig{
		LoadBalance: config.LoadBalanceConfig{
			ServiceTiers: []config.ServiceTier{
				{
					Name:          serviceTier2Name,
					UpstreamHosts: []string{upstreamRadiusHost2},
				},
			},
			Canaries: []config.Canary{
				{
					Name: canary2Name,
					Routing: config.TierRouting{
						Routes: []config.ListenerRoute{
							{
								Listener:    listener2Name,
								ServiceTier: serviceTier2Name,
							},
						},
					},
				},
			},
		},
	}
	state := session.State{Tier: canary2Name}
	doTestLBAllocateFails(t, &serverConfig, &state, listener1Name)
}

func TestLBAllocateFromCanaryWithNoMatchingCanaryDefinition(t *testing.T) {
	serverConfig := config.ServerConfig{
		LoadBalance: config.LoadBalanceConfig{
			ServiceTiers: []config.ServiceTier{
				{
					Name:          serviceTier2Name,
					UpstreamHosts: []string{upstreamRadiusHost2},
				},
			},
			Canaries: []config.Canary{
				{
					Name: canary2Name,
					Routing: config.TierRouting{
						Routes: []config.ListenerRoute{
							{
								Listener:    listener2Name,
								ServiceTier: serviceTier2Name,
							},
						},
					},
				},
			},
		},
	}
	state := session.State{Tier: canary3Name}
	doTestLBAllocateFails(t, &serverConfig, &state, listener2Name)
}

func doTestLBAllocateFails(t *testing.T, serverConfig *config.ServerConfig, state *session.State, listenerName string) {
	// Arrange
	var (
		sessionID    = "sessionID"
		genSessionID = "genSessionID"
	)
	Init(serverConfig)

	logger, _ := zap.NewDevelopment()

	// Test a case of initially missing Acct-Session-ID (generated SessionID is used) and then provided Acct-Session-ID
	// for a session with generated SID already in the storage...
	globalStorage := session.NewMultiSessionMemoryStorage()
	// Missing Acct-Session-ID -> use genSessionID to add session to storage
	sessionStorage := session.NewSessionStorageExt(globalStorage, genSessionID, genSessionID)
	sessionStorage.Set(*state)

	// Act
	err := Process(
		&modules.RequestContext{
			RequestID:      0,
			Logger:         logger,
			SessionStorage: sessionStorage,
		},
		listenerName,
		createRadiusRequest(),
	)

	// Assert
	require.NotNil(t, err)

	// Provided Acct-Session-ID -> use storage with previously genSessionID mapping
	sessionStorage = session.NewSessionStorageExt(globalStorage, sessionID, genSessionID)
	// Act
	err = Process(
		&modules.RequestContext{
			RequestID:      0,
			Logger:         logger,
			SessionStorage: sessionStorage,
		},
		listenerName,
		createRadiusRequest(),
	)

	// Assert
	require.NotNil(t, err)
}

func TestLBAllocateEvenDistribution(t *testing.T) {
	// Arrange
	const hostCount = 20
	const invocations = 100000
	const portOffset = 20000
	var sessionID = "sessionID"
	hostInvocations := make([]int, hostCount)
	upstreamHosts := make([]string, hostCount)
	errors := make([]error, invocations)
	for i := 0; i < hostCount; i++ {
		hostInvocations[i] = 0
		upstreamHosts[i] = fmt.Sprintf("192.168.200.200:%d", i+portOffset)
	}
	serverConfig := config.ServerConfig{
		LoadBalance: config.LoadBalanceConfig{
			ServiceTiers: []config.ServiceTier{
				{
					Name:          serviceTier3Name,
					UpstreamHosts: upstreamHosts,
				},
			},
			LiveTier: config.TierRouting{
				Routes: []config.ListenerRoute{
					{
						Listener:    listener3Name,
						ServiceTier: serviceTier3Name,
					},
				},
			},
		},
	}
	Init(&serverConfig)
	logger, _ := zap.NewDevelopment()

	// Act
	for i := 0; i < invocations; i++ {
		sessionStorage := session.NewSessionStorage(session.NewMultiSessionMemoryStorage(), sessionID)
		state := session.State{Tier: config.LiveTier}
		sessionStorage.Set(state)

		err := Process(
			&modules.RequestContext{
				RequestID:      0,
				Logger:         logger,
				SessionStorage: sessionStorage,
			},
			listener3Name,
			createRadiusRequest(),
		)

		s, _ := sessionStorage.Get()

		var dummyHost string
		var port int
		fmt.Sscanf(s.UpstreamHost, "%16s%d", &dummyHost, &port)
		hostIndex := port - portOffset
		hostInvocations[hostIndex]++

		errors[i] = err
	}

	// Assert
	for i := 0; i < invocations; i++ {
		require.Nil(t, errors[i])
	}
	expectedProportion := 100.0 / float64(hostCount)
	for i := 0; i < hostCount; i++ {
		actualProportion := 100.0 * float64(hostInvocations[i]) / float64(invocations)
		deviation := math.Abs(actualProportion - expectedProportion)
		require.True(t, deviation < 3.0)
	}
}

func TestLBAllocateMaintainsStickiness(t *testing.T) {
	// Arrange
	const hostCount = 120
	const invocations = 10
	const portOffset = 30000
	var sessionID = "sessionID"
	hostInvocations := make([]int, hostCount)
	upstreamHosts := make([]string, hostCount)
	errors := make([]error, invocations)
	for i := 0; i < hostCount; i++ {
		hostInvocations[i] = 0
		upstreamHosts[i] = fmt.Sprintf("192.168.200.200:%d", i+portOffset)
	}
	serverConfig := config.ServerConfig{
		LoadBalance: config.LoadBalanceConfig{
			ServiceTiers: []config.ServiceTier{
				{
					Name:          serviceTier3Name,
					UpstreamHosts: upstreamHosts,
				},
			},
			LiveTier: config.TierRouting{
				Routes: []config.ListenerRoute{
					{
						Listener:    listener3Name,
						ServiceTier: serviceTier3Name,
					},
				},
			},
		},
	}
	Init(&serverConfig)
	logger, _ := zap.NewDevelopment()

	// Act
	for i := 0; i < invocations; i++ {
		state := session.State{Tier: config.LiveTier, UpstreamHost: upstreamHosts[73]}

		sessionStorage := session.NewSessionStorage(session.NewMultiSessionMemoryStorage(), sessionID)
		sessionStorage.Set(state)

		err := Process(
			&modules.RequestContext{
				RequestID:      0,
				Logger:         logger,
				SessionStorage: sessionStorage,
			},
			listener3Name,
			createRadiusRequest(),
		)

		s, _ := sessionStorage.Get()

		var dummyHost string
		var port int
		fmt.Sscanf(s.UpstreamHost, "%16s%d", &dummyHost, &port)
		hostIndex := port - portOffset
		hostInvocations[hostIndex]++

		errors[i] = err
	}

	// Assert
	for i := 0; i < invocations; i++ {
		require.Nil(t, errors[i])
	}
	for i := 0; i < hostCount; i++ {
		if i == 73 {
			require.Equal(t, invocations, hostInvocations[i])
		} else {
			require.Equal(t, 0, hostInvocations[i])
		}
	}
}

func createRadiusRequest() *radius.Request {
	packet := radius.New(radius.CodeAccessRequest, []byte{0x01, 0x02, 0x03, 0x4, 0x05, 0x06})
	packet.Attributes[rfc2865.CallingStationID_Type] = []radius.Attribute{radius.Attribute("calling")}
	packet.Attributes[rfc2865.CalledStationID_Type] = []radius.Attribute{radius.Attribute("called")}
	req := &radius.Request{}
	req = req.WithContext(context.Background())
	req.Packet = packet
	return req
}
