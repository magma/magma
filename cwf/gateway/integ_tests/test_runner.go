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
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/magma/milenage"
	"github.com/stretchr/testify/assert"

	"fbc/lib/go/radius"
	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/cwf/gateway/registry"
	"magma/cwf/gateway/services/uesim"
	fegprotos "magma/feg/cloud/go/protos"
	lteprotos "magma/lte/cloud/go/protos"
)

// todo make Op configurable, or export it in the UESimServer.
const (
	Op               = "\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"
	Secret           = "123456"
	MockHSSRemote    = "HSS_REMOTE"
	MockPCRFRemote   = "PCRF_REMOTE"
	MockOCSRemote    = "OCS_REMOTE"
	MockPCRFRemote2  = "PCRF_REMOTE2"
	MockOCSRemote2   = "OCS_REMOTE2"
	PipelinedRemote  = "pipelined.local"
	DirectorydRemote = "DIRECTORYD"
	RedisRemote      = "REDIS"
	CwagIP           = "192.168.70.101"
	TrafficCltIP     = "192.168.128.2"
	IPDRControllerIP = "192.168.40.11"
	OCSPort          = 9201
	PCRFPort         = 9202
	OCSPort2         = 9205
	PCRFPort2        = 9206
	HSSPort          = 9204
	PipelinedPort    = 8443
	RedisPort        = 6380
	DirectorydPort   = 8443

	// If updating these, also update the ipfix exported hex values
	defaultMSISDN          = "5100001234"
	defaultCalledStationID = "98-DE-D0-84-B5-47:CWF-TP-LINK_B547_5G"

	ipfixMSISDN        = "0x35313030303031323334000000000000"
	ipfixApnMacAddress = "0x98ded084b547"
	ipfixApnName       = "0x4357462d54502d4c494e4b5f423534375f35470000000000"

	KiloBytes                = 1024
	MegaBytes                = 1024 * KiloBytes
	Buffer                   = 100 * KiloBytes
	RevalidationTimeoutEvent = 17

	ReAuthMaxUsageBytes   = 5 * MegaBytes
	ReAuthMaxUsageTimeSec = 1000 // in second
	ReAuthValidityTime    = 60   // in second

	GyMaxUsageBytes = 5 * MegaBytes
	GyMaxUsageTime  = 1000 // in second
	GyValidityTime  = 60   // in second
)

// TestRunner helps setting up all associated services
type TestRunner struct {
	t           *testing.T
	imsis       map[string]bool
	activePCRFs []string
	activeOCSs  []string
	startTime   time.Time
}

// imsi -> ruleID -> record
type RecordByIMSI map[string]map[string]*lteprotos.RuleRecord

// NewTestRunner initializes a new TestRunner by making a UESim client and
// and setting the next IMSI.
func NewTestRunner(t *testing.T) *TestRunner {
	startTime := time.Now()
	fmt.Println("************* TestRunner setup")

	fmt.Printf("Adding Mock HSS service at %s:%d\n", CwagIP, HSSPort)
	registry.AddService(MockHSSRemote, CwagIP, HSSPort)
	fmt.Printf("Adding Mock PCRF service at %s:%d\n", CwagIP, PCRFPort)
	registry.AddService(MockPCRFRemote, CwagIP, PCRFPort)
	fmt.Printf("Adding Mock OCS service at %s:%d\n", CwagIP, OCSPort)
	registry.AddService(MockOCSRemote, CwagIP, OCSPort)
	fmt.Printf("Adding Pipelined service at %s:%d\n", CwagIP, PipelinedPort)
	registry.AddService(PipelinedRemote, CwagIP, PipelinedPort)
	fmt.Printf("Adding Redis service at %s:%d\n", CwagIP, RedisPort)
	registry.AddService(RedisRemote, CwagIP, RedisPort)
	fmt.Printf("Adding Directoryd service at %s:%d\n", CwagIP, DirectorydPort)
	registry.AddService(DirectorydRemote, CwagIP, DirectorydPort)

	testRunner := &TestRunner{t: t,
		activePCRFs: []string{MockPCRFRemote},
		activeOCSs:  []string{MockOCSRemote},
		startTime:   startTime,
	}
	testRunner.imsis = make(map[string]bool)
	return testRunner
}

