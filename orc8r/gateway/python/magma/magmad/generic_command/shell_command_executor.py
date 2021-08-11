"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
import asyncio
import logging
import shlex
from functools import partial
from typing import Any, Dict

from magma.magmad.generic_command.command_executor import (
    CommandExecutor,
    ExecutorFuncT,
    ParamValueT,
)


class ShellCommandExecutor(CommandExecutor):
    """
    The shell command executor stores shell commands from the service config
    into its dispatch table
    """

    def __init__(
            self,
            config: Dict[str, Any],
            loop: asyncio.AbstractEventLoop,
    ) -> None:
        super().__init__(config, loop)
        self._dispatch_table = get_shell_commands_from_config(config)

    def get_command_dispatch(
            self,
    ) -> Dict[str, ExecutorFuncT]:
        return self._dispatch_table


def get_shell_commands_from_config(
        config: Dict[str, Any],
) -> Dict[str, ExecutorFuncT]:
    """
    Gets a list of shell commands from the config. Creates subprocess
    coroutines for each command and stores it in a dictionary.
    """
    shell_commands = config\
        .get('generic_command_config', {})\
        .get('shell_commands', {})

    command_dispatch_table = {}

    for shell_command in shell_commands:
        name = shell_command.get('name')
        command = shell_command.get('command')
        if not name or not command:
            continue

        allow_params = shell_command.get('allow_params', False)

        logging.debug("Loading command %s", name)
        command_dispatch_table[name] = partial(
            _run_subprocess,
            command,
            allow_params,
        )
    return command_dispatch_table


async def _run_subprocess(
        cmd: str,
        allow_params: bool,
        params: Dict[str, ParamValueT],
) -> Dict[str, Any]:
    """
    Runs a command given params (optional), and returns the return code,
    stdout, and stderr.
    """
    cmd_str = cmd
    if allow_params:
        cmd_str = cmd.format(*params.get('shell_params', []))

    logging.info("Running command: %s", cmd_str)

    cmd_list = shlex.split(cmd_str)

    process = await asyncio.create_subprocess_exec(
        *cmd_list,
        stdout=asyncio.subprocess.PIPE,
        stderr=asyncio.subprocess.PIPE,
    )

    stdout, stderr = await process.communicate()

    return {
        "returncode": process.returncode,
        "stdout": stdout.decode(),
        "stderr": stderr.decode(),
    }
