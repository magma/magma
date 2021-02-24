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

package statemachines

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"magma/fbinternal/cloud/go/services/testcontroller/obsidian/models"
	"magma/fbinternal/cloud/go/services/testcontroller/storage"
	"magma/fbinternal/cloud/go/services/testcontroller/utils"
	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	ltemodels "magma/lte/cloud/go/services/lte/obsidian/models"
	subscribermodels "magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/magmad"
	models2 "magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/wrappers"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

const (
	checkForUpgradeState = "check_for_upgrade"

	verifyUpgradeStateFmt      = "verify_upgrade_%d"
	maxVerifyUpgradeStateCount = 3
	verifyUpgrade1State        = "verify_upgrade_1"
	verifyUpgrade2State        = "verify_upgrade_2"
	verifyUpgrade3State        = "verify_upgrade_3"

	maxTrafficStateCount = 3
	trafficTestStateFmt  = "traffic_test%d_%d"
	trafficTest1State1   = "traffic_test1_1"
	trafficTest1State2   = "traffic_test1_2"
	trafficTest1State3   = "traffic_test1_3"

	rebootEnodebStateFmt      = "reboot_enodeb_%d"
	maxRebootEnodebStateCount = 3
	rebootEnodeb1State        = "reboot_enodeb_1"
	rebootEnodeb2State        = "reboot_enodeb_2"
	rebootEnodeb3State        = "reboot_enodeb_3"

	verifyConnectivityState = "verify_conn"

	trafficTest2State1 = "traffic_test2_1"
	trafficTest2State2 = "traffic_test2_2"
	trafficTest2State3 = "traffic_test2_3"

	maxConfigStateCount    = 3
	reconfigEnodebStateFmt = "reconfig_enodeb%d"
	reconfigEnodebState1   = "reconfig_enodeb1"
	reconfigEnodebState2   = "reconfig_enodeb2"
	reconfigEnodebState3   = "reconfig_enodeb3"
	verifyConfig1State     = "verify_config1"

	trafficTest3State1 = "traffic_test3_1"
	trafficTest3State2 = "traffic_test3_2"
	trafficTest3State3 = "traffic_test3_3"

	restoreEnodebConfigStateFmt = "restore_enodeb%d"
	restoreEnodebConfigState1   = "restore_enodeb1"
	restoreEnodebConfigState2   = "restore_enodeb2"
	restoreEnodebConfigState3   = "restore_enodeb3"
	verifyConfig2State          = "verify_config2"

	trafficTest4State1 = "traffic_test4_1"
	trafficTest4State2 = "traffic_test4_2"
	trafficTest4State3 = "traffic_test4_3"

	subscriberInactiveState = "subscriber_inactive"

	trafficTest5State1 = "traffic_test5_1"
	trafficTest5State2 = "traffic_test5_2"
	trafficTest5State3 = "traffic_test5_3"

	subscriberActiveState = "subscriber_active"

	trafficTest6State1 = "traffic_test6_1"
	trafficTest6State2 = "traffic_test6_2"
	trafficTest6State3 = "traffic_test6_3"
)

// GatewayClient defines an interface which is used to switch between
// implementations of its methods between a real implementation and a mock up for unit testing
type GatewayClient interface {
	GenerateTraffic(networdId string, gatewayId string, ssid string, pw string) (*protos.GenericCommandResponse, error)
	RebootEnodeb(networkdId string, gatewayId string, enodebSerial string) (*protos.GenericCommandResponse, error)
}

type MagmadClient struct{}

func (m *MagmadClient) GenerateTraffic(networkId string, gatewayId string, ssid string, pw string) (*protos.GenericCommandResponse, error) {
	stringVal := fmt.Sprintf("-c 'python3 /usr/local/bin/traffic_cli.py gen_traffic %s %s http://www.google.com'", ssid, pw)
	params := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"shell_params": {Kind: &structpb.Value_ListValue{
				ListValue: &structpb.ListValue{
					Values: []*structpb.Value{{Kind: &structpb.Value_StringValue{StringValue: stringVal}}},
				},
			}},
		},
	}
	trafficScriptCmd := &protos.GenericCommandParams{
		Command: "bash",
		Params:  params,
	}
	resp, err := magmad.GatewayGenericCommand(networkId, gatewayId, trafficScriptCmd)
	return resp, err
}

func (m *MagmadClient) RebootEnodeb(networkId string, gatewayId string, enodebSerial string) (*protos.GenericCommandResponse, error) {
	params := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"shell_params": {Kind: &structpb.Value_StringValue{StringValue: enodebSerial}},
		},
	}
	rebootEndCmd := &protos.GenericCommandParams{
		Command: "reboot_enodeb",
		Params:  params,
	}
	resp, err := magmad.GatewayGenericCommand(networkId, gatewayId, rebootEndCmd)
	return resp, err
}

