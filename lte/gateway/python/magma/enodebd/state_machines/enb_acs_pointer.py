"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

from magma.enodebd.state_machines.enb_acs import EnodebAcsStateMachine


class StateMachinePointer:
    """
    This is a hack to deal with the possibility that the specified data model
    doesn't match the eNB device enodebd ends up connecting to.

    When the data model doesn't match, the state machine is replaced with one
    that matches the data model.
    """

    def __init__(self, acs_state_machine: EnodebAcsStateMachine):
        self._acs_state_machine = acs_state_machine

    @property
    def state_machine(self):
        return self._acs_state_machine

    @state_machine.setter
    def state_machine(
        self,
        acs_state_machine: EnodebAcsStateMachine,
    ) -> None:
        self._acs_state_machine = acs_state_machine
