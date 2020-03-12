"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import ctypes
import ipaddress
import logging
import os
import threading
from queue import Queue

import s1ap_types
from integ_tests.gateway.rpc import get_rpc_channel
from lte.protos.policydb_pb2 import (
    FlowDescription,
    FlowMatch,
    FlowQos,
    PolicyRule,
    QosArp,
)
from lte.protos.spgw_service_pb2 import CreateBearerRequest, DeleteBearerRequest
from lte.protos.spgw_service_pb2_grpc import SpgwServiceStub
from magma.subscriberdb.sid import SIDUtils


class S1ApUtil(object):
    """
    Helper class to wrap the initialization and API interface of S1APTester
    Note that some of the values that are not that interesting are set
    through config files, that this class doesn't override. Examples include
    the various interface timeout params.
    """

    # Extracted from TestCntlrApp/src/ueApp/ue_esm.h
    CM_ESM_PDN_IPV4 = 0b01
    CM_ESM_PDN_IPV6 = 0b10
    CM_ESM_PDN_IPV4V6 = 0b11

    lib_name = "libtfw.so"

    _cond = threading.Condition()
    _msg = Queue()

    class Msg(object):
        def __init__(self, msg_type, msg_p, msg_len):
            self.msg_type = msg_type
            self.msg_p = ctypes.create_string_buffer(msg_len)
            ctypes.memmove(self.msg_p, msg_p, msg_len)
            self.msg_len = msg_len

        def cast(self, msg_class):
            return ctypes.cast(self.msg_p, ctypes.POINTER(msg_class)).contents

    @staticmethod
    def s1ap_callback(msg_type, msg_p, msg_len):
        """ S1ap tester compatible callback"""
        with S1ApUtil._cond:
            S1ApUtil._msg.put(S1ApUtil.Msg(msg_type, msg_p, msg_len))
            S1ApUtil._cond.notify_all()

    def __init__(self):
        """
        Initialize the s1aplibrary and its callbacks.
        """
        lib_path = os.environ["S1AP_TESTER_ROOT"]
        lib = os.path.join(lib_path, "bin", S1ApUtil.lib_name)
        os.chdir(lib_path)
        self._test_lib = ctypes.cdll.LoadLibrary(lib)
        self._callback_type = ctypes.CFUNCTYPE(
            None, ctypes.c_short, ctypes.c_void_p, ctypes.c_short
        )
        # Maintain a reference to the function object so GC doesn't release it.
        self._callback_fn = self._callback_type(S1ApUtil.s1ap_callback)
        self._test_lib.initTestFrameWork(self._callback_fn)
        self._test_api = self._test_lib.tfwApi
        self._test_api.restype = ctypes.c_int16
        self._test_api.argtypes = [ctypes.c_uint16, ctypes.c_void_p]

        # Mutex for state change operations
        self._lock = threading.RLock()

        # Maintain a map of UE IDs to IPs
        self._ue_ip_map = {}

    def cleanup(self):
        """
        Cleanup the dll loaded explicitly so the next run doesn't reuse the
        same globals as ctypes LoadLibrary uses dlopen under the covers

        Also clear out the UE ID: IP mappings
        """
        # self._test_lib.dlclose(self._test_lib._handle)
        self._test_lib = None
        self._ue_ip_map = {}

    def issue_cmd(self, cmd_type, req):
        """
        Issue a command to the s1aptester and blocks until response is recvd.
        Args:
            cmd_type: The cmd type enum
            req: The request Structure
        Returns:
            None
        """
        c_req = None
        if req:
            # For non NULL requests obtain the address.
            c_req = ctypes.byref(req)
        with self._cond:
            rc = self._test_api(cmd_type.value, c_req)
            if rc:
                logging.error("Error executing command %s" % repr(cmd_type))
                return rc
        return 0

    def get_ip(self, ue_id):
        """ Returns the IP assigned to a given UE ID

        Args:
            ue_id: the ue_id to query

        Returns an ipaddress.ip_address for the given UE ID, or None if no IP
        has been observed to be assigned to this IP
        """
        with self._lock:
            if ue_id in self._ue_ip_map:
                return self._ue_ip_map[ue_id]
            return None

    def get_response(self):
        # Wait until callback is invoked.
        return self._msg.get(True)

    def attach(
        self,
        ue_id,
        attach_type,
        resp_type,
        resp_msg_type,
        sec_ctxt=s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT,
        id_type=s1ap_types.TFW_MID_TYPE_IMSI,
        eps_type=s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH,
    ):
        """
        Given a UE issue the attach request of specified type

        Caches the assigned IP address, if any is assigned

        Args:
            ue_id: The eNB ue_id
            attach_type: The type of attach e.g. UE_END_TO_END_ATTACH_REQUEST
            resp_type: enum type of the expected response
            sec_ctxt: Optional param allows for the reuse of the security
                context, defaults to creating a new security context.
            id_type: Optional param allows for changing up the ID type,
                defaults to s1ap_types.TFW_MID_TYPE_IMSI.
            eps_type: Optional param allows for variation in the EPS attach
                type, defaults to s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH.
        """
        attach_req = s1ap_types.ueAttachRequest_t()
        attach_req.ue_Id = ue_id
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt

        assert self.issue_cmd(attach_type, attach_req) == 0
        response = self.get_response()

        # The MME actually sends INT_CTX_SETUP_IND and UE_ATTACH_ACCEPT_IND in
        # one message, but the s1aptester splits it and sends the tests 2
        # messages. Usually context setup comes before attach accept, but
        # it's possible it may happen the other way
        if s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value == response.msg_type:
            response = self.get_response()
        elif s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND.value == response.msg_type:
            context_setup = self.get_response()
            assert context_setup.msg_type == s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value

        logging.debug(
            "s1ap response expected, received: %d, %d",
            resp_type.value,
            response.msg_type,
        )
        assert resp_type.value == response.msg_type

        msg = response.cast(resp_msg_type)

        # We only support IPv4 right now, as max PDN address in S1AP tester is
        # currently 13 bytes, which is too short for IPv6 (which requires 16)
        if resp_msg_type == s1ap_types.ueAttachAccept_t:
            pdn_type = msg.esmInfo.pAddr.pdnType
            addr = msg.esmInfo.pAddr.addrInfo
            if S1ApUtil.CM_ESM_PDN_IPV4 == pdn_type:
                # Cast and cache the IPv4 address
                ip = ipaddress.ip_address(bytes(addr[:4]))
                with self._lock:
                    self._ue_ip_map[ue_id] = ip
            else:
                raise ValueError("PDN TYPE %s not supported" % pdn_type)
        return msg

    def receive_emm_info(self):
        response = self.get_response()
        logging.debug(
            "s1ap message expected, received: %d, %d",
            s1ap_types.tfwCmd.UE_EMM_INFORMATION.value,
            response.msg_type,
        )
        assert response.msg_type == s1ap_types.tfwCmd.UE_EMM_INFORMATION.value

    def detach(self, ue_id, reason_type, wait_for_s1_ctxt_release=True):
        """ Given a UE issue a detach request """
        detach_req = s1ap_types.uedetachReq_t()
        detach_req.ue_Id = ue_id
        detach_req.ueDetType = reason_type
        assert self.issue_cmd(s1ap_types.tfwCmd.UE_DETACH_REQUEST, detach_req) == 0
        if reason_type == s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value:
            response = self.get_response()
            assert s1ap_types.tfwCmd.UE_DETACH_ACCEPT_IND.value == response.msg_type

        # Now wait for the context release response
        if wait_for_s1_ctxt_release:
            response = self.get_response()
            assert s1ap_types.tfwCmd.UE_CTX_REL_IND.value == response.msg_type

        with self._lock:
            del self._ue_ip_map[ue_id]


