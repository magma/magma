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
import re
import sys
import time
from enum import Enum

from fabric import Connection, task

sys.path.append('../../orc8r')
from tools.fab.hosts import ansible_setup, vagrant_connection

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


@task
def integ_test(
    c, gateway_host=None, test_host=None, trf_host=None,
    gateway_vm="cwag", gateway_ansible_file="cwag_dev.yml",
    transfer_images=False, skip_docker_load=False, tar_path="/tmp/cwf-images",
    destroy_vm=False, build=True,
    tests_to_run="all", skip_unit_tests=False, test_re=None,
    test_result_xml=None, run_tests=True, count="1", provision_vm=True,
    rerun_fails="1",
):
    """
    Run the integration tests. This defaults to running on local vagrant
    machines, but can also be pointed to an arbitrary host (e.g. amazon) by
    passing "address:port" as arguments

    Args:
        c: Fabric connection.
        gateway_host: The ssh address string of the machine to run the gateway
            services on. Formatted as "host:port". If not specified, defaults
            to the `cwag` vagrant box.
        test_host: The ssh address string of the machine to run the tests on.
            Formatted as "host:port". If not specified, defaults to the
            `cwag_test` vagrant box.
        trf_host: The ssh address string of the machine to run the tests on.
            Formatted as "host:port". If not specified, defaults to the
            `magma_trfserver` vagrant box.
        gateway_vm: The name of the vagrant VM to use as the gateway.
        gateway_ansible_file: The ansible file to use to provision the gateway.
        destroy_vm: When set to true, all VMs will be destroyed before running
            the tests.
        provision_vm: When set to False, this script will not provision the VMs
            before running the tests.
        build: When set to false, this script will skip rebuilding all docker images
            in the CWAG VM
        tests_to_run: The tests to run. Valid values are in the SubTests enum.
        transfer_images: When set to true, the script will transfer all cwf_*
            docker images from the host machine to the CWAG VM to use in the
            test.
        skip_docker_load: When set to true, /tmp/cwf_* will be copied into the
            CWAG VM instead of loading the docker images then copying.
            This option only is valid if transfer_images is set.
        tar_path: The location where the tarred docker images will be copied
            from. Only valid if transfer_images is set.
        skip_unit_tests: When set to true, only integration tests will be run.
        run_tests: When set to false, no integration tests will be run.
        test_re: When set to a value, integrations tests that match the
            expression will be run.
            (Ex: test_re=TestAuth will run all tests that start with TestAuth)
        count: When set to a number, the integrations tests will be run that
            many times.
        test_result_xml: When set to a path, a JUnit style test summary in XML
            will be produced at the path.
        rerun_fails: Number of times to re-run a test on failure.
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

    # Set up the gateway: use the provided gateway if given, else default to the
    # vagrant machine
    c_cwf = _set_up_vm(
        c,
        gateway_host,
        gateway_vm,
        gateway_ansible_file,
        destroy_vm,
        provision_vm,
    )

    # We will direct coredumps to be placed in this directory
    # Clean up before every run
    with c_cwf:
        if c_cwf.run("test -e /var/opt/magma/cores", warn=True).ok:
            c_cwf.run("sudo rm /var/opt/magma/cores/*", warn=True, hide='err')
        else:
            c_cwf.run("sudo mkdir -p /var/opt/magma/cores", warn=True)

        if not skip_unit_tests:
            _run_unit_tests(c_cwf)

        _set_cwag_configs(c_cwf, "gateway.mconfig")
        _add_networkhost_docker(c_cwf)
        cwag_br_mac = _get_br_mac(c_cwf, CWAG_BR_NAME)

        # Transfer built images from local machine to CWAG host
        if gateway_host or transfer_images:
            _transfer_docker_images(c_cwf, skip_docker_load, tar_path)
        else:
            _stop_gateway(c_cwf)
            if build:
                _build_gateway(c_cwf)

        _run_gateway(c_cwf)

    # Set up the trfserver: use the provided trfserver if given, else default to
    # the vagrant machine
    with c.cd(LTE_AGW_ROOT):
        c_trf = _set_up_vm(
            c, gateway_host, "magma_trfserver",
            "magma_trfserver.yml", destroy_vm, provision_vm,
        )
        with c_trf:
            _start_trfserver(c_trf)

    # Run the tests: use the provided test machine if given, else default to
    # the vagrant machine
    c_test = _set_up_vm(
        c,
        gateway_host,
        "cwag_test",
        "cwag_test.yml",
        destroy_vm,
        provision_vm,
    )
    with c_test:
        cwag_test_br_mac = _get_br_mac(c_test, CWAG_TEST_BR_NAME)
        _set_cwag_test_configs(c_test)
        _set_cwag_configs(c_test, "gateway.mconfig")
        _start_ipfix_controller(c_test)

    # Get back to the gateway vm to set up static arp
    with c_cwf:
        _set_cwag_networking(c_cwf, cwag_test_br_mac)

        # check if docker services are alive except for OCS2 and PCRF2
        ignore_list = ["ocs2", "pcrf2"]
        _check_docker_services(c_cwf, ignore_list)

    with c_test:
        _start_ue_simulator(c_test)
        _set_cwag_test_networking(c_test, cwag_br_mac)

        if not run_tests:
            _add_docker_host_remote_network_envvar(c_test)
            print(
                "run_tests was set to false. Test will not be run\n"
                "You can now run the tests manually from cwag_test",
            )
            sys.exit(0)

    # HSSLESS tests are to be executed from gateway_host VM
    if tests_to_run.value == SubTests.HSSLESS.value:
        _run_integ_tests(
            gateway_host, trf_host, tests_to_run, test_re, count,
            test_result_xml, rerun_fails, c_cwf, c_trf,
        )
    else:
        _run_integ_tests(
            test_host, trf_host, tests_to_run, test_re, count,
            test_result_xml, rerun_fails, c_test, c_trf,
        )


@task
def transfer_artifacts(
    c, gateway_vm="cwag", gateway_ansible_file="cwag_dev.yml",
    services="sessiond session_proxy", get_core_dump=False,
):
    """
    Fetches service logs from Docker and optionally gets core dumps
    Args:
        c: Fabric connection
        gateway_vm: VM to fetch logs from
        gateway_ansible_file: Ansible file to use for VM
        services: A list of services for which services logs are requested
        get_core_dump: When set to True, it will fetch a tar of the core dumps
    """
    services = services.strip().split(' ')
    print("Transferring logs for " + str(services))

    c_cwf = _set_up_vm(
        c, None, gateway_vm, gateway_ansible_file, destroy_vm=False,
        provision_vm=False,
    )
    with c_cwf:
        with c_cwf.cd(CWAG_ROOT):
            for service in services:
                c_cwf.run("docker logs -t " + service + " &> " + service + ".log")
                # For vagrant the files should already be in CWAG_ROOT
            if get_core_dump:
                _tar_coredump(c_cwf)


def _tar_coredump(c_cwf):
    core_dump_dir = c_cwf.run('ls /var/opt/magma/cores/')
    num_of_dumps = len(core_dump_dir.split())
    print(f'Found {num_of_dumps} core dumps')
    if num_of_dumps > 0:
        c_cwf.run(
            "sudo tar -czvf coredump.tar.gz /var/opt/magma/cores/*",
            warn=True,
        )


def _set_up_vm(c, addr, host_name, ansible_file, destroy_vm, provision_vm):
    if not addr:
        return vagrant_connection(c, host_name, destroy_vm, provision_vm)
    else:
        ansible_setup(c, addr, host_name, ansible_file)
        return Connection(addr)


def _transfer_docker_images(c_cwf, skip_docker_load, tar_path):
    if skip_docker_load:
        print("Skipping docker save step and grabbing whatever is in " + tar_path)
    else:
        print("Loading cwf_* docker images into " + tar_path)
        c_cwf.local("rm -rf " + tar_path)
        c_cwf.local("mkdir -p " + tar_path)
        output = c_cwf.local("docker images cwf_*").stdout
        for line in output.splitlines():
            if not line.startswith('cwf'):
                continue
            line = line.rstrip("\n")
            image = line.split(" ")[0]
            c_cwf.local(f"docker save -o {tar_path}/{image}.tar {image}")

    output = c_cwf.local(f"ls {tar_path}/cwf_*.tar").stdout
    for line in output.splitlines():
        regex = f'{tar_path}/(.+?).tar'
        match = re.search(regex, line)
        if match:
            image = match.group(1)
            c_cwf.put(f"{tar_path}/{image}.tar", f"{image}.tar")
            c_cwf.run(f'docker load -i {image}.tar')
        else:
            print("We should not be here, but " + line + " does not match regex: " + regex)


def _set_cwag_configs(c_vm, configfile):
    """ Set the necessary config overrides """
    c_vm.run('sudo mkdir -p /var/opt/magma/configs')
    c_vm.run(f"sudo cp {CWAG_INTEG_ROOT}/{configfile} /var/opt/magma/configs/gateway.mconfig")


def _set_cwag_networking(c_cwf, mac):
    c_cwf.run(f'sudo arp -s {CWAG_TEST_IP} {mac}')


def _get_br_mac(c_vm, bridge_name):
    mac = c_vm.run(f"cat /sys/class/net/{bridge_name}/address").stdout.strip()
    return mac


def _set_cwag_test_configs(c_test):
    """ Set the necessary test configs """
    c_test.run('sudo mkdir -p /etc/magma')
    # Create empty uesim config
    c_test.run('sudo touch /etc/magma/uesim.yml')
    c_test.run('sudo cp -r $MAGMA_ROOT/cwf/gateway/configs /etc/magma/')


def _start_ipfix_controller(c_test):
    """ Start the IPFIX collector"""
    with c_test.cd("$MAGMA_ROOT"):
        c_test.run('mkdir -p records')
        c_test.run('rm -rf records/*')
        c_test.run('sudo pkill ipfix_collector > /dev/null &', pty=False, warn=True)
        c_test.run('sudo tmux new -d \'/usr/bin/ipfix_collector -4tu -p 4740 -o /home/vagrant/magma/records\'')


def _set_cwag_test_networking(c_test, mac):
    # Don't error if route already exists
    c_test.run(
        f'sudo ip route add {TRF_SERVER_SUBNET}/24 '
        f'dev {CWAG_TEST_BR_NAME} proto static scope link',
        warn=True,
    )
    c_test.run(f'sudo arp -s {TRF_SERVER_IP} {mac}')


def _add_networkhost_docker(c_cwf):
    """ Add network host to docker engine """

    local_host = "unix:///var/run/docker.sock"
    nw_host = "tcp://0.0.0.0:2375"
    tmp_daemon_json_fn = "/tmp/daemon.json"
    docker_cfg_dir = "/etc/docker/"

    # add daemon json file
    host_cfg = '\'{"hosts": [\"%s\", \"%s\"]}\'' % (local_host, nw_host)
    c_cwf.run("""printf %s > %s""" % (host_cfg, tmp_daemon_json_fn))
    c_cwf.run(f'sudo mkdir -p {docker_cfg_dir}')
    c_cwf.run(f"sudo mv {tmp_daemon_json_fn} {docker_cfg_dir}")

    # modify docker service cmd to remove hosts
    # https://docs.docker.com/config/daemon/#troubleshoot-conflicts-between-the-daemonjson-and-startup-scripts
    tmp_docker_svc_fn = "/tmp/docker.conf"
    svc_cmd = "'[Service]\nExecStart=\nExecStart=/usr/bin/dockerd'"
    svc_cfg_dir = "/etc/systemd/system/docker.service.d/"
    c_cwf.run(f"sudo mkdir -p {svc_cfg_dir}")
    c_cwf.run(f"printf {svc_cmd} > {tmp_docker_svc_fn}")
    c_cwf.run(f"sudo mv {tmp_docker_svc_fn} {svc_cfg_dir}")

    # restart systemd and docker service
    c_cwf.run("sudo systemctl daemon-reload")
    c_cwf.run("sudo systemctl restart docker")


def _stop_gateway(c_cwf):
    """ Stop the gateway docker images """
    with c_cwf.cd(CWAG_ROOT + '/docker'):
        c_cwf.run(
            ' docker compose'
            ' -f docker-compose.yml'
            ' -f docker-compose.override.yml'
            ' -f docker-compose.integ-test.yml'
            ' down',
        )


def _build_gateway(c_cwf):
    """ Builds the gateway docker images """
    with c_cwf.cd(CWAG_ROOT + '/docker'):
        c_cwf.run(
            ' docker compose'
            ' --compatibility'
            ' -f docker-compose.yml'
            ' -f docker-compose.override.yml'
            ' -f docker-compose.nginx.yml'
            ' -f docker-compose.integ-test.yml'
            ' build',
        )


def _run_gateway(c_cwf):
    """ Runs the gateway's docker images """
    with c_cwf.cd(CWAG_ROOT + '/docker'):
        c_cwf.run(
            ' docker compose'
            ' --compatibility'
            ' -f docker-compose.yml'
            ' -f docker-compose.override.yml'
            ' -f docker-compose.integ-test.yml'
            ' up -d ',
        )


