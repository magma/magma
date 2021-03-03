/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
#include <stdio.h>
#include <stdint.h>
#include "log.h"
#include "common_defs.h"
#include "intertask_interface.h"

int sgw_s8_handle_s11_create_session_request(
    const itti_s11_create_session_request_t* const session_req_pP,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  OAILOG_INFO_UE(
      LOG_SGW_S8, imsi64, "Received S11 CREATE SESSION REQUEST from MME_APP\n");
  OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNok);
}
