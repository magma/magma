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
from socket import htons
from unittest.mock import MagicMock

from magma.kernsnoopd.handlers import PacketCounter


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


class PacketCounterTests(unittest.TestCase):
    """Tests for PacketCounter eBPF handler class"""

    def setUp(self) -> None:
        registry = MockServiceRegistry()
        self.packet_counter = PacketCounter(registry)

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
            self.packet_counter._get_source_service(key),
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
            self.packet_counter._get_source_service(key),
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
            ValueError, self.packet_counter._get_source_service,
            key,
        )

    @unittest.mock.patch('magma.kernsnoopd.metrics.MAGMA_BYTES_SENT_TOTAL')
    @unittest.mock.patch('magma.kernsnoopd.metrics.MAGMA_PACKETS_SENT_TOTAL')
    def test_handle_magma_counters(self, packets_count_mock, bytes_count_mock):
        """
        Test handle with Magma service to Magma service traffic
        """
        packets_count_mock.labels = MagicMock(return_value=MagicMock())
        bytes_count_mock.labels = MagicMock(return_value=MagicMock())

        key = MagicMock()
        key.pid, key.comm = 0, b'subscriberdb'
        # 16777343 is "127.0.0.1" packed as a 4 byte int
        key.daddr, key.dport = 16777343, htons(80)

        counter = MagicMock()
        counter.bytes, counter.packets = 100, 10
        bpf = {'dest_counters': {key: counter}}

        self.packet_counter.handle(bpf)

        packets_count_mock.labels.assert_called_once_with(
            service_name='subscriberdb', dest_service='sessiond',
        )
        packets_count_mock.labels.return_value.inc.assert_called_once_with(10)
        bytes_count_mock.labels.assert_called_once_with(
            service_name='subscriberdb', dest_service='sessiond',
        )
        bytes_count_mock.labels.return_value.inc.assert_called_once_with(100)

    @unittest.mock.patch('magma.kernsnoopd.metrics.LINUX_BYTES_SENT_TOTAL')
    @unittest.mock.patch('magma.kernsnoopd.metrics.LINUX_PACKETS_SENT_TOTAL')
    def test_handle_linux_counters(self, packets_count_mock, bytes_count_mock):
        """
        Test handle with Linux binary traffic
        """
        packets_count_mock.labels = MagicMock(return_value=MagicMock())
        bytes_count_mock.labels = MagicMock(return_value=MagicMock())

        key = MagicMock()
        key.pid, key.comm = 0, b'sshd'
        # 16777343 is "127.0.0.1" packed as a 4 byte int
        key.daddr, key.dport = 16777343, htons(443)

        counter = MagicMock()
        counter.bytes, counter.packets = 100, 10
        bpf = {'dest_counters': {key: counter}}

        self.packet_counter.handle(bpf)

        packets_count_mock.labels.assert_called_once_with('sshd')
        packets_count_mock.labels.return_value.inc.assert_called_once_with(10)
        bytes_count_mock.labels.assert_called_once_with('sshd')
        bytes_count_mock.labels.return_value.inc.assert_called_once_with(100)
