"""
Copyright (c) 2017-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import s1ap_types


def attach_ue(ue, s1ap_wrapper):
    print("************************* Running End to End attach for ",
          "UE id ", ue.ue_id)
    s1ap_wrapper.s1_util.attach(
        ue.ue_id, s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
        s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
        s1ap_types.ueAttachAccept_t)
    # Wait on EMM Information from MME
    s1ap_wrapper.s1_util.receive_emm_info()


def detach_ue(ue, s1ap_wrapper):
    print("************************* Running UE detach for UE id ",
          ue.ue_id)
    s1ap_wrapper.s1_util.detach(
        ue.ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
        True)
