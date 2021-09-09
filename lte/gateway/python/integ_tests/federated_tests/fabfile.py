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
from distutils.util import strtobool

from fabric.api import cd, env, execute, hide, run

sys.path.append('../../../../../orc8r')


magma_path = "../../../../../"
orc8_docker_path = magma_path + "orc8r/cloud/docker/"
feg_path = magma_path + "feg/gateway/"
feg_docker_path = feg_path + "docker/"
agw_path = magma_path + "lte/gateway/"

vagrant_agw_path = "~/lte/gateway"


def build_all(clear_orc8r='False', provision_vm='False'):
    """ Builds AGW, FEG and Orc8r and starts them

        clear_orc8r: removes all contents from orc8r database like gw configs

        force_provision: forces the reprovision of the magma VM
    """

    clear_orc8r = bool(strtobool(clear_orc8r))
    execute(build_agw, provision_vm=provision_vm)
    execute(build_feg)
    execute(build_orc8r)
    execute(start_orc8r)
    if clear_orc8r:
        execute(clear_orc8r)
        execute(start_orc8r)

    execute(build_test_vm, provision_vm=provision_vm)
    execute(build_magma_trf, provision_vm=provision_vm)

    # do this at the end to make sure or8cr is running
    execute(configure_orc8r)


def build_orc8r():
    """build_orc8r builds orc8r locally on the host VM"""
    subprocess.check_call('./build.py -a', shell=True, cwd=orc8_docker_path)


def start_orc8r():
    subprocess.check_call(['./run.py'], shell=True, cwd=orc8_docker_path)


def configure_orc8r():
    print(f'#### Configuring orc8r ####')
    subprocess.check_call(
        'fab --fabfile=dev_tools.py register_federated_vm',
        shell=True, cwd=agw_path,
    )
    subprocess.check_call(
        'fab --fabfile=dev_tools.py register_feg_on_magma_vm', shell=True, cwd=agw_path,
    )


def reconfigure_orc8r():
    print(f'#### Removing VMs from orc8r ####')
    subprocess.check_call(
        'fab --fabfile=dev_tools.py deregister_federated_agw', shell=True, cwd=agw_path,
    )
    subprocess.check_call(
        'fab --fabfile=dev_tools.py deregister_feg_gw_on_magma_vm',
        shell=True, cwd=agw_path,
    )
    execute(configure_orc8r)


def clear_orc8r():
    print(f'#### Clearing swagger database from Orc8r ####')
    subprocess.check_call(['./run.py --clear_db'], shell=True, cwd=orc8_docker_path)
    print(
        f'#### Remember you may need to delete '
        f'gateway certs from the AGW and FEG ####',
    )


def build_agw(provision_vm='False'):
    print(f'#### Building AGW ####')
    subprocess.check_call('vagrant up magma', shell=True, cwd=agw_path)
    subprocess.check_call(
        'fab build_and_start_magma:provision_vm=%s'
        % provision_vm, shell=True, cwd=agw_path,
    )


def build_feg(provision_vm='False'):
    print(f'#### Building FEG ####')
    subprocess.check_call(
        'docker-compose down', shell=True,
        cwd=feg_docker_path,
    )
    subprocess.check_call(
        'docker-compose build', shell=True,
        cwd=feg_docker_path,
    )
    subprocess.check_call(
        'docker-compose up -d', shell=True,
        cwd=feg_docker_path,
    )


def build_test_vm(provision_vm='False'):
    print(f'#### Building test vm ####')
    subprocess.check_call('vagrant up magma_test', shell=True, cwd=agw_path)
    subprocess.check_call(
        'fab make_integ_tests:provision_vm=%s'
        % provision_vm, shell=True, cwd=agw_path,
    )


def build_magma_trf(provision_vm='False'):
    print(f'#### Building Traffic vm ####')
    subprocess.check_call('vagrant up magma_trf', shell=True, cwd=agw_path)
    subprocess.check_call(
        'fab build_and_start_magma_trf:provision_vm=%s'
        % provision_vm, shell=True, cwd=agw_path,
    )


def start_all(provision_vm='False'):
    subprocess.check_call(['./run.py'], shell=True, cwd=orc8_docker_path)
    subprocess.check_call('docker-compose up', shell=True, cwd=feg_docker_path)
    subprocess.check_call(
        'fab start_magma:provision_vm=%s' % provision_vm,
        shell=True, cwd=agw_path,
    )


def stop_all():
    subprocess.check_call(
        ['./run.py --down'], shell=True,
        cwd=orc8_docker_path,
    )
    subprocess.check_call(
        'docker-compose down', shell=True,
        cwd=feg_docker_path,
    )
    subprocess.check_call(
        'vagrant halt magma', shell=True,
        cwd=agw_path,
    )
