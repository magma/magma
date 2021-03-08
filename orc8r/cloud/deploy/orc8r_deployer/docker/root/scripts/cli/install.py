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
from typing import List
import os
import sys
import argparse
import subprocess
import pprint
import yaml
import click

from ansible.cli.playbook import PlaybookCLI

def run_playbook(args):
    pb_cli = PlaybookCLI(args)
    pb_cli.parse()
    pb_cli.run()

@click.group()
@click.pass_context
def install(ctx):
    """
    Install command enables user to run subcommands in context of Orc8r installation.
    """
    pass

@install.command()
@click.pass_context
def precheck(ctx):
    """
    Precheck which runs various checks to ensure successful installation
    """
    run_playbook([
        "ansible-playbook",
        "-v",
        "-e",
        "@/root/config.yml",
        "-t",
        "install_precheck",
        "%s/main.yml" % ctx.obj["playbooks"]])

@install.command()
@click.pass_context
@click.option('--self-signed', is_flag=True)
def addcerts(ctx, self_signed):
    """
    Addcerts which lets user create certificates for the installation
    --self-signed option enables creation of self signed root certs
    By default only application certs are created.
    """
    if self_signed:
        run_playbook([
            "ansible-playbook",
            "-v",
            "-e",
            "@/root/config.yml",
            "-t",
            "addcerts",
            "%s/main.yml" % ctx.obj["playbooks"]])
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
            "%s/main.yml" % ctx.obj["playbooks"]])