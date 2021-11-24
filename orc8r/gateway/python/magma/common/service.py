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

import asyncio
import faulthandler
import functools
import logging
import os
import signal
import time
from concurrent import futures
from typing import Any, Dict, List, Optional

import grpc
import pkg_resources
from magma.common.log_count_handler import MsgCounterHandler
from magma.common.log_counter import ServiceLogErrorReporter
from magma.common.metrics_export import get_metrics
from magma.common.service_registry import ServiceRegistry
from magma.configuration.exceptions import LoadConfigError
from magma.configuration.mconfig_managers import get_mconfig_manager
from magma.configuration.service_configs import load_service_config
from orc8r.protos.common_pb2 import LogLevel, Void
from orc8r.protos.mconfig import mconfigs_pb2
from orc8r.protos.metricsd_pb2 import MetricsContainer
from orc8r.protos.service303_pb2 import (
    GetOperationalStatesResponse,
    ReloadConfigResponse,
    ServiceInfo,
    State,
)
from orc8r.protos.service303_pb2_grpc import (
    Service303Servicer,
    Service303Stub,
    add_Service303Servicer_to_server,
)

MAX_DEFAULT_WORKER = 10


async def loop_exit():
    """
    Stop the loop in an async context
    """
    loop = asyncio.get_event_loop()
    loop.stop()


