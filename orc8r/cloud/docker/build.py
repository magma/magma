#!/usr/bin/env python3.7

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

# build.py creates the build context for the orc8r Docker builds.
# It first creates a tmp directory, and then copies the cloud directories
# for all modules into it.

import argparse
import glob
import os
import shutil
import subprocess
from collections import namedtuple
from typing import Iterable, List

HOST_BUILD_CTX = '/tmp/magma_orc8r_build'
HOST_MAGMA_ROOT = '../../../.'
IMAGE_MAGMA_ROOT = os.path.join('src', 'magma')

MODULES = [
    'orc8r',
    'lte',
    'feg',
    'cwf',
    'wifi',
    'fbinternal',
]

DEPLOYMENT_TO_MODULES = {
    'all': MODULES,
    'orc8r': [],
    'orc8r-f': ['fbinternal'],
    'fwa': ['lte'],
    'fwa-f': ['lte', 'fbinternal'],
    'ffwa': ['lte', 'feg'],
    'ffwa-f': ['lte', 'feg', 'fbinternal'],
    'cwf': ['lte', 'feg', 'cwf'],
    'cwf-f': ['lte', 'feg', 'cwf', 'fbinternal'],
    'wifi': ['wifi'],
    'wifi-f': ['wifi', 'fbinternal'],
}

DEPLOYMENTS = DEPLOYMENT_TO_MODULES.keys()

EXTRA_COMPOSE_FILES = [
    'docker-compose.metrics.yml',
    # For now, logging is left out of the build because the fluentd daemonset
    # and forwarder pod shouldn't change very frequently - we can build and
    # push locally when they need to be updated.
    # We can integrate this into the CI pipeline if/when we see the need for it
    # 'docker-compose.logging.yml',
]

MagmaModule = namedtuple('MagmaModule', ['name', 'host_path'])


def main() -> None:
    args = _parse_args()
    mods = _get_modules(DEPLOYMENT_TO_MODULES[args.deployment])

    if not args.extras:
        _create_build_context(mods)

    if args.mount:
        _run(['run', '--rm'] + _get_mnt_vols(mods) + ['test', 'bash'])
        _down(args)
    elif args.generate:
        _run(['run', '--rm'] + _get_mnt_vols(mods) + ['test', 'make gen'])
        _down(args)
    elif args.tests:
        _run(['up', '-d', 'postgres_test'])
        _run(['build', 'test'])
        _run(['run', '--rm', 'test', 'make test'])
        _down(args)
    else:
        d_args = _get_default_file_args(args) + _get_default_build_args(args)
        _run(d_args)


def _get_modules(mods: Iterable[str]) -> Iterable[MagmaModule]:
    """
    Read the modules config file and return all modules specified.
    """
    modules = []
    for m in mods:
        abspath = os.path.abspath(os.path.join(HOST_MAGMA_ROOT, m))
        module = MagmaModule(name=m, host_path=abspath)
        modules.append(module)
    return modules


def _create_build_context(modules: Iterable[MagmaModule]) -> None:
    """ Clear out the build context from the previous run """
    shutil.rmtree(HOST_BUILD_CTX, ignore_errors=True)
    os.mkdir(HOST_BUILD_CTX)

    print("Creating build context in '%s'..." % HOST_BUILD_CTX)
    for m in modules:
        _copy_module(m)


def _down(args: argparse.Namespace) -> None:
    if not args.leave:
        _run(['down'])


def _run(cmd: List[str]) -> None:
    """ Run the required docker-compose command """
    cmd = ['docker-compose'] + cmd
    print("Running '%s'..." % ' '.join(cmd))
    try:
        subprocess.run(cmd, check=True)
    except subprocess.CalledProcessError as err:
        exit(err.returncode)


def _get_mnt_vols(modules: Iterable[MagmaModule]) -> List[str]:
    """ Return the volumes argument for docker-compose commands """
    vols = []
    for m in modules:
        vols.extend(['-v', '%s:%s' % (m.host_path, _get_module_image_dst(m))])
    return vols


def _get_default_file_args(args: argparse.Namespace) -> List[str]:
    def make_file_args(fs: List[str]) -> List[str]:
        fs = ['docker-compose.yml'] + fs + ['docker-compose.override.yml']
        ret = []
        for f in fs:
            ret.extend(['-f', f])
        return ret

    def get_files_for_modules(ms: Iterable[str]) -> List[str]:
        return ['docker-compose.%s.yml' % m for m in ms if m != 'orc8r']

    mods = DEPLOYMENT_TO_MODULES[args.deployment]

    if args.all:
        all_files = get_files_for_modules(mods) + EXTRA_COMPOSE_FILES
        return make_file_args(all_files)

    if args.extras:
        return make_file_args(EXTRA_COMPOSE_FILES)

    # Default to docker-compose.yml + all modules + docker-compose.override.yml
    return make_file_args(get_files_for_modules(mods))


