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
import logging
import os
import queue
import signal
from typing import List

import grpc
import snowflake
from google.protobuf import json_format
from google.protobuf.json_format import MessageToJson
from google.protobuf.struct_pb2 import Struct
from magma.common.rpc_utils import return_void, set_grpc_err
from magma.common.sentry import (
    SEND_TO_ERROR_MONITORING,
    SentryStatus,
    get_sentry_status,
    send_uncaught_errors_to_monitoring,
)
from magma.common.service import MagmaService
from magma.common.service_registry import ServiceRegistry
from magma.common.stateless_agw import (
    check_stateless_agw,
    disable_stateless_agw,
    enable_stateless_agw,
)
from magma.configuration.mconfig_managers import MconfigManager
from magma.magmad.check.network_check import ping, traceroute
from magma.magmad.generic_command.command_executor import CommandExecutor
from magma.magmad.service_manager import ServiceManager
from orc8r.protos import magmad_pb2, magmad_pb2_grpc

enable_sentry_wrapper = get_sentry_status("magmad") == SentryStatus.SEND_SELECTED_ERRORS


class MagmadRpcServicer(magmad_pb2_grpc.MagmadServicer):
    """
    gRPC based server for Magmad.
    """

    def __init__(
        self,
        magma_service: MagmaService,
        services: List[str],
        service_manager: ServiceManager,
        mconfig_manager: MconfigManager,
        command_executor: CommandExecutor,
        loop: asyncio.AbstractEventLoop,
        print_grpc_payload: bool = False,
    ):
        """
        Constructor for the magmad RPC servicer

        Args:
            magma_service:
                MagmaService servicer (service303) which runs on the same
                server as this servicer
            services:
                List of services that magmad manages

            service_manager: ServiceManager instance
            mconfig_manager: MconfigManager instance
            loop: event loop
        """
        self._print_grpc_payload = print_grpc_payload
        self._service_manager = service_manager
        self._services = services
        self._mconfig_manager = mconfig_manager
        self._magma_service = magma_service
        self._command_executor = command_executor
        self._loop = loop

    def add_to_server(self, server):
        """
        Add the servicer to a gRPC server
        """
        magmad_pb2_grpc.add_MagmadServicer_to_server(self, server)

    @return_void
    @send_uncaught_errors_to_monitoring(enable_sentry_wrapper)
    def StartServices(self, _, context):
        """
        Start all magma services
        """
        self._loop.create_task(self._service_manager.start_services())

    @return_void
    @send_uncaught_errors_to_monitoring(enable_sentry_wrapper)
    def StopServices(self, _, context):
        """
        Stop all magma services
        """
        self._loop.create_task(self._service_manager.stop_services())

    @return_void
    @send_uncaught_errors_to_monitoring(enable_sentry_wrapper)
    def Reboot(self, _, context):
        """
        Reboot the gateway device
        """
        async def run_reboot():
            await asyncio.sleep(1)
            os.system('reboot')

        logging.info("Remote reboot triggered! Rebooting gateway...")
        self._loop.create_task(run_reboot())

    @return_void
    @send_uncaught_errors_to_monitoring(enable_sentry_wrapper)
    def RestartServices(self, request, context):
        """
        Restart specified magma services.
        If no services specified, restart all services.
        """
        self._print_grpc(request)

        async def run_restart():
            await asyncio.sleep(1)
            await self._service_manager.restart_services(request.services)

        logging.info("Restarting following services: %s", request.services)
        self._loop.create_task(run_restart())

    @return_void
    @send_uncaught_errors_to_monitoring(enable_sentry_wrapper)
    def SetConfigs(self, request, context):
        """
        Set the stored mconfig to a specific value. Restarts all services
        other than magmad and reloads the magmad mconfig in-place.

        If the gateway configs provided are empty, this will delete the
        managed configurations, reverting the stored mconfigs to the image
        defaults.
        """
        self._print_grpc(request)
        if request.configs_by_key is None or \
                len(request.configs_by_key) == 0:
            self._mconfig_manager.delete_stored_mconfig()
        else:
            # TODO: support streaming mconfig manager impl
            self._mconfig_manager.update_stored_mconfig(
                json_format.MessageToJson(request),
            )

        self._loop.create_task(
            self._service_manager.restart_services(self._services),
        )
        self._magma_service.reload_mconfig()

    @send_uncaught_errors_to_monitoring(enable_sentry_wrapper)
    def GetConfigs(self, _, context):
        # TODO: support streaming mconfig manager impl
        return self._mconfig_manager.load_mconfig()

    @send_uncaught_errors_to_monitoring(enable_sentry_wrapper)
    def RunNetworkTests(self, request, context):
        """
        Execute some specified network commands to check gateway network health
        """
        self._print_grpc(request)
        ping_results = self.__ping_specified_hosts(request.pings)
        traceroute_results = self.__traceroute_specified_hosts(
            request.traceroutes,
        )

        return magmad_pb2.NetworkTestResponse(
            pings=ping_results,
            traceroutes=traceroute_results,
        )

    @send_uncaught_errors_to_monitoring(enable_sentry_wrapper)
    def GetGatewayId(self, _, context):
        """
        Get gateway hardware ID
        """
        return magmad_pb2.GetGatewayIdResponse(
            gateway_id=snowflake.snowflake(),
        )

    @send_uncaught_errors_to_monitoring(enable_sentry_wrapper)
    def GenericCommand(self, request, context):
        """
        Execute generic command. This method will run the command with params
        as specified in the command executor's command table, then return
        the response of the command.
        """
        self._print_grpc(request)
        if 'generic_command_config' not in self._magma_service.config:
            set_grpc_err(
                context,
                grpc.StatusCode.NOT_FOUND,
                'Generic command config not found',
            )
            return magmad_pb2.GenericCommandResponse()

        params = json_format.MessageToDict(request.params)

        # Run the execute command coroutine. Return an error if it times out or
        # if an exception occurs.
        logging.info(
            'Running generic command %s with parameters %s',
            request.command, params,
        )
        future = asyncio.run_coroutine_threadsafe(
            self._command_executor.execute_command(request.command, params),
            self._loop,
        )

        timeout = self._magma_service.config['generic_command_config']\
            .get('timeout_secs', 15)

        response = magmad_pb2.GenericCommandResponse()
        try:
            result = future.result(timeout=timeout)
            logging.debug('Command was successful')
            response.response.MergeFrom(
                json_format.ParseDict(result, Struct()),
            )
        except asyncio.TimeoutError:
            logging.error(
                'Error running command %s! Command timed out',
                request.command,
            )
            future.cancel()
            set_grpc_err(
                context,
                grpc.StatusCode.DEADLINE_EXCEEDED,
                'Command timed out',
            )
        except Exception as e:  # pylint: disable=broad-except
            logging.error(
                'Error running command %s! %s: %s',
                request.command, e.__class__.__name__, e,
                extra=SEND_TO_ERROR_MONITORING,
            )
            set_grpc_err(
                context,
                grpc.StatusCode.UNKNOWN,
                '{}: {}'.format(e.__class__.__name__, str(e)),
            )

        return response

    @send_uncaught_errors_to_monitoring(enable_sentry_wrapper)
    def TailLogs(self, request, context):
        """
        Provides an infinite stream of logs to the client. The client can stop
        the stream by closing the connection.
        """
        self._print_grpc(request)
        if request.service and \
                request.service not in ServiceRegistry.list_services():
            set_grpc_err(
                context,
                grpc.StatusCode.NOT_FOUND,
                'Service {} not found'.format(request.service),
            )
            return

        if not request.service:
            exec_list = ['sudo', 'tail', '-f', '/var/log/syslog']
        else:
            exec_list = [
                'sudo', 'journalctl', '-fu',
                'magma@{}'.format(request.service),
            ]

        logging.debug('Tailing logs')
        log_queue = queue.Queue()

        async def enqueue_log_lines():
            #  https://stackoverflow.com/a/32222971
            proc = await asyncio.create_subprocess_exec(
                *exec_list,
                stdout=asyncio.subprocess.PIPE,
                preexec_fn=os.setsid,
            )
            try:
                while context.is_active():
                    try:
                        line = await asyncio.wait_for(
                            proc.stdout.readline(),
                            timeout=10.0,
                        )
                        log_queue.put(line)
                    except asyncio.TimeoutError:
                        pass
            finally:
                logging.debug('Terminating log stream')
                os.killpg(os.getpgid(proc.pid), signal.SIGTERM)

        self._loop.create_task(enqueue_log_lines())

        while context.is_active():
            try:
                log_line = log_queue.get(block=True, timeout=10.0)
                yield magmad_pb2.LogLine(line=log_line)
            except queue.Empty:
                pass

    @send_uncaught_errors_to_monitoring(enable_sentry_wrapper)
    def CheckStateless(self, _, context):
        """
        Check the stateless mode on AGW
        """
        status = check_stateless_agw()
        logging.debug(
            "AGW mode is %s",
            magmad_pb2.CheckStatelessResponse.AGWMode.Name(status),
        )
        return magmad_pb2.CheckStatelessResponse(agw_mode=status)

    @return_void
    @send_uncaught_errors_to_monitoring(enable_sentry_wrapper)
    def ConfigureStateless(self, request, context):
        """
        Change the stateless mode on AGW, with one of the following:
        enable: Modify AGW config to be stateless
        disable: Modify AGW config to be stateful
        """
        self._print_grpc(request)
        if request.config_cmd == magmad_pb2.ConfigureStatelessRequest.ENABLE:
            logging.info("RPC: config command enable")
            enable_stateless_agw()
        elif request.config_cmd == magmad_pb2.ConfigureStatelessRequest.DISABLE:
            logging.info("RPC: config command disable")
            disable_stateless_agw()

    @staticmethod
    def __ping_specified_hosts(ping_param_protos):
        def create_ping_result_proto(ping_result):
            if ping_result.error:
                return magmad_pb2.PingResult(
                    error=ping_result.error,
                    host_or_ip=ping_result.host_or_ip,
                    num_packets=ping_result.num_packets,
                )
            else:
                return magmad_pb2.PingResult(
                    host_or_ip=ping_result.host_or_ip,
                    num_packets=ping_result.num_packets,
                    packets_transmitted=ping_result.stats.packets_transmitted,
                    packets_received=ping_result.stats.packets_received,
                    avg_response_ms=ping_result.stats.rtt_avg,
                )

        pings_to_exec = [
            ping.PingCommandParams(
                host_or_ip=p.host_or_ip,
                num_packets=p.num_packets,
                timeout_secs=None,
            ) for p in ping_param_protos
        ]
        ping_results = ping.ping(pings_to_exec)
        return map(create_ping_result_proto, ping_results)

    @staticmethod
    def __traceroute_specified_hosts(traceroute_param_protos):
        def create_result_proto(result):
            if result.error:
                return magmad_pb2.TracerouteResult(
                    error=result.error,
                    host_or_ip=result.host_or_ip,
                )
            else:
                return magmad_pb2.TracerouteResult(
                    host_or_ip=result.host_or_ip,
                    hops=create_hop_protos(result.stats.hops),
                )

        def create_hop_protos(hops):
            hop_protos = []
            for hop in hops:
                hop_protos.append(
                    magmad_pb2.TracerouteHop(
                        idx=hop.idx,
                        probes=create_probe_protos(hop.probes),
                    ),
                )
            return hop_protos

        def create_probe_protos(probes):
            return [
                magmad_pb2.TracerouteProbe(
                    hostname=probe.hostname,
                    ip=probe.ip_addr,
                    rtt_ms=probe.rtt_ms,
                ) for probe in probes
            ]

        traceroutes_to_exec = [
            traceroute.TracerouteParams(
                host_or_ip=param.host_or_ip,
                max_hops=param.max_hops,
                bytes_per_packet=param.bytes_per_packet,
            ) for param in traceroute_param_protos
        ]
        traceroute_results = traceroute.traceroute(traceroutes_to_exec)
        return map(create_result_proto, traceroute_results)

    def _print_grpc(self, message):
        if self._print_grpc_payload:
            log_msg = "{} {}".format(
                message.DESCRIPTOR.full_name,
                MessageToJson(message),
            )
            # add indentation
            padding = 2 * ' '
            log_msg = ''.join(
                "{}{}".format(padding, line)
                for line in log_msg.splitlines(True)
            )

            log_msg = "GRPC message:\n{}".format(log_msg)
            logging.info(log_msg)
