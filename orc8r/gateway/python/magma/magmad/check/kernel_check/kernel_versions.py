"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Util module for executing multiple `dpkg` commands via subprocess.
"""

import asyncio
import re
from collections import namedtuple

from magma.magmad.check import subprocess_workflow

DpkgCommandParams = namedtuple('DpkgCommandParams', [])
DpkgCommandResult = namedtuple(
    'DpkgCommandResult',
    ['error', 'kernel_versions_installed'],
)


def get_kernel_versions():
    """
    Execute dpkg commands via subprocess. Blocks while waiting for output.

    Returns:
        [DpkgCommandResult]: stats from the executed dpkg commands
    """
    return subprocess_workflow.exec_and_parse_subprocesses(
        [DpkgCommandParams()],
        _get_dpkg_command_args_list,
        parse_dpkg_output,
    )


@asyncio.coroutine
def get_kernel_versions_async(loop=None):
    """
    Execute dpkg commands asynchronously.

    Args:
        loop: asyncio event loop (optional)

    Returns:
        [DpkgCommandResult]: stats from the executed dpkg commands
    """
    return subprocess_workflow.exec_and_parse_subprocesses_async(
        [DpkgCommandParams()],
        _get_dpkg_command_args_list,
        parse_dpkg_output,
        loop,
    )


def _get_dpkg_command_args_list(_):
    return ['dpkg', '--list']


def parse_dpkg_output(stdout, stderr, _):
    """
    Parse stdout output from a dpkg command.
    """
    if stderr:
        return DpkgCommandResult(
            kernel_versions_installed=None,
            error=str(stderr),
        )
    else:
        installed = re.findall(r'\S*linux-image\S*', str(stdout))
        return DpkgCommandResult(
            kernel_versions_installed=installed,
            error=None,
        )
