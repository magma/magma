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

# pylint: disable=protected-access
from time import sleep
from unittest import TestCase, mock

from dp.protos.enodebd_dp_pb2 import CBSDRequest, CBSDStateResult, LteChannel
from magma.enodebd.data_models.data_model import ParameterName
from magma.enodebd.device_config.configuration_init import build_desired_config
from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration
from magma.enodebd.devices.baicells_qrtb import (
    BaicellsQRTBNotifyDPState,
    BaicellsQRTBQueuedEventsWaitState,
    BaicellsQRTBTrDataModel,
    BaicellsQRTBWaitInformRebootState,
    qrtb_update_desired_config_from_cbsd_state,
)
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.exceptions import ConfigurationError
from magma.enodebd.state_machines.enb_acs_states import (
    WaitEmptyMessageState,
    WaitInformState,
)
from magma.enodebd.tests.test_utils.enb_acs_builder import (
    EnodebAcsStateMachineBuilder,
)
from magma.enodebd.tests.test_utils.enodeb_handler import EnodebHandlerTestCase
from magma.enodebd.tests.test_utils.tr069_msg_builder import (
    Param,
    Tr069MessageBuilder,
)
from magma.enodebd.tr069 import models

DEFAULT_INFORM_PARAMS = [
    Param('Device.DeviceInfo.HardwareVersion', 'string', 'E01'),
    Param('Device.DeviceInfo.SAS.RadioEnable', 'boolean', 'false'),
    Param('Device.DeviceInfo.SoftwareVersion', 'string', 'BaiBS_QRTB_2.6.2'),
    Param('Device.DeviceInfo.X_COM_STATION_RUN_Time', 'string', '0d 0h 2m 3s'),
    Param('Device.RootDataModelVersion', 'float', '2.8'),
    Param('Device.Services.FAPService.1.FAPControl.LTE.OpState', 'boolean', 'false'),
    Param('Device.Services.FAPService.1.FAPControl.LTE.RFTxStatus', 'boolean', 'false'),
    Param('Device.Services.FAPService.1.FAPControl.LTE.X_COM_HSS.HaloBEnableState', 'int', '0'),
    Param('Device.Services.FAPService.1.FAPControl.LTE.X_COM_HSS.HaloBMode', 'int', '0'),
    Param('Device.Services.FAPService.1.FAPControl.X_RADISYS_COM_AlarmStatus', 'string', 'Cleared'),
    Param('Device.Services.FAPService.2.FAPControl.LTE.OpState', 'boolean', 'false'),
    Param('Device.Services.FAPService.2.FAPControl.LTE.RFTxStatus', 'boolean', 'false'),
    Param('Device.Services.FAPService.2.FAPControl.LTE.X_COM_HSS.HaloBEnableState', 'int', '0'),
    Param('Device.Services.FAPService.2.FAPControl.LTE.X_COM_HSS.HaloBMode', 'int', '0'),
]

GET_TRANSIENT_PARAMS_RESPONSE_PARAMS = [
    Param('Device.DeviceInfo.X_COM_1588_Status', 'int', '0'),
    Param('Device.DeviceInfo.X_COM_GPS_Status', 'int', '1'),
    Param('Device.DeviceInfo.X_COM_MME_Status', 'int', '0'),
    Param('Device.FAP.GPS.LockedLatitude', 'int', '40002899'),
    Param('Device.FAP.GPS.LockedLongitude', 'unsigned', '-105287994'),
    Param('Device.Services.FAPService.1.FAPControl.LTE.OpState', 'boolean', 'false'),
    Param('Device.Services.FAPService.1.FAPControl.LTE.RFTxStatus', 'boolean', 'false'),
]

