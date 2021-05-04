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

from .common import (
    print_error_msg,
    print_success_msg,
    run_command,
    run_playbook,
)


@click.group()
@click.pass_context
def certs(ctx):
    """
    Manage certs in orc8r deployment
    """
    pass


@certs.command()
@click.pass_context
@click.option('--self-signed', is_flag=True)
def add(ctx, self_signed):
    """
    Add creates application and self signed(optional) certs
    """
    if self_signed:
        run_playbook([
            "ansible-playbook",
            "-v",
            "-e",
            "@/root/config.yml",
            "-t",
            "addcerts",
            "%s/certs.yml" % ctx.obj["playbooks"]])
    else:
        run_playbook([
            "ansible-playbook",
            "-v",
            "-e",
            "@/root/config.yml",
            "-t",
            "addcerts",
            "--skip-tags",
            "self_signed",
            "%s/certs.yml" % ctx.obj["playbooks"]])
