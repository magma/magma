# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


add_compile_options(-std=c++14)
set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

include_directories("${PROJECT_SOURCE_DIR}/../common/config")
include_directories("${PROJECT_SOURCE_DIR}/../common/service303")

include_directories("${PROJECT_SOURCE_DIR}/../common/protobuf")

add_library(COMMON_TEST_LIB)

target_link_libraries(COMMON_TEST_LIB
  CONFIG
  gmock_main
  yaml-cpp
  pthread
  rt
  ${GCOV_LIB}
  SERVICE303_LIB
)

foreach(common_test yaml_utils magma_service service_registry service303 metrics)
  add_executable(${common_test}_test test_${common_test}.cpp)
  target_link_libraries(${common_test}_test COMMON_TEST_LIB)
  add_test(test_${common_test} ${common_test}_test)
endforeach(common_test)
