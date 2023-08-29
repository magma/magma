#!/usr/bin/env python3

"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import argparse
import json
import os
import platform
import subprocess
from typing import List

HOST_MAGMA_ROOT = '../../../.'


def main() -> None:
    """ Run main"""
    args = _parse_args()

    if args.mount:
        _run(['up', '-d', 'test'])
        _run(['exec', 'test', 'bash'])
        _down(args)

    elif args.lint:
        _run(['up', '-d', 'test'])
        _run(['exec', 'test', 'make', 'lint'])
        _down(args)

    elif args.precommit:
        _run(['up', '-d', 'test'])
        _run(['exec', 'test', 'make', 'precommit'])
        _down(args)

    elif args.coverage:
        _run(['up', '-d', 'test'])
        _run(['exec', 'test', 'make', 'cover'])
        _down(args)

    elif args.tests:
        _run(['up', '-d', 'test'])
        _run(['exec', 'test', 'make', 'test'])
        _down(args)

    elif args.health:
        # _set_mac_env_vars is needed to override LOG_DRIVER for mac
        _set_mac_env_vars()
        _run(['-f', 'docker-compose.yml', '-f', 'docker-compose.override.yml', '-f', 'docker-compose.health.override.yml', 'up', '-d'])
        _run_health()
        _down(args)

    elif args.git:
        print(json.dumps(_run_get_git_vars(), indent=4, sort_keys=True))

    else:
        _run(['build'] + _get_default_build_args(args))
        _down(args)


def _run(cmd: List[str]) -> None:
    """ Run the required docker compose command """
    cmd = ['docker', 'compose', '--compatibility'] + cmd
    print("Running '%s'..." % ' '.join(cmd))
    try:
        subprocess.run(cmd, check=True)  # noqa: S603
    except subprocess.CalledProcessError as err:
        exit(err.returncode)


def _down(args: argparse.Namespace) -> None:
    if args.down:
        _run(['down'])


def _get_default_build_args(args: argparse.Namespace) -> List[str]:
    ret = []
    git_info = _run_get_git_vars()

    for arg, val in git_info.items():
        ret.append("--build-arg")
        ret.append("{0}={1}".format(arg, val))

    if args.nocache:
        ret.append('--no-cache')
    return ret


def _run_get_git_vars():
    try:
        cmd = "tools/get_version_info.sh"
        cmd_res = \
            subprocess.run(cmd, check=True, capture_output=True)  # noqa: S603
    except subprocess.CalledProcessError as err:
        print("Error _run_get_git_vars")
        exit(err.returncode)
    return json.loads(cmd_res.stdout)


def _run_health():
    try:
        cmd = "tools/docker_ps_healthcheck.sh"
        subprocess.run(cmd, check=True)  # noqa: S603
    except subprocess.CalledProcessError as err:
        print("Error _run_health")
        exit(err.returncode)


def _set_mac_env_vars():
    if (platform.system().lower() == "darwin"):
        os.environ['LOG_DRIVER'] = "json-file"


def _parse_args() -> argparse.Namespace:
    """ Parse the command line args """

    # There are multiple ways to invoke finer-grained control over which
    # images are built.
    #
    # (1) How many images to build
    #
    # all: all images
    # default: images required for minimum functionality
    #   - excluding metrics images
    #   - including postgres, proxy, etc
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
        '--precommit', '-c',
        action='store_true',
        help='Mount the source code and run pre-commit checks',
    )

    parser.add_argument(
        '--coverage', '-o',
        action='store_true',
        help='Generate test coverage statistics',
    )

    parser.add_argument(
        '--lint', '-l',
        action='store_true',
        help='Run lint test',
    )

    parser.add_argument(
        '--health', '-e',
        action='store_true',
        help='Run health test',
    )

    # Run something
    parser.add_argument(
        '--git', '-g',
        action='store_true',
        help='Get git info',
    )

    # How to do it
    parser.add_argument(
        '--nocache', '-n',
        action='store_true',
        help='Build the images with no Docker layer caching',
    )
    parser.add_argument(
        '--down', '-down',
        action='store_true',
        default=False,
        help='Leave containers up after running tests',
    )

    return parser.parse_args()


if __name__ == '__main__':
    main()
