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
import importlib
import logging
import typing

import snowflake
from magma.common.grpc_client_manager import GRPCClientManager
from magma.common.sdwatchdog import SDWatchdog
from magma.common.sentry import sentry_init
from magma.common.service import MagmaService
from magma.common.streamer import StreamerClient
from magma.configuration.mconfig_managers import (
    MconfigManagerImpl,
    get_mconfig_manager,
)
from magma.magmad.bootstrap_manager import BootstrapManager
from magma.magmad.config_manager import CONFIG_STREAM_NAME, ConfigManager
from magma.magmad.gateway_status import (
    GatewayStatusFactory,
    KernelVersionsPoller,
)
from magma.magmad.generic_command.command_executor import (
    get_command_executor_impl,
)
from magma.magmad.metrics import (
    metrics_collection_loop,
    monitor_unattended_upgrade_status,
)
from magma.magmad.metrics_collector import MetricsCollector, ScrapeTarget
from magma.magmad.rpc_servicer import MagmadRpcServicer
from magma.magmad.service_health_watchdog import ServiceHealthWatchdog
from magma.magmad.service_manager import ServiceManager
from magma.magmad.service_poller import ServicePoller
from magma.magmad.state_reporter import StateReporter
from magma.magmad.sync_rpc_client import SyncRPCClient
from magma.magmad.upgrade.upgrader import UpgraderFactory, start_upgrade_loop
from orc8r.protos.mconfig import mconfigs_pb2
from orc8r.protos.state_pb2_grpc import StateServiceStub


