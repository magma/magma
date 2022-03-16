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

import copy
import ctypes
import ipaddress
import os
import shlex
import socket
import subprocess
import threading

import iperf3
import pyroute2
import s1ap_types
from util.traffic_messages import (
    TrafficMessage,
    TrafficRequest,
    TrafficRequestType,
    TrafficResponseType,
    TrafficTestInstance,
)

# Tests shouldn't take longer than a few minutes
TRAFFIC_TEST_TIMEOUT_SEC = 180

"""
Using TrafficUtil
=================

TrafficUtil is designed to have one main entry point: generate_traffic_test.
This function sets up the necessary legwork to configuring the trfgen framework
in the S1AP tester and generating a TrafficTest object that represents the
configurations and constraints of the traffic that is to be generated.

Once generated, the TrafficTest object can be run -- either directly with the
start() function or as a context, using the `with' keyword. The wait() function
gives the tester the option to wait on the test completing before continuing.

Essentially, TrafficUtil is just a bridge for packaging together the parameters
of a given test. Once packaged, the actual testing is done via the TrafficTest
API.
"""


class TrafficUtil(object):
    """ Utility wrapper for tests requiring traffic generation """

    # Trfgen library setup
    _trfgen_lib_name = "libtrfgen.so"
    _trfgen_tests = ()

    # Traffic setup
    _remote_ip = ipaddress.IPv4Address("192.168.129.42")

    def __init__(self):
        """ Initialize the trfgen library and its callbacks """
        # _test_lib is the private variable containing the ctypes reference to
        # the trfgen library.
        self._test_lib = None
        self._init_lib()

        # _config_test is the private variable containing the ctypes reference
        # to the trfgen_configure_test() function in trfgen. This function is
        # called to inform the S1AP tester of the parameters of a test suite,
        # and is used to pass along configuration options to the tester.
        self._config_test = None
        self._setup_configure_test()

        # _start_test is the private variable containing the ctypes reference
        # to the trfgen_start_test() function in trfgen. This function is
        # called to begin a single trfgen instance on a given address, using
        # the predefined configuration options set with configure_test().
        self._start_test = None
        self._setup_start_test()

        # We collect references to the data we pass into ctypes to prevent
        # Python's garbage collection system from coming in and cleaning up the
        # memory used, which can result in unspecified behavior.
        self._data = ()

        # Configuration for triggering shell commands in TRF server VM
        self._cmd_data = {
            "user": "vagrant",
            "host": "192.168.60.144",
            "password": "vagrant",
            "command": "test",
        }

        self._command = (
            "sshpass -p {password} ssh "
            "-o UserKnownHostsFile=/dev/null "
            "-o StrictHostKeyChecking=no "
            "{user}@{host} {command}"
        )

    def exec_command(self, command):
        """
        Run a command remotely on magma_trfserver VM.

        Args:
            command: command (str) to be executed on remote host
            e.g. 'sed -i \'s/str1/str2/g\' /usr/local/bin/traffic_server.py'

        """
        data = self._cmd_data
        data["command"] = '"' + command + '"'
        param_list = shlex.split(self._command.format(**data))
        return subprocess.call(
            param_list,
            shell=False,
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
        )

    def update_dl_route(self, ue_ip_block):
        """ Update downlink route in TRF server """
        ret_code = self.exec_command(
            "sudo ip route flush via 192.168.129.1 && sudo ip route "
            "replace " + ue_ip_block + " via 192.168.129.1 dev eth2",
        )
        if ret_code != 0:
            return False
        return True

    def _init_lib(self):
        """ Initialize the trfgen library by loading in binary compiled from C
        """
        lib_path = os.environ["S1AP_TESTER_ROOT"]
        lib = os.path.join(lib_path, "bin", TrafficUtil._trfgen_lib_name)
        os.chdir(lib_path)
        self._test_lib = ctypes.cdll.LoadLibrary(lib)
        self._test_lib.trfgen_init()

    def _setup_configure_test(self):
        """ Set up the call to trfgen_configure_test

        The function prototype is:
            void trfgen_configure_test(int test_id, struct_test test_parms)

        This function call caches the test configurations specified in the
        struct to be called upon and run from the S1AP tester binary.
        """
        self._config_test = self._test_lib.trfgen_configure_test
        self._config_test.restype = None
        self._config_test.argtypes = (ctypes.c_int32, s1ap_types.struct_test)

    def _setup_start_test(self):
        """ Set up the call to trfgen_start_test

        The function prototype is:
            void trfgen_start_test(
                int test_id, char *host_ip, char *bind_ip, char *host_port)

        This function provides a configuration ID and bind address to the S1AP
        tester for it to start a trfgen test. This function returns practically
        immediately, as the iperf3 process is called on a separate fork.
        """
        self._start_test = self._test_lib.trfgen_start_test
        self._start_test.restype = None
        self._start_test.argtypes = (
            ctypes.c_int,
            ctypes.c_char_p,
            ctypes.c_char_p,
            ctypes.c_char_p,
        )

    def cleanup(self):
        """ Cleanup the dll loaded explicitly so the next run doesn't reuse the
        same globals as ctypes LoadLibrary uses dlopen under the covers """
        # self._test_lib.dlclose(self._test_lib._handle)
        self._test_lib = None
        self._data = None

    def configure_test(self, is_uplink, duration, is_udp):
        """ Returns the test configuration index for the configurations
        provided. This is the index that is in the trfgen internal state. If a
        configuration is new, will attempt to create a new one in trfgen

        Args:
            is_uplink (bool): uplink if True, downlink if False
            duration (int): test duration, in seconds
            is_udp (bool): use UDP if True, TCP if False

        Returns: an int, the index of the test configuration in trfgen, a.k.a.
            the test_id

        Raises MemoryError if return test index would exceed
            s1ap_types.MAX_TEST_CFG
        """
        test = s1ap_types.struct_test()
        test.trfgen_type = (
            s1ap_types.trfgen_type.CLIENT.value
            if is_uplink
            else s1ap_types.trfgen_type.SERVER.value
        )
        test.traffic_type = (
            s1ap_types.trf_type.UDP.value
            if is_udp
            else s1ap_types.trf_type.TCP.value
        )
        test.duration = duration
        test.server_timeout = duration

        # First we see if this test has already been configured. If so, just
        # reuse that configuration
        for t in self._trfgen_tests:
            if (
                t.trfgen_type == test.trfgen_type
                and t.traffic_type == test.traffic_type
                and t.duration == test.duration
                and t.server_timeout == test.server_timeout
            ):
                return t.test_id

        # Otherwise, we just create the new test
        if s1ap_types.MAX_TEST_CFG >= len(self._trfgen_tests):
            test.test_id = len(self._trfgen_tests)
            self._trfgen_tests += (test,)
            self._config_test(test.test_id, test)
            return test.test_id

        # If we get here, then we've reached the limit on the number of tests
        # that we can configure, so send an error. Eventually, come up with an
        # eviction scheme
        raise MemoryError(
            "Reached limit on number of configurable tests: %d"
            % s1ap_types.MAX_TEST_CFG,
        )

    def generate_traffic_test(
        self, ips, is_uplink=False, duration=120, is_udp=False,
    ):
        """ Creates a TrafficTest object for the given UE IPs and test type

        Args:
            ips (list(ipaddress.ip_address)): the IP addresses of the UEs to
                which to connect
            is_uplink (bool): whether to do an uplink test. Defaults to False
            duration (int): duration, in seconds, of the test. Defaults to 120
            is_udp (bool): whether to use UDP. If False, uses TCP. Defaults to
                False

        Returns: a TrafficTest object, which is used to interact with the
            trfgen test
        """
        test_id = self.configure_test(is_uplink, duration, is_udp)
        instances = tuple(
            TrafficTestInstance(is_uplink, is_udp, duration, ip, 0)
            for ip in ips
        )
        return TrafficTest(self._start_test, instances, (test_id,) * len(ips))


