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
import tools.fab.pkg as pkg
from tools.fab.dev_utils import connect_gateway_to_cloud
from tools.fab.hosts import ansible_setup, vagrant_connection, vagrant_setup

"""
Magma Gateway packaging tool:

Magma packages released to different channels have different version schemes.

    - dev: used for development.

    - release: used for Continuous Integration (CI). Packages in the `release`
            channel should be built and released automatically.

# HOWTO build magma.deb

1. From your laptop, update the magma version number in `release/build-magma.sh`

2. From the dev VM, build the magma package. Dependency packages are recorded
in `release/magma.lockfile`

    fab dev package
    # optionally upload to aws (if you are configured for it)
    fab dev package upload_to_aws
"""

GATEWAY_IP_ADDRESS = "192.168.60.142"
AGW_ROOT = "$MAGMA_ROOT/lte/gateway"
AGW_PYTHON_ROOT = "$MAGMA_ROOT/lte/gateway/python"
FEG_INTEG_TEST_ROOT = AGW_PYTHON_ROOT + "/integ_tests/federated_tests"
FEG_INTEG_TEST_DOCKER_ROOT = FEG_INTEG_TEST_ROOT + "/docker"
ORC8R_AGW_PYTHON_ROOT = "$MAGMA_ROOT/orc8r/gateway/python"
AGW_INTEG_ROOT = "$MAGMA_ROOT/lte/gateway/python/integ_tests"
DEFAULT_CERT = "$MAGMA_ROOT/.cache/test_certs/rootCA.pem"
DEFAULT_PROXY = "$MAGMA_ROOT/lte/gateway/configs/control_proxy.yml"
TEST_SUMMARY_GLOB = "/var/tmp/test_results/*.xml"
debug_mode = None


@task
def dev(c):
    global debug_mode
    debug_mode = True


@task
def release(c):
    """Set debug_mode to False, should be used for producing a production AGW package"""
    global debug_mode
    debug_mode = False


@task
def package(
    c, all_deps=False,
    cert_file=DEFAULT_CERT, proxy_config=DEFAULT_PROXY,
    destroy_vm=False,
    vm='magma', os="ubuntu",
):
    """ Builds the magma package """

    global debug_mode
    if debug_mode is None:
        raise RuntimeError(
            "Error: The Deploy target isn't specified. Specify one with\n\n" +
            "\tfab [dev|release] package",
        )

    hash = pkg.get_commit_hash(c)
    commit_count = pkg.get_commit_count(c)

    with vagrant_connection(c, vm, destroy_vm=destroy_vm) as c_gw:
        print('Uninstalling dev dependencies of the VM')
        c_gw.run('sudo pip uninstall --yes mypy-protobuf grpcio-tools grpcio protobuf')

        with c_gw.cd('~/magma/lte/gateway'):
            c_gw.run('mkdir -p ~/magma-deps')
            print(
                'Generating lte/setup.py and orc8r/setup.py magma dependency packages',
            )
            c_gw.run(
                './release/pydep finddep --install-from-repo -b --build-output '
                + '~/magma-deps'
                + f' -l ./release/magma.lockfile.{os}'
                + ' python/setup.py'
                + f' {ORC8R_AGW_PYTHON_ROOT}/setup.py',
            )

            print(f'Building magma package, picking up commit {hash}...')
            c_gw.run('make clean')
            build_type = "Debug" if debug_mode else "RelWithDebInfo"

            build_cmd = f'./release/build-magma.sh --hash {hash}' \
                        f' --commit-count {commit_count} --type {build_type}' \
                        f' --cert {cert_file} --proxy {proxy_config} --os {os}'
            # set '/usr/bin/' in PATH to ensure that the correct version of
            # python is used
            c_gw.run(
                build_cmd,
                env={'PATH': '/usr/bin:/usr/local/go/bin:/home/vagrant/go/bin:/usr/lib/ccache:$PATH'},
            )

            c_gw.run('rm -rf ~/magma-packages')
            c_gw.run('mkdir -p ~/magma-packages')
            c_gw.run('cp -f ~/magma-deps/*.deb ~/magma-packages', warn=True)
            c_gw.run('mv *.deb ~/magma-packages')

            with c_gw.cd('release'):
                mirrored_packages_file = 'mirrored_packages'
                if os == "ubuntu":
                    mirrored_packages_file += '_focal'
                if vm and vm.startswith('magma_'):
                    mirrored_packages_file += vm[5:]

                c_gw.run(
                    f'cat {mirrored_packages_file}'
                    + ' | xargs -I% sudo aptitude download -q2 %',
                )
                c_gw.run('cp *.deb ~/magma-packages')
                c_gw.run('sudo rm -f *.deb')

            if all_deps:
                pkg.download_all_pkgs(c_gw)
                c_gw.run('cp /var/cache/apt/archives/*.deb ~/magma-packages')

            # Copy out C executables into magma-packages as well
            _copy_out_c_execs_in_magma_vm(c_gw)


