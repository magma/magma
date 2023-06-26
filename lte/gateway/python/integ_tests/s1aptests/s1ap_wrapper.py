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
import inspect
import os
import re
import time
from typing import List

import s1ap_types
from integ_tests.common.magmad_client import MagmadServiceGrpc
from integ_tests.common.mobility_service_client import MobilityServiceGrpc
from integ_tests.common.service303_utils import GatewayServicesUtil
from integ_tests.common.subscriber_db_client import (
    HSSGrpc,
    SubscriberDbCassandra,
    SubscriberDbGrpc,
)
from integ_tests.s1aptests.s1ap_utils import (
    MagmadUtil,
    MobilityUtil,
    S1ApUtil,
    SubscriberUtil,
)
from integ_tests.s1aptests.util.traffic_util import TrafficUtil


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
    TEST_CASE_EXECUTION_COUNT = 0
    TEST_ERROR_TRACEBACKS: List[str] = []

    def __init__(
        self,
        stateless_mode=MagmadUtil.stateless_cmds.ENABLE,
        apn_correction=MagmadUtil.apn_correction_cmds.DISABLE,
        health_service=MagmadUtil.health_service_cmds.DISABLE,
        federated_mode=False,
        ip_version=4,
    ):
        """
        Initialize the various classes required by the tests and setup.
        """
        t = time.localtime()
        current_time = time.strftime("%H:%M:%S", t)
        if TestWrapper.TEST_CASE_EXECUTION_COUNT != 0:
            print("\n**Running the test case again to identify flaky behavior")
        TestWrapper.TEST_CASE_EXECUTION_COUNT += 1
        print(
            "\nTest Case Execution Count:",
            TestWrapper.TEST_CASE_EXECUTION_COUNT,
            "[Start time: " + str(current_time) + "]",
        )

        federated_mode = (os.environ.get("FEDERATED_MODE") == "True")
        print(
            f"\n *** Running the test in {'Non-' if not federated_mode else ''}"
            "Federated Mode\n",
        )

        if self._test_oai_upstream:
            subscriber_client = SubscriberDbCassandra()
            self.wait_gateway_healthy = False
        elif federated_mode:
            subscriber_client = HSSGrpc()
            self.wait_gateway_healthy = True
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
        self._magmad_util.config_health_service(health_service)

        if not self._magmad_util.is_nat_enabled():
            self._magmad_util.enable_nat()

        self._magmad_util._wait_for_pipelined_to_initialize()

        self._s1_util = S1ApUtil()  # calls get_datapath, i.e., should run after we ensured pipelined is started
        self._enBConfig(ip_version)

        # gateway tests don't require restart, just wait for healthy now
        self._gateway_services = GatewayServicesUtil()
        if not self.wait_gateway_healthy:
            self.init_s1ap_tester()

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

    def _enBConfig(self, ip_version=4):
        """Configure the eNB in S1APTester"""
        # Using exaggerated prints makes the stdout easier to read.
        print("************************* Enb tester config")
        req = s1ap_types.FwNbConfigReq_t()
        req.cellId_pr.pres = True
        req.cellId_pr.cell_id = 10
        req.ip_version = ip_version
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
        """Perform S1 setup to the EPC"""
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
            ues: (list(s1ap_types.ueAppConfig_t)): the UEs whose IPs we want

        Returns:
            A list of ipaddress.ip_address objects, corresponding in order
            with the input UE parameters
        """
        return [self._s1_util.get_ip(ue.ue_id) for ue in ues]

    def configIpBlock(self):
        """ Remove any existing allocated blocks, then adds the ones used for
        testing """
        print("************************* Configuring IP block")
        self._mobility_util.remove_all_ip_blocks()
        self._mobility_util.add_ip_block(self.TEST_IP_BLOCK)
        print("************************* Waiting for IP changes to propagate")
        self._mobility_util.wait_for_changes()

    def configUEDevice(self, num_ues, req_data=None, static_ips=None):
        """ Configure the device on the UE side """
        if req_data is None:
            req_data = []
        if static_ips is None:
            static_ips = []
        reqs = self._sub_util.add_sub(num_ues=num_ues)
        for i in range(num_ues):
            print(
                "************************* UE device config for ue_id ",
                reqs[i].ue_id,
            )
            if req_data and bool(req_data[i]):
                if req_data[i].ueNwCap_pr.pres:
                    reqs[i].ueNwCap_pr.pres = req_data[i].ueNwCap_pr.pres
                    reqs[i].ueNwCap_pr.eea2_128 = req_data[
                        i
                    ].ueNwCap_pr.eea2_128
                    reqs[i].ueNwCap_pr.eea1_128 = req_data[
                        i
                    ].ueNwCap_pr.eea1_128
                    reqs[i].ueNwCap_pr.eea0 = req_data[i].ueNwCap_pr.eea0
                    reqs[i].ueNwCap_pr.eia2_128 = req_data[
                        i
                    ].ueNwCap_pr.eia2_128
                    reqs[i].ueNwCap_pr.eia1_128 = req_data[
                        i
                    ].ueNwCap_pr.eia1_128
                    reqs[i].ueNwCap_pr.eia0 = req_data[i].ueNwCap_pr.eia0

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
            if i < len(static_ips):
                ue_ip = static_ips[i]
            else:
                ue_ip = None
            self.configAPN(
                imsi="IMSI" + "".join([str(j) for j in reqs[i].imsi]),
                apn_list=None,
                default=True,
                static_ip=ue_ip,
            )
            self._configuredUes.append(reqs[i])

        self.check_gw_health_after_ue_load()

    def configUEDevice_ues_same_imsi(self, num_ues):
        """ Configure the device on the UE side with same IMSI and
        having different ue-id

        Args:
            num_ues: Count of UEs to be configured
        """
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
                "IMSI" + "".join([str(j) for j in reqs[i].imsi]),
                None,
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
                "IMSI" + "".join([str(j) for j in reqs[i].imsi]),
                None,
            )
            self._configuredUes.append(reqs[i])

    def configAPN(self, imsi, apn_list, default=True, static_ip=None):
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
                "static_ip": static_ip,
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

        Returns:
            A TrafficTest object, the traffic test generated based on the
            given UEs

        Raises:
            ValueError: If no valid IP is found
        """
        # Configure downlink route in TRF server
        assert self._trf_util.update_dl_route(self.TEST_IP_BLOCK)

        ips = self._getAddresses(*ues)
        for ip, ue in zip(ips, ues):
            if not ip:
                raise ValueError(
                    "Encountered invalid IP for UE ID %s."
                    " Are you sure the UE is attached?" % ue,
                )
        return self._trf_util.generate_traffic_test(
            ips,
            is_uplink=False,
            **kwargs,
        )

    def configMtuSize(self, set_mtu=False):
        """ Config MTU size for DL ipv6 data """
        assert self._trf_util.update_mtu_size(set_mtu)

    def configUplinkTest(self, *ues, **kwargs):
        """ Set up an uplink test, returning a TrafficTest object

        Args:
            ues (s1ap_types.ueConfig_t): the UEs to test
            kwargs: the keyword args to pass into generate_uplink_test

        Returns:
            A TrafficTest object, the traffic test generated based on the
            given UEs

        Raises:
            ValueError: If no valid IP is found
        """
        ips = self._getAddresses(*ues)
        for ip, ue in zip(ips, ues):
            if not ip:
                raise ValueError(
                    "Encountered invalid IP for UE ID %s."
                    " Are you sure the UE is attached?" % ue,
                )
        return self._trf_util.generate_traffic_test(
            ips,
            is_uplink=True,
            **kwargs,
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
        """Get s1_util instance"""
        return self._s1_util

    @property
    def mobility_util(self):
        """Get mobility_util instance"""
        return self._mobility_util

    @property
    def traffic_util(self):
        """Get traffic_util instance"""
        return self._trf_util

    @property
    def magmad_util(self):
        """Get magmad_util instance"""
        return self._magmad_util

    @classmethod
    def generate_flaky_summary(cls):
        """Print the flaky report summary"""
        if TestWrapper.TEST_ERROR_TRACEBACKS:
            print("\n===Flaky Test Report===\n")
            for traceback in TestWrapper.TEST_ERROR_TRACEBACKS:
                print(traceback)
            print("===End Flaky Test Report===")

    @classmethod
    def is_test_successful(cls, test) -> bool:
        """Get current test case execution status"""
        if test is None:
            test = inspect.currentframe().f_back.f_back.f_locals["self"]  # type: ignore
        if test is not None and hasattr(test, "_outcome"):
            result = test.defaultTestResult()
            test._feedErrorsToResult(result, test._outcome.errors)
            test_contains_error = (
                result.errors and result.errors[-1][0] is test
            )
            test_contains_failure = (
                result.failures and result.failures[-1][0] is test
            )
            if test_contains_error or test_contains_failure:
                TestWrapper.TEST_ERROR_TRACEBACKS.append(
                    str(test)
                    + " failed (Execution Count: "
                    + str(TestWrapper.TEST_CASE_EXECUTION_COUNT)
                    + ")."
                    + result.failures[0][1]
                    if test_contains_failure
                    else result.errors[0][1],
                )
            return not (test_contains_error or test_contains_failure)
        return False

    def cleanup(self, test=None):
        """Cleanup test setup after testcase execution"""
        is_test_successful = self.is_test_successful(test)
        time.sleep(0.5)
        print("************************* send SCTP SHUTDOWN")
        self._s1_util.issue_cmd(s1ap_types.tfwCmd.SCTP_SHUTDOWN_REQ, None)

        print("************************* Cleaning up TFW")
        self._s1_util.issue_cmd(s1ap_types.tfwCmd.TFW_CLEANUP, None)

        if not is_test_successful:
            self._s1_util.delete_ovs_flow_rules()

        self._s1_util.cleanup()
        self._sub_util.cleanup()
        self._trf_util.cleanup()
        self._mobility_util.cleanup()
        self._magmad_util.print_redis_state()

        if not is_test_successful:
            print("The test has failed. Restarting Sctpd for cleanup")
            self.magmad_util.restart_services(['sctpd'], wait_time=30)
            self.magmad_util.print_redis_state()
            if TestWrapper.TEST_CASE_EXECUTION_COUNT == 3:
                self.generate_flaky_summary()

        elif TestWrapper.TEST_CASE_EXECUTION_COUNT > 1:
            self.generate_flaky_summary()

        if not self.magmad_util.is_redis_empty():
            print("************************* Redis not empty, initiating cleanup")
            self.magmad_util.restart_services(['sctpd'], wait_time=30)
            self.magmad_util.print_redis_state()

    def multiEnbConfig(self, num_of_enbs, enb_list=None):
        """Configure multiple eNB in S1APTester"""
        if enb_list is None:
            enb_list = []
        req = s1ap_types.multiEnbConfigReq_t()
        req.numOfEnbs = num_of_enbs
        # ENB Parameter column index initialization
        cellid_col_idx = 0
        tac_col_idx = 1
        enbtype_col_idx = 2
        plmnid_col_idx = 3
        plmn_length_idx = 4

        for idx1 in range(num_of_enbs):
            req.multiEnbCfgParam[idx1].cell_id = enb_list[idx1][cellid_col_idx]
            req.multiEnbCfgParam[idx1].tac = enb_list[idx1][tac_col_idx]
            req.multiEnbCfgParam[idx1].enbType = enb_list[idx1][
                enbtype_col_idx
            ]
            req.multiEnbCfgParam[idx1].plmn_length = enb_list[idx1][
                plmn_length_idx
            ]
            for idx2 in range(req.multiEnbCfgParam[idx1].plmn_length):
                val = enb_list[idx1][plmnid_col_idx][idx2]
                req.multiEnbCfgParam[idx1].plmn_id[idx2] = int(val)

        print("***************** Sending Multiple Enb Config Request\n")
        assert (
            self._s1_util.issue_cmd(
                s1ap_types.tfwCmd.MULTIPLE_ENB_CONFIG_REQ,
                req,
            )
            == 0
        )

    def sendActDedicatedBearerAccept(self, ue_id, bearer_id):
        """Send Activate Dedicated Bearer Accept message"""
        act_ded_bearer_acc = s1ap_types.UeActDedBearCtxtAcc_t()
        act_ded_bearer_acc.ue_Id = ue_id
        act_ded_bearer_acc.bearerId = bearer_id
        self._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ACT_DED_BER_ACC,
            act_ded_bearer_acc,
        )
        print(
            "************** Sending activate dedicated EPS bearer "
            "context accept\n",
        )

    def sendDeactDedicatedBearerAccept(self, ue_id, bearer_id):
        """Send Deactivate Dedicated Bearer Accept message"""
        deact_ded_bearer_acc = s1ap_types.UeDeActvBearCtxtAcc_t()
        deact_ded_bearer_acc.ue_Id = ue_id
        deact_ded_bearer_acc.bearerId = bearer_id
        self._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_DEACTIVATE_BER_ACC,
            deact_ded_bearer_acc,
        )
        print("************* Sending deactivate EPS bearer context accept\n")

    def sendPdnConnectivityReq(
        self,
        ue_id,
        apn,
        pdn_type=1,
        pcscf_addr_type=None,
        dns_ipv6_addr=False,
    ):
        """Send PDN Connectivity Request message"""
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
            *[ctypes.c_ubyte(ord(c)) for c in apn[:100]],
        )
        print("********* PDN type", pdn_type)
        # Populate PCO if pcscf_addr_type is set
        if pcscf_addr_type or dns_ipv6_addr:
            self._s1_util.populate_pco(
                req.protCfgOpts_pr,
                pcscf_addr_type,
                dns_ipv6_addr,
            )

        self.s1_util.issue_cmd(s1ap_types.tfwCmd.UE_PDN_CONN_REQ, req)

        print("************* Sending Standalone PDN Connectivity Request\n")

    def flush_arp(self):
        """Flush all ARP entries"""
        self._magmad_util.exec_command_output(
            "sudo ip neigh flush all",
        )
        print("magma-dev: ARP flushed")

    def enable_disable_ipv6_iface(self, cmd):
        """Enable or disable eth3 (ipv6) interface as nat_iface"""
        self._magmad_util.config_ipv6_iface(cmd)
