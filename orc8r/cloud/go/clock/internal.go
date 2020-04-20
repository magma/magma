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

package clock

import (
	"time"
)

const (
	minSleep = 100 * time.Millisecond
)

// clock is an interface for getting the current time
type clock interface {
	// now returns the current time (or what it's been set to)
	now() time.Time
}

// defaultClock is a clock implementation which wraps time.Now.
type defaultClock struct{}

func (d *defaultClock) now() time.Time {
	return time.Now()
}

// mockClock is a clock implementation which always returns a fixed time.
type mockClock struct {
	mockTime time.Time
}

func (m *mockClock) now() time.Time {
	return m.mockTime
}

type sleeper interface {
	sleep(d time.Duration)
}

// defaultSleep is a clock implementation which wraps time.Sleep.
type defaultSleep struct{}

func (s *defaultSleep) sleep(d time.Duration) {
	time.Sleep(d)
}

// mockSleep is a sleep implementation which always sleeps a small, negligible duration.
type mockSleep struct{}

func (s *mockSleep) sleep(d time.Duration) {
	time.Sleep(minSleep)
}
