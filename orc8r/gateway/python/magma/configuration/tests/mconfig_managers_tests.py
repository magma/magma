"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import unittest
from unittest.mock import Mock, patch

from orc8r.protos.mconfig.mconfigs_pb2 import MagmaD

from magma.configuration.exceptions import LoadConfigError
from magma.configuration.mconfig_managers import MconfigManagerImpl, \
    get_mconfig_manager


class GetMconfigManagerTest(unittest.TestCase):

    @patch('magma.configuration.mconfig_managers.StreamedMconfigManager')
    def test_happy_path_feature_on(self, new_manager_mock):
        """
        Test happy path with feature flag turned on
        """
        manager_instance = Mock()
        manager_instance.load_service_mconfig.return_value = MagmaD(
            feature_flags={'kafka_config_streamer': True},
        )
        new_manager_mock.return_value = manager_instance

        actual = get_mconfig_manager()
        # We should have instantiated a StreamedMconfigManager, which in the
        # context of this test method will return the mock instance.
        self.assertIs(actual, manager_instance)
        manager_instance.load_service_mconfig.assert_called_once_with('magmad')

    @patch('magma.configuration.mconfig_managers.StreamedMconfigManager')
    def test_happy_path_feature_off(self, new_manager_mock):
        """
        Test happy path with feature flag turned off explicitly
        """
        manager_instance = Mock()
        manager_instance.load_service_mconfig.return_value = MagmaD(
            feature_flags={'kafka_config_streamer': False},
        )
        new_manager_mock.return_value = manager_instance

        actual = get_mconfig_manager()
        self.assertIsInstance(actual, MconfigManagerImpl)
        manager_instance.load_service_mconfig.assert_called_once_with('magmad')

    @patch('magma.configuration.mconfig_managers.StreamedMconfigManager')
    def test_happy_path_feature_unspecified(self, new_manager_mock):
        manager_instance = Mock()
        manager_instance.load_service_mconfig.return_value = MagmaD()
        new_manager_mock.return_value = manager_instance

        actual = get_mconfig_manager()
        self.assertIsInstance(actual, MconfigManagerImpl)
        manager_instance.load_service_mconfig.assert_called_once_with('magmad')

    @patch('magma.configuration.mconfig_managers.MconfigManagerImpl')
    @patch('magma.configuration.mconfig_managers.StreamedMconfigManager')
    def test_new_mconfig_load_error(self, new_manager_mock, old_manager_mock):
        """
        Test feature flag on, but new mconfig manager errors out on load
        """
        new_manager_instance = Mock()
        new_manager_instance.load_service_mconfig = Mock(
            side_effect=LoadConfigError('mock'),
        )
        new_manager_mock.return_value = new_manager_instance

        old_manager_instance = Mock()
        old_manager_instance.load_service_mconfig.return_value = MagmaD(
            feature_flags={'kafka_config_streamer': True},
        )
        old_manager_mock.return_value = old_manager_instance

        actual = get_mconfig_manager()
        self.assertIs(actual, new_manager_instance)
        new_manager_instance.load_service_mconfig\
            .assert_called_once_with('magmad')
        old_manager_instance.load_service_mconfig\
            .assert_called_once_with('magmad')

    @patch('magma.configuration.mconfig_managers.MconfigManagerImpl')
    @patch('magma.configuration.mconfig_managers.StreamedMconfigManager')
    def test_both_mconfigs_load_error(self, new_manager_mock, old_manager_mock):
        """
        Test both mconfig managers erroring out on load
        """
        new_manager_instance = Mock()
        new_manager_instance.load_service_mconfig = Mock(
            side_effect=LoadConfigError('mock'),
        )
        new_manager_mock.return_value = new_manager_instance

        old_manager_instance = Mock()
        old_manager_instance.load_service_mconfig = Mock(
            side_effect=LoadConfigError('mock2'),
        )
        old_manager_mock.return_value = old_manager_instance

        actual = get_mconfig_manager()
        self.assertIs(actual, old_manager_instance)
        new_manager_instance.load_service_mconfig \
            .assert_called_once_with('magmad')
        old_manager_instance.load_service_mconfig \
            .assert_called_once_with('magmad')
