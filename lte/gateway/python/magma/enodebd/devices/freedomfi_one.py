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
from typing import Any, Callable, Dict, List, Optional, Type

from magma.common.service import MagmaService
from magma.enodebd.data_models import transform_for_magma
from magma.enodebd.data_models.data_model import (
    DataModel,
    InvalidTrParamPath,
    TrParam,
)
from magma.enodebd.data_models.data_model_parameters import (
    ParameterName,
    TrParameterType,
)
from magma.enodebd.device_config.configuration_init import build_desired_config
from magma.enodebd.device_config.enodeb_config_postprocessor import (
    EnodebConfigurationPostProcessor,
)
from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.logger import EnodebdLogger
from magma.enodebd.state_machines.acs_state_utils import (
    get_all_objects_to_add,
    get_all_objects_to_delete,
    get_all_param_values_to_set,
    get_params_to_get,
    parse_get_parameter_values_response,
)
from magma.enodebd.state_machines.enb_acs import EnodebAcsStateMachine
from magma.enodebd.state_machines.enb_acs_impl import BasicEnodebAcsStateMachine
from magma.enodebd.state_machines.enb_acs_states import (
    AcsMsgAndTransition,
    AcsReadMsgResult,
    AddObjectsState,
    DeleteObjectsState,
    EnbSendRebootState,
    EndSessionState,
    EnodebAcsState,
    ErrorState,
    GetParametersState,
    SetParameterValuesState,
    WaitGetParametersState,
    WaitInformMRebootState,
    WaitInformState,
    WaitRebootResponseState,
    WaitSetParameterValuesState,
)
from magma.enodebd.tr069 import models


class SASParameters(object):
    """ Class modeling the SAS parameters and their TR path"""
    # SAS parameters for FreedomFiOne
    FAP_CONTROL = 'Device.Services.FAPService.1.FAPControl.'
    FAPSERVICE_PATH = 'Device.Services.FAPService.1.'

    # Sas management parameters
    SAS_ENABLE = "sas_enabled"
    SAS_SERVER_URL = "sas_server_url"
    SAS_UID = "sas_uid"
    SAS_CATEGORY = "sas_category"
    SAS_CHANNEL_TYPE = "sas_channel_type"
    SAS_CERT_SUBJECT = "sas_cert_subject"
    SAS_IC_GROUP_ID = "sas_icg_group_id"
    SAS_LOCATION = "sas_location"
    SAS_HEIGHT_TYPE = "sas_height_type"
    SAS_CPI_ENABLE = "sas_cpi_enable"
    SAS_CPI_IPE = "sas_cpi_ipe"  # Install param supplied enable
    FREQ_BAND_1 = "freq_band_1"
    FREQ_BAND_2 = "freq_band_2"
    # For CBRS radios we set this to the limit and the SAS can reduce the
    # power if needed.
    TX_POWER_CONFIG = "tx_power_config"

    SAS_PARAMETERS = {
        SAS_ENABLE: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.Enable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        SAS_SERVER_URL: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.Server', is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SAS_UID: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.UserContactInformation', is_invasive=False,
            type=TrParameterType.STRING, is_optional=False,
        ),
        SAS_CATEGORY: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.Category', is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SAS_CHANNEL_TYPE: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.ProtectionLevel', is_invasive=False,
            type=TrParameterType.STRING, is_optional=False,
        ),
        SAS_CERT_SUBJECT: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.CertSubject', is_invasive=False,
            type=TrParameterType.STRING, is_optional=False,
        ),
        # SAS_IC_GROUP_ID: TrParam(
        #     FAP_CONTROL + 'LTE.X_000E8F_SAS.ICGGroupId', is_invasive=False,
        #     type=TrParameterType.STRING, False),
        SAS_LOCATION: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.Location', is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SAS_HEIGHT_TYPE: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.HeightType', is_invasive=False,
            type=TrParameterType.STRING, is_optional=False,
        ),
        SAS_CPI_ENABLE: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.CPIEnable', is_invasive=False,
            type=TrParameterType.BOOLEAN, is_optional=False,
        ),
        SAS_CPI_IPE: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.CPIInstallParamSuppliedEnable',
            is_invasive=False, type=TrParameterType.BOOLEAN, is_optional=False,
        ),
        FREQ_BAND_1: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.FreqBandIndicator',
            is_invasive=False, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        ),
        FREQ_BAND_2: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.X_000E8F_FreqBandIndicator2',
            is_invasive=False, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        ),
        TX_POWER_CONFIG: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.X_000E8F_TxPowerConfig',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
    }


