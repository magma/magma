# Copyright 2020 The Magma Authors.
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

add_subdirectory(mme_app_task)
add_subdirectory(mobility_client)
add_subdirectory(ngap)
add_subdirectory(itti)
add_subdirectory(amf)
add_subdirectory(n11)
add_subdirectory(s1ap_task)
add_subdirectory(mock_tasks)
add_subdirectory(nas)
add_subdirectory(lib)
if (EMBEDDED_SGW)
  add_subdirectory(openflow)
  add_subdirectory(spgw_task)
  add_subdirectory(pipelined_client)
  add_subdirectory(sgw_s8_task)
else (EMBEDDED_SGW)
  add_subdirectory(gtpv2-c)
endif (EMBEDDED_SGW)
