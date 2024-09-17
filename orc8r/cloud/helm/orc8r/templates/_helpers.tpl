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
{{/* Generate basic labels */}}
{{- define "labels" -}}
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: helm
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
{{- end -}}

{{/* Label for orc8r non-applications */}}
{{- define "orc8r-labels" -}}
app.kubernetes.io/part-of: orc8r
{{- end -}}

{{/* Label for orc8r applications */}}
{{- define "orc8r-app-labels" -}}
app.kubernetes.io/part-of: orc8r-app
{{- end -}}

{{/* Generate selector labels */}}
{{- define "selector-labels" -}}
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/* Generate selector labels */}}
{{- define "nginx-image-version-label" -}}
{{- end -}}