class StatusParameters(object):
    """
    Stateful class that converts eNB status to Magma understood status.
    FreedomFiOne has many status params that interact with each other.
    This class maintains the business logic of parsing these status params
    and converting it to Magma understood fields.
    """
    STATUS_PATH = "Device.X_000E8F_DeviceFeature.X_000E8F_NEStatus."

    # Status parameters
    DEFAULT_GW = 'defaultGW'
    SYNC_STATUS = 'syncStatus'
    SAS_STATUS = 'sasStatus'
    ENB_STATUS = 'enbStatus'
    GPS_SCAN_STATUS = 'gpsScanStatus'

    STATUS_PARAMETERS = {
        # Status nodes
        DEFAULT_GW: TrParam(
            STATUS_PATH + 'X_000E8F_DEFGW_Status', is_invasive=False,
            type=TrParameterType.STRING, is_optional=False,
        ),
        SYNC_STATUS: TrParam(
            STATUS_PATH + 'X_000E8F_Sync_Status', is_invasive=False,
            type=TrParameterType.STRING, is_optional=False,
        ),
        SAS_STATUS: TrParam(
            STATUS_PATH + 'X_000E8F_SAS_Status', is_invasive=False,
            type=TrParameterType.STRING, is_optional=False,
        ),
        ENB_STATUS: TrParam(
            STATUS_PATH + 'X_000E8F_eNB_Status', is_invasive=False,
            type=TrParameterType.STRING, is_optional=False,
        ),

        # GPS status, lat, long
        GPS_SCAN_STATUS: TrParam(
            'Device.FAP.GPS.ScanStatus',
            is_invasive=False, type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.GPS_LAT: TrParam(
            'Device.FAP.GPS.LockedLatitude',
            is_invasive=False, type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.GPS_LONG: TrParam(
            'Device.FAP.GPS.LockedLongitude',
            is_invasive=False, type=TrParameterType.STRING, is_optional=False,
        ),
    }

    # Derived status params that don't have tr69 representation.
    DERIVED_STATUS_PARAMETERS = {
        ParameterName.RF_TX_STATUS: TrParam(
            InvalidTrParamPath, is_invasive=False,
            type=TrParameterType.BOOLEAN, is_optional=False,
        ),
        ParameterName.GPS_STATUS: TrParam(
            InvalidTrParamPath, is_invasive=False,
            type=TrParameterType.BOOLEAN, is_optional=False,
        ),
        ParameterName.PTP_STATUS: TrParam(
            InvalidTrParamPath, is_invasive=False,
            type=TrParameterType.BOOLEAN, is_optional=False,
        ),
        ParameterName.MME_STATUS: TrParam(
            InvalidTrParamPath, is_invasive=False,
            type=TrParameterType.BOOLEAN, is_optional=False,
        ),
        ParameterName.OP_STATE: TrParam(
            InvalidTrParamPath, is_invasive=False,
            type=TrParameterType.BOOLEAN, is_optional=False,
        ),
    }

    @classmethod
    def set_magma_device_cfg(
        cls, name_to_val: Dict,
        device_cfg: EnodebConfiguration,
    ):
        """
        Convert FreedomFiOne name_to_val representation to magma device_cfg
        """
        success_str = "SUCCESS"  # String constant returned by radio
        insync_str = "INSYNC"

        if name_to_val.get(cls.DEFAULT_GW) \
                and name_to_val[cls.DEFAULT_GW].upper() != success_str:
            # Nothing will proceed if the eNB doesn't have an IP on the WAN
            serial_num = "unknown"
            if device_cfg.has_parameter(ParameterName.SERIAL_NUMBER):
                serial_num = device_cfg.get_parameter(
                    ParameterName.SERIAL_NUMBER,
                )
            EnodebdLogger.error(
                "Radio with serial number %s doesn't have IP address "
                "on WAN", serial_num,
            )
            device_cfg.set_parameter(
                param_name=ParameterName.RF_TX_STATUS,
                value=False,
            )
            device_cfg.set_parameter(
                param_name=ParameterName.GPS_STATUS,
                value=False,
            )
            device_cfg.set_parameter(
                param_name=ParameterName.PTP_STATUS,
                value=False,
            )
            device_cfg.set_parameter(
                param_name=ParameterName.MME_STATUS,
                value=False,
            )
            device_cfg.set_parameter(
                param_name=ParameterName.OP_STATE,
                value=False,
            )
            return

        if name_to_val.get(cls.SAS_STATUS) \
                and name_to_val[cls.SAS_STATUS].upper() == success_str:
            device_cfg.set_parameter(
                param_name=ParameterName.RF_TX_STATUS,
                value=True,
            )
        else:
            # No sas grant so not transmitting. There is no explicit node for Tx
            # in FreedomFiOne
            device_cfg.set_parameter(
                param_name=ParameterName.RF_TX_STATUS,
                value=False,
            )

        if name_to_val.get(cls.GPS_SCAN_STATUS) \
                and name_to_val[cls.GPS_SCAN_STATUS].upper() == success_str:
            device_cfg.set_parameter(
                param_name=ParameterName.GPS_STATUS,
                value=True,
            )
            # Time comes through GPS so can only be insync with GPS is
            # in sync, we use PTP_STATUS field to overload timer is in Sync.
            if name_to_val.get(cls.SYNC_STATUS) \
                    and name_to_val[cls.SYNC_STATUS].upper() == insync_str:
                device_cfg.set_parameter(
                    param_name=ParameterName.PTP_STATUS,
                    value=True,
                )
            else:
                device_cfg.set_parameter(
                    param_name=ParameterName.PTP_STATUS, value=False,
                )
        else:
            device_cfg.set_parameter(
                param_name=ParameterName.GPS_STATUS,
                value=False,
            )
            device_cfg.set_parameter(
                param_name=ParameterName.PTP_STATUS,
                value=False,
            )

        if name_to_val.get(cls.DEFAULT_GW) \
                and name_to_val[cls.DEFAULT_GW].upper() == success_str:
            device_cfg.set_parameter(
                param_name=ParameterName.MME_STATUS,
                value=True,
            )
        else:
            device_cfg.set_parameter(
                param_name=ParameterName.MME_STATUS,
                value=False,
            )

        if name_to_val.get(cls.ENB_STATUS) \
                and name_to_val[cls.ENB_STATUS].upper() == success_str:
            device_cfg.set_parameter(
                param_name=ParameterName.OP_STATE,
                value=True,
            )
        else:
            device_cfg.set_parameter(
                param_name=ParameterName.OP_STATE,
                value=False,
            )

        pass_through_params = [ParameterName.GPS_LAT, ParameterName.GPS_LONG]
        for name in pass_through_params:
            device_cfg.set_parameter(name, name_to_val[name])


class FreedomFiOneMiscParameters(object):
    """
    Default set of parameters that enable carrier aggregation and other
    miscellaneous properties
    """
    FAP_CONTROL = 'Device.Services.FAPService.1.FAPControl.'
    FAPSERVICE_PATH = 'Device.Services.FAPService.1.'

    # Tunnel ref format clobber it to non IPSEC as we don't support
    # IPSEC
    TUNNEL_REF = "tunnel_ref"
    PRIM_SOURCE = "prim_src"

    # Carrier aggregation
    CARRIER_AGG_ENABLE = "carrier_agg_enable"
    CARRIER_NUMBER = "carrier_number"  # Carrier aggregation params
    CONTIGUOUS_CC = "contiguous_cc"
    WEB_UI_ENABLE = "web_ui_enable"  # Enable or disable local enb UI

    MISC_PARAMETERS = {
        WEB_UI_ENABLE: TrParam(
            'Device.X_000E8F_DeviceFeature.X_000E8F_WebServerEnable', is_invasive=False,
            type=TrParameterType.BOOLEAN, is_optional=False,
        ),
        CARRIER_AGG_ENABLE: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_RRMConfig.X_000E8F_CA_Enable', is_invasive=False,
            type=TrParameterType.BOOLEAN, is_optional=False,
        ),
        CARRIER_NUMBER: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_RRMConfig.X_000E8F_Cell_Number', is_invasive=False,
            type=TrParameterType.INT, is_optional=False,
        ),
        CONTIGUOUS_CC: TrParam(
            FAP_CONTROL
            + 'LTE.X_000E8F_RRMConfig.X_000E8F_CELL_Freq_Contiguous',
            is_invasive=False,
            type=TrParameterType.INT, is_optional=False,
        ),
        TUNNEL_REF: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.Tunnel.1.TunnelRef', is_invasive=False,
            type=TrParameterType.STRING, is_optional=False,
        ),
        PRIM_SOURCE: TrParam(
            FAPSERVICE_PATH + 'REM.X_000E8F_tfcsManagerConfig.primSrc',
            is_invasive=False,
            type=TrParameterType.STRING, is_optional=False,
        ),
    }

    # Hardcoded defaults
    defaults = {
        # Use IPV4 only
        TUNNEL_REF: "Device.IP.Interface.1.IPv4Address.1.",
        # Only synchronize with GPS
        PRIM_SOURCE: "GNSS",
        # Always enable carrier aggregation for the CBRS bands
        CARRIER_AGG_ENABLE: True,
        CARRIER_NUMBER: 2,  # CBRS has two carriers
        CONTIGUOUS_CC: 0,  # Its not contiguous carrier
        WEB_UI_ENABLE: False,  # Disable WebUI by default
    }


