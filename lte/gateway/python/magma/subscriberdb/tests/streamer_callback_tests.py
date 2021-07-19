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

# pylint: disable=protected-access

import asyncio
import tempfile
import unittest
import unittest.mock

from lte.protos.s6a_service_pb2 import DeleteSubscriberRequest
from magma.common.service_registry import ServiceRegistry
from magma.subscriberdb.store.sqlite import SqliteStore
from magma.subscriberdb.streamer_callback import SubscriberDBStreamerCallback


class MockFuture(object):
    is_error = True

    def __init__(self, is_error):
        self.is_error = is_error

    def exception(self):
        if self.is_error:
            return self.MockException()

    class MockException(object):
        def details(self):
            return ''

        def code(self):
            return 0


class SubscriberDBStreamerCallbackTests(unittest.TestCase):
    """
    Tests for the SubscriberDBStreamerCallback detach_deleted_subscribers
    """

    def setUp(self):
        # Create sqlite3 database for testing
        self._tmpfile = tempfile.TemporaryDirectory()
        store = SqliteStore(self._tmpfile.name + '/')
        self._streamer_callback = \
            SubscriberDBStreamerCallback(store, loop=asyncio.new_event_loop())
        ServiceRegistry.add_service('test', '0.0.0.0', 0)
        ServiceRegistry._PROXY_CONFIG = {
            'local_port': 1234,
            'cloud_address': '',
            'proxy_cloud_connections': False,
        }
        ServiceRegistry._REGISTRY = {
            "services": {
                "s6a_service":
                {
                    "ip_address": "0.0.0.0",
                    "port": 2345,
                },
            },
        }

    def tearDown(self):
        self._tmpfile.cleanup()

    @unittest.mock.patch('magma.subscriberdb.streamer_callback.S6aServiceStub')
    def test_detach_deleted_subscribers(self, s6a_service_mock_stub):
        """
        Test if the streamer_callback detach deleted subscribers.
        """
        # Mock out DeleteSubscriber.future
        mock = unittest.mock.Mock()
        mock.DeleteSubscriber.future.side_effect = [unittest.mock.Mock()]
        s6a_service_mock_stub.side_effect = [mock]

        # Call with no samples
        old_sub_ids = ["IMSI202", "IMSI101"]
        new_sub_ids = ["IMSI101", "IMSI202"]
        self._streamer_callback.detach_deleted_subscribers(
            old_sub_ids,
            new_sub_ids,
        )
        s6a_service_mock_stub.DeleteSubscriber.future.assert_not_called()
        self._streamer_callback._loop.stop()

        # Call with one subscriber id deleted
        old_sub_ids = ["IMSI202", "IMSI101", "IMSI303"]
        new_sub_ids = ["IMSI202"]
        self._streamer_callback.detach_deleted_subscribers(
            old_sub_ids,
            new_sub_ids,
        )

        mock.DeleteSubscriber.future.assert_called_once_with(
            DeleteSubscriberRequest(imsi_list=["101", "303"]),
        )


if __name__ == "__main__":
    unittest.main()
