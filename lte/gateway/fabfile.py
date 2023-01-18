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
from time import sleep

import invoke
# import fab tasks from dev_tools, so they can be called via fab in the command line
from dev_tools import (  # noqa: F401
    check_agw_cloud_connectivity,
    check_agw_feg_connectivity,
    deregister_agw,
    deregister_federated_agw,
    deregister_feg_gw,
    register_federated_vm,
    register_feg_gw,
    register_vm,
    register_vm_remote,
)
from fabric import Connection, task

sys.path.append('../../orc8r')
from tools.fab import pkg  # noqa: E402
from tools.fab.hosts import (  # noqa: E402
    ansible_setup,
    vagrant_connection,
    vagrant_setup,
)

GATEWAY_IP_ADDRESS = "192.168.60.142"
AGW_ROOT = "$MAGMA_ROOT/lte/gateway"
AGW_PYTHON_ROOT = "$MAGMA_ROOT/lte/gateway/python"
FEG_INTEG_TEST_ROOT = AGW_PYTHON_ROOT + "/integ_tests/federated_tests"
FEG_INTEG_TEST_DOCKER_ROOT = FEG_INTEG_TEST_ROOT + "/docker"
ORC8R_AGW_PYTHON_ROOT = "$MAGMA_ROOT/orc8r/gateway/python"
AGW_INTEG_ROOT = "$MAGMA_ROOT/lte/gateway/python/integ_tests"
TEST_SUMMARY_GLOB = "/var/tmp/test_results/*.xml"
debug_mode = None


@task
def open_orc8r_port_in_vagrant(c):
    """
    Add a line to the Vagrantfile file to open port 9443 on the magma_deb VM.
    Note that localhost request to 9443 will be sent to a Vagrant vm.
    Remove this line manually if you intend to run orc8r on your host.
    """
    is_port_open = "grep -q 'guest: 9443, host: 9443' Vagrantfile"

    pattern_template = "^  config.vm.define :{},.*$"
    open_port_template = '    {}.vm.network \"forwarded_port\", guest: 9443, host: 9443'
    # workaround:
    # 1. duplicate line ("p") and then override - because "a" (append text
    #    after line) works different on linux and macos (but would be more elegant)
    # 2. don't use -i (inline replace), but copy .bak file - because syntax
    #    differs between linux and macos
    cmd_open_port_template = "sed '/{}/p; s/{}/{}/' Vagrantfile > v.bak && cp v.bak Vagrantfile && rm v.bak"

    def create_cmd(machine):
        pattern = pattern_template.format(machine)
        open_port_line = open_port_template.format(machine)
        return cmd_open_port_template.format(pattern, pattern, open_port_line)

    cmd_open_port = create_cmd('magma_deb')

    c.run(f"{is_port_open} || ({cmd_open_port})")


def _redirect_feg_agw_to_vagrant_orc8r(c_gw):
    """
    Modifies feg docker-compose.override.yml hosts and AGW /etc/hosts
    to point to localhost when Orc8r runs inside Vagrant
    """
    # This is only run in CI:
    # on macOS
    c_gw.local(
        f"sed -i '' 's/:10.0.2.2/:127.0.0.1/' "
        f"{FEG_INTEG_TEST_DOCKER_ROOT}/docker-compose.override.yml",
    )
    # on ubuntu
    c_gw.run("sudo sed -i 's/10.0.2.2/127.0.0.1/' '/etc/hosts'")


@task
def federated_integ_test(
    c, build_all=False, clear_orc8r=False, provision_vm=False,
    destroy_vm=False, orc8r_on_vagrant=False,
):

    if orc8r_on_vagrant:
        # Modify Vagrantfile to allow access to Orc8r running inside Vagrant
        open_orc8r_port_in_vagrant(c)

    if build_all:
        _run_build_all(c, clear_orc8r, orc8r_on_vagrant, provision_vm)

    with c.cd(FEG_INTEG_TEST_ROOT):
        _run_start_all(c, orc8r_on_vagrant)
        c.run("fab configure-orc8r")
        sleep(20)
        c.run("fab test-connectivity --timeout=200")

    test_vm_data = _build_test_vms(c, destroy_vm, provision_vm)
    sleep(20)
    # run this on the host, not on the vm, as it will connect to the vm via ssh
    _run_integ_tests(c, test_vm_data, test_mode="federated_integ_test")


