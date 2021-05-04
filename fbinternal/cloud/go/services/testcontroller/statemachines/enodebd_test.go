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

// NOTE: to run these tests outside the testing environment, e.g. from IntelliJ,
// ensure postgres_test container is running, and use the following environment
// variables to point to the relevant DB endpoints:
//	- DATABASE_SOURCE=host=localhost port=5433 dbname=magma_test user=magma_test password=magma_test sslmode=disable

package statemachines_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"magma/fbinternal/cloud/go/services/testcontroller/obsidian/models"
	"magma/fbinternal/cloud/go/services/testcontroller/statemachines"
	storage2 "magma/fbinternal/cloud/go/services/testcontroller/storage"
	tcTestInit "magma/fbinternal/cloud/go/services/testcontroller/test_init"
	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	ltemodels "magma/lte/cloud/go/services/lte/obsidian/models"
	subscribermodels "magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	cfgTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/device"
	deviceTestInit "magma/orc8r/cloud/go/services/device/test_init"
	models2 "magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/go-openapi/swag"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// don't test intermediate failure conditions (e.g. unexpected config types,
// service errors)
func Test_EnodebdE2ETestStateMachine_HappyPath(t *testing.T) {
	SetupTests(t, "testcontroller__statemachines__enodebd_happy")
	RegisterAGW(t)
	cli := &mockClient{}
	testConfig := GetEnodebTestConfig()
	mockMagmad, mockGenericCommandResp := GetMockObjects()

	mockMagmad.On("RebootEnodeb", "n1", "g1", "1202000038269KP0037").Return(mockGenericCommandResp, nil)
	mockMagmad.On("GenerateTraffic", "n1", "g2", "magmawifi", "magmamagma").Return(mockGenericCommandResp, nil).Times(4)

	// New test
	sm := statemachines.NewEnodebdE2ETestStateMachine(tcTestInit.GetTestTestcontrollerStorage(t), cli, mockMagmad)
	actualState, actualDuration, err := sm.Run(storage2.CommonStartState, testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	// ---
	// Check for upgrade, find version equal to what tier is configured to; expect epsilon transition, 20 minute delay
	// ---
	testdata, err := ioutil.ReadFile("../testdata/testdata")
	assert.NoError(t, err)
	mockResp := &http.Response{Status: "200", Body: ioutil.NopCloser(bytes.NewBuffer(testdata))}
	cli.On("Get", mock.AnythingOfType("string")).Return(mockResp, nil).Times(1)

	actualState, actualDuration, err = sm.Run("check_for_upgrade", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, 20*time.Minute, actualDuration)

	// ---
	// Check for upgrade find version ahead of what tier is configured to
	// ---
	err = configurator.CreateOrUpdateEntityConfig("n1", orc8r.UpgradeTierEntityType, "t1", &models2.Tier{Version: "0.0.0-0-abcdefg"}, serdes.Entity)
	assert.NoError(t, err)
	mockResp = &http.Response{Status: "200", Body: ioutil.NopCloser(bytes.NewBuffer(testdata))}
	cli.On("Get", mock.AnythingOfType("string")).Return(mockResp, nil).Times(1)

	actualState, actualDuration, err = sm.Run("check_for_upgrade", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "verify_upgrade_1", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	// Tier should get updated
	actualTierCfg, err := configurator.LoadEntityConfig("n1", orc8r.UpgradeTierEntityType, "t1", serdes.Entity)
	assert.NoError(t, err)
	assert.Equal(t, &models2.Tier{Version: "0.3.74-1560824953-b50f1bab"}, actualTierCfg)

	// ---
	// Check upgrade status, gateway hasn't upgraded yet
	// ---
	gatewayRecord := &models2.GatewayDevice{HardwareID: "hw1", Key: &models2.ChallengeKey{KeyType: "ECHO"}}
	err = device.RegisterDevice("n1", orc8r.AccessGatewayRecordType, "hw1", gatewayRecord, serdes.Device)
	assert.NoError(t, err)
	ctx := test_utils.GetContextWithCertificate(t, "hw1")
	test_utils.ReportGatewayStatus(t, ctx, &models2.GatewayStatus{
		HardwareID: "hw1",
		PlatformInfo: &models2.PlatformInfo{
			Packages: []*models2.Package{
				{Name: "magma", Version: "0.0.0-0-abcdefg"},
			},
		},
	})

	mockStatus := ltemodels.NewDefaultEnodebStatus()
	mockStatus.RfTxOn = swag.Bool(true)
	mockStatus.RfTxDesired = swag.Bool(true)
	reportEnodebState(t, ctx, "1202000038269KP0037", mockStatus)

	actualState, actualDuration, err = sm.Run("verify_upgrade_1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "verify_upgrade_2", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	// ---
	// Upgrade successful
	// ---
	mockResp = &http.Response{Status: "200", StatusCode: 200}
	// Should test for the payload eventually
	cli.On("Post", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return(mockResp, nil)
	test_utils.ReportGatewayStatus(t, ctx, &models2.GatewayStatus{
		HardwareID: "hw1",
		PlatformInfo: &models2.PlatformInfo{
			Packages: []*models2.Package{
				{Name: "magma", Version: "0.3.74-1560824953-b50f1bab"},
			},
		},
	})

	actualState, actualDuration, err = sm.Run("verify_upgrade_1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "traffic_test1_1", actualState)
	assert.Equal(t, 20*time.Minute, actualDuration)

	// ---
	// Traffic test 1
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test1_1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "reboot_enodeb_1", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	// ---
	// Reboot enodeb
	// ---
	actualState, actualDuration, err = sm.Run("reboot_enodeb_1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "verify_conn", actualState)
	assert.Equal(t, 15*time.Minute, actualDuration)

	// ---
	// Verify enodeb connectivity
	// ---
	actualState, actualDuration, err = sm.Run("verify_conn", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "traffic_test2_1", actualState)
	assert.Equal(t, 15*time.Minute, actualDuration)

	// ---
	// Traffic test 2
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test2_1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "reconfig_enodeb1", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	// ---
	// Reconfig Enodeb
	// ---
	actualState, actualDuration, err = sm.Run("reconfig_enodeb1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "verify_config1", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	// ---
	// Verify Config 1
	// ---
	actualState, actualDuration, err = sm.Run("verify_config1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "traffic_test3_1", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	// ---
	// Traffic Test 3
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test3_1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "restore_enodeb1", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	// ---
	// Restore Enodeb config
	// ---
	actualState, actualDuration, err = sm.Run("restore_enodeb1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "verify_config2", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	// ---
	// Verify Config 2
	// ---
	actualState, actualDuration, err = sm.Run("verify_config2", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "traffic_test4_1", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	// ---
	// Traffic Test 4
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test4_1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "subscriber_inactive", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	// ---
	// Detach subscriber
	// ---
	actualState, actualDuration, err = sm.Run("subscriber_inactive", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "traffic_test5_1", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	mockMagmad.On("GenerateTraffic", "n1", "g2", "magmawifi", "magmamagma").Return(mockGenericCommandResp, errors.New("")).Once()
	// ---
	// Traffic Test 5
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test5_1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "subscriber_active", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	// ---
	// Attach subscriber
	// ---
	actualState, actualDuration, err = sm.Run("subscriber_active", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "traffic_test6_1", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	mockMagmad.On("GenerateTraffic", "n1", "g2", "magmawifi", "magmamagma").Return(mockGenericCommandResp, nil)
	// ---
	// Traffic Test 6
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test6_1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	cli.On("Post", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return(mockResp, nil)
	// ---
	// Upgrade unsuccessful
	// ---
	test_utils.ReportGatewayStatus(t, ctx, &models2.GatewayStatus{
		HardwareID: "hw1",
		PlatformInfo: &models2.PlatformInfo{
			Packages: []*models2.Package{
				{Name: "magma", Version: "0.0.0-0-abcdefg"},
			},
		},
	})
	actualState, actualDuration, err = sm.Run("verify_upgrade_3", testConfig, nil)
	assert.EqualError(t, err, "gateway g1 did not upgrade within 3 tries")
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, 20*time.Minute, actualDuration)

	cli.AssertExpectations(t)
	mockMagmad.AssertExpectations(t)
}

func Test_EnodebdE2ETestStateMachine_VerifyConnection(t *testing.T) {
	SetupTests(t, "testcontroller__statemachines__enodebd_verify")
	RegisterAGW(t)
	cli := &mockClient{}
	testConfig := GetEnodebTestConfig()
	mockMagmad, mockGenericCommandResp := GetMockObjects()

	mockResp := &http.Response{Status: "200", StatusCode: 200}
	// Should test for the payload eventually
	cli.On("Post", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return(mockResp, nil)

	// New test
	sm := statemachines.NewEnodebdE2ETestStateMachine(tcTestInit.GetTestTestcontrollerStorage(t), cli, mockMagmad)

	mockMagmad.On("RebootEnodeb", "n1", "g1", "1202000038269KP0037").Return(mockGenericCommandResp, errors.New("")).Twice()
	mockStatus := ltemodels.NewDefaultEnodebStatus()
	mockStatus.RfTxOn = swag.Bool(false)
	mockStatus.RfTxDesired = swag.Bool(false)
	ctx := test_utils.GetContextWithCertificate(t, "hw1")
	reportEnodebState(t, ctx, "1202000038269KP0037", mockStatus)

	// ---
	// reboot_enodeb_1 transition to reboot_enodeb_2
	// --
	actualState, actualDuration, err := sm.Run("reboot_enodeb_1", testConfig, nil)
	assert.EqualError(t, err, "")
	assert.Equal(t, "reboot_enodeb_2", actualState)
	assert.Equal(t, 5*time.Minute, actualDuration)

	// ---
	// Reboot unsuccessful
	// --
	actualState, actualDuration, err = sm.Run("reboot_enodeb_3", testConfig, nil)
	assert.EqualError(t, err, "enodeb 1202000038269KP0037 did not reboot within 3 tries")
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, 15*time.Minute, actualDuration)

	// ---
	// Verify enodeb connectivity unsuccessful due to RfTx mismatch
	// ---
	actualState, actualDuration, err = sm.Run("verify_conn", testConfig, nil)
	assert.EqualError(t, err, "Error RF TX on/desired are not both set to true")
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, 5*time.Minute, actualDuration)

	mockMagmad.On("RebootEnodeb", "n1", "g1", "1202000038269KP0037").Return(mockGenericCommandResp, nil)
	mockStatus.RfTxOn = swag.Bool(true)
	mockStatus.RfTxDesired = swag.Bool(true)
	reportEnodebState(t, ctx, "1202000038269KP0037", mockStatus)

	// ---
	// Reboot enodeb
	// ---
	actualState, actualDuration, err = sm.Run("reboot_enodeb_1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "verify_conn", actualState)
	assert.Equal(t, 15*time.Minute, actualDuration)

	// ---
	// Verify enodeb connectivity
	// ---
	actualState, actualDuration, err = sm.Run("verify_conn", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "traffic_test2_1", actualState)
	assert.Equal(t, 15*time.Minute, actualDuration)

	cli.AssertExpectations(t)
	mockMagmad.AssertExpectations(t)
}

func Test_EnodebdE2ETestStateMachine_TrafficScript(t *testing.T) {
	SetupTests(t, "testcontroller__statemachines__enodebd_traffic")
	RegisterAGW(t)
	cli := &mockClient{}
	testConfig := GetEnodebTestConfig()
	mockMagmad, mockGenericCommandResp := GetMockObjects()

	// New test
	sm := statemachines.NewEnodebdE2ETestStateMachine(tcTestInit.GetTestTestcontrollerStorage(t), cli, mockMagmad)

	mockResp := &http.Response{Status: "200", StatusCode: 200}
	// Should test for the payload eventually
	cli.On("Post", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return(mockResp, nil)

	mockMagmad.On("GenerateTraffic", "n1", "g2", "magmawifi", "magmamagma").Return(mockGenericCommandResp, errors.New("")).Times(4)
	// ---
	// Traffic test 1 fails, cannot connect to AGW
	// ---
	actualState, actualDuration, err := sm.Run("traffic_test1_1", testConfig, nil)
	assert.EqualError(t, err, "")
	assert.Equal(t, "traffic_test1_2", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	actualState, actualDuration, err = sm.Run("traffic_test1_3", testConfig, nil)
	assert.EqualError(t, err, "Traffic test number 1 failed on gwID g2 after 3 tries")
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	// ---
	// Traffic test 2 fails, cannot connect to AGW
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test2_1", testConfig, nil)
	assert.EqualError(t, err, "")
	assert.Equal(t, "traffic_test2_2", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	actualState, actualDuration, err = sm.Run("traffic_test2_3", testConfig, nil)
	assert.EqualError(t, err, "Traffic test number 2 failed on gwID g2 after 3 tries")
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	mockMagmad.On("GenerateTraffic", "n1", "g2", "magmawifi", "magmamagma").Return(mockGenericCommandResp, nil)
	mockGenericCommandResp.Response.Fields = map[string]*structpb.Value{
		"returncode": {Kind: &structpb.Value_NumberValue{NumberValue: float64(1)}},
		"stdout":     {Kind: &structpb.Value_StringValue{StringValue: ""}},
		"stderr":     {Kind: &structpb.Value_StringValue{StringValue: "unable to connect to magmawifi"}},
	}
	// ---
	// Traffic test 1 fails, traffic script failing
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test1_1", testConfig, nil)
	assert.EqualError(t, err, "Traffic script failed. Return Code: 1, Stdout: , Stderr: unable to connect to magmawifi")
	assert.Equal(t, "traffic_test1_2", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	actualState, actualDuration, err = sm.Run("traffic_test1_3", testConfig, nil)
	assert.EqualError(t, err, "Traffic test number 1 failed on gwID g2 after 3 tries")
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	// ---
	// Traffic test 2 fails, traffic script failing
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test2_1", testConfig, nil)
	assert.EqualError(t, err, "Traffic script failed. Return Code: 1, Stdout: , Stderr: unable to connect to magmawifi")
	assert.Equal(t, "traffic_test2_2", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	actualState, actualDuration, err = sm.Run("traffic_test2_3", testConfig, nil)
	assert.EqualError(t, err, "Traffic test number 2 failed on gwID g2 after 3 tries")
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	mockGenericCommandResp.Response.Fields = map[string]*structpb.Value{
		"returncode": {Kind: &structpb.Value_NumberValue{NumberValue: float64(0)}},
		"stdout":     {Kind: &structpb.Value_StringValue{StringValue: ""}},
		"stderr":     {Kind: &structpb.Value_StringValue{StringValue: ""}},
	}
	// ---
	// Traffic Test 1
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test1_1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "reboot_enodeb_1", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	// ---
	// Traffic Test 2
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test2_1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "reconfig_enodeb1", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	// ---
	// Successful traffic test in state 3
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test1_3", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "reboot_enodeb_1", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	actualState, actualDuration, err = sm.Run("traffic_test2_3", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "reconfig_enodeb1", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	cli.AssertExpectations(t)
	mockMagmad.AssertExpectations(t)
}

func Test_EnodebdE2ETestStateMachine_ReconfigEnb(t *testing.T) {
	SetupTests(t, "testcontroller__statemachines__enodebd_reconfig")
	RegisterAGW(t)
	cli := &mockClient{}
	testConfig := GetEnodebTestConfig()
	mockMagmad, mockGenericCommandResp := GetMockObjects()
	mockMagmad.On("GenerateTraffic", "n1", "g2", "magmawifi", "magmamagma").Return(mockGenericCommandResp, errors.New("")).Times(4)
	mockResp := &http.Response{Status: "200", StatusCode: 200}
	// Should test for the payload eventually
	cli.On("Post", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return(mockResp, nil)

	ctx := test_utils.GetContextWithCertificate(t, "hw1")
	mockStatus := ltemodels.NewDefaultEnodebStatus()
	mockStatus.EnodebConfigured = swag.Bool(false)
	reportEnodebState(t, ctx, "1202000038269KP0037", mockStatus)

	// New test
	sm := statemachines.NewEnodebdE2ETestStateMachine(tcTestInit.GetTestTestcontrollerStorage(t), cli, mockMagmad)

	testConfig.EnodebConfig.ManagedConfig.Pci = 261
	// ---
	// Reconfig Enodeb
	// ---
	actualState, actualDuration, err := sm.Run("reconfig_enodeb1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "verify_config1", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	// ---
	// Verify Enb Config failing
	// ---
	actualState, actualDuration, err = sm.Run("verify_config1", testConfig, nil)
	assert.EqualError(t, err, "error enodeb 1202000038269KP0037 is not configured")
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	// ---
	// Traffic Test 3 fails, cannot connect to AGW
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test3_1", testConfig, nil)
	assert.EqualError(t, err, "")
	assert.Equal(t, "traffic_test3_2", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	actualState, actualDuration, err = sm.Run("traffic_test3_3", testConfig, nil)
	assert.EqualError(t, err, "Traffic test number 3 failed on gwID g2 after 3 tries")
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	testConfig.EnodebConfig.ManagedConfig.Pci = 260
	// ---
	// Restore Enodeb config
	// ---
	actualState, actualDuration, err = sm.Run("restore_enodeb1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "verify_config2", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	// ---
	// Verify Enb Config 2 failing
	// ---
	actualState, actualDuration, err = sm.Run("verify_config2", testConfig, nil)
	assert.EqualError(t, err, "error enodeb 1202000038269KP0037 is not configured")
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	// ---
	// Traffic Test 4 fails, cannot connect to AGW
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test4_1", testConfig, nil)
	assert.EqualError(t, err, "")
	assert.Equal(t, "traffic_test4_2", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	actualState, actualDuration, err = sm.Run("traffic_test4_3", testConfig, nil)
	assert.EqualError(t, err, "Traffic test number 4 failed on gwID g2 after 3 tries")
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	mockStatus.EnodebConfigured = swag.Bool(true)
	reportEnodebState(t, ctx, "1202000038269KP0037", mockStatus)
	mockMagmad.On("GenerateTraffic", "n1", "g2", "magmawifi", "magmamagma").Return(mockGenericCommandResp, nil)
	testConfig.EnodebConfig.ManagedConfig.Pci = 261
	mockGenericCommandResp.Response.Fields = map[string]*structpb.Value{
		"returncode": {Kind: &structpb.Value_NumberValue{NumberValue: float64(1)}},
		"stdout":     {Kind: &structpb.Value_StringValue{StringValue: ""}},
		"stderr":     {Kind: &structpb.Value_StringValue{StringValue: "unable to connect to magmawifi"}},
	}

	// ---
	// Traffic Test 3 fails, traffic script failing
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test3_1", testConfig, nil)
	assert.EqualError(t, err, "Traffic script failed. Return Code: 1, Stdout: , Stderr: unable to connect to magmawifi")
	assert.Equal(t, "traffic_test3_2", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	actualState, actualDuration, err = sm.Run("traffic_test3_3", testConfig, nil)
	assert.EqualError(t, err, "Traffic test number 3 failed on gwID g2 after 3 tries")
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	// ---
	// Traffic Test 4 fails, traffic script failing
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test4_1", testConfig, nil)
	assert.EqualError(t, err, "Traffic script failed. Return Code: 1, Stdout: , Stderr: unable to connect to magmawifi")
	assert.Equal(t, "traffic_test4_2", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	actualState, actualDuration, err = sm.Run("traffic_test4_3", testConfig, nil)
	assert.EqualError(t, err, "Traffic test number 4 failed on gwID g2 after 3 tries")
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	mockGenericCommandResp.Response.Fields = map[string]*structpb.Value{
		"returncode": {Kind: &structpb.Value_NumberValue{NumberValue: float64(0)}},
		"stdout":     {Kind: &structpb.Value_StringValue{StringValue: ""}},
		"stderr":     {Kind: &structpb.Value_StringValue{StringValue: ""}},
	}

	// ---
	// Reconfig Enodeb
	// ---
	actualState, actualDuration, err = sm.Run("reconfig_enodeb1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "verify_config1", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	// ---
	// Verify Enb Config from original to new config succeeding
	// ---
	actualState, actualDuration, err = sm.Run("verify_config1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "traffic_test3_1", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	// ---
	// Traffic Test 3
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test3_1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "restore_enodeb1", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	testConfig.EnodebConfig.ManagedConfig.Pci = 260
	// ---
	// Restore Enodeb config
	// ---
	actualState, actualDuration, err = sm.Run("restore_enodeb1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "verify_config2", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	// ---
	// Verify Enb Config 2 succeeding
	// ---
	actualState, actualDuration, err = sm.Run("verify_config2", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "traffic_test4_1", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	// ---
	// Traffic Test 4
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test4_1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "subscriber_inactive", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	cli.AssertExpectations(t)
	mockMagmad.AssertExpectations(t)
}

func Test_EnodebdE2ETestStateMachine_SubscriberState(t *testing.T) {
	SetupTests(t, "testcontroller__statemachines__enodebd_subscriber_state")
	RegisterAGW(t)
	testConfig := GetEnodebTestConfig()
	cli := &mockClient{}
	mockMagmad, mockGenericCommandResp := GetMockObjects()

	mockResp := &http.Response{Status: "200", StatusCode: 200}
	// Should test for the payload eventually
	cli.On("Post", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return(mockResp, nil)
	mockMagmad.On("GenerateTraffic", "n1", "g2", "magmawifi", "magmamagma").Return(mockGenericCommandResp, nil).Once()

	// New test
	sm := statemachines.NewEnodebdE2ETestStateMachine(tcTestInit.GetTestTestcontrollerStorage(t), cli, mockMagmad)

	testConfig.SubscriberID = swag.String("IMSI0987654321")
	// ---
	// Subscriber state set to inactive fail
	// ---
	actualState, actualDuration, err := sm.Run("subscriber_inactive", testConfig, nil)
	assert.EqualError(t, err, "Not found")
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	// ---
	// Traffic test 5 fail, traffic test ran when it should not have
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test5_1", testConfig, nil)
	assert.EqualError(t, err, "Traffic test number 5 should not have succeeded on gwID g2")
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	// ---
	// Subscriber state set to active fail
	// ---
	actualState, actualDuration, err = sm.Run("subscriber_active", testConfig, nil)
	assert.EqualError(t, err, "Not found")
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	testConfig.SubscriberID = swag.String("IMSI1234567890")
	mockMagmad.On("GenerateTraffic", "n1", "g2", "magmawifi", "magmamagma").Return(mockGenericCommandResp, errors.New("")).Times(3)

	// ---
	// Traffic test 6 fail, cannot connect to AGW
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test6_1", testConfig, nil)
	assert.EqualError(t, err, "")
	assert.Equal(t, "traffic_test6_2", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	actualState, actualDuration, err = sm.Run("traffic_test6_3", testConfig, nil)
	assert.EqualError(t, err, "Traffic test number 6 failed on gwID g2 after 3 tries")
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	// ---
	// Subscriber state set to inactive success
	// ---
	actualState, actualDuration, err = sm.Run("subscriber_inactive", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "traffic_test5_1", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	// ---
	// Traffic test should not succeed due to inactive subscriber. This special case does not return an error
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test5_1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "subscriber_active", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	mockMagmad.On("GenerateTraffic", "n1", "g2", "magmawifi", "magmamagma").Return(mockGenericCommandResp, nil)
	// ---
	// Subscriber state set to active success
	// ---
	actualState, actualDuration, err = sm.Run("subscriber_active", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "traffic_test6_1", actualState)
	assert.Equal(t, 10*time.Minute, actualDuration)

	// ---
	// Traffic test should succeed after flipping subscriber to active
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test6_1", testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	mockGenericCommandResp.Response.Fields = map[string]*structpb.Value{
		"returncode": {Kind: &structpb.Value_NumberValue{NumberValue: float64(1)}},
		"stdout":     {Kind: &structpb.Value_StringValue{StringValue: ""}},
		"stderr":     {Kind: &structpb.Value_StringValue{StringValue: "unable to connect to magmawifi"}},
	}
	// ---
	// Traffic test 6 fail, traffic script failing
	// ---
	actualState, actualDuration, err = sm.Run("traffic_test6_1", testConfig, nil)
	assert.EqualError(t, err, "Traffic script failed. Return Code: 1, Stdout: , Stderr: unable to connect to magmawifi")
	assert.Equal(t, "traffic_test6_2", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	actualState, actualDuration, err = sm.Run("traffic_test6_3", testConfig, nil)
	assert.EqualError(t, err, "Traffic test number 6 failed on gwID g2 after 3 tries")
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	cli.AssertExpectations(t)
	mockMagmad.AssertExpectations(t)
}

func Test_EnodebdE2ETestStateMachine_DynamicStartState(t *testing.T) {
	SetupTests(t, "testcontroller__statemachines__enodebd_dynamic_start_state")
	RegisterAGW(t)
	testConfig := GetEnodebTestConfig()
	testConfig.StartState = "reboot_enodeb_1"
	cli := &mockClient{}
	mockMagmad, _ := GetMockObjects()

	// New test
	sm := statemachines.NewEnodebdE2ETestStateMachine(tcTestInit.GetTestTestcontrollerStorage(t), cli, mockMagmad)
	actualState, actualDuration, err := sm.Run(storage2.CommonStartState, testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "reboot_enodeb_1", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	testConfig.StartState = "reconfig_enodeb1"
	actualState, actualDuration, err = sm.Run(storage2.CommonStartState, testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "reconfig_enodeb1", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	testConfig.StartState = "subscriber_inactive"
	actualState, actualDuration, err = sm.Run(storage2.CommonStartState, testConfig, nil)
	assert.NoError(t, err)
	assert.Equal(t, "subscriber_inactive", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	testConfig.StartState = "INVALID STATE"
	// Invalid state should default to check_for_upgrade
	actualState, actualDuration, err = sm.Run(storage2.CommonStartState, testConfig, nil)
	assert.EqualError(t, err, "Invalid starting state. Defaulting to check_for_upgrade")
	assert.Equal(t, "check_for_upgrade", actualState)
	assert.Equal(t, time.Minute, actualDuration)

	cli.AssertExpectations(t)
	mockMagmad.AssertExpectations(t)
}
func SetupTests(t *testing.T, dbName string) {
	tcTestInit.StartTestServiceWithDB(t, dbName)
	cfgTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)

	frozenClock := 1000 * time.Hour
	clock.SetAndFreezeClock(t, time.Unix(0, 0).Add(frozenClock))
	defer clock.UnfreezeClock(t)
}

func RegisterAGW(t *testing.T) {
	// Register an AGW
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{Type: orc8r.UpgradeTierEntityType, Key: "t1", Config: &models2.Tier{Name: "t1", Version: "0.3.74-1560824953-b50f1bab"}},
		serdes.Entity,
	)
	assert.NoError(t, err)
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type:         orc8r.MagmadGatewayType,
				Key:          "g1",
				Config:       &models2.MagmadGatewayConfigs{},
				PhysicalID:   "hw1",
				Associations: []storage.TypeAndKey{{Type: orc8r.UpgradeTierEntityType, Key: "t1"}},
			},
			{
				Type:       lte.CellularEnodebEntityType,
				Key:        "1202000038269KP0037",
				PhysicalID: "1202000038269KP0037",
				Config: &ltemodels.EnodebConfiguration{
					BandwidthMhz:           20,
					CellID:                 swag.Uint32(1234),
					DeviceClass:            "Baicells Nova-233 G2 OD FDD",
					Earfcndl:               39450,
					Pci:                    260,
					SpecialSubframePattern: 7,
					SubframeAssignment:     2,
					Tac:                    1,
					TransmitEnabled:        swag.Bool(true),
				},
			},
			{
				Type:        lte.SubscriberEntityType,
				Key:         "IMSI1234567890",
				Name:        "subscriber1",
				Description: "mock subscriber",
				Config: &subscribermodels.SubscriberConfig{
					Lte: &subscribermodels.LteSubscription{
						AuthAlgo:   "MILENAGE",
						AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
						AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
						State:      "ACTIVE",
						SubProfile: "default",
					},
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
}

func GetEnodebTestConfig() *models.EnodebdTestConfig {
	testConfig := &models.EnodebdTestConfig{
		AgwConfig: &models.AgwTestConfig{
			PackageRepo:     swag.String("https://artifactory.magmacore.org/artifactory/debian"),
			ReleaseChannel:  swag.String("stretch-1.5.0"),
			SlackWebhook:    swag.String("https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"),
			TargetGatewayID: swag.String("g1"),
			TargetTier:      swag.String("t1"),
		},
		EnodebSN:    swag.String("1202000038269KP0037"),
		NetworkID:   swag.String("n1"),
		Ssid:        ("magmawifi"),
		SsidPw:      ("magmamagma"),
		TrafficGwID: swag.String("g2"),
		EnodebConfig: &ltemodels.EnodebConfig{
			ConfigType: "MANAGED",
			ManagedConfig: &ltemodels.EnodebConfiguration{
				BandwidthMhz:           20,
				CellID:                 swag.Uint32(138777000),
				DeviceClass:            "Baicells ID TDD/FDD",
				Earfcndl:               44590,
				Pci:                    260,
				SpecialSubframePattern: 7,
				SubframeAssignment:     2,
				Tac:                    1,
				TransmitEnabled:        swag.Bool(true),
			},
		},
		SubscriberID: swag.String("IMSI1234567890"),
		StartState:   "check_for_upgrade",
	}
	return testConfig
}

func GetMockObjects() (*mockMagmadClient, *protos.GenericCommandResponse) {
	mockMagmad := &mockMagmadClient{}
	mockResponse := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"returncode": {Kind: &structpb.Value_NumberValue{NumberValue: float64(0)}},
			"stdout":     {Kind: &structpb.Value_StringValue{StringValue: ""}},
			"stderr":     {Kind: &structpb.Value_StringValue{StringValue: ""}},
		},
	}
	mockGenericCommandResp := &protos.GenericCommandResponse{
		Response: mockResponse,
	}
	return mockMagmad, mockGenericCommandResp
}

type mockMagmadClient struct {
	mock.Mock
}

func (m *mockMagmadClient) GenerateTraffic(networkId string, trafficGatewayId string, ssid string, pw string) (*protos.GenericCommandResponse, error) {
	args := m.Called(networkId, trafficGatewayId, ssid, pw)
	return args.Get(0).(*protos.GenericCommandResponse), args.Error(1)
}

func (m *mockMagmadClient) RebootEnodeb(networkId string, gatewayId string, enodebSerial string) (*protos.GenericCommandResponse, error) {
	args := m.Called(networkId, gatewayId, enodebSerial)
	return args.Get(0).(*protos.GenericCommandResponse), args.Error(1)
}

func reportEnodebState(t *testing.T, ctx context.Context, enodebSerial string, req *ltemodels.EnodebState) {
	client, err := state.GetStateClient()
	assert.NoError(t, err)

	serializedEnodebState, err := serde.Serialize(req, lte.EnodebStateType, serdes.State)
	assert.NoError(t, err)
	states := []*protos.State{
		{
			Type:     lte.EnodebStateType,
			DeviceID: enodebSerial,
			Value:    serializedEnodebState,
		},
	}
	_, err = client.ReportStates(
		ctx,
		&protos.ReportStatesRequest{States: states},
	)
	assert.NoError(t, err)
}

type mockClient struct {
	mock.Mock
}

func (client *mockClient) Get(url string) (resp *http.Response, err error) {
	args := client.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (client *mockClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	args := client.Called(url, contentType, body)
	return args.Get(0).(*http.Response), args.Error(1)
}
