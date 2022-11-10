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
from datetime import datetime
from time import sleep

from fabric.api import cd, env, execute, lcd, local, run, settings, sudo
from fabric.contrib.files import exists
from fabric.operations import get
from fabric.utils import error, fastprint, puts

sys.path.append('../../orc8r')
import tools.fab.pkg as pkg
from tools.fab.dev_utils import connect_gateway_to_cloud
from tools.fab.hosts import ansible_setup, split_hoststring, vagrant_setup
from tools.fab.python_utils import strtobool
from tools.fab.vagrant import setup_env_vagrant

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

# Look for keys as specified in our ~/.ssh/config
env.use_ssh_config = True
# Disable ssh known hosts to resolve key errors
# with multiple vagrant boxes in use.
env.disable_known_hosts = True


def dev():
    env.debug_mode = True


def release():
    """Set debug_mode to False, should be used for producing a production AGW package"""
    env.debug_mode = False


def package(
    all_deps='False',
    cert_file=DEFAULT_CERT, proxy_config=DEFAULT_PROXY,
    destroy_vm='False',
    vm='magma', os="ubuntu",
):
    """ Builds the magma package """
    all_deps = strtobool(all_deps)
    destroy_vm = strtobool(destroy_vm)

    # If a host list isn't specified, default to the magma vagrant vm
    if not env.hosts:
        vagrant_setup(vm, destroy_vm=destroy_vm)

    if not hasattr(env, 'debug_mode'):
        error(
            "Error: The Deploy target isn't specified. Specify one with\n\n" +
            "\tfab [dev|release] package",
        )

    hash = pkg.get_commit_hash()
    commit_count = pkg.get_commit_count()
    puts('Uninstalling dev dependencies of the VM')
    run('sudo pip uninstall --yes mypy-protobuf grpcio-tools grpcio protobuf')

    with cd('~/magma/lte/gateway'):
        run('mkdir -p ~/magma-deps')
        puts(
            'Generating lte/setup.py and orc8r/setup.py magma dependency packages',
        )
        run(
            './release/pydep finddep --install-from-repo -b --build-output '
            + '~/magma-deps'
            + (' -l ./release/magma.lockfile.%s' % os)
            + ' python/setup.py'
            + (' %s/setup.py' % ORC8R_AGW_PYTHON_ROOT),
        )

        puts('Building magma package, picking up commit %s...' % hash)
        run('make clean')
        build_type = "Debug" if env.debug_mode else "RelWithDebInfo"

        run(
            './release/build-magma.sh -h %s --commit-count %s -t %s --cert %s --proxy %s --os %s' %
            (hash, commit_count, build_type, cert_file, proxy_config, os),
        )

        run('rm -rf ~/magma-packages')
        run('mkdir -p ~/magma-packages')
        with settings(warn_only=True):
            run('cp -f ~/magma-deps/*.deb ~/magma-packages')
        run('mv *.deb ~/magma-packages')

        with cd('release'):
            mirrored_packages_file = 'mirrored_packages'
            if os == "ubuntu":
                mirrored_packages_file += '_focal'
            if vm and vm.startswith('magma_'):
                mirrored_packages_file += vm[5:]

            run(
                'cat {}'.format(mirrored_packages_file)
                + ' | xargs -I% sudo aptitude download -q2 %',
            )
            run('cp *.deb ~/magma-packages')
            run('sudo rm -f *.deb')

        if all_deps:
            pkg.download_all_pkgs()
            run('cp /var/cache/apt/archives/*.deb ~/magma-packages')

        # Copy out C executables into magma-packages as well
        _copy_out_c_execs_in_magma_vm()


def openvswitch(destroy_vm='False', destdir='~/magma-packages/'):
    destroy_vm = strtobool(destroy_vm)
    # If a host list isn't specified, default to the magma vagrant vm
    if not env.hosts:
        vagrant_setup('magma', destroy_vm=destroy_vm)
    run('~/magma/third_party/gtp_ovs/ovs-gtp-patches/2.15/build.sh ' + destdir)


