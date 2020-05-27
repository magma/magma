/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package indexers

import (
	"fmt"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/directoryd"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/indexer"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

const (
	sidID                      = "directoryd_session_id"
	sidVersion indexer.Version = 1
)

type sessionIDToIMSI struct{}

// NewSessionIDToIMSI returns a new state indexer that maps a session ID to an IMSI, for consumption via the directoryd service.
//
// Directoryd records are indexed as {IMSI -> {records...}}.
// SessionID indexer is an online generation of the derived reverse map, producing {session ID -> IMSI}.
//
// NOTE: the indexer provides a best-effort generation of the session ID -> IMSI mapping, meaning
//	- a {session ID -> IMSI} mapping may be missing even though the IMSI has a session ID record
//	- a {session ID -> IMSI} mapping may be stale
func NewSessionIDToIMSI() indexer.Indexer {
	return &sessionIDToIMSI{}
}

func (s *sessionIDToIMSI) GetID() string {
	return sidID
}

func (s *sessionIDToIMSI) GetVersion() indexer.Version {
	return sidVersion
}

func (s *sessionIDToIMSI) GetSubscriptions() []indexer.Subscription {
	return []indexer.Subscription{
		{Type: orc8r.DirectoryRecordType, KeyMatcher: indexer.MatchAll},
	}
}

// PrepareReindex needs no action since all storage is handled by directoryd service.
func (s *sessionIDToIMSI) PrepareReindex(from, to indexer.Version, isFirstReindex bool) error {
	return nil
}

func (s *sessionIDToIMSI) CompleteReindex(from, to indexer.Version) error {
	if from == 0 && to == 1 {
		return nil
	}
	return fmt.Errorf("unsupported from/to for CompleteReindex: %v to %v", from, to)
}

func (s *sessionIDToIMSI) Index(networkID string, states state.StatesByID) (indexer.StateErrors, error) {
	sessionIDToIMSI := map[string]string{}
	errs := indexer.StateErrors{}

	for id, st := range states {
		sessionID, imsi, err := getSessionIDAndIMSI(id, st)
		if err != nil {
			errs[id] = err
			continue
		}
		if sessionID == "" {
			glog.V(2).Infof("Session ID not found for record from %s", imsi)
			continue
		}

		sessionIDToIMSI[sessionID] = imsi
	}

	if len(sessionIDToIMSI) == 0 {
		return errs, nil
	}

	err := directoryd.MapSessionIDsToIMSIs(networkID, sessionIDToIMSI)
	if err != nil {
		return errs, errors.Wrapf(err, "update directoryd mapping of session IDs to IMSIs %+v", sessionIDToIMSI)
	}

	return errs, nil
}

// getSessionIDAndIMSI extracts session ID and IMSI from the state.
// Returns (session ID, IMSI, error).
func getSessionIDAndIMSI(id state.ID, st state.State) (string, string, error) {
	imsi := id.DeviceID

	// Cast to directory record
	record, ok := st.ReportedState.(*directoryd.DirectoryRecord)
	if !ok {
		return "", "", fmt.Errorf(
			"convert reported state (id: <%+v>, state: <%+v>) to type %s",
			id, st, orc8r.DirectoryRecordType,
		)
	}

	// Get session ID
	sessionID, err := record.GetSessionID()
	if err != nil {
		return "", "", errors.Wrap(err, "extract session ID from record")
	}

	return sessionID, imsi, nil
}
