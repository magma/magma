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

// Package interceptd provides the lawful interception service.
// Interceptd service primarily aggregates all the events based on
// intercept configuration reported by points of interception based
// on li configuration, encode records and export them to LIMS.
package interceptd

// ServiceName provides the service name
const ServiceName = "interceptd"
