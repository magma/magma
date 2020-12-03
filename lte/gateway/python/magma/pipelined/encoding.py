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

import codecs
import hashlib
from Crypto.Cipher import ARC4

from lte.protos.mconfig.mconfigs_pb2 import PipelineD


def encrypt_str(s: str, key: str, encryption_algorithm):
    ret = ""
    if encryption_algorithm == PipelineD.HEConfig.RC4:
        cipher = ARC4.new(key)
        ret = cipher.encrypt(s).hex()
    return ret


def get_hash(s: str, hash_function):
    hash_hex = ""
    if hash_function == PipelineD.HEConfig.MD5:
        m = hashlib.md5()
        m.update(s.encode('utf-8'))
        hash_hex = m.hexdigest()
    return hash_hex


def encode_str(s: str, encoding_type):
    encoded = ""
    if encoding_type == PipelineD.HEConfig.BASE64:
        encoded = codecs.encode(codecs.decode(s, 'hex'), 'base64').decode()
    return encoded.strip()
