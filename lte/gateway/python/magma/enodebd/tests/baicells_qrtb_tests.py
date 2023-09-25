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
import logging
from copy import deepcopy
from time import sleep
from typing import Any, Dict
from unittest import TestCase, mock

from dp.protos.cbsd_pb2 import CBSDStateResult, LteChannel
from magma.enodebd.data_models.data_model import ParameterName
from magma.enodebd.device_config.configuration_init import build_desired_config
from magma.enodebd.device_config.configuration_util import (
    calc_bandwidth_mhz,
    calc_bandwidth_rbs,
    calc_earfcn,
)
from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration
from magma.enodebd.devices.baicells_qrtb.data_model import (
    BaicellsQRTBTrDataModel,
)
from magma.enodebd.devices.baicells_qrtb.params import (
    CarrierAggregationParameters,
)
from magma.enodebd.devices.baicells_qrtb.states import (
    BaicellsQRTBNotifyDPState,
    BaicellsQRTBQueuedEventsWaitState,
    BaicellsQRTBWaitInformRebootState,
    qrtb_update_desired_config_from_cbsd_state,
)
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.dp_client import build_enodebd_update_cbsd_request
from magma.enodebd.exceptions import ConfigurationError
from magma.enodebd.state_machines.acs_state_utils import (
    get_firmware_upgrade_download_config,
)
from magma.enodebd.state_machines.enb_acs_states import (
    FirmwareUpgradeDownloadState,
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
from parameterized import parameterized

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
    Param('Device.DeviceInfo.cbsdCategory', 'string', 'A'),
    Param('Device.DeviceInfo.indoorDeployment', 'boolean', 'true'),
    Param('Device.DeviceInfo.AntennaInfo.Gain', 'int', '10'),
    Param('Device.DeviceInfo.AntennaInfo.Height', 'int', '8'),
    Param('Device.DeviceInfo.AntennaInfo.HeightType', 'string', 'AGL'),
    Param('Device.DeviceInfo.SAS.FccId', 'str', '2AG32MBS3100196N'),
    Param('Device.DeviceInfo.SAS.FccId', 'str', 'SAS.FccId'),
    Param('Device.DeviceInfo.SAS.UserId', 'str', 'SAS.UserId'),
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
    Param('Device.Services.FAPService.Ipsec.IPSEC_ENABLE', 'bool', 'false'),

    Param('Device.Services.FAPService.1.CellConfig.LTE.RAN.CA.CaEnable', 'bool', 'false'),
    Param('FAPService.1.CellConfig.LTE.RAN.CA.PARAMS.NumOfCells', 'int', '1'),
    Param('Device.Services.FAPService.2.CellConfig.LTE.RAN.Common.CellIdentity', 'int', '138777000'),
    Param('Device.Services.FAPService.2.CellConfig.LTE.RAN.RF.FreqBandIndicator', 'int', '48'),
    Param('Device.Services.FAPService.2.CellConfig.LTE.RAN.RF.DLBandwidth', 'str', '100'),
    Param('Device.Services.FAPService.2.CellConfig.LTE.RAN.RF.ULBandwidth', 'int', '100'),
    Param('Device.Services.FAPService.2.CellConfig.LTE.RAN.RF.EARFCNDL', 'int', '39150'),
    Param('Device.Services.FAPService.2.CellConfig.LTE.RAN.RF.EARFCNUL', 'int', '39150'),
    Param('Device.Services.FAPService.2.CellConfig.LTE.RAN.RF.PhyCellID', 'int', '260'),
    Param('Device.Services.FAPService.2.CellConfig.LTE.RAN.RF.X_COM_RadioEnable', 'bool', 'false'),
    Param('Device.Services.FAPService.2.FAPControl.LTE.AdminState', 'bool', 'true'),
    Param('Device.Services.FAPService.2.FAPControl.LTE.OpState', 'bool', 'true'),
    Param('Device.Services.FAPService.2.FAPControl.LTE.RFTxStatus', 'boolean', 'false'),
    Param('Device.Services.FAPService.2.CellConfig.LTE.EPC.PLMNList.1.CellReservedForOperatorUse', 'bool', 'true'),
    Param('Device.Services.FAPService.2.CellConfig.LTE.EPC.PLMNList.1.Enable', 'bool', 'true'),
    Param('Device.Services.FAPService.2.CellConfig.LTE.EPC.PLMNList.1.IsPrimary', 'bool', 'false'),
    Param('Device.Services.FAPService.2.CellConfig.LTE.EPC.PLMNList.1.PLMNID', 'int', '00102'),
]

