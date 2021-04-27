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
from typing import Optional

from magma.pipelined.openflow.registers import (
    DIRECTION_REG,
    DPI_REG,
    IMSI_REG,
    PASSTHROUGH_REG,
    PROXY_TAG_REG,
    RULE_NUM_REG,
    RULE_VERSION_REG,
    INGRESS_TUN_ID_REG,
    VLAN_TAG_REG,
    Direction,
    is_valid_direction,
)


class MagmaMatch(object):
    """
    MagmaMatch is a wrapper class for Ryu's OFPMatch. It provides enforcement
    and validation for matching against global registers such as imsi and
    direction.
    """

    def __init__(self, imsi: int = None, direction: Optional[Direction] = None,
                 rule_num: int = None, rule_version: int = None,
                 passthrough: int = None, vlan_tag: int = None,
                 app_id: int = None, proxy_tag: int = None,
                 teid: int = None, **kwargs):
        self.imsi = imsi
        self.direction = direction
        self.rule_num = rule_num
        self.rule_version = rule_version
        self.passthrough = passthrough
        self.vlan_tag = vlan_tag
        self.app_id = app_id
        self.proxy_tag = proxy_tag
        self.teid = teid
        self._match_kwargs = kwargs
        self._check_args()

    def update(self, v):
        self._match_kwargs.update(v)
        self._check_args()

    @property
    def ryu_match(self):
        """
        Convert the MagmaMatch object into a dict.
        """
        ryu_match = self._match_kwargs.copy()
        if self.direction is not None:
            ryu_match[DIRECTION_REG] = self.direction.value
        if self.imsi is not None:
            ryu_match[IMSI_REG] = self.imsi
        if self.rule_num is not None:
            ryu_match[RULE_NUM_REG] = self.rule_num
        if self.rule_version is not None:
            ryu_match[RULE_VERSION_REG] = self.rule_version
        if self.passthrough is not None:
            ryu_match[PASSTHROUGH_REG] = self.passthrough
        if self.vlan_tag is not None:
            ryu_match[VLAN_TAG_REG] = self.vlan_tag
        if self.app_id is not None:
            ryu_match[DPI_REG] = self.app_id
        if self.proxy_tag is not None:
            ryu_match[PROXY_TAG_REG] = self.proxy_tag
        if self.teid is not None and self.teid != 0:
            ryu_match[INGRESS_TUN_ID_REG] = self.teid
        return ryu_match

    def _check_args(self):
        if self.direction is not None and \
                not is_valid_direction(self.direction):
            raise Exception("Invalid direction: %s" % self.direction)

        # Avoid double register sets to ease debuggability
        for k in self._match_kwargs:
            if k == DIRECTION_REG and self.direction:
                raise Exception("Register %s should not be directly set" % k)
            if k == IMSI_REG and self.imsi:
                raise Exception("Register %s should not be directly set" % k)
