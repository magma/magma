---
#
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# nghttpx config will be generated here and used
nghttpx_config_location: /var/tmp/nghttpx.conf

# Location for certs
rootca_cert: /var/opt/magma/certs/rootCA.pem
gateway_cert: /var/opt/magma/certs/gateway.crt
gateway_key: /var/opt/magma/certs/gateway.key

# Listening port of the proxy for local services. The port would be closed
# for the rest of the world.
local_port: {{ .Values.cwf.proxy.local_port }}

# Cloud address for reaching out to the cloud.
cloud_address: {{ .Values.cwf.proxy.cloud_address }}
cloud_port: {{ .Values.cwf.proxy.cloud_port }}

bootstrap_address: {{ .Values.cwf.proxy.bootstrap_address }}
bootstrap_port: {{ .Values.cwf.proxy.bootstrap_port }}

fluentd_address: {{ .Values.cwf.log_aggregation.fluentd_address }}
fluentd_port: {{ .Values.cwf.log_aggregation.fluentd_port }}

# Option to use nghttpx for proxying. If disabled, the individual
# services would establish the TLS connections themselves.
proxy_cloud_connections: True

# Allows http_proxy usage if the environment variable is present
allow_http_proxy: True
