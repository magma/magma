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

/*! \file mme_app_ue_context.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdio.h>
#include <string.h>
#include <stdbool.h>
#include <stdint.h>
#include <inttypes.h>

#include "conversions.h"
#include "mme_app_ue_context.h"
#include "3gpp_23.003.h"
#include "3gpp_36.401.h"

/**
 * @brief mme_app_copy_imsi: copies an mme imsi to another mme imsi
 * @param imsi_dst
 * @param imsi_src
 */

void mme_app_copy_imsi(
    mme_app_imsi_t* imsi_dst, const mme_app_imsi_t* imsi_src) {
  strncpy(imsi_dst->data, imsi_src->data, IMSI_BCD_DIGITS_MAX + 1);
  imsi_dst->length = imsi_src->length;
}

/**
 * @brief mme_app_imsi_compare: compares to imsis returns true if the same else
 * false
 * @param imsi_a
 * @param imsi_b
 * @return
 */

bool mme_app_imsi_compare(
    mme_app_imsi_t const* imsi_a, mme_app_imsi_t const* imsi_b) {
  if ((strncmp(imsi_a->data, imsi_b->data, IMSI_BCD_DIGITS_MAX) == 0) &&
      imsi_a->length == imsi_b->length) {
    return true;
  } else
    return false;
}

/**
 * @brief mme_app_string_to_imsi converst the a string to the imsi mme structure
 * @param imsi_dst
 * @param imsi_string_src
 */

void mme_app_string_to_imsi(
    mme_app_imsi_t* const imsi_dst, char const* const imsi_string_src) {
  strncpy(imsi_dst->data, imsi_string_src, IMSI_BCD_DIGITS_MAX + 1);
  imsi_dst->length = strlen(imsi_dst->data);
  return;
}

/**
 * @brief mme_app_imsi_to_string converts imsi structure to a string
 * @param imsi_dst
 * @param imsi_src
 */

void mme_app_imsi_to_string(
    char* const imsi_dst, mme_app_imsi_t const* const imsi_src) {
  strncpy(imsi_dst, imsi_src->data, IMSI_BCD_DIGITS_MAX + 1);
  return;
}

/**
 * @brief mme_app_is_imsi_empty: checks if an imsi struct is empty returns true
 * if it is empty
 * @param imsi
 * @return
 */
bool mme_app_is_imsi_empty(mme_app_imsi_t const* imsi) {
  return (imsi->length == 0) ? true : false;
}

/**
 * @brief mme_app_imsi_to_u64: converts imsi to uint64 (be careful leading 00
 * will be cut off)
 * @param imsi_src
 * @return
 */

uint64_t mme_app_imsi_to_u64(mme_app_imsi_t imsi_src) {
  uint64_t uint_imsi;
  sscanf(imsi_src.data, "%" SCNu64, &uint_imsi);
  return uint_imsi;
}

mme_ue_s1ap_id_t mme_app_ctx_get_new_ue_id(
    mme_ue_s1ap_id_t* mme_app_ue_s1ap_id_generator_p) {
  mme_ue_s1ap_id_t tmp = 0;
  tmp = __sync_fetch_and_add(mme_app_ue_s1ap_id_generator_p, 1);
  return tmp;
}
