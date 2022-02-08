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

import unittest
from unittest import mock

from magma.common.sentry import (
    Event,
    Hint,
    SharedSentryConfig,
    _filter_excluded_messages,
    _get_before_send_hook,
    _get_shared_sentry_config,
    _ignore_if_marked,
)
from orc8r.protos.mconfig import mconfigs_pb2


class SentryTests(unittest.TestCase):
    """
    Tests for Sentry monitoring
    """

    def test_ignore_if_marked_as_excluded(self):
        """Test marked events that are not sent to Sentry"""
        self.assertIsNone(_ignore_if_marked({"key": "value", "extra": {"exclude_from_error_monitoring": True}}))

    def test_do_not_ignore_if_not_marked_as_excluded(self):
        """Test normal events that are sent to Sentry"""
        event = {"key": "value", "extra": {"something_else": 3}}
        returned_event = _ignore_if_marked(event)
        self.assertEqual(event, returned_event)

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

    def test_get_before_send_hook_returns_exclusion_and_filter_hook(self):
        """Test that the returned hook excludes marked events and applies the filter"""
        hook = _get_before_send_hook(["message to exclude"])
        self.assertEqual(hook.__name__, 'filter_excluded_and_marked_messages')

        excluded_event = {"key": "value", "extra": {"exclude_from_error_monitoring": True}}
        self.assertIsNone(hook(excluded_event, {}))

        filtered_event = _event_with_log_message("some message to exclude")
        self.assertIsNone(hook(filtered_event, {}))

        unfiltered_event = _event_with_log_message("another message")
        self.assertEqual(unfiltered_event, hook(unfiltered_event, {}))

    def test_get_before_send_hook_handles_missing_filters_correctly(self):
        """Test that the returned hook excludes only marked events if there's no filter"""
        hook = _get_before_send_hook([])
        self.assertEqual(hook.__name__, 'filter_marked_messages')

        excluded_event = {"key": "value", "extra": {"exclude_from_error_monitoring": True}}
        self.assertIsNone(hook(excluded_event, {}))

        unfiltered_event = _event_with_log_message("another message")
        self.assertEqual(unfiltered_event, hook(unfiltered_event, {}))


def _event_with_log_message(message: str) -> Event:
    return {"logentry": {"message": message}}


def _event_with_explicit_message(message: str) -> Event:
    return {"message": message}


def _hint_with_exception_message(message: str) -> Hint:
    return {"exc_info": (None, ValueError(message), None)}
