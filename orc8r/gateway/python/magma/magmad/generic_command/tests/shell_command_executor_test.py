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
    COMMAND_ALLOWLIST,
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
                    {'name': 'health', 'command': '/bin/true'},
                    {'name': 'get_flows', 'command': '/bin/false'},
                ],
            },
        }
        table = get_shell_commands_from_config(self.service.config)
        self.assertIn('health', table)
        self.assertIn('get_flows', table)

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
                    {'name': 'health'},
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
                    {'name': 'health', 'command': 'echo'},
                ],
            },
        }

        table = get_shell_commands_from_config(self.service.config)
        func_1 = table['health']
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
                        'name': 'health',
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
        func_1 = table['health']
        result_1 = asyncio.get_event_loop().run_until_complete(
            func_1(params),
        )

        self.assertEqual(
            result_1,
            {'stdout': 'Hello world!\n', 'stderr': '', 'returncode': 0},
        )

    def test_get_shell_cmds_nonexistent_func(self):
        self.service.config = {
            'generic_command_config': {
                'shell_commands': [
                    {'name': 'health', 'command': 'nonexistent/command foo'},
                ],
            },
        }

        table = get_shell_commands_from_config(self.service.config)
        func_1 = table['health']
        try:
            asyncio.get_event_loop().run_until_complete(
                func_1(mock.Mock()),
            )
            self.fail('An exception should have been raised')
        except FileNotFoundError:
            pass

    def test_allowlist_rejects_bash(self):
        """Commands not in COMMAND_ALLOWLIST should be silently skipped."""
        self.service.config = {
            'generic_command_config': {
                'shell_commands': [
                    {'name': 'bash', 'command': 'bash {}', 'allow_params': True},
                    {'name': 'fab', 'command': 'fab {}', 'allow_params': True},
                    {'name': 'echo', 'command': 'echo {}', 'allow_params': True},
                ],
            },
        }
        table = get_shell_commands_from_config(self.service.config)
        self.assertEqual(table, {})

    def test_allowlist_permits_domain_commands(self):
        """All COMMAND_ALLOWLIST entries should be registerable."""
        commands = [
            {'name': name, 'command': f'/usr/bin/{name}'}
            for name in COMMAND_ALLOWLIST
        ]
        self.service.config = {
            'generic_command_config': {
                'shell_commands': commands,
            },
        }
        table = get_shell_commands_from_config(self.service.config)
        self.assertEqual(set(table.keys()), COMMAND_ALLOWLIST)

    def test_metacharacter_rejection(self):
        """Parameters with shell metacharacters should raise ValueError."""
        self.service.config = {
            'generic_command_config': {
                'shell_commands': [
                    {
                        'name': 'health',
                        'command': 'health_cli.py {}',
                        'allow_params': True,
                    },
                ],
            },
        }
        table = get_shell_commands_from_config(self.service.config)
        func = table['health']

        for bad_param in ['; rm -rf /', '| cat /etc/passwd', '& bg', '$(id)', '`id`', 'foo\nbar']:
            params = {"shell_params": [bad_param]}
            with self.assertRaises(ValueError, msg=f"Should reject: {bad_param!r}"):
                asyncio.get_event_loop().run_until_complete(func(params))
