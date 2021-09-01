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

import re
from typing import NamedTuple, Optional

from magma.magmad.check import subprocess_workflow

LscpuCommandParams = NamedTuple('LscpuCommandParams', [])
LscpuCommandResult = NamedTuple(
    'LscpuCommandResult',
    [
        ('error', Optional[str]),
        ('core_count', Optional[int]),
        ('threads_per_core', Optional[int]),
        ('architecture', Optional[str]),
        ('model_name', Optional[str]),
    ],
)


def get_cpu_info() -> LscpuCommandResult:
    """
    Execute lscpu command via subprocess. Blocks while waiting for output.
    """
    return list(
        subprocess_workflow.exec_and_parse_subprocesses(
            [LscpuCommandParams()],
            _get_lscpu_command_args_list,
            parse_lscpu_output,
        ),
    )[0]


def _get_lscpu_command_args_list(_):
    return ['lscpu']


def parse_lscpu_output(stdout, stderr, _):
    """
    Parse stdout output from a lscpu command.
    """

    def _create_error_result(err_msg):
        return LscpuCommandResult(
            error=err_msg, core_count=None,
            threads_per_core=None, architecture=None,
            model_name=None,
        )
    if stderr:
        return _create_error_result(stderr)

    stdout_decoded = stdout.decode()
    try:
        cores_per_socket = int(
            re.search(
                r'Core\(s\) per socket:\s*(.*)\n',
                str(stdout_decoded),
            ).group(1),
        )
        num_sockets = int(
            re.search(
                r'Socket\(s\):\s*(.*)\n',
                str(stdout_decoded),
            ).group(1),
        )
        threads_per_core = int(
            re.search(
                r'Thread\(s\) per core:\s*(.*)\n',
                str(stdout_decoded),
            ).group(1),
        )
        architecture = re.search(
            r'Architecture:\s*(.*)\n',
            str(stdout_decoded),
        ).group(1)
        model_name = re.search(
            r'Model name:\s*(.*)\n',
            str(stdout_decoded),
        ).group(1)
        return LscpuCommandResult(
            error=None,
            core_count=cores_per_socket * num_sockets,
            threads_per_core=threads_per_core,
            architecture=architecture,
            model_name=model_name,
        )
    except (AttributeError, IndexError, ValueError) as e:
        return _create_error_result(
            'Parsing failed: %s\n%s' % (e, stdout_decoded),
        )
