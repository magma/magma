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
import textwrap
from typing import Optional, Union

from magma.enodebd.exceptions import ConfigurationError
from magma.enodebd.logger import EnodebdLogger as logger

DUPLEX_MAP = {
    '01': 'TDDMode',
    '02': 'FDDMode',
}

BANDWIDTH_RBS_TO_MHZ_MAP = {
    'n6': 1.4,
    'n15': 3,
    'n25': 5,
    'n50': 10,
    'n75': 15,
    'n100': 20,
}

BANDWIDTH_MHZ_LIST = {1.4, 3, 5, 10, 15, 20}


def duplex_mode(value: str) -> Optional[str]:
    return DUPLEX_MAP.get(value)


def band_capability(value: str) -> str:
    return ','.join([str(int(b, 16)) for b in textwrap.wrap(value, 2)])


def gps_tr181(value: str) -> str:
    """Convert GPS value (lat or lng) to float

    Per TR-181 specification, coordinates are returned in degrees,
    multiplied by 1,000,000.

    Args:
        value (string): GPS value (latitude or longitude)
    Returns:
        str: GPS value (latitude/longitude) in degrees
    """
    try:
        return str(float(value) / 1e6)
    except Exception:  # pylint: disable=broad-except
        return value


def bandwidth(bandwidth_rbs: Union[str, int, float]) -> float:
    """
    Map bandwidth in number of RBs to MHz
    TODO: TR-196 spec says this should be '6' rather than 'n6', but
    BaiCells eNodeB uses 'n6'. Need to resolve this.

    Args:
        bandwidth_rbs (str): Bandwidth in number of RBs
    Returns:
        str: Bandwidth in MHz
    """
    if bandwidth_rbs in BANDWIDTH_RBS_TO_MHZ_MAP:
        return BANDWIDTH_RBS_TO_MHZ_MAP[bandwidth_rbs]

    logger.warning('Unknown bandwidth_rbs (%s)', str(bandwidth_rbs))
    if bandwidth_rbs in BANDWIDTH_MHZ_LIST:
        return bandwidth_rbs
    elif isinstance(bandwidth_rbs, str):
        mhz = None
        if bandwidth_rbs.isdigit():
            mhz = int(bandwidth_rbs)
        elif bandwidth_rbs.replace('.', '', 1).isdigit():
            mhz = float(bandwidth_rbs)
        if mhz in BANDWIDTH_MHZ_LIST:
            return mhz
    raise ConfigurationError(
        'Unknown bandwidth specification (%s)' %
        str(bandwidth_rbs),
    )
