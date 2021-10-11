/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
#pragma once

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#include "common_defs.h"
#include "common_types.h"
#include "nas_timer.h"

void initialize_mme_ue_id_timer_id_map(void);

void mme_app_upsert_mme_ue_id_timer_id(
    mme_ue_s1ap_id_t mme_ue_id, long timer_id);
void mme_app_remove_mme_ue_id_timer_id(mme_ue_s1ap_id_t mme_ue_id);

long mme_app_get_timer_id_from_mme_ue_id(mme_ue_s1ap_id_t mme_ue_id);

#ifdef __cplusplus
}
#endif
