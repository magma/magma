/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package models_test

import (
	"testing"

	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb/obsidian/models"

	"github.com/stretchr/testify/assert"
)

func TestProtoToModel(t *testing.T) {
	rule1 := &protos.PolicyRule{
		Id:           "rule1",
		Priority:     100,
		RatingGroup:  200,
		TrackingType: protos.PolicyRule_ONLY_OCS,
	}
	model1 := &models.PolicyRule{}
	assert.NoError(t, model1.FromProto(rule1))
	assert.Equal(t, model1.ID, rule1.Id)
	assert.Equal(t, *model1.Priority, rule1.Priority)
	assert.Equal(t, *model1.RatingGroup, rule1.RatingGroup)
	assert.Equal(t, *model1.MonitoringKey, "")
	assert.Equal(t, model1.TrackingType, "ONLY_OCS")
	assert.Nil(t, model1.Redirect)
	assert.Nil(t, model1.Qos)

	rule2 := &protos.PolicyRule{
		Id:            "rule3",
		Priority:      100,
		RatingGroup:   200,
		MonitoringKey: "1234",
		TrackingType:  protos.PolicyRule_OCS_AND_PCRF,
		Redirect: &protos.RedirectInformation{
			Support:       protos.RedirectInformation_ENABLED,
			AddressType:   protos.RedirectInformation_IPv4,
			ServerAddress: "127.0.0.1",
		},
		Qos: &protos.FlowQos{
			MaxReqBwUl: 1000,
			MaxReqBwDl: 2000,
		},
	}
	model2 := &models.PolicyRule{}
	assert.NoError(t, model2.FromProto(rule2))
	assert.Equal(t, *model2.MonitoringKey, rule2.MonitoringKey)
	assert.Equal(t, model2.TrackingType, "OCS_AND_PCRF")
	assert.Equal(t, model2.Redirect.Support, "ENABLED")
	assert.Equal(t, model2.Redirect.AddressType, "IPv4")
	assert.Equal(t, model2.Redirect.ServerAddress, rule2.Redirect.ServerAddress)
	assert.Equal(t, model2.Qos.MaxReqBwDl, rule2.Qos.MaxReqBwDl)
	assert.Equal(t, model2.Qos.MaxReqBwUl, rule2.Qos.MaxReqBwUl)

	flowProto := &protos.FlowDescription{
		Action: protos.FlowDescription_DENY,
		Match: &protos.FlowMatch{
			Ipv4Src:   "127.0.0.1",
			IpProto:   protos.FlowMatch_IPPROTO_UDP,
			UdpDst:    123,
			Direction: protos.FlowMatch_DOWNLINK,
		},
	}
	rule3 := &protos.PolicyRule{
		Id:       "rule3",
		FlowList: []*protos.FlowDescription{flowProto},
	}
	model3 := &models.PolicyRule{}
	assert.NoError(t, model3.FromProto(rule3))
	assert.Equal(t, len(model3.FlowList), 1)
	flowModel := model3.FlowList[0]
	assert.Equal(t, *flowModel.Action, "DENY")
	assert.Equal(t, flowModel.Match.IPV4Dst, flowProto.Match.Ipv4Dst)
	assert.Equal(t, *flowModel.Match.IPProto, "IPPROTO_UDP")
	assert.Equal(t, flowModel.Match.UDPDst, flowProto.Match.UdpDst)
	assert.Equal(t, flowModel.Match.Direction, "DOWNLINK")
}

func TestModelToProto(t *testing.T) {
	var rg uint32 = 1
	model1 := &models.PolicyRule{
		ID:           "rule1",
		RatingGroup:  &rg,
		TrackingType: "OCS_AND_PCRF",
	}
	rule1 := &protos.PolicyRule{}
	assert.NoError(t, model1.ToProto(rule1))
	assert.Equal(t, rule1.Id, model1.ID)
	assert.Equal(t, rule1.RatingGroup, *(model1.RatingGroup))
	assert.Equal(t, rule1.TrackingType, protos.PolicyRule_OCS_AND_PCRF)
	assert.Nil(t, rule1.Qos)
	assert.Nil(t, rule1.Redirect)

	mkey := "mkey1"
	action := "PERMIT"
	flowModel := &models.FlowDescription{
		Action: &action,
		Match: &models.FlowMatch{
			Direction: "UPLINK",
			IPV4Dst:   "127.0.0.1",
		},
	}
	model2 := &models.PolicyRule{
		ID:            "rule2",
		MonitoringKey: &mkey,
		TrackingType:  "ONLY_PCRF",
		Qos: &models.FlowQos{
			MaxReqBwDl: 1000,
			MaxReqBwUl: 2000,
		},
		Redirect: &models.RedirectInformation{
			AddressType:   "IPv6",
			ServerAddress: "127.0.0.1",
			Support:       "ENABLED",
		},
		FlowList: []*models.FlowDescription{flowModel},
	}
	rule2 := &protos.PolicyRule{}
	assert.NoError(t, model2.ToProto(rule2))
	assert.Equal(t, len(rule2.FlowList), 1)
	flowProto := rule2.FlowList[0]
	assert.Equal(t, flowProto.Action, protos.FlowDescription_PERMIT)
	assert.Equal(t, flowProto.Match.Direction, protos.FlowMatch_UPLINK)
	assert.Equal(t, flowProto.Match.Ipv4Dst, flowModel.Match.IPV4Dst)
	assert.Equal(t, rule2.Qos.MaxReqBwDl, model2.Qos.MaxReqBwDl)
	assert.Equal(t, rule2.Qos.MaxReqBwUl, model2.Qos.MaxReqBwUl)
	assert.Equal(t, rule2.Redirect.Support, protos.RedirectInformation_ENABLED)
	assert.Equal(t, rule2.Redirect.AddressType, protos.RedirectInformation_IPv6)
	assert.Equal(t, rule2.Redirect.ServerAddress, model2.Redirect.ServerAddress)
}
