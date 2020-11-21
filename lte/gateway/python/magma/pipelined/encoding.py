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

import hashlib
import base64
from Crypto.Cipher import ARC4


def encrypt_str(s: str, key: str, encryption_algorithm):
    cipher = ARC4.new(key)
    return cipher.encrypt(s).hex()


def get_hash(s: str, hash_function):
    m = hashlib.md5()
    m.update(s.encode('utf-8'))
    return m.hexdigest()


def encode_str(s: str, encoding_type):
    encoded = base64.b64encode(s.encode('utf-8'))
    return encoded
