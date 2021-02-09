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
RUN apt update && apt install openssl

# prepare a place for nginx ssl keys
RUN mkdir -p /etc/nginx/ssl
RUN chown root:nginx /etc/nginx/ssl
RUN chmod g+rx /etc/nginx/ssl
RUN chmod g-w /etc/nginx/ssl
RUN chmod o-rwx /etc/nginx/ssl

# create a default ssl key for localhost
RUN openssl req -newkey rsa:2048 -nodes -x509 -days 3650 \
    -keyout /etc/nginx/ssl/aptly.key \
    -out /etc/nginx/ssl/aptly.crt \
    -subj "/C=US/ST=California/L=Menlo Park/O=Test Organization/OU=TEST/CN=localhost"
