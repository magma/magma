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
	// ErrIndexPerState indicates an Index error occurred for specific keys.
	ErrIndexPerState Error = "state index error: per-state errors"

	// ErrIndex indicates error source is indexer Index call.
	ErrIndex Error = "state index error: error from Index"

	maxRetry          = 3
	nIndexWorkers     = 5
	defaultIndexSleep = 10 * time.Second
)

// MustIndex forwards states to all registered indexers, according to their
// subscriptions.
// Per-state indexing errors are logged and reported as metrics.
// Overarching indexing errors are retried, then eventually logged.
// Returns after completing attempt at indexing states.
func MustIndex(networkID string, states state_types.StatesByID) {
	errs, err := Index(networkID, states)
	if err != nil {
		// Since we don't have a good way of recovering from failed goroutines
		// right now, fail immediately to restart the service to a good state
		glog.Fatalf("Error getting indexers during Index goroutine: %v", err)
	}
	for _, e := range errs {
		glog.Error(e)
	}
	glog.V(2).Infof("Completed state index for network %s with %d states", networkID, len(states))
}

// Index makes index calls via worker goroutines.
//	- each indexer gets up to maxRetry attempts
//	- returns after all goroutines have completed
// Prefer MustIndex except where receiving the returned errors is relevant.
func Index(networkID string, states state_types.StatesByID) ([]error, error) {
	index := func(indexers chan indexer.Indexer, out chan error) {
		for x := range indexers {
			var indexErr error
			for i := 0; i < maxRetry; i++ {
				indexErr = indexOne(networkID, x, states)
				if indexErr == nil {
					break
				}
				clock.Sleep(defaultIndexSleep)
			}
			out <- indexErr
		}
	}
	in := make(chan indexer.Indexer)
	out := make(chan error)
	for i := 0; i < nIndexWorkers; i++ {
		go index(in, out)
	}

	indexers, err := indexer.GetIndexers()
	if err != nil {
		return nil, err
	}
	go func() {
		for _, x := range indexers {
			in <- x
		}
		close(in)
	}()

	var indexErrs []error
	for i := 0; i < len(indexers); i++ {
		if e := <-out; e != nil {
			indexErrs = append(indexErrs, e)
		}
	}

	return indexErrs, nil
}

func indexOne(networkID string, idx indexer.Indexer, states state_types.StatesByID) error {
	filtered := indexer.FilterStates(idx.GetTypes(), states)
	if len(filtered) == 0 {
		return nil
	}

	id := idx.GetID()
	version := getVersion(idx)

	indexErrs, err := idx.Index(networkID, filtered)
	if err != nil {
		return wrap(err, ErrIndex, id)
	}
	if len(indexErrs) == len(filtered) {
		err := errors.New("all state IDs experienced per-state index errors")
		return wrap(err, ErrIndex, id)
	} else if len(indexErrs) != 0 {
		metrics.IndexErrors.WithLabelValues(id, version, metrics.SourceValueIndex).Add(float64(len(indexErrs)))
		err := wrap(fmt.Errorf("%s", indexErrs), ErrIndexPerState, id)
		glog.Warning(err)
		return nil
	}

	return nil
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
