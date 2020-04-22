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
import time
from queue import Queue
import grpc

import s1ap_types
from integ_tests.gateway.rpc import get_rpc_channel
from lte.protos.policydb_pb2 import (
    FlowDescription,
    FlowMatch,
    FlowQos,
    PolicyRule,
    QosArp,
)
from lte.protos.session_manager_pb2 import (
    DynamicRuleInstall,
    PolicyReAuthRequest,
    QoSInformation,
)
from lte.protos.spgw_service_pb2 import CreateBearerRequest, DeleteBearerRequest
from lte.protos.spgw_service_pb2_grpc import SpgwServiceStub
from magma.subscriberdb.sid import SIDUtils
from lte.protos.session_manager_pb2_grpc import SessionProxyResponderStub
from orc8r.protos.directoryd_pb2 import GetDirectoryFieldRequest
from orc8r.protos.directoryd_pb2_grpc import GatewayDirectoryServiceStub

DEFAULT_GRPC_TIMEOUT = 10


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

    def populate_pco(self, protCfgOpts_pr, volte_attach_type):
        """
        Populates the PCO values.
        Args:
            protCfgOpts_pr: PCO structure
            volte_attach_type: ipv4/ipv6 flag
        Returns:
            None
        """
        protCfgOpts_pr.pres = 1
        protCfgOpts_pr.len = 4
        protCfgOpts_pr.cfgProt = 0
        protCfgOpts_pr.ext = 1
        protCfgOpts_pr.numProtId = 0

        if volte_attach_type == "ipv4":
            protCfgOpts_pr.numContId = 1
            protCfgOpts_pr.c[0].cid = 0x000c

        elif volte_attach_type == "ipv6":
            protCfgOpts_pr.numContId = 1
            protCfgOpts_pr.c[0].cid = 0x0001

        elif volte_attach_type == "ipv4v6":
            protCfgOpts_pr.numContId = 2
            protCfgOpts_pr.c[0].cid = 0x000C
            protCfgOpts_pr.c[1].cid = 0x0001


    def attach(
        self,
        ue_id,
        attach_type,
        resp_type,
        resp_msg_type,
        sec_ctxt=s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT,
        id_type=s1ap_types.TFW_MID_TYPE_IMSI,
        eps_type=s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH,
        volte_attach_type=None):
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

        # Populate PCO only for VoLTE
        if volte_attach_type:
            self.populate_pco(attach_req.protCfgOpts_pr, volte_attach_type)
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

    def config_apn_data(self, imsi, apn_list):
        """ Add APN details """
        self._subscriber_client.config_apn_details(imsi, apn_list)

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

        self._data = {
            "user": "vagrant",
            "host": "192.168.60.142",
            "password": "vagrant",
            "command": "test",
        }

        self._command = "sshpass -p {password} ssh " \
                        "-o UserKnownHostsFile=/dev/null " \
                        "-o StrictHostKeyChecking=no " \
                        "{user}@{host} {command}"

    def exec_command(self, command):
        """
        Run a command remotly on magma_dev VM.

        Args:
            command: command (str) to be executed on remote host
            e.g. 'sed -i \'s/config1/config2/g\' /etc/magma/mme.yml'

        """
        data = self._data
        data["command"] = '"' + command + '"'
        os.system(self._command.format(**data))

    def set_config_stateless(self, enabled):
        """
            Sets the use_stateless flag in mme.yml file

            Args:
                enabled: sets the flag to true if enabled

            """
        if enabled:
            self.exec_command(
                "sed -i 's/use_stateless: false/use_stateless: true/g' "
                "/etc/magma/mme.yml"
            )
        else:
            self.exec_command(
                "sed -i 's/use_stateless: true/use_stateless: false/g' "
                "/etc/magma/mme.yml"
            )

    def restart_all_services(self):
        """
            Restart all magma services on magma_dev VM
            """
        self.exec_command("sudo service magma@* stop ; "
                          "sudo service magma@magmad start")
        time.sleep(10)

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

    def create_bearer(self, imsi, lbi, qci_val=1):
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
                        qci=qci_val,
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


