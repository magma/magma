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

include_directories(${CMAKE_CURRENT_SOURCE_DIR})
include_directories(${CMAKE_CURRENT_BINARY_DIR})
include_directories(${ORC8R_CPP_OUT_DIR})

add_library(redis_utils redis_client.cpp)
target_link_libraries(redis_utils ${CONFIG} COMMON cpp_redis tacopie protobuf)


target_include_directories(redis_utils PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${MAGMA_LIB_DIR}/config
)
