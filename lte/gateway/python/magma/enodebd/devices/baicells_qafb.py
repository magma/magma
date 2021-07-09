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
from magma.enodebd.device_config.configuration_init import build_desired_config
from magma.enodebd.device_config.enodeb_config_postprocessor import (
    EnodebConfigurationPostProcessor,
)
from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration
from magma.enodebd.devices.device_utils import EnodebDeviceName
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
    AddObjectsState,
    BaicellsSendRebootState,
    DeleteObjectsState,
    EndSessionState,
    EnodebAcsState,
    ErrorState,
    GetParametersState,
    GetRPCMethodsState,
    SendGetTransientParametersState,
    SetParameterValuesState,
    WaitEmptyMessageState,
    WaitGetParametersState,
    WaitInformMRebootState,
    WaitInformState,
    WaitRebootResponseState,
    WaitSetParameterValuesState,
)
from magma.enodebd.tr069 import models


class BaicellsQAFBHandler(BasicEnodebAcsStateMachine):
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
            'wait_get_transient_params': BaicellsQafbWaitGetTransientParametersState(self, when_get='get_params', when_get_obj_params='get_obj_params', when_delete='delete_objs', when_add='add_objs', when_set='set_params', when_skip='end_session'),
            'get_params': GetParametersState(self, when_done='wait_get_params'),
            'wait_get_params': WaitGetParametersState(self, when_done='get_obj_params'),
            'get_obj_params': BaicellsQafbGetObjectParametersState(self, when_delete='delete_objs', when_add='add_objs', when_set='set_params', when_skip='end_session'),
            'delete_objs': DeleteObjectsState(self, when_add='add_objs', when_skip='set_params'),
            'add_objs': AddObjectsState(self, when_done='set_params'),
            'set_params': SetParameterValuesState(self, when_done='wait_set_params'),
            'wait_set_params': WaitSetParameterValuesState(self, when_done='check_get_params', when_apply_invasive='check_get_params'),
            'check_get_params': GetParametersState(self, when_done='check_wait_get_params', request_all_params=True),
            'check_wait_get_params': WaitGetParametersState(self, when_done='end_session'),
            'end_session': EndSessionState(self),
            # These states are only entered through manual user intervention
            'reboot': BaicellsSendRebootState(self, when_done='wait_reboot'),
            'wait_reboot': WaitRebootResponseState(self, when_done='wait_post_reboot_inform'),
            'wait_post_reboot_inform': WaitInformMRebootState(self, when_done='wait_empty', when_timeout='wait_inform'),
            # The states below are entered when an unexpected message type is
            # received
            'unexpected_fault': ErrorState(self, inform_transition_target='wait_inform'),
        }

    @property
    def device_name(self) -> str:
        return EnodebDeviceName.BAICELLS_QAFB

    @property
    def data_model_class(self) -> Type[DataModel]:
        return BaicellsQAFBTrDataModel

    @property
    def config_postprocessor(self) -> EnodebConfigurationPostProcessor:
        return BaicellsQAFBTrConfigurationInitializer()

    @property
    def state_map(self) -> Dict[str, EnodebAcsState]:
        return self._state_map

    @property
    def disconnected_state_name(self) -> str:
        return 'wait_inform'

    @property
    def unexpected_fault_state_name(self) -> str:
        return 'unexpected_fault'


def _get_object_params_to_get(
    device_cfg: EnodebConfiguration,
    data_model: DataModel,
) -> List[ParameterName]:
    """
    Returns a list of parameter names for object parameters we don't know the
    current value of.

    Since there is no parameter for tracking the number of PLMNs, then we
    make the assumption that if any PLMN object exists, then we've already
    fetched the object parameter values.
    """
    if device_cfg.has_object(ParameterName.PLMN_N % 1):
        return []

    names = []
    num_plmns = data_model.get_num_plmns()
    obj_to_params = data_model.get_numbered_param_names()
    for i in range(1, num_plmns + 1):
        obj_name = ParameterName.PLMN_N % i
        desired = obj_to_params[obj_name]
        names = names + desired
    return names


