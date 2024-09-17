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

import abc

from .utils import CryptoError


class GSMA3A8Algo(metaclass=abc.ABCMeta):
    """
    Abstract class for the GSM A3/A8 algorithms. The A3/A8 algos take
    the key and random variable as input, and produce an auth tuple as output.
    """

    @abc.abstractmethod
    def generate_auth_tuple(self, key):
        """
        Args:
            key - secret key for a subscriberdb
        Returns:
            (rand, sres, cipher_key) auth tuple
        Raises:
            CryptoError on any error
        """
        raise NotImplementedError()


class UnsafePreComputedA3A8(GSMA3A8Algo):
    """
    Sample implementation of the A3/A8 algo. This algo expects the auth
    tuple to be stored directly as the key for the subscriber.

    Essentially this algo doesn't do any random number generation or crypto
    operation, but provides a dummy implementation of the A3/A8 interfaces.
    """

    def generate_auth_tuple(self, key):
        """
        Args:
            key - 28 byte long auth tuple
        Returns:
            (rand, sres, cipher_key) tuple
        Raises:
            CryptoError if the key is not 28 byte long
        """
        if len(key) != 28:
            raise CryptoError('Invalid auth vector: %s' % key)
        rand = key[:16]
        sres = key[16:20]
        cipher_key = key[20:]
        return (rand, sres, cipher_key)
