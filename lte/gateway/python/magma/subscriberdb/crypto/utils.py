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

from typing import List


def xor_buf(b1: List, b2: List) -> bytes:
    """
    xor_buf - xor two lists

    Args:
        b1: List
        b2: List

    Returns:
        bytes: the result

    """
    return bytes([b1[i] ^ b2[i] for i in range(0, min(len(b1), len(b2)))])


def int_from_bytes(b) -> int:
    """
    int_from_bytes - converts bytes to int

    Args:
        b: bytes

    Returns:
        int: result
    """
    return int.from_bytes(b, 'big')


# CryptoMobile-wide Exception handler
class CMException(Exception):
    """CryptoMobile specific exception
    """
    pass            # noqa: WPS604


class CryptoError(Exception):
    """
    Represents any error triggered during a crypto operation.
    """
    pass            # noqa: WPS604
