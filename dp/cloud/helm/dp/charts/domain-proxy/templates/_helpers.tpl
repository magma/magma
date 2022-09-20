{{/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/}}

{{/*
Expand the name of the chart.
*/}}
{{- define "domain-proxy.name" -}}
{{- default .Chart.Name .Values.dp.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "domain-proxy.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Match labels
*/}}
{{- define "domain-proxy.common.matchLabels" -}}
app.kubernetes.io/name: {{ include "domain-proxy.name" . }}
app.kubernetes.io/release: {{ .Release.Name }}
{{- end }}

{{/*
Meta labels
*/}}
{{- define "domain-proxy.common.metaLabels" -}}
helm.sh/chart: {{ include "domain-proxy.chart" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Configuration controller match labels
*/}}
{{- define "domain-proxy.configuration_controller.matchLabels" -}}
component: {{ .Values.dp.configuration_controller.name | quote }}
{{ include "domain-proxy.common.matchLabels" . }}
{{- end -}}

{{/*
Configuration controller labels
*/}}
{{- define "domain-proxy.configuration_controller.labels" -}}
{{ include "domain-proxy.configuration_controller.matchLabels" . }}
{{ include "domain-proxy.common.metaLabels" . }}
{{- end -}}

{{/*
Radio controller match labels
*/}}
{{- define "domain-proxy.radio_controller.matchLabels" -}}
component: {{ .Values.dp.radio_controller.name | quote }}
{{ include "domain-proxy.common.matchLabels" . }}
{{- end -}}

{{/*
Radio controller labels
*/}}
{{- define "domain-proxy.radio_controller.labels" -}}
{{ include "domain-proxy.radio_controller.matchLabels" . }}
{{ include "domain-proxy.common.metaLabels" . }}
{{- end -}}

{{/*
DB service labels
*/}}
{{- define "domain-proxy.db_service.labels" -}}
{{ include "domain-proxy.common.metaLabels" . }}
{{- end -}}

{{/*
fluentd match labels
*/}}
{{- define "domain-proxy.fluentd.matchLabels" -}}
component: {{ .Values.dp.fluentd.name | quote }}
{{ include "domain-proxy.common.matchLabels" . }}
{{- end -}}

{{/*
fluentd labels
*/}}
{{- define "domain-proxy.fluentd.labels" -}}
{{ include "domain-proxy.fluentd.matchLabels" . }}
{{ include "domain-proxy.common.metaLabels" . }}
{{- end -}}

{{/*
grafana match labels
*/}}
{{- define "domain-proxy.grafana.matchLabels" -}}
component: {{ .Values.dp.grafana.name | quote }}
{{ include "domain-proxy.common.matchLabels" . }}
{{- end -}}

{{/*
grafana labels
*/}}
{{- define "domain-proxy.grafana.labels" -}}
{{ include "domain-proxy.grafana.matchLabels" . }}
{{ include "domain-proxy.common.metaLabels" . }}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "domain-proxy.fullname" -}}
{{- if .Values.dp.fullnameOverride }}
{{- .Values.dp.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.dp.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}


{{/*
Create a fully qualified configuration_controller name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}

{{- define "domain-proxy.configuration_controller.fullname" -}}
{{- if .Values.dp.configuration_controller.fullnameOverride -}}
{{- .Values.dp.configuration_controller.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.dp.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- printf "%s-%s" .Release.Name .Values.dp.configuration_controller.name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-%s" .Release.Name $name .Values.dp.configuration_controller.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create a fully qualified radio_controller name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}

{{- define "domain-proxy.radio_controller.fullname" -}}
{{- if .Values.dp.radio_controller.fullnameOverride -}}
{{- .Values.dp.radio_controller.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.dp.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- printf "%s-%s" .Release.Name .Values.dp.radio_controller.name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-%s" .Release.Name $name .Values.dp.radio_controller.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create a fully qualified db_service name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}

{{- define "domain-proxy.db_service.fullname" -}}
{{- if .Values.dp.db_service.fullnameOverride -}}
{{- .Values.dp.db_service.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.dp.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- printf "%s-%s" .Release.Name .Values.dp.db_service.name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-%s" .Release.Name $name .Values.dp.db_service.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}


{{/*
Create a fully qualified fluentd name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}

{{- define "domain-proxy.fluentd.fullname" -}}
{{- if .Values.dp.fluentd.fullnameOverride -}}
{{- .Values.dp.fluentd.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.dp.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- printf "%s-%s" .Release.Name .Values.dp.fluentd.name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-%s" .Release.Name $name .Values.dp.fluentd.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create a fully qualified grafana name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}

{{- define "domain-proxy.grafana.fullname" -}}
{{- if .Values.dp.grafana.fullnameOverride -}}
{{- .Values.dp.grafana.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.dp.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- printf "%s-%s" .Release.Name .Values.dp.grafana.name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-%s" .Release.Name $name .Values.dp.grafana.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Return the appropriate apiVersion for deployment.
*/}}
{{- define "domain-proxy.deployment.apiVersion" -}}
{{- print "apps/v1" -}}
{{- end -}}

{{/*
Return the appropriate apiVersion for ingress.
*/}}
{{- define "domain-proxy.ingress.apiVersion" -}}
{{- print "networking.k8s.io/v1beta1" -}}
{{- end -}}

{{/*
Return the appropriate apiVersion for job.
*/}}
{{- define "domain-proxy.job.apiVersion" -}}
{{- print "batch/v1" -}}
{{- end -}}

{{/*
Return the appropriate apiVersion for ingress.
*/}}
{{- define "ingress.apiVersion" -}}
{{- print "networking.k8s.io/v1" -}}
{{- end -}}

{{/*
Return the appropriate apiVersion for HTTPProxy.
*/}}
{{- define "httpproxy.apiVersion" -}}
{{- print "projectcontour.io/v1" -}}
{{- end -}}

{{/*
Create the name of the service account to use for configuration controller
*/}}
{{- define "domain-proxy.configuration_controller.serviceAccountName" -}}
{{- if .Values.dp.configuration_controller.serviceAccount.create }}
{{- default (include "domain-proxy.fullname" .) .Values.dp.configuration_controller.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.dp.configuration_controller.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the service account to use for radio controller
*/}}
{{- define "domain-proxy.radio_controller.serviceAccountName" -}}
{{- if .Values.dp.radio_controller.serviceAccount.create }}
{{- default (include "domain-proxy.fullname" .) .Values.dp.radio_controller.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.dp.radio_controller.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the service account to use for db service
*/}}
{{- define "domain-proxy.db_service.serviceAccountName" -}}
{{- if .Values.dp.db_service.serviceAccount.create }}
{{- default (include "domain-proxy.fullname" .) .Values.dp.db_service.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.dp.db_service.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the service account to use for fluentd
*/}}
{{- define "domain-proxy.fluentd.serviceAccountName" -}}
{{- if .Values.dp.fluentd.serviceAccount.create }}
{{- default (include "domain-proxy.fullname" .) .Values.dp.fluentd.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.dp.fluentd.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the service account to use for grafana
*/}}
{{- define "domain-proxy.grafana.serviceAccountName" -}}
{{- if .Values.dp.grafana.serviceAccount.create }}
{{- default (include "domain-proxy.fullname" .) .Values.dp.grafana.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.dp.grafana.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Define the domain-proxy.namespace template
*/}}
{{- define "domain-proxy.namespace" -}}
{{ printf "namespace: %s" .Release.Namespace }}
{{- end -}}
