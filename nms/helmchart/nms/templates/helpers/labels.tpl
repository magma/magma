{{/*
Copyright 2019 Mirantis Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/}}

{{- define "labels" -}}
{{- $envAll := index . 0 -}}
{{- $application := index . 1 -}}
{{- $component := index . 2 -}}
release_group: {{ $envAll.Values.release_group | default $envAll.Release.Name }}
app.kubernetes.io/name: {{ $application }}
app.kubernetes.io/component: {{ $component }}
app.kubernetes.io/instance: {{ $envAll.Release.Name }}
app.kubernetes.io/managed-by: helm
app.kubernetes.io/part-of: magma
{{- end -}}

{{/* Generate selector labels */}}
{{- define "selector-labels" -}}
{{- $envAll := index . 0 -}}
{{- $application := index . 1 -}}
{{- $component := index . 2 -}}
release_group: {{ $envAll.Values.release_group | default $envAll.Release.Name }}
app.kubernetes.io/name: {{ $application }}
app.kubernetes.io/component: {{ $component }}
app.kubernetes.io/instance: {{ $envAll.Release.Name }}
{{- end -}}
