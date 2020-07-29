/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*! \file sgs_service_request.c
   \brief Sends Service Request message to MSC/VLR
   \author
   \version
   \company
   \email:
*/

#include <stdio.h>
#include <stdint.h>

#include "intertask_interface.h"
#include "log.h"
#include "sgs_defs.h"
#include "assertions.h"
#include "sgs_messages_types.h"
#include "sgs_messages.h"
#include "conversions.h"

int sgs_send_service_request(
    itti_sgsap_service_request_t* const sgs_service_request_p) {
  imsi64_t imsi64 = INVALID_IMSI64;
  IMSI_TO_IMSI64(&(sgs_service_request_p->imsi), imsi64);
  OAILOG_DEBUG(
      LOG_SGS,
      "Received SGS_SERVICE_REQUEST from task MME_APP for (IMSI = " IMSI_64_FMT
      ") \n",
      imsi64);
  return 0;
}
