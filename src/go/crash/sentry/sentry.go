// Copyright 2021 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sentry

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/magma/magma/src/go/crash"
	"github.com/pkg/errors"
)

//go:generate go run github.com/golang/mock/mockgen -source=sentry.go -package mock_sentry -destination mock_sentry/mock_sentry_hub.go sentryHub

// sentryHub defines a subset sentry.Hub functions that we use.
type sentryHub interface {
	AddBreadcrumb(breadcrumb *sentry.Breadcrumb, hint *sentry.BreadcrumbHint)
	Recover(err interface{}) *sentry.EventID
	Flush(timeout time.Duration) bool
}

type Crash struct {
	sentryHub
}

// NewCrash returns a Crash with a singleton sentry Hub.
// All calls on that Crash will be preformed on that Hub.
func NewCrash(options sentry.ClientOptions) crash.Crash {
	client, err := sentry.NewClient(options)
	if err != nil {
		panic(errors.Wrapf(err, "sentry.ClientOptions=%+v", options))
	}
	newHub := sentry.NewHub(client, sentry.NewScope())
	return &Crash{
		sentryHub: newHub,
	}
}

// AddBreadcrumb records a new breadcrumb.
func (c *Crash) AddBreadcrumb(breadcrumb crash.Breadcrumb) {
	sentryBreadCrumb := &sentry.Breadcrumb{
		Type:      breadcrumb.Type,
		Category:  breadcrumb.Category,
		Message:   breadcrumb.Message,
		Data:      breadcrumb.Data,
		Level:     sentry.Level(breadcrumb.Level),
		Timestamp: breadcrumb.Timestamp,
	}
	c.sentryHub.AddBreadcrumb(sentryBreadCrumb, nil)
}

// Recover calls Sentry's Recover function.
// Returns EventID if successfully, or empty Event ID if Sentry returns nil.
func (c *Crash) Recover(err interface{}) crash.EventID {
	if err == nil {
		return ""
	}
	return convertSentryEventID(c.sentryHub.Recover(err))
}

// Flush waits until the underlying Transport sends any buffered events to the Sentry server, blocking for at most the given timeout.
// It returns false if the timeout was reached. In that case, some events may not have been sent.
// Flush should be called before terminating the program to avoid unintentionally dropping events.
// Do not call Flush indiscriminately after every call to Capture.
func (c *Crash) Flush(timeout time.Duration) bool {
	return c.sentryHub.Flush(timeout)
}

// convertSentryEventID is a helper function to convert EventIDs
func convertSentryEventID(sEID *sentry.EventID) crash.EventID {
	if sEID == nil {
		return ""
	}
	rEID := crash.EventID(*sEID)
	return rEID
}