// NewTestRunnerWithTwoPCRFandOCS does the same as NewTestRunner but it inclides 2 PCRF and 2 OCS
// Used in scenarios that run 2 PCRFs and 2 OCSs
func NewTestRunnerWithTwoPCRFandOCS(t *testing.T) *TestRunner {
	tr := NewTestRunner(t)

	fmt.Printf("Adding Mock PCRF #2 service at %s:%d\n", CwagIP, PCRFPort2)
	registry.AddService(MockPCRFRemote2, CwagIP, PCRFPort2)
	fmt.Printf("Adding Mock OCS #2 service at %s:%d\n", CwagIP, OCSPort2)
	registry.AddService(MockOCSRemote2, CwagIP, OCSPort2)

	// add the extra two servers for clean up
	tr.activePCRFs = append(tr.activePCRFs, MockPCRFRemote2)
	tr.activeOCSs = append(tr.activeOCSs, MockOCSRemote2)

	return tr
}

// ConfigUEs creates and adds the specified number of UEs and Subscribers
// to the UE Simulator and the HSS.
func (tr *TestRunner) ConfigUEs(numUEs int) ([]*cwfprotos.UEConfig, error) {
	IMSIs := make([]string, 0, numUEs)
	for i := 0; i < numUEs; i++ {
		imsi := ""
		for {
			imsi = getRandomIMSI()
			_, present := tr.imsis[imsi]
			if !present {
				break
			}
		}
		IMSIs = append(IMSIs, imsi)
	}
	return tr.ConfigUEsPerInstance(IMSIs, MockPCRFRemote, MockOCSRemote)
}

// ConfigUEsPerInstance same as ConfigUEs but per specific PCRF and OCS instance
func (tr *TestRunner) ConfigUEsPerInstance(IMSIs []string, pcrfInstance, ocsInstance string) ([]*cwfprotos.UEConfig, error) {
	fmt.Printf("************* Configuring %d UE(s), PCRF instance: %s\n", len(IMSIs), pcrfInstance)
	ues := make([]*cwfprotos.UEConfig, 0)
	for _, imsi := range IMSIs {
		// If IMSIs were generated properly they should never give an error here
		if _, present := tr.imsis[imsi]; present {
			return nil, fmt.Errorf("IMSI %s already exist in database, use generateRandomIMSIS(num, tr.imsis) to create unique list", imsi)
		}
		key, opc, err := getRandKeyOpcFromOp([]byte(Op))
		if err != nil {
			return nil, err
		}
		seq := getRandSeq()

		ue := makeUE(imsi, key, opc, seq)
		sub := makeSubscriber(imsi, key, opc, seq+1)

		err = uesim.AddUE(ue)
		if err != nil {
			return nil, fmt.Errorf("Error adding UE to UESimServer: %w", err)
		}
		err = addSubscriberToHSS(sub)
		if err != nil {
			return nil, fmt.Errorf("Error adding Subscriber to HSS: %w", err)
		}
		err = addSubscriberToPCRFPerInstance(pcrfInstance, sub.GetSid())
		if err != nil {
			return nil, fmt.Errorf("Error adding Subscriber to PCRF: %w", err)
		}
		err = addSubscriberToOCSPerInstance(ocsInstance, sub.GetSid())
		if err != nil {
			return nil, fmt.Errorf("Error adding Subscriber to OCS: %w", err)
		}

		ues = append(ues, ue)
		fmt.Printf("Added UE to Simulator, %s, %s, and %s:\n"+
			"\tIMSI: %s\tKey: %x\tOpc: %x\tSeq: %d\n", MockHSSRemote, pcrfInstance, ocsInstance, imsi, key, opc, seq)
		tr.imsis[imsi] = true
	}
	fmt.Println("Successfully configured UE(s)")
	return ues, nil
}

// Authenticate simulates an authentication between the UE and the HSS with the specified
// IMSI and CalledStationID, and returns the resulting Radius packet.
func (tr *TestRunner) Authenticate(imsi, calledStationID string) (*radius.Packet, error) {
	fmt.Printf("************* Authenticating UE with IMSI: %s\n", imsi)
	res, err := uesim.Authenticate(&cwfprotos.AuthenticateRequest{Imsi: imsi, CalledStationID: calledStationID})
	if err != nil {
		fmt.Println(err)
		return &radius.Packet{}, err
	}
	encoded := res.GetRadiusPacket()
	radiusP, err := radius.Parse(encoded, []byte(Secret))
	if err != nil {
		err = fmt.Errorf("Error while parsing encoded Radius packet: %w", err)
		fmt.Println(err)
		return &radius.Packet{}, err
	}
	fmt.Println("Finished Authenticating UE")
	return radiusP, nil
}

