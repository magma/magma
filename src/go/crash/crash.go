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

// Package crash provides a generic crash reporting abstraction for Magma.
package crash

import (
	"time"
)

// EventID is a hexadecimal string representing a unique uuid4 for an Event.
// An EventID must be 32 characters long, lowercase and not have any dashes.
type EventID string

// Level marks the severity of the event.
type Level string

// Describes the severity of the event.
// Log descriptions are based on sentry-go log levels defined here
// https://github.com/getsentry/sentry-go/blob/master/interfaces.go#L23
const (
	LevelDebug   Level = "debug"
	LevelInfo    Level = "info"
	LevelWarning Level = "warning"
	LevelError   Level = "error"
	LevelFatal   Level = "fatal"
)

// Breadcrumb specifies an application event that occurred before a reported event.
// An event may contain one or more breadcrumbs.
type Breadcrumb struct {
	Type      string
	Category  string
	Message   string
	Data      map[string]interface{}
	Level     Level
	Timestamp time.Time
}

// Crash provides an interface for interacting with crash reporters.
type Crash interface {
	AddBreadcrumb(breadcrumb Breadcrumb)
	Recover(err interface{}) EventID
	Flush(timeout time.Duration) bool
}
