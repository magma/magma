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
from Crypto.Cipher import AES
from Crypto.Hash import HMAC
from Crypto.Random import get_random_bytes
from lte.protos.mconfig.mconfigs_pb2 import PipelineD


def pad(m):
    return m + ' ' * (16 - len(m) % 16)


def encrypt_str(s: str, key: bytes, encryption_algorithm, mac: bytes = None):
    ret = ""
    if encryption_algorithm == PipelineD.HEConfig.RC4:
        cipher = ARC4.new(key)
        ret = cipher.encrypt(s).hex()
    elif encryption_algorithm == PipelineD.HEConfig.AES256_CBC_HMAC_MD5:
        # Assuming provided key is valid 256 bytes
        iv = get_random_bytes(16)
        key_val = key
        key_mac = mac

        cipher = AES.new(key_val, AES.MODE_CBC, iv)
        enc = cipher.encrypt(pad(s))

        hmac = HMAC.new(key_mac)
        hmac.update(iv + enc)

        ret = hmac.hexdigest() + iv.hex() + enc.hex()
    elif encryption_algorithm == PipelineD.HEConfig.AES256_ECB_HMAC_MD5:
        # Assuming provided key is valid 256 bytes
        key_val = key
        key_mac = mac

        cipher = AES.new(key_val, AES.MODE_ECB)
        enc = cipher.encrypt(pad(s))

        hmac = HMAC.new(key_mac)
        hmac.update(enc)

        ret = hmac.hexdigest() + enc.hex()
        print("This be ecb")
        print(ret)
    elif encryption_algorithm == PipelineD.HEConfig.GZIPPED_AES256_ECB_SHA1:
        # Assuming provided key is valid 256 bytes
        key_val = key
        key_mac = mac

        cipher = AES.new(key_val, AES.MODE_ECB)
        enc = cipher.encrypt(pad(s))

        hmac = HMAC.new(key_mac)
        hmac.update(enc)

        ret = hmac.hexdigest() + enc.hex()
    return ret


def decrypt_str(data: str, key: bytes, encryption_algorithm, mac):
    if encryption_algorithm == PipelineD.HEConfig.RC4:
        return data
    elif encryption_algorithm == PipelineD.HEConfig.AES256_CBC_HMAC_MD5:
        verify = data[0:32]
        print(verify)
        hmac = HMAC.new(mac)
        hmac.update(codecs.decode(data[32:], 'hex_codec'))

        print("decrypted ->")
        print(hmac.hexdigest())
        if hmac.hexdigest() != verify:
            return ""

        # decrypt
        iv = codecs.decode(data[32:64], 'hex_codec')
        print(data[32:64])
        print(iv)
        cipher = AES.new(key, AES.MODE_CBC, iv)
        decrypted = cipher.decrypt(codecs.decode(data[64:], 'hex_codec'))
        print(codecs.decode(data[64:], 'hex_codec'))
        print(decrypted)
        return decrypted.decode("utf-8").strip()
    elif encryption_algorithm == PipelineD.HEConfig.AES256_ECB_HMAC_MD5:
        verify = data[0:32]
        print(verify)
        hmac = HMAC.new(mac)
        hmac.update(codecs.decode(data[32:], 'hex_codec'))

        print("decrypted ->")
        print(hmac.hexdigest())
        if hmac.hexdigest() != verify:
            return ""

        cipher = AES.new(key, AES.MODE_ECB)
        decrypted = cipher.decrypt(codecs.decode(data[32:], 'hex_codec'))
        print(codecs.decode(data[32:], 'hex_codec'))
        print(decrypted)
        return decrypted.decode("utf-8").strip()


def get_hash(s: str, hash_function) -> bytes:
    hash_bin = ""
    if hash_function == PipelineD.HEConfig.MD5:
        m = hashlib.md5()
        m.update(s.encode('utf-8'))
        hash_bin = m.digest()
    elif hash_function == PipelineD.HEConfig.HEX:
        hexlify = codecs.getencoder('hex')
        hash_bin = hexlify(s.encode('utf-8'))[0]
    elif hash_function == PipelineD.HEConfig.SHA256:
        m = hashlib.sha256()
        m.update(s.encode('utf-8'))
        hash_bin = m.digest()
    return hash_bin


def encode_str(s: str, encoding_type):
    if encoding_type == PipelineD.HEConfig.BASE64:
        s = codecs.encode(codecs.decode(s, 'hex'), 'base64').decode()
    elif encoding_type == PipelineD.HEConfig.HEX2BIN:
        bits = len(s) * 4
        s = bin(int(s, 16))[2:].zfill(bits)
    return s.strip()
