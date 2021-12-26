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
 
include_directories("/usr/src/googletest/googlemock/include/")
link_directories(/usr/src/googletest/googlemock/lib/)

add_library(MOCK_TASKS
    ha_mock.cpp
    s1ap_mock.cpp
    s6a_mock.cpp
    s11_mock.cpp
    service303_mock.cpp
    sgs_mock.cpp
    sgw_s8_mock.cpp
    spgw_mock.cpp
    sms_orc8r_mock.cpp
    mme_app_mock.cpp
    sctp_mock.cpp
    grpc_mock.cpp
)

target_link_libraries(MOCK_TASKS
    LIB_ITTI gmock_main gtest gtest_main gmock
    TASK_GRPC_SERVICE
)
