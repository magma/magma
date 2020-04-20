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
// package service_manager defines and implements API for service management
package service_manager

import "os"

type ServiceState int

const (
	_ ServiceState = iota
	Active
	Activating
	Deactivating
	Inactive
	Failed
	Unknown
	Error
)

// ServiceController defines service controller API for service manager providers
type ServiceController interface {
	// Name returns the type of the init system used by the GW, it should match magmad.yml "init_system" value
	Name() string
	// Start starts service and returns error if unsuccessful
	Start(service string) error
	// Stop stops service and returns error if unsuccessful
	Stop(service string) error
	// Restart restarts service and returns error if unsuccessful
	Restart(service string) error
	// GetState returns the given service state or error if unsuccessful
	GetState(service string) (ServiceState, error)

	// TailLogs executes command to start tailing service logs and returns string chan to receive log strings
	// closing the chan will terminate tailing
	TailLogs(service string) (chan string, *os.Process, error)
}