class SubscriberUtil(object):
    """
    Helper class to manage subscriber data for the tests.
    """

    SID_PREFIX = "IMSI00101"
    IMSI_LEN = 15

    def __init__(self, subscriber_client):
        """
        Initialize subscriber util.

        Args:
            subscriber_client (subscriber_db_client.SubscriberDbClient):
                client interacting with our subscriber APIs
        """
        self._sid_idx = 1
        self._ue_id = 1
        # Maintain references to UE configs to prevent GC
        self._ue_cfgs = []

        self._subscriber_client = subscriber_client

    def _gen_next_sid(self):
        """
        Generate the sid based on index offset and prefix
        """
        idx = str(self._sid_idx)
        # Find the 0 padding we need to add
        padding = self.IMSI_LEN - len(idx) - len(self.SID_PREFIX[4:])
        sid = self.SID_PREFIX + "0" * padding + idx
        self._sid_idx += 1
        print("Using subscriber IMSI %s" % sid)
        return sid

    def _get_s1ap_sub(self, sid):
        """
        Get the subscriber data in s1aptester format.
        Args:
            The string representation of the subscriber id
        """
        ue_cfg = s1ap_types.ueConfig_t()
        ue_cfg.ue_id = self._ue_id
        ue_cfg.auth_key = 1
        # Some s1ap silliness, the char field is modelled as an int and then
        # cast into a uint8.
        for i in range(0, 15):
            ue_cfg.imsi[i] = ctypes.c_ubyte(int(sid[4 + i]))
            ue_cfg.imei[i] = ctypes.c_ubyte(int("1"))
        ue_cfg.imei[15] = ctypes.c_ubyte(int("1"))
        ue_cfg.imsiLen = self.IMSI_LEN
        self._ue_cfgs.append(ue_cfg)
        self._ue_id += 1
        return ue_cfg

    def add_sub(self, num_ues=1):
        """ Add subscribers to the EPC, is blocking """
        # Add the default IMSI used for the tests
        subscribers = []
        for _ in range(num_ues):
            sid = self._gen_next_sid()
            self._subscriber_client.add_subscriber(sid)
            subscribers.append(self._get_s1ap_sub(sid))
        self._subscriber_client.wait_for_changes()
        return subscribers

    def cleanup(self):
        """ Cleanup added subscriber from subscriberdb """
        self._subscriber_client.clean_up()
        # block until changes propagate
        self._subscriber_client.wait_for_changes()