def _get_default_build_args(args: argparse.Namespace) -> List[str]:
    mods = DEPLOYMENT_TO_MODULES[args.deployment]
    ret = [
        'build',
        '--build-arg', 'MAGMA_MODULES=%s' % ' '.join(mods),
    ]
    if args.parallel:
        ret.append('--parallel')
    if args.nocache:
        ret.append('--no-cache')
    return ret


def _copy_module(module: MagmaModule) -> None:
    """ Copy module directory into the build context  """
    build_ctx = _get_module_host_dst(module)

    def copy_to_ctx(d: str) -> None:
        shutil.copytree(
            os.path.join(module.host_path, d),
            os.path.join(build_ctx, d),
        )

    copy_to_ctx('cloud')

    # Orc8r module also has lib/ and gateway/
    if module.name == 'orc8r':
        copy_to_ctx('lib')
        copy_to_ctx('gateway')

    # Optionally copy cloud/configs/
    # Converts e.g. lte/cloud/configs/ to configs/lte/
    if os.path.isdir(os.path.join(module.host_path, 'cloud', 'configs')):
        shutil.copytree(
            os.path.join(module.host_path, 'cloud', 'configs'),
            os.path.join(HOST_BUILD_CTX, 'configs', module.name),
        )

    # Copy the go.mod file for caching the go downloads
    # Preserves relative paths between modules
    for f in glob.iglob(build_ctx + '/**/go.mod', recursive=True):
        gomod = f.replace(
            HOST_BUILD_CTX, os.path.join(HOST_BUILD_CTX, 'gomod'),
        )
        print(gomod)
        os.makedirs(os.path.dirname(gomod))
        shutil.copyfile(f, gomod)


def _get_module_image_dst(module: MagmaModule) -> str:
    """
    Given a path to a module on the host, return the intended destination
    in the final image.
    """
    return os.path.join(os.sep, IMAGE_MAGMA_ROOT, module.name)


def _get_module_host_dst(module: MagmaModule) -> str:
    """
    Given a path to a module on the host, return the intended destination
    in the build context.
    """
    return os.path.join(HOST_BUILD_CTX, IMAGE_MAGMA_ROOT, module.name)


def _parse_args() -> argparse.Namespace:
    """ Parse the command line args """

    # There are multiple ways to invoke finer-grained control over which
    # images are built.
    #
    # (1) How many images to build
    #
    # all: all images
    # extras: inverse of default
    # default: images required for minimum functionality
    #   - excluding metrics images
    #   - including postgres, proxy, etc
    # core: only orc8r-specific images, across all modules
    #   - excluding postgres, proxy, etc
    #   - including orc8r, lte, cwf, etc
    #
    # (2) Of the core orc8r images, which modules to build
    #
    # Defaults to all modules, but can be further specified by targeting a
    # deployment type.

    parser = argparse.ArgumentParser(description='Orc8r build tool')

    # Run something
    parser.add_argument(
        '--tests', '-t',
        action='store_true',
        help='Run unit tests',
    )
    parser.add_argument(
        '--mount', '-m',
        action='store_true',
        help='Mount the source code and create a bash shell',
    )
    parser.add_argument(
        '--generate', '-g',
        action='store_true',
        help='Mount the source code and regenerate generated files',
    )

    # Build something
    parser.add_argument(
        '--all', '-a',
        action='store_true',
        help='Build all containers: core and extras',
    )
    parser.add_argument(
        '--extras', '-e',
        action='store_true',
        help='Build extras (non-essential) images (i.e. no proxy or lte)',
    )
    parser.add_argument(
        '--deployment', '-d',
        action='store',
        default='all',
        help='Build deployment type: %s' % ','.join(DEPLOYMENTS),
    )

    # How to do it
    parser.add_argument(
        '--nocache', '-n',
        action='store_true',
        help='Build the images with no Docker layer caching',
    )
    parser.add_argument(
        '--parallel', '-p',
        action='store_true',
        default=False,
        help='Build containers in parallel',
    )
    parser.add_argument(
        '--leave', '-l',
        action='store_true',
        help='Leave containers running after running tests',
    )

    args = parser.parse_args()
    return args


if __name__ == '__main__':
    main()