// Authenticate simulates an authentication between the UE and the HSS with the specified
// IMSI and CalledStationID, and returns the resulting Radius packet.
func (tr *TestRunner) Disconnect(imsi, calledStationID string) (*radius.Packet, error) {
	fmt.Printf("************* Sending a disconnect request UE with IMSI: %s\n", imsi)
	res, err := uesim.Disconnect(&cwfprotos.DisconnectRequest{Imsi: imsi, CalledStationID: calledStationID})
	if err != nil {
		return &radius.Packet{}, err
	}
	encoded := res.GetRadiusPacket()
	radiusP, err := radius.Parse(encoded, []byte(Secret))
	if err != nil {
		err = fmt.Errorf("Error while parsing encoded Radius packet: %w", err)
		fmt.Println(err)
		return &radius.Packet{}, err
	}
	fmt.Println("Finished Disconnecting UE")
	return radiusP, nil
}

// GenULTraffic simulates the UE sending traffic through the CWAG to the Internet
// by running an iperf3 client on the UE simulator and an iperf3 server on the
// Magma traffic server.
func (tr *TestRunner) GenULTraffic(req *cwfprotos.GenTrafficRequest) (*cwfprotos.GenTrafficResponse, error) {
	fmt.Printf("************* Generating Traffic for UE with Req: %v\n", req)
	res, err := uesim.GenTraffic(req)
	fmt.Printf("  ==========> Total Sent: %d bytes\n", res.GetEndOutput().GetSumSent().GetBytes())
	return res, err
}

// WARNING this function only works for ammounts smaller than 1 to 2 MB
// GenULTrafficBasedOnPolicyUsage uses GenULTraffic to send small chucks of data until specific rule has
// received enough quota. To avoid going over for too much quota, this function will try to
// adjust the amount sent in every iteration.
// Function will return nil error if we reach totalVolume (or more)
// Function will return error if it hasn't completed by waitFor time.
// - Policy doesn't exist before start sending data
// - If we spend more than wait for time, and we haven't reached totalVolume
//
// Arguments
// - req: request to pas to UEsim
// - ruleID: name of the rule to monitor
// - totalVolume: total used by that rule. Note that if hte rule was used before, you have to add it's previous usage
//	 So if the rule already used 1Mb and you want to send 1Mb more, you will have to use 2M as min
// - waitFor time out the UE will be sending data
func (tr *TestRunner) GenULTrafficBasedOnPolicyUsage(req *cwfprotos.GenTrafficRequest,
	ruleID string, totalVolume uint64, waitFor time.Duration) (*cwfprotos.GenTrafficResponse, error) {
	fmt.Printf("************* Checking rule %s exists before generating traffic for UE\n", ruleID)
	if !assert.Eventually(tr.t, tr.WaitForEnforcementStatsForRule(req.Imsi, ruleID),
		10*time.Second, 1*time.Second) {
		return nil, fmt.Errorf("GenULTrafficBasedOnPolicyUsage can not send traffic. Rule %s not installed", ruleID)
	}
	// Initial iteration will just send few bytes
	req.Volume = &wrappers.StringValue{Value: "100K"}
	req.Bitrate = &wrappers.StringValue{Value: "4M"}
	req.DisableServerReachabilityCheck = true
	fmt.Printf("************* Generating Traffic for UE in chuncks to fullfil request\n")
	var res *cwfprotos.GenTrafficResponse
	var err error
	var waitNeeded = false
	for start := time.Now(); time.Since(start) < waitFor; {
		startGenTrafficTime := time.Now()
		res, err = uesim.GenTrafficWithReatempts(req)
		if err != nil {
			return res, fmt.Errorf("GenULTrafficBasedOnPolicyUsage failed during GenTraffic: %s", err)
		}
		completeCycle := time.Now().Sub(startGenTrafficTime) % time.Second
		if waitNeeded {
			// this wait makes sure the time passes is exact in seconds. This way
			// we make sure policies had time to sync
			time.Sleep(completeCycle + 200*time.Millisecond)
			waitNeeded = false
		}
		time.Sleep(2200 * time.Millisecond)
		record, metEnforcerCondition := tr.WaitForEnforcementStatsForRuleGreaterThanOrDoesNotExistFunc(req.Imsi, ruleID, totalVolume)
		if err != nil || metEnforcerCondition {
			fmt.Printf("Done generating traffic\n")
			return res, err
		}
		// adjust volume
		// we will take around 95% of what is left in every iteration, but not bigger than
		// 5MB to make sure bandwidth th is under control. We will not apply a factor if remaining
		// is small enough (100k)
		remaining := totalVolume - record.BytesTx
		newVolume := remaining
		req.Bitrate = &wrappers.StringValue{Value: "10M"}
		if remaining > 100*KiloBytes {
			newVolume = uint64(float64(remaining) * 0.95)
			if newVolume > 5*MegaBytes {
				newVolume = 5 * MegaBytes
			}
			// only add wait for the bigger chunks
			waitNeeded = true
		}
		if newVolume < 1000 {
			newVolume = 1000
		}
		newVolumeStr := fmt.Sprintf("%dK", newVolume/1000)

		req.Volume = &wrappers.StringValue{Value: newVolumeStr}
		fmt.Printf("- not enough traffic genereted, sending %dKB more. Will be around %d%% of volume requested\n",
			newVolume/1000,
			100*(record.BytesTx+newVolume)/totalVolume,
		)
	}
	return res, fmt.Errorf("error: not enough traffic generated to fullfil GenULTrafficBasedOnPolicyUsage requieriment: %s", err)
}

