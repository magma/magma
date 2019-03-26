"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from typing import Any
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
    def get_baicells_qafb_inform(cls) -> models.Inform:

        pass

    @classmethod
    def get_inform(
        cls,
        oui: str = '48BF74',
        sw_version: str = 'BaiBS_RTS_3.1.6',
        enb_serial: str = '120200002618AGP0003',
    ) -> models.Inform:
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
        msg.Event.EventStruct = []

        # ParameterList
        val_list = []
        val_list.append(cls.get_parameter_value_struct(
            name='Device.DeviceInfo.HardwareVersion',
            val_type='string',
            data='VER.C',
        ))
        val_list.append(cls.get_parameter_value_struct(
            name='Device.DeviceInfo.ManufacturerOUI',
            val_type='string',
            data=oui,
        ))
        val_list.append(cls.get_parameter_value_struct(
            name='Device.DeviceInfo.SoftwareVersion',
            val_type='string',
            data=sw_version,
        ))
        val_list.append(cls.get_parameter_value_struct(
            name='Device.ManagementServer.ConnectionRequestURL',
            val_type='string',
            data='http://192.168.60.248:7547/25dbc91d31276f0cb03391160531ecae',
        ))
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = val_list

        return msg

    @classmethod
    def get_read_only_param_values_response(
            cls,
    ) -> models.GetParameterValuesResponse:
        msg = models.GetParameterValuesResponse()
        param_val_list = []
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.FAPControl.LTE.OpState',
            val_type='boolean',
            data='false',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.FAPControl.LTE.RFTxStatus',
            val_type='boolean',
            data='false',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.DeviceInfo.X_BAICELLS_COM_GPS_Status',
            val_type='boolean',
            data='0',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.DeviceInfo.X_BAICELLS_COM_1588_Status',
            val_type='boolean',
            data='0',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.DeviceInfo.X_BAICELLS_COM_MME_Status',
            val_type='boolean',
            data='false',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.FAP.GPS.LockedLatitude',
            val_type='int',
            data='0',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.FAP.GPS.LockedLongitude',
            val_type='int',
            data='0',
        ))
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = param_val_list
        return msg

    @classmethod
    def get_regular_param_values_response(
            cls,
    ) -> models.GetParameterValuesResponse:
        msg = models.GetParameterValuesResponse()
        param_val_list = []
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.DLBandwidth',
            val_type='string',
            data='n100',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.FreqBandIndicator',
            val_type='string',
            data='40',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.ManagementServer.PeriodicInformInterval',
            val_type='int',
            data='5',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.RAN.CellRestriction.CellReservedForOperatorUse',
            val_type='boolean',
            data='false',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.ULBandwidth',
            val_type='string',
            data='n100',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.X_BAICELLS_COM_LTE.EARFCNDLInUse',
            val_type='string',
            data='n100',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.RAN.PHY.TDDFrame.SpecialSubframePatterns',
            val_type='int',
            data='7',
        ))
        # MME IP
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.FAPControl.LTE.Gateway.S1SigLinkServerList',
            val_type='string',
            data='"192.168.60.142"',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNListNumberOfEntries',
            val_type='int',
            data='1'
        ))
        # perf mgmt enable
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.FAP.PerfMgmt.Config.1.Enable',
            val_type='boolean',
            data='true'
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.RAN.CellRestriction.CellBarred',
            val_type='boolean',
            data='false'
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.FAP.PerfMgmt.Config.1.PeriodicUploadInterval',
            val_type='int',
            data='300'
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.FAPControl.LTE.AdminState',
            val_type='boolean',
            data='false'
        ))
        # Local gateway enable
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.DeviceInfo.X_BAICELLS_COM_LTE_LGW_Switch',
            val_type='boolean',
            data='0'
        ))
        # Perf mgmt upload url
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.FAP.PerfMgmt.Config.1.URL',
            val_type='string',
            data='http://192.168.60.142:8081/'
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.EPC.TAC',
            val_type='int',
            data='1',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.FAPControl.LTE.Gateway.X_BAICELLS_COM_MmePool.Enable',
            val_type='boolean',
            data='false',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.RAN.PHY.TDDFrame.SubFrameAssignment',
            val_type='int',
            data='2',
        ))
        # PCI
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.RAN.RF.PhyCellID',
            val_type='int',
            data='260',
        ))
        # MME port
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.FAPControl.LTE.Gateway.S1SigLinkPort',
            val_type='int',
            data='36412',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.Ipsec.IPSEC_ENABLE',
            val_type='boolean',
            data='false',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.X_BAICELLS_COM_LTE.EARFCNULInUse',
            val_type='int',
            data='39150',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.Capabilities.LTE.DuplexMode',
            val_type='string',
            data='TDDMode',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.Capabilities.LTE.BandsSupported',
            val_type='string',
            data='40',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.ManagementServer.PeriodicInformEnable',
            val_type='int',
            data='5',
        ))
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = param_val_list
        return msg

    @classmethod
    def get_object_param_values_response(
            cls,
    ) -> models.GetParameterValuesResponse:
        msg = models.GetParameterValuesResponse()
        param_val_list = []
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.IsPrimary',
            val_type='boolean',
            data='true',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.CellReservedForOperatorUse',
            val_type='boolean',
            data='false',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.PLMNID',
            val_type='string',
            data='00101',
        ))
        param_val_list.append(cls.get_parameter_value_struct(
            name='Device.Services.FAPService.1.CellConfig.LTE.EPC.PLMNList.1.Enable',
            val_type='boolean',
            data='true',
        ))
        msg.ParameterList = models.ParameterValueList()
        msg.ParameterList.ParameterValueStruct = param_val_list
        return msg

    @classmethod
    def get_reboot_response(cls) -> models.RebootResponse:
        return models.RebootResponse()
