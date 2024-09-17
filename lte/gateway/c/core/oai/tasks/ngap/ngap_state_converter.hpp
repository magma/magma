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
/****************************************************************************
  Source      ngap_state_converter.h
  Date        2020/07/28
  Author      Ashish Prajapati
  Subsystem   Access and Mobility Management Function
  Description Defines NG Application Protocol Messages

*****************************************************************************/

#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <cstdint>

#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/lib/hashtable/hashtable.h"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_types.h"

#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/state_converter.hpp"
#include "lte/protos/oai/ngap_state.pb.h"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_state.hpp"
using namespace std;
using namespace magma::lte;
using namespace magma::lte::oai;

namespace magma5g {

class NgapStateConverter : magma::lte::StateConverter {
 public:
  static void state_to_proto(ngap_state_t* state, oai::NgapState* proto);

  static void proto_to_state(const oai::NgapState& proto, ngap_state_t* state);

  /**
   * Serializes ngap_imsi_map_t to NgapImsiMap proto
   */
  static void ngap_imsi_map_to_proto(const ngap_imsi_map_t* ngap_imsi_map,
                                     oai::NgapImsiMap* ngap_imsi_proto);

  /**
   * Deserializes ngap_imsi_map_t from NgapImsiMap proto
   */
  static void proto_to_ngap_imsi_map(const oai::NgapImsiMap& ngap_imsi_proto,
                                     ngap_imsi_map_t* ngap_imsi_map);

  /**
   * Serializes m5g_supported_ta_list_t to Ngap_SupportedTaList proto
   */
  static void supported_ta_list_to_proto(
      const m5g_supported_ta_list_t* supported_ta_list,
      oai::Ngap_SupportedTaList* supported_ta_list_proto);

  static void proto_to_supported_ta_list(
      m5g_supported_ta_list_t* supported_ta_list_state,
      const oai::Ngap_SupportedTaList& supported_ta_list_proto);

  /**
   * Serializes m5g_supported_tai_items_t to supported_tai_item proto
   */
  static void supported_tai_item_to_proto(
      const m5g_supported_tai_items_t* state_supported_tai_item,
      oai::Ngap_SupportedTaiItems* supported_tai_item_proto);

  static void proto_to_supported_tai_items(
      m5g_supported_tai_items_t* supported_tai_item_state,
      const oai::Ngap_SupportedTaiItems& supported_tai_item_proto);

  static void gnb_to_proto(gnb_description_t* gnb, oai::GnbDescription* proto);

  static void proto_to_gnb(const oai::GnbDescription& proto,
                           gnb_description_t* gnb);

  static void ue_to_proto(const m5g_ue_description_t* ue,
                          oai::Ngap_UeDescription* proto);

  static void proto_to_ue(const oai::Ngap_UeDescription& proto,
                          m5g_ue_description_t* ue);

 private:
  NgapStateConverter();
  ~NgapStateConverter();
};
}  // namespace magma5g
