#!/usr/bin/env bash

# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# create_application_certs.sh generates application certs for orc8r from
# existing certificates.
#
# Generated secrets
#   - bootstrapper.key -- private key for controller's bootstrapper service,
#     used in gateway bootstrapping challenges
#   - fluentd.key, fluentd.pem -- certs for fluentd endpoint, allowing
#     gateways to securely send logs (fluentd is currently outside orc8r proxy)
#   - certifier.key, certifier.pem -- certs for the controller's certifier
#     service, providing more fine-grained access to controller services
#   - admin_operator.key.pem, admin_operator.pem -- client certs for the
#     initial admin operator (e.g. whoever's deploying orc8r) to authenticate
#     to the orc8r proxy
#
# NOTE: extension naming for certs and keys is non-normalized here due to
# expected naming conventions from other parts of the code. For reference,
# all outputs are PEM-encoded, *.key is a private key, and *.crt and missing
# indicator are both certificates.

usage() {
  echo "Usage: $0 DOMAIN_NAME"
  exit 2
}

domain="$1"
if [[ ${domain} == "" ]]; then
    usage
fi

echo ""
echo "#########################"
echo "Creating bootstrapper key"
echo "#########################"
openssl genrsa -out bootstrapper.key 2048

echo ""
echo "######################"
echo "Creating fluentd certs"
echo "######################"
openssl genrsa -out fluentd.key 2048
openssl req -x509 -new -nodes -key fluentd.key -sha256 -days 3650 -out fluentd.pem -subj "/C=US/CN=fluentd.$domain"

echo ""
echo "#####################"
echo "Creating certifier CA"
echo "#####################"
openssl genrsa -out certifier.key 2048
openssl req -x509 -new -nodes -key certifier.key -sha256 -days 3650 -out certifier.pem -subj "/C=US/CN=certifier.$domain"

echo ""
echo "############################"
echo "Creating admin_operator cert"
echo "############################"
openssl genrsa -out admin_operator.key.pem 2048
openssl req -new -key admin_operator.key.pem -out admin_operator.csr -subj "/C=US/CN=admin_operator"
openssl x509 -req -in admin_operator.csr -CA certifier.pem -CAkey certifier.key -CAcreateserial -out admin_operator.pem -days 3650 -sha256

# Export to password-protected PKCS12 bundle, e.g. for import into client
# keychain, with the following command
# openssl pkcs12 -export -inkey admin_operator.key.pem -in admin_operator.pem -out admin_operator.pfx

echo ""
echo "###########################"
echo "Deleting intermediate files"
echo "###########################"
rm -f admin_operator.csr certifier.srl
