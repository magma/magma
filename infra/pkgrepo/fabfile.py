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

import fnmatch
import getpass
import json
import os
import re
import sys
import time
from contextlib import contextmanager
from distutils.version import LooseVersion
from enum import Enum

import requests
from fabric.api import cd, env, execute, hide, hosts, lcd, local, run, settings
from fabric.operations import get, prompt, put
from github import GithubException

if not env.get('magma_root'):
    env.magma_root = os.environ.get('MAGMA_ROOT', '../..')
env.magma_root = os.path.abspath(os.path.expanduser(env.magma_root))

sys.path.append(env.magma_root + '/orc8r')
from tools.fab.vagrant import setup_env_vagrant

env.pkgfmt = "deb"

# Look for keys as specified in our ~/.ssh/config
env.use_ssh_config = True


class DeployTarget(Enum):
    STABLE = 1
    BETA = 2
    DEV = 3
    TEST = 4


def _setup_env(target, channel):
    # If a host list isn't specified, default to the magma vagrant vm
    if not env.hosts:
        with lcd(env.magma_root + '/lte/gateway'):
            setup_env_vagrant()

    env.deploy_target = target
    env._release_channel = channel


def dev():
    """ [deploy] dev deploy settings """
    _setup_env(DeployTarget.DEV, "stretch-dev")


def beta():
    """ [deploy] beta deploy settings """
    _setup_env(DeployTarget.BETA, "stretch-beta")


def test():
    """ [deploy] Test deploy settings """
    _setup_env(DeployTarget.TEST, "stretch-test")


def stable():
    _setup_env(DeployTarget.STABLE, "stretch-stable")


def shipit():
    """Takes packages on a local dev VM and pushes them to repo."""

    if env.pkgfmt != "deb":
        # Since we don't support pushing packages to non-deb repos, just fail
        # early. This can be removed when _push_packages_to_repo has CentOS
        # support.
        print("Only shipping deb packages is supported, not shipping.")
        return

    with hide('running', 'warnings', 'output'), settings(warn_only=True):
        execute(_get_packages)

        pkgs = local("ls -lrth /tmp/magma-packages-deploy/*.deb", capture=True)
        magma_pkgs = local("ls /tmp/magma-packages-deploy | grep '^magma-[0-9].*.deb'",
                           capture=True).split()
        if len(pkgs.strip()) == 0:
            print("No packages to release!")
            execute(cleanup_package_deploy)
            exit(0)

    with hide('running', 'warnings', 'output'), settings(warn_only=True):
        release = env._release_channel
        if release not in ['stretch-test', 'stretch-dev', 'stretch-beta',
                           'stretch-stable']:
            release = prompt('Specify release branch',
                             validate='^(stretch-test|stretch-dev|stretch-beta'
                             + '|stretch-stable)$')
        print("Releasing to '%s'." % release)

        env.release_success = False

    with apply_custom_env_options() as opts:
        execute(push_packages_to_repo, release, **opts)
    execute(cleanup_package_deploy)


def _get_packages():
    """Get all the packages from a remote machine and put them in a local temp
    deploy directory.

    We check if the package we're getting has a version less than or equal to
    what's already in the remote VM's apt cache; if so, we print a warning and
    we don't pull the package in.
    """
    local('mkdir -p /tmp/magma-packages-deploy')
    result = run('ls ~/magma-packages/*.%s' % env.pkgfmt)
    if result.return_code != 0:
        local('ls /tmp/magma-packages-deploy')
        return

    pkgs = result.stdout.split()

    for p in pkgs:
        with hide('running', 'warnings', 'output'), settings(warn_only=True):
            put(local_path='scripts/pkgavail',
                remote_path='/tmp/')
            run('chmod a+x /tmp/pkgavail')
            is_avail = run("/tmp/pkgavail --release %s %s" %
                           (env._release_channel, p)).strip()
            if is_avail == "False":
                get(remote_path='%s' % p,
                    local_path='/tmp/magma-packages-deploy/')
            else:
                print(("WARNING: %s is not a newer version than what is "
                       "available already in release branch '%s', ignoring." %
                       (p, env._release_channel)))
    local('ls /tmp/magma-packages-deploy')


