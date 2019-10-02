"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import sys
from distutils.util import strtobool

from fabric.api import cd, env, execute, lcd, local, run, sudo, warn_only

sys.path.append('../../orc8r')
from tools.fab.hosts import ansible_setup, vagrant_setup

CWAG_ROOT = "$MAGMA_ROOT/cwf/gateway"
CWAG_INTEG_ROOT = "$MAGMA_ROOT/cwf/gateway/integ_tests"
LTE_AGW_ROOT = "../../lte/gateway"


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
    execute(_copy_config)
    execute(_start_gateway)

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

    execute(_start_ue_simulator)
    execute(_run_integ_tests)


def _copy_config():
    """ Copy necessary config overrides to /var/opt/magma/configs """

    with cd(CWAG_INTEG_ROOT):
        sudo('mkdir -p /var/opt/magma')
        sudo('mkdir -p /var/opt/magma/configs')
        sudo('cp gateway.mconfig /var/opt/magma/configs')
        sudo('cp sessiond.yml /var/opt/magma/configs')


def _start_gateway():
    """ Starts the gateway """
    with cd(CWAG_ROOT + '/docker'):
        run(' docker-compose'
            ' -f docker-compose.yml'
            ' -f docker-compose.override.yml'
            ' -f docker-compose.integ-test.yml'
            ' build --parallel')
        run(' docker-compose'
            ' -f docker-compose.yml'
            ' -f docker-compose.override.yml'
            ' -f docker-compose.integ-test.yml'
            ' up -d')


def _start_ue_simulator():
    """ Starts the UE Sim Service """
    with cd(CWAG_ROOT + '/services/uesim/uesim'):
        run('tmux new -d \'go run main.go\'')


def _start_trfserver():
    """ Starts the traffic gen server"""
    run("nohup iperf3 -s > /dev/null &", pty=False)


def _run_unit_tests():
    """ Run the cwag unit tests """

    with cd(CWAG_ROOT):
        run('make test')


def _run_integ_tests():
    """ Run the integration tests """

    with cd(CWAG_INTEG_ROOT):
        run('make integ_test')
