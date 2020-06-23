"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging
import time

import abc
import base64
import grpc
from orc8r.protos.common_pb2 import Void
from lte.protos.subscriberdb_pb2 import (
    LTESubscription,
    SubscriberData,
    SubscriberState,
    SubscriberID,
    SubscriberUpdate,
)
from lte.protos.subscriberdb_pb2_grpc import SubscriberDBStub

from integ_tests.gateway.rpc import get_gateway_hw_id, get_rpc_channel
from magma.subscriberdb.sid import SIDUtils

try:
    import MySQLdb
    mysql_import_error = None
except ImportError as e:
    mysql_import_error = e

KEY = '000102030405060708090A0B0C0D0E0F'
RETRY_COUNT = 4
RETRY_INTERVAL = 1  # seconds


class S1apTimeoutError(Exception):
    """ Indicate that a test-related check has timed out. """
    pass


class SubscriberDbClient(metaclass=abc.ABCMeta):
    """ Interface for the Subscriber DB. """

    @abc.abstractmethod
    def add_subscriber(self, sid):
        """
        Add a subscriber to the EPC by :sid:.
        Args:
            sid (str): the SID to add
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def delete_subscriber(self, sid):
        """
        Delete a subscriber from the EPC by :sid:.
        Args:
            sid (str): the SID to delete
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def list_subscriber_sids(self):
        """
        List all stored subscribers. Is blocking.
        Returns:
            sids (str[]): list of subscriber SIDs
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def clean_up(self):
        """ Clean up, delete all subscribers. """
        raise NotImplementedError()

    @abc.abstractmethod
    def wait_for_changes(self):
        """
        Blocks until changes go through. This is really only implemented on
        the cloud side, where subscriber changes can take a while to propagate
        from cloud to gateway
        """
        raise NotImplementedError()


class SubscriberDbGrpc(SubscriberDbClient):
    """
    Handle subscriber actions by making calls over gRPC directly to the
    gateway.
    """

    def __init__(self):
        """ Init the gRPC stub.  """
        self._added_sids = set()
        self._subscriber_stub = SubscriberDBStub(
            get_rpc_channel("subscriberdb"))

    @staticmethod
    def _try_to_call(grpc_call):
        """ Attempt to call into SubscriberDB and retry if unavailable """
        for i in range(RETRY_COUNT):
            try:
                return grpc_call()
            except grpc.RpcError as error:
                err_code = error.exception().code()
                # If unavailable, try again
                if (err_code == grpc.StatusCode.UNAVAILABLE):
                    logging.warning("Subscriberdb unavailable, retrying...")
                    time.sleep(RETRY_INTERVAL * (2 ** i))
                    continue
                logging.error("Subscriberdb grpc call failed with error : %s",
                              error)
                raise

    @staticmethod
    def _get_subscriberdb_data(sid):
        """
        Get subscriber data in protobuf format.

        Args:
            sid (str): string representation of the subscriber id
        Returns:
            subscriber_data (protos.subscriberdb_pb2.SubscriberData):
                full subscriber information for :sid: in protobuf format.
        """
        sub_db_sid = SIDUtils.to_pb(sid)
        lte = LTESubscription()
        lte.state = LTESubscription.ACTIVE
        lte.auth_key = bytes.fromhex(KEY)
        state = SubscriberState()
        state.lte_auth_next_seq = 1
        return SubscriberData(sid=sub_db_sid, lte=lte, state=state)

    @staticmethod
    def _get_apn_data(sid, apn_list):
        """
        Get APN data in protobuf format.

        Args:
            apn_list : list of APN configuration
        Returns:
            update (protos.subscriberdb_pb2.SubscriberUpdate)
        """
        # APN
        update = SubscriberUpdate()
        update.data.sid.CopyFrom(sid)
        non_3gpp = update.data.non_3gpp
        for apn in apn_list:
            apn_config = non_3gpp.apn_config.add()
            apn_config.service_selection = apn["apn_name"]
            apn_config.qos_profile.class_id = apn["qci"]
            apn_config.qos_profile.priority_level = apn["priority"]
            apn_config.qos_profile.preemption_capability = apn["pre_cap"]
            apn_config.qos_profile.preemption_vulnerability = apn["pre_vul"]
            apn_config.ambr.max_bandwidth_ul = apn["mbr_ul"]
            apn_config.ambr.max_bandwidth_dl = apn["mbr_dl"]
            apn_config.pdn = apn["pdn_type"] if "pdn_type" in apn else 0
        return update

    def _check_invariants(self):
        """
        Assert preservation of invariants.

        Raises:
            AssertionError: when invariants do not hold
        """
        sids_eq_len = len(self._added_sids) == len(self.list_subscriber_sids())
        assert sids_eq_len

    def add_subscriber(self, sid):
        logging.info("Adding subscriber : %s", sid)
        self._added_sids.add(sid)
        sub_data = self._get_subscriberdb_data(sid)
        SubscriberDbGrpc._try_to_call(
            lambda: self._subscriber_stub.AddSubscriber(sub_data)
        )
        self._check_invariants()

    def delete_subscriber(self, sid):
        logging.info("Deleting subscriber : %s", sid)
        self._added_sids.discard(sid)
        sid_pb = SubscriberID(id=sid[4:])
        SubscriberDbGrpc._try_to_call(
            lambda: self._subscriber_stub.DeleteSubscriber(sid_pb))

    def list_subscriber_sids(self):
        sids_pb = SubscriberDbGrpc._try_to_call(
            lambda: self._subscriber_stub.ListSubscribers(Void()).sids)
        sids = ['IMSI' + sid.id for sid in sids_pb]
        return sids

    def config_apn_details(self, imsi, apn_list):
        sid = SIDUtils.to_pb(imsi)
        update_sub = self._get_apn_data(sid, apn_list)
        fields = update_sub.mask.paths
        fields.append('non_3gpp')
        SubscriberDbGrpc._try_to_call(
            lambda: self._subscriber_stub.UpdateSubscriber(update_sub)
        )

    def clean_up(self):
        # Remove all sids
        for sid in self.list_subscriber_sids():
            self.delete_subscriber(sid)
        assert not self.list_subscriber_sids()
        assert not self._added_sids

    def wait_for_changes(self):
        # On gateway, changes propagate immediately
        return


#class SubscriberDbRest(SubscriberDbClient):
#    """
#    Handle subscriber actions by making calls to the REST API endpoints.
#    """
#
#    POLL_CEILING = 12
#    POLL_INTERVAL_SECS = 10
#
#    def __init__(self, cloud_manager, network_id=NETWORK_ID,
#                 gateway_id=GATEWAY_ID, use_gateway_grpc=True):
#        """ Init the REST wrapper by getting a network ID. """
#        self._added_sids = set()
#        self._cloud_manager = cloud_manager
#        self._network_id = network_id
#        self._gateway_id = gateway_id
#        self._use_gateway_grpc = use_gateway_grpc
#        if self._use_gateway_grpc:
#            self._cloud_manager.create_network(self._network_id)
#            self._cloud_manager.register_gateway(
#                self._network_id, self._gateway_id, get_gateway_hw_id())
#            # Create a gRPC subscriber client to communicate directly with
#            # the gateway
#            self._subscriber_grpc = SubscriberDbGrpc()
#        else:
#            self._subscriber_grpc = None
#
#        for sid in self.list_subscriber_sids():
#            self._added_sids.add(sid)
#
#    def _check_invariants(self):
#        """
#        Assert preservation of invariants.
#
#        Raises:
#            AssertionError: when invariants do not hold
#        """
#        sids_eq_len = len(self._added_sids) == len(self.list_subscriber_sids())
#        assert sids_eq_len, "SIDs length does not match expectations"
#
#    def _wait_for_cloud_match_gateway_sids(self):
#        """ Ensure that the gateway sids and cloud sids match. """
#        loops = 0
#        if self._subscriber_grpc is None:
#            return
#        while True:
#            cloud_sids = self.list_subscriber_sids()
#            gateway_sids = self._subscriber_grpc.list_subscriber_sids()
#            print('Waiting for sids to match:')
#            print('\t> cloud sids:', cloud_sids)
#            print('\t> gateway sids:', gateway_sids)
#            if set(cloud_sids) == set(gateway_sids):
#                print('Cloud and gateway subscribers match.')
#                break
#            loops += 1
#            if loops > self.POLL_CEILING:
#                raise S1apTimeoutError()
#            # TODO: add backoff
#            time.sleep(self.POLL_INTERVAL_SECS)
#
#    @staticmethod
#    def _get_subscriberdb_data(sid):
#        """
#        Get subscriber data in swagger-codegen compatible format.
#
#        Args:
#            sid (str): string representation of the subscriber id
#        Returns:
#            subscriber_data (swagger_client.Subscriber):
#                full subscriber information for :sid:, formatted for
#                use with our cloud API as generated by swagger-codegen
#        """
#        # The auth key is stored here as a string of hex values. The cloud
#        # expects the auth key as a string of base64-encoded values.
#        auth_key_bytes = base64.b64encode(bytes.fromhex(KEY))
#        auth_key = str(auth_key_bytes, 'utf-8')
#        lte = swagger_client.LteSubscription(
#            state='ACTIVE', auth_key=auth_key)
#        subscriber = swagger_client.Subscriber(id=sid, lte=lte)
#        return subscriber
#
#    def add_subscriber(self, sid):
#        logging.info("Adding subscriber : %s", sid)
#        self._added_sids.add(sid)
#        subscriber = self._get_subscriberdb_data(sid)
#        self._cloud_manager \
#            .subscribers_api \
#            .networks_network_id_subscribers_post(
#                network_id=self._network_id, subscriber=subscriber)
#
#        self._check_invariants()
#
#    def delete_subscriber(self, sid):
#        logging.info("Deleting subscriber : %s", sid)
#        self._added_sids.discard(sid)
#        self._cloud_manager \
#            .subscribers_api \
#            .networks_network_id_subscribers_subscriber_id_delete(
#                network_id=self._network_id, subscriber_id=sid)
#
#    def list_subscriber_sids(self):
#        sids = self._cloud_manager \
#            .subscribers_api \
#            .networks_network_id_subscribers_get(
#                self._network_id)
#        return sids
#
#    def clean_up(self):
#        # Delete all sids
#        for sid in self.list_subscriber_sids():
#            self.delete_subscriber(sid)
#        # Check sids-related invariants
#        self._check_invariants()
#        assert not self._added_sids
#
#    def wait_for_changes(self):
#        self._wait_for_cloud_match_gateway_sids()
#

def dbtransaction(func):
    """ Wrapper to rollback db changes if there's an exception """
    def wrapper(self, *args, **kwargs):
        try:
            func(self, *args, **kwargs)
            self._db.commit()
        except Exception as e:
            self._db.rollback()
            raise e
    return wrapper
