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
from magma.enodebd.data_models.data_model import (
    DataModel,
    InvalidTrParamPath,
    TrParam,
)
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
    DeleteObjectsState,
    EnbSendDownloadState,
    EnbSendRebootState,
    EndSessionState,
    EnodebAcsState,
    ErrorState,
    GetParametersState,
    GetRPCMethodsState,
    SendFactoryResetState,
    SendGetTransientParametersState,
    SetParameterValuesState,
    WaitDownloadResponseState,
    WaitEmptyMessageState,
    WaitFactoryResetResponseState,
    WaitGetParametersState,
    WaitInformMRebootState,
    WaitInformState,
    WaitRebootResponseState,
    WaitSetParameterValuesState,
)
from magma.enodebd.tr069 import models

logger = EnodebdLogger


class BaicellsQAFBHandler(BasicEnodebAcsStateMachine):
    """
    BaicellsQAFB State Machine
    """

    def __init__(
        self,
        service: MagmaService,
    ) -> None:
        self._state_map = {}
        super().__init__(service=service, use_param_key=False)

    def reboot_asap(self) -> None:
        """
        Transition to 'reboot' state
        """
        self.transition('reboot')

    def download_asap(
        self, url: str, user_name: str, password: str, target_file_name: str, file_size: int,
        md5: str,
    ) -> None:
        """
        Transition to 'download' state
        Args:
            url:
            user_name:
            password:
            target_file_name:
            file_size:
            md5:

        Returns:

        """
        if url is not None:
            self.desired_cfg.set_parameter(ParameterName.DOWNLOAD_URL, url)
            self.desired_cfg.set_parameter(ParameterName.DOWNLOAD_USER, user_name)
            self.desired_cfg.set_parameter(ParameterName.DOWNLOAD_PASSWORD, password)
            self.desired_cfg.set_parameter(ParameterName.DOWNLOAD_FILENAME, target_file_name)
            self.desired_cfg.set_parameter(ParameterName.DOWNLOAD_FILESIZE, file_size)
            self.desired_cfg.set_parameter(ParameterName.DOWNLOAD_MD5, md5)
        self.transition('download')

    def factory_reset_asap(self) -> None:
        """
        Impl to send a request to factoryReset the eNodeB ASAP
        The eNB will factory reset from this method.
        """
        self.transition('factory_reset')

    def is_enodeb_connected(self) -> bool:
        """
        Check if enodebd has received an Inform from the enodeb

        Returns:
            bool
        """
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
            'add_objs': BaicellsQafbAddObjectsState(self, when_done='set_params'),
            'set_params': SetParameterValuesState(self, when_done='wait_set_params'),
            'wait_set_params': WaitSetParameterValuesState(self, when_done='check_get_params', when_apply_invasive='check_get_params'),
            'check_get_params': GetParametersState(self, when_done='check_wait_get_params', request_all_params=True),
            'check_wait_get_params': WaitGetParametersState(self, when_done='end_session'),
            'end_session': EndSessionState(self),
            # These states are only entered through manual user intervention
            'reboot': EnbSendRebootState(self, when_done='wait_reboot'),
            'wait_reboot': WaitRebootResponseState(self, when_done='wait_post_reboot_inform'),
            'wait_post_reboot_inform': WaitInformMRebootState(self, when_done='wait_empty', when_timeout='wait_inform'),
            'download': EnbSendDownloadState(self, when_done='wait_download'),
            'wait_download': WaitDownloadResponseState(self, when_done='wait_inform_post_download'),
            'wait_inform_post_download': WaitInformState(self, when_done='wait_empty_post_download', when_boot=None),
            'wait_empty_post_download': WaitEmptyMessageState(
                self, when_done='get_transient_params',
                when_missing='check_optional_params',
            ),
            'factory_reset': SendFactoryResetState(self, when_done='wait_factory_reset'),
            'wait_factory_reset': WaitFactoryResetResponseState(self, when_done='wait_inform'),

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


