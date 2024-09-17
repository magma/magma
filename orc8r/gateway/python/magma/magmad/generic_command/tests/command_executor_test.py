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
import unittest
from unittest import mock

from magma.magmad.generic_command.command_executor import CommandExecutor


class FakeCommandExecutor(CommandExecutor):
    def __init__(self, _1, _2):
        super().__init__(_1, _2)

        async def test_func_1(_):
            return {
                "success": True,
            }

        async def test_func_2(params):
            return params

        self._dispatch_table = {
            "test_1": test_func_1,
            "test_2": test_func_2,
        }

    def get_command_dispatch(self):
        return self._dispatch_table


class CommandExecutorTest(unittest.TestCase):

    def setUp(self):
        asyncio.set_event_loop(asyncio.new_event_loop())
        self.command_executor = FakeCommandExecutor(mock.Mock(), mock.Mock())

    def test_execute_command(self):
        result = asyncio.get_event_loop().run_until_complete(
            self.command_executor.execute_command("test_1", {}),
        )
        self.assertEqual(result["success"], True)

    def test_execute_command_receives_params(self):
        params = {
            "a": 1,
            "b": "c",
        }
        result = asyncio.get_event_loop().run_until_complete(
            self.command_executor.execute_command("test_2", params),
        )
        self.assertEqual(result, params)
