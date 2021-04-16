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
from magma.enodebd.data_models.data_model_parameters import ParameterName
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.tests.test_utils.enb_acs_builder import (
    EnodebAcsStateMachineBuilder,
)
from magma.enodebd.tests.test_utils.enodeb_handler import EnodebHandlerTestCase
from magma.enodebd.tests.test_utils.tr069_msg_builder import Tr069MessageBuilder
from magma.enodebd.tr069 import models


class BaicellsHandlerTests(EnodebHandlerTestCase):
    def test_initial_enb_bootup(self) -> None:
        """
        Baicells does not support configuration during initial bootup of
        eNB device. This is because it is in a REM process, and we just need
        to wait for this process to finish, ~10 minutes. Attempting to
        configure the device during this period will cause undefined
        behavior.
        As a result of this, end any provisoning sessions, which we can do
        by just sending empty HTTP responses, not even using an
        InformResponse.
        """
        acs_state_machine = \
            EnodebAcsStateMachineBuilder \
            .build_acs_state_machine(EnodebDeviceName.BAICELLS)

        # Send an Inform message
        inform_msg = Tr069MessageBuilder.get_inform(
            '48BF74',
            'BaiBS_RTS_3.1.6',
            '120200002618AGP0003',
            ['1 BOOT'],
        )
        resp = acs_state_machine.handle_tr069_message(inform_msg)

        self.assertTrue(
            isinstance(resp, models.DummyInput),
            'Should respond with an InformResponse',
        )

    def test_manual_reboot(self) -> None:
        """
        Test a scenario where a Magma user goes through the enodebd CLI to
        reboot the Baicells eNodeB.

        This checks the scenario where the command is not sent in the middle
        of a TR-069 provisioning session.
        """
        acs_state_machine = \
            EnodebAcsStateMachineBuilder \
            .build_acs_state_machine(EnodebDeviceName.BAICELLS)

        # User uses the CLI tool to get eNodeB to reboot
        acs_state_machine.reboot_asap()

        # And now the Inform message arrives from the eNodeB
        inform_msg = Tr069MessageBuilder.get_inform(
            '48BF74',
            'BaiBS_RTS_3.1.6',
            '120200002618AGP0003',
            ['2 PERIODIC'],
        )
        resp = acs_state_machine.handle_tr069_message(inform_msg)
        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'In reboot sequence, state machine should still '
            'respond to an Inform with InformResponse.',
        )
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.Reboot),
            'In reboot sequence, state machine should send a '
            'Reboot message.',
        )
        req = Tr069MessageBuilder.get_reboot_response()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.DummyInput),
            'State machine should end TR-069 session after '
            'receiving a RebootResponse',
        )

    def test_gps_coords(self) -> None:
        """ Check GPS coordinates are processed and stored correctly """
        acs_state_machine = \
            EnodebAcsStateMachineBuilder \
            .build_acs_state_machine(EnodebDeviceName.BAICELLS)

        # Send an Inform message, wait for an InformResponse
        inform_msg = Tr069MessageBuilder.get_inform(
            '48BF74',
            'BaiBS_RTS_3.1.6',
            '120200002618AGP0003',
            ['2 PERIODIC'],
        )
        resp = acs_state_machine.handle_tr069_message(inform_msg)
        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'Should respond with an InformResponse',
        )

        # Send an empty http request to kick off the rest of provisioning
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for an optional parameter, three times
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = models.GetParameterValuesResponse()
        param_val_list = [
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.X_BAICELLS_COM_GpsSyncEnable',
                val_type='boolean',
                data='true',
            ),
        ]
        req.ParameterList = models.ParameterValueList()
        req.ParameterList.ParameterValueStruct = param_val_list
        resp = acs_state_machine.handle_tr069_message(req)

        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = models.GetParameterValuesResponse()
        param_val_list = [
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.FAP.GPS.LockedLatitude',
                val_type='int',
                data='37483629',
            ),
        ]
        req.ParameterList = models.ParameterValueList()
        req.ParameterList.ParameterValueStruct = param_val_list
        resp = acs_state_machine.handle_tr069_message(req)

        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = models.GetParameterValuesResponse()
        param_val_list = [
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.FAP.GPS.LockedLongitude',
                val_type='int',
                data='-122150583',
            ),
        ]
        req.ParameterList = models.ParameterValueList()
        req.ParameterList.ParameterValueStruct = param_val_list
        acs_state_machine.handle_tr069_message(req)

        gps_long = acs_state_machine.get_parameter(ParameterName.GPS_LONG)
        gps_lat = acs_state_machine.get_parameter(ParameterName.GPS_LAT)

        self.assertTrue(gps_long == '-122.150583', 'Should be valid longitude')
        self.assertTrue(gps_lat == '37.483629', 'Should be valid latitude')

    def test_manual_reboot_during_provisioning(self) -> None:
        """
        Test a scenario where a Magma user goes through the enodebd CLI to
        reboot the Baicells eNodeB.

        This checks the scenario where the command is sent in the middle
        of a TR-069 provisioning session.
        """
        acs_state_machine = \
            EnodebAcsStateMachineBuilder \
            .build_acs_state_machine(EnodebDeviceName.BAICELLS)

        # Send an Inform message, wait for an InformResponse
        inform_msg = Tr069MessageBuilder.get_inform(
            '48BF74',
            'BaiBS_RTS_3.1.6',
            '120200002618AGP0003',
            ['2 PERIODIC'],
        )
        resp = acs_state_machine.handle_tr069_message(inform_msg)
        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'Should respond with an InformResponse',
        )

        # Send an empty http request to kick off the rest of provisioning
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for an optional parameter, three times
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_fault()

        # User uses the CLI tool to get eNodeB to reboot
        acs_state_machine.reboot_asap()

        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.Reboot),
            'In reboot sequence, state machine should send a '
            'Reboot message.',
        )
        req = Tr069MessageBuilder.get_reboot_response()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.DummyInput),
            'State machine should end TR-069 session after '
            'receiving a RebootResponse',
        )

    def test_missing_param_during_provisioning(self) -> None:
        """
        Test the scenario where:
        - enodebd is configuring the eNodeB
        - eNB does not send all parameters due to bug
        """
        acs_state_machine = \
            EnodebAcsStateMachineBuilder \
            .build_acs_state_machine(EnodebDeviceName.BAICELLS)

        # Send an Inform message, wait for an InformResponse
        inform_msg = Tr069MessageBuilder.get_inform()
        resp = acs_state_machine.handle_tr069_message(inform_msg)
        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'Should respond with an InformResponse',
        )

        # Send an empty http request to kick off the rest of provisioning
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for an optional parameter, three times
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_fault()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_fault()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_fault()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for read-only params
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_read_only_param_values_response()

        # Send back some typical values
        # And then SM should request regular parameter values
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )

        # Send back typical values for the regular parameters
        # Pretend that here the NumPLMNs was not sent because of a Baicells bug
        req = Tr069MessageBuilder.\
            get_regular_param_values_response(
                admin_state=False,
                earfcndl=39150,
                exclude_num_plmns=True,
            )
        resp = acs_state_machine.handle_tr069_message(req)

        # The state machine will fail and go into an error state.
        # It will send an empty http response to end the session.
        # Regularly, the SM should be using info on the number
        # of PLMNs to figure out which object parameter values
        # to fetch.
        self.assertTrue(
            isinstance(resp, models.DummyInput),
            'State machine should be ending session',
        )

    def test_provision_multi_without_invasive_changes(self) -> None:
        """
        Test the scenario where:
        - eNodeB has already been powered for 10 minutes without configuration
        - Setting parameters which are 'non-invasive' on the eNodeB
        - Using enodebd mconfig which has old style config with addition
          of eNodeB config tied to a serial number

        'Invasive' parameters are those which require special behavior to apply
        the changes for the eNodeB.
        """
        acs_state_machine = \
            EnodebAcsStateMachineBuilder \
            .build_multi_enb_acs_state_machine(EnodebDeviceName.BAICELLS)

        # Send an Inform message, wait for an InformResponse
        inform_msg = Tr069MessageBuilder.get_inform()
        resp = acs_state_machine.handle_tr069_message(inform_msg)
        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'Should respond with an InformResponse',
        )

        # Send an empty http request to kick off the rest of provisioning
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for an optional parameter, three times
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_fault()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_fault()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_fault()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for read-only params
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_read_only_param_values_response()

        # Send back some typical values
        # And then SM should request regular parameter values
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )

        # Send back typical values for the regular parameters
        req = Tr069MessageBuilder.\
            get_regular_param_values_response(
                admin_state=False,
                earfcndl=39150,
            )
        resp = acs_state_machine.handle_tr069_message(req)

        # SM will be requesting object parameter values
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting object param vals',
        )

        # Send back some typical values for object parameters
        req = Tr069MessageBuilder.get_object_param_values_response()
        resp = acs_state_machine.handle_tr069_message(req)

        # In this scenario, the ACS and thus state machine will not need
        # to delete or add objects to the eNB configuration.
        # SM should then just be attempting to set parameter values
        self.assertTrue(
            isinstance(resp, models.SetParameterValues),
            'State machine should be setting param values',
        )

        isEnablingAdminState = False
        param = 'Device.Services.FAPService.1.FAPControl.LTE.AdminState'
        for name_value in resp.ParameterList.ParameterValueStruct:
            if name_value.Name == param:
                isEnablingAdminState = True
        self.assertTrue(
            isEnablingAdminState,
            'eNB config is set to enable transmit, '
            'while old enodebd config does not '
            'enable transmit. Use eNB config.',
        )

    def test_provision_without_invasive_changes(self) -> None:
        """
        Test the scenario where:
        - eNodeB has already been powered for 10 minutes without configuration
        - Setting parameters which are 'non-invasive' on the eNodeB

        'Invasive' parameters are those which require special behavior to apply
        the changes for the eNodeB.
        """
        acs_state_machine = \
            EnodebAcsStateMachineBuilder \
            .build_acs_state_machine(EnodebDeviceName.BAICELLS)

        # Send an Inform message, wait for an InformResponse
        inform_msg = Tr069MessageBuilder.get_inform()
        resp = acs_state_machine.handle_tr069_message(inform_msg)
        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'Should respond with an InformResponse',
        )

        # Send an empty http request to kick off the rest of provisioning
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for an optional parameter, three times
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_fault()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_fault()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_fault()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for read-only params
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_read_only_param_values_response()

        # Send back some typical values
        # And then SM should request regular parameter values
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )

        # Send back typical values for the regular parameters
        req = Tr069MessageBuilder.\
            get_regular_param_values_response(
                admin_state=False,
                earfcndl=39150,
            )
        resp = acs_state_machine.handle_tr069_message(req)

        # SM will be requesting object parameter values
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting object param vals',
        )

        # Send back some typical values for object parameters
        req = Tr069MessageBuilder.get_object_param_values_response()
        resp = acs_state_machine.handle_tr069_message(req)

        # In this scenario, the ACS and thus state machine will not need
        # to delete or add objects to the eNB configuration.
        # SM should then just be attempting to set parameter values
        self.assertTrue(
            isinstance(resp, models.SetParameterValues),
            'State machine should be setting param values',
        )

        # Send back confirmation that the parameters were successfully set
        req = models.SetParameterValuesResponse()
        req.Status = 0
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for read-only params
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_read_only_param_values_response()

        # Send back some typical values
        # And then SM should continue polling the read-only params
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.DummyInput),
            'State machine should be ending session',
        )

        # If a different eNB is suddenly plugged in, or the same eNB sends a
        # new Inform, enodebd should be able to handle it.
        # Send an Inform message, wait for an InformResponse
        inform_msg = Tr069MessageBuilder.get_inform()
        resp = acs_state_machine.handle_tr069_message(inform_msg)
        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'Should respond with an InformResponse',
        )

        # Send an empty http request to kick off the rest of provisioning
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for an optional parameter, three times
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )

    def test_reboot_after_invasive_changes(self) -> None:
        """
        Test the scenario where:
        - eNodeB has already been powered for 10 minutes without configuration
        - Setting parameters which are 'invasive' on the eNodeB
        - Simulate the scenario up until reboot, and test that enodebd does
          not try to complete configuration after reboot, because it is
          waiting for REM process to finish running
        - This test does not wait the ten minutes to simulate REM process
          finishing on the Baicells eNodeB

        'Invasive' parameters are those which require special behavior to apply
        the changes for the eNodeB.

        In the case of the Baicells eNodeB, properly applying changes to
        invasive parameters requires rebooting the device.
        """
        acs_state_machine = \
            EnodebAcsStateMachineBuilder\
            .build_acs_state_machine(EnodebDeviceName.BAICELLS)
        # Since the test utils pretend the eNB is set to 20MHz, we force this
        # to 10 MHz, so the state machine sets this value.
        acs_state_machine.mconfig.bandwidth_mhz = 10

        # Send an Inform message, wait for an InformResponse
        inform_msg = Tr069MessageBuilder.get_inform(
            '48BF74',
            'BaiBS_RTS_3.1.6',
            '120200002618AGP0003',
            ['2 PERIODIC'],
        )
        resp = acs_state_machine.handle_tr069_message(inform_msg)
        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'Should respond with an InformResponse',
        )

        # Send an empty http request to kick off the rest of provisioning
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for an optional parameter, three times
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_fault()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_fault()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_fault()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for read-only params
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_read_only_param_values_response()

        # Send back some typical values
        # And then SM should request regular parameter values
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )

        # Send back typical values for the regular parameters
        req = Tr069MessageBuilder.get_regular_param_values_response()
        resp = acs_state_machine.handle_tr069_message(req)

        # SM will be requesting object parameter values
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting object param vals',
        )

        # Send back some typical values for object parameters
        req = Tr069MessageBuilder.get_object_param_values_response()
        resp = acs_state_machine.handle_tr069_message(req)

        # In this scenario, the ACS and thus state machine will not need
        # to delete or add objects to the eNB configuration.
        # SM should then just be attempting to set parameter values
        self.assertTrue(
            isinstance(resp, models.SetParameterValues),
            'State machine should be setting param values',
        )

        # Send back confirmation that the parameters were successfully set
        req = models.SetParameterValuesResponse()
        req.Status = 0
        resp = acs_state_machine.handle_tr069_message(req)

        # Since invasive parameters have been set, then to apply the changes
        # to the Baicells eNodeB, we need to reboot the device
        self.assertTrue(isinstance(resp, models.Reboot))
        req = Tr069MessageBuilder.get_reboot_response()
        resp = acs_state_machine.handle_tr069_message(req)

        # After the reboot has been received, enodebd should end the
        # provisioning session
        self.assertTrue(
            isinstance(resp, models.DummyInput),
            'After sending command to reboot the Baicells eNodeB, '
            'enodeb should end the TR-069 session.',
        )

        # At this point, sometime after the eNodeB reboots, we expect it to
        # send an Inform indicating reboot. Since it should be in REM process,
        # we hold off on finishing configuration, and end TR-069 sessions.
        req = Tr069MessageBuilder.get_inform(
            '48BF74', 'BaiBS_RTS_3.1.6',
            '120200002618AGP0003',
            ['1 BOOT', 'M Reboot'],
        )
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.DummyInput),
            'After receiving a post-reboot Inform, enodebd '
            'should end TR-069 sessions for 10 minutes to wait '
            'for REM process to finish.',
        )

        # Pretend that we have waited, and now we are in normal operation again
        acs_state_machine.transition('wait_inform_post_reboot')
        req = Tr069MessageBuilder.get_inform(
            '48BF74', 'BaiBS_RTS_3.1.6',
            '120200002618AGP0003',
            ['2 PERIODIC'],
        )
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'After receiving a post-reboot Inform, enodebd '
            'should end TR-069 sessions for 10 minutes to wait '
            'for REM process to finish.',
        )

        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'enodebd should be requesting params',
        )
        self.assertTrue(
            len(resp.ParameterNames.string) > 1,
            'Should be requesting transient params.',
        )

    def test_reboot_without_getting_optional(self) -> None:
        """
        The state machine should not skip figuring out which optional
        parameters are present.
        """
        acs_state_machine = \
            EnodebAcsStateMachineBuilder \
            .build_acs_state_machine(EnodebDeviceName.BAICELLS)

        # Send an Inform message, wait for an InformResponse
        inform_msg = Tr069MessageBuilder.get_inform(
            '48BF74',
            'BaiBS_RTS_3.1.6',
            '120200002618AGP0003',
            ['2 PERIODIC'],
        )
        resp = acs_state_machine.handle_tr069_message(inform_msg)
        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'Should respond with an InformResponse',
        )

        # And now reboot the eNodeB
        acs_state_machine.transition('reboot')
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(isinstance(resp, models.Reboot))
        req = Tr069MessageBuilder.get_reboot_response()
        resp = acs_state_machine.handle_tr069_message(req)

        # After the reboot has been received, enodebd should end the
        # provisioning session
        self.assertTrue(
            isinstance(resp, models.DummyInput),
            'After sending command to reboot the Baicells eNodeB, '
            'enodeb should end the TR-069 session.',
        )

        # At this point, sometime after the eNodeB reboots, we expect it to
        # send an Inform indicating reboot. Since it should be in REM process,
        # we hold off on finishing configuration, and end TR-069 sessions.
        req = Tr069MessageBuilder.get_inform(
            '48BF74', 'BaiBS_RTS_3.1.6',
            '120200002618AGP0003',
            ['1 BOOT', 'M Reboot'],
        )
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.DummyInput),
            'After receiving a post-reboot Inform, enodebd '
            'should end TR-069 sessions for 10 minutes to wait '
            'for REM process to finish.',
        )

        # Pretend that we have waited, and now we are in normal operation again
        acs_state_machine.transition('wait_inform_post_reboot')
        req = Tr069MessageBuilder.get_inform(
            '48BF74', 'BaiBS_RTS_3.1.6',
            '120200002618AGP0003',
            ['2 PERIODIC'],
        )
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'After receiving a post-reboot Inform, enodebd '
            'should end TR-069 sessions for 10 minutes to wait '
            'for REM process to finish.',
        )

        # Since we haven't figured out the presence of optional parameters, the
        # state machine should be requesting them now. There are three for the
        # Baicells state machine.
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'enodebd should be requesting params',
        )
        self.assertTrue(
            len(resp.ParameterNames.string) == 1,
            'Should be requesting optional params.',
        )
        req = Tr069MessageBuilder.get_fault()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param',
        )
        self.assertTrue(
            len(resp.ParameterNames.string) == 1,
            'Should be requesting optional params.',
        )
        req = Tr069MessageBuilder.get_fault()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param',
        )
        self.assertTrue(
            len(resp.ParameterNames.string) == 1,
            'Should be requesting optional params.',
        )

    def test_missing_mme_timeout_handler(self) -> None:
        acs_state_machine = \
            EnodebAcsStateMachineBuilder \
            .build_acs_state_machine(EnodebDeviceName.BAICELLS)

        # Send an Inform message, wait for an InformResponse
        inform_msg = Tr069MessageBuilder.get_inform(
            '48BF74',
            'BaiBS_RTS_3.1.6',
            '120200002618AGP0003',
            ['2 PERIODIC'],
        )
        acs_state_machine.handle_tr069_message(inform_msg)
        # Send an empty http request to kick off the rest of provisioning
        req = models.DummyInput()
        acs_state_machine.handle_tr069_message(req)

        acs_state_machine.mme_timeout_handler.cancel()
        acs_state_machine.mme_timeout_handler = None

        # Send an Inform message, wait for an InformResponse
        inform_msg = Tr069MessageBuilder.get_inform(
            '48BF74',
            'BaiBS_RTS_3.1.6',
            '120200002618AGP0003',
            ['2 PERIODIC'],
        )
        acs_state_machine.handle_tr069_message(inform_msg)

    def test_fault_after_set_parameters(self) -> None:
        acs_state_machine = \
            EnodebAcsStateMachineBuilder \
            .build_acs_state_machine(EnodebDeviceName.BAICELLS)

        # Send an Inform message, wait for an InformResponse
        inform_msg = Tr069MessageBuilder.get_inform(
            '48BF74',
            'BaiBS_RTS_3.1.6',
            '120200002618AGP0003',
            ['2 PERIODIC'],
        )
        resp = acs_state_machine.handle_tr069_message(inform_msg)
        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'Should respond with an InformResponse',
        )

        # Send an empty http request to kick off the rest of provisioning
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for an optional parameter, three times
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_fault()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_fault()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_fault()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for read-only params
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = Tr069MessageBuilder.get_read_only_param_values_response()

        # Send back some typical values
        # And then SM should request regular parameter values
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )

        # Send back typical values for the regular parameters
        req = Tr069MessageBuilder.get_regular_param_values_response()
        resp = acs_state_machine.handle_tr069_message(req)

        # SM will be requesting object parameter values
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting object param vals',
        )

        # Send back some typical values for object parameters
        req = Tr069MessageBuilder.get_object_param_values_response()
        resp = acs_state_machine.handle_tr069_message(req)

        # In this scenario, the ACS and thus state machine will not need
        # to delete or add objects to the eNB configuration.
        # SM should then just be attempting to set parameter values
        self.assertTrue(
            isinstance(resp, models.SetParameterValues),
            'State machine should be setting param values',
        )

        req = models.Fault()
        req.FaultCode = 12345
        req.FaultString = 'Test FaultString'
        acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            'Error' in acs_state_machine.get_state(),
            'Should be in error state',
        )

    def test_autoremediation_from_fault(self):
        """
        Transition the state machine into the unexpected fault state, then
        verify that it transitions itself back to WaitInform after an Inform
        is received.
        """
        sm = EnodebAcsStateMachineBuilder.build_acs_state_machine(
            EnodebDeviceName.BAICELLS,
        )

        # Send an initial inform
        inform_msg = Tr069MessageBuilder.get_inform(
            '48BF74',
            'BaiBS_RTS_3.1.6',
            '120200002618AGP0003',
            ['2 PERIODIC'],
        )
        resp = sm.handle_tr069_message(inform_msg)
        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'Should respond with an InformResponse',
        )

        # Now send a fault
        req = models.Fault()
        req.FaultCode = 12345
        req.FaultString = 'Test FaultString'
        sm.handle_tr069_message(req)
        self.assertTrue('Error' in sm.get_state(), 'Should be in error state')

        # Send the Inform again, verify SM transitions out of fault
        resp = sm.handle_tr069_message(inform_msg)
        self.assertTrue(isinstance(resp, models.DummyInput))
        self.assertEqual('Waiting for an Inform', sm.get_state())