def depclean():
    '''Remove all generated packaged for dependencies'''
    # If a host list isn't specified, default to the magma vagrant vm
    if not env.hosts:
        setup_env_vagrant()
    run('rm -rf ~/magma-deps')


def upload_to_aws():
    # If a host list isn't specified, default to the magma vagrant vm
    if not env.hosts:
        setup_env_vagrant()

    pkg.upload_pkgs_to_aws()


def copy_packages():
    if not env.hosts:
        setup_env_vagrant()
    pkg.copy_packages()


def s1ap_setup_cloud():
    """ Prepare VMs for s1ap tests touching the cloud. """
    # Use the local cloud for integ tests
    setup_env_vagrant()
    connect_gateway_to_cloud(None, DEFAULT_CERT)

    # Update the gateway's streamer timeout and restart services
    run("sudo mkdir -p /var/opt/magma/configs")
    _set_service_config_var('streamer', 'reconnect_sec', 3)

    # Update the gateway's metricsd collect/sync intervals
    _set_service_config_var('metricsd', 'collect_interval', 5)
    _set_service_config_var('metricsd', 'sync_interval', 5)

    run("sudo systemctl stop magma@*")
    run("sudo systemctl restart magma@magmad")


def open_orc8r_port_in_vagrant():
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

    local("{} || ({})".format(cmd_yes_if_exists, cmd_insert_line))


def redirect_feg_agw_to_vagrant_orc8r():
    """
    Modifies feg docker-compose.override.yml hosts and AGW /etc/hosts
    to point to localhost when Orc8r runs on inside Vagrant
    """
    local(
        "sed -i '' 's/:10.0.2.2/:127.0.0.1/' {}/docker-compose.override.yml"
        .format(FEG_INTEG_TEST_DOCKER_ROOT),
    )

    vagrant_setup(
        'magma', destroy_vm=False, force_provision=False,
    )
    sudo("sed -i 's/10.0.2.2/127.0.0.1/' /etc/hosts")


def federated_integ_test(
        build_all='False', clear_orc8r='False', provision_vm='False',
        destroy_vm='False', orc8r_on_vagrant='False',
):
    build_all = strtobool(build_all)
    clear_orc8r = strtobool(clear_orc8r)
    provision_vm = strtobool(provision_vm)
    destroy_vm = strtobool(destroy_vm)
    orc8r_on_vagrant = strtobool(orc8r_on_vagrant)

    if orc8r_on_vagrant:
        # Modify Vagrantfile to allow access to Orc8r running inside Vagrant
        execute(open_orc8r_port_in_vagrant)

    with lcd(FEG_INTEG_TEST_ROOT):
        if build_all:
            local(
                "fab build_all:clear_orc8r={},provision_vm={},"
                "orc8r_on_vagrant={}".format(
                    clear_orc8r,
                    provision_vm,
                    orc8r_on_vagrant,
                ),
            )

        if orc8r_on_vagrant:
            # modify dns entries to find Orc8r from inside Vagrant
            execute(redirect_feg_agw_to_vagrant_orc8r)

        local("fab start_all:orc8r_on_vagrant={}".format(orc8r_on_vagrant))

        if orc8r_on_vagrant:
            fastprint("Wait for orc8r to be available")
            sleep(60)

        local("fab configure_orc8r")
        sleep(20)
        local("fab test_connectivity:timeout=200")

    vagrant_setup(
        'magma_trfserver', destroy_vm, force_provision=provision_vm,
    )

    vagrant_setup(
        'magma_test', destroy_vm, force_provision=provision_vm,
    )
    execute(_make_integ_tests)
    sleep(20)
    execute(_run_integ_tests, test_mode="federated_integ_test")


def _modify_for_bazel_services():
    """ Modify the service definitions to use the bazel-built executables """
    run(r"sudo cp $MAGMA_ROOT/lte/gateway/deploy/roles/magma/files/systemd_bazel/* /etc/systemd/system/")
    run("sudo systemctl daemon-reload")


