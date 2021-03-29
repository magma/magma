#!/usr/local/bin/python
# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from typing import List
import os
import sys
import argparse
import subprocess
import pprint
import yaml
import click
from ansible.cli.playbook import PlaybookCLI
from configlib import ConfigManager

constants = {}
def init_cli():
    global constants
    try:
        with open("/root/config.yml") as f:
            constants = yaml.load(f, Loader=yaml.FullLoader)
    except OSError:
        click.echo("Failed opening config.yml file")

    if not os.path.isdir(constants["config_dir"]):
        try:
            os.makedirs(constants["config_dir"])
        except OSError as error:
            click.echo("failed creating config directories ", error)
            sys.exit(1)

    if not os.path.isdir(constants["secret_dir"]):
        try:
            os.makedirs(constants["secret_dir"])
        except OSError as error:
            click.echo("failed creating secret directories ", error)
            sys.exit(1)
    return constants

# initiialize CLI
init_cli()

@click.group()
@click.version_option()
def cli():
    """Orchestrator Deployment CLI"""

########################## CONFIGURE COMMANDS ###################################
@cli.group(invoke_without_command=True)
@click.pass_context
@click.option('-c', '--component',
              type=click.Choice(constants['components'],case_sensitive=False),
              multiple=True,
              default=constants['components'])
def configure(ctx, component):
    """
    Configure enables user to configure values needed for Orc8r deployment.
    It additionally provides subcommands to see configured values, check
    if all required values are configured and finally provide detailed
    description on all options available to be configured
    """
    if ctx.invoked_subcommand is None:
        mgr = ConfigManager(constants)
        for c in component:
            mgr.configure(c)

@configure.command()
@click.option('-c', '--component',
              type=click.Choice(constants['components'],case_sensitive=False),
              multiple=True,
              default=constants['components'])
def show(component):
    mgr = ConfigManager(constants)
    for c in component:
        mgr.show(c)

@configure.command()
@click.option('-c', '--component',
              type=click.Choice(constants['components'],case_sensitive=False),
              multiple=True,
              default=constants['components'])
def info(component):
    mgr = ConfigManager(constants)
    for c in component:
        mgr.info(c)

@configure.command()
@click.option('-c', '--component',
              type=click.Choice(constants['components'],case_sensitive=False),
              multiple=True,
              default=constants['components'])
def check(component):
    mgr = ConfigManager(constants)
    valid = True
    for c in component:
        valid = mgr.check(c)
    if not valid:
        sys.exit(1)

@configure.command()
@click.option('-c', '--component',
              type=click.Choice(constants['components'],case_sensitive=False),
              prompt='component to configure')
@click.option('-k', '--key', prompt='name of the variable')
@click.option('-v', '--value', prompt='value of the variable')
def set(component, key, value):
    ConfigManager(constants).set(component, key, value)

########################## INSTALL COMMANDS ###################################
@cli.group()
def install():
    """
    Install command enables user to run subcommands in context of Orc8r installation.
    It provides subcommands like
    - precheck which runs various checks to ensure successful installation
    - addcerts which lets user create certificates for the installation
    """
    pass

@install.command('precheck')
def install_precheck():
    run_playbook([
        "ansible-playbook",
        "-v",
        "-e",
        "@/root/config.yml",
        "-t",
        "install_precheck",
        "%s/main.yml" % constants["playbooks"]])

@install.command()
@click.option('--self-signed', is_flag=True)
def addcerts(self_signed):
    if self_signed:
        run_playbook([
            "ansible-playbook",
            "-v",
            "-e",
            "@/root/config.yml",
            "-t",
            "addcerts",
            "%s/main.yml" % constants["playbooks"]])
    else:
        run_playbook([
            "ansible-playbook",
            "-v",
            "-e",
            "@/root/config.yml",
            "-t",
            "addcerts",
            "--skip-tags",
            "self_signed",
            "%s/main.yml" % constants["playbooks"]])


@cli.group()
def upgrade():
    """
    Upgrade command enables user to run subcommands in context of Orc8r upgrade.
    It provides subcommands like
    - precheck which runs various checks to ensure successful upgrade
    """
    pass

@upgrade.command('precheck')
def upgrade_precheck():
    run_playbook([
        "ansible-playbook",
        "-v",
        "-e",
        "@/root/config.yml",
        "-t",
        "upgrade_precheck",
        "%s/main.yml" % constants["playbooks"]])

def run_playbook(args):
    pb_cli = PlaybookCLI(args)
    pb_cli.parse()
    pb_cli.run()