class FreedomFiOneHandler(BasicEnodebAcsStateMachine):
    def __init__(
            self,
            service: MagmaService,
    ) -> None:
        self._state_map = {}
        super().__init__(service=service, use_param_key=True)

    def reboot_asap(self) -> None:
        self.transition('reboot')

    def is_enodeb_connected(self) -> bool:
        return not isinstance(self.state, WaitInformState)

    def _init_state_map(self) -> None:
        self._state_map = {
            # Inform comes in -> Respond with InformResponse
            'wait_inform': WaitInformState(self, when_done='get_rpc_methods'),
            # If first inform after boot -> GetRpc request comes in, if not
            # empty request comes in => Transition to get_transient_params
            'get_rpc_methods': FreedomFiOneGetInitState(
                self,
                when_done='get_transient_params',
            ),

            # Read transient readonly params.
            'get_transient_params': FreedomFiOneSendGetTransientParametersState(
                self,
                when_done='get_params',
            ),

            'get_params':
                FreedomFiOneGetObjectParametersState(
                    self,
                    when_delete='delete_objs',
                    when_add='add_objs',
                    when_set='set_params',
                    when_skip='end_session',
                ),

            'delete_objs': DeleteObjectsState(
                self, when_add='add_objs',
                when_skip='set_params',
            ),
            'add_objs': AddObjectsState(self, when_done='set_params'),
            'set_params': SetParameterValuesState(
                self,
                when_done='wait_set_params',
            ),
            'wait_set_params': WaitSetParameterValuesState(
                self,
                when_done='check_get_params',
                when_apply_invasive='check_get_params',
                status_non_zero_allowed=True,
            ),
            'check_get_params': GetParametersState(
                self,
                when_done='check_wait_get_params',
                request_all_params=True,
            ),
            'check_wait_get_params': WaitGetParametersState(
                self,
                when_done='end_session',
            ),
            'end_session': EndSessionState(self),

            # These states are only entered through manual user intervention
            'reboot': EnbSendRebootState(self, when_done='wait_reboot'),
            'wait_reboot': WaitRebootResponseState(
                self,
                when_done='wait_post_reboot_inform',
            ),
            'wait_post_reboot_inform': WaitInformMRebootState(
                self,
                when_done='wait_empty',
                when_timeout='wait_inform',
            ),
            # The states below are entered when an unexpected message type is
            # received
            'unexpected_fault': ErrorState(
                self,
                inform_transition_target='wait_inform',
            ),
        }

    @property
    def device_name(self) -> str:
        return EnodebDeviceName.FREEDOMFI_ONE

    @property
    def data_model_class(self) -> Type[DataModel]:
        return FreedomFiOneTrDataModel

    @property
    def config_postprocessor(self) -> EnodebConfigurationPostProcessor:
        return FreedomFiOneConfigurationInitializer(self)

    @property
    def state_map(self) -> Dict[str, EnodebAcsState]:
        return self._state_map

    @property
    def disconnected_state_name(self) -> str:
        return 'wait_inform'

    @property
    def unexpected_fault_state_name(self) -> str:
        return 'unexpected_fault'


