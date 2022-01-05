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

# this is a wrapper around the cryptography Python library and its elliptic curve
# module for ECDH computation
# cryptography is a wrapper around openssl, recent versions support both X25519 and secp256r1
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives import hashes  # noqa: E501
from cryptography.hazmat.primitives import serialization  # noqa: E501
from cryptography.hazmat.primitives.asymmetric import ec  # noqa: E501
from cryptography.hazmat.primitives.asymmetric.x25519 import (
    X25519PrivateKey,
    X25519PublicKey,
)
from cryptography.hazmat.primitives.kdf.x963kdf import X963KDF  # noqa: E501
from magma.subscriberdb.crypto.utils import (  # noqa: E501
    CMException,
    int_from_bytes,
)

_backend = default_backend()


def KDF(sharedinfo, sharedkey):
    """
    Create object of X963KDF

    Args:
        sharedinfo: sharedinfo
        sharedkey: sharedkey

    Returns:
        Object of X963KDF
    """
    return X963KDF(
         algorithm=hashes.SHA256(),     # noqa: E126
         length=64,  # 16 bytes AES key, 16 bytes AES CTR IV, 32 bytes HMAC-SHA-256 key
         sharedinfo=sharedinfo,
         backend=_backend,
    ).derive(sharedkey)


class X25519(object):
    """wrapper around Python cryptography library to handle a Diffie-Hellman
    exchange over a Curve25519 elliptic curve

    private key and public key are handle as simple bytes buffer
    """

    def __init__(self, loc_privkey=None):
        """
        Create the class object

        Args:
            loc_privkey: loc_privkey

        Returns: None
        """
        if not loc_privkey:
            self.generate_keypair()
        else:
            self.priv_key = X25519PrivateKey.from_private_bytes(loc_privkey)

    def generate_keypair(self):
        """
        generate_keypair- creates the public/priv key pair

        Args: None

        Returns: None
        """
        self.priv_key = X25519PrivateKey.generate()

    def get_pubkey(self):
        """
        get_pubkey - get the public key
        """
        return self.priv_key.public_key().public_bytes(
            encoding=serialization.Encoding.Raw,
            format=serialization.PublicFormat.Raw,
        )

    def get_privkey(self):
        """
        get_privkey - get the private key
        """
        return self.priv_key.private_bytes(
            encoding=serialization.Encoding.Raw,
            format=serialization.PrivateFormat.Raw,
            encryption_algorithm=serialization.NoEncryption(),
        )

    def generate_sharedkey(self, ext_pubkey):
        """
        generate_sharedkey - get the shared key
        """
        extpubkey = X25519PublicKey.from_public_bytes(ext_pubkey)
        return self.priv_key.exchange(extpubkey)


class ECDH_SECP256R1(object):       # noqa: N801
    """wrapper around Python cryptography library to handle an ECDH  exchange over
    a NIST secp256r1 elliptic curve
    private key is handled within a DER-encoded PKCS8 structure
    public key is handled as a compressed point bytes buffer according to ANSI X9.62
    """

    def __init__(self, loc_privkey=None, raw_keypair=None):
        """
        Init - Initializes the class object

        Args:
            loc_privkey: loc_privkey
            raw_keypair: raw_keypair

        Returns: None

        Raises:
            CMException: invalid secp256r1 private key
        """
        if loc_privkey:
            self.PrivKey = serialization.load_der_private_key(
                loc_privkey,
                password=None,
                backend=_backend,
            )
            if not hasattr(self.PrivKey, 'curve') or not isinstance(self.PrivKey.curve, ec.SECP256R1):
                raise CMException('invalid secp256r1 private key')
        elif isinstance(raw_keypair, (tuple, list)) and len(raw_keypair) == 2:
            self._load_raw_keypair(*raw_keypair)
        else:
            self.generate_keypair()

    def _load_raw_keypair(self, privkey, pubkey):
        """
        _load_raw_keypair - Initializes the priv key
        """
        self.PrivKey = ec.EllipticCurvePrivateNumbers(
                        private_value=int_from_bytes(privkey),  # noqa: E126
                        public_numbers=ec.EllipticCurvePublicKey.from_encoded_point(
                        curve=ec.SECP256R1(),           # noqa: E122
                        data=pubkey,                    # noqa: E122
                        ).public_numbers(),
        ).private_key(backend=_backend)

    def generate_keypair(self):
        """
        generate_keypair - Generates public/private keypair
        """
        self.PrivKey = ec.generate_private_key(
            curve=ec.SECP256R1(),
            backend=_backend,
        )

    def get_pubkey(self):
        """
        get_pubkey - Generates public key
        """
        return self.PrivKey.public_key().public_bytes(
            format=serialization.PublicFormat.CompressedPoint,
            encoding=serialization.Encoding.X962,
        )

    def get_privkey(self):
        """
        get_privkey - Generates private key
        """
        return self.PrivKey.private_bytes(
            encoding=serialization.Encoding.DER,
            format=serialization.PrivateFormat.PKCS8,
            encryption_algorithm=serialization.NoEncryption(),
        )

    def generate_sharedkey(self, ext_pubkey):
        """
        generate_sharedkey - Generates shared key
        """
        extpubkey = ec.EllipticCurvePublicKey.from_encoded_point(
            curve=ec.SECP256R1(),
            data=ext_pubkey,
        )
        return self.PrivKey.exchange(ec.ECDH(), extpubkey)
