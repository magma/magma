#!/usr/bin/env bash

# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# self_sign_certs.sh generates a set of keys and self-signed certificates.
#
# Generated secrets
#   - rootCA.key, rootCA.pem -- certs for trusted root CA
#   - controller.key, controller.crt -- certs for orc8r deployment's public
#     domain name, signed by rootCA.key
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
echo "################"
echo "Creating root CA"
echo "################"
openssl genrsa -out rootCA.key 2048
openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 3650 -out rootCA.pem -subj "/C=US/CN=rootca.$domain"

echo ""
echo "########################"
echo "Creating controller cert"
echo "########################"
openssl genrsa -out controller.key 2048
openssl req -new -key controller.key -out controller.csr -subj "/C=US/CN=*.$domain"


# Create an extension config file
> ${domain}.ext cat <<-EOF
basicConstraints=CA:FALSE
subjectAltName = @alt_names
[alt_names]
DNS.1 = *.$domain
DNS.2 = *.nms.$domain
EOF
openssl x509 -req -in controller.csr -CA rootCA.pem -CAkey rootCA.key -CAcreateserial -out controller.crt -days 825 -sha256 -extfile ${domain}.ext

echo ""
echo "###########################"
echo "Deleting intermediate files"
echo "###########################"
rm -f controller.csr rootCA.srl ${domain}.ext