@contextmanager
def apply_custom_env_options(**kwargs):
    extra_args = {}
    extra_settings = {}
    extra_settings.update(kwargs)
    pkgrepo_env = setup_env_vagrant(machine=env.get('local_pkgrepo')
                                    or extra_settings.get('local_pkgrepo', 'aptly'),
                                    apply_to_env=False)
    extra_args["hosts"] = pkgrepo_env["hosts"]
    extra_settings["hosts"] = pkgrepo_env["hosts"]
    extra_settings["key_filename"] = pkgrepo_env["key_filename"]
    extra_settings["host_string"] = pkgrepo_env["host_string"]
    
    with settings(**extra_settings):
        yield extra_args


def promote(src, dest, version):
    """
    Copy a magma package and its dependencies from one channel to another

    src: The channel to promote from
    dest: The channel to promote to
    version: The version of magma to promote. It should look something like
             0.3.31-1508456917-bdbaa7c2
    """

    REPO_PREFIX = 'stretch'
    if (src == 'test' and dest != 'beta') \
       or (src == 'beta' and dest != 'stable') \
       or (src != 'test' and src != 'beta'):
        print("Supported promotions are:\n\n"
              "\tfab promote:test,beta,VERSION\n"
              "\tfab promote:beta,stable,VERSION\n")
        return

    # FIXME: don't use aws here                                       

    # Grab the list of dependencies from the aws bucket
    os.environ['AWS_PROFILE'] = 'fbinfra'
    local('aws s3 cp s3://magma-images/gateway/%s.deplist'
          ' /tmp/deplist' % version)
    pkgs = local('cat /tmp/deplist | awk \'{print $3}\'', capture=True)
    pkgs = pkgs.split()

    print("Promoting magma version `%s` from '%s' to '%s'"
          % (version, src, dest))

    env.release_success = False
    try:
        repo_src = "%s-%s" % (REPO_PREFIX, src)
        repo_dest = "%s-%s" % (REPO_PREFIX, dest)
        with apply_custom_env_options():
            promote_pkgs(repo_src, repo_dest, pkgs)
        env.release_success = True
    finally:
        pass
    return


def fab_cmd(prefix):
    return lambda cmd, **kwargs: run(prefix + ' ' + cmd, **kwargs)


run_aptly = fab_cmd('docker-compose exec aptly aptly ')
rsudo = fab_cmd('sudo')


def push_packages_to_repo(release_channel):
    """Push local deploy directory of packages to actual repo, and refresh the
    repo.
    """

    if env.pkgfmt != "deb":
        # We only support freight, which is only for deb packages. We'd need to
        # add something that understands RPM repos as well if we want to add
        # support for CentOS here.
        print("Only pushing deb packages is supported, not pushing.")
        return

    magma_filename = local("ls /tmp/magma-packages-deploy/magma_*.deb",
                           capture=True).split("/")[-1]

    # quick check that version is in the right format
    # -- [version]-[timestamp]-[hashprefix]
    version, ts, hashprefix = magma_filename[6:-4].split("-")

    build_id = "-".join([version, ts, hashprefix])

    run("rm -rf /tmp/magma-packages-deploy")
    run('mkdir -p /tmp/magma-packages-deploy')
    put(local_path='/tmp/magma-packages-deploy/*.deb',
        remote_path='/tmp/magma-packages-deploy/')

    with cd("~/docker"):
        run_aptly("repo create " + release_channel, quiet=True)
        run("docker-compose exec aptly mkdir -p upload/" + build_id,
            quiet=True)

        aptly_image_ps = run("docker-compose ps -q aptly")

        destdir = "/home/aptly-user/upload/" + build_id
        run("docker cp /tmp/magma-packages-deploy/ "
            + aptly_image_ps + ":" + destdir)

        run("docker-compose exec -u root aptly chown -R aptly-user:aptly-user "
            + destdir)

        run_aptly("repo add -remove-files " + release_channel + " "
                  + destdir + "/magma-packages-deploy/" + magma_filename)
        run_aptly("repo add -remove-files " + release_channel + " "
                  + destdir + "/magma-packages-deploy/", warn_only=True)

        run("docker-compose exec aptly rm -r -- '" + destdir + "'")
        run_aptly("publish -architectures=amd64 -distribution "
                  + release_channel + " repo " + release_channel,
                  warn_only=True)
        run_aptly("publish update " + release_channel)

    env.release_success = True


def cleanup_package_deploy():
    """Delete local temp deploy directory."""
    local('rm -r /tmp/magma-packages-deploy')


