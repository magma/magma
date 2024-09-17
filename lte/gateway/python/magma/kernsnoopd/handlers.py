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

import abc
import logging
from ctypes import c_ulong
from functools import lru_cache
from socket import AF_INET, AF_INET6, inet_ntop
from struct import pack
from typing import Tuple

import psutil
from magma.kernsnoopd import metrics

# TASK_COMM_LEN is the string length of binary names that the kernel reports.
# Value should be the same as found in <linux/sched.h>
TASK_COMM_LEN = 16


class EBPFHandler(abc.ABC):
    """
    EBPFHandler class defines the interface for front-end programs
    corresponding to loaded eBPF programs.

    Method handle() must be implemented by a sub-class. Snooper will call the
    handle() method of registered front-end programs periodically.
    """

    def __init__(self, service_registry):
        self._registry = service_registry

        # only the first TASK_COMM_LEN letters of the service name are relevant
        # here as the kernel is only sending those in task->comm
        self._services = [
            s[:TASK_COMM_LEN] for s in service_registry.list_services()
        ]

    @abc.abstractmethod
    def handle(self, bpf) -> None:
        """
        Handle() should serve as the entry point of the front-end program
        performing tasks such as reading metrics collected from the kernel and
        storing them into Prometheus.

        Args:
            bpf: the bcc.BPF instance that was used to load the eBPF program

        Raises:
            NotImplementedError: Implement in sub-class
        """
        raise NotImplementedError()


class ByteCounter(EBPFHandler):
    """
    ByteCounter is the front-end program for ebpf/byte_count.bpf.c
    """

    def __init__(self, service_registry):
        super().__init__(service_registry)
        # Addr is a ctypes array of two 64-bit ints. It is used to hold an IPv6
        # address of int128 type. This type can be converted to tuple and back
        # to make it hashable for caching.
        self.Addr = c_ulong * 2

    @lru_cache(maxsize=1024)
    def _get_cmdline(self, pid: int) -> list:
        """
        _get_cmdline returns the command line arguments that were password to
        process with the given pid. It caches results in an LRU cache to reduce
        cost of reading /proc every time.

        Args:
            pid: process id

        Returns:
            list of strings that make up the command line arguments

        Raises:
            psutil.NoSuchProcess when process with given pid does not exist.
            Process may have already exited.
        """
        return psutil.Process(pid=pid).cmdline()

    @lru_cache(maxsize=1024)
    def _ip_addr_to_str(self, family: int, daddr: Tuple[int, int]) -> str:
        """
        _ip_addr_to_str returns a string representation of an IPv4 or IPv6
        address. It caches results in an LRU cache to reduce cost of conversion

        Args:
            family: socket.AF_INET (v4) or socket.AF_INET6 (v6)
            daddr: For IPv4, uint32 representation of address as the first item
            in a tuple. For IPv6, 16-byte array representation of address.

        Returns:
            String representation of IP address, e.g., '127.0.0.1'
        """
        if family == AF_INET:
            return inet_ntop(AF_INET, pack('I', daddr[0]))
        elif family == AF_INET6:
            # noinspection PyTypeChecker
            return inet_ntop(AF_INET6, self.Addr(*daddr))
        else:
            raise Exception("No valid socket family given!")

    def handle(self, bpf):
        """
        Handle() reads counters from the loaded byte_count program stored as
        a dict in 'dest_counters' with key type key_t and value type counter_t
        defined in ebpf/common.bpf.h

        Args:
            bpf: bcc.BPF object that was used to load eBPF program into kernel
        """
        table = bpf['dest_counters']
        for key, count in table.items():
            d_host = self._ip_addr_to_str(key.family, tuple(key.daddr))
            service_name = None

            try:
                service_name = self._get_source_service(key)
                # TODO: destination service name inference does not work
                # get destination service from host and port
                logging.debug(
                    '%s sent %s bytes to (%s, %s)',
                    service_name,
                    count.value,
                    d_host,
                    key.dport,
                )
                _inc_service_counter(service_name, '', count.value)
            except ValueError:
                # use binary name if source service name was not inferred
                binary_name = service_name or key.comm.decode()
                _inc_linux_counter(binary_name, count.value)

        # clear eBPF counters
        table.clear()

    def _get_source_service(self, key) -> str:
        """
        _get_source_service attempts to get Magma service from command line
        arguments of running process or binary name

        Args:
            key: struct of type key_t from which service name is inferred

        Returns:
            Magma service name inferred from key

        Raises:
            ValueError: Could not infer service name from key
        """

        try:
            # get python service name from command line args
            # e.g. "python3 -m magma.state.main"
            cmdline = self._get_cmdline(key.pid)
            python_service = self._get_service_from_cmdline(cmdline)
            if python_service:
                return python_service
        # key.pid process has exited or was not a Python service
        except (psutil.NoSuchProcess, IndexError):
            binary_name = key.comm.decode()
            if binary_name in self._services:
                # was a non-Python service
                return binary_name
        raise ValueError('Could not infer service name from key %s' % key.comm)

    def _get_service_from_cmdline(self, cmdline):
        if cmdline[2].startswith('magma.'):
            return cmdline[2].split('.')[1]
        return None


def _inc_service_counter(source_service, dest_service, count) -> None:
    """
    _inc_service_counter increments Prometheus byte counters for traffic
    between gateway and cloud Magma services

    Args:
        source_service: traffic source service name used as label
        dest_service: traffic destination service name used as label
        count: byte count to increment
    """
    metrics.MAGMA_BYTES_SENT_TOTAL.labels(
        service_name=source_service,
        dest_service=dest_service,
    ).inc(count)


def _inc_linux_counter(binary_name, count) -> None:
    """
    _inc_linux_counter increments Prometheus byte counters for traffic
    originating from arbitrary linux binaries

    Args:
        binary_name: traffic source binary name used as label
        count: byte count to increment
    """
    metrics.LINUX_BYTES_SENT_TOTAL.labels(binary_name).inc(count)


# ebpf_handlers provides the mapping from ebpf source files
# (e.g. epbf/packet_count.bpf.c) to front-end program class
ebpf_handlers = {
    'byte_count': ByteCounter,
}
