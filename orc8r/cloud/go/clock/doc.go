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

// Package clock provides a simple abstraction around the standard time package.
//	- Set and "freeze" the wall clock in test code, with provided wrappers for
//		- time.Now
//		- time.Since
//	- Skip sleeps in test code, with provided wrappers for
//		- time.Sleep
package clock
