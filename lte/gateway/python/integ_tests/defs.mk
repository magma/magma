# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

PROTO_LIST:=orc8r_protos lte_protos feg_protos

# Add the s1aptester integration tests
MANDATORY_TESTS = s1aptests/test_modify_mme_config_for_sanity.py \
s1aptests/test_attach_detach.py \
s1aptests/test_gateway_metrics_attach_detach.py \
s1aptests/test_attach_detach_looped.py  \
s1aptests/test_attach_emergency.py \
s1aptests/test_attach_combined_eps_imsi.py \
s1aptests/test_attach_via_guti.py \
s1aptests/test_attach_without_ips_available.py \
s1aptests/test_attach_detach_after_ue_context_release.py \
s1aptests/test_attach_detach_duplicate_nas_resp_messages.py \
s1aptests/test_attach_security_mode_reject.py \
s1aptests/test_attach_esm_information.py \
s1aptests/test_attach_esm_information_wrong_apn.py \
s1aptests/test_attach_ue_ctxt_release_cmp_delay.py \
s1aptests/test_attach_auth_failure.py \
s1aptests/test_nas_non_delivery_for_smc.py \
s1aptests/test_nas_non_delivery_for_identity_req.py \
s1aptests/test_attach_no_initial_context_resp.py \
s1aptests/test_attach_detach_no_ueContext_release_comp.py \
s1aptests/test_no_attach_complete.py \
s1aptests/test_attach_detach_with_mobilityd_restart.py \
s1aptests/test_attach_detach_multiple_ip_blocks_mobilityd_restart.py \
s1aptests/test_idle_mode_with_mme_restart.py \
s1aptests/test_attach_ul_udp_data_with_mme_restart.py \
s1aptests/test_attach_ul_udp_data_with_mobilityd_restart.py \
s1aptests/test_attach_ul_udp_data_with_multiple_service_restart.py \
s1aptests/test_attach_ul_udp_data_with_pipelined_restart.py \
s1aptests/test_attach_ul_udp_data_with_sessiond_restart.py \
s1aptests/test_service_req_ul_udp_data_with_mme_restart.py \
s1aptests/test_restore_mme_config_after_sanity.py

# Enable these tests once the CI job time-out has increased
# s1aptests/test_mobile_reachability_timer_with_mme_restart.py \
# s1aptests/test_implicit_detach_timer_with_mme_restart.py \
# s1aptests/test_ics_timer_expiry_with_mme_restart.py \

# These test cases pass without memory leaks, but needs DL-route in TRF server
# sudo /sbin/route add -net 192.168.128.0 gw 192.168.60.142
#     netmask 255.255.255.0 dev eth1
# s1aptests/test_attach_dl_udp_data.py \
# s1aptests/test_attach_dl_tcp_data.py \
# s1aptests/test_attach_detach_attach_dl_tcp_data.py

# TODO flaky tests we should look at
# s1aptests/test_attach_detach_ps_service_not_available.py \
# s1aptests/test_enb_complete_reset.py \
# s1aptests/test_attach_detach_multi_ue_looped.py \

CLOUD_TESTS = cloud_tests/checkin_test.py \
cloud_tests/metrics_export_test.py \
cloud_tests/config_test.py

S1AP_TESTER_CFG=$(MAGMA_ROOT)/lte/gateway/python/integ_tests/data/s1ap_tester_cfg
S1AP_TESTER_PYTHON_PATH=$(S1AP_TESTER_ROOT)/bin

# Local integ tests are run on the magma access gateway, not the test VM
LOCAL_INTEG_TESTS = gxgy_tests
