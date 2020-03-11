"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import sys
from distutils.util import strtobool

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


def integ_test(gateway_host=None, test_host=None, trf_host=None,
               destroy_vm="False"):
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
    """

    destroy_vm = bool(strtobool(destroy_vm))

    # Setup the gateway: use the provided gateway if given, else default to the
    # vagrant machine
    if not gateway_host:
        vagrant_setup("cwag", destroy_vm)
    else:
        ansible_setup(gateway_host, "cwag", "cwag_dev.yml")

    execute(_run_unit_tests)
    execute(_set_cwag_configs)
    cwag_host_to_mac = execute(_get_br_mac, CWAG_BR_NAME)
    host = env.hosts[0]
    cwag_br_mac = cwag_host_to_mac[host]

    # Transfer built images from local machine to CWAG host
    if gateway_host:
        execute(_transfer_docker_images)
    else:
        execute(_stop_gateway)
        execute(_build_gateway)
    execute(_run_gateway)

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
        vagrant_setup("cwag", destroy_vm)
    else:
        ansible_setup(gateway_host, "cwag", "cwag_dev.yml")
    execute(_set_cwag_networking, cwag_test_br_mac)

    # Start tests
    if not test_host:
        vagrant_setup("cwag_test", destroy_vm)
    else:
        ansible_setup(test_host, "cwag_test", "cwag_test.yml")
    execute(_start_ue_simulator)
    execute(_set_cwag_test_networking, cwag_br_mac)
    execute(_run_integ_tests, test_host, trf_host)


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


def _set_cwag_configs():
    """ Set the necessary config overrides """

    with cd(CWAG_INTEG_ROOT):
        sudo('mkdir -p /var/opt/magma')
        sudo('mkdir -p /var/opt/magma/configs')
        sudo('cp gateway.mconfig /var/opt/magma/configs/')


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


def _start_ue_simulator():
    """ Starts the UE Sim Service """
    with cd(CWAG_ROOT + '/services/uesim/uesim'):
        run('tmux new -d \'go run main.go\'')


def _start_trfserver():
    """ Starts the traffic gen server"""
    run('nohup iperf3 -s -B %s > /dev/null &' % TRF_SERVER_IP, pty=False)


def _run_unit_tests():
    """ Run the cwag unit tests """

    with cd(CWAG_ROOT):
        run('make test')


def _run_integ_tests(test_host, trf_host):
    """ Run the integration tests """
    with cd(CWAG_INTEG_ROOT):
        result = run('make integ_test', warn_only=True)

    if not test_host and not trf_host:
        # Clean up only for now when running locally
        execute(_clean_up)
    if result.return_code == 0:
        print("Integration Test Passed!")
        sys.exit(0)
    else:
        print("Integration Test returned ", result.return_code)
        sys.exit(result.return_code)


def _clean_up():
    # already in cwag test vm at this point
    # Kill uesim service
    run('pkill go', warn_only=True)

    with lcd(LTE_AGW_ROOT):
        vagrant_setup("magma_trfserver", False)
        run('pkill iperf3 > /dev/null &', pty=False, warn_only=True)
