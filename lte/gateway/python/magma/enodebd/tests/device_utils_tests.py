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

# pylint: disable=protected-access

from unittest import TestCase

from magma.enodebd.devices.device_utils import (
    EnodebDeviceName,
    _parse_sw_version,
    get_device_name,
)
from magma.enodebd.exceptions import UnrecognizedEnodebError
from parameterized import parameterized


class EnodebConfigUtilsTest(TestCase):
    @parameterized.expand([
        ('34ED0B', 'BaiStation_V100R001C00B110SPC003', '', '', EnodebDeviceName.BAICELLS),
        ('34ED0B', 'BaiStation_V100R001C00B110SPC002', '', '', EnodebDeviceName.BAICELLS_OLD),
        ('48BF74', 'BaiBS_QAFB_some_version', '', '', EnodebDeviceName.BAICELLS_QAFB),
        # baicells 436Q
        ('48BF74', 'BaiBS_QRTB_some_version', 'E01', 'FAP/mBS31001/SC', EnodebDeviceName.BAICELLS_QRTB),
        # baicells 430
        ('48BF74', 'BaiBS_QRTB_some_version', 'A01', 'FAP/pBS3101S/SC', EnodebDeviceName.BAICELLS_QRTB),
        ('000FB7', 'Some version of Cavium', '', '', EnodebDeviceName.CAVIUM),
        ('000E8F', 'Some version of Sercomm', '', '', EnodebDeviceName.FREEDOMFI_ONE),
    ])
    def test_get_device_name(self, oui, sw_version, hw_version, product_class, expected) -> None:
        oui = oui
        sw_version = sw_version
        hw_version = hw_version
        product_class = product_class
        device_name = get_device_name(oui, sw_version, hw_version, product_class)
        self.assertEqual(device_name, expected, 'Incorrect device name')

    @parameterized.expand([
        ('foo', 'boopboopboop', '', ''),
        ('34ED0B', 'blingblangblong', '', ''),
        # 430's hw_version mixed with 436Q's product_class
        ('48BF74', 'BaiBS_QRTB_some_version', 'A01', 'FAP/mBS31001/SC'),
        # 436Q's hw_version mixed with 430's product_class
        ('48BF74', 'BaiBS_QRTB_some_version', 'E01', 'FAP/pBS3101S/SC'),
    ])
    def test_get_device_name_incorrect_data(self, oui, sw_version, hw_version, product_class):
        oui = oui
        sw_version = sw_version
        hw_version = hw_version
        product_class = product_class
        with self.assertRaises(UnrecognizedEnodebError):
            get_device_name(oui, sw_version, hw_version, product_class)

    def test_parse_version(self):
        """ Test that version string is parsed correctly """
        self.assertEqual(
            _parse_sw_version('BaiStation_V100R001C00B110SPC003'),
            [100, 1, 0, 110, 3],
        )
        self.assertEqual(
            _parse_sw_version('BaiStation_V100R001C00B060SPC012'),
            [100, 1, 0, 60, 12],
        )
        self.assertEqual(
            _parse_sw_version('BaiStation_V100R001C00B060SPC012_FB_3'),
            [100, 1, 0, 60, 12],
        )
        # Incorrect number of digits
        self.assertEqual(
            _parse_sw_version('BaiStation_V10R001C00B060SPC012'),
            None,
        )
        self.assertEqual(
            _parse_sw_version('XYZ123'),
            None,
        )
        self.assertEqual(
            _parse_sw_version(''),
            None,
        )
