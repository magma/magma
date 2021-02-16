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

include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

set(MAGMA_LIB_DIR $ENV{C_BUILD}/magma_common)
include_directories("${OUTPUT_DIR}")

include_directories(${MAGMA_LIB_DIR}/async_grpc)
include_directories(${MAGMA_LIB_DIR}/eventd)

link_directories(
        ${MAGMA_LIB_DIR}/async_grpc
        ${MAGMA_LIB_DIR}/eventd
)

add_library(LIB_EVENT_CLIENT
        EventClientAPI.cpp
        )

target_link_libraries(LIB_EVENT_CLIENT
        COMMON
        LIB_BSTR ${EVENTD}
        ${ASYNC_GRPC}
        )

target_include_directories(LIB_EVENT_CLIENT PUBLIC
        ${CMAKE_CURRENT_SOURCE_DIR}
        ASYNC_GRPC COMMON SERVICE_REGISTRY EVENTD
        )
