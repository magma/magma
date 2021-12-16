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

import hashlib  # for SHA-256
import hmac  # HMAC–SHA-256
from struct import pack, unpack
from typing import Optional

from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from magma.subscriberdb.crypto.EC import ECDH_SECP256R1, KDF, X25519
from magma.subscriberdb.crypto.utils import CMException

########################################################
# CryptoMobile python toolkit
#
# ECIES: Elliptic Curve Integrated Encryption Scheme
# as defined by 3GPP to protect subscriber's fixed identity SUPI into SUCI
# TS 33.501, section 6.12 and annex C
#
# Based on an elliptic curve algorithm:
# - profile A: Curve25519 / X25519
# - profile B: NIST secp256r1
# and ANSI-X9.63-KDF, SHA-256 hash function, HMAC–SHA-256 MAC function and
# AES–128 in CTR mode encryption function
#######################################################

_backend = default_backend()


class AES_CTR_cryptography(object):     # noqa: N801
    """AES in CTR mode"""

    block_size = 16

    def __init__(self, key, nonce, cnt=0):
        """
        Init - Initialize AES in ECB mode with the given key and nonce buffer

        Args:
            key  : 16 bytes buffer
            nonce: 8 most significant bytes buffer of the counter initial value
                counter will be incremented starting at 0
            cnt  : uint64, 8 least significant bytes value of the counter
                default is 0

        Returns: None
        """
        self.aes = Cipher(
            algorithms.AES(key),
            modes.CTR(nonce + pack('>Q', cnt)),
            backend=_backend,
        ).encryptor()

    def encrypt(self, data):
        """Encrypt / decrypt data with the key and IV set at initialization"""
        return self.aes.update(data)

    decrypt = encrypt


class ECIES_UE(object):         # noqa: N801
    """ECIES_UE handles the ECIES computation required on the UE side to
    protect its fixed identity SUPI into a SUCI
    """

    def __init__(self, profile):
        if profile == 'A':
            self.EC = X25519()
        elif profile == 'B':
            self.EC = ECDH_SECP256R1()
        else:
            raise CMException('unknown ECIES profile %s' % profile)
        self.ephemeral_key = None
        self.shared_key = None

    def generate_sharedkey(self, hn_pub_key, fresh=True):
        """
        generate_sharedkey - Generates a shared keystream based on a UE ephemeral keypair (regenerated
        if fresh is True) and the HN public key

        Args:
            hn_pub_key: hn_pub_key
            fresh: True

        Returns: None
        """
        if fresh:
            # regenerate a new UE ephemeral keypair
            self.EC.generate_keypair()
        # get the UE ephemeral pubkey
        self.ephemeral_key = self.EC.get_pubkey()
        # generate the shared keystream by mixing the UE ephemeral key with HN pubkey
        self.shared_key = KDF(self.ephemeral_key, self.EC.generate_sharedkey(hn_pub_key))

    def protect(self, plaintext):
        """
        Protect - Protects the given plaintext and returns a 3-tuple of bytes:
        UE ephemeral pubkey, ciphertext, MAC

        Args:
            plaintext: plaintext

        Returns:
            ephemeral_key: public key
            ciphertext: ciphertext to decrypt
            mac: encrypted mac
        """
        aes_key, aes_nonce, aes_cnt, mac_key = (
            self.shared_key[:16],
            self.shared_key[16:24],
            unpack('>Q', self.shared_key[24:32])[0],
            self.shared_key[32:64],
        )
        # encryption
        aes = AES_CTR_cryptography(aes_key, aes_nonce, aes_cnt)
        ciphertext = aes.encrypt(plaintext)
        mac = hmac.new(mac_key, ciphertext, hashlib.sha256).digest()
        return self.ephemeral_key, ciphertext, mac[0:8]         # noqa: WPS349


class ECIES_HN(object):             # noqa: N801
    """ECIES_HN handles the ECIES computation required on the Home Network side
    to unprotect a subscriber's SUCI into a fixed identity SUPI
    """

    def __init__(self, hn_priv_key, profile='A', raw_keypair=None):
        if profile == 'A':
            self.EC = X25519(loc_privkey=hn_priv_key)
        elif profile == 'B':
            if isinstance(raw_keypair, (tuple, list)) and len(raw_keypair) == 2:
                self.EC = ECDH_SECP256R1(raw_keypair=raw_keypair)
            else:
                self.EC = ECDH_SECP256R1(loc_privkey=hn_priv_key)
        else:
            raise CMException('unknown ECIES profile %s' % profile)

    def unprotect(self, ue_pubkey, ciphertext, mac) -> Optional[bytes]:
        """
        Unprotect - Unprotects the given ciphertext using associated MAC and UE ephemeral
        public key returns the decrypted cleartext bytes buffer or None if MAC verification
        failed

        Args:
            ue_pubkey: emphereal public key
            ciphertext: ciphertext to decode
            mac: encrypted mac

        Returns:
            cleartext: decrypted mac
        """
        shared_key = KDF(ue_pubkey, self.EC.generate_sharedkey(ue_pubkey))
        aes_key, aes_nonce, aes_cnt, mac_key = (
            shared_key[:16],
            shared_key[16:24],
            unpack('>Q', shared_key[24:32])[0],
            shared_key[32:64],
        )
        # verify MAC
        mac_hn = hmac.new(mac_key, ciphertext, hashlib.sha256).digest()
        mac_verif = hmac.compare_digest(mac_hn[0:8], mac)       # noqa: WPS349
        # decrypt
        aes = AES_CTR_cryptography(aes_key, aes_nonce, aes_cnt)
        cleartext = aes.decrypt(ciphertext)
        if mac_verif:
            return cleartext
