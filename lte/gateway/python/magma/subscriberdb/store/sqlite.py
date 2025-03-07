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
import sqlite3
from collections import defaultdict
from contextlib import contextmanager
from datetime import datetime
from typing import List, NamedTuple

import psutil
from lte.protos.subscriberdb_pb2 import SubscriberData
from magma.subscriberdb.sid import SIDUtils
from magma.subscriberdb.store.base import (
    BaseStore,
    DuplicateSubscriberError,
    SubscriberNotFoundError,
    SubscriberServerTooBusy,
)
from magma.subscriberdb.store.onready import OnDataReady, OnDigestsReady
from orc8r.protos.digest_pb2 import Digest, LeafDigest


class DigestDBInfo(NamedTuple):
    """
    Named tuple representing the locations of the root and leaf digest databases.

    Attributes:
        root_digest_db_location (str): The location of the root digest database.
        leaf_digests_db_location (str): The location of the leaf digests database.
    """
    root_digest_db_location: str
    leaf_digests_db_location: str


class SqliteStore(BaseStore):
    """
    A thread-safe sqlite based implementation of the subscriber database.

    Processes using this store shouldn't be forked since the sqlite connections
    can't be shared by multiple processes.
    """

    def __init__(self, db_location, loop=None, sid_digits=2):
        self._sid_digits = sid_digits  # last digits to be included from subscriber id
        self._n_shards = 10**sid_digits
        self._db_locations = self._create_db_locations(db_location, self._n_shards)

        digest_db_info = self._create_digest_db_locations(db_location)
        self._root_digest_db_location = digest_db_info.root_digest_db_location
        self._leaf_digests_db_location = digest_db_info.leaf_digests_db_location

        self._create_store()
        self._on_ready = OnDataReady(loop=loop)
        self._on_digests_ready = OnDigestsReady(loop=loop)

    def _create_db_locations(self, db_location: str, n_shards: int) -> List[str]:
        # in memory if db_location is not specified
        if not db_location:
            db_location = "/var/opt/magma/"

        # construct db_location items as:
        # file:<path>subscriber<shard>.db?cache=shared
        db_location_list = []

        # file name is passed, use it as a base
        for shard in range(n_shards):
            db_location_list.append(
                'file:'
                + db_location
                + 'subscriber'
                + str(shard)
                + ".db?cache=shared",
            )
            logging.info("db location: %s", db_location_list[shard])

        return db_location_list

    def _create_digest_db_locations(self, db_location: str) -> DigestDBInfo:
        root_digest_db_location = 'file:' + db_location + \
            'subscriber-root-digest.db?cache=shared'
        logging.info("root digest db location: %s", root_digest_db_location)

        leaf_digests_db_location = 'file:' + db_location + \
            'subscriber-leaf-digests.db?cache=shared'
        logging.info(
            "leaf digests db location: %s",
            leaf_digests_db_location,
        )

        digest_db_info = DigestDBInfo(
            root_digest_db_location=root_digest_db_location,
            leaf_digests_db_location=leaf_digests_db_location,
        )
        return digest_db_info

    def _create_store(self) -> None:
        """
        Create the sqlite table for subscribers and digest if they don't exist
        already.
        """
        for db_location in self._db_locations:
            with sqlite3.connect(db_location, uri=True) as conn:
                conn.execute(
                    "CREATE TABLE IF NOT EXISTS subscriberdb"
                    "(subscriber_id text PRIMARY KEY, data text)",
                )

        with sqlite3.connect(self._root_digest_db_location, uri=True) as conn:
            conn.execute(
                "CREATE TABLE IF NOT EXISTS subscriber_root_digest"
                "(digest string PRIMARY KEY, updated_at timestamp)",
            )

        with sqlite3.connect(self._leaf_digests_db_location, uri=True) as conn:
            conn.execute(
                "CREATE TABLE IF NOT EXISTS subscriber_leaf_digests"
                "(sid string PRIMARY KEY, digest string)",
            )

    def add_subscriber(self, subscriber_data: SubscriberData):
        """
        Add the subscriber to store.
        """
        sid = SIDUtils.to_str(subscriber_data.sid)
        data_str = subscriber_data.SerializeToString()
        db_location = self._db_locations[self._sid2bucket(sid)]
        with sqlite3.connect(db_location, uri=True) as conn:
            res = conn.execute(
                "SELECT data FROM subscriberdb WHERE "
                "subscriber_id = ?", (sid,),
            )
            if res.fetchone():
                raise DuplicateSubscriberError(sid)

            conn.execute(
                "INSERT INTO subscriberdb(subscriber_id, data) "
                "VALUES (?, ?)", (sid, data_str),
            )

        self._on_ready.add_subscriber(subscriber_data)

    @contextmanager
    def edit_subscriber(self, subscriber_id):
        """
        Context manager to modify the subscriber data.
        """
        db_location = self._db_locations[self._sid2bucket(subscriber_id)]
        conn = sqlite3.connect(db_location, uri=True)
        try:
            with conn:
                res = conn.execute(
                    "SELECT data FROM subscriberdb WHERE subscriber_id = ?",
                    (subscriber_id,),
                )
                row = res.fetchone()
                if not row:
                    raise SubscriberNotFoundError(subscriber_id)
                subscriber_data = SubscriberData()
                subscriber_data.ParseFromString(row[0])
                yield subscriber_data
                data_str = subscriber_data.SerializeToString()
                conn.execute(
                    "UPDATE subscriberdb SET data = ? "
                    "WHERE subscriber_id = ?",
                    (data_str, subscriber_id),
                )
        finally:
            conn.close()

    def upsert_subscriber(self, subscriber_data: SubscriberData) -> None:
        """
        Check if the given subscriber exists in store. If so, update subscriber
        data; otherwise, add subscriber.

        Args:
            subscriber_data: the data of the subscriber to be upserted.
        """
        sid = SIDUtils.to_str(subscriber_data.sid)
        data_str = subscriber_data.SerializeToString()
        db_location = self._db_locations[self._sid2bucket(sid)]
        with sqlite3.connect(db_location, uri=True) as conn:
            res = conn.execute(
                "SELECT subscriber_id FROM subscriberdb WHERE "
                "subscriber_id = ?", (sid,),
            )
            row = res.fetchone()
            if row is None:
                conn.execute(
                    "INSERT INTO subscriberdb(subscriber_id, data) "
                    "VALUES (?, ?)", (sid, data_str),
                )
            else:
                conn.execute(
                    "UPDATE subscriberdb SET data = ? "
                    "WHERE subscriber_id = ?", (data_str, sid),
                )

        self._on_ready.upsert_subscriber(subscriber_data)

    def delete_subscriber(self, subscriber_id) -> None:
        """
        Delete a subscriber, if present.

        Deleting a subscriber also deletes all digest information.
        This is because deleting a subscriber invalidates the root digest
        and changes the leaf digests. However, because deleting a subscriber
        only happens during testing or debugging, it's easier to just blow
        away all digest data.

        Args:
            subscriber_id: The subscriber ID to delete
        """
        db_location = self._db_locations[self._sid2bucket(subscriber_id)]
        with sqlite3.connect(db_location, uri=True) as conn:
            self.clear_digests()
            conn.execute(
                "DELETE FROM subscriberdb WHERE subscriber_id = ?",
                (subscriber_id,),
            )

        self._on_ready.delete_subscriber(subscriber_id)

    def delete_all_subscribers(self):
        """
        Remove all the subscribers from the store
        """
        for db_location in self._db_locations:
            self.clear_digests()
            with sqlite3.connect(db_location, uri=True) as conn:
                conn.execute("DELETE FROM subscriberdb")

    def get_subscriber_data(self, subscriber_id):
        """
        Return the auth key for the subscriber.
        """
        db_location = self._db_locations[self._sid2bucket(subscriber_id)]
        try:
            with sqlite3.connect(db_location, uri=True) as conn:
                res = conn.execute(
                    "SELECT data FROM subscriberdb WHERE "
                    "subscriber_id = ?", (subscriber_id,),
                )
                row = res.fetchone()
                if not row:
                    raise SubscriberNotFoundError(subscriber_id)
        except sqlite3.OperationalError as exc:
            self._log_db_lock_info(db_location)
            raise SubscriberServerTooBusy(subscriber_id) from exc

        subscriber_data = SubscriberData()
        subscriber_data.ParseFromString(row[0])
        return subscriber_data

    def _log_db_lock_info(self, db_location):
        """
        Log information about processes holding locks on the database file.
        """
        db_parts = db_location.split(":", 1)
        if len(db_parts) != 2 or not db_parts[1]:
            return

        path_str = db_parts[1].split("?")[0]
        for proc in psutil.process_iter(['pid', 'name', 'open_files']):
            try:
                for file in proc.open_files():
                    if file.path == path_str:
                        logging.info("Process holding lock: PID=%d, Name=%s", proc.pid, proc.name())
            except (psutil.NoSuchProcess, psutil.AccessDenied, psutil.ZombieProcess):
                logging.warning("Unable to access process information")

    def list_subscribers(self):
        """
        Return the list of subscribers stored
        """
        sub_list = []
        for db_location in self._db_locations:
            with sqlite3.connect(db_location, uri=True) as conn:
                res = conn.execute(
                    "SELECT subscriber_id FROM subscriberdb",
                )
                sub_list.extend([row[0] for row in res])

        return sub_list

    def update_subscriber(self, subscriber_data):
        """
        Update the subscriber. edit_subscriber should be generally used since that guarantees
        the read/update/write atomicity, but this can be used if the application can guarantee
        the atomicity using a lock.

        Args:
            subscriber_data: SubscriberData protobuf message

        Raises:
            SubscriberNotFoundError: If the subscriber is not present
        """
        sid = SIDUtils.to_str(subscriber_data.sid)
        data_str = subscriber_data.SerializeToString()
        db_location = self._db_locations[self._sid2bucket(sid)]
        with sqlite3.connect(db_location, uri=True) as conn:
            res = conn.execute(
                "UPDATE subscriberdb SET data = ? "
                "WHERE subscriber_id = ?", (data_str, sid),
            )
            if not res.rowcount:
                raise SubscriberNotFoundError(sid)

    def resync(self, subscribers):
        """
        Resync the store with the mentioned list of subscribers.

        This method takes a list of `SubscriberData` objects and resynchronizes the store with the
        provided subscribers. It first groups the subscribers by their bucket using the
        `_sid2bucket` method. Then, for each bucket, it connects to the corresponding database,
        captures the current state of the subscribers, clears the subscriber table, and adds the
        subscribers back with their current state.

        Args:
            subscribers: list of subscribers to be in the store.

        Raises:
            SubscriberNotFoundError: If a subscriber is not found in the current state.
        """
        bucket_subs = defaultdict(list)
        for sub in subscribers:
            sid = SIDUtils.to_str(sub.sid)
            bucket_subs[self._sid2bucket(sid)].append(sub)

        for i, db_location in enumerate(self._db_locations):
            with sqlite3.connect(db_location, uri=True) as conn:
                # Capture the current state of the subscribers
                res = conn.execute(
                    "SELECT subscriber_id, data FROM subscriberdb",
                )
                current_state = {
                    row[0]: SubscriberData().ParseFromString(row[1]).state
                    for row in res
                }

                # Clear all subscribers
                conn.execute("DELETE FROM subscriberdb")

                # Add the subscribers with the current state
                for sub in bucket_subs[i]:
                    sid = SIDUtils.to_str(sub.sid)
                    if sid in current_state:
                        sub.state.CopyFrom(current_state.get(sid))
                    data_str = sub.SerializeToString()
                    conn.execute(
                        "INSERT INTO subscriberdb(subscriber_id, data) "
                        "VALUES (?, ?)", (sid, data_str),
                    )

        self._on_ready.resync(subscribers)

    def get_current_root_digest(self) -> str:
        """
        Retrieve the current root digest from the subscriber database.

        Returns:
            A string containing the current root digest.

        Description:
            This function connects to the root digest database, retrieves the latest digest
            and returns it.
        """
        with sqlite3.connect(self._root_digest_db_location, uri=True) as conn:
            res = conn.execute(
                "SELECT digest, updated_at FROM subscriber_root_digest "
                "ORDER BY updated_at DESC",
            )
            row = res.fetchone()
            if not row:
                row = ["", None]

        digest = str(row[0])
        logging.info("get digest stored in gateway: %s", digest)
        return digest

    def update_root_digest(self, new_digest: str) -> None:
        """
        Update the root digest in the subscriber database with the new digest provided.

        Args:
            new_digest (str): The new digest to update.

        Description:
            This function connects to the root digest database, deletes the existing digest, and
            inserts the new one along with the current datetime, then calls the `update_root_digest`
            method of the `_on_digests_ready` object with the new digest.

        Note:
            - The function assumes that the gRPC client and command-line arguments are valid.
            - The function does not perform any input validation.
            - The function does not handle any exceptions.
        """
        with sqlite3.connect(self._root_digest_db_location, uri=True) as conn:
            conn.execute("DELETE FROM subscriber_root_digest")

            conn.execute(
                "INSERT INTO subscriber_root_digest(digest, updated_at) "
                "VALUES (?, ?)", (new_digest, datetime.now()),
            )

        logging.info("update root digest stored in gateway: %s", new_digest)
        self._on_digests_ready.update_root_digest(new_digest)

    def get_current_leaf_digests(self) -> List[LeafDigest]:
        """
        Retrieve the current leaf digests from the subscriber leaf digests database.

        Returns:
            A list of LeafDigest objects containing the current leaf digests.
        """
        digests = []
        with sqlite3.connect(self._leaf_digests_db_location, uri=True) as conn:
            res = conn.execute(
                "SELECT sid, digest FROM subscriber_leaf_digests ",
            )

            for row in res:
                digest = LeafDigest(
                    id=row[0],
                    digest=Digest(md5_base64_digest=row[1]),
                )
                digests.append(digest)

        return digests

    def update_leaf_digests(self, new_digests: List[LeafDigest]) -> None:
        """
        Update the leaf digests in the subscriber database with the new digests provided.

        Args:
            new_digests: A list of LeafDigest objects containing the new digests to update.
        """
        with sqlite3.connect(self._leaf_digests_db_location, uri=True) as conn:
            conn.execute(
                "DELETE FROM subscriber_leaf_digests",
            )
            for leaf_digest in new_digests:
                sid = leaf_digest.id
                digest = leaf_digest.digest.md5_base64_digest
                conn.execute(
                    "INSERT INTO subscriber_leaf_digests(sid, digest)"
                    "VALUES (?, ?)", (sid, digest),
                )

        self._on_digests_ready.update_leaf_digests(new_digests)

    def clear_digests(self):
        """
        Clear the digests stored in the database by updating the root digest to an empty string and
        updating the leaf digests to an empty list.
        """
        self.update_root_digest("")
        self.update_leaf_digests([])

    async def on_ready(self):
        """
        Wait asynchronously for the `_on_ready` event to be set.

        Returns:
            Awaitable[None]: An awaitable that resolves when the `_on_ready` event is set.
        """
        return await self._on_ready.event.wait()

    async def on_digests_ready(self):
        """
        Wait asynchronously for the `_on_digests_ready` event to be set.

        Returns:
            Awaitable[None]: An awaitable that resolves when the `_on_digests_ready` event is set.
        """
        return await self._on_digests_ready.event.wait()

    def _update_apn(self, apn_config, apn_data):
        """
        Update the APN configuration based on the provided APN data.

        Args:
            apn_config: The APN configuration to be updated.
            apn_data: The APN data containing the new values.
        """
        apn_config.is_default = apn_data.is_default
        apn_config.service_selection = apn_data.service_selection
        apn_config.qos_profile.class_id = apn_data.qos_profile.class_id
        apn_config.qos_profile.priority_level = (
            apn_data.qos_profile.priority_level
        )
        apn_config.qos_profile.preemption_capability = (
            apn_data.qos_profile.preemption_capability
        )
        apn_config.qos_profile.preemption_vulnerability = (
            apn_data.qos_profile.preemption_vulnerability
        )
        apn_config.ambr.max_bandwidth_ul = apn_data.ambr.max_bandwidth_ul
        apn_config.ambr.max_bandwidth_dl = apn_data.ambr.max_bandwidth_dl

    def _sid2bucket(self, subscriber_id):
        """
        Calculate the bucket number based on the last `self._sid_digits` digits of the `subscriber_id`.

        Args:
            subscriber_id (str): The ID of the subscriber.

        Returns:
            int: The bucket number corresponding to the last `self._sid_digits` digits of the `subscriber_id`.
                If the conversion to an integer fails, a default value of 0 is returned.
        """
        try:
            bucket = int(subscriber_id[-self._sid_digits:])
        except (TypeError, ValueError):
            logging.info(
                "Last %d digits of subscriber id %s cannot mapped to a bucket:"
                " default to bucket 0", self._sid_digits, subscriber_id,
            )
            bucket = 0
        return bucket
