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
from abc import ABC, abstractmethod

from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration


class EnodebConfigurationPostProcessor(ABC):
    """
    Overrides the desired configuration for the eNodeB, with subclass per
    device/sw-version that requires non-standard configuration behavior.
    """

    @abstractmethod
    def postprocess(self, desired_cfg: EnodebConfiguration) -> None:
        """
        Implementation of function which overrides the desired configuration
        for the eNodeB
        """
        pass
