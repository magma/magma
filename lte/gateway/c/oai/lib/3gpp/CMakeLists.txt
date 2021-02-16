# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

add_library(LIB_3GPP
    3gpp_24.008_cc_ies.c
    3gpp_24.008_common_ies.c
    3gpp_24.008_gmm_ies.c
    3gpp_24.008_gprs_common_ies.c
    3gpp_24.008_mm_ies.c
    3gpp_24.008_sm_ies.c
)
target_link_libraries(LIB_3GPP
    LIB_BSTR LIB_S1AP
)
target_include_directories(LIB_3GPP PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
)
