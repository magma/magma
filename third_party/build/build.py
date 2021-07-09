#!/usr/bin/env python3
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


import argparse
import copy
import logging
import os
import re
import subprocess
import sys

log = logging.getLogger(__name__)

SUDO = 'sudo'


"""
sudo apt install ruby ruby-dev build-essential rubygems libffi-dev python3-setuptools python-setuptools libssl-dev

x E: Unable to locate package grpc-dev
E: Unable to locate package magma-cpp-redis
x E: Unable to locate package prometheus-cpp-dev
x E: Unable to locate package libfolly-dev
x E: Unable to locate package magma-libfluid

"""


def find_magma_root():
    path = os.path.realpath(__file__)
    m = re.match(r'(?P<magma>.*/magma/).*', path)
    if m:
        return m.group('magma')
    return None


def os_release():
    release_info = {}
    with open('/etc/os-release', 'r') as f:
        for line in f:
            try:
                k, v = line.rstrip().split('=')
                release_info[k] = v.strip('"')
            except Exception:
                pass
    return release_info


def packagemanager():
    release_info = os_release()
    os_type = release_info['ID']
    if os_type in ['centos', 'redhat']:
        return 'yum'
    else:
        return 'apt'


def pkgfmt():
    return 'deb' if packagemanager() == 'apt' else 'rpm'


def strsplitbytes(the_bytes):
    decoded = the_bytes.decode('utf-8')
    return [item for item in re.split(r'\s+', decoded) if item]


def buildscript(package_name):
    return './bin/' + package_name + '_build.sh'


def buildafter(package_name, env=None):
    script = buildscript(package_name)
    pre = strsplitbytes(
        subprocess.check_output(
            [script, '-A'],
            env=env,
        ),
    )
    return pre


def buildrequires(package_name, env=None):
    script = buildscript(package_name)
    req = strsplitbytes(
        subprocess.check_output(
            [script, '-B'],
            env=env,
        ),
    )
    return req


def build(package_name, env=None, install=True):
    script = buildscript(package_name)
    outputfilename = subprocess.check_output(
        [script, '-F'], env=env,
    ).decode('utf-8').strip()
    if not os.path.exists('./' + outputfilename):
        subprocess.run([script], check=True, env=env)
    else:
        log.info('found {}; skipping'.format(outputfilename))
    if install:
        subprocess.run(
            [SUDO, packagemanager(), 'install', '-y', './' + outputfilename],
            check=True,
        )


def main(args):
    env = copy.copy(os.environ)
    release_info = os_release()
    arch_map = {
        'x86_64': 'amd64',
        'aarch64': 'arm64',
    }
    arch = subprocess.check_output(['uname', '-m'])
    arch = arch.strip().decode('utf-8')
    if arch in arch_map:
        arch = arch_map[arch]
        if arch == 'amd64' and pkgfmt() == 'rpm':
            arch = 'x86_64'
    if 'MAGMA_ROOT' not in env:
        magma_root = find_magma_root()
        if magma_root:
            env['MAGMA_ROOT'] = magma_root

    env['ARCH'] = arch
    env['PKGFMT'] = pkgfmt()
    env['OS_RELEASE'] = ''
    os_id = release_info.get('ID', '')
    os_version = release_info.get('VERSION_ID', '')
    if os_id:
        env['OS_RELEASE'] = os_id + os_version

    packages = args.package
    to_install = set()
    all_packages = set()
    all_packages.update(packages)
    depmap = {}

    for package in packages:
        items = set(buildafter(package))
        depmap[package] = items
        all_packages.update(items)

    ordered_packages = []

    order_updated = True
    while order_updated:
        order_updated = False
        # packages with no outstanding deps may be added to build order
        ready = [p for p in all_packages if not depmap.get(p)]
        all_packages -= set(ready)
        if ready:
            order_updated = True
            ordered_packages.extend(sorted(ready))
            for p in ready:
                depmap.pop(p, None)
            for a in depmap.keys():
                for p in ready:
                    depmap[a].discard(p)

    if depmap:
        raise Exception('unprocessed dependencies: {}'.format(depmap))

    for package in ordered_packages:
        to_install.update(buildrequires(package))
    subprocess.run(
        [SUDO, packagemanager(), 'install', '-y'] + list(to_install),
    )

    for package in ordered_packages:
        build(package, env=env, install=not args.no_install)


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('package', nargs='*')
    parser.add_argument(
        '-N', '--no-install', action='store_true',
        help='Skip install of resulting packages to build system',
    )

    args = parser.parse_args()
    main(args)
