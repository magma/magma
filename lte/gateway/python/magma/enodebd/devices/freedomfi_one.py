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
from typing import Any, Callable, Dict, Iterable, List, Optional

from dp.protos.cbsd_pb2 import CBSDStateResult
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
from magma.enodebd.device_config.cbrs_consts import (
    BAND,
    SAS_MAX_POWER_SPECTRAL_DENSITY,
    SAS_MIN_POWER_SPECTRAL_DENSITY,
)
from magma.enodebd.device_config.configuration_init import build_desired_config
from magma.enodebd.device_config.configuration_util import (
    calc_bandwidth_mhz,
    calc_bandwidth_rbs,
    calc_earfcn,
)
from magma.enodebd.device_config.enodeb_config_postprocessor import (
    EnodebConfigurationPostProcessor,
)
from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.dp_client import (
    build_enodebd_update_cbsd_request,
    enodebd_update_cbsd,
)
from magma.enodebd.exceptions import ConfigurationError
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
    CheckFirmwareUpgradeDownloadState,
    DeleteObjectsState,
    EnbSendRebootState,
    EndSessionState,
    EnodebAcsState,
    ErrorState,
    FirmwareUpgradeDownloadState,
    GetParametersState,
    NotifyDPState,
    SetParameterValuesState,
    WaitForFirmwareUpgradeDownloadResponse,
    WaitGetParametersState,
    WaitInformMRebootState,
    WaitInformState,
    WaitRebootResponseState,
    WaitSetParameterValuesState,
)
from magma.enodebd.tr069 import models

SAS_KEY = 'sas'
WEB_UI_ENABLE_LIST_KEY = 'web_ui_enable_list'
DP_MODE_KEY = 'dp_mode'

RADIO_MIN_POWER = 0
RADIO_MAX_POWER = 24
ANTENNA_HEIGHT = 0


