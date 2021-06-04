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

import json
import os
import pathlib
import subprocess
import sys

import click
import yaml


def init():
    constants = None
    try:
        with open("/root/config.yml") as f:
            constants = yaml.load(f, Loader=yaml.FullLoader)
    except OSError:
        print("Failed opening config.yml file")

    dirnames = (constants["config_dir"], constants["secret_dir"])
    for dirname in dirnames:
        pathlib.Path(dirname).mkdir(parents=True, exist_ok=True)

    return constants


def execute_command(cmd: list[str], cwd=None, env=None) -> int:
    """Execute command and `return error code

    Args:
        cmd (list[str]): list describing command to be run

    Returns:
        int: return code for the executed command
    """
    if env:
        env.update(os.environ)
    with subprocess.Popen(cmd, stdout=subprocess.PIPE, cwd=cwd, env=env) as p:
        for output in p.stdout:
            click.echo(output, nl=False)
        return p.wait()
    return 1


def run_command(cmd: list[str]) -> subprocess.CompletedProcess:
    """Run command and capture output with string encoding

    Args:
        cmd (list[str]): list describing command to be run

    Returns:
        subprocess.CompletedProcess: return value from run
    """
    return subprocess.run(cmd, encoding='utf-8', capture_output=True)


def get_json(fn: str) -> dict:
    try:
        with open(fn) as f:
            return json.load(f)
    except OSError:
        pass
    return {}


def put_json(fn: str, cfgs: dict):
    with open(fn, 'w') as outfile:
        json.dump(cfgs, outfile)
