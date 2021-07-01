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
from concurrent import futures
from unittest.mock import MagicMock

import grpc
from lte.protos.s6a_service_pb2 import DeleteSubscriberRequest
from lte.protos.subscriberdb_pb2 import (
    Digest,
    ListSubscribersResponse,
    SubscriberData,
    SubscriberID,
)
from lte.protos.subscriberdb_pb2_grpc import (
    SubscriberDBCloudServicer,
    SubscriberDBCloudStub,
    add_SubscriberDBCloudServicer_to_server,
)
from magma.common.grpc_client_manager import GRPCClientManager
from magma.common.service_registry import ServiceRegistry
from magma.subscriberdb.client import SubscriberDBCloudClient
from magma.subscriberdb.store.sqlite import SqliteStore


class MockSubscriberDBServer(SubscriberDBCloudServicer):
    """MockSubscriberDBServer mocks calls from SubscriberDBCloudClient"""

    def __init__(self):
        """Init"""
        pass

    def add_to_server(self, server):
        """
        Add servicer to gRPC server

        Args:
            server: gRPC server
        """
        add_SubscriberDBCloudServicer_to_server(
            self, server,
        )

    def ListSubscribers(self, request, context):  # noqa: N802
        """
        List subscribers is a mock to trigger various test cases

        Args:
            request: ListSubscribersRequest
            context: request context

        Raises:
            RpcError: If page size is 1

        Returns:
            ListSubscribersResponse
        """
        # Add in logic to allow error handling testing
        flat_digest = Digest(md5_base64_digest="digest_pear")
        if request.page_size == 1:
            raise grpc.RpcError("Test Exception")
        if request.page_token == "":
            next_page_token = "aaa"  # noqa: S105
            subscribers = [
                SubscriberData(sid=SubscriberID(id="IMSI111")),
                SubscriberData(sid=SubscriberID(id="IMSI222")),
            ]
        elif request.page_token == "aaa":
            next_page_token = "bbb"  # noqa: S105
            subscribers = [
                SubscriberData(sid=SubscriberID(id="IMSI333")),
                SubscriberData(sid=SubscriberID(id="IMSI444")),
            ]
        else:
            next_page_token = ""  # noqa: S105
            subscribers = [
                SubscriberData(sid=SubscriberID(id="IMSI555")),
                SubscriberData(sid=SubscriberID(id="IMSI666")),
            ]
        return ListSubscribersResponse(
            subscribers=subscribers,
            next_page_token=next_page_token,
            flat_digest=flat_digest,
        )