def provision_magma_dev_vm(
    gateway_host=None, destroy_vm='True', provision_vm='True',
):
    """
    Prepare to run the integration tests on the bazel build services.
    This defaults to running on local vagrant machines, but can also be
    pointed to an arbitrary host (e.g. amazon) by passing "address:port"
    as arguments

    gateway_host: The ssh address string of the machine to run the gateway
        services on. Formatted as "host:port". If not specified, defaults to
        the `magma` vagrant box.
    """
    destroy_vm = strtobool(destroy_vm)
    provision_vm = strtobool(provision_vm)

    if not gateway_host:
        gateway_host = vagrant_setup(
            'magma', destroy_vm, force_provision=provision_vm,
        )
    else:
        ansible_setup(gateway_host, "dev", "magma_dev.yml")


def bazel_integ_test_post_build(
    gateway_host=None, test_host=None, trf_host=None,
    destroy_vm='True', provision_vm='True',
):
    """
    Run the integration tests on the bazel build services. This defaults to
    running on local vagrant machines, but can also be pointed to an
    arbitrary host (e.g. amazon) by passing "address:port" as arguments

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
    destroy_vm = strtobool(destroy_vm)
    provision_vm = strtobool(provision_vm)

    gateway_ip = '192.168.60.142'

    if not gateway_host:
        gateway_host = vagrant_setup(
            'magma', False, force_provision=False,
        )
    else:
        ansible_setup(gateway_host, "dev", "magma_dev.yml")
        gateway_ip = gateway_host.split('@')[1].split(':')[0]

    execute(_modify_for_bazel_services)
    execute(_start_gateway)

    # Set up the trfserver: use the provided trfserver if given, else default to the
    # vagrant machine
    if not trf_host:
        trf_host = vagrant_setup(
            'magma_trfserver', destroy_vm, force_provision=provision_vm,
        )
    else:
        ansible_setup(trf_host, "trfserver", "magma_trfserver.yml")
    execute(_start_trfserver)

    # Run the tests: use the provided test machine if given, else default to
    # the vagrant machine
    if not test_host:
        test_host = vagrant_setup(
            'magma_test', destroy_vm, force_provision=provision_vm,
        )
    else:
        ansible_setup(test_host, "test", "magma_test.yml")

    execute(_make_integ_tests)
    execute(_run_integ_tests, gateway_ip)

    if not gateway_host:
        setup_env_vagrant()
    else:
        env.hosts = [gateway_host]


def _setup_vm(host, name, ansible_role, ansible_file, destroy_vm, provision_vm):
    ip = None
    if not host:
        host = vagrant_setup(
            name, destroy_vm, force_provision=provision_vm,
        )
    else:
        ansible_setup(host, ansible_role, ansible_file)
        ip = host.split('@')[1].split(':')[0]
    return host, ip


def _setup_gateway(gateway_host, name, ansible_role, ansible_file, destroy_vm, provision_vm):
    gateway_host, gateway_ip = _setup_vm(gateway_host, name, ansible_role, ansible_file, destroy_vm, provision_vm)
    if gateway_ip is None:
        gateway_ip = GATEWAY_IP_ADDRESS
    return gateway_host, gateway_ip


def integ_test(
    gateway_host=None, test_host=None, trf_host=None,
    destroy_vm='True', provision_vm='True',
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

    destroy_vm = strtobool(destroy_vm)
    provision_vm = strtobool(provision_vm)

    # Set up the gateway: use the provided gateway if given, else default to the
    # vagrant machine
    gateway_host, gateway_ip = _setup_gateway(gateway_host, "magma", "dev", "magma_dev.yml", destroy_vm, provision_vm)
    execute(_build_magma)
    execute(_start_gateway)

    # Set up the trfserver: use the provided trfserver if given, else default to the
    # vagrant machine
    _setup_vm(trf_host, "magma_trfserver", "trfserver", "magma_trfserver.yml", destroy_vm, provision_vm)
    execute(_start_trfserver)

    # Run the tests: use the provided test machine if given, else default to
    # the vagrant machine
    _setup_vm(test_host, "magma_test", "test", "magma_test.yml", destroy_vm, provision_vm)
    execute(_make_integ_tests)
    execute(_run_integ_tests, gateway_ip)

    if not gateway_host:
        setup_env_vagrant()
    else:
        env.hosts = [gateway_host]


def integ_test_deb_installation(
    gateway_host=None, test_host=None, trf_host=None,
    destroy_vm='True', provision_vm='True',
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

    destroy_vm = strtobool(destroy_vm)
    provision_vm = strtobool(provision_vm)

    # Set up the gateway: use the provided gateway if given, else default to the
    # vagrant machine
    _, gateway_ip = _setup_gateway(gateway_host, "magma_deb", "deb", "magma_deb.yml", destroy_vm, provision_vm)
    execute(_start_gateway)

    # Set up the trfserver: use the provided trfserver if given, else default to the
    # vagrant machine
    _setup_vm(trf_host, "magma_trfserver", "trfserver", "magma_trfserver.yml", destroy_vm, provision_vm)
    execute(_start_trfserver)

    # Run the tests: use the provided test machine if given, else default to
    # the vagrant machine
    _setup_vm(test_host, "magma_test", "test", "magma_test.yml", destroy_vm, provision_vm)
    execute(_make_integ_tests)
    execute(_run_integ_tests, gateway_ip)


def integ_test_containerized(
        gateway_host=None, test_host=None, trf_host=None,
        destroy_vm='True', provision_vm='True',
        test_mode='integ_test_containerized',
        tests='',
):
    """
    Run the integration tests against the containerized AGW.
    Other than that the same as `integ_test`.
    """

    destroy_vm = bool(strtobool(destroy_vm))
    provision_vm = bool(strtobool(provision_vm))

    # Set up the gateway: use the provided gateway if given, else default to the
    # vagrant machine
    gateway_host, gateway_ip = _setup_gateway(gateway_host, "magma", "dev", "magma_dev.yml", destroy_vm, provision_vm)
    execute(_start_gateway_containerized)

    # Set up the trfserver: use the provided trfserver if given, else default to the
    # vagrant machine
    _setup_vm(trf_host, "magma_trfserver", "trfserver", "magma_trfserver.yml", destroy_vm, provision_vm)
    execute(_start_trfserver)

    # Run the tests: use the provided test machine if given, else default to
    # the vagrant machine
    _setup_vm(test_host, "magma_test", "test", "magma_test.yml", destroy_vm, provision_vm)
    execute(_make_integ_tests)
    execute(_run_integ_tests, gateway_ip, test_mode=test_mode, tests=tests)


def _start_gateway_containerized():
    """ Starts the containerized AGW """
    with cd(AGW_PYTHON_ROOT):
        run('make buildenv')

    with cd(AGW_ROOT):
        run('for component in redis nghttpx td-agent-bit; do cp "${MAGMA_ROOT}"/{orc8r,lte}/gateway/configs/templates/${component}.conf.template; done')

    with cd(AGW_ROOT + "/docker"):
        # The `docker-compose up` times are machine dependent, such that a retry is needed here for resilience.
        run_with_retry('DOCKER_REGISTRY=%s docker-compose -f docker-compose.yaml up -d --quiet-pull' % (env.DOCKER_REGISTRY))


def run_with_retry(command, retries=10):
    iteration = 0
    while iteration < retries:
        iteration += 1
        try:
            run(command)
            break
        except:
            fastprint(f"ERROR: Failed on retry {iteration} of \n$ {command}\n")
            sleep(3)
    else:
        run("docker ps")  # It is _not_ docker-compose by intention to see the container ID.
        raise Exception(f"ERROR: Failed after {retries} retries of \n$ {command}")


def get_test_summaries(
        gateway_host=None,
        test_host=None,
        dst_path="/tmp",
        integration_tests='True',
        sudo_tests='True',
        dev_vm_name="magma",
):
    local('mkdir -p ' + dst_path)

    vm_name_to_yaml = {
        "magma": "magma_dev.yml",
        "magma_deb": "magma_deb.yml",
    }

    if strtobool(sudo_tests):
        _switch_to_vm_no_provision(gateway_host, dev_vm_name, vm_name_to_yaml[dev_vm_name])
        with settings(warn_only=True):
            get(remote_path=TEST_SUMMARY_GLOB, local_path=dst_path)
    if strtobool(integration_tests):
        _switch_to_vm_no_provision(test_host, "magma_test", "magma_test.yml")
        with settings(warn_only=True):
            get(remote_path=TEST_SUMMARY_GLOB, local_path=dst_path)


def get_test_logs(
    gateway_host=None,
    gateway_host_name='magma',
    test_host=None,
    trf_host=None,
    dst_path="/tmp/build_logs.tar.gz",
):
    """
    Download the relevant magma logs from the given gateway and test machines.
    Place the logs in a path specified in 'dst_path' or
    "/tmp/build_logs.tar.gz" by default.

    Args:
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
    local('rm -rf /tmp/build_logs')
    local('mkdir /tmp/build_logs')
    local('mkdir /tmp/build_logs/dev')
    local('mkdir /tmp/build_logs/test')
    local('mkdir /tmp/build_logs/trfserver')
    dev_files = [
        '/var/log/mme.log',
        '/var/log/MME.magma*log*',
        '/var/log/syslog',
        '/var/log/envoy.log',
        '/var/log/openvswitch/ovs*.log',
    ]
    test_files = ['/var/log/syslog', '/tmp/fw/']
    trf_files = ['/home/vagrant/trfserver.log']

    # Set up to enter the gateway host
    env.host_string = gateway_host
    if not gateway_host:
        setup_env_vagrant(gateway_host_name)
        gateway_host = env.hosts[0]
    (env.user, _, _) = split_hoststring(gateway_host)

    # Don't fail if the logs don't exists
    for p in dev_files:
        with settings(warn_only=True):
            get(
                remote_path=p, local_path='/tmp/build_logs/dev/',
                use_sudo=True,
            )

    # Set up to enter the trfserver host
    env.host_string = trf_host
    if not trf_host:
        setup_env_vagrant("magma_trfserver")
        trf_host = env.hosts[0]
    (env.user, _, _) = split_hoststring(trf_host)

    # Don't fail if the logs don't exists
    for p in trf_files:
        with settings(warn_only=True):
            get(
                remote_path=p, local_path='/tmp/build_logs/trfserver/',
                use_sudo=True,
            )

    # Set up to enter the test host
    env.host_string = test_host
    if not test_host:
        setup_env_vagrant("magma_test")
        test_host = env.hosts[0]
    (env.user, _, _) = split_hoststring(test_host)

    # Fix the permissions on the fw directory -- it has permissions 000
    # otherwise
    with settings(warn_only=True):
        run('sudo chmod 755 /tmp/fw')

    # Don't fail if the logs don't exists
    for p in test_files:
        with settings(warn_only=True):
            get(
                remote_path=p, local_path='/tmp/build_logs/test/',
                use_sudo=True,
            )

    local("tar -czvf /tmp/build_logs.tar.gz /tmp/build_logs/*")
    local(f'mv /tmp/build_logs.tar.gz {dst_path}')
    local('rm -rf /tmp/build_logs')


