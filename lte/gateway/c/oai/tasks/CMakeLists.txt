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

add_subdirectory(mme_app) # TASK_MME_APP
add_subdirectory(nas) # TASK_NAS
add_subdirectory(s1ap) # TASK_S1AP, LIB_S1AP
add_subdirectory(s6a) # TASK_S6A
add_subdirectory(sctp) # TASK_SCTP_SERVER
add_subdirectory(service303) # TASK_SERVICE303
add_subdirectory(sgs) # TASK_SGS
add_subdirectory(sms_orc8r) # TASK_SMS_ORC8R
add_subdirectory(grpc_service) # TASK_GRPC_SERVICE
add_subdirectory(ha) # TASK_HA

if (EMBEDDED_SGW)
add_subdirectory(gtpv1-u) # TASK_GTPV1U
add_subdirectory(sgw) # TASK_SGW
else (EMBEDDED_SGW)
add_subdirectory(s11) # TASK_S11
add_subdirectory(udp) #TASK_UDP
endif (EMBEDDED_SGW)
