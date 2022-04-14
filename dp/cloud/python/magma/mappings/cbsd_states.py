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
    Granted,
    Registered,
    Unregistered,
    Unsync,
)
from magma.mappings.types import CbsdStates, GrantStates

cbsd_state_mapping = {
    CbsdStates.UNREGISTERED.value: Unregistered,
    CbsdStates.REGISTERED.value: Registered,
}

grant_state_mapping = {
    GrantStates.GRANTED.value: Granted,
    GrantStates.AUTHORIZED.value: Authorized,
    GrantStates.UNSYNC.value: Unsync,
}
