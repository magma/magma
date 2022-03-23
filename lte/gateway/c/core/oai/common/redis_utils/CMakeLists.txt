# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

add_library(redis_utils redis_client.cpp)

file(GLOB ABSL_LIBS "/usr/local/lib/libabsl*.a")

target_link_libraries(redis_utils MAGMA_CONFIG COMMON cpp_redis tacopie 
    ssl crypto z gpr
    ${ABSL_LIBS}
    ${ABSL_LIBS}
    absl_debugging_internal
    absl_demangle_internal
    upb
    re2
    cares
    address_sorting
	/usr/local/lib/libprotobuf.a
    )


target_include_directories(redis_utils PUBLIC
    $ENV{MAGMA_ROOT}
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${MAGMA_COMMON_DIR}
    )
