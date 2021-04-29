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

package integration

import (
	"context"
	"fmt"

	"magma/cwf/gateway/registry"
	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/object_store"
	"magma/feg/gateway/policydb"
	"magma/feg/gateway/services/testcore/hss"
	lteprotos "magma/lte/cloud/go/protos"
	registryTestUtils "magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
	orc8rprotos "magma/orc8r/lib/go/protos"

	"github.com/go-redis/redis"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	// PCRFinstances List of all possible PCRFs on the system (currently 2 are supported)
	PCRFinstances = []string{MockPCRFRemote, MockPCRFRemote2}
	// OCSinstances List of all possible OCSs on the system (currently 2 are supported)
	OCSinstances = []string{MockOCSRemote, MockOCSRemote2}
)

// Wrapper for GRPC Client functionality
type hssClient struct {
	fegprotos.HSSConfiguratorClient
	cc *grpc.ClientConn
}

// Wrapper for GRPC Client functionality
type pcrfClient struct {
	fegprotos.MockPCRFClient
	cc *grpc.ClientConn
}

// Wrapper for GRPC Client functionality
type ocsClient struct {
	fegprotos.MockOCSClient
	cc *grpc.ClientConn
}

// Wrapper for GRPC Client functionality
type pipelinedClient struct {
	lteprotos.PipelinedClient
	cc *grpc.ClientConn
}

// Wrapper for GRPC Client functionality
type directorydClient struct {
	orc8rprotos.GatewayDirectoryServiceClient
	cc *grpc.ClientConn
}

// Wrapper for PolicyDB objects
type policyDBWrapper struct {
	redisClient      object_store.RedisClient
	policyMap        object_store.ObjectMap
	baseNameMap      object_store.ObjectMap
	omniPresentRules object_store.ObjectMap
}

// TODO: convert all these helpers into methods for each type

