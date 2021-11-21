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
import re
from typing import Any, Callable, Dict, List, Optional, Type

from magma.common.service import MagmaService
from magma.enodebd.data_models import transform_for_enb, transform_for_magma
from magma.enodebd.data_models.data_model import DataModel, TrParam
from magma.enodebd.data_models.data_model_parameters import (
    BaicellsParameterName,
    ParameterName,
    TrParameterType,
)
from magma.enodebd.device_config.configuration_init import build_desired_config
from magma.enodebd.device_config.enodeb_config_postprocessor import (
    EnodebConfigurationPostProcessor,
)
from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.exceptions import Tr069Error
from magma.enodebd.logger import EnodebdLogger as logger
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
    BaicellsRemWaitState,
    CheckOptionalParamsState,
    EnbSendRebootState,
    EndSessionState,
    EnodebAcsState,
    ErrorState,
    GetParametersState,
    SendGetTransientParametersState,
    WaitEmptyMessageState,
    WaitGetParametersState,
    WaitInformMRebootState,
    WaitInformState,
    WaitRebootResponseState,
    WaitSetParameterValuesState,
)
from magma.enodebd.tr069 import models


class BaicellsRTSHandler(BasicEnodebAcsStateMachine):
    def __init__(
            self,
            service: MagmaService,
    ) -> None:
        self._state_map = {}
        super().__init__(service=service, use_param_key=False)

    def reboot_asap(self) -> None:
        self.transition('reboot')

    def is_enodeb_connected(self) -> bool:
        return not isinstance(self.state, WaitInformState)

    def _init_state_map(self) -> None:
        self._state_map = {
            'wait_inform': WaitInformState(self, when_done='wait_empty', when_boot='wait_rem'),
            'wait_rem': BaicellsRemWaitState(self, when_done='wait_inform'),
            'wait_empty': WaitEmptyMessageState(
                self, when_done='get_transient_params',
                when_missing='check_optional_params',
            ),
            'check_optional_params': CheckOptionalParamsState(self, when_done='get_transient_params'),
            'get_transient_params': SendGetTransientParametersState(self, when_done='wait_get_transient_params'),
            'wait_get_transient_params': BaicellsWaitGetTransientParametersState(
                self, when_get='get_params',
                when_get_obj_params='get_obj_params',
                when_delete='delete_objs',
                when_add='add_objs',
                when_set='set_params',
                when_skip='end_session',
            ),
            'get_params': GetParametersState(self, when_done='wait_get_params'),
            'wait_get_params': WaitGetParametersState(self, when_done='get_obj_params'),
            'get_obj_params': BaicellsGetObjectParametersState(self, when_done='wait_get_obj_params'),
            'wait_get_obj_params': BaicellsWaitGetObjectParametersState(
                self, when_delete='delete_objs',
                when_add='add_objs', when_set='set_params',
                when_skip='end_session',
            ),
            'delete_objs': BaicellsDeleteObjectsState(
                self, when_add='add_objs', when_skip='set_params',
            ),
            'add_objs': BaicellsAddObjectsState(self, when_done='set_params'),
            'set_params': BaicellsSetParameterValuesState(self, when_done='wait_set_params'),
            'wait_set_params': WaitSetParameterValuesState(
                self, when_done='check_get_params',
                when_apply_invasive='reboot',
            ),
            'check_get_params': GetParametersState(self, when_done='check_wait_get_params', request_all_params=True),
            'check_wait_get_params': WaitGetParametersState(self, when_done='end_session'),
            'end_session': EndSessionState(self),
            'reboot': EnbSendRebootState(self, when_done='wait_reboot'),
            'wait_reboot': WaitRebootResponseState(self, when_done='wait_post_reboot_inform'),
            'wait_post_reboot_inform': WaitInformMRebootState(
                self, when_done='wait_empty_post_reboot',
                when_timeout='wait_inform_post_reboot',
            ),
            'wait_inform_post_reboot': WaitInformState(
                self, when_done='wait_empty_post_reboot',
                when_boot='wait_rem_post_reboot',
            ),
            'wait_rem_post_reboot': BaicellsRemWaitState(self, when_done='wait_inform_post_reboot'),
            'wait_empty_post_reboot': WaitEmptyMessageState(
                self, when_done='get_transient_params',
                when_missing='check_optional_params',
            ),
            # The states below are entered when an unexpected message type is
            # received
            'unexpected_fault': ErrorState(self, inform_transition_target='wait_inform'),
        }

    @property
    def device_name(self) -> str:
        return EnodebDeviceName.BAICELLS_RTS

    @property
    def data_model_class(self) -> Type[DataModel]:
        return BaicellsRTSTrDataModel

    @property
    def config_postprocessor(self) -> EnodebConfigurationPostProcessor:
        return BaicellsRTSTrConfigurationInitializer()

    @property
    def state_map(self) -> Dict[str, EnodebAcsState]:
        return self._state_map

    @property
    def disconnected_state_name(self) -> str:
        return 'wait_inform'

    @property
    def unexpected_fault_state_name(self) -> str:
        return 'unexpected_fault'