def _check_docker_services(c_cwf, ignore_list):
    with c_cwf.cd(CWAG_ROOT + "/docker"):
        grep_ignore = "| grep --invert-match '" + \
            '\\|'.join(ignore_list) + "'" if ignore_list else ""
        count = 0
        while count < 5:
            # force wait to make sure docker logs are up
            time.sleep(1)
            result = c_cwf.run(
                " docker ps --format \"{{.Names}}\t{{.Status}}\" | "
                "grep Restarting" + grep_ignore,
                hide=True, warn=True,
            )

            if result.return_code == 1:
                # grep returns code 1 when empty string
                return
            print("Container restarting detected. Trying one more time")
            count += 1
    # if we got here, that means all attempts failed
    print("ERROR: Test NOT started due to docker container restarting")
    sys.exit(1)


def _start_ue_simulator(c_test):
    """ Starts the UE Sim Service and logs into uesim.log"""
    with c_test.cd(CWAG_ROOT + '/services/uesim/uesim'):
        c_test.run(
            env={'PATH': '$PATH:/usr/local/go/bin:/home/vagrant/go/bin'},
            command=f'tmux new -d \'go run main.go -logtostderr=true -v=9 &> {CWAG_ROOT}/uesim.log\'',
        )


def _start_trfserver(c_trf):
    """ Starts the traffic gen server"""
    trf_cmd = f'nohup iperf3 -s --json -B {TRF_SERVER_IP} > /dev/null &'
    c_trf.run('sudo apt-get install -y dtach')
    c_trf.run(f"sudo dtach -n `mktemp -u /tmp/dtach.XXXX` {trf_cmd}")