func getEnodebStatus(networkID string, enodebSN string) (*ltemodels.EnodebState, error) {
	st, err := state.GetState(networkID, lte.EnodebStateType, enodebSN, serdes.State)
	if err != nil {
		return nil, err
	}
	enodebState := st.ReportedState.(*ltemodels.EnodebState)
	enodebState.TimeReported = st.TimeMs
	ent, err := configurator.LoadEntityForPhysicalID(st.ReporterID, configurator.EntityLoadCriteria{}, serdes.Entity)
	if err == nil {
		enodebState.ReportingGatewayID = ent.Key
	}
	return enodebState, err
}

var (
	magmaPackageVersionRegex = regexp.MustCompile(`(?:Version:\s)(.*)`)
)

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

func NewEnodebdE2ETestStateMachine(store storage.TestControllerStorage, client HttpClient, gatewayClient GatewayClient) TestMachine {
	return &enodebdE2ETestStateMachine{
		store: store,
		stateHandlers: map[string]handlerFunc{
			storage.CommonStartState: startNewTest,
			checkForUpgradeState:     checkForUpgrade,

			verifyUpgrade1State: makeVerifyUpgradeStateHandler(1, trafficTest1State1),
			verifyUpgrade2State: makeVerifyUpgradeStateHandler(2, trafficTest1State1),
			verifyUpgrade3State: makeVerifyUpgradeStateHandler(3, trafficTest1State1),

			trafficTest1State1: makeTrafficTestStateHandler(1, 1, gatewayClient, rebootEnodeb1State),
			trafficTest1State2: makeTrafficTestStateHandler(1, 2, gatewayClient, rebootEnodeb1State),
			trafficTest1State3: makeTrafficTestStateHandler(1, 3, gatewayClient, rebootEnodeb1State),

			rebootEnodeb1State:      makeRebootEnodebStateHandler(1, gatewayClient),
			rebootEnodeb2State:      makeRebootEnodebStateHandler(2, gatewayClient),
			rebootEnodeb3State:      makeRebootEnodebStateHandler(3, gatewayClient),
			verifyConnectivityState: makeVerifyConnectivityHandler(trafficTest2State1),

			trafficTest2State1: makeTrafficTestStateHandler(2, 1, gatewayClient, reconfigEnodebState1),
			trafficTest2State2: makeTrafficTestStateHandler(2, 2, gatewayClient, reconfigEnodebState1),
			trafficTest2State3: makeTrafficTestStateHandler(2, 3, gatewayClient, reconfigEnodebState1),

			reconfigEnodebState1: makeConfigEnodebStateHandler(1, verifyConfig1State),
			reconfigEnodebState2: makeConfigEnodebStateHandler(2, verifyConfig1State),
			reconfigEnodebState3: makeConfigEnodebStateHandler(3, verifyConfig1State),
			verifyConfig1State:   makeVerifyConfigStateHandler(trafficTest3State1),

			trafficTest3State1: makeTrafficTestStateHandler(3, 1, gatewayClient, restoreEnodebConfigState1),
			trafficTest3State2: makeTrafficTestStateHandler(3, 2, gatewayClient, restoreEnodebConfigState1),
			trafficTest3State3: makeTrafficTestStateHandler(3, 3, gatewayClient, restoreEnodebConfigState1),

			restoreEnodebConfigState1: makeConfigEnodebStateHandler(1, verifyConfig2State),
			restoreEnodebConfigState2: makeConfigEnodebStateHandler(2, verifyConfig2State),
			restoreEnodebConfigState3: makeConfigEnodebStateHandler(3, verifyConfig2State),
			verifyConfig2State:        makeVerifyConfigStateHandler(trafficTest4State1),

			trafficTest4State1: makeTrafficTestStateHandler(4, 1, gatewayClient, subscriberInactiveState),
			trafficTest4State2: makeTrafficTestStateHandler(4, 2, gatewayClient, subscriberInactiveState),
			trafficTest4State3: makeTrafficTestStateHandler(4, 3, gatewayClient, subscriberInactiveState),

			subscriberInactiveState: makeSubscriberStateHandler(subscribermodels.LteSubscriptionStateINACTIVE, trafficTest5State1),

			trafficTest5State1: makeTrafficTestStateHandler(5, 1, gatewayClient, subscriberActiveState),
			trafficTest5State2: makeTrafficTestStateHandler(5, 2, gatewayClient, subscriberActiveState),
			trafficTest5State3: makeTrafficTestStateHandler(5, 3, gatewayClient, subscriberActiveState),

			subscriberActiveState: makeSubscriberStateHandler(subscribermodels.LteSubscriptionStateACTIVE, trafficTest6State1),

			trafficTest6State1: makeTrafficTestStateHandler(6, 1, gatewayClient, checkForUpgradeState),
			trafficTest6State2: makeTrafficTestStateHandler(6, 2, gatewayClient, checkForUpgradeState),
			trafficTest6State3: makeTrafficTestStateHandler(6, 3, gatewayClient, checkForUpgradeState),
		},
		client: client,
	}
}

