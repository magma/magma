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
	"strings"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/magma/magma/reporter"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewReporter(t *testing.T) {
	options := sentry.ClientOptions{
		Dsn: "http://whatever@really.com/1337",
	}
	r := NewReporter(options)
	reporter, ok := r.(*Reporter)
	assert.True(
		t, ok, "NewReporter must return *sentry.Reporter, have=%+v", r)
	assert.NotNil(t, reporter.hub)
	assert.Equal(t, options.Dsn, reporter.hub.Client().Options().Dsn)
}

func TestNewReporter_Error(t *testing.T) {
	options := sentry.ClientOptions{
		Dsn: "\"bad_input",
	}

	// we expect error from sentry.NewReporter() to panic
	defer func() {
		r := recover()
		if r == nil {
			t.Error("The code did not panic")
			return
		}
		err, ok := r.(error)
		if !ok {
			t.Errorf("panic value not err, r=%+v", r)
			return
		}
		assert.True(t, strings.Contains(err.Error(), "DsnParseError:"))
	}()
	NewReporter(options)

}

func TestReporter_CaptureMessage(t *testing.T) {
	msg := "test message"

	options := sentry.ClientOptions{}
	r := NewReporter(options)

	r.CaptureMessage(msg)
}

func TestReporter_CaptureMessage_Error(t *testing.T) {
	msg := "test message"

	options := sentry.ClientOptions{}
	r := NewReporter(options)
	reporter, _ := r.(*Reporter)
	reporter.hub = nil
	// we expect error from sentry.CaptureMessage() to panic
	defer func() {
		r := recover()
		if r == nil {
			t.Error("The code did not panic")
			return
		}
		err, ok := r.(error)
		if !ok {
			t.Errorf("panic value not err, r=%+v", r)
			return
		}
		assert.True(t, strings.Contains(err.Error(), "invalid memory address or nil pointer dereference"))
	}()

	r.CaptureMessage(msg)

}

func TestReporter_CaptureException(t *testing.T) {
	exception := errors.New("error")

	options := sentry.ClientOptions{}
	r := NewReporter(options)

	r.CaptureException(exception)
}

func TestReporter_CaptureException_Error(t *testing.T) {
	exception := errors.New("error")

	options := sentry.ClientOptions{}
	r := NewReporter(options)
	reporter, _ := r.(*Reporter)
	reporter.hub = nil
	// we expect error from sentry.CaptureException() to panic
	defer func() {
		r := recover()
		if r == nil {
			t.Error("The code did not panic")
			return
		}
		err, ok := r.(error)
		if !ok {
			t.Errorf("panic value not err, r=%+v", r)
			return
		}
		assert.True(t, strings.Contains(err.Error(), "invalid memory address or nil pointer dereference"))
	}()

	r.CaptureException(exception)
}

func TestReporter_AddBreadcrumb(t *testing.T) {
	breadcrumb := &reporter.Breadcrumb{}

	options := sentry.ClientOptions{}
	r := NewReporter(options)

	r.AddBreadcrumb(breadcrumb)
}

func TestReporter_AddBreadcrumb_Error(t *testing.T) {
	breadcrumb := &reporter.Breadcrumb{}

	options := sentry.ClientOptions{}
	r := NewReporter(options)
	reporter, _ := r.(*Reporter)
	reporter.hub = nil
	// we expect error from sentry.AddBreadcrumb() to panic
	defer func() {
		r := recover()
		if r == nil {
			t.Error("The code did not panic")
			return
		}
		err, ok := r.(error)
		if !ok {
			t.Errorf("panic value not err, r=%+v", r)
			return
		}
		assert.True(t, strings.Contains(err.Error(), "invalid memory address or nil pointer dereference"))
	}()

	r.AddBreadcrumb(breadcrumb)
}

func TestReporter_Recover(t *testing.T) {
	options := sentry.ClientOptions{}
	r := NewReporter(options)

	r.Recover()
}

func TestRepoerter_Flush(t *testing.T) {
	timeout := time.Second
	options := sentry.ClientOptions{}
	r := NewReporter(options)

	r.Flush(timeout)
}
