"""
Copyright (c) 2017-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""


def run_tcp_downlink(ue, s1ap_wrapper, duration=1):
    """
    Given a configured UE, run a TCP downlink test for the specified duration.
    """
    print("************************* Running UE downlink (TCP) for UE id ",
          ue.ue_id)
    with s1ap_wrapper.configDownlinkTest(ue, duration=duration) as test:
        test.verify()


def run_tcp_uplink(ue, s1ap_wrapper, duration=1):
    """
    Given a configured UE, run a TCP uplink test for the specified duration.
    """
    print("************************* Running UE uplink (TCP) for UE id ",
          ue.ue_id)
    with s1ap_wrapper.configUplinkTest(ue, duration=duration) as test:
        test.verify()


def run_udp_downlink(ue, s1ap_wrapper, duration=1):
    """
    Given a configured UE, run a UDP downlink test for the specified duration.
    """
    print("************************* Running UE downlink (UDP) for UE id ",
          ue.ue_id)
    with s1ap_wrapper.configDownlinkTest(ue, duration=duration,
                                         is_udp=True) as test:
        test.verify()


def run_udp_uplink(ue, s1ap_wrapper, duration=1):
    """
    Given a configured UE, run a UDP uplink test for the specified duration.
    """
    print("************************* Running UE uplink (UDP) for UE id ",
          ue.ue_id)
    with s1ap_wrapper.configUplinkTest(ue, duration=duration,
                                       is_udp=True) as test:
        test.verify()
