"""
Devops tasks: builds, deployments and migrations.

Copyright (c) 2017-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import sys
from distutils.util import strtobool

from fabric.api import cd, env, lcd, local, run
from fabric.operations import get
from fabric.utils import abort, puts

sys.path.append('../../orc8r/tools')
import fab.pkg as pkg  # NOQA
from fab.hosts import vagrant_setup, ansible_setup, split_hoststring # NOQA


# List of service tiers
SERVICES = ["controller", "proxy", "osquery", "prometheus"]

# Look for keys as specified in our ~/.ssh/config
env.use_ssh_config = True


def build():
    """ [deploy] Build cloud binaries in VM """
    # TODO: build into a different directory so we don't package tools like
    # swagger and protoc
    # run("rm -rf cloud/go/bin/")
    run("make -C magma/orc8r/cloud build")


def package(service, cloud_host=None, vcs="hg", force=False):
    """
    Create deploy package and push to S3. This defaults to running on local
    vagrant cloud VM machines, but can also be pointed to an arbitrary host
    (e.g. amazon) by specifying a VM.

    cloud_host: The ssh address string of the machine to run the package
        command. Formatted as "<user>@<host>:<port>". If not specified,
        defaults to the `cloud` vagrant VM.

    vcs: version control system used, "hg" or "git".

    force: Bypass local commits or changes check if set to True.
    """
    # Check that we have no local changes or commits at this point
    if not force and pkg.check_commit_changes():
        abort("Local changes or commits not allowed")

    _validate_service(service)

    if cloud_host:
        env.host_string = cloud_host
        (env.user, _, _) = split_hoststring(cloud_host)
    else:
        _vagrant()

    # Use same temp folder name for local and VM operations
    folder = "/tmp/magmadeploy_%s" % service
    local("rm -rf %s" % folder)
    local("mkdir -p %s" % folder)
    run("rm -rf %s" % folder)
    run("mkdir -p %s" % folder)

    with cd('magma/orc8r/cloud/deploy'):
        run('cp -pr aws/%s_appspec.yml %s/appspec.yml' % (service, folder))
        run('cp -pr aws/scripts %s/.' % folder)
        run("mkdir -p %s/ansible/roles" % folder)
        run('cp -pr %s.yml %s/ansible/main.yml' % (service, folder))
        run('cp -pr roles/%s %s/ansible/roles/.' % (service, folder))
        run('cp -pr roles/aws_setup %s/ansible/roles/.' % folder)
        run('cp -pr roles/osquery %s/ansible/roles/.' % folder)

        if service == "controller":
            run('cp -pr /etc/magma %s/configs' % folder)
            run('cp -pr files/scripts/setup_swagger_ui %s/scripts/.' % folder)
            run('cp -pr files/static/apidocs %s/.' % folder)
        if service == "proxy":
            run('cp -pr /etc/magma %s/configs' % folder)
            run('cp -pr roles/disk_metrics %s/ansible/roles/.' % folder)
            run('cp -pr ../../../orc8r/tools/ansible/roles/pkgrepo '
                '%s/ansible/roles/.' % folder)
        if service == 'prometheus':
            run('cp -pr roles/prometheus %s/ansible/roles/.' % folder)
            run('mkdir -p %s/bin' % folder)  # To make CodeDeploy happy

    # Build Go binaries and plugins
    build()
    run('cp -pr go/plugins %s' % folder)
    _copy_go_binaries(service, folder)

    # Zip and push to s3
    pkg_name = "magma_%s_%s" % (service, pkg.get_commit_hash(vcs))
    _push_archive_to_s3(service, folder, pkg_name)
    run('rm -rf %s' % folder)
    local('rm -rf %s' % folder)
    return "%s.zip" % pkg_name


def _copy_go_binaries(service, folder):
    if service == 'proxy':
        run('mkdir -p %s/bin' % folder)
        run('cp go/bin/metricsd %s/bin/.' % folder)
        run('cp go/bin/logger %s/bin/.' % folder)
    if service == 'controller':
        run('cp -pr go/bin %s' % folder)


def _push_archive_to_s3(service, folder, pkg_name):
    with cd(folder):
        run('zip -r %s *' % (pkg_name))
    get('%s/%s.zip' % (folder, pkg_name), '%s/%s.zip' % (folder, pkg_name))
    with lcd(folder):
        local('aws s3 cp %s.zip s3://magma-images/cloud/' % pkg_name)
    puts("Deployment bundle: s3://magma-images/cloud/%s.zip" % pkg_name)
    puts("To deploy, use 'fab staging deploy:%s,%s.zip'"
         % (service, pkg_name))


def _validate_service(service):
    if service not in SERVICES:
        raise ValueError(
            "Invalid service '%s'. Valid service tiers are: %s" %
            (service, SERVICES)
        )


def _vagrant():
    """ Host config for local Vagrant VM. """
    machine = "cloud"
    host = local(
        'vagrant ssh-config %s | grep HostName' % (machine), capture=True
    ).split()[1]
    port = local(
        'vagrant ssh-config %s | grep Port' % (machine), capture=True
    ).split()[1]
    env.host_string = 'vagrant@%s:%s' % (host, port)
    identity_file = local(
        'vagrant ssh-config %s | grep IdentityFile' % (machine), capture=True
    )
    # add Vagrant identity file to any values passed on command line
    if env.key_filename is None:
        env.key_filename = []
    # some installations seem to have quotes around the file location
    env.key_filename.append(identity_file.split()[1].strip('"'))


def cloud_test(cloud=None, datastore=None, destroy_vm="True"):
    """
    Run the cloud tests. This defaults to running on local vagrant
    machines, but can also be pointed to an arbitrary host (e.g. amazon) by
    passing "address:port" as arguments

    cloud: The ssh address string of the machine to run the cloud
        on. Formatted as "host:port". If not specified, defaults to
        the `cloud` vagrant box.

    datastore: The ssh address string of the machine to run the datastore on
        on. Formatted as "host:port". If not specified, defaults to the
        `datastore` vagrant box.
    """
    destroy_vm = bool(strtobool(destroy_vm))

    # Setup the datastore: use the provided test machine if given, else default
    # to the vagrant machine
    if not datastore:
        datastore = vagrant_setup("datastore", destroy_vm)
    else:
        ansible_setup(datastore, "datastore", "datastore.dev.yml")

    # Setup the cloud: use the provided address if given, else default to the
    # vagrant machine
    if not cloud:
        cloud = vagrant_setup("cloud", destroy_vm)
    else:
        ansible_setup(cloud, "magma-cloud-dev", "cloud.dev.yml")
        env.host_string = cloud
        (env.user, _, _) = split_hoststring(cloud)

    with cd('~/magma/orc8r/cloud'):
        run('make cover')
