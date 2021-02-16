# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

include_directories(${PROJECT_SOURCE_DIR}/common)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../lib/3gpp)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../lib/itti)
include_directories(${PROJECT_BINARY_DIR}/s1ap/${ASN1RELDIR})

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../lib/bstr)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../lib/hashtable)

set(S1AP_C_DIR ${PROJECT_BINARY_DIR}/tasks/s1ap/r15)
include_directories(${S1AP_C_DIR})

add_library(TASK_UDP udp_primitives_server.c)
