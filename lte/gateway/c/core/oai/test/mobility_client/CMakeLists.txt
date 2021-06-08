# Copyright 2020 The Magma Authors.
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

include_directories("/usr/src/googletest/googlemock/include/")

add_executable(mobility_client_test test_mobility_client.cpp)

target_link_libraries(mobility_client_test
    LIB_MOBILITY_CLIENT protobuf grpc++ dl stdc++ m
    )

# TODO add support for integration tests
# add_test(test_rpc_client_integration rpc_client_test)
