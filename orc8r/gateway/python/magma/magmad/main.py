"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import importlib
import logging

import snowflake
from magma.common.sdwatchdog import SDWatchdog
from magma.common.service import MagmaService
from magma.common.streamer import StreamerClient
from magma.configuration.mconfig_managers import MconfigManagerImpl, \
    get_mconfig_manager
from magma.magmad.logging.systemd_tailer import start_systemd_tailer
from magma.magmad.upgrade.upgrader import UpgraderFactory, start_upgrade_loop

from .bootstrap_manager import BootstrapManager
from .checkin_manager import CheckinManager
from .config_manager import CONFIG_STREAM_NAME, ConfigManager
from .metrics import metrics_collection_loop, monitor_unattended_upgrade_status
from .rpc_servicer import MagmadRpcServicer
from .service_manager import ServiceManager
from .service_poller import ServicePoller
from .sync_rpc_client import SyncRPCClient


def main():
    """
    Main magmad function
    """
    service = MagmaService('magmad')

    logging.info('Starting magmad for UUID: %s', snowflake.make_snowflake())

    # Create service manager
    services = service.config['magma_services']
    init_system = service.config.get('init_system', 'systemd')
    registered_dynamic_services = service.config.get(
        'registered_dynamic_services', [])
    enabled_dynamic_services = []
    if service.mconfig is not None:
        enabled_dynamic_services = service.mconfig.dynamic_services

    # Poll the services' Service303 interface
    service_poller = ServicePoller(service.loop, service.config)
    service_poller.start()

    service_manager = ServiceManager(services, init_system, service_poller,
                                     registered_dynamic_services,
                                     enabled_dynamic_services)

    # Start a background thread to stream updates from the cloud
    stream_client = None
    if service.config.get('enable_config_streamer', False):
        stream_client = StreamerClient(
            {
                CONFIG_STREAM_NAME: ConfigManager(
                    services, service_manager,
                    service,
                    MconfigManagerImpl(),
                ),
            },
            service.loop,
        )

    # Schedule periodic checkins
    checkin_manager = CheckinManager(service, service_poller)

    # Create sync rpc client with a timeout of 60 seconds
    sync_rpc_client = None
    if service.config.get('enable_sync_rpc', False):
        sync_rpc_client = SyncRPCClient(service.loop, 60)

    first_time_bootstrap = True

    # This is called when bootstrap succeeds and when _bootstrap_check is
    # invoked but bootstrap is not needed. If it's invoked right after certs
    # are generated, certs_generated is true, control_proxy will restart.
    async def bootstrap_success_cb(certs_generated):
        nonlocal first_time_bootstrap
        if first_time_bootstrap:
            if stream_client:
                stream_client.start()
            await checkin_manager.try_checkin()
            if sync_rpc_client:
                sync_rpc_client.start()
            first_time_bootstrap = False
        if certs_generated and 'control_proxy' in services:
            service.loop.create_task(
                service_manager.restart_services(services=['control_proxy'])
            )

    # Create bootstrap manager
    bootstrap_manager = BootstrapManager(service, bootstrap_success_cb)

    async def checkin_failure_cb(err_code):
        await bootstrap_manager.on_checkin_fail(err_code)
    checkin_manager.set_failure_cb(checkin_failure_cb)

    # Start bootstrap_manager after checkin_manager's callback is set
    bootstrap_manager.start_bootstrap_manager()

    # Start all services when magmad comes up
    service.loop.create_task(service_manager.start_services())

    # Start upgrade manager loop
    if service.config.get('enable_upgrade_manager', False):
        upgrader = _get_upgrader_impl(service)
        service.loop.create_task(start_upgrade_loop(service, upgrader))

    # Start network health metric collection
    if service.config.get('enable_network_monitor', False):
        service.loop.create_task(metrics_collection_loop(service.config))

    if service.config.get('enable_systemd_tailer', False):
        service.loop.create_task(start_systemd_tailer(service.config))

    # Start loop to monitor unattended upgrade status
    service.loop.create_task(monitor_unattended_upgrade_status(service.loop))

    # Add all servicers to the server
    magmad_servicer = MagmadRpcServicer(
        service,
        services, service_manager, get_mconfig_manager(),
        service.loop,
    )
    magmad_servicer.add_to_server(service.rpc_server)

    if SDWatchdog.has_notify():
        # Create systemd watchdog
        sdwatchdog = SDWatchdog(
            tasks=[bootstrap_manager, checkin_manager],
            update_status=True)
        # Start watchdog loop
        service.loop.create_task(sdwatchdog.run())

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


def _get_upgrader_impl(service):
    """
    Instantiate the MagmaUpgrader implementation to be used for magmad, as
    specified in the service YML config.

    Args:
        service (MagmaService): magmad service

    Returns:
        magma.magmad.upgrade.upgrader.Upgrader: specified upgrader instance
    """
    # Get the factory class from yml config
    factory_cfg = service.config.get('upgrader_factory', None)
    assert factory_cfg is not None, \
        'upgrader_factory is required in magmad service config'

    factory_module = factory_cfg.get('module', '')
    factory_clsname = factory_cfg.get('class', '')
    assert factory_module and factory_clsname, \
        'module and class are required in upgrader_factory config'

    # Instantiate factory class
    FactoryClass = getattr(importlib.import_module(factory_module),
                           factory_clsname)
    factory_impl = FactoryClass()
    assert isinstance(factory_impl, UpgraderFactory),\
        'upgrader_factory must be a subclass of UpgraderFactory'

    return factory_impl.create_upgrader(service, service.loop)


if __name__ == "__main__":
    main()
