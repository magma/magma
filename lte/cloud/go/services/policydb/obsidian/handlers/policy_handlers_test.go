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
	"fmt"
	"testing"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/policydb/obsidian/handlers"
	policyModels "magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

// Basic API workflow tests
func TestPolicyDBHandlersBasic(t *testing.T) {
	configurator_test_init.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	err := configurator.CreateNetwork(configurator.Network{ID: "n1", Type: lte.NetworkType}, serdes.Network)
	assert.NoError(t, err)

	listPolicies := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/rules", obsidian.GET).HandlerFunc
	createPolicy := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/rules", obsidian.POST).HandlerFunc
	getPolicy := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/rules/:rule_id", obsidian.GET).HandlerFunc
	updatePolicy := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/rules/:rule_id", obsidian.PUT).HandlerFunc
	deletePolicy := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/rules/:rule_id", obsidian.DELETE).HandlerFunc
	listNames := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/base_names", obsidian.GET).HandlerFunc
	createName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/base_names", obsidian.POST).HandlerFunc
	getName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/base_names/:base_name", obsidian.GET).HandlerFunc
	updateName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/base_names/:base_name", obsidian.PUT).HandlerFunc
	deleteName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/base_names/:base_name", obsidian.DELETE).HandlerFunc

	// Test empty response
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules?view=full",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listPolicies,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*policyModels.PolicyRule{}),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)
	tc.URL = "/magma/v1/networks/n1/policies/rules"
	tc.ExpectedResult = tests.JSONMarshaler([]string{})
	tests.RunUnitTest(t, e, tc)

	// Test add policy rule
	testRule := &policyModels.PolicyRule{
		ID: "PolicyRule1",
		FlowList: []*policyModels.FlowDescription{
			{
				Action: swag.String("PERMIT"),
				Match: &policyModels.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_ICMP"),
					IPDst: &policyModels.IPAddress{
						Version: policyModels.IPAddressVersionIPV4,
						Address: "42.42.42.42",
					},
					IPSrc: &policyModels.IPAddress{
						Version: policyModels.IPAddressVersionIPV4,
						Address: "192.168.0.1/24",
					},
					TCPDst: 2,
					TCPSrc: 1,
					UDPDst: 4,
					UDPSrc: 3,
				},
			},
		},
		Priority:     swag.Uint32(5),
		RatingGroup:  *swag.Uint32(2),
		TrackingType: "ONLY_OCS",
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/rules",
		Payload:        testRule,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createPolicy,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Check that policy rule was added
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules?view=full",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listPolicies,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*policyModels.PolicyRule{
			"PolicyRule1": testRule,
		}),
	}
	tests.RunUnitTest(t, e, tc)
	tc.URL = "/magma/v1/networks/n1/policies/rules"
	tc.ExpectedResult = tests.JSONMarshaler([]string{"PolicyRule1"})
	tests.RunUnitTest(t, e, tc)

	// Test Read Rule Using URL based ID
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules/PolicyRule1",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "PolicyRule1"},
		Handler:        getPolicy,
		ExpectedStatus: 200,
		ExpectedResult: testRule,
	}
	tests.RunUnitTest(t, e, tc)

	// Test Update Rule Using URL based ID
	testRule.FlowList = []*policyModels.FlowDescription{
		{Action: swag.String("PERMIT"), Match: &policyModels.FlowMatch{IPProto: swag.String("IPPROTO_ICMP"), Direction: swag.String("DOWNLINK")}},
	}
	testRule.Priority, testRule.RatingGroup, testRule.TrackingType = swag.Uint32(10), *swag.Uint32(3), "ONLY_OCS"
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/networks/n1/policies/rules/PolicyRule1",
		Payload:        testRule,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "PolicyRule1"},
		Handler:        updatePolicy,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Verify update results
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules/PolicyRule1",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "PolicyRule1"},
		Handler:        getPolicy,
		ExpectedStatus: 200,
		ExpectedResult: testRule,
	}
	tests.RunUnitTest(t, e, tc)

	// Test header enrichment targets
	testRule.HeaderEnrichmentTargets = []string{"http://example.com", "http://example.net"}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/networks/n1/policies/rules/PolicyRule1",
		Payload:        testRule,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "PolicyRule1"},
		Handler:        updatePolicy,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Verify update results
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules/PolicyRule1",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "PolicyRule1"},
		Handler:        getPolicy,
		ExpectedStatus: 200,
		ExpectedResult: testRule,
	}
	tests.RunUnitTest(t, e, tc)

	// Delete a rule
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/networks/n1/policies/rules/PolicyRule1",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "PolicyRule1"},
		Handler:        deletePolicy,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Confirm delete
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules?view=full",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listPolicies,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*policyModels.PolicyRule{}),
	}
	tests.RunUnitTest(t, e, tc)
	tc.URL = "/magma/v1/networks/n1/policies/rules"
	tc.ExpectedResult = tests.JSONMarshaler([]string{})
	tests.RunUnitTest(t, e, tc)

	// Test add invalid policy rule
	testRule_invalid_rule := &policyModels.PolicyRule{
		ID: "PolicyRule_invalid",
		FlowList: []*policyModels.FlowDescription{
			{
				Action: swag.String("PERMIT"),
				Match: &policyModels.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_ICMP"),
					IPDst: &policyModels.IPAddress{
						Version: policyModels.IPAddressVersionIPV4,
						Address: "42.42.42.42",
					},
					IPSrc: &policyModels.IPAddress{
						Version: policyModels.IPAddressVersionIPV4,
						Address: "192.168.0.1/24",
					},
					IPV4Dst: "42.42.42.42",
					IPV4Src: "192.168.0.1/24",
					TCPDst:  2,
				},
			},
		},
		Priority:     swag.Uint32(5),
		RatingGroup:  *swag.Uint32(2),
		TrackingType: "ONLY_OCS",
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/rules",
		Payload:        testRule_invalid_rule,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createPolicy,
		ExpectedStatus: 400,
		ExpectedError:  "Invalid Argument: Can't mix old ipv4_src/ipv4_dst type with the new ip_src/ip_dst",
	}
	tests.RunUnitTest(t, e, tc)

	// Test old ip(ipv4_src/ipvr_dst) is properly converted to new ip_src/ip_dst
	test_old_ip_policy := &policyModels.PolicyRule{
		ID: "test_old_ip_policy",
		FlowList: []*policyModels.FlowDescription{
			{
				Action: swag.String("PERMIT"),
				Match: &policyModels.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_ICMP"),
					IPV4Dst:   "42.42.42.42",
					IPV4Src:   "192.168.0.1/24",
					TCPDst:    2,
				},
			},
		},
		Priority:     swag.Uint32(5),
		RatingGroup:  *swag.Uint32(2),
		TrackingType: "ONLY_OCS",
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/rules",
		Payload:        test_old_ip_policy,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createPolicy,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	test_modified_policy := &policyModels.PolicyRule{
		ID: "test_old_ip_policy",
		FlowList: []*policyModels.FlowDescription{
			{
				Action: swag.String("PERMIT"),
				Match: &policyModels.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_ICMP"),
					IPDst: &policyModels.IPAddress{
						Version: policyModels.IPAddressVersionIPV4,
						Address: "42.42.42.42",
					},
					IPSrc: &policyModels.IPAddress{
						Version: policyModels.IPAddressVersionIPV4,
						Address: "192.168.0.1/24",
					},
					TCPDst: 2,
				},
			},
		},
		Priority:     swag.Uint32(5),
		RatingGroup:  *swag.Uint32(2),
		TrackingType: "ONLY_OCS",
	}

	// Test Read Rule Using URL based ID
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules/test_old_ip_policy",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "test_old_ip_policy"},
		Handler:        getPolicy,
		ExpectedStatus: 200,
		ExpectedResult: test_modified_policy,
	}
	tests.RunUnitTest(t, e, tc)

	// Test Multi Match Add Rule
	testRule = &policyModels.PolicyRule{
		ID: "Test_mult",
		FlowList: []*policyModels.FlowDescription{
			{
				Action: swag.String("DENY"),
				Match: &policyModels.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_TCP"),
					TCPDst:    2,
					TCPSrc:    1,
				},
			},
			{
				Action: swag.String("PERMIT"),
				Match: &policyModels.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_ICMP"),
					IPDst: &policyModels.IPAddress{
						Version: policyModels.IPAddressVersionIPV4,
						Address: "42.42.42.42",
					},
					IPSrc: &policyModels.IPAddress{
						Version: policyModels.IPAddressVersionIPV4,
						Address: "192.168.0.1/24",
					},
					TCPDst: 2,
					TCPSrc: 1,
					UDPDst: 4,
					UDPSrc: 3,
				},
			},
		},
		Priority:     swag.Uint32(5),
		RatingGroup:  *swag.Uint32(2),
		TrackingType: "ONLY_OCS",
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/rules",
		Payload:        testRule,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createPolicy,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Test Read Rule Using URL based ID
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules/Test_mult",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "Test_mult"},
		Handler:        getPolicy,
		ExpectedStatus: 200,
		ExpectedResult: testRule,
	}
	tests.RunUnitTest(t, e, tc)

	// Test adding a rule with QoS
	testRule = &policyModels.PolicyRule{
		ID:           "Test_qos",
		FlowList:     []*policyModels.FlowDescription{},
		Priority:     swag.Uint32(5),
		RatingGroup:  *swag.Uint32(2),
		TrackingType: "ONLY_OCS",
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/rules",
		Payload:        testRule,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createPolicy,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Check that QoS rule was added
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules/Test_qos",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "Test_qos"},
		Handler:        getPolicy,
		ExpectedStatus: 200,
		ExpectedResult: testRule,
	}
	tests.RunUnitTest(t, e, tc)

	// Test adding rule with redirect information
	testRule = &policyModels.PolicyRule{
		ID:           "Test_redirect",
		FlowList:     []*policyModels.FlowDescription{},
		Priority:     swag.Uint32(5),
		RatingGroup:  *swag.Uint32(2),
		TrackingType: "ONLY_OCS",
		Redirect: &policyModels.RedirectInformation{
			Support:       swag.String("ENABLED"),
			AddressType:   swag.String("URL"),
			ServerAddress: swag.String("127.0.0.1"),
		},
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/rules",
		Payload:        testRule,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createPolicy,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Check that rule with redirect was added
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules/Test_redirect",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "Test_redirect"},
		Handler:        getPolicy,
		ExpectedStatus: 200,
		ExpectedResult: testRule,
	}
	tests.RunUnitTest(t, e, tc)

	// Test add rule with app name match
	testRule = &policyModels.PolicyRule{
		ID: "test_app_policy",
		FlowList: []*policyModels.FlowDescription{
			{
				Action: swag.String("PERMIT"),
				Match: &policyModels.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_ICMP"),
					IPDst: &policyModels.IPAddress{
						Version: policyModels.IPAddressVersionIPV4,
						Address: "42.42.42.42",
					},
					IPSrc: &policyModels.IPAddress{
						Version: policyModels.IPAddressVersionIPV4,
						Address: "192.168.0.1/24",
					},
					TCPDst: 2,
					TCPSrc: 1,
					UDPDst: 4,
					UDPSrc: 3,
				},
			},
		},
		Priority:       swag.Uint32(5),
		RatingGroup:    *swag.Uint32(2),
		TrackingType:   "ONLY_OCS",
		AppName:        "INSTAGRAM",
		AppServiceType: "VIDEO",
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/rules",
		Payload:        testRule,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createPolicy,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Check that rule with app name was added
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules/test_app_policy",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "test_app_policy"},
		Handler:        getPolicy,
		ExpectedStatus: 200,
		ExpectedResult: testRule,
	}
	tests.RunUnitTest(t, e, tc)

	// Now run base name test cases using the rules created above

	// Test Listing All Base Names
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/base_names?view=full",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listNames,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*policyModels.BaseNameRecord{}),
	}
	tests.RunUnitTest(t, e, tc)
	tc.URL = "/magma/v1/networks/n1/policies/base_names"
	tc.ExpectedResult = tests.JSONMarshaler([]string{})
	tests.RunUnitTest(t, e, tc)

	// Test Add BaseName
	baseNameRecord := &policyModels.BaseNameRecord{
		Name:      "Test",
		RuleNames: policyModels.RuleNames{"Test_qos", "Test_redirect"},
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/base_names",
		Payload:        baseNameRecord,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createName,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Test Read BaseName Using URL based name
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/base_names/Test",
		Payload:        nil,
		ParamNames:     []string{"network_id", "base_name"},
		ParamValues:    []string{"n1", "Test"},
		Handler:        getName,
		ExpectedStatus: 200,
		ExpectedResult: &policyModels.BaseNameRecord{
			Name:      "Test",
			RuleNames: []string{"Test_qos", "Test_redirect"},
		},
	}
	tests.RunUnitTest(t, e, tc)

	// Test Update BaseName Using URL based name
	tc = tests.Test{
		Method: "PUT",
		URL:    "/magma/v1/networks/n1/policies/base_names/Test",
		Payload: &policyModels.BaseNameRecord{
			Name:      "Test",
			RuleNames: []string{"Test_qos"},
		},
		ParamNames:     []string{"network_id", "base_name"},
		ParamValues:    []string{"n1", "Test"},
		Handler:        updateName,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Verify update BaseName
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/base_names/Test",
		Payload:        nil,
		ParamNames:     []string{"network_id", "base_name"},
		ParamValues:    []string{"n1", "Test"},
		Handler:        getName,
		ExpectedStatus: 200,
		ExpectedResult: &policyModels.BaseNameRecord{
			Name:      "Test",
			RuleNames: []string{"Test_qos"},
		},
	}
	tests.RunUnitTest(t, e, tc)

	// Get all BaseNames
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/base_names?view=full",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listNames,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*policyModels.BaseNameRecord{
			"Test": {
				Name:      "Test",
				RuleNames: []string{"Test_qos"},
			},
		}),
	}
	tests.RunUnitTest(t, e, tc)
	tc.URL = "/magma/v1/networks/n1/policies/base_names"
	tc.ExpectedResult = tests.JSONMarshaler([]string{"Test"})
	tests.RunUnitTest(t, e, tc)

	// Delete a BaseName
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/networks/n1/policies/base_names/Test",
		Payload:        nil,
		ParamNames:     []string{"network_id", "base_name"},
		ParamValues:    []string{"n1", "Test"},
		Handler:        deleteName,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Confirm delete
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/base_names?view=full",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listNames,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*policyModels.BaseNameRecord{}),
	}
	tests.RunUnitTest(t, e, tc)
	tc.URL = "/magma/v1/networks/n1/policies/base_names"
	tc.ExpectedResult = tests.JSONMarshaler([]string{})
	tests.RunUnitTest(t, e, tc)
}

