#!/usr/bin/env python3
"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
# pylint: skip-file
# pylint does not play well with aioeventlet, as it uses asyncio.async which
# produces a parse error

import asyncio
import logging
import threading


import aioeventlet
from ryu import cfg
from ryu.base.app_manager import AppManager

from magma.common.misc_utils import call_process
from magma.common.service import MagmaService
from magma.configuration import environment
from magma.pipelined.app import of_rest_server
from magma.pipelined.check_quota_server import run_flask
from magma.pipelined.service_manager import ServiceManager
from magma.pipelined.ifaces import monitor_ifaces
from magma.pipelined.rpc_servicer import PipelinedRpcServicer
from lte.protos.mconfig import mconfigs_pb2


def main():
    """
    Loads the Ryu apps we want to run from the config file.
    This should exit on keyboard interrupt.
    """

    # Run asyncio loop in a greenthread so we can evaluate other eventlets
    # TODO: Remove once Ryu migrates to asyncio
    asyncio.set_event_loop_policy(aioeventlet.EventLoopPolicy())

    service = MagmaService('pipelined', mconfigs_pb2.PipelineD())
    service_config = service.config

    if environment.is_dev_mode():
        of_rest_server.configure(service_config)

    # Set Ryu config params
    cfg.CONF.ofp_listen_host = "127.0.0.1"

    # Load the ryu apps
    service_manager = ServiceManager(service)
    service_manager.load()

    def callback(returncode):
        if returncode != 0:
            logging.error(
                "Failed to set MASQUERADE: %d", returncode
            )

    if service.mconfig.nat_enabled:
        call_process('iptables -t nat -A POSTROUTING -o %s -j MASQUERADE'
                     % service.config['nat_iface'],
                     callback,
                     service.loop
                     )

    service.loop.create_task(monitor_ifaces(
        service.config['monitored_ifaces'],
        service.loop),
    )

    manager = AppManager.get_instance()
    # Add pipelined rpc servicer
    pipelined_srv = PipelinedRpcServicer(
        service.loop,
        manager.applications.get('GYController', None),
        manager.applications.get('EnforcementController', None),
        manager.applications.get('EnforcementStatsController', None),
        manager.applications.get('DPIController', None),
        manager.applications.get('UEMacAddressController', None),
        manager.applications.get('CheckQuotaController', None),
        manager.applications.get('IPFIXController', None),
        manager.applications.get('VlanLearnController', None),
        manager.applications.get('TunnelLearnController', None),
        service_manager)
    pipelined_srv.add_to_server(service.rpc_server)

    if service.config['setup_type'] == 'CWF':
        bridge_ip = service.config['bridge_ip_address']
        has_quota_port = service.config['has_quota_port']
        no_quota_port = service.config['no_quota_port']

        def on_exit_server_thread():
            service.StopService(None, None)

        # For CWF start quota check servers
        start_check_quota_server(run_flask, bridge_ip, has_quota_port, True,
                                 on_exit_server_thread)
        start_check_quota_server(run_flask, bridge_ip, no_quota_port, False,
                                 on_exit_server_thread)

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


def start_check_quota_server(target, ip, port, response, exit_callback):
    """ Starts service server threads """
    thread = threading.Thread(
        target=target,
        args=(ip, port, response, exit_callback))
    thread.daemon = True
    thread.start()


if __name__ == "__main__":
    main()
