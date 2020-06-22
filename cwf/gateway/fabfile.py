"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
from enum import Enum

import sys
from fabric.api import (cd, env, execute, lcd, local, put, run, settings,
                        sudo, shell_env)
from fabric.contrib import files
sys.path.append('../../orc8r')

from tools.fab.hosts import ansible_setup, vagrant_setup

CWAG_ROOT = "$MAGMA_ROOT/cwf/gateway"
CWAG_INTEG_ROOT = "$MAGMA_ROOT/cwf/gateway/integ_tests"
LTE_AGW_ROOT = "../../lte/gateway"

CWAG_IP = "192.168.70.101"
CWAG_TEST_IP = "192.168.128.2"
TRF_SERVER_IP = "192.168.129.42"
TRF_SERVER_SUBNET = "192.168.129.0"
CWAG_BR_NAME = "cwag_br0"
CWAG_TEST_BR_NAME = "cwag_test_br0"


class SubTests(Enum):
    ALL = "all"
    AUTH = "authenticate"
    GX = "gx"
    GY = "gy"
    QOS = "qos"
    MULTISESSIONPROXY = "multi_session_proxy"

    @staticmethod
    def list():
        return list(map(lambda t: t.value, SubTests))


def integ_test(gateway_host=None, test_host=None, trf_host=None,
               transfer_images=False, destroy_vm=False, no_build=False,
               tests_to_run="all", skip_unit_tests=False, test_re=None,
               run_tests=True):
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
    _switch_to_vm(gateway_host, "cwag", "cwag_dev.yml", destroy_vm)

    # We will direct coredumps to be placed in this directory
    # Clean up before every run
    if files.exists("/var/opt/magma/cores/"):
        run("sudo rm /var/opt/magma/cores/*", warn_only=True)
    else:
        run("sudo mkdir -p /var/opt/magma/cores", warn_only=True)

    if not skip_unit_tests:
        execute(_run_unit_tests)

    execute(_set_cwag_configs, "gateway.mconfig")
    execute(_add_networkhost_docker)
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

    # Setup the trfserver: use the provided trfserver if given, else default to
    # the vagrant machine
    with lcd(LTE_AGW_ROOT):
        _switch_to_vm(gateway_host, "magma_trfserver",
                      "magma_trfserver.yml", destroy_vm)

    execute(_start_trfserver)

    # Run the tests: use the provided test machine if given, else default to
    # the vagrant machine
    _switch_to_vm(gateway_host, "cwag_test", "cwag_test.yml", destroy_vm)

    cwag_test_host_to_mac = execute(_get_br_mac, CWAG_TEST_BR_NAME)
    host = env.hosts[0]
    cwag_test_br_mac = cwag_test_host_to_mac[host]
    execute(_set_cwag_test_configs)
    execute(_start_ipfix_controller)

    # Get back to the gateway vm to setup static arp
    _switch_to_vm_no_destroy(gateway_host, "cwag", "cwag_dev.yml")
    execute(_set_cwag_networking, cwag_test_br_mac)
    execute(_check_docker_services)

    _switch_to_vm_no_destroy(gateway_host, "cwag_test", "cwag_test.yml")
    execute(_start_ue_simulator)
    execute(_set_cwag_test_networking, cwag_br_mac)

    if run_tests == "False":
        execute(_add_docker_host_remote_network_envvar)
        print("run_test was set to false. Test will not be run\n"
              "You can now run the tests manually from cwag_test")
        sys.exit(0)

    if tests_to_run.value not in SubTests.MULTISESSIONPROXY.value:
        execute(_run_integ_tests, test_host, trf_host, tests_to_run, test_re)

    # Setup environment for multi service proxy if required
    if tests_to_run.value in (SubTests.ALL.value,
                              SubTests.MULTISESSIONPROXY.value):
        _switch_to_vm_no_destroy(gateway_host, "cwag", "cwag_dev.yml")

        # copy new config and restart the impacted services
        execute(_set_cwag_configs, "gateway.mconfig.multi_session_proxy")
        execute(_restart_docker_services, ["session_proxy", "pcrf", "ocs",
                                           "pcrf2", "ocs2", "ingress"])
        execute(_check_docker_services)

        _switch_to_vm_no_destroy(gateway_host, "cwag_test", "cwag_test.yml")
        execute(_run_integ_tests, test_host, trf_host,
                SubTests.MULTISESSIONPROXY, test_re)

    # If we got here means everything work well!!
    if not test_host and not trf_host:
        # Clean up only for now when running locally
        execute(_clean_up)
    print('Integration Test Passed for "{}"!'.format(tests_to_run.value))
    sys.exit(0)


