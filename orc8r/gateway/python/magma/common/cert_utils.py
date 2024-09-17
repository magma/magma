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

from cryptography import x509
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.x509.oid import NameOID
from magma.common.serialization_utils import write_to_file_atomically


def load_key(key_file):
    """Load a private key encoded in PEM format

    Args:
        key_file: path to the key file

    Returns:
        RSAPrivateKey or EllipticCurvePrivateKey depending on the contents of key_file

    Raises:
        IOError: If file cannot be opened
        ValueError: If the file content cannot be decoded successfully
        TypeError: If the key_file is encrypted
    """
    with open(key_file, 'rb') as f:
        key_bytes = f.read()
    return serialization.load_pem_private_key(
        key_bytes, None, default_backend(),
    )


def write_key(key, key_file):
    """Write key object to file in PEM format atomically

    Args:
        key: RSAPrivateKey or EllipticCurvePrivateKey object
        key_file: path to the key file
    """
    key_pem = key.private_bytes(
        serialization.Encoding.PEM,
        serialization.PrivateFormat.TraditionalOpenSSL,
        serialization.NoEncryption(),
    )
    write_to_file_atomically(key_file, key_pem.decode("utf-8"))


def load_public_key_to_base64der(key_file):
    """Load the public key of a private key and convert to base64 encoded DER
    The return value can be used directly for device registration.

    Args:
        key_file: path to the private key file, pem encoded

    Returns:
        base64 encoded public key in DER format

    Raises:
        IOError: If file cannot be opened
        ValueError: If the file content cannot be decoded successfully
        TypeError: If the key_file is encrypted
    """

    key = load_key(key_file)
    pub_key = key.public_key()
    pub_bytes = pub_key.public_bytes(
        encoding=serialization.Encoding.DER,
        format=serialization.PublicFormat.SubjectPublicKeyInfo,
    )
    encoded = base64.b64encode(pub_bytes)
    return encoded


def create_csr(
    key, common_name,
    country=None, state=None, city=None, org=None,
    org_unit=None, email_address=None,
):
    """Create csr and sign it with key.

    Args:
        key: RSAPrivateKey or EllipticCurvePrivateKey object
        common_name: common name
        country: country
        state: state or province
        city: city
        org: organization
        org_unit: organizational unit
        email_address: email address

    Returns:
        csr: x509.CertificateSigningRequest
    """
    name_attrs = [x509.NameAttribute(NameOID.COMMON_NAME, common_name)]
    if country:
        name_attrs.append(x509.NameAttribute(NameOID.COUNTRY_NAME, country))
    if state:
        name_attrs.append(
            x509.NameAttribute(NameOID.STATE_OR_PROVINCE_NAME, state),
        )
    if city:
        name_attrs.append(x509.NameAttribute(NameOID.LOCALITY_NAME, city))
    if org:
        name_attrs.append(x509.NameAttribute(NameOID.ORGANIZATION_NAME, org))
    if org_unit:
        name_attrs.append(
            x509.NameAttribute(NameOID.ORGANIZATIONAL_UNIT_NAME, org_unit),
        )
    if email_address:
        name_attrs.append(
            x509.NameAttribute(NameOID.EMAIL_ADDRESS, email_address),
        )

    csr = x509.CertificateSigningRequestBuilder().subject_name(
        x509.Name(name_attrs),
    ).sign(key, hashes.SHA256(), default_backend())

    return csr


def load_cert(cert_file):
    """Load certificate from a file

    Args:
        cert_file: path to file storing the cert in PEM format

    Returns:
        cert: an instance of x509.Certificate

    Raises:
        IOError: If file cannot be opened
        ValueError: If the file content cannot be decoded successfully
    """
    with open(cert_file, 'rb') as f:
        cert_pem = f.read()
    cert = x509.load_pem_x509_certificate(cert_pem, default_backend())
    return cert


def write_cert(cert_der, cert_file):
    """Write DER encoded cert to file in PEM format

    Args:
        cert_der: certificate encoded in DER format
        cert_file: path to certificate
    """
    cert = x509.load_der_x509_certificate(cert_der, default_backend())
    cert_pem = cert.public_bytes(serialization.Encoding.PEM)
    write_to_file_atomically(cert_file, cert_pem.decode("utf-8"))
