/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
	globalStorage      GlobalStorage
	sessionID          string
	generatedSessionID string
}

var (
	// ErrInvalidDataFormat indicate we have an invalid data as state
	ErrInvalidDataFormat = errors.New("invalid data format found in storage")
)

func (s *sessionStorage) Get() (*State, error) {
	state, err := s.globalStorage.Get(s.sessionID)
	if (err != nil || state == nil) && len(s.generatedSessionID) > 0 && s.generatedSessionID != s.sessionID {
		stateGen, errGen := s.globalStorage.Get(s.generatedSessionID)
		if errGen == nil && stateGen != nil {
			state, err = stateGen, errGen
		}
	}
	return state, err
}

func (s *sessionStorage) Set(state State) error {
	return s.globalStorage.Set(s.sessionID, state)
}

func (s *sessionStorage) Reset() error {
	return s.globalStorage.Reset(s.sessionID)
}

// NewSessionStorageExt returns a session-specific storage with generated Session ID "backup" key for use by listeners
func NewSessionStorageExt(globalStorage GlobalStorage, sessionID, generatedSessionID string) Storage {
	return &sessionStorage{
		globalStorage:      globalStorage,
		sessionID:          sessionID,
		generatedSessionID: generatedSessionID,
	}
}

// NewSessionStorage returns a session-specific storage for use by listeners
func NewSessionStorage(globalStorage GlobalStorage, sessionID string) Storage {
	return NewSessionStorageExt(globalStorage, sessionID, "")
}
