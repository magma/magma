/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package packet

// EAPType as defined in RFC3748 section 5
// For extended list of protocols supported by EAP, see
// https://www.vocal.com/secure-communication/eap-types/
type EAPType int

// EAPType values
const (
	// Used internally to indicate 'type is not used' (not EAP-Request or -Response)
	EAPTypeNONE EAPType = 0

	// By spec
	EAPTypeIDENTITY     EAPType = 1
	EAPTypeNOTIFICATION EAPType = 2
	EAPTypeNAK          EAPType = 3
	EAPTypeMD5CHALLENGE EAPType = 4
	EAPTypeOTP          EAPType = 5
	EAPTypeGENTOKENCARD EAPType = 6
	EAPTypeCISCOLEAP    EAPType = 17
	EAPTypeSIM          EAPType = 18
	EAPTypeAKA          EAPType = 23
	EAPTypeEAPMSCHAPV2  EAPType = 26
	EAPTypeEXPANDED     EAPType = 254
	EAPTypeEXPERIMENTAL EAPType = 255
)

// IsValid Verify if the value is a valid Type
// (may be coming from external source like incoming EAP packet)
func (t EAPType) IsValid() bool {
	switch t {
	case
		EAPTypeIDENTITY,
		EAPTypeNOTIFICATION,
		EAPTypeNAK,
		EAPTypeMD5CHALLENGE,
		EAPTypeOTP,
		EAPTypeGENTOKENCARD,
		EAPTypeCISCOLEAP,
		EAPTypeSIM,
		EAPTypeAKA,
		EAPTypeEAPMSCHAPV2,
		EAPTypeEXPANDED,
		EAPTypeEXPERIMENTAL:
		return true
	}
	return false
}
