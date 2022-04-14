# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PROJECT(MagmaEbpf)
cmake_minimum_required(VERSION 3.7.2)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

add_library(MAGMA_EBPF INTERFACE)
# Specify include path for chained dependencies
target_include_directories(MAGMA_EBPF INTERFACE $ENV{MAGMA_ROOT} ${CMAKE_CURRENT_SOURCE_DIR})

install(TARGETS MAGMA_EBPF EXPORT MAGMA_EBPF
    INCLUDES DESTINATION ""
    ARCHIVE DESTINATION lib)

# Export package
set(CMAKE_EXPORT_PACKAGE_REGISTRY ON)
export(TARGETS MAGMA_EBPF FILE MAGMA_EBPFConfig.cmake)
export(PACKAGE MAGMA_EBPF)
