/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package indexer

import (
	"magma/orc8r/cloud/go/services/state/protos"
)

func MakeProtoInfos(vs []*Versions) map[string]*protos.IndexerInfo {
	ret := map[string]*protos.IndexerInfo{}
	for _, v := range vs {
		ret[v.IndexerID] = MakeProtoInfo(v)
	}
	return ret
}

func MakeProtoInfo(v *Versions) *protos.IndexerInfo {
	return &protos.IndexerInfo{
		IndexerId:      v.IndexerID,
		ActualVersion:  uint32(v.Actual),
		DesiredVersion: uint32(v.Desired),
	}
}

func MakeVersions(ps map[string]*protos.IndexerInfo) []*Versions {
	var ret []*Versions
	for _, p := range ps {
		ret = append(ret, MakeVersion(p))
	}
	return ret
}

func MakeVersion(p *protos.IndexerInfo) *Versions {
	return &Versions{
		IndexerID: p.IndexerId,
		Actual:    Version(p.ActualVersion),
		Desired:   Version(p.DesiredVersion),
	}
}
