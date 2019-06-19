package session

import "bytes"

type (
	// State the data to store per session
	State struct {
		MACAddress        string
		MSISDN            string
		UpstreamHost      string
		Tier              string
		RadiusSessionFBID uint64 // the FBID of the XWFEntRadiusSession created for this RADIUS session
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

// CreateSessionIDStrings format the session key from its 2 constituents
func CreateSessionIDStrings(callingStationID string, calledStationID string) string {
	return CreateSessionID([]byte(callingStationID), []byte(calledStationID))
}

// CreateSessionID format the session key from its 2 constituents
func CreateSessionID(callingStationID []byte, calledStationID []byte) string {
	var sessionID bytes.Buffer

	if calledStationID != nil {
		sessionID.Write(calledStationID)
	}
	if callingStationID != nil {
		sessionID.Write(callingStationID)
	}
	return sessionID.String()
}