class SessionManagerUtil(object):
    """
    Helper class to communicate with session manager for the tests.
    """

    def __init__(self):
        """
        Initialize sessionManager util.
        """
        self._session_stub = SessionProxyResponderStub(
            get_rpc_channel("sessiond")
        )
        self._directorydstub = GatewayDirectoryServiceStub(
            get_rpc_channel("directoryd")
        )

    def get_flow_match(self, flow_list, flow_match_list):
        """
        Populates flow match list
        """
        for flow in flow_list:
            flow_direction = (
                FlowMatch.UPLINK
                if flow["direction"] == "UL"
                else FlowMatch.DOWNLINK
            )
            ip_protocol = flow["ip_proto"]
            if ip_protocol == "TCP":
                ip_protocol = FlowMatch.IPPROTO_TCP
                udp_src_port = 0
                udp_dst_port = 0
                tcp_src_port = (
                    int(flow["tcp_src_port"]) if "tcp_src_port" in flow else 0
                )
                tcp_dst_port = (
                    int(flow["tcp_dst_port"]) if "tcp_dst_port" in flow else 0
                )
            elif ip_protocol == "UDP":
                ip_protocol = FlowMatch.IPPROTO_UDP
                tcp_src_port = 0
                tcp_dst_port = 0
                udp_src_port = (
                    int(flow["udp_src_port"]) if "udp_src_port" in flow else 0
                )
                udp_dst_port = (
                    int(flow["udp_dst_port"]) if "udp_dst_port" in flow else 0
                )
            else:
                udp_src_port = 0
                udp_dst_port = 0
                tcp_src_port = 0
                tcp_dst_port = 0

            ipv4_src_addr = flow.get("ipv4_src", None)
            ipv4_dst_addr = flow.get("ipv4_dst", None)

            flow_match_list.append(
                FlowDescription(
                    match=FlowMatch(
                        ipv4_dst=ipv4_dst_addr,
                        ipv4_src=ipv4_src_addr,
                        tcp_src=tcp_src_port,
                        tcp_dst=tcp_dst_port,
                        udp_src=udp_src_port,
                        udp_dst=udp_dst_port,
                        ip_proto=ip_protocol,
                        direction=flow_direction,
                    ),
                    action=FlowDescription.PERMIT,
                )
            )

    def create_ReAuthRequest(self, imsi, policy_id, flow_list, qos):
        """
        Sends Policy RAR message to session manager
        """
        print("Sending Policy RAR message to session manager")
        flow_match_list = []
        res = None
        self.get_flow_match(flow_list, flow_match_list)

        policy_qos = FlowQos(
            qci=qos["qci"],
            max_req_bw_ul=qos["max_req_bw_ul"],
            max_req_bw_dl=qos["max_req_bw_dl"],
            gbr_ul=qos["gbr_ul"],
            gbr_dl=qos["gbr_dl"],
            arp=QosArp(
                priority_level=qos["arp_prio"],
                pre_capability=qos["pre_cap"],
                pre_vulnerability=qos["pre_vul"],
            ),
        )

        policy_rule = PolicyRule(
            id=policy_id,
            priority=qos["priority"],
            flow_list=flow_match_list,
            tracking_type=PolicyRule.NO_TRACKING,
            rating_group=1,
            monitoring_key=None,
            qos=policy_qos,
        )

        qos = QoSInformation(qci=qos["qci"])

        # Get sessionid
        req = GetDirectoryFieldRequest(id=imsi, field_key="session_id")
        try:
            res = self._directorydstub.GetDirectoryField(
                req, DEFAULT_GRPC_TIMEOUT
            )
        except grpc.RpcError as err:
            logging.error(
                "GetDirectoryFieldRequest error for id: %s! [%s] %s",
                imsi,
                err.code(),
                err.details(),
            )

        self._session_stub.PolicyReAuth(
            PolicyReAuthRequest(
                session_id=res.value,
                imsi=imsi,
                rules_to_remove=[],
                rules_to_install=[],
                dynamic_rules_to_install=[
                    DynamicRuleInstall(policy_rule=policy_rule)
                ],
                event_triggers=[],
                revalidation_time=None,
                usage_monitoring_credits=[],
                qos_info=qos,
            )
        )