MOCK_CBSD_STATE = CBSDStateResult(
    radio_enabled=True,
    channels=[
        LteChannel(
            low_frequency_hz=3550_000_000,
            high_frequency_hz=3570_000_000,
            max_eirp_dbm_mhz=34,
        ),
        LteChannel(
            low_frequency_hz=3570_000_000,
            high_frequency_hz=3590_000_000,
            max_eirp_dbm_mhz=34,
        ),
    ],
)


class SasToRfConfigTests(TestCase):
    def test_bandwidth_20MHz(self) -> None:
        config = EnodebConfiguration(BaicellsQRTBTrDataModel())
        channel = LteChannel(
            low_frequency_hz=3550_000_000,
            high_frequency_hz=3570_000_000,
            max_eirp_dbm_mhz=34,
        )
        channels = [channel]
        state = CBSDStateResult(
            radio_enabled=True,
            channels=channels,
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
            channels=[channel],
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
            channels=[channel],
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
            channels=[channel],
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
            channels=[channel],
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
            channels=[channel],
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
            channels=[channel],
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
            channels=[channel],
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
            channels=[channel],
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
            channels=[channel],
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
            channels=[channel],
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
    @mock.patch('magma.enodebd.devices.baicells_qrtb.states.enodebd_update_cbsd')
    def test_enodebd_update_cbsd_not_called_when_gps_unavailable(self, mock_enodebd_update_cbsd) -> None:
        test_serial_number = '120200024019APP0105'

        acs_state_machine = EnodebAcsStateMachineBuilder.build_acs_state_machine(EnodebDeviceName.BAICELLS_QRTB)

        acs_state_machine.desired_cfg = EnodebConfiguration(BaicellsQRTBTrDataModel())

        req = Tr069MessageBuilder.get_qrtb_inform(
            params=DEFAULT_INFORM_PARAMS,
            oui='48BF74',
            enb_serial=test_serial_number,
            event_codes=['2 PERIODIC'],
        )
        acs_state_machine.device_cfg.set_parameter(ParameterName.GPS_STATUS, '0')
        acs_state_machine.handle_tr069_message(req)
        acs_state_machine.transition('notify_dp')
        mock_enodebd_update_cbsd.assert_not_called()

    @mock.patch('magma.enodebd.devices.baicells_qrtb.states.enodebd_update_cbsd')
    def test_notify_dp(self, mock_enodebd_update_cbsd) -> None:
        expected_final_param_values = {
            ParameterName.UL_BANDWIDTH: '100',
            ParameterName.DL_BANDWIDTH: '100',
            ParameterName.EARFCNUL: 55340,
            ParameterName.EARFCNDL: 55340,
            ParameterName.POWER_SPECTRAL_DENSITY: 34,
        }

        # This serial needs to match the one defined in GET_PARAMS_RESPONSE_PARAMS
        test_serial_number = '120200024019APP0105'

        acs_state_machine = EnodebAcsStateMachineBuilder.build_acs_state_machine(EnodebDeviceName.BAICELLS_QRTB)

        acs_state_machine.desired_cfg = EnodebConfiguration(BaicellsQRTBTrDataModel())

        for param in expected_final_param_values:
            with self.assertRaises(KeyError):
                acs_state_machine.desired_cfg.get_parameter(param)

        # Do basic provisioning
        # Do not check, provisioning check already does that
        # We just want parameters to load
        req = Tr069MessageBuilder.get_qrtb_inform(
            params=DEFAULT_INFORM_PARAMS,
            oui='48BF74',
            enb_serial=test_serial_number,
            event_codes=['2 PERIODIC'],
        )
        acs_state_machine.handle_tr069_message(req)

        req = models.DummyInput()
        acs_state_machine.handle_tr069_message(req)

        req = Tr069MessageBuilder.param_values_qrtb_response(
            GET_TRANSIENT_PARAMS_RESPONSE_PARAMS, models.GetParameterValuesResponse,
        )
        acs_state_machine.handle_tr069_message(req)

        req = Tr069MessageBuilder.param_values_qrtb_response(
            GET_PARAMS_RESPONSE_PARAMS,
            models.GetParameterValuesResponse,
        )
        acs_state_machine.handle_tr069_message(req)

        # Transition to final get params check state
        acs_state_machine.transition('check_wait_get_params')
        req = Tr069MessageBuilder.param_values_qrtb_response(
            [], models.GetParameterValuesResponse,
        )

        mock_enodebd_update_cbsd.return_value = MOCK_CBSD_STATE

        resp = acs_state_machine.handle_tr069_message(req)

        enodebd_update_cbsd_request = build_enodebd_update_cbsd_request(
            serial_number=acs_state_machine.device_cfg.get_parameter(ParameterName.SERIAL_NUMBER),
            latitude_deg=acs_state_machine.device_cfg.get_parameter(ParameterName.GPS_LAT),
            longitude_deg=acs_state_machine.device_cfg.get_parameter(ParameterName.GPS_LONG),
            indoor_deployment=acs_state_machine.device_cfg.get_parameter(ParameterName.INDOOR_DEPLOYMENT),
            antenna_height=acs_state_machine.device_cfg.get_parameter(ParameterName.ANTENNA_HEIGHT),
            antenna_height_type=acs_state_machine.device_cfg.get_parameter(ParameterName.ANTENNA_HEIGHT_TYPE),
            cbsd_category=acs_state_machine.device_cfg.get_parameter(ParameterName.CBSD_CATEGORY),
        )
        mock_enodebd_update_cbsd.assert_called_with(enodebd_update_cbsd_request)

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

    @mock.patch('magma.enodebd.devices.baicells_qrtb.states.enodebd_update_cbsd')
    def test_provision(self, mock_enodebd_update_cbsd) -> None:
        self.maxDiff = None
        mock_enodebd_update_cbsd.return_value = MOCK_CBSD_STATE

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
            f'State machine should send back an empty message, got {resp} instead.',
        )

        self.assertIsInstance(
            acs_state_machine.state,
            BaicellsQRTBNotifyDPState,
        )

    def verify_acs_asking_enb_for_params(self, should_ask_for, response):
        self.maxDiff = None
        param_values_that_enodebd_actually_asked_for = response.ParameterNames.string

        self.assertEqual(
            sorted(should_ask_for),
            sorted(param_values_that_enodebd_actually_asked_for),
        )


