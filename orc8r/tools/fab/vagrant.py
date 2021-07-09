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
import os.path

from fabric.api import env, local


def __ensure_in_vagrant_dir():
    """
    Error out if there is not Vagrant instance associated with this directory
    """
    pwd = local('pwd', capture=True)
    if not os.path.isfile(pwd + '/Vagrantfile'):
        print("Error: Vagrantfile not found. Try executing from fbcode/magma")
        exit(1)

    return


def setup_env_vagrant(machine='magma', apply_to_env=True, force_provision=False):
    """ Host config for local Vagrant VM.

    Sets the environment to point at the local vagrant machine. Used
    whenever we need to run commands on the vagrant machine.
    """

    __ensure_in_vagrant_dir()

    # Ensure that VM is running
    isUp = local('vagrant status %s' % machine, capture=True) \
        .find('running') < 0
    if isUp:
        # The machine isn't running. Most likely it's just not up. Let's
        # first try the simple thing of bringing it up, and if that doesn't
        # work then we ask the user to fix it.
        print(
            "VM %s is not running... Attempting to bring it up."
            % machine,
        )
        local('vagrant up %s' % machine)
        isUp = local('vagrant status %s' % machine, capture=True) \
            .find('running')

        if isUp < 0:
            print(
                "Error: VM: %s is still not running...\n"
                " Failed to bring up %s'"
                % (machine, machine),
            )
            exit(1)
    elif force_provision:
        local('vagrant provision %s' % machine)

    ssh_config = local('vagrant ssh-config %s' % machine, capture=True)

    ssh_lines = [line.strip() for line in ssh_config.split("\n")]
    ssh_params = {
        key: val for key, val
        in [line.split(" ", 1) for line in ssh_lines]
    }

    host = ssh_params.get("HostName", "").strip()
    port = ssh_params.get("Port", "").strip()
    # some installations seem to have quotes around the file location
    identity_file = ssh_params.get("IdentityFile", "").strip().strip('"')
    host_string = 'vagrant@%s:%s' % (host, port)

    if apply_to_env:
        env.host_string = host_string
        env.hosts = [env.host_string]
        env.key_filename = identity_file
        env.disable_known_hosts = True
    else:
        return {
            "hosts": [host_string],
            "host_string": host_string,
            "key_filename": identity_file,
            "disable_known_hosts": True,
        }


def teardown_vagrant(machine):
    """ Destroy a vagrant machine so that we get a clean environment to work
        in
    """

    __ensure_in_vagrant_dir()

    # Destroy if vm if it exists
    created = local('vagrant status %s' % machine, capture=True) \
        .find('not created') < 0

    if created:
        local('vagrant destroy -f %s' % machine)
