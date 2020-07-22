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

/*! \file s11_common.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdbool.h>
#include <stdint.h>

#include "bstrlib.h"

#include "NwGtpv2c.h"
#include "s11_common.h"
#include "log.h"

nw_rc_t s11_ie_indication_generic(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg) {
  OAILOG_DEBUG(
      LOG_S11,
      "Received IE Parse Indication for of type %u, length %u, "
      "instance %u!\n",
      ieType, ieLength, ieInstance);
  return NW_OK;
}
