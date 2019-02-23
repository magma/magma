"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""
import os.path
from fabric.api import local
from fabric.api import env


def __ensure_in_vagrant_dir():
    """
    Error out if there is not Vagrant instance associated with this directory
    """
    pwd = local('pwd', capture=True)
    if not os.path.isfile(pwd + '/Vagrantfile'):
        print("Error: Vagrantfile not found. Try executing from fbcode/magma")
        exit(1)


def setup_env_vagrant(machine='magma'):
    """ Host config for local Vagrant VM.

    Sets the environment to point at the local vagrant machine. Used
    whenever we need to run commands on the vagrant machine.
    """

    __ensure_in_vagrant_dir()

    # Ensure that VM is running
    isUp = local('vagrant status %s' % machine, capture=True)\
        .find('running') < 0
    if isUp:
        # The machine isn't running. Most likely it's just not up. Let's
        # first try the simple thing of bringing it up, and if that doesn't
        # work then we ask the user to fix it.
        print("VM %s is not running... Attempting to bring it up."
              % machine)
        local('vagrant up %s' % machine)
        isUp = local('vagrant status %s' % machine, capture=True)\
            .find('running')

        if isUp < 0:
            print("Error: VM: %s is still not running...\n"
                  " Failed to bring up %s'"
                  % (machine, machine))
            exit(1)

    ssh_config = local('vagrant ssh-config %s' % machine, capture=True)
    host = local('echo "%s" | grep HostName' % ssh_config,
                 capture=True).split()[1]
    port = local('echo "%s" | grep Port' % ssh_config,
                 capture=True).split()[1]
    env.host_string = 'vagrant@%s:%s' % (host, port)
    env.hosts = [env.host_string]
    identity_file = local('echo "%s" | grep IdentityFile'
                          % ssh_config, capture=True)
    # some installations seem to have quotes around the file location
    env.key_filename = identity_file.split()[1].strip('"')


def teardown_vagrant(machine):
    """ Destroy a vagrant machine so that we get a clean environment to work
        in
    """

    __ensure_in_vagrant_dir()

    # Destroy if vm if it exists
    created = local('vagrant status %s' % machine, capture=True)\
        .find('not created') < 0

    if created:
        local('vagrant destroy -f %s' % machine)
