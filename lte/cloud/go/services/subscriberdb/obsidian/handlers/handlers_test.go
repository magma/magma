/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handlers_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	lteHandlers "magma/lte/cloud/go/services/lte/obsidian/handlers"
	lteModels "magma/lte/cloud/go/services/lte/obsidian/models"
	policydbHandlers "magma/lte/cloud/go/services/policydb/obsidian/handlers"
	policydbModels "magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/handlers"
	subscriberModels "magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	subscriberdbTestInit "magma/lte/cloud/go/services/subscriberdb/test_init"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
	testUtilConfigurator "magma/orc8r/cloud/go/services/configurator/test_utils"
	deviceTestInit "magma/orc8r/cloud/go/services/device/test_init"
	directorydTypes "magma/orc8r/cloud/go/services/directoryd/types"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"
	stateTypes "magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/cloud/go/storage"
)

func TestCreateSubscribers(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	networkConfigs := map[string]interface{}{
		lte.CellularNetworkConfigType: &lteModels.NetworkCellularConfigs{
			Epc: &lteModels.NetworkEpcConfigs{SubProfiles: map[string]lteModels.NetworkEpcConfigsSubProfilesAnon{"present-profile": {}}},
		},
	}
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1", Configs: networkConfigs}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/subscribers"
	handlers := handlers.GetHandlers()
	createSubscriber := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.POST).HandlerFunc
	listSubscribers := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc

	// Pre: seed 2 apns
	_, err = configurator.CreateEntities(context.Background(), "n1", []configurator.NetworkEntity{
		{Type: lte.APNEntityType, Key: "apn0"},
		{Type: lte.APNEntityType, Key: "apn1"},
	}, serdes.Entity)
	assert.NoError(t, err)

	// Pass: happy path
	sub0 := newMutableSubscriber("IMSI0000000000")
	sub1 := newMutableSubscriber("IMSI0000000001")
	sub1.Name = "Johnny Appleseed"
	payload := subscriberModels.MutableSubscribers{sub0, sub1}
	tc := tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(payload),
		Handler:        createSubscriber,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)
	expected := subscriberModels.PaginatedSubscribers{
		TotalCount:    2,
		NextPageToken: "",
		Subscribers: map[string]*subscriberModels.Subscriber{
			"IMSI0000000000": sub0.ToSubscriber(),
			"IMSI0000000001": sub1.ToSubscriber(),
		},
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expected),
	}
	tests.RunUnitTest(t, e, tc)

	// Fail: create subs that already exists
	sub2 := newMutableSubscriber("IMSI0000000002")
	payload = subscriberModels.MutableSubscribers{sub0, sub1, sub2}
	tc = tests.Test{
		Method:                 "POST",
		URL:                    testURLRoot,
		Payload:                tests.JSONMarshaler(payload),
		Handler:                createSubscriber,
		ParamNames:             []string{"network_id"},
		ParamValues:            []string{"n1"},
		ExpectedStatus:         400,
		ExpectedErrorSubstring: "found 2 existing subscribers which would have been overwritten",
	}
	tests.RunUnitTest(t, e, tc)

	// Fail: create two subs with same IMSI
	sub3 := newMutableSubscriber("IMSI0000000003")
	sub4 := newMutableSubscriber("IMSI0000000004")
	sub5 := newMutableSubscriber("IMSI0000000004") // same as sub4
	payload = subscriberModels.MutableSubscribers{sub3, sub4, sub5}
	tc = tests.Test{
		Method:                 "POST",
		URL:                    testURLRoot,
		Payload:                tests.JSONMarshaler(payload),
		Handler:                createSubscriber,
		ParamNames:             []string{"network_id"},
		ParamValues:            []string{"n1"},
		ExpectedStatus:         400,
		ExpectedErrorSubstring: "found multiple subscriber models for IDs",
	}
	tests.RunUnitTest(t, e, tc)

	// Fail: create sub with non-default sub profile that's missing from network config
	sub6 := newMutableSubscriber("IMSI0000000006")
	sub6.Lte.SubProfile = "missing-profile"
	payload = subscriberModels.MutableSubscribers{sub6}
	tc = tests.Test{
		Method:                 "POST",
		URL:                    testURLRoot,
		Payload:                tests.JSONMarshaler(payload),
		Handler:                createSubscriber,
		ParamNames:             []string{"network_id"},
		ParamValues:            []string{"n1"},
		ExpectedStatus:         400,
		ExpectedErrorSubstring: "subscriber profile 'missing-profile' does not exist for the network",
	}
	tests.RunUnitTest(t, e, tc)

	// Pass: create sub with non-default sub profile that's present in network config
	sub7 := newMutableSubscriber("IMSI0000000007")
	sub8 := newMutableSubscriber("IMSI0000000008")
	sub7.Lte.SubProfile = "present-profile"
	sub8.Lte.SubProfile = "present-profile"
	payload = subscriberModels.MutableSubscribers{sub7, sub8}
	tc = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(payload),
		Handler:        createSubscriber,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

}

