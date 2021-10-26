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

from fabric.api import execute

sys.path.append('../../../../../orc8r')


magma_path = "../../../../../"
orc8_docker_path = magma_path + "orc8r/cloud/docker/"
agw_path = magma_path + "lte/gateway/"
feg_path = magma_path + "feg/gateway/"
feg_docker_path = feg_path + "docker/"
feg_docker_integ_test_path = agw_path + "python/integ_tests/federated_tests/docker/"


vagrant_agw_path = "~/lte/gateway"


def build_all(clear_orc8r='False', provision_vm='False'):
    """
    Build, start and configure AGW, FEG and Orc8r
    Args:
        clear_orc8r: removes all contents from orc8r database like gw configs
        provision_vm: forces the reprovision of the magma VM
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
    """
    build orc8r locally on the host VM
    """
    subprocess.check_call('./build.py -a', shell=True, cwd=orc8_docker_path)


def start_orc8r():
    """
    start orc8r locally on the host VM
    """
    subprocess.check_call(['./run.py'], shell=True, cwd=orc8_docker_path)


def configure_orc8r():
    """
    configure orc8r with a federated AGW and FEG
    """
    print('#### Configuring orc8r ####')
    subprocess.check_call(
        'fab --fabfile=dev_tools.py register_federated_vm',
        shell=True, cwd=agw_path,
    )
    subprocess.check_call(
        'fab register_feg_gw', shell=True, cwd=feg_path,
    )


def clear_gateways():
    """
    delete AGW and FEG gateways from orc8r
    """
    print('#### Removing federated agw from orc8r and deleting certs ####')
    subprocess.check_call(
        'fab --fabfile=dev_tools.py deregister_federated_agw',
        shell=True, cwd=agw_path,
    )
    print('#### Removing feg gw from orc8r and deleting certs####')
    subprocess.check_call('fab deregister_feg_gw', shell=True, cwd=feg_path)


def clear_orc8r():
    """
    delete orc8r database. Requieres orc8r to be stopped
    """
    print('#### Clearing swagger database from Orc8r ####')
    subprocess.check_call(['./run.py --clear-db'], shell=True, cwd=orc8_docker_path)
    print(
        '#### Remember you may need to delete '
        'gateway certs from the AGW and FEG ####',
    )


def build_agw(provision_vm='False'):
    """build magma on AGW on magma Vagrant VM

       provision_vm: forces the reprovision of the magma VM
    """
    print('#### Building AGW ####')
    subprocess.check_call('vagrant up magma', shell=True, cwd=agw_path)
    subprocess.check_call(
        'fab build_and_start_magma:provision_vm=%s'
        % provision_vm, shell=True, cwd=agw_path,
    )


def build_feg():
    """
    build FEG on current Host using local docker
    """
    print('#### Building FEG ####')
    subprocess.check_call(
        'docker-compose down', shell=True,
        cwd=feg_docker_integ_test_path,
    )
    subprocess.check_call(
        'docker-compose build', shell=True,
        cwd=feg_docker_integ_test_path,
    )
    subprocess.check_call(
        'docker-compose up -d', shell=True,
        cwd=feg_docker_integ_test_path,
    )


def build_test_vm(provision_vm='False'):
    print('#### Building test vm ####')
    subprocess.check_call('vagrant up magma_test', shell=True, cwd=agw_path)
    subprocess.check_call(
        'fab make_integ_tests:provision_vm=%s'
        % provision_vm, shell=True, cwd=agw_path,
    )


def build_magma_trf(provision_vm='False'):
    print('#### Building Traffic vm ####')
    subprocess.check_call('vagrant up magma_trf', shell=True, cwd=agw_path)
    subprocess.check_call(
        'fab build_and_start_magma_trf:provision_vm=%s'
        % provision_vm, shell=True, cwd=agw_path,
    )


def start_all(provision_vm='False'):
    """
    start AGW, FEG and Orc8r
    Args:
        provision_vm: forces the reprovision of the magma VM
    """
    subprocess.check_call(['./run.py'], shell=True, cwd=orc8_docker_path)
    subprocess.check_call('docker-compose up -d', shell=True, cwd=feg_docker_integ_test_path)
    subprocess.check_call(
        'fab start_magma:provision_vm=%s' % provision_vm,
        shell=True, cwd=agw_path,
    )


def stop_all():
    """
    stop AGW, FEG and Orc8r
    """
    subprocess.check_call(
        ['./run.py --down'], shell=True,
        cwd=orc8_docker_path,
    )
    subprocess.check_call(
        'docker-compose down', shell=True,
        cwd=feg_docker_integ_test_path,
    )
    subprocess.check_call(
        'vagrant halt magma', shell=True,
        cwd=agw_path,
    )
