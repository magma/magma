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
    Event,
    Hint,
    SentryStatus,
    SharedSentryConfig,
    _filter_excluded_messages,
    _get_before_send_hook,
    _get_shared_sentry_config,
    _ignore_if_not_marked,
    send_uncaught_errors_to_monitoring,
)
from orc8r.protos.mconfig import mconfigs_pb2


class SentryTests(unittest.TestCase):
    """
    Tests for Sentry monitoring
    """

    def test_ignore_if_not_marked(self):
        """Test ignored events that are not sent to Sentry"""
        self.assertIsNone(_ignore_if_not_marked({"key": "value", "extra": {"something_else": 3}}))

    def test_do_not_ignore_and_remove_tag_if_marked(self):
        """Test marked events that are sent to Sentry"""
        returned_event = _ignore_if_not_marked({"key": "value", "extra": {"send_to_error_monitoring": True}})
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
    def test_get_shared_sentry_config_from_control_proxy(self, get_service_config_value_mock):
        """Test if control_proxy.yml overrides mconfig (except exclusion).
        """
        get_service_config_value_mock.side_effect = ['https://test.me', 0.5, ["Excluded"]]
        sentry_mconfig = mconfigs_pb2.SharedSentryConfig()
        sentry_mconfig.exclusion_patterns.append("mconfig")
        sentry_mconfig.exclusion_patterns.append("also mconfig")
        self.assertEqual(
            SharedSentryConfig('https://test.me', 0.5, ["mconfig", "also mconfig"]),
            _get_shared_sentry_config(sentry_mconfig),
        )

    @mock.patch('magma.common.sentry.get_service_config_value')
    def test_get_shared_sentry_config_from_mconfig(self, get_service_config_value_mock):
        """Test if mconfig is used if control_proxy.yml is empty.
        """
        get_service_config_value_mock.side_effect = ['', 0.5, ["Excluded"]]
        sentry_mconfig = mconfigs_pb2.SharedSentryConfig()
        sentry_mconfig.dsn_python = 'https://test.me'
        sentry_mconfig.sample_rate = 1
        sentry_mconfig.exclusion_patterns.append("another error")
        self.assertEqual(
            SharedSentryConfig('https://test.me', 1, ["another error"]),
            _get_shared_sentry_config(sentry_mconfig),
        )

    def test_exclusion_pattern_for_log_messages(self):
        """Test events written by error logs can be filtered"""
        event = _event_with_log_message("Error in system1")
        result = _filter_excluded_messages(event, {}, ["system1"])
        self.assertEqual(result, None)

        result = _filter_excluded_messages(event, {}, ["system2"])
        self.assertEqual(result, event)

    def test_exclusion_pattern_for_explicit_messages(self):
        """Test events written explicitly can be filtered"""
        event = _event_with_explicit_message("Error in system1")
        result = _filter_excluded_messages(event, {}, ["Error.*1"])
        self.assertEqual(result, None)

        result = _filter_excluded_messages(event, {}, ["IN.*1"])
        self.assertEqual(result, event)

    def test_exclusion_pattern_for_exception_messages(self):
        """Test events written explicitly can be filtered"""
        event = _event_with_log_message("Caught exception")
        hint = _hint_with_exception_message("Error in system1")
        result = _filter_excluded_messages(event, hint, ["something", "system1"])
        self.assertEqual(result, None)

        result = _filter_excluded_messages(event, hint, ["something", "system2"])
        self.assertEqual(result, event)

    def test_get_before_send_hook_returns_none(self):
        """Test hook is not set if unnecessary"""
        self.assertIsNone(_get_before_send_hook(SentryStatus.DISABLED, []))
        self.assertIsNone(_get_before_send_hook(SentryStatus.SEND_ALL_ERRORS, []))

    def test_get_before_send_hook_returns_callback(self):
        """Test hook is set correctly"""
        hook = _get_before_send_hook(SentryStatus.SEND_SELECTED_ERRORS, [])
        self.assertTrue(callable(hook))
        self.assertEqual(hook.__name__, 'filter_excluded_and_unmarked_messages')

        hook = _get_before_send_hook(SentryStatus.SEND_ALL_ERRORS, ["a"])
        self.assertTrue(callable(hook))
        self.assertEqual(hook.__name__, 'filter_excluded_messages')


def _event_with_log_message(message: str) -> Event:
    return {"logentry": {"message": message}}


def _event_with_explicit_message(message: str) -> Event:
    return {"message": message}


def _hint_with_exception_message(message: str) -> Hint:
    return {"exc_info": (None, ValueError(message), None)}
