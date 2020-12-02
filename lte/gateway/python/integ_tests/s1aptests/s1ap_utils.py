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

import ctypes
import ipaddress
import logging
import os
import shlex
import threading
import time
from enum import Enum
from queue import Queue
from typing import Optional

import grpc
import subprocess

import s1ap_types
from integ_tests.gateway.rpc import get_rpc_channel
from lte.protos.policydb_pb2 import (
    FlowDescription,
    FlowMatch,
    FlowQos,
    PolicyRule,
    QosArp,
)
from lte.protos.mobilityd_pb2 import IPAddress
from lte.protos.session_manager_pb2 import (
    DynamicRuleInstall,
    PolicyReAuthRequest,
    QoSInformation,
)
from lte.protos.abort_session_pb2 import (
    AbortSessionRequest,
    AbortSessionResult,
)
from lte.protos.spgw_service_pb2 import (
    CreateBearerRequest,
    DeleteBearerRequest,
)
from lte.protos.spgw_service_pb2_grpc import SpgwServiceStub
from magma.subscriberdb.sid import SIDUtils
from lte.protos.session_manager_pb2_grpc import SessionProxyResponderStub
from lte.protos.abort_session_pb2_grpc import AbortSessionResponderStub
from orc8r.protos.directoryd_pb2 import GetDirectoryFieldRequest
from orc8r.protos.directoryd_pb2_grpc import GatewayDirectoryServiceStub
from integ_tests.s1aptests.ovs.rest_api import get_datapath, get_flows

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

    PROT_CFG_CID_PCSCF_IPV6_ADDR_REQUEST = 0x0001
    PROT_CFG_CID_PCSCF_IPV4_ADDR_REQUEST = 0x000C
    PROT_CFG_CID_DNS_SERVER_IPV6_ADDR_REQUEST = 0x0003

    lib_name = "libtfw.so"

    _cond = threading.Condition()
    _msg = Queue()

    MAX_NUM_RETRIES = 5
    datapath = get_datapath()
    SPGW_TABLE = 0
    LOCAL_PORT = "LOCAL"

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
        self.gtpBridgeUtil = GTPBridgeUtils()

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

    def populate_pco(
        self, protCfgOpts_pr, pcscf_addr_type=None, dns_ipv6_addr=False
    ):
        """
        Populates the PCO values.
        Args:
            protCfgOpts_pr: PCO structure
            pcscf_addr_type: ipv4/ipv6/ipv4v6 flag
            dns_ipv6_addr: True/False flag
        Returns:
            None
        """
        # PCO parameters
        # Presence mask
        protCfgOpts_pr.pres = 1
        # Length
        protCfgOpts_pr.len = 4
        # Configuration protocol
        protCfgOpts_pr.cfgProt = 0
        # Extension bit for the additional parameters
        protCfgOpts_pr.ext = 1
        # Number of protocol IDs
        protCfgOpts_pr.numProtId = 0

        # Fill Number of container IDs and Container ID
        idx = 0
        if pcscf_addr_type == "ipv4":
            protCfgOpts_pr.numContId += 1
            protCfgOpts_pr.c[
                idx
            ].cid = S1ApUtil.PROT_CFG_CID_PCSCF_IPV4_ADDR_REQUEST
            idx += 1

        elif pcscf_addr_type == "ipv6":
            protCfgOpts_pr.numContId += 1
            protCfgOpts_pr.c[
                idx
            ].cid = S1ApUtil.PROT_CFG_CID_PCSCF_IPV6_ADDR_REQUEST
            idx += 1

        elif pcscf_addr_type == "ipv4v6":
            protCfgOpts_pr.numContId += 2
            protCfgOpts_pr.c[
                idx
            ].cid = S1ApUtil.PROT_CFG_CID_PCSCF_IPV4_ADDR_REQUEST
            idx += 1
            protCfgOpts_pr.c[
                idx
            ].cid = S1ApUtil.PROT_CFG_CID_PCSCF_IPV6_ADDR_REQUEST
            idx += 1

        if dns_ipv6_addr:
            protCfgOpts_pr.numContId += 1
            protCfgOpts_pr.c[
                idx
            ].cid = S1ApUtil.PROT_CFG_CID_DNS_SERVER_IPV6_ADDR_REQUEST

    def attach(
        self,
        ue_id,
        attach_type,
        resp_type,
        resp_msg_type,
        sec_ctxt=s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT,
        id_type=s1ap_types.TFW_MID_TYPE_IMSI,
        eps_type=s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH,
        pdn_type=1,
        pcscf_addr_type=None,
        dns_ipv6_addr=False,
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
            pdn_type:1 for IPv4, 2 for IPv6 and 3 for IPv4v6
            pcscf_addr_type:IPv4/IPv6/IPv4v6
        """
        attach_req = s1ap_types.ueAttachRequest_t()
        attach_req.ue_Id = ue_id
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt
        attach_req.pdnType_pr.pres = True
        attach_req.pdnType_pr.pdn_type = pdn_type

        # Populate PCO only if pcscf_addr_type is set
        if pcscf_addr_type or dns_ipv6_addr:
            self.populate_pco(
                attach_req.protCfgOpts_pr, pcscf_addr_type, dns_ipv6_addr
            )
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
            assert (
                context_setup.msg_type
                == s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value
            )

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
            elif S1ApUtil.CM_ESM_PDN_IPV6 == pdn_type:
                print("IPv6 PDN type received")
            elif S1ApUtil.CM_ESM_PDN_IPV4V6 == pdn_type:
                print("IPv4v6 PDN type received")
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
        assert (
            self.issue_cmd(s1ap_types.tfwCmd.UE_DETACH_REQUEST, detach_req)
            == 0
        )
        if reason_type == s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value:
            response = self.get_response()
            assert (
                s1ap_types.tfwCmd.UE_DETACH_ACCEPT_IND.value
                == response.msg_type
            )

        # Now wait for the context release response
        if wait_for_s1_ctxt_release:
            response = self.get_response()
            assert s1ap_types.tfwCmd.UE_CTX_REL_IND.value == response.msg_type

        with self._lock:
            del self._ue_ip_map[ue_id]

    def _verify_dl_flow(self, dl_flow_rules=None):
        # try at least 5 times before failing as gateway
        # might take some time to install the flows in ovs

        # Verify the total number of DL flows for this UE ip address
        num_dl_flows = 1
        for key, value in dl_flow_rules.items():
            tcp_src_port = 0
            ip_proto = 0
            ue_ip_str = str(key)
            dst_addr = "nw_dst" if key.version == 4 else "ipv6_dst"
            eth_typ = 2048 if key.version == 4 else 34525

            # Set to 1 for the default bearer
            total_num_dl_flows_to_be_verified = 1
            for item in value:
                for flow in item:
                    if flow["direction"] == FlowMatch.DOWNLINK:
                        total_num_dl_flows_to_be_verified += 1
            total_dl_ovs_flows_created = get_flows(
                self.datapath,
                {
                    "table_id": self.SPGW_TABLE,
                    "match": {
                        dst_addr: ue_ip_str,
                        "eth_type": eth_typ,
                        "in_port": self.LOCAL_PORT,
                    },
                },
            )
            assert (
                len(total_dl_ovs_flows_created)
                == total_num_dl_flows_to_be_verified
            )

            # Now verify the rules for every flow
            for item in value:
                for flow in item:
                    if flow["direction"] == FlowMatch.DOWNLINK:
                        ip_src_addr = (
                            flow["ipv4_src"]
                            if key.version == 4
                            else flow["ipv6_src"]
                        )
                        ip_src = "ipv4_src" if key.version == 4 else "ipv6_src"
                        ip_dst = "ipv4_dst" if key.version == 4 else "ipv6_dst"
                        tcp_src_port = flow["tcp_src_port"]
                        ip_proto = flow["ip_proto"]
                        for i in range(self.MAX_NUM_RETRIES):
                            print("Get downlink flows: attempt ", i)
                            downlink_flows = get_flows(
                                self.datapath,
                                {
                                    "table_id": self.SPGW_TABLE,
                                    "match": {
                                        dst_addr: ue_ip_str,
                                        "eth_type": eth_typ,
                                        "in_port": self.LOCAL_PORT,
                                        ip_src: ip_src_addr,
                                        "tcp_src": tcp_src_port,
                                        "ip_proto": ip_proto,
                                    },
                                },
                            )
                            if len(downlink_flows) >= num_dl_flows:
                                break
                            time.sleep(
                                5
                            )  # sleep for 5 seconds before retrying
                        assert (
                            len(downlink_flows) >= num_dl_flows
                        ), "Downlink flow missing for UE"
                        assert downlink_flows[0]["match"][ip_dst] == ue_ip_str
                        actions = downlink_flows[0]["instructions"][0][
                            "actions"
                        ]
                        has_tunnel_action = any(
                            action
                            for action in actions
                            if action["field"] == "tunnel_id"
                            and action["type"] == "SET_FIELD"
                        )
                        assert bool(has_tunnel_action)

    def verify_flow_rules(self, num_ul_flows, dl_flow_rules=None):
        GTP_PORT = self.gtpBridgeUtil.get_gtp_port_no()
        # Check if UL and DL OVS flows are created
        print("************ Verifying flow rules")
        # UPLINK
        print("Checking for uplink flow")
        # try at least 5 times before failing as gateway
        # might take some time to install the flows in ovs
        for i in range(self.MAX_NUM_RETRIES):
            print("Get uplink flows: attempt ", i)
            uplink_flows = get_flows(
                self.datapath,
                {
                   "table_id": self.SPGW_TABLE,
                   "match": {
                       "in_port": GTP_PORT,
                   }
                },
            )
            if len(uplink_flows) == num_ul_flows:
                break
            time.sleep(5)  # sleep for 5 seconds before retrying
        assert len(uplink_flows) == num_ul_flows, "Uplink flow missing for UE"
        assert uplink_flows[0]["match"]["tunnel_id"] is not None

        # DOWNLINK
        print("Checking for downlink flow")
        self._verify_dl_flow(dl_flow_rules)

    def verify_paging_flow_rules(self, ip_list):
        # Check if paging flows are created
        print("************ Verifying paging flow rules")
        num_paging_flows_to_be_verified = 1
        for ip in ip_list:
            ue_ip_str = str(ip)
            print("Verifying paging flow for ip", ue_ip_str)
            for i in range(self.MAX_NUM_RETRIES):
                print("Get paging flows: attempt ", i)
                paging_flows = get_flows(
                    self.datapath,
                    {
                        "table_id": self.SPGW_TABLE,
                        "match": {
                            "nw_dst": ue_ip_str,
                            "eth_type": 2048,
                            "priority": 5,
                        },
                    },
                )
                if len(paging_flows) == num_paging_flows_to_be_verified:
                    break
                time.sleep(5)  # sleep for 5 seconds before retrying
            assert (
                len(paging_flows) == num_paging_flows_to_be_verified
            ), "Paging flow missing for UE"

            # TODO - Verify that the action is to send to controller
            """controller_port = 4294967293
            actions = paging_flows[0]["instructions"][0]["actions"]
            has_tunnel_action = any(
                action
                for action in actions
                if action["type"] == "OUTPUT"
                and action["port"] == controller_port
            )
            assert bool(has_tunnel_action)"""


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
    stateless_cmds = Enum("stateless_cmds", "CHECK DISABLE ENABLE")
    config_update_cmds = Enum("config_update_cmds", "MODIFY RESTORE")
    apn_correction_cmds = Enum("apn_correction_cmds", "DISABLE ENABLE")

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

        self._command = (
            "sshpass -p {password} ssh "
            "-o UserKnownHostsFile=/dev/null "
            "-o StrictHostKeyChecking=no "
            "{user}@{host} {command}"
        )

    def exec_command(self, command):
        """
        Run a command remotly on magma_dev VM.

        Args:
            command: command (str) to be executed on remote host
            e.g. 'sed -i \'s/config1/config2/g\' /etc/magma/mme.yml'

        """
        data = self._data
        data["command"] = '"' + command + '"'
        param_list = shlex.split(self._command.format(**data))
        return subprocess.call(
            param_list,
            shell=False,
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
        )

    def exec_command_output(self, command):
        """
        Run a command remotly on magma_dev VM.

        Args:
            command: command (str) to be executed on remote host
            e.g. 'sed -i \'s/config1/config2/g\' /etc/magma/mme.yml'

        """
        data = self._data
        data["command"] = '"' + command + '"'
        param_list = shlex.split(self._command.format(**data))
        return subprocess.check_output(
            param_list,
            shell=False,
        ).decode("utf-8")

    def config_stateless(self, cmd):
        """
        Configure the stateless mode on the access gateway

        Args:
          cmd: Specify how to configure stateless mode on AGW,
          should be one of
            check: Run a check whether AGW is stateless or not
            enable: Enable stateless mode, do nothing if already stateless
            disable: Disable stateless mode, do nothing if already stateful

        """
        magtivate_cmd = "source /home/vagrant/build/python/bin/activate"
        venvsudo_cmd = "sudo -E PATH=$PATH PYTHONPATH=$PYTHONPATH env"
        config_stateless_script = "/usr/local/bin/config_stateless_agw.py"

        ret_code = self.exec_command(
            magtivate_cmd + " && " + venvsudo_cmd + " python3 " +
            config_stateless_script + " " + cmd.name.lower()
            )

        if ret_code == 0:
            print("AGW is stateless")
        elif ret_code == 1:
            print("AGW is stateful")
        elif ret_code == 2:
            print("AGW is in a mixed config, check gateway")
        else:
            print("Unknown command")

    def restart_all_services(self):
        """
            Restart all magma services on magma_dev VM
            """
        self.exec_command(
            "sudo service magma@* stop ; sudo service magma@magmad start"
        )
        print(
            "Waiting for all services to restart. Sleeping for 60 seconds.."
        )
        timeSlept = 0
        while timeSlept < 60:
            time.sleep(5)
            timeSlept += 5
            print("*********** Slept for " + str(timeSlept) + " seconds")

    def restart_services(self, services):
        """
        Restart a list of magmad services.

        Args:
            services: List of (str) services names

        """
        self._magmad_client.restart_services(services)

    def update_mme_config_for_sanity(self, cmd):
        mme_config_update_script = (
            "/home/vagrant/magma/lte/gateway/deploy/roles/magma/files/"
            "update_mme_config_for_sanity.sh"
        )

        action = cmd.name.lower()
        ret_code = self.exec_command(
            "sudo -E " + mme_config_update_script + " " + action
        )

        if ret_code == 0:
            print("MME configuration is updated successfully")
        elif ret_code == 1:
            assert False, (
                "Failed to "
                + action
                + " MME configuration. Error: Invalid command"
            )
        elif ret_code == 2:
            assert False, (
                "Failed to "
                + action
                + " MME configuration. Error: MME configuration file is "
                + "missing"
            )
        elif ret_code == 3:
            assert False, (
                "Failed to "
                + action
                + " MME configuration. Error: MME configuration's backup file "
                + "is missing"
            )
        else:
            assert False, (
                "Failed to "
                + action
                + " MME configuration. Error: Unknown error"
            )

    def config_apn_correction(self, cmd):
        """
        Configure the apn correction mode on the access gateway

        Args:
          cmd: Specify how to configure apn correction mode on AGW,
          should be one of
            enable: Enable apn correction feature, do nothing if already enabled
            disable: Disable apn correction feature, do nothing if already disabled

        """
        apn_correction_cmd = ""
        if cmd.name == MagmadUtil.apn_correction_cmds.ENABLE.name:
            apn_correction_cmd = "sed -i \'s/correction: false/correction: true/g\' /etc/magma/mme.yml"
        else:
            apn_correction_cmd = "sed -i \'s/correction: true/correction: false/g\' /etc/magma/mme.yml"

        ret_code = self.exec_command(
            "sudo " + apn_correction_cmd)

        if ret_code == 0:
            print("APN Correction configured")
        else:
            print("APN Correction failed")

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
                    id="rar_rule_1",
                    qos=FlowQos(
                        qci=qci_val,
                        gbr_ul=10000000,
                        gbr_dl=10000000,
                        max_req_bw_ul=10000000,
                        max_req_bw_dl=10000000,
                        arp=QosArp(
                            priority_level=1,
                            pre_capability=1,
                            pre_vulnerability=0,
                        ),
                    ),
                    flow_list=[
                        FlowDescription(
                            match=FlowMatch(
                                ip_dst=IPAddress(
                                    version=IPAddress.IPV4,
                                    address="0.0.0.0/0".encode('utf-8')),
                                tcp_dst=5001,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                direction=FlowMatch.UPLINK,
                            ),
                            action=FlowDescription.PERMIT,
                        ),
                        FlowDescription(
                            match=FlowMatch(
                                ip_dst=IPAddress(
                                    version=IPAddress.IPV4,
                                    address="192.168.129.42/24".encode('utf-8')
                                ),
                                tcp_dst=5002,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                direction=FlowMatch.UPLINK,
                            ),
                            action=FlowDescription.PERMIT,
                        ),
                        FlowDescription(
                            match=FlowMatch(
                                ip_dst=IPAddress(
                                    version=IPAddress.IPV4,
                                    address="192.168.129.42".encode('utf-8')),
                                tcp_dst=5003,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                direction=FlowMatch.UPLINK,
                            ),
                            action=FlowDescription.PERMIT,
                        ),
                        FlowDescription(
                            match=FlowMatch(
                                ip_dst=IPAddress(
                                    version=IPAddress.IPV4,
                                    address="192.168.129.42".encode('utf-8')),
                                tcp_dst=5004,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                direction=FlowMatch.UPLINK,
                            ),
                            action=FlowDescription.PERMIT,
                        ),
                        FlowDescription(
                            match=FlowMatch(
                                ip_dst=IPAddress(
                                    version=IPAddress.IPV4,
                                    address="192.168.129.42".encode('utf-8')),
                                tcp_dst=5005,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                direction=FlowMatch.UPLINK,
                            ),
                            action=FlowDescription.DENY,
                        ),
                        FlowDescription(
                            match=FlowMatch(
                                ip_src=IPAddress(
                                    version=IPAddress.IPV4,
                                    address="192.168.129.42".encode('utf-8')),
                                tcp_src=5001,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                direction=FlowMatch.DOWNLINK,
                            ),
                            action=FlowDescription.PERMIT,
                        ),
                        FlowDescription(
                            match=FlowMatch(
                                ip_src=IPAddress(version=IPAddress.IPV4,
                                                 address="".encode('utf-8')),
                                tcp_dst=5002,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                direction=FlowMatch.DOWNLINK,
                            ),
                            action=FlowDescription.PERMIT,
                        ),
                        FlowDescription(
                            match=FlowMatch(
                                ip_src=IPAddress(
                                    version=IPAddress.IPV4,
                                    address="192.168.129.64/26".encode('utf-8')
                                ),
                                tcp_src=5003,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                direction=FlowMatch.DOWNLINK,
                            ),
                            action=FlowDescription.PERMIT,
                        ),
                        FlowDescription(
                            match=FlowMatch(
                                ip_src=IPAddress(
                                    version=IPAddress.IPV4,
                                    address="192.168.129.42/16".encode('utf-8')
                                ),
                                tcp_src=5004,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                direction=FlowMatch.DOWNLINK,
                            ),
                            action=FlowDescription.PERMIT,
                        ),
                        FlowDescription(
                            match=FlowMatch(
                                ip_src=IPAddress(
                                    version=IPAddress.IPV4,
                                    address="192.168.129.42".encode('utf-8')),
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

    def create_bearer_ipv4v6(
        self, imsi, lbi, qci_val=1, ipv4=False, ipv6=False
    ):
        """
        Sends a CreateBearer Request with ipv4/ipv6/ipv4v6 packet """
        """ filters to SPGW service """
        print("Sending CreateBearer request to spgw service")
        flow_match_list = []
        if ipv4:
            flow_match_list.append(
                FlowDescription(
                    match=FlowMatch(
                        ipv4_dst="192.168.129.42/24",
                        tcp_src=5001,
                        ip_proto=FlowMatch.IPPROTO_TCP,
                        direction=FlowMatch.UPLINK,
                    ),
                    action=FlowDescription.PERMIT,
                )
            )
            flow_match_list.append(
                FlowDescription(
                    match=FlowMatch(
                        ipv4_src="192.168.129.42",
                        tcp_src=5001,
                        ip_proto=FlowMatch.IPPROTO_TCP,
                        direction=FlowMatch.DOWNLINK,
                    ),
                    action=FlowDescription.PERMIT,
                )
            )

        if ipv6:
            flow_match_list.append(
                FlowDescription(
                    match=FlowMatch(
                        ipv6_dst="5546:222:2259::226",
                        ip_proto=FlowMatch.IPPROTO_UDP,
                        direction=FlowMatch.UPLINK,
                    ),
                    action=FlowDescription.PERMIT,
                )
            )
            flow_match_list.append(
                FlowDescription(
                    match=FlowMatch(
                        ip_proto=FlowMatch.IPPROTO_UDP,
                        direction=FlowMatch.DOWNLINK,
                    ),
                    action=FlowDescription.PERMIT,
                )
            )

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
                            priority_level=1,
                            pre_capability=1,
                            pre_vulnerability=0,
                        ),
                    ),
                    flow_list=flow_match_list,
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
        self._abort_session_stub = AbortSessionResponderStub(
            get_rpc_channel("abort_session_service")
        )
        self._directorydstub = GatewayDirectoryServiceStub(
            get_rpc_channel("directoryd")
        )

    def get_flow_match(self, flow_list, flow_match_list):
        """
        Populates flow match list
        """
        for flow in flow_list:
            flow_direction = flow["direction"]
            ip_protocol = flow["ip_proto"]
            if ip_protocol == FlowMatch.IPPROTO_TCP:
                udp_src_port = 0
                udp_dst_port = 0
                tcp_src_port = (
                    int(flow["tcp_src_port"]) if "tcp_src_port" in flow else 0
                )
                tcp_dst_port = (
                    int(flow["tcp_dst_port"]) if "tcp_dst_port" in flow else 0
                )
            elif ip_protocol == FlowMatch.IPPROTO_UDP:
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

            src_addr = None
            if flow.get("ipv4_src", None):
                src_addr = IPAddress(
                    version=IPAddress.IPV4,
                    address=flow.get("ipv4_src").encode('utf-8'))
            elif flow.get("ipv6_src", None):
                src_addr = IPAddress(
                    version=IPAddress.IPV6,
                    address=flow.get("ipv6_src").encode('utf-8'))

            dst_addr = None
            if flow.get("ipv4_dst", None):
                dst_addr = IPAddress(
                    version=IPAddress.IPV4,
                    address=flow.get("ipv4_dst").encode('utf-8'))
            elif flow.get("ipv6_dst", None):
                dst_addr = IPAddress(
                    version=IPAddress.IPV6,
                    address=flow.get("ipv6_dst").encode('utf-8'))


            flow_match_list.append(
                FlowDescription(
                    match=FlowMatch(
                        ip_dst=dst_addr,
                        ip_src=src_addr,
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

    def create_AbortSessionRequest(self, imsi: str) -> AbortSessionResult:
        # Get SessionID
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
        return self._abort_session_stub.AbortSession(
            AbortSessionRequest(
                session_id=res.value,
                user_name=imsi,
            )
        )


class GTPBridgeUtils:
    def __init__(self):
        self.magma_utils = MagmadUtil(None)
        ret = self.magma_utils.exec_command_output("sudo grep ovs_multi_tunnel  /etc/magma/spgw.yml")
        if "false" in ret:
            self.gtp_port_name = "gtp0"
        else:
            self.gtp_port_name = "g_8d3ca8c0"

    def get_gtp_port_no(self) -> Optional[int]:
        output = self.magma_utils.exec_command_output("sudo ovsdb-client dump Interface name ofport")
        for line in output.split('\n'):
            if self.gtp_port_name in line:
                port_info = line.split()
                return port_info[1]