class BaicellsQafbWaitGetTransientParametersState(EnodebAcsState):
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
        """ Process GetParameterValuesResponse """
        if not isinstance(message, models.GetParameterValuesResponse):
            return AcsReadMsgResult(False, None)
        # Current values of the fetched parameters
        name_to_val = parse_get_parameter_values_response(
            self.acs.data_model,
            message,
        )
        logger.debug('Received Parameters: %s', str(name_to_val))

        # Update device configuration
        for name in name_to_val:
            magma_value = \
                self.acs.data_model.transform_for_magma(name, name_to_val[name])
            self.acs.device_cfg.set_parameter(name, magma_value)

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
                _get_object_params_to_get(
                    self.acs.device_cfg,
                    self.acs.data_model,
                ),
            ) > 0
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


class BaicellsQafbGetObjectParametersState(EnodebAcsState):
    """
    Get information on parameters belonging to objects that can be added or
    removed from the configuration.

    Baicells QAFB will report a parameter value as None if it does not exist
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

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        """ Respond with GetParameterValuesRequest """
        names = _get_object_params_to_get(
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

        return AcsMsgAndTransition(request, None)

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        """
        Process GetParameterValuesResponse

        Object parameters that have a reported value of None indicate that
        the object is not in the eNB's configuration. Most eNB devices will
        reply with a Fault message if we try to get values of parameters that
        don't exist on the data model, so this is an idiosyncrasy of Baicells
        QAFB.
        """
        if not isinstance(message, models.GetParameterValuesResponse):
            return AcsReadMsgResult(False, None)

        path_to_val = {}
        for param_value_struct in message.ParameterList.ParameterValueStruct:
            path_to_val[param_value_struct.Name] = \
                param_value_struct.Value.Data

        logger.debug('Received object parameters: %s', str(path_to_val))

        num_plmns = self.acs.data_model.get_num_plmns()
        for i in range(1, num_plmns + 1):
            obj_name = ParameterName.PLMN_N % i
            obj_to_params = self.acs.data_model.get_numbered_param_names()
            param_name_list = obj_to_params[obj_name]
            for name in param_name_list:
                path = self.acs.data_model.get_parameter(name).path
                if path in path_to_val:
                    value = path_to_val[path]
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


class BaicellsQAFBTrDataModel(DataModel):
    """
    Class to represent relevant data model parameters from TR-196/TR-098.
    This class is effectively read-only.

    This model specifically targets Qualcomm-based BaiCells units running
    QAFB firmware.

    These models have these idiosyncrasies (on account of running TR098):

    - Parameter content root is different (InternetGatewayDevice)
    - GetParameter queries with a wildcard e.g. InternetGatewayDevice. do
      not respond with the full tree (we have to query all parameters)
    - MME status is not exposed - we assume the MME is connected if
      the eNodeB is transmitting (OpState=true)
    - Parameters such as band capability/duplex config
      are rooted under `boardconf.` and not the device config root
    - Parameters like Admin state, CellReservedForOperatorUse,
      Duplex mode, DL bandwidth and Band capability have different
      formats from Intel-based Baicells units, necessitating,
      formatting before configuration and transforming values
      read from eNodeB state.
    - Num PLMNs is not reported by these units
    """
    # Mapping of TR parameter paths to aliases
    DEVICE_PATH = 'InternetGatewayDevice.'
    FAPSERVICE_PATH = DEVICE_PATH + 'Services.FAPService.1.'
    EEPROM_PATH = 'boardconf.status.eepromInfo.'
    PARAMETERS = {
        # Top-level objects
        ParameterName.DEVICE: TrParam(DEVICE_PATH, True, TrParameterType.OBJECT, False),
        ParameterName.FAP_SERVICE: TrParam(FAPSERVICE_PATH, True, TrParameterType.OBJECT, False),

        # Qualcomm units do not expose MME_Status (We assume that the eNB is broadcasting state is connected to the MME)
        ParameterName.MME_STATUS: TrParam(FAPSERVICE_PATH + 'CellConfig.1.LTE.X_QUALCOMM_FAPControl.OpState', True, TrParameterType.BOOLEAN, False),
        ParameterName.GPS_LAT: TrParam(DEVICE_PATH + 'FAP.GPS.latitude', True, TrParameterType.STRING, False),
        ParameterName.GPS_LONG: TrParam(DEVICE_PATH + 'FAP.GPS.longitude', True, TrParameterType.STRING, False),
        ParameterName.SW_VERSION: TrParam(DEVICE_PATH + 'DeviceInfo.SoftwareVersion', True, TrParameterType.STRING, False),
        ParameterName.SERIAL_NUMBER: TrParam(DEVICE_PATH + 'DeviceInfo.SerialNumber', True, TrParameterType.STRING, False),

        # Capabilities
        ParameterName.DUPLEX_MODE_CAPABILITY: TrParam(EEPROM_PATH + 'div_multiple', True, TrParameterType.STRING, False),
        ParameterName.BAND_CAPABILITY: TrParam(EEPROM_PATH + 'work_mode', True, TrParameterType.STRING, False),

        # RF-related parameters
        ParameterName.EARFCNDL: TrParam(FAPSERVICE_PATH + 'CellConfig.1.LTE.RAN.RF.EARFCNDL', True, TrParameterType.INT, False),
        ParameterName.PCI: TrParam(FAPSERVICE_PATH + 'CellConfig.1.LTE.RAN.RF.PhyCellID', True, TrParameterType.INT, False),
        ParameterName.DL_BANDWIDTH: TrParam(DEVICE_PATH + 'Services.RfConfig.1.RfCarrierCommon.carrierBwMhz', True, TrParameterType.INT, False),
        ParameterName.SUBFRAME_ASSIGNMENT: TrParam(FAPSERVICE_PATH + 'CellConfig.1.LTE.RAN.PHY.TDDFrame.SubFrameAssignment', True, 'bool', False),
        ParameterName.SPECIAL_SUBFRAME_PATTERN: TrParam(FAPSERVICE_PATH + 'CellConfig.1.LTE.RAN.PHY.TDDFrame.SpecialSubframePatterns', True, TrParameterType.INT, False),
        ParameterName.CELL_ID: TrParam(FAPSERVICE_PATH + 'CellConfig.1.LTE.RAN.Common.CellIdentity', True, TrParameterType.UNSIGNED_INT, False),

        # Other LTE parameters
        ParameterName.ADMIN_STATE: TrParam(FAPSERVICE_PATH + 'CellConfig.1.LTE.X_QUALCOMM_FAPControl.AdminState', False, TrParameterType.STRING, False),
        ParameterName.OP_STATE: TrParam(FAPSERVICE_PATH + 'CellConfig.1.LTE.X_QUALCOMM_FAPControl.OpState', True, TrParameterType.BOOLEAN, False),
        ParameterName.RF_TX_STATUS: TrParam(FAPSERVICE_PATH + 'CellConfig.1.LTE.X_QUALCOMM_FAPControl.OpState', True, TrParameterType.BOOLEAN, False),

        # Core network parameters
        ParameterName.MME_IP: TrParam(FAPSERVICE_PATH + 'FAPControl.LTE.Gateway.S1SigLinkServerList', True, TrParameterType.STRING, False),
        ParameterName.MME_PORT: TrParam(FAPSERVICE_PATH + 'FAPControl.LTE.Gateway.S1SigLinkPort', True, TrParameterType.INT, False),
        # This parameter is standard but doesn't exist
        # ParameterName.NUM_PLMNS: TrParam(FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNListNumberOfEntries', True, TrParameterType.INT, False),
        ParameterName.TAC: TrParam(FAPSERVICE_PATH + 'CellConfig.1.LTE.EPC.TAC', True, TrParameterType.INT, False),
        ParameterName.IP_SEC_ENABLE: TrParam('boardconf.ipsec.ipsecConfig.onBoot', False, TrParameterType.BOOLEAN, False),

        # Management server parameters
        ParameterName.PERIODIC_INFORM_ENABLE: TrParam(DEVICE_PATH + 'ManagementServer.PeriodicInformEnable', False, TrParameterType.BOOLEAN, False),
        ParameterName.PERIODIC_INFORM_INTERVAL: TrParam(DEVICE_PATH + 'ManagementServer.PeriodicInformInterval', False, TrParameterType.INT, False),

        # Performance management parameters
        ParameterName.PERF_MGMT_ENABLE: TrParam(FAPSERVICE_PATH + 'CellConfig.1.X_QUALCOMM_PerfMgmt.Config.Enable', False, TrParameterType.BOOLEAN, False),
        ParameterName.PERF_MGMT_UPLOAD_INTERVAL: TrParam(DEVICE_PATH + 'FAP.PerfMgmt.Config.PeriodicUploadInterval', False, TrParameterType.INT, False),
        ParameterName.PERF_MGMT_UPLOAD_URL: TrParam(DEVICE_PATH + 'FAP.PerfMgmt.Config.URL', False, TrParameterType.STRING, False),
    }

    NUM_PLMNS_IN_CONFIG = 6
    TRANSFORMS_FOR_ENB = {
        ParameterName.CELL_BARRED: transform_for_enb.invert_cell_barred,
    }
    for i in range(1, NUM_PLMNS_IN_CONFIG + 1):
        TRANSFORMS_FOR_ENB[ParameterName.PLMN_N_CELL_RESERVED % i] = transform_for_enb.cell_reserved
        PARAMETERS[ParameterName.PLMN_N % i] = TrParam(FAPSERVICE_PATH + 'CellConfig.1.LTE.EPC.PLMNList.%d.' % i, True, TrParameterType.STRING, False)
        PARAMETERS[ParameterName.PLMN_N_CELL_RESERVED % i] = TrParam(FAPSERVICE_PATH + 'CellConfig.1.LTE.EPC.PLMNList.%d.CellReservedForOperatorUse' % i, True, TrParameterType.STRING, False)
        PARAMETERS[ParameterName.PLMN_N_ENABLE % i] = TrParam(FAPSERVICE_PATH + 'CellConfig.1.LTE.EPC.PLMNList.%d.Enable' % i, True, TrParameterType.BOOLEAN, False)
        PARAMETERS[ParameterName.PLMN_N_PRIMARY % i] = TrParam(FAPSERVICE_PATH + 'CellConfig.1.LTE.EPC.PLMNList.%d.IsPrimary' % i, True, TrParameterType.BOOLEAN, False)
        PARAMETERS[ParameterName.PLMN_N_PLMNID % i] = TrParam(FAPSERVICE_PATH + 'CellConfig.1.LTE.EPC.PLMNList.%d.PLMNID' % i, True, TrParameterType.STRING, False)

    TRANSFORMS_FOR_ENB[ParameterName.ADMIN_STATE] = transform_for_enb.admin_state
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
            params = []
            params.append(ParameterName.PLMN_N_CELL_RESERVED % i)
            params.append(ParameterName.PLMN_N_ENABLE % i)
            params.append(ParameterName.PLMN_N_PRIMARY % i)
            params.append(ParameterName.PLMN_N_PLMNID % i)
            names[ParameterName.PLMN_N % i] = params

        return names


class BaicellsQAFBTrConfigurationInitializer(EnodebConfigurationPostProcessor):
    def postprocess(self, desired_cfg: EnodebConfiguration) -> None:
        # We don't set this parameter for this device, it should be
        # auto-configured by the device.
        desired_cfg.delete_parameter(ParameterName.ADMIN_STATE)
        return