def get_object_params_to_get(
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

    names = []
    obj_to_params = data_model.get_numbered_param_names()
    num_neighbor_freq = data_model.get_num_neighbor_freq()
    for i in range(1, num_neighbor_freq + 1):
        obj_name = BaicellsParameterName.NEGIH_FREQ_LIST % i
        desired = obj_to_params[obj_name]
        names = names + desired
    num_neighbor_cell = data_model.get_num_neighbor_cell()
    for i in range(1, num_neighbor_cell + 1):
        obj_name = BaicellsParameterName.NEIGHBOR_CELL_LIST_N % i
        desired = obj_to_params[obj_name]
        names = names + desired
    if device_cfg.has_object(ParameterName.PLMN_N % 1):
        return names
    num_plmns = data_model.get_num_plmns()
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
                get_object_params_to_get(
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
        names = get_object_params_to_get(
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

        obj_to_params = self.acs.data_model.get_numbered_param_names()
        logger.info('enb obj_to_params= %s', obj_to_params)
        num_plmns = self.acs.data_model.get_num_plmns()
        for i in range(1, num_plmns + 1):
            obj_name = ParameterName.PLMN_N % i
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
        num_neighbor_freq = self.acs.data_model.get_num_neighbor_freq()
        for i in range(1, num_neighbor_freq + 1):
            obj_name = BaicellsParameterName.NEGIH_FREQ_LIST % i
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
        num_neighbor_cell = self.acs.data_model.get_num_neighbor_cell()
        for i in range(1, num_neighbor_cell + 1):
            obj_name = BaicellsParameterName.NEIGHBOR_CELL_LIST_N % i
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


class BaicellsQafbAddObjectsState(EnodebAcsState):
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
    # Parameters to query when reading eNodeB config
    LOAD_PARAMETERS = [ParameterName.DEVICE]
    # Mapping of TR parameter paths to aliases
    DEVICE_PATH = 'InternetGatewayDevice.'
    FAPSERVICE_PATH = DEVICE_PATH + 'Services.FAPService.1.'
    EEPROM_PATH = 'boardconf.status.eepromInfo.'
    PTP_PATH = 'boardconf.ptp.ptpConfig.'
    PARAMETERS = {
        # Top-level objects
        ParameterName.DEVICE: TrParam(
            path=DEVICE_PATH,
            is_invasive=True, type=TrParameterType.OBJECT, is_optional=False,
        ),
        ParameterName.FAP_SERVICE: TrParam(
            path=FAPSERVICE_PATH,
            is_invasive=True, type=TrParameterType.OBJECT, is_optional=False,
        ),

        # Device info parameters
        # Qualcomm units do not expose MME_Status (We assume that the eNB is broadcasting state is connected to the MME)
        ParameterName.MME_STATUS: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.1.LTE.X_QUALCOMM_FAPControl.OpState',
            is_invasive=True, type=TrParameterType.BOOLEAN, is_optional=False,
        ),
        ParameterName.GPS_LAT: TrParam(
            path=DEVICE_PATH + 'FAP.GPS.latitude',
            is_invasive=True, type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.GPS_LONG: TrParam(
            path=DEVICE_PATH + 'FAP.GPS.longitude',
            is_invasive=True, type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.GPS_ALTI: TrParam(
            path=DEVICE_PATH + 'FAP.GPS.altitudeWrtMeanSeaLevel',
            is_invasive=True, type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.SW_VERSION: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.SoftwareVersion',
            is_invasive=True, type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.SERIAL_NUMBER: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.SerialNumber',
            is_invasive=True, type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.VENDOR: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.ManufacturerOUI',
            is_invasive=True, type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.MODEL_NAME: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.ModuleType',
            is_invasive=True, type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.RF_STATE: TrParam(
            path=DEVICE_PATH + 'Services.RfConfig.1.RfCarrierCommon.adminState',
            is_invasive=True, type=TrParameterType.BOOLEAN, is_optional=False,
        ),
        ParameterName.UPTIME: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.X_BAICELLS_COM_STATION_RUN_Time',
            is_invasive=True, type=TrParameterType.STRING, is_optional=False,
        ),

        # Capabilities
        ParameterName.DUPLEX_MODE_CAPABILITY: TrParam(
            path=EEPROM_PATH + 'div_multiple',
            is_invasive=True, type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.BAND_CAPABILITY: TrParam(
            path=EEPROM_PATH + 'work_mode',
            is_invasive=True, type=TrParameterType.STRING, is_optional=False,
        ),

        # RF-related parameters
        ParameterName.EARFCNDL: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.1.LTE.RAN.RF.EARFCNDL',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        ParameterName.PCI: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.1.LTE.RAN.RF.PhyCellID',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        ParameterName.DL_BANDWIDTH: TrParam(
            path=DEVICE_PATH + 'Services.RfConfig.1.RfCarrierCommon.carrierBwMhz',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        ParameterName.SUBFRAME_ASSIGNMENT: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.1.LTE.RAN.PHY.TDDFrame.SubFrameAssignment',
            is_invasive=True, type='bool', is_optional=False,
        ),
        ParameterName.SPECIAL_SUBFRAME_PATTERN: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.1.LTE.RAN.PHY.TDDFrame.SpecialSubframePatterns',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        ParameterName.CELL_ID: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.1.LTE.RAN.Common.CellIdentity',
            is_invasive=True, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        ),

        # Other LTE parameters
        ParameterName.ADMIN_STATE: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.1.LTE.X_QUALCOMM_FAPControl.AdminState',
            is_invasive=False, type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.OP_STATE: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.1.LTE.X_QUALCOMM_FAPControl.OpState',
            is_invasive=True, type=TrParameterType.BOOLEAN, is_optional=False,
        ),
        ParameterName.RF_TX_STATUS: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.1.LTE.X_QUALCOMM_FAPControl.OpState',
            is_invasive=True, type=TrParameterType.BOOLEAN, is_optional=False,
        ),

        # RAN parameters
        BaicellsParameterName.X2_ENABLE_DISABLE: TrParam(
            path=FAPSERVICE_PATH + 'X_QUALCOMM_ENB_CONFIG.ULTRASON_ENB_CONFIG.UsonEnbSelfConfig.X2ConnectionEnabled',
            is_invasive=True, type=TrParameterType.BOOLEAN, is_optional=False,
        ),

        # Core network parameters
        ParameterName.MME_IP: TrParam(
            path=FAPSERVICE_PATH + 'FAPControl.LTE.Gateway.S1SigLinkServerList',
            is_invasive=True, type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.MME_PORT: TrParam(
            path=FAPSERVICE_PATH + 'FAPControl.LTE.Gateway.S1SigLinkPort',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        # This parameter is standard but doesn't exist
        # ParameterName.NUM_PLMNS: TrParam(FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNListNumberOfEntries', True, TrParameterType.INT, False),
        ParameterName.NUM_PLMNS: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNListNumberOfEntries',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        ParameterName.TAC: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.1.LTE.EPC.TAC',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        ParameterName.IP_SEC_ENABLE: TrParam(
            path='boardconf.ipsec.ipsecConfig.onBoot',
            is_invasive=False, type=TrParameterType.BOOLEAN, is_optional=False,
        ),
        BaicellsParameterName.NUM_LTE_NEIGHBOR_FREQ: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.CarrierNumberOfEntries',
            is_invasive=False, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.NUM_LTE_NEIGHBOR_CELL: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECellNumberOfEntries',
            is_invasive=False, type=TrParameterType.INT, is_optional=False,
        ),

        # Management server parameters
        ParameterName.PERIODIC_INFORM_ENABLE: TrParam(
            path=DEVICE_PATH + 'ManagementServer.PeriodicInformEnable',
            is_invasive=False, type=TrParameterType.BOOLEAN, is_optional=False,
        ),
        ParameterName.PERIODIC_INFORM_INTERVAL: TrParam(
            path=DEVICE_PATH + 'ManagementServer.PeriodicInformInterval',
            is_invasive=False, type=TrParameterType.INT, is_optional=False,
        ),

        # Performance management parameters
        ParameterName.PERF_MGMT_ENABLE: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.1.X_QUALCOMM_PerfMgmt.Config.Enable',
            is_invasive=False, type=TrParameterType.BOOLEAN, is_optional=False,
        ),
        ParameterName.PERF_MGMT_UPLOAD_INTERVAL: TrParam(
            path=DEVICE_PATH + 'FAP.PerfMgmt.Config.PeriodicUploadInterval',
            is_invasive=False, type=TrParameterType.INT, is_optional=False,
        ),
        ParameterName.PERF_MGMT_UPLOAD_URL: TrParam(
            path=DEVICE_PATH + 'FAP.PerfMgmt.Config.URL',
            is_invasive=False, type=TrParameterType.STRING, is_optional=False,
        ),

        # download params that don't have tr69 representation.
        ParameterName.DOWNLOAD_URL: TrParam(
            path=InvalidTrParamPath,
            is_invasive=False, type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.DOWNLOAD_USER: TrParam(
            path=InvalidTrParamPath,
            is_invasive=False, type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.DOWNLOAD_PASSWORD: TrParam(
            path=InvalidTrParamPath,
            is_invasive=False, type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.DOWNLOAD_FILENAME: TrParam(
            path=InvalidTrParamPath,
            is_invasive=False, type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.DOWNLOAD_FILESIZE: TrParam(
            path=InvalidTrParamPath,
            is_invasive=False, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        ),
        ParameterName.DOWNLOAD_MD5: TrParam(
            path=InvalidTrParamPath,
            is_invasive=False, type=TrParameterType.STRING, is_optional=False,
        ),

        # Radio Power config
        BaicellsParameterName.REFERENCE_SIGNAL_POWER: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.ReferenceSignalPower',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.POWER_CLASS: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.X_QUALCOMM_EXPANDED_POWER_PARAMS.MaxTxPowerExpanded',
            is_invasive=True, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        ),
        BaicellsParameterName.PA: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.PHY.PDSCH.Pa',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.PB: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.PHY.PDSCH.Pb',
            is_invasive=True, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        ),

        # Management server
        BaicellsParameterName.MANAGEMENT_SERVER: TrParam(
            path=DEVICE_PATH + 'ManagementServer.URL',
            is_invasive=False, type=TrParameterType.STRING, is_optional=False,
        ),
        BaicellsParameterName.MANAGEMENT_SERVER_PORT: TrParam(
            path=InvalidTrParamPath,
            is_invasive=False, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.MANAGEMENT_SERVER_SSL_ENABLE: TrParam(
            path=InvalidTrParamPath,
            is_invasive=False, type=TrParameterType.BOOLEAN, is_optional=False,
        ),

        # SYNC
        BaicellsParameterName.SYNC_1588_SWITCH: TrParam(
            path=PTP_PATH + 'ptpEnable',
            is_invasive=True, type=TrParameterType.BOOLEAN, is_optional=False,
        ),
        BaicellsParameterName.SYNC_1588_DOMAIN: TrParam(
            path=PTP_PATH + 'q_Domain',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.SYNC_1588_SYNC_MSG_INTREVAL: TrParam(
            path=PTP_PATH + 'q_SyncInterval',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.SYNC_1588_DELAY_REQUEST_MSG_INTERVAL: TrParam(
            path=PTP_PATH + 'q_DelayInterval',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.SYNC_1588_HOLDOVER: TrParam(
            path=InvalidTrParamPath,
            is_invasive=False, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.SYNC_1588_ASYMMETRY: TrParam(
            path=PTP_PATH + 'q_Asymmetry',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.SYNC_1588_UNICAST_ENABLE: TrParam(
            path=PTP_PATH + 'q_ModeSwitch',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.SYNC_1588_UNICAST_SERVERIP: TrParam(
            path=PTP_PATH + 'unicastAddr',
            is_invasive=True, type=TrParameterType.STRING, is_optional=False,
        ),

        # Ho algorithm
        BaicellsParameterName.HO_A1_THRESHOLD_RSRP: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.ConnMode.EUTRA.A1ThresholdRSRP',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.HO_A2_THRESHOLD_RSRP: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.ConnMode.EUTRA.A2ThresholdRSRP',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.HO_A3_OFFSET: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.ConnMode.EUTRA.A3Offset',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.HO_A3_OFFSET_ANR: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.X_QUALCOMM_ULTRASON_CONFIG.SelfConfig.AnrA3NormalThresholdInDb',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.HO_A4_THRESHOLD_RSRP: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.ConnMode.EUTRA.A4ThresholdRSRP',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.HO_LTE_INTRA_A5_THRESHOLD_1_RSRP: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.ConnMode.EUTRA.A5Threshold1RSRP',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.HO_LTE_INTRA_A5_THRESHOLD_2_RSRP: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.ConnMode.EUTRA.A5Threshold2RSRP',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.HO_LTE_INTER_ANR_A5_THRESHOLD_1_RSRP: TrParam(
            path=InvalidTrParamPath,
            is_invasive=False, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.HO_LTE_INTER_ANR_A5_THRESHOLD_2_RSRP: TrParam(
            path=InvalidTrParamPath,
            is_invasive=False, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.HO_B2_THRESHOLD1_RSRP: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.ConnMode.IRAT.B2Threshold1RSRP',
            is_invasive=True, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        ),
        BaicellsParameterName.HO_B2_THRESHOLD2_RSRP: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.ConnMode.IRAT.B2Threshold2UTRARSCP',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.HO_B2_GERAN_IRAT_THRESHOLD: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.ConnMode.IRAT.B2Threshold2GERAN',
            is_invasive=True, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        ),
        BaicellsParameterName.HO_QRXLEVMIN_SELECTION: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.IntraFreq.QRxLevMinSIB1',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.HO_QRXLEVMINOFFSET: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.IntraFreq.QRxLevMinOffset',
            is_invasive=True, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        ),
        BaicellsParameterName.HO_S_INTRASEARCH: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.IntraFreq.SIntraSearch',
            is_invasive=True, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        ),
        BaicellsParameterName.HO_S_NONINTRASEARCH: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.IntraFreq.SNonIntraSearch',
            is_invasive=True, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        ),
        BaicellsParameterName.HO_QRXLEVMIN_RESELECTION: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.IntraFreq.QRxLevMinSIB3',
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        ),
        BaicellsParameterName.HO_RESELECTION_PRIORITY: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.IntraFreq.CellReselectionPriority', True,
            TrParameterType.UNSIGNED_INT,
            False,
        ),
        BaicellsParameterName.HO_THRESHSERVINGLOW: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.IntraFreq.ThreshServingLow', True,
            TrParameterType.UNSIGNED_INT,
            False,
        ),
        BaicellsParameterName.HO_CIPHERING_ALGORITHM: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.AllowedCipheringAlgorithmList', True, TrParameterType.STRING,
            False,
        ),
        BaicellsParameterName.HO_INTEGRITY_ALGORITHM: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.AllowedIntegrityProtectionAlgorithmList', True,
            TrParameterType.STRING,
            False,
        ),
    }

    TRANSFORMS_FOR_ENB = {
        ParameterName.CELL_BARRED: transform_for_enb.invert_cell_barred,
        ParameterName.ADMIN_STATE: transform_for_enb.admin_state,
    }
    TRANSFORMS_FOR_MAGMA = {
        # We don't set these parameters
        ParameterName.BAND_CAPABILITY: transform_for_magma.band_capability,
        ParameterName.DUPLEX_MODE_CAPABILITY: transform_for_magma.duplex_mode,
    }
    NUM_PLMNS_IN_CONFIG = 6
    for i in range(1, NUM_PLMNS_IN_CONFIG + 1):
        TRANSFORMS_FOR_ENB[ParameterName.PLMN_N_CELL_RESERVED % i] = transform_for_enb.cell_reserved
        TRANSFORMS_FOR_MAGMA[ParameterName.PLMN_N_CELL_RESERVED % i] = transform_for_magma.cell_reserved
        PARAMETERS[ParameterName.PLMN_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.1.LTE.EPC.PLMNList.%d.' % i,
            is_invasive=True, type=TrParameterType.STRING, is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_CELL_RESERVED % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.1.LTE.EPC.PLMNList.%d.CellReservedForOperatorUse' % i,
            is_invasive=True, type=TrParameterType.STRING, is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_ENABLE % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.1.LTE.EPC.PLMNList.%d.Enable' % i,
            is_invasive=True, type=TrParameterType.BOOLEAN, is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_PRIMARY % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.1.LTE.EPC.PLMNList.%d.IsPrimary' % i,
            is_invasive=True, type=TrParameterType.BOOLEAN, is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_PLMNID % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.1.LTE.EPC.PLMNList.%d.PLMNID' % i,
            is_invasive=True, type=TrParameterType.STRING, is_optional=False,
        )

    NUM_NEIGHBOR_CELL_CONFIG = 16
    for i in range(1, NUM_NEIGHBOR_CELL_CONFIG + 1):
        PARAMETERS[BaicellsParameterName.NEIGHBOR_CELL_LIST_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECell.%d.' % i,
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_CELL_CELL_ID_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECell.%d.CID' % i,
            is_invasive=True, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_CELL_PLMN_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECell.%d.PLMNID' % i,
            is_invasive=True, type=TrParameterType.STRING, is_optional=False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_CELL_EARFCN_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECell.%d.EUTRACarrierARFCN' % i,
            is_invasive=True, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_CELL_PCI_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECell.%d.PhyCellID' % i,
            is_invasive=True, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_CELL_TAC_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECell.%d.X_QUALCOMM_TAC' % i,
            is_invasive=True, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_CELL_QOFFSET_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECell.%d.QOffset' % i,
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_CELL_CIO_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECell.%d.CIO' % i,
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_CELL_ENABLE_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.NeighborList.LTECell.%d.Enable' % i,
            is_invasive=True, type=TrParameterType.BOOLEAN, is_optional=False,
        )

    NUM_NEIGHBOR_FREQ_CONFIG = 8
    for i in range(1, NUM_NEIGHBOR_FREQ_CONFIG + 1):
        PARAMETERS[BaicellsParameterName.NEGIH_FREQ_LIST % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.' % i,
            is_invasive=True, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_EARFCN_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.EUTRACarrierARFCN' % i,
            is_invasive=True, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_QRXLEVMINSIB5_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.QRxLevMinSIB5' % i,
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_Q_OFFSETRANGE_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.QOffsetFreq' % i,
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_TRESELECTIONEUTRA_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.TReselectionEUTRA' % i,
            is_invasive=True, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_RESELECTIONPRIORITY_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.CellReselectionPriority' % i,
            is_invasive=True, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_RESELTHRESHHIGH_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.ThreshXHigh' % i,
            is_invasive=True, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_RESELTHRESHLOW_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.ThreshXLow' % i,
            is_invasive=True, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_PMAX_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.PMax' % i,
            is_invasive=True, type=TrParameterType.INT, is_optional=False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_TRESELECTIONEUTRASFMEDIUM_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.TReselectionEUTRASFMedium' % i,
            is_invasive=True, type=TrParameterType.UNSIGNED_INT, is_optional=False,
        )
        PARAMETERS[BaicellsParameterName.NEIGHBOR_FREQ_ENABLE_N % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Mobility.IdleMode.InterFreq.Carrier.%d.Enable' % i,
            is_invasive=True, type=TrParameterType.BOOLEAN, is_optional=False,
        )

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
    def get_num_neighbor_freq(cls) -> int:
        """ Get the neighbor freq number """
        return cls.NUM_NEIGHBOR_FREQ_CONFIG

    @classmethod
    def get_num_neighbor_cell(cls) -> int:
        """ Get the neighbor cell number """
        return cls.NUM_NEIGHBOR_CELL_CONFIG

    @classmethod
    def get_parameter_names(cls) -> List[ParameterName]:
        excluded_params = [
            str(ParameterName.DEVICE),
            str(ParameterName.FAP_SERVICE),
        ]
        names = list(
            filter(
                lambda x: (not str(x).startswith('PLMN')) and (not str(x).startswith('Download')) and (not str(x).startswith('neighbor')) and (
                    str(x) not in excluded_params
                ), cls.PARAMETERS.keys(),
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


class BaicellsQAFBTrConfigurationInitializer(EnodebConfigurationPostProcessor):
    def postprocess(self, mconfig: Any, service_cfg: Any, desired_cfg: EnodebConfiguration) -> None:
        # We don't set this parameter for this device, it should be
        # auto-configured by the device.
        desired_cfg.delete_parameter(ParameterName.ADMIN_STATE)
        desired_cfg.set_parameter(ParameterName.PERF_MGMT_UPLOAD_INTERVAL, 900)
