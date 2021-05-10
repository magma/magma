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

from lte.protos.subscriberdb_pb2 import SubscriberData
from magma.subscriberdb.sid import SIDUtils

from .base import BaseStore, DuplicateSubscriberError, SubscriberNotFoundError
from .onready import OnDataReady


class SqliteStore(BaseStore):
    """
    A thread-safe sqlite based implementation of the subscriber database.

    Processes using this store shouldn't be forked since the sqlite connections
    can't be shared by multiple processes.
    """

    def __init__(self, db_location, loop=None, sid_digits=2):
        self._sid_digits = sid_digits # last digits to be included from subscriber id
        self._n_shards = 10**sid_digits
        self._db_locations = self._create_db_locations(db_location, self._n_shards)
        self._create_store()
        self._on_ready = OnDataReady(loop=loop)


    def _create_db_locations(self, db_location, n_shards):
        # in memory if db_location is not specified
        if not db_location:
            db_location = "/var/opt/magma/"

        # construct db_location items as: file:<path>subscriber<shard>.db?cache=shared
        db_location_list = []

        # file name is passed, use it as a base
        for shard in range(n_shards):
            db_location_list.append('file:' + db_location + 'subscriber' + str(shard) + ".db?cache=shared")
            logging.info("db location: %s", db_location_list[shard])
        return db_location_list

    def _create_store(self):
        """
        Create the sqlite table if it doesn't exist already.
        """
        for db_location in self._db_locations:
            conn = sqlite3.connect(db_location, uri=True)
            try:
                with conn:
                    conn.execute("CREATE TABLE IF NOT EXISTS subscriberdb"
                                      "(subscriber_id text PRIMARY KEY, data text)")
            finally:
                conn.close()

    def add_subscriber(self, subscriber_data):
        """
        Method that adds the subscriber.
        """
        sid = SIDUtils.to_str(subscriber_data.sid)
        data_str = subscriber_data.SerializeToString()
        db_location = self._db_locations[self._sid2bucket(sid)]
        conn = sqlite3.connect(db_location, uri=True)
        try:
            with conn:
                res = conn.execute("SELECT data FROM subscriberdb WHERE "
                                        "subscriber_id = ?", (sid, ))
                if res.fetchone():
                    raise DuplicateSubscriberError(sid)

                conn.execute("INSERT INTO subscriberdb(subscriber_id, data) "
                                "VALUES (?, ?)", (sid, data_str))
        finally:
            conn.close()
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
                    "SELECT data FROM subscriberdb WHERE " "subscriber_id = ?",
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
                    "UPDATE subscriberdb SET data = ? " "WHERE subscriber_id = ?",
                    (data_str, subscriber_id),
                )
        finally:
            conn.close()

    def delete_subscriber(self, subscriber_id):
        """
        Method that deletes a subscriber, if present.
        """
        db_location = self._db_locations[self._sid2bucket(subscriber_id)]
        conn = sqlite3.connect(db_location, uri=True)
        try:
            with conn:
                conn.execute(
                    "DELETE FROM subscriberdb WHERE " "subscriber_id = ?",
                    (subscriber_id,),
                )
        finally:
            conn.close()

    def delete_all_subscribers(self):
        """
        Method that removes all the subscribers from the store
        """
        for db_location in self._db_locations:
            conn = sqlite3.connect(db_location, uri=True)
            try:
                with conn:
                    conn.execute("DELETE FROM subscriberdb")
            finally:
                conn.close()

    def get_subscriber_data(self, subscriber_id):
        """
        Method that returns the auth key for the subscriber.
        """
        db_location = self._db_locations[self._sid2bucket(subscriber_id)]
        conn = sqlite3.connect(db_location, uri=True)
        try:
            with conn:
                res = conn.execute("SELECT data FROM subscriberdb WHERE "
                                        "subscriber_id = ?", (subscriber_id, ))
                row = res.fetchone()
                if not row:
                    raise SubscriberNotFoundError(subscriber_id)
        finally:
            conn.close()
        subscriber_data = SubscriberData()
        subscriber_data.ParseFromString(row[0])
        return subscriber_data

    def list_subscribers(self):
        """
        Method that returns the list of subscribers stored
        """
        sub_list = []
        for db_location in self._db_locations:
            conn = sqlite3.connect(db_location, uri=True)
            try:
                with conn:
                    res = conn.execute("SELECT subscriber_id FROM subscriberdb")
                    sub_list.extend([row[0] for row in res])
            finally:
                conn.close()
        return sub_list

    def update_subscriber(self, subscriber_data):
        """
        Method that updates the subscriber. edit_subscriber should
        be generally used since that guarantees the read/update/write
        atomicity, but this can be used if the application can
        guarantee the atomicity using a lock.

        Args:
            subscriber_data - SubscriberData protobuf message
        Raises:
            SubscriberNotFoundError if the subscriber is not present

        """
        sid = SIDUtils.to_str(subscriber_data.sid)
        data_str = subscriber_data.SerializeToString()
        db_location = self._db_locations[self._sid2bucket(sid)]
        conn = sqlite3.connect(db_location, uri=True)
        try:
            with conn:
                res = conn.execute("UPDATE subscriberdb SET data = ? "
                                        "WHERE subscriber_id = ?", (data_str, sid))
                if not res.rowcount:
                    raise SubscriberNotFoundError(sid)
        finally:
            conn.close()

    def resync(self, subscribers):
        """
        Method that should resync the store with the mentioned list of
        subscribers. The resync leaves the current state of subscribers
        intact.

        Args:
            subscribers - list of subscribers to be in the store.
        """
        bucket_subs = defaultdict(list)
        for sub in subscribers:
            sid = SIDUtils.to_str(sub.sid)
            bucket_subs[self._sid2bucket(sid)].append(sub)

        for i, db_location in enumerate(self._db_locations):
            conn = sqlite3.connect(db_location, uri=True)
            try:
                with conn:
                    # Capture the current state of the subscribers
                    res = conn.execute("SELECT subscriber_id, data FROM subscriberdb")
                    current_state = {}
                    for row in res:
                        sub = SubscriberData()
                        sub.ParseFromString(row[1])
                        current_state[row[0]] = sub.state

                    # Clear all subscribers
                    conn.execute("DELETE FROM subscriberdb")

                    # Add the subscribers with the current state
                    for sub in bucket_subs[i]:
                        sid = SIDUtils.to_str(sub.sid)
                        if sid in current_state:
                            sub.state.CopyFrom(current_state[sid])
                        data_str = sub.SerializeToString()
                        conn.execute("INSERT INTO subscriberdb(subscriber_id, data) "
                                     "VALUES (?, ?)", (sid, data_str))
            finally:
                conn.close()
        self._on_ready.resync(subscribers)

    async def on_ready(self):
        return await self._on_ready.event.wait()

    def _update_apn(self, apn_config, apn_data):
        """
        Method that populates apn data.
        """
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
        Maps Subscriber ID to bucket
        """
        try:
            bucket = int(subscriber_id[-self._sid_digits:])
        except (TypeError, ValueError):
            logging.info("Last %d digits of subscriber id %s cannot mapped to a bucket:"
                         " default to bucket 0", self._sid_digits, subscriber_id)
            bucket = 0
        return bucket