def promote_pkgs(src, dest, packages):
    """
    Copy a list of packages from one channel to another
    """
    packages = [package[:-4] if package.endswith(".deb") else package
                for package in packages]
    # magma package filename not consistent with .deb version format
    # convert package name to package query for current version
    # in similar form to 'magma (= 0.3.74-1560475061-fb43abf4)'
    # to conform to the expected filename format,
    # magma-0.3.74-1560475061-fb43abf4.deb
    # should be
    # magma_0.3.74-1560475061-fb43abf4_amd64.deb
    packages = [re.sub("^magma-([0-9]+.*)", r"'magma (= \1)'", package)
                for package in packages]

    with cd("docker"):
        pkg_list = " ".join(packages)
        run_aptly("repo create {dest}".format(dest=dest), warn_only=True)
        run_aptly("repo copy -with-deps -architectures=amd64 "
                  "{src} {dest} {pkg_list}".format(src=src, dest=dest,
                                                   pkg_list=pkg_list))
        run_aptly("publish repo -architectures=amd64 "
                  "-distribution={dest} {dest}".format(dest=dest),
                  warn_only=True)
        run_aptly("publish update " + dest)
    return


def _get_latest_version(version_list):
    loose_versions = [LooseVersion(x) for x in version_list]
    return str(max(loose_versions))


def as_bool(v):
    return str(v).lower() in ['true', 't', '1', 'y', 'yes']


def get_repo_list(release):
    upstream_repos = {
        'bionic': [
            'deb http://security.ubuntu.com/ubuntu bionic-security main restricted',
            'deb http://security.ubuntu.com/ubuntu bionic-security multiverse',
            'deb http://security.ubuntu.com/ubuntu bionic-security universe',
            ('deb http://us.archive.ubuntu.com/ubuntu bionic-backports main '
             'restricted universe multiverse'),
            'deb http://us.archive.ubuntu.com/ubuntu bionic main restricted',
            'deb http://us.archive.ubuntu.com/ubuntu bionic multiverse',
            'deb http://us.archive.ubuntu.com/ubuntu bionic universe',
            'deb http://us.archive.ubuntu.com/ubuntu bionic-updates main restricted/',
            'deb http://us.archive.ubuntu.com/ubuntu bionic-updates multiverse',
            'deb http://us.archive.ubuntu.com/ubuntu bionic-updates universe',
        ],
        'xenial': [
            'deb http://security.ubuntu.com/ubuntu xenial-security main restricted',
            'deb http://security.ubuntu.com/ubuntu xenial-security multiverse',
            'deb http://security.ubuntu.com/ubuntu xenial-security universe',
            ('deb http://us.archive.ubuntu.com/ubuntu/ xenial-backports main '
             'restricted universe multiverse'),
            'deb http://us.archive.ubuntu.com/ubuntu/ xenial main restricted',
            'deb http://us.archive.ubuntu.com/ubuntu/ xenial multiverse',
            'deb http://us.archive.ubuntu.com/ubuntu/ xenial universe',
            'deb http://us.archive.ubuntu.com/ubuntu/ xenial-updates main restricted',
            'deb http://us.archive.ubuntu.com/ubuntu/ xenial-updates multiverse',
            'deb http://us.archive.ubuntu.com/ubuntu/ xenial-updates universe',
        ]
    }
    return upstream_repos.get(release)


def get_required_packages(release):
    per_release = {}
    required_packages = [
        'autoconf',
        'automake',
        'build-essential',
        'bzip2',
        'bzr',
        'curl',
        'daemontools',
        'debhelper',
        'fakeroot',
        'git',
        'graphviz',
        'libseccomp2',
        'libssl-dev',
        'libtool',
        'netcat',
        'openssl',
        'python-all',
        'python-cffi-backend',
        'python-six',
        'python-twisted-conch',
        'python-zope.interface',
        'python3-pip',
        'supervisor',
        'unzip',
        'vim',
        'wget',
    ]

    return required_packages + per_release.get(release, [])


