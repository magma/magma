//
// Copyright (c) Facebook, Inc. and its affiliates.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//

// +build link_local_service

// package aka implements EAP-AKA provider
package provider

import (
	"errors"
	"log"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap/providers"
	"magma/feg/gateway/services/eap/providers/aka"
	"magma/feg/gateway/services/eap/providers/aka/servicers"
	managed_configs "magma/gateway/mconfig"
)

func NewServiced(srvsr *servicers.EapAkaSrv) providers.Method {
	return &providerImpl{EapAkaSrv: srvsr}
}

// Handle handles passed EAP-AKA payload & returns corresponding result
// this Handle implementation is using GRPC based AKA provider service
func (prov *providerImpl) Handle(msg *protos.Eap) (*protos.Eap, error) {
	if msg == nil {
		return nil, errors.New("Invalid EAP AKA Message")
	}
	prov.RLock()
	if prov.EapAkaSrv == nil {
		// servicer is not initialized, relock, recheck, create
		prov.RUnlock()
		prov.Lock()
		if prov.EapAkaSrv == nil {
			akaConfigs := &mconfig.EapAkaConfig{}
			err := managed_configs.GetServiceConfigs(aka.EapAkaServiceName, akaConfigs)
			if err != nil {
				log.Printf("Error getting EAP AKA service configs: %s", err)
				akaConfigs = nil
			}
			prov.EapAkaSrv, err = servicers.NewEapAkaService(akaConfigs)
			if err != nil || prov.EapAkaSrv == nil {
				log.Fatalf("failed to create EAP AKA Service: %v", err) // should never happen
			}
		}
		prov.Unlock()
		prov.RLock()
	}
	defer prov.RUnlock()
	return prov.EapAkaSrv.HandleImpl(msg)
}
