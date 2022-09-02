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
import logging
import os
from copy import deepcopy
from typing import Any, Dict
from unittest import TestCase
from unittest.mock import Mock, call, patch

from dp.protos.cbsd_pb2 import CBSDStateResult, LteChannel
from lte.protos.mconfig import mconfigs_pb2
from magma.common.service import MagmaService
from magma.enodebd.data_models.data_model_parameters import ParameterName
from magma.enodebd.device_config.cbrs_consts import BAND
from magma.enodebd.device_config.configuration_init import build_desired_config
from magma.enodebd.device_config.configuration_util import (
    calc_bandwidth_mhz,
    calc_bandwidth_rbs,
    calc_earfcn,
)
from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.devices.freedomfi_one.data_model import (
    FreedomFiOneTrDataModel,
)
from magma.enodebd.devices.freedomfi_one.impl import (
    SAS_KEY,
    FreedomFiOneConfigurationInitializer,
)
from magma.enodebd.devices.freedomfi_one.params import (
    CarrierAggregationParameters,
    FreedomFiOneMiscParameters,
    SASParameters,
    StatusParameters,
)
from magma.enodebd.devices.freedomfi_one.states import (
    FreedomFiOneEndSessionState,
    FreedomFiOneGetInitState,
    FreedomFiOneNotifyDPState,
    ff_one_update_desired_config_from_cbsd_state,
)
from magma.enodebd.dp_client import build_enodebd_update_cbsd_request
from magma.enodebd.exceptions import ConfigurationError
from magma.enodebd.state_machines.acs_state_utils import (
    get_firmware_upgrade_download_config,
)
from magma.enodebd.state_machines.enb_acs_states import (
    FirmwareUpgradeDownloadState,
    WaitInformMRebootState,
)
from magma.enodebd.tests.test_utils.config_builder import EnodebConfigBuilder
from magma.enodebd.tests.test_utils.enb_acs_builder import (
    EnodebAcsStateMachineBuilder,
)
from magma.enodebd.tests.test_utils.enodeb_handler import EnodebHandlerTestCase
from magma.enodebd.tests.test_utils.tr069_msg_builder import Tr069MessageBuilder
from magma.enodebd.tr069 import models
from parameterized import parameterized

magma_root = os.environ.get('MAGMA_ROOT')
if magma_root:
    SRC_CONFIG_DIR = os.path.join(
        magma_root,
        'lte/gateway/configs',
    )

MOCK_CHANNEL = LteChannel(
    low_frequency_hz=3550000000,
    high_frequency_hz=3570000000,
    max_eirp_dbm_mhz=15,
)
MOCK_CBSD_STATE = CBSDStateResult(
    radio_enabled=True,
    channel=MOCK_CHANNEL,
    channels=[MOCK_CHANNEL],
    carrier_aggregation_enabled=False,
)

TEST_SAS_URL = 'test_sas_url'
TEST_SAS_CERT_SUBJECT = 'test_sas_cert_subject'


