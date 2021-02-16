# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

###############################################################################
# Security library
###############################################################################

add_library(LIB_SECU
    kdf.c
    key_nas_deriver.c
    key_nas_encryption.c
    nas_stream_eea1.c
    nas_stream_eea2.c
    nas_stream_eia1.c
    nas_stream_eia2.c
    rijndael.c
    snow3g.c
)
target_link_libraries(LIB_SECU
    ${NETTLE_LIBRARIES}
    COMMON
    LIB_BSTR
)
target_include_directories(LIB_SECU PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    )
