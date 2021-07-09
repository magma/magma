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

from typing import Any, Callable, Dict, List, Optional, Type

from magma.common.service import MagmaService
from magma.enodebd.data_models import transform_for_enb, transform_for_magma
from magma.enodebd.data_models.data_model import DataModel, TrParam
from magma.enodebd.data_models.data_model_parameters import (
    ParameterName,
    TrParameterType,
)
from magma.enodebd.device_config.enodeb_config_postprocessor import (
    EnodebConfigurationPostProcessor,
)
from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.state_machines.enb_acs_impl import BasicEnodebAcsStateMachine
from magma.enodebd.state_machines.enb_acs_states import (
    AddObjectsState,
    BaicellsRemWaitState,
    BaicellsSendRebootState,
    CheckOptionalParamsState,
    DeleteObjectsState,
    EndSessionState,
    EnodebAcsState,
    ErrorState,
    GetObjectParametersState,
    GetParametersState,
    SendGetTransientParametersState,
    SetParameterValuesState,
    WaitEmptyMessageState,
    WaitGetObjectParametersState,
    WaitGetParametersState,
    WaitGetTransientParametersState,
    WaitInformMRebootState,
    WaitInformState,
    WaitRebootResponseState,
    WaitSetParameterValuesState,
)


class BaicellsHandler(BasicEnodebAcsStateMachine):
    def __init__(
        self,
        service: MagmaService,
    ) -> None:
        self._state_map = {}
        super().__init__(service)

    def reboot_asap(self) -> None:
        self.transition('reboot')

    def is_enodeb_connected(self) -> bool:
        return not isinstance(self.state, WaitInformState)

    def _init_state_map(self) -> None:
        self._state_map = {
            'wait_inform': WaitInformState(self, when_done='wait_empty', when_boot='wait_rem'),
            'wait_rem': BaicellsRemWaitState(self, when_done='wait_inform'),
            'wait_empty': WaitEmptyMessageState(self, when_done='get_transient_params', when_missing='check_optional_params'),
            'check_optional_params': CheckOptionalParamsState(self, when_done='get_transient_params'),
            'get_transient_params': SendGetTransientParametersState(self, when_done='wait_get_transient_params'),
            'wait_get_transient_params': WaitGetTransientParametersState(self, when_get='get_params', when_get_obj_params='get_obj_params', when_delete='delete_objs', when_add='add_objs', when_set='set_params', when_skip='end_session'),
            'get_params': GetParametersState(self, when_done='wait_get_params'),
            'wait_get_params': WaitGetParametersState(self, when_done='get_obj_params'),
            'get_obj_params': GetObjectParametersState(self, when_done='wait_get_obj_params'),
            'wait_get_obj_params': WaitGetObjectParametersState(self, when_delete='delete_objs', when_add='add_objs', when_set='set_params', when_skip='end_session'),
            'delete_objs': DeleteObjectsState(self, when_add='add_objs', when_skip='set_params'),
            'add_objs': AddObjectsState(self, when_done='set_params'),
            'set_params': SetParameterValuesState(self, when_done='wait_set_params'),
            'wait_set_params': WaitSetParameterValuesState(self, when_done='check_get_params', when_apply_invasive='reboot'),
            'check_get_params': GetParametersState(self, when_done='check_wait_get_params', request_all_params=True),
            'check_wait_get_params': WaitGetParametersState(self, when_done='end_session'),
            'end_session': EndSessionState(self),
            'reboot': BaicellsSendRebootState(self, when_done='wait_reboot'),
            'wait_reboot': WaitRebootResponseState(self, when_done='wait_post_reboot_inform'),
            'wait_post_reboot_inform': WaitInformMRebootState(self, when_done='wait_rem_post_reboot', when_timeout='wait_inform_post_reboot'),
            # After rebooting, we don't need to query optional params again.
            'wait_inform_post_reboot': WaitInformState(self, when_done='wait_empty_post_reboot', when_boot='wait_rem_post_reboot'),
            'wait_rem_post_reboot': BaicellsRemWaitState(self, when_done='wait_inform_post_reboot'),
            'wait_empty_post_reboot': WaitEmptyMessageState(self, when_done='get_transient_params', when_missing='check_optional_params'),
            # The states below are entered when an unexpected message type is
            # received
            'unexpected_fault': ErrorState(self, inform_transition_target='wait_inform'),
        }

    @property
    def device_name(self) -> str:
        return EnodebDeviceName.BAICELLS

    @property
    def data_model_class(self) -> Type[DataModel]:
        return BaicellsTrDataModel

    @property
    def config_postprocessor(self) -> EnodebConfigurationPostProcessor:
        return BaicellsTrConfigurationInitializer()

    @property
    def state_map(self) -> Dict[str, EnodebAcsState]:
        return self._state_map

    @property
    def disconnected_state_name(self) -> str:
        return 'wait_inform'

    @property
    def unexpected_fault_state_name(self) -> str:
        return 'unexpected_fault'


