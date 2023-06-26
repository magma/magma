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

/*! \file s6a_error.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <string.h>
#include <stdint.h>
#include <errno.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/tasks/s6a/s6a_defs.hpp"

struct avp;

int s6a_parse_experimental_result(struct avp* avp,
                                  s6a_experimental_result_t* ptr) {
  struct avp_hdr* hdr;
  struct avp* child_avp = NULL;

  if (!avp) {
    return (status_code_e)EINVAL;
  }

  CHECK_FCT(fd_msg_avp_hdr(avp, &hdr));
  DevAssert(hdr->avp_code == AVP_CODE_EXPERIMENTAL_RESULT);
  CHECK_FCT(fd_msg_browse_internal(avp, MSG_BRW_FIRST_CHILD,
                                   reinterpret_cast<msg_or_avp**>(&child_avp),
                                   NULL));

  while (child_avp) {
    CHECK_FCT(fd_msg_avp_hdr(child_avp, &hdr));

    switch (hdr->avp_code) {
      case AVP_CODE_EXPERIMENTAL_RESULT_CODE:
        OAILOG_ERROR(LOG_S6A, "Got experimental error %u:%s\n",
                     hdr->avp_value->u32,
                     experimental_retcode_2_string(hdr->avp_value->u32));

        if (ptr) {
          *ptr = (s6a_experimental_result_t)hdr->avp_value->u32;
        }

        break;

      case AVP_CODE_VENDOR_ID:
        DevCheck(hdr->avp_value->u32 == 10415, hdr->avp_value->u32,
                 AVP_CODE_VENDOR_ID, 10415);
        break;

      default:
        return RETURNerror;
    }

    /*
     * Go to next AVP in the grouped AVP
     */
    CHECK_FCT(fd_msg_browse_internal(child_avp, MSG_BRW_NEXT,
                                     reinterpret_cast<msg_or_avp**>(&child_avp),
                                     NULL));
  }

  return RETURNok;
}

char* experimental_retcode_2_string(uint32_t ret_code) {
  switch (ret_code) {
      /*
       * Experimental-Result-Codes
       */
    case DIAMETER_ERROR_USER_UNKNOWN:
      return const_cast<char*>("DIAMETER_ERROR_USER_UNKNOWN");

    case DIAMETER_ERROR_ROAMING_NOT_ALLOWED:
      return const_cast<char*>("DIAMETER_ERROR_ROAMING_NOT_ALLOWED");

    case DIAMETER_ERROR_UNKNOWN_EPS_SUBSCRIPTION:
      return const_cast<char*>("DIAMETER_ERROR_UNKNOWN_EPS_SUBSCRIPTION");

    case DIAMETER_ERROR_RAT_NOT_ALLOWED:
      return const_cast<char*>("DIAMETER_ERROR_RAT_NOT_ALLOWED");

    case DIAMETER_ERROR_EQUIPMENT_UNKNOWN:
      return const_cast<char*>("DIAMETER_ERROR_EQUIPMENT_UNKNOWN");

    case DIAMETER_ERROR_UNKNOWN_SERVING_NODE:
      return const_cast<char*>("DIAMETER_ERROR_UNKNOWN_SERVING_NODE");

    case DIAMETER_AUTHENTICATION_DATA_UNAVAILABLE:
      return const_cast<char*>("DIAMETER_AUTHENTICATION_DATA_UNAVAILABLE");

    default:
      break;
  }

  return const_cast<char*>("DIAMETER_AVP_UNSUPPORTED");
}

char* retcode_2_string(uint32_t ret_code) {
  switch (ret_code) {
    case ER_DIAMETER_SUCCESS:
      return const_cast<char*>("DIAMETER_SUCCESS");

    case ER_DIAMETER_MISSING_AVP:
      return const_cast<char*>("DIAMETER_MISSING_AVP");

    case ER_DIAMETER_INVALID_AVP_VALUE:
      return const_cast<char*>("DIAMETER_INVALID_AVP_VALUE");

    case ER_DIAMETER_AUTHORIZATION_REJECTED:
      return const_cast<char*>("DIAMETER_AUTHORIZATION_REJECTED");

    case ER_DIAMETER_COMMAND_UNSUPPORTED:
      return const_cast<char*>("DIAMETER_COMMAND_UNSUPPORTED");

    case ER_DIAMETER_UNABLE_TO_DELIVER:
      return const_cast<char*>("DIAMETER_UNABLE_TO_DELIVER");

    case ER_DIAMETER_UNKNOWN_PEER:
      return const_cast<char*>("DIAMETER_UNKNOWN_PEER");

    default:
      break;
  }

  return const_cast<char*>("DIAMETER_AVP_UNSUPPORTED");
}
