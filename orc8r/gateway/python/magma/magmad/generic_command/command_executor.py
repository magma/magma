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
import importlib
from abc import ABC, abstractmethod
from typing import Any, Awaitable, Callable, Dict, List, Union

ParamValueT = Union[str, int, float, bool, List[Union[str, int, float, bool]]]
ExecutorFuncT = Callable[[Dict[str, ParamValueT]], Awaitable[Dict[str, Any]]]


class CommandExecutor(ABC):
    """
    Abstract class for command executors
    """

    def __init__(
            self,
            config: Dict[str, Any],
            loop: asyncio.AbstractEventLoop,
    ) -> None:
        self._config = config
        self._loop = loop

    async def execute_command(
            self,
            command: str,
            params: Dict[str, ParamValueT],
    ) -> Dict[str, Any]:
        """
        Run the command from the dispatch table with params
        """
        result = await self.get_command_dispatch()[command](params)
        return result

    @abstractmethod
    def get_command_dispatch(self) -> Dict[str, ExecutorFuncT]:
        """
        Returns the command dispatch table for this command executor
        """
        pass


def get_command_executor_impl(service):
    """
    Gets the command executor impl from the service config
    """
    config = service.config.get('generic_command_config', None)
    assert config is not None, 'generic_command_config not found'
    module = config.get('module', None)
    impl_class = config.get('class', None)
    assert module is not None, 'generic command module not found'
    assert impl_class is not None, 'generic command class not found'
    command_executor_class = getattr(
        importlib.import_module(module),
        impl_class,
    )
    command_executor = command_executor_class(service.config, service.loop)
    assert isinstance(command_executor, CommandExecutor), \
        'command_executor is not an instance of CommandExecutor'
    return command_executor
