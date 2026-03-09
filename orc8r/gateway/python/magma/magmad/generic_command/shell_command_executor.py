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
import re
import shlex
from functools import partial
from typing import Any, Dict

from magma.magmad.generic_command.command_executor import (
    CommandExecutor,
    ExecutorFuncT,
    ParamValueT,
)

# Only these command names may be registered via config.
# Generic shell access (bash, fab, echo) is intentionally excluded
# to prevent arbitrary command execution via the Orchestrator API.
COMMAND_ALLOWLIST = frozenset({
    'reboot_enodeb',
    'reboot_all_enodeb',
    'health',
    'agw_health',
    'get_flows',
    'get_subscriber_table',
    'check_stateless',
    'configure_stateless',
})

# Shell metacharacters that must not appear in command parameters.
_SHELL_METACHAR_RE = re.compile(r'[;|&$`\n]')


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
    Only commands whose name is in COMMAND_ALLOWLIST are registered.
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

        if name not in COMMAND_ALLOWLIST:
            logging.warning(
                "Skipping command '%s': not in COMMAND_ALLOWLIST", name,
            )
            continue

        allow_params = shell_command.get('allow_params', False)

        logging.debug("Loading command %s", name)
        command_dispatch_table[name] = partial(
            _run_subprocess,
            command,
            allow_params,
        )
    return command_dispatch_table


def _validate_shell_params(params: list) -> None:
    """Reject parameters containing shell metacharacters."""
    for param in params:
        if isinstance(param, str) and _SHELL_METACHAR_RE.search(param):
            raise ValueError(
                "shell_params contains forbidden metacharacter: "
                f"{param!r}",
            )


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
        shell_params = params.get('shell_params', [''])
        if shell_params == []:
            shell_params = ['']
        _validate_shell_params(shell_params)
        cmd_str = cmd.format(*shell_params)

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