GET_PARAMS_RESPONSE_PARAMS = [
    Param('Device.DeviceInfo.SAS.FccId', 'str', '2AG32MBS3100196N'),
    Param('Device.DeviceInfo.SAS.UserId', 'str', 'M0LK4T'),
    Param('Device.DeviceInfo.SerialNumber', 'str', '120200024019APP0105'),
    Param('Device.DeviceInfo.X_COM_GpsSyncEnable', 'boolean', 'true'),
    Param('Device.DeviceInfo.X_COM_LTE_LGW_Switch', 'int', '0'),
    Param('Device.DeviceInfo.X_COM_REM_Status', 'int', '0'),
    Param('Device.DeviceInfo.SAS.enableMode', 'int', '1'),
    Param('Device.FAP.PerfMgmt.Config.1.Enable', 'boolean', 'true'),
    Param('Device.FAP.PerfMgmt.Config.1.PeriodicUploadInterval', 'int', '300'),
    Param('Device.FAP.PerfMgmt.Config.1.URL', 'string', 'http://192.168.60.142:8081/'),
    Param('Device.ManagementServer.PeriodicInformEnable', 'bool', 'true'),
    Param('Device.ManagementServer.PeriodicInformInterval', 'int', '5'),
    Param('Device.Services.FAPService.1.Capabilities.LTE.BandsSupported', 'int', '7'),
    Param('Device.Services.FAPService.1.Capabilities.LTE.DuplexMode', 'string', 'TDDMode'),
    Param('Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNListNumberOfEntries', 'int', '1'),
    Param('Device.Services.FAPService.1.CellConfig.LTE.EPC.TAC', 'int', '1'),
    Param('Device.Services.FAPService.1.CellConfig.LTE.RAN.CellRestriction.CellReservedForOperatorUse', 'int', '0'),
    Param('Device.Services.FAPService.1.CellConfig.LTE.RAN.Common.CellIdentity', 'int', '138777000'),
    Param('Device.Services.FAPService.1.CellConfig.LTE.RAN.PHY.TDDFrame.SpecialSubframePatterns', 'int', '7'),
    Param('Device.Services.FAPService.1.CellConfig.LTE.RAN.PHY.TDDFrame.SubFrameAssignment', 'int', '2'),
    Param('Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.DLBandwidth', 'str', '100'),
    Param('Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.EARFCNDL', 'int', '39150'),
    Param('Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.EARFCNUL', 'int', '39150'),
    Param('Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.FreqBandIndicator', 'int', '40'),
    Param('Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.PhyCellID', 'int', '260'),
    Param('Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.ReferenceSignalPower', 'unsigned', '-24'),
    Param('Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.ULBandwidth', 'int', '100'),
    Param('Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.X_COM_RadioEnable', 'bool', 'false'),
    Param('Device.Services.FAPService.1.FAPControl.LTE.AdminState', 'bool', 'false'),
    Param('Device.Services.FAPService.1.FAPControl.LTE.Gateway.S1SigLinkPort', 'int', '36412'),
    Param('Device.Services.FAPService.1.FAPControl.LTE.Gateway.S1SigLinkServerList', 'string', '192.168.60.142'),
    Param('Device.Services.FAPService.1.FAPControl.LTE.Gateway.X_COM_MmePool.Enable', 'boolean', 'false'),
    Param('Device.Services.FAPService.Ipsec.IPSEC_ENABLE', 'boolean', 'false'),
]

MOCK_CBSD_STATE = CBSDStateResult(
    radio_enabled=True,
    channel=LteChannel(
        low_frequency_hz=3550_000_000,
        high_frequency_hz=3570_000_000,
        max_eirp_dbm_mhz=34,
    ),
)


