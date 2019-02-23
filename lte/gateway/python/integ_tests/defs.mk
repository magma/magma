# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
#

PROTO_LIST:=orc8r_protos lte_protos feg_protos

# Add the s1aptester integration tests
MANDATORY_TESTS = s1aptests/test_attach_detach.py \
s1aptests/test_gateway_metrics_attach_detach.py \
s1aptests/test_attach_detach_multi_ue.py \
s1aptests/test_attach_detach_looped.py  \
s1aptests/test_attach_emergency.py \
s1aptests/test_attach_combined_eps_imsi.py \
s1aptests/test_attach_via_guti.py \
s1aptests/test_attach_without_ips_available.py \
s1aptests/test_attach_detach_after_ue_context_release.py \
s1aptests/test_no_attach_complete.py \
s1aptests/test_no_auth_response.py \
s1aptests/test_no_security_mode_complete.py \
s1aptests/test_tau_periodic_inactive.py \
s1aptests/test_tau_periodic_active.py \
s1aptests/test_attach_service.py \
s1aptests/test_attach_detach_service.py \
s1aptests/test_attach_service_ue_radio_capability.py \
s1aptests/test_attach_service_multi_ue.py \
s1aptests/test_attach_ipv4v6_pdn_type.py \
s1aptests/test_service_info.py \
s1aptests/test_attach_detach_with_ovs.py \
s1aptests/test_resync.py \
s1aptests/test_standalone_pdn_conn_req.py
#s1aptests/test_attach_ul_udp_data.py \
#s1aptests/test_attach_ul_tcp_data.py \

# TODO Disabled because MME wont run without UEs in HSS
#s1aptests/test_attach_missing_imsi.py \

# TODO flaky tests we should look at
# s1aptests/test_attach_detach_multi_ue_looped.py \

CLOUD_TESTS = cloud_tests/checkin_test.py \
cloud_tests/metrics_export_test.py \
cloud_tests/config_test.py

S1AP_TESTER_CFG=$(MAGMA_ROOT)/lte/gateway/python/integ_tests/data/s1ap_tester_cfg
S1AP_PYTHON_PATH=$(S1AP_ROOT)/bin

# Local integ tests are run on the magma access gateway, not the test VM
LOCAL_INTEG_TESTS = gxgy_tests
