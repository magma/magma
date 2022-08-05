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
#include "lte/gateway/c/core/oai/test/spgw_task/state_creators.hpp"

#include <gtest/gtest.h>

extern "C" {
#include "lte/gateway/c/core/oai/common/conversions.h"
}

#include "lte/gateway/c/core/oai/include/sgw_context_manager.hpp"
#include "lte/gateway/c/core/oai/include/state_converter.hpp"

namespace magma {

gtpv1u_data_t make_gtpv1u_data(int fd0, int fd1u) {
  gtpv1u_data_t data;
  data.fd0 = fd0;
  data.fd1u = fd1u;
  return data;
}

spgw_state_t make_spgw_state(uint32_t gtpv1u_teid, int fd0, int fd1u) {
  spgw_state_t result;
  result.gtpv1u_teid = gtpv1u_teid;
  result.gtpv1u_data = make_gtpv1u_data(fd0, fd1u);
  return result;
}

// make_bearer_context creates a test bearer context with default values.
s_plus_p_gw_eps_bearer_context_information_t* make_bearer_context(imsi64_t imsi,
                                                                  teid_t teid) {
  // Insert into hashtable
  spgw_create_or_get_ue_context(imsi);
  spgw_update_teid_in_ue_context(imsi, teid);

  // Create underlying object
  auto ctx = sgw_cm_create_bearer_context_information_in_collection(teid);
  auto sgw = &ctx->sgw_eps_bearer_context_information;
  auto pgw = &ctx->pgw_eps_bearer_context_information;

  auto msisdn = "some msisdn";

  // Set PGW context values
  IMSI64_TO_STRING(imsi, (char*)(&pgw->imsi.digit), IMSI_BCD_DIGITS_MAX);
  pgw->imsi_unauthenticated_indicator = 100;
  strncpy(pgw->msisdn, msisdn, strlen(msisdn));

  // Set SGW context values
  sgw->imsi64 = imsi;
  IMSI64_TO_STRING(imsi, (char*)(&sgw->imsi.digit), IMSI_BCD_DIGITS_MAX);
  sgw->imsi_unauthenticated_indicator = 20;
  strncpy(pgw->msisdn, msisdn, strlen(msisdn));
  sgw->mme_teid_S11 = 300;
  sgw->s_gw_teid_S11_S4 = 400;

  std::string ip_str = "191.1.3.0";
  bstring ip_bstr;
  STRING_TO_BSTRING(ip_str, ip_bstr);
  bstring_to_ip_address(ip_bstr, &sgw->mme_ip_address_S11);
  bdestroy(ip_bstr);

  ip_str = "192.0.2.1";
  STRING_TO_BSTRING(ip_str, ip_bstr);
  bstring_to_ip_address(ip_bstr, &sgw->s_gw_ip_address_S11_S4);
  bdestroy(ip_bstr);

  strncpy((char*)&sgw->last_known_cell_Id.plmn, "\x01\x02\x03\x04\x05\x06", 6);
  sgw->last_known_cell_Id.cell_identity.enb_id = 500;
  sgw->last_known_cell_Id.cell_identity.cell_id = 60;
  sgw->last_known_cell_Id.cell_identity.empty = 7;

  return ctx;
}

}  // namespace magma
