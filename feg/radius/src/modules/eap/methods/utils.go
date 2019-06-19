package methods

import (
	"fbc/cwf/radius/modules/eap/packet"
	"fbc/lib/go/radius"
)

// ToRadiusCode returs the RADIUS packet code which, as per RFCxxxx
// should carry the EAP payload of the given EAP Code
func ToRadiusCode(eapCode packet.Code) radius.Code {
	switch eapCode {
	case packet.CodeFAILURE:
		return radius.CodeAccessReject
	case packet.CodeSUCCESS:
		return radius.CodeAccessAccept
	case packet.CodeRESPONSE:
	case packet.CodeREQUEST:
		return radius.CodeAccessChallenge
	}
	return radius.CodeAccessReject
}
