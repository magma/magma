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
# pylint: disable=broad-except
import os
import subprocess
from collections import OrderedDict

import psutil
from magma.common.health.service_state_wrapper import ServiceStateWrapper
from magma.common.service import MagmaService
from magma.magmad.check.network_check import ping
from orc8r.protos.mconfig import mconfigs_pb2
from prometheus_client import Counter, Gauge

POLL_INTERVAL_SECONDS = 100

MAGMAD_PING_STATS = Gauge(
    'magmad_ping_rtt_ms',
    'Gateway ping metrics',
    ['host', 'metric'],
)
CPU_PERCENT = Gauge(
    'cpu_percent',
    'System-wide CPU utilization as a percentage over 1 sec',
)
SWAP_MEMORY_PERCENT = Gauge(
    'swap_memory_percent', 'Percent of memory that can'
    ' be assigned to processes',
)
VIRTUAL_MEMORY_PERCENT = Gauge(
    'virtual_memory_percent',
    'Percent of memory that can be assigned to '
    'processes without the system going into swap',
)
MEM_TOTAL = Gauge('mem_total', 'memory total')
MEM_AVAILABLE = Gauge('mem_available', 'memory available')
MEM_USED = Gauge('mem_used', 'memory used')
MEM_FREE = Gauge('mem_free', 'memory free')
DISK_PERCENT = Gauge(
    'disk_percent',
    'Percent of disk space used for the '
    'volume mounted at root',
)
BYTES_SENT = Gauge('bytes_sent', 'System-wide network I/O bytes sent')
BYTES_RECEIVED = Gauge(
    'bytes_received',
    'System-wide network I/O bytes received',
)
TEMPERATURE = Gauge(
    'temperature', 'Temperature readings from system sensors',
    ['sensor'],
)
CHECKIN_STATUS = Gauge(
    'checkin_status',
    '1 for checkin success, and 0 for failure',
)
BOOTSTRAP_EXCEPTION = Counter(
    'bootstrap_exception',
    'Count for exceptions raised by bootstrapper',
    ['cause'],
)
UNEXPECTED_SERVICE_RESTARTS = Counter(
    'unexpected_service_restarts',
    'Count of unexpected restarts',
    ['service_name'],
)
UNATTENDED_UPGRADE_STATUS = Gauge(
    'unattended_upgrade_status',
    'Unattended Upgrade update status'
    '1 for active, 0 for inactive',
)


SERVICE_RESTART_STATUS = Gauge(
    'service_restart_status',
    'Count of service restarts',
    ['service_name', 'status'],
)


SERVICE_CPU_PERCENTAGE = Gauge(
    'service_cpu_percentage',
    'Service CPU Percentage',
    ['service_name'],
)


SERVICE_MEMORY_USAGE = Gauge(
    'service_memory_usage',
    'Service Memory Usage',
    ['service_name'],
)


SERVICE_MEMORY_PERCENTAGE = Gauge(
    'service_memory_percentage',
    'Service Memory Percentage',
    ['service_name'],
)


def _get_ping_params(config):
    ping_params = []
    if 'ping_config' in config and 'hosts' in config['ping_config']:
        ping_params = [
            ping.PingCommandParams(
                host,
                config['ping_config']['num_packets'],
                config['ping_config']['timeout_secs'],
            ) for host in config['ping_config']['hosts']
        ]
    return ping_params


@asyncio.coroutine
def metrics_collection_loop(service_config, loop=None):
    if 'network_monitor_config' not in service_config:
        return

    config = service_config['network_monitor_config']
    ping_params = _get_ping_params(config)

    while True:
        logging.debug("Running metrics collections loop")
        if len(ping_params):
            yield from _collect_ping_metrics(ping_params, loop=loop)
        yield from _collect_load_metrics()
        yield from _collect_service_restart_stats()
        yield from _collect_service_metrics()
        yield from asyncio.sleep(int(config['sampling_period']))


@asyncio.coroutine
def _collect_service_restart_stats():
    """
    Collect the success and failure restarts for services
    """
    try:
        service_dict = ServiceStateWrapper().get_all_services_status()
    except Exception as e:
        logging.error("Could not fetch service status: %s", e)
        return
    for service_name, status in service_dict.items():
        SERVICE_RESTART_STATUS.labels(
            service_name=service_name,
            status="Failure",
        ).set(status.num_fail_exits)
        SERVICE_RESTART_STATUS.labels(
            service_name=service_name,
            status="Success",
        ).set(status.num_clean_exits)


