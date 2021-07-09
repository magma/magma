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

from unittest import TestCase, mock
from xml.etree import ElementTree

import pkg_resources
from magma.enodebd import metrics
from magma.enodebd.data_models.data_model_parameters import ParameterName
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.state_machines.enb_acs_manager import StateMachineManager
from magma.enodebd.stats_manager import StatsManager
from magma.enodebd.tests.test_utils.config_builder import EnodebConfigBuilder
from magma.enodebd.tests.test_utils.enb_acs_builder import (
    EnodebAcsStateMachineBuilder,
)


class StatsManagerTest(TestCase):
    """
    Tests for eNodeB statistics manager
    """

    def setUp(self) -> None:
        service = EnodebConfigBuilder.get_service_config()
        self.enb_acs_manager = StateMachineManager(service)
        self.mgr = StatsManager(self.enb_acs_manager)
        self.is_clear_stats_called = False

    def tearDown(self):
        self.mgr = None

    def test_check_rf_tx(self):
        """ Check that stats are cleared when transmit is disabled on eNB """
        handler = EnodebAcsStateMachineBuilder \
            .build_acs_state_machine(EnodebDeviceName.BAICELLS)
        with mock.patch(
                'magma.enodebd.devices.baicells.BaicellsHandler.is_enodeb_connected',
                return_value=True,
        ):
            handler.device_cfg.set_parameter(ParameterName.RF_TX_STATUS, True)
            handler.device_cfg.set_parameter(
                ParameterName.SERIAL_NUMBER,
                '123454',
            )
            with mock.patch(
                'magma.enodebd.stats_manager.StatsManager'
                '._clear_stats',
            ) as func:
                self.mgr._check_rf_tx_for_handler(handler)
                func.assert_not_called()
                handler.device_cfg.set_parameter(
                    ParameterName.RF_TX_STATUS,
                    False,
                )
                self.mgr._check_rf_tx_for_handler(handler)
                func.assert_any_call()

    def test_parse_stats(self):
        """ Test that example statistics from eNodeB can be parsed, and metrics
            updated """
        # Example performance metrics structure, sent by eNodeB
        pm_file_example = pkg_resources.resource_string(
            __name__,
            'pm_file_example.xml',
        )

        root = ElementTree.fromstring(pm_file_example)
        self.mgr._parse_pm_xml('1234', root)

        # Check that metrics were correctly populated
        # See '<V i="5">123</V>' in pm_file_example
        rrc_estab_attempts = metrics.STAT_RRC_ESTAB_ATT.collect()
        self.assertEqual(rrc_estab_attempts[0].samples[0][2], 123)
        # See '<V i="7">99</V>' in pm_file_example
        rrc_estab_successes = metrics.STAT_RRC_ESTAB_SUCC.collect()
        self.assertEqual(rrc_estab_successes[0].samples[0][2], 99)
        # See '<SV>654</SV>' in pm_file_example
        rrc_reestab_att_reconf_fail = \
            metrics.STAT_RRC_REESTAB_ATT_RECONF_FAIL.collect()
        self.assertEqual(rrc_reestab_att_reconf_fail[0].samples[0][2], 654)
        # See '<SV>65537</SV>' in pm_file_example
        erab_rel_req_radio_conn_lost = \
            metrics.STAT_ERAB_REL_REQ_RADIO_CONN_LOST.collect()
        self.assertEqual(erab_rel_req_radio_conn_lost[0].samples[0][2], 65537)

        pdcp_user_plane_bytes_ul = \
            metrics.STAT_PDCP_USER_PLANE_BYTES_UL.collect()
        pdcp_user_plane_bytes_dl = \
            metrics.STAT_PDCP_USER_PLANE_BYTES_DL.collect()
        self.assertEqual(pdcp_user_plane_bytes_ul[0].samples[0][1], {'enodeb': '1234'})
        self.assertEqual(pdcp_user_plane_bytes_dl[0].samples[0][1], {'enodeb': '1234'})
        self.assertEqual(pdcp_user_plane_bytes_ul[0].samples[0][2], 1000)
        self.assertEqual(pdcp_user_plane_bytes_dl[0].samples[0][2], 500)
