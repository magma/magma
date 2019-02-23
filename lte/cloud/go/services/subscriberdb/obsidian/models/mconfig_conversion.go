/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package models

import (
	"fmt"
	"reflect"

	lteprotos "magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/obsidian/models"
	"magma/orc8r/cloud/go/protos"

	"github.com/go-openapi/strfmt"
	"github.com/golang/protobuf/proto"
)

var sharedFormatsRegistry = strfmt.NewFormats()

func init() {
	// Echo encodes/decodes base64 encoded byte arrays, no verification needed
	b64 := strfmt.Base64([]byte(nil))
	sharedFormatsRegistry.Add(
		"byte", &b64, func(_ string) bool { return true })
}

// Subscriber's FromMconfig fills in models.Subscriber struct from
// passed protos.SubscriberData
func (sub *Subscriber) FromMconfig(psm proto.Message) error {
	psub, ok := psm.(*lteprotos.SubscriberData)
	if !ok {
		return fmt.Errorf(
			"Invalid Source Type %s, *protos.SubscriberData expected",
			reflect.TypeOf(psm))
	}
	if sub != nil {
		sub.Lte = nil
		if psub != nil {
			protos.FillIn(psub, sub)
			sub.ID = SubscriberID(lteprotos.SidString(psub.Sid))
			if sub.Lte != nil && psub.Lte != nil {
				t, ok := lteprotos.LTESubscription_LTESubscriptionState_name[int32(psub.Lte.State)]
				if ok {
					sub.Lte.State = t
				} else {
					sub.Lte.State = "INACTIVE"
				}
				t, ok = lteprotos.LTESubscription_LTEAuthAlgo_name[int32(psub.Lte.AuthAlgo)]
				if ok {
					sub.Lte.AuthAlgo = t
				} else {
					sub.Lte.AuthAlgo = "MILENAGE"
				}
				if len(psub.Lte.AuthKey) > 0 {
					sub.Lte.AuthKey = (*strfmt.Base64)(&psub.Lte.AuthKey)
				} else {
					sub.Lte.AuthKey = nil
				}
				if len(psub.Lte.AuthOpc) > 0 {
					sub.Lte.AuthOpc = (*strfmt.Base64)(&psub.Lte.AuthOpc)
				} else {
					sub.Lte.AuthOpc = nil
				}
			}
			return sub.Verify()
		}
	}
	return nil
}

// Subscriber's ToMconfig fills in passed protos.SubscriberData struct from
// receiver's protos.SubscriberData
func (sub *Subscriber) ToMconfig(psm proto.Message) error {
	psub, ok := psm.(*lteprotos.SubscriberData)
	if !ok {
		return fmt.Errorf(
			"Invalid Destination Type %s, *protos.SubscriberData expected",
			reflect.TypeOf(psm))
	}
	if sub != nil || psub != nil {
		protos.FillIn(sub, psub)
		t, err := lteprotos.SidProto(string(sub.ID))
		if err != nil {
			return err
		}
		if sub.Lte != nil {
			if psub.Lte == nil {
				psub.Lte = new(lteprotos.LTESubscription)
			}
			t, ok := lteprotos.LTESubscription_LTESubscriptionState_value[sub.Lte.State]
			if ok {
				psub.Lte.State = lteprotos.LTESubscription_LTESubscriptionState(t)
			} else {
				psub.Lte.State = 0
			}
			t, ok = lteprotos.LTESubscription_LTEAuthAlgo_value[sub.Lte.AuthAlgo]
			if ok {
				psub.Lte.AuthAlgo = lteprotos.LTESubscription_LTEAuthAlgo(t)
			} else {
				psub.Lte.AuthAlgo = 0
			}
			if sub.Lte.AuthKey != nil {
				psub.Lte.AuthKey = []byte(*sub.Lte.AuthKey)
			}
			if sub.Lte.AuthOpc != nil {
				psub.Lte.AuthOpc = []byte(*sub.Lte.AuthOpc)
			}
		}
		psub.Sid = t
	}
	return nil
}

// Verify validates given Subscriber
func (sub *Subscriber) Verify() error {
	if sub == nil {
		return fmt.Errorf("Nil Subscriber pointer")
	}
	err := sub.Validate(sharedFormatsRegistry)
	if err != nil {
		err = models.ValidateErrorf("Subscriber Validation Error: %s", err)
	}
	return err
}

// Verify validates given SubscriberID
func (sid *SubscriberID) Verify() error {
	if sid == nil {
		return fmt.Errorf("Nil SubscriberID pointer")
	}
	err := sid.Validate(sharedFormatsRegistry)
	if err != nil {
		err = models.ValidateErrorf("SubscriberID Validation Error: %s", err)
	}
	return err
}