class BaicellsQRTBStatesTests(EnodebHandlerTestCase):
    """Testing Baicells QRTB specific states"""

    @mock.patch('magma.enodebd.devices.baicells_qrtb.states.enodebd_update_cbsd')
    def test_end_session_and_notify_dp_transition(self, mock_enodebd_update_cbsd):
        """Testing if SM steps in and out of BaicellsQRTBWaitNotifyDPState as per state map"""

        mock_enodebd_update_cbsd.return_value = MOCK_CBSD_STATE

        acs_state_machine = provision_clean_sm(
            state='wait_get_transient_params',
        )

        msg = Tr069MessageBuilder.param_values_qrtb_response(
            GET_TRANSIENT_PARAMS_RESPONSE_PARAMS,
            models.GetParameterValuesResponse,
        )
        acs_state_machine.handle_tr069_message(msg)
        msg = Tr069MessageBuilder.param_values_qrtb_response(
            GET_PARAMS_RESPONSE_PARAMS,
            models.GetParameterValuesResponse,
        )
        acs_state_machine.handle_tr069_message(msg)
        # Transition to final get params check state
        acs_state_machine.transition('check_wait_get_params')

        # SM should transition from check_wait_get_params to end_session -> notify_dp automatically
        # upon receiving response from the radio
        msg = Tr069MessageBuilder.param_values_qrtb_response(
            [], models.GetParameterValuesResponse,
        )
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

    def test_cells_plmn_configured_in_postprocessor(self):
        acs_state_machine = provision_clean_sm()
        acs_state_machine.device_cfg.set_parameter(ParameterName.IP_SEC_ENABLE, 'false')

        desired_cfg = build_desired_config(
            acs_state_machine.mconfig,
            acs_state_machine.service_config,
            acs_state_machine.device_cfg,
            acs_state_machine.data_model,
            acs_state_machine.config_postprocessor,
        )

        self.assertEqual(desired_cfg.get_parameter_for_object(ParameterName.PLMN_N_CELL_RESERVED % 1, ParameterName.PLMN_N % 1), True)
        self.assertEqual(desired_cfg.get_parameter_for_object(ParameterName.PLMN_N_ENABLE % 1, ParameterName.PLMN_N % 1), True)
        self.assertEqual(desired_cfg.get_parameter_for_object(ParameterName.PLMN_N_PRIMARY % 1, ParameterName.PLMN_N % 1), True)
        self.assertEqual(desired_cfg.get_parameter(CarrierAggregationParameters.CA_PLMN_CELL_RESERVED), True)
        self.assertEqual(desired_cfg.get_parameter(CarrierAggregationParameters.CA_PLMN_ENABLE), True)
        self.assertEqual(desired_cfg.get_parameter(CarrierAggregationParameters.CA_PLMN_PRIMARY), False)

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


