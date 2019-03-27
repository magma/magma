"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""
import sys

from fabric.api import cd, env, get, lcd, local, puts, run
from fabric.utils import abort

sys.path.append('../../orc8r/tools')
import fab.dev_utils as dev_utils  # NOQA
import fab.pkg as pkg  # NOQA
from fab.hosts import split_hoststring # NOQA
from fab.vagrant import setup_env_vagrant  # NOQA

AWS = 'aws --region eu-west-1'
DEFAULT_CERT = "$MAGMA_ROOT/.cache/test_certs/rootCA.pem"

def register_feg_vm():
    """ Provisions the feg vm with the cloud vm """
    dev_utils.register_vm(vm_type="feg", admin_cert=(
        "../../.cache/test_certs/admin_operator.pem",
        "../../.cache/test_certs/admin_operator.key.pem"))


def _package_python(target):
    with cd("magma/orc8r/gateway/python"):
        # Copy python source code
        run("mkdir -p %s/python" % target)
        run("python3 setup.py build  --build-lib %s/python/lib --build-scripts"
            " %s/python/scripts" % (target, target))
        run("cp -pr setup.py %s/python/." % target)
        # Generate protobufs
        run("PROTO_LIST='orc8r_protos feg_protos lte_protos' make protos")
        run("cp -pr ~/build/python/gen/* %s/python/lib/." % target)
        # Create a requires.txt file
        run("python3 setup.py egg_info -e ~/build")
        run("sed -i '/dev/,$d' ~/build/orc8r.egg-info/requires.txt")
        run("cp -pr ~/build/orc8r.egg-info/requires.txt %s/python/." % target)


def _package_go(target):
    # TODO: build into different directory so we don't package tools
    with cd("magma/feg/gateway"):
        run("make build")
        run('mkdir -p %s/bin' % target)
        run("cp -Tpr /var/opt/magma/bin/ %s/bin/" % target)


def _package_scripts(target):
    with cd('magma/feg/gateway'):
        run('cp -pr scripts/install.sh %s/' % target)
        run("mkdir -p %s/ansible/roles" % target)
        run('cp -pr deploy/feg.yml %s/ansible/main.yml' % target)
        run('cp -pr deploy/roles/feg_services %s/ansible/roles/.' % target)
        run('cp -pr configs %s/config' % target)
    with cd('magma/orc8r/gateway'):
        run("mkdir -p %s/config/templates" % target)
        run('cp -pr configs/templates/* %s/config/templates/.' % target)
    with cd('magma/orc8r/tools/ansible/roles'):
        run('cp -pr pkgrepo %s/ansible/roles/.' % target)
        run('cp -pr gateway_services %s/ansible/roles/.' % target)
    with cd('magma/fb/config'):
        run('cp service/control_proxy.yml %s/config/.' % target)
        run('mkdir %s/certs' % target)
        run('cp certs/rootCA.pem %s/certs/.' % target)


def _push_archive_to_s3(vcs, target):
    pkg_name = "magma_feg_%s" % pkg.get_commit_hash(vcs)
    with cd(target):
        run('zip -r %s *' % (pkg_name))
    zip_name = "%s.zip" % pkg_name
    local("rm -rf %s" % target)
    local("mkdir -p %s" % target)
    get('%s/%s' % (target, zip_name), '%s/%s' % (target, zip_name))
    with lcd(target):
        local('%s s3 cp %s s3://magma-images/feg/' % (AWS, zip_name))
    puts("Deployment bundle: s3://magma-images/feg/%s" % zip_name)
    return zip_name


def package(feg_host=None, vcs="hg", force=False):
    """
    Create deploy package and push to S3. This defaults to running on local
    vagrant feg VM machines, but can also be pointed to an arbitrary host
    (e.g. amazon) by specifying a VM.

    feg_host: The ssh address string of the machine to run the package
        command. Formatted as "<user>@<host>:<port>". If not specified,
        defaults to the `feg` vagrant VM.

    vcs: version control system used, "hg" or "git".

    force: Bypass local commits or changes check if set to True.
    """
    if not force and pkg.check_commit_changes():
        abort("Local changes or commits not allowed")

    if feg_host:
        env.host_string = feg_host
        (env.user, _, _) = split_hoststring(feg_host)
    else:
        setup_env_vagrant("feg")

    target = "/tmp/magmadeploy_feg"
    run("rm -rf %s" % target)
    run("mkdir -p %s" % target)

    _package_python(target)
    _package_go(target)
    _package_scripts(target)
    return _push_archive_to_s3(vcs, target)


def connect_gateway_to_cloud(control_proxy_setting_path=None, cert_path=DEFAULT_CERT):
    """
    Setup the feg gateway VM to connect to the cloud
    Path to control_proxy.yml and rootCA.pem could be specified to use
    non-default control proxy setting and certificates
    """
    setup_env_vagrant("feg")
    dev_utils.connect_gateway_to_cloud(control_proxy_setting_path, cert_path)
