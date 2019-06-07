{{/*
# Copyright (c) 2018-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.
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
