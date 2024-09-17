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
import sys

import click
from cli.style import print_error_msg, print_success_msg, print_warning_msg
from utils.ansiblelib import AnsiblePlay, run_playbook
from utils.common import execute_command


@click.group(invoke_without_command=True)
@click.pass_context
def upgrade(ctx):
    """
    Upgrade existing orc8r deployment
    """
    tf_cmds = [
        ["terraform", "init", "--upgrade"],
        ["terraform", "refresh"],
        ["terraform", "apply", "-auto-approve"],
    ]

    if ctx.invoked_subcommand is None:
        if click.confirm('Do you want to run upgrade prechecks?'):
            ctx.invoke(precheck)
        else:
            print_warning_msg(f"Skipping upgrade prechecks")

        click.echo(
            "Following commands will be run during upgrade\n%s" % (
                "\n".join((map(" ".join, tf_cmds)))
            ),
        )
        for cmd in tf_cmds:
            if click.confirm(
                'Do you want to continue with %s?' %
                " ".join(cmd),
            ):
                rc = execute_command(cmd)
                if rc != 0:
                    print_error_msg("Upgrade Failed!!!")
                    return


def precheck_cmd(constants: dict) -> list:
    """Get the arg list to run prechecks

    Args:
        constants ([dict]): config dict
    """
    playbook_dir = constants["playbooks"]
    return AnsiblePlay(
        playbook=f"{playbook_dir}/main.yml",
        tags=['upgrade_precheck'],
        extra_vars=constants,
    )


@upgrade.command()
@click.pass_context
def precheck(ctx):
    """
    Performs various checks to ensure successful upgrade
    """
    rc = run_playbook(precheck_cmd(ctx.obj))
    if rc != 0:
        print_error_msg("Upgrade prechecks failed!!!")
        sys.exit(1)
    print_success_msg("Upgrade prechecks ran successfully")
