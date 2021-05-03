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

from .common import (
    print_error_msg,
    print_info_msg,
    print_success_msg,
    print_warning_msg,
    run_command,
    run_playbook,
)


@click.group(invoke_without_command=True)
@click.pass_context
def install(ctx):
    """
    Deploy new instance of orc8r
    """
    constants = ctx.obj

    tf_init = ["terraform", "init"]
    tf_orc8r = ["terraform", "apply", "-target=module.orc8r", "-auto-approve"]
    tf_secrets = [
        "terraform",
        "apply",
        "-target=module.orc8r-app.null_resource.orc8r_seed_secrets",
        "-auto-approve"]
    tf_orc8r_app = ["terraform", "apply", "-auto-approve"]

    if ctx.invoked_subcommand is None:
        if click.confirm('Do you want to run installation prechecks?'):
            ctx.invoke(precheck)
        else:
            print_warning_msg(f"Skipping installation prechecks")

        for tf_cmd in [tf_init, tf_orc8r, tf_secrets, tf_orc8r_app]:
            cmd = " ".join(tf_cmd)
            if click.confirm(f'Do you want to continue with {cmd}?'):
                rc = run_command(tf_cmd)
                if rc != 0:
                    print_error_msg(f"Install failed when running {cmd} !!!")
                    return

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
                        f'For accessing kubernetes cluster, set'
                        ' `export KUBECONFIG={kubeconfig}`')

                print_success_msg(f"Command {cmd} ran successfully")
            else:
                print_warning_msg(f"Skipping Command {cmd}")


@install.command()
@click.pass_context
def precheck(ctx):
    """
    Performs various checks to ensure successful installation
    """
    rc = run_playbook([
        "ansible-playbook",
        "-v",
        "-e",
        "@/root/config.yml",
        "-t",
        "install_precheck",
        "%s/main.yml" % ctx.obj["playbooks"]])
    if rc != 0:
        print_error_msg("Install prechecks failed!!!")
        sys.exit(1)
    print_success_msg("Install prechecks ran successfully")