class TrafficTest(object):
    ''' Class for representing a trfgen test with which to interact

    This is the class that directly interacts with the TrafficTestServer via a
    socketed connection, when the test starts (i.e. the "client" for the
    "server").
    '''

    _alias_counter = 0
    _alias_lock = threading.Lock()
    _iproute = pyroute2.IPRoute()
    _net_iface = 'eth2'
    _port = 7000
    _port_lock = threading.Lock()

    # Remote iperf3 superserver (IP, port) tuple. Port 62462 is chosen because
    # 'MAGMA' translates to 62462 on a 12-key phone pad
    _remote_server = ('192.168.60.144', 62462)

    def __init__(self, test_runner, instances, test_ids):
        ''' Creates a new TrafficTest object for running the test instance(s)
        with the associated test_ids

        Ports will be assigned when the test is run by communicating with the
        test server responsible for iperf3 test servers

        Args:
            test_runner: the ctypes hook into the traffic gen trfgen_start_test
                function
            instances (list(TrafficTestInstance)): the instances to run
            test_ids (list(int)): the associated trfgen test configuration
                indices; must be the same length as instances
        '''
        assert len(instances) is len(test_ids)
        self._done = threading.Event()
        self._instances = tuple(instances)
        self._results = None  # Cached list(iperf3.TestResult) objects
        self._runner = test_runner
        self._test_ids = tuple(test_ids)
        self._test_lock = threading.RLock()  # Provide mutex between tests

    def __enter__(self):
        ''' Starts execution of the test '''
        self.start()
        return self

    def __exit__(self, *_):
        ''' Waits for test to end '''
        self.wait()

    @staticmethod
    def _get_port():
        ''' Returns the next port for testing '''
        with TrafficTest._port_lock:
            TrafficTest._port += 1
            return TrafficTest._port

    @staticmethod
    def _iface_up(ip):
        ''' Brings up an iface for the given IP

        Args:
            ip (ipaddress.ip_address): the IP address to use for bringing up
                the iface

        Returns the iface name with alias that was brought up
        '''
        # Generate a unique alias
        with TrafficTest._alias_lock:
            TrafficTest._alias_counter += 1
            net_iface = TrafficTest._net_iface
            alias = TrafficTest._alias_counter
        net_alias = '%s:UE%d' % (net_iface, alias)

        # Bring up the iface alias
        net_iface_index = TrafficTest._iproute.link_lookup(
            ifname=TrafficTest._net_iface,
        )[0]
        TrafficTest._iproute.addr(
            'add', index=net_iface_index, label=net_alias, address=ip.exploded,
        )

        return net_alias

    @staticmethod
    def _network_from_ip(ip, mask_len):
        ''' Returns the ipaddress.ip_network with the given mask that contains
        the given IP address

        Args:
            ip (ipaddress.ip_address): the IP address for which we want to find
                the network
            mask_len (int): the number of bits to mask

        Returns an ipaddress.ip_network; works agnostic to IPv4 or IPv6
        '''
        # Convert to int to make bit shifting easier
        ip_int = int.from_bytes(ip.packed, 'big')  # Packed is big-endian
        ip_masked = ipaddress.ip_address(ip_int >> mask_len << mask_len)

        # Compute the appropriate prefix length
        prefix_len = ip.max_prefixlen - mask_len

        return ipaddress.ip_network('%s/%d' % (ip_masked.exploded, prefix_len))

    def _run(self):
        ''' Run the traffic test

        Sets up traffic test with remote traffic server and local ifaces, then
        runs the runner hook into the trfgen binary and collects the results to
        cache

        Will block until the test ends
        '''
        # Create a snapshot of the test's states, in case they get changed or
        # wiped in a later operation. Basically, render tests immune to later
        # operations after the test has started.
        with self._test_lock:
            instances = copy.deepcopy(self._instances)
            test_ids = copy.deepcopy(self._test_ids)

        try:
            # Set up sockets and associated streams
            sc = socket.create_connection(self._remote_server)
            sc_in = sc.makefile('rb')
            sc_out = sc.makefile('wb')

            # Flush all the addresses left by previous failed tests
            net_iface_index = TrafficTest._iproute.link_lookup(
                ifname=TrafficTest._net_iface,
            )[0]
            for instance in instances:
                TrafficTest._iproute.flush_addr(
                    index=net_iface_index,
                    address=instance.ip.exploded,
                )

            # Set up network ifaces and get UL port assignments for DL
            aliases = ()
            for instance in instances:
                aliases += (TrafficTest._iface_up(instance.ip),)
                if not instance.is_uplink:
                    # Assign a local port for the downlink UE server
                    instance.port = TrafficTest._get_port()

            # Create and send TEST message
            msg = TrafficRequest(
                TrafficRequestType.TEST, payload=instances,
            )
            msg.send(sc_out)

            # Receive SERVER message and update test instances
            msg = TrafficMessage.recv(sc_in)
            assert msg.message is TrafficResponseType.SERVER
            r_id = msg.id  # Remote server test identifier
            server_instances = msg.payload  # (TrafficServerInstance, ...)

            # Locally keep references to arguments passed into trfgen
            args = [None] * len(instances)

            # Post-SERVER, pre-START logic
            for i in range(len(instances)):
                instance = instances[i]
                server_instance = server_instances[i]

                # Add ip network route
                net_iface_index = TrafficTest._iproute.link_lookup(
                    ifname=TrafficTest._net_iface,
                )[0]
                server_instance_network = \
                    TrafficTest._network_from_ip(server_instance.ip, 8)
                TrafficTest._iproute.route(
                    'replace', dst=server_instance_network.exploded,
                    iif=net_iface_index, oif=net_iface_index, scope='link',
                )

                # Add arp table entry
                os.system(
                    '/usr/sbin/arp -s %s %s' % (
                        server_instance.ip.exploded, server_instance.mac,
                    ),
                )

                if instance.is_uplink:
                    # Port should be the port of the remote for uplink
                    instance.port = server_instance.port
                else:
                    args[i] = self._run_test(
                        test_ids[i], server_instance.ip, instance.ip,
                        instance.port,
                    )

            # Send START for the given r_id
            msg = TrafficRequest(
                TrafficRequestType.START, identifier=r_id,
            )
            msg.send(sc_out)

            # Wait for STARTED response
            msg = TrafficMessage.recv(sc_in)
            assert msg.message is TrafficResponseType.STARTED
            assert msg.id == r_id

            # Post-STARTED, pre-RESULTS logic
            for i in range(len(instances)):
                instance = instances[i]
                if instance.is_uplink:
                    args[i] = self._run_test(
                        test_ids[i], server_instances[i].ip, instance.ip,
                        server_instances[i].port,
                    )

            # Wait for RESULTS message
            msg = TrafficMessage.recv(sc_in)
            assert msg.message is TrafficResponseType.RESULTS
            assert msg.id == r_id
            results = msg.payload

            # Signal to end connection
            msg = TrafficRequest(TrafficRequestType.EXIT)
            msg.send(sc_out)

            # Close out network ifaces
            net_iface_index = TrafficTest._iproute.link_lookup(
                ifname=TrafficTest._net_iface,
            )[0]
            # For some reason the first call to flush this address flushes all
            # the addresses brought up during testing. But subsequent flushes
            # do nothing if the address doesn't exist
            for instance in instances:
                TrafficTest._iproute.flush_addr(
                    index=net_iface_index,
                    address=instance.ip.exploded,
                )
            # Do socket cleanup
            sc_in.close()
            sc_out.close()
            sc.shutdown(socket.SHUT_RDWR)  # Ensures safe socket closure
            sc.close()

            # Cache results after cleanup
            with self._test_lock:
                self._results = results
        finally:
            # Signal that we're done
            self._done.set()

    def _run_test(self, test_id, host_ip, ue_ip, port):
        ''' Run the test at the given index by calling the test runner on the
        test parameters for the instance at the given index and port

        Args:
            test_id (int): the trfgen configuration index to use
            host_ip (ipaddress.ip_address): the remote iperf3 server's IP
                address [-c, for uplink]
            ue_ip (ipaddress.ip_address): the local UE's IP address to which to
                bind [-B]
            port (int): the UE's port (downlink) or the remote server's port
                (uplink) [-p]

        Returns the raw arguments passed into the trfgen binary, for the caller
            to keep track of and avoid garbage collection
        '''
        args = (
            test_id, host_ip.exploded.encode(), ue_ip.exploded.encode(),
            str(port).encode(),
        )
        self._runner(*args)
        return args

    @staticmethod
    def combine(test, *tests):
        ''' Combines TrafficTest objects to produce a single test object that
        will run the parameters given in the tests all at the same time

        All tests in the argument will become unrunnable, as their instances
        will be stripped!

        Args:
            test (TrafficTest): a test, included to force at least one test to
                be passed as an argument
            tests (list(TrafficTest)): any remaining tests to combine

        Return a single TrafficTest that will run all the instances together
        '''
        runner = test._runner

        tests = (test,) + tests
        instances = ()
        test_ids = ()
        for test in tests:
            with test._test_lock:
                instances += test._instances
                test_ids += test._test_ids

                # Now disable the test from later runs
                test._instances = ()
                test._test_ids = ()

        # Create and return the new test
        return TrafficTest(runner, instances, test_ids)

    @property
    def results(self):
        return self._results

    def start(self):
        ''' Start this test by spinning off runner thread '''
        self._done.clear()
        threading.Thread(target=self._run).start()

    def verify(self):
        ''' Verify the results of this test

        Raises a RuntimeError if any tests returned with an error message
        '''
        self.wait()
        with self._test_lock:
            if not isinstance(self.results, tuple):
                raise RuntimeError(
                    'Cached results object is not a tuple : {}'.format(
                        self.results,
                    ),
                )
            for result in self.results:
                if not isinstance(result, iperf3.TestResult):
                    raise RuntimeError(
                        'Cached results are not iperf3.TestResult objects',
                    )
                if result.error:
                    # iPerf dumps out-of-order packet information on stderr,
                    # ignore these while verifying the test results
                    if "OUT OF ORDER" not in result.error:
                        raise RuntimeError(result.error)

    def wait(self):
        ''' Wait for this test to complete '''
        self._done.wait(timeout=TRAFFIC_TEST_TIMEOUT_SEC)
