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
import subprocess

import click
from ansible.cli.playbook import PlaybookCLI


def run_playbook(args):
    pb_cli = PlaybookCLI(args)
    pb_cli.parse()
    return pb_cli.run()


def run_command(cmd):
    with subprocess.Popen(cmd, stdout=subprocess.PIPE) as p:
        for output in p.stdout:
            click.echo(output, nl=False)
        return p.wait()
    return 1


def print_error_msg(msg):
    click.echo(click.style(msg, fg='red'))


def print_success_msg(msg):
    click.echo(click.style(msg, fg='green'))


def print_warning_msg(msg):
    click.echo(click.style(msg, fg='yellow'))


def print_info_msg(msg):
    click.echo(click.style(msg, fg='blue'))
