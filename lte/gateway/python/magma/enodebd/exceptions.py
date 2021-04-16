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


class ConfigurationError(Exception):
    """ Indicates that the eNodeB could not be configured correctly. """
    pass


class Tr069Error(Exception):
    pass


class IncorrectDeviceHandlerError(Exception):
    """ Indicates that we're using the wrong data model for configuration. """

    def __init__(self, device_name: str):
        """
        device_name: What device we actually are dealing with
        """
        super().__init__()
        self.device_name = device_name


class UnrecognizedEnodebError(Exception):
    """
    Indicates that the Access Gateway does not recognize the eNodeB.
    The Access Gateway will not interact with the eNodeB in question.
    """
    pass
