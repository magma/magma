#!/usr/bin/env python3
#
# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# This script creates the build context for the orc8r docker builds.
# It first creates a tmp directory, and then copies the cloud directories
# for all modules into it.

import argparse
import glob
import subprocess
from collections import namedtuple
from subprocess import PIPE
from typing import List

import os
import shutil
import yaml

BUILD_CONTEXT = '/tmp/magma_orc8r_build'
SRC_ROOT = 'src'
HOST_MAGMA_ROOT = '../../../.'
DEFAULT_MODULES_FILE = os.path.join(HOST_MAGMA_ROOT, 'modules.yml')
FB_MODULES_FILE = os.path.join(HOST_MAGMA_ROOT, 'fb/config/modules.yml')

COMPOSE_FILES = [
    'docker-compose.yml',
    'docker-compose.metrics.yml',

    # For now, this file is left out of the build because the fluentd daemonset
    # and forwarder pod shouldn't change very frequently - we can build and
    # push locally when they need to be updated.
    # We can integrate this into the CI pipeline if/when we see the need for it
    # 'docker-compose.logging.yml',

    'docker-compose.override.yml',
]

# Root directory where external modules will be mounted
GUEST_MODULE_ROOT = 'modules'
GUEST_MAGMA_ROOT = 'magma'

MagmaModule = namedtuple('MagmaModule', ['is_external', 'host_path', 'name'])


def main() -> None:
    args = _parse_args()

    # Create build context only if we're not mounting
    # If we're building, we always need to build controller first because proxy
    # copies metricsd binary and plugins from controller
    files_args = _get_docker_files_command_args(args)
    if not args.mount:
        _create_build_context()
    _build_cache_if_necessary(args)
    _build_controller_if_necessary(args)

    if args.mount:
        # Mount the source code and run a container with bash shell
        _run_docker(['run', '--rm'] + _get_mount_volumes() + ['test', 'bash'])
    elif args.tests:
        # Run unit tests
        _run_docker(['build', 'test'])
        _run_docker(['run', '--rm', 'test', 'make test'])
    elif args.nocache:
        # Build containers without go-cache in base image
        _run_docker(files_args + ['build'])
    else:
        # Build images using go-cache base image
        _run_docker(files_args + ['build',
                                  '--build-arg', 'baseImage=orc8r_cache'])


def _get_docker_files_command_args(args: argparse.Namespace) -> List[str]:
    if args.all:
        ret = []
        for f in COMPOSE_FILES:
            ret.append('-f')
            ret.append(f)
        return ret

    # docker-compose uses docker-compose.yml and docker-compose.override.yml
    # by default
    return []


def _build_cache_if_necessary(args: argparse.Namespace) -> None:
    if args.nocache or args.mount or args.tests:
        return

    # Check if orc8r_cache image exists
    result = subprocess.run(['docker', 'images', '-q', 'orc8r_cache'],
                            stdout=PIPE, stderr=PIPE)
    if result.stdout == b'':
        print("Orc8r_cache image does not exist. Building...")
        _run_docker(['-f', 'docker-compose.cache.yml', 'build'])


def _build_controller_if_necessary(args: argparse.Namespace) -> None:
    # We don't build the controller container if we're running tests or
    # generating code
    if args.mount or args.tests:
        return

    # controller will always only use docker-compose.yml and override so we
    # don't need to worry about file args (-f)
    if args.nocache:
        _run_docker(['build', 'controller'])
    else:
        _run_docker(['build', '--build-arg', 'baseImage=orc8r_cache',
                     'controller'])


def _run_docker(cmd: List[str]) -> None:
    """ Run the required docker-compose command """
    print("Running 'docker-compose %s'..." % " ".join(cmd))
    try:
        subprocess.run(['docker-compose'] + cmd, check=True)
    except subprocess.CalledProcessError as err:
        exit(err.returncode)


def _create_build_context() -> None:
    """ Clear out the build context from the previous run """
    if os.path.exists(BUILD_CONTEXT):
        shutil.rmtree(BUILD_CONTEXT)
    os.mkdir(BUILD_CONTEXT)

    print("Creating build context in '%s'..." % BUILD_CONTEXT)
    modules = []
    for module in _get_modules():
        _copy_module(module)
        modules.append(module.name)
    print('Context created for modules: %s' % ', '.join(modules))


