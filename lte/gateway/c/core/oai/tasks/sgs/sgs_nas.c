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

/*! \file sgs_nas.c
  \brief
  \author
  \company
  \email:
*/
#define SGG
#define SGS_NAS_C

#include "log.h"
#include "intertask_interface.h"
#include "sgs_messages.h"
#include "assertions.h"

int sgs_send_uplink_unitdata(
    itti_sgsap_uplink_unitdata_t* sgs_uplink_unitdata_p) {
  DevAssert(sgs_uplink_unitdata_p);

  OAILOG_DEBUG(
      LOG_SGS, "Received SGS_UPLINK_UNITDATA from task NAS for IMSI : %s \n",
      sgs_uplink_unitdata_p->imsi);
  /* TODO: Add the code for sending this message to sgs service*/

  return 0;
}
