"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from distutils.util import strtobool
from enum import Enum

import sys
from fabric.api import cd, env, execute, lcd, local, put, run, settings, sudo

sys.path.append('../../orc8r')
from tools.fab.hosts import ansible_setup, vagrant_setup

CWAG_ROOT = "$MAGMA_ROOT/cwf/gateway"
CWAG_INTEG_ROOT = "$MAGMA_ROOT/cwf/gateway/integ_tests"
LTE_AGW_ROOT = "../../lte/gateway"

CWAG_TEST_IP = "192.168.128.2"
TRF_SERVER_IP = "192.168.129.42"
TRF_SERVER_SUBNET = "192.168.129.0"
CWAG_BR_NAME = "cwag_br0"
CWAG_TEST_BR_NAME = "cwag_test_br0"


class SubTests(Enum):
    ALL = "integ_test"
    AUTH = "authenticate"
    GX = "gx"
    GY = "gy"
    MULTISESSIONPROXY = "multi_session_proxy"

    @staticmethod
    def list():
        return list(map(lambda t: t.value, SubTests))


def integ_test(gateway_host=None, test_host=None, trf_host=None,
               transfer_images=False, destroy_vm=False, no_build=False,
               tests_to_run="integ_test", skip_unit_tests=False, test_re=None):
    """
    Run the integration tests. This defaults to running on local vagrant
    machines, but can also be pointed to an arbitrary host (e.g. amazon) by
    passing "address:port" as arguments

    gateway_host: The ssh address string of the machine to run the gateway
        services on. Formatted as "host:port". If not specified, defaults to
        the `cwag` vagrant box.

    test_host: The ssh address string of the machine to run the tests on
        on. Formatted as "host:port". If not specified, defaults to the
        `cwag_test` vagrant box.

    trf_host: The ssh address string of the machine to run the tests on
        on. Formatted as "host:port". If not specified, defaults to the
        `magma_trfserver` vagrant box.

    no_build: When set to true, this script will NOT rebuild all docker images.
    """
    try:
        tests_to_run = SubTests(tests_to_run)
    except ValueError:
        print("{} is not a valid value. We support {}".format(
            tests_to_run, SubTests.list()))
        return

    # Setup the gateway: use the provided gateway if given, else default to the
    # vagrant machine
    if not gateway_host:
        vagrant_setup("cwag", destroy_vm)
    else:
        ansible_setup(gateway_host, "cwag", "cwag_dev.yml")

    if not skip_unit_tests:
        execute(_run_unit_tests)

    execute(_set_cwag_configs, "gateway.mconfig")
    cwag_host_to_mac = execute(_get_br_mac, CWAG_BR_NAME)
    host = env.hosts[0]
    cwag_br_mac = cwag_host_to_mac[host]

    # Transfer built images from local machine to CWAG host
    if gateway_host or transfer_images:
        execute(_transfer_docker_images)
    else:
        execute(_stop_gateway)
        if not no_build:
            execute(_build_gateway)
    execute(_run_gateway)
    # Stop not necessary services for this test case
    execute(_stop_docker_services, ["pcrf2", "ocs2"])

    # Setup the trfserver: use the provided trfserver if given, else default to the
    # vagrant machine
    with lcd(LTE_AGW_ROOT):
        if not trf_host:
            vagrant_setup("magma_trfserver", destroy_vm)
        else:
            ansible_setup(trf_host, "trfserver", "magma_trfserver.yml")

    execute(_start_trfserver)

    # Run the tests: use the provided test machine if given, else default to
    # the vagrant machine
    if not test_host:
        vagrant_setup("cwag_test", destroy_vm)
    else:
        ansible_setup(test_host, "cwag_test", "cwag_test.yml")

    cwag_test_host_to_mac = execute(_get_br_mac, CWAG_TEST_BR_NAME)
    host = env.hosts[0]
    cwag_test_br_mac = cwag_test_host_to_mac[host]
    execute(_set_cwag_test_configs)

    # Get back to the gateway vm to setup static arp
    if not gateway_host:
        # We do NOT want to destroy this VM after we just set it up...
        vagrant_setup("cwag", False)
    else:
        ansible_setup(gateway_host, "cwag", "cwag_dev.yml")
    execute(_set_cwag_networking, cwag_test_br_mac)

    # Start main tests - except for multi session proxy
    if not test_host:
        # No, definitely do NOT destroy this VM
        vagrant_setup("cwag_test", False)
    else:
        ansible_setup(test_host, "cwag_test", "cwag_test.yml")
    execute(_start_ue_simulator)
    execute(_set_cwag_test_networking, cwag_br_mac)

    if tests_to_run.value not in [SubTests.MULTISESSIONPROXY.value]:
        execute(_run_integ_tests, test_host, trf_host, tests_to_run, test_re)

    # Setup environment and run test for multi service proxy if required
    if tests_to_run.value in [SubTests.MULTISESSIONPROXY.value, SubTests.ALL.value]:

        # CWAG VM
        if not gateway_host:
            vagrant_setup("cwag", False)
        else:
            ansible_setup(gateway_host, "cwag", "cwag_dev.yml")
        # copy new config and restart the impacted services
        execute(_set_cwag_configs, "gateway.mconfig.multi_session_proxy")
        execute(_restart_docker_services, ["session_proxy", "pcrf", "ocs",
                                           "pcrf2", "ocs2"])

        # CWAG_TEST VM
        if not test_host:
            vagrant_setup("cwag_test", False)
        else:
            ansible_setup(test_host, "cwag_test", "cwag_test.yml")
        execute(
            _run_integ_tests, test_host, trf_host, SubTests.MULTISESSIONPROXY, test_re
        )

    # If we got here means everything work well!!
    if not test_host and not trf_host:
        # Clean up only for now when running locally
        execute(_clean_up)
    print('Integration Test Passed for "{}"!'.format(tests_to_run.value))
    sys.exit(0)


