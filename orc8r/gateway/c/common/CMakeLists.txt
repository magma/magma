# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

PROJECT(MagmaCommon)

include(CMakeProtoMacros.txt)

ADD_SUBDIRECTORY(service303)
ADD_SUBDIRECTORY(service_registry)
ADD_SUBDIRECTORY(async_grpc)
ADD_SUBDIRECTORY(config)
ADD_SUBDIRECTORY(datastore)
ADD_SUBDIRECTORY(policydb)
ADD_SUBDIRECTORY(logging)
ADD_SUBDIRECTORY(eventd)

if (BUILD_TESTS)
  ENABLE_TESTING()
  ADD_SUBDIRECTORY(test)
endif (BUILD_TESTS)
