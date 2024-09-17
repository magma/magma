/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package clock

import (
	"testing"
	"time"
)

var (
	c clock   = &defaultClock{}
	s sleeper = &defaultSleep{}
)

// Now returns the current time or the time to which it's been set.
func Now() time.Time {
	return c.now()
}

// Sleep either the specified duration or a small, negligible duration.
func Sleep(d time.Duration) {
	s.sleep(d)
}

// Since returns the time elapsed since t, where the current time may have
// been frozen.
func Since(t time.Time) time.Duration {
	return Now().Sub(t)
}

// SetAndFreezeClock will set the value to be returned by Now.
// This should only be called by test code.
func SetAndFreezeClock(t *testing.T, ti time.Time) {
	if t == nil {
		panic("for tests only")
	}
	c = &mockClock{mockTime: ti}
}

// UnfreezeClock will revert clock.Now's behavior to delegating to time.Now.
// This should only be called by test code.
func UnfreezeClock(t *testing.T) {
	r := recover()
	if t == nil {
		panic("for tests only")
	}
	c = &defaultClock{}
	if r != nil {
		panic(r)
	}
}

// SkipSleeps causes time.Sleep to sleep for only a small, negligible duration.
// This should only be called by test code.
func SkipSleeps(t *testing.T) {
	if t == nil {
		panic("for tests only")
	}
	s = &mockSleep{}
}

// ResumeSleeps causes time.Sleep to resume default behavior.
// This should only be called by test code.
func ResumeSleeps(t *testing.T) {
	r := recover()
	if t == nil {
		panic("for tests only")
	}
	s = &defaultSleep{}
	if r != nil {
		panic(r)
	}
}
