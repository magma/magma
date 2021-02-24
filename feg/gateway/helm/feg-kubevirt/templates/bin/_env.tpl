# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

COMPOSE_PROJECT_NAME=feg
DOCKER_REGISTRY={{ .Values.feg.image.docker_registry }}
DOCKER_USERNAME={{ .Values.feg.image.username }}
DOCKER_PASSWORD={{ .Values.feg.image.password }}
IMAGE_VERSION={{ .Values.feg.image.tag }}
BUILD_CONTEXT={{ .Values.feg.repo.url }}#{{ .Values.feg.repo.branch }}

ROOTCA_PATH=/var/opt/magma/certs/rootCA.pem
CONTROL_PROXY_PATH=/etc/magma/control_proxy.yml
SNOWFLAKE_PATH=/etc/snowflake
CONFIGS_DEFAULT_VOLUME=/etc/magma
CONFIGS_TEMPLATES_PATH=/etc/magma/templates
CERTS_VOLUME=/var/opt/magma/certs
CONFIGS_VOLUME=/var/opt/magma/configs

{{ if .Values.feg.log_aggregation.enabled }}
LOG_DRIVER=fluentd
{{ else }}
LOG_DRIVER=journald
{{- end }}

{{ if .Values.feg.env }}
{{ .Values.feg.env }}
{{- end }}
