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

#include <stdint.h>

#include "log.h"
#include "TLVEncoder.h"
#include "TLVDecoder.h"
#include "TrackingAreaIdentityList.h"
#include "common_defs.h"

//------------------------------------------------------------------------------
int decode_tracking_area_identity_list(
    tai_list_t* trackingareaidentitylist, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded   = 0;
  uint8_t ielen = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  int partial_item = 0;
  while ((decoded < ielen) &&
         (partial_item < TRACKING_AREA_IDENTITY_LIST_MAXIMUM_NUM_TAI)) {
    trackingareaidentitylist->partial_tai_list[partial_item].typeoflist =
        (*(buffer + decoded) >> 5) & 0x3;
    trackingareaidentitylist->partial_tai_list[partial_item].numberofelements =
        *(buffer + decoded) & 0x1f;
    if (TRACKING_AREA_IDENTITY_LIST_MAXIMUM_NUM_TAI <=
        ((trackingareaidentitylist->partial_tai_list[partial_item]
              .numberofelements +
          1))) {
      return TLV_VALUE_DOESNT_MATCH;
    }
    decoded++;
    if (TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_CONSECUTIVE_TACS ==
        trackingareaidentitylist->partial_tai_list[partial_item].typeoflist) {
      trackingareaidentitylist->partial_tai_list[partial_item]
          .u.tai_one_plmn_consecutive_tacs.plmn.mcc_digit2 =
          (*(buffer + decoded) >> 4) & 0xf;
      trackingareaidentitylist->partial_tai_list[partial_item]
          .u.tai_one_plmn_consecutive_tacs.plmn.mcc_digit1 =
          *(buffer + decoded) & 0xf;
      decoded++;
      trackingareaidentitylist->partial_tai_list[partial_item]
          .u.tai_one_plmn_consecutive_tacs.plmn.mnc_digit3 =
          (*(buffer + decoded) >> 4) & 0xf;
      trackingareaidentitylist->partial_tai_list[partial_item]
          .u.tai_one_plmn_consecutive_tacs.plmn.mcc_digit3 =
          *(buffer + decoded) & 0xf;
      decoded++;
      trackingareaidentitylist->partial_tai_list[partial_item]
          .u.tai_one_plmn_consecutive_tacs.plmn.mnc_digit2 =
          (*(buffer + decoded) >> 4) & 0xf;
      trackingareaidentitylist->partial_tai_list[partial_item]
          .u.tai_one_plmn_consecutive_tacs.plmn.mnc_digit1 =
          *(buffer + decoded) & 0xf;
      decoded++;
      IES_DECODE_U16(
          buffer, decoded,
          trackingareaidentitylist->partial_tai_list[partial_item]
              .u.tai_one_plmn_consecutive_tacs.tac);
    } else if (
        TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_NON_CONSECUTIVE_TACS ==
        trackingareaidentitylist->partial_tai_list[partial_item].typeoflist) {
      int i;
      trackingareaidentitylist->partial_tai_list[partial_item]
          .u.tai_one_plmn_non_consecutive_tacs.plmn.mcc_digit2 =
          (*(buffer + decoded) >> 4) & 0xf;
      trackingareaidentitylist->partial_tai_list[partial_item]
          .u.tai_one_plmn_non_consecutive_tacs.plmn.mcc_digit1 =
          *(buffer + decoded) & 0xf;
      decoded++;
      trackingareaidentitylist->partial_tai_list[partial_item]
          .u.tai_one_plmn_non_consecutive_tacs.plmn.mnc_digit3 =
          (*(buffer + decoded) >> 4) & 0xf;
      trackingareaidentitylist->partial_tai_list[partial_item]
          .u.tai_one_plmn_non_consecutive_tacs.plmn.mcc_digit3 =
          *(buffer + decoded) & 0xf;
      decoded++;
      trackingareaidentitylist->partial_tai_list[partial_item]
          .u.tai_one_plmn_non_consecutive_tacs.plmn.mnc_digit2 =
          (*(buffer + decoded) >> 4) & 0xf;
      trackingareaidentitylist->partial_tai_list[partial_item]
          .u.tai_one_plmn_non_consecutive_tacs.plmn.mnc_digit1 =
          *(buffer + decoded) & 0xf;
      decoded++;
      for (i = 0; i <= trackingareaidentitylist->partial_tai_list[partial_item]
                               .numberofelements +
                           1;
           i++) {
        IES_DECODE_U16(
            buffer, decoded,
            trackingareaidentitylist->partial_tai_list[partial_item]
                .u.tai_one_plmn_non_consecutive_tacs.tac[i]);
      }
    } else if (
        TRACKING_AREA_IDENTITY_LIST_MANY_PLMNS ==
        trackingareaidentitylist->partial_tai_list[partial_item].typeoflist) {
      int i;
      for (i = 0; i <= trackingareaidentitylist->partial_tai_list[partial_item]
                               .numberofelements +
                           1;
           i++) {
        trackingareaidentitylist->partial_tai_list[partial_item]
            .u.tai_many_plmn[i]
            .plmn.mcc_digit2 = (*(buffer + decoded) >> 4) & 0xf;
        trackingareaidentitylist->partial_tai_list[partial_item]
            .u.tai_many_plmn[i]
            .plmn.mcc_digit1 = *(buffer + decoded) & 0xf;
        decoded++;
        trackingareaidentitylist->partial_tai_list[partial_item]
            .u.tai_many_plmn[i]
            .plmn.mnc_digit3 = (*(buffer + decoded) >> 4) & 0xf;
        trackingareaidentitylist->partial_tai_list[partial_item]
            .u.tai_many_plmn[i]
            .plmn.mcc_digit3 = *(buffer + decoded) & 0xf;
        decoded++;
        trackingareaidentitylist->partial_tai_list[partial_item]
            .u.tai_many_plmn[i]
            .plmn.mnc_digit2 = (*(buffer + decoded) >> 4) & 0xf;
        trackingareaidentitylist->partial_tai_list[partial_item]
            .u.tai_many_plmn[i]
            .plmn.mnc_digit1 = *(buffer + decoded) & 0xf;
        decoded++;
        IES_DECODE_U16(
            buffer, decoded,
            trackingareaidentitylist->partial_tai_list[partial_item]
                .u.tai_many_plmn[i]
                .tac);
      }
    } else {
      OAILOG_DEBUG(
          LOG_NAS, "Type of TAIL list not handled %d",
          trackingareaidentitylist->partial_tai_list[partial_item].typeoflist);
      return TLV_VALUE_DOESNT_MATCH;
    }
    partial_item++;
  }
  return decoded;
}

