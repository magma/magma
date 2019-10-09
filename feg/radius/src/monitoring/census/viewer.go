package census

import (
	"errors"
	"sync"

	"go.opencensus.io/stats/view"
)

// Viewer acts as a views provider.
type Viewer interface {
	Views() []*view.View
}

var registeredViewers sync.Map

func init() {
	if err := RegisterViewer("proc", Views{}); err != nil {
		panic(err)
	}
}

// ErrViewerExist is returned by RegisterViewer on name collision.
var ErrViewerExist = errors.New("oc: viewer already exist")

// RegisterViewer registers the provided Viewer with oc by its name. It will
// be accessed when instantiating census from configuration.
func RegisterViewer(name string, viewer Viewer) error {
	if _, loaded := registeredViewers.LoadOrStore(name, viewer); loaded {
		return ErrViewerExist
	}
	return nil
}

// GetViewer returns Viewer for the given viewer name.
func GetViewer(name string) Viewer {
	if v, ok := registeredViewers.Load(name); ok {
		return v.(Viewer)
	}
	return nil
}

// Views attaches the methods of Viewer to []*view.View.
type Views []*view.View

// Views implements Viewer interface.
func (v Views) Views() []*view.View {
	return v
}