class BaicellsQRTBFirmwareUpgradeDownloadTests(EnodebHandlerTestCase):
    """
    Class for testing firmware upgrade download flow.

    Firmware upgrade download request should initiate in certain configurations.
    When initiated, a sequence of TR069 exchange needs to happen in order to
    schedule a download on the eNB.

    Firmware upgrade procedure on Baicells QRTB eNB starts at any time after
    eNB reports TransferComplete, at which point the eNB will reboot on its own.
    So we only test the TR069 message sequencing and configuration interpretation.
    TransferComplete message should appear after eNB reboot.
    """
    # helper variables
    _enb_serial = "baicells_serial_123"
    _enb_sw_version = "baicells_firmware_v0.0"
    _new_sw_version = "baicells_firmware_v1.0"
    _no_url_sw_version = "baicells_no_url_firmware"
    _sw_url = "http://fw_url/fw_file.ffw"

    _firmwares = {
        _enb_sw_version: {'url': _sw_url, 'md5': "12345678901234567890123456789012", 'rawmode': False},
        _new_sw_version: {'url': _sw_url, 'md5': "12345678901234567890123456789013", 'rawmode': False},
        _no_url_sw_version: {},
    }

    # configs which should not lead to firmware upgrade download flow
    config_empty: Dict[str, Dict[Any, Any]] = {'firmwares': {}, 'enbs': {}, 'models': {}}

    config_just_firmwares = deepcopy(config_empty)
    config_just_firmwares['firmwares'] = _firmwares

    config_just_enbs = deepcopy(config_empty)
    config_just_enbs['enbs'][_enb_serial] = _new_sw_version

    config_just_models = deepcopy(config_empty)
    config_just_models['models'][EnodebDeviceName.BAICELLS_QRTB] = _new_sw_version

    config_enb_fw_up_to_date = deepcopy(config_just_firmwares)
    config_enb_fw_up_to_date['enbs'][_enb_serial] = _enb_sw_version

    config_model_fw_up_to_date = deepcopy(config_just_firmwares)
    config_model_fw_up_to_date['models'][EnodebDeviceName.BAICELLS_QRTB] = _enb_sw_version

    config_enb_has_fw_version_without_url = deepcopy(config_just_firmwares)
    config_enb_has_fw_version_without_url['enbs'][_enb_serial] = _no_url_sw_version

    config_model_has_fw_version_without_url = deepcopy(config_just_firmwares)
    config_model_has_fw_version_without_url['models'][EnodebDeviceName.BAICELLS_QRTB] = _no_url_sw_version

    config_enb_fw_up_to_date_but_model_has_upgrade = deepcopy(
        config_enb_fw_up_to_date,
    )
    config_enb_fw_up_to_date_but_model_has_upgrade['models'][
        EnodebDeviceName.BAICELLS_QRTB
    ] = _new_sw_version

    # valid configs which should initiate fw upgrade
    config_enb_fw_upgrade = deepcopy(config_just_firmwares)
    config_enb_fw_upgrade['enbs'][_enb_serial] = _new_sw_version

    config_model_fw_upgrade = deepcopy(config_just_firmwares)
    config_model_fw_upgrade['models'][EnodebDeviceName.BAICELLS_QRTB] = _new_sw_version

    config_enb_fw_upgrade_but_model_fw_up_to_date = deepcopy(
        config_enb_fw_upgrade,
    )
    config_enb_fw_upgrade_but_model_fw_up_to_date['models'][
        EnodebDeviceName.BAICELLS_QRTB
    ] = _enb_sw_version

    @parameterized.expand([
        (config_empty,),
        (config_just_firmwares,),
        (config_just_enbs,),
        (config_just_models,),
        (config_enb_fw_up_to_date,),
        (config_model_fw_up_to_date,),
        (config_enb_has_fw_version_without_url,),
        (config_model_has_fw_version_without_url,),
        (config_enb_fw_up_to_date_but_model_has_upgrade,),
    ])
    def test_firmware_upgrade_download_flow_skip_on_config(self, fw_upgrade_download_config):
        """
        Test skipping firmware upgrade download flow.
        Skip should happen on certain firmware upgrade download configuraion conditions
        and eNB SW VERSION state.
        """
        logging.root.level = logging.DEBUG
        acs_state_machine = EnodebAcsStateMachineBuilder.build_acs_state_machine(
            EnodebDeviceName.BAICELLS_QRTB,
        )
        acs_state_machine._service.config = _get_service_config()
        acs_state_machine._service.config['firmware_upgrade_download'] = fw_upgrade_download_config

        # eNB sends Inform message, we wait for an InformResponse
        inform = Tr069MessageBuilder.get_inform(
            oui="48BF74",
            sw_version=self._enb_sw_version,
            enb_serial=self._enb_serial,
        )
        resp = acs_state_machine.handle_tr069_message(inform)
        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'Should respond with an InformResponse',
        )

        # eNB sends an empty http request
        # State machine should detect that no firmware upgrade config exists and so
        # should transition to getting param values skipping download flow
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request parameter values
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )

        # Firmware upgrade timeout timer should not be started
        self.assertFalse(acs_state_machine.is_fw_upgrade_in_progress())

    @parameterized.expand([
        (config_enb_fw_upgrade,),
        (config_model_fw_upgrade,),
        (config_enb_fw_upgrade_but_model_fw_up_to_date,),
    ])
    def test_firmware_upgrade_download_flow_skip_on_download_fault9017(self, fw_upgrade_download_config):
        """
        Test firmware upgrade download flow skip due to TR069 fault on Download request.
        State machine should resume normal operation when Fault code 9017 is received.
        """
        logging.root.level = logging.DEBUG
        acs_state_machine = EnodebAcsStateMachineBuilder.build_acs_state_machine(
            EnodebDeviceName.BAICELLS_QRTB,
        )
        acs_state_machine._service.config = _get_service_config()
        acs_state_machine._service.config['firmware_upgrade_download'] = fw_upgrade_download_config

        logging.info(f'{fw_upgrade_download_config=}')

        # eNB sends Inform message, we wait for an InformResponse
        inform = Tr069MessageBuilder.get_inform(
            oui="48BF74",
            sw_version=self._enb_sw_version,
            enb_serial=self._enb_serial,
        )
        resp = acs_state_machine.handle_tr069_message(inform)
        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'Should respond with an InformResponse',
        )

        # eNB sends an empty http request
        # State machine should detect that firmware upgrade config exists and so
        # should transition to sending Download message
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)
        self._assert_download_message(
            acs=acs_state_machine,
            message=resp,
        )

        # eNB may reply with a Fault code 9017 which appearts to mean that a Download
        # is already in progress on the eNB side (for instance the same CommandKey)
        # In such case, state machine should resume normal operation
        req = models.Fault()
        req.FaultCode = 9017
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )

        # Firmware upgrade timeout timer should not be started
        self.assertFalse(acs_state_machine.is_fw_upgrade_in_progress())

    @parameterized.expand([
        (config_enb_fw_upgrade,),
        (config_model_fw_upgrade,),
        (config_enb_fw_upgrade_but_model_fw_up_to_date,),
    ])
    def test_firmware_upgrade_download_flow(self, fw_upgrade_download_config):
        """
        Test firmware upgrade download flow.
        Download sequence should initiate on certain
        on
        """
        logging.root.level = logging.DEBUG
        acs_state_machine = EnodebAcsStateMachineBuilder.build_acs_state_machine(
            EnodebDeviceName.BAICELLS_QRTB,
        )
        acs_state_machine._service.config = _get_service_config()
        acs_state_machine._service.config['firmware_upgrade_download'] = fw_upgrade_download_config

        logging.info(f'{fw_upgrade_download_config=}')

        # eNB sends Inform message, we wait for an InformResponse
        inform = Tr069MessageBuilder.get_inform(
            oui="48BF74",
            sw_version=self._enb_sw_version,
            enb_serial=self._enb_serial,
        )
        resp = acs_state_machine.handle_tr069_message(inform)
        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'Should respond with an InformResponse',
        )

        # eNB sends an empty http request
        # State machine should detect that firmware upgrade config exists and so
        # should transition to sending Download message
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)
        self._assert_download_message(
            acs=acs_state_machine,
            message=resp,
        )

        # When eNB replies with a DownloadResponse, all is good.
        # eNB should transition to get params state and should start the upgrade
        # timeout timer
        req = models.DownloadResponse()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        self.assertTrue(acs_state_machine.is_fw_upgrade_in_progress())

    def _assert_download_message(
        self,
        acs,
        message,
    ):
        # Expect a dowload message
        self.assertTrue(
            isinstance(message, models.Download),
            'Expecting Download message',
        )
        # eNB firmware upgrade config should be obtainable
        fw_upgrade_config = get_firmware_upgrade_download_config(acs)
        self.assertTrue(
            fw_upgrade_config,
            'Firmware Upgrade config should not be empty',
        )

        # Explicitly set params should have correct values
        self.assertEqual(message.CommandKey, fw_upgrade_config['version'])
        self.assertEqual(
            message.FileType,
            FirmwareUpgradeDownloadState.FIRMWARE_FILE_TYPE,
        )
        self.assertEqual(message.URL, fw_upgrade_config['url'])

        # Optional params should have default values
        self.assertEqual(
            message.Username,
            fw_upgrade_config.get('username', ""),
        )
        self.assertEqual(
            message.Password,
            fw_upgrade_config.get('password', ""),
        )

        # Constant/Fixed params should have default values
        self.assertEqual(message.FileSize, 0)
        self.assertEqual(message.TargetFileName, "")
        self.assertEqual(message.DelaySeconds, 0)
        self.assertEqual(message.SuccessURL, "")
        self.assertEqual(message.FailureURL, "")


