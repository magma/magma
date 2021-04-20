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

file(GLOB PROTO_SRCS "orc8r/protos/*.pb.*")
file(GLOB MCONFIG_SRCS "orc8r/protos/mconfig/*.pb.*")
file(GLOB PROMETHEUS_SRCS "orc8r/protos/prometheus/*.pb.*")

include_directories(".")

add_library(MAGMA_PROTOBUF
    ${PROTO_SRCS}
    ${MCONFIG_SRCS}
    ${PROMETHEUS_SRCS}
    )
target_link_libraries(MAGMA_PROTOBUF
    protobuf grpc++ grpc dl prometheus-cpp
    )
target_include_directories(MAGMA_PROTOBUF SYSTEM PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    )
