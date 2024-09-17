# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

PROJECT(LteGatewayCCommon)

###############################################################################
# Core Network Utils
###############################################################################

set(LTE_GATEWAY_C_COMMON_SRC
    assertions.h
    backtrace.c
    common_defs.h
    dynamic_memory_check.cpp
    dynamic_memory_check.h
    )

add_library(LTE_GATEWAY_C_COMMON ${LTE_GATEWAY_C_COMMON_SRC})

# Must include MAGMA_ROOT so that full path includes will succeed.
target_include_directories(LTE_GATEWAY_C_COMMON PUBLIC
    $ENV{MAGMA_ROOT}
    )

set(CMAKE_EXPORT_PACKAGE_REGISTRY ON)
export(TARGETS LTE_GATEWAY_C_COMMON FILE LTE_GATEWAY_C_COMMONConfig.cmake)
export(PACKAGE LTE_GATEWAY_C_COMMON)