class SasToRfConfigTests(TestCase):
    def test_bandwidth_20MHz(self) -> None:
        config = EnodebConfiguration(BaicellsQRTBTrDataModel())
        channel = LteChannel(
            low_frequency_hz=3550_000_000,
            high_frequency_hz=3570_000_000,
            max_eirp_dbm_mhz=34,
        )
        state = CBSDStateResult(
            radio_enabled=True,
            channel=channel,
        )
        qrtb_update_desired_config_from_cbsd_state(state, config)
        self.assert_config_updated(config, '100', 55340, 34)

    def test_bandwidth_15MHz(self) -> None:
        config = EnodebConfiguration(BaicellsQRTBTrDataModel())
        channel = LteChannel(
            low_frequency_hz=3555_000_000,
            high_frequency_hz=3570_000_000,
            max_eirp_dbm_mhz=37,
        )
        state = CBSDStateResult(
            radio_enabled=True,
            channel=channel,
        )
        qrtb_update_desired_config_from_cbsd_state(state, config)
        self.assert_config_updated(config, '75', 55365, 37)

    def test_bandwidth_10MHz(self) -> None:
        config = EnodebConfiguration(BaicellsQRTBTrDataModel())
        channel = LteChannel(
            low_frequency_hz=3600_000_000,
            high_frequency_hz=3610_000_000,
            max_eirp_dbm_mhz=24,
        )
        state = CBSDStateResult(
            radio_enabled=True,
            channel=channel,
        )
        qrtb_update_desired_config_from_cbsd_state(state, config)
        self.assert_config_updated(config, '50', 55790, 24)

    def test_bandwidth_5MHz(self) -> None:
        config = EnodebConfiguration(BaicellsQRTBTrDataModel())
        channel = LteChannel(
            low_frequency_hz=3600_000_000,
            high_frequency_hz=3605_000_000,
            max_eirp_dbm_mhz=22,
        )
        state = CBSDStateResult(
            radio_enabled=True,
            channel=channel,
        )
        qrtb_update_desired_config_from_cbsd_state(state, config)
        self.assert_config_updated(config, '25', 55765, 22)

    def test_unsupported_bandwidth_3MHz(self) -> None:
        config = EnodebConfiguration(BaicellsQRTBTrDataModel())
        channel = LteChannel(
            low_frequency_hz=3550_000_000,
            high_frequency_hz=3553_000_000,
            max_eirp_dbm_mhz=17,
        )
        state = CBSDStateResult(
            radio_enabled=True,
            channel=channel,
        )
        with self.assertRaises(ConfigurationError):
            qrtb_update_desired_config_from_cbsd_state(state, config)

    def test_unsupported_bandwidth_1_4MHz(self) -> None:
        config = EnodebConfiguration(BaicellsQRTBTrDataModel())
        channel = LteChannel(
            low_frequency_hz=3550_000_000,
            high_frequency_hz=3551_400_000,
            max_eirp_dbm_mhz=10,
        )
        state = CBSDStateResult(
            radio_enabled=True,
            channel=channel,
        )
        with self.assertRaises(ConfigurationError):
            qrtb_update_desired_config_from_cbsd_state(state, config)

    def test_bandwidth_non_default_MHz(self) -> None:
        config = EnodebConfiguration(BaicellsQRTBTrDataModel())
        channel = LteChannel(
            low_frequency_hz=3690_000_000,
            high_frequency_hz=3698_000_000,
            max_eirp_dbm_mhz=1,
        )
        state = CBSDStateResult(
            radio_enabled=True,
            channel=channel,
        )
        qrtb_update_desired_config_from_cbsd_state(state, config)
        self.assert_config_updated(config, '25', 56680, 1)

    def test_too_low_bandwidth(self) -> None:
        config = EnodebConfiguration(BaicellsQRTBTrDataModel())
        channel = LteChannel(
            low_frequency_hz=3550_000_000,
            high_frequency_hz=3551_000_000,
            max_eirp_dbm_mhz=10,
        )
        state = CBSDStateResult(
            radio_enabled=True,
            channel=channel,
        )
        with self.assertRaises(ConfigurationError):
            qrtb_update_desired_config_from_cbsd_state(state, config)

    def test_omit_other_params_when_radio_disabled(self) -> None:
        config = EnodebConfiguration(BaicellsQRTBTrDataModel())
        channel = LteChannel(
            low_frequency_hz=3550_000_000,
            high_frequency_hz=3560_000_000,
            max_eirp_dbm_mhz=-100,
        )
        state = CBSDStateResult(
            radio_enabled=False,
            channel=channel,
        )
        qrtb_update_desired_config_from_cbsd_state(state, config)
        self.assertEqual(
            config.get_parameter(
                ParameterName.SAS_RADIO_ENABLE,
            ), False,
        )

    def test_power_spectral_density_below_range(self) -> None:
        config = EnodebConfiguration(BaicellsQRTBTrDataModel())
        channel = LteChannel(
            low_frequency_hz=3550_000_000,
            high_frequency_hz=3560_000_000,
            max_eirp_dbm_mhz=-138,
        )
        state = CBSDStateResult(
            radio_enabled=True,
            channel=channel,
        )
        with self.assertRaises(ConfigurationError):
            qrtb_update_desired_config_from_cbsd_state(state, config)

    def test_power_spectral_density_above_range(self) -> None:
        config = EnodebConfiguration(BaicellsQRTBTrDataModel())
        channel = LteChannel(
            low_frequency_hz=3550_000_000,
            high_frequency_hz=3560_000_000,
            max_eirp_dbm_mhz=38,
        )
        state = CBSDStateResult(
            radio_enabled=True,
            channel=channel,
        )
        with self.assertRaises(ConfigurationError):
            qrtb_update_desired_config_from_cbsd_state(state, config)

    def assert_config_updated(self, config: EnodebConfiguration, bandwidth: str, earfcn: int, eirp: int) -> None:
        expected_values = {
            ParameterName.SAS_RADIO_ENABLE: True,
            ParameterName.DL_BANDWIDTH: bandwidth,
            ParameterName.UL_BANDWIDTH: bandwidth,
            ParameterName.EARFCNDL: earfcn,
            ParameterName.EARFCNUL: earfcn,
            ParameterName.POWER_SPECTRAL_DENSITY: eirp,
            ParameterName.BAND: 48,
        }
        for key, value in expected_values.items():
            self.assertEqual(config.get_parameter(key), value)


