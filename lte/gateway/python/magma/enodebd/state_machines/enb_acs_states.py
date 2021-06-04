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

from abc import ABC, abstractmethod
from collections import namedtuple
from typing import Any, Optional

from magma.enodebd.data_models.data_model_parameters import ParameterName
from magma.enodebd.device_config.configuration_init import build_desired_config
from magma.enodebd.exceptions import ConfigurationError, Tr069Error
from magma.enodebd.logger import EnodebdLogger as logger
from magma.enodebd.state_machines.acs_state_utils import (
    does_inform_have_event,
    get_all_objects_to_add,
    get_all_objects_to_delete,
    get_all_param_values_to_set,
    get_obj_param_values_to_set,
    get_object_params_to_get,
    get_optional_param_to_check,
    get_param_values_to_set,
    get_params_to_get,
    parse_get_parameter_values_response,
    process_inform_message,
)
from magma.enodebd.state_machines.enb_acs import EnodebAcsStateMachine
from magma.enodebd.state_machines.timer import StateMachineTimer
from magma.enodebd.tr069 import models

AcsMsgAndTransition = namedtuple(
    'AcsMsgAndTransition', ['msg', 'next_state'],
)

AcsReadMsgResult = namedtuple(
    'AcsReadMsgResult', ['msg_handled', 'next_state'],
)


class EnodebAcsState(ABC):
    """
    State class for the Enodeb state machine

    States can transition after reading a message from the eNB, sending a
    message out to the eNB, or when a timer completes. As such, some states
    are only responsible for message sending, and others are only responsible
    for reading incoming messages.

    In the constructor, set up state transitions.
    """

    def __init__(self):
        self._acs = None

    def enter(self) -> None:
        """
        Set up your timers here. Call transition(..) on the ACS when the timer
        completes or throw an error
        """
        pass

    def exit(self) -> None:
        """Destroy timers here"""
        pass

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        """
        Args: message: tr069 message
        Returns: name of the next state, if transition required
        """
        raise ConfigurationError(
            '%s should implement read_msg() if it '
            'needs to handle message reading' % self.__class__.__name__,
        )

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        """
        Produce a message to send back to the eNB.

        Args:
            message: TR-069 message which was already processed by read_msg

        Returns: Message and possible transition
        """
        raise ConfigurationError(
            '%s should implement get_msg() if it '
            'needs to produce messages' % self.__class__.__name__,
        )

    @property
    def acs(self) -> EnodebAcsStateMachine:
        return self._acs

    @acs.setter
    def acs(self, val: EnodebAcsStateMachine) -> None:
        self._acs = val

    @abstractmethod
    def state_description(self) -> str:
        """ Provide a few words about what the state represents """
        pass


class WaitInformState(EnodebAcsState):
    """
    This state indicates that no Inform message has been received yet, or
    that no Inform message has been received for a long time.

    This state is used to handle an Inform message that arrived when enodebd
    already believes that the eNB is connected. As such, it is unclear to
    enodebd whether the eNB is just sending another Inform, or if a different
    eNB was plugged into the same interface.
    """

    def __init__(
        self,
        acs: EnodebAcsStateMachine,
        when_done: str,
        when_boot: Optional[str] = None,
    ):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done
        self.boot_transition = when_boot
        self.has_enb_just_booted = False

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        """
        Args:
            message: models.Inform Tr069 Inform message
        """
        if not isinstance(message, models.Inform):
            return AcsReadMsgResult(False, None)
        process_inform_message(
            message, self.acs.data_model,
            self.acs.device_cfg,
        )
        if does_inform_have_event(message, '1 BOOT'):
            return AcsReadMsgResult(True, self.boot_transition)
        return AcsReadMsgResult(True, None)

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        """ Reply with InformResponse """
        response = models.InformResponse()
        # Set maxEnvelopes to 1, as per TR-069 spec
        response.MaxEnvelopes = 1
        return AcsMsgAndTransition(response, self.done_transition)

    def state_description(self) -> str:
        return 'Waiting for an Inform'


