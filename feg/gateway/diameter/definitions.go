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

package diameter

const (
	// SuccessCode is the result code returned from a successful diameter call
	SuccessCode                = 2001
	LimitedSuccessCode         = 2002
	DiameterRatingFailed       = 5031
	DiameterCreditLimitReached = 4012
)

var diamCodeToNameMap = map[uint32]string{
	0:    "UNDEFINED",
	1001: "MULTI_ROUND_AUTH",
	2001: "SUCCESS",
	2002: "LIMITED_SUCCESS",
	3001: "COMMAND_UNSUPPORTED",
	3002: "UNABLE_TO_DELIVER",
	3003: "REALM_NOT_SERVED",
	3004: "TOO_BUSY",
	3005: "LOOP_DETECTED",
	3006: "REDIRECT_INDICATION",
	3007: "APPLICATION_UNSUPPORTED",
	3008: "INVALIDH_DR_BITS",
	3009: "INVALID_AVP_BITS",
	3010: "UNKNOWN_PEER",
	4001: "AUTHENTICATION_REJECTED",
	4002: "OUT_OF_SPACE",
	4003: "ELECTION_LOST",
	4012: "DIAMETER_CREDIT_LIMIT_REACHED",
	4181: "AUTHENTICATION_DATA_UNAVAILABLE",
	5001: "USER_UNKNOWN",
	5003: "IDENTITY_NOT_REGISTERED",
	5004: "ROAMING_NOT_ALLOWED",
	5005: "IDENTITY_ALREADY_REGISTERED",
	5420: "UNKNOWN_EPS_SUBSCRIPTION",
	5421: "RAT_NOT_ALLOWED",
	5422: "EQUIPMENT_UNKNOWN",
	5423: "UNKNOWN_SERVING_NODE",
	5450: "USER_NO_NON_3GPP_SUBSCRIPTION",
	5451: "USER_NO_APN_SUBSCRIPTION",
	5452: "RAT_TYPE_NOT_ALLOWED",
}

const (
	ServiceContextIDDefault = "32251@3gpp.org" // Packet-Switch service context
	ServiceIDDefault        = 0
)
