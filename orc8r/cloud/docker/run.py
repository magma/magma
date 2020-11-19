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

import argparse
import fileinput
import subprocess
import sys

from typing import List

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


def main() -> None:
    args = _parse_args()

    if args.clear:
        _clear_line('.env', 'COMPOSE_FILE=')

    files = DEPLOYMENT_TO_MODULES[args.deployment]
    if args.metrics:
        files.append('metrics')
    if args.thanos:
        files.append('thanos')

    file_args = _make_file_args(files)
    compose_line = 'COMPOSE_FILE=%s' % ':'.join(file_args)
    _add_or_replace_line('.env', 'COMPOSE_FILE=', compose_line)

    if args.print:
        return

    cmd = ['docker-compose'] + _make_cmd_args(file_args) + ['up', '-d']
    print("Running '%s'..." % ' '.join(cmd))
    try:
        subprocess.run(cmd, check=True)
    except subprocess.CalledProcessError as err:
        exit(err.returncode)


def _clear_line(file: str, search: str) -> None:
    with open(file, 'r') as f:
        lines = f.readlines()
    with open(file, 'w') as f:
        for l in lines:
            if search not in l:
                f.write(l)


def _add_or_replace_line(file: str, search: str, replace: str) -> None:
    found = False
    for l in fileinput.input(file, inplace=1):
        if search in l:
            found = True
            l = replace
        sys.stdout.write(l)
    if not found:
        with open(file, 'a') as f:
            f.write(replace)


def _make_file_args(files: List[str]) -> List[str]:
    files = ['docker-compose.yml'] + \
            ['docker-compose.%s.yml' % f for f in files] + \
            ['docker-compose.override.yml']
    return files


def _make_cmd_args(files: List[str]) -> List[str]:
    args = []
    for f in files:
        args.extend(['-f', f])
    return args


def _parse_args() -> argparse.Namespace:
    """ Parse the command line args """
    parser = argparse.ArgumentParser(description='Orc8r run tool')

    # Config-only
    parser.add_argument(
        '--print', '-p',
        action='store_true',
        help="Only update COMPOSE_FILE line from the .env file",
    )

    # Set docker-compose files
    parser.add_argument(
        '--clear', '-c',
        action='store_true',
        help='Clear COMPOSE_FILE line from the .env file',
    )
    parser.add_argument(
        '--deployment', '-d',
        action='store',
        default='all',
        help='Activate deployment type: %s' % ','.join(DEPLOYMENTS),
    )
    parser.add_argument(
        '--metrics',
        action='store_true',
        help='Include docker-compose.metrics.yml',
    )
    parser.add_argument(
        '--thanos',
        action='store_true',
        help='Include docker-compose.thanos.yml',
    )

    args = parser.parse_args()
    return args


if __name__ == '__main__':
    main()