@task
def openvswitch(c, destroy_vm=False, destdir='~/magma-packages/'):
    # If a host list isn't specified, default to the magma vagrant vm
    with vagrant_connection(c, 'magma', destroy_vm=destroy_vm) as c_gw:
        c_gw.run('~/magma/third_party/gtp_ovs/ovs-gtp-patches/2.15/build.sh ' + destdir)


@task
def depclean(c):
    '''Remove all generated packaged for dependencies'''
    # If a host list isn't specified, default to the magma vagrant vm
    with vagrant_connection(c, 'magma') as c_gw:
        c_gw.run('rm -rf ~/magma-deps')


@task
def upload_to_aws(c):
    # If a host list isn't specified, default to the magma vagrant vm
    with vagrant_connection(c, 'magma') as c_gw:
        pkg.upload_pkgs_to_aws(c_gw)


@task
def copy_packages(c):
    with vagrant_connection(c, 'magma') as c_gw:
        pkg.copy_packages(c_gw)


@task
def s1ap_setup_cloud(c):
    """ Prepare VMs for s1ap tests touching the cloud. """
    # Use the local cloud for integ tests
    with vagrant_connection(c, "magma") as c_gw:
        connect_gateway_to_cloud(c_gw, None, DEFAULT_CERT)

        # Update the gateway's streamer timeout and restart services
        c_gw.run("sudo mkdir -p /var/opt/magma/configs")
        _set_service_config_var(c_gw, 'streamer', 'reconnect_sec', 3)

        # Update the gateway's metricsd collect/sync intervals
        _set_service_config_var(c_gw, 'metricsd', 'collect_interval', 5)
        _set_service_config_var(c_gw, 'metricsd', 'sync_interval', 5)

        c_gw.run("sudo systemctl stop magma@*")
        c_gw.run("sudo systemctl restart magma@magmad")


@task
def open_orc8r_port_in_vagrant(c):
    """
    Add a line to Vagrantfile file to open 9445 port on Vagrant.
    Note that localhost request to 9443 will be sent to Vagrant vm.
    Remove this line manually if you intend to run orc8r on your host
    """
    cmd_yes_if_exists = """grep -q 'guest: 9443, host: 9443' Vagrantfile"""

    # Insert line after a specific line
    cmd_insert_line = \
        "awk '/config.vm.define :magma,/{print;print " \
        "\"    magma.vm.network \\\"forwarded_port\\\", " \
        "guest: 9443, host: 9443\"" \
        ";next}1' Vagrantfile >> Vagrantfile.bak00 && " \
        "cp Vagrantfile.bak00 Vagrantfile && rm Vagrantfile.bak00"

    c.run(f"{cmd_yes_if_exists} || ({cmd_insert_line})")


