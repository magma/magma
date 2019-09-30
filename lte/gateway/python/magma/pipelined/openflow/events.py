"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""

import copy

from ryu.controller import event


class EventSendPacket(event.EventBase):
    def __init__(self, pkt, imsi=None):
        super().__init__()
        self.packet = copy.copy(pkt)
        self.imsi = imsi
