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
from bisect import bisect_right
from typing import NamedTuple, Optional

from lte.protos.mconfig.mconfigs_pb2 import EnodebD
from magma.enodebd.device_config.cbrs_consts import (
    BAND48_LOW_FREQ_MHZ,
    BAND48_NOFFS_DL,
    SAS_BANDWIDTHS,
)
from magma.enodebd.exceptions import ConfigurationError


class EnodebConfig(NamedTuple):
    serial_num: str
    config: EnodebD.EnodebConfig


def get_enb_rf_tx_desired(mconfig: EnodebD, enb_serial: str) -> bool:
    """
    Get transmit enabled for eNB.

    Args:
        mconfig: enodebd mconfig
        enb_serial: eNB serial number

    Raises:
        KeyError: when eNB is missing from enodebd mconfig

    Returns:
        True if the mconfig specifies to enable transmit on the eNB
    """
    mconfig_serials_no = len(mconfig.enb_configs_by_serial)
    if mconfig.enb_configs_by_serial is not None and \
            mconfig_serials_no > 0:
        if enb_serial in mconfig.enb_configs_by_serial:
            enb_config = mconfig.enb_configs_by_serial.get(enb_serial)
            return enb_config.transmit_enabled
        raise KeyError('Missing eNB from mconfig: %s' % enb_serial)
    return mconfig.allow_enodeb_transmit


def is_enb_registered(mconfig: EnodebD, enb_serial: str) -> bool:
    """
    Check if enb is known to enodebd.

    Args:
        mconfig: enodebd mconfig
        enb_serial: eNB serial number

    Returns:
        True if either:
            - the eNodeB is registered by serial to the Access Gateway
            or
            - the Access Gateway accepts all eNodeB devices
    """
    mconfig_serials_no = len(mconfig.enb_configs_by_serial)
    if mconfig.enb_configs_by_serial is not None and \
            mconfig_serials_no > 0:
        return enb_serial in mconfig.enb_configs_by_serial
    return True


def find_enb_by_cell_id(mconfig: EnodebD, cell_id: int) \
        -> Optional[EnodebConfig]:
    """
    Find eNB by Cell ID

    Args:
        mconfig: enodebd mconfig
        cell_id: Cell ID

    Returns:
        eNB config if:
            - the eNodeB is registered by serial to the Access Gateway
            - cell ID is found in eNB status by serial
        else: returns None
    """
    mconfig_serials_no = len(mconfig.enb_configs_by_serial)
    if mconfig.enb_configs_by_serial is not None and \
            mconfig_serials_no > 0:
        for sn, enb in mconfig.enb_configs_by_serial.items():
            if cell_id == enb.cell_id:
                config = EnodebConfig(serial_num=sn, config=enb)
                return config
    return None


def calc_bandwidth_mhz(low_freq_hz: int, high_freq_hz: int) -> float:
    """
    Calculate bandwidth in mhz for CBRS

    Args:
        low_freq_hz: int, Low frequency limit taken from available channel
        high_freq_hz: int, High frequency limit taken from available channel

    Returns:
        Bandwidth in mhz

    Raises:
        ConfigurationError: if bandwidth is not supported by the device
    """
    bandwidth_mhz = (high_freq_hz - low_freq_hz) / 1e6
    i = bisect_right(SAS_BANDWIDTHS, bandwidth_mhz)
    if not i:
        raise ConfigurationError('Unknown/unsupported bandwidth specification (%f)' % bandwidth_mhz)
    return SAS_BANDWIDTHS[i - 1]


def calc_bandwidth_rbs(bandwidth_mhz: float) -> str:
    """
    Convert bandwidth in mhz to rbs

    Args:
        bandwidth_mhz: float, Bandwidth in mhz

    Returns:
        Bandwidth in rbs
    """
    bandwidth_rbs = int(5 * bandwidth_mhz)
    return str(bandwidth_rbs)


def calc_earfcn(low_freq_hz: int, high_freq_hz: int) -> int:
    """
    Calculate EARFCN in mhz for CBRS

    Args:
        low_freq_hz: int, Low frequency limit taken from available channel
        high_freq_hz: int, High frequency limit taken from available channel

    Returns:
        EARFCN in mhz
    """
    mid_frequency_mhz = (high_freq_hz + low_freq_hz) / 2e6
    earfcn = 10 * (mid_frequency_mhz - BAND48_LOW_FREQ_MHZ) + BAND48_NOFFS_DL
    return int(earfcn)