class FreedomFiOneTrDataModel(DataModel):
    """
    Class to represent relevant data model parameters from TR-196/TR-098.
    This class is effectively read-only.

    These models have these idiosyncrasies (on account of running TR098):

    - Parameter content root is different (InternetGatewayDevice)
    - GetParameter queries with a wildcard e.g. InternetGatewayDevice. do
      not respond with the full tree (we have to query all parameters)
    - MME status is not exposed - we assume the MME is connected if
      the eNodeB is transmitting (OpState=true)
    - Parameters such as band capability/duplex config
      are rooted under `boardconf.` and not the device config root
    - Num PLMNs is not reported by these units
    """
    # Mapping of TR parameter paths to aliases
    DEVICE_PATH = 'Device.'
    FAPSERVICE_PATH = DEVICE_PATH + 'Services.FAPService.1.'
    FAP_CONTROL = FAPSERVICE_PATH + 'FAPControl.'
    BCCH = FAPSERVICE_PATH + 'REM.LTE.Cell.1.BCCH.'

    PARAMETERS = {
        # Top-level objects
        ParameterName.DEVICE: TrParam(
            DEVICE_PATH, is_invasive=False, type=TrParameterType.OBJECT,
            is_optional=False,
        ),
        ParameterName.FAP_SERVICE: TrParam(
            FAP_CONTROL, is_invasive=False,
            type=TrParameterType.OBJECT, is_optional=False,
        ),

        # Device info
        ParameterName.SW_VERSION: TrParam(
            DEVICE_PATH + 'DeviceInfo.SoftwareVersion', is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.SERIAL_NUMBER: TrParam(
            DEVICE_PATH + 'DeviceInfo.SerialNumber', is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),

        # RF-related parameters
        ParameterName.EARFCNDL: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.EARFCNDL', is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.DL_BANDWIDTH: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.DLBandwidth', is_invasive=False,
            type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.UL_BANDWIDTH: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.ULBandwidth', is_invasive=False,
            type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.PCI: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.PhyCellID', is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.SUBFRAME_ASSIGNMENT: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.PHY.TDDFrame.SubFrameAssignment',
            is_invasive=False, type=TrParameterType.BOOLEAN, is_optional=False,
        ),
        ParameterName.SPECIAL_SUBFRAME_PATTERN: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.PHY.TDDFrame.SpecialSubframePatterns',
            is_invasive=False, type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.CELL_ID: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Common.CellIdentity', is_invasive=False,
            type=TrParameterType.UNSIGNED_INT, is_optional=False,
        ),

        # Readonly LTE state
        ParameterName.ADMIN_STATE: TrParam(
            FAP_CONTROL + 'LTE.AdminState',
            is_invasive=False, type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.GPS_ENABLE: TrParam(
            DEVICE_PATH + 'FAP.GPS.ScanOnBoot',
            is_invasive=False, type=TrParameterType.BOOLEAN, is_optional=False,
        ),

        # Core network parameters
        ParameterName.MME_IP: TrParam(
            FAP_CONTROL + 'LTE.Gateway.S1SigLinkServerList', is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.MME_PORT: TrParam(
            FAP_CONTROL + 'LTE.Gateway.S1SigLinkPort', is_invasive=False,
            type=TrParameterType.INT, is_optional=False,
        ),

        ParameterName.NUM_PLMNS: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNListNumberOfEntries',
            is_invasive=False,
            type=TrParameterType.INT, is_optional=False,
        ),

        ParameterName.TAC: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.TAC', is_invasive=False,
            type=TrParameterType.INT, is_optional=False,
        ),
        # Management server parameters
        ParameterName.PERIODIC_INFORM_ENABLE: TrParam(
            DEVICE_PATH + 'ManagementServer.PeriodicInformEnable', is_invasive=False,
            type=TrParameterType.BOOLEAN, is_optional=False,
        ),
        ParameterName.PERIODIC_INFORM_INTERVAL: TrParam(
            DEVICE_PATH + 'ManagementServer.PeriodicInformInterval', is_invasive=False,
            type=TrParameterType.INT, is_optional=False,
        ),

        # Performance management parameters
        ParameterName.PERF_MGMT_ENABLE: TrParam(
            DEVICE_PATH + 'FAP.PerfMgmt.Config.1.Enable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN, is_optional=False,
        ),
        ParameterName.PERF_MGMT_UPLOAD_INTERVAL: TrParam(
            DEVICE_PATH + 'FAP.PerfMgmt.Config.1.PeriodicUploadInterval',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.PERF_MGMT_UPLOAD_URL: TrParam(
            DEVICE_PATH + 'FAP.PerfMgmt.Config.1.URL',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
    }
    TRANSFORMS_FOR_ENB = {}
    NUM_PLMNS_IN_CONFIG = 1
    for i in range(1, NUM_PLMNS_IN_CONFIG + 1):
        PARAMETERS[ParameterName.PLMN_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.' % i, is_invasive=False,
            type=TrParameterType.STRING, is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_CELL_RESERVED % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.CellReservedForOperatorUse' % i,
            is_invasive=False,
            type=TrParameterType.BOOLEAN, is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_ENABLE % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.Enable' % i, is_invasive=False,
            type=TrParameterType.BOOLEAN, is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_PRIMARY % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.IsPrimary' % i,
            is_invasive=False, type=TrParameterType.BOOLEAN, is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_PLMNID % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.PLMNID' % i, is_invasive=False,
            type=TrParameterType.STRING, is_optional=False,
        )

    PARAMETERS.update(SASParameters.SAS_PARAMETERS)
    PARAMETERS.update(FreedomFiOneMiscParameters.MISC_PARAMETERS)
    PARAMETERS.update(StatusParameters.STATUS_PARAMETERS)
    # These are stateful parameters that have no tr-69 representation
    PARAMETERS.update(StatusParameters.DERIVED_STATUS_PARAMETERS)

    TRANSFORMS_FOR_MAGMA = {
        # We don't set these parameters
        ParameterName.BAND_CAPABILITY: transform_for_magma.band_capability,
        ParameterName.DUPLEX_MODE_CAPABILITY: transform_for_magma.duplex_mode,
    }

    @classmethod
    def get_parameter(cls, param_name: ParameterName) -> Optional[TrParam]:
        return cls.PARAMETERS.get(param_name)

    @classmethod
    def _get_magma_transforms(
            cls,
    ) -> Dict[ParameterName, Callable[[Any], Any]]:
        return cls.TRANSFORMS_FOR_MAGMA

    @classmethod
    def _get_enb_transforms(cls) -> Dict[ParameterName, Callable[[Any], Any]]:
        return cls.TRANSFORMS_FOR_ENB

    @classmethod
    def get_load_parameters(cls) -> List[ParameterName]:
        """
        Load all the parameters instead of a subset.
        """
        return list(cls.PARAMETERS.keys())

    @classmethod
    def get_num_plmns(cls) -> int:
        return cls.NUM_PLMNS_IN_CONFIG

    @classmethod
    def get_parameter_names(cls) -> List[ParameterName]:
        excluded_params = [
            str(ParameterName.DEVICE),
            str(ParameterName.FAP_SERVICE),
        ]
        names = list(
            filter(
                lambda x: (not str(x).startswith('PLMN'))
                and (str(x) not in excluded_params),
                cls.PARAMETERS.keys(),
            ),
        )
        return names

    @classmethod
    def get_numbered_param_names(
            cls,
    ) -> Dict[ParameterName, List[ParameterName]]:
        names = {}
        for i in range(1, cls.NUM_PLMNS_IN_CONFIG + 1):
            params = [
                ParameterName.PLMN_N_CELL_RESERVED % i,
                ParameterName.PLMN_N_ENABLE % i,
                ParameterName.PLMN_N_PRIMARY % i,
                ParameterName.PLMN_N_PLMNID % i,
            ]
            names[ParameterName.PLMN_N % i] = params

        return names

    @classmethod
    def get_sas_param_names(cls) -> List[ParameterName]:
        return SASParameters.SAS_PARAMETERS.keys()


class FreedomFiOneConfigurationInitializer(EnodebConfigurationPostProcessor):
    """
    Class to add the sas related parameters to the desired config.
    """
    SAS_KEY = 'sas'
    WEB_UI_ENABLE_LIST_KEY = 'web_ui_enable_list'

    def __init__(self, acs: EnodebAcsStateMachine):
        super().__init__()
        self.acs = acs

    def postprocess(
        self, mconfig: Any, service_cfg: Any,
        desired_cfg: EnodebConfiguration,
    ) -> None:
        # TODO: Get this config from the domain proxy
        # TODO @amarpad, set these when DProxy integration is done.
        # For now the radio will directly talk to the SAS and get these
        # attributes.
        desired_cfg.delete_parameter(ParameterName.EARFCNDL)
        desired_cfg.delete_parameter(ParameterName.DL_BANDWIDTH)
        desired_cfg.delete_parameter(ParameterName.UL_BANDWIDTH)

        # go through misc parameters and set them to default.
        for name, val in FreedomFiOneMiscParameters.defaults.items():
            desired_cfg.set_parameter(name, val)

        # Bump up the parameter key version
        self.acs.parameter_version_inc()

        if self.WEB_UI_ENABLE_LIST_KEY in service_cfg:
            serial_nos = service_cfg.get(self.WEB_UI_ENABLE_LIST_KEY)
            if self.acs.device_cfg.has_parameter(
                    ParameterName.SERIAL_NUMBER,
            ):
                if self.acs.get_parameter(ParameterName.SERIAL_NUMBER) in \
                        serial_nos:
                    desired_cfg.set_parameter(
                        FreedomFiOneMiscParameters.WEB_UI_ENABLE,
                        True,
                    )
            else:
                # This should not happen
                EnodebdLogger.error("Serial number unknown for device")

        if self.SAS_KEY not in service_cfg:
            return

        sas_cfg = service_cfg[self.SAS_KEY]
        sas_param_names = self.acs.data_model.get_sas_param_names()
        for name, val in sas_cfg.items():
            if name not in sas_param_names:
                EnodebdLogger.warning("Ignoring attribute %s", name)
                continue
            desired_cfg.set_parameter(name, val)


class FreedomFiOneSendGetTransientParametersState(EnodebAcsState):
    """
    Periodically read eNodeB status. Note: keep frequency low to avoid
    backing up large numbers of read operations if enodebd is busy.
    Some eNB parameters are read only and updated by the eNB itself.
    """

    def __init__(self, acs: EnodebAcsStateMachine, when_done: str):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        request = models.GetParameterValues()
        request.ParameterNames = models.ParameterNames()
        request.ParameterNames.string = []
        for _, tr_param in StatusParameters.STATUS_PARAMETERS.items():
            path = tr_param.path
            request.ParameterNames.string.append(path)
        request.ParameterNames.arrayType = \
            'xsd:string[%d]' % len(request.ParameterNames.string)

        return AcsMsgAndTransition(msg=request, next_state=None)

    def read_msg(self, message: Any) -> AcsReadMsgResult:

        if not isinstance(message, models.GetParameterValuesResponse):
            return AcsReadMsgResult(msg_handled=False, next_state=None)
        # Current values of the fetched parameters
        name_to_val = parse_get_parameter_values_response(
            self.acs.data_model,
            message,
        )
        EnodebdLogger.debug('Received Parameters: %s', str(name_to_val))

        # Update device configuration
        StatusParameters.set_magma_device_cfg(
            name_to_val,
            self.acs.device_cfg,
        )

        return AcsReadMsgResult(msg_handled=True, next_state=self.done_transition)

    def state_description(self) -> str:
        return 'Getting transient read-only parameters'


class FreedomFiOneGetInitState(EnodebAcsState):
    """
    After the first Inform message the following can happen:
    1 - eNB can try to learn the RPC method of the ACS, reply back with the
    RPC response (happens right after boot)
    2 - eNB can send an empty message -> This means that the eNB is already
    provisioned so transition to next state. Only transition to next state
    after this message.
    3 - Some other method call that we don't care about so ignore.
    expected that the eNB -> This is an unhandled state so unlikely
    """

    def __init__(self, acs: EnodebAcsStateMachine, when_done):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done
        self._is_rpc_request = False

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        """
        Return empty message response if care about this
        message type otherwise return empty RPC methods response.
        """
        if not self._is_rpc_request:
            resp = models.DummyInput()
            return AcsMsgAndTransition(msg=resp, next_state=None)

        resp = models.GetRPCMethodsResponse()
        resp.MethodList = models.MethodList()
        RPC_METHODS = ['Inform', 'GetRPCMethods', 'TransferComplete']
        resp.MethodList.arrayType = 'xsd:string[%d]' \
                                    % len(RPC_METHODS)
        resp.MethodList.string = RPC_METHODS
        # Don't transition to next state wait for the empty HTTP post
        return AcsMsgAndTransition(msg=resp, next_state=None)

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        # If this is a regular Inform, not after a reboot we'll get an empty
        # message, in this case transition to the next state. We consider
        # this phase as "initialized"
        if isinstance(message, models.DummyInput):
            return AcsReadMsgResult(msg_handled=True, next_state=self.done_transition)
        if not isinstance(message, models.GetRPCMethods):
            # Unexpected, just don't die, ignore message.
            logging.error("Ignoring message %s", str(type(message)))
            # Set this so get_msg will return an empty message
            self._is_rpc_request = False
        else:
            # Return a valid RPC response
            self._is_rpc_request = True
        return AcsReadMsgResult(msg_handled=True, next_state=None)

    def state_description(self) -> str:
        return 'Initializing the post boot sequence for eNB'


class FreedomFiOneGetObjectParametersState(EnodebAcsState):
    """
    Get information on parameters belonging to objects that can be added or
    removed from the configuration.

    Englewood will report a parameter value as None if it does not exist
    in the data model, rather than replying with a Fault message like most
    eNB devices.
    """

    def __init__(
            self,
            acs: EnodebAcsStateMachine,
            when_delete: str,
            when_add: str,
            when_set: str,
            when_skip: str,
    ):
        super().__init__()
        self.acs = acs
        self.rm_obj_transition = when_delete
        self.add_obj_transition = when_add
        self.set_params_transition = when_set
        self.skip_transition = when_skip

    def get_params_to_get(
        self,
        data_model: DataModel,
    ) -> List[ParameterName]:
        names = []

        # First get base params
        names = get_params_to_get(
            self.acs.device_cfg, self.acs.data_model, request_all_params=True,
        )
        # Add object params.
        num_plmns = data_model.get_num_plmns()
        obj_to_params = data_model.get_numbered_param_names()
        for i in range(1, num_plmns + 1):
            obj_name = ParameterName.PLMN_N % i
            desired = obj_to_params[obj_name]
            names += desired
        return names

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        """ Respond with GetParameterValuesRequest """
        names = self.get_params_to_get(
            self.acs.data_model,
        )

        # Generate the request
        request = models.GetParameterValues()
        request.ParameterNames = models.ParameterNames()
        request.ParameterNames.arrayType = 'xsd:string[%d]' \
                                           % len(names)
        request.ParameterNames.string = []
        for name in names:
            path = self.acs.data_model.get_parameter(name).path
            if path is not InvalidTrParamPath:
                # Only get data elements backed by tr69 path
                request.ParameterNames.string.append(path)

        return AcsMsgAndTransition(msg=request, next_state=None)

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        """
        Process GetParameterValuesResponse
        """
        if not isinstance(message, models.GetParameterValuesResponse):
            return AcsReadMsgResult(msg_handled=False, next_state=None)

        path_to_val = {}
        for param_value_struct in message.ParameterList.ParameterValueStruct:
            path_to_val[param_value_struct.Name] = \
                param_value_struct.Value.Data

        EnodebdLogger.debug('Received object parameters: %s', str(path_to_val))

        # Parse simple params
        param_name_list = self.acs.data_model.get_parameter_names()
        for name in param_name_list:
            path = self.acs.data_model.get_parameter(name).path
            if path in path_to_val:
                value = path_to_val.get(path)
                magma_val = \
                    self.acs.data_model.transform_for_magma(
                        name,
                        value,
                    )
                self.acs.device_cfg.set_parameter(name, magma_val)

        # Parse object params
        num_plmns = self.acs.data_model.get_num_plmns()
        for i in range(1, num_plmns + 1):
            obj_name = ParameterName.PLMN_N % i
            obj_to_params = self.acs.data_model.get_numbered_param_names()
            param_name_list = obj_to_params[obj_name]
            for name in param_name_list:
                path = self.acs.data_model.get_parameter(name).path
                if path in path_to_val:
                    value = path_to_val.get(path)
                    if value is None:
                        continue
                    if obj_name not in self.acs.device_cfg.get_object_names():
                        self.acs.device_cfg.add_object(obj_name)
                    magma_value = \
                        self.acs.data_model.transform_for_magma(name, value)
                    self.acs.device_cfg.set_parameter_for_object(
                        name,
                        magma_value,
                        obj_name,
                    )
        # Now we have enough information to build the desired configuration
        if self.acs.desired_cfg is None:
            self.acs.desired_cfg = build_desired_config(
                self.acs.mconfig,
                self.acs.service_config,
                self.acs.device_cfg,
                self.acs.data_model,
                self.acs.config_postprocessor,
            )

        if len(
                get_all_objects_to_delete(
                    self.acs.desired_cfg,
                    self.acs.device_cfg,
                ),
        ) > 0:
            return AcsReadMsgResult(
                msg_handled=True,
                next_state=self.rm_obj_transition,
            )
        elif len(
                get_all_objects_to_add(
                    self.acs.desired_cfg,
                    self.acs.device_cfg,
                ),
        ) > 0:
            return AcsReadMsgResult(
                msg_handled=True,
                next_state=self.add_obj_transition,
            )
        elif len(
                get_all_param_values_to_set(
                    self.acs.desired_cfg,
                    self.acs.device_cfg,
                    self.acs.data_model,
                ),
        ) > 0:
            return AcsReadMsgResult(
                msg_handled=True,
                next_state=self.set_params_transition,
            )
        return AcsReadMsgResult(msg_handled=True, next_state=self.skip_transition)

    def state_description(self) -> str:
        return 'Getting well known parameters'