class FreedomFiOneTests(EnodebHandlerTestCase):
    """Testing FreedomfiOne state machine"""

    def _get_ff_one_read_only_param_values_resp(
            self,
    ) -> models.GetParameterValuesResponse:
        msg = models.GetParameterValuesResponse()
        param_val_list = []
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.X_000E8F_DeviceFeature.X_000E8F_NEStatus'
                     '.X_000E8F_Sync_Status',
                val_type='string',
                data='InSync',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.X_000E8F_DeviceFeature.X_000E8F_NEStatus'
                     '.X_000E8F_SAS_Status',
                val_type='string',
                data='SUCCESS',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.X_000E8F_DeviceFeature.X_000E8F_NEStatus'
                     '.X_000E8F_eNB_Status',
                val_type='string',
                data='SUCCESS',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.X_000E8F_DeviceFeature.X_000E8F_NEStatus'
                     '.X_000E8F_DEFGW_Status',
                val_type='string',
                data='SUCCESS',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.FAP.GPS.ScanStatus',
                val_type='string',
                data='SUCCESS',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.FAP.GPS.LockedLongitude',
                val_type='int',
                data='-105272892',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.FAP.GPS.LockedLatitude',
                val_type='int',
                data='40019606',
            ),
        )
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = param_val_list
        return msg

    def _get_freedomfi_one_param_values_response(self):
        msg = models.GetParameterValuesResponse()
        param_val_list = []
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = param_val_list

        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.EARFCNDL',
                val_type='int',
                data='56240',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.FAP.GPS.ScanOnBoot',
                val_type='boolean',
                data='1',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.AdminState',
                val_type='boolean',
                data='1',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.FAP.PerfMgmt.Config.1.Enable',
                val_type='boolean',
                data='1',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.Gateway.S1SigLinkServerList',
                val_type='string',
                data='10.0.2.1',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.EPC.TAC',
                val_type='int',
                data='1',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.FAP.PerfMgmt.Config.1.PeriodicUploadInterval',
                val_type='int',
                data='60',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.DeviceInfo.SoftwareVersion',
                val_type='string',
                data='TEST3920@210901',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.X_000E8F_SAS.FCCIdentificationNumber',
                val_type='string',
                data='P27-SCE4255W',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.X_000E8F_SAS.UserContactInformation',
                val_type='string',
                data='M0LK4T',  # TODO do not take it from the radio. Embed it in config somewhere
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.X_000E8F_SAS.ProtectionLevel',
                val_type='string',
                data='GAA',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.X_000E8F_SAS.CertSubject',
                val_type='string',
                data=TEST_SAS_CERT_SUBJECT,
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.X_000E8F_SAS.HeightType',
                val_type='string',
                data='AMSL',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.X_000E8F_SAS.Category',
                val_type='string',
                data='A',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.X_000E8F_SAS.AntennaGain',
                val_type='int',
                data='5',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.FAP.GPS.ScanStatus',
                val_type='string',
                data='Success',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.ManagementServer.PeriodicInformInterval',
                val_type='int',
                data='60',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.FreqBandIndicator',
                val_type='unsignedInt',
                data='48',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.Common.CellIdentity',
                val_type='unsignedInt',
                data='101',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.FAP.GPS.LockedLongitude',
                val_type='string',
                data='-105272892',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.FAP.GPS.LockedLatitude',
                val_type='string',
                data='40019606',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.X_000E8F_SAS.CPIEnable',
                val_type='boolean',
                data='0',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.X_000E8F_RRMConfig.X_000E8F_CA_Enable',
                val_type='boolean',
                data='0',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.X_000E8F_RRMConfig.X_000E8F_Cell_Number',
                val_type='int',
                data='1',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.FAP.GPS.ScanOnBoot',
                val_type='boolean',
                data='1',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.X_000E8F_DeviceFeature.X_000E8F_NEStatus.X_000E8F_Sync_Status',
                val_type='string',
                data='InSync',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.PhyCellID',
                val_type='string',
                data='101,102',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.FAP.PerfMgmt.Config.1.URL',
                val_type='string',
                data=None,
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.X_000E8F_SAS.Location',
                val_type='string',
                data='indoor',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.PHY.TDDFrame.SubFrameAssignment',
                val_type='boolean',
                data='2',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.IsPrimary',
                val_type='boolean',
                data='1',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.Enable',
                val_type='boolean',
                data='1',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.X_000E8F_SAS.Server',
                val_type='string',
                data=TEST_SAS_URL,
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.X_000E8F_DeviceFeature.X_000E8F_WebServerEnable',
                val_type='boolean',
                data='1',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.CellReservedForOperatorUse',
                val_type='boolean',
                data='0',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.Tunnel.1.TunnelRef',
                val_type='string',
                data='Device.IP.Interface.1.IPv4Address.1.',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.REM.X_000E8F_tfcsManagerConfig.primSrc',
                val_type='string',
                data='GNSS',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.X_000E8F_SAS.Enable',
                val_type='boolean',
                data='1',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.X_000E8F_SAS.Method',
                val_type='boolean',
                data='0',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.ManagementServer.PeriodicInformEnable',
                val_type='boolean',
                data='1',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNListNumberOfEntries',
                val_type='int',
                data='1',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.PHY.TDDFrame.SpecialSubframePatterns',
                val_type='int',
                data='7',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.X_000E8F_RRMConfig.X_000E8F_CELL_Freq_Contiguous',
                val_type='int',
                data='0',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.Gateway.S1SigLinkPort',
                val_type='int',
                data='36412',
            ),
        )
        param_val_list.append(
            Tr069MessageBuilder.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.X_000E8F_SAS.MaxEirpMHz_Carrier1',
                val_type='int',
                data='24',
            ),
        )
        return msg

    @patch('magma.enodebd.devices.freedomfi_one.states.enodebd_update_cbsd')
    def test_provision(self, mock_enodebd_update_cbsd) -> None:
        """
        Test the basic provisioning workflow
        1 - enodeb sends Inform, enodebd sends InformResponse
        2 - enodeb sends empty HTTP message,
        3 - enodebd sends get transient params, updates the device state.
        4 - enodebd sends get param values, updates the device state
        5 - enodebd, sets fields including SAS fields.

        Args:
            mock_enodebd_update_cbsd (Any): mocking enodebd_update_cbsd method
        """
        self.maxDiff = None
        test_serial_number = "2006CW5000023"
        mock_enodebd_update_cbsd.return_value = MOCK_CBSD_STATE

        logging.root.level = logging.DEBUG
        acs_state_machine = EnodebAcsStateMachineBuilder.build_acs_state_machine(
            EnodebDeviceName.FREEDOMFI_ONE,
        )
        acs_state_machine._service.config = _get_service_config()
        acs_state_machine.desired_cfg = build_desired_config(
            acs_state_machine.mconfig,
            acs_state_machine.service_config,
            acs_state_machine.device_cfg,
            acs_state_machine.data_model,
            acs_state_machine.config_postprocessor,
        )

        inform = Tr069MessageBuilder.get_inform(
            oui="000E8F",
            sw_version="TEST3920@210901",
            enb_serial=test_serial_number,
        )
        resp = acs_state_machine.handle_tr069_message(inform)
        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'Should respond with an InformResponse',
        )

        # Send an empty http request
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for read-only params
        self.assertTrue(
            isinstance(resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        for tr69nodes in StatusParameters.STATUS_PARAMETERS.values():
            self.assertIn(tr69nodes.path, resp.ParameterNames.string)

        req = self._get_ff_one_read_only_param_values_resp()
        get_resp = acs_state_machine.handle_tr069_message(req)

        self.assertTrue(
            isinstance(get_resp, models.GetParameterValues),
            'State machine should be requesting param values',
        )
        req = self._get_freedomfi_one_param_values_response()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.SetParameterValues),
            'State machine should be setting parameters',
        )
        self.assertIsNotNone(
            resp.ParameterKey.Data,
            'ParameterKey should be set for FreedomFiOne eNB',
        )

        msg = models.SetParameterValuesResponse()
        msg.Status = 1
        get_resp = acs_state_machine.handle_tr069_message(msg)
        self.assertTrue(
            isinstance(get_resp, models.GetParameterValues),
            'We should read back all parameters',
        )

        req = self._get_freedomfi_one_param_values_response()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(
            isinstance(resp, models.DummyInput),
            'Provisioning completed with Dummy response',
        )

        enodebd_update_cbsd_request = build_enodebd_update_cbsd_request(
            serial_number=acs_state_machine.device_cfg.get_parameter(ParameterName.SERIAL_NUMBER),
            latitude_deg=acs_state_machine.device_cfg.get_parameter(ParameterName.GPS_LAT),
            longitude_deg=acs_state_machine.device_cfg.get_parameter(ParameterName.GPS_LONG),
            indoor_deployment=acs_state_machine.device_cfg.get_parameter(SASParameters.SAS_LOCATION),
            antenna_height='0',
            antenna_height_type=acs_state_machine.device_cfg.get_parameter(SASParameters.SAS_HEIGHT_TYPE),
            cbsd_category=acs_state_machine.device_cfg.get_parameter(SASParameters.SAS_CATEGORY),
        )

        mock_enodebd_update_cbsd.assert_called_with(enodebd_update_cbsd_request)

    @patch('magma.enodebd.devices.freedomfi_one.states.enodebd_update_cbsd')
    def test_enodebd_update_cbsd_not_called_when_gps_unavailable(self, mock_enodebd_update_cbsd) -> None:
        """
        Test enodebd does not call Domain Proxy API for parameter update when
        not all parameters are yet available from the eNB
        """
        self.maxDiff = None
        test_serial_number = "2006CW5000023"
        mock_enodebd_update_cbsd.return_value = MOCK_CBSD_STATE

        acs_state_machine = EnodebAcsStateMachineBuilder.build_acs_state_machine(
            EnodebDeviceName.FREEDOMFI_ONE,
        )
        acs_state_machine._service.config = _get_service_config()
        acs_state_machine.desired_cfg = build_desired_config(
            acs_state_machine.mconfig,
            acs_state_machine.service_config,
            acs_state_machine.device_cfg,
            acs_state_machine.data_model,
            acs_state_machine.config_postprocessor,
        )

        inform = Tr069MessageBuilder.get_inform(
            oui="000E8F",
            sw_version="TEST3920@210901",
            enb_serial=test_serial_number,
        )
        resp = acs_state_machine.handle_tr069_message(inform)
        self.assertTrue(
            isinstance(resp, models.InformResponse),
            'Should respond with an InformResponse',
        )

        # No GPS locked means LAT and LONG are not available yet
        acs_state_machine.device_cfg.set_parameter(ParameterName.GPS_STATUS, False)
        acs_state_machine.transition('notify_dp')
        mock_enodebd_update_cbsd.assert_not_called()


