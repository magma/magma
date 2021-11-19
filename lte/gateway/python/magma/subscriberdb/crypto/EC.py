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

__all__ = ['X25519', 'ECDH_SECP256R1', 'KDF']

from .utils import CMException, int_from_bytes

# this is a wrapper around the cryptography Python librar and its elliptic curve
# module for ECDH computation
# cryptography is a wrapper around openssl, recent versions support both X25519 and secp256r1

try:
    from cryptography.hazmat.backends                       import default_backend
    from cryptography.hazmat.primitives                     import serialization
    from cryptography.hazmat.primitives.asymmetric          import ec
    from cryptography.hazmat.primitives.asymmetric.x25519   import X25519PrivateKey, X25519PublicKey
    from cryptography.hazmat.primitives.kdf.x963kdf         import X963KDF
    from cryptography.hazmat.primitives                     import hashes
except ImportError:
    raise(ImportError('missing ECC backend: requires cryptography Python library'))
else:
    _backend = default_backend()


class X25519(object):
    """wrapper around Python cryptography library to handle a Diffie-Hellman
    exchange over a Curve25519 elliptic curve
    
    private key and public key are handle as simple bytes buffer
    """
    
    def __init__(self, loc_privkey=None):
        if not loc_privkey:
            self.generate_keypair()
        else:
            self.PrivKey = X25519PrivateKey.from_private_bytes(loc_privkey)
    
    def generate_keypair(self):
        self.PrivKey = X25519PrivateKey.generate()
    
    def get_pubkey(self):
        return self.PrivKey.public_key().public_bytes(
            encoding=serialization.Encoding.Raw,
            format=serialization.PublicFormat.Raw)
    
    def get_privkey(self):
        return self.PrivKey.private_bytes(
            encoding=serialization.Encoding.Raw,
            format=serialization.PrivateFormat.Raw,
            encryption_algorithm=serialization.NoEncryption())
    
    def generate_sharedkey(self, ext_pubkey):
        ExtPubKey = X25519PublicKey.from_public_bytes(ext_pubkey)
        return self.PrivKey.exchange(ExtPubKey)


class ECDH_SECP256R1(object):
    """wrapper around Python cryptography library to handle an ECDH  exchange over 
    a NIST secp256r1 elliptic curve
    
    private key is handled within a DER-encoded PKCS8 structure
    public key is handled as a compressed point bytes buffer according to ANSI X9.62
    """
    
    def __init__(self, loc_privkey=None, _raw_keypair=None):
        if loc_privkey:
            self.PrivKey = serialization.load_der_private_key(
                loc_privkey,
                password=None,
                backend=_backend)
            if not hasattr(self.PrivKey, 'curve') or not isinstance(self.PrivKey.curve, ec.SECP256R1):
                raise(CMException('invalid secp256r1 private key'))
        elif isinstance(_raw_keypair, (tuple, list)) and len(_raw_keypair) == 2:
            self._load_raw_keypair(*_raw_keypair)
        else:
            self.generate_keypair()
    
    def _load_raw_keypair(self, privkey, pubkey):
        self.PrivKey = ec.EllipticCurvePrivateNumbers(
            private_value=int_from_bytes(privkey),
            public_numbers=ec.EllipticCurvePublicKey.from_encoded_point(
                curve=ec.SECP256R1(),
                data=pubkey).public_numbers()).private_key(backend=_backend)
    
    def generate_keypair(self):
        self.PrivKey = ec.generate_private_key(
            curve=ec.SECP256R1(),
            backend=_backend)
    
    def get_pubkey(self):
        return self.PrivKey.public_key().public_bytes(
            format=serialization.PublicFormat.CompressedPoint,
            encoding=serialization.Encoding.X962)
    
    def get_privkey(self):
        return self.PrivKey.private_bytes(
            encoding=serialization.Encoding.DER,
            format=serialization.PrivateFormat.PKCS8,
            encryption_algorithm=serialization.NoEncryption())
    
    def generate_sharedkey(self, ext_pubkey):
        ExtPubKey = ec.EllipticCurvePublicKey.from_encoded_point(
            curve=ec.SECP256R1(),
            data=ext_pubkey)
        return self.PrivKey.exchange(ec.ECDH(), ExtPubKey)


def KDF(sharedinfo, sharedkey):
    return X963KDF(
         algorithm=hashes.SHA256(),
         length=64, # 16 bytes AES key, 16 bytes AES CTR IV, 32 bytes HMAC-SHA-256 key
         sharedinfo=sharedinfo,
         backend=_backend).derive(sharedkey)

