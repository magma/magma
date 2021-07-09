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
from magma.enodebd.exceptions import Tr069Error
from magma.enodebd.logger import EnodebdLogger as logger
from magma.enodebd.state_machines.acs_state_utils import (
    get_all_objects_to_add,
    get_all_objects_to_delete,
)
from magma.enodebd.state_machines.enb_acs import EnodebAcsStateMachine
from magma.enodebd.state_machines.enb_acs_impl import BasicEnodebAcsStateMachine
from magma.enodebd.state_machines.enb_acs_states import (
    AcsMsgAndTransition,
    AcsReadMsgResult,
    AddObjectsState,
    DeleteObjectsState,
    EndSessionState,
    EnodebAcsState,
    ErrorState,
    GetParametersState,
    GetRPCMethodsState,
    SendGetTransientParametersState,
    SendRebootState,
    SetParameterValuesNotAdminState,
    WaitEmptyMessageState,
    WaitGetObjectParametersState,
    WaitGetParametersState,
    WaitGetTransientParametersState,
    WaitInformMRebootState,
    WaitInformState,
    WaitRebootResponseState,
    WaitSetParameterValuesState,
)
from magma.enodebd.tr069 import models


class CaviumHandler(BasicEnodebAcsStateMachine):
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
            'wait_inform': WaitInformState(self, when_done='get_rpc_methods'),
            'get_rpc_methods': GetRPCMethodsState(self, when_done='wait_empty', when_skip='get_transient_params'),
            'wait_empty': WaitEmptyMessageState(self, when_done='get_transient_params'),
            'get_transient_params': SendGetTransientParametersState(self, when_done='wait_get_transient_params'),
            'wait_get_transient_params': WaitGetTransientParametersState(self, when_get='get_params', when_get_obj_params='get_obj_params', when_delete='delete_objs', when_add='add_objs', when_set='set_params', when_skip='end_session'),
            'get_params': GetParametersState(self, when_done='wait_get_params'),
            'wait_get_params': WaitGetParametersState(self, when_done='get_obj_params'),
            'get_obj_params': CaviumGetObjectParametersState(self, when_done='wait_get_obj_params'),
            'wait_get_obj_params': CaviumWaitGetObjectParametersState(self, when_edit='disable_admin', when_skip='get_transient_params'),
            'disable_admin': CaviumDisableAdminEnableState(self, admin_value=False, when_done='wait_disable_admin'),
            'wait_disable_admin': CaviumWaitDisableAdminEnableState(self, admin_value=False, when_add='add_objs', when_delete='delete_objs', when_done='set_params'),
            'delete_objs': DeleteObjectsState(self, when_add='add_objs', when_skip='set_params'),
            'add_objs': AddObjectsState(self, when_done='set_params'),
            'set_params': SetParameterValuesNotAdminState(self, when_done='wait_set_params'),
            'wait_set_params': WaitSetParameterValuesState(self, when_done='enable_admin', when_apply_invasive='enable_admin'),
            'enable_admin': CaviumDisableAdminEnableState(self, admin_value=True, when_done='wait_enable_admin'),
            'wait_enable_admin': CaviumWaitDisableAdminEnableState(self, admin_value=True, when_done='check_get_params', when_add='check_get_params', when_delete='check_get_params'),
            'check_get_params': GetParametersState(self, when_done='check_wait_get_params', request_all_params=True),
            'check_wait_get_params': WaitGetParametersState(self, when_done='end_session'),
            'end_session': EndSessionState(self),
            # Below states only entered through manual user intervention
            'reboot': SendRebootState(self, when_done='wait_reboot'),
            'wait_reboot': WaitRebootResponseState(self, when_done='wait_post_reboot_inform'),
            'wait_post_reboot_inform': WaitInformMRebootState(self, when_done='wait_reboot_delay', when_timeout='wait_inform'),
            # The states below are entered when an unexpected message type is
            # received
            'unexpected_fault': ErrorState(self, inform_transition_target='wait_inform'),
        }

    @property
    def device_name(self) -> str:
        return EnodebDeviceName.CAVIUM

    @property
    def data_model_class(self) -> Type[DataModel]:
        return CaviumTrDataModel

    @property
    def config_postprocessor(self) -> EnodebConfigurationPostProcessor:
        return CaviumTrConfigurationInitializer()

    @property
    def state_map(self) -> Dict[str, EnodebAcsState]:
        return self._state_map

    @property
    def disconnected_state_name(self) -> str:
        return 'wait_inform'

    @property
    def unexpected_fault_state_name(self) -> str:
        return 'unexpected_fault'


