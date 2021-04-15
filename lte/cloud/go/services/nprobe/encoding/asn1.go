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

package encoding

import "encoding/asn1"

// ETSI TS 133 108 R16 [B9] ASN.1 OID
func GetOID() asn1.ObjectIdentifier {
	return []int{0, 4, 0, 2, 2, 4, 8, 15, 4}
}

// ASN.1 Payload struct as definied in ETSI TS 133 108 R16 [B9]
type EpsIRIContent IRIParameter

type IRIParameter struct {
	Hi2epsDomainID        asn1.ObjectIdentifier `asn1:"tag:0"`
	LawInterceptID        []byte                `asn1:"optional,tag:1"`
	TimeStamp             Timestamp             `asn1:"tag:3"`
	Initiator             asn1.Enumerated       `asn1:"tag:4"`
	PartyInformation      []PartyInformation    `asn1:"set,optional,tag:9"`
	EPSCorrelationNumber  []byte                `asn1:"optional,tag:18"`
	EPSEvent              asn1.Enumerated       `asn1:"optional,tag:20"`
	NetworkIdentifier     NetworkIdentifier     `asn1:"optional,tag:26"`
	EPSSpecificParameters EPSSpecificParameters `asn1:"optional,tag:36"`
}

type Timestamp struct {
	LocalTime LocalTimestamp `asn1:"tag:0"`
}

type LocalTimestamp struct {
	GeneralizedTime        []byte          `asn1:"tag:0"`
	WinterSummerIndication asn1.Enumerated `asn1:"tag:1"`
}

type PartyInformation struct {
	PartyQualified asn1.Enumerated `asn1:"optional,tag:0"`
	PartyIdentity  PartyIdentity   `asn1:"optional,tag:1"`
}

type PartyIdentity struct {
	IMEI   []byte `asn1:"optional,tag:1"`
	IMSI   []byte `asn1:"optional,tag:3"`
	MSISDN []byte `asn1:"optional,tag:6"`
}

type NetworkIdentifier struct {
	OperatorIdentifier       []byte                   `asn1:"tag:0"`
	NetworkElementIdentifier NetworkElementIdentifier `asn1:"optional,tag:1"`
}

type NetworkElementIdentifier struct {
	IPAddress IPAddress `asn1:"tag:5"`
}

type IPAddress struct {
	IPType  asn1.Enumerated `asn1:"tag:1"`
	IPValue IPValue         `asn1:"tag:2"`
}

type IPValue struct {
	IPBinaryAddress []byte `asn1:"tag:1"`
}

type EPSLocation struct {
	UserLocationInfo []byte `asn1:"optional,tag:1"`
}

type EPSSpecificParameters struct {
	PDNAddressAllocation   []byte          `asn1:"optional,tag:1"`
	APN                    []byte          `asn1:"optional,tag:2"`
	EPSBearerIdentity      []byte          `asn1:"optional,tag:5"`
	DetachType             []byte          `asn1:"optional,tag:6"`
	RATType                []byte          `asn1:"optional,tag:7"`
	FailedBearerActReason  []byte          `asn1:"optional,tag:8"`
	EPSBearerQoS           []byte          `asn1:"optional,tag:9"`
	BearerActivationType   asn1.Enumerated `asn1:"optional,tag:10"`
	ApnAmbr                []byte          `asn1:"optional,tag:11"`
	BearerDeactivationType asn1.Enumerated `asn1:"optional,tag:21"`
	EPSLocationOfTheTarget EPSLocation     `asn1:"optional,tag:23"`
}
