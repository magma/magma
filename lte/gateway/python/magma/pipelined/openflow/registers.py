"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
from enum import IntEnum

# Register names
# Global registers:
IMSI_REG = 'metadata'
DIRECTION_REG = 'reg1'
DPI_REG = 'reg10'
TEST_PACKET_REG = 'reg5'
PASSTHROUGH_REG = 'reg6'
VLAN_TAG_REG = 'reg7'

# xxreg3 allow us to specify 16 bytes vakue to describe APN
# according to http://man7.org/linux/man-pages/man7/ovs-fields.7.html 
# xxreg3 will allocate [xreg2, xreg3] or [reg12, reg13, reg14, reg15] 
# under the hood
APN_TAG_REG = 'xxreg3' 

XXREGISTERS_MAP = {
    'xxreg0': ['reg0', 'reg1', 'reg2', 'reg3'],
    'xxreg1': ['reg4', 'reg5', 'reg6', 'reg7'],
    'xxreg2': ['reg8', 'reg9', 'reg10', 'reg11'],
    'xxreg3': ['reg12', 'reg13', 'reg14', 'reg15']
}

# Local scratch registers (These registers are reset when submitting to
# another app):
SCRATCH_REGS = ['reg0']
RULE_VERSION_REG = 'reg4'

# Register values
REG_ZERO_VAL = 0x0
PASSTHROUGH_REG_VAL = 0x1


class Direction(IntEnum):
    """
    Direction bits for direction reg
    """
    OUT = 0x01
    IN = 0x10


class TestPacket(IntEnum):
    ON = 0x1
    OFF = 0x0


def load_passthrough(parser, passthrough=PASSTHROUGH_REG_VAL):
    """
    Wrapper for loading the direction register
    """
    return parser.NXActionRegLoad2(dst=PASSTHROUGH_REG, value=passthrough)


def load_direction(parser, direction: Direction):
    """
    Wrapper for loading the direction register
    """
    if not is_valid_direction(direction):
        raise Exception("Invalid direction")
    return parser.NXActionRegLoad2(dst=DIRECTION_REG, value=direction.value)


def is_valid_direction(direction: Direction):
    return isinstance(direction, Direction)


def load_trace_packet(parser, test: TestPacket):
    """
    Wrapper for loading the test-packet register
    """
    if not isinstance(test, TestPacket):
        raise Exception('Invalid test object')
    return parser.NXActionRegLoad2(dst=TEST_PACKET_REG, value=test.value)
