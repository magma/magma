# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

set(NWGTPV2C_DIR  ${PROJECT_SOURCE_DIR}/lib/gtpv2-c/nwgtpv2c-0.11/src)
set(NWGTPV2C_IE_FORMATTER_DIR ${PROJECT_SOURCE_DIR}/lib/gtpv2-c/gtpv2c_ie_formatter/src)

add_library(LIB_GTPV2C
    ${NWGTPV2C_DIR}/NwGtpv2cTrxn.c
    ${NWGTPV2C_DIR}/NwGtpv2cTunnel.c
    ${NWGTPV2C_DIR}/NwGtpv2cMsg.c
    ${NWGTPV2C_DIR}/NwGtpv2cMsgIeParseInfo.c
    ${NWGTPV2C_DIR}/NwGtpv2cMsgParser.c
    ${NWGTPV2C_DIR}/NwGtpv2c.c
    ${NWGTPV2C_IE_FORMATTER_DIR}/gtpv2c_ie_formatter.c
)
target_link_libraries(LIB_GTPV2C
    COMMON
    LIB_BSTR
 )

target_include_directories(LIB_GTPV2C PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}/nwgtpv2c-0.11/include/
    ${CMAKE_CURRENT_SOURCE_DIR}/nwgtpv2c-0.11/shared/
    ${CMAKE_CURRENT_SOURCE_DIR}/nwgtpv2c-0.11/src/
    ${PROJECT_SOURCE_DIR}/tasks/nas/ies/
    ${CMAKE_CURRENT_SOURCE_DIR}/
    ${CMAKE_CURRENT_SOURCE_DIR}/gtpv2c_ie_formatter/shared/
)