def _run_build_all(c, clear_orc8r, orc8r_on_vagrant, provision_vm):
    with c.cd(FEG_INTEG_TEST_ROOT):
        cmd = "fab build-all"
        if clear_orc8r:
            cmd += " --clear-orc8r"
        if orc8r_on_vagrant:
            cmd += " --orc8r-on-vagrant"
        if provision_vm:
            cmd += " --provision-vm"
        c.run(cmd)


def _run_start_all(c, orc8r_on_vagrant):
    start_all_cmd = "fab start-all"
    if orc8r_on_vagrant:
        start_all_cmd += " --orc8r-on-vagrant"
        # modify dns entries to find Orc8r from inside Vagrant
        with vagrant_connection(c, 'magma_deb') as c_gw:
            _redirect_feg_agw_to_vagrant_orc8r(c_gw)

    c.run(start_all_cmd)

    if orc8r_on_vagrant:
        print("Wait for orc8r to be available")
        sleep(60)


@task
def provision_magma_dev_vm(
    c, gateway_host=None, destroy_vm=True, provision_vm=True,
):
    """
    Prepare to run the integration tests on the bazel build services.
    This defaults to running on local vagrant machines, but can also be
    pointed to an arbitrary host (e.g. amazon) by passing "address:port"
    as arguments
    """
    if not gateway_host:
        vagrant_setup(c, 'magma', destroy_vm, force_provision=provision_vm)
    else:
        ansible_setup(gateway_host, "dev", "magma_dev.yml")


def _setup_vm(c, host, name, ansible_role, ansible_file, destroy_vm, provision_vm, max_retries=1):
    if not host:
        connection, host_data = vagrant_setup(
            c, name, destroy_vm, force_provision=provision_vm, max_retries=max_retries,
        )
    else:
        ansible_setup(host, ansible_role, ansible_file)
        host_data = {
            'host_string': host,
        }
        connection = Connection(host_data.get('host_string'))
    return connection, host_data


def _setup_gateway(
        c, gateway_host, name, ansible_role, ansible_file, destroy_vm,
        provision_vm, max_retries=1,
):
    gateway_connection, _ = _setup_vm(
        c, gateway_host, name, ansible_role, ansible_file, destroy_vm,
        provision_vm, max_retries=max_retries,
    )
    if gateway_host is None:
        gateway_ip = GATEWAY_IP_ADDRESS
    else:
        gateway_ip = gateway_host.split('@')[1].split(':')[0]
    return gateway_connection, gateway_ip


@task
def integ_test_deb_installation(
    c, gateway_host=None, test_host=None, trf_host=None,
    destroy_vm=True, provision_vm=True,
):
    """
    Run the integration tests. This defaults to running on local vagrant
    machines, but can also be pointed to an arbitrary host (e.g. amazon) by
    passing "address:port" as arguments

    Args:
        c: Fabric connection.
        gateway_host: The ssh address string of the machine to run the gateway
            services on. Formatted as "host:port". If not specified, defaults to
            the `magma_deb` vagrant box.
        test_host: The ssh address string of the machine to run the tests on.
            Formatted as "host:port". If not specified, defaults to the
            `magma_test` vagrant box.
        trf_host: The ssh address string of the machine to run the TrafficServer
            on. Formatted as "host:port". If not specified, defaults to the
            `magma_trfserver` vagrant box.
        destroy_vm: If True, destroy the magma deb VM before running the tests.
        provision_vm: If True, provision the magma deb VM before running the
            tests.
    """

    # Set up the gateway: use the provided gateway if given, else default to the
    # vagrant machine
    gateway_ip = _build_and_start_magma(
        c, destroy_vm, provision_vm, gateway_host=gateway_host,
        build_magma=False,
    )

    test_vm_data = _build_test_vms(
        c, destroy_vm, provision_vm, start_trfserver=True, test_host=test_host,
        trf_host=trf_host,
    )

    # run this on the host, not on the vm, as it will connect to the vm via ssh
    _run_integ_tests(c, test_vm_data, gateway_ip=gateway_ip)