class GetRPCMethodsState(EnodebAcsState):
    """
    After the first Inform message from boot, it is expected that the eNB
    will try to learn the RPC methods of the ACS.
    """

    def __init__(self, acs: EnodebAcsStateMachine, when_done: str, when_skip: str):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done
        self.skip_transition = when_skip

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        # If this is a regular Inform, not after a reboot we'll get an empty
        if isinstance(message, models.DummyInput):
            return AcsReadMsgResult(True, self.skip_transition)
        if not isinstance(message, models.GetRPCMethods):
            return AcsReadMsgResult(False, self.done_transition)
        return AcsReadMsgResult(True, None)

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        resp = models.GetRPCMethodsResponse()
        resp.MethodList = models.MethodList()
        RPC_METHODS = ['Inform', 'GetRPCMethods', 'TransferComplete']
        resp.MethodList.arrayType = 'xsd:string[%d]' \
                                          % len(RPC_METHODS)
        resp.MethodList.string = RPC_METHODS
        return AcsMsgAndTransition(resp, self.done_transition)

    def state_description(self) -> str:
        return 'Waiting for incoming GetRPC Methods after boot'


class BaicellsRemWaitState(EnodebAcsState):
    """
    We've already received an Inform message. This state is to handle a
    Baicells eNodeB issue.

    After eNodeB is rebooted, hold off configuring it for some time to give
    time for REM to run. This is a BaiCells eNodeB issue that doesn't support
    enabling the eNodeB during initial REM.

    In this state, just hang at responding to Inform, and then ending the
    TR-069 session.
    """

    CONFIG_DELAY_AFTER_BOOT = 600

    def __init__(self, acs: EnodebAcsStateMachine, when_done: str):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done
        self.rem_timer = None

    def enter(self):
        self.rem_timer = StateMachineTimer(self.CONFIG_DELAY_AFTER_BOOT)
        logger.info(
            'Holding off of eNB configuration for %s seconds. '
            'Will resume after eNB REM process has finished. ',
            self.CONFIG_DELAY_AFTER_BOOT,
        )

    def exit(self):
        self.rem_timer = None

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        if not isinstance(message, models.Inform):
            return AcsReadMsgResult(False, None)
        process_inform_message(
            message, self.acs.data_model,
            self.acs.device_cfg,
        )
        return AcsReadMsgResult(True, None)

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        if self.rem_timer.is_done():
            return AcsMsgAndTransition(
                models.DummyInput(),
                self.done_transition,
            )
        return AcsMsgAndTransition(models.DummyInput(), None)

    def state_description(self) -> str:
        remaining = self.rem_timer.seconds_remaining()
        return 'Waiting for eNB REM to run for %d more seconds before ' \
               'resuming with configuration.' % remaining


class WaitEmptyMessageState(EnodebAcsState):
    def __init__(
        self,
        acs: EnodebAcsStateMachine,
        when_done: str,
        when_missing: Optional[str] = None,
    ):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done
        self.unknown_param_transition = when_missing

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        """
        It's expected that we transition into this state right after receiving
        an Inform message and replying with an InformResponse. At that point,
        the eNB sends an empty HTTP request (aka DummyInput) to initiate the
        rest of the provisioning process
        """
        if not isinstance(message, models.DummyInput):
            return AcsReadMsgResult(False, None)
        if get_optional_param_to_check(self.acs.data_model) is None:
            return AcsReadMsgResult(True, self.done_transition)
        return AcsReadMsgResult(True, self.unknown_param_transition)

    def state_description(self) -> str:
        return 'Waiting for empty message from eNodeB'


