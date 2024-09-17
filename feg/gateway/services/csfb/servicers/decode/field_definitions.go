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

// InformationElementIdentifier is a byte for identifying an element in a message
type InformationElementIdentifier byte

// Information Element Identifier mapping
const (
	IEIIMSI                            InformationElementIdentifier = 0x01
	IEIVLRName                         InformationElementIdentifier = 0x02
	IEITMSI                            InformationElementIdentifier = 0x03
	IEILocationAreaIdentifier          InformationElementIdentifier = 0x04
	IEIChannelNeeded                   InformationElementIdentifier = 0x05
	IEIeMLPPPriority                   InformationElementIdentifier = 0x06
	IEITMSIStatus                      InformationElementIdentifier = 0x07
	IEISGsCause                        InformationElementIdentifier = 0x08
	IEIMMEName                         InformationElementIdentifier = 0x09
	IEIEPSLocationUpdateType           InformationElementIdentifier = 0x0A
	IEIGlobalCNId                      InformationElementIdentifier = 0x0B
	IEIMobileIdentity                  InformationElementIdentifier = 0x0E
	IEIRejectCause                     InformationElementIdentifier = 0x0F
	IEIIMSIDetachFromEPSServiceType    InformationElementIdentifier = 0x10
	IEIIMSIDetachFromNonEPSServiceType InformationElementIdentifier = 0x11
	IEIIMEISV                          InformationElementIdentifier = 0x15
	IEINASMessageContainer             InformationElementIdentifier = 0x16
	IEIMMInformation                   InformationElementIdentifier = 0x17
	IEIErroneousMessage                InformationElementIdentifier = 0x1B
	IEICLI                             InformationElementIdentifier = 0x1C
	IEILCSClientIdentity               InformationElementIdentifier = 0x1D
	IEILCSIndicator                    InformationElementIdentifier = 0x1E
	IEISSCode                          InformationElementIdentifier = 0x1F
	IEIServiceIndicator                InformationElementIdentifier = 0x20
	IEIUETimeZone                      InformationElementIdentifier = 0x21
	IEIMobileStationClassmark2         InformationElementIdentifier = 0x22
	IEITAI                             InformationElementIdentifier = 0x23
	IEIEUTRANCellGlobalIdentity        InformationElementIdentifier = 0x24
	IEIUEEMMMode                       InformationElementIdentifier = 0x25
)

// length in number of bytes
const (
	LengthIEI             = 1
	LengthLengthIndicator = 1
	LengthRejectCause     = 3
)

// length of information element in number of bytes
const (
	IELengthMessageType                     = 1
	IELengthTMSI                            = 6
	IELengthIMSIMin                         = 6
	IELengthIMSIMax                         = 10
	IELengthMMEName                         = 57
	IELengthVLRNameMin                      = 3
	IELengthMMInformationMin                = 3
	IELengthLocationAreaIdentifier          = 7
	IELengthSGsCause                        = 3
	IELengthCLIMin                          = 3
	IELengthCLIMax                          = 14
	IELengthGlobalCNId                      = 7
	IELengthLCSIndicator                    = 3
	IELengthLCSClientIdentityMin            = 3
	IELengthSSCode                          = 3
	IELengthChannelNeeded                   = 3
	IELengthEMLPPPriority                   = 3
	IELengthServiceIndicator                = 3
	IELengthNASMessageContainerMin          = 4
	IELengthNASMessageContainerMax          = 253
	IELengthErroneousMessageMin             = 3
	IELengthIMSIDetachFromEPSServiceType    = 3
	IELengthIMSIDetachFromNonEPSServiceType = 3
	IELengthEPSLocationUpdateType           = 3
	IELengthTMSIStatus                      = 3
	IELengthIMEISV                          = 10
	IELengthTAI                             = 7
	IELengthEUTRANCellGlobalIdentity        = 9
	IELengthUETimeZone                      = 3
	IELengthMobileStationClassmark2         = 5
	IELengthUEEMMMode                       = 3
	IELengthMobileIdentityMin               = 6
	IELengthMobileIdentityMax               = 10
	IELengthRejectCause                     = 3
)

// special encoding for fields
const (
	MobileIdentityTypeMask      = 0x07
	MobileIdentityIMSI          = 0x01
	MobileIdentityTMSI          = 0x04
	MobileIdentityTMSIFirstByte = 0xF4
)

var IEINamesByCode = map[InformationElementIdentifier]string{
	IEIIMSI:                            "IMSI",
	IEIVLRName:                         "VLRName",
	IEITMSI:                            "TMSI",
	IEILocationAreaIdentifier:          "LocationAreaIdentifier",
	IEIChannelNeeded:                   "ChannelNeeded",
	IEIeMLPPPriority:                   "eMLPPPriority",
	IEITMSIStatus:                      "TMSIStatus",
	IEISGsCause:                        "SGsCause",
	IEIMMEName:                         "MMEName",
	IEIEPSLocationUpdateType:           "LocationUpdateType",
	IEIGlobalCNId:                      "GlobalCNId",
	IEIMobileIdentity:                  "MobileIdentity",
	IEIRejectCause:                     "RejectCause",
	IEIIMSIDetachFromEPSServiceType:    "IMSIDetachFromEPSServiceType",
	IEIIMSIDetachFromNonEPSServiceType: "IMSIDetachFromNonEPSServiceType",
	IEIIMEISV:                          "IMEISV",
	IEINASMessageContainer:             "NASMessageContainer",
	IEIMMInformation:                   "MMInformation",
	IEIErroneousMessage:                "ErroneousMessage",
	IEICLI:                             "CLI",
	IEILCSClientIdentity:               "LCSClientIdentity",
	IEILCSIndicator:                    "IEILCSIndicator",
	IEISSCode:                          "SSCode",
	IEIServiceIndicator:                "ServiceIndicator",
	IEIUETimeZone:                      "UETimeZone",
	IEIMobileStationClassmark2:         "MobileStationClassmark2",
	IEITAI:                             "TAI",
	IEIEUTRANCellGlobalIdentity:        "EUTRANCellGlobalIdentity",
	IEIUEEMMMode:                       "UEEMMMode",
}
