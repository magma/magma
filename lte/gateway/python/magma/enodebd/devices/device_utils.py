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

import re

from magma.enodebd.exceptions import UnrecognizedEnodebError
from magma.enodebd.logger import EnodebdLogger as logger


class EnodebDeviceName():
    """
    This exists only to break a circular dependency. Otherwise there's no
    point of having these names for the devices
    """
    BAICELLS = 'Baicells'
    BAICELLS_OLD = 'Baicells Old'
    BAICELLS_QAFA = 'Baicells QAFA'
    BAICELLS_QAFB = 'Baicells QAFB'
    BAICELLS_RTS = 'Baicells RTS'
    CAVIUM = 'Cavium'


def get_device_name(
    device_oui: str,
    sw_version: str,
) -> str:
    """
    Use the manufacturer organization unique identifier read during INFORM
    to select the TR data model used for configuration and status reports

    Qualcomm-based Baicells eNodeBs use a TR098-based model different
    from the Intel units. The software version on the Qualcomm models
    also further limits the model usable by that device.

    Args:
        device_oui: string, OUI representing device vendor
        sw_version: string, firmware version of eNodeB device

    Returns:
        DataModel
    """
    if device_oui in {'34ED0B', '48BF74'}:
        if sw_version.startswith('BaiBS_QAFB'):
            return EnodebDeviceName.BAICELLS_QAFB
        elif sw_version.startswith('BaiBS_QAFA'):
            return EnodebDeviceName.BAICELLS_QAFA
        elif sw_version.startswith('BaiStation_'):
            # Note: to disable flag inversion completely (for all builds),
            # set to BaiStation_V000R000C00B000SPC000
            # Note: to force flag inversion always (for all builds),
            # set to BaiStation_V999R999C99B999SPC999
            invert_before_version = \
                _parse_sw_version('BaiStation_V100R001C00B110SPC003')
            if _parse_sw_version(sw_version) < invert_before_version:
                return EnodebDeviceName.BAICELLS_OLD
            return EnodebDeviceName.BAICELLS
        elif sw_version.startswith('BaiBS_RTS_'):
            return EnodebDeviceName.BAICELLS_RTS
        elif sw_version.startswith('BaiBS_RTSH_'):
            return EnodebDeviceName.BAICELLS_RTS
        else:
            raise UnrecognizedEnodebError(
                "Device %s unsupported: Software (%s)"
                % (device_oui, sw_version),
            )
    elif device_oui in {'000FB7', '744D28'}:
        return EnodebDeviceName.CAVIUM
    else:
        raise UnrecognizedEnodebError("Device %s unsupported" % device_oui)


def _parse_sw_version(version_str):
    """
    Parse SW version string.
    Expects format: BaiStation_V100R001C00B110SPC003
    For the above version string, returns: [100, 1, 0, 110, 3]
    Note: trailing characters (for dev builds) are ignored. Null is returned
    for version strings that don't match the above format.
    """
    logger.debug('Got firmware version: %s', version_str)

    version = re.findall(
        r'BaiStation_V(\d{3})R(\d{3})C(\d{2})B(\d{3})SPC(\d{3})', version_str,
    )
    if not version:
        return None
    elif len(version) > 1:
        logger.warning(
            'SW version (%s) not formatted as expected',
            version_str,
        )
    version_int = []
    for num in version[0]:
        try:
            version_int.append(int(num))
        except ValueError:
            logger.warning(
                'SW version (%s) not formatted as expected',
                version_str,
            )
            return None

    logger.debug('Parsed firmware version: %s', version_int)

    return version_int
