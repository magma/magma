/**
 * copyright 2020 The Magma Authors.
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

#pragma once
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_38.401.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/amf/amf_smfDefs.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_asDefs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GIdentityResponse.hpp"

namespace magma5g {
#define AMF_CTXT_MEMBER_IMEI ((uint32_t)1 << 1)
#define MOBILE_IDENTITY_PROTECTION_SCHEME_NULL \
  0x0  // SUCI protection scheme as 0

// Processes Identity Response message
int amf_cn_identity_res(amf_ue_ngap_id_t ue_id, M5GSMobileIdentityMsg* msg,
                        int* amf_cause,
                        const amf_nas_message_decode_status_t* status);

// Performs the identification completion procedure
status_code_e amf_proc_identification_complete(
    const amf_ue_ngap_id_t ue_id, imsi_t* const imsi, imei_t* const imei,
    imeisv_t* const imeisv, uint32_t* const tmsii, guti_m5_t* amf_ctx_guti);

typedef struct amf_guamfi_s {
  amf_plmn_t plmn; /*MCC + MNC*/
  uint8_t amf_regionid;
  uint16_t amf_set_id : 10;
  uint16_t amf_pointer : 6;
} amf_guamfi_t;

// 5G-GUTI
typedef struct amf_guti_m5g_s {
  guamfi_t guamfi;
  uint32_t m_tmsi;
} amf_guti_m5g_t;

typedef struct supi_as_imsi_s {
  // 12 bits for MCC and 12 bits for MNC
  amf_plmn_t plmn;                /*MCC + MNC*/
#define MSIN_MAX_LENGTH 5         // for 10 digits or nibbel
  uint8_t msin[MSIN_MAX_LENGTH];  // last one would be '\0'
} supi_as_imsi_t;

/* Structure for MAP to be used as value against key IMSI
 */
typedef struct guti_and_amf_id_s {
  amf_guti_m5g_t amf_guti;         /* GUTI for SUPI from UE */
  amf_ue_ngap_id_t amf_ue_ngap_id; /* AMF ID to be used to fetch amf_context*/
} guti_and_amf_id_t;

// Generating GUTI based on SUPI/IMSI received from identity message.
void amf_app_generate_guti_on_supi(amf_guti_m5g_t* amf_guti,
                                   supi_as_imsi_t* supi_imsi);

}  // namespace magma5g
