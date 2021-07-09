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

from enum import Enum
from typing import Optional


class DuplexMode(Enum):
    FDD = 1
    TDD = 2


class LTEBandInfo:
    """ Class for holding information related to LTE band.
        Takes advantage of the following properties of LTE EARFCN assignment:
            - EARFCN spacing within a band is always 0.1MHz
            - 1:1 mapping between EARFCNDL and EARFCNUL (for FDD)
    """

    def __init__(
        self, duplex_mode, earfcndl, start_freq_dl_mhz,
        start_earfcnul=None, start_freq_ul_mhz=None,
    ):
        """
        Inputs:
        - Duplex mode - type = DuplexMode
        - earfcndl - List/generator of integer EARFCNDLs in band, from lowest to
                     highest
        - start_freq_dl_mhz - Frequency of lowest numbered EARFCNDL in band
                            (MHz)
        - start_earfcnul - Lowest numbered EARFCNUL in band
                            (or None if band is TDD)
        - start_freq_ul_mhz - Frequency of lowest numbered EARFCNUL in band
                              (MHz) (or None if band is TDD)
        """
        # Validate inputs
        assert type(earfcndl) == list or type(earfcndl) == range
        assert type(start_freq_dl_mhz) == int
        assert type(duplex_mode) == DuplexMode
        if duplex_mode == DuplexMode.FDD:
            assert type(start_earfcnul) == int
            assert type(start_freq_ul_mhz) == int
        else:
            assert start_earfcnul is None
            assert start_freq_ul_mhz is None

        # DuplexMode.TDD or DuplexMode.FDD
        self.duplex_mode = duplex_mode
        # Array of EARFCNDL values
        self.earfcndl = earfcndl
        # Array of DL frequencies in MHz, one per EARFCNDL
        self.freq_mhz_dl = [
            start_freq_dl_mhz + 0.1 * (earfcn - earfcndl[0])
            for earfcn in earfcndl
        ]

        if duplex_mode == DuplexMode.FDD:
            # Array of EARFCNUL values that map to EARFCNDL
            self.earfcnul = range(
                start_earfcnul,
                start_earfcnul + len(earfcndl),
            )
            # Array of UL frequencies in MHz, one per EARFCNUL
            self.freq_mhz_ul = [
                start_freq_ul_mhz + 0.1 * (
                    earfcn
                    - self.earfcnul[0]
                ) for earfcn in self.earfcnul
            ]
        else:
            self.earfcnul = None
            self.freq_mhz_ul = None


# See, for example, http://niviuk.free.fr/lte_band.php for LTE band info
LTE_BAND_INFO = {
    # FDD bands
    # duplex_mode, EARFCNDL, start_freq_dl, start_EARCNUL, start_freq_ul
    1: LTEBandInfo(DuplexMode.FDD, range(0, 600), 2110, 18000, 1920),
    2: LTEBandInfo(DuplexMode.FDD, range(600, 1200), 1930, 18600, 1850),
    3: LTEBandInfo(DuplexMode.FDD, range(1200, 1950), 1805, 19200, 1710),
    4: LTEBandInfo(DuplexMode.FDD, range(1950, 2400), 2110, 19950, 1710),
    5: LTEBandInfo(DuplexMode.FDD, range(2400, 2649), 869, 20400, 824),
    28: LTEBandInfo(DuplexMode.FDD, range(9210, 9660), 758, 27210, 703),
    # TDD bands
    # duplex_mode, EARFCNDL, start_freq_dl
    33: LTEBandInfo(DuplexMode.TDD, range(36000, 36199), 1900),
    34: LTEBandInfo(DuplexMode.TDD, range(36200, 36349), 2010),
    35: LTEBandInfo(DuplexMode.TDD, range(36350, 36949), 1850),
    36: LTEBandInfo(DuplexMode.TDD, range(36950, 37549), 1930),
    37: LTEBandInfo(DuplexMode.TDD, range(37550, 37750), 1910),
    38: LTEBandInfo(DuplexMode.TDD, range(37750, 38250), 2570),
    39: LTEBandInfo(DuplexMode.TDD, range(38250, 38650), 1880),
    40: LTEBandInfo(DuplexMode.TDD, range(38650, 39650), 2300),
    41: LTEBandInfo(DuplexMode.TDD, range(39650, 41590), 2496),
    42: LTEBandInfo(DuplexMode.TDD, range(41590, 43590), 3400),
    43: LTEBandInfo(DuplexMode.TDD, range(43590, 45590), 3600),
    44: LTEBandInfo(DuplexMode.TDD, range(45590, 46589), 703),
    45: LTEBandInfo(DuplexMode.TDD, range(46590, 46789), 1447),
    46: LTEBandInfo(DuplexMode.TDD, range(46790, 54539), 5150),
    47: LTEBandInfo(DuplexMode.TDD, range(54540, 55239), 5855),
    48: LTEBandInfo(DuplexMode.TDD, range(55240, 56740), 3550),
    49: LTEBandInfo(DuplexMode.TDD, range(56740, 58239), 3550),
    50: LTEBandInfo(DuplexMode.TDD, range(58240, 59089), 1432),
    51: LTEBandInfo(DuplexMode.TDD, range(59090, 59139), 1427),
    52: LTEBandInfo(DuplexMode.TDD, range(59140, 60139), 3300),
    # For the band #53 start_freq_dl is float value which require some changes
    # in the code
    # 53: LTEBandInfo(DuplexMode.TDD, range(60140, 60254), 2483.5),

}
# TODO - add remaining FDD LTE bands


def map_earfcndl_to_duplex_mode(earfcndl: int) -> Optional[DuplexMode]:
    """
    Returns None if we do not support the EARFCNDL
    """
    for band in LTE_BAND_INFO.keys():
        if earfcndl in LTE_BAND_INFO[band].earfcndl:
            return LTE_BAND_INFO[band].duplex_mode
    return None


def map_earfcndl_to_band_earfcnul_mode(earfcndl):
    """ Inputs:
            - EARFCNDL (integer)
        Outputs:
            - Band (integer)
            - Mode (DuplexMode)
            - EARFCNUL (integer - or None if TDD)
    """
    for band in LTE_BAND_INFO.keys():
        if earfcndl in LTE_BAND_INFO[band].earfcndl:
            if LTE_BAND_INFO[band].duplex_mode == DuplexMode.FDD:
                index = LTE_BAND_INFO[band].earfcndl.index(earfcndl)
                earfcnul = LTE_BAND_INFO[band].earfcnul[index]
            else:
                earfcnul = None

            return band, LTE_BAND_INFO[band].duplex_mode, earfcnul

    raise ValueError('EARFCNDL %d not found in LTE band info' % earfcndl)