// Associate base names and policies to subscribers
func TestPolicyHandlersAssociations(t *testing.T) {
	configurator_test_init.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	err := configurator.CreateNetwork(configurator.Network{ID: "n1", Type: lte.NetworkType}, serdes.Network)
	assert.NoError(t, err)

	createPolicy := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/rules", obsidian.POST).HandlerFunc
	getPolicy := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/rules/:rule_id", obsidian.GET).HandlerFunc
	updatePolicy := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/rules/:rule_id", obsidian.PUT).HandlerFunc

	createName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/base_names", obsidian.POST).HandlerFunc
	getName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/base_names/:base_name", obsidian.GET).HandlerFunc
	updateName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/base_names/:base_name", obsidian.PUT).HandlerFunc

	// preseed 3 subscribers
	imsi1, imsi2, imsi3 := "IMSI1234567890", "IMSI0987654321", "IMSI1111111111"
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.SubscriberEntityType, Key: imsi1},
			{Type: lte.SubscriberEntityType, Key: imsi2},
			{Type: lte.SubscriberEntityType, Key: imsi3},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	// Create rule assigned to s1, s2
	expectedP1 := &policyModels.PolicyRule{
		AssignedSubscribers: []policyModels.SubscriberID{policyModels.SubscriberID(imsi2), policyModels.SubscriberID(imsi1)},
		ID:                  "p1",
		FlowList: []*policyModels.FlowDescription{
			{
				Action: swag.String("PERMIT"),
				Match: &policyModels.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_IP"),
				},
			},
		},
		Priority: swag.Uint32(1),
	}
	tc := tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/rules",
		Payload:        expectedP1,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createPolicy,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)
	validatePolicy(
		t, e, getPolicy,
		expectedP1,
		configurator.NetworkEntity{
			NetworkID:          "n1",
			Type:               lte.PolicyRuleEntityType,
			Key:                "p1",
			GraphID:            "2",
			ParentAssociations: []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: imsi2}, {Type: lte.SubscriberEntityType, Key: imsi1}},
			Version:            0,
		},
	)

	// Update rule to assign to s3
	expectedP1.AssignedSubscribers = []policyModels.SubscriberID{policyModels.SubscriberID(imsi3)}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/networks/n1/policies/rules/p1",
		Payload:        expectedP1,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "p1"},
		Handler:        updatePolicy,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	validatePolicy(
		t, e, getPolicy,
		expectedP1,
		configurator.NetworkEntity{
			NetworkID:          "n1",
			Type:               lte.PolicyRuleEntityType,
			Key:                "p1",
			GraphID:            "2",
			ParentAssociations: []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: imsi3}},
			Version:            1,
		},
	)

	// Create another policy p2 unbound to subs
	expectedP2 := &policyModels.PolicyRule{
		ID: "p2",
		FlowList: []*policyModels.FlowDescription{
			{
				Action: swag.String("PERMIT"),
				Match: &policyModels.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_IP"),
				},
			},
		},
		Priority: swag.Uint32(1),
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/rules",
		Payload:        expectedP2,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createPolicy,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Create base name bound to p1 and imsi1, imsi2
	expectedBN := &policyModels.BaseNameRecord{
		Name:                "b1",
		RuleNames:           policyModels.RuleNames{"p1"},
		AssignedSubscribers: []policyModels.SubscriberID{policyModels.SubscriberID(imsi2), policyModels.SubscriberID(imsi1)},
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/base_names",
		Payload:        expectedBN,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createName,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)
	validateBaseName(
		t, e, getName,
		expectedBN,
		configurator.NetworkEntity{
			NetworkID: "n1",
			Type:      lte.BaseNameEntityType,
			Key:       "b1",
			GraphID:   "10",
			Associations: storage.TKs{
				{Type: lte.PolicyRuleEntityType, Key: "p1"},
			},
			ParentAssociations: storage.TKs{
				{Type: lte.SubscriberEntityType, Key: imsi2},
				{Type: lte.SubscriberEntityType, Key: imsi1},
			},
		},
	)

	// Update base name to bind to p2 and s3
	expectedBN = &policyModels.BaseNameRecord{
		Name:                "b1",
		RuleNames:           policyModels.RuleNames{"p2"},
		AssignedSubscribers: []policyModels.SubscriberID{policyModels.SubscriberID(imsi3)},
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/networks/n1/policies/base_names/b1",
		Payload:        expectedBN,
		ParamNames:     []string{"network_id", "base_name"},
		ParamValues:    []string{"n1", "b1"},
		Handler:        updateName,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	validateBaseName(
		t, e, getName,
		expectedBN,
		configurator.NetworkEntity{
			NetworkID: "n1",
			Type:      lte.BaseNameEntityType,
			Key:       "b1",
			GraphID:   "10",
			Associations: storage.TKs{
				{Type: lte.PolicyRuleEntityType, Key: "p2"},
			},
			ParentAssociations: storage.TKs{
				{Type: lte.SubscriberEntityType, Key: imsi3},
			},
			Version: 1,
		},
	)

	// Update base name to bind to p1, p2, s1, s2, s3
	expectedBN = &policyModels.BaseNameRecord{
		Name:                "b1",
		RuleNames:           policyModels.RuleNames{"p1", "p2"},
		AssignedSubscribers: []policyModels.SubscriberID{policyModels.SubscriberID(imsi2), policyModels.SubscriberID(imsi3), policyModels.SubscriberID(imsi1)},
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/networks/n1/policies/base_names/b1",
		Payload:        expectedBN,
		ParamNames:     []string{"network_id", "base_name"},
		ParamValues:    []string{"n1", "b1"},
		Handler:        updateName,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	validateBaseName(
		t, e, getName,
		expectedBN,
		configurator.NetworkEntity{
			NetworkID: "n1",
			Type:      lte.BaseNameEntityType,
			Key:       "b1",
			GraphID:   "10",
			Associations: storage.TKs{
				{Type: lte.PolicyRuleEntityType, Key: "p1"},
				{Type: lte.PolicyRuleEntityType, Key: "p2"},
			},
			ParentAssociations: storage.TKs{
				{Type: lte.SubscriberEntityType, Key: imsi2},
				{Type: lte.SubscriberEntityType, Key: imsi3},
				{Type: lte.SubscriberEntityType, Key: imsi1},
			},
			Version: 2,
		},
	)
}

func TestQoSProfile(t *testing.T) {
	configurator_test_init.StartTestService(t)
	e := echo.New()

	policydbHandlers := handlers.GetHandlers()
	err := configurator.CreateNetwork(configurator.Network{ID: "n1", Type: lte.NetworkType}, serdes.Network)
	assert.NoError(t, err)

	getAllProfiles := tests.GetHandlerByPathAndMethod(t, policydbHandlers, "/magma/v1/lte/:network_id/policy_qos_profiles", obsidian.GET).HandlerFunc
	postProfile := tests.GetHandlerByPathAndMethod(t, policydbHandlers, "/magma/v1/lte/:network_id/policy_qos_profiles", obsidian.POST).HandlerFunc
	putProfile := tests.GetHandlerByPathAndMethod(t, policydbHandlers, "/magma/v1/lte/:network_id/policy_qos_profiles/:profile_id", obsidian.PUT).HandlerFunc
	getProfile := tests.GetHandlerByPathAndMethod(t, policydbHandlers, "/magma/v1/lte/:network_id/policy_qos_profiles/:profile_id", obsidian.GET).HandlerFunc
	deleteProfile := tests.GetHandlerByPathAndMethod(t, policydbHandlers, "/magma/v1/lte/:network_id/policy_qos_profiles/:profile_id", obsidian.DELETE).HandlerFunc

	// Get all profiles, initially empty
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n1/policy_qos_profiles",
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getAllProfiles,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*policyModels.PolicyQosProfile{}),
	}
	tests.RunUnitTest(t, e, tc)

	// Get nonexistent profile
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n1/policy_qos_profiles/profile0",
		ParamNames:     []string{"network_id", "profile_id"},
		ParamValues:    []string{"n1", "profile0"},
		Handler:        getProfile,
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)

	// Put nonexistent profile
	profileX := newTestQoSProfile()
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/lte/n1/policy_qos_profiles/profile0",
		Payload:        profileX,
		ParamNames:     []string{"network_id", "profile_id"},
		ParamValues:    []string{"n1", "profile0"},
		Handler:        putProfile,
		ExpectedStatus: 400,
		ExpectedError:  "profile does not exist",
	}
	tests.RunUnitTest(t, e, tc)

	// Delete nonexistent profile
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/lte/n1/policy_qos_profiles/profile0",
		ParamNames:     []string{"network_id", "profile_id"},
		ParamValues:    []string{"n1", "profile0"},
		Handler:        deleteProfile,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Post malformed profile
	profileXa := newTestQoSProfile()
	tc = tests.Test{
		Method:                 "POST",
		URL:                    "/magma/v1/lte/n1/policy_qos_profiles",
		Payload:                profileXa,
		MalformedPayload:       true,
		ParamNames:             []string{"network_id"},
		ParamValues:            []string{"n1"},
		Handler:                postProfile,
		ExpectedStatus:         400,
		ExpectedErrorSubstring: "Syntax error",
	}
	tests.RunUnitTest(t, e, tc)

	// Post invalid profile
	profileXb := newTestQoSProfile()
	profileXb.Arp.PriorityLevel = swag.Uint32(16) // invalid
	tc = tests.Test{
		Method:                 "POST",
		URL:                    "/magma/v1/lte/n1/policy_qos_profiles",
		Payload:                profileXb,
		ParamNames:             []string{"network_id"},
		ParamValues:            []string{"n1"},
		Handler:                postProfile,
		ExpectedStatus:         400,
		ExpectedErrorSubstring: "validation failure list",
	}
	tests.RunUnitTest(t, e, tc)

	// Post profile
	profile0 := newTestQoSProfile()
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte/n1/policy_qos_profiles",
		Payload:        profile0,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        postProfile,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Post existing profile
	tc = tests.Test{
		Method:                 "POST",
		URL:                    "/magma/v1/lte/n1/policy_qos_profiles",
		Payload:                profile0,
		ParamNames:             []string{"network_id"},
		ParamValues:            []string{"n1"},
		Handler:                postProfile,
		ExpectedStatus:         400,
		ExpectedErrorSubstring: "Bad Request",
	}
	tests.RunUnitTest(t, e, tc)

	// Get existing profile
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n1/policy_qos_profiles/profile0",
		ParamNames:     []string{"network_id", "profile_id"},
		ParamValues:    []string{"n1", "profile0"},
		Handler:        getProfile,
		ExpectedStatus: 200,
		ExpectedResult: profile0,
	}
	tests.RunUnitTest(t, e, tc)

	// Put existing profile
	profile0a := newTestQoSProfile()
	profile0a.ClassID = 5
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/lte/n1/policy_qos_profiles/profile0",
		Payload:        profile0a,
		ParamNames:     []string{"network_id", "profile_id"},
		ParamValues:    []string{"n1", "profile0"},
		Handler:        putProfile,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Put invalid ID
	profileX = newTestQoSProfile()
	profileX.ID = "xxx"
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/lte/n1/policy_qos_profiles/profile0",
		Payload:        profileX,
		ParamNames:     []string{"network_id", "profile_id"},
		ParamValues:    []string{"n1", "profile0"},
		Handler:        putProfile,
		ExpectedStatus: 400,
		ExpectedError:  "id field is read-only",
	}
	tests.RunUnitTest(t, e, tc)

	// Get existing profile, put succeeded
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n1/policy_qos_profiles/profile0",
		ParamNames:     []string{"network_id", "profile_id"},
		ParamValues:    []string{"n1", "profile0"},
		Handler:        getProfile,
		ExpectedStatus: 200,
		ExpectedResult: profile0a,
	}
	tests.RunUnitTest(t, e, tc)

	// Get all profiles, no longer empty
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n1/policy_qos_profiles",
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getAllProfiles,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*policyModels.PolicyQosProfile{"profile0": profile0a}),
	}
	tests.RunUnitTest(t, e, tc)

	// Delete profile
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/lte/n1/policy_qos_profiles/profile0",
		ParamNames:     []string{"network_id", "profile_id"},
		ParamValues:    []string{"n1", "profile0"},
		Handler:        deleteProfile,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Get existing profile, delete succeeded
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n1/policy_qos_profiles/profile0",
		ParamNames:     []string{"network_id", "profile_id"},
		ParamValues:    []string{"n1", "profile0"},
		Handler:        getProfile,
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)

	// Get all profiles, empty
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/lte/n1/policy_qos_profiles",
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getAllProfiles,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*policyModels.PolicyQosProfile{}),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestPolicyWithQoSProfile(t *testing.T) {
	configurator_test_init.StartTestService(t)
	e := echo.New()

	policydbHandlers := handlers.GetHandlers()
	err := configurator.CreateNetwork(configurator.Network{ID: "n1", Type: lte.NetworkType}, serdes.Network)
	assert.NoError(t, err)

	postProfile := tests.GetHandlerByPathAndMethod(t, policydbHandlers, "/magma/v1/lte/:network_id/policy_qos_profiles", obsidian.POST).HandlerFunc
	deleteProfile := tests.GetHandlerByPathAndMethod(t, policydbHandlers, "/magma/v1/lte/:network_id/policy_qos_profiles/:profile_id", obsidian.DELETE).HandlerFunc

	getAllRules := tests.GetHandlerByPathAndMethod(t, policydbHandlers, "/magma/v1/networks/:network_id/policies/rules", obsidian.GET).HandlerFunc
	postRule := tests.GetHandlerByPathAndMethod(t, policydbHandlers, "/magma/v1/networks/:network_id/policies/rules", obsidian.POST).HandlerFunc
	getRule := tests.GetHandlerByPathAndMethod(t, policydbHandlers, "/magma/v1/networks/:network_id/policies/rules/:rule_id", obsidian.GET).HandlerFunc
	putRule := tests.GetHandlerByPathAndMethod(t, policydbHandlers, "/magma/v1/networks/:network_id/policies/rules/:rule_id", obsidian.PUT).HandlerFunc

	// Post profile
	tc := tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/lte/n1/policy_qos_profiles",
		Payload:        newTestQoSProfile(),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        postProfile,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Post rule
	rule := newTestPolicy("rule0")
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/rules",
		Payload:        rule,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        postRule,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Get rule, no profile
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules/rule0",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "rule0"},
		Handler:        getRule,
		ExpectedStatus: 200,
		ExpectedResult: rule,
	}
	tests.RunUnitTest(t, e, tc)

	// Put rule, leave profile empty
	rule.Priority = swag.Uint32(14)
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/networks/n1/policies/rules/profile1",
		Payload:        rule,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "rule0"},
		Handler:        putRule,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Put rule, associate with profile
	rule.QosProfile = "profile0"
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/networks/n1/policies/rules/profile1",
		Payload:        rule,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "rule0"},
		Handler:        putRule,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Get rule, profile found
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules/rule0",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "rule0"},
		Handler:        getRule,
		ExpectedStatus: 200,
		ExpectedResult: rule,
	}
	tests.RunUnitTest(t, e, tc)

	// Get all rules, profile found
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules?view=full",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getAllRules,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*policyModels.PolicyRule{"rule0": rule}),
	}
	tests.RunUnitTest(t, e, tc)

	// Delete profile
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/lte/n1/policy_qos_profiles/profile0",
		ParamNames:     []string{"network_id", "profile_id"},
		ParamValues:    []string{"n1", "profile0"},
		Handler:        deleteProfile,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Get rule, profile is gone
	rule.QosProfile = ""
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules/rule0",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "rule0"},
		Handler:        getRule,
		ExpectedStatus: 200,
		ExpectedResult: rule,
	}
	tests.RunUnitTest(t, e, tc)
}

