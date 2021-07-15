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
import sys
import time
from enum import Enum

from fabric.api import (
    cd,
    env,
    execute,
    hide,
    lcd,
    local,
    put,
    run,
    settings,
    shell_env,
    sudo,
)
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
    HSSLESS = "hssless"

    @staticmethod
    def list():
        return list(map(lambda t: t.value, SubTests))


def integ_test(
    gateway_host=None, test_host=None, trf_host=None,
    gateway_vm="cwag", gateway_ansible_file="cwag_dev.yml",
    transfer_images=False, destroy_vm=False, no_build=False,
    tests_to_run="all", skip_unit_tests=False, test_re=None,
    test_result_xml=None, run_tests=True, count="1", provision_vm=True,
):
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
        print(
            "{} is not a valid value. We support {}".format(
                tests_to_run, SubTests.list(),
            ),
        )
        return
    provision_vm = False if provision_vm == "False" else True

    # Setup the gateway: use the provided gateway if given, else default to the
    # vagrant machine
    _switch_to_vm(
        gateway_host,
        gateway_vm,
        gateway_ansible_file,
        destroy_vm,
        provision_vm,
    )

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

    # Setup the trfserver: use the provided trfserver if given, else default to
    # the vagrant machine
    with lcd(LTE_AGW_ROOT):
        _switch_to_vm(
            gateway_host, "magma_trfserver",
            "magma_trfserver.yml", destroy_vm, provision_vm,
        )

    execute(_start_trfserver)

    # Run the tests: use the provided test machine if given, else default to
    # the vagrant machine
    _switch_to_vm(
        gateway_host,
        "cwag_test",
        "cwag_test.yml",
        destroy_vm,
        provision_vm,
    )

    cwag_test_host_to_mac = execute(_get_br_mac, CWAG_TEST_BR_NAME)
    host = env.hosts[0]
    cwag_test_br_mac = cwag_test_host_to_mac[host]
    execute(_set_cwag_test_configs)
    execute(_start_ipfix_controller)

    # Get back to the gateway vm to setup static arp
    _switch_to_vm_no_destroy(gateway_host, gateway_vm, gateway_ansible_file)
    execute(_set_cwag_networking, cwag_test_br_mac)

    # check if docker services are alive except for OCS2 and PCRF2
    ignore_list = ["ocs2", "pcrf2"]
    execute(_check_docker_services, ignore_list)

    _switch_to_vm_no_destroy(gateway_host, "cwag_test", "cwag_test.yml")
    execute(_start_ue_simulator)
    execute(_set_cwag_test_networking, cwag_br_mac)

    if run_tests == "False":
        execute(_add_docker_host_remote_network_envvar)
        print(
            "run_test was set to false. Test will not be run\n"
            "You can now run the tests manually from cwag_test",
        )
        sys.exit(0)

    # HSSLESS tests are to be executed from gateway_host VM
    if tests_to_run.value == SubTests.HSSLESS.value:
        _switch_to_vm_no_destroy(
            gateway_host,
            gateway_vm,
            gateway_ansible_file,
        )
        execute(
            _run_integ_tests, gateway_host, trf_host,
            tests_to_run, test_re, count, test_result_xml,
        )
    else:
        execute(
            _run_integ_tests, test_host, trf_host,
            tests_to_run, test_re, count, test_result_xml,
        )

    # If we got here means everything work well!!
    if not test_host and not trf_host:
        # Clean up only for now when running locally
        execute(_clean_up)
    print('Integration Test Passed for "{}"!'.format(tests_to_run.value))
    sys.exit(0)


def transfer_artifacts(
    gateway_vm="cwag", gateway_ansible_file="cwag_dev.yml",
    services="sessiond session_proxy", get_core_dump=False,
):
    """
    Fetches service logs from Docker and optionally gets core dumps
    Args:
        services: A list of services for which services logs are requested
        get_core_dump: When set to True, it will fetch a tar of the core dumps
    """
    services = services.strip().split(' ')
    print("Transferring logs for " + str(services))

    # We do NOT want to destroy this VM after we just set it up...
    vagrant_setup(gateway_vm, False, False)
    with cd(CWAG_ROOT):
        for service in services:
            run("docker logs -t " + service + " &> " + service + ".log")
            # For vagrant the files should already be in CWAG_ROOT
    if get_core_dump == "True":
        execute(
            _tar_coredump, gateway_vm=gateway_vm,
            gateway_ansible_file=gateway_ansible_file,
        )


def _tar_coredump(gateway_vm="cwag", gateway_ansible_file="cwag_dev.yml"):
    _switch_to_vm_no_destroy(None, gateway_vm, gateway_ansible_file)
    with cd(CWAG_ROOT):
        core_dump_dir = run('ls /var/opt/magma/cores/')
        num_of_dumps = len(core_dump_dir.split())
        print(f'Found {num_of_dumps} core dumps')
        if num_of_dumps > 0:
            run(
                "sudo tar -czvf coredump.tar.gz /var/opt/magma/cores/*",
                warn_only=True,
            )


