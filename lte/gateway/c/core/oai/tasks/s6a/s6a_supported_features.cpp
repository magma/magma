/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/****************************************************************************
  Subsystem   MME
  Description Handles the Supported Features Diameter AVP
*****************************************************************************/

#include <stdint.h>
#include <netinet/in.h>
#include <stdio.h>
#include <string.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/tasks/s6a/s6a_defs.hpp"

struct avp;

int s6a_parse_supported_features(struct avp* avp_supported_features,
                                 supported_features_t* subscription_data) {
  struct avp* avp = NULL;
  struct avp_hdr* hdr;
  uint32_t feature_list_id = 0;
  uint32_t feature_list = 0;

  CHECK_FCT(fd_msg_browse_internal(avp_supported_features, MSG_BRW_FIRST_CHILD,
                                   reinterpret_cast<msg_or_avp**>(&avp), NULL));

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
          OAILOG_DEBUG(LOG_S6A, "Unknown AVP code %d not processed\n",
                       hdr->avp_code);
          return RETURNerror;
      }
    } else {
      OAILOG_DEBUG(LOG_S6A, "Supported Features parsing Error\n");
      return RETURNerror;
    }

    /*
     * Go to next AVP in the grouped AVP
     */
    CHECK_FCT(fd_msg_browse_internal(
        avp, MSG_BRW_NEXT, reinterpret_cast<msg_or_avp**>(&avp), NULL));
  }

  return RETURNok;
}
