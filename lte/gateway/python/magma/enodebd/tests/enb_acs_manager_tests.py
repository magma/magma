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

# pylint: disable=protected-access
from unittest import TestCase

from magma.enodebd.state_machines.enb_acs_manager import StateMachineManager
from magma.enodebd.tests.test_utils.enb_acs_builder import (
    EnodebAcsStateMachineBuilder,
)
from magma.enodebd.tests.test_utils.spyne_builder import (
    get_spyne_context_with_ip,
)
from magma.enodebd.tests.test_utils.tr069_msg_builder import Tr069MessageBuilder
from magma.enodebd.tr069 import models


class StateMachineManagerTests(TestCase):
    def test_handle_one_ip(self):
        manager = self._get_manager()

        # Send in an Inform message, and we should get an InformResponse
        ctx = get_spyne_context_with_ip()
        inform = Tr069MessageBuilder.get_inform()
        req = manager.handle_tr069_message(ctx, inform)
        self.assertTrue(
            isinstance(req, models.InformResponse),
            'State machine handler should reply with an '
            'InformResponse',
        )

    def test_serial_not_found(self):
        """
        Test that the SM manager doesn't crash if serial number is not found
        in an Inform message under any expected param path.
        """
        manager = self._get_manager()
        ctx = get_spyne_context_with_ip("192.168.60.145")
        inform_msg = models.Inform(
            DeviceId=models.DeviceIdStruct(
                Manufacturer='Unused',
                OUI='48BF74',
                ProductClass='Unused',
            ),
            Event=models.EventList(EventStruct=[]),
            ParameterList=models.ParameterValueList(
                ParameterValueStruct=[
                    Tr069MessageBuilder.get_parameter_value_struct(
                        name='Device.DeviceInfo.HardwareVersion',
                        val_type='string',
                        data='VER.C',
                    ),
                    Tr069MessageBuilder.get_parameter_value_struct(
                        name='Device.DeviceInfo.ManufacturerOUI',
                        val_type='string',
                        data='48BF74',
                    ),
                    Tr069MessageBuilder.get_parameter_value_struct(
                        name='Device.DeviceInfo.SoftwareVersion',
                        val_type='string',
                        data='BaiBS_RTS_3.1.6',
                    ),
                ],
            ),
        )

        # No exception should be thrown, and we should return an empty response
        resp = manager.handle_tr069_message(ctx, inform_msg)
        self.assertTrue(isinstance(resp, models.DummyInput))

    def test_handle_two_ips(self):
        manager = self._get_manager()
        ctx1 = get_spyne_context_with_ip("192.168.60.145")
        ctx2 = get_spyne_context_with_ip("192.168.60.99")

        ##### Start session for the first IP #####
        # Send an Inform message, wait for an InformResponse
        inform_msg = Tr069MessageBuilder.get_inform(
            '48BF74',
            'BaiBS_RTS_3.1.6',
            '120200002618AGP0001',
        )
        resp1 = manager.handle_tr069_message(ctx1, inform_msg)
        self.assertTrue(
            isinstance(resp1, models.InformResponse),
            'Should respond with an InformResponse',
        )

        # Send an empty http request to kick off the rest of provisioning
        req1 = models.DummyInput()
        resp1 = manager.handle_tr069_message(ctx1, req1)

        # Expect a request for an optional parameter, three times
        self.assertTrue(
            isinstance(resp1, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req1 = Tr069MessageBuilder.get_fault()
        resp1 = manager.handle_tr069_message(ctx1, req1)
        self.assertTrue(
            isinstance(resp1, models.GetParameterValues),
            'State machine should be requesting param values',
        )

        ##### Start session for the second IP #####
        # Send an Inform message, wait for an InformResponse
        inform_msg = Tr069MessageBuilder.get_inform(
            '48BF74',
            'BaiBS_RTS_3.1.6',
            '120200002618AGP0002',
        )
        resp2 = manager.handle_tr069_message(ctx2, inform_msg)
        self.assertTrue(
            isinstance(resp2, models.InformResponse),
            'Should respond with an InformResponse',
        )

        ##### Continue session for the first IP #####
        req1 = Tr069MessageBuilder.get_fault()
        resp1 = manager.handle_tr069_message(ctx1, req1)
        self.assertTrue(
            isinstance(resp1, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req1 = Tr069MessageBuilder.get_fault()
        resp1 = manager.handle_tr069_message(ctx1, req1)
        # Expect a request for read-only params
        self.assertTrue(
            isinstance(resp1, models.GetParameterValues),
            'State machine should be requesting param values',
        )

        ##### Continue session for the second IP #####
        # Send an empty http request to kick off the rest of provisioning
        req2 = models.DummyInput()
        resp2 = manager.handle_tr069_message(ctx2, req2)
        # Expect a request for an optional parameter, three times
        self.assertTrue(
            isinstance(resp2, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req2 = Tr069MessageBuilder.get_fault()
        resp2 = manager.handle_tr069_message(ctx2, req2)
        self.assertTrue(
            isinstance(resp2, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req2 = Tr069MessageBuilder.get_fault()
        resp2 = manager.handle_tr069_message(ctx2, req2)
        self.assertTrue(
            isinstance(resp2, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req2 = Tr069MessageBuilder.get_fault()
        resp2 = manager.handle_tr069_message(ctx2, req2)
        # Expect a request for read-only params
        self.assertTrue(
            isinstance(resp2, models.GetParameterValues),
            'State machine should be requesting param values',
        )

    def test_handle_registered_enb(self):
        """
        When we have a config with eNB registered per serial, we should accept
        TR-069 sessions from any registered eNB, and ereject from unregistered
        eNB devices.
        """
        manager = self._get_manager_multi_enb()
        ip1 = "192.168.60.145"
        ctx1 = get_spyne_context_with_ip(ip1)
        inform_msg = Tr069MessageBuilder.get_inform(
            '48BF74',
            'BaiBS_RTS_3.1.6',
            '120200002618AGP0003',
        )
        resp1 = manager.handle_tr069_message(ctx1, inform_msg)
        self.assertTrue(
            isinstance(resp1, models.InformResponse),
            'Should respond with an InformResponse',
        )

        ip2 = "192.168.60.146"
        ctx2 = get_spyne_context_with_ip(ip2)
        inform_msg = Tr069MessageBuilder.get_inform(
            '48BF74',
            'BaiBS_RTS_3.1.6',
            'unregistered_serial',
        )

        resp2 = manager.handle_tr069_message(ctx2, inform_msg)
        self.assertTrue(
            isinstance(resp2, models.DummyInput),
            'Should respond with an empty HTTP response',
        )

    def test_ip_change(self) -> None:
        manager = self._get_manager()

        # Send an Inform
        ip1 = "192.168.60.145"
        ctx1 = get_spyne_context_with_ip(ip1)
        inform_msg = Tr069MessageBuilder.get_inform(
            '48BF74',
            'BaiBS_RTS_3.1.6',
            '120200002618AGP0003',
        )
        resp1 = manager.handle_tr069_message(ctx1, inform_msg)
        self.assertTrue(
            isinstance(resp1, models.InformResponse),
            'Should respond with an InformResponse',
        )
        handler1 = manager.get_handler_by_ip(ip1)

        # Send an Inform from the same serial, but different IP
        ip2 = "192.168.60.99"
        ctx2 = get_spyne_context_with_ip(ip2)
        inform_msg = Tr069MessageBuilder.get_inform(
            '48BF74',
            'BaiBS_RTS_3.1.6',
            '120200002618AGP0003',
        )
        resp2 = manager.handle_tr069_message(ctx2, inform_msg)
        self.assertTrue(
            isinstance(resp2, models.InformResponse),
            'Should respond with an InformResponse',
        )
        handler2 = manager.get_handler_by_ip(ip2)

        # Now check that the serial is associated with the second ip
        self.assertTrue(
            (handler1 is handler2),
            'After an IP switch, the manager should have moved '
            'the handler to a new IP',
        )

    def test_serial_change(self) -> None:
        manager = self._get_manager()
        ip = "192.168.60.145"

        # Send an Inform
        ctx1 = get_spyne_context_with_ip(ip)
        inform_msg = Tr069MessageBuilder.get_inform(
            '48BF74',
            'BaiBS_RTS_3.1.6',
            '120200002618AGP0001',
        )
        resp1 = manager.handle_tr069_message(ctx1, inform_msg)
        self.assertTrue(
            isinstance(resp1, models.InformResponse),
            'Should respond with an InformResponse',
        )
        handler1 = manager.get_handler_by_ip(ip)

        # Send an Inform from the same serial, but different IP
        ctx2 = get_spyne_context_with_ip(ip)
        inform_msg = Tr069MessageBuilder.get_inform(
            '48BF74',
            'BaiBS_RTS_3.1.6',
            '120200002618AGP0002',
        )
        resp2 = manager.handle_tr069_message(ctx2, inform_msg)
        self.assertTrue(
            isinstance(resp2, models.InformResponse),
            'Should respond with an InformResponse',
        )
        handler2 = manager.get_handler_by_ip(ip)

        # Now check that the serial is associated with the second ip
        self.assertTrue(
            (handler1 is not handler2),
            'After an IP switch, the manager should have moved '
            'the handler to a new IP',
        )

    def test_inform_from_baicells_qafb(self) -> None:
        manager = self._get_manager()
        ip = "192.168.60.145"

        # Send an Inform
        ctx1 = get_spyne_context_with_ip(ip)
        inform_msg = Tr069MessageBuilder.get_qafb_inform(
            '48BF74',
            'BaiBS_QAFB_v1234',
            '120200002618AGP0001',
        )
        resp1 = manager.handle_tr069_message(ctx1, inform_msg)
        self.assertTrue(
            isinstance(resp1, models.InformResponse),
            'Should respond with an InformResponse',
        )

    def test_inform_from_unrecognized(self) -> None:
        manager = self._get_manager()
        ip = "192.168.60.145"

        # Send an Inform
        ctx1 = get_spyne_context_with_ip(ip)
        inform_msg = Tr069MessageBuilder.get_qafb_inform(
            '48BF74',
            'Unrecognized device',
            '120200002618AGP0001',
        )
        resp1 = manager.handle_tr069_message(ctx1, inform_msg)
        self.assertTrue(
            isinstance(resp1, models.DummyInput),
            'Should end provisioninng session with empty response',
        )

    def _get_manager(self) -> StateMachineManager:
        service = EnodebAcsStateMachineBuilder.build_magma_service()
        return StateMachineManager(service)

    def _get_manager_multi_enb(self) -> StateMachineManager:
        service = EnodebAcsStateMachineBuilder.build_multi_enb_magma_service()
        return StateMachineManager(service)
