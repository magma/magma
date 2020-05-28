/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package protos

import (
	"magma/orc8r/cloud/go/services/state/indexer/reindex"
)

func MakeProtoInfos(versions []*reindex.Version) map[string]*IndexerInfo {
	ret := map[string]*IndexerInfo{}
	for _, v := range versions {
		ret[v.IndexerID] = MakeProtoInfo(v)
	}
	return ret
}

func MakeProtoInfo(v *reindex.Version) *IndexerInfo {
	return &IndexerInfo{
		IndexerId:      v.IndexerID,
		ActualVersion:  uint32(v.Actual),
		DesiredVersion: uint32(v.Desired),
	}
}
