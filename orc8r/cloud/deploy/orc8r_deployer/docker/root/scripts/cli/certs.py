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
import click
from utils.ansiblelib import AnsiblePlay, run_playbook


def certs_cmd(constants: dict, self_signed: bool = False) -> list:
    """Provide the arg list to add certs

    Args:
        constants (dict): config constants
        self_signed (bool): add self_signed certs. Defaults to False.

    Returns:
        list: command list
    """
    playbook_dir = constants['playbooks']
    skip_tags = []
    if not self_signed:
        skip_tags = ['self_signed']

    return AnsiblePlay(
        playbook=f"{playbook_dir}/certs.yml",
        tags=['addcerts'],
        skip_tags=skip_tags,
        extra_vars=constants,
    )


@click.group()
@click.pass_context
def certs(ctx):
    """Manage certs in orc8r deployment

    Args:
        ctx ([type]): Click context
    """
    pass


@certs.command()
@click.pass_context
@click.option('--self-signed', is_flag=True)
def add(ctx, self_signed):
    """Add creates application and self signed(optional) certs

    Args:
        ctx ([type]): Click context
        self_signed ([type]): add self_signed certs. Defaults to False.
    """
    run_playbook(certs_cmd(ctx.obj, self_signed))
