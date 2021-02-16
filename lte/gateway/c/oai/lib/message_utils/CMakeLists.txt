# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

add_library(LIB_MESSAGE_UTILS
    bytes_to_ie.c
    ie_to_bytes.c
    service303_message_utils.c
)
target_include_directories(LIB_MESSAGE_UTILS
    PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    PRIVATE
    ${LIB_DIR}/itti
)
target_link_libraries(LIB_MESSAGE_UTILS
    COMMON
)