class BaicellsQRTBCarrierAggregationTests(EnodebHandlerTestCase):
    channel_20_1 = LteChannel(low_frequency_hz=3550_000_000, high_frequency_hz=3570_000_000, max_eirp_dbm_mhz=34)
    channel_20_2 = LteChannel(low_frequency_hz=3570_000_000, high_frequency_hz=3590_000_000, max_eirp_dbm_mhz=34)
    channel_20_3 = LteChannel(low_frequency_hz=3670_000_000, high_frequency_hz=3690_000_000, max_eirp_dbm_mhz=14)
    channel_15_1 = LteChannel(low_frequency_hz=3570_000_000, high_frequency_hz=3585_000_000, max_eirp_dbm_mhz=34)
    channel_15_2 = LteChannel(low_frequency_hz=3585_000_000, high_frequency_hz=3600_000_000, max_eirp_dbm_mhz=34)
    channel_10_1 = LteChannel(low_frequency_hz=3570_000_000, high_frequency_hz=3580_000_000, max_eirp_dbm_mhz=34)
    channel_10_2 = LteChannel(low_frequency_hz=3580_000_000, high_frequency_hz=3590_000_000, max_eirp_dbm_mhz=34)
    channel_10_3 = LteChannel(low_frequency_hz=3680_000_000, high_frequency_hz=3690_000_000, max_eirp_dbm_mhz=34)
    channel_5_1 = LteChannel(low_frequency_hz=3550_000_000, high_frequency_hz=3555_000_000, max_eirp_dbm_mhz=34)
    channel_5_2 = LteChannel(low_frequency_hz=3555_000_000, high_frequency_hz=3560_000_000, max_eirp_dbm_mhz=34)
    channel_5_3 = LteChannel(low_frequency_hz=3655_000_000, high_frequency_hz=3660_000_000, max_eirp_dbm_mhz=34)

    # Domain Proxy state results, which should be supported and accepted by QRTB. Carrier Aggregation should be set on eNB.
    state_ca_enabled_2_channels_20 = CBSDStateResult(channels=[channel_20_1, channel_20_2], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_enabled_2_channels_10 = CBSDStateResult(channels=[channel_10_1, channel_10_2], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_enabled_2_channels_5 = CBSDStateResult(channels=[channel_5_1, channel_5_2], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_enabled_2_channels_20_10 = CBSDStateResult(channels=[channel_20_1, channel_10_1], radio_enabled=True, carrier_aggregation_enabled=True)

    # Domain Proxy state results, which result in disabling the radio.
    state_ca_enabled_0_channels = CBSDStateResult(channels=[], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_disabled_0_channels = CBSDStateResult(channels=[], radio_enabled=True, carrier_aggregation_enabled=False)

    # Domain Proxy state results, with just 1 channel or 2 channels but CA disabled. Should result in eNB switching to Single Carrier
    state_ca_enabled_1_channel_20 = CBSDStateResult(channels=[channel_20_1], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_disabled_2_channels_20 = CBSDStateResult(channels=[channel_20_1, channel_20_2], radio_enabled=True, carrier_aggregation_enabled=False)
    state_ca_disabled_1_channel_20 = CBSDStateResult(channels=[channel_20_1], radio_enabled=True, carrier_aggregation_enabled=False)

    # Domain Proxy state results, with 2 channels that exceed max allowed IBW. Should result in eNB switching to Single Carrier
    state_ca_enabled_2_channels_20_ibw_over_100 = CBSDStateResult(channels=[channel_20_1, channel_20_3], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_enabled_2_channels_10_ibw_over_100 = CBSDStateResult(channels=[channel_10_1, channel_10_3], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_enabled_2_channels_5_ibw_over_100 = CBSDStateResult(channels=[channel_5_1, channel_5_3], radio_enabled=True, carrier_aggregation_enabled=True)

    # Domain Proxy state results, with 2 channels that are not in supported bandwidth configurations. Should result in eNB switching to Single Carrier
    state_ca_enabled_2_channels_15 = CBSDStateResult(channels=[channel_15_1, channel_15_2], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_enabled_2_channels_20_15 = CBSDStateResult(channels=[channel_20_1, channel_15_2], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_enabled_2_channels_20_5 = CBSDStateResult(channels=[channel_20_1, channel_5_2], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_enabled_2_channels_15_10 = CBSDStateResult(channels=[channel_15_1, channel_10_2], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_enabled_2_channels_15_5 = CBSDStateResult(channels=[channel_15_1, channel_5_2], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_enabled_2_channels_10_5 = CBSDStateResult(channels=[channel_10_1, channel_5_2], radio_enabled=True, carrier_aggregation_enabled=True)

    @parameterized.expand([
        (state_ca_enabled_2_channels_20,),
        (state_ca_enabled_2_channels_10,),
        (state_ca_enabled_2_channels_5,),
        (state_ca_enabled_2_channels_20_10,),
    ])
    def test_ca_enabled_and_fapservice_2_configured_based_on_domain_proxy_state(self, state) -> None:
        """
        Test eNB configuration set to Carrier Aggregation when Domain Proxy
        state response contains applicable channels
        """
        config = EnodebConfiguration(BaicellsQRTBTrDataModel())

        # Set required params in device configuration
        _pci = 260
        _cell_id = 138777000
        config.set_parameter(ParameterName.PCI, 260)
        config.set_parameter(ParameterName.CELL_ID, 138777000)

        # Update eNB configuration based on Domain Proxy state response
        qrtb_update_desired_config_from_cbsd_state(state, config)

        # Check eNB set to Carrier Aggregation
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_ENABLE), 1)
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_NUM_OF_CELLS), 2)

        # Check FAPService.2 set to parameters derived from second channel
        ca_channel_low_freq = state.channels[1].low_frequency_hz
        ca_channel_high_freq = state.channels[1].high_frequency_hz
        ca_bw_mhz = calc_bandwidth_mhz(ca_channel_low_freq, ca_channel_high_freq)
        ca_bw_rbs = calc_bandwidth_rbs(ca_bw_mhz)
        ca_earfcn = calc_earfcn(ca_channel_low_freq, ca_channel_high_freq)
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_DL_BANDWIDTH), ca_bw_rbs)
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_UL_BANDWIDTH), ca_bw_rbs)
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_EARFCNDL), ca_earfcn)
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_EARFCNDL), ca_earfcn)
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_BAND), 48)

        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_RADIO_ENABLE), True)
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_PCI), _pci + 1)
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_CELL_ID), _cell_id + 1)

    @parameterized.expand([
        (state_ca_enabled_0_channels,),
        (state_ca_disabled_0_channels,),
    ])
    def test_radio_disabled_based_on_domain_proxy_state(self, state) -> None:
        """
        Test eNB configuration set to disable radio tranmission and not touch
        any eNB parameters, when Domain Proxy state response contains no channels
        """
        config = EnodebConfiguration(BaicellsQRTBTrDataModel())
        qrtb_update_desired_config_from_cbsd_state(state, config)
        self.assertEqual(config.get_parameter(ParameterName.SAS_RADIO_ENABLE), False)

    @parameterized.expand([
        (state_ca_enabled_1_channel_20,),
        (state_ca_disabled_2_channels_20,),
        (state_ca_disabled_1_channel_20,),
        (state_ca_enabled_2_channels_20_ibw_over_100,),
        (state_ca_enabled_2_channels_10_ibw_over_100,),
        (state_ca_enabled_2_channels_5_ibw_over_100,),
        (state_ca_enabled_2_channels_20_15,),
        (state_ca_enabled_2_channels_20_5,),
        (state_ca_enabled_2_channels_15,),
        (state_ca_enabled_2_channels_15_10,),
        (state_ca_enabled_2_channels_15_5,),
        (state_ca_enabled_2_channels_10_5,),
    ])
    def test_single_carrier_enabled_based_on_domain_proxy_state(self, state) -> None:
        """
        Test eNB configuration set to Single Carrier when Domain Proxy
        state response contains one channel or 2 channels with incompatible
        channel configuration
        """
        config = EnodebConfiguration(BaicellsQRTBTrDataModel())

        # Set required params in device configuration
        config.set_parameter(ParameterName.GPS_STATUS, True)

        # Update eNB configuration based on Domain Proxy state response
        qrtb_update_desired_config_from_cbsd_state(state, config)

        # Check eNB set to Single Carrier
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_ENABLE), 0)
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_NUM_OF_CELLS), 1)

        # Check eNB FAPService.1 set to parameters derived from first channel
        sc_channel_low_freq = state.channels[0].low_frequency_hz
        sc_channel_high_freq = state.channels[0].high_frequency_hz
        sc_bw_mhz = calc_bandwidth_mhz(sc_channel_low_freq, sc_channel_high_freq)
        sc_bw_rbs = calc_bandwidth_rbs(sc_bw_mhz)
        sc_earfcn = calc_earfcn(sc_channel_low_freq, sc_channel_high_freq)
        self.assertEqual(config.get_parameter(ParameterName.DL_BANDWIDTH), sc_bw_rbs)
        self.assertEqual(config.get_parameter(ParameterName.UL_BANDWIDTH), sc_bw_rbs)
        self.assertEqual(config.get_parameter(ParameterName.EARFCNDL), sc_earfcn)
        self.assertEqual(config.get_parameter(ParameterName.EARFCNDL), sc_earfcn)
        self.assertEqual(config.get_parameter(ParameterName.SAS_RADIO_ENABLE), True)


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


def _get_service_config():
    return {
        "firmware_upgrade_download": {
            "enbs": {},
            "firmwares": {},
            "models": {},
        },
    }
