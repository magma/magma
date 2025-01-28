{{/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/}}
{{- define "labels" -}}
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: helm
app.kubernetes.io/part-of: magma
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
{{- end -}}

{{- define "orchestrator-config-template" -}}
useGRPCExporter: true
prometheusGRPCPushAddress: "{{ .Release.Name }}-prometheus-cache:9092"
# Comment out above line, uncomment below, and set useGRPCExporter to false
# to switch to HTTP metric pushing (less efficient)
# prometheusPushAddresses:
#  - "http://{{ .Release.Name }}-prometheus-cache:9091/metrics"
{{- end -}}

{{- define "metricsd-thanos-config-template" -}}
profile: "prometheus"
prometheusQueryAddress: "http://{{ .Release.Name }}-thanos-query-http:10902"
alertmanagerApiURL: "http://{{ .Release.Name }}-alertmanager:9093/api/v2"
prometheusConfigServiceURL: "http://{{ .Release.Name }}-prometheus-configurer:9100/v1"
alertmanagerConfigServiceURL: "http://{{ .Release.Name }}-alertmanager-configurer:9101/v1"
{{- end -}}

{{- define "metricsd-config-template" -}}
profile: "prometheus"
prometheusQueryAddress: "http://{{ .Release.Name }}-prometheus:9090"
alertmanagerApiURL: "http://{{ .Release.Name }}-alertmanager:9093/api/v2"
prometheusConfigServiceURL: "http://{{ .Release.Name }}-prometheus-configurer:9100/v1"
alertmanagerConfigServiceURL: "http://{{ .Release.Name }}-alertmanager-configurer:9101/v1"
{{- end -}}