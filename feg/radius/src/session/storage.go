/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package session

import (
	"errors"
)

type (
	// State the data to store per session
	State struct {
		NextCoAIdentifier byte
		MACAddress        string
		MSISDN            string
		UpstreamHost      string
		Tier              string
		RadiusSessionFBID uint64 // the FBID of the XWFEntRadiusSession created for this RADIUS session
		AcctSessionID     string
		CalledStationID   string
	}

	// GlobalStorage an interface for session-level storage, which allows
	// access to any session state, by using sessionID as key
	GlobalStorage interface {
		Get(sessionID string) (*State, error)
		Set(sessionID string, state State) error
		Reset(sessionID string) error
	}

	// Storage an interface for session-level storage, which allows access
	// to one specific session state. This interface is to be used on
	// session-specific flows, like accounting
	Storage interface {
		Get() (*State, error)
		Set(state State) error
		Reset() error
	}
)

// sessionStorage a wrapper implementation for globalStorage to get/set
// in a session-specific state
type sessionStorage struct {
	globalStorage GlobalStorage
	sessionID     string
}

var (
	// ErrInvalidDataFormat indicate we have an invalid data as state
	ErrInvalidDataFormat = errors.New("invalid data format found in storage")
)

func (s *sessionStorage) Get() (*State, error) {
	return s.globalStorage.Get(s.sessionID)
}

func (s *sessionStorage) Set(state State) error {
	return s.globalStorage.Set(s.sessionID, state)
}

func (s *sessionStorage) Reset() error {
	return s.globalStorage.Reset(s.sessionID)
}

// NewSessionStorage returns a session-specific storage for use by listeners
func NewSessionStorage(globalStorage GlobalStorage, sessionID string) Storage {
	return &sessionStorage{
		globalStorage: globalStorage,
		sessionID:     sessionID,
	}
}
