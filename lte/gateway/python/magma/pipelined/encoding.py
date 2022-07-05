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
from typing import Optional

from Crypto.Cipher import AES, ARC4
from Crypto.Hash import HMAC
from Crypto.Random import get_random_bytes
from lte.protos.mconfig.mconfigs_pb2 import PipelineD


def pad(m):
    return m + ' ' * (16 - len(m) % 16)


def encrypt_str(s: str, key: bytes, encryption_algorithm, mac: Optional[bytes] = None):
    if encryption_algorithm == PipelineD.HEConfig.RC4:
        cipher = ARC4.new(key)
        return cipher.encrypt(s.encode('utf-8')).hex()

    if mac is not None:
        key_val = key
        key_mac = mac
        hmac = HMAC.new(key_mac)

        if encryption_algorithm == PipelineD.HEConfig.AES256_CBC_HMAC_MD5:
            iv = get_random_bytes(16)
            aes_cipher = AES.new(key_val, AES.MODE_CBC, iv)
            enc = aes_cipher.encrypt(pad(s).encode('utf-8'))
            hmac.update(iv + enc)
            return hmac.hexdigest() + iv.hex() + enc.hex()
        elif encryption_algorithm == PipelineD.HEConfig.AES256_ECB_HMAC_MD5:
            aes_cipher = AES.new(key_val, AES.MODE_ECB)
            enc = aes_cipher.encrypt(pad(s).encode('utf-8'))
            hmac.update(enc)
            return hmac.hexdigest() + enc.hex()
        elif encryption_algorithm == PipelineD.HEConfig.GZIPPED_AES256_ECB_SHA1:
            aes_cipher = AES.new(key_val, AES.MODE_ECB)
            enc = aes_cipher.encrypt(pad(s).encode('utf-8'))
            hmac.update(enc)
            return gzip.compress(hmac.digest() + enc)

    raise ValueError("Unsupported encryption algorithm")


def decrypt_str(data, key: bytes, encryption_algorithm, mac) -> str:
    if encryption_algorithm == PipelineD.HEConfig.RC4:
        cipher = ARC4.new(key)
        return cipher.decrypt(data).hex()

    hmac = HMAC.new(mac)

    if encryption_algorithm == PipelineD.HEConfig.AES256_CBC_HMAC_MD5:
        verify = data[0:32]
        hmac.update(codecs.decode(data[32:], 'hex_codec'))

        if hmac.hexdigest() != verify:
            return ""

        iv = codecs.decode(data[32:64], 'hex_codec')
        aes_cipher = AES.new(key, AES.MODE_CBC, iv)
        decrypted = aes_cipher.decrypt(codecs.decode(data[64:], 'hex_codec'))
        return decrypted.decode("utf-8").strip()

    elif encryption_algorithm == PipelineD.HEConfig.AES256_ECB_HMAC_MD5:
        verify = data[0:32]
        hmac.update(codecs.decode(data[32:], 'hex_codec'))

        if hmac.hexdigest() != verify:
            return ""

        aes_cipher = AES.new(key, AES.MODE_ECB)
        decrypted = aes_cipher.decrypt(codecs.decode(data[32:], 'hex_codec'))
        return decrypted.decode("utf-8").strip()

    elif encryption_algorithm == PipelineD.HEConfig.GZIPPED_AES256_ECB_SHA1:
        # Convert to hex str
        data = gzip.decompress(data).hex()
        verify = data[0:32]
        hmac.update(codecs.decode(data[32:], 'hex_codec'))

        if hmac.hexdigest() != verify:
            return ""

        aes_cipher = AES.new(key, AES.MODE_ECB)
        decrypted = aes_cipher.decrypt(codecs.decode(data[32:], 'hex_codec'))
        return decrypted.decode("utf-8").strip()
    raise ValueError("Unsupported encryption algorithm")


def get_hash(s, hash_function) -> bytes:
    hash_bytes: bytes
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


def encode_str(s: str, encoding_type) -> str:
    if encoding_type == PipelineD.HEConfig.BASE64:
        s = codecs.encode(codecs.decode(s, 'hex'), 'base64').decode()  # type: ignore
    elif encoding_type == PipelineD.HEConfig.HEX2BIN:
        bits = len(s) * 4
        s = bin(int(s, 16))[2:].zfill(bits)
    else:
        logging.error("Unsupported encoding type")
    return s.strip()
