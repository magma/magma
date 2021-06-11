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


def run_tcp_downlink(ue, s1ap_wrapper, duration=1):
    """
    Given a configured UE, run a TCP downlink test for the specified duration.
    """
    print(
        "************************* Running UE downlink (TCP) for UE id ",
        ue.ue_id,
    )
    with s1ap_wrapper.configDownlinkTest(ue, duration=duration) as test:
        test.verify()


def run_tcp_uplink(ue, s1ap_wrapper, duration=1):
    """
    Given a configured UE, run a TCP uplink test for the specified duration.
    """
    print(
        "************************* Running UE uplink (TCP) for UE id ",
        ue.ue_id,
    )
    with s1ap_wrapper.configUplinkTest(ue, duration=duration) as test:
        test.verify()


def run_udp_downlink(ue, s1ap_wrapper, duration=1):
    """
    Given a configured UE, run a UDP downlink test for the specified duration.
    """
    print(
        "************************* Running UE downlink (UDP) for UE id ",
        ue.ue_id,
    )
    with s1ap_wrapper.configDownlinkTest(
        ue, duration=duration,
        is_udp=True,
    ) as test:
        test.verify()


def run_udp_uplink(ue, s1ap_wrapper, duration=1):
    """
    Given a configured UE, run a UDP uplink test for the specified duration.
    """
    print(
        "************************* Running UE uplink (UDP) for UE id ",
        ue.ue_id,
    )
    with s1ap_wrapper.configUplinkTest(
        ue, duration=duration,
        is_udp=True,
    ) as test:
        test.verify()