class BaicellsQRTBHandlerTests(EnodebHandlerTestCase):

    @mock.patch('magma.enodebd.devices.baicells_qrtb.get_cbsd_state')
    def test_notify_dp_sets_values_received_by_dp_in_desired_config(self, mock_get_state) -> None:
        expected_final_param_values = {
            ParameterName.UL_BANDWIDTH: '100',
            ParameterName.DL_BANDWIDTH: '100',
            ParameterName.EARFCNUL: 55340,
            ParameterName.EARFCNDL: 55340,
            ParameterName.POWER_SPECTRAL_DENSITY: 34,
        }
        test_user = 'test_user'
        test_fcc_id = 'fcc_id'
        test_serial_number = '123'

        acs_state_machine = EnodebAcsStateMachineBuilder.build_acs_state_machine(EnodebDeviceName.BAICELLS_QRTB)

        acs_state_machine.desired_cfg = EnodebConfiguration(BaicellsQRTBTrDataModel())

        acs_state_machine.device_cfg.set_parameter(ParameterName.SAS_USER_ID, test_user)
        acs_state_machine.device_cfg.set_parameter(ParameterName.SAS_FCC_ID, test_fcc_id)
        acs_state_machine.device_cfg.set_parameter(ParameterName.SERIAL_NUMBER, test_serial_number)

        for param in expected_final_param_values:
            with self.assertRaises(KeyError):
                acs_state_machine.desired_cfg.get_parameter(param)

        # Skip previous steps not to duplicate the code
        acs_state_machine.transition('check_wait_get_params')
        req = Tr069MessageBuilder.param_values_qrtb_response(
            [], models.GetParameterValuesResponse,
        )

        mock_get_state.return_value = MOCK_CBSD_STATE

        resp = acs_state_machine.handle_tr069_message(req)

        mock_get_state.assert_called_with(
            CBSDRequest(serial_number=test_serial_number),
        )

        self.assertTrue(isinstance(resp, models.DummyInput))

        for param, value in expected_final_param_values.items():
            self.assertEqual(
                value, acs_state_machine.desired_cfg.get_parameter(param),
            )

    def test_manual_reboot(self) -> None:
        """
        Test a scenario where a Magma user goes through the enodebd CLI to
        reboot the Baicells eNodeB.

        This checks the scenario where the command is not sent in the middle
        of a TR-069 provisioning session.
        """
        acs_state_machine = EnodebAcsStateMachineBuilder.build_acs_state_machine(EnodebDeviceName.BAICELLS_QRTB)

        # User uses the CLI tool to get eNodeB to reboot
        acs_state_machine.reboot_asap()

        # And now the Inform message arrives from the eNodeB
        inform_msg = \
            Tr069MessageBuilder.get_qrtb_inform(
                params=DEFAULT_INFORM_PARAMS,
                oui='48BF74',
                enb_serial='1202000181186TB0006',
                event_codes=['2 PERIODIC'],
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

    def test_manual_reboot_during_provisioning(self) -> None:
        """
        Test a scenario where a Magma user goes through the enodebd CLI to
        reboot the Baicells eNodeB.

        This checks the scenario where the command is sent in the middle
        of a TR-069 provisioning session.
        """
        acs_state_machine = EnodebAcsStateMachineBuilder.build_acs_state_machine(EnodebDeviceName.BAICELLS_QRTB)

        # Send an Inform message, wait for an InformResponse
        inform_msg = Tr069MessageBuilder.get_qrtb_inform(
            params=DEFAULT_INFORM_PARAMS,
            oui='48BF74',
            enb_serial='1202000181186TB0006',
            event_codes=['2 PERIODIC'],
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

    @mock.patch('magma.enodebd.devices.baicells_qrtb.get_cbsd_state')
    def test_provision(self, mock_get_state) -> None:
        mock_get_state.return_value = MOCK_CBSD_STATE

        acs_state_machine = EnodebAcsStateMachineBuilder.build_acs_state_machine(EnodebDeviceName.BAICELLS_QRTB)
        data_model = BaicellsQRTBTrDataModel()

        # Send an Inform message, wait for an InformResponse
        inform_msg = Tr069MessageBuilder.get_qrtb_inform(
            params=DEFAULT_INFORM_PARAMS,
            oui='48BF74',
            enb_serial='1202000181186TB0006',
            event_codes=['2 PERIODIC'],
        )
        resp = acs_state_machine.handle_tr069_message(inform_msg)
        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'Should respond with an InformResponse',
        )

        # Send an empty http request to kick off the rest of provisioning
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for transient params
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )

        should_ask_for = [p.name for p in GET_TRANSIENT_PARAMS_RESPONSE_PARAMS]
        self.verify_acs_asking_enb_for_params(should_ask_for, resp)

        # Send back some typical values
        req = Tr069MessageBuilder.param_values_qrtb_response(
            GET_TRANSIENT_PARAMS_RESPONSE_PARAMS, models.GetParameterValuesResponse,
        )

        # And then SM should request all parameter values from the data model
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )

        should_ask_for = [
            data_model.get_parameter(
                pn,
            ).path for pn in data_model.get_parameter_names()
        ]
        self.verify_acs_asking_enb_for_params(should_ask_for, resp)

        # Send back typical values
        req = Tr069MessageBuilder.param_values_qrtb_response(
            GET_PARAMS_RESPONSE_PARAMS,
            models.GetParameterValuesResponse,
        )

        # SM will be requesting object parameter values
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting object param vals',
        )

        should_ask_for = [
            'Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.CellReservedForOperatorUse',
            'Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.Enable',
            'Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.IsPrimary',
            'Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.PLMNID',
        ]

        self.verify_acs_asking_enb_for_params(should_ask_for, resp)

        # Send back some typical values for object parameters
        req = Tr069MessageBuilder.get_object_param_values_response(
            cell_reserved_for_operator_use='true', enable='true', is_primary='true', plmnid='00101',
        )
        resp = acs_state_machine.handle_tr069_message(req)

        # All the radio responses were intentionally crafted so that they match enodebd desired config.
        # Therefore the provisioning ends here. The radio goes directly into end_session -> notify_dp state
        self.assertTrue(
            isinstance(resp, models.DummyInput),
            'State machine should send back an empty message',
        )

        self.assertIsInstance(
            acs_state_machine.state,
            BaicellsQRTBNotifyDPState,
        )

    def verify_acs_asking_enb_for_params(self, should_ask_for, response):
        param_values_that_enodebd_actually_asked_for = response.ParameterNames.string

        self.assertEqual(
            sorted(should_ask_for),
            sorted(param_values_that_enodebd_actually_asked_for),
        )