class SubscriberDBCloudClientTests(unittest.TestCase):
    """Tests for the SubscriberDBCloudClient"""

    def setUp(self):
        """Initialize client tests"""
        # Create sqlite3 database for testing
        self._tmpfile = tempfile.TemporaryDirectory()
        self.loop = asyncio.new_event_loop()
        asyncio.set_event_loop(self.loop)
        store = SqliteStore(
            '{filename}{slash}'.format(
                filename=self._tmpfile.name,
                slash='/',
            ),
        )

        ServiceRegistry.add_service('test', '0.0.0.0', 0)  # noqa: S104
        ServiceRegistry._PROXY_CONFIG = {
            'local_port': 1234,
            'cloud_address': '',
            'proxy_cloud_connections': False,
        }
        ServiceRegistry._REGISTRY = {
            "services": {
                "s6a_service":
                {
                    "ip_address": "0.0.0.0",  # noqa: S104
                    "port": 2345,
                },
            },
        }

        self.service = MagicMock()
        self.service.loop = self.loop

        # Bind the rpc server to a free port
        self._rpc_server = grpc.server(
            futures.ThreadPoolExecutor(max_workers=10),
        )
        port = self._rpc_server.add_insecure_port('0.0.0.0:0')
        # Add the servicer
        self._servicer = MockSubscriberDBServer()
        self._servicer.add_to_server(self._rpc_server)
        self._rpc_server.start()
        # Create a rpc stub
        self.channel = grpc.insecure_channel(
            '0.0.0.0:{port}'.format(
                port=port,
            ),
        )
        grpc_client_manager = GRPCClientManager(
            service_name="subscriberdb",
            service_stub=SubscriberDBCloudStub,
            max_client_reuse=60,
        )
        self.subscriberdb_cloud_client = SubscriberDBCloudClient(
            loop=self.service.loop,
            store=store,
            subscriber_page_size=2,
            sync_interval=10,
            grpc_client_manager=grpc_client_manager,
        )
        self.subscriberdb_cloud_client.start()

    def tearDown(self):
        """Clean up test setup"""
        self._tmpfile.cleanup()
        self._rpc_server.stop(None)
        self.subscriberdb_cloud_client.stop()

    def get_all_subscribers(self):
        return [
            SubscriberData(
                sid=SubscriberID(
                    id="IMSI111",
                ),
            ), SubscriberData(
                sid=SubscriberID(
                    id="IMSI222",
                ),
            ), SubscriberData(
                sid=SubscriberID(
                    id="IMSI333",
                ),
            ), SubscriberData(
                sid=SubscriberID(
                    id="IMSI444",
                ),
            ), SubscriberData(
                sid=SubscriberID(
                    id="IMSI555",
                ),
            ), SubscriberData(
                sid=SubscriberID(
                    id="IMSI666",
                ),
            ),
        ]

    @ unittest.mock.patch(
        'magma.common.service_registry.ServiceRegistry.get_rpc_channel',
    )
    def test_get_all_subscribers_success(self, get_grpc_mock):
        """
        Test ListSubscribers RPC happy path

        Args:
            get_grpc_mock: mock for service registry

        """
        async def test():  # noqa: WPS430
            get_grpc_mock.return_value = self.channel
            ret = (
                await self.subscriberdb_cloud_client._get_all_subscribers()
            )
            self.assertTrue(ret is not None)
            self.assertEqual(self.get_all_subscribers(), ret[0])
            self.assertEqual("digest_pear", ret[1].md5_base64_digest)

        # Cancel the client's loop so there are no other activities
        self.subscriberdb_cloud_client._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @ unittest.mock.patch(
        'magma.common.service_registry.ServiceRegistry.get_rpc_channel',
    )
    def test_get_all_subscribers_failure(self, get_grpc_mock):
        """
        Test ListSubscribers RPC failures

        Args:
            get_grpc_mock: mock for service registry

        """
        async def test():  # noqa: WPS430
            get_grpc_mock.return_value = self.channel
            # update page size to special value to trigger gRPC error
            self.subscriberdb_cloud_client._subscriber_page_size = 1
            ret = (
                await self.subscriberdb_cloud_client._get_all_subscribers()
            )
            self.assertEqual((None, None), ret)

        # Cancel the client's loop so there are no other activities
        self.subscriberdb_cloud_client._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @ unittest.mock.patch('magma.subscriberdb.client.S6aServiceStub')
    def test_detach_deleted_subscribers(self, s6a_service_mock_stub):
        """
        Test if the subscriberdb cloud client detaches deleted subscribers

        Args:
            s6a_service_mock_stub: mocked s6a stub
        """
        # Mock out DeleteSubscriber.future
        mock = unittest.mock.Mock()
        mock.DeleteSubscriber.future.side_effect = [unittest.mock.Mock()]
        s6a_service_mock_stub.side_effect = [mock]

        # Call with no samples
        old_sub_ids = ["IMSI202", "IMSI101"]
        new_sub_ids = ["IMSI101", "IMSI202"]
        self.subscriberdb_cloud_client._detach_deleted_subscribers(
            old_sub_ids,
            new_sub_ids,
        )
        s6a_service_mock_stub.DeleteSubscriber.future.assert_not_called()
        self.subscriberdb_cloud_client._loop.stop()

        # Call with one subscriber id deleted
        old_sub_ids = ["IMSI202", "IMSI101", "IMSI303"]
        new_sub_ids = ["IMSI202"]
        self.subscriberdb_cloud_client._detach_deleted_subscribers(
            old_sub_ids,
            new_sub_ids,
        )

        mock.DeleteSubscriber.future.assert_called_once_with(
            DeleteSubscriberRequest(imsi_list=["101", "303"]),
        )
