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
import os
from typing import Dict, Tuple

from fabric import Connection


def _ensure_in_vagrant_dir(c: Connection) -> None:
    """
    Error out if there is not Vagrant instance associated with this directory
    """
    pwd = c.run('pwd', hide=True).stdout.strip()
    if not os.path.isfile(pwd + '/Vagrantfile'):
        # check if we are on a vagrant subdirectory
        ret_code = c.run('vagrant validate', hide=True, warn=True).return_code
        if ret_code != 0:
            print("Error: Vagrantfile not found or not valid")
            exit(ret_code)


def setup_env_vagrant(
    c: Connection, machine: str = 'magma', force_provision: bool = False,
    max_retries: int = 1,
) -> Tuple[Connection, Dict[str, str]]:
    """ Host config for local Vagrant VM.

    Sets the environment to point at the local vagrant machine. Used
    whenever we need to run commands on the vagrant machine.

    Parameters:
        c: The connection context for execution
        machine: The machine that should be setup
        force_provision: Whether provisioning should be done
        max_retries: Max attempts to setup machine

    Returns:
        Connection to the set up vm and meta information
    """
    _ensure_in_vagrant_dir(c)

    # An empty local ssh config file stops warnings from being printed.
    c.run('touch ~/.ssh/config')

    # Ensure that VM is not already running
    if not _is_vm_running(c, machine):
        _vagrant_up_with_retry(c, machine, max_retries)
    elif force_provision:
        c.run(f'vagrant provision {machine}')

    ssh_config = c.run(f'vagrant ssh-config {machine}', hide=True).stdout.strip()

    ssh_lines = [line.strip() for line in ssh_config.split("\n")]
    ssh_key_val = [line.split(" ", 1) for line in ssh_lines]
    ssh_params = {
        key: val
        for key, val in ssh_key_val
    }

    host = ssh_params.get("HostName", "").strip()
    port = ssh_params.get("Port", "").strip()
    # some installations seem to have quotes around the file location
    identity_file = ssh_params.get("IdentityFile", "").strip().strip('"')
    host_string = f'vagrant@{host}:{port}'

    return Connection(
        host_string, connect_kwargs={"key_filename": identity_file},
    ), {
        "host_string": host_string,
        "key_filename": identity_file,
    }


def _is_vm_running(c: Connection, machine: str):
    return c.run(f'vagrant status {machine}', hide=True).stdout.find('running') >= 0


def _vagrant_up_with_retry(c: Connection, machine: str, max_retries: int):
    print(f'VM {machine} is not running - attempting to bring it up ...')
    for attempt in range(max_retries):
        inc_attempt = attempt + 1
        print(f'... attempt {inc_attempt}/{max_retries} to bring up VM {machine}')

        try:
            c.run(f'vagrant up {machine}')
        except Exception:
            print(f'... VM {machine} failed during "vagrant up" - VM will be destroyed')
            teardown_vagrant(c, machine)

        if _is_vm_running(c, machine):
            print(f'... VM {machine} is running.')
            break  # success

    if not _is_vm_running(c, machine):
        print(f'Error: VM {machine} is still not running after {max_retries} attempt(s).')
        exit(1)


def teardown_vagrant(c: Connection, machine: str) -> None:
    """
    Destroy a vagrant machine so that we get a clean environment to work in.

    Parameters:
        c: The connection context for execution
        machine: The machine to be destroyed
    """
    _ensure_in_vagrant_dir(c)

    # Destroy vm if it exists
    created = c.run(f'vagrant status {machine}', hide=True).stdout.find('not created') < 0

    if created:
        c.run(f'vagrant destroy -f {machine}')
