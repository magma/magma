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

from fabric.api import cd, hide, local, run, settings
from fabric.operations import get

# Local changes are only allowed in files specified in the EXCLUDE_FILE_LIST
EXCLUDE_FILE_LIST = [
    os.path.realpath(x) for x in [
        'release/build-magma.sh',
    ]
]


def check_commit_changes():
    """ Compare agaist remote/master, ensure there is no local modifications or
    commits """
    with hide('running', 'warnings', 'output'), settings(warn_only=True):
        uncommitted_change = local(
            'hg summary | grep -q "commit: (clean)";'
            'echo $?', capture=True,
        ) == "1"
        local_commit_hash = local(
            'hg identify -i',
            capture=True,
        ).replace('+', '').split()[0]
        remote_commit_hash = local(
            'hg identify -r remote/master',
            capture=True,
        ).split()[0]
        if uncommitted_change:
            changes = [
                os.path.realpath(x.split()[1])
                for x in local('hg status', capture=True).split('\n')
            ]
            if set(changes) <= set(EXCLUDE_FILE_LIST):
                print(
                    "Local changes detected in the following files from the "
                    "EXCLUDE_FILE_LIST are ignored:\n%s" % '\n'.join(changes),
                )
                return False
            else:
                print("Warning: uncommitted changes found!!!")
                return True
        elif local_commit_hash != remote_commit_hash:
            print("Warning: local commits found compared against remote/master\
                   !!!")
            return True
        else:
            return False


def get_commit_hash(vcs='hg'):
    with hide('running', 'warnings', 'output'), settings(warn_only=True):
        if vcs == 'hg':
            local_commit_hash = local(
                'hg identify -i',
                capture=True,
            ).replace('+', '').split()[0]
        elif vcs == 'git':
            local_commit_hash = local('git rev-parse HEAD', capture=True)
        else:
            print('Unknown vcs: %s' % vcs)
            exit(1)
    return local_commit_hash[0:8]


def download_all_pkgs():
    '''
    Figure out the list of installed packages on the system and download the
    packages locally into the apt cache
    '''
    # Get a list of all the installed packages
    packages = run(
        'dpkg-query -l'
        ' | tail -n +6'
        ' | awk \'{print $2}\''
        ' | sed "s/:.*$//"',
    )
    # Get the list of the packages we have local .debs for
    have_pkgs = run(
        'ls /var/cache/apt/archives/*.deb'
        ' | xargs -I "%" dpkg -I "%"'
        ' | grep " Package: "'
        ' | awk \'{print $2}\'',
    )
    # Figure out the set difference of the two
    have_not_pkgs = run(
        'echo \'' + packages + "\n" + have_pkgs + '\''
        ' | sort'
        ' | uniq -u'
        ' | tr "\\n\\r" " "',
    )
    # Download all the packages we don't have
    with cd('/var/cache/apt/archives'):
        for p in have_not_pkgs.split():
            # Some of the packages aren't available on the repo, so
            # download them one at a time so we can swallow errors
            run('sudo aptitude download -q2 %s || true' % p)


def upload_pkgs_to_aws():
    """
    Upload the dependencies in the apt cache to aws. This allows us to record
    and retrieve the versions of the dependencies a specific version of magma
    was tested against.

    This creates three files:
       - VERSION.deps.tar.gz -- A tar ball of all installed packages on the
                                machine
       - VERSION.deplist     -- Text metadata of the packages in the tar
                                ball, formatted as:

                                package_name package_version package_file_name

       - VERSION.lockfile    -- A python lock file listing the installed python
                                dependencies
    """

    # Get the version of magma we are releasing
    magma_version = get_magma_version()
    copy_packages()

    # Upload to AWS
    s3_path = 's3://magma-images/gateway/' + magma_version
    local('aws s3 cp /tmp/packages.txt ' + s3_path + '.deplist')
    local('aws s3 cp release/magma.lockfile.debian ' + s3_path + '.lockfile.debian')
    local('aws s3 cp /tmp/packages.tar.gz ' + s3_path + '.deps.tar.gz')

    # Clean up
    run('rm -r /tmp/packages')
    run('rm /tmp/packages.tar.gz')
    local('rm /tmp/packages.tar.gz')
    local('rm /tmp/packages.txt')


def get_magma_version():
    return run(
        'ls ~/magma-packages'
        ' | grep "^magma_[0-9].*"'
        ' | xargs -I "%" dpkg -I ~/magma-packages/%'
        ' | grep "Version"'
        ' | awk \'{print $2}\'',
    )


def copy_packages():
    """
    Copy the dependencies in the apt cache to /tmp on the local machine.
    """
    # Build a list of package metadata were each row looks like:
    #
    #   pkg_name version deb_name
    run('rm -rf /tmp/packages')
    run('rm -rf /tmp/packages.txt')
    run('mkdir -p /tmp/packages')
    run('cp ~/magma-packages/*.deb /tmp/packages')
    packages = run('ls /tmp/packages')

    with cd('/tmp/packages'):
        for f in packages.split():
            pkg_name = run(
                'dpkg -I %s'
                ' | grep " Package:"'
                ' | awk \'{print $2}\'' % f,
            ).strip()
            version = run(
                'dpkg -I %s'
                ' | grep " Version:"'
                ' | awk \'{print $2}\'' % f,
            ).strip()
            run(
                'echo "%s %s %s" >> /tmp/packages.txt'
                % (pkg_name, version, f),
            )

    # Tar up the packages
    with cd('/tmp/packages'):
        run('tar czf /tmp/packages.tar.gz *.deb')

    # Pull the artifacts onto the local machine
    get('/tmp/packages.tar.gz', '/tmp/packages.tar.gz')
    get('/tmp/packages.txt', '/tmp/packages.txt')
    magma_version = get_magma_version()
    local(f'echo "{magma_version}" > /tmp/magma_version')
