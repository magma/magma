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

from magma.magmad.generic_command.shell_command_executor import (
    get_shell_commands_from_config,
)


class ShellCommandExecutorTest(unittest.TestCase):
    def setUp(self):
        self.service = mock.Mock()
        asyncio.set_event_loop(asyncio.new_event_loop())

    def test_get_shell_cmds(self):
        self.service.config = {
            'generic_command_config': {
                'shell_commands': [
                    {'name': 'test_1', 'command': '/bin/true'},
                    {'name': 'test_2', 'command': '/bin/false'},
                ],
            },
        }
        table = get_shell_commands_from_config(self.service.config)
        self.assertIn('test_1', table)
        self.assertIn('test_2', table)

    def test_get_shell_cmds_missing_config(self):
        self.service.config = {}
        table = get_shell_commands_from_config(self.service.config)
        self.assertEqual(table, {})

    def test_get_shell_cmds_missing_shell_commands_list(self):
        self.service.config = {
            'generic_command_config': {},
        }
        table = get_shell_commands_from_config(self.service.config)
        self.assertEqual(table, {})

    def test_get_shell_cmds_missing_fields(self):
        self.service.config = {
            'generic_command_config': {
                'shell_commands': [
                    {'name': 'test_1'},
                    {'command': '/bin/false'},
                ],
            },
        }
        table = get_shell_commands_from_config(self.service.config)
        self.assertEqual(table, {})

    def test_get_shell_cmds_funcs(self):
        self.service.config = {
            'generic_command_config': {
                'shell_commands': [
                    {'name': 'test_1', 'command': 'echo'},
                ],
            },
        }

        table = get_shell_commands_from_config(self.service.config)
        func_1 = table['test_1']
        result_1 = asyncio.get_event_loop().run_until_complete(
            func_1(mock.Mock()),
        )

        self.assertEqual(
            result_1,
            {'stdout': '\n', 'stderr': '', 'returncode': 0},
        )

    def test_get_shell_cmds_funcs_with_params(self):
        self.service.config = {
            'generic_command_config': {
                'shell_commands': [
                    {
                        'name': 'test_1',
                        'command': 'echo {} {}',
                        'allow_params': True,
                    },
                ],
            },
        }
        params = {
            "shell_params": ["Hello", "world!"],
        }

        table = get_shell_commands_from_config(self.service.config)
        func_1 = table['test_1']
        result_1 = asyncio.get_event_loop().run_until_complete(
            func_1(params),
        )

        self.assertEqual(
            result_1,
            {'stdout': 'Hello world!\n', 'stderr': '', 'returncode': 0},
        )
        pass

    def test_get_shell_cmds_nonexistent_func(self):
        self.service.config = {
            'generic_command_config': {
                'shell_commands': [
                    {'name': 'test_1', 'command': 'nonexistent/command foo'},
                ],
            },
        }

        table = get_shell_commands_from_config(self.service.config)
        func_1 = table['test_1']
        try:
            asyncio.get_event_loop().run_until_complete(
                func_1(mock.Mock()),
            )
            self.fail('An exception should have been raised')
        except FileNotFoundError:
            pass
