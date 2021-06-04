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
from magma.enodebd.device_config.configuration_init import (
    _get_enb_config,
    _set_bandwidth,
    _set_earfcn_freq_band_mode,
    _set_management_server,
    _set_misc_static_params,
    _set_pci,
    _set_perf_mgmt,
    _set_plmnids_tac,
    _set_s1_connection,
    _set_tdd_subframe_config,
)
from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration
from magma.enodebd.devices.baicells import BaicellsTrDataModel
from magma.enodebd.exceptions import ConfigurationError
from magma.enodebd.tests.test_utils.config_builder import EnodebConfigBuilder


class EnodebConfigurationFactoryTest(TestCase):

    def setUp(self):
        self.data_model = BaicellsTrDataModel()
        self.cfg = EnodebConfiguration(BaicellsTrDataModel())
        self.device_cfg = EnodebConfiguration(BaicellsTrDataModel())

    def tearDown(self):
        self.data_model = None
        self.cfg = None
        self.device_cfg = None

    def test_set_pci(self):
        pci = 3
        _set_pci(self.cfg, pci)
        self.assertEqual(
            self.cfg.get_parameter(ParameterName.PCI), pci,
            'PCI value should be same as what was set',
        )
        with self.assertRaises(ConfigurationError):
            _set_pci(self.cfg, 505)

    def test_set_bandwidth(self):
        mhz = 15
        _set_bandwidth(self.cfg, self.data_model, mhz)
        self.assertEqual(
            self.cfg.get_parameter(ParameterName.DL_BANDWIDTH),
            mhz,
            'Should have set %s' % ParameterName.DL_BANDWIDTH,
        )
        self.assertEqual(
            self.cfg.get_parameter(ParameterName.UL_BANDWIDTH),
            mhz,
            'Should have set %s' % ParameterName.UL_BANDWIDTH,
        )

    def test_set_tdd_subframe_config(self):
        # Not TDD mode, should not try to set anything
        self.device_cfg.set_parameter(
            ParameterName.DUPLEX_MODE_CAPABILITY, 'Not TDDMode',
        )
        subframe = 0
        special_subframe = 0
        _set_tdd_subframe_config(
            self.device_cfg, self.cfg, subframe,
            special_subframe,
        )
        self.assertTrue(
            ParameterName.SUBFRAME_ASSIGNMENT not in
            self.cfg.get_parameter_names(),
        )

        # Invalid subframe assignment
        self.device_cfg.set_parameter(
            ParameterName.DUPLEX_MODE_CAPABILITY, 'TDDMode',
        )
        _set_tdd_subframe_config(
            self.device_cfg, self.cfg, subframe,
            special_subframe,
        )
        self.assertIn(
            ParameterName.SUBFRAME_ASSIGNMENT,
            self.cfg.get_parameter_names(),
            'Expected a subframe assignment',
        )

    def test_set_management_server(self):
        _set_management_server(self.cfg)
        self.assertEqual(
            self.cfg.get_parameter(ParameterName.PERIODIC_INFORM_ENABLE),
            True, 'Expected periodic inform enable to be true',
        )
        self.assertTrue(
            isinstance(
                self.cfg.get_parameter(ParameterName.PERIODIC_INFORM_INTERVAL),
                int,
            ),
            'Expected periodic inform interval to ani integer',
        )

    def test_set_s1_connection(self):
        invalid_mme_ip = 1234
        invalid_mme_port = '8080'
        mme_ip = '192.168.0.1'
        mme_port = 8080

        # MME IP should be a string
        with self.assertRaises(ConfigurationError):
            _set_s1_connection(self.cfg, invalid_mme_ip, mme_port)

        # MME Port should be an integer
        with self.assertRaises(ConfigurationError):
            _set_s1_connection(self.cfg, mme_ip, invalid_mme_port)

        # Check the ip and port are sort properly
        _set_s1_connection(self.cfg, mme_ip, mme_port)
        self.assertEqual(
            self.cfg.get_parameter(ParameterName.MME_IP), mme_ip,
            'Expected mme ip to be set',
        )
        self.assertEqual(
            self.cfg.get_parameter(ParameterName.MME_PORT), mme_port,
            'Expected mme port to be set',
        )

    def test_set_perf_mgmt(self):
        mgmt_ip = '192.168.0.1'
        mgmt_upload_interval = 300
        mgmt_port = 8080
        _set_perf_mgmt(self.cfg, mgmt_ip, mgmt_port)
        self.assertTrue(
            self.cfg.get_parameter(ParameterName.PERF_MGMT_ENABLE),
            'Expected perf mgmt to be enabled',
        )
        self.assertEqual(
            self.cfg.get_parameter(ParameterName.PERF_MGMT_UPLOAD_INTERVAL),
            mgmt_upload_interval, 'Expected upload interval to be set',
        )
        expected_url = 'http://192.168.0.1:8080/'
        self.assertEqual(
            self.cfg.get_parameter(ParameterName.PERF_MGMT_UPLOAD_URL),
            expected_url, 'Incorrect Url',
        )

    def test_set_misc_static_params(self):
        # IPSec enable as integer
        self.device_cfg.set_parameter(ParameterName.IP_SEC_ENABLE, 0)
        self.data_model.set_parameter_presence(ParameterName.GPS_ENABLE, True)
        _set_misc_static_params(self.device_cfg, self.cfg, self.data_model)
        self.assertTrue(
            isinstance(
                self.cfg.get_parameter(ParameterName.IP_SEC_ENABLE), int,
            ),
            'Should support an integer IP_SEC_ENABLE parameter',
        )

        # IPSec enable as boolean
        self.device_cfg.set_parameter(ParameterName.IP_SEC_ENABLE, 'False')
        _set_misc_static_params(self.device_cfg, self.cfg, self.data_model)
        self.assertTrue(
            isinstance(
                self.cfg.get_parameter(ParameterName.IP_SEC_ENABLE), bool,
            ),
            'Should support a boolean IP_SEC_ENABLE parameter',
        )
        self.assertEqual(
            self.cfg.get_parameter(ParameterName.LOCAL_GATEWAY_ENABLE), 0,
            'Should be disabled',
        )
        self.assertEqual(
            self.cfg.get_parameter(ParameterName.CELL_RESERVED), False,
            'Should be disabled',
        )
        self.assertEqual(
            self.cfg.get_parameter(ParameterName.MME_POOL_ENABLE), False,
            'Should be disabled',
        )

    def test_set_plmnids_tac(self):
        # We only handle a single PLMNID for now
        plmnids = '1, 2, 3, 4'
        tac = 1
        with self.assertRaises(ConfigurationError):
            _set_plmnids_tac(self.cfg, plmnids, tac)

        # Max PLMNID length is 6 characters
        plmnids = '1234567'
        with self.assertRaises(ConfigurationError):
            _set_plmnids_tac(self.cfg, plmnids, tac)

        # Check that only one PLMN element is enabled
        plmnids = '1'
        _set_plmnids_tac(self.cfg, plmnids, tac)
        self.assertTrue(
            self.cfg.get_parameter_for_object(
                ParameterName.PLMN_N_ENABLE % 1, ParameterName.PLMN_N % 1,
            ),
            'First PLMN should be enabled',
        )
        self.assertFalse(
            self.cfg.has_object(ParameterName.PLMN_N % 2),
            'Second PLMN should be disabled',
        )

    def test_set_earafcn_freq_band_mode(self):
        # Invalid earfcndl
        with self.assertRaises(ConfigurationError):
            invalid_earfcndl = -1
            _set_earfcn_freq_band_mode(
                self.device_cfg, self.cfg,
                self.data_model, invalid_earfcndl,
            )

        # Duplex_mode is TDD but capability is FDD
        with self.assertRaises(ConfigurationError):
            self.device_cfg.set_parameter(
                ParameterName.DUPLEX_MODE_CAPABILITY, 'FDDMode',
            )
            earfcndl = 38650  # Corresponds to TDD
            _set_earfcn_freq_band_mode(
                self.device_cfg, self.cfg,
                self.data_model, earfcndl,
            )

        # Duplex_mode is FDD but capability is TDD
        with self.assertRaises(ConfigurationError):
            self.device_cfg.set_parameter(
                ParameterName.DUPLEX_MODE_CAPABILITY, 'TDDMode',
            )
            earfcndl = 0  # Corresponds to FDD
            _set_earfcn_freq_band_mode(
                self.device_cfg, self.cfg,
                self.data_model, earfcndl,
            )

    def test_get_enb_config(self):
        mconfig = EnodebConfigBuilder.get_mconfig()
        enb_config = _get_enb_config(mconfig, self.device_cfg)
        self.assertTrue(
            enb_config.earfcndl == 39150,
            "Should give earfcndl from default eNB config",
        )

        mconfig = EnodebConfigBuilder.get_multi_enb_mconfig()
        self.device_cfg.set_parameter(
            ParameterName.SERIAL_NUMBER,
            '120200002618AGP0003',
        )
        enb_config = _get_enb_config(mconfig, self.device_cfg)
        self.assertTrue(
            enb_config.earfcndl == 39151,
            "Should give earfcndl from specific eNB config",
        )