class MagmadUtil(object):
    def __init__(self, magmad_client):
        """
        Init magmad util.

        Args:
            magmad_client: MagmadServiceClient
        """
        self._magmad_client = magmad_client

    def restart_services(self, services):
        """
        Restart a list of magmad services.

        Args:
            services: List of (str) services names

        """
        self._magmad_client.restart_services(services)


class MobilityUtil(object):
    """ Utility wrapper for interacting with mobilityd """

    def __init__(self, mobility_client):
        """
        Initialize mobility util.

        Args:
            mobility_client (mobility_service_client.MobilityServiceClient):
                client interacting with our mobility APIs
        """
        self._mobility_client = mobility_client

    def add_ip_block(self, ip_block):
        """ Add an ip block

        Args:
            ip_block (str | ipaddress.ip_network): the IP block to add
        """
        ip_network_block = ipaddress.ip_network(ip_block)
        self._mobility_client.add_ip_block(ip_network_block)

    def remove_all_ip_blocks(self):
        """ Delete all allocated IP blocks. """
        self._mobility_client.remove_all_ip_blocks()

    def get_subscriber_table(self):
        """ Retrieve subscriber table from mobilityd """
        table = self._mobility_client.get_subscriber_ip_table()
        return table

    def list_ip_blocks(self):
        """ List all IP blocks in mobilityd """
        blocks = self._mobility_client.list_added_blocks()
        return blocks

    def remove_ip_blocks(self, blocks):
        """ Attempt to remove the given blocks from mobilityd

        Args:
            blocks (tuple(ip_network)): tuple of ipaddress.ip_network objects
                representing the IP blocks to remove.
        Returns:
            removed_blocks (tuple(ip_network)): tuple of ipaddress.ip_netework
                objects representing the removed IP blocks.
        """
        removed_blocks = self._mobility_client.remove_ip_blocks(blocks)
        return removed_blocks

    def cleanup(self):
        """ Cleanup added IP blocks """
        blocks = self.list_ip_blocks()
        self.remove_ip_blocks(blocks)

    def wait_for_changes(self):
        self._mobility_client.wait_for_changes()