class BaicellsQRTBStatesTests(EnodebHandlerTestCase):
    """Testing Baicells QRTB specific states"""

    @mock.patch('magma.enodebd.devices.baicells_qrtb.get_cbsd_state')
    def test_end_session_and_notify_dp_transition(self, mock_get_state):
        """Testing if SM steps in and out of BaicellsQRTBWaitNotifyDPState as per state map"""

        mock_get_state.return_value = MOCK_CBSD_STATE

        acs_state_machine = provision_clean_sm(
            state='check_wait_get_params',
        )

        msg = Tr069MessageBuilder.param_values_qrtb_response(
            GET_PARAMS_RESPONSE_PARAMS,
            models.GetParameterValuesResponse,
        )

        # SM should transition from check_wait_get_params to end_session -> notify_dp automatically
        # upon receiving response from the radio
        acs_state_machine.handle_tr069_message(msg)

        self.assertIsInstance(
            acs_state_machine.state,
            BaicellsQRTBNotifyDPState,
        )

        msg = Tr069MessageBuilder.get_inform(event_codes=['1 BOOT'])

        # SM should go into wait_inform state, respond with Inform response and transition to wait_empty
        acs_state_machine.handle_tr069_message(msg)

        self.assertIsInstance(acs_state_machine.state, WaitEmptyMessageState)

    def test_wait_post_reboot_inform_transition(self):
        """Testing if SM steps in and out of BaicellsQRTBWaitInformRebootState as per state map"""
        acs_state_machine = provision_clean_sm(state='wait_reboot')

        acs_state_machine.handle_tr069_message(models.RebootResponse())

        self.assertIsInstance(
            acs_state_machine.state,
            BaicellsQRTBWaitInformRebootState,
        )

        msg = Tr069MessageBuilder.get_inform(event_codes=['1 BOOT'])

        acs_state_machine.handle_tr069_message(msg)

        self.assertIsInstance(
            acs_state_machine.state,
            BaicellsQRTBQueuedEventsWaitState,
        )

    @mock.patch.object(BaicellsQRTBQueuedEventsWaitState, 'CONFIG_DELAY_AFTER_BOOT', new=0.5)
    def test_wait_queued_events_post_reboot_transition(self):
        """Testing if we step in and out of BaicellsQRTBQueuedEventsWaitState as per state map"""
        acs_state_machine = provision_clean_sm(
            state='wait_post_reboot_inform',
        )

        msg = Tr069MessageBuilder.get_inform(event_codes=['1 BOOT'])

        acs_state_machine.handle_tr069_message(msg)

        self.assertIsInstance(
            acs_state_machine.state,
            BaicellsQRTBQueuedEventsWaitState,
        )

        msg = Tr069MessageBuilder.get_inform()

        # Need to wait after reboot, timeout set to 0.5 sec
        sleep(1)

        acs_state_machine.handle_tr069_message(msg)

        self.assertIsInstance(acs_state_machine.state, WaitInformState)


