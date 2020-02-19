"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import sqlite3
import threading
from contextlib import contextmanager

from lte.protos.subscriberdb_pb2 import SubscriberData, Non3GPPUserProfile

from magma.subscriberdb.sid import SIDUtils
from .base import (
    BaseStore,
    DuplicateSubscriberError,
    SubscriberNotFoundError,
    ApnNotFoundError,
    DuplicateApnError,
)
from .onready import OnDataReady


class SqliteStore(BaseStore):
    """
    A thread-safe sqlite based implementation of the subscriber database.

    Processes using this store shouldn't be forked since the sqlite connections
    can't be shared by multiple processes.
    """

    def __init__(self, db_location, loop=None):
        self._db_location = db_location
        self._tlocal = threading.local()
        self._create_store()
        self._on_ready = OnDataReady(loop=loop)

    @property
    def conn(self):
        """
        Returns a thread local connection to the sqlite db.
        """
        if not getattr(self._tlocal, 'conn', None):
            self._tlocal.conn = sqlite3.connect(self._db_location, uri=True)
        return self._tlocal.conn

    def _create_store(self):
        """
        Create the sqlite table if it doesn't exist already.
        """
        with self.conn:
            self.conn.execute("CREATE TABLE IF NOT EXISTS subscriberdb"
                              "(subscriber_id text PRIMARY KEY, data text)")
        with self.conn:
            self.conn.execute(
                "CREATE TABLE IF NOT EXISTS apndb"
                "(apn_name text PRIMARY KEY, data text)"
            )

    def add_subscriber(self, subscriber_data):
        """
        Method that adds the subscriber.
        """
        sid = SIDUtils.to_str(subscriber_data.sid)
        data_str = subscriber_data.SerializeToString()
        with self.conn:
            res = self.conn.execute("SELECT data FROM subscriberdb WHERE "
                                    "subscriber_id = ?", (sid, ))
            if res.fetchone():
                raise DuplicateSubscriberError(sid)

            if subscriber_data.non_3gpp.apn_config:
                for idx in range(len(subscriber_data.non_3gpp.apn_config)):
                    res = self.conn.execute(
                        "SELECT data FROM apndb WHERE " "apn_name = ?",
                        (
                            subscriber_data.non_3gpp.apn_config[
                                idx
                            ].service_selection,
                        ),
                    )

                    row = res.fetchone()
                    if not row:
                        raise ApnNotFoundError()

            self.conn.execute("INSERT INTO subscriberdb(subscriber_id, data) "
                              "VALUES (?, ?)", (sid, data_str))
        self._on_ready.add_subscriber(subscriber_data)

    @contextmanager
    def edit_subscriber(self, subscriber_id):
        """
        Context manager to modify the subscriber data.
        """
        with self.conn:
            res = self.conn.execute("SELECT data FROM subscriberdb WHERE "
                                    "subscriber_id = ?", (subscriber_id, ))
            row = res.fetchone()
            if not row:
                raise SubscriberNotFoundError(subscriber_id)
            subscriber_data = SubscriberData()
            subscriber_data.ParseFromString(row[0])
            yield subscriber_data
            data_str = subscriber_data.SerializeToString()
            self.conn.execute("UPDATE subscriberdb SET data = ? "
                              "WHERE subscriber_id = ?",
                              (data_str, subscriber_id))

    @contextmanager
    def edit_subscriber_apn(self, sid, request):
        """
        Context manager to modify the apn data for a subscriber.
        """
        with self.conn:
            res = self.conn.execute(
                "SELECT data FROM subscriberdb WHERE " "subscriber_id = ?",
                (sid,),
            )
            row = res.fetchone()
            if not row:
                raise SubscriberNotFoundError(sid)
            non_3gpp = Non3GPPUserProfile()
            sub_data = self.get_subscriber_data(sid)
            num_apn = len(request.non_3gpp.apn_config)
            for idx in range(num_apn):
                res = self.conn.execute(
                    "SELECT data FROM apndb WHERE " "apn_name = ?",
                    (request.non_3gpp.apn_config[idx].service_selection,),
                )

                row = res.fetchone()
                if not row:
                    raise ApnNotFoundError()

                apn_config = non_3gpp.apn_config.add()
                apn_config.service_selection = request.non_3gpp.apn_config[
                    idx
                ].service_selection
            # Re-populate subscriber data with the APN parameters
            new_sub_data = SubscriberData(
                sid=SIDUtils.to_pb(sid),
                gsm=sub_data.gsm,
                lte=sub_data.lte,
                state=sub_data.state,
                non_3gpp=non_3gpp,
            )
            data_str = new_sub_data.SerializeToString()
            self.conn.execute(
                "UPDATE subscriberdb SET data = ? " "WHERE subscriber_id = ?",
                (data_str, sid),
            )

    def delete_subscriber(self, subscriber_id):
        """
        Method that deletes a subscriber, if present.
        """
        with self.conn:
            self.conn.execute("DELETE FROM subscriberdb WHERE "
                              "subscriber_id = ?", (subscriber_id, ))

    def delete_all_subscribers(self):
        """
        Method that removes all the subscribers from the store
        """
        with self.conn:
            self.conn.execute("DELETE FROM subscriberdb")

    def get_subscriber_data(self, subscriber_id):
        """
        Method that returns the auth key for the subscriber.
        """
        with self.conn:
            res = self.conn.execute("SELECT data FROM subscriberdb WHERE "
                                    "subscriber_id = ?", (subscriber_id, ))
            row = res.fetchone()
            if not row:
                raise SubscriberNotFoundError(subscriber_id)
        subscriber_data = SubscriberData()
        subscriber_data.ParseFromString(row[0])
        return subscriber_data

    def list_subscribers(self):
        """
        Method that returns the list of subscribers stored
        """
        with self.conn:
            res = self.conn.execute("SELECT subscriber_id FROM subscriberdb")
            return [row[0] for row in res]

    def list_apns(self):
        """
        Method that returns the list of apns stored
        """
        with self.conn:
            res = self.conn.execute("SELECT apn_name FROM apndb")
            return [row[0] for row in res]

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
        with self.conn:
            res = self.conn.execute("UPDATE subscriberdb SET data = ? "
                                    "WHERE subscriber_id = ?", (data_str, sid))
            if not res.rowcount:
                raise SubscriberNotFoundError(sid)

    def resync(self, subscribers):
        """
        Method that should resync the store with the mentioned list of
        subscribers. The resync leaves the current state of subscribers
        intact.

        Args:
            subscribers - list of subscribers to be in the store.
        """
        with self.conn:
            # Capture the current state of the subscribers
            res = self.conn.execute("SELECT subscriber_id, data FROM subscriberdb")
            current_state = {}
            for row in res:
                sub = SubscriberData()
                sub.ParseFromString(row[1])
                current_state[row[0]] = sub.state

            # Clear all subscribers
            self.conn.execute("DELETE FROM subscriberdb")

            # Add the subscribers with the current state
            for sub in subscribers:
                sid = SIDUtils.to_str(sub.sid)
                if sid in current_state:
                    sub.state.CopyFrom(current_state[sid])
                data_str = sub.SerializeToString()
                self.conn.execute("INSERT INTO subscriberdb(subscriber_id, data) "
                                  "VALUES (?, ?)", (sid, data_str))
        self._on_ready.resync(subscribers)

    def on_ready(self):
        return self._on_ready.event.wait()

    def get_apn_config(self, apn_data):
        """
        Method that returns the APN
        """
        with self.conn:
            res = self.conn.execute(
                "SELECT data FROM apndb WHERE " "apn_name = ?",
                (apn_data.non_3gpp.apn_config[0].service_selection,),
            )
            row = res.fetchone()
            if not row:
                raise ApnNotFoundError()
        sub_data = SubscriberData()
        sub_data.ParseFromString(row[0])
        num_apn = len(sub_data.non_3gpp.apn_config)
        for idx in range(num_apn):
            if (
                sub_data.non_3gpp.apn_config[idx].service_selection
                == apn_data.non_3gpp.apn_config[0].service_selection
            ):
                return sub_data.non_3gpp.apn_config[idx]
        return sub_data.non_3gpp.apn_config[idx]

    def _populate_apn(self, apn_data, apn_config, idx):
        """
        Method that populates apn data.
        """
        apn_config.service_selection = apn_data.non_3gpp.apn_config[
            idx
        ].service_selection
        apn_config.qos_profile.class_id = apn_data.non_3gpp.apn_config[
            idx
        ].qos_profile.class_id
        apn_config.qos_profile.priority_level = apn_data.non_3gpp.apn_config[
            idx
        ].qos_profile.priority_level
        apn_config.qos_profile.preemption_capability = apn_data.non_3gpp.apn_config[
            idx
        ].qos_profile.preemption_capability
        apn_config.qos_profile.preemption_vulnerability = apn_data.non_3gpp.apn_config[
            idx
        ].qos_profile.preemption_vulnerability
        apn_config.ambr.max_bandwidth_ul = apn_data.non_3gpp.apn_config[
            idx
        ].ambr.max_bandwidth_ul
        apn_config.ambr.max_bandwidth_dl = apn_data.non_3gpp.apn_config[
            idx
        ].ambr.max_bandwidth_dl

    def _update_apn(self, sub_data, apn_data):
        """
        Method that populates apn data.
        """
        num_apn = len(sub_data.non_3gpp.apn_config)
        # Only one APN config will be received at a time.
        # Hence fetching from index 0
        for idx in range(num_apn):
            if (
                sub_data.non_3gpp.apn_config[idx].service_selection
                == apn_data.non_3gpp.apn_config[0].service_selection
            ):
                if apn_data.non_3gpp.apn_config[0].qos_profile.class_id:
                    sub_data.non_3gpp.apn_config[
                        idx
                    ].qos_profile.class_id = apn_data.non_3gpp.apn_config[
                        0
                    ].qos_profile.class_id
                if apn_data.non_3gpp.apn_config[0].qos_profile.priority_level:
                    sub_data.non_3gpp.apn_config[
                        idx
                    ].qos_profile.priority_level = apn_data.non_3gpp.apn_config[
                        0
                    ].qos_profile.priority_level
                # preemption_capability and preemption_vulnerability are bool
                # type and cannot be checked for non-zero. Hence they are
                # mandatory parameters
                sub_data.non_3gpp.apn_config[
                    idx
                ].qos_profile.preemption_capability = apn_data.non_3gpp.apn_config[
                    0
                ].qos_profile.preemption_capability
                sub_data.non_3gpp.apn_config[
                    idx
                ].qos_profile.preemption_vulnerability = apn_data.non_3gpp.apn_config[
                    0
                ].qos_profile.preemption_vulnerability
                if apn_data.non_3gpp.apn_config[0].ambr.max_bandwidth_ul:
                    sub_data.non_3gpp.apn_config[
                        idx
                    ].ambr.max_bandwidth_ul = apn_data.non_3gpp.apn_config[
                        0
                    ].ambr.max_bandwidth_ul
                if apn_data.non_3gpp.apn_config[0].ambr.max_bandwidth_dl:
                    sub_data.non_3gpp.apn_config[
                        idx
                    ].ambr.max_bandwidth_dl = apn_data.non_3gpp.apn_config[
                        0
                    ].ambr.max_bandwidth_dl

    def add_apn_config(self, apn_data):
        """
        Method that adds apn data.
        """
        num_new_apn = len(apn_data.non_3gpp.apn_config)
        non_3gpp = Non3GPPUserProfile()

        with self.conn:
            for idx in range(num_new_apn):
                apn_config = non_3gpp.apn_config.add()
                self._populate_apn(apn_data, apn_config, idx)

                sub = SubscriberData(non_3gpp=non_3gpp)
                data_str = sub.SerializeToString()

                res = self.conn.execute(
                    "SELECT data FROM apndb WHERE " "apn_name = ?",
                    (apn_data.non_3gpp.apn_config[idx].service_selection,),
                )
                if res.fetchone():
                    raise DuplicateApnError()

                self.conn.execute(
                    "INSERT INTO apndb(apn_name,data) " "VALUES (?, ?)",
                    (
                        apn_data.non_3gpp.apn_config[idx].service_selection,
                        data_str,
                    ),
                )

    def delete_apn_config(self, apn_data):
        """
        Method that deletes an apn, if present.
        """
        num_apn = 0
        with self.conn:
            res = self.conn.execute(
                "SELECT data FROM apndb WHERE " "apn_name = ?",
                (apn_data.non_3gpp.apn_config[0].service_selection,),
            )
            row = res.fetchone()
            if not row:
                raise ApnNotFoundError()
            self.conn.execute(
                "DELETE FROM apndb WHERE " "apn_name = ?",
                (apn_data.non_3gpp.apn_config[0].service_selection,),
            )

        # Delete the corresponding APN from subscriberdb table
        res = self.conn.execute("SELECT subscriber_id, data FROM subscriberdb")
        for row in res:
            sub = SubscriberData()
            sub.ParseFromString(row[1])
            if sub.non_3gpp.apn_config:
                num_apn = len(sub.non_3gpp.apn_config)
            for idx in range(num_apn):
                sid = SIDUtils.to_str(sub.sid)
                if (
                    sub.non_3gpp.apn_config[idx].service_selection
                    == apn_data.non_3gpp.apn_config[0].service_selection
                ):
                    del sub.non_3gpp.apn_config[idx]
                    data_str = sub.SerializeToString()
                    # Repopulate subscriberdb after deleting.
                    with self.conn:
                        res = self.conn.execute(
                            "UPDATE subscriberdb SET data = ? "
                            "WHERE subscriber_id = ?",
                            (data_str, sid),
                        )
                    break

    def edit_apn_config(self, apn_data):
        """
        Context manager to modify the APN data.
        """

        with self.conn:
            res = self.conn.execute(
                "SELECT data FROM apndb WHERE " "apn_name = ?",
                (apn_data.non_3gpp.apn_config[0].service_selection,),
            )
            row = res.fetchone()
            if not row:
                raise ApnNotFoundError()

        sub_data = SubscriberData()
        sub_data.ParseFromString(row[0])
        self._update_apn(sub_data, apn_data)
        data_str = sub_data.SerializeToString()
        with self.conn:
            res = self.conn.execute(
                "UPDATE apndb SET data = ? " "WHERE apn_name = ?",
                (data_str, apn_data.non_3gpp.apn_config[0].service_selection),
            )
            if not res.rowcount:
                raise ApnNotFoundError()
