#!/bin/bash

# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

echo "Creating postgresql DB dp_test"
psql -U "${POSTGRES_USER}" -c "CREATE DATABASE dp_test"
echo "Grant privileges to ${POSTGRES_USER} on dp_test..."
psql -U "${POSTGRES_USER}" -c "GRANT ALL PRIVILEGES ON DATABASE dp_test TO ${POSTGRES_USER}"
