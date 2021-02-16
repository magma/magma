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

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

set (CMAKE_CXX_FLAGS "-Wno-write-strings -Wno-literal-suffix")
add_definitions(-DDEBUG_IS_ON=1)

add_library(LIB_OPENFLOW_CONTROLLER
  ControllerMain.cpp
  OpenflowController.cpp
  PagingApplication.cpp
  ControllerEvents.cpp
  BaseApplication.cpp
  OpenflowMessenger.cpp
  GTPApplication.cpp
  IMSIEncoder.cpp
  )
target_link_libraries(LIB_OPENFLOW_CONTROLLER
  COMMON
  fluid_base fluid_msg
  LIB_BSTR
  TASK_SGW
)
target_include_directories(LIB_OPENFLOW_CONTROLLER PUBLIC
  ${CMAKE_CURRENT_SOURCE_DIR}
)
