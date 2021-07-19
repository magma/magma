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

from lte.protos.session_manager_pb2 import (
    PolicyReAuthAnswer,
    PolicyReAuthRequest,
    ReAuthResult,
    SessionRules,
)
from orc8r.protos.common_pb2 import Void


class MockLocalSessionManagerStub:
    """
    This Mock LocalSessionManagerStub will always respond with a Void
    """

    def __init__(self):
        pass

    def SetSessionRules(self, _: SessionRules, timeout: float) -> Void:
        return Void()


class MockSessionProxyResponderStub1:
    """
    This Mock SessionProxyResponderStub will always respond with a success to
    a received RAR
    """

    def __init__(self):
        pass

    def PolicyReAuth(self, _: PolicyReAuthRequest) -> PolicyReAuthAnswer:
        return PolicyReAuthAnswer(
            result=ReAuthResult.Value('UPDATE_INITIATED'),
        )


class MockSessionProxyResponderStub2:
    """
    This Mock SessionProxyResponderStub will always respond with a failure to
    a received RAR
    """

    def __init__(self):
        pass

    def PolicyReAuth(self, _: PolicyReAuthRequest) -> PolicyReAuthAnswer:
        return PolicyReAuthAnswer(
            result=ReAuthResult.Value('OTHER_FAILURE'),
        )


class MockSessionProxyResponderStub3:
    """
    This Mock SessionProxyResponderStub will always fail to install rule p2
    """

    def __init__(self):
        pass

    def PolicyReAuth(self, _: PolicyReAuthRequest) -> PolicyReAuthAnswer:
        return PolicyReAuthAnswer(
            result=ReAuthResult.Value('UPDATE_INITIATED'),
            failed_rules={
                "p2": PolicyReAuthAnswer.FailureCode.Value("UNKNOWN_RULE_NAME"),
            },
        )
