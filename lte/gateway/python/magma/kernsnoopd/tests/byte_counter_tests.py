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
import unittest
from socket import AF_INET, AF_INET6, htons
from unittest.mock import MagicMock

from magma.kernsnoopd.handlers import ByteCounter


class MockServiceRegistry:
    @staticmethod
    def list_services():
        return ['sessiond', 'subscriberdb']

    @staticmethod
    def get_service_name(host, port):
        service_map = {
            ('127.0.0.1', 80): 'sessiond',
            ('127.0.0.1', 443): 'subscriberdb',
        }
        return service_map[(host, port)]


class ByteCounterTests(unittest.TestCase):
    # pylint: disable=protected-access
    """Tests for ByteCounter eBPF handler class"""

    def setUp(self) -> None:
        registry = MockServiceRegistry()
        self.byte_counter = ByteCounter(registry)

    @unittest.mock.patch('psutil.Process')
    def test_get_source_service_python(self, process_mock):
        """
        Test _get_source_service Python service happy path
        """
        cmdline_value = 'python3 -m magma.subscriberdb.main'.split(' ')
        process_mock.return_value.cmdline.return_value = cmdline_value
        key = MagicMock()
        self.assertEqual(
            'subscriberdb',
            self.byte_counter._get_source_service(key),
        )

    @unittest.mock.patch('psutil.Process')
    def test_get_source_service_index_error_with_fallback(self, process_mock):
        """
        Test _get_source_service sessiond in service list
        """
        cmdline_value = 'sessiond'.split(' ')
        process_mock.return_value.cmdline.return_value = cmdline_value
        key = MagicMock()
        key.comm = b'sessiond'
        self.assertEqual(
            'sessiond',
            self.byte_counter._get_source_service(key),
        )

    @unittest.mock.patch('psutil.Process')
    def test_get_source_service_index_error_fail(self, process_mock):
        """
        Test _get_source_service sshd not in service list
        """
        cmdline_value = 'sshd'.split(' ')
        process_mock.return_value.cmdline.return_value = cmdline_value
        key = MagicMock()
        key.comm = b'sshd'
        self.assertRaises(
            ValueError, self.byte_counter._get_source_service,
            key,
        )

    def test_check_cmdline_for_magma_service_python_service_found(self):
        """
        Test _check_cmdline_for_magma_service python service name found in cmdline
        """
        cmdline_value = 'python3 -m magma.subscriberdb.main'.split(' ')
        self.assertEqual(
            'subscriberdb',
            self.byte_counter._get_service_from_cmdline(cmdline_value),
        )

    def test_check_cmdline_for_magma_service_native_service_found(self):
        """
        Test _check_cmdline_for_magma_service native service name found in cmdline
        """
        cmdline_value = 'foo bar magma.sessiond'.split(' ')
        self.assertEqual(
            'sessiond',
            self.byte_counter._get_service_from_cmdline(cmdline_value),
        )

    def test_check_cmdline_for_magma_service_index_error(self):
        """
        Test _check_cmdline_for_magma_service index error
        """
        cmdline_value = 'sshd'.split(' ')
        self.assertRaises(
            IndexError, self.byte_counter._get_service_from_cmdline,
            cmdline_value,
        )

    def test_check_cmdline_for_magma_service_magma_not_found(self):
        """
        Test _check_cmdline_for_magma_service no magma. in cmdline
        """
        cmdline_value = 'python3 -m subscriberdb.main'.split(' ')
        self.assertIsNone(self.byte_counter._get_service_from_cmdline(cmdline_value))

    @unittest.mock.patch('magma.kernsnoopd.metrics.MAGMA_BYTES_SENT_TOTAL')
    @unittest.mock.patch('psutil.Process')
    def test_handle_magma_counters(self, process_mock, bytes_count_mock):
        """
        Test handle with Magma service to Magma service traffic
        """
        bytes_count_mock.labels = MagicMock(return_value=MagicMock())
        cmdline_value = 'python3 -m magma.subscriberdb.main'.split(' ')
        process_mock.return_value.cmdline.return_value = cmdline_value

        key = MagicMock()
        key.family = AF_INET
        # 16777343 is "127.0.0.1" packed as a 4 byte int
        key.daddr = self.byte_counter.Addr(16777343, 0)
        key.dport = htons(80)

        count = MagicMock()
        count.value = 100
        bpf = {'dest_counters': {key: count}}

        self.byte_counter.handle(bpf)

        bytes_count_mock.labels.assert_called_once_with(
            service_name='subscriberdb', dest_service='',
        )
        bytes_count_mock.labels.return_value.inc.assert_called_once_with(100)

    @unittest.mock.patch('magma.kernsnoopd.metrics.LINUX_BYTES_SENT_TOTAL')
    @unittest.mock.patch('psutil.Process')
    def test_handle_linux_counters(self, process_mock, bytes_count_mock):
        """
        Test handle with Linux binary traffic
        """
        bytes_count_mock.labels = MagicMock(return_value=MagicMock())
        cmdline_value = 'sshd'.split(' ')
        process_mock.return_value.cmdline.return_value = cmdline_value

        key = MagicMock()
        key.comm = b'sshd'
        key.family = AF_INET6
        # localhost in IPv6 with embedded IPv4
        # ::ffff:127.0.0.1 = 0x0100007FFFFF0000
        key.daddr = self.byte_counter.Addr(0, 0x0100007FFFFF0000)

        count = MagicMock()
        count.value = 100
        bpf = {'dest_counters': {key: count}}

        self.byte_counter.handle(bpf)

        bytes_count_mock.labels.assert_called_once_with('sshd')
        bytes_count_mock.labels.return_value.inc.assert_called_once_with(100)
