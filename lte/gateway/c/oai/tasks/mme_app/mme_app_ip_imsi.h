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
void initialize_ipv4_map(void);
int mme_app_insert_ue_ipv4_addr(uint32_t ipv4_addr, imsi64_t imsi64);
int mme_app_get_imsi_from_ipv4(uint32_t ipv4_addr, imsi64_t** imsi_list);
void mme_app_remove_ue_ipv4_addr(uint32_t ipv4_addr, imsi64_t imsi64);
#ifdef __cplusplus
}
#endif