def transfer_service_logs(services="sessiond session_proxy"):
    services = services.strip().split(' ')
    print("Transferring logs for " + str(services))

    # We do NOT want to destroy this VM after we just set it up...
    vagrant_setup("cwag", False)
    with cd(CWAG_ROOT):
        for service in services:
            run("docker logs -t " + service + " 2> " + service + ".log")
            # For vagrant the files should already be in CWAG_ROOT

def _transfer_docker_images():
    output = local("docker images cwf_*", capture=True)
    for line in output.splitlines():
        if not line.startswith('cwf'):
            continue
        line = line.rstrip("\n")
        image = line.split(" ")[0]

        local("docker save -o /tmp/%s.tar %s" % (image, image))
        put("/tmp/%s.tar" % image, "%s.tar" % image)
        local("rm -f /tmp/%s.tar" % image)

        run('docker load -i %s.tar' % image)


def _set_cwag_configs(configfile):
    """ Set the necessary config overrides """

    with cd(CWAG_INTEG_ROOT):
        sudo('mkdir -p /var/opt/magma')
        sudo('mkdir -p /var/opt/magma/configs')
        sudo("cp {} /var/opt/magma/configs/gateway.mconfig".format(configfile))


def _set_cwag_networking(mac):
    sudo('arp -s %s %s' % (CWAG_TEST_IP, mac))


def _get_br_mac(bridge_name):
    mac = run("cat /sys/class/net/%s/address" % bridge_name)
    return mac


def _set_cwag_test_configs():
    """ Set the necessary test configs """

    sudo('mkdir -p /etc/magma')
    # Create empty uesim config
    sudo('touch /etc/magma/uesim.yml')


def _set_cwag_test_networking(mac):
    # Don't error if route already exists
    with settings(warn_only=True):
        sudo('ip route add %s/24 dev %s proto static scope link' %
             (TRF_SERVER_SUBNET, CWAG_TEST_BR_NAME))
    sudo('arp -s %s %s' % (TRF_SERVER_IP, mac))


def _stop_gateway():
    """ Stop the gateway docker images """
    with cd(CWAG_ROOT + '/docker'):
        sudo(' docker-compose'
             ' -f docker-compose.yml'
             ' -f docker-compose.override.yml'
             ' -f docker-compose.integ-test.yml'
             ' down')


def _build_gateway():
    """ Builds the gateway docker images """
    with cd(CWAG_ROOT + '/docker'):
        sudo(' docker-compose'
             ' -f docker-compose.yml'
             ' -f docker-compose.override.yml'
             ' -f docker-compose.nginx.yml'
             ' -f docker-compose.integ-test.yml'
             ' build --parallel')


def _run_gateway():
    """ Runs the gateway's docker images """
    with cd(CWAG_ROOT + '/docker'):
        sudo(' docker-compose'
             ' -f docker-compose.yml'
             ' -f docker-compose.override.yml'
             ' -f docker-compose.integ-test.yml'
             ' up -d ')


def _restart_docker_services(services):
    with cd(CWAG_ROOT + "/docker"):
        sudo(
            " docker-compose"
            " -f docker-compose.yml"
            " -f docker-compose.override.yml"
            " -f docker-compose.nginx.yml"
            " -f docker-compose.integ-test.yml"
            " restart {}".format(" ".join(services))
        )


def _stop_docker_services(services):
    with cd(CWAG_ROOT + "/docker"):
        sudo(
            " docker-compose"
            " -f docker-compose.yml"
            " -f docker-compose.override.yml"
            " -f docker-compose.nginx.yml"
            " -f docker-compose.integ-test.yml"
            " stop {}".format(" ".join(services))
        )


def _start_ue_simulator():
    """ Starts the UE Sim Service """
    with cd(CWAG_ROOT + '/services/uesim/uesim'):
        run('tmux new -d \'go run main.go\'')


def _start_trfserver():
    """ Starts the traffic gen server"""
    run('nohup iperf3 -s --json -B %s > /dev/null &' % TRF_SERVER_IP, pty=False)


def _run_unit_tests():
    """ Run the cwag unit tests """
    with cd(CWAG_ROOT):
        run('make test')


def _run_integ_tests(test_host, trf_host, tests_to_run: SubTests, testRe=None):
    """ Run the integration tests """
    with cd(CWAG_INTEG_ROOT):
        if testRe:
            command = "TESTS=" + testRe + " make " + str(tests_to_run.value)
        else:
            command = "make " + str(tests_to_run.value)
        result = run(command, warn_only=True)
    if result.return_code != 0:
        if not test_host and not trf_host:
            # Clean up only for now when running locally
            execute(_clean_up)
        print("Integration Test returned ", result.return_code)
        sys.exit(result.return_code)


def _clean_up():
    # already in cwag test vm at this point
    # Kill uesim service
    run('pkill go', warn_only=True)
    with lcd(LTE_AGW_ROOT):
        vagrant_setup("magma_trfserver", False)
        run('pkill iperf3 > /dev/null &', pty=False, warn_only=True)
