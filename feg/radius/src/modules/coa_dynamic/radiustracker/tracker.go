package radiustracker

import "fbc/lib/go/radius"

// RadiusTracker an interface for radius packet tracking
type RadiusTracker interface {
	Set(r *radius.Request) error
	Get(r *radius.Request) (string, error)
}
