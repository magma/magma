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
import pathlib
import subprocess
import sys
from typing import List

import click
import yaml

from .certs import certs
from .cleanup import cleanup
from .configure import configure
from .install import install
from .upgrade import upgrade
from .verify import verify


def init():
    constants = None
    try:
        with open("/root/config.yml") as f:
            constants = yaml.load(f, Loader=yaml.FullLoader)
    except OSError:
        click.echo("Failed opening config.yml file")

    dirnames = (constants["config_dir"], constants["secret_dir"])
    for dirname in dirnames:
        try:
            pathlib.Path(dirname).mkdir(parents=True, exist_ok=True)
        except OSError as error:
            click.echo(f"failed creating dir {dirname} error {error}")
            sys.exit(1)
    return constants


@click.group()
@click.version_option()
@click.pass_context
def cli(ctx):
    """
    Orchestrator Deployment CLI
    """
    ctx.obj = init()


cli.add_command(configure)
cli.add_command(certs)
cli.add_command(install)
cli.add_command(upgrade)
cli.add_command(verify)
cli.add_command(cleanup)