func NewEnodebdE2ETestStateMachineNoTraffic(store storage.TestControllerStorage, client HttpClient, gatewayClient GatewayClient) TestMachine {
	return &enodebdE2ETestStateMachine{
		store: store,
		stateHandlers: map[string]handlerFunc{
			storage.CommonStartState: startNewTest,
			checkForUpgradeState:     checkForUpgrade,

			verifyUpgrade1State: makeVerifyUpgradeStateHandler(1, rebootEnodeb1State),
			verifyUpgrade2State: makeVerifyUpgradeStateHandler(2, rebootEnodeb1State),
			verifyUpgrade3State: makeVerifyUpgradeStateHandler(3, rebootEnodeb1State),

			rebootEnodeb1State:      makeRebootEnodebStateHandler(1, gatewayClient),
			rebootEnodeb2State:      makeRebootEnodebStateHandler(2, gatewayClient),
			rebootEnodeb3State:      makeRebootEnodebStateHandler(3, gatewayClient),
			verifyConnectivityState: makeVerifyConnectivityHandler(reconfigEnodebState1),

			reconfigEnodebState1: makeConfigEnodebStateHandler(1, verifyConfig1State),
			reconfigEnodebState2: makeConfigEnodebStateHandler(2, verifyConfig1State),
			reconfigEnodebState3: makeConfigEnodebStateHandler(3, verifyConfig1State),
			verifyConfig1State:   makeVerifyConfigStateHandler(restoreEnodebConfigState1),

			restoreEnodebConfigState1: makeConfigEnodebStateHandler(1, verifyConfig2State),
			restoreEnodebConfigState2: makeConfigEnodebStateHandler(2, verifyConfig2State),
			restoreEnodebConfigState3: makeConfigEnodebStateHandler(3, verifyConfig2State),
			verifyConfig2State:        makeVerifyConfigStateHandler(subscriberInactiveState),

			subscriberInactiveState: makeSubscriberStateHandler(subscribermodels.LteSubscriptionStateINACTIVE, subscriberActiveState),
			subscriberActiveState:   makeSubscriberStateHandler(subscribermodels.LteSubscriptionStateACTIVE, checkForUpgradeState),
		},
		client: client,
	}
}

type handlerFunc func(*enodebdE2ETestStateMachine, *models.EnodebdTestConfig) (string, time.Duration, error)

type enodebdE2ETestStateMachine struct {
	store         storage.TestControllerStorage
	stateHandlers map[string]handlerFunc
	client        HttpClient
}

func (e *enodebdE2ETestStateMachine) Run(state string, config interface{}, previousErr error) (string, time.Duration, error) {
	// TODO: notify slack if previousErr is non-nil?
	configCasted, ok := config.(*models.EnodebdTestConfig)
	if !ok {
		return "", 1 * time.Hour, errors.Errorf("expected config *models.EnodebdTestConfig, got %T", config)
	}
	handler, found := e.stateHandlers[state]
	if !found {
		return "", 1 * time.Hour, errors.Errorf("no handler registered for test case state %s", state)
	}
	return handler(e, configCasted)
}

// Handlers (consider making these instance methods)
// TODO: ASCII art diagram of the state machine
// TODO: refactor out the AGW autoupgrade handlers if/when we make more test cases