def _run_unit_tests(c_cwf):
    """ Run the cwag unit tests """
    with c_cwf.cd(CWAG_ROOT):
        c_cwf.run(
            'make test',
            env={'PATH': '$PATH:/usr/local/go/bin:/home/vagrant/go/bin'},
        )


def _add_docker_host_remote_network_envvar(c_test):
    c_test.run(
        f"grep -q 'DOCKER_HOST=tcp://{CWAG_IP}:2375' /etc/environment || "
        f"echo 'DOCKER_HOST=tcp://{CWAG_IP}:2375' | sudo tee -a /etc/environment > /dev/null && "
        "echo 'DOCKER_API_VERSION=1.40' | sudo tee -a /etc/environment > /dev/null",
    )


def _run_integ_tests(
    test_host, trf_host, tests_to_run: SubTests, test_re, count,
    test_result_xml, rerun_fails, c_test_vm, c_trf,
):
    """ Run the integration tests """
    # add docker host environment as well
    shell_env_vars = {
        "DOCKER_HOST": f"tcp://{CWAG_IP}:2375",
        "DOCKER_API_VERSION": "1.40",
        "PATH": "$PATH:/usr/local/go/bin:/home/vagrant/go/bin",
    }
    if test_re:
        shell_env_vars["TESTS"] = test_re

    # QOS take a while to run. Increasing the timeout to 50m
    go_test_cmd = "gotestsum --format=standard-verbose "
    # Set retry count on failure
    go_test_cmd += "--rerun-fails=" + rerun_fails + " --packages='./...' "
    if test_result_xml:  # generate test result XML in cwf/gateway directory
        go_test_cmd += "--junitfile ../" + test_result_xml + " "
    go_test_cmd += " -- -test.short -timeout 50m -count " + count  # go test args
    go_test_cmd += " -tags=" + tests_to_run.value
    if test_re:
        go_test_cmd += " -run=" + test_re

    with c_test_vm:
        with c_test_vm.cd(CWAG_INTEG_ROOT):
            result = c_test_vm.run(go_test_cmd, env=shell_env_vars, warn=True)

        if result.return_code != 0:
            if not test_host and not trf_host:
                # Clean up only for now when running locally
                _clean_up(c_test_vm, c_trf)

        print("Integration Test returned exit code ", result.return_code)
        sys.exit(result.return_code)


def _clean_up(c_test_vm, c_trf):
    # Kill uesim service
    c_test_vm.run('sudo pkill go', warn=True)
    with c_test_vm.cd(LTE_AGW_ROOT):
        with c_trf:
            c_trf.run('sudo pkill iperf3 > /dev/null &', pty=False, warn=True)


class FabricException(Exception):
    pass
