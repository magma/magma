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
import glob

import click

from .common import (
    run_playbook,
    print_error_msg,
    print_success_msg)

@click.group()
@click.pass_context
def verify(ctx):
    """
    Run post deployment checks on orc8r
    """
    pass

@verify.command('sanity')
@click.pass_context
def verify_sanity(ctx):
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
                print_error_msg("multiple kubeconfigs found %s!!!" % repr(kubeconfigs))
            return
        kubeconfig = kubeconfigs[0]

    os.environ["KUBECONFIG"] = kubeconfig
    os.environ["K8S_AUTH_KUBECONFIG"] = kubeconfig

    rc = run_playbook([
            "ansible-playbook",
            "-v",
            "-e",
            "@/root/config.yml",
            "-t",
            "verify_sanity",
            "%s/main.yml" % constants["playbooks"]])
    if rc != 0:
        print_error_msg("Post deployment verification checks failed!!!")
        sys.exit(1)
    print_success_msg("Post deployment verification ran successfully")