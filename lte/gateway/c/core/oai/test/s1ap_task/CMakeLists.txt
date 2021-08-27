# Copyright 2021 The Magma Authors.
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

add_executable(s1ap_test
        s1ap_test.cpp
        test_s1ap_handle_new_association.cpp
        test_s1ap_state_manager.cpp)

target_link_libraries(s1ap_test
        TASK_S1AP
        gtest gtest_main
        )

add_test(test_s1ap s1ap_test)
