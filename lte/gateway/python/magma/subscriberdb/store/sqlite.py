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

from lte.protos.subscriberdb_pb2 import (
    SubscriberData,
    Non3GPPUserProfile,
    SubscriberIDSet,
    APNConfiguration,
)

from magma.subscriberdb.sid import SIDUtils
from .base import (
    BaseStore,
    DuplicateSubscriberError,
    SubscriberNotFoundError,
    ApnNotFoundError,
    DuplicateApnError,
)
from .onready import OnDataReady
import logging


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
                "(apn_name text PRIMARY KEY, data text, subscriberids text)"
            )

    def _add_sub_to_apndb(self, subscriber_data):
        """
        Method that adds the subscriber ids to apndb.
        """
        subscriber_id = SIDUtils.to_str(subscriber_data.sid)
        for apn in subscriber_data.non_3gpp.apn_config:
            duplicate_id = False
            # Add the subscriber ID in apnDB
            sids = SubscriberIDSet()
            res = self.conn.execute(
                "SELECT subscriberids FROM apndb WHERE " "apn_name = ?",
                (apn.service_selection,),
            )
            row = res.fetchone()
            if not row:
                raise ApnNotFoundError()
            stored_sid_data = SubscriberIDSet()
            # If sids are already present for the apn, add them first
            # and then add the new sids, else add only the new sids
            if row[0]:
                stored_sid_data.ParseFromString(row[0])

                # Repopulate sids
                # Add the sids already present in stored_sid_data
                for sid in stored_sid_data.sids:
                    sid_data = sids.sids.add()
                    sid_data.id = sid.id
                # Now add the new sids received from test script/cli
                for sid_idx in stored_sid_data.sids:
                    if sid_idx.id == subscriber_id:
                        duplicate_id = True
                        break
                if not duplicate_id:
                    sid_data = sids.sids.add()
                    sid_data.id = subscriber_id
            else:
                sid_data = sids.sids.add()
                sid_data.id = subscriber_id
            data = sids
            data_str = data.SerializeToString()
            res = self.conn.execute(
                "UPDATE apndb SET subscriberids = ? " "WHERE apn_name = ?",
                (data_str, apn.service_selection,),
            )

    def _fill_apn_from_apndb(self, non_3gpp, subscriber_data):
        """
        Method that retrieves apn data from apndb and updates the subscriberdb.
        """
        # Fetch the APN data from apndb and add to subscriberdb
        for apn in subscriber_data.non_3gpp.apn_config:
            res = self.conn.execute(
                "SELECT data FROM apndb WHERE " "apn_name = ?",
                (apn.service_selection,),
            )
            row = res.fetchone()
            if not row:
                raise ApnNotFoundError()

            apn_data = APNConfiguration()
            apn_data.ParseFromString(row[0])
            apn_config = non_3gpp.apn_config.add()
            self._populate_apn(apn_data, apn_config)

    def add_subscriber(self, subscriber_data):
        """
        Method that adds the subscriber.
        """
        sid = SIDUtils.to_str(subscriber_data.sid)
        with self.conn:
            res = self.conn.execute(
                "SELECT data FROM subscriberdb WHERE " "subscriber_id = ?",
                (sid,),
            )
            if res.fetchone():
                raise DuplicateSubscriberError(sid)

            if subscriber_data.non_3gpp.apn_config:
                non_3gpp = Non3GPPUserProfile()
                self._fill_apn_from_apndb(non_3gpp, subscriber_data)

                new_sub_data = SubscriberData(
                    sid=SIDUtils.to_pb(sid),
                    gsm=subscriber_data.gsm,
                    lte=subscriber_data.lte,
                    state=subscriber_data.state,
                    non_3gpp=non_3gpp,
                )
                data_str = new_sub_data.SerializeToString()
                # Add the sid to apndb
                self._add_sub_to_apndb(subscriber_data)
            else:
                data_str = subscriber_data.SerializeToString()
            self.conn.execute(
                "INSERT INTO subscriberdb(subscriber_id, data) "
                "VALUES (?, ?)",
                (sid, data_str),
            )
        self._on_ready.add_subscriber(subscriber_data)

    @contextmanager
    def edit_subscriber(self, subscriber_id):
        """
        Context manager to modify the subscriber data.
        """
        with self.conn:
            res = self.conn.execute(
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
            self.conn.execute(
                "UPDATE subscriberdb SET data = ? " "WHERE subscriber_id = ?",
                (data_str, subscriber_id),
            )

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

            # Fetch apn data from apndb and update subscriberdb
            self._fill_apn_from_apndb(non_3gpp, request)
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
            # Add sids to apndb
            self._add_sub_to_apndb(request)

    def _delete_sid_from_apndb(self, apn, sid):
        """
        Method that deletes subscriberId from apndb.
        """
        with self.conn:
            res = self.conn.execute(
                "SELECT subscriberids FROM apndb WHERE " "apn_name = ?",
                (apn,),
            )
            row = res.fetchone()
            if not row:
                raise ApnNotFoundError()
            stored_sid_data = SubscriberIDSet()
            if row[0]:
                stored_sid_data.ParseFromString(row[0])
                for idx in range(len(stored_sid_data.sids)):
                    if sid == stored_sid_data.sids[idx].id:
                        del stored_sid_data.sids[idx]
                        # Repopulate the db after deletion
                        data_str = stored_sid_data.SerializeToString()
                        res = self.conn.execute(
                            "UPDATE apndb SET subscriberids = ? "
                            "WHERE apn_name = ?",
                            (data_str, apn,),
                        )
                        break

    @contextmanager
    def delete_subscriber_apn(self, sid, request):
        """
        Context manager to delete the apn data of a subscriber.
        """
        with self.conn:
            res = self.conn.execute(
                "SELECT data FROM subscriberdb WHERE " "subscriber_id = ?",
                (sid,),
            )
            row = res.fetchone()
            if not row:
                raise SubscriberNotFoundError(sid)
            sub_data = self.get_subscriber_data(sid)
            num_stored_apn = len(sub_data.non_3gpp.apn_config)
            for rcvd_apn in request.non_3gpp.apn_config:
                for idx in range(num_stored_apn):
                    if (
                        sub_data.non_3gpp.apn_config[idx].service_selection
                        == rcvd_apn.service_selection
                    ):
                        del sub_data.non_3gpp.apn_config[idx]
                        # Delete the sid entry from apndb
                        try:
                            self._delete_sid_from_apndb(
                                rcvd_apn.service_selection, sid,
                            )

                        except ApnNotFoundError as e:
                            logging.warning(
                                "APN not found : %s %s",
                                e,
                                rcvd_apn.service_selection,
                            )

                        break
            data_str = sub_data.SerializeToString()
            # Re-populate subscriber data with the APN parameters
            self.conn.execute(
                "UPDATE subscriberdb SET data = ? " "WHERE subscriber_id = ?",
                (data_str, sid),
            )

    def delete_subscriber(self, subscriber_id):
        """
        Method that deletes a subscriber, if present.
        """
        with self.conn:
            # Delete the subscriber id from the apndb
            try:
                sub_data = self.get_subscriber_data(subscriber_id)
            except SubscriberNotFoundError:
                raise SubscriberNotFoundError()
            for apn in sub_data.non_3gpp.apn_config:
                try:
                    self._delete_sid_from_apndb(
                        apn.service_selection, subscriber_id,
                    )
                except ApnNotFoundError:
                    continue
            self.conn.execute(
                "DELETE FROM subscriberdb WHERE " "subscriber_id = ?",
                (subscriber_id,),
            )

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
                (apn_data.service_selection,),
            )
            row = res.fetchone()
            if not row:
                raise ApnNotFoundError()
        apn_config = APNConfiguration()
        apn_config.ParseFromString(row[0])
        return apn_config

    def _populate_apn(self, apn_data, apn_config):
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

    def _update_apn(self, apn_config, apn_data):
        """
        Method that populates apn data.
        """
        if apn_data.qos_profile.class_id:
            apn_config.qos_profile.class_id = apn_data.qos_profile.class_id
        if apn_data.qos_profile.priority_level:
            apn_config.qos_profile.priority_level = (
                apn_data.qos_profile.priority_level
            )
        # preemption_capability and preemption_vulnerability are bool
        # type and cannot be checked for non-zero. Hence they are
        # mandatory parameters
        apn_config.qos_profile.preemption_capability = (
            apn_data.qos_profile.preemption_capability
        )
        apn_config.qos_profile.preemption_vulnerability = (
            apn_data.qos_profile.preemption_vulnerability
        )
        if apn_data.ambr.max_bandwidth_ul:
            apn_config.ambr.max_bandwidth_ul = apn_data.ambr.max_bandwidth_ul
        if apn_data.ambr.max_bandwidth_dl:
            apn_config.ambr.max_bandwidth_dl = apn_data.ambr.max_bandwidth_dl

    def add_apn_config(self, apn_data):
        """
        Method that adds apn data.
        """
        with self.conn:
            res = self.conn.execute(
                "SELECT data FROM apndb WHERE " "apn_name = ?",
                (apn_data.service_selection,),
            )
            if res.fetchone():
                raise DuplicateApnError()

            data_str = apn_data.SerializeToString()

            self.conn.execute(
                "INSERT INTO apndb(apn_name,data,subscriberids) "
                "VALUES (?, ?, ?)",
                (apn_data.service_selection, data_str, None,),
            )

    def delete_apn_config(self, apn_data):
        """
        Method that deletes an apn, if present.
        """
        num_apn = 0

        with self.conn:
            # Check if APN is valid
            res = self.conn.execute(
                "SELECT data FROM apndb WHERE " "apn_name = ?",
                (apn_data.service_selection,),
            )
            row = res.fetchone()
            if not row:
                raise ApnNotFoundError()
            # First delete the APN entry from subscriberdb table
            res = self.conn.execute(
                "SELECT subscriberids FROM apndb WHERE " "apn_name = ?",
                (apn_data.service_selection,),
            )
            row = res.fetchone()
            if not row:
                raise ApnNotFoundError()
            stored_sid_data = SubscriberIDSet()
            if row[0]:
                stored_sid_data.ParseFromString(row[0])
                for sid in stored_sid_data.sids:
                    try:
                        sub_data = self.get_subscriber_data(sid.id)
                    except SubscriberNotFoundError as e:
                        logging.warning(
                            "Subscriber not found for apn: %s %s",
                            e,
                            apn_data.service_selection,
                        )
                        continue
                    num_apn = len(sub_data.non_3gpp.apn_config)
                    for apn_idx in range(num_apn):
                        if (
                            sub_data.non_3gpp.apn_config[
                                apn_idx
                            ].service_selection
                            == apn_data.service_selection
                        ):
                            # Delete the apn entry in subscriberdb
                            del sub_data.non_3gpp.apn_config[apn_idx]
                            data_str = sub_data.SerializeToString()
                            # Repopulate subscriberdb after deleting.
                            res = self.conn.execute(
                                "UPDATE subscriberdb SET data = ? "
                                "WHERE subscriber_id = ?",
                                (data_str, SIDUtils.to_str(sub_data.sid)),
                            )
                            break
                # Delete all the sids  for the apn in apndb
                del stored_sid_data.sids[:]

            # Now delete the APN entry
            self.conn.execute(
                "DELETE FROM apndb WHERE " "apn_name = ?",
                (apn_data.service_selection,),
            )

    def edit_apn_config(self, apn_data):
        """
        Context manager to modify the APN data.
        """

        with self.conn:
            res = self.conn.execute(
                "SELECT data FROM apndb WHERE " "apn_name = ?",
                (apn_data.service_selection,),
            )
            row = res.fetchone()
            if not row:
                raise ApnNotFoundError()
        # Update the apn data in apndb
        apn_config = APNConfiguration()
        apn_config.ParseFromString(row[0])
        self._update_apn(apn_config, apn_data)
        data_str = apn_config.SerializeToString()
        with self.conn:
            res = self.conn.execute(
                "UPDATE apndb SET data = ? " "WHERE apn_name = ?",
                (data_str, apn_data.service_selection),
            )
            if not res.rowcount:
                raise ApnNotFoundError()
            # Update the apn data for the corresponding subscriber
            res = self.conn.execute(
                "SELECT subscriberids FROM apndb WHERE " "apn_name = ?",
                (apn_data.service_selection,),
            )
            row = res.fetchone()
            if not row:
                raise ApnNotFoundError()
            stored_sid_data = SubscriberIDSet()
            if row[0]:
                stored_sid_data.ParseFromString(row[0])
                for sid in stored_sid_data.sids:
                    subs_data = self.get_subscriber_data(sid.id)
                    for apn in subs_data.non_3gpp.apn_config:
                        if apn.service_selection == apn_data.service_selection:
                            self._update_apn(apn, apn_data)
                            break
                    data_str = subs_data.SerializeToString()
                    # Repopulate subscriberdb after updating.
                    res = self.conn.execute(
                        "UPDATE subscriberdb SET data = ? "
                        "WHERE subscriber_id = ?",
                        (data_str, SIDUtils.to_str(subs_data.sid)),
                    )

    def list_sids_for_apn(self, apn_data):
        """
        Method that returns the sids for the given APN
        """
        with self.conn:
            res = self.conn.execute(
                "SELECT subscriberids FROM apndb WHERE " "apn_name = ?",
                (apn_data.service_selection,),
            )
            row = res.fetchone()
            if not row:
                raise ApnNotFoundError()

            stored_sid_data = SubscriberIDSet()
            if row[0]:
                stored_sid_data.ParseFromString(row[0])
                return stored_sid_data
            return SubscriberIDSet()
