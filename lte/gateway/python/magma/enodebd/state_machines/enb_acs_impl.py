"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging
from typing import Any, Dict
from abc import abstractmethod
from magma.common.service import MagmaService
from magma.enodebd import metrics
from magma.enodebd.device_config.enodeb_configuration import \
    EnodebConfiguration
from magma.enodebd.enodeb_status import get_enodeb_status
from magma.enodebd.exceptions import ConfigurationError
from magma.enodebd.state_machines.enb_acs import EnodebAcsStateMachine
from magma.enodebd.state_machines.enb_acs_states import EnodebAcsState
from magma.enodebd.state_machines.timer import StateMachineTimer
from magma.enodebd.stats_manager import StatsManager
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
            stats_mgr: StatsManager,
    ) -> None:
        super().__init__()
        self.state = None
        self.timeout_handler = None
        self.mme_timeout_handler = None
        self.mme_timer = None
        self._start_state_machine(service, stats_mgr)

    def get_state(self) -> str:
        if self.state is None:
            logging.warning('ACS State machine is not in any state.')
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
        self._read_tr069_msg(message)
        return self._get_tr069_msg()

    def transition(self, next_state: str) -> Any:
        logging.debug('State transition to <%s>', next_state)
        self.state.exit()
        self.state = self.state_map[next_state]
        self.state.enter()

    def stop_state_machine(self) -> None:
        """ Clean up anything the state machine is tracking or doing """
        self.state.exit()
        if self.timeout_handler is not None:
            self.mme_timeout_handler.cancel()
            self.timeout_handler.cancel()
        self._service = None
        self._stats_manager = None
        self._desired_cfg = None
        self._device_cfg = None
        self._data_model = None

        self.timeout_handler = None
        self.mme_timeout_handler = None
        self.mme_timer = None

    def _start_state_machine(
            self,
            service: MagmaService,
            stats_mgr: StatsManager,
    ):
        self.service = service
        self.stats_manager = stats_mgr
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
        stats_mgr: StatsManager,
    ):
        self.stop_state_machine()
        self._start_state_machine(service, stats_mgr)

    def _read_tr069_msg(self, message: Any) -> None:
        """ Process incoming message and maybe transition state """
        self._reset_timeout()
        msg_handled, next_state = self.state.read_msg(message)
        if not msg_handled:
            self._transition_for_unexpected_msg(message)
            _msg_handled, next_state = self.state.read_msg(message)
        if next_state is not None:
            self.transition(next_state)

    def _get_tr069_msg(self) -> Any:
        """ Get a new message to send, and maybe transition state """
        msg_and_transition = self.state.get_msg()
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
            logging.debug('ACS in (%s) state. Received an Inform message',
                          self.state.state_description())
            self._reset_state_machine(self.service, self.stats_manager)
        elif isinstance(message, models.Fault):
            logging.debug('ACS in (%s) state. Received a Fault <%s>',
                          self.state.state_description(), message.FaultString)
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
        """
        logging.info('Checking mme connection')
        status = get_enodeb_status(self)

        reboot_disabled = \
            not self.is_enodeb_connected() \
            or not self.is_enodeb_configured() \
            or status['mme_connected'] == '1' \
            or not self.mconfig.allow_enodeb_transmit

        if reboot_disabled:
            if self.mme_timer is not None:
                logging.info('Clearing eNodeB reboot timer')
            metrics.STAT_ENODEB_REBOOT_TIMER_ACTIVE.set(0)
            self.mme_timer = None
            return

        if self.mme_timer is None:
            logging.info('Set eNodeB reboot timer: %s',
                         self.MME_DISCONNECT_ENODEB_REBOOT_TIMER)
            metrics.STAT_ENODEB_REBOOT_TIMER_ACTIVE.set(1)
            self.mme_timer = \
                StateMachineTimer(self.MME_DISCONNECT_ENODEB_REBOOT_TIMER)
        elif self.mme_timer.is_done():
            logging.warning('eNodeB reboot timer expired - rebooting!')
            metrics.STAT_ENODEB_REBOOTS.labels(cause='MME disconnect').inc()
            metrics.STAT_ENODEB_REBOOT_TIMER_ACTIVE.set(0)
            self.mme_timer = None
            self.reboot_asap()
        else:
            # eNB is not connected to MME, but we're still waiting to see if
            # it will connect within the timeout period.
            # Take no action for now.
            pass

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
