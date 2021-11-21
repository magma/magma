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
from unittest.mock import Mock, call, patch

from lte.protos.mconfig import mconfigs_pb2
from magma.common.service import MagmaService
from magma.enodebd.data_models.data_model_parameters import ParameterName
from magma.enodebd.devices.device_map import get_device_handler_from_name
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.devices.freedomfi_one import (
    FreedomFiOneConfigurationInitializer,
    SASParameters,
    StatusParameters,
)
from magma.enodebd.tests.test_utils.config_builder import EnodebConfigBuilder
from magma.enodebd.tests.test_utils.enb_acs_builder import (
    EnodebAcsStateMachineBuilder,
)
from magma.enodebd.tests.test_utils.enodeb_handler import EnodebHandlerTestCase
from magma.enodebd.tests.test_utils.tr069_msg_builder import Tr069MessageBuilder
from magma.enodebd.tr069 import models

SRC_CONFIG_DIR = os.path.join(
    os.environ.get('MAGMA_ROOT'),
    'lte/gateway/configs',
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
                name='Device.Services.FAPService.1.FAPControl.LTE.X_000E8F_SAS.UserContactInformation',
                val_type='string',
                data='M0LK4T',
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

    def _get_service_config(self):
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
                "sas_enabled": True,
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
            "sentry": "disabled",
        }

    def build_freedomfi_one_acs_state_machine(self):
        service = EnodebAcsStateMachineBuilder.build_magma_service(
            mconfig=EnodebConfigBuilder.get_mconfig(),
            service_config=self._get_service_config(),
        )
        handler_class = get_device_handler_from_name(
            EnodebDeviceName.FREEDOMFI_ONE,
        )
        acs_state_machine = handler_class(service)
        return acs_state_machine

    def test_provision(self) -> None:
        """
        Test the basic provisioning workflow
        1 - enodeb sends Inform, enodebd sends InformResponse
        2 - enodeb sends empty HTTP message,
        3 - enodebd sends get transient params, updates the device state.
        4 - enodebd sends get param values, updates the device state
        5 - enodebd, sets fields including SAS fields.
        """
        logging.root.level = logging.DEBUG
        acs_state_machine = self.build_freedomfi_one_acs_state_machine()

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

    def test_manual_reboot_during_provisioning(self) -> None:
        """
        Test a scenario where a Magma user goes through the enodebd CLI to
        reboot the Baicells eNodeB.

        This checks the scenario where the command is sent in the middle
        of a TR-069 provisioning session.
        """
        logging.root.level = logging.DEBUG
        acs_state_machine = self.build_freedomfi_one_acs_state_machine()

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

    def test_post_processing(self) -> None:
        """ Test FreedomFi One specific post processing functionality"""

        acs_state_machine = self.build_freedomfi_one_acs_state_machine()
        cfg_desired = Mock()
        acs_state_machine.device_cfg.set_parameter(
            ParameterName.SERIAL_NUMBER,
            "2006CW5000023",
        )

        cfg_init = FreedomFiOneConfigurationInitializer(acs_state_machine)
        cfg_init.postprocess(
            EnodebConfigBuilder.get_mconfig(),
            self._get_service_config(), cfg_desired,
        )
        expected = [
            call.delete_parameter('EARFCNDL'),
            call.delete_parameter('DL bandwidth'),
            call.delete_parameter('UL bandwidth'),
            call.set_parameter(
                'tunnel_ref',
                'Device.IP.Interface.1.IPv4Address.1.',
            ),
            call.set_parameter('prim_src', 'GNSS'),
            call.set_parameter('carrier_agg_enable', True),
            call.set_parameter('carrier_number', 2),
            call.set_parameter('contiguous_cc', 0),
            call.set_parameter('web_ui_enable', False),
            call.set_parameter('sas_enabled', True),
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
            call.set_parameter_for_object(
                param_name='PLMN 1 cell reserved',
                value=True, object_name='PLMN 1',
            ),
        ]
        cfg_desired.mock_calls.sort()
        expected.sort()
        self.assertEqual(cfg_desired.mock_calls, expected)

        # Check without sas config
        service_cfg = {
            "tr069": {
                "interface": "eth1",
                "port": 48080,
                "perf_mgmt_port": 8081,
                "public_ip": "192.88.99.142",
            },
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
                'Device.IP.Interface.1.IPv4Address.1.',
            ),
            call.set_parameter('prim_src', 'GNSS'),
            call.set_parameter('carrier_agg_enable', True),
            call.set_parameter('carrier_number', 2),
            call.set_parameter('contiguous_cc', 0),
            call.set_parameter('web_ui_enable', False),
            call.set_parameter_for_object(
                param_name='PLMN 1 cell reserved',
                value=True, object_name='PLMN 1',
            ),
        ]
        cfg_desired.mock_calls.sort()
        expected.sort()
        self.assertEqual(cfg_desired.mock_calls, expected)

        service_cfg['web_ui_enable_list'] = ["2006CW5000023"]

        expected = [
            call.delete_parameter('EARFCNDL'),
            call.delete_parameter('DL bandwidth'),
            call.delete_parameter('UL bandwidth'),
            call.set_parameter(
                'tunnel_ref',
                'Device.IP.Interface.1.IPv4Address.1.',
            ),
            call.set_parameter('prim_src', 'GNSS'),
            call.set_parameter('carrier_agg_enable', True),
            call.set_parameter('carrier_number', 2),
            call.set_parameter('contiguous_cc', 0),
            call.set_parameter('web_ui_enable', False),
            call.set_parameter('web_ui_enable', True),
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
        cfg_desired.mock_calls.sort()
        expected.sort()
        self.assertEqual(cfg_desired.mock_calls, expected)

    @patch('magma.configuration.service_configs.CONFIG_DIR', SRC_CONFIG_DIR)
    def test_service_cfg_parsing(self):
        """ Test the parsing of the service config file for enodebd.yml"""
        service = MagmaService('enodebd', mconfigs_pb2.EnodebD())
        service_cfg = service.config
        service_cfg_1 = self._get_service_config()
        service_cfg_1['web_ui_enable_list'] = []
        service_cfg_1[FreedomFiOneConfigurationInitializer.SAS_KEY][
            SASParameters.SAS_UID
        ] = "INVALID_ID"
        service_cfg_1[FreedomFiOneConfigurationInitializer.SAS_KEY][
            SASParameters.SAS_CERT_SUBJECT
        ] = "INVALID_CERT_SUBJECT"
        self.assertEqual(service_cfg, service_cfg_1)

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
