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

package nprobe

const (
	AttachSuccess        = "attach_success"
	DetachSuccess        = "detach_success"
	SessionCreated       = "session_created"
	SessionUpdated       = "session_updated"
	SessionTerminated    = "session_terminated"
	SessionCreateFailure = "session_create_failure"
	S1SetupSuccess       = "s1_setup_success"

	// ExporterTlsTimeout tls session timeout in milliseconds
	ExporterTlsTimeout = 1000

	// Record types
	IRIRecord    = "IRIRecord"
	NProbeRecord = "NProbeRecord"
)

// GetESStreams returns the list of Intercepted streams
func GetESStreams() []string {
	return []string{"mme", "sessiond"}
}

// GetESEventTypes returns the list of Intercepted events
func GetESEventTypes() []string {
	return []string{
		AttachSuccess, DetachSuccess, SessionCreated, SessionUpdated,
		SessionTerminated, SessionCreateFailure, S1SetupSuccess}
}