class FreedomFiOneStatesTests(EnodebHandlerTestCase):
    """Testing FreedomfiOne specific states"""

    @parameterized.expand([
        (True, FreedomFiOneNotifyDPState),
        (False, FreedomFiOneEndSessionState),
    ])
    @patch('magma.enodebd.devices.freedomfi_one.states.enodebd_update_cbsd')
    def test_transition_depending_on_sas_enabled_flag(
            self, dp_mode, expected_state, mock_enodebd_update_cbsd,
    ):
        """Testing if SM steps in and out of FreedomFiOneWaitNotifyDPState as per state map depending on whether
        sas_enabled param is set to True or False in the service config

        Args:
            dp_mode: bool flag to enable or disable dp mode
            expected_state (Any): State
        """

        mock_enodebd_update_cbsd.return_value = MOCK_CBSD_STATE

        acs_state_machine = EnodebAcsStateMachineBuilder.build_acs_state_machine(
            EnodebDeviceName.FREEDOMFI_ONE,
        )
        acs_state_machine._service.config = _get_service_config(
            dp_mode=dp_mode,
        )
        acs_state_machine.desired_cfg = build_desired_config(
            acs_state_machine.mconfig,
            acs_state_machine.service_config,
            acs_state_machine.device_cfg,
            acs_state_machine.data_model,
            acs_state_machine.config_postprocessor,
        )

        # Need to fill these values in the device_cfg if we're going to transition to notify_dp state
        acs_state_machine.device_cfg.set_parameter(
            SASParameters.SAS_USER_ID, 'test_user',
        )
        acs_state_machine.device_cfg.set_parameter(
            SASParameters.SAS_FCC_ID, 'test_fcc',
        )
        acs_state_machine.device_cfg.set_parameter(
            ParameterName.SERIAL_NUMBER, 'test_sn',
        )
        acs_state_machine.device_cfg.set_parameter(
            ParameterName.GPS_LAT, '10',
        )
        acs_state_machine.device_cfg.set_parameter(
            ParameterName.GPS_LONG, '-10',
        )
        acs_state_machine.device_cfg.set_parameter(
            SASParameters.SAS_CATEGORY, 'A',
        )
        acs_state_machine.device_cfg.set_parameter(
            SASParameters.SAS_HEIGHT_TYPE, 'AMSL',
        )
        acs_state_machine.device_cfg.set_parameter(
            SASParameters.SAS_LOCATION, 'indoor',
        )
        acs_state_machine.device_cfg.set_parameter(
            SASParameters.SAS_ANTENNA_GAIN, '5',
        )
        acs_state_machine.device_cfg.set_parameter(
            ParameterName.GPS_STATUS, 'True',
        )
        acs_state_machine.transition('check_wait_get_params')

        msg = Tr069MessageBuilder.param_values_qrtb_response(
            [], models.GetParameterValuesResponse,
        )

        # SM should transition from check_wait_get_params to end_session -> notify_dp automatically
        # upon receiving response from the radio
        acs_state_machine.handle_tr069_message(msg)

        self.assertIsInstance(acs_state_machine.state, expected_state)

        msg = Tr069MessageBuilder.get_inform(event_codes=['1 BOOT'])

        # SM should go into wait_inform state, respond with Inform response and transition to FreedomFiOneGetInitState
        acs_state_machine.handle_tr069_message(msg)

        self.assertIsInstance(
            acs_state_machine.state,
            FreedomFiOneGetInitState,
        )

    def test_manual_reboot_during_provisioning(self) -> None:
        """
        Test a scenario where a Magma user goes through the enodebd CLI to
        reboot the Sercomm eNodeB.

        This checks the scenario where the command is sent in the middle
        of a TR-069 provisioning session.
        """
        logging.root.level = logging.DEBUG
        acs_state_machine = EnodebAcsStateMachineBuilder.build_acs_state_machine(
            EnodebDeviceName.FREEDOMFI_ONE,
        )

        # Send an Inform message, wait for an InformResponse
        inform = Tr069MessageBuilder.get_inform(
            oui="000E8F",
            sw_version="TEST3920@210901",
            enb_serial="2006CW5000023",
        )
        resp = acs_state_machine.handle_tr069_message(inform)
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
        self.assertIsInstance(acs_state_machine.state, WaitInformMRebootState)
        inform = Tr069MessageBuilder.get_inform(
            event_codes=["M Reboot"],
        )
        resp = acs_state_machine.handle_tr069_message(inform)
        self.assertIsInstance(resp, models.InformResponse)
        self.assertIsInstance(
            acs_state_machine.state,
            FreedomFiOneGetInitState,
        )

    def test_post_processing_in_dp_mode(self) -> None:
        """ Test FreedomFi One specific post processing functionality in Domain Proxy mode"""

        service_cfg = _get_service_config()
        expected = [
            call.delete_parameter(ParameterName.EARFCNDL),
            call.delete_parameter(ParameterName.DL_BANDWIDTH),
            call.delete_parameter(ParameterName.UL_BANDWIDTH),
            call.set_parameter(
                param_name=FreedomFiOneMiscParameters.TUNNEL_REF,
                value='Device.IP.Interface.1.IPv4Address.1.',
            ),
            call.set_parameter(
                param_name=FreedomFiOneMiscParameters.CONTIGUOUS_CC, value=0,
            ),
            call.set_parameter(
                param_name=FreedomFiOneMiscParameters.WEB_UI_ENABLE, value=False,
            ),
            call.set_parameter(
                param_name=SASParameters.SAS_ENABLE, value=True,
            ),
            call.set_parameter(
                param_name=SASParameters.SAS_METHOD, value=False,
            ),
            call.set_parameter(
                param_name=CarrierAggregationParameters.CA_ENABLE, value=False,
            ),
            call.set_parameter(
                param_name=CarrierAggregationParameters.CA_CARRIER_NUMBER, value=1,
            ),
            call.set_parameter_for_object(
                param_name='PLMN 1 cell reserved',
                value=True, object_name='PLMN 1',
            ),
            call.set_parameter(SASParameters.SAS_METHOD, value=True),
            call.set_parameter(FreedomFiOneMiscParameters.PRIM_SOURCE, 'GNSS'),
        ]
        self._check_postprocessing(expected=expected, service_cfg=service_cfg)

    def test_post_processing_in_non_dp_mode(self) -> None:
        """ Test FreedomFi One specific post processing functionality in standalone mode"""
        service_cfg = _get_service_config(dp_mode=False)
        expected = [
            call.delete_parameter(ParameterName.EARFCNDL),
            call.delete_parameter(ParameterName.DL_BANDWIDTH),
            call.delete_parameter(ParameterName.UL_BANDWIDTH),
            call.set_parameter(
                param_name=FreedomFiOneMiscParameters.TUNNEL_REF,
                value='Device.IP.Interface.1.IPv4Address.1.',
            ),
            call.set_parameter(
                param_name=FreedomFiOneMiscParameters.CONTIGUOUS_CC, value=0,
            ),
            call.set_parameter(
                param_name=FreedomFiOneMiscParameters.WEB_UI_ENABLE, value=False,
            ),
            call.set_parameter(
                param_name=SASParameters.SAS_ENABLE, value=True,
            ),
            call.set_parameter(
                param_name=SASParameters.SAS_METHOD, value=False,
            ),
            call.set_parameter(
                param_name=CarrierAggregationParameters.CA_ENABLE, value=False,
            ),
            call.set_parameter(
                param_name=CarrierAggregationParameters.CA_CARRIER_NUMBER, value=1,
            ),
            call.set_parameter_for_object(
                param_name='PLMN 1 cell reserved',
                value=True, object_name='PLMN 1',
            ),
            call.set_parameter(
                SASParameters.SAS_SERVER_URL,
                TEST_SAS_URL,
            ),
            call.set_parameter(SASParameters.SAS_UID, 'M0LK4T'),
            call.set_parameter(SASParameters.SAS_CATEGORY, 'A'),
            call.set_parameter(SASParameters.SAS_CHANNEL_TYPE, 'GAA'),
            call.set_parameter(
                SASParameters.SAS_CERT_SUBJECT,
                TEST_SAS_CERT_SUBJECT,
            ),
            call.set_parameter(SASParameters.SAS_LOCATION, 'indoor'),
            call.set_parameter(SASParameters.SAS_HEIGHT_TYPE, 'AMSL'),

            call.set_parameter(FreedomFiOneMiscParameters.PRIM_SOURCE, 'GNSS'),
        ]

        self._check_postprocessing(expected=expected, service_cfg=service_cfg)

    def test_post_processing_without_sas_config(self) -> None:
        """ Test FreedomFi One specific post processing functionality without sas config"""
        service_cfg = {
            "tr069": {
                "interface": "eth1",
                "port": 48080,
                "perf_mgmt_port": 8081,
                "public_ip": "192.88.99.142",
            },
            "prim_src": 'GNSS',
            "reboot_enodeb_on_mme_disconnected": True,
            "s1_interface": "eth1",
        }
        expected = [
            call.delete_parameter(ParameterName.BAND),
            call.delete_parameter(ParameterName.EARFCNDL),
            call.delete_parameter(ParameterName.DL_BANDWIDTH),
            call.delete_parameter(ParameterName.UL_BANDWIDTH),
            call.set_parameter(
                param_name=FreedomFiOneMiscParameters.TUNNEL_REF,
                value='Device.IP.Interface.1.IPv4Address.1.',
            ),
            call.set_parameter(
                param_name=FreedomFiOneMiscParameters.CONTIGUOUS_CC, value=0,
            ),
            call.set_parameter(
                param_name=FreedomFiOneMiscParameters.WEB_UI_ENABLE, value=False,
            ),
            call.set_parameter(
                param_name=SASParameters.SAS_ENABLE, value=True,
            ),
            call.set_parameter(
                param_name=SASParameters.SAS_METHOD, value=False,
            ),
            call.set_parameter(
                param_name=CarrierAggregationParameters.CA_ENABLE, value=False,
            ),
            call.set_parameter(
                param_name=CarrierAggregationParameters.CA_CARRIER_NUMBER, value=1,
            ),
            call.set_parameter_for_object(
                param_name='PLMN 1 cell reserved',
                value=True, object_name='PLMN 1',
            ),
            call.set_parameter(SASParameters.SAS_METHOD, value=True),
            call.set_parameter(FreedomFiOneMiscParameters.PRIM_SOURCE, 'GNSS'),
        ]

        self._check_postprocessing(expected=expected, service_cfg=service_cfg)

    def test_post_process_without_sas_cfg_with_ui(self) -> None:
        """ Test FreedomFi One specific post processing functionality without sas config with ui enabled"""
        service_cfg = {
            "tr069": {
                "interface": "eth1",
                "port": 48080,
                "perf_mgmt_port": 8081,
                "public_ip": "192.88.99.142",
            },
            "prim_src": 'GNSS',
            "reboot_enodeb_on_mme_disconnected": True,
            "s1_interface": "eth1",
            "web_ui_enable_list": ["2006CW5000023"],
        }

        expected = [
            call.delete_parameter(ParameterName.BAND),
            call.delete_parameter(ParameterName.EARFCNDL),
            call.delete_parameter(ParameterName.DL_BANDWIDTH),
            call.delete_parameter(ParameterName.UL_BANDWIDTH),
            call.set_parameter(
                param_name=FreedomFiOneMiscParameters.TUNNEL_REF,
                value='Device.IP.Interface.1.IPv4Address.1.',
            ),
            call.set_parameter(
                param_name=FreedomFiOneMiscParameters.CONTIGUOUS_CC, value=0,
            ),
            call.set_parameter(
                param_name=FreedomFiOneMiscParameters.WEB_UI_ENABLE, value=False,
            ),
            call.set_parameter(
                param_name=SASParameters.SAS_ENABLE, value=True,
            ),
            call.set_parameter(
                param_name=SASParameters.SAS_METHOD, value=False,
            ),
            call.set_parameter(
                param_name=CarrierAggregationParameters.CA_ENABLE, value=False,
            ),
            call.set_parameter(
                param_name=CarrierAggregationParameters.CA_CARRIER_NUMBER, value=1,
            ),
            call.set_parameter_for_object(
                param_name='PLMN 1 cell reserved',
                value=True, object_name='PLMN 1',
            ),
            call.set_parameter(
                FreedomFiOneMiscParameters.WEB_UI_ENABLE, value=True,
            ),
            call.set_parameter(SASParameters.SAS_METHOD, value=True),
            call.set_parameter(FreedomFiOneMiscParameters.PRIM_SOURCE, 'GNSS'),
        ]

        self._check_postprocessing(expected=expected, service_cfg=service_cfg)

    def _check_postprocessing(self, expected, service_cfg):
        cfg_desired = Mock()
        acs_state_machine = EnodebAcsStateMachineBuilder.build_acs_state_machine(
            EnodebDeviceName.FREEDOMFI_ONE,
        )
        acs_state_machine.device_cfg.set_parameter(
            ParameterName.SERIAL_NUMBER,
            "2006CW5000023",
        )

        cfg_init = FreedomFiOneConfigurationInitializer(acs_state_machine)
        cfg_init.postprocess(
            EnodebConfigBuilder.get_mconfig(),
            service_cfg,
            cfg_desired,
        )

        cfg_desired.assert_has_calls(expected)

    @patch('magma.configuration.service_configs.CONFIG_DIR', SRC_CONFIG_DIR)
    def test_service_cfg_parsing(self):
        """ Test the parsing of the service config file for enodebd.yml"""
        self.maxDiff = None
        service = MagmaService('enodebd', mconfigs_pb2.EnodebD())
        service_cfg = service.config
        service_cfg["sas"]["sas_server_url"] = TEST_SAS_URL
        service_cfg1 = _get_service_config()
        service_cfg1['web_ui_enable_list'] = []
        service_cfg1['prim_src'] = 'GNSS'
        service_cfg1[SAS_KEY][SASParameters.SAS_UID] = 'INVALID_ID'
        service_cfg1[SAS_KEY][SASParameters.SAS_CERT_SUBJECT] = 'INVALID_CERT_SUBJECT'
        service_cfg1['print_grpc_payload'] = False
        self.assertDictEqual(service_cfg, service_cfg1)

    def test_status_nodes(self):
        """ Test that the status of the node is valid"""
        status = StatusParameters()

        # Happy path
        n1 = {
            StatusParameters.DEFAULT_GW: "SUCCESS",
            StatusParameters.SYNC_STATUS: "InSync",
            StatusParameters.ENB_STATUS: "Success",
            StatusParameters.SAS_STATUS: "Success",
            StatusParameters.GPS_SCAN_STATUS: "SUCCESS",
            ParameterName.GPS_LONG: "1",
            ParameterName.GPS_LAT: "1",
        }
        device_config = Mock()
        status.set_magma_device_cfg(n1, device_config)
        expected = [
            call.set_parameter(param_name='RF TX status', value=True),
            call.set_parameter(param_name='GPS status', value=True),
            call.set_parameter(param_name='PTP status', value=True),
            call.set_parameter(param_name='MME status', value=True),
            call.set_parameter(param_name='Opstate', value=True),
            call.set_parameter('GPS lat', '1'),
            call.set_parameter('GPS long', '1'),
        ]
        self.assertEqual(expected, device_config.mock_calls)

        n2 = n1.copy()
        # Verify we can handle specific none params
        n2[StatusParameters.DEFAULT_GW] = None
        n3 = n1.copy()
        n3[StatusParameters.SYNC_STATUS] = None
        n4 = n1.copy()
        n4[StatusParameters.ENB_STATUS] = None
        n5 = n1.copy()
        n5[StatusParameters.SAS_STATUS] = None
        n6 = n1.copy()
        n6[StatusParameters.GPS_SCAN_STATUS] = None
        n7 = n1.copy()
        n7[ParameterName.GPS_LONG] = None
        n8 = n1.copy()
        n8[ParameterName.GPS_LAT] = None

        device_config = Mock()
        expected = [
            call.set_parameter(param_name='RF TX status', value=True),
            call.set_parameter(param_name='GPS status', value=True),
            call.set_parameter(param_name='PTP status', value=True),
            call.set_parameter(param_name='MME status', value=False),
            call.set_parameter(param_name='Opstate', value=True),
            call.set_parameter('GPS lat', '1'),
            call.set_parameter('GPS long', '1'),
        ]
        status.set_magma_device_cfg(n2, device_config)
        self.assertEqual(expected, device_config.mock_calls)

        device_config = Mock()
        expected = [
            call.set_parameter(param_name='RF TX status', value=True),
            call.set_parameter(param_name='GPS status', value=True),
            call.set_parameter(param_name='PTP status', value=False),
            call.set_parameter(param_name='MME status', value=True),
            call.set_parameter(param_name='Opstate', value=True),
            call.set_parameter('GPS lat', '1'),
            call.set_parameter('GPS long', '1'),
        ]
        status.set_magma_device_cfg(n3, device_config)
        self.assertEqual(expected, device_config.mock_calls)

        device_config = Mock()
        expected = [
            call.set_parameter(param_name='RF TX status', value=True),
            call.set_parameter(param_name='GPS status', value=True),
            call.set_parameter(param_name='PTP status', value=True),
            call.set_parameter(param_name='MME status', value=True),
            call.set_parameter(param_name='Opstate', value=False),
            call.set_parameter('GPS lat', '1'),
            call.set_parameter('GPS long', '1'),
        ]
        status.set_magma_device_cfg(n4, device_config)
        self.assertEqual(expected, device_config.mock_calls)

        device_config = Mock()
        expected = [
            call.set_parameter(param_name='RF TX status', value=False),
            call.set_parameter(param_name='GPS status', value=True),
            call.set_parameter(param_name='PTP status', value=True),
            call.set_parameter(param_name='MME status', value=True),
            call.set_parameter(param_name='Opstate', value=True),
            call.set_parameter('GPS lat', '1'),
            call.set_parameter('GPS long', '1'),
        ]
        status.set_magma_device_cfg(n5, device_config)
        self.assertEqual(expected, device_config.mock_calls)

        device_config = Mock()
        expected = [
            call.set_parameter(param_name='RF TX status', value=True),
            call.set_parameter(param_name='GPS status', value=False),
            call.set_parameter(param_name='PTP status', value=False),
            call.set_parameter(param_name='MME status', value=True),
            call.set_parameter(param_name='Opstate', value=True),
            call.set_parameter('GPS lat', '1'),
            call.set_parameter('GPS long', '1'),
        ]
        status.set_magma_device_cfg(n6, device_config)
        self.assertEqual(expected, device_config.mock_calls)

        device_config = Mock()
        expected = [
            call.set_parameter(param_name='RF TX status', value=True),
            call.set_parameter(param_name='GPS status', value=True),
            call.set_parameter(param_name='PTP status', value=True),
            call.set_parameter(param_name='MME status', value=True),
            call.set_parameter(param_name='Opstate', value=True),
            call.set_parameter('GPS lat', '1'),
            call.set_parameter('GPS long', None),
        ]
        status.set_magma_device_cfg(n7, device_config)
        self.assertEqual(expected, device_config.mock_calls)

        device_config = Mock()
        expected = [
            call.set_parameter(param_name='RF TX status', value=True),
            call.set_parameter(param_name='GPS status', value=True),
            call.set_parameter(param_name='PTP status', value=True),
            call.set_parameter(param_name='MME status', value=True),
            call.set_parameter(param_name='Opstate', value=True),
            call.set_parameter('GPS lat', None),
            call.set_parameter('GPS long', '1'),
        ]
        status.set_magma_device_cfg(n8, device_config)
        self.assertEqual(expected, device_config.mock_calls)