def main():
    """
    Main magmad function
    """
    service = MagmaService('magmad', mconfigs_pb2.MagmaD())

    # Optionally pipe errors to Sentry
    sentry_init(service_name=service.name, sentry_mconfig=service.shared_mconfig.sentry_config)

    logging.info('Starting magmad for UUID: %s', snowflake.make_snowflake())

    # Create service manager
    services = service.config.get('magma_services')
    init_system = service.config.get('init_system', 'systemd')
    registered_dynamic_services = service.config.get(
        'registered_dynamic_services', [],
    )
    enabled_dynamic_services = []
    if service.mconfig is not None:
        enabled_dynamic_services = service.mconfig.dynamic_services

    # Poll the services' Service303 interface
    service_poller = ServicePoller(
        service.loop, service.config,
        enabled_dynamic_services,
    )
    service_poller.start()

    service_manager = ServiceManager(
        services, init_system, service_poller,
        registered_dynamic_services,
        enabled_dynamic_services,
    )

    # Get metrics service config
    metrics_config = service.config.get('metricsd')
    metrics_services = metrics_config['services']
    collect_interval = metrics_config['collect_interval']
    sync_interval = metrics_config['sync_interval']
    grpc_timeout = metrics_config['grpc_timeout']
    grpc_msg_size = metrics_config.get('max_grpc_msg_size_mb', 4)
    metrics_post_processor_fn = metrics_config.get('post_processing_fn')

    metric_scrape_targets = [
        ScrapeTarget(t['url'], t['name'], t['interval'])
        for t in
        metrics_config.get('metric_scrape_targets', [])
    ]

    # Create local metrics collector
    metrics_collector = MetricsCollector(
        services=metrics_services,
        collect_interval=collect_interval,
        sync_interval=sync_interval,
        grpc_timeout=grpc_timeout,
        grpc_max_msg_size_mb=grpc_msg_size,
        loop=service.loop,
        post_processing_fn=get_metrics_postprocessor_fn(
            metrics_post_processor_fn,
        ),
        scrape_targets=metric_scrape_targets,
    )

    # Poll and sync the metrics collector loops
    metrics_collector.run()

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

    # Create sync rpc client with a heartbeat of 30 seconds (timeout = 60s)
    sync_rpc_client = None
    if service.config.get('enable_sync_rpc', False):
        sync_rpc_client = SyncRPCClient(
            service.loop, 30,
            service.config.get('print_grpc_payload', False),
        )

    first_time_bootstrap = True

    # This is called when bootstrap succeeds and when _bootstrap_check is
    # invoked but bootstrap is not needed. If it's invoked right after certs
    # are generated, certs_generated is true, control_proxy will restart.
    async def bootstrap_success_cb(certs_generated: bool):
        nonlocal first_time_bootstrap
        if first_time_bootstrap:
            if stream_client:
                stream_client.start()
            if sync_rpc_client:
                sync_rpc_client.start()
            first_time_bootstrap = False
        if certs_generated:
            svcs_to_restart = []
            if 'control_proxy' in services:
                svcs_to_restart.append('control_proxy')

            # fluent-bit caches TLS client certs in memory, so we need to
            # restart it whenever the certs change
            fresh_mconfig = get_mconfig_manager().load_service_mconfig(
                'magmad', mconfigs_pb2.MagmaD(),
            )
            dynamic_svcs = fresh_mconfig.dynamic_services or []
            if 'td-agent-bit' in dynamic_svcs:
                svcs_to_restart.append('td-agent-bit')

            await service_manager.restart_services(services=svcs_to_restart)

    # Create bootstrap manager
    bootstrap_manager = BootstrapManager(service, bootstrap_success_cb)

    # Initialize kernel version poller if it is enabled
    kernel_version_poller = None
    if service.config.get('enable_kernel_version_checking', False):
        kernel_version_poller = KernelVersionsPoller(service)
        kernel_version_poller.start()

    # gateway status generator to bundle various information about this
    # gateway into an object.
    gateway_status_factory = GatewayStatusFactory(
        service=service,
        service_poller=service_poller,
        kernel_version_poller=kernel_version_poller,
    )

    # _grpc_client_manager to manage grpc client recycling
    grpc_client_manager = GRPCClientManager(
        service_name="state",
        service_stub=StateServiceStub,
        max_client_reuse=60,
    )

    # Initialize StateReporter
    state_reporter = StateReporter(
        config=service.config,
        mconfig=service.mconfig,
        loop=service.loop,
        bootstrap_manager=bootstrap_manager,
        gw_status_factory=gateway_status_factory,
        grpc_client_manager=grpc_client_manager,
    )

    # Initialize ServiceHealthWatchdog
    service_health_watchdog = ServiceHealthWatchdog(
        config=service.config,
        loop=service.loop,
        service_poller=service_poller,
        service_manager=service_manager,
    )

    # Start _bootstrap_manager
    bootstrap_manager.start_bootstrap_manager()

    # Start all services when magmad comes up
    service.loop.create_task(service_manager.start_services())

    # Start state reporting loop
    state_reporter.start()

    # Start service timeout health check loop
    service_health_watchdog.start()

    # Start upgrade manager loop
    if service.config.get('enable_upgrade_manager', False):
        upgrader = _get_upgrader_impl(service)
        service.loop.create_task(start_upgrade_loop(service, upgrader))

    # Start network health metric collection
    if service.config.get('enable_network_monitor', False):
        service.loop.create_task(metrics_collection_loop(service.config))

    # Create generic command executor
    command_executor = None
    if service.config.get('generic_command_config', None):
        command_executor = get_command_executor_impl(service)

    # Start loop to monitor unattended upgrade status
    service.loop.create_task(monitor_unattended_upgrade_status())

    # Add all servicers to the server
    magmad_servicer = MagmadRpcServicer(
        service,
        services, service_manager, get_mconfig_manager(), command_executor,
        service.loop,
        service.config.get('print_grpc_payload', False),
    )
    magmad_servicer.add_to_server(service.rpc_server)

    if SDWatchdog.has_notify():
        # Create systemd watchdog
        sdwatchdog = SDWatchdog(
            tasks=[bootstrap_manager, state_reporter],
            update_status=True,
        )
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
    FactoryClass = getattr(
        importlib.import_module(factory_module),
        factory_clsname,
    )
    factory_impl = FactoryClass()
    assert isinstance(factory_impl, UpgraderFactory), (
        'upgrader_factory '
        'must be a subclass '
        'of UpgraderFactory'
    )

    return factory_impl.create_upgrader(service, service.loop)


def get_metrics_postprocessor_fn(module_fn_name: typing.Optional[str]):
    """Load and validate the metricsd post processing function"""
    if not module_fn_name:
        return None
    module, sep, fn_name = module_fn_name.rpartition(".")
    assert all(
        (module, sep, fn_name),
    ), "metrics_postprocessor_fn needs to be import.path.func_name"
    fn = getattr(importlib.import_module(module), fn_name)
    assert callable(fn), "metrics_postprocessor_fn is not callable"
    fn([])  # One way to check args
    return fn


if __name__ == "__main__":
    main()
