# Copyright 2021 The Magma Authors.
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

add_executable(ngap_test
        ngap_test.cpp
        util_ngap_setup_failure.cpp
        util_ngap_initiate_ue.cpp
        util_ngap_amf_nas_procedures.cpp
        test_ngap_encode_decode.cpp
        test_ngap_handle_new_association.cpp
        test_ngap_flows.cpp
	test_ngap_state_converter.cpp)

target_link_libraries(ngap_test
        TASK_NGAP TASK_AMF_APP LIB_BSTR
        gtest gtest_main
        )

add_test(test_ngap ngap_test)