def build_and_start_magma(gateway_host=None, destroy_vm='False', provision_vm='False'):
    """
    Build Magma AGW and starts magma
    Args:
        gateway_host: name of host in case is not Vagrant
        destroy_vm: if set to True it will destroy Magma Vagrant VM
        provision_vm: if set to true it will reprovision Magma VM

    Returns:

    """
    provision_vm = strtobool(provision_vm)
    destroy_vm = strtobool(destroy_vm)
    if gateway_host:
        ansible_setup(gateway_host, 'dev', 'magma_dev.yml')
    else:
        vagrant_setup('magma', destroy_vm, provision_vm)
    sudo('service magma@* stop')
    execute(_build_magma)
    sudo('service magma@magmad start')


def make_integ_tests(test_host=None, destroy_vm='False', provision_vm='False'):
    destroy_vm = strtobool(destroy_vm)
    provision_vm = strtobool(provision_vm)
    if not test_host:
        vagrant_setup('magma_test', destroy_vm, force_provision=provision_vm)
    else:
        ansible_setup(test_host, "test", "magma_test.yml")
    execute(_make_integ_tests)


def build_and_start_magma_trf(test_host=None, destroy_vm='False', provision_vm='False'):
    destroy_vm = strtobool(destroy_vm)
    provision_vm = strtobool(provision_vm)
    if not test_host:
        vagrant_setup('magma_trfserver', destroy_vm, force_provision=provision_vm)
    else:
        ansible_setup(test_host, "test", "magma_test.yml")
    execute(_start_trfserver)