//------------------------------------------------------------------------------
int encode_tracking_area_identity_list(
    tai_list_t* trackingareaidentitylist, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, TRACKING_AREA_IDENTITY_LIST_MINIMUM_LENGTH, len);

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;

  int partial_item = 0;
  while ((encoded < len) &&
         (partial_item < trackingareaidentitylist->numberoflists)) {
    *(buffer + encoded) =
        0x00 |
        ((trackingareaidentitylist->partial_tai_list[partial_item].typeoflist &
          0x3)
         << 5) |
        (trackingareaidentitylist->partial_tai_list[partial_item]
             .numberofelements &
         0x1f);
    encoded++;

    if (TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_CONSECUTIVE_TACS ==
        trackingareaidentitylist->partial_tai_list[partial_item].typeoflist) {
      *(buffer + encoded) =
          0x00 |
          ((trackingareaidentitylist->partial_tai_list[partial_item]
                .u.tai_one_plmn_consecutive_tacs.plmn.mcc_digit2 &
            0xf)
           << 4) |
          (trackingareaidentitylist->partial_tai_list[partial_item]
               .u.tai_one_plmn_consecutive_tacs.plmn.mcc_digit1 &
           0xf);
      encoded++;
      *(buffer + encoded) =
          0x00 |
          ((trackingareaidentitylist->partial_tai_list[partial_item]
                .u.tai_one_plmn_consecutive_tacs.plmn.mnc_digit3 &
            0xf)
           << 4) |
          (trackingareaidentitylist->partial_tai_list[partial_item]
               .u.tai_one_plmn_consecutive_tacs.plmn.mcc_digit3 &
           0xf);
      encoded++;
      *(buffer + encoded) =
          0x00 |
          ((trackingareaidentitylist->partial_tai_list[partial_item]
                .u.tai_one_plmn_consecutive_tacs.plmn.mnc_digit2 &
            0xf)
           << 4) |
          (trackingareaidentitylist->partial_tai_list[partial_item]
               .u.tai_one_plmn_consecutive_tacs.plmn.mnc_digit1 &
           0xf);
      encoded++;
      IES_ENCODE_U16(
          buffer, encoded,
          trackingareaidentitylist->partial_tai_list[partial_item]
              .u.tai_one_plmn_consecutive_tacs.tac);
    } else if (
        TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_NON_CONSECUTIVE_TACS ==
        trackingareaidentitylist->partial_tai_list[partial_item].typeoflist) {
      int i;
      *(buffer + encoded) =
          0x00 |
          ((trackingareaidentitylist->partial_tai_list[partial_item]
                .u.tai_one_plmn_non_consecutive_tacs.plmn.mcc_digit2 &
            0xf)
           << 4) |
          (trackingareaidentitylist->partial_tai_list[partial_item]
               .u.tai_one_plmn_non_consecutive_tacs.plmn.mcc_digit1 &
           0xf);
      encoded++;
      *(buffer + encoded) =
          0x00 |
          ((trackingareaidentitylist->partial_tai_list[partial_item]
                .u.tai_one_plmn_non_consecutive_tacs.plmn.mnc_digit3 &
            0xf)
           << 4) |
          (trackingareaidentitylist->partial_tai_list[partial_item]
               .u.tai_one_plmn_non_consecutive_tacs.plmn.mcc_digit3 &
           0xf);
      encoded++;
      *(buffer + encoded) =
          0x00 |
          ((trackingareaidentitylist->partial_tai_list[partial_item]
                .u.tai_one_plmn_non_consecutive_tacs.plmn.mnc_digit2 &
            0xf)
           << 4) |
          (trackingareaidentitylist->partial_tai_list[partial_item]
               .u.tai_one_plmn_non_consecutive_tacs.plmn.mnc_digit1 &
           0xf);
      encoded++;
      for (i = 0; i <= trackingareaidentitylist->partial_tai_list[partial_item]
                           .numberofelements;
           i++) {
        IES_ENCODE_U16(
            buffer, encoded,
            trackingareaidentitylist->partial_tai_list[partial_item]
                .u.tai_one_plmn_non_consecutive_tacs.tac[i]);
      }
    } else if (
        TRACKING_AREA_IDENTITY_LIST_MANY_PLMNS ==
        trackingareaidentitylist->partial_tai_list[partial_item].typeoflist) {
      int i;
      for (i = 0; i <= trackingareaidentitylist->partial_tai_list[partial_item]
                           .numberofelements;
           i++) {
        *(buffer + encoded) =
            0x00 |
            ((trackingareaidentitylist->partial_tai_list[partial_item]
                  .u.tai_many_plmn[i]
                  .plmn.mcc_digit2 &
              0xf)
             << 4) |
            (trackingareaidentitylist->partial_tai_list[partial_item]
                 .u.tai_many_plmn[i]
                 .plmn.mcc_digit1 &
             0xf);
        encoded++;
        *(buffer + encoded) =
            0x00 |
            ((trackingareaidentitylist->partial_tai_list[partial_item]
                  .u.tai_many_plmn[i]
                  .plmn.mnc_digit3 &
              0xf)
             << 4) |
            (trackingareaidentitylist->partial_tai_list[partial_item]
                 .u.tai_many_plmn[i]
                 .plmn.mcc_digit3 &
             0xf);
        encoded++;
        *(buffer + encoded) =
            0x00 |
            ((trackingareaidentitylist->partial_tai_list[partial_item]
                  .u.tai_many_plmn[i]
                  .plmn.mnc_digit2 &
              0xf)
             << 4) |
            (trackingareaidentitylist->partial_tai_list[partial_item]
                 .u.tai_many_plmn[i]
                 .plmn.mnc_digit1 &
             0xf);
        encoded++;
        IES_ENCODE_U16(
            buffer, encoded,
            trackingareaidentitylist->partial_tai_list[partial_item]
                .u.tai_many_plmn[i]
                .tac);
      }
    } else {
      OAILOG_DEBUG(
          LOG_NAS, "Type of TAIL list not handled %d",
          trackingareaidentitylist->partial_tai_list[partial_item].typeoflist);
      return TLV_VALUE_DOESNT_MATCH;
    }
    partial_item++;
  }
  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);
  return encoded;
}
