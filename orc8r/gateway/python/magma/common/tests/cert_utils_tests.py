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

import base64
import datetime
import os
from tempfile import TemporaryDirectory
from unittest import TestCase

import magma.common.cert_utils as cu
from cryptography import x509
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import ec


class CertUtilsTest(TestCase):
    def test_key(self):
        with TemporaryDirectory(prefix='/tmp/test_cert_utils') as temp_dir:
            key = ec.generate_private_key(ec.SECP384R1(), default_backend())
            cu.write_key(key, os.path.join(temp_dir, 'test.key'))
            key_load = cu.load_key(os.path.join(temp_dir, 'test.key'))

        key_bytes = key.private_bytes(
            serialization.Encoding.PEM,
            serialization.PrivateFormat.TraditionalOpenSSL,
            serialization.NoEncryption(),
        )
        key_load_bytes = key_load.private_bytes(
            serialization.Encoding.PEM,
            serialization.PrivateFormat.TraditionalOpenSSL,
            serialization.NoEncryption(),
        )
        self.assertEqual(key_bytes, key_load_bytes)

    def load_public_key_to_base64der(self):
        with TemporaryDirectory(prefix='/tmp/test_cert_utils') as temp_dir:
            key = ec.generate_private_key(ec.SECP384R1(), default_backend())
            cu.write_key(key, os.path.join(temp_dir, 'test.key'))
            base64der = cu.load_public_key_to_base64der(
                os.path.join(temp_dir, 'test.key'),
            )
            der = base64.b64decode(base64der)
            pub_key = serialization.load_der_public_key(der, default_backend())
            self.assertEqual(pub_key, key.public_key())

    def test_csr(self):
        key = ec.generate_private_key(ec.SECP384R1(), default_backend())
        csr = cu.create_csr(
            key, 'i am dummy test',
            'US', 'CA', 'MPK', 'FB', 'magma', 'magma@fb.com',
        )
        self.assertTrue(csr.is_signature_valid)
        public_key_bytes = key.public_key().public_bytes(
            serialization.Encoding.OpenSSH,
            serialization.PublicFormat.OpenSSH,
        )
        csr_public_key_bytes = csr.public_key().public_bytes(
            serialization.Encoding.OpenSSH,
            serialization.PublicFormat.OpenSSH,
        )
        self.assertEqual(public_key_bytes, csr_public_key_bytes)

    def test_cert(self):
        with TemporaryDirectory(prefix='/tmp/test_cert_utils') as temp_dir:
            cert = _create_dummy_cert()
            cert_file = os.path.join(temp_dir, 'test.cert')
            cu.write_cert(
                cert.public_bytes(
                    serialization.Encoding.DER,
                ), cert_file,
            )
            cert_load = cu.load_cert(cert_file)
        self.assertEqual(cert, cert_load)


def _create_dummy_cert():
    key = ec.generate_private_key(ec.SECP384R1(), default_backend())
    subject = issuer = x509.Name([
        x509.NameAttribute(x509.oid.NameOID.COUNTRY_NAME, u"US"),
        x509.NameAttribute(x509.oid.NameOID.STATE_OR_PROVINCE_NAME, u"CA"),
        x509.NameAttribute(x509.oid.NameOID.LOCALITY_NAME, u"San Francisco"),
        x509.NameAttribute(x509.oid.NameOID.ORGANIZATION_NAME, u"My Company"),
        x509.NameAttribute(x509.oid.NameOID.COMMON_NAME, u"mysite.com"),
    ])
    cert = x509.CertificateBuilder().subject_name(
        subject,
    ).issuer_name(
        issuer,
    ).public_key(
        key.public_key(),
    ).serial_number(
        x509.random_serial_number(),
    ).not_valid_before(
        datetime.datetime.utcnow(),
    ).not_valid_after(
        datetime.datetime.utcnow() + datetime.timedelta(days=10),
    ).sign(key, hashes.SHA256(), default_backend())
    return cert
