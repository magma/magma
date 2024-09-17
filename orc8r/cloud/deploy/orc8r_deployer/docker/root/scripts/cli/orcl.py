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

import click
from cli.certs import certs
from cli.cleanup import cleanup
from cli.configure import configure
from cli.install import install
from cli.style import print_error_msg
from cli.upgrade import upgrade
from cli.verify import verify
from utils.common import init


@click.group()
@click.version_option()
@click.pass_context
def cli(ctx):
    """
    Orchestrator Deployment CLI
    """
    ctx.obj = init()


cli.add_command(configure)
cli.add_command(certs)
cli.add_command(install)
cli.add_command(upgrade)
cli.add_command(verify)
cli.add_command(cleanup)
