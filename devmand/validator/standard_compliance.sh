#!/bin/bash

# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# open config has some non-standard complaint yang models. This script takes in
# a model and makes it compliant.

function compliance {
  local pattern="$1"
  sed -i'' "${pattern}" "${model}"
}

model="$1"

# no when statement on state data
compliance 's;when "oc-if:state/oc-if:type;when "oc-if:config/oc-if:type;g'

# Use wc3 regex not posix
compliance "s;\(\s'\)\^;\1;g"
compliance "s;\$\('\;\);\1;g"
