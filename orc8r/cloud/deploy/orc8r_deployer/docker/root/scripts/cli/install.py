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
import glob
import os
import sys

import click
from cli.style import (
    print_error_msg,
    print_info_msg,
    print_success_msg,
    print_warning_msg,
)
from utils.ansiblelib import AnsiblePlay, run_playbook
from utils.common import execute_command


def tf_install(constants: dict, warn: bool = True,
               max_retries: int = 2) -> int:
    """Run through terraform installation

    Args:
        constants (dict): config dict
        warn (bool, optional): require user confirmation. Defaults to True.
        max_retries (int): Number of times to retry in case of a failure.

    Returns:
        int: return code
    """

    tf_init = ["terraform", "init"]
    tf_orc8r = ["terraform", "apply", "-target=module.orc8r", "-auto-approve"]
    tf_secrets = [
        "terraform",
        "apply",
        "-target=module.orc8r-app.null_resource.orc8r_seed_secrets",
        "-auto-approve"]
    tf_orc8r_app = ["terraform", "apply", "-auto-approve"]

    for tf_cmd in [tf_init, tf_orc8r, tf_secrets, tf_orc8r_app]:
        cmd = " ".join(tf_cmd)
        if warn and not click.confirm(f'Do you want to continue with {cmd}?'):
            continue

        for i in range(max_retries):
            # terraform fails randomly due to timeouts
            click.echo(f"Running {tf_cmd}, iteration {i}")
            rc = execute_command(tf_cmd, cwd=constants['project_dir'])
            if rc == 0:
                break
            print_error_msg(f"Install failed when running {cmd} !!!")
            if i == (max_retries - 1):
                print_error_msg(f"Max retries exceeded!!!")
                return 1

        # set the kubectl after bringing up the infra
        if tf_cmd in (tf_orc8r, tf_orc8r_app):
            kubeconfigs = glob.glob(
                constants['project_dir'] + "/kubeconfig_*")
            if len(kubeconfigs) != 1:
                print_error_msg(
                    "zero or multiple kubeconfigs found %s!!!" %
                    repr(kubeconfigs))
                return
            kubeconfig = kubeconfigs[0]
            os.environ['KUBECONFIG'] = kubeconfig
            print_info_msg(
                'For accessing kubernetes cluster, set'
                f' `export KUBECONFIG={kubeconfig}`')

        print_success_msg(f"Command {cmd} ran successfully")
    else:
        print_warning_msg(f"Skipping Command {cmd}")
    return 0


@click.group(invoke_without_command=True)
@click.pass_context
def install(ctx):
    """
    Deploy new instance of orc8r
    """
    if ctx.invoked_subcommand is None:
        if click.confirm('Do you want to run installation prechecks?'):
            ctx.invoke(precheck)
        else:
            print_warning_msg(f"Skipping installation prechecks")

        tf_install(ctx.obj)


def precheck_cmd(constants: dict) -> list:
    """Get the arg list to run prechecks

    Args:
        constants ([dict]): config dict
    """
    playbook_dir = constants["playbooks"]
    return AnsiblePlay(
        playbook=f"{playbook_dir}/main.yml",
        tags=['install_precheck'],
        extra_vars=constants)


@install.command()
@click.pass_context
def precheck(ctx):
    """
    Performs various checks to ensure successful installation
    """
    rc = run_playbook(precheck_cmd(ctx.obj))
    if rc != 0:
        print_error_msg("Install prechecks failed!!!")
        sys.exit(1)
    print_success_msg("Install prechecks ran successfully")
