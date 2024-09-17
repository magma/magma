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
import subprocess
import unittest
from unittest import mock

from magma.magmad.check import subprocess_workflow


class SubprocessWorkflowTests(unittest.TestCase):

    def setUp(self):
        self.mock_params = ['param1', 'param2']
        self.mock_parser_callback = mock.Mock()
        self.mock_parser_callback.return_value = 42

        self.process_communicate_mock = mock.Mock(
            wraps=self._mock_process_communicate,
        )
        self.process_mock = mock.Mock()
        self.process_mock.configure_mock(
            communicate=self.process_communicate_mock,
        )

        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)

    @staticmethod
    def _mock_arg_factory(param):
        return [param, 'a', 'b']

    @asyncio.coroutine
    # pylint: disable=unused-argument
    def _mock_process_communicate(self, *args, **kwargs):
        return 'stdout', 'stderr'

    @asyncio.coroutine
    # pylint: disable=unused-argument
    def _mock_async_subprocess(self, *args, **kwargs):
        return self.process_mock

    @mock.patch('subprocess.Popen')
    def test_subprocess_exec(self, mock_popen):
        process_mock = mock.Mock()
        mock_config = {'communicate.return_value': ('stdout', 'stderr')}
        process_mock.configure_mock(**mock_config)
        mock_popen.return_value = process_mock

        actual = subprocess_workflow.exec_and_parse_subprocesses(
            self.mock_params,
            self._mock_arg_factory,
            self.mock_parser_callback,
        )
        self.assertEqual([42, 42], list(actual))

        mock_popen.assert_has_calls([
            mock.call(
                ['param1', 'a', 'b'],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
            ),
            mock.call(
                ['param2', 'a', 'b'],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
            ),
        ])
        self.mock_parser_callback.assert_has_calls([
            mock.call('stdout', 'stderr', 'param1'),
            mock.call('stdout', 'stderr', 'param2'),
        ])

    def test_async_subprocess_exec(self):
        with mock.patch.object(
                asyncio, 'create_subprocess_exec',
                mock.Mock(wraps=self._mock_async_subprocess),
        ) as m:
            loop = asyncio.get_event_loop()
            actual = loop.run_until_complete(
                subprocess_workflow.exec_and_parse_subprocesses_async(
                    self.mock_params,
                    self._mock_arg_factory,
                    self.mock_parser_callback,
                ),
            )

            self.assertEqual([42, 42], list(actual))
            m.assert_has_calls([
                mock.call(
                    'param1', 'a', 'b',
                    stdout=subprocess.PIPE,
                    stderr=subprocess.PIPE,
                ),
                mock.call(
                    'param2', 'a', 'b',
                    stdout=subprocess.PIPE,
                    stderr=subprocess.PIPE,
                ),
            ])
            self.mock_parser_callback.assert_has_calls([
                mock.call('stdout', 'stderr', 'param1'),
                mock.call('stdout', 'stderr', 'param2'),
            ])

    def test_async_subprocess_exec_exception(self):
        @asyncio.coroutine
        # pylint: disable=unused-argument
        def mock_subprocess_raises(*args, **kwargs):
            raise ValueError('oops')

        with mock.patch.object(
            asyncio, 'create_subprocess_exec',
            mock.Mock(wraps=mock_subprocess_raises),
        ):
            loop = asyncio.get_event_loop()
            with self.assertRaises(ValueError) as e:
                loop.run_until_complete(
                    subprocess_workflow.exec_and_parse_subprocesses_async(
                        self.mock_params,
                        self._mock_arg_factory,
                        self.mock_parser_callback,
                    ),
                )
            self.assertEqual('oops', e.exception.args[0])

    def test_async_subprocess_communicate_exception(self):
        @asyncio.coroutine
        # pylint: disable=unused-argument
        def mock_communicate_raises(*args, **kwargs):
            raise ValueError('oops')
        self.process_mock.configure_mock(communicate=mock_communicate_raises)

        with mock.patch.object(
                asyncio,
                'create_subprocess_exec',
                mock.Mock(wraps=self._mock_async_subprocess),
        ) as m:
            loop = asyncio.get_event_loop()
            with self.assertRaises(ValueError) as e:
                loop.run_until_complete(
                    subprocess_workflow.exec_and_parse_subprocesses_async(
                        self.mock_params,
                        self._mock_arg_factory,
                        self.mock_parser_callback,
                    ),
                )
            self.assertEqual('oops', e.exception.args[0])
            m.assert_has_calls([
                mock.call(
                    'param1', 'a', 'b',
                    stdout=subprocess.PIPE,
                    stderr=subprocess.PIPE,
                ),
                mock.call(
                    'param2', 'a', 'b',
                    stdout=subprocess.PIPE,
                    stderr=subprocess.PIPE,
                ),
            ])

    def test_async_subprocess_empty_params(self):
        with mock.patch.object(
                asyncio, 'create_subprocess_exec',
                mock.Mock(wraps=self._mock_async_subprocess),
        ) as m:
            loop = asyncio.get_event_loop()
            actual = loop.run_until_complete(
                subprocess_workflow.exec_and_parse_subprocesses_async(
                    [],
                    self._mock_arg_factory,
                    self.mock_parser_callback,
                ),
            )

            self.assertEqual([], list(actual))
            m.assert_not_called()
            self.mock_parser_callback.assert_not_called()
