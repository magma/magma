# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# Add the s1aptester integration tests with federation gateway
FEDERATED_TESTS = s1aptests/test_attach_detach.py \
s1aptests/test_attach_detach_multi_ue.py \
s1aptests/test_attach_auth_failure.py \
s1aptests/test_no_auth_response.py \
s1aptests/test_nas_non_delivery_for_auth.py \
s1aptests/test_sctp_abort_after_auth_req.py \
s1aptests/test_sctp_shutdown_after_auth_req.py \
s1aptests/test_no_auth_resp_with_mme_restart_reattach.py \
s1aptests/test_send_error_ind_for_dl_nas_with_auth_req.py \
s1aptests/test_attach_auth_mac_failure.py \
s1aptests/test_no_auth_response_with_mme_restart.py \
s1aptests/test_no_security_mode_complete.py \
s1aptests/test_no_attach_complete.py \
s1aptests/test_attach_detach_security_algo_eea0_eia0.py \
s1aptests/test_attach_detach_security_algo_eea1_eia1.py \
s1aptests/test_attach_detach_security_algo_eea2_eia2.py \
s1aptests/test_attach_security_mode_reject.py \
s1aptests/test_attach_missing_imsi.py \
s1aptests/test_duplicate_attach.py \
s1aptests/test_attach_emergency.py \
s1aptests/test_attach_detach_after_ue_context_release.py \
s1aptests/test_attach_esm_information_wrong_apn.py \
s1aptests/test_attach_detach_secondary_pdn_invalid_apn.py \
s1aptests/test_standalone_pdn_conn_req_with_apn_correction.py


CLOUD_TESTS = cloud_tests/checkin_test.py \
cloud_tests/metrics_export_test.py \
cloud_tests/config_test.py

S1AP_TESTER_CFG=$(MAGMA_ROOT)/lte/gateway/python/integ_tests/data/s1ap_tester_cfg
S1AP_TESTER_PYTHON_PATH=$(S1AP_TESTER_ROOT)/bin