func TestListSubscribers(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/subscribers"
	handlers := handlers.GetHandlers()
	listSubscribers := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc

	// preseed 2 apns
	apn1, apn2 := "foo", "bar"
	_, err = configurator.CreateEntities(context.Background(), "n1", []configurator.NetworkEntity{
		{Type: lte.APNEntityType, Key: apn1},
		{Type: lte.APNEntityType, Key: apn2},
	}, serdes.Entity)
	assert.NoError(t, err)
	expectedResult := subscriberModels.PaginatedSubscribers{
		TotalCount:    int64(0),
		NextPageToken: "",
		Subscribers:   map[string]*subscriberModels.Subscriber{},
	}
	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedResult),
	}
	tests.RunUnitTest(t, e, tc)

	// Set the total expected count to 3, the number of subscribers created below
	expectedResult.TotalCount = 3
	_, err = configurator.CreateEntities(context.Background(), "n1", []configurator.NetworkEntity{
		{
			Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
			Config: &subscriberModels.SubscriberConfig{
				Lte: &subscriberModels.LteSubscription{
					AuthAlgo: "MILENAGE",
					AuthKey:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					AuthOpc:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					State:    "ACTIVE",
				},
				StaticIps:             subscriberModels.SubscriberStaticIps{apn1: "192.168.100.1", apn2: "10.10.10.5"},
				ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"5GC"},
			},
			Associations: storage.TKs{{Type: lte.APNEntityType, Key: apn2}, {Type: lte.APNEntityType, Key: apn1}},
		},
		{
			Type: lte.SubscriberEntityType, Key: "IMSI0987654321",
			Config: &subscriberModels.SubscriberConfig{
				Lte: &subscriberModels.LteSubscription{
					AuthAlgo:   "MILENAGE",
					AuthKey:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
					AuthOpc:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
					State:      "ACTIVE",
					SubProfile: "foo",
				},
				ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"EPC"},
			},
			Associations: storage.TKs{{Type: lte.APNEntityType, Key: apn1}},
		},
		{
			Type: lte.SubscriberEntityType, Key: "IMSI0987654322",
			Config: &subscriberModels.SubscriberConfig{
				Lte: &subscriberModels.LteSubscription{
					AuthAlgo:   "MILENAGE",
					AuthKey:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
					AuthOpc:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
					State:      "ACTIVE",
					SubProfile: "foo",
				},
				ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"EPC"},
			},
			Associations: storage.TKs{{Type: lte.APNEntityType, Key: apn2}},
		},
	}, serdes.Entity)
	assert.NoError(t, err)

	_, err = configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g1", Config: &models.MagmadGatewayConfigs{}, PhysicalID: "hw1"}, serdes.Entity)
	assert.NoError(t, err)
	frozenClock := int64(1000000)
	clock.SetAndFreezeClock(t, time.Unix(frozenClock, 0))
	defer clock.UnfreezeClock(t)

	icmpStatus := &subscriberModels.IcmpStatus{LatencyMs: f32Ptr(12.34)}
	ctx := test_utils.GetContextWithCertificate(t, "hw1")
	test_utils.ReportState(t, ctx, lte.ICMPStateType, "IMSI1234567890", icmpStatus, serdes.State)
	mmeState := state.ArbitraryJSON{"mme": "foo"}
	test_utils.ReportState(t, ctx, lte.MMEStateType, "IMSI1234567890", &mmeState, serdes.State)
	spgwState := state.ArbitraryJSON{"spgw": "foo"}
	test_utils.ReportState(t, ctx, lte.SPGWStateType, "IMSI1234567890", &spgwState, serdes.State)
	s1apState := state.ArbitraryJSON{"s1ap": "foo"}
	test_utils.ReportState(t, ctx, lte.S1APStateType, "IMSI1234567890", &s1apState, serdes.State)
	// Report 2 allocated IP addresses for the subscriber
	mobilitydState1 := state.ArbitraryJSON{
		"ip": map[string]interface{}{
			"address": "wKiArg==",
		},
	}
	mobilitydState2 := state.ArbitraryJSON{
		"ip": map[string]interface{}{
			"address": "wKiAhg==",
		},
	}
	test_utils.ReportState(t, ctx, lte.MobilitydStateType, "IMSI1234567890.oai.ipv4", &mobilitydState1, serdes.State)
	test_utils.ReportState(t, ctx, lte.MobilitydStateType, "IMSI1234567890.magma.apn", &mobilitydState2, serdes.State)
	directoryState := directorydTypes.DirectoryRecord{LocationHistory: []string{"foo", "bar"}}
	test_utils.ReportState(t, ctx, orc8r.DirectoryRecordType, "IMSI1234567890", &directoryState, serdes.State)

	expectedResult.Subscribers = map[string]*subscriberModels.Subscriber{
		"IMSI0987654321": {
			ID: "IMSI0987654321",
			Lte: &subscriberModels.LteSubscription{
				AuthAlgo:   "MILENAGE",
				AuthKey:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
				AuthOpc:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
				State:      "ACTIVE",
				SubProfile: "foo",
			},
			ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"EPC"},
			Config: &subscriberModels.SubscriberConfig{
				Lte: &subscriberModels.LteSubscription{
					AuthAlgo:   "MILENAGE",
					AuthKey:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
					AuthOpc:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
					State:      "ACTIVE",
					SubProfile: "foo",
				},
				ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"EPC"},
			},
			ActiveApns: subscriberModels.ApnList{apn1},
		},
		"IMSI0987654322": {
			ID: "IMSI0987654322",
			Lte: &subscriberModels.LteSubscription{
				AuthAlgo:   "MILENAGE",
				AuthKey:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
				AuthOpc:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
				State:      "ACTIVE",
				SubProfile: "foo",
			},
			ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"EPC"},
			Config: &subscriberModels.SubscriberConfig{
				Lte: &subscriberModels.LteSubscription{
					AuthAlgo:   "MILENAGE",
					AuthKey:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
					AuthOpc:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
					State:      "ACTIVE",
					SubProfile: "foo",
				},
				ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"EPC"},
			},
			ActiveApns: subscriberModels.ApnList{apn2},
		},
	}
	expectedResult.NextPageToken = "Cg5JTVNJMDk4NzY1NDMyMg=="

	// Test paginated requests
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot + "?page_size=2&page_token=",
		Handler:        listSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedResult),
	}
	tests.RunUnitTest(t, e, tc)

	// Ensure the same request returns the same output
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot + "?page_size=2&page_token=",
		Handler:        listSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedResult),
	}
	tests.RunUnitTest(t, e, tc)

	// No page token should return the same results
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot + "?page_size=2",
		Handler:        listSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedResult),
	}
	tests.RunUnitTest(t, e, tc)

	expectedResult.Subscribers = map[string]*subscriberModels.Subscriber{
		"IMSI1234567890": {
			ID: "IMSI1234567890",
			Lte: &subscriberModels.LteSubscription{
				AuthAlgo:   "MILENAGE",
				AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:      "ACTIVE",
				SubProfile: "default",
			},
			ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"5GC"},
			Config: &subscriberModels.SubscriberConfig{
				Lte: &subscriberModels.LteSubscription{
					AuthAlgo:   "MILENAGE",
					AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					State:      "ACTIVE",
					SubProfile: "default",
				},
				StaticIps:             subscriberModels.SubscriberStaticIps{apn1: "192.168.100.1", apn2: "10.10.10.5"},
				ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"5GC"},
			},
			ActiveApns: subscriberModels.ApnList{apn2, apn1},
			Monitoring: &subscriberModels.SubscriberStatus{
				Icmp: &subscriberModels.IcmpStatus{
					// LastReportedTime is calculated in milliseconds
					LastReportedTime: frozenClock * 1000,
					LatencyMs:        f32Ptr(12.34),
				},
			},
			State: &subscriberModels.SubscriberState{
				Mme:  mmeState,
				S1ap: s1apState,
				Spgw: spgwState,
				Mobility: []*subscriberModels.SubscriberIPAllocation{
					{
						Apn: "magma.apn",
						IP:  "192.168.128.134",
					},
					{
						Apn: "oai.ipv4",
						IP:  "192.168.128.174",
					},
				},
				Directory: &subscriberModels.SubscriberDirectoryRecord{
					LocationHistory: []string{"foo", "bar"},
				},
			},
		},
	}
	expectedResult.NextPageToken = ""
	// Get last page of subscribers
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot + "?page_size=2&page_token=Cg5JTVNJMDk4NzY1NDMyMg==",
		Handler:        listSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedResult),
	}
	tests.RunUnitTest(t, e, tc)

	expectedResultAbbreviated := subscriberModels.PaginatedSubscriberIds{
		TotalCount:    expectedResult.TotalCount,
		NextPageToken: expectedResult.NextPageToken,
	}

	for k := range expectedResult.Subscribers {
		expectedResultAbbreviated.Subscribers = append(expectedResultAbbreviated.Subscribers, k)
	}

	assert.NotEqual(t, 0, len(expectedResult.Subscribers), "Subscriber list empty! Make sure there's at least one subscriber")
	assert.NotEqual(t, 0, expectedResult.TotalCount, "Total count empty! Make sure there's at least one subscriber")

	// Get last page of subscribers
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot + "?page_size=2&verbose=false&page_token=Cg5JTVNJMDk4NzY1NDMyMg==",
		Handler:        listSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedResultAbbreviated),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestGetSubscriber(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/subscribers/:subscriber_id"
	handlers := handlers.GetHandlers()
	getSubscriber := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc

	//preseed 2 apns
	apn1, apn2 := "foo", "bar"
	_, err = configurator.CreateEntities(context.Background(), "n1", []configurator.NetworkEntity{
		{Type: lte.APNEntityType, Key: apn1},
		{Type: lte.APNEntityType, Key: apn2},
	}, serdes.Entity)
	assert.NoError(t, err)

	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getSubscriber,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	// No sub profile configured, we should return "default"
	_, err = configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{
		Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
		Name: "Jane Doe",
		Config: &subscriberModels.SubscriberConfig{
			Lte: &subscriberModels.LteSubscription{
				AuthAlgo: "MILENAGE",
				AuthKey:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:    "ACTIVE",
			},
			ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"5GC"},
			StaticIps:             subscriberModels.SubscriberStaticIps{apn1: "192.168.100.1"},
		},
		Associations: storage.TKs{{Type: lte.APNEntityType, Key: apn2}, {Type: lte.APNEntityType, Key: apn1}},
	}, serdes.Entity)
	assert.NoError(t, err)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getSubscriber,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 200,
		ExpectedResult: &subscriberModels.Subscriber{
			ID:   "IMSI1234567890",
			Name: "Jane Doe",
			Lte: &subscriberModels.LteSubscription{
				AuthAlgo:   "MILENAGE",
				AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:      "ACTIVE",
				SubProfile: "default",
			},
			ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"5GC"},
			Config: &subscriberModels.SubscriberConfig{
				Lte: &subscriberModels.LteSubscription{
					AuthAlgo:   "MILENAGE",
					AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					State:      "ACTIVE",
					SubProfile: "default",
				},
				ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"5GC"},
				StaticIps:             subscriberModels.SubscriberStaticIps{apn1: "192.168.100.1"},
			},
			ActiveApns: subscriberModels.ApnList{apn2, apn1},
		},
	}
	tests.RunUnitTest(t, e, tc)

	// Now create AGW
	// First we need to register a gateway which can report state
	_, err = configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g1", Config: &models.MagmadGatewayConfigs{}, PhysicalID: "hw1"}, serdes.Entity)
	assert.NoError(t, err)
	frozenClock := int64(1000000)
	clock.SetAndFreezeClock(t, time.Unix(frozenClock, 0))
	defer clock.UnfreezeClock(t)
	icmpStatus := &subscriberModels.IcmpStatus{LatencyMs: f32Ptr(12.34)}
	ctx := test_utils.GetContextWithCertificate(t, "hw1")
	test_utils.ReportState(t, ctx, lte.ICMPStateType, "IMSI1234567890", icmpStatus, serdes.State)
	mmeState := state.ArbitraryJSON{"mme": "foo"}
	test_utils.ReportState(t, ctx, lte.MMEStateType, "IMSI1234567890", &mmeState, serdes.State)
	spgwState := state.ArbitraryJSON{"spgw": "foo"}
	test_utils.ReportState(t, ctx, lte.SPGWStateType, "IMSI1234567890", &spgwState, serdes.State)
	s1apState := state.ArbitraryJSON{"s1ap": "foo"}
	test_utils.ReportState(t, ctx, lte.S1APStateType, "IMSI1234567890", &s1apState, serdes.State)
	// Report 2 allocated IP addresses for the subscriber
	mobilitydState1 := state.ArbitraryJSON{
		"ip": map[string]interface{}{
			"address": "wKiArg==",
		},
	}
	mobilitydState2 := state.ArbitraryJSON{
		"ip": map[string]interface{}{
			"address": "wKiAhg==",
		},
	}
	test_utils.ReportState(t, ctx, lte.MobilitydStateType, "IMSI1234567890.oai.ipv4", &mobilitydState1, serdes.State)
	test_utils.ReportState(t, ctx, lte.MobilitydStateType, "IMSI1234567890.magma.apn", &mobilitydState2, serdes.State)
	directoryState := directorydTypes.DirectoryRecord{LocationHistory: []string{"foo", "bar"}}
	test_utils.ReportState(t, ctx, orc8r.DirectoryRecordType, "IMSI1234567890", &directoryState, serdes.State)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getSubscriber,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 200,
		ExpectedResult: &subscriberModels.Subscriber{
			ID:   "IMSI1234567890",
			Name: "Jane Doe",
			Lte: &subscriberModels.LteSubscription{
				AuthAlgo:   "MILENAGE",
				AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:      "ACTIVE",
				SubProfile: "default",
			},
			ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"5GC"},
			Config: &subscriberModels.SubscriberConfig{
				Lte: &subscriberModels.LteSubscription{
					AuthAlgo:   "MILENAGE",
					AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					State:      "ACTIVE",
					SubProfile: "default",
				},
				ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"5GC"},
				StaticIps:             subscriberModels.SubscriberStaticIps{apn1: "192.168.100.1"},
			},
			ActiveApns: subscriberModels.ApnList{apn2, apn1},
			Monitoring: &subscriberModels.SubscriberStatus{
				Icmp: &subscriberModels.IcmpStatus{
					// LastReportedTime is calculated in milliseconds
					LastReportedTime: frozenClock * 1000,
					LatencyMs:        f32Ptr(12.34),
				},
			},
			State: &subscriberModels.SubscriberState{
				Mme:  mmeState,
				S1ap: s1apState,
				Spgw: spgwState,
				Mobility: []*subscriberModels.SubscriberIPAllocation{
					{
						Apn: "magma.apn",
						IP:  "192.168.128.134",
					},
					{
						Apn: "oai.ipv4",
						IP:  "192.168.128.174",
					},
				},
				Directory: &subscriberModels.SubscriberDirectoryRecord{
					LocationHistory: []string{"foo", "bar"},
				},
			},
		},
	}
	tests.RunUnitTest(t, e, tc)
}

