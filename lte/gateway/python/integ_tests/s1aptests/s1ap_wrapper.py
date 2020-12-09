""""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import os
import time
import ctypes
import redis
import json
import jsonpickle
from typing import Union
from magma.common.redis.serializers import get_proto_deserializer

import s1ap_types
from integ_tests.common.magmad_client import MagmadServiceGrpc

# from integ_tests.cloud.cloud_manager import CloudManager
from integ_tests.common.mobility_service_client import MobilityServiceGrpc
from integ_tests.common.service303_utils import GatewayServicesUtil
from integ_tests.common.subscriber_db_client import (
    SubscriberDbGrpc,
    SubscriberDbCassandra,
)
from integ_tests.s1aptests.s1ap_utils import (
    MagmadUtil,
    MobilityUtil,
    S1ApUtil,
    SubscriberUtil,
)
from integ_tests.s1aptests.util.traffic_util import TrafficUtil

from lte.protos.keyval_pb2 import IPDesc
from lte.protos.policydb_pb2 import PolicyRule, InstalledPolicies
from lte.protos.oai.mme_nas_state_pb2 import MmeNasState, UeContext
from lte.protos.oai.spgw_state_pb2 import SpgwState, S11BearerContext
from lte.protos.oai.s1ap_state_pb2 import S1apState, UeDescription
from magma.common.redis.serializers import get_json_deserializer, \
    get_proto_deserializer
from magma.mobilityd.serialize_utils import deserialize_ip_block, \
    deserialize_ip_desc

def _deserialize_generic_json(
        element: Union[str, dict, list])-> Union[str, dict, list]:
    """
    Helper function to deserialize dictionaries or list with nested
    json strings
    :param element
    """
    if isinstance(element, str):
        # try to deserialize as json string
        try:
            element = ast.literal_eval(element)
        except:
            try:
                element = jsonpickle.decode(element)
            except:
                return element

    if isinstance(element, dict):
        keys = element.keys()
    elif isinstance(element, list):
        keys = range(len(element))
    else:
        # in case it is neither of the know elements, just return as is
        return element

    for k in keys:
        element[k] = _deserialize_generic_json(element[k])
    return element

def _deserialize_session_json(serialized_json_str: bytes) -> str:
    """
    Helper function to deserialize sessiond:sessions hash list values
    :param serialized_json_str
    """
    res = _deserialize_generic_json(str(serialized_json_str, 'utf-8', 'ignore'))
    dumped = json.dumps(res, indent=2, sort_keys=True)
    return dumped

STATE_DESERIALIZERS = {
    'assigned_ip_blocks': deserialize_ip_block,
    'ip_states': deserialize_ip_desc,
    'sessions': _deserialize_session_json,
    'rule_names': get_json_deserializer(),
    'rule_ids': get_json_deserializer(),
    'rule_versions': get_json_deserializer(),
}

STATE_PROTOS = {
    'mme_nas_state': MmeNasState,
    'spgw_state': SpgwState,
    's1ap_state': S1apState,
    'mme': UeContext,
    'spgw': S11BearerContext,
    's1ap': UeDescription,
    'mobilityd_ipdesc_record': IPDesc,
    'rules': PolicyRule,
    'installed': InstalledPolicies,
}

class TestWrapper(object):
    """
    Module wrapping boiler plate code for all test setups and cleanups.
    """

    # With the current mask value of 24 in TEST_IP_BLOCK, we can allocate a
    # maximum of 255 UE IP addresses only. Moreover, magma has reserved 12 IP
    # addresses for testing purpose, hence maximum allowed free IP addresses
    # are 243. We need to change this mask value in order to allocate more than
    # 243 UE IP addresses. Therefore, with the mask value of n, the maximum
    # number of UE IP addresses allowed will be ((2^(32-n)) - 13).
    # Example:
    #  mask value 24, max allowed UE IP addresses = ((2^(32-24)) - 13) = 243
    #  mask value 20, max allowed UE IP addresses = ((2^(32-20)) - 13) = 4083
    #  mask value 17, max allowed UE IP addresses = ((2^(32-17)) - 13) = 32755
    # Decreasing the mask value will provide more UE IP addresses in the free
    # IP address pool
    TEST_IP_BLOCK = "192.168.128.0/24"
    MSX_S1_RETRY = 2

    def __init__(
        self,
        stateless_mode=MagmadUtil.stateless_cmds.ENABLE,
        apn_correction=MagmadUtil.apn_correction_cmds.DISABLE,
    ):
        """
        Initialize the various classes required by the tests and setup.
        """
        self._s1_util = S1ApUtil()
        self._enBConfig()

        if self._test_oai_upstream:
            subscriber_client = SubscriberDbCassandra()
            self.wait_gateway_healthy = False
        else:
            subscriber_client = SubscriberDbGrpc()
            self.wait_gateway_healthy = True

        mobility_client = MobilityServiceGrpc()
        magmad_client = MagmadServiceGrpc()
        self._sub_util = SubscriberUtil(subscriber_client)
        # Remove existing subscribers to start
        self._sub_util.cleanup()
        self._mobility_util = MobilityUtil(mobility_client)
        self._mobility_util.cleanup()
        self._magmad_util = MagmadUtil(magmad_client)
        self._magmad_util.config_stateless(stateless_mode)
        self._magmad_util.config_apn_correction(apn_correction)
        # gateway tests don't require restart, just wait for healthy now
        self._gateway_services = GatewayServicesUtil()
        if not self.wait_gateway_healthy:
            self.init_s1ap_tester()

        self.redis_client = redis.Redis(
            host=os.environ.get('GATEWAY_IP', '192.168.60.142'),
            port=6380)
        self.redis_client.flushall()
        self._configuredUes = []
        self._ue_idx = 0  # Index of UEs already used in test
        self._trf_util = TrafficUtil()

    def init_s1ap_tester(self):
        """
        Initialize the s1ap tester and the UEApp.

        Doing this separately allows initialization to occur during
        tests rather than during setup stage.
        """
        # config ip first, because cloud tests will restart gateway
        self.configIpBlock()

        self._s1setup()
        self._configUEApp()

    @property
    def _test_cloud(self):
        test_cloud = os.getenv("MAGMA_S1APTEST_USE_CLOUD") is not None
        return test_cloud

    @property
    def _test_oai_upstream(self):
        return os.getenv("TEST_OAI_UPSTREAM") is not None

    def _enBConfig(self):
        """Helper to configure the eNB"""
        # Using exaggerated prints makes the stdout easier to read.
        print("************************* Enb tester config")
        req = s1ap_types.FwNbConfigReq_t()
        req.cellId_pr.pres = True
        req.cellId_pr.cell_id = 10
        assert self._s1_util.issue_cmd(s1ap_types.tfwCmd.ENB_CONFIG, req) == 0
        response = self._s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.ENB_CONFIG_CONFIRM.value
        res = response.cast(s1ap_types.FwNbConfigCfm_t)
        assert res.status == s1ap_types.CfgStatus.CFG_DONE.value

    def _issue_s1setup_req(self):
        """ Issue the actual setup request and get the response"""
        req = None
        assert (
            self._s1_util.issue_cmd(s1ap_types.tfwCmd.ENB_S1_SETUP_REQ, req)
            == 0
        )
        response = self._s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.ENB_S1_SETUP_RESP.value
        return response.cast(s1ap_types.FwNbS1setupRsp_t)

    def _s1setup(self):
        """Helper to setup s1 to the EPC"""
        print("************************* S1 setup")
        res = self._issue_s1setup_req()

        retry = 0
        while retry < TestWrapper.MSX_S1_RETRY:
            if (
                res.res == s1ap_types.S1_setp_Result.S1_SETUP_FAILED.value
                and res.waitIe.pres == 1
            ):
                print(
                    "Received time to wait in S1-Setup-Failure" " message is",
                    res.waitIe.val,
                )
                time.sleep(res.waitIe.val)
                res = self._issue_s1setup_req()
                retry += 1
            else:
                # Not a failure in setup.
                break

        assert res.res == s1ap_types.S1_setp_Result.S1_SETUP_SUCCESS.value

    def _configUEApp(self):
        """ Update the internal configuration of the UEApp"""
        print("************************* UE App config")
        req = s1ap_types.ueAppConfig_t()
        req.nasProcGuardTimer_pr.pres = True
        req.nasProcGuardTimer_pr.nas_proc_guard_timer = 5
        assert (
            self._s1_util.issue_cmd(s1ap_types.tfwCmd.UE_APPL_CONFIG, req) == 0
        )
        response = self._s1_util.get_response()
        assert (
            s1ap_types.tfwCmd.UE_APP_CONFIG_COMPLETE_IND.value
            == response.msg_type
        )

    def _getAddresses(self, *ues):
        """ Retrieve IP addresses for the given UEs

        Will put None for IPs in the cases where a UE has been included that
        doesn't have a cached IP (e.g. the UE has not yet been attached)

        Args:
            ues (list(s1ap_types.ueAppConfig_t)): the UEs whose IPs we want

        Returns a list of ipaddress.ip_address objects, corresponding in order
            with the input UE parameters
        """
        return [self._s1_util.get_ip(ue.ue_id) for ue in ues]

    def configIpBlock(self):
        """ Removes any existing allocated blocks, then adds the ones used for
        testing """
        print("************************* Configuring IP block")
        self._mobility_util.remove_all_ip_blocks()
        self._mobility_util.add_ip_block(self.TEST_IP_BLOCK)
        print("************************* Waiting for IP changes to propagate")
        self._mobility_util.wait_for_changes()

    def configUEDevice(self, num_ues, reqData=[]):
        """ Configure the device on the UE side """
        reqs = self._sub_util.add_sub(num_ues=num_ues)
        for i in range(num_ues):
            print(
                "************************* UE device config for ue_id ",
                reqs[i].ue_id,
            )
            if reqData and bool(reqData[i]):
                if reqData[i].ueNwCap_pr.pres:
                    reqs[i].ueNwCap_pr.pres = reqData[i].ueNwCap_pr.pres
                    reqs[i].ueNwCap_pr.eea2_128 = reqData[
                        i
                    ].ueNwCap_pr.eea2_128
                    reqs[i].ueNwCap_pr.eea1_128 = reqData[
                        i
                    ].ueNwCap_pr.eea1_128
                    reqs[i].ueNwCap_pr.eea0 = reqData[i].ueNwCap_pr.eea0
                    reqs[i].ueNwCap_pr.eia2_128 = reqData[
                        i
                    ].ueNwCap_pr.eia2_128
                    reqs[i].ueNwCap_pr.eia1_128 = reqData[
                        i
                    ].ueNwCap_pr.eia1_128
                    reqs[i].ueNwCap_pr.eia0 = reqData[i].ueNwCap_pr.eia0

            assert (
                self._s1_util.issue_cmd(s1ap_types.tfwCmd.UE_CONFIG, reqs[i])
                == 0
            )
            response = self._s1_util.get_response()
            assert (
                s1ap_types.tfwCmd.UE_CONFIG_COMPLETE_IND.value
                == response.msg_type
            )
            # APN configuration below can be overwritten in the test case
            # after configuring UE device.
            self.configAPN(
                "IMSI" + "".join([str(j) for j in reqs[i].imsi]), None
            )
            self._configuredUes.append(reqs[i])

        self.check_gw_health_after_ue_load()

    def configUEDevice_ues_same_imsi(self, num_ues):
        """ Configure the device on the UE side with same IMSI and
        having different ue-id"""
        reqs = self._sub_util.add_sub(num_ues=num_ues)
        for i in range(num_ues):
            print(
                "************************* UE device config for ue_id ",
                reqs[i].ue_id,
            )
            assert (
                self._s1_util.issue_cmd(s1ap_types.tfwCmd.UE_CONFIG, reqs[i])
                == 0
            )
            response = self._s1_util.get_response()
            assert (
                s1ap_types.tfwCmd.UE_CONFIG_COMPLETE_IND.value
                == response.msg_type
            )
            # APN configuration below can be overwritten in the test case
            # after configuring UE device.
            self.configAPN(
                "IMSI" + "".join([str(j) for j in reqs[i].imsi]), None
            )
            self._configuredUes.append(reqs[i])
        for i in range(num_ues):
            reqs[i].ue_id = 2
            print(
                "************************* UE device config for ue_id ",
                reqs[i].ue_id,
            )
            assert (
                self._s1_util.issue_cmd(s1ap_types.tfwCmd.UE_CONFIG, reqs[i])
                == 0
            )
            response = self._s1_util.get_response()
            assert (
                s1ap_types.tfwCmd.UE_CONFIG_COMPLETE_IND.value
                == response.msg_type
            )
            self._configuredUes.append(reqs[i])

        self.check_gw_health_after_ue_load()

    def configUEDevice_without_checking_gw_health(self, num_ues):
        """ Configure the device on the UE side """
        reqs = self._sub_util.add_sub(num_ues=num_ues)
        for i in range(num_ues):
            print(
                "************************* UE device config for ue_id ",
                reqs[i].ue_id,
            )
            assert (
                self._s1_util.issue_cmd(s1ap_types.tfwCmd.UE_CONFIG, reqs[i])
                == 0
            )
            response = self._s1_util.get_response()
            assert (
                s1ap_types.tfwCmd.UE_CONFIG_COMPLETE_IND.value
                == response.msg_type
            )
            # APN configuration below can be overwritten in the test case
            # after configuring UE device.
            self.configAPN(
                "IMSI" + "".join([str(j) for j in reqs[i].imsi]), None
            )
            self._configuredUes.append(reqs[i])

    def configAPN(self, imsi, apn_list, default=True):
        """ Configure the APN """
        # add a default APN to be used in attach requests
        if default:
            magma_default_apn = {
                "apn_name": "magma.ipv4",  # APN-name
                "qci": 9,  # qci
                "priority": 15,  # priority
                "pre_cap": 1,  # preemption-capability
                "pre_vul": 0,  # preemption-vulnerability
                "mbr_ul": 200000000,  # MBR UL
                "mbr_dl": 100000000,  # MBR DL
            }
            # APN list to be configured
            if apn_list is not None:
                apn_list.insert(0, magma_default_apn)
            else:
                apn_list = [magma_default_apn]
        self._sub_util.config_apn_data(imsi, apn_list)

    def check_gw_health_after_ue_load(self):
        """ Wait for the MME only after adding entries to HSS """
        if self.wait_gateway_healthy:
            self._gateway_services.wait_for_healthy_gateway()
            self.init_s1ap_tester()
            self.wait_gateway_healthy = False

    def configDownlinkTest(self, *ues, **kwargs):
        """ Set up an downlink test, returning a TrafficTest object

        Args:
            ues (s1ap_types.ueConfig_t): the UEs to test
            kwargs: the keyword args to pass into generate_downlink_test

        Returns: a TrafficTest object, the traffic test generated based on the
            given UEs
        """
        # Configure downlink route in TRF server
        assert self._trf_util.update_dl_route(self.TEST_IP_BLOCK)

        ips = self._getAddresses(*ues)
        for ip, ue in zip(ips, ues):
            if not ip:
                raise ValueError(
                    "Encountered invalid IP for UE ID %s."
                    " Are you sure the UE is attached?" % ue
                )
        return self._trf_util.generate_traffic_test(
            ips, is_uplink=False, **kwargs
        )

    def configUplinkTest(self, *ues, **kwargs):
        """ Set up an uplink test, returning a TrafficTest object

        Args:
            ues (s1ap_types.ueConfig_t): the UEs to test
            kwargs: the keyword args to pass into generate_uplink_test

        Returns: a TrafficTest object, the traffic test generated based on the
            given UEs
        """
        ips = self._getAddresses(*ues)
        for ip, ue in zip(ips, ues):
            if not ip:
                raise ValueError(
                    "Encountered invalid IP for UE ID %s."
                    " Are you sure the UE is attached?" % ue
                )
        return self._trf_util.generate_traffic_test(
            ips, is_uplink=True, **kwargs
        )

    def get_gateway_services_util(self):
        """ Not a property, so return object is callable """
        return self._gateway_services

    @property
    def ue_req(self):
        """ Get a configured UE """
        req = self._configuredUes[self._ue_idx]
        self._ue_idx += 1
        return req

    @property
    def s1_util(self):
        return self._s1_util

    @property
    def mobility_util(self):
        return self._mobility_util

    @property
    def traffic_util(self):
        return self._trf_util

    @property
    def magmad_util(self):
        return self._magmad_util

    def cleanup(self):
        time.sleep(0.5)
        print("************************* send SCTP SHUTDOWN")

        self._s1_util.issue_cmd(s1ap_types.tfwCmd.SCTP_SHUTDOWN_REQ, None)
        self._s1_util.cleanup()
        self._sub_util.cleanup()
        self._trf_util.cleanup()
        self._mobility_util.cleanup()

        time.sleep(5)
        print("************************* going to do redis :)")
        keys = self.redis_client.keys()
        print(keys)

        for key in keys:
            print("\n")
            self.parse(str(key, 'utf-8'), str(self.redis_client.type(key), 'utf-8'))

        # Cloud cleanup needs to happen after cleanup for
        # subscriber util and mobility util
        # if self._test_cloud:
        #    self._cloud_manager.clean_up()

    def multiEnbConfig(self, num_of_enbs, enb_list=None):
        if enb_list is None:
            enb_list = []
        req = s1ap_types.multiEnbConfigReq_t()
        req.numOfEnbs = num_of_enbs
        # ENB Parameter column index initialization
        CELLID_COL_IDX = 0
        TAC_COL_IDX = 1
        ENBTYPE_COL_IDX = 2
        PLMNID_COL_IDX = 3
        PLMN_LENGTH_IDX = 4

        for idx1 in range(num_of_enbs):
            req.multiEnbCfgParam[idx1].cell_id = enb_list[idx1][CELLID_COL_IDX]
            req.multiEnbCfgParam[idx1].tac = enb_list[idx1][TAC_COL_IDX]
            req.multiEnbCfgParam[idx1].enbType = enb_list[idx1][
                ENBTYPE_COL_IDX
            ]
            req.multiEnbCfgParam[idx1].plmn_length = enb_list[idx1][
                PLMN_LENGTH_IDX
            ]
            for idx2 in range(req.multiEnbCfgParam[idx1].plmn_length):
                val = enb_list[idx1][PLMNID_COL_IDX][idx2]
                req.multiEnbCfgParam[idx1].plmn_id[idx2] = int(val)

        print("***************** Sending Multiple Enb Config Request\n")
        assert (
            self._s1_util.issue_cmd(
                s1ap_types.tfwCmd.MULTIPLE_ENB_CONFIG_REQ, req
            )
            == 0
        )

    def sendActDedicatedBearerAccept(self, ue_id, bearerId):
        act_ded_bearer_acc = s1ap_types.UeActDedBearCtxtAcc_t()
        act_ded_bearer_acc.ue_Id = ue_id
        act_ded_bearer_acc.bearerId = bearerId
        self._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ACT_DED_BER_ACC, act_ded_bearer_acc
        )
        print(
            "************** Sending activate dedicated EPS bearer "
            "context accept\n"
        )

    def sendDeactDedicatedBearerAccept(self, ue_id, bearerId):
        deact_ded_bearer_acc = s1ap_types.UeDeActvBearCtxtAcc_t()
        deact_ded_bearer_acc.ue_Id = ue_id
        deact_ded_bearer_acc.bearerId = bearerId
        self._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_DEACTIVATE_BER_ACC, deact_ded_bearer_acc
        )
        print("************* Sending deactivate EPS bearer context accept\n")

    def sendPdnConnectivityReq(
        self, ue_id, apn, pdn_type=1, pcscf_addr_type=None, dns_ipv6_addr=False
    ):
        req = s1ap_types.uepdnConReq_t()
        req.ue_Id = ue_id
        # Initial Request
        req.reqType = 1
        req.pdnType_pr.pres = 1
        # PDN Type 1 = IPv4, 2 = IPv6, 3 = IPv4v6
        req.pdnType_pr.pdn_type = pdn_type
        req.pdnAPN_pr.pres = 1
        req.pdnAPN_pr.len = len(apn)
        req.pdnAPN_pr.pdn_apn = (ctypes.c_ubyte * 100)(
            *[ctypes.c_ubyte(ord(c)) for c in apn[:100]]
        )
        print("********* PDN type", pdn_type)
        # Populate PCO if pcscf_addr_type is set
        if pcscf_addr_type or dns_ipv6_addr:
            self._s1_util.populate_pco(
                req.protCfgOpts_pr, pcscf_addr_type, dns_ipv6_addr
            )

        self.s1_util.issue_cmd(s1ap_types.tfwCmd.UE_PDN_CONN_REQ, req)

        print("************* Sending Standalone PDN Connectivity Request\n")

    def parse(self, key: str, redis_type: str):
        """
        Parse value of redis key on redis for encoded HASH, SET types, or
        JSON / Protobuf encoded state-wrapped types and prints it

        Args:
            key: key on redis

        """
        print(key, redis_type)
        key_type = key
        if ":" in key:
            key_type = key.split(":")[1]
        key_type.strip('\'')
        if 'hash' in redis_type:
            print("in hash!")
            deserializer = STATE_DESERIALIZERS.get(key_type)
            if not deserializer:
                print("no desreializer1 \n")
                return
            self._parse_hash_type(deserializer, key)
        elif 'set' in redis_type:
            print("in set!")
            deserializer = STATE_DESERIALIZERS.get(key_type)
            if not deserializer:
                print("no desreializer2 \n")
                return
            self._parse_set_type(deserializer, key)
        else:
            print("in else :(!")
            value =self.redis_client.get(key)
            # Try parsing as json first, if there's decoding error, parse proto
            try:
                self._parse_state_json(value)
            except Exception as e:
                print(e)
                try:
                    self._parse_state_proto(key_type, value)
                except Exception as e:
                    print(e)
                    print(value)

    def _parse_state_json(self, value):
        if value:
            deserializer = get_json_deserializer()
            value = json.loads(jsonpickle.encode(deserializer(value)))
            print(json.dumps(value, indent=2, sort_keys=True))
        else:
            raise AttributeError('Key not found on redis')

    def _parse_state_proto(self, key_type, value):
        proto = STATE_PROTOS.get(key_type.lower())
        if proto:
            deserializer = get_proto_deserializer(proto)
            print(deserializer(value))
        else:
            print("no key found \n")
            return

    def _parse_set_type(self, deserializer, key):
        set_values = self.redis_client.smembers(key)
        for value in set_values:
            print(deserializer(value))

    def _parse_hash_type(self, deserializer, key):
        value = self.redis_client.hgetall(key)
        for key, val in value.items():
            print(key.decode('utf-8'))
            print(deserializer(val))