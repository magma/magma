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
#pragma once

#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"
#include "lte/gateway/c/core/oai/include/map.h"

namespace magma5g {

class AmfUeContextStorage {
 private:
  AmfUeContextStorage()                           = default;
  AmfUeContextStorage(const AmfUeContextStorage&) = delete;
  AmfUeContextStorage& operator=(const AmfUeContextStorage&) = delete;

  // std::vector< std::shared_ptr<ue_m5gmm_context_t> > ue_context_pool;
  amf_ue_ngap_id_t amf_app_ue_ngap_id_generator;

  // hash table to fetch fast access context.
  // id-AMF-UE-NGAP-ID <--> ue_m5gmm_context_t map
  magma::map_s<amf_ue_ngap_id_t, std::shared_ptr<ue_m5gmm_context_t> >
      amfid_ue_context_map;
  // id-GNB-UE-NGAP-ID <---> ue_m5gmm_context_t map
  magma::map_s<gnb_ue_ngap_id_t, std::shared_ptr<ue_m5gmm_context_t> >
      gnbid_ue_context_map;
  // GUTI <---> ue_m5gmm_context_t map
  magma::map_s<guti_m5_t, std::shared_ptr<ue_m5gmm_context_t> >
      guti_ue_context_map;
  // SUPI  <---> ue_m5gmm_context_t map
  magma::map_s<imsi64_t, std::shared_ptr<ue_m5gmm_context_t> >
      supi_ue_context_map;

  amf_ue_ngap_id_t generate_amf_ue_ngap_id() {
    return amf_app_ue_ngap_id_generator + 1;
  }

 public:
  std::shared_ptr<ue_m5gmm_context_t> amf_create_new_ue_context();

  AmfUeContextStorage& getUeContextStorage() {
    static AmfUeContextStorage ue_context_storage;
    return ue_context_storage;
  }

  // id-AMF-UE-NGAP-ID <--> ue_m5gmm_context_t map
  bool amf_insert_into_amfid_ue_context_map(
      std::shared_ptr<ue_m5gmm_context_t> pContext);
  bool amf_remove_from_amfid_ue_context_map(amf_ue_ngap_id_t ue_amf_id);
  std::shared_ptr<ue_m5gmm_context_t> amf_get_from_amfid_ue_context_map(
      amf_ue_ngap_id_t ue_amf_id);

  // id-GNB-UE-NGAP-ID <---> ue_m5gmm_context_t map
  bool amf_insert_into_gnbid_ue_context_map(
      std::shared_ptr<ue_m5gmm_context_t> pContext);
  bool amf_remove_from_gnbid_ue_context_map(gnb_ue_ngap_id_t ue_gnb_id);
  std::shared_ptr<ue_m5gmm_context_t> amf_get_from_gnbid_ue_context_map(
      gnb_ue_ngap_id_t ue_gnb_id);

  // GUTI <---> ue_m5gmm_context_t map
  bool amf_insert_into_guti_ue_context_map(
      std::shared_ptr<ue_m5gmm_context_t> pContext);
  bool amf_remove_from_guti_ue_context_map(guti_m5_t guti);
  std::shared_ptr<ue_m5gmm_context_t> amf_get_from_guti_ue_context_map(
      guti_m5_t guti);

  // SUPI  <---> ue_m5gmm_context_t map
  bool amf_insert_into_supi_ue_context_map(
      std::shared_ptr<ue_m5gmm_context_t> pContext);
  bool amf_remove_from_supi_ue_context_map(imsi64_t supi);
  std::shared_ptr<ue_m5gmm_context_t> amf_get_from_supi_ue_context_map(
      imsi64_t supi);

  bool amf_remove_ue_context_from_cache(amf_ue_ngap_id_t ue_amf_id);
  bool amf_add_ue_context_in_cache(
      std::shared_ptr<ue_m5gmm_context_t> ue_ctxt_p);
};

}  // namespace magma5g
