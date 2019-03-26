"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

# pylint: disable=protected-access
from unittest import TestCase
from magma.enodebd.state_machines.timer import StateMachineTimer


class StateMachineTimerTests(TestCase):
    def test_is_done(self):
        timer_a = StateMachineTimer(0)
        self.assertTrue(timer_a.is_done(), 'Timer should be done')

        timer_b = StateMachineTimer(600)
        self.assertFalse(timer_b.is_done(), 'Timer should not be done')

