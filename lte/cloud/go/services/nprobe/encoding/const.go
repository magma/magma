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

const (
	// IRI ASN.1 options as defined in ETSI TS 103 221-2.
	IRIBeginRecord    string = "tag:1"
	IRIEndRecord      string = "tag:2"
	IRIContinueRecord string = "tag:3"
	IRIReportRecord   string = "tag:4"
)

const (
	// Event types as defined in ETSI TS 133 108 R16 [B9].
	UnsupportedEvent                 asn1.Enumerated = 0
	EutranAttach                     asn1.Enumerated = 16
	EutranDetach                     asn1.Enumerated = 17
	BearerActivation                 asn1.Enumerated = 18
	StartInterceptWithActiveBearer   asn1.Enumerated = 19
	BearerModification               asn1.Enumerated = 20
	BearerDeactivation               asn1.Enumerated = 21
	UERequestedPDNConnectivity       asn1.Enumerated = 23
	UERequestedPDNDisconnection      asn1.Enumerated = 24
	LocationUpdate                   asn1.Enumerated = 25 // trackingAreaEpsLocationUpdate
	StartInterceptWithEutranAttached asn1.Enumerated = 41
)

const (
	// Bearer types as defined in ETSI TS 133 108 R16 [B9].
	DefaultBearer   asn1.Enumerated = 0x01
	DedicatedBearer asn1.Enumerated = 0x02

	// Detach, bearer activation, modification or deactivation initiator
	// as defined in ETSI TS 133 108 R16 [B9]
	InitiatorNotAvailable asn1.Enumerated = 0x00
	OriginatingTarget     asn1.Enumerated = 0x01 // UE requested
	TerminatingTarget     asn1.Enumerated = 0x02 // Network initiated

	// Winter/Summer Indication as defined in ETSI TS 133 108 R16 [B9]
	IndicationNotAvailable asn1.Enumerated = 0x00
	WinterTime             asn1.Enumerated = 0x01
	SummerTime             asn1.Enumerated = 0x02

	// IP Address type
	IPV4Type asn1.Enumerated = 0x00
	IPV6Type asn1.Enumerated = 0x01

	PartyQualifierTarget asn1.Enumerated = 0x03 // gPRSorEPS-Target
	RatTypeEutran        uint8           = 0x06 // Ran access type EUTRAN
)
