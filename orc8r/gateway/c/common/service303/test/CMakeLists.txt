# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)
PROJECT(MagmaCommonTests)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

include_directories("/usr/src/googletest/googlemock/include/")
link_directories("/usr/src/googletest/googlemock/lib/")

file(GLOB ABSL_LIBS "/usr/local/lib/libabsl*.a")

foreach (service303_test magma_service service303 metrics)
  add_executable(${service303_test}_test test_${service303_test}.cpp)
  target_link_libraries(${service303_test}_test
    -Wl,--start-group
    SERVICE303_LIB
    SERVICE_REGISTRY
    gmock_main gtest gtest_main gmock
    yaml-cpp pthread rt
    ${GCOV_LIB}
    ${ABSL_LIBS}
    ${ABSL_LIBS}
    absl_debugging_internal
    absl_demangle_internal
    z
    gpr
    upb
    re2
    cares
    ssl
    crypto
    address_sorting
    /usr/local/lib/libprotobuf.a
    -Wl,--end-group 
  )
  add_test(test_${service303_test} ${service303_test}_test)
endforeach (service303_test)