class CheckOptionalParamsState(EnodebAcsState):
    def __init__(
            self,
            acs: EnodebAcsStateMachine,
            when_done: str,
    ):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done
        self.optional_param = None

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        self.optional_param = get_optional_param_to_check(self.acs.data_model)
        if self.optional_param is None:
            raise Tr069Error('Invalid State')
        # Generate the request
        request = models.GetParameterValues()
        request.ParameterNames = models.ParameterNames()
        request.ParameterNames.arrayType = 'xsd:string[1]'
        request.ParameterNames.string = []
        path = self.acs.data_model.get_parameter(self.optional_param).path
        request.ParameterNames.string.append(path)
        return AcsMsgAndTransition(request, None)

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        """ Process either GetParameterValuesResponse or a Fault """
        if type(message) == models.Fault:
            self.acs.data_model.set_parameter_presence(
                self.optional_param,
                False,
            )
        elif type(message) == models.GetParameterValuesResponse:
            name_to_val = parse_get_parameter_values_response(
                self.acs.data_model,
                message,
            )
            logger.debug(
                'Received CPE parameter values: %s',
                str(name_to_val),
            )
            for name, val in name_to_val.items():
                self.acs.data_model.set_parameter_presence(
                    self.optional_param,
                    True,
                )
                magma_val = self.acs.data_model.transform_for_magma(name, val)
                self.acs.device_cfg.set_parameter(name, magma_val)
        else:
            return AcsReadMsgResult(False, None)

        if get_optional_param_to_check(self.acs.data_model) is not None:
            return AcsReadMsgResult(True, None)
        return AcsReadMsgResult(True, self.done_transition)

    def state_description(self) -> str:
        return 'Checking if some optional parameters exist in data model'


class SendGetTransientParametersState(EnodebAcsState):
    """
    Periodically read eNodeB status. Note: keep frequency low to avoid
    backing up large numbers of read operations if enodebd is busy.
    Some eNB parameters are read only and updated by the eNB itself.
    """
    PARAMETERS = [
        ParameterName.OP_STATE,
        ParameterName.RF_TX_STATUS,
        ParameterName.GPS_STATUS,
        ParameterName.PTP_STATUS,
        ParameterName.MME_STATUS,
        ParameterName.GPS_LAT,
        ParameterName.GPS_LONG,
    ]

    def __init__(self, acs: EnodebAcsStateMachine, when_done: str):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        if not isinstance(message, models.DummyInput):
            return AcsReadMsgResult(False, None)
        return AcsReadMsgResult(True, None)

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        request = models.GetParameterValues()
        request.ParameterNames = models.ParameterNames()
        request.ParameterNames.string = []
        for name in self.PARAMETERS:
            # Not all data models have these parameters
            if self.acs.data_model.is_parameter_present(name):
                path = self.acs.data_model.get_parameter(name).path
                request.ParameterNames.string.append(path)
        request.ParameterNames.arrayType = \
            'xsd:string[%d]' % len(request.ParameterNames.string)

        return AcsMsgAndTransition(request, self.done_transition)

    def state_description(self) -> str:
        return 'Getting transient read-only parameters'


class WaitGetTransientParametersState(EnodebAcsState):
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