class BaicellsRTSTrDataModel(DataModel):
    """
    Class to represent relevant data model parameters from TR-196/TR-098/TR-181.
    This class is effectively read-only

    This is for any version beginning with BaiBS_ or after
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
        ParameterName.GPS_STATUS: TrParam(
            DEVICE_PATH + 'DeviceInfo.X_BAICELLS_COM_GPS_Status', True,
            TrParameterType.BOOLEAN, False,
        ),
        ParameterName.PTP_STATUS: TrParam(
            DEVICE_PATH + 'DeviceInfo.X_BAICELLS_COM_1588_Status', True,
            TrParameterType.BOOLEAN, False,
        ),
        ParameterName.MME_STATUS: TrParam(
            DEVICE_PATH + 'DeviceInfo.X_BAICELLS_COM_MME_Status', True,
            TrParameterType.BOOLEAN, False,
        ),
        ParameterName.REM_STATUS: TrParam(
            FAPSERVICE_PATH + 'REM.X_BAICELLS_COM_REM_Status', True,
            TrParameterType.BOOLEAN, False,
        ),
        ParameterName.LOCAL_GATEWAY_ENABLE:
            TrParam(DEVICE_PATH + 'DeviceInfo.X_BAICELLS_COM_LTE_LGW_Switch', False, TrParameterType.BOOLEAN, False),
        # Tested Baicells devices were missing this parameter
        ParameterName.GPS_ENABLE: TrParam(
            DEVICE_PATH + 'X_BAICELLS_COM_GpsSyncEnable', False, TrParameterType.BOOLEAN,
            True,
        ),
        ParameterName.GPS_LAT: TrParam(DEVICE_PATH + 'FAP.GPS.LockedLatitude', True, TrParameterType.INT, True),
        ParameterName.GPS_LONG: TrParam(DEVICE_PATH + 'FAP.GPS.LockedLongitude', True, TrParameterType.INT, True),
        ParameterName.SW_VERSION: TrParam(
            DEVICE_PATH + 'DeviceInfo.SoftwareVersion', True, TrParameterType.STRING,
            False,
        ),
        ParameterName.SERIAL_NUMBER: TrParam(
            DEVICE_PATH + 'DeviceInfo.SerialNumber', True, TrParameterType.STRING,
            False,
        ),

        # Capabilities
        ParameterName.DUPLEX_MODE_CAPABILITY: TrParam(
            FAPSERVICE_PATH + 'Capabilities.LTE.DuplexMode', True, TrParameterType.STRING, False,
        ),
        ParameterName.BAND_CAPABILITY: TrParam(
            FAPSERVICE_PATH + 'Capabilities.LTE.BandsSupported', True,
            TrParameterType.STRING, False,
        ),

        # RF-related parameters
        ParameterName.EARFCNDL: TrParam(
            FAPSERVICE_PATH + 'X_BAICELLS_COM_LTE.EARFCNDLInUse', True, TrParameterType.INT,
            False,
        ),
        ParameterName.EARFCNUL: TrParam(
            FAPSERVICE_PATH + 'X_BAICELLS_COM_LTE.EARFCNULInUse', True, TrParameterType.INT,
            False,
        ),
        ParameterName.BAND: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.FreqBandIndicator', True,
            TrParameterType.INT, False,
        ),
        ParameterName.PCI: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.PhyCellID', False, TrParameterType.INT,
            False,
        ),
        ParameterName.DL_BANDWIDTH: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.DLBandwidth', True,
            TrParameterType.STRING, False,
        ),
        ParameterName.UL_BANDWIDTH: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.ULBandwidth', True,
            TrParameterType.STRING, False,
        ),
        ParameterName.SUBFRAME_ASSIGNMENT: TrParam(
            FAPSERVICE_PATH
            + 'CellConfig.LTE.RAN.PHY.TDDFrame.SubFrameAssignment', True, TrParameterType.INT, False,
        ),
        ParameterName.SPECIAL_SUBFRAME_PATTERN: TrParam(
            FAPSERVICE_PATH
            + 'CellConfig.LTE.RAN.PHY.TDDFrame.SpecialSubframePatterns', True, TrParameterType.INT, False,
        ),
        ParameterName.CELL_ID: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Common.CellIdentity', True,
            TrParameterType.UNSIGNED_INT, False,
        ),

        # Other LTE parameters
        ParameterName.ADMIN_STATE: TrParam(
            FAPSERVICE_PATH + 'FAPControl.LTE.AdminState', False,
            TrParameterType.BOOLEAN, False,
        ),
        ParameterName.OP_STATE: TrParam(
            FAPSERVICE_PATH + 'FAPControl.LTE.OpState', True, TrParameterType.BOOLEAN,
            False,
        ),
        ParameterName.RF_TX_STATUS: TrParam(
            FAPSERVICE_PATH + 'FAPControl.LTE.RFTxStatus', True,
            TrParameterType.BOOLEAN, False,
        ),

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
        ParameterName.MME_PORT: TrParam(
            FAPSERVICE_PATH + 'FAPControl.LTE.Gateway.S1SigLinkPort', True,
            TrParameterType.INT, False,
        ),
        ParameterName.NUM_PLMNS: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNListNumberOfEntries', True, TrParameterType.INT, False,
        ),
        BaicellsParameterName.NUM_LTE_NEIGHBOR_FREQ: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.CarrierNumberOfEntries', False,
            TrParameterType.INT, False,
        ),
        BaicellsParameterName.NUM_LTE_NEIGHBOR_CELL: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECellNumberOfEntries', False, TrParameterType.INT,
            False,
        ),
        ParameterName.PLMN: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.', True, TrParameterType.STRING,
            False,
        ),
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
        PARAMETERS[ParameterName.PLMN_N % i] = TrParam(
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

    NUM_NEIGHBOR_CELL_CONFIG = 16
    for i in range(1, NUM_NEIGHBOR_CELL_CONFIG + 1):
        PARAMETERS[BaicellsParameterName.NEIGHBOR_CELL_LIST_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECell.%d.' % i, False, TrParameterType.INT, False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_CELL_CELL_ID_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECell.%d.CID' % i, False, TrParameterType.UNSIGNED_INT,
            False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_CELL_PLMN_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECell.%d.PLMNID' % i, False, TrParameterType.STRING,
            False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_CELL_EARFCN_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECell.%d.EUTRACarrierARFCN' % i, False,
            TrParameterType.UNSIGNED_INT, False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_CELL_PCI_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECell.%d.PhyCellID' % i, False,
            TrParameterType.UNSIGNED_INT, False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_CELL_TAC_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECell.%d.X_BAICELLS_COM_TAC' % i, False,
            TrParameterType.UNSIGNED_INT, False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_CELL_QOFFSET_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECell.%d.QOffset' % i, False, TrParameterType.INT,
            False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_CELL_CIO_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECell.%d.CIO' % i, False, TrParameterType.INT, False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_CELL_ENABLE_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECell.%d.Enable' % i, False, TrParameterType.BOOLEAN,
            False,
        )

    NUM_NEIGHBOR_FREQ_CONFIG = 8
    for i in range(1, NUM_NEIGHBOR_FREQ_CONFIG + 1):
        PARAMETERS[BaicellsParameterName.NEGIH_FREQ_LIST % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.' % i, False,
            TrParameterType.UNSIGNED_INT, False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_EARFCN_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.EUTRACarrierARFCN' % i, False,
            TrParameterType.UNSIGNED_INT, False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_QRXLEVMINSIB5_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.QRxLevMinSIB5' % i, False,
            TrParameterType.INT, False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_Q_OFFSETRANGE_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.QOffsetFreq' % i, False,
            TrParameterType.INT, False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_TRESELECTIONEUTRA_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.TReselectionEUTRA' % i, False,
            TrParameterType.UNSIGNED_INT, False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_RESELECTIONPRIORITY_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.CellReselectionPriority' % i,
            True,
            TrParameterType.UNSIGNED_INT, False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_RESELTHRESHHIGH_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.ThreshXHigh' % i, False,
            TrParameterType.UNSIGNED_INT, False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_RESELTHRESHLOW_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.ThreshXLow' % i, False,
            TrParameterType.UNSIGNED_INT, False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_PMAX_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.PMax' % i, False,
            TrParameterType.INT, False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_TRESELECTIONEUTRASFMEDIUM_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.TReselectionEUTRASFMedium' % i,
            False,
            TrParameterType.UNSIGNED_INT, False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_ENABLE_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.Enable' % i, False,
            TrParameterType.BOOLEAN, False,

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
    def get_num_neighbor_freq(cls) -> int:
        return cls.NUM_NEIGHBOR_FREQ_CONFIG

    @classmethod
    def get_num_neighbor_cell(cls) -> int:
        return cls.NUM_NEIGHBOR_CELL_CONFIG

    @classmethod
    def get_parameter_names(cls) -> List[ParameterName]:
        excluded_params = [
            str(ParameterName.DEVICE),
            str(ParameterName.FAP_SERVICE),
        ]
        names = list(
            filter(
                lambda x: (not str(x).startswith('PLMN')) and (not str(x).startswith('neighbor'))
                and (str(x) not in excluded_params),
                cls.PARAMETERS.keys(),
            ),
        )
        return names

    @classmethod
    def get_numbered_param_names(cls) -> Dict[ParameterName, List[ParameterName]]:
        names = {}
        for i in range(1, cls.NUM_PLMNS_IN_CONFIG + 1):
            params = [
                ParameterName.PLMN_N_CELL_RESERVED % i,
                ParameterName.PLMN_N_ENABLE % i,
                ParameterName.PLMN_N_PRIMARY % i,
                ParameterName.PLMN_N_PLMNID % i,
            ]
            names[ParameterName.PLMN_N % i] = params
        for i in range(1, cls.NUM_NEIGHBOR_FREQ_CONFIG + 1):
            params = [
                BaicellsParameterName.NEIGHBOR_FREQ_ENABLE_N % i,
                BaicellsParameterName.NEIGHBOR_FREQ_EARFCN_N % i,
                BaicellsParameterName.NEIGHBOR_FREQ_PMAX_N % i,
                BaicellsParameterName.NEIGHBOR_FREQ_Q_OFFSETRANGE_N % i,
                BaicellsParameterName.NEIGHBOR_FREQ_Q_OFFSETRANGE_N % i,
                BaicellsParameterName.NEIGHBOR_FREQ_RESELTHRESHLOW_N % i,
                BaicellsParameterName.NEIGHBOR_FREQ_RESELTHRESHHIGH_N % i,
                BaicellsParameterName.NEIGHBOR_FREQ_RESELECTIONPRIORITY_N % i,
                BaicellsParameterName.NEIGHBOR_FREQ_QRXLEVMINSIB5_N % i,
                BaicellsParameterName.NEIGHBOR_FREQ_TRESELECTIONEUTRA_N % i,
            ]
            names[BaicellsParameterName.NEGIH_FREQ_LIST % i] = params
        for i in range(1, cls.NUM_NEIGHBOR_CELL_CONFIG + 1):
            params = [
                BaicellsParameterName.NEIGHBOR_CELL_ENABLE_N % i,
                BaicellsParameterName.NEIGHBOR_CELL_PLMN_N % i,
                BaicellsParameterName.NEIGHBOR_CELL_CELL_ID_N % i,
                BaicellsParameterName.NEIGHBOR_CELL_EARFCN_N % i,
                BaicellsParameterName.NEIGHBOR_CELL_PCI_N % i,
                BaicellsParameterName.NEIGHBOR_CELL_TAC_N % i,
                BaicellsParameterName.NEIGHBOR_CELL_QOFFSET_N % i,
                BaicellsParameterName.NEIGHBOR_CELL_CIO_N % i,
            ]
            names[BaicellsParameterName.NEIGHBOR_CELL_LIST_N % i] = params
        return names


class BaicellsRTSTrConfigurationInitializer(EnodebConfigurationPostProcessor):
    def postprocess(self, mconfig: Any, service_cfg: Any, desired_cfg: EnodebConfiguration) -> None:
        desired_cfg.set_parameter(ParameterName.CELL_BARRED, False)


class BaicellsWaitGetObjectParametersState(EnodebAcsState):
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

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        """ Process GetParameterValuesResponse """
        if not isinstance(
                message,
                models.GetParameterValuesResponse,
        ):
            return AcsReadMsgResult(False, None)

        path_to_val = {}
        if hasattr(message.ParameterList, 'ParameterValueStruct') and \
                message.ParameterList.ParameterValueStruct is not None:
            for param_value_struct in message.ParameterList.ParameterValueStruct:
                path_to_val[param_value_struct.Name] = \
                    param_value_struct.Value.Data
        logger.debug('Received object parameters: %s', str(path_to_val))

        # Number of PLMN objects reported can be incorrect. Let's count them
        num_plmns = 0
        obj_to_params = self.acs.data_model.get_numbered_param_names()
        logger.info('enb obj_to_params= %s', obj_to_params)
        while True:
            obj_name = ParameterName.PLMN_N % (num_plmns + 1)
            if obj_name not in obj_to_params or len(obj_to_params[obj_name]) == 0:
                logger.warning(
                    "eNB has PLMN %s but not defined in model",
                    obj_name,
                )
                break
            param_name_list = obj_to_params[obj_name]
            obj_path = self.acs.data_model.get_parameter(param_name_list[0]).path
            if obj_path not in path_to_val:
                break
            if not self.acs.device_cfg.has_object(obj_name):
                self.acs.device_cfg.add_object(obj_name)
            num_plmns += 1
            for name in param_name_list:
                path = self.acs.data_model.get_parameter(name).path
                value = path_to_val[path]
                magma_val = \
                    self.acs.data_model.transform_for_magma(name, value)
                self.acs.device_cfg.set_parameter_for_object(
                    name, magma_val,
                    obj_name,
                )
        num_plmns_reported = \
            int(self.acs.device_cfg.get_parameter(ParameterName.NUM_PLMNS))
        if num_plmns != num_plmns_reported:
            logger.warning(
                "eNB reported %d PLMNs but found %d",
                num_plmns_reported, num_plmns,
            )
            self.acs.device_cfg.set_parameter(
                ParameterName.NUM_PLMNS,
                num_plmns,
            )
        # Number of Neighbor Freq objects reported can be incorrect. Let's count them
        num_neighbor = 0
        while True:
            obj_name = BaicellsParameterName.NEGIH_FREQ_LIST % (num_neighbor + 1)
            logger.info('enb obj_name= %s', obj_name)
            if obj_name not in obj_to_params or len(obj_to_params[obj_name]) == 0:
                logger.warning(
                    "eNB has Neighbor %s but not defined in model",
                    obj_name,
                )
                break
            param_name_list = obj_to_params[obj_name]
            obj_path = self.acs.data_model.get_parameter(param_name_list[0]).path
            if obj_path not in path_to_val:
                break
            if not self.acs.device_cfg.has_object(obj_name):
                self.acs.device_cfg.add_object(obj_name)
            num_neighbor = num_neighbor + 1
            for name in param_name_list:
                path = self.acs.data_model.get_parameter(name).path
                value = path_to_val[path]
                magma_val = \
                    self.acs.data_model.transform_for_magma(name, value)
                self.acs.device_cfg.set_parameter_for_object(
                    name, magma_val, obj_name,
                )
        num_neighbor_reported = \
            int(self.acs.device_cfg.get_parameter(BaicellsParameterName.NUM_LTE_NEIGHBOR_FREQ))
        if num_neighbor != num_neighbor_reported:
            logger.warning(
                "eNB reported %d Neighbor but found %d",
                num_neighbor_reported, num_neighbor,
            )
            self.acs.device_cfg.set_parameter(
                BaicellsParameterName.NUM_LTE_NEIGHBOR_FREQ,
                num_neighbor,
            )
        # Number of Neighbor Cell objects reported can be incorrect. Let's count them
        num_neighbor_cell = 0
        while True:
            obj_name = BaicellsParameterName.NEIGHBOR_CELL_LIST_N % (num_neighbor_cell + 1)
            if obj_name not in obj_to_params or len(obj_to_params[obj_name]) == 0:
                logger.warning(
                    "eNB has Neighbor %s but not defined in model",
                    obj_name,
                )
                break
            param_name_list = obj_to_params[obj_name]
            obj_path = self.acs.data_model.get_parameter(param_name_list[0]).path
            if obj_path not in path_to_val:
                break
            if not self.acs.device_cfg.has_object(obj_name):
                self.acs.device_cfg.add_object(obj_name)
            num_neighbor_cell = num_neighbor_cell + 1
            for name in param_name_list:
                path = self.acs.data_model.get_parameter(name).path
                value = path_to_val[path]
                magma_val = \
                    self.acs.data_model.transform_for_magma(name, value)
                self.acs.device_cfg.set_parameter_for_object(
                    name, magma_val, obj_name,
                )
        num_neighbor_cell_reported = int(self.acs.device_cfg.get_parameter(BaicellsParameterName.NUM_LTE_NEIGHBOR_CELL))
        if num_neighbor_cell != num_neighbor_cell_reported:
            logger.warning(
                "eNB reported %d neighbor cell but found %d",
                num_neighbor_cell_reported, num_neighbor_cell,
            )
            self.acs.device_cfg.set_parameter(
                BaicellsParameterName.NUM_LTE_NEIGHBOR_CELL,
                num_neighbor_cell,
            )

        # Now we can have the desired state
        if self.acs.desired_cfg is None:
            self.acs.desired_cfg = build_desired_config(
                self.acs.mconfig,
                self.acs.service_config,
                self.acs.device_cfg,
                self.acs.data_model,
                self.acs.config_postprocessor,
            )
        logger.debug('the building desired config is %s', self.acs.desired_cfg.get_object_names())
        logger.debug('the building device_cfg  is------%s-------', self.acs.device_cfg.get_object_names())
        if len(
                get_all_objects_to_delete(
                    self.acs.desired_cfg,
                    self.acs.device_cfg,
                ),
        ) > 0:

            return AcsReadMsgResult(True, self.rm_obj_transition)
        elif len(
                get_all_objects_to_add(
                    self.acs.desired_cfg,
                    self.acs.device_cfg,
                ),
        ) > 0:
            return AcsReadMsgResult(True, self.add_obj_transition)
        elif len(
                get_all_param_values_to_set(
                    self.acs.desired_cfg,
                    self.acs.device_cfg,
                    self.acs.data_model,
                ),
        ) > 0:
            return AcsReadMsgResult(True, self.set_params_transition)
        return AcsReadMsgResult(True, self.skip_transition)

    def state_description(self) -> str:
        return 'Getting object parameters'


class BaicellsDeleteObjectsState(EnodebAcsState):
    def __init__(
            self,
            acs: EnodebAcsStateMachine,
            when_add: str,
            when_skip: str,
    ):
        super().__init__()
        self.acs = acs
        self.deleted_param = None
        self.add_obj_transition = when_add
        self.skip_transition = when_skip

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        """
        Send DeleteObject message to TR-069 and poll for response(s).
        Input:
            - Object name (string)
        """
        request = models.DeleteObject()
        self.deleted_param = get_all_objects_to_delete(
            self.acs.desired_cfg,
            self.acs.device_cfg,
        )[0]
        logger.debug('get obj to delete %s', self.deleted_param)
        request.ObjectName = \
            self.acs.data_model.get_parameter(self.deleted_param).path
        return AcsMsgAndTransition(request, None)

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        """
        Send DeleteObject message to TR-069 and poll for response(s).
        Input:
            - Object name (string)
        """
        if type(message) == models.DeleteObjectResponse:
            if message.Status != 0:
                raise Tr069Error(
                    'Received DeleteObjectResponse with '
                    'Status=%d' % message.Status,
                )
        elif type(message) == models.Fault:
            raise Tr069Error(
                'Received Fault in response to DeleteObject '
                '(faultstring = %s)' % message.FaultString,
            )
        else:
            return AcsReadMsgResult(False, None)

        self.acs.device_cfg.delete_object(self.deleted_param)
        obj_list_to_delete = get_all_objects_to_delete(
            self.acs.desired_cfg,
            self.acs.device_cfg,
        )
        if len(obj_list_to_delete) > 0:
            return AcsReadMsgResult(True, None)
        if len(
                get_all_objects_to_add(
                    self.acs.desired_cfg,
                    self.acs.device_cfg,
                ),
        ) == 0:
            return AcsReadMsgResult(True, self.skip_transition)
        return AcsReadMsgResult(True, self.add_obj_transition)

    def state_description(self) -> str:
        return 'Deleting objects'


class BaicellsGetObjectParametersState(EnodebAcsState):
    def __init__(self, acs: EnodebAcsStateMachine, when_done: str):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        """ Respond with GetParameterValuesRequest """
        names = get_object_params_to_get(
            self.acs.desired_cfg,
            self.acs.device_cfg,
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
            request.ParameterNames.string.append(path)

        return AcsMsgAndTransition(request, self.done_transition)

    def state_description(self) -> str:
        return 'Getting object parameters'


class BaicellsWaitGetTransientParametersState(EnodebAcsState):
    """
    Periodically read eNodeB status. Note: keep frequency low to avoid
    backing up large numbers of read operations if enodebd is busy
    """

    def __init__(
            self,
            acs: EnodebAcsStateMachine,
            when_get: str,
            when_get_obj_params: str,
            when_delete: str,
            when_add: str,
            when_set: str,
            when_skip: str,
    ):
        super().__init__()
        self.acs = acs
        self.done_transition = when_get
        self.get_obj_params_transition = when_get_obj_params
        self.rm_obj_transition = when_delete
        self.add_obj_transition = when_add
        self.set_transition = when_set
        self.skip_transition = when_skip

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        if not isinstance(message, models.GetParameterValuesResponse):
            return AcsReadMsgResult(False, None)
        # Current values of the fetched parameters
        name_to_val = parse_get_parameter_values_response(
            self.acs.data_model,
            message,
        )
        logger.debug('Fetched Transient Params: %s', str(name_to_val))

        # Update device configuration
        for name in name_to_val:
            magma_val = \
                self.acs.data_model.transform_for_magma(
                    name,
                    name_to_val[name],
                )
            self.acs.device_cfg.set_parameter(name, magma_val)

        return AcsReadMsgResult(True, self.get_next_state())

    def get_next_state(self) -> str:
        should_get_params = \
            len(
                get_params_to_get(
                    self.acs.device_cfg,
                    self.acs.data_model,
                ),
            ) > 0
        if should_get_params:
            return self.done_transition
        should_get_obj_params = \
            len(
                get_object_params_to_get(
                    self.acs.desired_cfg,
                    self.acs.device_cfg,
                    self.acs.data_model,
                ),
            ) > 0
        logger.debug(
            'get object param to get %s', get_object_params_to_get(
                self.acs.desired_cfg,
                self.acs.device_cfg,
                self.acs.data_model,
            ),
        )
        if should_get_obj_params:
            return self.get_obj_params_transition
        elif len(
                get_all_objects_to_delete(
                    self.acs.desired_cfg,
                    self.acs.device_cfg,
                ),
        ) > 0:
            return self.rm_obj_transition
        elif len(
                get_all_objects_to_add(
                    self.acs.desired_cfg,
                    self.acs.device_cfg,
                ),
        ) > 0:
            return self.add_obj_transition
        return self.skip_transition

    def state_description(self) -> str:
        return 'Getting transient read-only parameters'


def get_object_params_to_get(
        desired_cfg: Optional[EnodebConfiguration],
        device_cfg: EnodebConfiguration,
        data_model: DataModel,
) -> List[ParameterName]:
    """
    Returns a list of parameter names for object parameters we don't know the
    current value of
    """
    names = []
    # TODO: This might a string for some strange reason, investigate why
    num_plmns = \
        int(device_cfg.get_parameter(ParameterName.NUM_PLMNS))
    for i in range(1, num_plmns + 1):
        obj_name = ParameterName.PLMN_N % i
        if not device_cfg.has_object(obj_name):
            device_cfg.add_object(obj_name)
        obj_to_params = data_model.get_numbered_param_names()
        desired = obj_to_params[obj_name]
        current = []
        if desired_cfg is not None:
            current = desired_cfg.get_parameter_names_for_object(obj_name)
        names_to_add = list(set(desired) - set(current))
        names = names + names_to_add
    num_neighbor_freq = int(device_cfg.get_parameter(BaicellsParameterName.NUM_LTE_NEIGHBOR_FREQ))
    for i in range(1, num_neighbor_freq + 1):
        obj_name = BaicellsParameterName.NEGIH_FREQ_LIST % i
        if not device_cfg.has_object(obj_name):
            device_cfg.add_object(obj_name)
        obj_to_params = data_model.get_numbered_param_names()
        desired = obj_to_params[obj_name]
        current = []
        if desired_cfg is not None and desired_cfg.has_object():
            current = desired_cfg.get_parameter_names_for_object(obj_name)
        names_to_add = list(set(desired) - set(current))
        names = names + names_to_add
    num_neighbor_cell = int(device_cfg.get_parameter(BaicellsParameterName.NUM_LTE_NEIGHBOR_CELL))
    for i in range(1, num_neighbor_cell + 1):
        obj_name = BaicellsParameterName.NEIGHBOR_CELL_LIST_N % i
        if not device_cfg.has_object(obj_name):
            device_cfg.add_object(obj_name)
        obj_to_params = data_model.get_numbered_param_names()
        desired = obj_to_params[obj_name]
        current = []
        if desired_cfg is not None and desired_cfg.has_object(obj_name):
            current = desired_cfg.get_parameter_names_for_object(obj_name)
        names_to_add = list(set(desired) - set(current))
        names = names + names_to_add
    return names


class BaicellsAddObjectsState(EnodebAcsState):
    def __init__(self, acs: EnodebAcsStateMachine, when_done: str):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done
        self.added_param = None

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        request = models.AddObject()
        self.added_param = get_all_objects_to_add(
            self.acs.desired_cfg,
            self.acs.device_cfg,
        )[0]
        desired_param = self.acs.data_model.get_parameter(self.added_param)
        desired_path = desired_param.path
        path_parts = desired_path.split('.')
        # If adding enumerated object, ie. XX.N. we should add it to the
        # parent object XX. so strip the index
        if len(path_parts) > 2 and \
                path_parts[-1] == '' and path_parts[-2].isnumeric():
            logger.debug('Stripping index from path=%s', desired_path)
            desired_path = '.'.join(path_parts[:-2]) + '.'
        request.ObjectName = desired_path
        return AcsMsgAndTransition(request, None)

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        if type(message) == models.AddObjectResponse:
            if message.Status != 0:
                raise Tr069Error(
                    'Received AddObjectResponse with '
                    'Status=%d' % message.Status,
                )
        elif type(message) == models.Fault:
            raise Tr069Error(
                'Received Fault in response to AddObject '
                '(faultstring = %s)' % message.FaultString,
            )
        else:
            return AcsReadMsgResult(False, None)
        instance_n = message.InstanceNumber
        self.added_param = re.sub(r'\d', str(instance_n), self.added_param)
        self.acs.device_cfg.add_object(self.added_param)
        obj_list_to_add = get_all_objects_to_add(
            self.acs.desired_cfg,
            self.acs.device_cfg,
        )
        if len(obj_list_to_add) > 0:
            return AcsReadMsgResult(True, None)
        return AcsReadMsgResult(True, self.done_transition)

    def state_description(self) -> str:
        return 'Adding objects'


class BaicellsSetParameterValuesState(EnodebAcsState):
    def __init__(self, acs: EnodebAcsStateMachine, when_done: str):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        request = models.SetParameterValues()
        request.ParameterList = models.ParameterValueList()
        param_values = get_all_param_values_to_set(
            self.acs.desired_cfg,
            self.acs.device_cfg,
            self.acs.data_model,
            exclude_admin=True,
        )
        request.ParameterList.arrayType = 'cwmp:ParameterValueStruct[%d]' \
                                          % len(param_values)
        request.ParameterList.ParameterValueStruct = []
        logger.debug(
            'Sending TR069 request to set CPE parameter values: %s',
            str(param_values),
        )
        # TODO: Match key response when we support having multiple outstanding
        # calls.
        if self.acs.has_version_key:
            request.ParameterKey = models.ParameterKeyType()
            request.ParameterKey.Data = \
                "SetParameter-{:10.0f}".format(self.acs.parameter_version_key)
            request.ParameterKey.type = 'xsd:string'

        for name, value in param_values.items():
            param_info = self.acs.data_model.get_parameter(name)
            type_ = param_info.type
            name_value = models.ParameterValueStruct()
            name_value.Value = models.anySimpleType()
            name_value.Name = param_info.path
            enb_value = self.acs.data_model.transform_for_enb(name, value)
            if type_ in ('int', 'unsignedInt'):
                name_value.Value.type = 'xsd:%s' % type_
                name_value.Value.Data = str(enb_value)
            elif type_ == 'boolean':
                # Boolean values have integral representations in spec
                name_value.Value.type = 'xsd:boolean'
                name_value.Value.Data = str(int(enb_value))
            elif type_ == 'string':
                name_value.Value.type = 'xsd:string'
                name_value.Value.Data = str(enb_value)
            else:
                raise Tr069Error(
                    'Unsupported type for %s: %s' %
                    (name, type_),
                )
            if param_info.is_invasive:
                self.acs.are_invasive_changes_applied = False
            request.ParameterList.ParameterValueStruct.append(name_value)

        return AcsMsgAndTransition(request, self.done_transition)

    def state_description(self) -> str:
        return 'Setting parameter values'