def start_magma(test_host=None, destroy_vm='False', provision_vm='False'):
    destroy_vm = strtobool(destroy_vm)
    provision_vm = strtobool(provision_vm)
    if not test_host:
        vagrant_setup('magma', destroy_vm, force_provision=provision_vm)
    else:
        ansible_setup(test_host, "test", "magma_test.yml")
    sudo('service magma@magmad start')


def build_test_vms(provision_vm='False', destroy_vm='False'):
    destroy_vm = strtobool(destroy_vm)
    provision_vm = strtobool(provision_vm)
    vagrant_setup(
        'magma_trfserver', destroy_vm, force_provision=provision_vm,
    )

    vagrant_setup(
        'magma_test', destroy_vm, force_provision=provision_vm,
    )
    execute(_make_integ_tests)


def _copy_out_c_execs_in_magma_vm():
    with settings(warn_only=True):
        exec_paths = [
            '/usr/local/bin/sessiond', '/usr/local/bin/mme',
            '/usr/local/sbin/sctpd', '/usr/local/bin/connectiond',
            '/usr/local/bin/liagentd',
        ]
        dest_path = '~/magma-packages/executables'
        run('mkdir -p ' + dest_path)
        for exec_path in exec_paths:
            if not exists(exec_path):
                fastprint(exec_path + " does not exist")
                continue
            run('cp ' + exec_path + ' ' + dest_path)


