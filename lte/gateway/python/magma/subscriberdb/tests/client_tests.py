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
from google.protobuf.any_pb2 import Any  # pylint: disable=no-name-in-module
from lte.protos.s6a_service_pb2 import DeleteSubscriberRequest
from lte.protos.subscriberdb_pb2 import (
    CheckInSyncRequest,
    CheckInSyncResponse,
    ListSubscribersRequest,
    ListSubscribersResponse,
    LTESubscription,
    SubscriberData,
    SubscriberID,
    SuciProfile,
    SuciProfileList,
    SyncRequest,
    SyncResponse,
)
from lte.protos.subscriberdb_pb2_grpc import (
    SubscriberDBCloudServicer,
    SubscriberDBCloudStub,
    add_SubscriberDBCloudServicer_to_server,
)
from magma.common.grpc_client_manager import GRPCClientManager
from magma.common.service_registry import ServiceRegistry
from magma.subscriberdb.client import SubscriberDBCloudClient
from magma.subscriberdb.sid import SIDUtils
from magma.subscriberdb.store.sqlite import SqliteStore
from orc8r.protos.digest_pb2 import Changeset, Digest, DigestTree, LeafDigest


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

    def CheckInSync(self, request: CheckInSyncRequest, context) -> CheckInSyncResponse:
        """
        Mock to trigger CheckInSync-related test cases

        Args:
            request: CheckInSyncRequest
            context: request context

        Returns:
            CheckInSyncResponse
        """
        in_sync = request.root_digest.md5_base64_digest == "root_digest_apple"
        return CheckInSyncResponse(in_sync=in_sync)

    def Sync(self, request: SyncRequest, context) -> SyncResponse:
        """
        Mock to trigger Sync-related test cases

        Args:
            request: SyncRequest
            context: request context

        Returns:
            SyncResponse
        """
        leaf_digests = [
            LeafDigest(
                id='IMSI11111',
                digest=Digest(md5_base64_digest="digest_apple"),
            ),
            LeafDigest(
                id='IMSI22222',
                digest=Digest(md5_base64_digest="digest_banana"),
            ),
            LeafDigest(
                id='IMSI33333',
                digest=Digest(md5_base64_digest="digest_cherry"),
            ),
        ]

        client_leaf_digest_ids = [
            digest.id for digest in request.leaf_digests
        ]
        to_renew = []
        deleted = []
        if 'IMSI11111' not in client_leaf_digest_ids:
            to_renew.append(subscriber_data_by_id('IMSI11111'))
        if 'IMSI22222' not in client_leaf_digest_ids:
            to_renew.append(subscriber_data_by_id('IMSI22222'))
        if 'IMSI33333' not in client_leaf_digest_ids:
            to_renew.append(subscriber_data_by_id('IMSI33333'))
        if 'IMSI00000' in client_leaf_digest_ids:
            deleted.append('IMSI00000')
        resync = len(to_renew) >= 3

        marshaled_to_renew = []
        for data in to_renew:
            anyVal = Any()
            anyVal.Pack(data)
            marshaled_to_renew.append(anyVal)

        return SyncResponse(
            resync=resync,
            digests=DigestTree(
                root_digest=Digest(md5_base64_digest="root_digest_apple"),
                leaf_digests=leaf_digests,
            ),
            changeset=Changeset(
                to_renew=marshaled_to_renew,
                deleted=deleted,
            ),
        )

    def ListSubscribers(self, request: ListSubscribersRequest, context) -> ListSubscribersResponse:  # noqa: N802
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
        root_digest = Digest(md5_base64_digest="")
        leaf_digests = []
        if request.page_size == 1:
            raise grpc.RpcError("Test Exception")
        if request.page_token == "":
            next_page_token = "aaa"  # noqa: S105
            subscribers = [
                SubscriberData(sid=SubscriberID(id="IMSI111")),
                SubscriberData(sid=SubscriberID(id="IMSI222")),
            ]
            root_digest = Digest(md5_base64_digest="root_digest_apple")
            leaf_digests = [
                LeafDigest(
                    id='IMSI11111',
                    digest=Digest(md5_base64_digest="leaf_digests_apple"),
                ),
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
            digests=DigestTree(
                root_digest=root_digest,
                leaf_digests=leaf_digests,
            ),
        )

    def ListSuciProfiles(self, request, context) -> SuciProfileList:
        suci_profiles = SuciProfile()
        if request.home_network_public_key_identifier < 0 or request.home_network_public_key_identifier > 255:
            raise grpc.RpcError("Invalid home network key identifier value ")
        else:
            return SuciProfileList(
                suci_profiles=suci_profiles,
            )


