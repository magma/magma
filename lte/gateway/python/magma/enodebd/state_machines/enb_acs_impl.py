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

import traceback
from abc import abstractmethod
from typing import Any, Dict

from magma.common.service import MagmaService
from magma.enodebd import metrics
from magma.enodebd.data_models.data_model_parameters import ParameterName
from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration
from magma.enodebd.exceptions import ConfigurationError
from magma.enodebd.logger import EnodebdLogger as logger
from magma.enodebd.state_machines.enb_acs import EnodebAcsStateMachine
from magma.enodebd.state_machines.enb_acs_states import EnodebAcsState
from magma.enodebd.state_machines.timer import StateMachineTimer
from magma.enodebd.tr069 import models
from magma.enodebd.tr069.models import Tr069ComplexModel


class BasicEnodebAcsStateMachine(EnodebAcsStateMachine):
    """
    Most of the EnodebAcsStateMachine classes for each device work about the
    same way. Differences lie mainly in the data model, desired configuration,
    and the state transition map.

    This class specifies the shared implementation between them.
    """

    # eNodeB connection timeout is used to determine whether or not eNodeB is
    # connected to enodebd based on time of last Inform message. By default,
    # periodic inform interval is 30secs, so timeout should be larger than
    # this.
    # Also set timer longer than reboot time, so that an eNodeB reboot does not
    # trigger a connection-timeout alarm.
    ENB_CONNECTION_TIMEOUT = 600  # In seconds

    # If eNodeB is disconnected from MME for an unknown reason for this time,
    # then reboot it. Set to a long time to ensure this doesn't interfere with
    # other enodebd configuration processes - it is just a measure of last
    # resort for an unlikely error case
    MME_DISCONNECT_ENODEB_REBOOT_TIMER = 15 * 60

    # Check the MME connection status every 15 seconds
    MME_CHECK_TIMER = 15

    def __init__(
            self,
            service: MagmaService,
    ) -> None:
        super().__init__()
        self.state = None
        self.timeout_handler = None
        self.mme_timeout_handler = None
        self.mme_timer = None
        self._start_state_machine(service)

    def get_state(self) -> str:
        if self.state is None:
            logger.warning('ACS State machine is not in any state.')
            return 'N/A'
        return self.state.state_description()

    def handle_tr069_message(
            self,
            message: Tr069ComplexModel,
    ) -> Tr069ComplexModel:
        """
        Accept the tr069 message from the eNB and produce a reply.

        States may transition after reading a message but BEFORE producing
        a reply. Most steps in the provisioning process are represented as
        beginning with enodebd sending a request to the eNB, and waiting for
        the reply from the eNB.
        """
        # TransferComplete messages come at random times, and we ignore them
        if isinstance(message, models.TransferComplete):
            return models.TransferCompleteResponse()
        try:
            self._read_tr069_msg(message)
            return self._get_tr069_msg(message)
        except Exception:  # pylint: disable=broad-except
            logger.error('Failed to handle tr069 message')
            logger.error(traceback.format_exc())
            self._dump_debug_info()
            self.transition(self.unexpected_fault_state_name)
            return self._get_tr069_msg(message)

    def transition(self, next_state: str) -> Any:
        logger.debug('State transition to <%s>', next_state)
        self.state.exit()
        self.state = self.state_map[next_state]
        self.state.enter()

    def stop_state_machine(self) -> None:
        """ Clean up anything the state machine is tracking or doing """
        self.state.exit()
        if self.timeout_handler is not None:
            self.timeout_handler.cancel()
            self.timeout_handler = None
        if self.mme_timeout_handler is not None:
            self.mme_timeout_handler.cancel()
            self.mme_timeout_handler = None
        self._service = None
        self._desired_cfg = None
        self._device_cfg = None
        self._data_model = None

        self.mme_timer = None

    def _start_state_machine(
            self,
            service: MagmaService,
    ):
        self.service = service
        self.data_model = self.data_model_class()
        # The current known device config has few known parameters
        # The desired configuration depends on what the current configuration
        # is. This we don't know fully, yet.
        self.device_cfg = EnodebConfiguration(self.data_model)

        self._init_state_map()
        self.state = self.state_map[self.disconnected_state_name]
        self.state.enter()
        self._reset_timeout()
        self._periodic_check_mme_connection()

    def _reset_state_machine(
        self,
        service: MagmaService,
    ):
        self.stop_state_machine()
        self._start_state_machine(service)

    def _read_tr069_msg(self, message: Any) -> None:
        """ Process incoming message and maybe transition state """
        self._reset_timeout()
        msg_handled, next_state = self.state.read_msg(message)
        if not msg_handled:
            self._transition_for_unexpected_msg(message)
            _msg_handled, next_state = self.state.read_msg(message)
        if next_state is not None:
            self.transition(next_state)

    def _get_tr069_msg(self, message: Any) -> Any:
        """ Get a new message to send, and maybe transition state """
        msg_and_transition = self.state.get_msg(message)
        if msg_and_transition.next_state:
            self.transition(msg_and_transition.next_state)
        msg = msg_and_transition.msg
        return msg

    def _transition_for_unexpected_msg(self, message: Any) -> None:
        """
        eNB devices may send an Inform message in the middle of a provisioning
        session. To deal with this, transition to a state that expects an
        Inform message, but also track the status of the eNB as not having
        been disconnected.
        """
        if isinstance(message, models.Inform):
            logger.debug(
                'ACS in (%s) state. Received an Inform message',
                self.state.state_description(),
            )
            self._reset_state_machine(self.service)
        elif isinstance(message, models.Fault):
            logger.debug(
                'ACS in (%s) state. Received a Fault <%s>',
                self.state.state_description(), message.FaultString,
            )
            self.transition(self.unexpected_fault_state_name)
        else:
            raise ConfigurationError('Cannot handle unexpected TR069 msg')

    def _reset_timeout(self) -> None:
        if self.timeout_handler is not None:
            self.timeout_handler.cancel()

        def timed_out():
            self.transition(self.disconnected_state_name)

        self.timeout_handler = self.event_loop.call_later(
            self.ENB_CONNECTION_TIMEOUT,
            timed_out,
        )

    def _periodic_check_mme_connection(self) -> None:
        self._check_mme_connection()
        self.mme_timeout_handler = self.event_loop.call_later(
            self.MME_CHECK_TIMER,
            self._periodic_check_mme_connection,
        )

    def _check_mme_connection(self) -> None:
        """
        Check if eNodeB should be connected to MME but isn't, and maybe reboot.

        If the eNB doesn't report connection to MME within a timeout period,
        get it to reboot in the hope that it will fix things.

        Usually, enodebd polls the eNodeB for whether it is connected to MME.
        This method checks the last polled MME connection status, and if
        eNodeB should be connected to MME but it isn't.
        """
        if self.device_cfg.has_parameter(ParameterName.MME_STATUS) and \
                self.device_cfg.get_parameter(ParameterName.MME_STATUS):
            is_mme_connected = 1
        else:
            is_mme_connected = 0

        # True if we would expect MME to be connected, but it isn't
        is_mme_unexpectedly_dc = \
            self.is_enodeb_connected() \
            and self.is_enodeb_configured() \
            and self.mconfig.allow_enodeb_transmit \
            and not is_mme_connected

        if is_mme_unexpectedly_dc:
            logger.warning(
                'eNodeB is connected to AGw, is configured, '
                'and has AdminState enabled for transmit. '
                'MME connection to eNB is missing.',
            )
            if self.mme_timer is None:
                logger.warning(
                    'eNodeB will be rebooted if MME connection '
                    'is not established in: %s seconds.',
                    self.MME_DISCONNECT_ENODEB_REBOOT_TIMER,
                )
                metrics.STAT_ENODEB_REBOOT_TIMER_ACTIVE.set(1)
                self.mme_timer = \
                    StateMachineTimer(self.MME_DISCONNECT_ENODEB_REBOOT_TIMER)
            elif self.mme_timer.is_done():
                logger.warning(
                    'eNodeB has not established MME connection '
                    'within %s seconds - rebooting!',
                    self.MME_DISCONNECT_ENODEB_REBOOT_TIMER,
                )
                metrics.STAT_ENODEB_REBOOTS.labels(cause='MME disconnect').inc()
                metrics.STAT_ENODEB_REBOOT_TIMER_ACTIVE.set(0)
                self.mme_timer = None
                self.reboot_asap()
            else:
                # eNB is not connected to MME, but we're still waiting to see
                # if it will connect within the timeout period.
                # Take no action for now.
                pass
        else:
            if self.mme_timer is not None:
                logger.info('eNodeB has established MME connection.')
                self.mme_timer = None
            metrics.STAT_ENODEB_REBOOT_TIMER_ACTIVE.set(0)

    def _dump_debug_info(self) -> None:
        if self.device_cfg is not None:
            logger.error(
                'Device configuration: %s',
                self.device_cfg.get_debug_info(),
            )
        else:
            logger.error('Device configuration: None')
        if self.desired_cfg is not None:
            logger.error(
                'Desired configuration: %s',
                self.desired_cfg.get_debug_info(),
            )
        else:
            logger.error('Desired configuration: None')

    @abstractmethod
    def _init_state_map(self) -> None:
        pass

    @property
    @abstractmethod
    def state_map(self) -> Dict[str, EnodebAcsState]:
        pass

    @property
    @abstractmethod
    def disconnected_state_name(self) -> str:
        pass

    @property
    @abstractmethod
    def unexpected_fault_state_name(self) -> str:
        """ State to handle unexpected Fault messages """
        pass
