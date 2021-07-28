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
import asyncio
import datetime
import logging
from typing import List, NamedTuple, Optional

import grpc
from lte.protos.s6a_service_pb2 import DeleteSubscriberRequest
from lte.protos.s6a_service_pb2_grpc import S6aServiceStub
from lte.protos.subscriberdb_pb2 import (
    CheckSubscribersInSyncRequest,
    Digest,
    ListSubscribersRequest,
    LTESubscription,
    SubscriberData,
    SubscriberDigestWithID,
    SyncSubscribersRequest,
)
from magma.common.grpc_client_manager import GRPCClientManager
from magma.common.rpc_utils import grpc_async_wrapper
from magma.common.sdwatchdog import SDWatchdogTask
from magma.common.service_registry import ServiceRegistry
from magma.subscriberdb.metrics import (
    SUBSCRIBER_SYNC_FAILURE_TOTAL,
    SUBSCRIBER_SYNC_LATENCY,
    SUBSCRIBER_SYNC_SUCCESS_TOTAL,
)
from magma.subscriberdb.store.sqlite import SqliteStore

CloudSubscribersInfo = NamedTuple(
    'CloudSubscribersInfo', [
        ('subscribers', List[SubscriberData]),
        ('flat_digest', Optional[Digest]),
        ('per_sub_digests', Optional[List[SubscriberDigestWithID]]),
    ],
)


