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
from unittest import TestCase, mock
from unittest.mock import Mock, call, patch

from dp.protos.enodebd_dp_pb2 import CBSDStateResult, LteChannel
from lte.protos.mconfig import mconfigs_pb2
from magma.common.service import MagmaService
from magma.enodebd.data_models.data_model_parameters import ParameterName
from magma.enodebd.device_config.cbrs_consts import BAND
from magma.enodebd.device_config.configuration_init import build_desired_config
from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.devices.freedomfi_one import (
    SAS_KEY,
    FreedomFiOneConfigurationInitializer,
    FreedomFiOneEndSessionState,
    FreedomFiOneGetInitState,
    FreedomFiOneNotifyDPState,
    FreedomFiOneTrDataModel,
    SASParameters,
    StatusParameters,
    ff_one_update_desired_config_from_cbsd_state,
)
from magma.enodebd.exceptions import ConfigurationError
from magma.enodebd.tests.test_utils.config_builder import EnodebConfigBuilder
from magma.enodebd.tests.test_utils.enb_acs_builder import (
    EnodebAcsStateMachineBuilder,
)
from magma.enodebd.tests.test_utils.enodeb_handler import EnodebHandlerTestCase
from magma.enodebd.tests.test_utils.tr069_msg_builder import Tr069MessageBuilder
from magma.enodebd.tr069 import models
from parameterized import parameterized

SRC_CONFIG_DIR = os.path.join(
    os.environ.get('MAGMA_ROOT'),
    'lte/gateway/configs',
)

MOCK_CBSD_STATE = CBSDStateResult(
    radio_enabled=True,
    channel=LteChannel(
        low_frequency_hz=3550_000_000,
        high_frequency_hz=3570_000_000,
        max_eirp_dbm_mhz=15,
    ),
)


