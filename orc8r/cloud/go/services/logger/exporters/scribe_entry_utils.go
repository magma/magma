/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package exporters

import (
	"encoding/json"
	"fmt"

	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
)

type ScribeLogEntry struct {
	Category string `json:"category"`
	Message  string `json:"message"`
}

type ScribeLogMessage struct {
	Int     map[string]int64  `json:"int,omitempty"`
	Normal  map[string]string `json:"normal,omitempty"`
	TagSet  []string          `json:"tagset,omitempty"`
	NormVec []string          `json:"normvector,omitempty"`
}

//convert a slice of protos.LogEntry into a slice of ScribeLogMessage.
// Add networkId and gatewayId into normal map of ScribeLogEntry if
// the original LogEntry had a valid hardware_id.
func ConvertToScribeLogEntries(entries []*protos.LogEntry) ([]*ScribeLogEntry, error) {
	scribeEntries := []*ScribeLogEntry{}
	for _, entry := range entries {
		if entry.Time == 0 {
			return nil, fmt.Errorf("ScribeLogEntry %v doesn't have time field set", entry)
		}
		scribeMsg := ScribeLogMessage{}
		// if any of the following fields are nil, they will be omitted when scribeMsg is marshalled into json.
		scribeMsg.Normal = entry.NormalMap
		scribeMsg.Int = entry.IntMap
		scribeMsg.TagSet = entry.TagSet
		scribeMsg.NormVec = entry.Normvector
		// append Time field to the int map
		if scribeMsg.Int == nil {
			scribeMsg.Int = map[string]int64{}
		}
		scribeMsg.Int["time"] = entry.Time
		// add gatewayId and networkId if it's a logEntry logged from a gateway
		if len(entry.HwId) != 0 {
			networkID, gatewayID, err := configurator.GetNetworkAndEntityIDForPhysicalID(entry.HwId)
			if err != nil {
				glog.Errorf("Error retrieving nwId and gwId for hwId %s in scribeExporter: %v\n", entry.HwId, err)
			}
			if scribeMsg.Normal == nil {
				scribeMsg.Normal = map[string]string{}
			}
			scribeMsg.Normal["networkId"] = networkID
			scribeMsg.Normal["gatewayId"] = gatewayID
		}
		// marshall scribeMsg into json
		msgJson, err := json.Marshal(scribeMsg)
		if err != nil {
			glog.Errorf("Error formatting scribeMsg %v in scribeExporter: %v\n", scribeMsg, err)
			continue
		}
		scribeEntries = append(scribeEntries, &ScribeLogEntry{Category: entry.Category, Message: string(msgJson)})
	}
	return scribeEntries, nil
}
