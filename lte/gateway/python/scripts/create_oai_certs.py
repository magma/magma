#!/usr/bin/env python

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

import argparse
import os
import shutil
import socket
import tempfile

import envoy

OPENSSL_BIN = "/usr/bin/openssl"

# Default mme/freediameter (s6a) cert and key file names.
OAI_MME_CA_CERT = "mme.cacert.pem"
OAI_MME_CSR = "mme.csr.pem"
OAI_MME_CA_KEY = "mme.cakey.pem"
OAI_MME_CERT = "mme.cert.pem"
OAI_MME_KEY = "mme.key.pem"


def generate_mme_certs(conf_dir):
    # hostname fqdn needs to match with the idenitity configured
    # in mme_fd.conf
    fqdn = socket.gethostname() + '.magma.com'

    # Check directory already exists, or create it
    os.makedirs(conf_dir, exist_ok=True)

    with tempfile.TemporaryDirectory(prefix="/tmp/cert") as temp_dir:
        # Setup a fake demo CA.
        demo_ca_dir = os.path.join(temp_dir, "demoCA")
        os.mkdir(demo_ca_dir)
        with open(os.path.join(demo_ca_dir, "serial"), 'w') as f:
            f.write('01')
        with open(os.path.join(demo_ca_dir, "index.txt"), 'w'):
            # Just touch the file.
            pass

        # Generate the root certificate required by freediameter.
        mme_cacert = os.path.join(temp_dir, OAI_MME_CA_CERT)
        mme_ca_key = os.path.join(temp_dir, OAI_MME_CA_KEY)
        cmd = "{ssl_bin} req -new -batch -x509 -days 3650 -nodes -newkey" \
              " rsa:1024 -out {mme_cacert} -keyout {mme_ca_key} -subj " \
              "/CN={fqdn}/C=FR/ST=PACA/L=Aix/O=Facebook/OU=CM".format(
            ssl_bin=OPENSSL_BIN, mme_cacert=mme_cacert, mme_ca_key=mme_ca_key,
            fqdn=fqdn,
              )
        envoy.run(cmd)

        # Generate the private key
        mme_key = os.path.join(temp_dir, OAI_MME_KEY)
        cmd = "%s genrsa -out %s 1024" % (OPENSSL_BIN, mme_key)
        envoy.run(cmd)

        # Create a self signed cert signing request
        mme_csr = os.path.join(temp_dir, OAI_MME_CSR)
        cmd = "{ssl_bin} req -new -batch -out {mme_csr} -key {mme_key} " \
              "-subj /CN={fqdn}/C=FR/ST=PACA/L=Aix/O=Facebook/OU=CM".format(
            ssl_bin=OPENSSL_BIN, mme_csr=mme_csr, mme_key=mme_key, fqdn=fqdn,
              )
        envoy.run(cmd)

        # No way to change the location of demoCA dir
        os.chdir(temp_dir)
        # Certify it
        mme_cert = os.path.join(temp_dir, OAI_MME_CERT)
        cmd = "{ssl_bin} ca -cert {mme_cacert} -keyfile {mme_ca_key} " \
              "-in {mme_csr} -out {mme_cert} -outdir {out_dir} " \
              "-batch".format(
                  ssl_bin=OPENSSL_BIN, mme_cacert=mme_cacert,
                  mme_ca_key=mme_ca_key, mme_csr=mme_csr,
                  mme_cert=mme_cert, out_dir=temp_dir,
              )
        envoy.run(cmd)

        # Copy the files to the free diameter directory.
        shutil.copy(mme_ca_key, conf_dir)
        shutil.copy(mme_cacert, conf_dir)
        shutil.copy(mme_cert, conf_dir)
        shutil.copy(mme_key, conf_dir)


def main():
    parser = argparse.ArgumentParser(
        description='Generate mme certs for free diameter',
    )

    parser.add_argument(
        "--conf-dir", "-c",
        default="/usr/local/etc/oai/freeDiameter",
    )
    opts = parser.parse_args()
    generate_mme_certs(opts.conf_dir)


if __name__ == "__main__":
    main()
