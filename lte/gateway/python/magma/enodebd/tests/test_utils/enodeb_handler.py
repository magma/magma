"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from unittest import TestCase, mock
import magma.enodebd.tests.test_utils.mock_functions as enb_mock


class EnodebHandlerTestCase(TestCase):
    """
    Sets up test class with a set of patches needed for eNodeB handlers
    """
    def setUp(self):
        self.patches = {
            enb_mock.GET_IP_FROM_IF_PATH:
                mock.Mock(side_effect=enb_mock.mock_get_ip_from_if),
            enb_mock.LOAD_SERVICE_MCONFIG_PATH:
                mock.Mock(
                    side_effect=enb_mock.mock_load_service_mconfig_as_json),
        }
        self.applied_patches = [mock.patch(patch, data) for patch, data in
                                self.patches.items()]
        for patch in self.applied_patches:
            patch.start()
        self.addCleanup(mock.patch.stopall)
