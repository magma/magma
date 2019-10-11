/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-----------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include "state_converter.h"

namespace magma {
namespace lte {

StateConverter::StateConverter() = default;
StateConverter::~StateConverter() = default;

/*************************************************/
/*        Common Types -> Proto                 */
/*************************************************/

void StateConverter::plmn_to_chars(const plmn_t &state_plmn, char *plmn_array)
{
  plmn_array[0] = (char) state_plmn.mcc_digit2;
  plmn_array[1] = (char) state_plmn.mcc_digit1;
  plmn_array[2] = (char) state_plmn.mnc_digit3;
  plmn_array[3] = (char) state_plmn.mcc_digit3;
  plmn_array[4] = (char) state_plmn.mnc_digit2;
  plmn_array[5] = (char) state_plmn.mnc_digit1;
}

void StateConverter::guti_to_proto(const guti_t &state_guti, Guti *guti_proto)
{
  guti_proto->Clear();

  char *plmn_array = (char *) malloc(sizeof(plmn_t));
  memcpy(plmn_array, &state_guti.gummei.plmn, sizeof(plmn_t));
  guti_proto->set_plmn(plmn_array);
  guti_proto->set_mme_gid(state_guti.gummei.mme_gid);
  guti_proto->set_mme_code(state_guti.gummei.mme_code);
  guti_proto->set_m_tmsi(state_guti.m_tmsi);
}

void StateConverter::ecgi_to_proto(const ecgi_t &state_ecgi, Ecgi *ecgi_proto)
{
  ecgi_proto->Clear();

  char plmn_array[PLMN_BYTES];
  plmn_to_chars(state_ecgi.plmn, plmn_array);
  ecgi_proto->set_plmn(plmn_array);
  ecgi_proto->set_enb_id(state_ecgi.cell_identity.enb_id);
  ecgi_proto->set_cell_id(state_ecgi.cell_identity.cell_id);
  ecgi_proto->set_empty(state_ecgi.cell_identity.empty);
}

void StateConverter::proto_to_ecgi(const Ecgi &ecgi_proto, ecgi_t *state_ecgi)
{
  strncpy((char *) &state_ecgi->plmn, ecgi_proto.plmn().c_str(), PLMN_BYTES);

  state_ecgi->cell_identity.enb_id = ecgi_proto.enb_id();
  state_ecgi->cell_identity.cell_id = ecgi_proto.cell_id();
  state_ecgi->cell_identity.empty = ecgi_proto.empty();
}

} // namespace lte
} // namespace magma