class GetParametersState(EnodebAcsState):
    """
    Get the value of most parameters of the eNB that are defined in the data
    model. Object parameters are excluded.
    """

    def __init__(
        self,
        acs: EnodebAcsStateMachine,
        when_done: str,
        request_all_params: bool = False,
    ):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done
        # Set to True if we want to request values of all parameters, even if
        # the ACS state machine already has recorded values of them.
        self.request_all_params = request_all_params

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        """
        It's expected that we transition into this state right after receiving
        an Inform message and replying with an InformResponse. At that point,
        the eNB sends an empty HTTP request (aka DummyInput) to initiate the
        rest of the provisioning process
        """
        if not isinstance(message, models.DummyInput):
            return AcsReadMsgResult(False, None)
        return AcsReadMsgResult(True, None)

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        """
        Respond with GetParameterValuesRequest

        Get the values of all parameters defined in the data model.
        Also check which addable objects are present, and what the values of
        parameters for those objects are.
        """

        # Get the names of regular parameters
        names = get_params_to_get(
            self.acs.device_cfg, self.acs.data_model,
            self.request_all_params,
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
        return 'Getting non-object parameters'


class WaitGetParametersState(EnodebAcsState):
    def __init__(self, acs: EnodebAcsStateMachine, when_done: str):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        """ Process GetParameterValuesResponse """
        if not isinstance(message, models.GetParameterValuesResponse):
            return AcsReadMsgResult(False, None)
        name_to_val = parse_get_parameter_values_response(
            self.acs.data_model,
            message,
        )
        logger.debug('Received CPE parameter values: %s', str(name_to_val))
        for name, val in name_to_val.items():
            magma_val = self.acs.data_model.transform_for_magma(name, val)
            self.acs.device_cfg.set_parameter(name, magma_val)
        return AcsReadMsgResult(True, self.done_transition)

    def state_description(self) -> str:
        return 'Getting non-object parameters'


class GetObjectParametersState(EnodebAcsState):
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


class WaitGetObjectParametersState(EnodebAcsState):
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
        if not isinstance(message, models.GetParameterValuesResponse):
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

        # Now we can have the desired state
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


class DeleteObjectsState(EnodebAcsState):
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


class AddObjectsState(EnodebAcsState):
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
        self.acs.device_cfg.add_object(self.added_param % instance_n)
        obj_list_to_add = get_all_objects_to_add(
            self.acs.desired_cfg,
            self.acs.device_cfg,
        )
        if len(obj_list_to_add) > 0:
            return AcsReadMsgResult(True, None)
        return AcsReadMsgResult(True, self.done_transition)

    def state_description(self) -> str:
        return 'Adding objects'


class SetParameterValuesState(EnodebAcsState):
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
        )
        request.ParameterList.arrayType = 'cwmp:ParameterValueStruct[%d]' \
                                          % len(param_values)
        request.ParameterList.ParameterValueStruct = []
        logger.debug(
            'Sending TR069 request to set CPE parameter values: %s',
            str(param_values),
        )
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


class SetParameterValuesNotAdminState(EnodebAcsState):
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
        return 'Setting parameter values excluding Admin Enable'


class WaitSetParameterValuesState(EnodebAcsState):
    def __init__(
        self,
        acs: EnodebAcsStateMachine,
        when_done: str,
        when_apply_invasive: str,
    ):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done
        self.apply_invasive_transition = when_apply_invasive

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        if type(message) == models.SetParameterValuesResponse:
            if message.Status != 0:
                raise Tr069Error(
                    'Received SetParameterValuesResponse with '
                    'Status=%d' % message.Status,
                )
            self._mark_as_configured()
            if not self.acs.are_invasive_changes_applied:
                return AcsReadMsgResult(True, self.apply_invasive_transition)
            return AcsReadMsgResult(True, self.done_transition)
        elif type(message) == models.Fault:
            logger.error(
                'Received Fault in response to SetParameterValues, '
                'Code (%s), Message (%s)', message.FaultCode,
                message.FaultString,
            )
            if message.SetParameterValuesFault is not None:
                for fault in message.SetParameterValuesFault:
                    logger.error(
                        'SetParameterValuesFault Param: %s, '
                        'Code: %s, String: %s', fault.ParameterName,
                        fault.FaultCode, fault.FaultString,
                    )
        return AcsReadMsgResult(False, None)

    def _mark_as_configured(self) -> None:
        """
        A successful attempt at setting parameter values means that we need to
        update what we think the eNB's configuration is to match what we just
        set the parameter values to.
        """
        # Values of parameters
        name_to_val = get_param_values_to_set(
            self.acs.desired_cfg,
            self.acs.device_cfg,
            self.acs.data_model,
        )
        for name, val in name_to_val.items():
            magma_val = self.acs.data_model.transform_for_magma(name, val)
            self.acs.device_cfg.set_parameter(name, magma_val)

        # Values of object parameters
        obj_to_name_to_val = get_obj_param_values_to_set(
            self.acs.desired_cfg,
            self.acs.device_cfg,
            self.acs.data_model,
        )
        for obj_name, name_to_val in obj_to_name_to_val.items():
            for name, val in name_to_val.items():
                logger.debug(
                    'Set obj: %s, name: %s, val: %s', str(obj_name),
                    str(name), str(val),
                )
                magma_val = self.acs.data_model.transform_for_magma(name, val)
                self.acs.device_cfg.set_parameter_for_object(
                    name, magma_val,
                    obj_name,
                )
        logger.info('Successfully configured CPE parameters!')

    def state_description(self) -> str:
        return 'Setting parameter values'


