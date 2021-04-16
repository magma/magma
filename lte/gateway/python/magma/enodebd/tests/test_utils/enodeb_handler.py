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
                    side_effect=enb_mock.mock_load_service_mconfig_as_json,
                ),
        }
        self.applied_patches = [
            mock.patch(patch, data) for patch, data in
            self.patches.items()
        ]
        for patch in self.applied_patches:
            patch.start()
        self.addCleanup(mock.patch.stopall)
