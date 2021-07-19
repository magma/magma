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
from cli.configlib import ConfigManager


def get_component_choices():
    return ['infra', 'platform', 'service']


@click.group(invoke_without_command=True)
@click.option(
    '-c', '--component',
    type=click.Choice(get_component_choices()),
    multiple=True,
    default=get_component_choices(),
)
@click.pass_context
def configure(ctx, component):
    """
    Configure orc8r deployment variables
    """
    if ctx.invoked_subcommand is None:
        mgr = ConfigManager(ctx.obj)
        for c in component:
            mgr.configure(c)
            mgr.commit(c)


@configure.command()
@click.pass_context
@click.option(
    '-c', '--component',
    type=click.Choice(get_component_choices()),
    multiple=True,
    default=get_component_choices(),
)
def show(ctx, component):
    """
    Display the current configuration
    """
    mgr = ConfigManager(ctx.obj)
    for c in component:
        mgr.show(c)


@configure.command()
@click.option(
    '-c', '--component',
    type=click.Choice(get_component_choices()),
    multiple=True,
    default=get_component_choices(),
)
@click.pass_context
def info(ctx, component):
    """
    Display all possible config options in detail
    """
    mgr = ConfigManager(ctx.obj)
    for c in component:
        mgr.info(c)


@configure.command()
@click.option(
    '-c', '--component',
    type=click.Choice(get_component_choices()),
    multiple=True,
    default=get_component_choices(),
)
@click.pass_context
def check(ctx, component):
    """
    Check if all mandatory configs are present
    """
    mgr = ConfigManager(ctx.obj)
    valid = True
    for c in component:
        valid = mgr.check(c)
    if not valid:
        sys.exit(1)


@configure.command()
@click.option(
    '-c', '--component',
    type=click.Choice(get_component_choices()),
    prompt='select component',
)
@click.option('-k', '--key', prompt='name of the variable')
@click.option('-v', '--value', prompt='value of the variable')
@click.pass_context
def set(ctx, component, key, value):
    """
    Set a specific configuration attribute
    """
    mgr = ConfigManager(ctx.obj)
    mgr.set(component, key, value)
    mgr.commit(component)