class EndSessionState(EnodebAcsState):
    """ To end a TR-069 session, send an empty HTTP response """

    def __init__(self, acs: EnodebAcsStateMachine):
        super().__init__()
        self.acs = acs

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        """
        No message is expected after enodebd sends the eNodeB
        an empty HTTP response.

        If a device sends an empty HTTP request, we can just
        ignore it and send another empty response.
        """
        if isinstance(message, models.DummyInput):
            return AcsReadMsgResult(True, None)
        return AcsReadMsgResult(False, None)

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        request = models.DummyInput()
        return AcsMsgAndTransition(request, None)

    def state_description(self) -> str:
        return 'Completed provisioning eNB. Awaiting new Inform.'


class BaicellsSendRebootState(EnodebAcsState):
    def __init__(self, acs: EnodebAcsStateMachine, when_done: str):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done
        self.prev_msg_was_inform = False

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        """
        This state can be transitioned into through user command.
        All messages received by enodebd will be ignored in this state.
        """
        if self.prev_msg_was_inform \
                and not isinstance(message, models.DummyInput):
            return AcsReadMsgResult(False, None)
        elif isinstance(message, models.Inform):
            self.prev_msg_was_inform = True
            process_inform_message(
                message, self.acs.data_model,
                self.acs.device_cfg,
            )
            return AcsReadMsgResult(True, None)
        self.prev_msg_was_inform = False
        return AcsReadMsgResult(True, None)

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        if self.prev_msg_was_inform:
            response = models.InformResponse()
            # Set maxEnvelopes to 1, as per TR-069 spec
            response.MaxEnvelopes = 1
            return AcsMsgAndTransition(response, None)
        logger.info('Sending reboot request to eNB')
        request = models.Reboot()
        request.CommandKey = ''
        self.acs.are_invasive_changes_applied = True
        return AcsMsgAndTransition(request, self.done_transition)

    def state_description(self) -> str:
        return 'Rebooting eNB'


class SendRebootState(EnodebAcsState):
    def __init__(self, acs: EnodebAcsStateMachine, when_done: str):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done
        self.prev_msg_was_inform = False

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        """
        This state can be transitioned into through user command.
        All messages received by enodebd will be ignored in this state.
        """
        if self.prev_msg_was_inform \
                and not isinstance(message, models.DummyInput):
            return AcsReadMsgResult(False, None)
        elif isinstance(message, models.Inform):
            self.prev_msg_was_inform = True
            process_inform_message(
                message, self.acs.data_model,
                self.acs.device_cfg,
            )
            return AcsReadMsgResult(True, None)
        self.prev_msg_was_inform = False
        return AcsReadMsgResult(True, None)

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        if self.prev_msg_was_inform:
            response = models.InformResponse()
            # Set maxEnvelopes to 1, as per TR-069 spec
            response.MaxEnvelopes = 1
            return AcsMsgAndTransition(response, None)
        logger.info('Sending reboot request to eNB')
        request = models.Reboot()
        request.CommandKey = ''
        return AcsMsgAndTransition(request, self.done_transition)

    def state_description(self) -> str:
        return 'Rebooting eNB'


class WaitRebootResponseState(EnodebAcsState):
    def __init__(self, acs: EnodebAcsStateMachine, when_done: str):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        if not isinstance(message, models.RebootResponse):
            return AcsReadMsgResult(False, None)
        return AcsReadMsgResult(True, None)

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        """ Reply with empty message """
        return AcsMsgAndTransition(models.DummyInput(), self.done_transition)

    def state_description(self) -> str:
        return 'Rebooting eNB'


