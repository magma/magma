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
    """Tests for ByteCounter eBPF handler class"""

    def setUp(self) -> None:
        registry = MockServiceRegistry()
        self.byte_counter = ByteCounter(registry)

    @unittest.mock.patch('psutil.Process.cmdline')
    def test_get_source_service_python(self, cmdline_mock):
        """
        Test _get_source_service Python service happy path
        """
        cmdline_mock.return_value = 'python3 -m magma.subscriberdb.main'.split(
            ' ',
        )
        key = MagicMock()
        key.pid = 0
        self.assertEqual(
            'subscriberdb',
            self.byte_counter._get_source_service(key),
        )

    @unittest.mock.patch('psutil.Process.cmdline')
    def test_get_source_service_native(self, cmdline_mock):
        """
        Test _get_source_service native service happy path
        """
        cmdline_mock.return_value = 'sessiond'.split(' ')
        key = MagicMock()
        key.pid, key.comm = 0, b'sessiond'
        self.assertEqual(
            'sessiond',
            self.byte_counter._get_source_service(key),
        )

    @unittest.mock.patch('psutil.Process.cmdline')
    def test_get_source_service_fail(self, cmdline_mock):
        """
        Test _get_source_service failure
        """
        cmdline_mock.return_value = 'sshd'.split(' ')
        key = MagicMock()
        key.pid, key.comm = 0, b'sshd'
        self.assertRaises(
            ValueError, self.byte_counter._get_source_service,
            key,
        )

    @unittest.mock.patch('magma.kernsnoopd.metrics.MAGMA_BYTES_SENT_TOTAL')
    def test_handle_magma_counters(self, bytes_count_mock):
        """
        Test handle with Magma service to Magma service traffic
        """
        bytes_count_mock.labels = MagicMock(return_value=MagicMock())

        key = MagicMock()
        key.pid, key.comm = 0, b'subscriberdb'
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
    def test_handle_linux_counters(self, bytes_count_mock):
        """
        Test handle with Linux binary traffic
        """
        bytes_count_mock.labels = MagicMock(return_value=MagicMock())

        key = MagicMock()
        key.pid, key.comm = 0, b'sshd'
        key.family = AF_INET6
        # localhost in IPv6 with embedded IPv4
        # ::ffff:127.0.0.1 = 0x0100007FFFFF0000
        key.daddr = self.byte_counter.Addr(0, 0x0100007FFFFF0000)
        key.dport = htons(443)

        count = MagicMock()
        count.value = 100
        bpf = {'dest_counters': {key: count}}

        self.byte_counter.handle(bpf)

        bytes_count_mock.labels.assert_called_once_with('sshd')
        bytes_count_mock.labels.return_value.inc.assert_called_once_with(100)
