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


def print_title(msg):
    click.echo(click.style(msg, underline=True))


def print_error_msg(msg):
    click.echo(click.style(msg, fg='red'))


def print_success_msg(msg):
    click.echo(click.style(msg, fg='green'))


def print_warning_msg(msg):
    click.echo(click.style(msg, fg='yellow'))


def print_info_msg(msg):
    click.echo(click.style(msg, fg='blue'))