/*
States with traffic tests enabled:
	- Check for upgrade: compare repo version and gateway reported version; change tier config if different
		> Epsilon if same, 20 minutes
		> "Verify upgrade 1" if different, 10 minutes
	- Upgrade gateway: change version of target tier
		> "Verify upgrade 1", 10 minutes
	- Verify upgrade N: check gateway's reported version against tier, ping slack if max attempts reached and unsuccessful; there are 3 of these states
		> "Traffic test 1 1" if equal, 20 minutes
		> "Verify upgrade N+1" if not equal and N < 3, 20 minutes
		> "Check for upgrade" if not equal and N >= 3 (after pinging slack), 20 minutes
	- Traffic Test 1 N: Send user traffic to an arbitrary end point, ping slack upon success/failure; there are 3 of these states
		> "Reboot Enodeb 1" traffic is sent, pings slack, 1 minute
		> "Traffic Test 1 N+1" traffic unable to be sent, 1 minute
		> "Check for upgrade" N >= 3 (ping slack), 1 minute
	- Reboot enodeb N: reboot a gateway's enodeb, ping slack if max attempts reached and unsuccessful; there are 3 of these states
		> "Reboot Enodeb N+1" if reboot fails and N < 3, 15 minutes
		> "Check for upgrade" if N >= 3, 15 minutes (ping slack)
		> "Verify Connectivity", 15 minutes
	- Verify Connectivity: after reboot, we check that the enodeb has successfully reconnected, then ping slack whether successful or unsuccessful
		> "Traffic Test 2 1" able to get enodeb status, pings slack, 15 minutes
		> "Check for upgrade" if cannot get enodeb status or hwID, 5 minutes
		> "Check for upgrade", 15 minutes
	- Traffic Test 2 N: Same deal with traffic test 1, only difference is success state transition
		> "Reconfig Enodeb 1" traffic is sent, pings slack, 1 minute
	- Reconfig Enodeb N: Changes config (for now PCI) of enodeb. Pings slack upon successful/fail
		> "Verify Config 1" pings slack, 10 minutes
		> "Reconfig Enodeb N+1" Error in reconfiguring enodeb, 5 minutes
		> "Check for upgrade" if N >= 3, 10 minutes (ping slack)
	- Verify Config 1: Makes sure enodeb has a configuration
		> "Traffic Test 3 1" Enodeb has a config, 10 minutes
		> "Check for upgrade" if cannot get enodeb status, 5 minutes, else 10 minutes
	- Traffic Test 3 N: Same deal with traffic test 1, only difference is success state transition
		> "Restore Enodeb 1" traffic is sent, pings slack, 1 minute
	- Restore Enodeb N: Return enodeb config to original state
		> "Verify Config 2" pings slack, 10 minutes
		> "Restore Enodeb N+1" Error in reconfiguring enodeb, 5 minutes
		> "Check for upgrade" if N >= 3, 10 minutes (ping slack)
	- Verify Config 2 Same deal with Verify Config 1, only difference is success state transition
		> "Traffic Test 4 1" Enodeb has a config, 10 minutes
	- Traffic Test 4 N: Same deal with traffic test 1, only difference is success state transition
		> "Subscriber Inactive" traffic is sent, pings slack, 1 minute
	- Subscriber Inactive: Flips subscriber state to INACTIVE
		> "Traffic Test 5 1" Subscriber state successfully set to INACTIVE, 10 minutes
		> "Check for upgrade" Unsuccessfully set subscriber state, 10 minutes
	- Traffic Test 5 N: Expected to fail since subscriber cannot receive data. The "successful" state transition will be triggered by a failure
		> "Restore Enodeb 1" traffic is NOT sent, pings slack, 1 minute (success)
		> "Check for upgrade" traffic gets sent, pings slack, 1 minute (fail)
	- Subscriber Active: Restores subscriber state to ACTIVE
		> "Traffic Test 6 1" Subscriber state successfully set to ACTIVE, 10 minutes
		> "Check for upgrade" Unsuccessfully set subscriber state, 10 minutes
	- Traffic Test 6 N: Same deal with traffic test 1, only difference is success state transition
		> "Check for upgrade" pings slack, 1 minute
*/

func startNewTest(machine *enodebdE2ETestStateMachine, config *models.EnodebdTestConfig) (string, time.Duration, error) {
	states := []string{
		checkForUpgradeState,
		rebootEnodeb1State,
		reconfigEnodebState1,
		restoreEnodebConfigState1,
		subscriberInactiveState,
		subscriberActiveState,
	}
	if config.StartState != "" {
		if !find(states, config.StartState) {
			// Invalid state, default to check_for_upgrade
			return checkForUpgradeState, time.Minute, errors.Errorf("Invalid starting state. Defaulting to check_for_upgrade")
		}
	} else {
		return checkForUpgradeState, time.Minute, nil
	}
	return config.StartState, time.Minute, nil
}

func makeSubscriberStateHandler(desiredState string, successState string) handlerFunc {
	return func(machine *enodebdE2ETestStateMachine, config *models.EnodebdTestConfig) (string, time.Duration, error) {
		return subscriberState(desiredState, successState, machine, config)
	}
}

func subscriberState(desiredState string, successState string, machine *enodebdE2ETestStateMachine, config *models.EnodebdTestConfig) (string, time.Duration, error) {
	pretext := fmt.Sprintf(subscriberPretextFmt, *config.SubscriberID, desiredState, "SUCCEEDED")
	fallback := "Subscriber state notification"
	cfg, err := configurator.LoadEntityConfig(*config.NetworkID, lte.SubscriberEntityType, *config.SubscriberID, serdes.Entity)
	if err != nil {
		pretext = fmt.Sprintf(subscriberPretextFmt, *config.SubscriberID, desiredState, "FAILED")
		postToSlack(machine.client, *config.AgwConfig.SlackWebhook, false, pretext, fallback, "", "")
		return checkForUpgradeState, 10 * time.Minute, err
	}

	newConfig, ok := cfg.(*subscribermodels.SubscriberConfig)
	if !ok {
		glog.Errorf("got data of type %T but wanted SubscriberConfig", cfg)
		return checkForUpgradeState, 10 * time.Minute, err
	}
	newConfig.Lte.State = desiredState
	err = configurator.CreateOrUpdateEntityConfig(*config.NetworkID, lte.SubscriberEntityType, *config.SubscriberID, newConfig, serdes.Entity)
	if err != nil {
		// Restore subscriber to original config before erroring out
		err = configurator.CreateOrUpdateEntityConfig(*config.NetworkID, lte.SubscriberEntityType, *config.SubscriberID, cfg, serdes.Entity)
		if err != nil {
			glog.Error(err)
		}
		pretext = fmt.Sprintf(subscriberPretextFmt, *config.SubscriberID, desiredState, "FAILED")
		postToSlack(machine.client, *config.AgwConfig.SlackWebhook, false, pretext, fallback, "", "")
		return checkForUpgradeState, 10 * time.Minute, err
	}
	postToSlack(machine.client, *config.AgwConfig.SlackWebhook, true, pretext, fallback, "", "")
	return successState, 10 * time.Minute, nil
}

