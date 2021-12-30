"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import logging
import unittest
from unittest import mock

from magma.common.sentry import (
    _ignore_if_not_marked,
    get_sentry_dsn_and_sample_rate,
    send_uncaught_errors_to_monitoring,
)
from orc8r.protos.mconfig import mconfigs_pb2


class SentryTests(unittest.TestCase):
    """
    Tests for Sentry monitoring
    """

    def test_ignore_if_not_marked(self):
        """Test ignored events that are not sent to Sentry"""
        self.assertIsNone(_ignore_if_not_marked({"key": "value", "extra": {"something_else": 3}}, {}))

    def test_do_not_ignore_and_remove_tag_if_marked(self):
        """Test marked events that are sent to Sentry"""
        returned_event = _ignore_if_not_marked({"key": "value", "extra": {"send_to_error_monitoring": True}}, {})
        self.assertEqual({"key": "value", "extra": {}}, returned_event)

    def test_uncaught_error_wrapper(self):
        """Test enabled wrapper logs an error"""
        @send_uncaught_errors_to_monitoring(True)
        def raise_error():
            raise ValueError("Something went wrong")

        with self.assertLogs() as captured:
            try:
                raise_error()
            except ValueError:
                pass

        self.assertEqual(1, len(captured.records))
        log_record = captured.records[0]
        self.assertEqual("Uncaught error", log_record.message)
        self.assertEqual(logging.ERROR, log_record.levelno)

    @mock.patch('magma.common.sentry.get_service_config_value')
    def test_get_sentry_dsn_and_sample_rate_from_control_proxy(self, get_service_config_value_mock):
        """Test if control_proxy.yml overrides mconfig.
        """
        get_service_config_value_mock.side_effect = ['https://test.me', 0.5]
        sentry_mconfig = mconfigs_pb2.SharedSentryConfig()
        self.assertEqual(('https://test.me', 0.5), get_sentry_dsn_and_sample_rate(sentry_mconfig))

    @mock.patch('magma.common.sentry.get_service_config_value')
    def test_get_sentry_dsn_and_sample_rate_from_mconfig(self, get_service_config_value_mock):
        """Test if mconfig is used if control_proxy.yml is empty.
        """
        get_service_config_value_mock.side_effect = ['', 0.5]
        sentry_mconfig = mconfigs_pb2.SharedSentryConfig()
        sentry_mconfig.dsn_python = 'https://test.me'
        sentry_mconfig.sample_rate = 1
        self.assertEqual(('https://test.me', 1), get_sentry_dsn_and_sample_rate(sentry_mconfig))