class CaviumGetObjectParametersState(EnodebAcsState):
    """
    When booted, the PLMN list is empty so we cannot get individual
    object parameters. Instead, get the parent object PLMN_LIST
    which will include any children if they exist.
    """

    def __init__(self, acs: EnodebAcsStateMachine, when_done: str):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        """ Respond with GetParameterValuesRequest """
        names = [ParameterName.PLMN_LIST]

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


class CaviumWaitGetObjectParametersState(WaitGetObjectParametersState):
    def __init__(
        self,
        acs: EnodebAcsStateMachine,
        when_edit: str,
        when_skip: str,
    ):
        super().__init__(
            acs=acs,
            when_add=when_edit,
            when_delete=when_edit,
            when_set=when_edit,
            when_skip=when_skip,
        )


class CaviumDisableAdminEnableState(EnodebAcsState):
    """
    Cavium requires that we disable 'Admin Enable' before configuring
    most parameters
    """

    def __init__(self, acs: EnodebAcsStateMachine, admin_value: bool, when_done: str):
        super().__init__()
        self.acs = acs
        self.admin_value = admin_value
        self.done_transition = when_done

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        if not isinstance(message, models.DummyInput):
            return AcsReadMsgResult(False, None)
        return AcsReadMsgResult(True, None)

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        """
        Returns:
            A SetParameterValueRequest for setting 'Admin Enable' to False
        """
        param_name = ParameterName.ADMIN_STATE
        # if we want the cell to be down don't force it up
        desired_admin_value = \
                self.acs.desired_cfg.get_parameter(param_name) \
                and self.admin_value
        admin_value = \
                self.acs.data_model.transform_for_enb(
                    param_name,
                    desired_admin_value,
                )
        admin_path = self.acs.data_model.get_parameter(param_name).path
        param_values = {admin_path: admin_value}

        request = models.SetParameterValues()
        request.ParameterList = models.ParameterValueList()
        request.ParameterList.arrayType = 'cwmp:ParameterValueStruct[%d]' \
                                          % len(param_values)

        name_value = models.ParameterValueStruct()
        name_value.Name = admin_path
        name_value.Value = models.anySimpleType()
        name_value.Value.type = 'xsd:string'
        name_value.Value.Data = str(admin_value)
        request.ParameterList.ParameterValueStruct = [name_value]

        return AcsMsgAndTransition(request, self.done_transition)

    def state_description(self) -> str:
        return 'Disabling admin_enable (Cavium only)'


