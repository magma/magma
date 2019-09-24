package radiustracker

import (
	"errors"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2866"
	"fmt"
	"net"
	"sync"
)

type memoryRadiusTracker struct {
	storage sync.Map
}

func (m *memoryRadiusTracker) Set(r *radius.Request) error {
	key, err := genKey(r)
	if nil != err {
		return err
	}

	// Obtain the source ip of the radius packet
	host, _, err := net.SplitHostPort(r.RemoteAddr.String())
	if err != nil {
		return err
	}

	m.storage.Store(key, host)
	return nil
}

func (m *memoryRadiusTracker) Get(r *radius.Request) (string, error) {
	key, err := genKey(r)
	if nil != err {
		return "", err
	}

	value, found := m.storage.Load(key)
	if !found {
		return "", fmt.Errorf("radius request wasn't tracked yet (key: %s)", key)
	}

	ip, ok := value.(string)
	if !ok {
		return "", errors.New("ip failed to deserialize")
	}

	return ip, nil
}

func genKey(r *radius.Request) (string, error) {
	sessionID := r.Get(rfc2866.AcctSessionID_Type)
	if sessionID == nil {
		return "", errors.New("accounting session id attribute not found")
	}

	return fmt.Sprintf("__%s", string(sessionID)), nil
}

// NewRadiusTracker Create a new radius packet tracker
func NewRadiusTracker() RadiusTracker {
	return &memoryRadiusTracker{
		storage: sync.Map{},
	}
}
