"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

# pylint: disable=protected-access
from unittest import TestCase
from magma.enodebd.enodeb_status import get_service_status_old, get_all_enb_status
from magma.enodebd.state_machines.enb_acs_manager import StateMachineManager
from magma.enodebd.tests.test_utils.tr069_msg_builder import \
    Tr069MessageBuilder
from magma.enodebd.tests.test_utils.enb_acs_builder import \
    EnodebAcsStateMachineBuilder
from magma.enodebd.tests.test_utils.spyne_builder import \
    get_spyne_context_with_ip


class EnodebStatusTests(TestCase):
    def test_get_service_status_old(self):
        manager = self._get_manager()
        status = get_service_status_old(manager)
        self.assertTrue(status['enodeb_connected'] == '0',
                        'Should report no eNB connected')

        ##### Start session for the first IP #####
        ctx1 = get_spyne_context_with_ip("192.168.60.145")
        # Send an Inform message, wait for an InformResponse
        inform_msg = Tr069MessageBuilder.get_inform('48BF74',
                                                    'BaiBS_RTS_3.1.6',
                                                    '120200002618AGP0001')
        manager.handle_tr069_message(ctx1, inform_msg)
        status = get_service_status_old(manager)
        self.assertTrue(status['enodeb_connected'] == '1',
                        'Should report an eNB as conencted')
        self.assertTrue(status['enodeb_serial'] == '120200002618AGP0001',
                        'eNodeB serial should match the earlier Inform')

    def test_get_enodeb_all_status(self):
        manager = self._get_manager()

        ##### Test Empty #####
        enb_status_by_serial = get_all_enb_status(manager)
        self.assertTrue(enb_status_by_serial == {}, "No eNB connected")

        ##### Start session for the first IP #####
        ctx1 = get_spyne_context_with_ip("192.168.60.145")
        # Send an Inform message, wait for an InformResponse
        inform_msg = Tr069MessageBuilder.get_inform('48BF74',
                                                    'BaiBS_RTS_3.1.6',
                                                    '120200002618AGP0001')
        manager.handle_tr069_message(ctx1, inform_msg)
        enb_status_by_serial = get_all_enb_status(manager)
        enb_status = enb_status_by_serial.get('120200002618AGP0001')
        self.assertTrue(enb_status.enodeb_connected)

    def _get_manager(self) -> StateMachineManager:
        service = EnodebAcsStateMachineBuilder.build_magma_service()
        return StateMachineManager(service)