class SubscriberDBCloudClientTests(unittest.TestCase):
    """Tests for the SubscriberDBCloudClient"""

    def setUp(self):
        """Initialize client tests"""
        # Create sqlite3 database for testing
        self._tmpfile = tempfile.TemporaryDirectory()
        self.loop = asyncio.new_event_loop()
        asyncio.set_event_loop(self.loop)
        suciprofile_db_dict = {},
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
            suciprofile_db_dict=suciprofile_db_dict,
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
            self.assertEqual(self.get_all_subscribers(), ret.subscribers)
            self.assertEqual("root_digest_apple", ret.root_digest.md5_base64_digest)
            self.assertEqual(1, len(ret.leaf_digests))
            self.assertEqual(
                ret.leaf_digests[0].digest.md5_base64_digest,
                "leaf_digests_apple",
            )
            self.assertEqual(ret.leaf_digests[0].id, "IMSI11111")

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
            self.assertTrue(ret is None)

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

    @ unittest.mock.patch(
        'magma.common.service_registry.ServiceRegistry.get_rpc_channel',
    )
    def test_check_subscribers_in_sync(self, get_grpc_mock):
        """
        Test CheckInSync RPC success

        Args:
            get_grpc_mock: mock for service registry
        """
        async def test():  # noqa: WPS430
            get_grpc_mock.return_value = self.channel
            in_sync = (
                await self.subscriberdb_cloud_client._check_subscribers_in_sync()
            )
            self.assertEqual(False, in_sync)

            self.subscriberdb_cloud_client._store.update_root_digest("root_digest_apple")
            in_sync = (
                await self.subscriberdb_cloud_client._check_subscribers_in_sync()
            )
            self.assertEqual(True, in_sync)

        # Cancel the client's loop so there are no other activities
        self.subscriberdb_cloud_client._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @unittest.mock.patch('magma.subscriberdb.client.S6aServiceStub')
    @unittest.mock.patch(
        'magma.common.service_registry.ServiceRegistry.get_rpc_channel',
    )
    def test_sync_subscribers(self, get_grpc_mock, s6a_service_mock_stub):
        """
        Test Sync RPC success

        Args:
            get_grpc_mock: mock for service registry
        """

        async def test():  # noqa: WPS430
            get_grpc_mock.return_value = self.channel

            mock = unittest.mock.Mock()
            mock.DeleteSubscriber.future.side_effect = [unittest.mock.Mock()]
            s6a_service_mock_stub.side_effect = [mock]

            # resync is True if the changeset is too big
            resync = (
                await self.subscriberdb_cloud_client._sync_subscribers()
            )
            self.assertEqual(True, resync)

            self.subscriberdb_cloud_client._store.update_leaf_digests([
                LeafDigest(
                    id='IMSI11111',
                    digest=Digest(md5_base64_digest="digest_apple"),
                ),
                LeafDigest(
                    id='IMSI00000',
                    digest=Digest(md5_base64_digest="digest_zebra"),
                ),
            ])
            self.subscriberdb_cloud_client._store.add_subscriber(
                subscriber_data_by_id('IMSI00000'),
            )
            self.subscriberdb_cloud_client._store.add_subscriber(
                subscriber_data_by_id('IMSI11111'),
            )
            mock.DeleteSubscriber.future.assert_not_called()

            # the client subscriber db and leaf digests db are updated
            # when resync is False
            expected_leaf_digests = [
                LeafDigest(
                    id='IMSI11111',
                    digest=Digest(md5_base64_digest="digest_apple"),
                ),
                LeafDigest(
                    id='IMSI22222',
                    digest=Digest(md5_base64_digest="digest_banana"),
                ),
                LeafDigest(
                    id='IMSI33333',
                    digest=Digest(md5_base64_digest="digest_cherry"),
                ),
            ]
            resync = (
                await self.subscriberdb_cloud_client._sync_subscribers()
            )
            self.assertEqual(False, resync)
            self.assertEqual(
                "root_digest_apple",
                self.subscriberdb_cloud_client._store.get_current_root_digest(),
            )
            self.assertEqual(
                ['IMSI11111', 'IMSI22222', 'IMSI33333'],
                self.subscriberdb_cloud_client._store.list_subscribers(),
            )
            self.assertEqual(
                expected_leaf_digests,
                self.subscriberdb_cloud_client._store.get_current_leaf_digests(),
            )

            # Deleted (IMSI00000) and inactive (IMSI22222) subs are detached
            mock.DeleteSubscriber.future.assert_called_once_with(
                DeleteSubscriberRequest(imsi_list=["00000", "22222"]),
            )

        # Cancel the client's loop so there are no other activities
        self.subscriberdb_cloud_client._periodic_task.cancel()

        self.loop.run_until_complete(test())


def subscriber_data_by_id(sid_str):
    sid = SIDUtils.to_pb(sid_str)
    data = SubscriberData(
        sid=sid, lte=LTESubscription(
            state=LTESubscription.ACTIVE
            if sid_str != "IMSI22222"
            else LTESubscription.INACTIVE,
        ),
    )
    return data
