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
import subprocess

#import swagger_client
from orc8r.protos.common_pb2 import Void
from lte.protos.subscriberdb_pb2 import LTESubscription, SubscriberData, \
    SubscriberState, SubscriberID
from lte.protos.subscriberdb_pb2_grpc import SubscriberDBStub

#from integ_tests.cloud.fixtures import GATEWAY_ID, NETWORK_ID
from integ_tests.gateway.rpc import get_gateway_hw_id, get_rpc_channel
from magma.subscriberdb.sid import SIDUtils

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
            lambda: self._subscriber_stub.AddSubscriber(sub_data))
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

    def clean_up(self):
        # Remove all sids
        for sid in self.list_subscriber_sids():
            self.delete_subscriber(sid)
        assert not self.list_subscriber_sids()
        assert not self._added_sids

    def wait_for_changes(self):
        # On gateway, changes propagate immediately
        return

class SubscriberDbCassandra(SubscriberDbClient):
    """
    Handle subscriber action by making calls to Cassandra database of OAI HSS
    """
    HSS_IP = '192.168.60.153'
    HSS_USER = 'vagrant'
    IDENTITY_FILE = '$HOME/.ssh/id_rsa'
    CASSANDRA_SERVER_IP = '127.0.0.1'

    def __init__(self):
        self._added_sids = set()
        print("*********Init SubscriberDbCassandra***********")
        add_mme_cmd = "$HOME/openair-cn/scripts/data_provisioning_mme --id 3 "\
            "--mme-identity mme.magma.com --realm magma.com "\
            "--ue-reachability 1 -C "+ self.CASSANDRA_SERVER_IP
        self._run_remote_cmd(add_mme_cmd)

    def _run_remote_cmd(self, cmd_str):
        ssh_args = "-o UserKnownHostsFile=/dev/null "\
            "-o StrictHostKeyChecking=no"
        ssh_cmd = "ssh -i {id_file} {args} {user}@{host} {cmd}".format(
            id_file=self.IDENTITY_FILE, args=ssh_args, user=self.HSS_USER,
            host=self.HSS_IP, cmd=cmd_str)
        output, error = subprocess.Popen(ssh_cmd, shell=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE).communicate()
        print("Output: ", output)
        print("Error: ", error)
        return output, error

    def add_subscriber(self, sid):
        sid = sid[4:]
        print("Adding subscriber", sid)
        # Insert into users
        add_usr_cmd = "$HOME/openair-cn/scripts/data_provisioning_users "\
            "--apn oai.ipv4 --apn2 internet --key " + KEY + \
            " --imsi-first " + sid + "--mme-identity mme.magma.com "\
            "--no-of-users 1 --realm magma.com "\
            "--cassandra-cluster " + self.CASSANDRA_SERVER_IP
        self._run_remote_cmd(add_usr_cmd)

    def delete_subscriber(self, sid):
        print("Removing subscriber", sid)

    def _delete_all_subscribers(self):
        print("Removing all subscribers")
        del_all_subs_cmd = "$HOME/openair-cn/scripts/data_provisioning_users "\
            "--verbose True --truncate True -n 0 "\
            "-C " + self.CASSANDRA_SERVER_IP
        del_mme_cmd = "$HOME/openair-cn/scripts/data_provisioning_mme "\
            "--verbose True --truncate True -C " + self.CASSANDRA_SERVER_IP

    def list_subscriber_sids(self):
        sids = []
        return sids

    def clean_up(self):
        self._delete_all_subscribers()

    def wait_for_changes(self):
        # On gateway, changes propagate immediately
        return
