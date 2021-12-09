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
from typing import NamedTuple

import grpc
from lte.protos import subscriberdb_pb2, subscriberdb_pb2_grpc
from magma.common.rpc_utils import print_grpc, return_void
from magma.common.sentry import (
    SentryStatus,
    get_sentry_status,
    send_uncaught_errors_to_monitoring,
)
from magma.subscriberdb.sid import SIDUtils
from magma.subscriberdb.store.base import (
    DuplicateSubscriberError,
    SubscriberNotFoundError,
)

enable_sentry_wrapper = get_sentry_status("subscriberdb") == SentryStatus.SEND_SELECTED_ERRORS
suci_profile_data = NamedTuple(
    'suci_profile_data', [
        ('protection_scheme', int),
        ('home_net_public_key', bytes),
        ('home_net_private_key', bytes),
    ],
)


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
    @send_uncaught_errors_to_monitoring(enable_sentry_wrapper)
    def AddSubscriber(self, request, context):
        """
        Add a subscriber to the store
        """
        print_grpc(
            request, self._print_grpc_payload,
            "Add Subscriber Request:",
        )
        sid = SIDUtils.to_str(request.sid)
        logging.debug("Add subscriber rpc for sid: %s", sid)
        try:
            self._store.add_subscriber(request)
        except DuplicateSubscriberError:
            context.set_details("Duplicate subscriber: %s" % sid)
            context.set_code(grpc.StatusCode.ALREADY_EXISTS)

    @return_void
    @send_uncaught_errors_to_monitoring(enable_sentry_wrapper)
    def DeleteSubscriber(self, request, context):
        """
        Delete a subscriber from the store
        """
        print_grpc(
            request, self._print_grpc_payload,
            "Delete Subscriber Request:",
        )
        sid = SIDUtils.to_str(request)
        logging.debug("Delete subscriber rpc for sid: %s", sid)
        self._store.delete_subscriber(sid)

    @return_void
    @send_uncaught_errors_to_monitoring(enable_sentry_wrapper)
    def UpdateSubscriber(self, request, context):
        """
        Update the subscription data
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

    @send_uncaught_errors_to_monitoring(enable_sentry_wrapper)
    def GetSubscriberData(self, request, context):
        """
        Return the subscription data for the subscriber
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

    @send_uncaught_errors_to_monitoring(enable_sentry_wrapper)
    def ListSubscribers(self, request, context):  # pylint:disable=unused-argument
        """
        Return a list of subscribers from the store
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


class SuciProfileDBRpcServicer(subscriberdb_pb2_grpc.SuciProfileDBServicer):
    """
    gRPC based server for the SubscriberDB.
    """

    def __init__(
        self, store, suciprofile_db_dict: dict,
        print_grpc_payload: bool = False,
    ):
        """
        Store should be thread-safe since we use a thread pool for requests.
        """
        self._store = store
        self.suciprofile_db_dict = suciprofile_db_dict
        self._print_grpc_payload = print_grpc_payload

    def add_to_server(self, server):
        """
        Add the servicer to a gRPC server
        """
        subscriberdb_pb2_grpc.add_SuciProfileDBServicer_to_server(self, server)

    @return_void
    def AddSuciProfile(self, request, context):
        """
        AddSuciProfile - Adds a suciprofile to the store        # noqa: D403

        Args:
            request: SuciProfile
            context: context

        Returns: None

        Raises:
            RpcError: Key already exists
        """
        print_grpc(
            request, self._print_grpc_payload,
            "Add SuciProfile Request:",
        )

        if request.home_net_public_key_id in self.suciprofile_db_dict.keys():
            logging.warning(
                "home_net_public_key_id:%d already exist",
                request.home_net_public_key_id,
            )
            context.set_details("Duplicate suciprofile")
            context.set_code(grpc.StatusCode.ALREADY_EXISTS)
            raise grpc.RpcError("Key already exists")

        self.suciprofile_db_dict[request.home_net_public_key_id] = suci_profile_data(
            request.protection_scheme, request.home_net_public_key,
            request.home_net_private_key,
        )

    @return_void
    def DeleteSuciProfile(self, request, context):
        """
        DeleteSuciProfile - Deletes suciprofile from the store      # noqa: D403

        Args:
            request: SuciProfile
            context: context

        Returns: None

        Raises:
            RpcError: Key is invalid

        """
        print_grpc(
            request, self._print_grpc_payload,
            "Delete SuciProfile Request:",
        )

        if request.home_net_public_key_id not in self.suciprofile_db_dict.keys():
            logging.warning(
                "The home_net_public_key_id:%d is not a valid key,"
                "try again", request.home_net_public_key_id,
            )
            context.set_details("Suciprofile not found")
            context.set_code(grpc.StatusCode.NOT_FOUND)
            raise grpc.RpcError("Key is invalid")
        else:
            del self.suciprofile_db_dict[request.home_net_public_key_id]

    def ListSuciProfile(self, request, context):
        """
        ListSuciProfile - Returns SuciProfile list      # noqa: D403

        Args:
            request: SuciProfile
            context: context

        Returns:
            response: SuciProfileList
        """
        print_grpc(
            request, self._print_grpc_payload,
            "List Suciprofile Request:",
        )
        suciprofiles = []
        for k, v in self.suciprofile_db_dict.items():
            suciprofiles.append(
                subscriberdb_pb2.SuciProfile(
                    home_net_public_key_id=int(k),
                    protection_scheme=v.protection_scheme,
                    home_net_public_key=v.home_net_public_key,
                    home_net_private_key=v.home_net_private_key,
                ),
            )

        response = subscriberdb_pb2.SuciProfileList(suci_profiles=suciprofiles)
        print_grpc(
            response, self._print_grpc_payload,
            "List SuciProfile Response:",
        )
        return response