@task
def integ_test_containerized(
    c, gateway_host=None, test_host=None, trf_host=None,
    destroy_vm=True, provision_vm=True,
    test_mode='integ_test_containerized',
    tests='', docker_registry=None,
):
    """
    Run the integration tests against the containerized AGW.
    Other than that the same as `integ_test`.
    """

    # Set up the gateway: use the provided gateway if given, else default to the
    # vagrant machine
    gateway_ip = _build_and_start_magma(
        c, destroy_vm, provision_vm, gateway_host=gateway_host,
        build_magma=False, containerized=True, docker_registry=docker_registry,
    )

    test_vm_data = _build_test_vms(
        c, destroy_vm, provision_vm, start_trfserver=True, test_host=test_host,
        trf_host=trf_host,
    )
    # run this on the host, not on the vm, as it will connect to the vm via ssh
    _run_integ_tests(
        c, test_vm_data, gateway_ip=gateway_ip, test_mode=test_mode,
        tests=tests,
    )


def _build_and_start_magma(
        c, destroy_vm, provision_vm, gateway_host=None, build_magma=False,
        containerized=False, docker_registry=None,
):
    # Set up the gateway: use the provided gateway if given, else default to
    # the vagrant machine
    c_gw, gateway_ip = _setup_gateway(
        c, gateway_host, "magma", "dev", "magma_dev.yml", destroy_vm,
        provision_vm,
    )
    with c_gw:
        if build_magma:
            _build_magma(c_gw)
        if containerized:
            _start_gateway_containerized(c_gw, docker_registry=docker_registry)
        else:
            _start_gateway(c_gw)
    return gateway_ip


def _start_gateway_containerized(c_gw, docker_registry=None):
    """ Starts the containerized AGW """

    c_gw.run('sudo rm -rf /etc/snowflake && sudo touch /etc/snowflake')
    with c_gw.cd("${MAGMA_ROOT}"):
        c_gw.run('bazel/scripts/link_scripts_for_bazel_integ_tests.sh')
        c_gw.run('bazel build `bazel query "attr(tags, util_script, kind(.*_binary,//orc8r/... union //lte/...))"`')
    c_gw.run(
        'for component in redis nghttpx td-agent-bit; do cp "${MAGMA_ROOT}"'
        '/{orc8r,lte}/gateway/configs/templates/${component}.conf.template;'
        ' done',
    )
    c_gw.run('sed -i \'s/init_system: systemd/init_system: docker/\' "${MAGMA_ROOT}"/lte/gateway/configs/magmad.yml')

    c_gw.run('sudo systemctl start magma_dp@envoy')

    with c_gw.cd(AGW_ROOT + "/docker"):
        # The `docker-compose up` times are machine dependent, such that a
        # retry is needed here for resilience.
        _run_with_retry(
            c_gw, f'DOCKER_REGISTRY={docker_registry} docker compose'
            f' --compatibility -f docker-compose.yaml up -d --quiet-pull',
        )


def _run_with_retry(c_gw, command, retries=10):
    iteration = 0
    while iteration < retries:
        iteration += 1
        try:
            c_gw.run(command)
            break
        except (
                invoke.exceptions.UnexpectedExit,
                invoke.exceptions.Failure,
                invoke.exceptions.ThreadException,
        ) as e:
            print(f"ERROR: Failed on retry {iteration} of \n$ {command}\n")
            print(f"ERROR: {e}\n")
            sleep(3)
    else:
        c_gw.run("docker ps")  # It is _not_ docker compose by intention to see the container ID.
        print(f"ERROR: Failed after {retries} retries of \n$ {command}")
        sys.exit(1)