func makeConfigEnodebStateHandler(stateNumber int, successState string) handlerFunc {
	return func(machine *enodebdE2ETestStateMachine, config *models.EnodebdTestConfig) (string, time.Duration, error) {
		return configEnodeb(stateNumber, successState, machine, config)
	}
}

func configEnodeb(stateNumber int, successState string, machine *enodebdE2ETestStateMachine, config *models.EnodebdTestConfig) (string, time.Duration, error) {
	pretext := fmt.Sprintf(reconfigPretextFmt, *config.EnodebSN, "SUCCEEDED")
	fallback := "Reconfig enodeb notification"
	_, err := configurator.UpdateEntity(
		*config.NetworkID,
		configurator.EntityUpdateCriteria{
			Type:      lte.CellularEnodebEntityType,
			Key:       *config.EnodebSN,
			NewConfig: config.EnodebConfig,
		},
		serdes.Entity,
	)

	if err != nil {
		if stateNumber >= maxConfigStateCount {
			// TODO Restore enodeb config to original state
			pretext = fmt.Sprintf(reconfigPretextFmt, *config.EnodebSN, "FAILED")
			if successState == trafficTest4State1 {
				pretext = fmt.Sprintf(restoreConfigPretextFmt, *config.EnodebSN, "FAILED")
			}
			postToSlack(machine.client, *config.AgwConfig.SlackWebhook, false, pretext, fallback, "", "")
			return checkForUpgradeState, 10 * time.Minute, err
		}
		switch successState {
		case trafficTest3State1:
			return fmt.Sprintf(reconfigEnodebStateFmt, stateNumber+1), 5 * time.Minute, err
		case trafficTest4State1:
			return fmt.Sprintf(restoreEnodebConfigStateFmt, stateNumber+1), 5 * time.Minute, err
		default:
			return checkForUpgradeState, 10 * time.Minute, err
		}
	}
	postToSlack(machine.client, *config.AgwConfig.SlackWebhook, true, pretext, fallback, "", "")
	return successState, 10 * time.Minute, nil
}

func makeVerifyConfigStateHandler(successState string) handlerFunc {
	return func(machine *enodebdE2ETestStateMachine, config *models.EnodebdTestConfig) (string, time.Duration, error) {
		return verifyConfig(successState, machine, config)
	}
}

func verifyConfig(successState string, machine *enodebdE2ETestStateMachine, config *models.EnodebdTestConfig) (string, time.Duration, error) {
	resp, err := getEnodebStatus(*config.NetworkID, *config.EnodebSN)
	if resp == nil || err != nil {
		return checkForUpgradeState, 5 * time.Minute, errors.Wrap(err, "error getting enodeb status")
	}

	if !*resp.EnodebConfigured {
		return checkForUpgradeState, 10 * time.Minute, errors.Errorf("error enodeb %s is not configured", *config.EnodebSN)
	}
	return successState, 10 * time.Minute, nil
}

func makeTrafficTestStateHandler(trafficTestNumber int, stateNumber int, gatewayClient GatewayClient, successState string) handlerFunc {
	return func(machine *enodebdE2ETestStateMachine, config *models.EnodebdTestConfig) (string, time.Duration, error) {
		return trafficTest(trafficTestNumber, stateNumber, gatewayClient, successState, machine, config)
	}
}