// TestGetSubscriberByExactIMSI is a regression test to ensure we are loading
// states with the exact same IMSI key as the subscriber
func TestGetSubscriberByExactIMSI(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	getSubscriberURL := "/magma/v1/lte/:network_id/subscribers/:subscriber_id"
	createSubscribersURL := "/magma/v1/lte/:network_id/subscribers"
	handlers := handlers.GetHandlers()
	getSubscriber := tests.GetHandlerByPathAndMethod(t, handlers, getSubscriberURL, obsidian.GET).HandlerFunc
	createSubscribers := tests.GetHandlerByPathAndMethod(t, handlers, createSubscribersURL, obsidian.POST).HandlerFunc

	// Pre: create two APNs
	_, err = configurator.CreateEntities(context.Background(), "n1", []configurator.NetworkEntity{
		{Type: lte.APNEntityType, Key: "apn0"},
		{Type: lte.APNEntityType, Key: "apn1"},
	}, serdes.Entity)
	assert.NoError(t, err)

	// Create two subscribers, one a prefix of the other
	sub1 := newMutableSubscriber("IMSI1234567890")
	sub2 := newMutableSubscriber("IMSI123456789000")
	sub2.Name = "John Doe"
	payload := subscriberModels.MutableSubscribers{sub1, sub2}

	tc := tests.Test{
		Method:         "POST",
		URL:            createSubscribersURL,
		Payload:        tests.JSONMarshaler(payload),
		Handler:        createSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            getSubscriberURL,
		Handler:        getSubscriber,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 200,
		ExpectedResult: sub1.ToSubscriber(),
	}
	tests.RunUnitTest(t, e, tc)

	// Now create some AGW-reported state for the subscribers
	// First we need to register a gateway which can report state
	testUtilConfigurator.RegisterGateway(t, "n1", "g1", &models.GatewayDevice{HardwareID: "hw1"})
	ctx := test_utils.GetContextWithCertificate(t, "hw1")

	// Sub1 and sub2 differ in their mme states
	mmeState1 := state.ArbitraryJSON{"mme": "foo"}
	test_utils.ReportState(t, ctx, lte.MMEStateType, "IMSI1234567890", &mmeState1, serdes.State)
	mmeState2 := state.ArbitraryJSON{"mme": "fee"}
	test_utils.ReportState(t, ctx, lte.MMEStateType, "IMSI123456789000", &mmeState2, serdes.State)

	sub1ExpectedResult := sub1.ToSubscriber()
	sub1ExpectedResult.Monitoring = &subscriberModels.SubscriberStatus{}
	sub1ExpectedResult.State = &subscriberModels.SubscriberState{
		Mme: mmeState1,
	}

	// Should only report states for IMSI1234567890, not IMSI123456789000
	tc = tests.Test{
		Method:         "GET",
		URL:            getSubscriberURL,
		Handler:        getSubscriber,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 200,
		ExpectedResult: sub1ExpectedResult,
	}
	tests.RunUnitTest(t, e, tc)
}

func TestListSubscriberStates(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n0"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/subscriber_state"
	listSubscribers := tests.GetHandlerByPathAndMethod(t, handlers.GetHandlers(), testURLRoot, obsidian.GET).HandlerFunc

	// Initially no state
	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*subscriberModels.SubscriberState{}),
	}
	tests.RunUnitTest(t, e, tc)

	// Create gateway and report states
	_, err = configurator.CreateEntity(context.Background(), "n0", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g0", Config: &models.MagmadGatewayConfigs{}, PhysicalID: "hw0"}, serdes.Entity)
	assert.NoError(t, err)
	frozenClock := int64(1000000)
	clock.SetAndFreezeClock(t, time.Unix(frozenClock, 0))
	defer clock.UnfreezeClock(t)

	icmpStatus := &subscriberModels.IcmpStatus{LatencyMs: f32Ptr(12.34)}
	ctx := test_utils.GetContextWithCertificate(t, "hw0")
	test_utils.ReportState(t, ctx, lte.ICMPStateType, "IMSI1234567890", icmpStatus, serdes.State)
	mmeState := state.ArbitraryJSON{"mme": "foo"}
	test_utils.ReportState(t, ctx, lte.MMEStateType, "IMSI1234567890", &mmeState, serdes.State)
	spgwState := state.ArbitraryJSON{"spgw": "foo"}
	test_utils.ReportState(t, ctx, lte.SPGWStateType, "IMSI1234567890", &spgwState, serdes.State)
	s1apState := state.ArbitraryJSON{"s1ap": "foo"}
	test_utils.ReportState(t, ctx, lte.S1APStateType, "IMSI1234567890", &s1apState, serdes.State)
	// Report 2 allocated IP addresses for the subscriber
	mobilitydState1 := state.ArbitraryJSON{
		"ip": map[string]interface{}{
			"address": "wKiArg==",
		},
	}
	mobilitydState2 := state.ArbitraryJSON{
		"ip": map[string]interface{}{
			"address": "wKiAhg==",
		},
	}
	test_utils.ReportState(t, ctx, lte.MobilitydStateType, "IMSI1234567890.oai.ipv4", &mobilitydState1, serdes.State)
	test_utils.ReportState(t, ctx, lte.MobilitydStateType, "IMSI1234567890.magma.apn", &mobilitydState2, serdes.State)
	directoryState := directorydTypes.DirectoryRecord{LocationHistory: []string{"foo", "bar"}}
	test_utils.ReportState(t, ctx, orc8r.DirectoryRecordType, "IMSI1234567890", &directoryState, serdes.State)

	subState0 := state.ArbitraryJSON{
		"subscriber_state": state.ArbitraryJSON{
			"imsi":     "IMSI1234567890",
			"sessions": []string{"p0", "p1"},
		},
		"session_state": state.ArbitraryJSON{
			"apn":  "apn0",
			"ipv4": "168.212.226.204",
		},
	}
	test_utils.ReportState(t, ctx, lte.SubscriberStateType, "IMSI1234567890", &subState0, serdes.State)
	subState1 := state.ArbitraryJSON{
		"subscriber_state": state.ArbitraryJSON{
			"imsi":     "IMSI0987654321",
			"sessions": []string{"p2"},
		},
		"session_state": state.ArbitraryJSON{
			"apn":  "apn1",
			"ipv4": "168.212.226.203",
		},
	}
	test_utils.ReportState(t, ctx, lte.SubscriberStateType, "IMSI0987654321", &subState1, serdes.State)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*subscriberModels.SubscriberState{
			"IMSI1234567890": {
				SubscriberState: subState0,
				Mme:             mmeState,
				S1ap:            s1apState,
				Spgw:            spgwState,
				Mobility: []*subscriberModels.SubscriberIPAllocation{
					{
						Apn: "magma.apn",
						IP:  "192.168.128.134",
					},
					{
						Apn: "oai.ipv4",
						IP:  "192.168.128.174",
					},
				},
				Directory: &subscriberModels.SubscriberDirectoryRecord{
					LocationHistory: []string{"foo", "bar"},
				},
			},
			"IMSI0987654321": {
				SubscriberState: subState1,
			},
		}),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestGetSubscriberState(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n0"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/subscriber_state/:subscriber_id"
	getSubscriber := tests.GetHandlerByPathAndMethod(t, handlers.GetHandlers(), testURLRoot, obsidian.GET).HandlerFunc

	// Initially no state
	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getSubscriber,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n0", "IMSI1234567890"},
		ExpectedStatus: 200,
		ExpectedResult: &subscriberModels.SubscriberState{},
	}
	tests.RunUnitTest(t, e, tc)

	// Create gateway and report states
	_, err = configurator.CreateEntity(context.Background(), "n0", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g0", Config: &models.MagmadGatewayConfigs{}, PhysicalID: "hw0"}, serdes.Entity)
	assert.NoError(t, err)
	frozenClock := int64(1000000)
	clock.SetAndFreezeClock(t, time.Unix(frozenClock, 0))
	defer clock.UnfreezeClock(t)
	icmpStatus := &subscriberModels.IcmpStatus{LatencyMs: f32Ptr(12.34)}
	ctx := test_utils.GetContextWithCertificate(t, "hw0")
	test_utils.ReportState(t, ctx, lte.ICMPStateType, "IMSI1234567890", icmpStatus, serdes.State)
	mmeState := state.ArbitraryJSON{"mme": "foo"}
	test_utils.ReportState(t, ctx, lte.MMEStateType, "IMSI1234567890", &mmeState, serdes.State)
	spgwState := state.ArbitraryJSON{"spgw": "foo"}
	test_utils.ReportState(t, ctx, lte.SPGWStateType, "IMSI1234567890", &spgwState, serdes.State)
	s1apState := state.ArbitraryJSON{"s1ap": "foo"}
	test_utils.ReportState(t, ctx, lte.S1APStateType, "IMSI1234567890", &s1apState, serdes.State)
	// Report 2 allocated IP addresses for the subscriber
	mobilitydState1 := state.ArbitraryJSON{
		"ip": map[string]interface{}{
			"address": "wKiArg==",
		},
	}
	mobilitydState2 := state.ArbitraryJSON{
		"ip": map[string]interface{}{
			"address": "wKiAhg==",
		},
	}
	test_utils.ReportState(t, ctx, lte.MobilitydStateType, "IMSI1234567890.oai.ipv4", &mobilitydState1, serdes.State)
	test_utils.ReportState(t, ctx, lte.MobilitydStateType, "IMSI1234567890.magma.apn", &mobilitydState2, serdes.State)
	directoryState := directorydTypes.DirectoryRecord{LocationHistory: []string{"foo", "bar"}}
	test_utils.ReportState(t, ctx, orc8r.DirectoryRecordType, "IMSI1234567890", &directoryState, serdes.State)
	subState := state.ArbitraryJSON{
		"subscriber_state": state.ArbitraryJSON{
			"imsi":     "IMSI1234567890",
			"sessions": []string{"p0", "p1"},
		},
		"session_state": state.ArbitraryJSON{
			"apn":  "apn0",
			"ipv4": "168.212.226.204",
		},
	}
	test_utils.ReportState(t, ctx, lte.SubscriberStateType, "IMSI1234567890", &subState, serdes.State)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getSubscriber,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n0", "IMSI1234567890"},
		ExpectedStatus: 200,
		ExpectedResult: &subscriberModels.SubscriberState{
			SubscriberState: subState,
			Mme:             mmeState,
			S1ap:            s1apState,
			Spgw:            spgwState,
			Mobility: []*subscriberModels.SubscriberIPAllocation{
				{
					Apn: "magma.apn",
					IP:  "192.168.128.134",
				},
				{
					Apn: "oai.ipv4",
					IP:  "192.168.128.174",
				},
			},
			Directory: &subscriberModels.SubscriberDirectoryRecord{
				LocationHistory: []string{"foo", "bar"},
			},
		},
	}
	tests.RunUnitTest(t, e, tc)
}

