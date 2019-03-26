"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

# pylint: disable=protected-access
from unittest import TestCase
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.tr069 import models
from magma.enodebd.tests.test_utils.tr069_msg_builder import \
    Tr069MessageBuilder
from magma.enodebd.tests.test_utils.enb_acs_builder import \
    EnodebAcsStateMachineBuilder


class BaicellsHandlerTests(TestCase):
    def test_provisioning(self) -> None:
        acs_state_machine = \
            EnodebAcsStateMachineBuilder\
            .build_acs_state_machine(EnodebDeviceName.BAICELLS)

        # Send an Inform message, wait for an InformResponse
        inform_msg = Tr069MessageBuilder.get_inform()
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
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting param values')
        req = Tr069MessageBuilder.get_fault()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting param values')
        req = Tr069MessageBuilder.get_fault()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for read-only params
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting param values')
        req = Tr069MessageBuilder.get_read_only_param_values_response()

        # Send back some typical values
        # And then SM should request regular parameter values
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting param values')

        # Send back typical values for the regular parameters
        req = Tr069MessageBuilder.get_regular_param_values_response()
        resp = acs_state_machine.handle_tr069_message(req)

        # SM will be requesting object parameter values
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting object param vals')

        # Send back some typical values for object parameters
        req = Tr069MessageBuilder.get_object_param_values_response()
        resp = acs_state_machine.handle_tr069_message(req)

        # In this scenario, the ACS and thus state machine will not need
        # to delete or add objects to the eNB configuration.
        # SM should then just be attempting to set parameter values
        self.assertTrue(isinstance(resp, models.SetParameterValues),
                        'State machine should be setting param values')

        # Send back confirmation that the parameters were successfully set
        req = models.SetParameterValuesResponse()
        req.Status = 0
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for read-only params
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting param values')
        req = Tr069MessageBuilder.get_read_only_param_values_response()

        # Send back some typical values
        # And then SM should continue polling the read-only params
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(isinstance(resp, models.DummyInput),
                        'State machine should be ending session')

        # If a different eNB is suddenly plugged in, or the same eNB sends a
        # new Inform, enodebd should be able to handle it.
        # Send an Inform message, wait for an InformResponse
        inform_msg = Tr069MessageBuilder.get_inform()
        resp = acs_state_machine.handle_tr069_message(inform_msg)
        self.assertTrue(isinstance(resp, models.InformResponse),
                        'Should respond with an InformResponse')

        # Send an empty http request to kick off the rest of provisioning
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for an optional parameter, three times
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting param values')