@task
def get_test_summaries(c, sudo_tests=False, integration_tests=False):
    results_folder = "test_results"
    results_dir = "/var/tmp/"

    c.run('mkdir -p ' + results_folder)

    if sudo_tests == integration_tests:
        print(
            "Specify either \'sudo-tests\' or \'integration-tests\'"
            "to get test summaries",
        )
        return
    if sudo_tests:
        vm_name = "magma"
    if integration_tests:
        vm_name = "magma_test"

    with vagrant_connection(c, vm_name) as c_vm:
        if c_vm.run(
            f"test -e {results_dir}{results_folder}", warn=True,
        ).ok:
            # Fix the permissions on the files -- they have permissions 000
            # otherwise
            c_vm.run(f'sudo chmod 755 {results_dir}{results_folder}')
            _get_folder(c_vm, results_folder, results_dir, results_folder)


@task
def get_test_logs(
    c,
    gateway_host_name='magma',
    gateway_host=None,
    test_host=None,
    trf_host=None,
    dst_path="/tmp/build_logs.tar.gz",
):
    """
    Download the relevant magma logs from the given gateway and test machines.
    Place the logs in a path specified in 'dst_path' or
    "/tmp/build_logs.tar.gz" by default.

    Args:
        c: fabric connection
        gateway_host_name: name of the gateway machine
        gateway_host: The ssh address string of the gateway machine formatted
            as "host:port". If not specified, defaults to the `magma` vagrant box.
        test_host: The ssh address string of the test machine formatted as
            "host:port". If not specified, defaults to the `magma_test` vagrant
        box.
        trf_host:  The ssh address string of the machine to run the
            TrafficServer on. Formatted as "host:port". If not specified,
         defaults to the `magma_trfserver` vagrant box.
        dst_path: The path where the tarred logs will be placed on the host
    """

    # Grab the build logs from the machines and bring them to the host
    dev_logs_location = '/tmp/build_logs/dev'
    test_logs_location = '/tmp/build_logs/test'
    trfserver_logs_location = '/tmp/build_logs/trfserver'

    c.run('rm -rf /tmp/build_logs')
    c.run('mkdir /tmp/build_logs')
    c.run(f'mkdir {dev_logs_location}')
    c.run(f'mkdir {test_logs_location}')
    c.run(f'mkdir {trfserver_logs_location}')

    dev_files = [
        '/var/log/mme.log',
        '/var/log/MME.magma*log*',
        '/var/log/syslog',
        '/var/log/envoy.log',
        '/var/log/openvswitch/ovs*.log',
    ]
    _get_files_from_vm(
        c, gateway_host, gateway_host_name, dev_files, dev_logs_location,
    )

    trf_files = ['/home/vagrant/trfserver.log']
    _get_files_from_vm(
        c, trf_host, 'magma_trfserver', trf_files, trfserver_logs_location,
    )

    test_files = ['/var/log/syslog', '/tmp/fw/*']
    _get_files_from_vm(
        c, test_host, 'magma_test', test_files, test_logs_location,
    )

    c.run("tar -czvf /tmp/build_logs.tar.gz /tmp/build_logs/*")
    if dst_path != "/tmp/build_logs.tar.gz":
        c.run(f'mv /tmp/build_logs.tar.gz {dst_path}', warn=True)
    c.run('rm -rf /tmp/build_logs')


def _get_files_from_vm(c, host, vm_name, files, logs_location):
    c_vm, host_data = vagrant_setup(c, vm_name)
    if host:
        c_vm = Connection(
            host,
            connect_kwargs={"key_filename": host_data.get("key_filename")},
        )

    with c_vm:
        for p in files:
            if c_vm.run(f"test -e {p}", warn=True).ok:
                # Fix the permissions on the files -- they have permissions 000
                # otherwise
                c_vm.run(f'sudo chmod 755 {p}')
                if p[-1] == '/':
                    folder = p.split('/')[-2]
                    path = p.split(folder)[0]
                    _get_folder(c_vm, folder, path, logs_location)
                else:
                    c_vm.get(p, local=f"{logs_location}/{p}")


