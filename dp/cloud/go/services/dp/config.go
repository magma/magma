/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dp

type Config struct {
	// TODO cleanup config (common fields, separate packages, etc...)
	DpBackend            *BackendConfig `yaml:"dp_backend"`
	ActiveModeController *AmcConfig     `yaml:"active_mode_controller"`
}

type BackendConfig struct {
	CbsdInactivityIntervalSec int    `yaml:"cbsd_inactivity_interval_sec"`
	LogConsumerUrl            string `yaml:"log_consumer_url"`
}

type AmcConfig struct {
	DialTimeoutSec               int    `yaml:"dial_timeout_sec"`
	HeartbeatSendTimeoutSec      int    `yaml:"heartbeat_send_timeout_sec"`
	RequestTimeoutSec            int    `yaml:"request_timeout_sec"`
	RequestProcessingIntervalSec int    `yaml:"request_processing_interval_sec"`
	PollingIntervalSec           int    `yaml:"polling_interval"` // TODO add sec to deployment scripts
	GrpcService                  string `yaml:"grpc_service"`
	GrpcPort                     int    `yaml:"grpc_port"`
	CbsdInactivityTimeoutSec     int    `yaml:"cbsd_inactivity_interval_sec"` // TODO temporary fix to make integration tests pass
}
