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
import os.path
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
    return


def setup_env_vagrant(
    c: Connection, machine: str = 'magma', force_provision: bool = False,
) -> Tuple[Connection, Dict[str, str]]:
    """ Host config for local Vagrant VM.

    Sets the environment to point at the local vagrant machine. Used
    whenever we need to run commands on the vagrant machine.
    """
    _ensure_in_vagrant_dir(c)

    # An empty local ssh config file stops warnings from being printed.
    c.run('touch ~/.ssh/config')

    # Ensure that VM is running
    down = c.run(f'vagrant status {machine}', hide=True).stdout \
        .find('running') < 0
    if down:
        # The machine isn't running. Most likely it's just not up. Let's
        # first try the simple thing of bringing it up, and if that doesn't
        # work then we ask the user to fix it.
        print(f"VM {machine} is not running... Attempting to bring it up.")
        c.run(f'vagrant up {machine}')
        down = c.run(f'vagrant status {machine}', hide=True).stdout \
            .find('running')

        if down < 0:
            print(
                f"Error: VM: {machine} is still not running...\n"
                f" Failed to bring up {machine}",
            )
            exit(1)
    elif force_provision:
        c.run(f'vagrant provision {machine}')

    ssh_config = c.run(f'vagrant ssh-config {machine}', hide=True).stdout.strip()

    ssh_lines = [line.strip() for line in ssh_config.split("\n")]
    ssh_params = {
        key: val for key, val
        in [line.split(" ", 1) for line in ssh_lines]
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


def teardown_vagrant(c: Connection, machine: str) -> None:
    """ Destroy a vagrant machine so that we get a clean environment to work
        in
    """
    _ensure_in_vagrant_dir(c)

    # Destroy vm if it exists
    created = c.run(f'vagrant status {machine}', hide=True).stdout.find('not created') < 0

    if created:
        c.run(f'vagrant destroy -f {machine}')
