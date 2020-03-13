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
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */
#define SERVICE303

#include <stddef.h>

#include "mme_app_state.h"
#include "service303.h"

static void service303_mme_statistics_read(void)
{
  size_t label = 0;
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  set_gauge("enb_connected", mme_app_desc_p->nb_enb_connected, label);
  set_gauge("ue_registered", mme_app_desc_p->nb_ue_attached, label);
  set_gauge("ue_connected", mme_app_desc_p->nb_ue_connected, label);
  put_mme_nas_state();
  return;
}

void service303_statistics_read(void)
{
  service303_mme_statistics_read();
  return;
}
