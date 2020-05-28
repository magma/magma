/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package index

import (
	"fmt"
	"time"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/metrics"
	state_types "magma/orc8r/cloud/go/services/state/types"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

type Error string

const (
	// ErrDefault is the default error reported for index failure.
	ErrDefault Error = "state index error"
	// ErrIndexPerState indicates an Index error occured for specific keys.
	ErrIndexPerState Error = "state index error: per-state errors"

	// ErrIndex indicates error source is indexer Index call.
	ErrIndex Error = "state index error: error from Index"

	maxRetry     = 3
	defaultSleep = 10 * time.Second
)

// Index forwards states to all registered indexers, according to their subscriptions.
// Returns after completing attempt at indexing states.
func Index(networkID string, states state_types.StatesByID) {
	errs := indexImpl(networkID, states)
	for _, e := range errs {
		glog.Error(e)
	}

	glog.V(2).Infof("Completed state index for network %s with %d states", networkID, len(states))
}

func indexImpl(networkID string, states state_types.StatesByID) []error {
	errByIdx := map[string]error{}
	indexers := indexer.GetAllIndexers()

	// Retry indexing, up to a max number of times
	for i := 0; len(indexers) != 0 && i < maxRetry; i++ {
		var failed []indexer.Indexer
		for _, x := range indexers {
			err := indexOne(networkID, x, states)
			if err != nil {
				errByIdx[x.GetID()] = err
				failed = append(failed, x)
			}
		}

		indexers = failed
		clock.Sleep(defaultSleep)
	}

	var errs []error
	for _, idx := range indexers {
		errs = append(errs, errByIdx[idx.GetID()])
	}

	return errs
}

func indexOne(networkID string, idx indexer.Indexer, states state_types.StatesByID) error {
	filtered := filterStates(idx, states)
	if len(filtered) == 0 {
		return nil
	}

	id := idx.GetID()
	version := getVersion(idx)

	errs, err := idx.Index(networkID, filtered)
	if err != nil {
		return wrap(err, ErrIndex, id)
	}
	if len(errs) == len(filtered) {
		err := errors.New("all state IDs experienced per-state index errors")
		return wrap(err, ErrIndex, id)
	} else if len(errs) != 0 {
		metrics.IndexErrors.WithLabelValues(id, version, metrics.SourceValueIndex).Add(float64(len(errs)))
		err := wrap(fmt.Errorf("%s", errs), ErrIndexPerState, id)
		glog.Warning(err)
		return nil
	}

	return nil
}

// Filter to states matching at least one subscription
func filterStates(idx indexer.Indexer, states state_types.StatesByID) state_types.StatesByID {
	ret := state_types.StatesByID{}
	subs := idx.GetSubscriptions()
	for id, st := range states {
		for _, s := range subs {
			if s.Match(id) {
				ret[id] = st
				break
			}
		}
	}
	return ret
}

func getVersion(idx indexer.Indexer) string {
	return fmt.Sprint(idx.GetVersion())
}

func wrap(err error, sentinel Error, indexerID string) error {
	var wrap string
	switch indexerID {
	case "":
		wrap = string(sentinel)
	default:
		wrap = fmt.Sprintf("%s for idx %s", sentinel, indexerID)
	}
	return errors.Wrap(err, wrap)
}
