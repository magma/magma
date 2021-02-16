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

add_library(TASK_SMS_ORC8R
    sms_orc8r_task.c
    sms_orc8r_service_handler.c
    )

target_link_libraries(TASK_SMS_ORC8R
    COMMON
    LIB_BSTR LIB_HASHTABLE LIB_SMS_ORC8R_CLIENT
    ${PROTOBUF_LIBRARIES}
)
target_include_directories(TASK_SMS_ORC8R PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
)
