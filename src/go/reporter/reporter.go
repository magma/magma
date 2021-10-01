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

package reporter

import (
	"time"
)

// EventID is a hexadecimal string representing a unique uuid4 for an Event.
// An EventID must be 32 characters long, lowercase and not have any dashes.
type EventID string

// Level marks the severity of the event.
type Level string

// Breadcrumb specifies an application event that occurred before a reported event.
// An event may contain one or more breadcrumbs.
type Breadcrumb struct {
	Type      string                 `json:"type,omitempty"`
	Category  string                 `json:"category,omitempty"`
	Message   string                 `json:"message,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Level     Level                  `json:"level,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

type Reporter interface {
	AddBreadcrumb(breadcrumb *Breadcrumb)
	CaptureException(exception error) *EventID
	CaptureMessage(message string) *EventID

	Recover() *EventID
	Flush(timeout time.Duration) bool
}