class FreedomFiOneFirmwareUpgradeDownloadTests(EnodebHandlerTestCase):
    """
    Class for testing firmware upgrade download flow.

    Firmware upgrade download request should initiate in certain configurations.
    When initiated, a sequence of TR069 exchange needs to happen in order to
    schedule a download on the eNB.

    Firmware upgrade procedure on FreedomFi one eNB starts at any time after
    eNB reports TransferComplete, at which point the eNB will reboot on its own.
    So we only test the TR069 message sequencing and configuration interpretation.
    """
    # helper variables
    _enb_serial = "sercomm_serial_123"
    _enb_sw_version = "sercomm_firmware_v0.0"
    _new_sw_version = "sercomm_firmware_v1.0"
    _no_url_sw_version = "sercomm_no_url_firmware"
    _sw_url = "http://fw_url/fw_file.ffw"

    _firmwares = {
        _enb_sw_version: {'url': _sw_url},
        _new_sw_version: {'url': _sw_url},
        _no_url_sw_version: {},
    }

    # configs which should not lead to firmware upgrade download flow
    config_empty: Dict[str, Dict[Any, Any]] = {'firmwares': {}, 'enbs': {}, 'models': {}}

    config_just_firmwares = deepcopy(config_empty)
    config_just_firmwares['firmwares'] = _firmwares

    config_just_enbs = deepcopy(config_empty)
    config_just_enbs['enbs'][_enb_serial] = _new_sw_version

    config_just_models = deepcopy(config_empty)
    config_just_models['models'][EnodebDeviceName.FREEDOMFI_ONE] = _new_sw_version

    config_enb_fw_up_to_date = deepcopy(config_just_firmwares)
    config_enb_fw_up_to_date['enbs'][_enb_serial] = _enb_sw_version

    config_model_fw_up_to_date = deepcopy(config_just_firmwares)
    config_model_fw_up_to_date['models'][EnodebDeviceName.FREEDOMFI_ONE] = _enb_sw_version

    config_enb_has_fw_version_without_url = deepcopy(config_just_firmwares)
    config_enb_has_fw_version_without_url['enbs'][_enb_serial] = _no_url_sw_version

    config_model_has_fw_version_without_url = deepcopy(config_just_firmwares)
    config_model_has_fw_version_without_url['models'][EnodebDeviceName.FREEDOMFI_ONE] = _no_url_sw_version

    config_enb_fw_up_to_date_but_model_has_upgrade = deepcopy(
        config_enb_fw_up_to_date,
    )
    config_enb_fw_up_to_date_but_model_has_upgrade['models'][
        EnodebDeviceName.FREEDOMFI_ONE
    ] = _new_sw_version

    # valid configs which should initiate fw upgrade
    config_enb_fw_upgrade = deepcopy(config_just_firmwares)
    config_enb_fw_upgrade['enbs'][_enb_serial] = _new_sw_version

    config_model_fw_upgrade = deepcopy(config_just_firmwares)
    config_model_fw_upgrade['models'][EnodebDeviceName.FREEDOMFI_ONE] = _new_sw_version

    config_enb_fw_upgrade_but_model_fw_up_to_date = deepcopy(
        config_enb_fw_upgrade,
    )
    config_enb_fw_upgrade_but_model_fw_up_to_date['models'][
        EnodebDeviceName.FREEDOMFI_ONE
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
            EnodebDeviceName.FREEDOMFI_ONE,
        )
        acs_state_machine._service.config = _get_service_config()
        acs_state_machine._service.config['firmware_upgrade_download'] = fw_upgrade_download_config

        # eNB sends Inform message, we wait for an InformResponse
        inform = Tr069MessageBuilder.get_inform(
            oui="000E8F",
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
            EnodebDeviceName.FREEDOMFI_ONE,
        )
        acs_state_machine._service.config = _get_service_config()
        acs_state_machine._service.config['firmware_upgrade_download'] = fw_upgrade_download_config

        logging.info(f'{fw_upgrade_download_config=}')

        # eNB sends Inform message, we wait for an InformResponse
        inform = Tr069MessageBuilder.get_inform(
            oui="000E8F",
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
            EnodebDeviceName.FREEDOMFI_ONE,
        )
        acs_state_machine._service.config = _get_service_config()
        acs_state_machine._service.config['firmware_upgrade_download'] = fw_upgrade_download_config

        logging.info(f'{fw_upgrade_download_config=}')

        # eNB sends Inform message, we wait for an InformResponse
        inform = Tr069MessageBuilder.get_inform(
            oui="000E8F",
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


class TXParamsTests(TestCase):
    """Testing TX parameters calculations"""
    @parameterized.expand([
        (3550000000, 3560000000, 19, '50', 55290),
        (3555000000, 3570000000, 17, '75', 55365),
        (3600000000, 3605000000, 19, '25', 55765),
    ])
    def test_max_eirp_with_eirp_within_range(
            self,
            low_frequency_hz,
            high_frequency_hz,
            max_eirp_dbm_mhz,
            expected_bw_rbs,
            expected_earfcn,
    ) -> None:
        """Test that max erip parameters of the enodeb are set correctly when Domain Proxy returns CBSD sate

        Args:
            low_frequency_hz (Any): low frequency in hz
            high_frequency_hz (Any): high frequency in hz
            max_eirp_dbm_mhz (Any): max eirp
            expected_bw_rbs (Any): expected bandwidth
            expected_earfcn (Any): expected earfcn
        """
        desired_config = EnodebConfiguration(FreedomFiOneTrDataModel())
        channel = LteChannel(
            low_frequency_hz=low_frequency_hz,
            high_frequency_hz=high_frequency_hz,
            max_eirp_dbm_mhz=max_eirp_dbm_mhz,
        )
        state = CBSDStateResult(
            radio_enabled=True,
            channel=channel,
            channels=[channel],
            carrier_aggregation_enabled=False,
        )

        ff_one_update_desired_config_from_cbsd_state(state, desired_config)
        self._assert_config_updated(
            config=desired_config,
            bandwidth=expected_bw_rbs,
            earfcn=expected_earfcn,
            max_eirp=max_eirp_dbm_mhz,
            radio_enabled=True,
        )

    @parameterized.expand([
        (-138,),
        (38,),
    ])
    def test_max_eirp_with_eirp_out_of_range(self, max_eirp_dbm_mhz) -> None:
        """Test setting max eirp to out of range values raises an exception

        Args:
            max_eirp_dbm_mhz (Any): max eirp
        """
        desired_config = EnodebConfiguration(FreedomFiOneTrDataModel())
        channel = LteChannel(
            low_frequency_hz=3550000000,
            high_frequency_hz=3570000000,
            max_eirp_dbm_mhz=max_eirp_dbm_mhz,
        )
        state = CBSDStateResult(
            radio_enabled=True,
            channel=channel,
            channels=[channel],
            carrier_aggregation_enabled=False,
        )
        with self.assertRaises(ConfigurationError):
            ff_one_update_desired_config_from_cbsd_state(state, desired_config)

    @parameterized.expand([
        (3550000000, 3551000000),
        (3550000000, 3552000000),
    ])
    def test_tx_params_with_unsupported_bandwidths(self, low_frequency_hz, high_frequency_hz) -> None:
        """Test that tx parameters calculations raise an exception for unsupported bandwidth ranges"""
        desired_config = EnodebConfiguration(FreedomFiOneTrDataModel())
        channel = LteChannel(
            low_frequency_hz=low_frequency_hz,
            high_frequency_hz=high_frequency_hz,
            max_eirp_dbm_mhz=5,
        )
        state = CBSDStateResult(
            radio_enabled=True,
            channel=channel,
            channels=[channel],
            carrier_aggregation_enabled=False,
        )
        with self.assertRaises(ConfigurationError):
            ff_one_update_desired_config_from_cbsd_state(state, desired_config)

    def test_tx_params_not_set_when_radio_disabled(self):
        """Test that tx parameters of the enodeb are not set when ADMIN_STATE is disabled on the radio"""
        desired_config = EnodebConfiguration(FreedomFiOneTrDataModel())
        channel = LteChannel(
            low_frequency_hz=3550000000,
            high_frequency_hz=3570000000,
            max_eirp_dbm_mhz=20,
        )
        state = CBSDStateResult(
            radio_enabled=False,
            channel=channel,
            channels=[channel],
            carrier_aggregation_enabled=False,
        )

        ff_one_update_desired_config_from_cbsd_state(state, desired_config)
        self.assertEqual(1, len(desired_config.get_parameter_names()))
        self.assertFalse(
            desired_config.get_parameter(
                ParameterName.ADMIN_STATE,
            ),
        )

    def _assert_config_updated(
            self, config: EnodebConfiguration, bandwidth: str, earfcn: int, max_eirp: int, radio_enabled: bool,
    ) -> None:
        expected_values = {
            ParameterName.ADMIN_STATE: radio_enabled,
            ParameterName.DL_BANDWIDTH: bandwidth,
            ParameterName.UL_BANDWIDTH: bandwidth,
            ParameterName.EARFCNDL: earfcn,
            ParameterName.EARFCNUL: earfcn,
            SASParameters.SAS_MAX_EIRP: max_eirp,
            ParameterName.BAND: BAND,
        }
        for key, value in expected_values.items():
            self.assertEqual(config.get_parameter(key), value)


class FreedomFiOneCarrierAggregationTests(EnodebHandlerTestCase):
    channel_20_1 = LteChannel(low_frequency_hz=3550_000_000, high_frequency_hz=3570_000_000, max_eirp_dbm_mhz=20)
    channel_20_2 = LteChannel(low_frequency_hz=3570_000_000, high_frequency_hz=3590_000_000, max_eirp_dbm_mhz=20)
    channel_15_1 = LteChannel(low_frequency_hz=3570_000_000, high_frequency_hz=3585_000_000, max_eirp_dbm_mhz=20)
    channel_15_2 = LteChannel(low_frequency_hz=3585_000_000, high_frequency_hz=3600_000_000, max_eirp_dbm_mhz=20)
    channel_10_1 = LteChannel(low_frequency_hz=3570_000_000, high_frequency_hz=3580_000_000, max_eirp_dbm_mhz=20)
    channel_10_2 = LteChannel(low_frequency_hz=3580_000_000, high_frequency_hz=3590_000_000, max_eirp_dbm_mhz=20)
    channel_5_1 = LteChannel(low_frequency_hz=3550_000_000, high_frequency_hz=3555_000_000, max_eirp_dbm_mhz=20)
    channel_5_2 = LteChannel(low_frequency_hz=3555_000_000, high_frequency_hz=3560_000_000, max_eirp_dbm_mhz=20)

    # Domain Proxy state results, which should be supported and accepted by FreedomFi One. Carrier Aggregation should be set on eNB.
    state_ca_enabled_2_channels_20 = CBSDStateResult(channels=[channel_20_1, channel_20_2], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_enabled_2_channels_10 = CBSDStateResult(channels=[channel_10_1, channel_10_2], radio_enabled=True, carrier_aggregation_enabled=True)

    # Domain Proxy state results, which result in disabling the radio.
    state_ca_enabled_0_channels = CBSDStateResult(channels=[], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_disabled_0_channels = CBSDStateResult(channels=[], radio_enabled=True, carrier_aggregation_enabled=False)

    # Domain Proxy state results, with just 1 channel or 2 channels but CA disabled. Should result in eNB switching to Single Carrier
    state_ca_enabled_1_channel_20 = CBSDStateResult(channels=[channel_20_1], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_disabled_2_channels_20 = CBSDStateResult(channels=[channel_20_1, channel_20_2], radio_enabled=True, carrier_aggregation_enabled=False)
    state_ca_disabled_1_channel_20 = CBSDStateResult(channels=[channel_20_1], radio_enabled=True, carrier_aggregation_enabled=False)

    # Domain Proxy state results, with 2 channels that are not in supported bandwidth configurations. Should result in eNB switching to Single Carrier
    state_ca_enabled_2_channels_20_15 = CBSDStateResult(channels=[channel_20_1, channel_15_2], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_enabled_2_channels_20_10 = CBSDStateResult(channels=[channel_20_1, channel_10_2], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_enabled_2_channels_20_5 = CBSDStateResult(channels=[channel_20_1, channel_5_2], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_enabled_2_channels_15 = CBSDStateResult(channels=[channel_15_1, channel_15_2], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_enabled_2_channels_15_10 = CBSDStateResult(channels=[channel_15_1, channel_10_2], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_enabled_2_channels_15_5 = CBSDStateResult(channels=[channel_15_1, channel_5_2], radio_enabled=True, carrier_aggregation_enabled=True)
    state_ca_enabled_2_channels_10_5 = CBSDStateResult(channels=[channel_10_1, channel_5_2], radio_enabled=True, carrier_aggregation_enabled=True)

    @parameterized.expand([
        (state_ca_enabled_2_channels_20,),
        (state_ca_enabled_2_channels_10,),
    ])
    def test_ca_enabled_based_on_domain_proxy_state(self, state) -> None:
        """
        Test eNB configuration set to Carrier Aggregation when Domain Proxy
        state response contains applicable channels
        """
        config = EnodebConfiguration(FreedomFiOneTrDataModel())

        # Set required params in device configuration
        _pci = 260
        _cell_id = 138777000
        _tac = 1
        config.set_parameter(ParameterName.TAC, _tac)
        config.set_parameter(ParameterName.PCI, _pci)
        config.set_parameter(ParameterName.CELL_ID, 138777000)
        config.set_parameter(ParameterName.GPS_STATUS, True)

        # Update eNB configuration based on Domain Proxy state response
        ff_one_update_desired_config_from_cbsd_state(state, config)

        # Check eNB set to Carrier Aggregation
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_ENABLE), True)
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_CARRIER_NUMBER), 2)

        # Check second carrier set to parameters derived from second channel
        ca_channel_low_freq = state.channels[1].low_frequency_hz
        ca_channel_high_freq = state.channels[1].high_frequency_hz
        ca_earfcn = calc_earfcn(ca_channel_low_freq, ca_channel_high_freq)
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_EARFCNDL), ca_earfcn)
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_EARFCNDL), ca_earfcn)
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_BAND), 48)
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_TAC), _tac)
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_CELL_ID), _cell_id + 1)
        self.assertEqual(config.get_parameter(ParameterName.ADMIN_STATE), True)

    @parameterized.expand([
        (state_ca_enabled_0_channels,),
        (state_ca_disabled_0_channels,),
    ])
    def test_radio_disabled_based_on_domain_proxy_state(self, state) -> None:
        """
        Test eNB configuration set to disable radio tranmission and not touch
        any eNB parameters, when Domain Proxy state response contains no channels
        """
        config = EnodebConfiguration(FreedomFiOneTrDataModel())
        ff_one_update_desired_config_from_cbsd_state(state, config)
        self.assertEqual(config.get_parameter(ParameterName.ADMIN_STATE), False)

    @parameterized.expand([
        (state_ca_enabled_1_channel_20,),
        (state_ca_disabled_2_channels_20,),
        (state_ca_disabled_1_channel_20,),
        (state_ca_enabled_2_channels_20_15,),
        (state_ca_enabled_2_channels_20_10,),
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
        config = EnodebConfiguration(FreedomFiOneTrDataModel())

        # Set required params in device configuration
        config.set_parameter(ParameterName.GPS_STATUS, True)

        # Update eNB configuration based on Domain Proxy state response
        ff_one_update_desired_config_from_cbsd_state(state, config)

        # Check eNB set to Single Carrier
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_ENABLE), False)
        self.assertEqual(config.get_parameter(CarrierAggregationParameters.CA_CARRIER_NUMBER), 1)

        # Check eNB first carrier set to parameters derived from first channel
        sc_channel_low_freq = state.channels[0].low_frequency_hz
        sc_channel_high_freq = state.channels[0].high_frequency_hz
        sc_bw_mhz = calc_bandwidth_mhz(sc_channel_low_freq, sc_channel_high_freq)
        sc_bw_rbs = calc_bandwidth_rbs(sc_bw_mhz)
        sc_earfcn = calc_earfcn(sc_channel_low_freq, sc_channel_high_freq)
        self.assertEqual(config.get_parameter(ParameterName.DL_BANDWIDTH), sc_bw_rbs)
        self.assertEqual(config.get_parameter(ParameterName.UL_BANDWIDTH), sc_bw_rbs)
        self.assertEqual(config.get_parameter(ParameterName.EARFCNDL), sc_earfcn)
        self.assertEqual(config.get_parameter(ParameterName.EARFCNDL), sc_earfcn)
        self.assertEqual(config.get_parameter(ParameterName.ADMIN_STATE), True)


def _get_service_config(dp_mode: bool = True, prim_src: str = "GNSS"):
    return {
        "tr069": {
            "interface": "eth1",
            "port": 48080,
            "perf_mgmt_port": 8081,
            "public_ip": "192.88.99.142",
        },
        "reboot_enodeb_on_mme_disconnected": True,
        "s1_interface": "eth1",
        "sas": {
            "dp_mode": dp_mode,
            "sas_server_url":
                TEST_SAS_URL,
            "sas_uid": "M0LK4T",
            "sas_category": "A",
            "sas_channel_type": "GAA",
            "sas_cert_subject": TEST_SAS_CERT_SUBJECT,
            "sas_icg_group_id": "",
            "sas_location": "indoor",
            "sas_height_type": "AMSL",
        },
        "prim_src": prim_src,
        "firmware_upgrade_download": {
            "enbs": {},
            "firmwares": {},
            "models": {},
        },
    }
