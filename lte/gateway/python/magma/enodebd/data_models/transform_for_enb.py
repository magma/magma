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

from magma.enodebd.exceptions import ConfigurationError

CELL_RESERVED_MAP = {
    True: 'reserved',
    False: 'notReserved',
}


INVERT_CELL_RESERVED_MAP = {
    True: 'notReserved',
    False: 'reserved',
}


def admin_state(flag):
    return 'UP' if flag else 'DOWN'


def cell_reserved(value):
    return CELL_RESERVED_MAP.get(value)


def invert_cell_reserved(value):
    """
    We need to handle Baicells bug which inverts the meaning of 'cell reserved'
    """
    return INVERT_CELL_RESERVED_MAP.get(value)


def invert_cell_barred(value: bool):
    """
    We need to handle Baicells bug which inverts the meaning of 'cell barred'
    """
    return not value


def bandwidth(bandwidth_mhz):
    """
    Map bandwidth in MHz to number of RBs
    TODO: TR-196 spec says this should be '6' rather than 'n6', but
    BaiCells eNodeB uses 'n6'. Need to resolve this.

    Args:
        bandwidth_mhz (int): Bandwidth in MHz
    Returns:
        str: Bandwidth in RBS
    """
    if bandwidth_mhz == 1.4:
        bandwidth_rbs = 'n6'
    elif bandwidth_mhz == 3:
        bandwidth_rbs = 'n15'
    elif bandwidth_mhz == 5:
        bandwidth_rbs = 'n25'
    elif bandwidth_mhz == 10:
        bandwidth_rbs = 'n50'
    elif bandwidth_mhz == 15:
        bandwidth_rbs = 'n75'
    elif bandwidth_mhz == 20:
        bandwidth_rbs = 'n100'
    else:
        raise ConfigurationError(
            'Unknown bandwidth_mhz (%s)' %
            str(bandwidth_mhz),
        )
    return bandwidth_rbs