class CaviumWaitDisableAdminEnableState(EnodebAcsState):
    def __init__(
            self,
            acs: EnodebAcsStateMachine,
            admin_value: bool,
            when_done: str,
            when_add: str,
            when_delete: str,
    ):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done
        self.add_obj_transition = when_add
        self.del_obj_transition = when_delete
        self.admin_value = admin_value

    def read_msg(self, message: Any) -> Optional[str]:
        if type(message) == models.Fault:
            logger.error('Received Fault in response to SetParameterValues')
            if message.SetParameterValuesFault is not None:
                for fault in message.SetParameterValuesFault:
                    logger.error(
                        'SetParameterValuesFault Param: %s, Code: %s, String: %s',
                        fault.ParameterName, fault.FaultCode, fault.FaultString,
                    )
            raise Tr069Error(
                'Received Fault in response to SetParameterValues '
                '(faultstring = %s)' % message.FaultString,
            )
        elif not isinstance(message, models.SetParameterValuesResponse):
            return AcsReadMsgResult(False, None)
        if message.Status != 0:
            raise Tr069Error(
                'Received SetParameterValuesResponse with '
                'Status=%d' % message.Status,
            )
        param_name = ParameterName.ADMIN_STATE
        desired_admin_value = \
                self.acs.desired_cfg.get_parameter(param_name) \
                and self.admin_value
        magma_value = \
                self.acs.data_model.transform_for_magma(
                    param_name,
                    desired_admin_value,
                )
        self.acs.device_cfg.set_parameter(param_name, magma_value)

        if len(
            get_all_objects_to_delete(
                self.acs.desired_cfg,
                self.acs.device_cfg,
            ),
        ) > 0:
            return AcsReadMsgResult(True, self.del_obj_transition)
        elif len(
            get_all_objects_to_add(
                self.acs.desired_cfg,
                self.acs.device_cfg,
            ),
        ) > 0:
            return AcsReadMsgResult(True, self.add_obj_transition)
        else:
            return AcsReadMsgResult(True, self.done_transition)

    def state_description(self) -> str:
        return 'Disabling admin_enable (Cavium only)'