def _build_magma():
    """
    Build magma on AGW
    """
    with cd(AGW_ROOT):
        run('make')


def _start_gateway():
    """ Starts the gateway """
    run('sudo service magma@magmad start')


def _set_service_config_var(service, var_name, value):
    """ Sets variable in config file by value """
    run(
        "echo '%s: %s' | sudo tee -a /var/opt/magma/configs/%s.yml" % (
            var_name, str(value), service,
        ),
    )


def _start_trfserver():
    """ Starts the traffic gen server"""
    host = env.hosts[0].split(':')[0]
    port = env.hosts[0].split(':')[1]
    key = env.key_filename

    def _call_trfserver_ssh_command(cmd):
        local(
            f'ssh -f -i {key} -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -tt {host} -p {port} "{cmd}"',
        )

    _call_trfserver_ssh_command('sudo ethtool --offload eth1 rx off tx off')
    _call_trfserver_ssh_command('sudo ethtool --offload eth2 rx off tx off')
    _call_trfserver_ssh_command('nohup sudo /usr/local/bin/traffic_server.py 192.168.60.144 62462 > trfserver.log 2>&1')


def _make_integ_tests():
    """ Build and run the integration tests """

    with cd(AGW_PYTHON_ROOT):
        run('make')
    with cd(AGW_INTEG_ROOT):
        run('make')


def _run_integ_tests(gateway_ip='192.168.60.142', test_mode='integ_test', tests=''):
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
    host = env.hosts[0].split(':')[0]
    port = env.hosts[0].split(':')[1]
    key = env.key_filename

    # We do not have a proper shell, so the `magtivate` alias is not available.
    # We instead directly source the activate file.
    local(
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


def _switch_to_vm_no_provision(addr, host_name, ansible_file):
    if not addr:
        vagrant_setup(host_name, destroy_vm=False, force_provision=False)
    else:
        ansible_setup(addr, host_name, ansible_file, full_provision='false')
