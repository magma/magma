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

package x2_interface

import "encoding/asn1"

const (
	// FixHeaderLength X2 header fix length
	FixHeaderLength uint32 = 40

	// X2 header values
	X2HeaderVersion       uint16 = 2
	X2HeaderPduType       uint16 = 1  // X2 PDU
	X2HeaderPayloadFormat uint16 = 14 // ETSI TS 133 108 [B.9] Defined Payload

	// X2 header values
	X2PayloadDirectionUnkown     uint16 = 1
	X2PayloadDirectionToTarger   uint16 = 2
	X2PayloadDirectionFromTarget uint16 = 3

	// X2 conditional attribute types
	X2AttrIDDomainID    uint16 = 5
	X2AttrIDNetworkFunc uint16 = 6
	X2AttrIDTimestamp   uint16 = 9
	X2AttrIDSeqNumber   uint16 = 8
	X2AttrIDTargetID    uint16 = 17
)

const (
	// IRI ASN.1 options
	IRIBeginRecord    string = "tag:1"
	IRIEndRecord      string = "tag:2"
	IRIContinueRecord string = "tag:3"
	IRIReportRecord   string = "tag:4"
)

const (
	// IRI event types
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
	// IRI Bearer types
	DefaultBearer   asn1.Enumerated = 0x01
	DedicatedBearer asn1.Enumerated = 0x02

	// IRI Initiator of detach, bearer activation, modification or deactivation
	InitiatorNotAvailable asn1.Enumerated = 0x00
	OriginatingTarget     asn1.Enumerated = 0x01 // UE requested
	TerminatingTarget     asn1.Enumerated = 0x02 // Network initiated

	// IRI Winter/Summer Indication
	IndicationNotAvailable asn1.Enumerated = 0x00
	WinterTime             asn1.Enumerated = 0x01
	SummerTime             asn1.Enumerated = 0x02

	// IP Address type
	IPV4Type asn1.Enumerated = 0x00
	IPV6Type asn1.Enumerated = 0x01

	PartyQualifierTarget asn1.Enumerated = 0x03 // gPRSorEPS-Target
	RatTypeEutran        int             = 0x06 // Ran access type EUTRAN
)