def _redirect_feg_agw_to_vagrant_orc8r(c_gw):
    """
    Modifies feg docker-compose.override.yml hosts and AGW /etc/hosts
    to point to localhost when Orc8r runs inside Vagrant
    """
    # This is only run in CI:
    # on macos
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

    start_all_cmd = "fab start-all"
    if orc8r_on_vagrant:
        start_all_cmd += " --orc8r-on-vagrant"
        # modify dns entries to find Orc8r from inside Vagrant
        with vagrant_connection(c, 'magma') as c_gw:
            _redirect_feg_agw_to_vagrant_orc8r(c_gw)

    with c.cd(FEG_INTEG_TEST_ROOT):
        c.run(start_all_cmd)

        if orc8r_on_vagrant:
            print("Wait for orc8r to be available")
            sleep(60)

        c.run("fab configure-orc8r")
        sleep(20)
        c.run("fab test-connectivity --timeout=200")

    # back at AGW_ROOT
    vagrant_setup(
        c, 'magma_trfserver', destroy_vm, force_provision=provision_vm,
    )

    c_test, vm_data = vagrant_setup(
        c, 'magma_test', destroy_vm, force_provision=provision_vm,
    )
    sleep(20)
    # run this on the host, not on the vm, as it will connect to the vm via ssh
    _run_integ_tests(c, vm_data, test_mode="federated_integ_test")


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


def _setup_vm(c, host, name, ansible_role, ansible_file, destroy_vm, provision_vm):
    if not host:
        connection, host_data = vagrant_setup(
            c, name, destroy_vm, force_provision=provision_vm,
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
        provision_vm,
):
    gateway_connection, _ = _setup_vm(
        c, gateway_host, name, ansible_role, ansible_file, destroy_vm,
        provision_vm,
    )
    if gateway_host is None:
        gateway_ip = GATEWAY_IP_ADDRESS
    else:
        gateway_ip = gateway_host.split('@')[1].split(':')[0]
    return gateway_connection, gateway_ip


@task
def integ_test(
    c, gateway_host=None, test_host=None, trf_host=None,
    destroy_vm=True, provision_vm=True,
):
    """
    Run the integration tests. This defaults to running on local vagrant
    machines, but can also be pointed to an arbitrary host (e.g. amazon) by
    passing "address:port" as arguments

    gateway_host: The ssh address string of the machine to run the gateway
        services on. Formatted as "host:port". If not specified, defaults to
        the `magma` vagrant box.

    test_host: The ssh address string of the machine to run the tests on
        on. Formatted as "host:port". If not specified, defaults to the
        `magma_test` vagrant box.

    trf_host: The ssh address string of the machine to run the TrafficServer
        on. Formatted as "host:port". If not specified, defaults to the
        `magma_trfserver` vagrant box.
    """

    # Set up the gateway: use the provided gateway if given, else default to the
    # vagrant machine
    c_gw, gateway_ip = _setup_gateway(
        c, gateway_host, "magma", "dev", "magma_dev.yml", destroy_vm,
        provision_vm,
    )
    with c_gw:
        _build_magma(c_gw)
        _start_gateway(c_gw)

    # Set up the trfserver: use the provided trfserver if given, else default to the
    # vagrant machine
    c_trf, _ = _setup_vm(
        c, trf_host, "magma_trfserver", "trfserver", "magma_trfserver.yml",
        destroy_vm, provision_vm,
    )
    with c_trf:
        _start_trfserver(c_trf)

    # Run the tests: use the provided test machine if given, else default to
    # the vagrant machine
    c_test, test_host_data = _setup_vm(
        c, test_host, "magma_test", "test", "magma_test.yml", destroy_vm,
        provision_vm,
    )
    # run this on the host, not on the vm, as it will connect to the vm via ssh
    _run_integ_tests(c, test_host_data, gateway_ip=gateway_ip)