func trafficTest(trafficTestNumber int, stateNumber int, gatewayClient GatewayClient, successState string, machine *enodebdE2ETestStateMachine, config *models.EnodebdTestConfig) (string, time.Duration, error) {
	trafficGWID := *config.TrafficGwID
	pretext := fmt.Sprintf(trafficPretextFmt, trafficTestNumber, *config.EnodebSN, *config.AgwConfig.TargetGatewayID, "SUCCEEDED")
	fallback := "Generate traffic notification"

	helper := &protos.GenericCommandResponse{
		Response: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"returncode": {Kind: &structpb.Value_NumberValue{NumberValue: float64(0)}},
				"stdout":     {Kind: &structpb.Value_StringValue{StringValue: ""}},
				"stderr":     {Kind: &structpb.Value_StringValue{StringValue: ""}},
			},
		},
	}
	resp, err := gatewayClient.GenerateTraffic(*config.NetworkID, trafficGWID, config.Ssid, config.SsidPw)
	// Any result that is not 0 is considered a failure on the traffic script's part
	if resp == nil || err != nil || !proto.Equal(resp.Response.Fields["returncode"], helper.Response.Fields["returncode"]) {
		if successState == subscriberActiveState {
			pretext = fmt.Sprintf(trafficInactiveSubPretextFmt, *config.SubscriberID, *config.EnodebSN, *config.AgwConfig.TargetGatewayID, "SUCCEEDED")
			postToSlack(machine.client, *config.AgwConfig.SlackWebhook, false, pretext, fallback, "", "")
			return successState, 1 * time.Minute, nil
		}
		if stateNumber >= maxTrafficStateCount {
			pretext = fmt.Sprintf(trafficPretextFmt, trafficTestNumber, *config.EnodebSN, *config.AgwConfig.TargetGatewayID, "FAILED")
			postToSlack(machine.client, *config.AgwConfig.SlackWebhook, false, pretext, fallback, "", "")
			return checkForUpgradeState, 1 * time.Minute, errors.Errorf("Traffic test number %d failed on gwID %s after %d tries", trafficTestNumber, trafficGWID, maxTrafficStateCount)
		}
		if err == nil {
			err = errors.Errorf("Traffic script failed. Return Code: %d, Stdout: %s, Stderr: %s",
				int(resp.Response.Fields["returncode"].GetNumberValue()),
				resp.Response.Fields["stdout"].GetStringValue(),
				resp.Response.Fields["stderr"].GetStringValue(),
			)
		}
		return fmt.Sprintf(trafficTestStateFmt, trafficTestNumber, stateNumber+1), 1 * time.Minute, err
	}
	if successState == subscriberActiveState {
		pretext = fmt.Sprintf(trafficInactiveSubPretextFmt, *config.SubscriberID, *config.EnodebSN, *config.AgwConfig.TargetGatewayID, "FAILED")
		postToSlack(machine.client, *config.AgwConfig.SlackWebhook, true, pretext, fallback, "", "")
		return checkForUpgradeState, 1 * time.Minute, errors.Errorf("Traffic test number %d should not have succeeded on gwID %s", trafficTestNumber, trafficGWID)
	}
	postToSlack(machine.client, *config.AgwConfig.SlackWebhook, true, pretext, fallback, "", "")
	return successState, 1 * time.Minute, nil
}

func makeRebootEnodebStateHandler(stateNumber int, gatewayClient GatewayClient) handlerFunc {
	return func(machine *enodebdE2ETestStateMachine, config *models.EnodebdTestConfig) (string, time.Duration, error) {
		return rebootEnodebStateHandler(stateNumber, gatewayClient, machine, config)
	}
}

func rebootEnodebStateHandler(stateNumber int, gatewayClient GatewayClient, machine *enodebdE2ETestStateMachine, config *models.EnodebdTestConfig) (string, time.Duration, error) {
	targetGWID := *config.AgwConfig.TargetGatewayID
	enodebSN := *config.EnodebSN

	resp, err := gatewayClient.RebootEnodeb(*config.NetworkID, targetGWID, enodebSN)
	if resp == nil || err != nil {
		if stateNumber >= maxRebootEnodebStateCount {
			pretext := fmt.Sprintf(rebootPretextFmt, enodebSN, targetGWID, "FAILED")
			fallback := "Reboot enodeb notification"
			postToSlack(machine.client, *config.AgwConfig.SlackWebhook, false, pretext, fallback, "", "")
			return checkForUpgradeState, 15 * time.Minute, errors.Errorf("enodeb %s did not reboot within %d tries", enodebSN, maxRebootEnodebStateCount)
		}
		return fmt.Sprintf(rebootEnodebStateFmt, stateNumber+1), 5 * time.Minute, err
	}
	return verifyConnectivityState, 15 * time.Minute, nil
}

func makeVerifyConnectivityHandler(successState string) handlerFunc {
	return func(machine *enodebdE2ETestStateMachine, config *models.EnodebdTestConfig) (string, time.Duration, error) {
		return verifyConnectivity(successState, machine, config)
	}
}

func verifyConnectivity(successState string, machine *enodebdE2ETestStateMachine, config *models.EnodebdTestConfig) (string, time.Duration, error) {
	targetGWID := *config.AgwConfig.TargetGatewayID
	enodebSN := *config.EnodebSN
	pretext := fmt.Sprintf(rebootPretextFmt, enodebSN, targetGWID, "FAILED")
	fallback := "Reboot enodeb notification"

	resp, err := getEnodebStatus(*config.NetworkID, enodebSN)
	if resp == nil || err != nil {
		postToSlack(machine.client, *config.AgwConfig.SlackWebhook, false, pretext, fallback, "", "")
		return checkForUpgradeState, 5 * time.Minute, errors.Wrap(err, "error getting enodeb status")
	}
	if !*resp.EnodebConnected {
		postToSlack(machine.client, *config.AgwConfig.SlackWebhook, false, pretext, fallback, "", "")
		return checkForUpgradeState, 5 * time.Minute, errors.Errorf("Error Enodeb is not connected")
	}
	if !*resp.RfTxDesired || !*resp.RfTxOn {
		postToSlack(machine.client, *config.AgwConfig.SlackWebhook, false, pretext, fallback, "", "")
		return checkForUpgradeState, 5 * time.Minute, errors.Errorf("Error RF TX on/desired are not both set to true")
	}

	pretext = fmt.Sprintf(rebootPretextFmt, enodebSN, targetGWID, "SUCCEEDED")
	postToSlack(machine.client, *config.AgwConfig.SlackWebhook, true, pretext, fallback, "", "")
	return successState, 15 * time.Minute, nil
}

