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
import os
import sys
from pathlib import Path
from shutil import copyfile

import click
from boto3 import Session
from cli.configlib import get_input
from cli.style import print_error_msg
from utils.ansiblelib import AnsiblePlay, run_playbook
from utils.common import execute_command


def setup_aws_environ():
    """Set up aws configuration attributes in environment"""
    session = Session()
    creds = session.get_credentials()
    if not creds or not session.region_name:
        print_error_msg('''
AWS credentials not configured.
configure through awscli or through orcl
orcl configure set -k aws_access_key_id -v <access_key_id>
orcl configure set -k aws_secret_access_key -v <aws_secret_access_key>
orcl configure set -k region -v <region>
''')
        sys.exit(1)

    frozen_creds = creds.get_frozen_credentials()
    os.environ["AWS_ACCESS_KEY_ID"] = frozen_creds.access_key
    os.environ["AWS_SECRET_ACCESS_KEY"] = frozen_creds.secret_key
    os.environ["AWS_REGION"] = session.region_name


def tf_state_fn(tf_dir):
    return f'{tf_dir}/terraform.tfstate'


def tf_backup_fn(tf_dir):
    return f'{tf_dir}/terraform.tfstate.golden'


def tf_destroy(constants: dict, warn: bool = True,
               max_retries: int = 2) -> int:
    """Run through terraform cleanup

    Args:
        constants (dict): Config definitions
        warn (bool): require user confirmation. Defaults to True.
        max_retries (int): Number of times to retry in case of a failure.
    Returns:
        int: Return code
    """
    if warn and not click.confirm(
            'Do you want to continue with cleanup?', abort=True):
        return 0

    # backup existing terraform state
    project_dir = constants['project_dir']
    try:
        copyfile(tf_state_fn(project_dir), tf_backup_fn(project_dir))
    except OSError:
        print_error_msg('Unable to backup terraform state')
        return 1

    tf_destroy_cmds = ["terraform", "destroy", "-auto-approve"]
    cmd = " ".join(tf_destroy_cmds)
    for i in range(max_retries):
        click.echo(f"Running {cmd}, iteration {i}")
        rc = execute_command(tf_destroy_cmds, cwd=project_dir)
        if rc == 0:
            break
        print_error_msg("Destroy Failed!!!")
        if i == (max_retries - 1):
            print_error_msg("Max retries exceeded!!! Attempt cleaning up using"
                            " 'orcl cleanup raw' subcommand")
            return 1
    return 0


@click.group(invoke_without_command=True)
@click.pass_context
def cleanup(ctx):
    """
    Removes resources deployed for orc8r
    """
    if ctx.invoked_subcommand is None:
        tf_destroy(ctx.obj)


def cleanup_cmd(constants: dict, dryrun: bool = False) -> list:
    """Get the arg list to run cleanup resources

    Args:
        constants (dict): config dict
        dryrun (bool): flag to indicate dryrun. Defaults to False.

    Returns:
        list: command list
    """
    playbook_dir = constants["playbooks"]
    return AnsiblePlay(
        playbook=f"{playbook_dir}/cleanup.yml",
        tags=['cleanup_dryrun'] if dryrun else ['cleanup'],
        extra_vars=constants)


def raw_cleanup(
        constants: dict,
        override_dict: dict = None,
        dryrun: bool = False,
        max_retries: int = 2):
    """Perform raw cleanup of resources using internal commands

    Args:
        constants (dict): config dict
        overrides (dict): overide dict
        dryrun (bool): flag to indicate dryrun. Defaults to False.
        max_retries (int): maximum number of retries
    Returns:
        list: command list
    """
    if not override_dict and not constants.get('cleanup_state'):
        backup_fn = tf_backup_fn(constants['project_dir'])
        if Path(backup_fn).exists():
            constants['cleanup_state'] = backup_fn
    if override_dict:
        constants.update(override_dict)

    # sometimes cleanups might not fully happen due to timing related
    # resource dependencies. Run it few times to eliminate all resources
    # completely
    for i in range(max_retries):
        rc = run_playbook(cleanup_cmd(constants, dryrun))
        if rc != 0:
            print_error_msg("Failed cleaning up resources!!!")


@cleanup.command()
@click.pass_context
@click.option('--dryrun', default=False, is_flag=True, help='Show resources '
              'to be cleaned up during raw cleanup')
@click.option('--state', help='Provide state file containing resource '
              'information e.g. terraform.tfstate or '
              'terraform.tfstate.backup')
@click.option('--override', default=False, is_flag=True, help='Provide values'
              'to cleanup the orc8r deployment')
def raw(ctx, dryrun, state, override):
    """
    Individually cleans up resources deployed for orc8r
    Attributes:
    ctx: Click context
    dryrun: knob to enable dryrun of the cleanup to be performed
    state: location of the terraform state file
    override: override any state information with custom values
    """
    # Few boto dependent modules in ansible require these values to be
    # setup as environment variables. Hence setting these up.
    setup_aws_environ()

    if not dryrun:
        click.confirm(click.style('This is irreversable!! Do you want to '
                      'continue with cleanup?', fg='red'), abort=True)
    if state:
        ctx.obj['cleanup_state'] = state

    override_dict = None
    if override:
        override_dict = {
            'orc8r_namespace': 'orc8r',
            'orc8r_secrets': 'orc8r-secrets',
            'orc8r_es_domain': 'orc8r-es',
            'orc8r_cluster_name': 'orc8r',
            'orc8r_db_id': 'orc8rdb',
            'orc8r_db_subnet': 'orc8r_vpc',
            'vpc_name': 'orc8r_vpc',
            'region_name': os.environ["AWS_REGION"],
            'efs_fs_targets': '',
            'efs_mount_targets': '',
            'domain_name': '',
        }
        for k, v in override_dict.items():
            inp = get_input(k, v)
            inp_entries = inp.split(',')
            if len(inp_entries) > 1:
                # mainly relevant for passing in list of mount and fs targets
                override_dict[k] = inp_entries
            else:
                override_dict[k] = inp_entries[0]
    raw_cleanup(ctx.obj, override_dict, dryrun)
