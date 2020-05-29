/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package protos

import (
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/reindex"
)

func MakeProtoInfos(vs []*reindex.Version) map[string]*IndexerInfo {
	ret := map[string]*IndexerInfo{}
	for _, v := range vs {
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

func MakeVersions(ps map[string]*IndexerInfo) []*reindex.Version {
	var ret []*reindex.Version
	for _, p := range ps {
		ret = append(ret, MakeVersion(p))
	}
	return ret
}

func MakeVersion(p *IndexerInfo) *reindex.Version {
	return &reindex.Version{
		IndexerID: p.IndexerId,
		Actual:    indexer.Version(p.ActualVersion),
		Desired:   indexer.Version(p.DesiredVersion),
	}
}