func checkForUpgrade(machine *enodebdE2ETestStateMachine, config *models.EnodebdTestConfig) (string, time.Duration, error) {
	repoVersion, err := getLatestRepoMagmaVersion(machine.client, *config.AgwConfig.PackageRepo, *config.AgwConfig.ReleaseChannel)
	if err != nil {
		return checkForUpgradeState, 20 * time.Minute, errors.Wrap(err, "error getting latest package version from repo")
	}

	tierCfg, err := getTargetTierConfig(config)
	if err != nil {
		return checkForUpgradeState, 20 * time.Minute, err
	}
	existingVersion := tierCfg.Version.ToString()
	if existingVersion == "" {
		existingVersion = "0.0.0"
	}

	newer, err := utils.IsNewerVersion(existingVersion, repoVersion)
	if err != nil {
		return checkForUpgradeState, 20 * time.Minute, errors.Wrapf(err, "bad versions encountered: %s, %s", repoVersion, existingVersion)
	}
	if !newer {
		return checkForUpgradeState, 20 * time.Minute, nil
	}

	// Update the tier config
	newTierCfg := tierCfg
	newTierCfg.Version = models2.TierVersion(repoVersion)
	_, err = configurator.UpdateEntity(
		*config.NetworkID,
		configurator.EntityUpdateCriteria{
			Key:       *config.AgwConfig.TargetTier,
			Type:      orc8r.UpgradeTierEntityType,
			NewConfig: newTierCfg,
		},
		serdes.Entity,
	)
	if err != nil {
		return checkForUpgradeState, 20 * time.Minute, errors.Wrap(err, "error updating target tier")
	}
	return verifyUpgrade1State, 10 * time.Minute, nil
}

func makeVerifyUpgradeStateHandler(stateNumber int, successState string) handlerFunc {
	return func(machine *enodebdE2ETestStateMachine, config *models.EnodebdTestConfig) (string, time.Duration, error) {
		return verifyUpgrade(stateNumber, successState, machine, config)
	}
}

func verifyUpgrade(stateNumber int, successState string, machine *enodebdE2ETestStateMachine, config *models.EnodebdTestConfig) (string, time.Duration, error) {
	targetGWID := *config.AgwConfig.TargetGatewayID
	fallback := "Gateway auto-upgrade notification"

	// Load target version
	tierCfg, err := getTargetTierConfig(config)
	if err != nil {
		return fmt.Sprintf(verifyUpgradeStateFmt, stateNumber+1), 10 * time.Minute, err
	}
	currentVersion, err := getCurrentAGWPackageVersion(config)
	if err != nil {
		return fmt.Sprintf(verifyUpgradeStateFmt, stateNumber+1), 10 * time.Minute, err
	}

	// If equal, transition to reboot enodeb state
	if string(tierCfg.Version) == currentVersion {
		pretext := fmt.Sprintf(autoupgradePretextFmt, targetGWID, "SUCCEEDED", "")
		postToSlack(machine.client, *config.AgwConfig.SlackWebhook, true, pretext, fallback, string(tierCfg.Version), "")
		return successState, 20 * time.Minute, nil
	}

	if stateNumber >= maxVerifyUpgradeStateCount {
		pretext := fmt.Sprintf(autoupgradePretextFmt, targetGWID, "FAILED", "")
		postToSlack(machine.client, *config.AgwConfig.SlackWebhook, false, pretext, fallback, string(tierCfg.Version), "")
		return checkForUpgradeState, 20 * time.Minute, errors.Errorf("gateway %s did not upgrade within %d tries", targetGWID, maxVerifyUpgradeStateCount)
	} else {
		return fmt.Sprintf(verifyUpgradeStateFmt, stateNumber+1), 10 * time.Minute, nil
	}
}

// State handler helpers
func find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func getLatestRepoMagmaVersion(client HttpClient, url string, releaseChannel string) (string, error) {
	url = fmt.Sprintf("%s/dists/%s/main/binary-amd64/Packages", url, releaseChannel)
	resp, err := client.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "unable to retrieve packages from url %s", url)
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "unable to read body from http get response at url %s", url)
	}

	if len(responseData) == 0 {
		return "", errors.Errorf("no packages found at url %s", url)
	}

	packages := strings.Split(string(responseData), "\n\n")

	latest := "0.0.0-0-0"
	for _, pkg := range packages {
		if !strings.Contains(pkg, "Package: magma\n") {
			continue
		}

		version := magmaPackageVersionRegex.FindStringSubmatch(pkg)
		if len(version) != 2 {
			glog.Warningf("incorrect regex match on package version. "+
				"should have one capturing and one non-capturing group: %s", version)
			continue
		}
		newer, err := utils.IsNewerVersion(latest, version[1])
		if err != nil {
			return "", errors.Wrap(err, "unable to compare package versions")
		}
		if newer {
			latest = version[1]
		}
	}
	if latest == "0.0.0-0-0" {
		return "", errors.Errorf("no latest magma version found")
	}
	return latest, nil
}

