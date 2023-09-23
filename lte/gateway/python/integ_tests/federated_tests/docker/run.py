#!/usr/bin/env python3

"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import argparse
import os
import subprocess
import sys

MAGMA_ROOT = os.environ["MAGMA_ROOT"]
SNOWFLAKE_PATH = MAGMA_ROOT + '/.cache/feg/'
SNOWFLAKE_FILE = MAGMA_ROOT + '/.cache/feg/snowflake'


def main() -> None:
    """ create a snowflake file if necessary, then start docker containers """
    args = _parse_args()
    if not os.path.isfile(SNOWFLAKE_FILE):
        _create_snowflake_file()
    _exec_docker_cmd(args)


def _create_snowflake_file() -> None:
    if os.path.isdir(SNOWFLAKE_FILE):
        _exec_cmd(['rm', '-r', SNOWFLAKE_FILE])
    print("Creating snowflake file")
    _exec_cmd(['mkdir', '-p', SNOWFLAKE_PATH])
    _exec_cmd(['touch', SNOWFLAKE_FILE])


def _exec_docker_cmd(args) -> None:
    cmd = ['docker', 'compose', '--compatibility', 'up', '-d']
    if args.down:
        cmd = ['docker', 'compose', 'down']
    print(f"Running {' '.join(cmd)}...")
    _exec_cmd(cmd)


def _exec_cmd(cmd) -> None:
    try:
        subprocess.run(cmd, check=True)
    except subprocess.CalledProcessError as err:
        sys.exit(err.returncode)


def _parse_args() -> argparse.Namespace:
    """ Parse the command line args """
    parser = argparse.ArgumentParser(description='FeG run tool')

    # Other actions
    parser.add_argument(
        '--down', '-d',
        action='store_true',
        help='Stop running containers',
    )
    args = parser.parse_args()
    return args


if __name__ == '__main__':
    main()