// CleanUp Remove subscribers, rules, flows, and monitors to clean up the state for
// consecutive test runs
func (tr *TestRunner) CleanUp() error {
	for imsi := range tr.imsis {
		err := deleteSubscribersFromHSS(imsi)
		if err != nil {
			return err
		}
	}
	for _, instance := range tr.activePCRFs {
		err := clearSubscribersFromPCRFPerInstance(instance)
		if err != nil {
			return err
		}
	}
	for _, instance := range tr.activeOCSs {
		err := clearSubscribersFromOCSPerInstance(instance)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetPolicyUsage is a wrapper around pipelined's GetPolicyUsage and returns
// the policy usage keyed by subscriber ID
func (tr *TestRunner) GetPolicyUsage() (RecordByIMSI, error) {
	recordsBySubID := RecordByIMSI{}
	table, err := getPolicyUsage()
	if err != nil {
		return recordsBySubID, err
	}
	for _, record := range table.Records {
		fmt.Printf("\tRecord %v\n", record)
		_, exists := recordsBySubID[record.Sid]
		if !exists {
			recordsBySubID[record.Sid] = map[string]*lteprotos.RuleRecord{}
		}
		recordsBySubID[record.Sid][record.RuleId] = record
	}
	return recordsBySubID, nil
}

func (tr *TestRunner) WaitForEnforcementStatsToSync() {
	// TODO load this value from pipelined.yml
	enforcementPollPeriod := 1 * time.Second
	time.Sleep(4 * enforcementPollPeriod)
}

func (tr *TestRunner) WaitForPoliciesToSync() {
	// TODO load this value from sessiond.yml (rule_update_interval_sec)
	ruleUpdatePeriod := 1 * time.Second
	time.Sleep(4 * ruleUpdatePeriod)
}

// WaitForEnforcementStatsForRule Wait until the ruleIDs show up for the IMSI
func (tr *TestRunner) WaitForEnforcementStatsForRule(imsi string, ruleIDs ...string) func() bool {
	return func() bool {
		fmt.Printf("\tWaiting until %s, %v shows up in enforcement stats...\n", imsi, ruleIDs)
		records, err := tr.GetPolicyUsage()
		if err != nil {
			return false
		}
		if records[prependIMSIPrefix(imsi)] == nil {
			return false
		}
		for _, ruleID := range ruleIDs {
			if records[prependIMSIPrefix(imsi)][ruleID] == nil {
				return false
			}
		}
		fmt.Printf("\t%s, %v are now in enforcement stats!\n", imsi, ruleIDs)
		return true
	}
}

// WaitForNoEnforcementStatsForRule Wait until the ruleIDs disappear for the IMSI
func (tr *TestRunner) WaitForNoEnforcementStatsForRule(imsi string, ruleIDs ...string) func() bool {
	return func() bool {
		fmt.Printf("\tWaiting until %s, %v disappear from enforcement stats...\n", imsi, ruleIDs)
		records, err := tr.GetPolicyUsage()
		if err != nil {
			return false
		}
		if records[prependIMSIPrefix(imsi)] == nil {
			fmt.Printf("%s are no longer in enforcement stats!\n", imsi)
			return true
		}
		for _, ruleID := range ruleIDs {
			if records[prependIMSIPrefix(imsi)][ruleID] != nil {
				return false
			}
		}
		fmt.Printf("%s, %v are no longer in enforcement stats!\n", imsi, ruleIDs)
		return true
	}
}

func (tr *TestRunner) WaitForEnforcementStatsForRuleGreaterThan(imsi, ruleID string, min uint64) func() bool {
	// Todo figure out the best way to figure out when RAR is processed
	return func() bool {
		fmt.Printf("\tWaiting until %s, %s has more than %d bytes in enforcement stats...\n", imsi, ruleID, min)
		records, err := tr.GetPolicyUsage()
		imsi = prependIMSIPrefix(imsi)
		if err != nil {
			return false
		}
		if records[imsi] == nil {
			return false
		}
		record := records[imsi][ruleID]
		if record == nil {
			return false
		}
		txBytes := record.BytesTx
		if record.BytesTx <= min {
			return false
		}
		fmt.Printf("\t\u2713 %s, %s now passed %d > %d in enforcement stats!\n", imsi, ruleID, txBytes, min)
		return true
	}
}

// WaitForEnforcementStatsForRuleGreaterThanOrDoesNotExist returns true if we have sent more data > min, if the
// session doesn't exist, or if rule doest exist
func (tr *TestRunner) WaitForEnforcementStatsForRuleGreaterThanOrDoesNotExist(imsi, ruleID string, min uint64) func() bool {
	return func() bool {
		fmt.Printf("\tWaiting until %s, %s has more than %d bytes in enforcement stats or rule does not exist ...\n", imsi, ruleID, min)
		records, err := tr.GetPolicyUsage()
		if err != nil {
			return false
		}
		imsi = prependIMSIPrefix(imsi)
		if records[imsi] == nil {
			// Session is gone
			fmt.Printf("\tSession for %s, does not exist...\n", imsi)
			return true
		}
		record := records[imsi][ruleID]
		if record == nil {
			// Session is gone
			fmt.Printf("\tRule %s for %s, does not exist...\n", ruleID, imsi)
			return true
		}
		txBytes := record.BytesTx
		if record.BytesTx <= min {
			return false
		}
		fmt.Printf("\t\u2713 %s, %s now passed %d > %d in enforcement stats!\n", imsi, ruleID, txBytes, min)
		return true
	}
}

func (tr *TestRunner) WaitForEnforcementStatsForRuleGreaterThanOrDoesNotExistFunc(imsi, ruleID string, min uint64) (*lteprotos.RuleRecord, bool) {
	fmt.Printf("\tWaiting until %s, %s has more than %d bytes in enforcement stats or rule does not exist ...\n", imsi, ruleID, min)
	records, err := tr.GetPolicyUsage()
	if err != nil {
		return nil, false
	}
	imsi = prependIMSIPrefix(imsi)
	if records[imsi] == nil {
		// Session is gone
		fmt.Printf("\tSession for %s, does not exist...\n", imsi)
		return nil, true
	}
	record := records[imsi][ruleID]
	if record == nil {
		// Session is gone
		fmt.Printf("\tRule %s for %s, does not exist...\n", ruleID, imsi)
		return nil, true
	}
	txBytes := record.BytesTx
	if record.BytesTx < min {
		return record, false
	}
	fmt.Printf("\t\u2713 %s, %s now passed %d > %d in enforcement stats!(%d%%)\n",
		imsi, ruleID, txBytes, min, 100*txBytes/min)
	return record, true
}

// WaitForPolicyReAuthToProcess returns a method which checks for reauth answer and
// if it has sessionID which contains the IMSI
func (tr *TestRunner) WaitForPolicyReAuthToProcess(raa *fegprotos.PolicyReAuthAnswer, imsi string) func() bool {
	// Todo figure out the best way to figure out when RAR is processed
	return func() bool {
		if raa != nil && strings.Contains(raa.SessionId, "IMSI"+imsi) {
			return true
		}
		return false
	}
}

// WaitForChargingReAuthToProcess returns a method which checks for reauth answer and
// if it has sessionID which contains the IMSI
func (tr *TestRunner) WaitForChargingReAuthToProcess(raa *fegprotos.ChargingReAuthAnswer, imsi string) func() bool {
	// Todo figure out the best way to figure out when RAR is processed
	return func() bool {
		if raa != nil && strings.Contains(raa.SessionId, "IMSI"+imsi) {
			return true
		}
		return false
	}
}

func (tr *TestRunner) PrintElapsedTime() {
	now := time.Now()
	fmt.Printf("Elapsed Time: %s\n", now.Sub(tr.startTime))
}

// generateRandomIMSIS creates a slice of unique Random IMSIs taking into consideration a previous list with IMSIS
func generateRandomIMSIS(numIMSIs int, preExistingIMSIS map[string]interface{}) []string {
	set := make(map[string]bool)
	IMSIs := make([]string, 0, numIMSIs)
	for i := 0; i < numIMSIs; i++ {
		imsi := ""
		for {
			imsi = getRandomIMSI()
			// Check if IMSI is in the preexisting list of IMSI or in the current generated list
			presentPreExistingIMSIs := false
			if preExistingIMSIS != nil {
				_, presentPreExistingIMSIs = preExistingIMSIS[imsi]
			}
			_, present := set[imsi]
			if !present && !presentPreExistingIMSIs {
				break
			}
		}
		set[imsi] = true
		IMSIs = append(IMSIs, imsi)
	}
	return IMSIs
}

// getRandomIMSI makes a random 15-digit IMSI that is not added to the UESim or HSS.
func getRandomIMSI() string {
	imsi := ""
	for len(imsi) < 15 {
		imsi += strconv.Itoa(rand.Intn(10))
	}
	return imsi
}

// RandKeyOpc makes a random 16-byte key and calculates the Opc based off the Op.
func getRandKeyOpcFromOp(op []byte) (key, opc []byte, err error) {
	key = make([]byte, 16)
	rand.Read(key)

	tempOpc, err := milenage.GenerateOpc(key, op)
	if err != nil {
		return nil, nil, err
	}
	opc = tempOpc[:]
	return
}

// getRandSeq makes a random 43-bit Seq.
func getRandSeq() uint64 {
	return rand.Uint64() >> 21
}

// makeUE creates a new UE using the given values.
func makeUE(imsi string, key []byte, opc []byte, seq uint64) *cwfprotos.UEConfig {
	return &cwfprotos.UEConfig{
		Imsi:    imsi,
		AuthKey: key,
		AuthOpc: opc,
		Seq:     seq,
	}
}

func prependIMSIPrefix(imsi string) string {
	if strings.HasPrefix(imsi, "IMSI") {
		return imsi
	} else {
		return "IMSI" + imsi
	}
}

// MakeSubcriber creates a new Subscriber using the given values.
func makeSubscriber(imsi string, key []byte, opc []byte, seq uint64) *lteprotos.SubscriberData {
	return &lteprotos.SubscriberData{
		Sid: &lteprotos.SubscriberID{
			Id:   imsi,
			Type: 1,
		},
		Lte: &lteprotos.LTESubscription{
			State:    1,
			AuthAlgo: 0,
			AuthKey:  key,
			AuthOpc:  opc,
		},
		State: &lteprotos.SubscriberState{
			LteAuthNextSeq: seq,
		},
		Non_3Gpp: &lteprotos.Non3GPPUserProfile{
			Msisdn:              defaultMSISDN,
			Non_3GppIpAccess:    lteprotos.Non3GPPUserProfile_NON_3GPP_SUBSCRIPTION_ALLOWED,
			Non_3GppIpAccessApn: lteprotos.Non3GPPUserProfile_NON_3GPP_APNS_ENABLE,
			ApnConfig:           []*lteprotos.APNConfiguration{{}},
		},
	}
}

// Get the Pipelined encoded version of IMSI (set in metadata register)
func getEncodedIMSI(imsiStr string) (string, error) {
	imsi, err := strconv.Atoi(imsiStr)
	if err != nil {
		return "", err
	}

	prefixLen := len(imsiStr) - len(strings.TrimLeft(imsiStr, "0"))
	compacted := (imsi << 2) | (prefixLen & 0x3)
	return fmt.Sprintf("0x%016x", compacted<<1|0x1), nil
}
