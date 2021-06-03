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
from typing import Any

from lte.protos.s6a_service_pb2 import DeleteSubscriberRequest
from lte.protos.s6a_service_pb2_grpc import S6aServiceStub
from lte.protos.subscriberdb_pb2 import LTESubscription, SubscriberData
from magma.common.service_registry import ServiceRegistry
from magma.common.streamer import StreamerClient


class SubscriberDBStreamerCallback(StreamerClient.Callback):
    """
    Callback implementation for the SubscriberDB StreamerClient instance.
    """

    def __init__(self, store, loop):
        self._store = store
        self._loop = loop

    def get_request_args(self, stream_name: str) -> Any:
        return None

    def process_update(self, stream_name, updates, resync):
        """
        The cloud streams ALL subscribers registered, both active and inactive.
        Since we don't have a good way of knowing whether a detach succeeds or
        fails to update the local database correctly, we have to send down all
        subscribers to keep trying to delete inactive subscribers.
        TODO we can optimize a bit on the MME side to not detach already
        detached subscribers.
        """
        logging.info(
            "Processing %d subscriber updates (resync=%s)",
            len(updates), resync,
        )

        if resync:
            # TODO:
            # - handle database exceptions
            keys = []
            subscribers = []
            active_subscriber_ids = []
            for update in updates:
                sub = SubscriberData()
                sub.ParseFromString(update.value)
                subscribers.append(sub)
                keys.append(update.key)
                if sub.lte.state == LTESubscription.ACTIVE:
                    active_subscriber_ids.append(update.key)
            old_sub_ids = self._store.list_subscribers()
            # Only compare active subscribers against the database to decide
            # what to detach.
            self.detach_deleted_subscribers(old_sub_ids, active_subscriber_ids)
            logging.debug("Resync with subscribers: %s", ','.join(keys))
            self._store.resync(subscribers)
        else:
            # TODO: implement updates
            pass

    def detach_deleted_subscribers(self, old_sub_ids, new_sub_ids):
        """
        Compares current subscriber ids and new subscriber ids list
        just streamed from the cloud to figure out the deleted subscribers.
        Then send grpc DeleteSubscriber request to mme to detach all the
        deleted subscribers.
        :param old_sub_ids: a list of old subscriber ids in the store.
        :param new_sub_ids: a list of new active subscriber ids
                just streamed from the cloud
        :return: n/a
        """
        # THIS IS A HACK UNTIL WE FIX THIS ON CLOUD
        # We accept IMSIs with or without 'IMSI' prepended on cloud, but we
        # always store IMSIs on local subscriberdb with IMSI prepended. If the
        # cloud streams down subscriber IDs without 'IMSI' prepended,
        # subscriberdb will try to delete all of the subscribers from MME every
        # time it streams from cloud because the set membership will fail
        # when comparing '12345' to 'IMSI12345'.
        new_sub_ids = set(
            map(
                lambda s: 'IMSI' + s if not s.startswith('IMSI') else s,
                new_sub_ids,
            ),
        )
        deleted_sub_ids = [
            sub_id for sub_id in old_sub_ids
            if sub_id not in set(new_sub_ids)
        ]
        if len(deleted_sub_ids) == 0:
            return
        # send detach request to mme for all deleted subscribers.
        chan = ServiceRegistry.get_rpc_channel(
            's6a_service',
            ServiceRegistry.LOCAL,
        )
        client = S6aServiceStub(chan)
        req = DeleteSubscriberRequest()

        # mme expects a list of IMSIs without "IMSI" prefix
        imsis_to_delete_without_prefix = [sub[4:] for sub in deleted_sub_ids]

        req.imsi_list.extend(imsis_to_delete_without_prefix)
        future = client.DeleteSubscriber.future(req)
        future.add_done_callback(
            lambda future:
            self._loop.call_soon_threadsafe(
                self.detach_deleted_subscribers_done,
                future,
            ),
        )

    @staticmethod
    def detach_deleted_subscribers_done(delete_future):
        """
        Detach deleted subscribers callback to handle exceptions
        """
        err = delete_future.exception()
        if err:
            logging.error(
                "Detach Deleted Subscribers Error! [%s] %s",
                err.code(), err.details(),
            )
