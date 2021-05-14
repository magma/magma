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
from kubernetes import client, config

from .common import print_error_msg, print_success_msg, run_playbook


@click.group()
@click.pass_context
def verify(ctx):
    """
    Run post deployment checks on orc8r
    """
    pass


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
                    repr(kubeconfigs))
            return
        kubeconfig = kubeconfigs[0]

    os.environ["KUBECONFIG"] = kubeconfig
    os.environ["K8S_AUTH_KUBECONFIG"] = kubeconfig

    # check if we have a valid namespace
    config.load_kube_config(kubeconfig)
    v1 = client.CoreV1Api()
    response = v1.list_namespace()
    all_namespaces = [item.metadata.name for item in response.items]
    if namespace not in all_namespaces:
        namespace = click.prompt('Provide orc8r namespace', abort=True)
        if namespace not in all_namespaces:
            print_error_msg(f"Orc8r namespace {namespace} not found")
            sys.exit(1)

    # add constants to the list of variables sent to ansible
    constants['orc8r_namespace'] = namespace

    rc = run_playbook([
        "ansible-playbook",
        "-v",
        "-e",
        json.dumps(constants),
        "-t",
        "verify_sanity",
        "%s/main.yml" % constants["playbooks"]])
    if rc != 0:
        print_error_msg("Post deployment verification checks failed!!!")
        sys.exit(1)
    print_success_msg("Post deployment verification ran successfully")