@asyncio.coroutine
def _collect_load_metrics():
    CPU_PERCENT.set(psutil.cpu_percent(interval=1))

    SWAP_MEMORY_PERCENT.set(psutil.swap_memory().percent)

    mem = psutil.virtual_memory()
    VIRTUAL_MEMORY_PERCENT.set(mem.percent)
    MEM_TOTAL.set(mem.total)
    MEM_AVAILABLE.set(mem.available)
    MEM_USED.set(mem.used)
    MEM_FREE.set(mem.free)

    DISK_PERCENT.set(psutil.disk_usage('/').percent)

    network_io = psutil.net_io_counters()
    BYTES_SENT.set(network_io.bytes_sent)
    BYTES_RECEIVED.set(network_io.bytes_recv)

    # sensors may not exist on all platforms, or error out.
    # try/except to avoid error out in metrics_collection_loop
    try:
        for sensor, values in psutil.sensors_temperatures().items():
            for idx, value in enumerate(values):
                TEMPERATURE.labels(
                    sensor='%s_%d' % (sensor, idx),
                ).set(value.current)
    except OSError as ex:
        logging.warning("sensors_temperatures error: %s", ex)


@asyncio.coroutine
def _collect_ping_metrics(ping_params, loop=None):
    ping_results = yield from ping.ping_async(ping_params, loop=loop)
    ping_results_list = list(ping_results)

    def extract_metrics(ping_stats):
        # metric name: (value, gauge method name)
        return OrderedDict([
            ('rtt_ms', (ping_stats.rtt_avg, 'set')),
            ('packets_sent', (ping_stats.packets_transmitted, 'inc')),
            (
                'packets_lost',
                (
                    ping_stats.packets_transmitted - ping_stats.packets_received,
                    'inc',
                ),
            ),
        ])

    for param, result in zip(ping_params, ping_results_list):
        if result.error:
            logging.debug(
                'Failed to ping %s with error: %s',
                param.host_or_ip, result.error,
            )
        else:
            host = param.host_or_ip
            metrics = extract_metrics(result.stats)
            for metric, value in metrics.items():
                label = MAGMAD_PING_STATS.labels(host=host, metric=metric)
                getattr(label, value[1])(value[0])

            logging.debug(
                'Pinged %s with %d packet(s). Average RTT ms: %s',
                result.host_or_ip, result.num_packets, result.stats.rtt_avg,
            )
    return ping_results_list


@asyncio.coroutine
def monitor_unattended_upgrade_status():
    """
    Call to poll the unattended upgrade status and set the corresponding metric
    """
    while True:
        status = 0
        auto_upgrade_file_name = '/etc/apt/apt.conf.d/20auto-upgrades'
        if os.path.isfile(auto_upgrade_file_name):
            with open(auto_upgrade_file_name, encoding='utf-8') as auto_upgrade_file:
                for line in auto_upgrade_file:
                    package_name, flag = line.strip().strip(';').split()
                    if package_name == "APT::Periodic::Unattended-Upgrade":
                        if flag == '"1"':
                            status = 1
                        break
        logging.debug('Unattended upgrade status is %d', status)
        UNATTENDED_UPGRADE_STATUS.set(status)
        yield from asyncio.sleep(POLL_INTERVAL_SECONDS)


@asyncio.coroutine
def _collect_service_metrics():
    config = MagmaService('magmad', mconfigs_pb2.MagmaD()).config
    magma_services = ["magma@" + service for service in config['magma_services']]
    non_magma_services = ["sctpd", "openvswitch-switch"]
    for service in magma_services + non_magma_services:
        cmd = ["systemctl", "show", service, "--property=MainPID,MemoryCurrent,MemoryAccounting,MemoryLimit"]
        # TODO(@wallyrb): Move away from subprocess and use psystemd
        output = subprocess.check_output(cmd)
        output_str = str(output, "utf-8").strip().replace("MainPID=", "").replace("MemoryCurrent=", "").replace("MemoryAccounting=", "").replace("MemoryLimit=", "")
        properties = output_str.split("\n")
        pid = int(properties[0])
        memory = properties[1]
        memory_accounting = properties[2]
        memory_limit = properties[3]

        if pid != 0:
            try:
                p = psutil.Process(pid=pid)
                cpu_percentage = p.cpu_percent(interval=1)
            except psutil.NoSuchProcess:
                logging.warning("When collecting CPU usage for service %s: Process with PID %d no longer exists.", service, pid)
                continue
            else:
                SERVICE_CPU_PERCENTAGE.labels(
                    service_name=service,
                ).set(cpu_percentage)

        if not memory.isnumeric():
            continue

        if memory_accounting == "yes":
            SERVICE_MEMORY_USAGE.labels(
                service_name=service,
            ).set(int(memory))

        if memory_limit.isnumeric():
            SERVICE_MEMORY_PERCENTAGE.labels(
                service_name=service,
            ).set(int(memory) / int(memory_limit))
