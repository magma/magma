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

from dp.protos.active_mode_pb2 import (
    Authorized,
    DeregistrationRequest,
    Granted,
    GrantRequest,
    HeartbeatRequest,
    Registered,
    RegistrationRequest,
    RelinquishmentRequest,
    SpectrumInquiryRequest,
    Unregistered,
    Unsync,
)
from magma.mappings.types import CbsdStates, GrantStates, RequestTypes

cbsd_state_mapping = {
    CbsdStates.UNREGISTERED.value: Unregistered,
    CbsdStates.REGISTERED.value: Registered,
}

grant_state_mapping = {
    GrantStates.GRANTED.value: Granted,
    GrantStates.AUTHORIZED.value: Authorized,
    GrantStates.UNSYNC.value: Unsync,
}

request_type_mapping = {
    RequestTypes.REGISTRATION.value: RegistrationRequest,
    RequestTypes.SPECTRUM_INQUIRY.value: SpectrumInquiryRequest,
    RequestTypes.GRANT.value: GrantRequest,
    RequestTypes.HEARTBEAT.value: HeartbeatRequest,
    RequestTypes.RELINQUISHMENT.value: RelinquishmentRequest,
    RequestTypes.DEREGISTRATION.value: DeregistrationRequest,
}
