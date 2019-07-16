#!/usr/bin/env python3

"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import asyncio
import subprocess
import sys

import apt
import fire as fire
from magma.common.service import MagmaService
from magma.common.service_registry import ServiceRegistry
from magma.configuration.mconfig_managers import load_service_mconfig_as_json
from magma.magmad.metrics import UNEXPECTED_SERVICE_RESTARTS
from orc8r.protos import common_pb2, magmad_pb2
from orc8r.protos.magmad_pb2_grpc import MagmadStub
from orc8r.protos.mconfig import mconfigs_pb2
from pystemd.systemd1 import Unit

from magma.magmad.service_poller import ServicePoller
from . import checkin_cli


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


class Errors:
    def __init__(self, log_level, error_count):
        self.log_level = log_level
        self.error_count = error_count

    def __str__(self):
        return '{}: {}'.format(self.log_level, self.error_count)


class RestartFrequency:
    def __init__(self, count, time_interval):
        self.count = count
        self.time_interval = time_interval

    def __str__(self):
        return 'Restarted {} times {}'.format(
            self.count,
            self.time_interval,
        )


class HealthStatus:
    DOWN = 'Down'
    UP = 'Up'
    UNKNOWN = 'Unknown'


class Version:
    def __init__(self, version_code, last_update_time):
        self.version_code = version_code
        self.last_update_time = last_update_time

    def __str__(self):
        return '{}, last updated: {}'.format(
            self.version_code,
            self.last_update_time,
        )


class ServiceHealth:
    def __init__(self, service_name, active_state, sub_state,
                 time_running, errors):
        self.service_name = service_name
        self.active_state = active_state
        self.sub_state = sub_state
        self.time_running = time_running
        self.errors = errors

    def __str__(self):
        return '{} {:20} {:10} {:15} {:10} {:>10} {:>10}'.format(
            ActiveState.state2status[self.active_state],
            self.service_name,
            self.active_state,
            self.sub_state,
            self.time_running,
            self.errors.log_level,
            self.errors.error_count,
        )


class HealthSummary:
    def __init__(self, version,
                 services_health,
                 internet_health, dns_health,
                 unexpected_restarts):
        self.version = version
        self.services_health = services_health
        self.internet_health = internet_health
        self.dns_health = dns_health
        self.unexpected_restarts = unexpected_restarts

    def __str__(self):
        any_restarts = any([restarts.count
                            for restarts in self.unexpected_restarts.values()])
        return """
Version: {}:
  {:20} {:10} {:15} {:10} {:>10} {:>10}
{}

Internet health: {}
DNS health: {}

Restart summary:
{}
        """.format(self.version,
                   'Service', 'Status', 'SubState', 'Running for', 'Log level',
                   'Errors since last restart',
                   '\n'.join([str(h) for h in self.services_health]),
                   self.internet_health, self.dns_health,
                   '\n'.join(['{:20} {}'.format(name, restarts)
                              for name, restarts
                              in self.unexpected_restarts.items()])
                   if any_restarts
                   else "No restarts since the gateway started",
                   )


def ping(host, num_packets=4):
    chan = ServiceRegistry.get_rpc_channel('magmad', ServiceRegistry.LOCAL)
    client = MagmadStub(chan)

    response = client.RunNetworkTests(magmad_pb2.NetworkTestRequest(
        pings=[magmad_pb2.PingParams(host_or_ip=host, num_packets=num_packets)]
    ))
    return response.pings


def ping_status(host):
    pings = ping(host=host, num_packets=4)[0]
    if pings.error:
        return HealthStatus.DOWN
    if pings.avg_response_ms:
        return HealthStatus.UP
    return HealthStatus.UNKNOWN


def get_error_summary(service_names):
    configs = {service_name: load_service_mconfig_as_json(service_name)
               for service_name in service_names}
    res = {service_name: Errors(log_level=configs[service_name]['logLevel'],
                                error_count=0)
           for service_name in service_names}

    with open('/var/log/syslog', 'r') as f:
        for line in f:
            for service_name in service_names:
                if service_name not in line:
                    continue
                # Reset the counter for restart/start
                if 'Starting {}...'.format(service_name) in line:
                    res[service_name].error_count = 0
                elif 'ERROR' in line:
                    res[service_name].error_count += 1
    return res


def get_magma_services_summary():
    """ Get health for all the running services """
    # DBus Unit objects: https://www.freedesktop.org/wiki/Software/systemd/dbus/
    chan = ServiceRegistry.get_rpc_channel('magmad', ServiceRegistry.LOCAL)
    client = MagmadStub(chan)

    configs = client.GetConfigs(common_pb2.Void())
    services_health_summary = []

    service_names = [str(name) for name in configs.configs_by_key]
    services_errors = get_error_summary(service_names=service_names)

    for service_name in service_names:
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

        services_health_summary.append(ServiceHealth(
            service_name=service_name,
            active_state=active_state, sub_state=sub_state,
            time_running=str(time_running, 'utf-8').strip(),
            errors=services_errors[service_name]
        ))
    return services_health_summary


def get_unexpected_restart_summary():
    service = MagmaService('magmad', mconfigs_pb2.MagmaD())
    service_poller = ServicePoller(service.loop, service.config)
    service_poller.start()

    asyncio.set_event_loop(service.loop)

    # noinspection PyProtectedMember
    async def fetch_info():
        restart_frequencies = {}
        await service_poller._get_service_info()
        for service_name in service_poller.service_info.keys():
            restarts = int(UNEXPECTED_SERVICE_RESTARTS
                           .labels(service_name=service_name)
                           ._value.get())
            restart_frequencies[service_name] = RestartFrequency(
                count=restarts,
                time_interval=''
            )

        return restart_frequencies

    return service.loop.run_until_complete(fetch_info())


def get_health_summary():
    """ Get health summary for the whole program """

    # Get magma version and when it was updated
    cache = apt.Cache()
    pkg = cache.get('magma', default=cache['python'])
    version = Version(version_code=pkg.versions[0],
                      last_update_time='19 Jun 2019')

    health_summary = HealthSummary(
        version=version,
        services_health=get_magma_services_summary(),
        internet_health=ping_status(host='8.8.8.8'),
        dns_health=ping_status(host='google.com'),
        unexpected_restarts=get_unexpected_restart_summary(),
    )

    ''' Check connection to the orchestrator '''
    # Call this last because it closes the event loop
    print('\nGateway <-> Controller connectivity')
    checkin_cli.main()

    return str(health_summary)


if __name__ == '__main__':
    print('Health Summary')
    if len(sys.argv) == 1:
        fire.Fire(get_health_summary)
    else:
        fire.Fire({
            'status': get_health_summary,
            'services_status': get_magma_services_summary,
            'restarts_status': get_unexpected_restart_summary,
            'error_status': get_error_summary,
        })
