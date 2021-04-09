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
import gzip
import hashlib
import logging

from Crypto.Cipher import AES, ARC4
from Crypto.Hash import HMAC
from Crypto.Random import get_random_bytes
from lte.protos.mconfig.mconfigs_pb2 import PipelineD


def pad(m):
    return m + ' ' * (16 - len(m) % 16)


def encrypt_str(s: str, key: bytes, encryption_algorithm, mac: bytes = None):
    ret = ""
    if encryption_algorithm == PipelineD.HEConfig.RC4:
        cipher = ARC4.new(key)
        ret = cipher.encrypt(s.encode('utf-8')).hex()
    elif encryption_algorithm == PipelineD.HEConfig.AES256_CBC_HMAC_MD5:
        iv = get_random_bytes(16)
        key_val = key
        key_mac = mac

        cipher = AES.new(key_val, AES.MODE_CBC, iv)
        enc = cipher.encrypt(pad(s).encode('utf-8'))

        hmac = HMAC.new(key_mac)
        hmac.update(iv + enc)

        ret = hmac.hexdigest() + iv.hex() + enc.hex()
    elif encryption_algorithm == PipelineD.HEConfig.AES256_ECB_HMAC_MD5:
        key_val = key
        key_mac = mac

        cipher = AES.new(key_val, AES.MODE_ECB)
        enc = cipher.encrypt(pad(s).encode('utf-8'))

        hmac = HMAC.new(key_mac)
        hmac.update(enc)

        ret = hmac.hexdigest() + enc.hex()
    elif encryption_algorithm == PipelineD.HEConfig.GZIPPED_AES256_ECB_SHA1:
        key_val = key
        key_mac = mac

        cipher = AES.new(key_val, AES.MODE_ECB)
        enc = cipher.encrypt(pad(s).encode('utf-8'))

        hmac = HMAC.new(key_mac)
        hmac.update(enc)
        ret = gzip.compress(hmac.digest() + enc)
    else:
        logging.error("Unsupported encryption algorithm")
    return ret


def decrypt_str(data: str, key: bytes, encryption_algorithm, mac) -> str:
    ret = ""
    if encryption_algorithm == PipelineD.HEConfig.RC4:
        cipher = ARC4.new(key)
        ret = cipher.decrypt(data).hex()
    elif encryption_algorithm == PipelineD.HEConfig.AES256_CBC_HMAC_MD5:
        verify = data[0:32]
        hmac = HMAC.new(mac)
        hmac.update(codecs.decode(data[32:], 'hex_codec'))

        if hmac.hexdigest() != verify:
            return ""

        iv = codecs.decode(data[32:64], 'hex_codec')
        cipher = AES.new(key, AES.MODE_CBC, iv)
        decrypted = cipher.decrypt(codecs.decode(data[64:], 'hex_codec'))
        ret = decrypted.decode("utf-8").strip()
    elif encryption_algorithm == PipelineD.HEConfig.AES256_ECB_HMAC_MD5:
        verify = data[0:32]
        hmac = HMAC.new(mac)
        hmac.update(codecs.decode(data[32:], 'hex_codec'))

        if hmac.hexdigest() != verify:
            return ""

        cipher = AES.new(key, AES.MODE_ECB)
        decrypted = cipher.decrypt(codecs.decode(data[32:], 'hex_codec'))
        ret = decrypted.decode("utf-8").strip()
    elif encryption_algorithm == PipelineD.HEConfig.GZIPPED_AES256_ECB_SHA1:
        # Convert to hex str
        hexlify = codecs.getencoder('hex')
        data = hexlify(gzip.decompress(data))[0].decode('utf-8')

        verify = data[0:32]
        hmac = HMAC.new(mac)
        hmac.update(codecs.decode(data[32:], 'hex_codec'))

        if hmac.hexdigest() != verify:
            return ""

        cipher = AES.new(key, AES.MODE_ECB)
        decrypted = cipher.decrypt(codecs.decode(data[32:], 'hex_codec'))
        ret = decrypted.decode("utf-8").strip()
    else:
        logging.error("Unsupported encryption algorithm")
    return ret


def get_hash(s: str, hash_function) -> bytes:
    hash_bytes = ""
    if hash_function == PipelineD.HEConfig.MD5:
        m = hashlib.md5()
        m.update(s.encode('utf-8'))
        hash_bytes = m.digest()
    elif hash_function == PipelineD.HEConfig.HEX:
        hexlify = codecs.getencoder('hex')
        hash_bytes = hexlify(s.encode('utf-8'))[0]
    elif hash_function == PipelineD.HEConfig.SHA256:
        m = hashlib.sha256()
        m.update(s.encode('utf-8'))
        hash_bytes = m.digest()
    else:
        logging.error("Unsupported hash function")
    return hash_bytes


def encode_str(s: str, encoding_type) -> bytes:
    if encoding_type == PipelineD.HEConfig.BASE64:
        s = codecs.encode(codecs.decode(s, 'hex'), 'base64').decode()
    elif encoding_type == PipelineD.HEConfig.HEX2BIN:
        bits = len(s) * 4
        s = bin(int(s, 16))[2:].zfill(bits)
    else:
        logging.error("Unsupported encoding type")
    return s.strip()