class SubscriberDBCloudClient(SDWatchdogTask):
    """
    SubscriberDBCloudClient for requesting subscribers from Orchestrator.

    SubscriberDBCloudClient calls the Orchestrator's SubscriberDBCloud service
    to fetch pages of subscribers stored in Orchestrator and updates the local
    subscriberdb store with the changes.
    """

    SUBSCRIBERDB_REQUEST_TIMEOUT = 10

    def __init__(
        self,
        loop: asyncio.AbstractEventLoop,
        store: SqliteStore,
        subscriber_page_size: int,
        sync_interval: int,
        grpc_client_manager: GRPCClientManager,
    ):
        """
        Initialize subscriberdb client

        Args:
            loop: asyncio event loop
            store: SqliteStore for subscribers
            subscriber_page_size: integer for page size
            sync_interval: integer for frequency of subscriber streaming
            grpc_client_manager: GRPCClientManager for gRPC client mgmt

        Returns: None
        """
        super().__init__(
            sync_interval,
            loop,
        )
        self._loop = loop
        self._subscriber_page_size = subscriber_page_size
        self._store = store

        # grpc_client_manager to manage grpc client recycling
        self._grpc_client_manager = grpc_client_manager

    async def _run(self) -> None:
        in_sync = await self._check_subscribers_in_sync()
        if in_sync:
            return

        resync = await self._sync_subscribers()
        if not resync:
            return

        subscribers_info = await self._get_all_subscribers()
        if subscribers_info is None:
            return

        # Process subscriber data before digest data, in case there's a gateway
        # failure between the calls
        self._process_subscribers(subscribers_info.subscribers)
        self._update_flat_digest(subscribers_info.flat_digest)
        self._update_per_sub_digests(subscribers_info.per_sub_digests)

    async def _check_subscribers_in_sync(self) -> bool:
        """
        Check if the local subscriber data is up-to-date with the cloud by
        comparing flat digests

        Returns:
            boolean value for whether the local data is in sync
        """
        subscriberdb_cloud_client = self._grpc_client_manager.get_client()
        req = CheckSubscribersInSyncRequest(
            flat_digest=Digest(
                md5_base64_digest=self._store.get_current_digest(),
            ),
        )
        try:
            res = await grpc_async_wrapper(
                subscriberdb_cloud_client.CheckSubscribersInSync.future(
                    req,
                    self.SUBSCRIBERDB_REQUEST_TIMEOUT,
                ),
                self._loop,
            )
        except grpc.RpcError as err:
            logging.error(
                "Check subscribers in sync request error! [%s] %s", err.code(),
                err.details(),
            )
            return False
        return res.in_sync

    async def _sync_subscribers(self) -> bool:
        """
        Sync local subscriber data and digests with the cloud if didn't receive
        resync signal.

        Returns:
            boolean value for whether a resync with cloud is needed
        """
        subscriberdb_cloud_client = self._grpc_client_manager.get_client()
        req = SyncSubscribersRequest(
            per_sub_digests=self._store.get_current_per_sub_digests(),
        )
        try:
            res = await grpc_async_wrapper(
                subscriberdb_cloud_client.SyncSubscribers.future(
                    req,
                    self.SUBSCRIBERDB_REQUEST_TIMEOUT,
                ),
                self._loop,
            )
        except grpc.RpcError as err:
            logging.error(
                "Sync subscribers request error! [%s] %s", err.code(),
                err.details(),
            )
            return True

        if not res.resync:
            self._update_flat_digest(res.flat_digest)
            self._update_per_sub_digests(res.per_sub_digests)

            # TODO(hcgatewood): switch to bulk editing subscriber data rows
            for subscriber_data in res.to_renew:
                self._store.upsert_subscriber(subscriber_data)
            for sid in res.deleted:
                self._store.delete_subscriber(sid)
            self._detach_subscribers_by_ids(res.deleted)

        return res.resync

    async def _get_all_subscribers(self) -> Optional[CloudSubscribersInfo]:
        subscriberdb_cloud_client = self._grpc_client_manager.get_client()
        subscribers = []
        flat_digest = None
        per_sub_digests = None
        req_page_token = ""  # noqa: S105
        found_empty_token = False
        sync_start = datetime.datetime.now()

        # Next page token empty means read all updates
        while not found_empty_token:  # noqa: S105
            try:
                req = ListSubscribersRequest(
                    page_size=self._subscriber_page_size,
                    page_token=req_page_token,
                )
                res = await grpc_async_wrapper(
                    subscriberdb_cloud_client.ListSubscribers.future(
                        req,
                        self.SUBSCRIBERDB_REQUEST_TIMEOUT,
                    ),
                    self._loop,
                )
                subscribers.extend(res.subscribers)
                # Cloud sends back digests during request for the first page
                if req_page_token == "":
                    flat_digest = res.flat_digest
                    per_sub_digests = res.per_sub_digests

                req_page_token = res.next_page_token
                found_empty_token = req_page_token == ""

            except grpc.RpcError as err:
                logging.error(
                    "Fetch subscribers error! [%s] %s", err.code(),
                    err.details(),
                )
                time_elapsed = datetime.datetime.now() - sync_start
                SUBSCRIBER_SYNC_LATENCY.observe(
                    time_elapsed.total_seconds() * 1000,
                )
                SUBSCRIBER_SYNC_FAILURE_TOTAL.inc()
                return None
        logging.info(
            "Successfully fetched all subscriber "
            "pages from the cloud!",
        )
        SUBSCRIBER_SYNC_SUCCESS_TOTAL.inc()
        time_elapsed = datetime.datetime.now() - sync_start
        SUBSCRIBER_SYNC_LATENCY.observe(
            time_elapsed.total_seconds() * 1000,
        )

        subscribers_info = CloudSubscribersInfo(
            subscribers=subscribers,
            flat_digest=flat_digest,
            per_sub_digests=per_sub_digests,
        )
        return subscribers_info

    def _update_flat_digest(self, flat_digest: Optional[Digest]) -> None:
        if Digest is None:
            return
        self._store.update_digest(flat_digest.md5_base64_digest)

    def _update_per_sub_digests(
            self,
            per_sub_digests: Optional[List[SubscriberDigestWithID]],
    ) -> None:
        if per_sub_digests is None:
            return
        self._store.update_per_sub_digests(per_sub_digests)

    def _process_subscribers(self, subscribers: List[SubscriberData]) -> None:
        active_subscriber_ids = []
        sids = []
        for subscriber in subscribers:
            sids.append(subscriber.sid.id)
            if subscriber.lte.state == LTESubscription.ACTIVE:
                active_subscriber_ids.append(subscriber.sid.id)
        old_sub_ids = self._store.list_subscribers()
        # Only compare active subscribers against the database to decide
        # what to detach.
        self._detach_deleted_subscribers(old_sub_ids, active_subscriber_ids)
        logging.debug("Resync with subscribers: %s", ','.join(sids))
        self._store.resync(subscribers)

    def _detach_deleted_subscribers(self, old_sub_ids, new_sub_ids):
        """
        Detach deleted subscribers from store and network.

        Compares current subscriber ids and new subscriber ids list
        just streamed from the cloud to figure out the deleted subscribers.
        Then send grpc DeleteSubscriber request to mme to detach all the
        deleted subscribers.

        Args:
            old_sub_ids: a list of old subscriber ids in the store
            new_sub_ids: a list of new active subscriber ids

        Returns:
            None

        """
        # THIS IS A HACK UNTIL WE FIX THIS ON CLOUD
        # We accept IMSIs with or without 'IMSI' prepended on cloud, but we
        # always store IMSIs on local subscriberdb with IMSI prepended. If the
        # cloud streams down subscriber IDs without 'IMSI' prepended,
        # subscriberdb will try to delete all of the subscribers from MME every
        # time it streams from cloud because the set membership will fail
        # when comparing '12345' to 'IMSI12345'.
        new_sub_ids = {
            'IMSI{imsiVal}'.format(imsiVal=s) if not s.startswith('IMSI')
            else s
            for s in new_sub_ids
        }
        deleted_sub_ids = [
            sub_id for sub_id in old_sub_ids
            if sub_id not in set(new_sub_ids)
        ]
        if not deleted_sub_ids:
            return
        self._detach_subscribers_by_ids(deleted_sub_ids)

    def _detach_subscribers_by_ids(self, deleted_sub_ids: List[str]):
        """
        Sends grpc DeleteSubscriber request to mme to detach all subscribers
        given as args.

        Args:
            deleted_sub_ids: a list of old subscriber ids in the store,
                             prefixed by subscriber type

        Returns:
            None

        """
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
                self._detach_deleted_subscribers_done,
                future,
            ),
        )

    def _detach_deleted_subscribers_done(self, delete_future):
        """
        Detach deleted subscribers callback to handle exceptions

        Args:
            delete_future: future of delete RPC call
        """
        err = delete_future.exception()
        if err:
            logging.error(
                "Detach Deleted Subscribers Error! [%s] %s",
                err.code(), err.details(),
            )
