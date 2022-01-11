"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

from magma.configuration_controller.response_processor.strategies.response_processing import (
    process_deregistration_response,
    process_grant_response,
    process_heartbeat_response,
    process_registration_response,
    process_relinquishment_response,
    process_spectrum_inquiry_response,
)

# TODO use enum and constants here
processor_strategies = {
    "registrationRequest": {
        "process_responses": process_registration_response,
    },
    "spectrumInquiryRequest": {
        "process_responses": process_spectrum_inquiry_response,
    },
    "grantRequest": {
        "process_responses": process_grant_response,
    },
    "heartbeatRequest": {
        "process_responses": process_heartbeat_response,
    },
    "relinquishmentRequest": {
        "process_responses": process_relinquishment_response,
    },
    "deregistrationRequest": {
        "process_responses": process_deregistration_response,
    },
}
