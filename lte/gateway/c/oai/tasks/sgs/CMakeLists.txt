# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

find_package(Protobuf REQUIRED)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

add_library(TASK_SGS
    sgs_task.c
    sgs_service_handler.c
    )

target_link_libraries(TASK_SGS
    COMMON
    LIB_BSTR LIB_HASHTABLE LIB_SGS_CLIENT
    ${PROTOBUF_LIBRARIES}
)
target_include_directories(TASK_SGS PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
)
