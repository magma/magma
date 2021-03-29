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

from .install import run_playbook

@click.group()
@click.pass_context
def upgrade(ctx):
    """
    Upgrade command enables user to run subcommands in context of Orc8r upgrade.
    """
    pass

@upgrade.command()
@click.pass_context
def precheck(ctx):
    """
    Precheck runs various checks to ensure successful upgrade
    """

    run_playbook([
        "ansible-playbook",
        "-v",
        "-e",
        "@/root/config.yml",
        "-t",
        "upgrade_precheck",
        "%s/main.yml" % ctx.obj["playbooks"]])