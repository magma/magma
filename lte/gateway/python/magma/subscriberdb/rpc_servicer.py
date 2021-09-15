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

import logging

import grpc
from lte.protos import subscriberdb_pb2, subscriberdb_pb2_grpc
from magma.common.rpc_utils import print_grpc, return_void
from magma.subscriberdb.sid import SIDUtils

from .store.base import DuplicateSubscriberError, SubscriberNotFoundError


class SubscriberDBRpcServicer(subscriberdb_pb2_grpc.SubscriberDBServicer):
    """
    gRPC based server for the SubscriberDB.
    """

    def __init__(self, store, print_grpc_payload: bool = False):
        """
        Store should be thread-safe since we use a thread pool for requests.
        """
        self._store = store
        self._print_grpc_payload = print_grpc_payload

    def add_to_server(self, server):
        """
        Add the servicer to a gRPC server
        """
        subscriberdb_pb2_grpc.add_SubscriberDBServicer_to_server(self, server)

    @return_void
    def AddSubscriber(self, request, context):
        """
        Adds a subscriber to the store
        """
        print_grpc(request, self._print_grpc_payload, "Add Subscriber Request:")
        sid = SIDUtils.to_str(request.sid)
        logging.debug("Add subscriber rpc for sid: %s", sid)
        try:
            self._store.add_subscriber(request)
        except DuplicateSubscriberError:
            context.set_details("Duplicate subscriber: %s" % sid)
            context.set_code(grpc.StatusCode.ALREADY_EXISTS)

    @return_void
    def DeleteSubscriber(self, request, context):
        """
        Deletes a subscriber from the store
        """
        print_grpc(
            request, self._print_grpc_payload,
            "Delete Subscriber Request:",
        )
        sid = SIDUtils.to_str(request)
        logging.debug("Delete subscriber rpc for sid: %s", sid)
        self._store.delete_subscriber(sid)

    @return_void
    def UpdateSubscriber(self, request, context):
        """
        Updates the subscription data
        """
        try:
            print_grpc(
                request, self._print_grpc_payload,
                "Update Subscriber Request",
            )
        except Exception as e:  # pylint: disable=broad-except
            logging.debug("Exception while trying to log GRPC: %s", e)
        sid = SIDUtils.to_str(request.data.sid)
        try:
            with self._store.edit_subscriber(sid) as subs:
                request.mask.MergeMessage(
                    request.data, subs, replace_message_field=True,
                )
        except SubscriberNotFoundError:
            context.set_details("Subscriber not found: %s" % sid)
            context.set_code(grpc.StatusCode.NOT_FOUND)

    def GetSubscriberData(self, request, context):
        """
        Returns the subscription data for the subscriber
        """
        print_grpc(
            request, self._print_grpc_payload,
            "Get Subscriber Data Request:",
        )
        sid = SIDUtils.to_str(request)
        try:
            response = self._store.get_subscriber_data(sid)
        except SubscriberNotFoundError:
            context.set_details("Subscriber not found: %s" % sid)
            context.set_code(grpc.StatusCode.NOT_FOUND)
            response = subscriberdb_pb2.SubscriberData()
        print_grpc(
            response, self._print_grpc_payload,
            "Get Subscriber Data Response:",
        )
        return response

    def ListSubscribers(self, request, context):  # pylint:disable=unused-argument
        """
        Returns a list of subscribers from the store
        """
        print_grpc(
            request, self._print_grpc_payload,
            "List Subscribers Request:",
        )
        sids = self._store.list_subscribers()
        sid_msgs = [SIDUtils.to_pb(sid) for sid in sids]
        response = subscriberdb_pb2.SubscriberIDSet(sids=sid_msgs)
        print_grpc(
            response, self._print_grpc_payload,
            "List Subscribers Response:",
        )
        return response
