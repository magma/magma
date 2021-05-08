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
	ESStreamMME      = "mme"
	ESStreamSessionD = "sessiond"

	AttachSuccess        = "attach_success"
	DetachSuccess        = "detach_success"
	SessionCreated       = "session_created"
	SessionUpdated       = "session_updated"
	SessionTerminated    = "session_terminated"
	S1SetupSuccess       = "s1_setup_success"
	SessionCreateFailure = "session_create_failure"
)

// GetESStreams returns the list of Intercepted streams
func GetESStreams() []string {
	return []string{ESStreamMME, ESStreamSessionD}
}

// GetESEventTypes returns the list of Intercepted events
func GetESEventTypes() []string {
	return []string{
		AttachSuccess,
		DetachSuccess,
		SessionCreated,
		SessionUpdated,
		SessionTerminated,
		SessionCreateFailure,
		S1SetupSuccess,
	}
}