class BaicellsTrDataModel(DataModel):
    """
    Class to represent relevant data model parameters from TR-196/TR-098/TR-181.
    This class is effectively read-only

    This is for any version BaiStation_V100R001C00B110SPC003 or after
    """
    # Parameters to query when reading eNodeB config
    LOAD_PARAMETERS = [ParameterName.DEVICE]
    # Mapping of TR parameter paths to aliases
    DEVICE_PATH = 'Device.'
    FAPSERVICE_PATH = DEVICE_PATH + 'Services.FAPService.1.'
    PARAMETERS = {
        # Top-level objects
        ParameterName.DEVICE: TrParam(DEVICE_PATH, True, TrParameterType.OBJECT, False),
        ParameterName.FAP_SERVICE: TrParam(FAPSERVICE_PATH, True, TrParameterType.OBJECT, False),

        # Device info parameters
        ParameterName.GPS_STATUS: TrParam(DEVICE_PATH + 'DeviceInfo.X_BAICELLS_COM_GPS_Status', True, TrParameterType.BOOLEAN, False),
        ParameterName.PTP_STATUS: TrParam(DEVICE_PATH + 'DeviceInfo.X_BAICELLS_COM_1588_Status', True, TrParameterType.BOOLEAN, False),
        ParameterName.MME_STATUS: TrParam(DEVICE_PATH + 'DeviceInfo.X_BAICELLS_COM_MME_Status', True, TrParameterType.BOOLEAN, False),
        ParameterName.REM_STATUS: TrParam(FAPSERVICE_PATH + 'REM.X_BAICELLS_COM_REM_Status', True, TrParameterType.BOOLEAN, False),
        ParameterName.LOCAL_GATEWAY_ENABLE:
            TrParam(DEVICE_PATH + 'DeviceInfo.X_BAICELLS_COM_LTE_LGW_Switch', False, TrParameterType.BOOLEAN, False),
        # Tested Baicells devices were missing this parameter
        ParameterName.GPS_ENABLE: TrParam(DEVICE_PATH + 'X_BAICELLS_COM_GpsSyncEnable', False, TrParameterType.BOOLEAN, True),
        ParameterName.GPS_LAT: TrParam(DEVICE_PATH + 'FAP.GPS.LockedLatitude', True, TrParameterType.INT, True),
        ParameterName.GPS_LONG: TrParam(DEVICE_PATH + 'FAP.GPS.LockedLongitude', True, TrParameterType.INT, True),
        ParameterName.SW_VERSION: TrParam(DEVICE_PATH + 'DeviceInfo.SoftwareVersion', True, TrParameterType.STRING, False),
        ParameterName.SERIAL_NUMBER: TrParam(DEVICE_PATH + 'DeviceInfo.SerialNumber', True, TrParameterType.STRING, False),

        # Capabilities
        ParameterName.DUPLEX_MODE_CAPABILITY: TrParam(
            FAPSERVICE_PATH + 'Capabilities.LTE.DuplexMode', True, TrParameterType.STRING, False,
        ),
        ParameterName.BAND_CAPABILITY: TrParam(FAPSERVICE_PATH + 'Capabilities.LTE.BandsSupported', True, TrParameterType.STRING, False),

        # RF-related parameters
        ParameterName.EARFCNDL: TrParam(FAPSERVICE_PATH + 'X_BAICELLS_COM_LTE.EARFCNDLInUse', True, TrParameterType.INT, False),
        ParameterName.EARFCNUL: TrParam(FAPSERVICE_PATH + 'X_BAICELLS_COM_LTE.EARFCNULInUse', True, TrParameterType.INT, False),
        ParameterName.BAND: TrParam(FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.FreqBandIndicator', True, TrParameterType.INT, False),
        ParameterName.PCI: TrParam(FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.PhyCellID', False, TrParameterType.INT, False),
        ParameterName.DL_BANDWIDTH: TrParam(FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.DLBandwidth', True, TrParameterType.STRING, False),
        ParameterName.UL_BANDWIDTH: TrParam(FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.ULBandwidth', True, TrParameterType.STRING, False),
        ParameterName.SUBFRAME_ASSIGNMENT: TrParam(
            FAPSERVICE_PATH
            + 'CellConfig.LTE.RAN.PHY.TDDFrame.SubFrameAssignment', True, TrParameterType.INT, False,
        ),
        ParameterName.SPECIAL_SUBFRAME_PATTERN: TrParam(
            FAPSERVICE_PATH
            + 'CellConfig.LTE.RAN.PHY.TDDFrame.SpecialSubframePatterns', True, TrParameterType.INT, False,
        ),
        ParameterName.CELL_ID: TrParam(FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Common.CellIdentity', True, TrParameterType.UNSIGNED_INT, False),

        # Other LTE parameters
        ParameterName.ADMIN_STATE: TrParam(FAPSERVICE_PATH + 'FAPControl.LTE.AdminState', False, TrParameterType.BOOLEAN, False),
        ParameterName.OP_STATE: TrParam(FAPSERVICE_PATH + 'FAPControl.LTE.OpState', True, TrParameterType.BOOLEAN, False),
        ParameterName.RF_TX_STATUS: TrParam(FAPSERVICE_PATH + 'FAPControl.LTE.RFTxStatus', True, TrParameterType.BOOLEAN, False),

        # RAN parameters
        ParameterName.CELL_RESERVED: TrParam(
            FAPSERVICE_PATH
            + 'CellConfig.LTE.RAN.CellRestriction.CellReservedForOperatorUse', True, TrParameterType.BOOLEAN, False,
        ),
        ParameterName.CELL_BARRED: TrParam(
            FAPSERVICE_PATH
            + 'CellConfig.LTE.RAN.CellRestriction.CellBarred', True, TrParameterType.BOOLEAN, False,
        ),

        # Core network parameters
        ParameterName.MME_IP: TrParam(
            FAPSERVICE_PATH + 'FAPControl.LTE.Gateway.S1SigLinkServerList', True, TrParameterType.STRING, False,
        ),
        ParameterName.MME_PORT: TrParam(FAPSERVICE_PATH + 'FAPControl.LTE.Gateway.S1SigLinkPort', True, TrParameterType.INT, False),
        ParameterName.NUM_PLMNS: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNListNumberOfEntries', True, TrParameterType.INT, False,
        ),
        ParameterName.PLMN: TrParam(FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.', True, TrParameterType.STRING, False),
        # PLMN arrays are added below
        ParameterName.TAC: TrParam(FAPSERVICE_PATH + 'CellConfig.LTE.EPC.TAC', True, TrParameterType.INT, False),
        ParameterName.IP_SEC_ENABLE: TrParam(
            DEVICE_PATH + 'Services.FAPService.Ipsec.IPSEC_ENABLE', False, TrParameterType.BOOLEAN, False,
        ),
        ParameterName.MME_POOL_ENABLE: TrParam(
            FAPSERVICE_PATH
            + 'FAPControl.LTE.Gateway.X_BAICELLS_COM_MmePool.Enable', True, TrParameterType.BOOLEAN, False,
        ),

        # Management server parameters
        ParameterName.PERIODIC_INFORM_ENABLE:
            TrParam(DEVICE_PATH + 'ManagementServer.PeriodicInformEnable', False, TrParameterType.BOOLEAN, False),
        ParameterName.PERIODIC_INFORM_INTERVAL:
            TrParam(DEVICE_PATH + 'ManagementServer.PeriodicInformInterval', False, TrParameterType.INT, False),

        # Performance management parameters
        ParameterName.PERF_MGMT_ENABLE: TrParam(
            DEVICE_PATH + 'FAP.PerfMgmt.Config.1.Enable', False, TrParameterType.BOOLEAN, False,
        ),
        ParameterName.PERF_MGMT_UPLOAD_INTERVAL: TrParam(
            DEVICE_PATH + 'FAP.PerfMgmt.Config.1.PeriodicUploadInterval', False, TrParameterType.INT, False,
        ),
        ParameterName.PERF_MGMT_UPLOAD_URL: TrParam(
            DEVICE_PATH + 'FAP.PerfMgmt.Config.1.URL', False, TrParameterType.STRING, False,
        ),

    }

    NUM_PLMNS_IN_CONFIG = 6
    for i in range(1, NUM_PLMNS_IN_CONFIG + 1):
        PARAMETERS[(ParameterName.PLMN_N) % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.' % i, True, TrParameterType.STRING, False,
        )
        PARAMETERS[ParameterName.PLMN_N_CELL_RESERVED % i] = TrParam(
            FAPSERVICE_PATH
            + 'CellConfig.LTE.EPC.PLMNList.%d.CellReservedForOperatorUse' % i, True, TrParameterType.BOOLEAN, False,
        )
        PARAMETERS[ParameterName.PLMN_N_ENABLE % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.Enable' % i, True, TrParameterType.BOOLEAN, False,
        )
        PARAMETERS[ParameterName.PLMN_N_PRIMARY % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.IsPrimary' % i, True, TrParameterType.BOOLEAN, False,
        )
        PARAMETERS[ParameterName.PLMN_N_PLMNID % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.PLMNID' % i, True, TrParameterType.STRING, False,
        )

    TRANSFORMS_FOR_ENB = {
        ParameterName.DL_BANDWIDTH: transform_for_enb.bandwidth,
        ParameterName.UL_BANDWIDTH: transform_for_enb.bandwidth,
    }
    TRANSFORMS_FOR_MAGMA = {
        ParameterName.DL_BANDWIDTH: transform_for_magma.bandwidth,
        ParameterName.UL_BANDWIDTH: transform_for_magma.bandwidth,
        # We don't set GPS, so we don't need transform for enb
        ParameterName.GPS_LAT: transform_for_magma.gps_tr181,
        ParameterName.GPS_LONG: transform_for_magma.gps_tr181,
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
        return cls.LOAD_PARAMETERS

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
    def get_numbered_param_names(cls) -> Dict[ParameterName, List[ParameterName]]:
        names = {}
        for i in range(1, cls.NUM_PLMNS_IN_CONFIG + 1):
            params = []
            params.append(ParameterName.PLMN_N_CELL_RESERVED % i)
            params.append(ParameterName.PLMN_N_ENABLE % i)
            params.append(ParameterName.PLMN_N_PRIMARY % i)
            params.append(ParameterName.PLMN_N_PLMNID % i)
            names[ParameterName.PLMN_N % i] = params
        return names


class BaicellsTrConfigurationInitializer(EnodebConfigurationPostProcessor):
    def postprocess(self, desired_cfg: EnodebConfiguration) -> None:
        desired_cfg.set_parameter(ParameterName.CELL_BARRED, False)
