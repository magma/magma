#!/bin/bash
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -u -e
set -o pipefail
apt-get update
apt-get install -y graphviz autoconf automake bzip2 debhelper dh-autoreconf libssl-dev libtool openssl procps python-all python-twisted-conch python-zopeinterface python-six build-essential fakeroot
{{ if .Values.virtlet.extra }}
{{ .Values.virtlet.extra }}
{{- end }}
