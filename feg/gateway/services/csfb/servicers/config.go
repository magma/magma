/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"fmt"
	"strconv"
	"strings"

	fegMconfig "magma/gateway/mconfig"
	"magma/lte/cloud/go/protos/mconfig"

	"github.com/golang/glog"
)

const (
	DefaultMMEName = ".mmec01.mmegi0001.mme.epc.mnc001.mcc001.3gppnetwork.org"
	MMECLength     = 2
	MMEGILength    = 4
	MNCLength      = 3
	MCCLength      = 3
	MMEServiceName = "mme"
)

// ConstructMMEName constructs MME name from mconfig
func ConstructMMEName() (string, error) {
	mmeConfig, err := getMMEConfig()
	if err != nil {
		glog.V(2).Infof(
			"Failed to retrieve MME config: %s, using default MME name: %s",
			err,
			DefaultMMEName,
		)
		return DefaultMMEName, nil
	}

	mnc := mmeConfig.GetCsfbMnc()
	mnc = fieldLengthCorrection(mnc, MNCLength)

	mcc := mmeConfig.GetCsfbMcc()
	mcc = fieldLengthCorrection(mcc, MCCLength)

	mmeCode := strconv.Itoa(int(mmeConfig.GetMmeCode()))
	mmeCode = fieldLengthCorrection(mmeCode, MMECLength)

	mmeGid := strconv.Itoa(int(mmeConfig.GetMmeGid()))
	mmeGid = fieldLengthCorrection(mmeGid, MMEGILength)

	mmeName := fmt.Sprintf(
		".mmec%s.mmegi%s.mme.epc.mnc%s.mcc%s.3gppnetwork.org",
		mmeCode,
		mmeGid,
		mnc,
		mcc,
	)

	return mmeName, nil
}

func getMMEConfig() (*mconfig.MME, error) {
	mmeConfig := &mconfig.MME{}
	err := fegMconfig.GetServiceConfigs(MMEServiceName, mmeConfig)
	if err != nil {
		return nil, err
	}
	return mmeConfig, nil
}

func fieldLengthCorrection(field string, requiredLength int) string {
	prefix := strings.Repeat("0", requiredLength-len(field))
	return prefix + field
}
