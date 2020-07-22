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

/*! \file oai_mme_log.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <freeDiameter/freeDiameter-host.h>
#include <freeDiameter/libfdcore.h>

#include "log.h"
#include "oai_mme.h"

// TODO: (amar) unused function check with OAI.
int oai_mme_log_specific(int log_level) {
  if (log_level == 1) {
    asn_debug      = 0;
    asn1_xer_print = 1;
    fd_g_debug_lvl = INFO;
  } else if (log_level == 2) {
    asn_debug      = 1;
    asn1_xer_print = 1;
    fd_g_debug_lvl = ANNOYING;
  } else {
    asn1_xer_print = 0;
    asn_debug      = 0;
    fd_g_debug_lvl = NONE;
  }

  return 0;
}