// config will be filled from the expected model
func validatePolicy(t *testing.T, e *echo.Echo, getRule echo.HandlerFunc, expectedModel *policyModels.PolicyRule, expectedEnt configurator.NetworkEntity) {
	expectedEnt.Config = getExpectedRuleConfig(expectedModel)

	actual, err := configurator.LoadEntity(
		"n1", lte.PolicyRuleEntityType, string(expectedModel.ID),
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEnt, actual)
	tc := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("/magma/v1/networks/n1/policies/rules/%s", expectedModel.ID),
		Payload:        expectedModel,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", string(expectedModel.ID)},
		Handler:        getRule,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)
}

func getExpectedRuleConfig(m *policyModels.PolicyRule) *policyModels.PolicyRuleConfig {
	return &policyModels.PolicyRuleConfig{
		FlowList:      m.FlowList,
		MonitoringKey: m.MonitoringKey,
		Priority:      m.Priority,
		RatingGroup:   m.RatingGroup,
		Redirect:      m.Redirect,
		TrackingType:  m.TrackingType,
	}
}

func validateBaseName(t *testing.T, e *echo.Echo, getName echo.HandlerFunc, expectedModel *policyModels.BaseNameRecord, expectedEnt configurator.NetworkEntity) {
	actual, err := configurator.LoadEntity(
		"n1", lte.BaseNameEntityType, string(expectedModel.Name),
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedEnt, actual)
	tc := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("/magma/v1/networks/n1/policies/base_names/%s", expectedModel.Name),
		Payload:        expectedModel,
		ParamNames:     []string{"network_id", "base_name"},
		ParamValues:    []string{"n1", string(expectedModel.Name)},
		Handler:        getName,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)
}

func newTestQoSProfile() *policyModels.PolicyQosProfile {
	profile := &policyModels.PolicyQosProfile{
		Arp: &policyModels.Arp{
			PreemptionCapability:    swag.Bool(true),
			PreemptionVulnerability: swag.Bool(false),
			PriorityLevel:           swag.Uint32(5),
		},
		ClassID: 3,
		Gbr: &policyModels.Gbr{
			Downlink: swag.Uint32(42),
			Uplink:   swag.Uint32(420),
		},
		ID:         "profile0",
		MaxReqBwDl: swag.Uint32(42),
		MaxReqBwUl: swag.Uint32(420),
	}
	return profile
}

func newTestPolicy(id string) *policyModels.PolicyRule {
	policy := &policyModels.PolicyRule{
		ID: policyModels.PolicyID(id),
		FlowList: []*policyModels.FlowDescription{
			{
				Action: swag.String("PERMIT"),
				Match: &policyModels.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_IP"),
				},
			},
		},
		Priority: swag.Uint32(1),
	}
	return policy
}
