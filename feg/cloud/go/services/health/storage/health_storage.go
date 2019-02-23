/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package storage

import (
	"fmt"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/protos"
)

type healthStorage struct {
	store datastore.Api
}

const HealthStatusTableName string = "health"

func healthTable(networkID string) string {
	return datastore.GetTableName(networkID, HealthStatusTableName)
}

// Create a new Health Store
func NewHealthStore(ds datastore.Api) (HealthStorage, error) {
	hs := &healthStorage{ds}
	if hs.store == nil {
		return nil, fmt.Errorf("Nil Health Store datastore")
	}
	return hs, nil
}

func (s *healthStorage) UpdateHealth(networkID string, gatewayID string, health *fegprotos.HealthStats) error {
	if health == nil {
		return fmt.Errorf("Nil Gateway Health provided")
	}

	marshaledHealth, err := protos.MarshalIntern(health)
	if err != nil {
		return fmt.Errorf("Health Store: Marshalling error: %s for GW: %s", err, gatewayID)
	}
	err = s.store.Put(healthTable(networkID), gatewayID, marshaledHealth)
	if err != nil {
		return fmt.Errorf("Health Store Write Error: %s for GW: %s", err, gatewayID)
	}
	return nil
}

func (s *healthStorage) GetHealth(networkID string, gatewayID string) (*fegprotos.HealthStats, error) {
	marshaledHealth, _, err := s.store.Get(healthTable(networkID), gatewayID)
	if err != nil {
		return nil, fmt.Errorf("Feg Get Health Error: %s for Network: %s, GW: %s", err, networkID, gatewayID)
	}

	healthData := new(fegprotos.HealthStats)
	err = protos.Unmarshal(marshaledHealth, healthData)
	if err != nil {
		return nil, fmt.Errorf("Feg Health Unmarshaling Error: %s", err)
	}
	return healthData, nil
}