class SASParameters(object):
    """ Class modeling the SAS parameters and their TR path"""
    # SAS parameters for FreedomFiOne
    FAP_CONTROL = 'Device.Services.FAPService.1.FAPControl.'
    FAPSERVICE_PATH = 'Device.Services.FAPService.1.'

    # Sas management parameters
    # TODO move param definitions to ParameterNames class or think of something to make them more generic across devices
    SAS_ENABLE = "sas_enabled"
    SAS_METHOD = "sas_method"  # 0 = SAS client, 1 = DP mode
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
    SAS_USER_ID = "sas_uid"  # TODO this should be set to a constant value in config
    SAS_FCC_ID = "sas_fcc_id"
    # For CBRS radios we set this to value returned by SAS, eNB can reduce the
    # power if needed.
    SAS_MAX_EIRP = "sas_max_eirp"
    SAS_ANTENNA_GAIN = "sas_antenna_gain"

    SAS_PARAMETERS = {
        SAS_ENABLE: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.Enable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        SAS_METHOD: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.Method',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        SAS_SERVER_URL: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.Server',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SAS_USER_ID: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.UserContactInformation',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SAS_FCC_ID: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.FCCIdentificationNumber',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SAS_CATEGORY: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.Category',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SAS_CHANNEL_TYPE: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.ProtectionLevel',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SAS_CERT_SUBJECT: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.CertSubject',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        # SAS_IC_GROUP_ID: TrParam(
        #     FAP_CONTROL + 'LTE.X_000E8F_SAS.ICGGroupId', is_invasive=False,
        #     type=TrParameterType.STRING, False),
        SAS_LOCATION: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.Location',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SAS_HEIGHT_TYPE: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.HeightType',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SAS_CPI_ENABLE: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.CPIEnable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        SAS_CPI_IPE: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.CPIInstallParamSuppliedEnable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        SAS_MAX_EIRP: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.MaxEirpMHz_Carrier1',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        SAS_ANTENNA_GAIN: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.AntennaGain',
            is_invasive=True,
            type=TrParameterType.INT,
            is_optional=False,
        ),
    }

    # Hardcoded defaults
    defaults = {
        SAS_ENABLE: True,
        SAS_METHOD: False,
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
            STATUS_PATH + 'X_000E8F_DEFGW_Status',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SYNC_STATUS: TrParam(
            STATUS_PATH + 'X_000E8F_Sync_Status',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SAS_STATUS: TrParam(
            STATUS_PATH + 'X_000E8F_SAS_Status',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ENB_STATUS: TrParam(
            STATUS_PATH + 'X_000E8F_eNB_Status',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
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
            InvalidTrParamPath,
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.GPS_STATUS: TrParam(
            InvalidTrParamPath,
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.PTP_STATUS: TrParam(
            InvalidTrParamPath,
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.MME_STATUS: TrParam(
            InvalidTrParamPath,
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.OP_STATE: TrParam(
            InvalidTrParamPath,
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
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

        if (
            name_to_val.get(cls.DEFAULT_GW)
            and name_to_val[cls.DEFAULT_GW].upper() != success_str
        ):
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

        if (
            name_to_val.get(cls.SAS_STATUS)
            and name_to_val[cls.SAS_STATUS].upper() == success_str
        ):
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

        if (
            name_to_val.get(cls.GPS_SCAN_STATUS)
            and name_to_val[cls.GPS_SCAN_STATUS].upper() == success_str
        ):
            device_cfg.set_parameter(
                param_name=ParameterName.GPS_STATUS,
                value=True,
            )
            # Time comes through GPS so can only be insync with GPS is
            # in sync, we use PTP_STATUS field to overload timer is in Sync.
            if (
                name_to_val.get(cls.SYNC_STATUS)
                and name_to_val[cls.SYNC_STATUS].upper() == insync_str
            ):
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

        if (
            name_to_val.get(cls.DEFAULT_GW)
            and name_to_val[cls.DEFAULT_GW].upper() == success_str
        ):
            device_cfg.set_parameter(
                param_name=ParameterName.MME_STATUS,
                value=True,
            )
        else:
            device_cfg.set_parameter(
                param_name=ParameterName.MME_STATUS,
                value=False,
            )

        if (
            name_to_val.get(cls.ENB_STATUS)
            and name_to_val[cls.ENB_STATUS].upper() == success_str
        ):
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

    CONTIGUOUS_CC = "contiguous_cc"
    WEB_UI_ENABLE = "web_ui_enable"  # Enable or disable local enb UI

    MISC_PARAMETERS = {
        WEB_UI_ENABLE: TrParam(
            'Device.X_000E8F_DeviceFeature.X_000E8F_WebServerEnable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        CONTIGUOUS_CC: TrParam(
            FAP_CONTROL
            + 'LTE.X_000E8F_RRMConfig.X_000E8F_CELL_Freq_Contiguous',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        TUNNEL_REF: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.Tunnel.1.TunnelRef',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        PRIM_SOURCE: TrParam(
            FAPSERVICE_PATH + 'REM.X_000E8F_tfcsManagerConfig.primSrc',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
    }

    # Hardcoded defaults
    defaults = {
        # Use IPV4 only
        TUNNEL_REF: "Device.IP.Interface.1.IPv4Address.1.",
        CONTIGUOUS_CC: 0,  # Its not contiguous carrier
        WEB_UI_ENABLE: False,  # Disable WebUI by default
    }


class CarrierAggregationParameters(object):
    """
    eNB parameters related to Carrier Aggregation
    """
    FAP_CONTROL = 'Device.Services.FAPService.1.FAPControl.'
    FAPSERVICE_PATH = 'Device.Services.FAPService.1.'

    CA_ENABLE = "CA Enable"
    CA_CARRIER_NUMBER = "CA Carrier Number"
    CA_EARFCNDL = "CA EARFCNDL"
    CA_EARFCNUL = "CA EARFCNUL"
    CA_CELL_ID = "CA Cell ID"
    CA_TAC = "CA TAC"
    CA_MAX_EIRP_MHZ = "CA Max EIRP"
    CA_BAND = "CA Band"

    CA_PARAMETERS = {
        CA_ENABLE: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_RRMConfig.X_000E8F_CA_Enable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        CA_CARRIER_NUMBER: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_RRMConfig.X_000E8F_Cell_Number',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        CA_EARFCNDL: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.X_000E8F_EARFCNDL2',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        CA_EARFCNUL: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.X_000E8F_EARFCNUL2',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        CA_CELL_ID: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Common.X_000E8F_CellIdentity2',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        CA_TAC: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.X_000E8F_TAC2',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        CA_MAX_EIRP_MHZ: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.MaxEirpMHz_Carrier2',
            is_invasive=False,
            type=TrParameterType.UNSIGNED_INT,
            is_optional=False,
        ),
        CA_BAND: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.X_000E8F_FreqBandIndicator2',
            is_invasive=False,
            type=TrParameterType.UNSIGNED_INT,
            is_optional=False,
        ),
    }

    # By default, disable Carrier Aggregation and switch to Single Carrier mode
    defaults = {
        CA_ENABLE: False,
        CA_CARRIER_NUMBER: 1,
    }


class FreedomFiOneHandler(BasicEnodebAcsStateMachine):
    """
    FreedomFi One State Machine
    """

    def __init__(
            self,
            service,
    ) -> None:
        self._state_map: Dict[str, Any] = {}
        super().__init__(service=service, use_param_key=True)

    def reboot_asap(self) -> None:
        """
        Transition to 'reboot' state
        """
        self.transition('reboot')

    def is_enodeb_connected(self) -> bool:
        """
        Check if enodebd has received an Inform from the enodeb

        Returns:
            bool
        """
        return not isinstance(self.state, WaitInformState)

    def _init_state_map(self) -> None:
        self._state_map = {
            # Inform comes in -> Respond with InformResponse
            'wait_inform': WaitInformState(self, when_done='get_rpc_methods'),
            # If first inform after boot -> GetRpc request comes in, if not
            # empty request comes in => Transition
            'get_rpc_methods': FreedomFiOneGetInitState(
                self,
                when_done='check_fw_upgrade_download',
            ),

            # Download flow
            'check_fw_upgrade_download': CheckFirmwareUpgradeDownloadState(
                self,
                when_download='fw_upgrade_download',
                when_skip='get_transient_params',
            ),
            'fw_upgrade_download': FirmwareUpgradeDownloadState(
                self,
                when_done='wait_fw_upgrade_download_response',
            ),
            'wait_fw_upgrade_download_response': WaitForFirmwareUpgradeDownloadResponse(
                self,
                when_done='get_transient_params',
                when_skip='get_transient_params',
            ),
            # Download flow ends

            # Read transient readonly params.
            'get_transient_params': FreedomFiOneSendGetTransientParametersState(
                self,
                when_done='get_params',
            ),

            'get_params': FreedomFiOneGetObjectParametersState(
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
            'end_session': FreedomFiOneEndSessionState(self, when_dp_mode='notify_dp'),
            'notify_dp': FreedomFiOneNotifyDPState(self, when_inform='wait_inform'),

            # These states are only entered through manual user intervention
            'reboot': EnbSendRebootState(self, when_done='wait_reboot'),
            'wait_reboot': WaitRebootResponseState(
                self,
                when_done='wait_post_reboot_inform',
            ),
            'wait_post_reboot_inform': WaitInformMRebootState(
                self,
                when_done='wait_inform',
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
        """
        Return the device name

        Returns:
            device name
        """
        return EnodebDeviceName.FREEDOMFI_ONE

    @property
    def data_model_class(self) -> DataModel:
        """
        Return the class of the data model

        Returns:
            DataModel
        """
        return FreedomFiOneTrDataModel

    @property
    def config_postprocessor(self) -> EnodebConfigurationPostProcessor:
        """
        Return the instance of config postprocessor

        Returns:
            EnodebConfigurationPostProcessor
        """
        return FreedomFiOneConfigurationInitializer(self)

    @property
    def state_map(self) -> Dict[str, EnodebAcsState]:
        """
        Return the state map for the State Machine

        Returns:
            Dict[str, EnodebAcsState]
        """
        return self._state_map

    @property
    def disconnected_state_name(self) -> str:
        """
        Return the string representation of a disconnected state

        Returns:
            str
        """
        return 'wait_inform'

    @property
    def unexpected_fault_state_name(self) -> str:
        """
        Return the string representation of an unexpected fault state

        Returns:
            str
        """
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
            DEVICE_PATH + 'DeviceInfo.SoftwareVersion',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.SERIAL_NUMBER: TrParam(
            DEVICE_PATH + 'DeviceInfo.SerialNumber',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),

        # RF-related parameters
        ParameterName.EARFCNDL: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.EARFCNDL',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.EARFCNUL: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.EARFCNUL',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.DL_BANDWIDTH: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.DLBandwidth',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.UL_BANDWIDTH: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.ULBandwidth',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.PCI: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.PhyCellID',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.SUBFRAME_ASSIGNMENT: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.PHY.TDDFrame.SubFrameAssignment',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.SPECIAL_SUBFRAME_PATTERN: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.PHY.TDDFrame.SpecialSubframePatterns',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.CELL_ID: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Common.CellIdentity',
            is_invasive=False,
            type=TrParameterType.UNSIGNED_INT,
            is_optional=False,
        ),
        ParameterName.BAND: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.FreqBandIndicator',
            is_invasive=False,
            type=TrParameterType.UNSIGNED_INT,
            is_optional=False,
        ),

        # Readonly LTE state
        ParameterName.ADMIN_STATE: TrParam(
            FAP_CONTROL + 'LTE.AdminState',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.GPS_ENABLE: TrParam(
            DEVICE_PATH + 'FAP.GPS.ScanOnBoot',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),

        # Core network parameters
        ParameterName.MME_IP: TrParam(
            FAP_CONTROL + 'LTE.Gateway.S1SigLinkServerList',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.MME_PORT: TrParam(
            FAP_CONTROL + 'LTE.Gateway.S1SigLinkPort',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),

        ParameterName.NUM_PLMNS: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNListNumberOfEntries',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),

        ParameterName.TAC: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.TAC', is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        # Management server parameters
        ParameterName.PERIODIC_INFORM_ENABLE: TrParam(
            DEVICE_PATH + 'ManagementServer.PeriodicInformEnable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.PERIODIC_INFORM_INTERVAL: TrParam(
            DEVICE_PATH + 'ManagementServer.PeriodicInformInterval',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),

        # Performance management parameters
        ParameterName.PERF_MGMT_ENABLE: TrParam(
            DEVICE_PATH + 'FAP.PerfMgmt.Config.1.Enable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
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
    TRANSFORMS_FOR_ENB: Dict[ParameterName, Callable[[Any], Any]] = {}
    NUM_PLMNS_IN_CONFIG = 1
    for i in range(1, NUM_PLMNS_IN_CONFIG + 1):  # noqa: WPS604
        PARAMETERS[ParameterName.PLMN_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.' % i,
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_CELL_RESERVED % i] = TrParam(
            FAPSERVICE_PATH
            + 'CellConfig.LTE.EPC.PLMNList.%d.CellReservedForOperatorUse' % i,
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_ENABLE % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.Enable' % i,
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_PRIMARY % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.IsPrimary' % i,
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_PLMNID % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.PLMNID' % i,
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        )

    PARAMETERS.update(SASParameters.SAS_PARAMETERS)  # noqa: WPS604
    PARAMETERS.update(FreedomFiOneMiscParameters.MISC_PARAMETERS)  # noqa: WPS604
    PARAMETERS.update(StatusParameters.STATUS_PARAMETERS)  # noqa: WPS604
    # These are stateful parameters that have no tr-69 representation
    PARAMETERS.update(StatusParameters.DERIVED_STATUS_PARAMETERS)  # noqa: WPS604
    PARAMETERS.update(CarrierAggregationParameters.CA_PARAMETERS)  # noqa: WPS604

    TRANSFORMS_FOR_MAGMA = {
        # We don't set these parameters
        ParameterName.BAND_CAPABILITY: transform_for_magma.band_capability,
        ParameterName.DUPLEX_MODE_CAPABILITY: transform_for_magma.duplex_mode,
        ParameterName.GPS_LAT: transform_for_magma.gps_tr181,
        ParameterName.GPS_LONG: transform_for_magma.gps_tr181,
    }

    @classmethod
    def get_parameter(cls, param_name: ParameterName) -> Optional[TrParam]:
        """
        Retrieve parameter by its name

        Args:
            param_name (ParameterName): parameter name to retrieve

        Returns:
            Optional[TrParam]
        """
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
        Retrieve all load parameters

        Returns:
             List[ParameterName]
        """
        return list(cls.PARAMETERS.keys())

    @classmethod
    def get_num_plmns(cls) -> int:
        """
        Retrieve the number of all PLMN parameters

        Returns:
            int
        """
        return cls.NUM_PLMNS_IN_CONFIG

    @classmethod
    def get_parameter_names(cls) -> List[ParameterName]:
        """
        Retrieve all parameter names

        Returns:
            List[ParameterName]
        """
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
        """
        Retrieve parameter names of all objects

        Returns:
            Dict[ParameterName, List[ParameterName]]
        """
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
    def get_sas_param_names(cls) -> Iterable[ParameterName]:
        """
        Retrieve names of SAS parameters

        Returns:
            List[ParameterName]
        """
        return SASParameters.SAS_PARAMETERS.keys()


class FreedomFiOneConfigurationInitializer(EnodebConfigurationPostProcessor):
    """
    Overrides desired config on the State Machine
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
        """
        Add some params to the desired config

        Args:
            mconfig (Any): mconfig
            service_cfg (Any): service config
            desired_cfg (EnodebConfiguration): desired config
        """
        desired_cfg.delete_parameter(ParameterName.BAND)
        desired_cfg.delete_parameter(ParameterName.EARFCNDL)
        desired_cfg.delete_parameter(ParameterName.DL_BANDWIDTH)
        desired_cfg.delete_parameter(ParameterName.UL_BANDWIDTH)

        self._set_default_params(desired_cfg)
        self._increment_param_version_key()
        self._verify_cell_reserved_param(desired_cfg)
        self._verify_ui_enable(service_cfg, desired_cfg)
        self._verify_sas_params(service_cfg, desired_cfg)
        self._set_misc_params_from_service_config(service_cfg, desired_cfg)

    def _set_default_params(self, desired_cfg):
        """Go through default params and set them in desired config"""
        defaults = {
            **FreedomFiOneMiscParameters.defaults,
            **SASParameters.defaults,
            **CarrierAggregationParameters.defaults,
        }
        for name, val in defaults.items():
            desired_cfg.set_parameter(param_name=name, value=val)

    def _increment_param_version_key(self):
        """Bump up the parameter key version"""
        self.acs.parameter_version_inc()

    def _verify_cell_reserved_param(self, desired_cfg):
        """
        Workaround a bug in Sercomm firmware in release 3920, 3921
        where the meaning of CellReservedForOperatorUse is wrong.
        Set to True to ensure the PLMN is not reserved

        Args:
            desired_cfg: desired eNB config
        """
        num_plmns = self.acs.data_model.get_num_plmns()
        for i in range(1, num_plmns + 1):
            object_name = ParameterName.PLMN_N % i
            desired_cfg.set_parameter_for_object(
                param_name=ParameterName.PLMN_N_CELL_RESERVED % i,
                value=True,
                object_name=object_name,
            )

    def _verify_ui_enable(self, service_cfg, desired_cfg):
        if WEB_UI_ENABLE_LIST_KEY in service_cfg:
            serial_nos = service_cfg.get(WEB_UI_ENABLE_LIST_KEY)
            if self.acs.device_cfg.has_parameter(
                    ParameterName.SERIAL_NUMBER,
            ):
                if self.acs.get_parameter(ParameterName.SERIAL_NUMBER) in \
                        serial_nos:
                    desired_cfg.set_parameter(
                        FreedomFiOneMiscParameters.WEB_UI_ENABLE,
                        value=True,
                    )
            else:
                # This should not happen
                EnodebdLogger.error("Serial number unknown for device")

    def _verify_sas_params(self, service_cfg, desired_cfg):
        sas_cfg = service_cfg.get(SAS_KEY)
        if not sas_cfg or sas_cfg[DP_MODE_KEY]:
            desired_cfg.set_parameter(SASParameters.SAS_METHOD, value=True)
            return

        sas_param_names = self.acs.data_model.get_sas_param_names()
        for name, val in sas_cfg.items():
            if name not in sas_param_names:
                EnodebdLogger.warning("Ignoring attribute %s", name)
                continue
            desired_cfg.set_parameter(name, val)

    def _set_misc_params_from_service_config(self, service_cfg, desired_cfg):
        prim_src_name = FreedomFiOneMiscParameters.PRIM_SOURCE
        prim_src_service_cfg_val = service_cfg.get(prim_src_name)
        desired_cfg.set_parameter(prim_src_name, prim_src_service_cfg_val)


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
        """
        Send back a message to enb

        Args:
            message (Any): TR069 message

        Returns:
            AcsMsgAndTransition
        """
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
        """
        Read incoming message

        Args:
            message (Any): TR069 message

        Returns:
            AcsReadMsgResult
        """
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

        return AcsReadMsgResult(
            msg_handled=True,
            next_state=self.done_transition,
        )

    def state_description(self) -> str:
        """
        Describe the state

        Returns:
            str
        """
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
        Send back a message to enb

        Args:
            message (Any): TR069 message

        Returns:
            AcsMsgAndTransition
        """
        if not self._is_rpc_request:
            resp = models.DummyInput()
            return AcsMsgAndTransition(msg=resp, next_state=None)

        resp = models.GetRPCMethodsResponse()
        resp.MethodList = models.MethodList()
        rpc_methods = ['Inform', 'GetRPCMethods', 'TransferComplete']
        resp.MethodList.arrayType = 'xsd:string[%d]' \
                                    % len(rpc_methods)
        resp.MethodList.string = rpc_methods
        # Don't transition to next state wait for the empty HTTP post
        return AcsMsgAndTransition(msg=resp, next_state=None)

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        """
        Read incoming message

        Args:
            message (Any): TR069 message

        Returns:
            AcsReadMsgResult
        """
        # If this is a regular Inform, not after a reboot we'll get an empty
        # message, in this case transition to the next state. We consider
        # this phase as "initialized"
        if isinstance(message, models.DummyInput):
            return AcsReadMsgResult(
                msg_handled=True,
                next_state=self.done_transition,
            )
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
        """
        Describe the state

        Returns:
            str
        """
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
        """
        Return the names of params not belonging to objects that are added/removed

        Args:
            data_model: Data model of eNB

        Returns:
            List[ParameterName]
        """
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
        """
        Send back a message to enb

        Args:
            message (Any): TR069 message

        Returns:
            AcsMsgAndTransition
        """
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
        Read incoming message

        Args:
            message (Any): TR069 message

        Returns:
            AcsReadMsgResult
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

        if len(  # noqa: WPS507
                get_all_objects_to_delete(
                    self.acs.desired_cfg,
                    self.acs.device_cfg,
                ),
        ) > 0:
            return AcsReadMsgResult(
                msg_handled=True,
                next_state=self.rm_obj_transition,
            )
        elif len(  # noqa: WPS507
                get_all_objects_to_add(
                    self.acs.desired_cfg,
                    self.acs.device_cfg,
                ),
        ) > 0:
            return AcsReadMsgResult(
                msg_handled=True,
                next_state=self.add_obj_transition,
            )
        elif len(  # noqa: WPS507
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
        return AcsReadMsgResult(
            msg_handled=True,
            next_state=self.skip_transition,
        )

    def state_description(self) -> str:
        """
        Describe the state

        Returns:
            str
        """
        return 'Getting well known parameters'


class FreedomFiOneEndSessionState(EndSessionState):
    """ To end a TR-069 session, send an empty HTTP response

    We can expect an inform message on
    End Session state, either a queued one or a periodic one
    """

    def __init__(self, acs: EnodebAcsStateMachine, when_dp_mode: str):
        super().__init__(acs)
        self.acs = acs
        self.notify_dp = when_dp_mode

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        """
        Send back a message to enb

        Args:
            message (Any): TR069 message

        Returns:
            AcsMsgAndTransition
        """
        request = models.DummyInput()
        if self.acs.desired_cfg and self.acs.desired_cfg.get_parameter(SASParameters.SAS_METHOD):
            return AcsMsgAndTransition(request, self.notify_dp)
        return AcsMsgAndTransition(request, None)

    def state_description(self) -> str:
        """
        Describe the state

        Returns:
            str
        """
        description = 'Completed provisioning eNB. Awaiting new Inform'
        if self.acs.desired_cfg and self.acs.desired_cfg.get_parameter(SASParameters.SAS_METHOD):
            description = 'Completed initial provisioning of the eNB. Awaiting update from DProxy'
        return description


class FreedomFiOneNotifyDPState(NotifyDPState):
    """
        FreedomFiOne NotifyDPState implementation
    """

    def enter(self):
        """
        Enter the state
        """
        serial_number = self.acs.device_cfg.get_parameter(ParameterName.SERIAL_NUMBER)
        # NOTE: Some params are not available in eNB Data Model, but still required by Domain Proxy
        # antenna_height: need to provide any value, SAS should not care about the value for CBSD cat A indoor.
        #                 Hardcode it.
        # NOTE: In case GPS scan is not completed, eNB reports LAT and LONG values as 0.
        #       Only update CBSD in Domain Proxy when all params are available.
        gps_status = self.acs.device_cfg.get_parameter(ParameterName.GPS_STATUS)
        if gps_status:
            enodebd_update_cbsd_request = build_enodebd_update_cbsd_request(
                serial_number=serial_number,
                latitude_deg=self.acs.device_cfg.get_parameter(ParameterName.GPS_LAT),
                longitude_deg=self.acs.device_cfg.get_parameter(ParameterName.GPS_LONG),
                indoor_deployment=self.acs.device_cfg.get_parameter(SASParameters.SAS_LOCATION),
                antenna_height=str(ANTENNA_HEIGHT),
                antenna_height_type=self.acs.device_cfg.get_parameter(SASParameters.SAS_HEIGHT_TYPE),
                cbsd_category=self.acs.device_cfg.get_parameter(SASParameters.SAS_CATEGORY),
            )
            state = enodebd_update_cbsd(enodebd_update_cbsd_request)
            ff_one_update_desired_config_from_cbsd_state(state, self.acs.desired_cfg)
        else:
            EnodebdLogger.debug("Waiting for GPS to sync, before updating CBSD params in Domain Proxy.")


def _ff_one_check_state_compatibility_with_ca(state: CBSDStateResult) -> bool:
    """
    Check if state returned by Domain Proxy contains data that can be applied
    to FreedomFi One BS in Carrier Aggregation Mode.
    FreedomFi One can apply carrier aggregation if:
    * 2 channels are available:
      * with symmetric bandwidths: 10+10 or 20+20

    Additionally, such channels may be available but Domain Proxy may explicitly disable CA

    Only check the first 2 channels (Domain Proxy may return more)
    """
    _CA_SUPPORTED_BANDWIDTHS_MHZ = (
        (10, 10),
        (20, 20),
    )
    num_of_channels = len(state.channels)
    # Check if CA explicitly disabled, or not enough channels
    if num_of_channels < 2 or not state.carrier_aggregation_enabled:
        EnodebdLogger.debug(f"Domain Proxy state {num_of_channels=}, {state.carrier_aggregation_enabled=}.")
        return False

    ch1 = state.channels[0]
    ch2 = state.channels[1]

    # Check supported bandwidths of the channels
    bw1 = calc_bandwidth_mhz(low_freq_hz=ch1.low_frequency_hz, high_freq_hz=ch1.high_frequency_hz)
    bw2 = calc_bandwidth_mhz(low_freq_hz=ch2.low_frequency_hz, high_freq_hz=ch2.high_frequency_hz)
    if not (bw1, bw2) in _CA_SUPPORTED_BANDWIDTHS_MHZ:
        EnodebdLogger.debug(f"Domain Proxy channel1 {ch1}, channel2 {ch2}, bandwidth configuration not in {_CA_SUPPORTED_BANDWIDTHS_MHZ}.")
        return False

    return True


def ff_one_update_desired_config_from_cbsd_state(state: CBSDStateResult, config: EnodebConfiguration) -> None:
    """
    Call grpc endpoint on the Domain Proxy to update the desired config based on sas grant

    Args:
        state (CBSDStateResult): state result as received from DP
        config (EnodebConfiguration): configuration to update
    """
    EnodebdLogger.debug("Updating desired config based on Domain Proxy state.")
    num_of_channels = len(state.channels)
    radio_enabled = num_of_channels > 0 and state.radio_enabled

    config.set_parameter(ParameterName.ADMIN_STATE, radio_enabled)
    if not radio_enabled:
        return

    # Carrier1
    channel = state.channels[0]
    earfcn = calc_earfcn(channel.low_frequency_hz, channel.high_frequency_hz)
    bandwidth_mhz = calc_bandwidth_mhz(channel.low_frequency_hz, channel.high_frequency_hz)
    bandwidth_rbs = calc_bandwidth_rbs(bandwidth_mhz)
    max_eirp = _calc_max_eirp_for_carrier(channel.max_eirp_dbm_mhz)
    EnodebdLogger.debug(f"Channel1: {earfcn=}, {bandwidth_rbs=}, {max_eirp=}")

    can_enable_carrier_aggregation = _ff_one_check_state_compatibility_with_ca(state)
    EnodebdLogger.debug(f"Should Carrier Aggregation be enabled on eNB: {can_enable_carrier_aggregation=}")

    # Enabling Carrier Aggregation on FreedomFi One eNB means:
    # 1. Set CA_ENABLED to True
    # 2. Set CA_CARRIER_NUMBER to 2
    # 3. Configure appropriate TR nodes for Carrier2 like EARFCNDL/UL etc
    # Otherwise we need to disable Carrier Aggregation on eNB and switch to Single Carrier configuration
    # 1. Set CA_ENABLED to False
    # 2. Set CA_CARRIER_NUMBER to 1
    # Those two nodes should handle everything else accordingly.
    # (NOTE: carrier aggregation may still be enabled on Domain Proxy, but Domain Proxy may not have 2 channels granted by SAS.
    #        In such case, we still have to switch eNB to Single Carrier)
    ca_carrier_number = 2 if can_enable_carrier_aggregation else 1
    ca_enable = True if can_enable_carrier_aggregation else False
    params_to_set = {
        ParameterName.DL_BANDWIDTH: bandwidth_rbs,
        ParameterName.UL_BANDWIDTH: bandwidth_rbs,
        ParameterName.EARFCNDL: earfcn,
        ParameterName.EARFCNUL: earfcn,
        SASParameters.SAS_MAX_EIRP: max_eirp,
        ParameterName.BAND: BAND,
        CarrierAggregationParameters.CA_ENABLE: ca_enable,
        CarrierAggregationParameters.CA_CARRIER_NUMBER: ca_carrier_number,


    }

    if can_enable_carrier_aggregation:
        # Configure Carrier2
        # NOTE: We set CELL_ID to the value of Carrier1 CELL_ID "+1"
        channel = state.channels[1]
        earfcn = calc_earfcn(channel.low_frequency_hz, channel.high_frequency_hz)
        bandwidth_mhz = calc_bandwidth_mhz(channel.low_frequency_hz, channel.high_frequency_hz)
        bandwidth_rbs = calc_bandwidth_rbs(bandwidth_mhz)
        max_eirp = _calc_max_eirp_for_carrier(channel.max_eirp_dbm_mhz)
        EnodebdLogger.debug(f"Channel2: {earfcn=}, {bandwidth_rbs=}, {max_eirp=}")
        params_to_set.update({
            CarrierAggregationParameters.CA_BAND: BAND,
            CarrierAggregationParameters.CA_EARFCNDL: earfcn,
            CarrierAggregationParameters.CA_EARFCNUL: earfcn,
            CarrierAggregationParameters.CA_TAC: config.get_parameter(ParameterName.TAC),
            CarrierAggregationParameters.CA_CELL_ID: config.get_parameter(ParameterName.CELL_ID) + 1,
        })

    for param, value in params_to_set.items():
        config.set_parameter(param, value)


def _calc_max_eirp_for_carrier(max_eirp_dbm_mhz: float) -> int:
    max_eirp = int(max_eirp_dbm_mhz)
    if not SAS_MIN_POWER_SPECTRAL_DENSITY <= max_eirp <= SAS_MAX_POWER_SPECTRAL_DENSITY:  # noqa: WPS508
        raise ConfigurationError(
            'Max EIRP %d exceeds allowed range [%d, %d]' %
            (max_eirp, SAS_MIN_POWER_SPECTRAL_DENSITY, SAS_MAX_POWER_SPECTRAL_DENSITY),
        )
    return max_eirp
