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

from magma.enodebd.data_models.data_model_parameters import ParameterName
from magma.enodebd.devices.baicells import BaicellsTrDataModel


class BaicellsTrDataModelTest(TestCase):
    """
    Tests for BaicellsTrDataModel
    """

    def test_is_parameter_present(self):
        data_model = BaicellsTrDataModel()
        with self.assertRaises(KeyError):
            data_model.is_parameter_present(ParameterName.GPS_LONG)

        result = data_model.is_parameter_present(ParameterName.DEVICE)
        self.assertTrue(result, "Should have the device parameter")

    def test_get_parameter(self):
        param_info =\
            BaicellsTrDataModel.get_parameter(ParameterName.GPS_STATUS)
        self.assertIsNotNone(
            param_info,
            'Cannot get parameter info on %s' % ParameterName.GPS_STATUS,
        )
        path = param_info.path
        expected_path = 'Device.DeviceInfo.X_BAICELLS_COM_GPS_Status'
        self.assertEqual(
            path,
            expected_path,
            'Path for parameter %s has incorrect value' %
            ParameterName.GPS_STATUS,
        )

    def test_get_num_plmns(self):
        n_plmns = BaicellsTrDataModel.get_num_plmns()
        expected_n_plmns = 6
        self.assertEqual(n_plmns, expected_n_plmns, 'Incorrect # of PLMNs')

    def test_get_parameter_names(self):
        name_list = BaicellsTrDataModel.get_parameter_names()

        # Check that some random parameter names we expect are there
        self.assertIn(
            ParameterName.MME_STATUS, name_list,
            'Should have %s in parameter name list' %
            ParameterName.MME_STATUS,
        )
        self.assertIn(
            ParameterName.DUPLEX_MODE_CAPABILITY, name_list,
            'Should have %s in parameter name list' %
            ParameterName.DUPLEX_MODE_CAPABILITY,
        )
        self.assertIn(
            ParameterName.EARFCNUL, name_list,
            'Should have %s in parameter name list' %
            ParameterName.EARFCNUL,
        )

        # Check that some other parameter names are missing
        self.assertNotIn(
            ParameterName.PLMN, name_list,
            'Should not have %s in parameter name list' %
            ParameterName.PLMN,
        )
        self.assertNotIn(
            ParameterName.PLMN_N % 1, name_list,
            'Should not have %s in parameter name list' %
            ParameterName.PLMN_N % 1,
        )

    def test_get_numbered_param_names(self):
        name_list = list(BaicellsTrDataModel.get_numbered_param_names().keys())

        # Check that unnumbered parameters are missing
        self.assertNotIn(
            ParameterName.EARFCNDL, name_list,
            'Should not have %s in parameter name list' %
            ParameterName.EARFCNDL,
        )
        self.assertNotIn(
            ParameterName.MME_PORT, name_list,
            'Should not have %s in parameter name list' %
            ParameterName.MME_PORT,
        )
        self.assertNotIn(
            ParameterName.PERIODIC_INFORM_ENABLE, name_list,
            'Should not have %s in parameter name list' %
            ParameterName.PERIODIC_INFORM_ENABLE,
        )

        # Check that some numbered parameters are present
        self.assertIn(
            ParameterName.PLMN_N % 1, name_list,
            'Should have %s in parameter name list' %
            ParameterName.PLMN_N % 1,
        )
        self.assertIn(
            ParameterName.PLMN_N % 6, name_list,
            'Should have %s in parameter name list' %
            ParameterName.PLMN_N % 6,
        )

    def test_transform_for_magma(self):
        gps_lat = str(10 * 1000000)
        gps_lat_magma = BaicellsTrDataModel.transform_for_magma(
            ParameterName.GPS_LAT, gps_lat,
        )
        expected = str(10.0)
        self.assertEqual(gps_lat_magma, expected)

    def test_transform_for_enb(self):
        dl_bandwidth = 15
        dl_bandwidth_enb = BaicellsTrDataModel.transform_for_enb(
            ParameterName.DL_BANDWIDTH, dl_bandwidth,
        )
        expected = 'n75'
        self.assertEqual(
            dl_bandwidth_enb, expected,
            'Transform for enb returning incorrect value',
        )