class WaitInformMRebootState(EnodebAcsState):
    """
    After sending a reboot request, we expect an Inform request with a
    specific 'inform event code'
    """

    # Time to wait for eNodeB reboot. The measured time
    # (on BaiCells indoor eNodeB)
    # is ~110secs, so add healthy padding on top of this.
    REBOOT_TIMEOUT = 300  # In seconds
    # We expect that the Inform we receive tells us the eNB has rebooted
    INFORM_EVENT_CODE = 'M Reboot'

    def __init__(
        self,
        acs: EnodebAcsStateMachine,
        when_done: str,
        when_timeout: str,
    ):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done
        self.timeout_transition = when_timeout
        self.timeout_timer = None
        self.timer_handle = None

    def enter(self):
        self.timeout_timer = StateMachineTimer(self.REBOOT_TIMEOUT)

        def check_timer() -> None:
            if self.timeout_timer.is_done():
                self.acs.transition(self.timeout_transition)
                raise Tr069Error(
                    'Did not receive Inform response after '
                    'rebooting',
                )

        self.timer_handle = \
            self.acs.event_loop.call_later(
                self.REBOOT_TIMEOUT,
                check_timer,
            )

    def exit(self):
        self.timer_handle.cancel()
        self.timeout_timer = None

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        if not isinstance(message, models.Inform):
            return AcsReadMsgResult(False, None)
        if not does_inform_have_event(message, self.INFORM_EVENT_CODE):
            raise Tr069Error(
                'Did not receive M Reboot event code in '
                'Inform',
            )
        process_inform_message(
            message, self.acs.data_model,
            self.acs.device_cfg,
        )
        return AcsReadMsgResult(True, self.done_transition)

    def state_description(self) -> str:
        return 'Waiting for M Reboot code from Inform'


class WaitRebootDelayState(EnodebAcsState):
    """
    After receiving the Inform notifying us that the eNodeB has successfully
    rebooted, wait a short duration to prevent unspecified race conditions
    that may occur w.r.t reboot
    """

    # Short delay timer to prevent race conditions w.r.t. reboot
    SHORT_CONFIG_DELAY = 10

    def __init__(self, acs: EnodebAcsStateMachine, when_done: str):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done
        self.config_timer = None
        self.timer_handle = None

    def enter(self):
        self.config_timer = StateMachineTimer(self.SHORT_CONFIG_DELAY)

        def check_timer() -> None:
            if self.config_timer.is_done():
                self.acs.transition(self.done_transition)

        self.timer_handle = \
            self.acs.event_loop.call_later(
                self.SHORT_CONFIG_DELAY,
                check_timer,
            )

    def exit(self):
        self.timer_handle.cancel()
        self.config_timer = None

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        return AcsReadMsgResult(True, None)

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        return AcsMsgAndTransition(models.DummyInput(), None)

    def state_description(self) -> str:
        return 'Waiting after eNB reboot to prevent race conditions'


class ErrorState(EnodebAcsState):
    """
    The eNB handler will enter this state when an unhandled Fault is received.

    If the inform_transition_target constructor parameter is non-null, this
    state will attempt to autoremediate by transitioning to the specified
    target state when an Inform is received.
    """

    def __init__(
        self, acs: EnodebAcsStateMachine,
        inform_transition_target: Optional[str] = None,
    ):
        super().__init__()
        self.acs = acs
        self.inform_transition_target = inform_transition_target

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        return AcsReadMsgResult(True, None)

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        if not self.inform_transition_target:
            return AcsMsgAndTransition(models.DummyInput(), None)

        if isinstance(message, models.Inform):
            return AcsMsgAndTransition(
                models.DummyInput(),
                self.inform_transition_target,
            )
        return AcsMsgAndTransition(models.DummyInput(), None)

    def state_description(self) -> str:
        return 'Error state - awaiting manual restart of enodebd service or ' \
               'an Inform to be received from the eNB'