@task
def integ_test_deb_installation(
    c, gateway_host=None, test_host=None, trf_host=None,
    destroy_vm=True, provision_vm=True,
):
    """
    Run the integration tests. This defaults to running on local vagrant
    machines, but can also be pointed to an arbitrary host (e.g. amazon) by
    passing "address:port" as arguments

    gateway_host: The ssh address string of the machine to run the gateway
        services on. Formatted as "host:port". If not specified, defaults to
        the `magma_deb` vagrant box.

    test_host: The ssh address string of the machine to run the tests on
        on. Formatted as "host:port". If not specified, defaults to the
        `magma_test` vagrant box.

    trf_host: The ssh address string of the machine to run the TrafficServer
        on. Formatted as "host:port". If not specified, defaults to the
        `magma_trfserver` vagrant box.
    """

    # Set up the gateway: use the provided gateway if given, else default to the
    # vagrant machine
    c_gw, gateway_ip = _setup_gateway(
        c, gateway_host, "magma_deb", "deb", "magma_deb.yml", destroy_vm,
        provision_vm,
    )
    with c_gw:
        _start_gateway(c_gw)

    # Set up the trfserver: use the provided trfserver if given, else default to the
    # vagrant machine
    c_trf, _ = _setup_vm(
        c, trf_host, "magma_trfserver", "trfserver", "magma_trfserver.yml",
        destroy_vm, provision_vm,
    )
    with c_trf:
        _start_trfserver(c_trf)

    # Run the tests: use the provided test machine if given, else default to
    # the vagrant machine
    c_test, test_host_data = _setup_vm(
        c, test_host, "magma_test", "test", "magma_test.yml", destroy_vm,
        provision_vm,
    )

    # run this on the host, not on the vm, as it will connect to the vm via ssh
    _run_integ_tests(c, test_host_data, gateway_ip=gateway_ip)


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
    c_gw, gateway_ip = _setup_gateway(
        c, gateway_host, "magma", "dev", "magma_dev.yml", destroy_vm,
        provision_vm,
    )
    with c_gw:
        _start_gateway_containerized(c_gw, docker_registry)

    # Set up the trfserver: use the provided trfserver if given, else default to the
    # vagrant machine
    c_trf, _ = _setup_vm(
        c, trf_host, "magma_trfserver", "trfserver", "magma_trfserver.yml",
        destroy_vm, provision_vm,
    )
    with c_trf:
        _start_trfserver(c_trf)

    # Run the tests: use the provided test machine if given, else default to
    # the vagrant machine
    c_test, test_host_data = _setup_vm(
        c, test_host, "magma_test", "test", "magma_test.yml", destroy_vm,
        provision_vm,
    )
    # run this on the host, not on the vm, as it will connect to the vm via ssh
    _run_integ_tests(c, test_host_data, gateway_ip=gateway_ip, test_mode=test_mode, tests=tests)


def _start_gateway_containerized(c_gw, docker_registry=None):
    """ Starts the containerized AGW """
    with c_gw.cd(AGW_PYTHON_ROOT):
        c_gw.run('make buildenv')

    with c_gw.cd(AGW_ROOT):
        c_gw.run(
            'for component in redis nghttpx td-agent-bit; do cp "${MAGMA_ROOT}"'
            '/{orc8r,lte}/gateway/configs/templates/${component}.conf.template;'
            ' done',
        )

    c_gw.run('sudo systemctl start magma_dp@envoy')

    with c_gw.cd(AGW_ROOT + "/docker"):
        # The `docker-compose up` times are machine dependent, such that a
        # retry is needed here for resilience.
        run_with_retry(
            c_gw, f'DOCKER_REGISTRY={docker_registry} docker compose'
            f' --compatibility -f docker-compose.yaml up -d --quiet-pull',
        )


@task
def run_with_retry(c_gw, command, retries=10):
    iteration = 0
    while iteration < retries:
        iteration += 1
        try:
            c_gw.run(command)
            break
        except:
            print(f"ERROR: Failed on retry {iteration} of \n$ {command}\n")
            sleep(3)
    else:
        c_gw.run("docker ps")  # It is _not_ docker compose by intention to see the container ID.
        raise Exception(f"ERROR: Failed after {retries} retries of \n$ {command}")