class CaviumTrDataModel(DataModel):
    """
    Class to represent relevant data model parameters from TR-196/TR-098/TR-181.
    This class is effectively read-only
    """
    # Mapping of TR parameter paths to aliases
    DEVICE_PATH = 'Device.'
    FAPSERVICE_PATH = DEVICE_PATH + 'Services.FAPService.1.'
    PARAMETERS = {
        # Top-level objects
        ParameterName.DEVICE: TrParam(DEVICE_PATH, True, TrParameterType.OBJECT, False),
        ParameterName.FAP_SERVICE: TrParam(FAPSERVICE_PATH, True, TrParameterType.OBJECT, False),

        # Device info parameters
        ParameterName.GPS_STATUS: TrParam(DEVICE_PATH + 'FAP.GPS.ContinuousGPSStatus.GotFix', True, TrParameterType.BOOLEAN, False),
        ParameterName.GPS_LAT: TrParam(DEVICE_PATH + 'FAP.GPS.LockedLatitude', True, TrParameterType.INT, False),
        ParameterName.GPS_LONG: TrParam(DEVICE_PATH + 'FAP.GPS.LockedLongitude', True, TrParameterType.INT, False),
        ParameterName.SW_VERSION: TrParam(DEVICE_PATH + 'DeviceInfo.SoftwareVersion', True, TrParameterType.STRING, False),
        ParameterName.SERIAL_NUMBER: TrParam(DEVICE_PATH + 'DeviceInfo.SerialNumber', True, TrParameterType.STRING, False),

        # Capabilities
        ParameterName.DUPLEX_MODE_CAPABILITY: TrParam(
            FAPSERVICE_PATH + 'Capabilities.LTE.DuplexMode', True, TrParameterType.STRING, False,
        ),
        ParameterName.BAND_CAPABILITY: TrParam(FAPSERVICE_PATH + 'Capabilities.LTE.BandsSupported', True, TrParameterType.UNSIGNED_INT, False),

        # RF-related parameters
        ParameterName.EARFCNDL: TrParam(FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.EARFCNDL', True, TrParameterType.UNSIGNED_INT, False),
        ParameterName.EARFCNUL: TrParam(FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.EARFCNUL', True, TrParameterType.UNSIGNED_INT, False),
        ParameterName.BAND: TrParam(FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.FreqBandIndicator', True, TrParameterType.UNSIGNED_INT, False),
        ParameterName.PCI: TrParam(FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.PhyCellID', True, TrParameterType.STRING, False),
        ParameterName.DL_BANDWIDTH: TrParam(FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.DLBandwidth', True, TrParameterType.STRING, False),
        ParameterName.UL_BANDWIDTH: TrParam(FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.ULBandwidth', True, TrParameterType.STRING, False),
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
        ParameterName.MME_PORT: TrParam(FAPSERVICE_PATH + 'FAPControl.LTE.Gateway.S1SigLinkPort', True, TrParameterType.UNSIGNED_INT, False),
        ParameterName.NUM_PLMNS: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNListNumberOfEntries', True, TrParameterType.UNSIGNED_INT, False,
        ),
        ParameterName.PLMN: TrParam(FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.', True, TrParameterType.OBJECT, False),
        # PLMN arrays are added below
        ParameterName.TAC: TrParam(FAPSERVICE_PATH + 'CellConfig.LTE.EPC.TAC', True, TrParameterType.UNSIGNED_INT, False),
        ParameterName.IP_SEC_ENABLE: TrParam(
            DEVICE_PATH + 'IPsec.Enable', False, TrParameterType.BOOLEAN, False,
        ),
        ParameterName.PERIODIC_INFORM_INTERVAL:
            TrParam(DEVICE_PATH + 'ManagementServer.PeriodicInformInterval', False, TrParameterType.UNSIGNED_INT, False),

        # Management server parameters
        ParameterName.PERIODIC_INFORM_ENABLE: TrParam(
                DEVICE_PATH + 'ManagementServer.PeriodicInformEnable',
                False, TrParameterType.BOOLEAN, False,
        ),
        ParameterName.PERIODIC_INFORM_INTERVAL: TrParam(
                DEVICE_PATH + 'ManagementServer.PeriodicInformInterval',
                False, TrParameterType.UNSIGNED_INT, False,
        ),

        # Performance management parameters
        ParameterName.PERF_MGMT_ENABLE: TrParam(
            FAPSERVICE_PATH + 'PerfMgmt.Config.1.Enable', False, TrParameterType.BOOLEAN, False,
        ),
        ParameterName.PERF_MGMT_UPLOAD_INTERVAL: TrParam(
            FAPSERVICE_PATH + 'PerfMgmt.Config.1.PeriodicUploadInterval', False, TrParameterType.UNSIGNED_INT, False,
        ),
        ParameterName.PERF_MGMT_UPLOAD_URL: TrParam(
            FAPSERVICE_PATH + 'PerfMgmt.Config.1.URL', False, TrParameterType.STRING, False,
        ),
        ParameterName.PERF_MGMT_USER: TrParam(
            FAPSERVICE_PATH + 'PerfMgmt.Config.1.Username',
            False, TrParameterType.STRING, False,
        ),
        ParameterName.PERF_MGMT_PASSWORD: TrParam(
            FAPSERVICE_PATH + 'PerfMgmt.Config.1.Password',
            False, TrParameterType.STRING, False,
        ),

        # PLMN Info
        ParameterName.PLMN_LIST: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.', False, TrParameterType.OBJECT, False,
        ),
    }

    NUM_PLMNS_IN_CONFIG = 6
    for i in range(1, NUM_PLMNS_IN_CONFIG + 1):
        PARAMETERS[ParameterName.PLMN_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.' % i, True, TrParameterType.OBJECT, False,
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
        """
        Load all the parameters instead of a subset.
        """
        return [ParameterName.DEVICE]

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
            params = []
            params.append(ParameterName.PLMN_N_CELL_RESERVED % i)
            params.append(ParameterName.PLMN_N_ENABLE % i)
            params.append(ParameterName.PLMN_N_PRIMARY % i)
            params.append(ParameterName.PLMN_N_PLMNID % i)
            names[ParameterName.PLMN_N % i] = params
        return names


class CaviumTrConfigurationInitializer(EnodebConfigurationPostProcessor):
    def postprocess(self, desired_cfg: EnodebConfiguration) -> None:
        desired_cfg.set_parameter(ParameterName.CELL_BARRED, True)
        desired_cfg.set_parameter(ParameterName.ADMIN_STATE, True)
