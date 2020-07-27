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

package oc

import (
	"errors"
	"sync"

	"fbc/lib/go/oc/ocstats"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
)

// Viewer acts as a views provider.
type Viewer interface {
	Views() []*view.View
}

var registeredViewers sync.Map

func init() {
	MustRegisterViewer("proc", Views{})
	MustRegisterViewer("http", Views{
		ochttp.ClientCompletedCount,
		ochttp.ClientSentBytesDistribution,
		ochttp.ClientReceivedBytesDistribution,
		ochttp.ClientRoundtripLatencyDistribution,
		ochttp.ClientCompletedCount,
		ochttp.ServerRequestCountView,
		ochttp.ServerRequestBytesView,
		ochttp.ServerResponseBytesView,
		ochttp.ServerLatencyView,
		ochttp.ServerRequestCountByMethod,
		ochttp.ServerResponseCountByStatusCode,
		ocstats.HTTPServerResponseCountByStatusAndPath,
	})
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

// MustRegisterViewer works like RegisterViewer but panics on error.
func MustRegisterViewer(name string, viewer Viewer) {
	if err := RegisterViewer(name, viewer); err != nil {
		panic(err)
	}
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
