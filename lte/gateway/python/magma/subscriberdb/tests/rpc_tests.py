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

import tempfile
import unittest
from concurrent import futures

import grpc
from lte.protos.subscriberdb_pb2 import SubscriberData, SubscriberUpdate
from lte.protos.subscriberdb_pb2_grpc import SubscriberDBStub
from magma.subscriberdb.rpc_servicer import SubscriberDBRpcServicer
from magma.subscriberdb.sid import SIDUtils
from magma.subscriberdb.store.sqlite import SqliteStore
from orc8r.protos.common_pb2 import Void


class RpcTests(unittest.TestCase):
    """
    Tests for the SubscriberDB rpc servicer and stub
    """

    def setUp(self):
        # Create an in-memory store
        self._tmpfile = tempfile.TemporaryDirectory()
        store = SqliteStore(self._tmpfile.name + '/')

        # Bind the rpc server to a free port
        self._rpc_server = grpc.server(
            futures.ThreadPoolExecutor(max_workers=10),
        )
        port = self._rpc_server.add_insecure_port('0.0.0.0:0')

        # Add the servicer
        self._servicer = SubscriberDBRpcServicer(store)
        self._servicer.add_to_server(self._rpc_server)
        self._rpc_server.start()

        # Create a rpc stub
        channel = grpc.insecure_channel('0.0.0.0:{}'.format(port))
        self._stub = SubscriberDBStub(channel)

    def tearDown(self):
        self._tmpfile.cleanup()
        self._rpc_server.stop(0)

    def test_get_invalid_subscriber(self):
        """
        Test if the rpc call returns NOT_FOUND
        """
        with self.assertRaises(grpc.RpcError) as err:
            self._stub.GetSubscriberData(SIDUtils.to_pb('IMSI123'))
        self.assertEqual(err.exception.code(), grpc.StatusCode.NOT_FOUND)

    def test_add_delete_subscriber(self):
        """
        Test if AddSubscriber and DeleteSubscriber rpc call works
        """
        sid = SIDUtils.to_pb('IMSI1')
        data = SubscriberData(sid=sid)

        # Add subscriber
        self._stub.AddSubscriber(data)

        # Add subscriber again
        with self.assertRaises(grpc.RpcError) as err:
            self._stub.AddSubscriber(data)
        self.assertEqual(err.exception.code(), grpc.StatusCode.ALREADY_EXISTS)

        # See if we can get the data for the subscriber
        self.assertEqual(self._stub.GetSubscriberData(sid), data)
        self.assertEqual(len(self._stub.ListSubscribers(Void()).sids), 1)
        self.assertEqual(self._stub.ListSubscribers(Void()).sids[0], sid)

        # Delete the subscriber
        self._stub.DeleteSubscriber(sid)
        self.assertEqual(len(self._stub.ListSubscribers(Void()).sids), 0)

    def test_update_subscriber(self):
        """
        Test if UpdateSubscriber rpc call works
        """
        sid = SIDUtils.to_pb('IMSI1')
        data = SubscriberData(sid=sid)

        # Add subscriber
        self._stub.AddSubscriber(data)

        sub = self._stub.GetSubscriberData(sid)
        self.assertEqual(sub.lte.auth_key, b'')
        self.assertEqual(sub.state.lte_auth_next_seq, 0)

        # Update subscriber
        update = SubscriberUpdate()
        update.data.sid.CopyFrom(sid)
        update.data.lte.auth_key = b'\xab\xcd'
        update.data.state.lte_auth_next_seq = 1
        update.mask.paths.append('lte.auth_key')  # only auth_key
        self._stub.UpdateSubscriber(update)

        sub = self._stub.GetSubscriberData(sid)
        self.assertEqual(sub.state.lte_auth_next_seq, 0)  # no change
        self.assertEqual(sub.lte.auth_key, b'\xab\xcd')

        update.data.state.lte_auth_next_seq = 1
        update.mask.paths.append('state.lte_auth_next_seq')
        self._stub.UpdateSubscriber(update)

        sub = self._stub.GetSubscriberData(sid)
        self.assertEqual(sub.state.lte_auth_next_seq, 1)

        # Delete the subscriber
        self._stub.DeleteSubscriber(sid)

        with self.assertRaises(grpc.RpcError) as err:
            self._stub.UpdateSubscriber(update)
        self.assertEqual(err.exception.code(), grpc.StatusCode.NOT_FOUND)


if __name__ == "__main__":
    unittest.main()