def transfer_artifacts(services="sessiond session_proxy", get_core_dump=False):
    """
    Fetches service logs from Docker and optionally gets core dumps
    Args:
        services: A list of services for which services logs are requested
        get_core_dump: When set to True, it will fetch a tar of the core dumps
    """
    services = services.strip().split(' ')
    print("Transferring logs for " + str(services))

    # We do NOT want to destroy this VM after we just set it up...
    vagrant_setup("cwag", False)
    with cd(CWAG_ROOT):
        for service in services:
            run("docker logs -t " + service + " &> " + service + ".log")
            # For vagrant the files should already be in CWAG_ROOT
    if get_core_dump == "True":
        execute(_tar_coredump)


def _tar_coredump():
    _switch_to_vm_no_destroy(None, "cwag", "cwag_dev.yml")
    with cd(CWAG_ROOT):
        core_dump_dir = run('ls /var/opt/magma/cores/')
        num_of_dumps = len(core_dump_dir.split())
        print(f'Found {num_of_dumps} core dumps')
        if num_of_dumps > 0:
            run("sudo tar -czvf coredump.tar.gz /var/opt/magma/cores/*",
                warn_only=True)


def _switch_to_vm(addr, host_name, ansible_file, destroy_vm):
    if not addr:
        vagrant_setup(host_name, destroy_vm)
    else:
        ansible_setup(addr, host_name, ansible_file)


def _switch_to_vm_no_destroy(addr, host_name, ansible_file):
    _switch_to_vm(addr, host_name, ansible_file, False)


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


def _start_ipfix_controller():
    """ Start the IPFIX collector"""
    with cd("$MAGMA_ROOT"):
        sudo('mkdir -p records')
        sudo('rm -rf records/*')
        sudo('pkill ipfix_collector > /dev/null &', pty=False, warn_only=True)
        sudo('tmux new -d \'/usr/bin/ipfix_collector -4tu -p 4740 -o /home/vagrant/magma/records\'')


def _set_cwag_test_networking(mac):
    # Don't error if route already exists
    with settings(warn_only=True):
        sudo('ip route add %s/24 dev %s proto static scope link' %
             (TRF_SERVER_SUBNET, CWAG_TEST_BR_NAME))
    sudo('arp -s %s %s' % (TRF_SERVER_IP, mac))


def _add_networkhost_docker():
    ''' Add network host to docker engine '''

    local_host = "unix:///var/run/docker.sock"
    nw_host = "tcp://0.0.0.0:2375"
    tmp_daemon_json_fn = "/tmp/daemon.json"
    docker_cfg_dir = "/etc/docker/"

    # add daemon json file
    host_cfg = '\'{"hosts": [\"%s\", \"%s\"]}\'' % (local_host, nw_host)
    run("""printf %s > %s""" % (host_cfg, tmp_daemon_json_fn))
    sudo("mv %s %s" % (tmp_daemon_json_fn, docker_cfg_dir))

    # modify docker service cmd to remove hosts
    # https://docs.docker.com/config/daemon/#troubleshoot-conflicts-between-the-daemonjson-and-startup-scripts
    tmp_docker_svc_fn = "/tmp/docker.conf"
    svc_cmd = "'[Service]\nExecStart=\nExecStart=/usr/bin/dockerd'"
    svc_cfg_dir = "/etc/systemd/system/docker.service.d/"
    sudo("mkdir -p %s" % svc_cfg_dir)
    run("printf %s > %s" % (svc_cmd, tmp_docker_svc_fn))
    sudo("mv %s %s" % (tmp_docker_svc_fn, svc_cfg_dir))

    # restart systemd and docker service
    sudo("systemctl daemon-reload")
    sudo("systemctl restart docker")


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


def _check_docker_services():
    with cd(CWAG_ROOT + "/docker"):
        run(
            " DCPS=$(docker ps --format \"{{.Names}}\t{{.Status}}\" | grep Restarting);"
            " [[ -z \"$DCPS\" ]] ||"
            " ( echo \"Container restarting detected.\" ; echo \"$DCPS\"; exit 1 )"
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


def _add_docker_host_remote_network_envvar():
    sudo(
        "grep -q 'DOCKER_HOST=tcp://%s:2375' /etc/environment || "
        "echo 'DOCKER_HOST=tcp://%s:2375' >> /etc/environment && "
        "echo 'DOCKER_API_VERSION=1.40' >> /etc/environment" % (CWAG_IP, CWAG_IP))


def _run_integ_tests(test_host, trf_host, tests_to_run: SubTests,
                     test_re=None):
    """ Run the integration tests """
    # add docker host environment as well
    shell_env_vars = {
        "DOCKER_HOST" : "tcp://%s:2375" % CWAG_IP,
        "DOCKER_API_VERSION" : "1.40",
    }
    if test_re:
        shell_env_vars["TESTS"] = test_re

    # QOS take a while to run. Increasing the timeout to 20m
    go_test_cmd = "go test -v -test.short -timeout 20m"
    go_test_cmd += " -tags " + tests_to_run.value
    if test_re:
        go_test_cmd += " -run " + test_re

    with cd(CWAG_INTEG_ROOT), shell_env(**shell_env_vars):
        result = run(go_test_cmd, warn_only=True)
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
