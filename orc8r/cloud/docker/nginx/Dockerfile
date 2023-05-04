# Copyright 2023 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM nginx:1.23

RUN apt-get update && \
  apt-get install -y python3-pip daemontools
RUN pip3 install --no-cache-dir PyYAML jinja2

RUN mkdir -p /var/opt/magma/envdir
COPY configs /etc/magma/configs

WORKDIR /tmp

ARG NGINX_FILES=src/magma/orc8r/cloud/docker/nginx

COPY ${NGINX_FILES}/templates /etc/magma/templates
COPY ${NGINX_FILES}/generate_nginx_configs.py /usr/local/bin/generate_nginx_configs.py
COPY ${NGINX_FILES}/create_test_proxy_certs /usr/local/bin/create_test_proxy_certs

CMD /usr/local/bin/create_test_proxy_certs && envdir /var/opt/magma/envdir /usr/local/bin/generate_nginx_configs.py && exec nginx -g 'daemon off;'
