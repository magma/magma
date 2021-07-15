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
import json
import os
import sys

import click
from cli.style import print_error_msg, print_success_msg
from utils.ansiblelib import AnsiblePlay, run_playbook
from utils.kubelib import get_all_namespaces, set_kubeconfig_environ


@click.group()
@click.pass_context
def verify(ctx):
    """
    Run post deployment checks on orc8r
    """
    pass


def verify_cmd(constants: dict, namespace: str) -> list:
    """Get the arg list to run prechecks

    Args:
        constants ([dict]): config dict
        namespace ([str]): orc8r namespace
    """
    playbook_dir = constants["playbooks"]
    extra_vars = {
        "orc8r_namespace": namespace,
    }
    return AnsiblePlay(
        playbook=f"{playbook_dir}/main.yml",
        tags=['verify_sanity'],
        extra_vars=extra_vars,
    )


@verify.command('sanity')
@click.option('-n', '--namespace', default='orc8r')
@click.pass_context
def verify_sanity(ctx, namespace):
    # check if KUBECONFIG is set else find kubeconfig file and set the
    # environment variable
    constants = ctx.obj

    # set kubeconfig
    kubeconfig = os.environ.get('KUBECONFIG')
    if not kubeconfig:
        kubeconfigs = glob.glob(constants['project_dir'] + "/kubeconfig_*")
        if len(kubeconfigs) != 1:
            if len(kubeconfigs) == 0:
                print_success_msg('No kubeconfig found!!!')
            else:
                print_error_msg(
                    "multiple kubeconfigs found %s!!!" %
                    repr(kubeconfigs),
                )
            return
        kubeconfig = kubeconfigs[0]
        set_kubeconfig_environ(kubeconfig)

    # check if we have a valid namespace
    all_namespaces = get_all_namespaces(kubeconfig)
    while namespace not in all_namespaces:
        namespace = click.prompt(
            'Provide orc8r namespace',
            type=click.Choice(all_namespaces),
        )

    rc = run_playbook(verify_cmd(ctx.obj, namespace))
    if rc != 0:
        print_error_msg("Post deployment verification checks failed!!!")
        sys.exit(1)
    print_success_msg("Post deployment verification ran successfully")