func TestGetSubscriberByMSISDN(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	subscriberdbTestInit.StartTestService(t)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n0"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	subscriberdbHandlers := handlers.GetHandlers()

	subURLBase := "/magma/v1/lte/:network_id/subscribers"
	getAllSubscribers := tests.GetHandlerByPathAndMethod(t, subscriberdbHandlers, subURLBase, obsidian.GET).HandlerFunc

	msisdnURLBase := "/magma/v1/lte/:network_id/msisdns"
	msisdnURLManage := "/magma/v1/lte/:network_id/msisdns/:msisdn"
	getAllMSISDNs := tests.GetHandlerByPathAndMethod(t, subscriberdbHandlers, msisdnURLBase, obsidian.GET).HandlerFunc
	postMSISDN := tests.GetHandlerByPathAndMethod(t, subscriberdbHandlers, msisdnURLBase, obsidian.POST).HandlerFunc
	getMSISDN := tests.GetHandlerByPathAndMethod(t, subscriberdbHandlers, msisdnURLManage, obsidian.GET).HandlerFunc
	deleteMSISDN := tests.GetHandlerByPathAndMethod(t, subscriberdbHandlers, msisdnURLManage, obsidian.DELETE).HandlerFunc

	// MSISDNs initially empty
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/msisdns",
		Handler:        getAllMSISDNs,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]string{}),
	}
	tests.RunUnitTest(t, e, tc)

	// List all => empty
	emptyPaginatedSub := subscriberModels.PaginatedSubscribers{
		TotalCount:    int64(0),
		NextPageToken: "",
		Subscribers:   map[string]*subscriberModels.Subscriber{},
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            subURLBase,
		Handler:        getAllSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(emptyPaginatedSub),
	}
	tests.RunUnitTest(t, e, tc)

	// List one => 404
	tc = tests.Test{
		Method:         "GET",
		URL:            subURLBase + "?msisdn=13109976224",
		Handler:        getAllSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	// Create default subscriber profile
	_, err = configurator.CreateEntity(context.Background(), "n0", configurator.NetworkEntity{
		Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
		Name: "Jane Doe",
	}, serdes.Entity)
	assert.NoError(t, err)

	// Create MSISDN->IMSI mapping
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n0/msisdns",
		Payload:        &subscriberModels.MsisdnAssignment{ID: "IMSI1234567890", Msisdn: "msisdn0"},
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		Handler:        postMSISDN,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Create another MSISDN->IMSI mapping
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n0/msisdns",
		Payload:        &subscriberModels.MsisdnAssignment{ID: "IMSI9999999999", Msisdn: "msisdn1"},
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		Handler:        postMSISDN,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Get MSISDN => MSISDN exists
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/msisdns/msisdn0",
		Handler:        getMSISDN,
		ParamNames:     []string{"network_id", "msisdn"},
		ParamValues:    []string{"n0", "msisdn0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler("IMSI1234567890"),
	}
	tests.RunUnitTest(t, e, tc)

	// Get subscriber by MSISDN
	tc = tests.Test{
		Method:         "GET",
		URL:            subURLBase + "?msisdn=msisdn0",
		Handler:        getAllSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*subscriberModels.Subscriber{
			"IMSI1234567890": {
				ID:     "IMSI1234567890",
				Name:   "Jane Doe",
				Config: &subscriberModels.SubscriberConfig{Lte: nil},
			},
		}),
	}
	tests.RunUnitTest(t, e, tc)

	// Delete MSISDN->IMSI mapping
	tc = tests.Test{
		Method:         "DELETE",
		URL:            msisdnURLManage,
		Handler:        deleteMSISDN,
		ParamNames:     []string{"network_id", "msisdn"},
		ParamValues:    []string{"n0", "msisdn0"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Get MSISDN => 404
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/msisdns/msisdn0",
		Handler:        getMSISDN,
		ParamNames:     []string{"network_id", "msisdn"},
		ParamValues:    []string{"n0", "msisdn0"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	// Get subscriber by MSISDN => 404
	tc = tests.Test{
		Method:         "GET",
		URL:            subURLBase + "?msisdn=msisdn0",
		Handler:        getAllSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestGetSubscriberByIP(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	subscriberdbTestInit.StartTestService(t)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n0"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	subscriberdbHandlers := handlers.GetHandlers()

	subURLBase := "/magma/v1/lte/:network_id/subscribers"
	getAllSubscribers := tests.GetHandlerByPathAndMethod(t, subscriberdbHandlers, subURLBase, obsidian.GET).HandlerFunc

	// List one => none found
	tc := tests.Test{
		Method:         "GET",
		URL:            subURLBase + "?ip=127.0.0.1",
		Handler:        getAllSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*subscriberModels.Subscriber{}),
	}
	tests.RunUnitTest(t, e, tc)

	// Create default subscriber profiles
	_, err = configurator.CreateEntities(context.Background(), "n0", []configurator.NetworkEntity{
		{
			Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
			Name: "Jane Doe",
		},
		{
			Type: lte.SubscriberEntityType, Key: "IMSI0000000001",
			Name: "John Doe",
		},
	}, serdes.Entity)
	assert.NoError(t, err)

	// List one => still not found
	tc = tests.Test{
		Method:         "GET",
		URL:            subURLBase + "?ip=127.0.0.1",
		Handler:        getAllSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*subscriberModels.Subscriber{}),
	}
	tests.RunUnitTest(t, e, tc)

	// Report IP state: Jane has an IP
	_, err = configurator.CreateEntity(context.Background(), "n0", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g0", Config: &models.MagmadGatewayConfigs{}, PhysicalID: "hw0"}, serdes.Entity)
	assert.NoError(t, err)
	ctx := test_utils.GetContextWithCertificate(t, "hw0")
	mobilitydState := &state.ArbitraryJSON{"ip": state.ArbitraryJSON{"address": "fwAAAQ=="}} // 127.0.0.1
	test_utils.ReportState(t, ctx, lte.MobilitydStateType, "IMSI1234567890.oai.ipv4", mobilitydState, serdes.State)
	// Wait for state to be indexed
	time.Sleep(time.Second)

	// List one => found
	tc = tests.Test{
		Method:         "GET",
		URL:            subURLBase + "?ip=127.0.0.1",
		Handler:        getAllSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*subscriberModels.Subscriber{
			"IMSI1234567890": {
				ID:         "IMSI1234567890",
				Name:       "Jane Doe",
				Config:     &subscriberModels.SubscriberConfig{Lte: nil},
				Monitoring: &subscriberModels.SubscriberStatus{},
				State: &subscriberModels.SubscriberState{
					Mobility: []*subscriberModels.SubscriberIPAllocation{{Apn: "oai.ipv4", IP: "127.0.0.1"}},
				},
			},
		}),
	}
	tests.RunUnitTest(t, e, tc)

	// Report IP state: Jane has new IP
	mobilitydState = &state.ArbitraryJSON{"ip": state.ArbitraryJSON{"address": "fwAAAg=="}} // 127.0.0.2
	test_utils.ReportState(t, ctx, lte.MobilitydStateType, "IMSI1234567890.oai.ipv4", mobilitydState, serdes.State)
	// Wait for state to be indexed
	time.Sleep(time.Second)

	// List one => IP changed means not found
	tc = tests.Test{
		Method:         "GET",
		URL:            subURLBase + "?ip=127.0.0.1",
		Handler:        getAllSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*subscriberModels.Subscriber{}),
	}
	tests.RunUnitTest(t, e, tc)

	// Report IP state: John has Jane's new IP
	mobilitydState = &state.ArbitraryJSON{"ip": state.ArbitraryJSON{"address": "fwAAAg=="}} // 127.0.0.2
	test_utils.ReportState(t, ctx, lte.MobilitydStateType, "IMSI0000000001.oai.ipv4", mobilitydState, serdes.State)
	// Wait for state to be indexed
	time.Sleep(time.Second)

	// Delete Jane's new IP
	err = state.DeleteStates(context.Background(), "n0", []stateTypes.ID{{Type: lte.MobilitydStateType, DeviceID: "IMSI1234567890.oai.ipv4"}})
	assert.NoError(t, err)

	// List one => found John (and only John -- Jane should be filtered out)
	tc = tests.Test{
		Method:         "GET",
		URL:            subURLBase + "?ip=127.0.0.2",
		Handler:        getAllSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*subscriberModels.Subscriber{
			"IMSI0000000001": {
				ID:         "IMSI0000000001",
				Name:       "John Doe",
				Config:     &subscriberModels.SubscriberConfig{Lte: nil},
				Monitoring: &subscriberModels.SubscriberStatus{},
				State: &subscriberModels.SubscriberState{
					Mobility: []*subscriberModels.SubscriberIPAllocation{{Apn: "oai.ipv4", IP: "127.0.0.2"}},
				},
			},
		}),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateSubscriber(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/subscribers/:subscriber_id"
	handlers := handlers.GetHandlers()
	updateSubscriber := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.PUT).HandlerFunc

	//preseed 2 apns
	apn1, apn2 := "foo", "bar"
	_, err = configurator.CreateEntities(context.Background(), "n1", []configurator.NetworkEntity{
		{Type: lte.APNEntityType, Key: apn1},
		{Type: lte.APNEntityType, Key: apn2},
	}, serdes.Entity)
	assert.NoError(t, err)

	// 404
	payload := &subscriberModels.MutableSubscriber{
		ID: "IMSI1234567890",
		Lte: &subscriberModels.LteSubscription{
			AuthAlgo:   "MILENAGE",
			AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			State:      "ACTIVE",
			SubProfile: "default",
		},
		ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"5GC", "EPC"},
		ActiveApns:            subscriberModels.ApnList{apn2, apn1},
	}
	tc := tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Handler:        updateSubscriber,
		Payload:        payload,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	// Happy path
	err = configurator.UpdateNetworkConfig(context.Background(), "n1", lte.CellularNetworkConfigType, &lteModels.NetworkCellularConfigs{
		Epc: &lteModels.NetworkEpcConfigs{
			SubProfiles: map[string]lteModels.NetworkEpcConfigsSubProfilesAnon{
				"foo": {
					MaxUlBitRate: 100,
					MaxDlBitRate: 100,
				},
			},
		},
	}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{
		Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
		Config: &subscriberModels.SubscriberConfig{
			Lte: &subscriberModels.LteSubscription{
				AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:      "ACTIVE",
				SubProfile: "default",
			},
			ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"5GC", "EPC"},
		},
		Associations: storage.TKs{{Type: lte.APNEntityType, Key: apn2}},
	}, serdes.Entity)
	assert.NoError(t, err)

	payload = &subscriberModels.MutableSubscriber{
		ID:   "IMSI1234567890",
		Name: "Jane Doe",
		Lte: &subscriberModels.LteSubscription{
			AuthAlgo:   "MILENAGE",
			AuthKey:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
			AuthOpc:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
			State:      "INACTIVE",
			SubProfile: "foo",
		},
		ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{},
		StaticIps:             subscriberModels.SubscriberStaticIps{apn1: "192.168.100.1"},
		ActiveApns:            subscriberModels.ApnList{apn2, apn1},
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Handler:        updateSubscriber,
		Payload:        payload,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadEntity(context.Background(), "n1", lte.SubscriberEntityType, "IMSI1234567890", configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.NoError(t, err)
	expected := configurator.NetworkEntity{
		NetworkID:    "n1",
		Type:         lte.SubscriberEntityType,
		Key:          "IMSI1234567890",
		Name:         "Jane Doe",
		Config:       &subscriberModels.SubscriberConfig{Lte: payload.Lte, StaticIps: payload.StaticIps},
		GraphID:      "2",
		Version:      1,
		Associations: storage.TKs{{Type: lte.APNEntityType, Key: apn2}, {Type: lte.APNEntityType, Key: apn1}},
	}
	assert.Equal(t, expected, actual)

	// No profile matching
	payload.Lte.SubProfile = "bar"
	tc = tests.Test{
		Method:                 "PUT",
		URL:                    testURLRoot,
		Handler:                updateSubscriber,
		Payload:                payload,
		ParamNames:             []string{"network_id", "subscriber_id"},
		ParamValues:            []string{"n1", "IMSI1234567890"},
		ExpectedStatus:         400,
		ExpectedErrorSubstring: "subscriber profile 'bar' does not exist for the network",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestDeleteSubscriber(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/subscribers/:subscriber_id"
	handlers := handlers.GetHandlers()
	deleteSubscriber := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.DELETE).HandlerFunc

	//preseed 2 apns
	apn1, apn2 := "foo", "bar"
	_, err = configurator.CreateEntities(context.Background(), "n1", []configurator.NetworkEntity{
		{Type: lte.APNEntityType, Key: apn1},
		{Type: lte.APNEntityType, Key: apn2},
	}, serdes.Entity)
	assert.NoError(t, err)

	_, err = configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{
		Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
		// Intentionally populate with invalid config
		Config: &subscriberModels.LteSubscription{
			AuthAlgo: "MILENAGE",
			AuthKey:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			AuthOpc:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			State:    "ACTIVE",
		},
		Associations: storage.TKs{{Type: lte.APNEntityType, Key: apn2}, {Type: lte.APNEntityType, Key: apn1}},
	}, serdes.Entity)
	assert.NoError(t, err)

	tc := tests.Test{
		Method:         "DELETE",
		URL:            testURLRoot,
		Handler:        deleteSubscriber,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actual, _, err := configurator.LoadAllEntitiesOfType(context.Background(), "n1", lte.SubscriberEntityType, configurator.EntityLoadCriteria{}, serdes.Entity)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(actual))
}

func TestActivateDeactivateSubscriber(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/subscribers/:subscriber_id"
	handlers := handlers.GetHandlers()
	activateSubscriber := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot+"/activate", obsidian.POST).HandlerFunc
	deactivateSubscriber := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot+"/deactivate", obsidian.POST).HandlerFunc

	//preseed 2 apns
	apn1, apn2 := "foo", "bar"
	_, err = configurator.CreateEntities(context.Background(), "n1", []configurator.NetworkEntity{
		{Type: lte.APNEntityType, Key: apn1},
		{Type: lte.APNEntityType, Key: apn2},
	}, serdes.Entity)
	assert.NoError(t, err)

	expected := configurator.NetworkEntity{
		Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
		Config: &subscriberModels.SubscriberConfig{
			Lte: &subscriberModels.LteSubscription{
				AuthAlgo: "MILENAGE",
				AuthKey:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:    "ACTIVE",
			},
		},
		Associations: storage.TKs{{Type: lte.APNEntityType, Key: apn2}, {Type: lte.APNEntityType, Key: apn1}},
	}
	_, err = configurator.CreateEntity(context.Background(), "n1", expected, serdes.Entity)
	assert.NoError(t, err)
	expected.NetworkID = "n1"
	expected.GraphID = "2"
	expected.Version = 1

	// activate already activated subscriber
	tc := tests.Test{
		Method:         "POST",
		URL:            testURLRoot + "/activate",
		Handler:        activateSubscriber,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadEntity(context.Background(), "n1", lte.SubscriberEntityType, "IMSI1234567890", configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// deactivate
	tc.URL = testURLRoot + "/deactivate"
	tc.Handler = deactivateSubscriber
	tests.RunUnitTest(t, e, tc)

	actual, err = configurator.LoadEntity(context.Background(), "n1", lte.SubscriberEntityType, "IMSI1234567890", configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.NoError(t, err)
	expected.Config.(*subscriberModels.SubscriberConfig).Lte.State = "INACTIVE"
	expected.Version = 2
	assert.Equal(t, expected, actual)

	// deactivate deactivated sub
	tests.RunUnitTest(t, e, tc)
	actual, err = configurator.LoadEntity(context.Background(), "n1", lte.SubscriberEntityType, "IMSI1234567890", configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.NoError(t, err)
	expected.Config.(*subscriberModels.SubscriberConfig).Lte.State = "INACTIVE"
	expected.Version = 3
	assert.Equal(t, expected, actual)

	// activate
	tc.URL = testURLRoot + "/activate"
	tc.Handler = activateSubscriber
	tests.RunUnitTest(t, e, tc)
	actual, err = configurator.LoadEntity(context.Background(), "n1", lte.SubscriberEntityType, "IMSI1234567890", configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.NoError(t, err)
	expected.Config.(*subscriberModels.SubscriberConfig).Lte.State = "ACTIVE"
	expected.Version = 4
	assert.Equal(t, expected, actual)
}

func TestUpdateSubscriberProfile(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)

	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	err = configurator.UpdateNetworkConfig(context.Background(), "n1", lte.CellularNetworkConfigType, &lteModels.NetworkCellularConfigs{
		Epc: &lteModels.NetworkEpcConfigs{
			SubProfiles: map[string]lteModels.NetworkEpcConfigsSubProfilesAnon{
				"foo": {
					MaxUlBitRate: 100,
					MaxDlBitRate: 100,
				},
			},
		},
	}, serdes.Network)
	assert.NoError(t, err)

	//preseed 2 apns
	apn1, apn2 := "foo", "bar"
	_, err = configurator.CreateEntities(context.Background(), "n1", []configurator.NetworkEntity{
		{Type: lte.APNEntityType, Key: apn1},
		{Type: lte.APNEntityType, Key: apn2},
	}, serdes.Entity)
	assert.NoError(t, err)

	_, err = configurator.CreateEntity(context.Background(), "n1", configurator.NetworkEntity{
		Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
		Config: &subscriberModels.SubscriberConfig{
			Lte: &subscriberModels.LteSubscription{
				AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:      "ACTIVE",
				SubProfile: "default",
			},
			ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"5GC"},
		},
		Associations: storage.TKs{{Type: lte.APNEntityType, Key: apn2}, {Type: lte.APNEntityType, Key: apn1}},
	}, serdes.Entity)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/subscribers/:subscriber_id/lte/sub_profile"
	handlers := handlers.GetHandlers()
	updateProfile := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.PUT).HandlerFunc

	// 404
	payload := "foo"
	tc := tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Handler:        updateProfile,
		Payload:        tests.JSONMarshaler(payload),
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI0987654321"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	// bad profile
	payload = "bar"
	tc = tests.Test{
		Method:                 "PUT",
		URL:                    testURLRoot,
		Handler:                updateProfile,
		Payload:                tests.JSONMarshaler(payload),
		ParamNames:             []string{"network_id", "subscriber_id"},
		ParamValues:            []string{"n1", "IMSI1234567890"},
		ExpectedStatus:         400,
		ExpectedErrorSubstring: "subscriber profile 'bar' does not exist for the network",
	}
	tests.RunUnitTest(t, e, tc)

	// happy path
	payload = "foo"
	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Handler:        updateProfile,
		Payload:        tests.JSONMarshaler(payload),
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadEntity(context.Background(), "n1", lte.SubscriberEntityType, "IMSI1234567890", configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.NoError(t, err)
	expected := configurator.NetworkEntity{
		NetworkID: "n1", Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
		Config: &subscriberModels.SubscriberConfig{
			Lte: &subscriberModels.LteSubscription{
				AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:      "ACTIVE",
				SubProfile: "foo",
			},
			ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"5GC"},
		},
		GraphID:      "2",
		Version:      1,
		Associations: storage.TKs{{Type: lte.APNEntityType, Key: apn2}, {Type: lte.APNEntityType, Key: apn1}},
	}
	assert.Equal(t, expected, actual)

	// set to default
	payload = "default"
	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Handler:        updateProfile,
		Payload:        tests.JSONMarshaler(payload),
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err = configurator.LoadEntity(context.Background(), "n1", lte.SubscriberEntityType, "IMSI1234567890", configurator.FullEntityLoadCriteria(), serdes.Entity)
	assert.NoError(t, err)
	expected = configurator.NetworkEntity{
		NetworkID: "n1", Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
		Config: &subscriberModels.SubscriberConfig{
			Lte: &subscriberModels.LteSubscription{
				AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:      "ACTIVE",
				SubProfile: "default",
			},
			ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"5GC"},
		},
		GraphID:      "2",
		Version:      2,
		Associations: storage.TKs{{Type: lte.APNEntityType, Key: apn2}, {Type: lte.APNEntityType, Key: apn1}},
	}
	assert.Equal(t, expected, actual)
}

