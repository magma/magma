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
import fileinput
import pathlib
import subprocess
import sys
from typing import List

HOST_BUILD_CTX = '/tmp/magma_orc8r_build'
DO_NOT_COMMIT = '# DO NOT COMMIT THIS CHANGE'


def main() -> None:
    args = _parse_args()

    if args.clear:
        _clear_line('.env', 'COMPOSE_FILE=')
        return

    default = not (args.metrics or args.thanos)
    if default:
        _clear_line('.env', 'COMPOSE_FILE=')
    else:
        f = []
        if args.metrics:
            f.append('metrics')
        if args.thanos:
            f.append('thanos')
        file_args = _make_file_args(f)
        compose_line = 'COMPOSE_FILE={}  {}\n'.format(
            ':'.join(file_args),
            DO_NOT_COMMIT,
        )
        _add_or_replace_line('.env', 'COMPOSE_FILE=', compose_line)

    if args.print:
        return

    # Ensure build context exists, otherwise docker-compose throws an error
    pathlib.Path(HOST_BUILD_CTX).mkdir(parents=True, exist_ok=True)

    cmd = ['docker-compose', 'up', '-d']
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
