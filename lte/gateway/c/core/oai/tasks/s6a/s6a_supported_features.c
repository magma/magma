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

/*! \file s6a_supported_features.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdint.h>
#include <netinet/in.h>
#include <stdio.h>
#include <string.h>

#include "assertions.h"
#include "3gpp_24.008.h"
#include "common_defs.h"
#include "common_types.h"
#include "s6a_defs.h"

struct avp;

int s6a_parse_supported_features(
    struct avp* avp_supported_features,
    supported_features_t* subscription_data) {
  struct avp* avp = NULL;
  struct avp_hdr* hdr;
  uint32_t feature_list_id = 0;
  uint32_t feature_list    = 0;

  CHECK_FCT(
      fd_msg_browse(avp_supported_features, MSG_BRW_FIRST_CHILD, &avp, NULL));

  while (avp) {
    hdr = NULL;
    CHECK_FCT(fd_msg_avp_hdr(avp, &hdr));

    if (hdr) {
      switch (hdr->avp_code) {
        case AVP_CODE_VENDOR_ID:
          // no check
          break;

        case AVP_CODE_FEATURE_LIST_ID:
          if (hdr->avp_value) {
            feature_list_id = hdr->avp_value->u32;
          } else {
            return RETURNerror;
          }
          break;

        case AVP_CODE_FEATURE_LIST:
          if (hdr->avp_value) {
            feature_list = hdr->avp_value->u32;
            if (feature_list_id == 2) {
              subscription_data->external_identifier = true;
              subscription_data->nr_as_secondary_rat =
                  FLAG_IS_SET(feature_list, FLID_NR_AS_SECONDARY_RAT);
            }
          } else {
            return RETURNerror;
          }
          break;

        default:
          OAILOG_DEBUG(
              LOG_S6A, "Unknown AVP code %d not processed\n", hdr->avp_code);
          return RETURNerror;
      }
    } else {
      OAILOG_DEBUG(LOG_S6A, "Supported Features parsing Error\n");
      return RETURNerror;
    }

    /*
     * Go to next AVP in the grouped AVP
     */
    CHECK_FCT(fd_msg_browse(avp, MSG_BRW_NEXT, &avp, NULL));
  }

  return RETURNok;
}