def make_ubuntu_snapshot(release, force_provision=False, clean=False,
                         extra_packages=''):
    extra_packages = extra_packages.split(',') if extra_packages else []
    clean = as_bool(clean)
    force_provision = as_bool(force_provision)

    setup_env_vagrant(machine='aptly', force_provision=force_provision)
    workdir = '/work/{}'.format(release)
    cachedir = '/cache/{}/deb'.format(release)
    localdir = os.path.expanduser('~/.magma/ubuntu_snapshots/{}'.format(release))
    snap_time = time.time()
    snap_name = time.strftime('{}_%Y%m%d%H%M%S'.format(release),
                              time.gmtime(snap_time))

    rsudo('mkdir -p {}'.format(workdir))
    rsudo('apt update')
    rsudo('apt install -y debootstrap')

    if clean:
        rsudo('find {} -name "*.deb" -delete'.format(cachedir))

    rsudo('rm -rf {}'.format(workdir))
    rsudo('mkdir -p {}'.format(workdir))

    _prepare_ubuntu_snapshot(release, cachedir, workdir, extra_packages)
    _create_ubuntu_archive(release, snap_name, workdir, localdir)


def _prepare_ubuntu_snapshot(release, cachedir, workdir, extra_packages):
    required_packages = get_required_packages(release)
    with cd(workdir):
        # sync from cache
        rsudo('mkdir -p {}'.format(cachedir))
        rsudo('mkdir -p {}/var/cache/apt/archives'.format(workdir))
        rsudo('rsync -r {}/ {}/var/cache/apt/archives/'.format(cachedir, workdir))

        rsudo(('debootstrap --arch amd64 --components=main,universe --download-only '
               '--include=software-properties-common {} '
               '{} http://archive.ubuntu.com/ubuntu/'
               '').format(release, workdir))

        # sync to cache
        rsudo('rsync -r {}/var/cache/apt/archives/ {}/'.format(workdir, cachedir))

        # install base system
        rsudo(('debootstrap --arch amd64 --components=main,universe '
               '--include=software-properties-common {} '
               '{} http://archive.ubuntu.com/ubuntu/'
               '').format(release, workdir))
        for repo in get_repo_list(release):
            rsudo('''chroot {} /bin/bash -c "add-apt-repository '{}'"'''.format(workdir,
                                                                                repo))

        rsudo('chroot /work/{} /bin/bash -c "apt update"'.format(release))

        rsudo(('chroot /work/{} /bin/bash -c "apt install -y --download-only {}"'
               '').format(release, ' '.join(extra_packages + required_packages)))

        # sync to cache
        rsudo('rsync -r {}/var/cache/apt/archives/ {}/'.format(workdir, cachedir))


def _create_ubuntu_archive(release, snap_name, workdir, localdir):
    with cd('~/docker'):
        archivedir = '/tmp/' + snap_name
        archivefile = snap_name + '.tar.gz'
        uploaddir = archivedir + '/upload'
        conf = archivedir + '/aptly.conf'
        run_compose = fab_cmd('docker-compose exec aptly')
        # not actually sudo, but has equivalent effect
        rsudo_compose = fab_cmd('docker-compose exec -u root aptly')

        try:
            aptly_image_ps = run('docker-compose ps -q aptly')

            run_compose('mkdir -p ' + archivedir, quiet=True)
            run_compose('mkdir -p ' + uploaddir)

            run('docker cp ~/template_aptly.conf '
                + aptly_image_ps + ':' + conf)
            rsudo_compose('chown aptly-user:aptly-user ' + conf)
            run_compose('sed -i "s|ROOT|{}|g" {}'.format(archivedir, conf))
            run('docker cp {}/var/cache/apt/archives/. {}:{}'.format(workdir,
                                                                     aptly_image_ps,
                                                                     uploaddir))

            rsudo_compose('chown -R aptly-user:aptly-user ' + uploaddir)

            run_aptly = fab_cmd('docker-compose exec aptly aptly -config=' + conf)
            run_aptly('repo create ' + snap_name)
            run_aptly('repo add {} {}'.format(snap_name, uploaddir))
            run_aptly('publish -distribution {} repo {}'.format(release, snap_name))

            run_compose('cp /aptly/public/key.gpg {}/public/'.format(archivedir))
            run_compose('mv {}/public {}/{}'.format(archivedir, archivedir, snap_name))
            run_compose('tar czf /tmp/{} -C {} {}'.format(archivefile, archivedir,
                                                          snap_name))
            run('docker cp {}:/tmp/{} /tmp/'.format(aptly_image_ps, archivefile))
            get('/tmp/{}'.format(archivefile), '{}/%(path)s'.format(localdir))
        finally:
            rsudo_compose('rm -rf -- ' + archivedir)
            rsudo('rm -f /tmp/{}'.format(archivefile))