class SpgwUtil(object):
    """
    Helper class to communicate with spgw for the tests.
    """

    def __init__(self):
        """
        Initialize spgw util.
        """
        self._stub = SpgwServiceStub(get_rpc_channel("spgw_service"))

    def create_bearer(self, imsi, lbi):
        """
        Sends a CreateBearer Request to SPGW service
        """
        print("Sending CreateBearer request to spgw service")
        req = CreateBearerRequest(
            sid=SIDUtils.to_pb(imsi),
            link_bearer_id=lbi,
            policy_rules=[
                PolicyRule(
                    qos=FlowQos(
                        qci=1,
                        gbr_ul=10000000,
                        gbr_dl=10000000,
                        max_req_bw_ul=10000000,
                        max_req_bw_dl=10000000,
                        arp=QosArp(
                            priority_level=1, pre_capability=1, pre_vulnerability=0
                        ),
                    ),
                    flow_list=[
                        FlowDescription(
                            match=FlowMatch(
                                ipv4_dst="0.0.0.0/0",
                                tcp_dst=5001,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                direction=FlowMatch.UPLINK,
                            ),
                            action=FlowDescription.PERMIT,
                        ),
                        FlowDescription(
                            match=FlowMatch(
                                ipv4_dst="192.168.129.42/24",
                                tcp_dst=5002,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                direction=FlowMatch.UPLINK,
                            ),
                            action=FlowDescription.PERMIT,
                        ),
                        FlowDescription(
                            match=FlowMatch(
                                ipv4_dst="192.168.129.42",
                                tcp_dst=5003,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                direction=FlowMatch.UPLINK,
                            ),
                            action=FlowDescription.PERMIT,
                        ),
                        FlowDescription(
                            match=FlowMatch(
                                ipv4_dst="192.168.129.42",
                                tcp_dst=5004,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                direction=FlowMatch.UPLINK,
                            ),
                            action=FlowDescription.PERMIT,
                        ),
                        FlowDescription(
                            match=FlowMatch(
                                ipv4_dst="192.168.129.42",
                                tcp_dst=5005,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                direction=FlowMatch.UPLINK,
                            ),
                            action=FlowDescription.DENY,
                        ),
                        FlowDescription(
                            match=FlowMatch(
                                ipv4_src="192.168.129.42",
                                tcp_src=5001,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                direction=FlowMatch.DOWNLINK,
                            ),
                            action=FlowDescription.PERMIT,
                        ),
                        FlowDescription(
                            match=FlowMatch(
                                ipv4_src="",
                                tcp_src=5002,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                direction=FlowMatch.DOWNLINK,
                            ),
                            action=FlowDescription.PERMIT,
                        ),
                        FlowDescription(
                            match=FlowMatch(
                                ipv4_src="192.168.129.64/26",
                                tcp_src=5003,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                direction=FlowMatch.DOWNLINK,
                            ),
                            action=FlowDescription.PERMIT,
                        ),
                        FlowDescription(
                            match=FlowMatch(
                                ipv4_src="192.168.129.42/16",
                                tcp_src=5004,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                direction=FlowMatch.DOWNLINK,
                            ),
                            action=FlowDescription.PERMIT,
                        ),
                        FlowDescription(
                            match=FlowMatch(
                                ipv4_src="192.168.129.42",
                                tcp_src=5005,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                direction=FlowMatch.DOWNLINK,
                            ),
                            action=FlowDescription.DENY,
                        ),
                    ],
                )
            ],
        )
        self._stub.CreateBearer(req)

    def delete_bearer(self, imsi, lbi, ebi):
        """
        Sends a DeleteBearer Request to SPGW service
        """
        print("Sending DeleteBearer request to spgw service")
        req = DeleteBearerRequest(
            sid=SIDUtils.to_pb(imsi), link_bearer_id=lbi, eps_bearer_ids=[ebi]
        )
        self._stub.DeleteBearer(req)
