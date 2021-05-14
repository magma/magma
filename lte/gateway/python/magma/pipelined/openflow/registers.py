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
from enum import IntEnum

from magma.pipelined.imsi import encode_imsi

# Register names
# Global registers:
IMSI_REG = 'metadata'
DIRECTION_REG = 'reg1'
RULE_NUM_REG = 'reg2'
DPI_REG = 'reg10'
TEST_PACKET_REG = 'reg5'
PASSTHROUGH_REG = 'reg6'
VLAN_TAG_REG = 'reg7'
TUN_PORT_REG = 'reg8'
INGRESS_TUN_ID_REG = 'reg9'
PROXY_TAG_REG = 'reg10'

# Local scratch registers (These registers are reset when submitting to
# another app):
SCRATCH_REGS = ['reg0', 'reg3']
RULE_VERSION_REG = 'reg4'

# Register values
REG_ZERO_VAL = 0x0
PASSTHROUGH_REG_VAL = 0x1

# values for PROXY_TAG_REG
PROXY_TAG_TO_PROXY = 0x1


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

def load_imsi(parser, imsi):
    """
    Wrapper for loading the direction register
    """
    return parser.NXActionRegLoad2(dst=IMSI_REG, value=encode_imsi(imsi))


def set_in_port(parser, port_no):
    """
    Wrapper for loading the direction register
    """

    return parser.NXActionRegLoad2(dst='in_port', value=port_no)


def set_proxy_tag(parser, value=PROXY_TAG_TO_PROXY):
    """
    Wrapper for setting proxy flow tag.
    """
    return parser.NXActionRegLoad2(dst=PROXY_TAG_REG, value=value)


def set_tun_id(parser, tun_id:str):
    """
    Wrapper for setting proxy flow tag.
    """
    return parser.OFPActionSetField(tunnel_id=tun_id)


def is_valid_direction(direction: Direction):
    return isinstance(direction, Direction)


def load_trace_packet(parser, test: TestPacket):
    """
    Wrapper for loading the test-packet register
    """
    if not isinstance(test, TestPacket):
        raise Exception('Invalid test object')
    return parser.NXActionRegLoad2(dst=TEST_PACKET_REG, value=test.value)
