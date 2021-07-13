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

from typing import Any, List, Optional

from magma.enodebd.tr069 import models


class Tr069MessageBuilder:
    @classmethod
    def get_parameter_value_struct(
            cls,
            name: str,
            val_type: str,
            data: Any,
    ) -> models.ParameterValueStruct:
        param_value = models.ParameterValueStruct()
        param_value.Name = name
        value = models.anySimpleType()
        value.type = val_type
        value.Data = data
        param_value.Value = value
        return param_value

    @classmethod
    def get_fault(cls) -> models.Fault:
        msg = models.Fault()
        msg.FaultCode = 0
        msg.FaultString = 'Some sort of fault'
        return msg

    @classmethod
    def get_reboot_inform(cls) -> models.Inform:
        msg = cls.get_inform()
        events = []

        event_boot = models.EventStruct()
        event_boot.EventCode = '1 BOOT'
        events.append(event_boot)

        event_reboot = models.EventStruct()
        event_reboot.EventCode = 'M Reboot'
        events.append(event_reboot)

        msg.Event.EventStruct = events
        return msg

    @classmethod
    def get_qafb_inform(
        cls,
        oui: str = '48BF74',
        sw_version: str = 'BaiBS_QAFB_1.6.4',
        enb_serial: str = '1202000181186TB0006',
        event_codes: Optional[List[str]] = None,
    ) -> models.Inform:
        if event_codes is None:
            event_codes = []
        msg = models.Inform()

        # DeviceId
        device_id = models.DeviceIdStruct()
        device_id.Manufacturer = 'Unused'
        device_id.OUI = oui
        device_id.ProductClass = 'Unused'
        device_id.SerialNumber = enb_serial
        msg.DeviceId = device_id

        # Event
        msg.Event = models.EventList()
        event_list = []
        for code in event_codes:
            event = models.EventStruct()
            event.EventCode = code
            event.CommandKey = ''
            event_list.append(event)
        msg.Event.EventStruct = event_list

        # ParameterList
        val_list = []
        val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.DeviceInfo.HardwareVersion',
                val_type='string',
                data='VER.C',
            ),
        )
        val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.DeviceInfo.ManufacturerOUI',
                val_type='string',
                data=oui,
            ),
        )
        val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.DeviceInfo.SoftwareVersion',
                val_type='string',
                data=sw_version,
            ),
        )
        val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.DeviceInfo.SerialNumber',
                val_type='string',
                data=enb_serial,
            ),
        )
        val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.ManagementServer.ConnectionRequestURL',
                val_type='string',
                data='http://192.168.60.248:7547/25dbc91d31276f0cb03391160531ecae',
            ),
        )
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = val_list

        return msg

        pass

    @classmethod
    def get_inform(
        cls,
        oui: str = '48BF74',
        sw_version: str = 'BaiBS_RTS_3.1.6',
        enb_serial: str = '120200002618AGP0003',
        event_codes: Optional[List[str]] = None,
    ) -> models.Inform:
        if event_codes is None:
            event_codes = []
        msg = models.Inform()

        # DeviceId
        device_id = models.DeviceIdStruct()
        device_id.Manufacturer = 'Unused'
        device_id.OUI = oui
        device_id.ProductClass = 'Unused'
        device_id.SerialNumber = enb_serial
        msg.DeviceId = device_id

        # Event
        msg.Event = models.EventList()
        event_list = []
        for code in event_codes:
            event = models.EventStruct()
            event.EventCode = code
            event.CommandKey = ''
            event_list.append(event)
        msg.Event.EventStruct = event_list

        # ParameterList
        val_list = []
        val_list.append(
            cls.get_parameter_value_struct(
                name='Device.DeviceInfo.HardwareVersion',
                val_type='string',
                data='VER.C',
            ),
        )
        val_list.append(
            cls.get_parameter_value_struct(
                name='Device.DeviceInfo.ManufacturerOUI',
                val_type='string',
                data=oui,
            ),
        )
        val_list.append(
            cls.get_parameter_value_struct(
                name='Device.DeviceInfo.SoftwareVersion',
                val_type='string',
                data=sw_version,
            ),
        )
        val_list.append(
            cls.get_parameter_value_struct(
                name='Device.DeviceInfo.SerialNumber',
                val_type='string',
                data=enb_serial,
            ),
        )
        val_list.append(
            cls.get_parameter_value_struct(
                name='Device.ManagementServer.ConnectionRequestURL',
                val_type='string',
                data='http://192.168.60.248:7547/25dbc91d31276f0cb03391160531ecae',
            ),
        )
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = val_list

        return msg

    @classmethod
    def get_qafb_read_only_param_values_response(
        cls,
    ) -> models.GetParameterValuesResponse:
        msg = models.GetParameterValuesResponse()
        param_val_list = []
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.CellConfig.1.LTE.X_QUALCOMM_FAPControl.OpState',
                val_type='boolean',
                data='false',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.CellConfig.1.LTE.X_QUALCOMM_FAPControl.OpState',
                val_type='boolean',
                data='false',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.CellConfig.1.LTE.X_QUALCOMM_FAPControl.OpState',
                val_type='boolean',
                data='false',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.FAP.GPS.latitude',
                val_type='int',
                data='0',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.FAP.GPS.longitude',
                val_type='int',
                data='0',
            ),
        )
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = param_val_list
        return msg

    @classmethod
    def get_read_only_param_values_response(
        cls,
    ) -> models.GetParameterValuesResponse:
        msg = models.GetParameterValuesResponse()
        param_val_list = []
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.OpState',
                val_type='boolean',
                data='false',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.RFTxStatus',
                val_type='boolean',
                data='false',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.DeviceInfo.X_BAICELLS_COM_GPS_Status',
                val_type='boolean',
                data='0',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.DeviceInfo.X_BAICELLS_COM_1588_Status',
                val_type='boolean',
                data='0',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.DeviceInfo.X_BAICELLS_COM_MME_Status',
                val_type='boolean',
                data='false',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.FAP.GPS.LockedLatitude',
                val_type='int',
                data='0',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.FAP.GPS.LockedLongitude',
                val_type='int',
                data='0',
            ),
        )
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = param_val_list
        return msg

    @classmethod
    def get_cavium_param_values_response(
        cls,
        admin_state: bool = False,
        earfcndl: int = 2405,
        num_plmns: int = 0,
    ) -> models.GetParameterValuesResponse:
        msg = models.GetParameterValuesResponse()
        param_val_list = []
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.DLBandwidth',
                val_type='string',
                data='20',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.FreqBandIndicator',
                val_type='string',
                data='5',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.ManagementServer.PeriodicInformInterval',
                val_type='int',
                data='5',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.CellRestriction.CellReservedForOperatorUse',
                val_type='boolean',
                data='false',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.ULBandwidth',
                val_type='string',
                data='n100',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.Common.CellIdentity',
                val_type='int',
                data='138777000',
            ),
        )
        # MME IP
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.Gateway.S1SigLinkServerList',
                val_type='string',
                data='"192.168.60.142"',
            ),
        )
        # perf mgmt enable
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.PerfMgmt.Config.1.Enable',
                val_type='boolean',
                data='true',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.CellRestriction.CellBarred',
                val_type='boolean',
                data='false',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.PerfMgmt.Config.1.PeriodicUploadInterval',
                val_type='int',
                data='600',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.AdminState',
                val_type='boolean',
                data=admin_state,
            ),
        )
        # Perf mgmt upload url
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.PerfMgmt.Config.1.URL',
                val_type='string',
                data='http://192.168.60.142:8081/',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.EPC.TAC',
                val_type='int',
                data='1',
            ),
        )
        # PCI
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.PhyCellID',
                val_type='int',
                data='260',
            ),
        )
        # MME port
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.Gateway.S1SigLinkPort',
                val_type='int',
                data='36412',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.IPsec.Enable',
                val_type='boolean',
                data='false',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.EARFCNDL',
                val_type='int',
                data='2405',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.EARFCNUL',
                val_type='int',
                data='20405',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.Capabilities.LTE.DuplexMode',
                val_type='string',
                data='FDDMode',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.Capabilities.LTE.BandsSupported',
                val_type='string',
                data='5',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.ManagementServer.PeriodicInformEnable',
                val_type='int',
                data='5',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNListNumberOfEntries',
                val_type='int',
                data=str(num_plmns),
            ),
        )
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = param_val_list
        return msg

    @classmethod
    def get_regular_param_values_response(
        cls,
        admin_state: bool = False,
        earfcndl: int = 39250,
        exclude_num_plmns: bool = False,
    ) -> models.GetParameterValuesResponse:
        msg = models.GetParameterValuesResponse()
        param_val_list = []
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.DLBandwidth',
                val_type='string',
                data='n100',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.FreqBandIndicator',
                val_type='string',
                data='40',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.ManagementServer.PeriodicInformInterval',
                val_type='int',
                data='5',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.CellRestriction.CellReservedForOperatorUse',
                val_type='boolean',
                data='false',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.ULBandwidth',
                val_type='string',
                data='20',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.X_BAICELLS_COM_LTE.EARFCNDLInUse',
                val_type='string',
                data=earfcndl,
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.PHY.TDDFrame.SpecialSubframePatterns',
                val_type='int',
                data='7',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.Common.CellIdentity',
                val_type='int',
                data='138777000',
            ),
        )
        # MME IP
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.Gateway.S1SigLinkServerList',
                val_type='string',
                data='"192.168.60.142"',
            ),
        )
        if not exclude_num_plmns:
            param_val_list.append(
                cls.get_parameter_value_struct(
                    name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNListNumberOfEntries',
                    val_type='int',
                    data='1',
                ),
            )
        # perf mgmt enable
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.FAP.PerfMgmt.Config.1.Enable',
                val_type='boolean',
                data='true',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.CellRestriction.CellBarred',
                val_type='boolean',
                data='false',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.FAP.PerfMgmt.Config.1.PeriodicUploadInterval',
                val_type='int',
                data='300',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.AdminState',
                val_type='boolean',
                data=admin_state,
            ),
        )
        # Local gateway enable
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.DeviceInfo.X_BAICELLS_COM_LTE_LGW_Switch',
                val_type='boolean',
                data='0',
            ),
        )
        # Perf mgmt upload url
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.FAP.PerfMgmt.Config.1.URL',
                val_type='string',
                data='http://192.168.60.142:8081/',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.EPC.TAC',
                val_type='int',
                data='1',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.Gateway.X_BAICELLS_COM_MmePool.Enable',
                val_type='boolean',
                data='false',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.PHY.TDDFrame.SubFrameAssignment',
                val_type='int',
                data='2',
            ),
        )
        # PCI
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.PhyCellID',
                val_type='int',
                data='260',
            ),
        )
        # MME port
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.FAPControl.LTE.Gateway.S1SigLinkPort',
                val_type='int',
                data='36412',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.Ipsec.IPSEC_ENABLE',
                val_type='boolean',
                data='false',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.X_BAICELLS_COM_LTE.EARFCNULInUse',
                val_type='int',
                data='39150',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.Capabilities.LTE.DuplexMode',
                val_type='string',
                data='TDDMode',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.Capabilities.LTE.BandsSupported',
                val_type='string',
                data='40',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.ManagementServer.PeriodicInformEnable',
                val_type='int',
                data='5',
            ),
        )
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = param_val_list
        return msg

    @classmethod
    def get_qafb_regular_param_values_response(
        cls,
        admin_state: bool = False,
        earfcndl: int = 39250,
    ) -> models.GetParameterValuesResponse:
        msg = models.GetParameterValuesResponse()
        param_val_list = []
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.CellConfig.LTE.RAN.RF.DLBandwidth',
                val_type='string',
                data='20',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.CellConfig.LTE.RAN.RF.FreqBandIndicator',
                val_type='string',
                data='40',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.ManagementServer.PeriodicInformInterval',
                val_type='int',
                data='5',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.CellConfig.LTE.RAN.CellRestriction.CellReservedForOperatorUse',
                val_type='boolean',
                data='false',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.CellConfig.LTE.RAN.RF.ULBandwidth',
                val_type='string',
                data='20',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.CellConfig.LTE.RAN.RF.ULBandwidth',
                val_type='int',
                data='1',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.X_BAICELLS_COM_LTE.EARFCNDLInUse',
                val_type='string',
                data=earfcndl,
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.CellConfig.LTE.RAN.PHY.TDDFrame.SpecialSubframePatterns',
                val_type='int',
                data='7',
            ),
        )
        # MME IP
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.FAPControl.LTE.Gateway.S1SigLinkServerList',
                val_type='string',
                data='"192.168.60.142"',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.CellConfig.LTE.EPC.PLMNListNumberOfEntries',
                val_type='int',
                data='1',
            ),
        )
        # perf mgmt enable
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.FAP.PerfMgmt.Config.1.Enable',
                val_type='boolean',
                data='true',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.CellConfig.LTE.RAN.CellRestriction.CellBarred',
                val_type='boolean',
                data='false',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.FAP.PerfMgmt.Config.1.PeriodicUploadInterval',
                val_type='int',
                data='300',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.FAPControl.LTE.AdminState',
                val_type='boolean',
                data='false',
            ),
        )
        # Local gateway enable
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.DeviceInfo.X_BAICELLS_COM_LTE_LGW_Switch',
                val_type='boolean',
                data='0',
            ),
        )
        # Perf mgmt upload url
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.FAP.PerfMgmt.Config.1.URL',
                val_type='string',
                data='http://192.168.60.142:8081/',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.CellConfig.LTE.EPC.TAC',
                val_type='int',
                data='1',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.FAPControl.LTE.Gateway.X_BAICELLS_COM_MmePool.Enable',
                val_type='boolean',
                data='false',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.CellConfig.LTE.RAN.PHY.TDDFrame.SubFrameAssignment',
                val_type='int',
                data='2',
            ),
        )
        # PCI
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.CellConfig.LTE.RAN.RF.PhyCellID',
                val_type='int',
                data='260',
            ),
        )
        # MME port
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.FAPControl.LTE.Gateway.S1SigLinkPort',
                val_type='int',
                data='36412',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='boardconf.ipsec.ipsecConfig.onBoot',
                val_type='boolean',
                data='false',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.X_BAICELLS_COM_LTE.EARFCNULInUse',
                val_type='int',
                data='9212',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='boardconf.status.eepromInfo.div_multiple',
                val_type='string',
                data='02',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='boardconf.status.eepromInfo.work_mode',
                val_type='string',
                data='1C000400',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.ManagementServer.PeriodicInformEnable',
                val_type='int',
                data='5',
            ),
        )
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = param_val_list
        return msg

    @classmethod
    def get_cavium_object_param_values_response(
            cls,
            num_plmns: int,
    ) -> models.GetParameterValuesResponse:
        msg = models.GetParameterValuesResponse()
        param_val_list = []
        for i in range(1, num_plmns + 1):
            param_val_list.append(
                cls.get_parameter_value_struct(
                    name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.%d.IsPrimary' % i,
                    val_type='boolean',
                    data='true',
                ),
            )
            param_val_list.append(
                cls.get_parameter_value_struct(
                    name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.%d.CellReservedForOperatorUse' % i,
                    val_type='boolean',
                    data='false',
                ),
            )
            param_val_list.append(
                cls.get_parameter_value_struct(
                    name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.%d.PLMNID' % i,
                    val_type='string',
                    data='00101',
                ),
            )
            param_val_list.append(
                cls.get_parameter_value_struct(
                    name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.%d.Enable' % i,
                    val_type='boolean',
                    data='true',
                ),
            )
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = param_val_list
        return msg

    @classmethod
    def get_object_param_values_response(
            cls,
    ) -> models.GetParameterValuesResponse:
        msg = models.GetParameterValuesResponse()
        param_val_list = []
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.IsPrimary',
                val_type='boolean',
                data='true',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.CellReservedForOperatorUse',
                val_type='boolean',
                data='false',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.PLMNID',
                val_type='string',
                data='00101',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.Enable',
                val_type='boolean',
                data='true',
            ),
        )
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = param_val_list
        return msg

    @classmethod
    def get_qafb_object_param_values_response(
            cls,
    ) -> models.GetParameterValuesResponse:
        msg = models.GetParameterValuesResponse()
        param_val_list = []
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.IsPrimary',
                val_type='boolean',
                data='true',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.CellReservedForOperatorUse',
                val_type='boolean',
                data='false',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.PLMNID',
                val_type='string',
                data='00101',
            ),
        )
        param_val_list.append(
            cls.get_parameter_value_struct(
                name='InternetGatewayDevice.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.Enable',
                val_type='boolean',
                data='true',
            ),
        )
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = param_val_list
        return msg

    @classmethod
    def get_reboot_response(cls) -> models.RebootResponse:
        return models.RebootResponse()
