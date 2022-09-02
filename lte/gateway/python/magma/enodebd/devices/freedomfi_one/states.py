"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
import logging
from typing import Any, List

from dp.protos.cbsd_pb2 import CBSDStateResult
from magma.enodebd.data_models.data_model import DataModel, InvalidTrParamPath
from magma.enodebd.data_models.data_model_parameters import ParameterName
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
from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration
from magma.enodebd.devices.freedomfi_one.params import (
    CarrierAggregationParameters,
    SASParameters,
    StatusParameters,
)
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
from magma.enodebd.state_machines.enb_acs_states import (
    AcsMsgAndTransition,
    AcsReadMsgResult,
    EndSessionState,
    EnodebAcsState,
    NotifyDPState,
)
from magma.enodebd.tr069 import models

ANTENNA_HEIGHT = 0


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
