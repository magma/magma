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
	"github.com/magma/magma/reporter"
	"github.com/pkg/errors"
)

type Reporter struct {
	hub *sentry.Hub
}

func NewReporter(options sentry.ClientOptions) reporter.Reporter {
	client, err := sentry.NewClient(options)
	if err != nil {
		panic(errors.Wrapf(err, "sentry.ClientOptions=%+v", options))
	}
	newHub := sentry.NewHub(client, sentry.NewScope())
	newHub.BindClient(client)

	return &Reporter{
		newHub,
	}
}

func (r *Reporter) CaptureException(exception error) *reporter.EventID {
	return convertSentryEventID(r.hub.CaptureException(exception))
}

func (r *Reporter) AddBreadcrumb(breadcrumb *reporter.Breadcrumb) {
	sentryBreadCrumb := &sentry.Breadcrumb{
		Type:      breadcrumb.Type,
		Category:  breadcrumb.Category,
		Message:   breadcrumb.Message,
		Data:      breadcrumb.Data,
		Level:     sentry.Level(breadcrumb.Level),
		Timestamp: time.Time{},
	}
	r.hub.AddBreadcrumb(sentryBreadCrumb, nil)
}

func (r *Reporter) CaptureMessage(message string) *reporter.EventID {
	return convertSentryEventID(r.hub.CaptureMessage(message))
}

func (r *Reporter) Recover() *reporter.EventID {
	if err := recover(); err != nil {
		return convertSentryEventID(r.hub.Recover(err))
	}
	return nil
}

func (r *Reporter) Flush(timeout time.Duration) bool {
	return r.hub.Flush(timeout)
}

func convertSentryEventID(sEID *sentry.EventID) *reporter.EventID {
	rEID := reporter.EventID(*sEID)
	return &rEID
}