/**  ========== HSS Helpers ========== **/
// getHSSClient is a utility function to getHSSClient a RPC connection to a
// remote HSS service.
func getHSSClient() (*hssClient, error) {
	var conn *grpc.ClientConn
	var err error
	conn, err = registry.GetConnection(MockHSSRemote)
	if err != nil {
		errMsg := fmt.Sprintf("HSS client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &hssClient{
		fegprotos.NewHSSConfiguratorClient(conn),
		conn,
	}, err
}

// addSubscriberToHSS tries to add this subscriber to the HSS server.
// This function returns an AlreadyExists error if the subscriber has already
// been added.
// Input: The subscriber data which will be added.
func addSubscriberToHSS(sub *lteprotos.SubscriberData) error {
	err := hss.VerifySubscriberData(sub)
	if err != nil {
		errMsg := fmt.Errorf("Invalid AddSubscriberRequest provided: %s", err)
		return errors.New(errMsg.Error())
	}
	cli, err := getHSSClient()
	if err != nil {
		return err
	}
	_, err = cli.AddSubscriber(context.Background(), sub)
	return err
}

// deleteSubscriberFromHSS tries to add delete subscriber from the HSS server.
// If the subscriber is not found, then this call is ignored.
// Input: The id of the subscriber to be deleted.
func deleteSubscribersFromHSS(subscriberID string) error {
	cli, err := getHSSClient()
	if err != nil {
		return err
	}
	_, err = cli.DeleteSubscriber(context.Background(), &lteprotos.SubscriberID{Id: subscriberID})
	return err
}

/**  ========== PCRF Helpers ========== **/
// getPCRFClient is a utility function to get an RPC connection to a
// remote PCRF service.
func getPCRFClient(instanceName string) (*pcrfClient, error) {
	var conn *grpc.ClientConn
	var err error
	if !contains(instanceName, PCRFinstances) {
		return nil,
			fmt.Errorf("mockPCRF Instance does not exist, use one of the existings")
	}
	conn, err = registry.GetConnection(instanceName)
	if err != nil {
		errMsg := fmt.Sprintf("PCRF client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &pcrfClient{
		MockPCRFClient: fegprotos.NewMockPCRFClient(conn),
		cc:             conn,
	}, err
}

// sendPolicyReAuthRequest initiates a RAR request from PCRF server to sessiond.
// Input: Policy RAR Target
func sendPolicyReAuthRequest(target *fegprotos.PolicyReAuthTarget) (*fegprotos.PolicyReAuthAnswer, error) {
	return sendPolicyReAuthRequestPerInstance(MockPCRFRemote, target)
}

func sendPolicyReAuthRequestPerInstance(instanceName string, target *fegprotos.PolicyReAuthTarget) (*fegprotos.PolicyReAuthAnswer, error) {
	cli, err := getPCRFClient(instanceName)
	if err != nil {
		return nil, err
	}
	raa, err := cli.ReAuth(context.Background(), target)
	return raa, err
}

// sendPolicyAbortSession initiates an abort request from PCRF server to sessiond.
// Input: Policy abort session request
func sendPolicyAbortSession(target *fegprotos.AbortSessionRequest) (*fegprotos.AbortSessionAnswer, error) {
	return sendPolicyAbortSessionPerInstance(MockPCRFRemote, target)
}

func sendPolicyAbortSessionPerInstance(instanceName string, target *fegprotos.AbortSessionRequest) (*fegprotos.AbortSessionAnswer, error) {
	cli, err := getPCRFClient(instanceName)
	if err != nil {
		return nil, err
	}
	raa, err := cli.AbortSession(context.Background(), target)
	return raa, err
}

// addSubscriberToPCRF tries to add this subscriber to the PCRF server.
// Input: The subscriber data which will be added.
func addSubscriberToPCRF(sub *lteprotos.SubscriberID) error {
	return addSubscriberToPCRFPerInstance(MockPCRFRemote, sub)
}

func addSubscriberToPCRFPerInstance(instanceName string, sub *lteprotos.SubscriberID) error {
	cli, err := getPCRFClient(instanceName)
	if err != nil {
		return err
	}
	_, err = cli.CreateAccount(context.Background(), sub)
	return err
}

// clearSubscriberToPCRF tries to clear all subscribers from the PCRF server.
func clearSubscribersFromPCRF() error {
	return clearSubscribersFromPCRFPerInstance(MockPCRFRemote)
}

func clearSubscribersFromPCRFPerInstance(instanceName string) error {
	cli, err := getPCRFClient(instanceName)
	if err != nil {
		return err
	}
	_, err = cli.ClearSubscribers(context.Background(), &protos.Void{})
	return err
}

func addPCRFRules(rules *fegprotos.AccountRules) error {
	return addPCRFRulesPerInstance(MockPCRFRemote, rules)
}

func addPCRFRulesPerInstance(instanceName string, rules *fegprotos.AccountRules) error {
	cli, err := getPCRFClient(instanceName)
	if err != nil {
		return err
	}
	_, err = cli.SetRules(context.Background(), rules)
	return err
}

func addPCRFUsageMonitors(monitorInfo *fegprotos.UsageMonitorConfiguration) error {
	return addPCRFUsageMonitorsPerInstance(MockPCRFRemote, monitorInfo)
}

func addPCRFUsageMonitorsPerInstance(instanceName string, monitorInfo *fegprotos.UsageMonitorConfiguration) error {
	cli, err := getPCRFClient(instanceName)
	if err != nil {
		return err
	}
	_, err = cli.SetUsageMonitors(context.Background(), monitorInfo)
	return err
}

// usePCRFMockDriver enable MockPCRFDriver
func usePCRFMockDriver() error {
	return usePCRFMockDriverPerInstance(MockPCRFRemote)
}

func usePCRFMockDriverPerInstance(instanceName string) error {
	cli, err := getPCRFClient(instanceName)
	if err != nil {
		return err
	}
	_, err = cli.SetPCRFConfigs(context.Background(), &fegprotos.PCRFConfigs{UseMockDriver: true})
	return err
}

// clearPCRFMockDriver disable MockPCRFDriver
func clearPCRFMockDriver() error {
	return clearPCRFMockDriverPerInstance(MockPCRFRemote)
}

func clearPCRFMockDriverPerInstance(instanceName string) error {
	cli, err := getPCRFClient(instanceName)
	if err != nil {
		return err
	}
	_, err = cli.SetPCRFConfigs(context.Background(), &fegprotos.PCRFConfigs{UseMockDriver: false})
	return err
}

// setPCRFExpectations allows to set the expectations in PCRF
func setPCRFExpectations(expectations []*fegprotos.GxCreditControlExpectation, defaultAnswer *fegprotos.GxCreditControlAnswer) error {
	return setPCRFExpectationsPerInstance(MockPCRFRemote, expectations, defaultAnswer)
}

func setPCRFExpectationsPerInstance(instanceName string, expectations []*fegprotos.GxCreditControlExpectation, defaultAnswer *fegprotos.GxCreditControlAnswer) error {
	cli, err := getPCRFClient(instanceName)
	if err != nil {
		return err
	}
	request := &fegprotos.GxCreditControlExpectations{
		Expectations: expectations,
		GxDefaultCca: defaultAnswer,
	}
	if defaultAnswer != nil {
		request.UnexpectedRequestBehavior = fegprotos.UnexpectedRequestBehavior_CONTINUE_WITH_DEFAULT_ANSWER
	}
	_, err = cli.SetExpectations(context.Background(), request)
	return err
}

// getPCRFAssertExpectationsResult allows to get the expectations results in PCRF
func getPCRFAssertExpectationsResult() ([]*fegprotos.ExpectationResult, []*fegprotos.ErrorByIndex, error) {
	return getPCRFAssertExpectationsResultPerInstance(MockPCRFRemote)
}

func getPCRFAssertExpectationsResultPerInstance(instanceName string) ([]*fegprotos.ExpectationResult, []*fegprotos.ErrorByIndex, error) {
	cli, err := getPCRFClient(instanceName)
	if err != nil {
		return nil, nil, nil
	}
	res, err := cli.AssertExpectations(context.Background(), &protos.Void{})
	if err != nil {
		return nil, nil, err
	}
	return res.Results, res.Errors, nil
}

/**  ========== OCS Helpers ========== **/
// getOCSClient is a utility function to an RPC connection to a
// remote OCS service.
func getOCSClient(instanceName string) (*ocsClient, error) {
	var conn *grpc.ClientConn
	var err error
	if !contains(instanceName, OCSinstances) {
		return nil,
			fmt.Errorf("mockOCS Instance does not exist, use one of the existings")
	}
	conn, err = registry.GetConnection(instanceName)
	if err != nil {
		errMsg := fmt.Sprintf("OCS client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &ocsClient{
		fegprotos.NewMockOCSClient(conn),
		conn,
	}, err
}

// setNewOCSConfig tries to override the default ocs settings
// Input: ocsConfig data
func setNewOCSConfig(ocsConfig *fegprotos.OCSConfig) error {
	return setNewOCSConfigPerInstance(MockOCSRemote, ocsConfig)
}

func setNewOCSConfigPerInstance(instanceName string, ocsConfig *fegprotos.OCSConfig) error {
	cli, err := getOCSClient(instanceName)
	if err != nil {
		return err
	}
	_, err = cli.SetOCSSettings(context.Background(), ocsConfig)
	return err
}

// addSubscriber tries to add this subscriber to the OCS server.
// Input: The subscriber data which will be added.
func addSubscriberToOCS(sub *lteprotos.SubscriberID) error {
	return addSubscriberToOCSPerInstance(MockOCSRemote, sub)
}

func addSubscriberToOCSPerInstance(instanceName string, sub *lteprotos.SubscriberID) error {
	cli, err := getOCSClient(instanceName)
	if err != nil {
		return err
	}
	_, err = cli.CreateAccount(context.Background(), sub)
	return err
}

func clearSubscribersFromOCS() error {
	return clearSubscribersFromOCSPerInstance(MockOCSRemote)
}

func clearSubscribersFromOCSPerInstance(instanceName string) error {
	cli, err := getOCSClient(instanceName)
	if err != nil {
		return err
	}
	_, err = cli.ClearSubscribers(context.Background(), &protos.Void{})
	return err
}

// setCreditOCS tries to set a credit for this subscriber to the OCS server
// Input: The credit info data which will be set
func setCreditOnOCS(creditInfo *fegprotos.CreditInfo) error {
	return setCreditOnOCSPerInstance(MockOCSRemote, creditInfo)
}

func setCreditOnOCSPerInstance(instanceName string, creditInfo *fegprotos.CreditInfo) error {
	cli, err := getOCSClient(instanceName)
	if err != nil {
		return err
	}
	_, err = cli.SetCredit(context.Background(), creditInfo)
	return err
}

// getCreditOnOCS tries to get credit for this subscriber to the OCS server
// Input: The credit info data which will be set
func getCreditOnOCS(imsi string) (*fegprotos.CreditInfos, error) {
	return getCreditOnOCSPerInstance(MockOCSRemote, imsi)
}

func getCreditOnOCSPerInstance(instanceName, imsi string) (*fegprotos.CreditInfos, error) {
	cli, err := getOCSClient(instanceName)
	if err != nil {
		return &fegprotos.CreditInfos{}, err
	}
	return cli.GetCredits(context.Background(), &lteprotos.SubscriberID{Id: imsi})
}

// sendChargingReAuthRequestEntireSession triggers a RAR from OCS to Sessiond for all
// the charging groups on a session (inclding the suspended ones)
func sendChargingReAuthRequestEntireSession(imsi string) (*fegprotos.ChargingReAuthAnswer, error) {
	return sendChargingReAuthRequestPerInstance(MockOCSRemote, imsi, 0)
}

// sendChargingReAuthRequest triggers a RAR from OCS to Sessiond
// Input: ChargingReAuthTarget
func sendChargingReAuthRequest(imsi string, ratingGroup uint32) (*fegprotos.ChargingReAuthAnswer, error) {
	return sendChargingReAuthRequestPerInstance(MockOCSRemote, imsi, ratingGroup)
}

func sendChargingReAuthRequestPerInstance(instanceName string, imsi string, ratingGroup uint32) (*fegprotos.ChargingReAuthAnswer, error) {
	cli, err := getOCSClient(instanceName)
	if err != nil {
		return nil, err
	}
	raa, err := cli.ReAuth(context.Background(), &fegprotos.ChargingReAuthTarget{Imsi: imsi, RatingGroup: ratingGroup})
	return raa, err
}

func sendChargingAbortSession(target *fegprotos.AbortSessionRequest) (*fegprotos.AbortSessionAnswer, error) {
	return sendChargingAbortSessionPerInstance(MockOCSRemote, target)
}

func sendChargingAbortSessionPerInstance(instanceName string, target *fegprotos.AbortSessionRequest) (*fegprotos.AbortSessionAnswer, error) {
	cli, err := getOCSClient(instanceName)
	if err != nil {
		return nil, err
	}
	asa, err := cli.AbortSession(context.Background(), target)
	return asa, err
}

// useOCSMockDriver enables MockOCSDriver
func useOCSMockDriver() error {
	return useOCSMockDriverPerInstance(MockOCSRemote)
}

func useOCSMockDriverPerInstance(instanceName string) error {
	cli, err := getOCSClient(instanceName)
	if err != nil {
		return err
	}
	_, err = cli.SetOCSSettings(context.Background(), &fegprotos.OCSConfig{UseMockDriver: true})
	return err
}

// clearOCSMockDriver disable MockOCSDriver
func clearOCSMockDriver() error {
	return clearOCSMockDriverPerInstance(MockOCSRemote)
}

func clearOCSMockDriverPerInstance(instanceName string) error {
	cli, err := getOCSClient(instanceName)
	if err != nil {
		return err
	}
	_, err = cli.SetOCSSettings(context.Background(), &fegprotos.OCSConfig{UseMockDriver: false})
	return err
}

// setOCSExpectations allows to set the expectations in OCS
func setOCSExpectations(expectations []*fegprotos.GyCreditControlExpectation, defaultAnswer *fegprotos.GyCreditControlAnswer) error {
	return setOCSExpectationsPerInstance(MockOCSRemote, expectations, defaultAnswer)
}

func setOCSExpectationsPerInstance(instanceName string, expectations []*fegprotos.GyCreditControlExpectation, defaultAnswer *fegprotos.GyCreditControlAnswer) error {
	cli, err := getOCSClient(instanceName)
	if err != nil {
		return err
	}
	request := &fegprotos.GyCreditControlExpectations{
		Expectations: expectations,
		GyDefaultCca: defaultAnswer,
	}
	if defaultAnswer != nil {
		request.UnexpectedRequestBehavior = fegprotos.UnexpectedRequestBehavior_CONTINUE_WITH_DEFAULT_ANSWER
	}
	_, err = cli.SetExpectations(context.Background(), request)
	return err
}

// getOCSAssertExpectationsResult allows to get the expectations results in OCS
func getOCSAssertExpectationsResult() ([]*fegprotos.ExpectationResult, []*fegprotos.ErrorByIndex, error) {
	return getOCSAssertExpectationsResultPerInstance(MockOCSRemote)
}

func getOCSAssertExpectationsResultPerInstance(instanceName string) ([]*fegprotos.ExpectationResult, []*fegprotos.ErrorByIndex, error) {
	cli, err := getOCSClient(instanceName)
	if err != nil {
		return nil, nil, nil
	}
	res, err := cli.AssertExpectations(context.Background(), &protos.Void{})
	if err != nil {
		return nil, nil, err
	}
	return res.Results, res.Errors, nil
}

/**  ========== Pipelined Helpers ========== **/
// getPipelinedClient is a utility function to an RPC connection to a
// remote Pipelined service.
func getPipelinedClient() (*pipelinedClient, error) {
	var conn *grpc.ClientConn
	var err error
	conn, err = registryTestUtils.GetConnectionWithAuthority(PipelinedRemote)
	if err != nil {
		errMsg := fmt.Sprintf("Pipelined client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &pipelinedClient{
		lteprotos.NewPipelinedClient(conn),
		conn,
	}, err
}

func deactivateAllFlowsPerSub(imsi string) error {
	cli, err := getPipelinedClient()
	if err == nil && cli != nil {
		_, err = cli.DeactivateFlows(context.Background(), &lteprotos.DeactivateFlowsRequest{
			Sid: &lteprotos.SubscriberID{Id: imsi},
		})
	}
	return err
}

func getPolicyUsage() (*lteprotos.RuleRecordTable, error) {
	client, _ := getPipelinedClient()
	return client.GetPolicyUsage(context.Background(), &protos.Void{})
}

/**  ========== PolicyDB related Helpers ========== **/
// In the actual CWAG setup, PolicyRules and BaseNames managed by PolicyDB are
// streamed down from the cloud. Since this integration test setup does not
// include the cloud, we will get around this by directly modifying the redis
// DB.
func initializePolicyDBWrapper() (*policyDBWrapper, error) {
	address, err := registry.GetServiceAddress(RedisRemote)
	if err != nil {
		return nil, err
	}
	redisClientImpl := &object_store.RedisClientImpl{
		RawClient: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: "",
			DB:       0,
		}),
	}
	policyMap := object_store.NewRedisMap(
		redisClientImpl,
		"policydb:rules",
		policydb.GetPolicySerializer(),
		policydb.GetPolicyDeserializer(),
	)
	baseNameMap := object_store.NewRedisMap(
		redisClientImpl,
		"policydb:base_names",
		policydb.GetBaseNameSerializer(),
		policydb.GetBaseNameDeserializer(),
	)
	omniPresentRules := object_store.NewRedisMap(
		redisClientImpl,
		"policydb:omnipresent_rules",
		policydb.GetRuleMappingSerializer(),
		policydb.GetRuleMappingDeserializer(),
	)
	return &policyDBWrapper{
		redisClient:      redisClientImpl,
		policyMap:        policyMap,
		baseNameMap:      baseNameMap,
		omniPresentRules: omniPresentRules,
	}, nil
}

func contains(wordToFind string, words []string) bool {
	for _, w := range words {
		if w == wordToFind {
			return true
		}
	}
	return false
}

/**  ========== Directoryd Helpers ========== **/
// getDirectorydClient is a utility function to an RPC connection to a
// remote Directoryd service.
func getDirectorydClient() (*directorydClient, error) {
	var conn *grpc.ClientConn
	var err error
	conn, err = registryTestUtils.GetConnectionWithAuthority(DirectorydRemote)
	if err != nil {
		errMsg := fmt.Sprintf("Directoryd client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &directorydClient{
		orc8rprotos.NewGatewayDirectoryServiceClient(conn),
		conn,
	}, err
}

func updateDirectorydRecord(imsi, field, value string) error {
	cli, err := getDirectorydClient()
	if err == nil && cli != nil {
		updateReg := &orc8rprotos.UpdateRecordRequest{Id: imsi}
		updateReg.Fields = make(map[string]string, 1)
		updateReg.Fields[field] = value
		_, err = cli.UpdateRecord(context.Background(), updateReg)
	}
	return err
}
