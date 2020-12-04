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

import os
import codecs
import hashlib
from Crypto.Cipher import ARC4
from Crypto.Cipher import AES
from Crypto.Hash import HMAC
from lte.protos.mconfig.mconfigs_pb2 import PipelineD


def encrypt_str(s: str, key: str, encryption_algorithm):
    ret = ""
    if encryption_algorithm == PipelineD.HEConfig.RC4:
        cipher = ARC4.new(key)
        ret = cipher.encrypt(s).hex()
    elif encryption_algorithm == PipelineD.HEConfig.AES256_CBC_HMAC_MD5:
        iv = os.urandom(16)
        cipher = AES.new(key, AES.MODE_CBC, iv)
        ret = cipher.encrypt(s).hex()
    elif encryption_algorithm == PipelineD.HEConfig.AES256_ECB_HMAC_MD5:
        iv = os.urandom(16)
        cipher = AES.new(key, AES.MODE_ECB, iv)
        ret = cipher.encrypt(s).hex()
    elif encryption_algorithm == PipelineD.HEConfig.GZIPPED_AES256_ECB_SHA1:
    return ret


def get_hash(s: str, hash_function):
    hash_hex = ""
    if hash_function == PipelineD.HEConfig.MD5:
        m = hashlib.md5()
        m.update(s.encode('utf-8'))
        hash_hex = m.hexdigest()
    elif hash_function == PipelineD.HEConfig.HEX:
        hash_hex = s.encode('utf-8').hex()
    elif hash_function == PipelineD.HEConfig.SHA256:
        m = hashlib.sha256()
        m.update(s.encode('utf-8'))
        hash_hex = m.hexdigest()
    return hash_hex


def encode_str(s: str, encoding_type):
    encoded = ""
    if encoding_type == PipelineD.HEConfig.BASE64:
        encoded = codecs.encode(codecs.decode(s, 'hex'), 'base64').decode()
    elif encoding_type == PipelineD.HEConfig.HEX2BIN:
        bits = len(s) * 4
        encoded = bin(int(s, 16))[2:].zfill(bits)
    return encoded.strip()
