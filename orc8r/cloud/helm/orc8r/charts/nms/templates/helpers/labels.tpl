{{/*
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
*/}}

{{- define "nms.labels" -}}
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
{{- define "nms.selector-labels" -}}
{{- $envAll := index . 0 -}}
{{- $application := index . 1 -}}
{{- $component := index . 2 -}}
release_group: {{ $envAll.Values.release_group | default $envAll.Release.Name }}
app.kubernetes.io/name: {{ $application }}
app.kubernetes.io/component: {{ $component }}
app.kubernetes.io/instance: {{ $envAll.Release.Name }}
{{- end -}}
