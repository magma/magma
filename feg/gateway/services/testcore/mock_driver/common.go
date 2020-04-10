/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package mock_driver

import (
	"fmt"

	"magma/feg/cloud/go/protos"
)

type CreditControlRequestPK struct {
	imsi        string
	requestType protos.CCRequestType
}

func NewCCRequestPK(imsi string, requestType protos.CCRequestType) CreditControlRequestPK {
	return CreditControlRequestPK{
		imsi:        imsi,
		requestType: requestType,
	}
}

func (r CreditControlRequestPK) String() string {
	return fmt.Sprintf("Imsi: %v, Type: %v", r.imsi, r.requestType)
}

func EqualWithinDelta(a, b, delta uint64) bool {
	if b >= a && b-a <= delta {
		return true
	}
	if a >= b && a-b <= delta {
		return true
	}
	return false
}