def _copy_module(module: MagmaModule) -> None:
    """ Copy the module dir into the build context  """
    module_dest = _get_module_destination(module)
    dst = os.path.join(BUILD_CONTEXT, module_dest)

    # Copy relevant parts of the module to the build context
    shutil.copytree(
        os.path.join(module.host_path, 'cloud'),
        os.path.join(dst, 'cloud'),
    )

    if os.path.isdir(os.path.join(module.host_path, 'tools')):
        shutil.copytree(
            os.path.join(module.host_path, 'tools'),
            os.path.join(dst, 'tools'),
        )

    if os.path.isdir(os.path.join(module.host_path, 'cloud', 'configs')):
        shutil.copytree(
            os.path.join(module.host_path, 'cloud', 'configs'),
            os.path.join(BUILD_CONTEXT, 'configs', module.name),
        )

    # Copy the go.mod file for caching the go downloads
    # Use module_dest to preserve relative paths between go modules
    for filename in glob.iglob(dst + '/**/go.mod', recursive=True):
        gomod = filename.replace(
            dst, os.path.join(BUILD_CONTEXT, 'gomod', module_dest),
        )
        os.makedirs(os.path.dirname(gomod))
        shutil.copyfile(filename, gomod)


def _get_mount_volumes() -> List[str]:
    """ Return the volumes argument for docker-compose commands """
    volumes = []
    for module in _get_modules():
        module_mount_path = _get_module_destination(module)
        dst = os.path.join('/', module_mount_path)
        volumes.extend(['-v', '%s:%s' % (module.host_path, dst)])
    return volumes


def _get_modules() -> List[MagmaModule]:
    """
    Read the modules config file and return all modules specified.
    """
    filename = os.environ.get('MAGMA_MODULES_FILE', DEFAULT_MODULES_FILE)
    # Use the FB modules file if the file exists
    if os.path.isfile(FB_MODULES_FILE):
        filename = FB_MODULES_FILE
    modules = []
    with open(filename) as file:
        conf = yaml.safe_load(file)
        for module in conf['native_modules']:
            mod_path = os.path.abspath(os.path.join(HOST_MAGMA_ROOT, module))
            modules.append(
                MagmaModule(
                    is_external=False,
                    host_path=mod_path,
                    name=os.path.basename(mod_path),
                ),
            )
        for ext_module in conf['external_modules']:
            # NOTE: host_path for external modules is relative to the
            # $MAGMA_ROOT/orc8r/cloud directory on the host for legacy reasons.
            module_abspath = os.path.abspath(
                os.path.join(HOST_MAGMA_ROOT, 'orc8r', 'cloud',
                             ext_module['host_path']),
            )
            modules.append(
                MagmaModule(
                    is_external=True,
                    host_path=module_abspath,
                    name=os.path.basename(module_abspath),
                ),
            )
    return modules


def _get_module_destination(module: MagmaModule) -> str:
    """
    Given a path to a module on the host, return the destination to copy or
    mount the module to in the build context or container.
    """
    # The parent directory of the module should be the same on the host and
    # the guest for external modules
    if module.is_external:
        module_parent_dir = os.path.basename(
            os.path.abspath(os.path.join(module.host_path, os.path.pardir))
        )
        return os.path.join(SRC_ROOT, GUEST_MODULE_ROOT,
                            module_parent_dir, module.name)
    # We mount internal modules straight to MAGMA_ROOT as-is
    else:
        return os.path.join(SRC_ROOT, GUEST_MAGMA_ROOT, module.name)


def _parse_args() -> argparse.Namespace:
    """ Parse the command line args """
    parser = argparse.ArgumentParser(description='Orc8r build tool')
    parser.add_argument('--tests', '-t', action='store_true',
                        help="Run unit tests")
    parser.add_argument('--mount', '-m', action='store_true',
                        help='Mount the source code and create a bash shell')
    parser.add_argument('--nocache', '-n', action='store_true',
                        help='Build the images without go cache base image')
    parser.add_argument('--all', '-a', action='store_true',
                        help='Build all containers')
    args = parser.parse_args()
    return args


if __name__ == '__main__':
    main()
