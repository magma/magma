"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

# pylint: disable=protected-access
import asyncio
from typing import Any
from unittest import TestCase, mock
from lte.protos.mconfig import mconfigs_pb2
from magma.enodebd.devices.baicells import BaicellsHandler
from magma.enodebd.stats_manager import StatsManager
from magma.enodebd.tr069 import models


class BaicellsHandlerTests(TestCase):
    def test_provisioning(self) -> None:
        acs_state_machine = self._build_acs_state_machine()

        # Send an Inform message, wait for an InformResponse
        inform_msg = self._get_inform()
        resp = acs_state_machine.handle_tr069_message(inform_msg)
        self.assertTrue(isinstance(resp, models.InformResponse),
                        'Should respond with an InformResponse')

        # Send an empty http request to kick off the rest of provisioning
        req = models.DummyInput()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for an optional parameter, three times
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting param values')
        req = self._get_fault()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting param values')
        req = self._get_fault()
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting param values')
        req = self._get_fault()
        resp = acs_state_machine.handle_tr069_message(req)

        # Expect a request for read-only params
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting param values')
        req = self._get_read_only_param_values_response()

        # Send back some typical values
        # And then SM should request regular parameter values
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting param values')

        # Send back typical values for the regular parameters
        req = self._get_regular_param_values_response()
        resp = acs_state_machine.handle_tr069_message(req)

        # SM will be requesting object parameter values
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting object param vals')

        # Send back some typical values for object parameters
        req = self._get_object_param_values_response()
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
        req = self._get_read_only_param_values_response()

        # Send back some typical values
        # And then SM should continue polling the read-only params
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting param values')

        # Expect a request for read-only params
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting param values')
        req = self._get_read_only_param_values_response()

        # Send back some typical values
        resp = acs_state_machine.handle_tr069_message(req)
        self.assertTrue(isinstance(resp, models.GetParameterValues),
                        'State machine should be requesting param values')
        return

    def _get_mconfig(self) -> mconfigs_pb2.EnodebD:
        mconfig = mconfigs_pb2.EnodebD()
        mconfig.bandwidth_mhz = 20
        mconfig.special_subframe_pattern = 7
        mconfig.earfcndl = 44490
        mconfig.log_level = 1
        mconfig.plmnid_list = "00101"
        mconfig.pci = 260
        mconfig.allow_enodeb_transmit = False
        mconfig.subframe_assignment = 2
        mconfig.tac = 1
        # tdd config
        mconfig.tdd_config.earfcndl = 39150
        mconfig.tdd_config.subframe_assignment = 2
        mconfig.tdd_config.special_subframe_pattern = 7
        return mconfig

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
        }

    def _build_acs_state_machine(self) -> BaicellsHandler:
        # Build the state_machine
        stats_mgr = StatsManager()
        event_loop = asyncio.get_event_loop()
        mconfig = self._get_mconfig()
        service_config = self._get_service_config()
        with mock.patch('magma.common.service.MagmaService') as MockService:
            MockService.config = service_config
            MockService.mconfig = mconfig
            MockService.loop = event_loop
            acs_state_machine = BaicellsHandler(MockService, stats_mgr)
            return acs_state_machine

    def _get_parameter_value_struct(
        self,
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

    def _get_fault(self) -> models.Fault:
        msg = models.Fault()
        msg.FaultCode = 0
        msg.FaultString = 'Some sort of fault'
        return msg

    def _get_reboot_inform(self) -> models.Inform:
        msg = self._get_inform()
        events = []

        event_boot = models.EventStruct()
        event_boot.EventCode = '1 BOOT'
        events.append(event_boot)

        event_reboot = models.EventStruct()
        event_reboot.EventCode = 'M Reboot'
        events.append(event_reboot)

        msg.Event.EventStruct = events
        return msg

    def _get_inform(self) -> models.Inform:
        msg = models.Inform()

        # DeviceId
        device_id = models.DeviceIdStruct()
        device_id.Manufacturer = 'Unused'
        device_id.OUI = '48BF74'
        device_id.ProductClass = 'Unused'
        device_id.SerialNumber = '120200002618AGP0003'
        msg.DeviceId = device_id

        # Event
        msg.Event = models.EventList()
        msg.Event.EventStruct = []

        # ParameterList
        val_list = []
        val_list.append(self._get_parameter_value_struct(
            name='Device.DeviceInfo.HardwareVersion',
            val_type='string',
            data='VER.C',
        ))
        val_list.append(self._get_parameter_value_struct(
            name='Device.DeviceInfo.ManufacturerOUI',
            val_type='string',
            data='48BF74',
        ))
        val_list.append(self._get_parameter_value_struct(
            name='Device.DeviceInfo.SoftwareVersion',
            val_type='string',
            data='BaiBS_RTS_3.1.6',
        ))
        val_list.append(self._get_parameter_value_struct(
            name='Device.ManagementServer.ConnectionRequestURL',
            val_type='string',
            data='http://192.168.60.248:7547/25dbc91d31276f0cb03391160531ecae',
        ))
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = val_list

        return msg

    def _get_read_only_param_values_response(
        self,
    ) -> models.GetParameterValuesResponse:
        msg = models.GetParameterValuesResponse()
        param_val_list = []
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.FAPControl.LTE.OpState',
            val_type='boolean',
            data='false',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.FAPControl.LTE.RFTxStatus',
            val_type='boolean',
            data='false',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.DeviceInfo.X_BAICELLS_COM_GPS_Status',
            val_type='boolean',
            data='0',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.DeviceInfo.X_BAICELLS_COM_1588_Status',
            val_type='boolean',
            data='0',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.DeviceInfo.X_BAICELLS_COM_MME_Status',
            val_type='boolean',
            data='false',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.FAP.GPS.LockedLatitude',
            val_type='int',
            data='0',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.FAP.GPS.LockedLongitude',
            val_type='int',
            data='0',
        ))
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = param_val_list
        return msg

    def _get_regular_param_values_response(
        self,
    ) -> models.GetParameterValuesResponse:
        msg = models.GetParameterValuesResponse()
        param_val_list = []
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.DLBandwidth',
            val_type='string',
            data='n100',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.FreqBandIndicator',
            val_type='string',
            data='40',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.ManagementServer.PeriodicInformInterval',
            val_type='int',
            data='5',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.RAN.CellRestriction.CellReservedForOperatorUse',
            val_type='boolean',
            data='false',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.ULBandwidth',
            val_type='string',
            data='n100',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.X_BAICELLS_COM_LTE.EARFCNDLInUse',
            val_type='string',
            data='n100',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.RAN.PHY.TDDFrame.SpecialSubframePatterns',
            val_type='int',
            data='7',
        ))
        # MME IP
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.FAPControl.LTE.Gateway.S1SigLinkServerList',
            val_type='string',
            data='"192.168.60.142"',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNListNumberOfEntries',
            val_type='int',
            data='1'
        ))
        # perf mgmt enable
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.FAP.PerfMgmt.Config.1.Enable',
            val_type='boolean',
            data='true'
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.RAN.CellRestriction.CellBarred',
            val_type='boolean',
            data='false'
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.FAP.PerfMgmt.Config.1.PeriodicUploadInterval',
            val_type='int',
            data='300'
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.FAPControl.LTE.AdminState',
            val_type='boolean',
            data='false'
        ))
        # Local gateway enable
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.DeviceInfo.X_BAICELLS_COM_LTE_LGW_Switch',
            val_type='boolean',
            data='0'
        ))
        # Perf mgmt upload url
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.FAP.PerfMgmt.Config.1.URL',
            val_type='string',
            data='http://192.168.60.142:8081/'
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.EPC.TAC',
            val_type='int',
            data='1',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.FAPControl.LTE.Gateway.X_BAICELLS_COM_MmePool.Enable',
            val_type='boolean',
            data='false',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.RAN.PHY.TDDFrame.SubFrameAssignment',
            val_type='int',
            data='2',
        ))
        # PCI
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.PhyCellID',
            val_type='int',
            data='260',
        ))
        # MME port
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.FAPControl.LTE.Gateway.S1SigLinkPort',
            val_type='int',
            data='36412',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.Ipsec.IPSEC_ENABLE',
            val_type='boolean',
            data='false',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.X_BAICELLS_COM_LTE.EARFCNULInUse',
            val_type='int',
            data='39150',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.Capabilities.LTE.DuplexMode',
            val_type='string',
            data='TDDMode',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.Capabilities.LTE.BandsSupported',
            val_type='string',
            data='40',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.ManagementServer.PeriodicInformEnable',
            val_type='int',
            data='5',
        ))
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = param_val_list
        return msg

    def _get_object_param_values_response(
        self,
    ) -> models.GetParameterValuesResponse:
        msg = models.GetParameterValuesResponse()
        param_val_list = []
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.IsPrimary',
            val_type='boolean',
            data='true',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.CellReservedForOperatorUse',
            val_type='boolean',
            data='false',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.PLMNID',
            val_type='string',
            data='00101',
        ))
        param_val_list.append(self._get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.Enable',
            val_type='boolean',
            data='true',
        ))
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = param_val_list
        return msg

    def _get_reboot_response(self) -> models.RebootResponse:
        return models.RebootResponse()