class FreedomFiOneTests(EnodebHandlerTestCase):

    def _get_freedomfi_one_read_only_param_values_response(
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
                name='Device.Services.FAPService.1.FAPControl.LTE.X_000E8F_RRMConfig.X_000E8F_Cell_Number',
                val_type='int',
                data='2',
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
                data='/C=TW/O=Sercomm/OU=WInnForum CBSD Certificate/CN=P27-SCE4255W:%s',
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
                data='https://spectrum-connect.federatedwireless.com/v1.2/',
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
        return msg

    @mock.patch('magma.enodebd.devices.freedomfi_one.get_cbsd_state')
    def test_provision(self, mock_get_state) -> None:
        """
        Test the basic provisioning workflow
        1 - enodeb sends Inform, enodebd sends InformResponse
        2 - enodeb sends empty HTTP message,
        3 - enodebd sends get transient params, updates the device state.
        4 - enodebd sends get param values, updates the device state
        5 - enodebd, sets fields including SAS fields.
        """

        mock_get_state.return_value = MOCK_CBSD_STATE

        logging.root.level = logging.DEBUG
        acs_state_machine = EnodebAcsStateMachineBuilder.build_acs_state_machine(EnodebDeviceName.FREEDOMFI_ONE)
        acs_state_machine._service.config = get_service_config()
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
            enb_serial="2006CW5000023",
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
        for tr_69_nodes in StatusParameters.STATUS_PARAMETERS.values():
            self.assertIn(tr_69_nodes.path, resp.ParameterNames.string)

        req = self._get_freedomfi_one_read_only_param_values_response()
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


class FreedomFiOneStatesTests(EnodebHandlerTestCase):
    """Testing FreedomfiOne specific states"""

    @parameterized.expand([
        (False, FreedomFiOneNotifyDPState),
        (True, FreedomFiOneEndSessionState),
    ])
    @mock.patch('magma.enodebd.devices.freedomfi_one.get_cbsd_state')
    def test_end_session_and_notify_dp_transition_depending_on_sas_enabled_flag(
            self, sas_enabled, expected_state, mock_get_state,
    ):
        """Testing if SM steps in and out of FreedomFiOneWaitNotifyDPState as per state map depending on whether
        sas_enabled param is set to True or False in the service config"""

        mock_get_state.return_value = MOCK_CBSD_STATE

        acs_state_machine = EnodebAcsStateMachineBuilder.build_acs_state_machine(EnodebDeviceName.FREEDOMFI_ONE)
        acs_state_machine._service.config = get_service_config(sas_enabled=sas_enabled)
        acs_state_machine.desired_cfg = build_desired_config(
            acs_state_machine.mconfig,
            acs_state_machine.service_config,
            acs_state_machine.device_cfg,
            acs_state_machine.data_model,
            acs_state_machine.config_postprocessor,
        )

        # Need to fill these values in the device_cfg if we're going to transition to notify_dp state
        acs_state_machine.device_cfg.set_parameter(SASParameters.SAS_USER_ID, 'test_user')
        acs_state_machine.device_cfg.set_parameter(SASParameters.SAS_FCC_ID, 'test_fcc')
        acs_state_machine.device_cfg.set_parameter(ParameterName.SERIAL_NUMBER, 'test_sn')
        acs_state_machine.transition('check_wait_get_params')

        msg = Tr069MessageBuilder.param_values_qrtb_response([], models.GetParameterValuesResponse)

        # SM should transition from check_wait_get_params to end_session -> notify_dp automatically
        # upon receiving response from the radio
        acs_state_machine.handle_tr069_message(msg)

        self.assertIsInstance(acs_state_machine.state, expected_state)

        msg = Tr069MessageBuilder.get_inform(event_codes=['1 BOOT'])

        # SM should go into wait_inform state, respond with Inform response and transition to FreedomFiOneGetInitState
        acs_state_machine.handle_tr069_message(msg)

        self.assertIsInstance(acs_state_machine.state, FreedomFiOneGetInitState)

    def test_manual_reboot_during_provisioning(self) -> None:
        """
        Test a scenario where a Magma user goes through the enodebd CLI to
        reboot the Baicells eNodeB.

        This checks the scenario where the command is sent in the middle
        of a TR-069 provisioning session.
        """
        logging.root.level = logging.DEBUG
        acs_state_machine = EnodebAcsStateMachineBuilder.build_acs_state_machine(EnodebDeviceName.FREEDOMFI_ONE)

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

    @parameterized.expand([
        (True, "GNSS"),
        (False, "GNSS"),
        (True, "some_other_value"),
        (False, "some_other_value"),
    ])
    def test_post_processing(self, sas_enabled, prim_src) -> None:
        """ Test FreedomFi One specific post processing functionality"""

        acs_state_machine = EnodebAcsStateMachineBuilder.build_acs_state_machine(EnodebDeviceName.FREEDOMFI_ONE)
        cfg_desired = Mock()
        acs_state_machine.device_cfg.set_parameter(
            ParameterName.SERIAL_NUMBER,
            "2006CW5000023",
        )

        cfg_init = FreedomFiOneConfigurationInitializer(acs_state_machine)
        cfg_init.postprocess(
            EnodebConfigBuilder.get_mconfig(),
            get_service_config(sas_enabled=sas_enabled, prim_src=prim_src),
            cfg_desired,
        )
        expected = [
            call.delete_parameter('EARFCNDL'),
            call.delete_parameter('DL bandwidth'),
            call.delete_parameter('UL bandwidth'),
            call.set_parameter(
                'tunnel_ref',
                value='Device.IP.Interface.1.IPv4Address.1.',
            ),
            call.set_parameter('prim_src', prim_src),
            call.set_parameter('carrier_agg_enable', value=True),
            call.set_parameter('carrier_number', value=2),
            call.set_parameter('contiguous_cc', value=0),
            call.set_parameter('web_ui_enable', value=False),
            call.set_parameter('sas_enabled', sas_enabled),
            call.set_parameter_for_object(
                param_name='PLMN 1 cell reserved',
                value=True, object_name='PLMN 1',
            ),
        ]
        if sas_enabled:
            expected += [
                call.set_parameter(
                    'sas_server_url',
                    'https://spectrum-connect.federatedwireless.com/v1.2/',
                ),
                call.set_parameter('sas_uid', 'M0LK4T'),
                call.set_parameter('sas_category', 'A'),
                call.set_parameter('sas_channel_type', 'GAA'),
                call.set_parameter(
                    'sas_cert_subject',
                    '/C=TW/O=Sercomm/OU=WInnForum CBSD Certificate/CN=P27-SCE4255W:%s',
                ),
                call.set_parameter('sas_location', 'indoor'),
                call.set_parameter('sas_height_type', 'AMSL'),
            ]

        cfg_desired.assert_has_calls(expected, any_order=True)

        # Check without sas config
        service_cfg = {
            "tr069": {
                "interface": "eth1",
                "port": 48080,
                "perf_mgmt_port": 8081,
                "public_ip": "192.88.99.142",
            },
            "prim_src": prim_src,
            "reboot_enodeb_on_mme_disconnected": True,
            "s1_interface": "eth1",
        }
        cfg_desired = Mock()
        cfg_init.postprocess(
            EnodebConfigBuilder.get_mconfig(),
            service_cfg, cfg_desired,
        )
        expected = [
            call.delete_parameter('EARFCNDL'),
            call.delete_parameter('DL bandwidth'),
            call.delete_parameter('UL bandwidth'),
            call.set_parameter(
                'tunnel_ref',
                value='Device.IP.Interface.1.IPv4Address.1.',
            ),
            call.set_parameter('sas_enabled', False),
            call.set_parameter('prim_src', prim_src),
            call.set_parameter('carrier_agg_enable', value=True),
            call.set_parameter('carrier_number', value=2),
            call.set_parameter('contiguous_cc', value=0),
            call.set_parameter('web_ui_enable', value=False),
            call.set_parameter_for_object(
                param_name='PLMN 1 cell reserved',
                value=True, object_name='PLMN 1',
            ),
        ]
        cfg_desired.assert_has_calls(expected, any_order=True)

        service_cfg['web_ui_enable_list'] = ["2006CW5000023"]

        expected = [
            call.delete_parameter('EARFCNDL'),
            call.delete_parameter('DL bandwidth'),
            call.delete_parameter('UL bandwidth'),
            call.set_parameter(
                'tunnel_ref',
                value='Device.IP.Interface.1.IPv4Address.1.',
            ),
            call.set_parameter('sas_enabled', False),
            call.set_parameter('prim_src', prim_src),
            call.set_parameter('carrier_agg_enable', value=True),
            call.set_parameter('carrier_number', value=2),
            call.set_parameter('contiguous_cc', value=0),
            call.set_parameter('web_ui_enable', value=False),
            call.set_parameter('web_ui_enable', value=True),
            call.set_parameter_for_object(
                param_name='PLMN 1 cell reserved',
                value=True, object_name='PLMN 1',
            ),
        ]
        cfg_desired = Mock()
        cfg_init.postprocess(
            EnodebConfigBuilder.get_mconfig(),
            service_cfg, cfg_desired,
        )
        cfg_desired.assert_has_calls(expected, any_order=True)

    @patch('magma.configuration.service_configs.CONFIG_DIR', SRC_CONFIG_DIR)
    def test_service_cfg_parsing(self):
        """ Test the parsing of the service config file for enodebd.yml"""
        self.maxDiff = None
        service = MagmaService('enodebd', mconfigs_pb2.EnodebD())
        service_cfg = service.config
        service_cfg_1 = get_service_config()
        service_cfg_1['web_ui_enable_list'] = []
        service_cfg_1["prim_src"] = "GNSS"
        service_cfg_1[SAS_KEY][SASParameters.SAS_UID] = "INVALID_ID"
        service_cfg_1[SAS_KEY][SASParameters.SAS_CERT_SUBJECT] = "INVALID_CERT_SUBJECT"
        self.assertDictEqual(service_cfg, service_cfg_1)

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


class TXParamsTests(TestCase):
    @parameterized.expand([
        (3550_000_000, 3560_000_000, 19, '50', 55290, 24),
        (3555_000_000, 3570_000_000, 17, '75', 55365, 23),
        (3600_000_000, 3605_000_000, 19, '25', 55765, 20),
    ])
    def test_tx_parameters_with_eirp_within_range(
            self,
            low_frequency_hz,
            high_frequency_hz,
            max_eirp_dbm_mhz,
            expected_bw_rbs,
            expected_earfcn,
            expected_tx_power,
    ) -> None:
        """Test that tx parameters of the enodeb are calculated correctly when eirp received from SAS
        is within acceptable range for the given bandwidth"""
        desired_config = EnodebConfiguration(FreedomFiOneTrDataModel())
        channel = LteChannel(
            low_frequency_hz=low_frequency_hz,
            high_frequency_hz=high_frequency_hz,
            max_eirp_dbm_mhz=max_eirp_dbm_mhz,
        )
        state = CBSDStateResult(
            radio_enabled=True,
            channel=channel,
        )

        ff_one_update_desired_config_from_cbsd_state(state, desired_config)
        self.assert_config_updated(
            config=desired_config,
            bandwidth=expected_bw_rbs,
            earfcn=expected_earfcn,
            tx_power=expected_tx_power,
            radio_enabled=True,
        )

    @parameterized.expand([
        (30,),
        (-10,),
    ])
    def test_tx_parameters_with_eirp_out_of_range(self, max_eirp_dbm_mhz) -> None:
        """Test that tx parameters calculations raise an exception when eirp received from SAS
        is outside of acceptable range for the given bandwidth"""
        desired_config = EnodebConfiguration(FreedomFiOneTrDataModel())
        channel = LteChannel(
            low_frequency_hz=3550_000_000,
            high_frequency_hz=3570_000_000,
            max_eirp_dbm_mhz=max_eirp_dbm_mhz,
        )
        state = CBSDStateResult(
            radio_enabled=True,
            channel=channel,
        )
        with self.assertRaises(ConfigurationError):
            ff_one_update_desired_config_from_cbsd_state(state, desired_config)

    @parameterized.expand([
        (3550_000_000, 3551_000_000),
        (3550_000_000, 3552_000_000),
    ])
    def test_tx_parameters_with_unsupported_bandwidths(self, low_frequency_hz, high_frequency_hz) -> None:
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
        )
        with self.assertRaises(ConfigurationError):
            ff_one_update_desired_config_from_cbsd_state(state, desired_config)

    def test_tx_parameters_not_set_when_radio_disabled(self):
        """Test that tx parameters of the enodeb are not set when ADMIN_STATE is disabled on the radio"""
        desired_config = EnodebConfiguration(FreedomFiOneTrDataModel())
        channel = LteChannel(
            low_frequency_hz=3550_000_000,
            high_frequency_hz=3570_000_000,
            max_eirp_dbm_mhz=20,
        )
        state = CBSDStateResult(
            radio_enabled=False,
            channel=channel,
        )

        ff_one_update_desired_config_from_cbsd_state(state, desired_config)
        self.assertEqual(1, len(desired_config.get_parameter_names()))
        self.assertEqual(False, desired_config.get_parameter(ParameterName.ADMIN_STATE))

    def assert_config_updated(
            self, config: EnodebConfiguration, bandwidth: str, earfcn: int, tx_power: int, radio_enabled: bool,
    ) -> None:
        expected_values = {
            ParameterName.ADMIN_STATE: radio_enabled,
            ParameterName.DL_BANDWIDTH: bandwidth,
            ParameterName.UL_BANDWIDTH: bandwidth,
            ParameterName.EARFCNDL: earfcn,
            ParameterName.EARFCNUL: earfcn,
            SASParameters.TX_POWER_CONFIG: tx_power,
            SASParameters.FREQ_BAND1: BAND,
            SASParameters.FREQ_BAND2: BAND,
        }
        for key, value in expected_values.items():
            self.assertEqual(config.get_parameter(key), value)


def get_service_config(sas_enabled: bool = True, prim_src: str = "GNSS"):
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
            "sas_enabled": sas_enabled,
            "sas_server_url":
                "https://spectrum-connect.federatedwireless.com/v1.2/",
            "sas_uid": "M0LK4T",
            "sas_category": "A",
            "sas_channel_type": "GAA",
            "sas_cert_subject": "/C=TW/O=Sercomm/OU=WInnForum CBSD "
                                "Certificate/CN=P27-SCE4255W:%s",
            "sas_icg_group_id": "",
            "sas_location": "indoor",
            "sas_height_type": "AMSL",
        },
        'sentry': 'disabled',
        'prim_src': prim_src,
    }
