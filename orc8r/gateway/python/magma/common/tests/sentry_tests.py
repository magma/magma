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

from magma.common.sentry import (
    _ignore_if_not_marked,
    send_uncaught_errors_to_monitoring,
)


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
        self.assertEqual(returned_event, {"key": "value", "extra": {}})

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

        self.assertEqual(len(captured.records), 1)
        log_record = captured.records[0]
        self.assertEqual(log_record.message, "Uncaught error")
        self.assertEqual(log_record.levelno, logging.ERROR)
