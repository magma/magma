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

package decode

// SGsMessageType is the first byte in the SGs message for identifying the type of a message
type SGsMessageType byte

// Message type mapping
const (
	SGsAPPagingRequest            SGsMessageType = 0x01
	SGsAPPagingReject             SGsMessageType = 0x02
	SGsAPServiceRequest           SGsMessageType = 0x06
	SGsAPDownlinkUnitdata         SGsMessageType = 0x07
	SGsAPUplinkUnitdata           SGsMessageType = 0x08
	SGsAPLocationUpdateRequest    SGsMessageType = 0x09
	SGsAPLocationUpdateAccept     SGsMessageType = 0x0A
	SGsAPLocationUpdateReject     SGsMessageType = 0x0B
	SGsAPTMSIReallocationComplete SGsMessageType = 0x0C
	SGsAPAlertRequest             SGsMessageType = 0x0D
	SGsAPAlertAck                 SGsMessageType = 0x0E
	SGsAPAlertReject              SGsMessageType = 0x0F
	SGsAPUEActivityIndication     SGsMessageType = 0x10
	SGsAPEPSDetachIndication      SGsMessageType = 0x11
	SGsAPEPSDetachAck             SGsMessageType = 0x12
	SGsAPIMSIDetachIndication     SGsMessageType = 0x13
	SGsAPIMSIDetachAck            SGsMessageType = 0x14
	SGsAPResetIndication          SGsMessageType = 0x15
	SGsAPResetAck                 SGsMessageType = 0x16
	SGsAPServiceAbortRequest      SGsMessageType = 0x17
	SGsAPMMInformationRequest     SGsMessageType = 0x1A
	SGsAPReleaseRequest           SGsMessageType = 0x1B
	SGsAPStatus                   SGsMessageType = 0x1D
	SGsAPUEUnreachable            SGsMessageType = 0x1F
)

var MsgTypeNameByCode = map[SGsMessageType]string{
	SGsAPLocationUpdateAccept: "SGsAPLocationUpdateAccept",
	SGsAPLocationUpdateReject: "SGsAPLocationUpdateReject",
	SGsAPIMSIDetachAck:        "SGsAPIMSIDetachAck",
	SGsAPPagingRequest:        "SGsAPPagingRequest",
	SGsAPEPSDetachAck:         "SGsAPEPSDetachAck",
	SGsAPAlertRequest:         "SGsAPAlertRequest",
	SGsAPDownlinkUnitdata:     "SGsAPDownlinkUnitdata",
	SGsAPMMInformationRequest: "SGsAPMMInformationRequest",
	SGsAPReleaseRequest:       "SGsAPReleaseRequest",
	SGsAPServiceAbortRequest:  "SGsAPServiceAbortRequest",
	SGsAPStatus:               "SGsAPStatus",
	SGsAPResetAck:             "SGsAPResetAck",
	SGsAPResetIndication:      "SGsAPResetIndication",
}