class BaicellsQRTBConfigTests(EnodebHandlerTestCase):
    def test_frequency_related_params_removed_in_postprocessor(self):
        acs_state_machine = provision_clean_sm()
        acs_state_machine.device_cfg.set_parameter(ParameterName.IP_SEC_ENABLE, 'false')

        parameters_to_delete = [
            ParameterName.RADIO_ENABLE, ParameterName.POWER_SPECTRAL_DENSITY,
            ParameterName.EARFCNDL, ParameterName.EARFCNUL, ParameterName.BAND,
            ParameterName.DL_BANDWIDTH, ParameterName.UL_BANDWIDTH,
            ParameterName.SAS_RADIO_ENABLE,
        ]
        for p in parameters_to_delete:
            acs_state_machine.device_cfg.set_parameter(p, 'some_value')

        desired_cfg = build_desired_config(
            acs_state_machine.mconfig,
            acs_state_machine.service_config,
            acs_state_machine.device_cfg,
            acs_state_machine.data_model,
            acs_state_machine.config_postprocessor,
        )
        for p in parameters_to_delete:
            self.assertFalse(desired_cfg.has_parameter(p))

    def test_device_and_desired_config_discrepancy_after_initial_configuration(self):
        """
        Testing a situation where device_cfg and desired_cfg are already present on the state machine,
        because the initial configuration of the radio has occurred, but then the configs have diverged
        (e.g. as a result of domain-proxy setting different values on the desired config)
        """
        # Skipping previous states
        acs_state_machine = provision_clean_sm('wait_get_transient_params')

        # Need to set this param on the device_cfg first, otherwise we won't be able to generate the desired_cfg
        # using 'build_desired_config' function
        acs_state_machine.device_cfg.set_parameter(ParameterName.IP_SEC_ENABLE, 'false')

        acs_state_machine.desired_cfg = build_desired_config(
            acs_state_machine.mconfig,
            acs_state_machine.service_config,
            acs_state_machine.device_cfg,
            acs_state_machine.data_model,
            acs_state_machine.config_postprocessor,
        )

        prepare_device_cfg_same_as_desired_cfg(acs_state_machine)

        # Let's say that while in 'wait_for_dp' state, the DP asked us to change this param's value in the desired_cfg
        acs_state_machine.desired_cfg.set_parameter(ParameterName.EARFCNDL, 5)

        req = Tr069MessageBuilder.param_values_qrtb_response(
            GET_TRANSIENT_PARAMS_RESPONSE_PARAMS,
            models.GetParameterValuesResponse,
        )

        # ACS asking for all params from data model
        acs_state_machine.handle_tr069_message(req)

        req = Tr069MessageBuilder.param_values_qrtb_response(
            GET_PARAMS_RESPONSE_PARAMS, models.GetParameterValuesResponse,
        )

        # ACS asking for object params
        acs_state_machine.handle_tr069_message(req)
        req = Tr069MessageBuilder.get_object_param_values_response(
            cell_reserved_for_operator_use='true', enable='true', is_primary='true', plmnid='00101',
        )

        # ACS should ask the radio to correct the only parameter that doesn't match the desired config - EARFCNDL
        resp = acs_state_machine.handle_tr069_message(req)

        self.assertIsInstance(resp, models.SetParameterValues)
        self.assertEqual(1, len(resp.ParameterList.ParameterValueStruct))
        self.assertEqual(
            'Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.EARFCNDL',
            resp.ParameterList.ParameterValueStruct[0].Name,
        )
        self.assertEqual('5', resp.ParameterList.ParameterValueStruct[0].Value.Data)

    def test_sas_enable_mode_is_enabled(self):
        acs_state_machine = provision_clean_sm()
        acs_state_machine.device_cfg.set_parameter(ParameterName.IP_SEC_ENABLE, 'false')

        acs_state_machine.desired_cfg = build_desired_config(
            acs_state_machine.mconfig,
            acs_state_machine.service_config,
            acs_state_machine.device_cfg,
            acs_state_machine.data_model,
            acs_state_machine.config_postprocessor,
        )

        self.assertEqual(1, acs_state_machine.desired_cfg.get_parameter(ParameterName.SAS_ENABLED))


def prepare_device_cfg_same_as_desired_cfg(acs_state_machine):
    # Setting all regular params in the device_cfg to the same values as those in desired_cfg
    for param in acs_state_machine.desired_cfg.get_parameter_names():
        acs_state_machine.device_cfg.set_parameter(param, acs_state_machine.desired_cfg.get_parameter(param))

    # Setting all object params in the device_cfg to the same values as those in desired_cfg
    for obj in acs_state_machine.desired_cfg.get_object_names():
        acs_state_machine.device_cfg.add_object(obj)
        for param in acs_state_machine.desired_cfg.get_parameter_names_for_object(obj):
            val = acs_state_machine.desired_cfg.get_parameter_for_object(param, obj)
            acs_state_machine.device_cfg.set_parameter_for_object(param, val, obj)


def provision_clean_sm(state=None):
    acs_state_machine = EnodebAcsStateMachineBuilder.build_acs_state_machine(EnodebDeviceName.BAICELLS_QRTB)
    acs_state_machine.desired_cfg = EnodebConfiguration(
        BaicellsQRTBTrDataModel(),
    )
    if state is not None:
        acs_state_machine.transition(state)
    return acs_state_machine
