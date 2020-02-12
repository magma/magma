/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package factory

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

type mconfigFactory struct {
	sync.RWMutex
	builders []MconfigBuilder
	clock    clock
}

type clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}

var factory = mconfigFactory{
	builders: []MconfigBuilder{},
	clock:    &realClock{},
}

func RegisterMconfigBuilders(builders ...MconfigBuilder) {
	factory.Lock()
	defer factory.Unlock()
	for _, builder := range builders {
		registerUnsafe(builder)
	}
}

func RegisterMconfigBuilder(builder MconfigBuilder) {
	factory.Lock()
	defer factory.Unlock()
	registerUnsafe(builder)
}

func registerUnsafe(builder MconfigBuilder) {
	factory.builders = append(factory.builders, builder)
}

// Create an mconfig by delegating to all builders that have been registered
// with the factory and append those results together.
// Note that the keys which builders return must be globally unique.
func CreateMconfig(networkId string, gatewayId string) (*protos.GatewayConfigs, error) {
	factory.RLock()
	defer factory.RUnlock()

	ret := map[string]*any.Any{}
	for _, builder := range factory.builders {
		subConfig, err := builder.Build(networkId, gatewayId)
		if err != nil {
			return nil, err
		}

		for k, v := range subConfig {
			_, ok := ret[k]
			if ok {
				return nil, fmt.Errorf("mconfig builder returned result for duplicate key %s", k)
			}

			vAny, err := ptypes.MarshalAny(v)
			if err != nil {
				return nil, fmt.Errorf("Error marshaling builder value to Any: %s", err)
			}
			ret[k] = vAny
		}
	}
	return &protos.GatewayConfigs{
		ConfigsByKey: ret,
		Metadata: &protos.GatewayConfigsMetadata{
			CreatedAt: uint64(factory.clock.Now().Unix()),
		},
	}, nil
}

// Methods below exist ONLY for testing - thus the required but unused *testing.T param
// DO NOT USE IN ANYTHING BUT TESTS
func ClearMconfigBuilders(_ *testing.T) {
	factory.Lock()
	factory.builders = factory.builders[:0]
	factory.Unlock()
}

func SetClock(_ *testing.T, clock clock) {
	factory.Lock()
	factory.clock = clock
	factory.Unlock()
}