func getTargetTierConfig(config *models.EnodebdTestConfig) (*models2.Tier, error) {
	tierEnt, err := configurator.LoadEntity(
		*config.NetworkID, orc8r.UpgradeTierEntityType, *config.AgwConfig.TargetTier,
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load target upgrade tier")
	}

	tierCfg, ok := tierEnt.Config.(*models2.Tier)
	if !ok {
		return nil, errors.Wrapf(err, "expected tier of type *models.Tier, got %T", tierEnt.Config)
	}
	return tierCfg, nil
}

func getCurrentAGWPackageVersion(config *models.EnodebdTestConfig) (string, error) {
	targetGWID := *config.AgwConfig.TargetGatewayID

	hwID, err := configurator.GetPhysicalIDOfEntity(*config.NetworkID, orc8r.MagmadGatewayType, *config.AgwConfig.TargetGatewayID)
	if err != nil {
		return "", errors.Wrapf(err, "failed to load hwID for target gateway %s", targetGWID)
	}
	agwState, err := wrappers.GetGatewayStatus(*config.NetworkID, hwID)
	if err != nil {
		return "", errors.Wrapf(err, "failed to load gateway status for %s", targetGWID)
	}
	if agwState == nil || agwState.PlatformInfo == nil {
		return "", errors.Wrapf(err, "gateway status not fully reported for %s", targetGWID)
	}
	magmaPackage := funk.Find(agwState.PlatformInfo.Packages, func(p *models2.Package) bool { return p.Name == "magma" })
	if magmaPackage == nil {
		return "", errors.Errorf("no magma package version reported for %s", targetGWID)
	}
	return magmaPackage.(*models2.Package).Version, nil
}

// swallow errors
func postToSlack(client HttpClient, slackURL string, success bool, pretext string, fallback string, targetVersion string, extraErrorText string) {
	payload, err := getSlackPayload(success, pretext, fallback, targetVersion, extraErrorText)
	if err != nil {
		glog.Errorf("failed to construct slack payload: %s", err)
		return
	}
	postPayload(client, slackURL, payload)
}

func postPayload(client HttpClient, slackURL string, payload io.Reader) {
	resp, err := client.Post(slackURL, "application/json", payload)
	if err != nil {
		glog.Errorf("slack webhook post failure: %s", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			glog.Errorf("failed to read non-200 response body: %s", err)
		}
		glog.Errorf("non-200 response %d from slack: %s", resp.StatusCode, string(respBody))
	}
}

const (
	colorRed                     = "#8b0902"
	colorGreen                   = "#36a64f"
	autoupgradePretextFmt        = "Auto-upgrade of gateway %s %s. %s"
	rebootPretextFmt             = "Enodeb reboot of enodeb %s of gateway %s %s"
	trafficPretextFmt            = "Generate traffic test %d for enodeb %s of gateway %s %s"
	trafficInactiveSubPretextFmt = "Generate traffic test for inactive subscriber %s for enodeb %s of gateway %s %s"
	reconfigPretextFmt           = "Reconfig enodeb %s %s"
	restoreConfigPretextFmt      = "Restored config of enodeb %s %s"
	subscriberPretextFmt         = "Subscriber %s set to %s %s"
)

type slackPayload struct {
	Attachments []slackAttachment `json:"attachments"`
}

type slackAttachment struct {
	Color    string                 `json:"color"`
	Pretext  string                 `json:"pretext"`
	Fallback string                 `json:"fallback"`
	Fields   []slackAttachmentField `json:"fields"`
}

type slackAttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

func getSlackPayload(success bool, pretext string, fallback string, targetVersion string, extraErrorText string) (io.Reader, error) {
	var color string
	if success {
		color = colorGreen
	} else {
		color = colorRed
	}

	if extraErrorText != "" {
		extraErrorText = fmt.Sprintf("Additional Error: %s", extraErrorText)
	}

	var fields []slackAttachmentField
	// targetVersion only not empty when producing autoupgrade payload
	if targetVersion != "" {
		fields = []slackAttachmentField{
			{
				Title: "Target Package Version",
				Value: targetVersion,
				Short: false,
			},
		}
	}

	payload := slackPayload{
		Attachments: []slackAttachment{
			{
				Color:    color,
				Pretext:  fmt.Sprintf("%s%s", pretext, extraErrorText),
				Fallback: fallback,
				Fields:   fields,
			},
		},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal slack payload")
	}
	return bytes.NewReader(jsonPayload), nil
}