class MagmaService(Service303Servicer):
    """
    MagmaService provides the framework for all Magma services.
    This class also implements the Service303 interface for external
    entities to interact with the service.
    """

    def __init__(self, name, empty_mconfig, loop=None):
        self._name = name
        self._port = 0
        self._get_status_callback = None
        self._get_operational_states_cb = None
        self._log_count_handler = MsgCounterHandler()

        # Init logging before doing anything
        logging.basicConfig(
            level=logging.INFO,
            format='[%(asctime)s %(levelname)s %(name)s] %(message)s',
        )
        # Add a handler to count errors
        logging.root.addHandler(self._log_count_handler)

        # Set gRPC polling strategy
        self._set_grpc_poll_strategy()

        # Load the managed config if present
        self._mconfig = empty_mconfig
        self._shared_mconfig = mconfigs_pb2.SharedMconfig()
        self._mconfig_metadata = None
        self._mconfig_manager = get_mconfig_manager()
        self.reload_mconfig()

        self._state = ServiceInfo.STARTING
        self._health = ServiceInfo.APP_UNHEALTHY
        if loop is None:
            loop = asyncio.get_event_loop()
        self._loop = loop
        self._start_time = int(time.time())
        self._register_signal_handlers()

        # Load the service config if present
        self._config = None
        self.reload_config()

        # Count errors
        self.log_counter = ServiceLogErrorReporter(
            loop=self._loop,
            service_config=self._config,
            handler=self._log_count_handler,
        )
        self.log_counter.start()

        # Operational States
        self._operational_states = []

        self._version = '0.0.0'
        # Load the service version if available
        try:
            # Check if service on docker
            if self._config and 'init_system' in self._config \
                    and self._config['init_system'] == 'docker':
                # image comes in form of "feg_gateway_python:<IMAGE_TAG>\n"
                # Skip the "feg_gateway_python:" part
                image = os.popen(
                    'docker ps --filter name=magmad --format "{{.Image}}" | '
                    'cut -d ":" -f 2',
                )
                image_tag = image.read().strip('\n')
                self._version = image_tag
            else:
                self._version = pkg_resources.get_distribution('orc8r').version
        except pkg_resources.ResolutionError as e:
            logging.info(e)

        if self._config and 'grpc_workers' in self._config:
            self._server = grpc.server(
                futures.ThreadPoolExecutor(
                    max_workers=self._config['grpc_workers'],
                ),
            )
        else:
            self._server = grpc.server(
                futures.ThreadPoolExecutor(max_workers=MAX_DEFAULT_WORKER),
            )
        add_Service303Servicer_to_server(self, self._server)

    @property
    def version(self):
        """Return the current running version of the Magma service"""
        return self._version

    @property
    def name(self):
        """Return the name of service

        Returns:
            tr: name of service
        """
        return self._name

    @property
    def rpc_server(self):
        """Return the RPC server used by the service"""
        return self._server

    @property
    def port(self):
        """Return the listening port of the service"""
        return self._port

    @property
    def loop(self):
        """Return the asyncio event loop used by the service"""
        return self._loop

    @property
    def state(self):
        """Return the state of the service"""
        return self._state

    @property
    def config(self) -> Dict[str, Any]:
        """Return the service config"""
        return self._config

    @property
    def mconfig(self):
        """Return the managed config"""
        return self._mconfig

    @property
    def shared_mconfig(self):
        """Return the shared managed config"""
        return self._shared_mconfig

    @property
    def mconfig_metadata(self):
        """Return the metadata of the managed config"""
        return self._mconfig_metadata

    @property
    def mconfig_manager(self):
        """Return the mconfig manager for this service"""
        return self._mconfig_manager

    def reload_config(self):
        """Reload the local config for the service"""
        try:
            self._config = load_service_config(self._name)
            self._setup_logging()
        except LoadConfigError as e:
            logging.warning(e)

    def reload_mconfig(self):
        """Reload the managed config for the service"""
        # reload mconfig manager in case feature flag for streaming changed
        self._mconfig_manager = get_mconfig_manager()
        try:
            self._mconfig_metadata = self._mconfig_manager.load_mconfig_metadata()
            self._mconfig = self._mconfig_manager.load_service_mconfig(
                self._name,
                self._mconfig,
            )
            self._shared_mconfig = self._mconfig_manager.load_shared_mconfig(self._shared_mconfig)
        except LoadConfigError as e:
            logging.warning(e)

    def add_operational_states(self, states: List[State]):
        """Add a list of states into the service

        Args:
            states (List[State]): [description]
        """
        self._operational_states.extend(states)

    def run(self):
        """
        Start the service and runs the event loop until a term signal
        is received or a StopService rpc call is made on the Service303
        interface.
        """
        logging.info("Starting %s...", self._name)
        (host, port) = ServiceRegistry.get_service_address(self._name)
        self._port = self._server.add_insecure_port('{}:{}'.format(host, port))
        logging.info("Listening on address %s:%d", host, self._port)
        self._state = ServiceInfo.ALIVE
        # Python services are healthy immediately when run
        self._health = ServiceInfo.APP_HEALTHY
        self._server.start()
        self._loop.run_forever()
        # Waiting for the term signal or StopService rpc call

    def close(self):
        """
        Clean up the service before termination. This needs to be
        called atleast once after the service has been inited.
        """
        self._loop.close()
        self._server.stop(0).wait()

    def register_get_status_callback(self, get_status_callback):
        """Register function for getting status

        Must return a map(string, string)
        """
        self._get_status_callback = get_status_callback

    def register_operational_states_callback(self, get_operational_states_cb):
        """Register the callback function that gets called on GetOperationalStates rpc

        Args:
            get_operational_states_cb ([type]): callback function
        """
        self._get_operational_states_cb = get_operational_states_cb

    def _stop(self, reason):
        """Stop the service gracefully"""
        logging.info("Stopping %s with reason %s...", self._name, reason)
        self._state = ServiceInfo.STOPPING
        self._server.stop(0)

        for pending_task in asyncio.Task.all_tasks(self._loop):
            pending_task.cancel()

        self._state = ServiceInfo.STOPPED
        self._health = ServiceInfo.APP_UNHEALTHY
        asyncio.ensure_future(loop_exit())

    def _set_grpc_poll_strategy(self):
        """
        The new default 'epollex' poll strategy is causing fd leaks, leading to
        service crashes after 1024 open fds.
        See https://github.com/grpc/grpc/issues/15759
        """
        os.environ['GRPC_POLL_STRATEGY'] = 'epoll1,poll'

    def _get_log_level_from_config(self) -> Optional[int]:
        if self._config is None:
            return None
        log_level = self._config.get('log_level', None)
        if log_level is None:
            return None
        # convert from log level string to LogLevel enum value
        try:
            proto_level = LogLevel.Value(log_level)
        except ValueError:
            logging.error(
                'Unknown logging level in config: %s, ignoring',
                log_level,
            )
            return None
        return proto_level

    def _get_log_level_from_mconfig(self) -> Optional[int]:
        if self._mconfig is None:
            return None
        return self._mconfig.log_level

    def _setup_logging(self):
        """Set up log level from config values

        The config file on the AGW takes precedence over the mconfig
        If neither config file nor mconfig has the log level config, default to INFO
        """
        log_level_from_config = self._get_log_level_from_config()
        log_level_from_mconfig = self._get_log_level_from_mconfig()

        if log_level_from_config is not None:
            log_level = log_level_from_config
        elif log_level_from_mconfig is not None:
            log_level = log_level_from_mconfig
        else:
            logging.warning(
                'logging level is not specified in either yml config '
                'or mconfig, defaulting to INFO',
            )
            log_level = LogLevel.Value('INFO')
        self._set_log_level(log_level)

    @staticmethod
    def _set_log_level(proto_level: int):
        """Set log level based on proto-enum level

        Args:
            proto_level (int): proto enum defined in common.proto
        """
        if proto_level == LogLevel.Value('DEBUG'):
            level = logging.DEBUG
        elif proto_level == LogLevel.Value('INFO'):
            level = logging.INFO
        elif proto_level == LogLevel.Value('WARNING'):
            level = logging.WARNING
        elif proto_level == LogLevel.Value('ERROR'):
            level = logging.ERROR
        elif proto_level == LogLevel.Value('FATAL'):
            level = logging.FATAL
        else:
            logging.error(
                'Unknown logging level: %d, defaulting to INFO',
                proto_level,
            )
            level = logging.INFO

        logging.info(
            "Setting logging level to %s",
            logging.getLevelName(level),
        )
        logger = logging.getLogger('')
        logger.setLevel(level)

    def _register_signal_handlers(self):
        """Register signal handlers

        Right now we just exit on SIGINT/SIGTERM/SIGQUIT
        """
        for signame in ['SIGINT', 'SIGTERM', 'SIGQUIT']:
            self._loop.add_signal_handler(
                getattr(signal, signame),
                functools.partial(self._stop, signame),
            )

        def _signal_handler():
            logging.info('Handling SIGHUP...')
            faulthandler.dump_traceback()
        self._loop.add_signal_handler(
            signal.SIGHUP, functools.partial(_signal_handler),
        )

    def GetServiceInfo(self, request, context):
        """Return the service info (name, version, state, meta, etc.)"""
        service_info = ServiceInfo(
            name=self._name,
            version=self._version,
            state=self._state,
            health=self._health,
            start_time_secs=self._start_time,
        )
        if self._get_status_callback is not None:
            status = self._get_status_callback()
            try:
                service_info.status.meta.update(status)
            except (TypeError, ValueError) as exp:
                logging.error("Error getting service status: %s", exp)
        return service_info

    def StopService(self, request, context):
        """Handle request to stop the service"""
        logging.info("Request to stop service.")
        self._loop.call_soon_threadsafe(self._stop, 'RPC')
        return Void()

    def GetMetrics(self, request, context):
        """
        Collects timeseries samples from prometheus python client on this
        process
        """
        metrics = MetricsContainer()
        metrics.family.extend(get_metrics())
        return metrics

    def SetLogLevel(self, request, context):
        """Handle request to set the log level"""
        self._set_log_level(request.level)
        return Void()

    def SetLogVerbosity(self, request, context):
        pass  # Not Implemented

    def ReloadServiceConfig(self, request, context):
        """Handle request to reload the service config file"""
        self.reload_config()
        return ReloadConfigResponse(result=ReloadConfigResponse.RELOAD_SUCCESS)

    def GetOperationalStates(self, request, context):
        """Return the  operational states of devices managed by this service."""
        res = GetOperationalStatesResponse()
        if self._get_operational_states_cb is not None:
            states = self._get_operational_states_cb()
            res.states.extend(states)
        return res


def get_service303_client(service_name: str, location: str) \
        -> Optional[Service303Stub]:
    """
    Return a grpc client attached to the given service
    name and location.
    Example Use: client = get_service303_client("state", ServiceRegistry.LOCAL)
    """
    try:
        chan = ServiceRegistry.get_rpc_channel(
            service_name,
            location,
        )
        return Service303Stub(chan)
    except ValueError:
        # Service can't be contacted
        logging.error('Failed to get RPC channel to %s', service_name)
        return None
