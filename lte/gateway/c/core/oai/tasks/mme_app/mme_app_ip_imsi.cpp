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

#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_ip_imsi.hpp"

#include <unordered_map>
#include <iostream>

#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_ueip_imsi_map.hpp"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_state_manager.hpp"

using magma::lte::MmeNasStateManager;
// Description: Logs the content of ueip_imsi map
void mme_app_log_ipv4_imsi_map() {
  OAILOG_FUNC_IN(LOG_MME_APP);
  UeIpImsiMap& ueip_imsi_map =
      MmeNasStateManager::getInstance().get_mme_ueip_imsi_map();
  for (const auto& itr_map : ueip_imsi_map) {
    for (const auto& it_vec : itr_map.second) {
      OAILOG_TRACE(LOG_MME_APP, "ue_ip: %s \t imsi:%lu \n",
                   itr_map.first.c_str(), it_vec);
    }
    OAILOG_TRACE(LOG_MME_APP, "\n");
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

void mme_app_log_ipv6_imsi_map() {
  OAILOG_FUNC_IN(LOG_MME_APP);
  UeIpImsiMap& ueip_imsi_map =
      MmeNasStateManager::getInstance().get_mme_ueip_imsi_map();
  for (const auto& itr_map : ueip_imsi_map) {
    for (const auto& it_vec : itr_map.second) {
      OAILOG_TRACE(LOG_MME_APP, "IPv6 ue_ip: %s \t imsi:%lu \n",
                   itr_map.first.c_str(), it_vec);
    }
    OAILOG_TRACE(LOG_MME_APP, "\n");
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

/* Description: ue_ip address is allocated by either roaming PGWs or mobilityd
 * So there is possibility to allocate same ue ip address for different UEs.
 * So defining ue_ip_imsi map with key as ue_ip and value as list of imsis
 * having same ue_ip
 */
int mme_app_insert_ue_ipv4_addr(uint32_t ipv4_addr, imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  UeIpImsiMap& ueip_imsi_map =
      MmeNasStateManager::getInstance().get_mme_ueip_imsi_map();
  char ipv4[INET_ADDRSTRLEN] = {0};
  inet_ntop(AF_INET, (void*)&ipv4_addr, ipv4, INET_ADDRSTRLEN);
  auto itr_map = ueip_imsi_map.find(ipv4);
  if (itr_map == ueip_imsi_map.end()) {
    std::vector<uint64_t> vec = {imsi64};
    ueip_imsi_map[ipv4] = vec;
    OAILOG_DEBUG_UE(LOG_MME_APP, imsi64, "Inserting ue_ipv4:%x \n", ipv4_addr);
  } else {
    OAILOG_DEBUG_UE(LOG_MME_APP, imsi64,
                    "Inserting imsi for existing ue_ip:%x \n", ipv4_addr);
    ueip_imsi_map[ipv4].push_back(imsi64);
  }
  MmeNasStateManager::getInstance().write_mme_ueip_imsi_map_to_db();
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

/* Description: ue_ip address is allocated by either roaming PGWs or mobilityd
 * So there is possibility to allocate same ue ip address for different UEs.
 * So defining ue_ip_imsi map with key as ue_ip and value as list of imsis
 * having same ue_ip
 */
int mme_app_insert_ue_ipv6_addr(struct in6_addr ipv6_addr, imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  UeIpImsiMap& ueip_imsi_map =
      MmeNasStateManager::getInstance().get_mme_ueip_imsi_map();
  char ipv6[INET6_ADDRSTRLEN] = {0};
  inet_ntop(AF_INET6, (void*)&ipv6_addr, ipv6, INET6_ADDRSTRLEN);
  auto itr_map = ueip_imsi_map.find(ipv6);
  if (itr_map == ueip_imsi_map.end()) {
    std::vector<uint64_t> vec = {imsi64};
    ueip_imsi_map[ipv6] = vec;
    OAILOG_DEBUG_UE(LOG_MME_APP, imsi64, "Inserting ue_ipv6:%s \n", ipv6);
  } else {
    OAILOG_DEBUG_UE(LOG_MME_APP, imsi64,
                    "Inserting imsi for existing ue_ipv6:%s \n", ipv6);
    ueip_imsi_map[ipv6].push_back(imsi64);
  }
  MmeNasStateManager::getInstance().write_mme_ueip_imsi_map_to_db();
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

/* Description: The function shall provide list of imsis allocated for
 * ue ip address; Imsi list is dynamically created and filled with imsis
 * The caller of function needs to free the memory allocated for imsi list
 */
int mme_app_get_imsi_from_ipv4(uint32_t ipv4_addr, imsi64_t** imsi_list) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  UeIpImsiMap& ueip_imsi_map =
      MmeNasStateManager::getInstance().get_mme_ueip_imsi_map();
  int num_imsis = 0;
  char ipv4[INET_ADDRSTRLEN] = {0};
  inet_ntop(AF_INET, (void*)&ipv4_addr, ipv4, INET_ADDRSTRLEN);
  auto itr_map = ueip_imsi_map.find(ipv4);
  if (itr_map == ueip_imsi_map.end()) {
    OAILOG_ERROR(LOG_MME_APP, " No imsi found for ip:%x \n", ipv4_addr);
  } else {
    uint8_t idx = 0;
    num_imsis = itr_map->second.size();
    (*imsi_list) = (imsi64_t*)calloc(num_imsis, sizeof(imsi64_t));

    for (const auto& vect_itr : itr_map->second) {
      (*imsi_list)[idx++] = vect_itr;
      OAILOG_DEBUG_UE(LOG_MME_APP, vect_itr, " Found imsi for ip:%x \n",
                      ipv4_addr);
    }
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, num_imsis);
}

/* Description: The function shall provide list of imsis allocated for
 * ue ipv6 address; Imsi list is dynamically created and filled with imsis
 * The caller of function needs to free the memory allocated for imsi list
 */
int mme_app_get_imsi_from_ipv6(struct in6_addr ipv6_addr,
                               imsi64_t** imsi_list) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  UeIpImsiMap& ueip_imsi_map =
      MmeNasStateManager::getInstance().get_mme_ueip_imsi_map();
  int num_imsis_ipv6 = 0;
  char ipv6[INET6_ADDRSTRLEN] = {0};
  inet_ntop(AF_INET6, (void*)&ipv6_addr, ipv6, INET6_ADDRSTRLEN);

  if (inet_ntop(AF_INET6, (void*)&ipv6_addr, ipv6, INET6_ADDRSTRLEN) == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "IPV6 address conversion IS NULL:%s \n", ipv6);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  auto itr_map = ueip_imsi_map.find(ipv6);

  if (itr_map == ueip_imsi_map.end()) {
    OAILOG_ERROR(LOG_MME_APP, " No imsi found for ip:%s \n", ipv6);
  } else {
    uint8_t idx = 0;
    num_imsis_ipv6 = itr_map->second.size();
    (*imsi_list) = (imsi64_t*)calloc(num_imsis_ipv6, sizeof(imsi64_t));

    for (const auto& vect_itr : itr_map->second) {
      (*imsi_list)[idx++] = vect_itr;
      OAILOG_DEBUG_UE(LOG_MME_APP, vect_itr, " Found imsi for ip:%s \n", ipv6);
    }
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, num_imsis_ipv6);
}

/* Description: Shall remove an entry from ueip_imsi map for matching
 *  ueip and imsi
 */
void mme_app_remove_ue_ipv4_addr(uint32_t ipv4_addr, imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  UeIpImsiMap& ueip_imsi_map =
      MmeNasStateManager::getInstance().get_mme_ueip_imsi_map();
  char ipv4[INET_ADDRSTRLEN] = {0};
  inet_ntop(AF_INET, (void*)&ipv4_addr, ipv4, INET_ADDRSTRLEN);
  auto itr_map = ueip_imsi_map.find(ipv4);
  if (itr_map == ueip_imsi_map.end()) {
    OAILOG_ERROR_UE(LOG_MME_APP, imsi64, "No imsi found for ip:%x \n",
                    ipv4_addr);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  } else {
    auto vec_it = itr_map->second.begin();
    while (vec_it != itr_map->second.end()) {
      if (*vec_it == imsi64) {
        OAILOG_DEBUG_UE(LOG_MME_APP, imsi64,
                        "Deleted ue ipv4:%x from ipv4_imsi map \n", ipv4_addr);
        vec_it = itr_map->second.erase(vec_it);
        if (itr_map->second.empty()) {
          ueip_imsi_map.erase(ipv4);
        }
        MmeNasStateManager::getInstance().write_mme_ueip_imsi_map_to_db();
        break;
      } else {
        vec_it++;
      }
    }
    if (ueip_imsi_map.find(ipv4) != ueip_imsi_map.end() &&
        vec_it == itr_map->second.end()) {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Failed to remove an entry for ue_ip:%x from ipv4_imsi map \n",
          ipv4_addr);
    }
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

/* Description: Shall remove an entry from ueip_imsi map for matching
 *  ueip and imsi
 */
void mme_app_remove_ue_ipv6_addr(struct in6_addr ipv6_addr, imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  UeIpImsiMap& ueip_imsi_map =
      MmeNasStateManager::getInstance().get_mme_ueip_imsi_map();
  char ipv6[INET6_ADDRSTRLEN] = {0};
  inet_ntop(AF_INET6, (void*)&ipv6_addr, ipv6, INET6_ADDRSTRLEN);
  auto itr_map = ueip_imsi_map.find(ipv6);
  if (itr_map == ueip_imsi_map.end()) {
    OAILOG_ERROR_UE(LOG_MME_APP, imsi64, "No imsi found for ip:%s \n", ipv6);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  } else {
    auto vec_it = itr_map->second.begin();
    while (vec_it != itr_map->second.end()) {
      if (*vec_it == imsi64) {
        OAILOG_DEBUG_UE(LOG_MME_APP, imsi64,
                        "Deleted ue ipv6:%s from ipv6_imsi map \n", ipv6);
        vec_it = itr_map->second.erase(vec_it);
        if (itr_map->second.empty()) {
          ueip_imsi_map.erase(ipv6);
        }
        MmeNasStateManager::getInstance().write_mme_ueip_imsi_map_to_db();
        break;
      } else {
        vec_it++;
      }
    }
    if (ueip_imsi_map.find(ipv6) != ueip_imsi_map.end() &&
        vec_it == itr_map->second.end()) {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Failed to remove an entry for ue_ipv6:%s from ipv6_imsi map \n",
          ipv6);
    }
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}
