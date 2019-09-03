{{/*
# Copyright (c) 2018-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.
*/}}

{{- define "kubernetes_resources" -}}
{{- $envAll := index . 0 -}}
{{- $component := index . 1 -}}
{{- if $envAll.Values.pod.resources.enabled -}}
resources:
  {{- if or $component.limits.cpu $component.limits.memory }}
  limits:
    {{- if $component.limits.cpu }}
    cpu: {{ $component.limits.cpu | quote }}
    {{- end }}
    {{- if $component.limits.memory }}
    memory: {{ $component.limits.memory | quote }}
    {{- end }}
  {{- end }}
  {{- if or $component.requests.cpu $component.requests.memory }}
  requests:
    {{- if $component.requests.cpu }}
    cpu: {{ $component.requests.cpu | quote }}
    {{- end }}
    {{- if $component.requests.memory }}
    memory: {{ $component.requests.memory | quote }}
    {{- end }}
  {{- end }}
{{- end -}}
{{- end -}}
