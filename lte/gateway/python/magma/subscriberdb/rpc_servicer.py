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
from google.protobuf.json_format import MessageToJson
from lte.protos import subscriberdb_pb2, subscriberdb_pb2_grpc
from magma.common.rpc_utils import return_void
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
        self._print_grpc(request)
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
        self._print_grpc(request)
        sid = SIDUtils.to_str(request)
        logging.debug("Delete subscriber rpc for sid: %s", sid)
        self._store.delete_subscriber(sid)

    @return_void
    def UpdateSubscriber(self, request, context):
        """
        Updates the subscription data
        """
        try:
            self._print_grpc(request)
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
        self._print_grpc(request)
        sid = SIDUtils.to_str(request)
        try:
            return self._store.get_subscriber_data(sid)
        except SubscriberNotFoundError:
            context.set_details("Subscriber not found: %s" % sid)
            context.set_code(grpc.StatusCode.NOT_FOUND)
            return subscriberdb_pb2.SubscriberData()

    def ListSubscribers(self, request, context):  # pylint:disable=unused-argument
        """
        Returns a list of subscribers from the store
        """
        self._print_grpc(request)
        sids = self._store.list_subscribers()
        sid_msgs = [SIDUtils.to_pb(sid) for sid in sids]
        return subscriberdb_pb2.SubscriberIDSet(sids=sid_msgs)

    def _print_grpc(self, message):
        if self._print_grpc_payload:
            try:
                log_msg = "{} {}".format(
                    message.DESCRIPTOR.full_name,
                    MessageToJson(message),
                )
                # add indentation
                padding = 2 * ' '
                log_msg = ''.join(
                    "{}{}".format(padding, line)
                    for line in log_msg.splitlines(True)
                )

                log_msg = "GRPC message:\n{}".format(log_msg)
                logging.info(log_msg)
            except Exception as e:  # pylint: disable=broad-except
                logging.debug("Exception while trying to log GRPC: %s", e)