@task
def get_test_summaries(
        c,
        dst_path="/tmp",
        integration_tests=False,
        sudo_tests=False,
        dev_vm_name="magma",
):
    c.run('mkdir -p ' + dst_path)

    if sudo_tests:
        _get_test_summaries_from_vm(c, dst_path, dev_vm_name)
    if integration_tests:
        _get_test_summaries_from_vm(c, dst_path, "magma_test")


def _get_test_summaries_from_vm(c, dst_path, vm_name):
    results_folder = "test_results"
    results_dir = "/var/tmp/"
    with vagrant_connection(c, vm_name) as c_gw:
        if c_gw.run(
            f"test -e {results_dir}{results_folder}", warn=True,
        ).ok:
            # Fix the permissions on the files -- they have permissions 000
            # otherwise
            c_gw.run(f'sudo chmod 755 {results_dir}{results_folder}')
            _get_folder(c_gw, results_folder, results_dir, dst_path)


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
def build_and_start_magma(c, destroy_vm=False, provision_vm=False):
    """
    Build Magma AGW and starts magma
    Args:
        destroy_vm: if set to True it will destroy Magma Vagrant VM
        provision_vm: if set to true it will reprovision Magma VM

    Returns:

    """
    with vagrant_connection(
        c, 'magma', destroy_vm=destroy_vm, force_provision=provision_vm,
    ) as c_gw:
        c_gw.run('sudo service magma@* stop')
        _build_magma(c_gw)
        c_gw.run('sudo service magma@magmad start')


@task
def build_and_start_magma_trf(c, destroy_vm=False, provision_vm=False):
    c_trf = vagrant_connection(
        c, 'magma_trfserver', destroy_vm=destroy_vm, force_provision=provision_vm,
    )
    with c_trf:
        _start_trfserver(c_trf)


@task
def start_magma(c, destroy_vm=False, provision_vm=False):
    with vagrant_connection(
        c, 'magma', destroy_vm=destroy_vm, force_provision=provision_vm,
    ) as c_gw:
        c_gw.run('sudo service magma@magmad start')


@task
def build_test_vms(c, provision_vm=False, destroy_vm=False):
    vagrant_connection(
        c, 'magma_trfserver', destroy_vm=destroy_vm, force_provision=provision_vm,
    )

    c_test = vagrant_connection(
        c, 'magma_test', destroy_vm=destroy_vm, force_provision=provision_vm,
    )


def _copy_out_c_execs_in_magma_vm(c_gw):
    exec_paths = [
        '/usr/local/bin/sessiond', '/usr/local/bin/mme',
        '/usr/local/sbin/sctpd', '/usr/local/bin/connectiond',
        '/usr/local/bin/liagentd',
    ]
    dest_path = '~/magma-packages/executables'
    c_gw.run('mkdir -p ' + dest_path, warn=True)
    for exec_path in exec_paths:
        if not c_gw.run(f"test -e {exec_path}", warn=True).ok:
            print(exec_path + " does not exist")
            continue
        c_gw.run('cp ' + exec_path + ' ' + dest_path, warn=True)


def _build_magma(c_gw):
    """
    Build magma on AGW
    """
    with c_gw.cd(AGW_ROOT):
        c_gw.run(
            env={'PATH': '$PATH:/usr/local/go/bin:/home/vagrant/go/bin'},
            command='make',
        )


def _start_gateway(c_gw):
    """ Starts the gateway """
    c_gw.run('sudo service magma@magmad start')


def _set_service_config_var(c_gw, service, var_name, value):
    """ Sets variable in config file by value """
    c_gw.run(
        f"echo '{var_name}: {str(value)}'"
        f" | sudo tee -a /var/opt/magma/configs/{service}.yml",
    )


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
