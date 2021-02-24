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

#ifndef TRACKING_AREA_IDENTITY_LIST_SEEN
#define TRACKING_AREA_IDENTITY_LIST_SEEN

#include <stdint.h>

#include "TrackingAreaIdentity.h"

#define TRACKING_AREA_IDENTITY_LIST_MINIMUM_LENGTH 8
#define TRACKING_AREA_IDENTITY_LIST_MAXIMUM_LENGTH 98

#define TRACKING_AREA_IDENTITY_LIST_MAXIMUM_NUM_TAI 16

typedef struct partial_tai_list_s {
  /* XXX - The only supported type of list is a list of TACs
   * belonging to one PLMN, with consecutive TAC values */
#define TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_NON_CONSECUTIVE_TACS 0b00
#define TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_CONSECUTIVE_TACS 0b01
#define TRACKING_AREA_IDENTITY_LIST_MANY_PLMNS 0b10
  uint8_t typeoflist;
  uint8_t numberofelements;
  // TODO not optimized
  union {
    tai_t tai_many_plmn[TRACKING_AREA_IDENTITY_LIST_MAXIMUM_NUM_TAI];
    tai_t tai_one_plmn_consecutive_tacs;
    struct {
      plmn_t plmn;
      tac_t tac[TRACKING_AREA_IDENTITY_LIST_MAXIMUM_NUM_TAI];
    } tai_one_plmn_non_consecutive_tacs;
  } u;
} partial_tai_list_t;

typedef struct tai_list_s {
  uint8_t numberoflists;
  partial_tai_list_t
      partial_tai_list[TRACKING_AREA_IDENTITY_LIST_MAXIMUM_NUM_TAI];
} tai_list_t;

int encode_tracking_area_identity_list(
    tai_list_t* trackingareaidentitylist, uint8_t iei, uint8_t* buffer,
    uint32_t len);

int decode_tracking_area_identity_list(
    tai_list_t* trackingareaidentitylist, uint8_t iei, uint8_t* buffer,
    uint32_t len);

#endif /* TRACKING AREA IDENTITY LIST_SEEN */
