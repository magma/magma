"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

# pylint: disable=protected-access
from unittest import TestCase, mock
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.tr069 import models
from magma.enodebd.tests.test_utils.tr069_msg_builder import \
    Tr069MessageBuilder
from magma.enodebd.tests.test_utils.enb_acs_builder import \
    EnodebAcsStateMachineBuilder
from magma.enodebd.tests.test_utils.mock_functions import \
    mock_get_ip_from_if, GET_IP_FROM_IF_PATH


class BaicellsQAFBHandlerTests(TestCase):
    def test_manual_reboot(self) -> None:
        """
        Test a scenario where a Magma user goes through the enodebd CLI to
        reboot the Baicells eNodeB.

        This checks the scenario where the command is not sent in the middle
        of a TR-069 provisioning session.
        """
        acs_state_machine = \
            EnodebAcsStateMachineBuilder \
                .build_acs_state_machine(EnodebDeviceName.BAICELLS_QAFB)

        # User uses the CLI tool to get eNodeB to reboot
        acs_state_machine.reboot_asap()

        # And now the Inform message arrives from the eNodeB
        inform_msg = \
            Tr069MessageBuilder.get_qafb_inform('48BF74',
                                                'BaiBS_QAFBv123',
                                                '1202000181186TB0006',
                                                ['2 PERIODIC'])
        resp = acs_state_machine.handle_tr069_message(inform_msg)
        self.assertTrue(isinstance(resp, models.InformResponse),
                        'In reboot sequence, state machine should still '
                        'respond to an Inform with InformResponse.')
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(isinstance(resp, models.Reboot),
                        'In reboot sequence, state machine should send a '
                        'Reboot message.')
        req = Tr069MessageBuilder.get_reboot_response()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(isinstance(resp, models.DummyInput),
                        'State machine should end TR-069 session after '
                        'receiving a RebootResponse')

    def test_manual_reboot_during_provisioning(self) -> None:
        """
        Test a scenario where a Magma user goes through the enodebd CLI to
        reboot the Baicells eNodeB.

        This checks the scenario where the command is sent in the middle
        of a TR-069 provisioning session.
        """
        acs_state_machine = \
            EnodebAcsStateMachineBuilder \
                .build_acs_state_machine(EnodebDeviceName.BAICELLS_QAFB)

        # Send an Inform message, wait for an InformResponse
        inform_msg = \
            Tr069MessageBuilder.get_qafb_inform('48BF74',
                                                'BaiBS_QAFBv123',
                                                '1202000181186TB0006',
                                                ['2 PERIODIC'])
        resp = acs_state_machine.handle_tr069_message(inform_msg)
        self.assertTrue(isinstance(resp, models.InformResponse),
                        'Should respond with an InformResponse')

        # Send an empty http request to kick off the rest of provisioning
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for an optional parameter, three times
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting param values')
        req = Tr069MessageBuilder.get_fault()

        # User uses the CLI tool to get eNodeB to reboot
        acs_state_machine.reboot_asap()

        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(isinstance(resp, models.Reboot),
                        'In reboot sequence, state machine should send a '
                        'Reboot message.')
        req = Tr069MessageBuilder.get_reboot_response()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(isinstance(resp, models.DummyInput),
                        'State machine should end TR-069 session after '
                        'receiving a RebootResponse')

    @mock.patch(GET_IP_FROM_IF_PATH, side_effect=mock_get_ip_from_if)
    def test_provision(self, _mock_func) -> None:
        acs_state_machine = \
            EnodebAcsStateMachineBuilder \
                .build_acs_state_machine(EnodebDeviceName.BAICELLS_QAFB)

        # Send an Inform message, wait for an InformResponse
        inform_msg = \
            Tr069MessageBuilder.get_qafb_inform('48BF74',
                                                'BaiBS_QAFBv123',
                                                '1202000181186TB0006',
                                                ['2 PERIODIC'])
        resp = acs_state_machine.handle_tr069_message(inform_msg)
        self.assertTrue(isinstance(resp, models.InformResponse),
                        'Should respond with an InformResponse')

        # Send an empty http request to kick off the rest of provisioning
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for read-only params
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting param values')
        req = Tr069MessageBuilder.get_qafb_read_only_param_values_response()

        # Send back some typical values
        # And then SM should request regular parameter values
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting param values')

        # Send back typical values for the regular parameters
        req = Tr069MessageBuilder.\
            get_qafb_regular_param_values_response(admin_state=False,
                                              earfcndl=39150)
        resp = acs_state_machine.handle_tr069_message(req)

        # SM will be requesting object parameter values
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting object param vals')

        # Send back some typical values for object parameters
        req = Tr069MessageBuilder.get_qafb_object_param_values_response()
        resp = acs_state_machine.handle_tr069_message(req)

        self.assertTrue(isinstance(resp, models.AddObject),
                        'State machine should be adding objects')

    def test_get_rpc_methods_cold(self) -> None:
        """
        Test the scenario where:
        - enodeB just booted
        - enodeB is cold and has no state of ACS RPCMethods
        - Simulate the enodeB performing the initial Inform and
          the call for the GetRPCMethods, and the subsequent Empty
          response for provisioning
          finishing on the Baicells eNodeB

        Verifies that the ACS will continue into provisioning
        """
        acs_state_machine = \
            EnodebAcsStateMachineBuilder\
                .build_acs_state_machine(EnodebDeviceName.BAICELLS_QAFB)

        # Send an Inform message, wait for an InformResponse
        inform_msg = \
            Tr069MessageBuilder.get_inform('48BF74',
                                           'BaiBS_QAFBv123',
                                           '120200002618AGP0003',
                                           ['1 BOOT'])
        resp = acs_state_machine.handle_tr069_message(inform_msg)
        self.assertTrue(isinstance(resp, models.InformResponse),
                        'Should respond with an InformResponse')

        # Send GetRPCMethods
        req = models.GetRPCMethods()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(isinstance(resp, models.GetRPCMethodsResponse),
                        'State machine should be sending RPC methods')

        # Send an empty http request to kick off the rest of provisioning
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for an optional parameter
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting param values')
