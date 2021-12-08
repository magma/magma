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
from lte.protos.subscriberdb_pb2 import (
    SubscriberData,
    SubscriberUpdate,
    SuciProfile,
)
from lte.protos.subscriberdb_pb2_grpc import SubscriberDBStub, SuciProfileDBStub
from magma.subscriberdb.rpc_servicer import (
    SubscriberDBRpcServicer,
    SuciProfileDBRpcServicer,
)
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


class RpcTestsSuciExt(unittest.TestCase):
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

        suciprofile_db = {}
        # Add the servicer
        self._servicer = SuciProfileDBRpcServicer(
            store,
            suciprofile_db,
        )

        self._servicer.add_to_server(self._rpc_server)
        self._rpc_server.start()

        # Create a rpc stub
        channel = grpc.insecure_channel('0.0.0.0:{}'.format(port))
        self._stub = SuciProfileDBStub(channel)

    def tearDown(self):
        self._tmpfile.cleanup()
        self._rpc_server.stop(0)

    def test_add_delete_suciprofile(self):
        """
        Test if AddSuciProfile and DeleteSuciProfile rpc call works
        """
        home_net_public_key_id = 2
        protection_scheme = SuciProfile.ProfileA
        home_net_public_key = bytes(b'\t\xd4(\x93O7\x15\x13\xa4\x1c\xf6\xef\x96\x01bwH\xf3wO\x1ds\x99\xd7\x8d{dc\x0c\x94Q\x08')
        home_net_private_key = bytes(b'\xc8\x0b\x16v\xa0\xc9\x83PI\x0f\xf1\xc3\x13\x08-\xedE\xbcY\xe9\xe7)\xb8x\x8e\xba\xc44\xb3\x8a\xb5m')

        request = SuciProfile(
            home_net_public_key_id=home_net_public_key_id,
            protection_scheme=protection_scheme,
            home_net_public_key=home_net_public_key,
            home_net_private_key=home_net_private_key,
        )
        # Add subscriber
        self._stub.AddSuciProfile(request)

        # Add subscriber again
        with self.assertRaises(grpc.RpcError) as err:
            self._stub.AddSuciProfile(request)
        self.assertEqual(err.exception.code(), grpc.StatusCode.ALREADY_EXISTS)

        # See if we can get the data for the subscriber
        self.assertEqual(len(self._stub.ListSuciProfile(Void()).suci_profiles), 1)
        self.assertEqual(self._stub.ListSuciProfile(Void()).suci_profiles[0].home_net_public_key_id, 2)
        self.assertEqual(self._stub.ListSuciProfile(Void()).suci_profiles[0].protection_scheme, SuciProfile.ProfileA)

        # Delete the subscriber
        request = SuciProfile(home_net_public_key_id=int(home_net_public_key_id))
        self._stub.DeleteSuciProfile(request)
        self.assertEqual(len(self._stub.ListSuciProfile(Void()).suci_profiles), 0)

        # Delete the subscriber
        with self.assertRaises(grpc.RpcError) as err:
            request = SuciProfile(home_net_public_key_id=int(home_net_public_key_id))
            self._stub.DeleteSuciProfile(request)
        self.assertEqual(err.exception.code(), grpc.StatusCode.NOT_FOUND)


if __name__ == "__main__":
    unittest.main()
