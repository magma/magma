/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSDstyle license found in the
LICENSE file in the root directory of this source tree.
*/

// package aka implements EAP-AKA EAP Method
package aka

import (
	"magma/feg/gateway/services/eap"
)

func NewIdentityReq(identifier uint8, attr eap.AttrType) eap.Packet {
	return []byte{
		eap.RequestCode,
		identifier,
		0, 12, // EAP Len
		TYPE,
		byte(SubtypeIdentity),
		0, 0,
		byte(attr),
		1,
		0, 0} // padding
}
