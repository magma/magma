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

from fabric import Connection

# Local changes are only allowed in files specified in the EXCLUDE_FILE_LIST
EXCLUDE_FILE_LIST = [
    os.path.realpath(x) for x in [
        'release/build-magma.sh',
    ]
]


def check_commit_changes(c: Connection):
    """ Compare against remote/master, ensure there is no local modifications or
    commits """
    uncommitted_change = c.run(
        'hg summary | grep -q "commit: (clean)";'
        'echo $?', hide=True, warn=True,
    ).stdout == "1"
    local_commit_hash = c.run(
        'hg identify -i',
        hide=True, warn=True,
    ).stdout.replace('+', '').split()[0]
    remote_commit_hash = c.run(
        'hg identify -r remote/master',
        hide=True, warn=True,
    ).stdout.split()[0]
    if uncommitted_change:
        changes = [
            os.path.realpath(x.split()[1])
            for x in c.run('hg status', capture=True).split('\n')
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


def get_commit_hash(c: Connection):
    return c.run('git rev-parse HEAD').stdout[0:8]


def get_commit_count(c: Connection):
    return c.run('git rev-list --count HEAD').stdout.strip()


def download_all_pkgs(c_gw: Connection):
    '''
    Figure out the list of installed packages on the system and download the
    packages locally into the apt cache
    '''
    # Get a list of all the installed packages
    packages = c_gw.run(
        'dpkg-query -l'
        ' | tail -n +6'
        ' | awk \'{print $2}\''
        ' | sed "s/:.*$//"',
    )
    # Get the list of the packages we have local .debs for
    have_pkgs = c_gw.run(
        'ls /var/cache/apt/archives/*.deb'
        ' | xargs -I "%" dpkg -I "%"'
        ' | grep " Package: "'
        ' | awk \'{print $2}\'',
    )
    # Figure out the set difference of the two
    have_not_pkgs = c_gw.run(
        'echo \'' + packages + "\n" + have_pkgs + '\''
        ' | sort'
        ' | uniq -u'
        ' | tr "\\n\\r" " "',
    )
    # Download all the packages we don't have
    with c_gw.cd('/var/cache/apt/archives'):
        for p in have_not_pkgs.split():
            # Some of the packages aren't available on the repo, so
            # download them one at a time so we can swallow errors
            c_gw.run(f'sudo aptitude download -q2 {p} || true')


def upload_pkgs_to_aws(c_gw: Connection):
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
    magma_version = get_magma_version(c_gw)
    copy_packages(c_gw)

    # Upload to AWS
    s3_path = 's3://magma-images/gateway/' + magma_version
    c_gw.run('aws s3 cp /tmp/packages.txt ' + s3_path + '.deplist')
    c_gw.run('aws s3 cp release/magma.lockfile.debian ' + s3_path + '.lockfile.debian')
    c_gw.run('aws s3 cp /tmp/packages.tar.gz ' + s3_path + '.deps.tar.gz')

    # Clean up
    c_gw.run('rm -r /tmp/packages')
    c_gw.run('rm /tmp/packages.tar.gz')
    c_gw.local('rm /tmp/packages.tar.gz')
    c_gw.local('rm /tmp/packages.txt')


def get_magma_version(c: Connection):
    return c.run(
        'directory=$(mktemp -d) &&'
        'dpkg-deb --extract ~/magma-packages/magma_[0-9]*.deb $directory &&'
        'source $directory/usr/local/share/magma/commit_hash &&'
        'echo $COMMIT_HASH',
    )


def copy_packages(c: Connection):
    """
    Copy the dependencies in the apt cache to /tmp on the local machine.
    """
    # Build a list of package metadata were each row looks like:
    #
    #   pkg_name version deb_name
    c.run('rm -rf /tmp/packages')
    c.run('rm -rf /tmp/packages.txt')
    c.run('mkdir -p /tmp/packages')
    c.run('cp ~/magma-packages/*.deb /tmp/packages')
    packages = c.run('ls /tmp/packages')

    with c.cd('/tmp/packages'):
        for f in packages.split():
            pkg_name = c.run(
                f'dpkg -I {f}'
                ' | grep " Package:"'
                ' | awk \'{print $2}\'',
            ).strip()
            version = c.run(
                f'dpkg -I {f}'
                ' | grep " Version:"'
                ' | awk \'{print $2}\'',
            ).strip()
            c.run(
                f'echo "{pkg_name} {version} {f}" >> /tmp/packages.txt',
            )

    c.run('mkdir -p /tmp/packages/executables')
    c.run('cp ~/magma-packages/executables/* /tmp/packages/executables')

    # Tar up the packages and executables
    with c.cd('/tmp/packages'):
        c.run('tar czf /tmp/packages.tar.gz ./*.deb ./executables')

    # Pull the artifacts onto the local machine
    c.get('/tmp/packages.tar.gz', '/tmp/packages.tar.gz')
    c.get('/tmp/packages.txt', '/tmp/packages.txt')
    magma_version = get_magma_version(c)
    c.run(f'echo "{magma_version}" > /tmp/magma_version')
