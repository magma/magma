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

from typing import Dict
from unittest import TestCase
from unittest.mock import patch

from magma.common.service import MagmaService
from magma.enodebd.data_models.data_model import DataModel
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.exceptions import Tr069Error
from magma.enodebd.state_machines.enb_acs_impl import BasicEnodebAcsStateMachine
from magma.enodebd.state_machines.enb_acs_states import (
    AcsReadMsgResult,
    EnodebAcsState,
    WaitInformState,
    WaitSetParameterValuesState,
)
from magma.enodebd.tests.test_utils.enb_acs_builder import (
    EnodebAcsStateMachineBuilder,
)
from magma.enodebd.tr069 import models


class DummyDataModel(DataModel):
    @classmethod
    def get_parameter(cls, param_name):
        return None

    @classmethod
    def _get_magma_transforms(cls):
        return {}

    @classmethod
    def _get_enb_transforms(cls):
        return {}

    @classmethod
    def get_load_parameters(cls):
        return []

    @classmethod
    def get_num_plmns(cls) -> int:
        return 1

    @classmethod
    def get_parameter_names(cls):
        return []

    @classmethod
    def get_numbered_param_names(cls):
        return {}


class DummyHandler(BasicEnodebAcsStateMachine):

    def __init__(
            self,
            service: MagmaService,
    ) -> None:
        self._state_map = {}
        super().__init__(service)

    def are_invasive_changes_applied(self) -> bool:
        return False

    def _init_state_map(self) -> None:
        self._state_map = {
            'wait_inform': WaitInformState(
                self,
                when_done='wait_empty',
                when_boot='wait_rem',
            ),
        }

    @property
    def state_map(self) -> Dict[str, EnodebAcsState]:
        return self._state_map

    @property
    def disconnected_state_name(self) -> str:
        return 'wait_inform'

    @property
    def unexpected_fault_state_name(self) -> str:
        """ State to handle unexpected Fault messages """
        return ''

    @property
    def device_name(self) -> EnodebDeviceName:
        return "dummy"

    @property
    def config_postprocessor(self):
        pass

    def reboot_asap(self) -> None:
        """
        Send a request to reboot the eNodeB ASAP
        """
        pass

    def is_enodeb_connected(self) -> bool:
        return True

    @property
    def data_model_class(self):
        return DummyDataModel


class EnodebStatusTests(TestCase):

    def _get_acs(self):
        """ Get a dummy ACS statemachine for tests"""
        service = EnodebAcsStateMachineBuilder.build_magma_service()
        return DummyHandler(service)

    @patch(
        'magma.enodebd.state_machines.enb_acs_states'
        '.get_param_values_to_set',
    )
    @patch(
        'magma.enodebd.state_machines.enb_acs_states.get_obj_param_values_to_set')
    def test_wait_set_parameter_values_state(
            self, mock_get_obj_param,
            mock_get_param,
    ):
        """ Test SetParameter return values"""
        mock_get_param.return_value = {}
        mock_get_obj_param.return_value = {}
        test_message_0 = models.SetParameterValuesResponse()
        test_message_0.Status = 0
        test_message_1 = models.SetParameterValuesResponse()
        test_message_1.Status = 1
        # TC-1: return value is 0. No fault
        acs_state = WaitSetParameterValuesState(
            self._get_acs(), 'done',
            'invasive',
        )

        rc = acs_state.read_msg(test_message_0)
        self.assertEqual(type(rc), AcsReadMsgResult)

        # It raises exception if we return 1
        self.assertRaises(
            Tr069Error,
            acs_state.read_msg, test_message_1,
        )

        # It passes if we return 1 and pass the non zero flag
        acs_state = WaitSetParameterValuesState(
            self._get_acs(), 'done',
            'invasive',
            status_non_zero_allowed=True,
        )
        rc = acs_state.read_msg(test_message_1)
        self.assertEqual(type(rc), AcsReadMsgResult)
        rc = acs_state.read_msg(test_message_0)
        self.assertEqual(type(rc), AcsReadMsgResult)