def _get_folder(c_vm, folder_name, remote_path, local_path):
    """
    Get a folder from the remote machine to the local machine
    """
    with c_vm.cd(remote_path):
        c_vm.run(f'tar -czvf /tmp/{folder_name}.tar.gz {folder_name}')
    c_vm.get(f'/tmp/{folder_name}.tar.gz', local=f'{local_path}/{folder_name}.tar.gz')
    c_vm.run(f'rm /tmp/{folder_name}.tar.gz')
    c_vm.local(f'sudo tar -xzf {local_path}/{folder_name}.tar.gz -C {local_path}')
    c_vm.local(f'sudo chmod 755 {local_path}/{folder_name}')
    c_vm.local(f'sudo rm {local_path}/{folder_name}.tar.gz')


@task
def build_and_start_magma_trf(c, destroy_vm=False, provision_vm=False):
    with vagrant_connection(
        c, 'magma_trfserver', destroy_vm=destroy_vm, force_provision=provision_vm,
    ) as c_trf:
        _start_trfserver(c_trf)


def _start_gateway(c_gw):
    """ Starts the gateway """
    c_gw.run('sudo service magma@magmad start')


def _build_test_vms(
        c, destroy_vm=False, provision_vm=False, start_trfserver=False,
        test_host=None, trf_host=None,
):
    _start_trfserver_vm(c, destroy_vm, provision_vm, start_trfserver, trf_host)
    test_host_data = _start_test_vm(c, destroy_vm, provision_vm, test_host)
    return test_host_data


def _start_trfserver_vm(c, destroy_vm, provision_vm, start_trfserver, trf_host):
    # Set up the trfserver: use the provided trfserver if given, else default to the
    # vagrant machine
    c_trf, _ = _setup_vm(
        c, trf_host, "magma_trfserver", "trfserver", "magma_trfserver.yml",
        destroy_vm, provision_vm,
    )
    if start_trfserver:
        with c_trf:
            _start_trfserver(c_trf)


def _start_test_vm(c, destroy_vm, provision_vm, test_host):
    # Run the tests: use the provided test machine if given, else default to
    # the vagrant machine
    _, test_host_data = _setup_vm(
        c, test_host, "magma_test", "test", "magma_test.yml", destroy_vm,
        provision_vm,
    )
    return test_host_data


def _start_trfserver(c_trf):
    """ Starts the traffic gen server"""

    c_trf.run('sudo ethtool --offload eth1 rx off tx off')
    c_trf.run('sudo ethtool --offload eth2 rx off tx off')
    trf_cmd = 'nohup /usr/local/bin/traffic_server.py 192.168.60.144 62462 > trfserver.log 2>&1'
    c_trf.run('sudo apt-get install -y dtach')
    c_trf.run(f"sudo dtach -n `mktemp -u /tmp/dtach.XXXX` {trf_cmd}")


def _run_integ_tests(
        c, vm_data, gateway_ip='192.168.60.142', test_mode='integ_test', tests='',
):
    """ Run the integration tests

    NOTE: The S1AP-tester produces a bunch of output which the python ssh
    library, and thus fab, has trouble processing quickly. Instead, we manually
    ssh into machine and run the tests.

    ssh switch reference:
        -i: identity file
        -tt: (really) force a pseudo tty -- The tests can't initialize logging
            without this
        -p: the port to connect to
        -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no: have ssh
         never prompt to confirm the host fingerprints
    """
    host = vm_data.get("host_string").split(':')[0]
    port = vm_data.get("host_string").split(':')[1]
    key = vm_data.get("key_filename")

    # We do not have a proper shell, so the `magtivate` alias is not available.
    # We instead directly source the activate file.
    c.run(
        f'ssh'
        f' -i {key}'
        f' -o UserKnownHostsFile=/dev/null'
        f' -o StrictHostKeyChecking=no'
        f' -tt {host}'
        f' -p {port}'
        f' \'cd $MAGMA_ROOT/lte/gateway/python/integ_tests;'
        f' sudo ethtool --offload eth1 rx off tx off;'
        f' sudo ethtool --offload eth2 rx off tx off;'
        f' source ~/build/python/bin/activate;'
        f' export GATEWAY_IP={gateway_ip};'
        f' make {test_mode} enable-flaky-retry=true {tests};'
        f' make evaluate_result\'',
    )