def _switch_to_vm(addr, host_name, ansible_file, destroy_vm, provision_vm):
    if not addr:
        vagrant_setup(host_name, destroy_vm, provision_vm)
    else:
        ansible_setup(addr, host_name, ansible_file)


def _switch_to_vm_no_destroy(addr, host_name, ansible_file):
    _switch_to_vm(addr, host_name, ansible_file, False, False)


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
        sudo(
            'ip route add %s/24 dev %s proto static scope link' %
            (TRF_SERVER_SUBNET, CWAG_TEST_BR_NAME),
        )
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
    sudo('mkdir -p {}'.format(docker_cfg_dir))
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
        sudo(
            ' docker-compose'
            ' -f docker-compose.yml'
            ' -f docker-compose.override.yml'
            ' -f docker-compose.integ-test.yml'
            ' down',
        )


def _build_gateway():
    """ Builds the gateway docker images """
    with cd(CWAG_ROOT + '/docker'):
        sudo(
            ' docker-compose'
            ' -f docker-compose.yml'
            ' -f docker-compose.override.yml'
            ' -f docker-compose.nginx.yml'
            ' -f docker-compose.integ-test.yml'
            ' build --parallel',
        )


def _run_gateway():
    """ Runs the gateway's docker images """
    with cd(CWAG_ROOT + '/docker'):
        sudo(
            ' docker-compose'
            ' -f docker-compose.yml'
            ' -f docker-compose.override.yml'
            ' -f docker-compose.integ-test.yml'
            ' up -d ',
        )


def _restart_docker_services(services):
    with cd(CWAG_ROOT + "/docker"):
        sudo(
            " docker-compose"
            " -f docker-compose.yml"
            " -f docker-compose.override.yml"
            " -f docker-compose.nginx.yml"
            " -f docker-compose.integ-test.yml"
            " restart {}".format(" ".join(services)),
        )


def _stop_docker_services(services):
    with cd(CWAG_ROOT + "/docker"):
        sudo(
            " docker-compose"
            " -f docker-compose.yml"
            " -f docker-compose.override.yml"
            " -f docker-compose.nginx.yml"
            " -f docker-compose.integ-test.yml"
            " stop {}".format(" ".join(services)),
        )


def _check_docker_services(ignore_list):
    with cd(CWAG_ROOT + "/docker"), settings(warn_only=True), hide("warnings"):

        grep_ignore = "| grep --invert-match '" + \
            '\\|'.join(ignore_list) + "'" if ignore_list else ""
        count = 0
        while (count < 5):
            # force wait to make sure docker logs are up
            time.sleep(1)
            result = run(
                " docker ps --format \"{{.Names}}\t{{.Status}}\" | "
                "grep Restarting" + grep_ignore,
            )

            if result.return_code == 1:
                # grep returns code 1 when empty string
                return
            print("Container restarting detected. Trying one more time")
            count += 1
    # if we got here, that means all attempts failed
    print("ERROR: Test NOT started due to docker container restarting")
    sys.exit(1)


def _start_ue_simulator():
    """ Starts the UE Sim Service and logs into uesim.log"""
    with cd(CWAG_ROOT + '/services/uesim/uesim'):
        run('tmux new -d \'go run main.go -logtostderr=true -v=9 &> %s/uesim.log\'' % CWAG_ROOT)


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
        "echo 'DOCKER_API_VERSION=1.40' >> /etc/environment" % (
            CWAG_IP, CWAG_IP,
        ),
    )


def _run_integ_tests(
    test_host, trf_host, tests_to_run: SubTests,
    test_re=None, count="1", test_result_xml=None,
):
    """ Run the integration tests """
    # add docker host environment as well
    shell_env_vars = {
        "DOCKER_HOST": "tcp://%s:2375" % CWAG_IP,
        "DOCKER_API_VERSION": "1.40",
    }
    if test_re:
        shell_env_vars["TESTS"] = test_re

    # QOS take a while to run. Increasing the timeout to 50m
    go_test_cmd = "gotestsum --format=standard-verbose "
    # Retry once on failure
    go_test_cmd += "--rerun-fails=1 --packages='./...' "
    if test_result_xml:  # generate test result XML in cwf/gateway directory
        go_test_cmd += "--junitfile ../" + test_result_xml + " "
    go_test_cmd += " -- -test.short -timeout 50m -count " + count  # go test args
    go_test_cmd += " -tags=" + tests_to_run.value
    if test_re:
        go_test_cmd += " -run=" + test_re

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
        vagrant_setup("magma_trfserver", False, False)
        run('pkill iperf3 > /dev/null &', pty=False, warn_only=True)


class FabricException(Exception):
    pass