func TestSubscriberBasename(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n0"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntities(context.Background(), "n0", []configurator.NetworkEntity{
		{Type: lte.APNEntityType, Key: "apn0"},
		{Type: lte.APNEntityType, Key: "apn1"},
		{Type: lte.BaseNameEntityType, Key: "basename0"},
		{Type: lte.BaseNameEntityType, Key: "basename2"},
	}, serdes.Entity)
	assert.NoError(t, err)

	e := echo.New()
	urlBase := "/magma/v1/lte/:network_id/subscribers"
	urlManage := urlBase + "/:subscriber_id"
	subscriberdbHandlers := handlers.GetHandlers()
	getAllSubscribers := tests.GetHandlerByPathAndMethod(t, subscriberdbHandlers, urlBase, obsidian.GET).HandlerFunc
	postSubscriber := tests.GetHandlerByPathAndMethod(t, subscriberdbHandlers, urlBase, obsidian.POST).HandlerFunc
	putSubscriber := tests.GetHandlerByPathAndMethod(t, subscriberdbHandlers, urlManage, obsidian.PUT).HandlerFunc

	imsi := "IMSI1234567890"
	mutableSub := newMutableSubscriber(imsi)

	// Post a policy association with a non existent policy
	mutableSub.ActiveBaseNames = policydbModels.BaseNames{"baseXXX"}
	tc := tests.Test{
		Method:                 "POST",
		URL:                    "/magma/v1/lte/n0/subscribers",
		Payload:                tests.JSONMarshaler(subscriberModels.MutableSubscribers{mutableSub}),
		Handler:                postSubscriber,
		ParamNames:             []string{"network_id"},
		ParamValues:            []string{"n0"},
		ExpectedStatus:         500, // would make more sense as 400
		ExpectedErrorSubstring: `code=500, message=could not find entities matching [type:"base_name" key:"baseXXX" ]`,
	}
	tests.RunUnitTest(t, e, tc)

	// Successful post request
	mutableSub.ActiveBaseNames = policydbModels.BaseNames{"basename0"}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte/n0/subscribers",
		Payload:        tests.JSONMarshaler(subscriberModels.MutableSubscribers{mutableSub}),
		Handler:        postSubscriber,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Get all, posted subscriber found
	expected := subscriberModels.PaginatedSubscribers{
		TotalCount:    1,
		NextPageToken: "",
		Subscribers: map[string]*subscriberModels.Subscriber{
			imsi: mutableSub.ToSubscriber(),
		},
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/subscribers",
		Handler:        getAllSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expected),
	}
	tests.RunUnitTest(t, e, tc)

	// Put request with non existent basename
	mutableSub.ActiveBaseNames = policydbModels.BaseNames{"baseXXX"}
	tc = tests.Test{
		Method:                 "PUT",
		URL:                    "/magma/v1/lte/n0/subscribers/" + imsi,
		Payload:                mutableSub,
		ParamNames:             []string{"network_id", "subscriber_id"},
		ParamValues:            []string{"n0", imsi},
		Handler:                putSubscriber,
		ExpectedStatus:         500, // would make more sense as 400
		ExpectedErrorSubstring: `could not find entities matching [type:"base_name" key:"baseXXX" ]`,
	}
	tests.RunUnitTest(t, e, tc)

	// Successful put request
	mutableSub.ActiveBaseNames = policydbModels.BaseNames{"basename2"}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/lte/n0/subscribers/" + imsi,
		Payload:        mutableSub,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n0", imsi},
		Handler:        putSubscriber,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Get all, updated subscriber matches the expected value
	expected = subscriberModels.PaginatedSubscribers{
		TotalCount:    1,
		NextPageToken: "",
		Subscribers: map[string]*subscriberModels.Subscriber{
			imsi: mutableSub.ToSubscriber(),
		},
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/subscribers",
		Handler:        getAllSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expected),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestSubscriberPolicy(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n0"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntities(context.Background(), "n0", []configurator.NetworkEntity{
		{Type: lte.APNEntityType, Key: "apn0"},
		{Type: lte.APNEntityType, Key: "apn1"},
		{Type: lte.PolicyRuleEntityType, Key: "rule0"},
		{Type: lte.PolicyRuleEntityType, Key: "rule1"},
		{Type: lte.PolicyRuleEntityType, Key: "rule2"},
	}, serdes.Entity)
	assert.NoError(t, err)

	e := echo.New()
	urlBase := "/magma/v1/lte/:network_id/subscribers"
	urlManage := urlBase + "/:subscriber_id"
	subscriberdbHandlers := handlers.GetHandlers()
	getAllSubscribers := tests.GetHandlerByPathAndMethod(t, subscriberdbHandlers, urlBase, obsidian.GET).HandlerFunc
	postSubscriber := tests.GetHandlerByPathAndMethod(t, subscriberdbHandlers, urlBase, obsidian.POST).HandlerFunc
	putSubscriber := tests.GetHandlerByPathAndMethod(t, subscriberdbHandlers, urlManage, obsidian.PUT).HandlerFunc

	imsi := "IMSI1234567890"
	mutableSub := newMutableSubscriber(imsi)

	// Post a policy association with a non existent policy
	mutableSub.ActivePolicies = policydbModels.PolicyIds{"ruleXXX"}
	tc := tests.Test{
		Method:                 "POST",
		URL:                    "/magma/v1/lte/n0/subscribers",
		Payload:                tests.JSONMarshaler(subscriberModels.MutableSubscribers{mutableSub}),
		Handler:                postSubscriber,
		ParamNames:             []string{"network_id"},
		ParamValues:            []string{"n0"},
		ExpectedStatus:         500, // would make more sense as 400
		ExpectedErrorSubstring: `code=500, message=could not find entities matching [type:"policy" key:"ruleXXX" ]`,
	}
	tests.RunUnitTest(t, e, tc)

	// Successful post request
	mutableSub.ActivePolicies = policydbModels.PolicyIds{"rule0"}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte/n0/subscribers",
		Payload:        tests.JSONMarshaler(subscriberModels.MutableSubscribers{mutableSub}),
		Handler:        postSubscriber,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Get all, posted subscriber found
	expected := subscriberModels.PaginatedSubscribers{
		TotalCount:    1,
		NextPageToken: "",
		Subscribers: map[string]*subscriberModels.Subscriber{
			imsi: mutableSub.ToSubscriber(),
		},
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/subscribers",
		Handler:        getAllSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expected),
	}
	tests.RunUnitTest(t, e, tc)

	// Put request with non existent policy
	mutableSub.ActivePolicies = policydbModels.PolicyIds{"ruleXXX"}
	tc = tests.Test{
		Method:                 "PUT",
		URL:                    "/magma/v1/lte/n0/subscribers/" + imsi,
		Payload:                mutableSub,
		ParamNames:             []string{"network_id", "subscriber_id"},
		ParamValues:            []string{"n0", imsi},
		Handler:                putSubscriber,
		ExpectedStatus:         500, // would make more sense as 400
		ExpectedErrorSubstring: `could not find entities matching [type:"policy" key:"ruleXXX" ]`,
	}
	tests.RunUnitTest(t, e, tc)

	// Successful put request
	mutableSub.ActivePolicies = policydbModels.PolicyIds{"rule2"}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/lte/n0/subscribers/" + imsi,
		Payload:        mutableSub,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n0", imsi},
		Handler:        putSubscriber,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Get all, updated subscriber matches the expected value
	expected = subscriberModels.PaginatedSubscribers{
		TotalCount:    1,
		NextPageToken: "",
		Subscribers: map[string]*subscriberModels.Subscriber{
			imsi: mutableSub.ToSubscriber(),
		}}
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/subscribers",
		Handler:        getAllSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expected),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestAPNPolicyProfile(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	err := configurator.CreateNetwork(context.Background(), configurator.Network{ID: "n0"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntities(context.Background(), "n0", []configurator.NetworkEntity{
		{Type: lte.APNEntityType, Key: "apn0"},
		{Type: lte.APNEntityType, Key: "apn1"},
		{Type: lte.PolicyRuleEntityType, Key: "rule0"},
		{Type: lte.PolicyRuleEntityType, Key: "rule1"},
		{Type: lte.PolicyRuleEntityType, Key: "rule2"},
	}, serdes.Entity)
	assert.NoError(t, err)

	e := echo.New()
	urlBase := "/magma/v1/lte/:network_id/subscribers"
	urlManage := urlBase + "/:subscriber_id"
	subscriberdbHandlers := handlers.GetHandlers()
	getAllSubscribers := tests.GetHandlerByPathAndMethod(t, subscriberdbHandlers, urlBase, obsidian.GET).HandlerFunc
	postSubscriber := tests.GetHandlerByPathAndMethod(t, subscriberdbHandlers, urlBase, obsidian.POST).HandlerFunc
	putSubscriber := tests.GetHandlerByPathAndMethod(t, subscriberdbHandlers, urlManage, obsidian.PUT).HandlerFunc
	getSubscriber := tests.GetHandlerByPathAndMethod(t, subscriberdbHandlers, urlManage, obsidian.GET).HandlerFunc
	deleteSubscriber := tests.GetHandlerByPathAndMethod(t, subscriberdbHandlers, urlManage, obsidian.DELETE).HandlerFunc

	deleteAPN := tests.GetHandlerByPathAndMethod(t, lteHandlers.GetHandlers(), "/magma/v1/lte/:network_id/apns/:apn_name", obsidian.DELETE).HandlerFunc
	postPolicy := tests.GetHandlerByPathAndMethod(t, policydbHandlers.GetHandlers(), "/magma/v1/networks/:network_id/policies/rules", obsidian.POST).HandlerFunc
	deletePolicy := tests.GetHandlerByPathAndMethod(t, policydbHandlers.GetHandlers(), "/magma/v1/networks/:network_id/policies/rules/:rule_id", obsidian.DELETE).HandlerFunc

	imsi := "IMSI1234567890"
	imsi1 := "IMSI1234567800"
	mutableSub := newMutableSubscriber(imsi)
	sub := mutableSub.ToSubscriber()

	t.Run("dangling apn_policy_profile regression", func(t *testing.T) {
		// Post policy
		policy := newPolicy("ruleXXX")
		tc := tests.Test{
			Method:         "POST",
			URL:            "/magma/v1/networks/n0/policies/rules",
			Payload:        policy,
			ParamNames:     []string{"network_id"},
			ParamValues:    []string{"n0"},
			Handler:        postPolicy,
			ExpectedStatus: 201,
		}
		tests.RunUnitTest(t, e, tc)

		// Post, sub with same policy both static and for specific APN
		mutableSub.ActivePolicies = policydbModels.PolicyIds{policy.ID}
		mutableSub.ActivePoliciesByApn = policydbModels.PolicyIdsByApn{
			"apn0": policydbModels.PolicyIds{"ruleXXX"},
		}
		tc = tests.Test{
			Method:         "POST",
			URL:            "/magma/v1/lte/n0/subscribers",
			Payload:        tests.JSONMarshaler(subscriberModels.MutableSubscribers{mutableSub}),
			Handler:        postSubscriber,
			ParamNames:     []string{"network_id"},
			ParamValues:    []string{"n0"},
			ExpectedStatus: 201,
		}
		tests.RunUnitTest(t, e, tc)

		// Configurator confirms apn_policy_profile exists
		profiles, err := configurator.ListEntityKeys(context.Background(), "n0", lte.APNPolicyProfileEntityType)
		assert.NoError(t, err)
		assert.Len(t, profiles, 1)

		// Put, remove policy
		mutableSub.ActivePolicies = policydbModels.PolicyIds{}
		mutableSub.ActivePoliciesByApn = policydbModels.PolicyIdsByApn{
			"apn0": policydbModels.PolicyIds{},
		}
		tc = tests.Test{
			Method:         "PUT",
			URL:            "/magma/v1/lte/n0/subscribers/" + imsi,
			Payload:        mutableSub,
			ParamNames:     []string{"network_id", "subscriber_id"},
			ParamValues:    []string{"n0", imsi},
			Handler:        putSubscriber,
			ExpectedStatus: 204,
		}
		tests.RunUnitTest(t, e, tc)

		// Configurator confirms apn_policy_profile still exists
		profiles, err = configurator.ListEntityKeys(context.Background(), "n0", lte.APNPolicyProfileEntityType)
		assert.NoError(t, err)
		assert.Len(t, profiles, 1)

		// Delete
		tc = tests.Test{
			Method:         "DELETE",
			URL:            "/magma/v1/lte/n0/subscribers/" + imsi,
			ParamNames:     []string{"network_id", "subscriber_id"},
			ParamValues:    []string{"n0", imsi},
			Handler:        deleteSubscriber,
			ExpectedStatus: 204,
		}
		tests.RunUnitTest(t, e, tc)

		// Configurator confirms subscriber no longer exists
		profiles, err = configurator.ListEntityKeys(context.Background(), "n0", lte.SubscriberEntityType)
		assert.NoError(t, err)
		assert.Len(t, profiles, 0)

		// Configurator confirms apn_policy_profile no longer exists
		profiles, err = configurator.ListEntityKeys(context.Background(), "n0", lte.APNPolicyProfileEntityType)
		assert.NoError(t, err)
		assert.Len(t, profiles, 0)

		// Clean up created policy rule
		tc = tests.Test{
			Method:         "DELETE",
			URL:            "/magma/v1/networks/n1/policies/rules/rule0",
			Payload:        nil,
			ParamNames:     []string{"network_id", "rule_id"},
			ParamValues:    []string{"n0", "ruleXXX"},
			Handler:        deletePolicy,
			ExpectedStatus: 204,
		}
		tests.RunUnitTest(t, e, tc)
	})

	// Get all, initially empty
	emptySub := subscriberModels.PaginatedSubscribers{
		TotalCount:    0,
		NextPageToken: "",
		Subscribers:   map[string]*subscriberModels.Subscriber{},
	}
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/subscribers",
		Handler:        getAllSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(emptySub),
	}
	tests.RunUnitTest(t, e, tc)

	// Post err, APN doesn't exist
	mutableSub.ActivePoliciesByApn = policydbModels.PolicyIdsByApn{"apnXXX": policydbModels.PolicyIds{"rule0"}}
	tc = tests.Test{
		Method:                 "POST",
		URL:                    "/magma/v1/lte/n0/subscribers",
		Payload:                tests.JSONMarshaler(subscriberModels.MutableSubscribers{mutableSub}),
		Handler:                postSubscriber,
		ParamNames:             []string{"network_id"},
		ParamValues:            []string{"n0"},
		ExpectedStatus:         500, // would make more sense as 400
		ExpectedErrorSubstring: `could not find entities matching [type:"apn" key:"apnXXX" ]`,
	}
	tests.RunUnitTest(t, e, tc)

	// Post err, rule doesn't exist
	mutableSub.ActivePoliciesByApn = policydbModels.PolicyIdsByApn{"apn0": policydbModels.PolicyIds{"ruleXXX"}}
	tc = tests.Test{
		Method:                 "POST",
		URL:                    "/magma/v1/lte/n0/subscribers",
		Payload:                tests.JSONMarshaler(subscriberModels.MutableSubscribers{mutableSub}),
		Handler:                postSubscriber,
		ParamNames:             []string{"network_id"},
		ParamValues:            []string{"n0"},
		ExpectedStatus:         500, // would make more sense as 400
		ExpectedErrorSubstring: `could not find entities matching [type:"policy" key:"ruleXXX" ]`,
	}
	tests.RunUnitTest(t, e, tc)

	// Post, successful
	mutableSub.ActivePoliciesByApn = policydbModels.PolicyIdsByApn{
		"apn0": policydbModels.PolicyIds{"rule0"},
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte/n0/subscribers",
		Payload:        tests.JSONMarshaler(subscriberModels.MutableSubscribers{mutableSub}),
		Handler:        postSubscriber,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Configurator confirms policy profile exists
	profiles, err := configurator.ListEntityKeys(context.Background(), "n0", lte.APNPolicyProfileEntityType)
	assert.NoError(t, err)
	assert.Len(t, profiles, 1)

	// Get all, posted subscriber found
	expected := subscriberModels.PaginatedSubscribers{
		TotalCount:    1,
		NextPageToken: "",
		Subscribers: map[string]*subscriberModels.Subscriber{
			imsi: mutableSub.ToSubscriber(),
		},
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/subscribers",
		Handler:        getAllSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expected),
	}
	tests.RunUnitTest(t, e, tc)

	// Put err, APN doesn't exist
	mutableSub.ActivePoliciesByApn = policydbModels.PolicyIdsByApn{"apnXXX": policydbModels.PolicyIds{"rule0"}}
	tc = tests.Test{
		Method:                 "PUT",
		URL:                    "/magma/v1/lte/n0/subscribers/" + imsi,
		Payload:                mutableSub,
		ParamNames:             []string{"network_id", "subscriber_id"},
		ParamValues:            []string{"n0", imsi},
		Handler:                putSubscriber,
		ExpectedStatus:         500, // would make more sense as 400
		ExpectedErrorSubstring: `could not find entities matching [type:"apn" key:"apnXXX" ]`,
	}
	tests.RunUnitTest(t, e, tc)

	// Put err, rule doesn't exist
	mutableSub.ActivePoliciesByApn = policydbModels.PolicyIdsByApn{"apn0": policydbModels.PolicyIds{"ruleXXX"}}
	tc = tests.Test{
		Method:                 "PUT",
		URL:                    "/magma/v1/lte/n0/subscribers/" + imsi,
		Payload:                mutableSub,
		ParamNames:             []string{"network_id", "subscriber_id"},
		ParamValues:            []string{"n0", imsi},
		Handler:                putSubscriber,
		ExpectedStatus:         500, // would make more sense as 400
		ExpectedErrorSubstring: `could not find entities matching [type:"policy" key:"ruleXXX" ]`,
	}
	tests.RunUnitTest(t, e, tc)

	// Put, add new mappings
	mutableSub.ActivePoliciesByApn = policydbModels.PolicyIdsByApn{
		"apn0": policydbModels.PolicyIds{"rule0", "rule1"},
		"apn1": policydbModels.PolicyIds{"rule1", "rule2"},
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/lte/n0/subscribers/" + imsi,
		Payload:        mutableSub,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n0", imsi},
		Handler:        putSubscriber,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Configurator confirms policy profiles exist
	profiles, err = configurator.ListEntityKeys(context.Background(), "n0", lte.APNPolicyProfileEntityType)
	assert.NoError(t, err)
	assert.Len(t, profiles, 2)

	// Get, changes are reflected
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/subscribers/" + imsi,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n0", imsi},
		Handler:        getSubscriber,
		ExpectedStatus: 200,
		ExpectedResult: mutableSub.ToSubscriber(),
	}
	tests.RunUnitTest(t, e, tc)

	// Delete
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/lte/n0/subscribers/" + imsi,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n0", imsi},
		Handler:        deleteSubscriber,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Delete, subsequent delete still "succeeds"
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/lte/n0/subscribers/" + imsi,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n0", imsi},
		Handler:        deleteSubscriber,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Get, delete confirmed
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/subscribers/" + imsi,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n0", imsi},
		Handler:        getSubscriber,
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	// Configurator confirms deletion
	profiles, err = configurator.ListEntityKeys(context.Background(), "n0", lte.APNPolicyProfileEntityType)
	assert.NoError(t, err)
	assert.Len(t, profiles, 0)

	// Post, add subscriber back
	mutableSub.ActivePoliciesByApn = policydbModels.PolicyIdsByApn{
		"apn0": policydbModels.PolicyIds{"rule0", "rule1"},
		"apn1": policydbModels.PolicyIds{"rule1", "rule2"},
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte/n0/subscribers",
		Payload:        tests.JSONMarshaler(subscriberModels.MutableSubscribers{mutableSub}),
		Handler:        postSubscriber,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Get, successfully added back
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/subscribers/" + imsi,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n0", imsi},
		Handler:        getSubscriber,
		ExpectedStatus: 200,
		ExpectedResult: mutableSub.ToSubscriber(),
	}
	tests.RunUnitTest(t, e, tc)

	// Delete linked policy rule
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/networks/n1/policies/rules/rule0",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n0", "rule0"},
		Handler:        deletePolicy,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Get, policy rule changes reflected
	sub.ActivePoliciesByApn = policydbModels.PolicyIdsByApn{
		"apn0": policydbModels.PolicyIds{"rule1"}, // rule0 deleted
		"apn1": policydbModels.PolicyIds{"rule1", "rule2"},
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/subscribers/" + imsi,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n0", imsi},
		Handler:        getSubscriber,
		ExpectedStatus: 200,
		ExpectedResult: sub,
	}
	tests.RunUnitTest(t, e, tc)

	// Delete linked APN
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/lte/n0/apns/apn1",
		Handler:        deleteAPN,
		ParamNames:     []string{"network_id", "apn_name"},
		ParamValues:    []string{"n0", "apn0"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Get, APN change reflected
	sub.ActivePoliciesByApn = policydbModels.PolicyIdsByApn{
		// DELETED: "apn0": policydbModels.PolicyIds{"rule1"},
		"apn1": policydbModels.PolicyIds{"rule1", "rule2"},
	}
	sub.ActiveApns = subscriberModels.ApnList{
		// DELETED: "apn0",
		"apn1",
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n0/subscribers/" + imsi,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n0", imsi},
		Handler:        getSubscriber,
		ExpectedStatus: 200,
		ExpectedResult: sub,
	}
	tests.RunUnitTest(t, e, tc)

	// Configurator confirms deletion
	profiles, err = configurator.ListEntityKeys(context.Background(), "n0", lte.APNPolicyProfileEntityType)
	assert.NoError(t, err)
	assert.Len(t, profiles, 1)

	// Post, add sub1, no namespacing issues
	mutableSub1 := newMutableSubscriber(imsi1)
	mutableSub1.ActivePoliciesByApn = policydbModels.PolicyIdsByApn{"apn1": policydbModels.PolicyIds{"rule1", "rule2"}}
	mutableSub1.ActiveApns = subscriberModels.ApnList{"apn1"}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte/n0/subscribers",
		Payload:        tests.JSONMarshaler(subscriberModels.MutableSubscribers{mutableSub1}),
		Handler:        postSubscriber,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n0"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Configurator non-shared apn_policy_profile
	profiles, err = configurator.ListEntityKeys(context.Background(), "n0", lte.APNPolicyProfileEntityType)
	assert.NoError(t, err)
	assert.Len(t, profiles, 2)
}

func f32Ptr(f float32) *float32 {
	return &f
}

func newMutableSubscriber(id string) *subscriberModels.MutableSubscriber {
	sub := &subscriberModels.MutableSubscriber{
		ID:   policydbModels.SubscriberID(id),
		Name: "Jane Doe",
		Lte: &subscriberModels.LteSubscription{
			AuthAlgo:   "MILENAGE",
			AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			State:      "ACTIVE",
			SubProfile: "default",
		},
		ForbiddenNetworkTypes: subscriberModels.CoreNetworkTypes{"5GC", "EPC"},
		StaticIps: subscriberModels.SubscriberStaticIps{
			"apn1": "192.168.100.1",
		},
		ActiveApns: subscriberModels.ApnList{"apn0", "apn1"},
	}
	return sub
}

func newPolicy(id string) *policydbModels.PolicyRule {
	policy := &policydbModels.PolicyRule{
		ID: policydbModels.PolicyID(id),
		FlowList: []*policydbModels.FlowDescription{
			{
				Action: swag.String("PERMIT"),
				Match: &policydbModels.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_IP"),
				},
			},
		},
		Priority: swag.Uint32(1),
	}
	return policy
}
