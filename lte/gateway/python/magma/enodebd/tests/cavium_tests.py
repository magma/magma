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

from magma.enodebd.data_models.data_model_parameters import ParameterName
# pylint: disable=protected-access
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.tests.test_utils.enb_acs_builder import (
    EnodebAcsStateMachineBuilder,
)
from magma.enodebd.tests.test_utils.enodeb_handler import EnodebHandlerTestCase
from magma.enodebd.tests.test_utils.tr069_msg_builder import Tr069MessageBuilder
from magma.enodebd.tr069 import models


class CaviumHandlerTests(EnodebHandlerTestCase):
    def test_count_plmns_less(self) -> None:
        """
        Tests the Cavium provisioning up to GetObjectParameters.

        In particular tests when the eNB reports NUM_PLMNS less
        than actually listed. The eNB says there are no PLMNs
        defined when actually there are two.

        Verifies that the number of PLMNs is correctly accounted.
        """
        acs_state_machine = \
            EnodebAcsStateMachineBuilder \
            .build_acs_state_machine(EnodebDeviceName.CAVIUM)

        # Send an Inform message
        inform_msg = Tr069MessageBuilder.get_inform(
            '000FB7',
            'OC-LTE',
            '120200002618AGP0003',
            ['1 BOOT'],
        )
        resp = acs_state_machine.handle_tr069_message(inform_msg)

        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'Should respond with an InformResponse',
        )

        # Send an empty http request to kick off the rest of provisioning
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values: %s' % resp,
        )

        # Transient config response and request for parameter values
        req = Tr069MessageBuilder.get_read_only_param_values_response()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values: %s' % resp,
        )

        # Send back typical values for the regular parameters
        req = Tr069MessageBuilder.get_cavium_param_values_response(num_plmns=0)
        resp = acs_state_machine.handle_tr069_message(req)

        # SM will be requesting object parameter values
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting object param vals',
        )

        # Send back some object parameters with TWO plmns
        req = Tr069MessageBuilder.get_cavium_object_param_values_response(
            num_plmns=2,
        )
        resp = acs_state_machine.handle_tr069_message(req)

        # In this scenario, the ACS and thus state machine will not need
        # to delete or add objects to the eNB configuration.
        # SM should then just be attempting to set parameter values
        self.assertTrue(
            isinstance(resp, models.SetParameterValues),
            'State machine should be setting param values',
        )

        # Number of PLMNs should reflect object count
        num_plmns_cur = \
            acs_state_machine \
            .device_cfg.get_parameter(ParameterName.NUM_PLMNS)
        self.assertEqual(num_plmns_cur, 2)

    def test_count_plmns_more_defined(self) -> None:
        """
        Tests the Cavium provisioning up to GetObjectParameters.

        In particular tests when the eNB has more PLMNs than is
        currently defined in our data model (NUM_PLMNS_IN_CONFIG)
        """
        acs_state_machine = \
            EnodebAcsStateMachineBuilder \
            .build_acs_state_machine(EnodebDeviceName.CAVIUM)

        # Send an Inform message
        inform_msg = Tr069MessageBuilder.get_inform(
            '000FB7',
            'OC-LTE',
            '120200002618AGP0003',
            ['1 BOOT'],
        )
        resp = acs_state_machine.handle_tr069_message(inform_msg)

        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'Should respond with an InformResponse',
        )

        # Send an empty http request to kick off the rest of provisioning
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values: %s' % resp,
        )

        # Transient config response and request for parameter values
        req = Tr069MessageBuilder.get_read_only_param_values_response()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values: %s' % resp,
        )

        # Send back regular parameters, and some absurd number of PLMNS
        req = Tr069MessageBuilder.get_cavium_param_values_response(
            num_plmns=100,
        )
        resp = acs_state_machine.handle_tr069_message(req)

        # SM will be requesting object parameter values
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting object param vals',
        )

        # Send back some object parameters with an absurd number of PLMNs
        req = Tr069MessageBuilder.get_cavium_object_param_values_response(
            num_plmns=100,
        )
        resp = acs_state_machine.handle_tr069_message(req)

        # In this scenario, the ACS and thus state machine will not need
        # to delete or add objects to the eNB configuration.
        # SM should then just be attempting to set parameter values
        self.assertTrue(
            isinstance(resp, models.SetParameterValues),
            'State machine should be setting param values',
        )

        # Number of PLMNs should reflect data model
        num_plmns_cur = \
            acs_state_machine \
            .device_cfg.get_parameter(ParameterName.NUM_PLMNS)
        self.assertEqual(num_plmns_cur, 6)
