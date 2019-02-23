"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from magma.enodebd.state_machines.enb_acs import EnodebAcsStateMachine


class StateMachinePointer:
    """
    This is a hack to deal with the possibility that the specified data model
    doesn't match the eNB device enodebd ends up connecting to.

    When the data model doesn't match, the state machine is replaced with one
    that matches the data model.
    """
    def __init__(self, acs_state_machine: EnodebAcsStateMachine):
        self._acs_state_machine = acs_state_machine

    @property
    def state_machine(self):
        return self._acs_state_machine

    @state_machine.setter
    def state_machine(
        self,
        acs_state_machine: EnodebAcsStateMachine,
    ) -> None:
        self._acs_state_machine = acs_state_machine
