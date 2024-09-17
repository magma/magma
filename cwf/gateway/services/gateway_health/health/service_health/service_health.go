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

package service_health

// ServiceHealth defines an interface to fetch unhealthy services and enable
// functionality necessary for promotion/demotions of the gateway.
type ServiceHealth interface {
	// GetUnhealthyServices return a list of services found to be in an
	// unhealthy state.
	GetUnhealthyServices() ([]string, error)

	// Restart restarts the provided service.
	Restart(service string) error

	// Stop stops the provided service.
	Stop(service string) error
}
