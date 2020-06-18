"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
from typing import Optional

from magma.pipelined.openflow.registers import IMSI_REG, DIRECTION_REG, \
    is_valid_direction, Direction, RULE_VERSION_REG, PASSTHROUGH_REG, \
    VLAN_TAG_REG, DPI_REG


class MagmaMatch(object):
    """
    MagmaMatch is a wrapper class for Ryu's OFPMatch. It provides enforcement
    and validation for matching against global registers such as imsi and
    direction.
    """

    def __init__(self, imsi: int = None, direction: Optional[Direction] = None,
                 rule_version: int = None, passthrough: int = None,
                 vlan_tag: int = None, app_id: int = None, **kwargs):
        self.imsi = imsi
        self.direction = direction
        self.rule_version = rule_version
        self.passthrough = passthrough
        self.vlan_tag = vlan_tag
        self.app_id = app_id
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
        if self.rule_version is not None:
            ryu_match[RULE_VERSION_REG] = self.rule_version
        if self.passthrough is not None:
            ryu_match[PASSTHROUGH_REG] = self.passthrough
        if self.vlan_tag is not None:
            ryu_match[VLAN_TAG_REG] = self.vlan_tag
        if self.app_id is not None:
            ryu_match[DPI_REG] = self.app_id
        return ryu_match

    def _check_args(self):
        if self.direction is not None and \
                not is_valid_direction(self.direction):
            raise Exception("Invalid direction: %s" % self.direction)

        for k in self._match_kwargs:
            if k in [DIRECTION_REG, IMSI_REG]:
                raise Exception("Register %s should not be directly set" % k)
