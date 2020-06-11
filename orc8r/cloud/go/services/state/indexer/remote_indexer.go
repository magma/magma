/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package indexer

import (
	"context"
	"strings"

	"magma/orc8r/cloud/go/services/state/protos"
	state_types "magma/orc8r/cloud/go/services/state/types"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
)

// remoteIndexer identifies a remote state indexer.
// The remote indexer's fields are cached at the state service.
type remoteIndexer struct {
	// service name of the indexer
	// should always be uppercase to match service registry convention
	service string
	// version of the indexer
	version Version
	// types is the types of state the indexer should receive
	types []string
}

func NewRemoteIndexer(serviceName string, version Version, types ...string) Indexer {
	return &remoteIndexer{service: strings.ToUpper(serviceName), version: version, types: types}
}

func (r *remoteIndexer) GetID() string {
	return r.service
}

func (r *remoteIndexer) GetVersion() Version {
	return r.version
}

func (r *remoteIndexer) GetTypes() []string {
	return r.types
}

func (r *remoteIndexer) PrepareReindex(from, to Version, isFirstReindex bool) error {
	c, err := r.getIndexerClient()
	if err != nil {
		return err
	}
	_, err = c.PrepareReindex(context.Background(), &protos.PrepareReindexRequest{
		IndexerId:   r.service,
		FromVersion: uint32(from),
		ToVersion:   uint32(to),
		IsFirst:     isFirstReindex,
	})
	return err
}

func (r *remoteIndexer) CompleteReindex(from, to Version) error {
	c, err := r.getIndexerClient()
	if err != nil {
		return err
	}
	_, err = c.CompleteReindex(context.Background(), &protos.CompleteReindexRequest{
		IndexerId:   r.service,
		FromVersion: uint32(from),
		ToVersion:   uint32(to),
	})
	return err
}

func (r *remoteIndexer) Index(networkID string, states state_types.StatesByID) (state_types.StateErrors, error) {
	if len(states) == 0 {
		return nil, nil
	}
	var reporterHWID string
	for _, st := range states {
		reporterHWID = st.ReporterID
		break
	}

	c, err := r.getIndexerClient()
	if err != nil {
		return nil, err
	}

	pStates, err := state_types.MakeProtoStates(states)
	if err != nil {
		return nil, err
	}
	res, err := c.Index(context.Background(), &protos.IndexRequest{
		States:       pStates,
		NetworkId:    networkID,
		ReporterHwid: reporterHWID,
	})
	if err != nil {
		return nil, err
	}

	return state_types.MakeStateErrors(res.StateErrors), nil
}

func (r *remoteIndexer) getIndexerClient() (protos.IndexerClient, error) {
	conn, err := registry.GetConnection(r.service)
	if err != nil {
		initErr := merrors.NewInitError(err, r.service)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewIndexerClient(conn), nil
}
