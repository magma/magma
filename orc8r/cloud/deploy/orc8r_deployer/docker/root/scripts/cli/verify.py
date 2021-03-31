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
import glob

from ansible.cli.playbook import PlaybookCLI

from .install import run_playbook

@click.group()
@click.pass_context
def verify(ctx):
    """
    Upgrade command enables user to run subcommands in context of Orc8r upgrade.
    It provides subcommands like
    - precheck which runs various checks to ensure successful upgrade
    """
    pass

@verify.command('sanity')
@click.pass_context
def verify_sanity(ctx):
    # check if KUBECONFIG is set else find kubeconfig file and set the
    # environment variable
    constants = ctx.obj
    kubeconfig = os.environ.get('KUBECONFIG')
    if not kubeconfig:
        kubeconfigs = glob.glob(constants['project_dir'] + "/kubeconfig_*")
        if len(kubeconfigs) > 1:
            click.echo("multiple kubeconfigs found %s, "
                "unable to determine kubeconfig" % repr(kubeconfigs))
            return
        kubeconfig = kubeconfigs[0]

    os.environ["KUBECONFIG"] = kubeconfig
    os.environ["K8S_AUTH_KUBECONFIG"] = kubeconfig

    run_playbook([
        "ansible-playbook",
        "-v",
        "-e",
        "@/root/config.yml",
        "-t",
        "verify_sanity",
        "%s/main.yml" % constants["playbooks"]])