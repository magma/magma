#!/usr/bin/env python3

"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import subprocess
import sys

import fire as fire
from orc8r.protos import common_pb2
from orc8r.protos.magmad_pb2_grpc import MagmadStub
from magma.common.service_registry import ServiceRegistry
from pystemd.systemd1 import Unit


class ActiveState:
    ACTIVE = 'active'
    RELOADING = 'reloading'
    INACTIVE = 'inactive'
    FAILED = 'failed'
    ACTIVATING = 'activating'
    DEACTIVATING = 'deactivating'

    dbus2state = {
        b'active': ACTIVE,
        b'reloading': RELOADING,
        b'inactive': INACTIVE,
        b'failed': FAILED,
        b'activating': ACTIVATING,
        b'deactivating': DEACTIVATING,
    }

    state2status = {
        ACTIVE: u'\u2714',
        RELOADING: u'\u27A4',
        INACTIVE: u'\u2717',
        FAILED: u'\u2717',
        ACTIVATING: u'\u27A4',
        DEACTIVATING: u'\u27A4',
    }


class ServiceHealth:
    def __init__(self, service_name, active_state, sub_state, time_running):
        self.service_name = service_name
        self.active_state = active_state
        self.sub_state = sub_state
        self.time_running = time_running

    def __str__(self):
        return '{} {:20} {:10} {:15} Running for {}'.format(
            ActiveState.state2status[self.active_state],
            self.service_name,
            self.active_state,
            self.sub_state,
            self.time_running
        )


def magma_services_status():
    # DBus Unit objects: https://www.freedesktop.org/wiki/Software/systemd/dbus/
    chan = ServiceRegistry.get_rpc_channel('magmad', ServiceRegistry.LOCAL)
    client = MagmadStub(chan)

    configs = client.GetConfigs(common_pb2.Void())
    health_summary = []

    for service_name in configs.configs_by_key:
        unit = Unit('magma@{}.service'.format(service_name), _autoload=True)
        active_state = ActiveState.dbus2state[unit.Unit.ActiveState]
        sub_state = str(unit.Unit.SubState, 'utf-8')
        if active_state == ActiveState.ACTIVE:
            pid = unit.Service.MainPID
            process = subprocess.Popen('ps -o etime= -p {}'.format(pid).split(),
                                       stdout=subprocess.PIPE)

            time_running, error = process.communicate()
        else:
            time_running = b'00'

        health_summary.append(ServiceHealth(
            service_name=service_name,
            active_state=active_state, sub_state=sub_state,
            time_running=str(time_running, 'utf-8').strip()
        ))

    return health_summary


if __name__ == '__main__':
    if len(sys.argv) == 1:
        print('\n'.join([str(h) for h in magma_services_status()]))
    else:
        fire.Fire({
            'status': magma_services_status,
        })
