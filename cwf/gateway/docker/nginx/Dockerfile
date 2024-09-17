################################################################################
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
################################################################################
FROM nginx

# install openssl to generate default key
RUN apt update -y && apt install openssl -y

# prepare a place for nginx ssl keys
RUN mkdir -p /etc/nginx/certs
RUN chown root:nginx /etc/nginx/certs
RUN chmod g+rx /etc/nginx/certs
RUN chmod g-w /etc/nginx/certs
RUN chmod o-rwx /etc/nginx/certs

# create root CA certificate
RUN openssl genrsa -out ngxRootCA.key 2048
RUN openssl req -x509 -new -nodes -key ngxRootCA.key -sha256 -days 3650 \
    -out ngxRootCA.pem -subj "/C=US/ST=California/L=Menlo Park/O=Facebook/OU=FBC/CN=rootca.magma.test"

# create a default key for localhost to be overridden in production
RUN openssl genrsa -out ingress_proxy.key 2048
RUN openssl req -new -key ingress_proxy.key -out ingress_proxy.csr \
    -subj "/C=US/ST=California/L=Menlo Park/O=Facebook/OU=FBC/CN=*.magma.test"
RUN openssl x509 -req -in ingress_proxy.csr -CA ngxRootCA.pem -CAkey ngxRootCA.key \
    -CAcreateserial -out ingress_proxy.crt -days 3650 -sha256

# clean up and copy root ca for controller_proxy
RUN rm -f ngxRootCA.key ingress_proxy.csr
RUN mv ingress_proxy.key ingress_proxy.crt ngxRootCA.pem /etc/nginx/certs/
