/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package models

import (
	"crypto/x509"
	"fmt"
	"reflect"

	"magma/orc8r/cloud/go/obsidian/models"
	"magma/orc8r/cloud/go/protos"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"

	"github.com/go-openapi/strfmt"
	"github.com/golang/protobuf/proto"
)

var formatsRegistry = strfmt.NewFormats()

// Config conversion for magmad

func (m *MagmadGatewayConfig) ValidateModel() error {
	if err := m.ValidateGatewayConfig(); err != nil {
		return err
	}
	return m.Validate(formatsRegistry)
}

func (m *MagmadGatewayConfig) ToServiceModel() (interface{}, error) {
	ret := &magmadprotos.MagmadGatewayConfig{}
	protos.FillIn(m, ret)
	if err := magmadprotos.ValidateGatewayConfig(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (m *MagmadGatewayConfig) FromServiceModel(magmadModel interface{}) error {
	_, ok := magmadModel.(*magmadprotos.MagmadGatewayConfig)
	if !ok {
		return fmt.Errorf(
			"Invalid magmad config type. Expected *MagmadGatewayConfig but got %s",
			reflect.TypeOf(magmadModel),
		)
	}
	protos.FillIn(magmadModel, m)
	return nil
}

// Other magmad model conversion

func NetworkRecordFromProto(msg *magmadprotos.MagmadNetworkRecord) *NetworkRecord {
	ret := &NetworkRecord{}
	protos.FillIn(msg, ret)
	return ret
}

func (m *NetworkRecord) ToProto() *magmadprotos.MagmadNetworkRecord {
	ret := &magmadprotos.MagmadNetworkRecord{}
	protos.FillIn(m, ret)
	return ret
}

func (m *NetworkRecord) ValidateModel() error {
	if err := m.ValidateNetworkRecord(); err != nil {
		return err
	}
	return m.Validate(formatsRegistry)
}

// Verify validates given GatewayConfigs
func (record *AccessGatewayRecord) Verify() error {
	if record == nil {
		return fmt.Errorf("Nil AccessGatewayRecord pointer")
	}
	err := record.Validate(formatsRegistry)
	if err != nil {
		err = models.ValidateErrorf("AccessGatewayRecord Validation Error: %s", err)
	}
	err = verifyKey(record.Key)
	if err != nil {
		err = models.ValidateErrorf("Key Validation Error: %s", err)
	}
	return err
}

// MutableGatewayRecord's FromMconfig fills in models.MutableGatewayRecord struct
// from passed protos.AccessGatewayRecord.
// For now the 'fill in" is automatic & relies on Go reflect based
// protos.FillIn function, so be carefull how you name & 'type' your fields
func (record *MutableGatewayRecord) FromMconfig(msg proto.Message) error {
	mrcrd, ok := msg.(*magmadprotos.AccessGatewayRecord)
	if !ok {
		return fmt.Errorf(
			"Invalid Source Type %s, *protos.AccessGatewayRecord expected",
			reflect.TypeOf(mrcrd))
	}
	if record != nil && mrcrd != nil {
		protos.FillIn(mrcrd, record)
		err := fillKeyFromMconfig(mrcrd.Key, record.Key)
		if err != nil {
			return err
		}
		return record.Verify()
	}
	return nil
}

// mutablegatewayrecord's ToMconfig fills in passed protos.AccessGatewayRecord
// struct from receiver's models.AccessGatewayRecord
// For now the 'fill in" is automatic & relies on Go reflect based
// protos.FillIn function, so be carefull how you name & 'type' your fields
func (record *MutableGatewayRecord) ToMconfig(msg proto.Message) error {
	mrcrd, ok := msg.(*magmadprotos.AccessGatewayRecord)
	if !ok {
		return fmt.Errorf(
			"Invalid Destination Type %s, *protos.AccessGatewayRecord expected",
			reflect.TypeOf(mrcrd))
	}
	if record != nil && mrcrd != nil {
		protos.FillIn(record, mrcrd)
		key, err := fillKeyToMconfig(record.Key)
		if err != nil {
			return fmt.Errorf("Failed to fill in the key")
		}
		mrcrd.Key = key
	}
	return nil
}

// Verify validates given GatewayConfigs
func (record *MutableGatewayRecord) Verify() error {
	if record == nil {
		return fmt.Errorf("Nil MutableGatewayRecord pointer")
	}
	err := record.Validate(formatsRegistry)
	if err != nil {
		err = models.ValidateErrorf("MutableGatewayRecord Validation Error: %s", err)
	}
	verifyKey(record.Key)
	if err != nil {
		err = models.ValidateErrorf("Key Validation Error: %s", err)
	}
	return err
}

func fillKeyFromMconfig(mkey *protos.ChallengeKey, key *ChallengeKey) error {
	if mkey == nil {
		return nil
	}
	t, ok := protos.ChallengeKey_KeyType_name[int32(mkey.KeyType)]
	if !ok {
		return fmt.Errorf("Unknown ChallengeKey Type: %s", mkey.KeyType)
	}
	key.KeyType = t

	if len(mkey.Key) > 0 {
		key.Key = (*strfmt.Base64)(&mkey.Key)
	} else {
		key.Key = nil
	}
	return nil
}

func fillKeyToMconfig(key *ChallengeKey) (*protos.ChallengeKey, error) {
	if key == nil {
		return nil, nil
	}
	mkey := new(protos.ChallengeKey)
	t, ok := protos.ChallengeKey_KeyType_value[key.KeyType]
	if !ok {
		return mkey, fmt.Errorf("Invalid ChallengeKey Type: %s", key.KeyType)
	}
	mkey.KeyType = protos.ChallengeKey_KeyType(t)

	if key.Key != nil {
		mkey.Key = []byte(*key.Key)
	}
	return mkey, nil
}

func verifyKey(key *ChallengeKey) error {
	if key == nil {
		return nil
	}
	switch key.KeyType {
	case "ECHO":
		if key.Key != nil {
			return fmt.Errorf("ECHO mode should not have key value")
		} else {
			return nil
		}
	case "SOFTWARE_ECDSA_SHA256":
		if key.Key == nil {
			return fmt.Errorf("No key supplied")
		}
		_, err := x509.ParsePKIXPublicKey([]byte(*key.Key))
		if err != nil {
			return fmt.Errorf("Failed to parse key: %s", err)
		}
		return nil
	default:
		return fmt.Errorf("Unknown key type: %s", key.KeyType)
	}
}
