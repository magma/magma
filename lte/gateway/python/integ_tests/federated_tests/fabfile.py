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

import subprocess
import sys

from fabric import Connection, task

sys.path.append('../../../../../orc8r')
import tools.fab.dev_utils as dev_utils
from tools.fab.hosts import vagrant_connection, vagrant_setup

magma_path = "../../../../../"
orc8_docker_path = magma_path + "orc8r/cloud/docker/"
agw_path = magma_path + "lte/gateway/"
feg_path = magma_path + "feg/gateway/"
feg_docker_path = feg_path + "docker/"
feg_docker_integ_test_path = agw_path + "python/integ_tests/federated_tests/docker/"
agw_vagrant_path = "magma/lte/gateway/"
feg_docker_integ_test_path_vagrant = agw_vagrant_path + "python/integ_tests/federated_tests/docker/"
feg_vagrant_path = "magma/feg/gateway/"
orc8r_vagrant_path = "magma/orc8r/cloud/docker/"


@task
def build_all_and_configure(c, clear_orc8r=False, provision_vm=False):
    """
    Build, start and configure AGW, FEG and Orc8r
    Args:
        c: fabric connection
        clear_orc8r: removes all contents from orc8r database like gw configs
        provision_vm: forces the re-provision of the magma VM
    """
    build_all(c, clear_orc8r, provision_vm)
    start_all(c)
    configure_orc8r(c)


@task
def build_all(c, clear_orc8r=False, provision_vm=False, orc8r_on_vagrant=False):
    """
    Build AGW, FEG and Orc8r
    Args:
        c: fabric connection
        clear_orc8r: removes all contents from orc8r database like gw configs
        provision_vm: forces the reprovision of the magma VM
        orc8r_on_vagrant: flag to build orc8r on vagrant or on host machine
    """
    # build all components
    build_orc8r(c, on_vagrant=orc8r_on_vagrant)
    print("#### Starting Orc8r (to generate certificates) ####")
    # need to start orc8r to generate certificates
    start_orc8r(c, on_vagrant=orc8r_on_vagrant)
    build_feg(c)
    install_agw(c, provision_vm=provision_vm)

    if clear_orc8r:
        clear_orc8r_db(c)

    # build other VMs
    build_test_vm(c)
    build_magma_trf(c, provision_vm=provision_vm)


@task
def build_orc8r(c, on_vagrant=False):
    """
    Build orc8r locally on the host VM
    """
    command = './build.py -a'
    _run_orc8r_command(c, command, on_vagrant)


@task
def start_orc8r(c, on_vagrant=False):
    """
    Start orc8r locally on Docker
    """
    command = './run.py'
    _run_orc8r_command(c, command, on_vagrant)


def _run_orc8r_command(c, command, on_vagrant):
    if not on_vagrant:
        subprocess.check_call(command, shell=True, cwd=orc8_docker_path)
    else:
        with c.cd(agw_path):
            with vagrant_connection(c, 'magma_deb') as c_gw:
                with c_gw.cd(orc8r_vagrant_path):
                    c_gw.run(command)


@task
def configure_orc8r(c, on_vagrant=False):
    """
    Configure orc8r with a federated AGW and FEG
    """
    print('#### Configuring orc8r ####')
    command_agw = 'fab register-federated-vm'
    command_feg = 'fab register-feg-gw'
    if not on_vagrant:
        subprocess.check_call(command_agw, shell=True, cwd=agw_path)
        subprocess.check_call(command_feg, shell=True, cwd=feg_path)
    else:
        with vagrant_connection(c, 'magma_deb') as c_gw:
            with c_gw.cd(agw_vagrant_path):
                c_gw.run(command_agw)
            with c_gw.cd(feg_vagrant_path):
                c_gw.run(command_feg)


@task
def clear_orc8r_db(c):
    """
    Delete orc8r database. Requires orc8r to be stopped
    """
    print('#### Clearing swagger database from Orc8r ####')
    subprocess.check_call(
        ['./run.py --clear-db'],
        shell=True, cwd=orc8_docker_path,
    )
    print(
        '#### Remember you may need to delete '
        'gateway certs from the AGW and FEG ####',
    )


@task
def install_agw(c, provision_vm=False):
    """
    Install a magma AGW debian package on the magma_deb Vagrant VM.
    Args:
        c: fabric connection.
       provision_vm: forces the reprovision of the magma VM
    """
    print('#### Installing AGW ####')
    with c.prefix("export INSTALL_DOCKER=true SETUP_TEST_CERTS=true"):
        vagrant_connection(c, 'magma_deb', force_provision=provision_vm)


@task
def start_agw(c, provision_vm=False):
    """
    start AGW on Vagrant VM
    """
    with vagrant_connection(c, 'magma_deb', force_provision=provision_vm) as c_gw:
        c_gw.run('sudo service magma@magmad start')


@task
def build_feg(c):
    """
    build FEG on magma_deb Vagrant vm using docker running in Vagrant
    """
    print('#### Building FEG on magma_deb Vagrant VM ####')
    with vagrant_connection(c, 'magma_deb') as c_gw:
        with c_gw.cd(feg_docker_integ_test_path_vagrant):
            c_gw.run('docker compose down')
            c_gw.run('docker compose --compatibility build')
            c_gw.run('./run.py')


@task
def start_feg(c):
    """
    start FEG on magma_deb Vagrant vm using docker running in vm
    """
    with vagrant_connection(c, 'magma_deb') as c_gw:
        with c_gw.cd(feg_docker_integ_test_path_vagrant):
            c_gw.run('./run.py')


@task
def build_test_vm(c):
    print('#### Building test vm ####')
    subprocess.check_call('vagrant up magma_test', shell=True, cwd=agw_path)


@task
def build_magma_trf(c, provision_vm=False):
    print('#### Building Traffic vm ####')
    subprocess.check_call('vagrant up magma_trfserver', shell=True, cwd=agw_path)
    cmd = 'fab build-and-start-magma-trf'
    if provision_vm:
        cmd += ' --provision-vm'
    subprocess.check_call(cmd, shell=True, cwd=agw_path)


@task
def start_all(c, provision_vm=False, orc8r_on_vagrant=False):
    """
    start AGW, FEG and Orc8r
    Args:
        c: fabric context
        provision_vm: forces the provision of the magma VM
        orc8r_on_vagrant: flag  to run orc8r on vagrant or on host machine
    """
    start_orc8r(c, on_vagrant=orc8r_on_vagrant)
    start_agw(c, provision_vm)
    start_feg(c)


@task
def test_connectivity(c, timeout=10):
    """
    Check if all running gateways have connectivity
    Args:
        timeout: amount of time the command will retry
    """
    # check AGW-Cloud connectivity
    print("### Checking AGW-Cloud connectivity ###")
    subprocess.check_call(
        f'fab check-agw-cloud-connectivity --timeout={timeout}',
        shell=True, cwd=agw_path,
    )

    # check FEG-cloud connectivity
    print("\n### Checking FEG-Cloud connectivity ###")
    with vagrant_connection(c, 'magma_deb') as c_gw:
        with c_gw.cd(feg_docker_integ_test_path_vagrant):
            dev_utils.run_remote_command_with_repetition(
                c_gw, 'docker compose exec -t magmad checkin_cli.py', timeout,
            )

    # check AGW-FEG connectivity
    print("\n### Checking AGW-FEG connectivity ###")
    subprocess.check_call(
        f'fab check-agw-feg-connectivity --timeout={timeout}',
        shell=True, cwd=agw_path,
    )
