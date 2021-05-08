#include <bits/stdc++.h>
#include <unordered_map>
#include <iostream>
#include "mme_app_ip_imsi.h"
using namespace std;
typedef unordered_map<string, vector<uint64_t>> Ipv4Map;

Ipv4Map ipv4map;

void initialize_ipv4_map() {
  OAILOG_FUNC_IN(LOG_MME_APP);
  ipv4map = Ipv4Map{};
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

void mme_app_log_ipv4_imsi_map() {
  OAILOG_FUNC_IN(LOG_MME_APP);
  for (auto itr = ipv4map.begin(); itr != ipv4map.end(); ++itr) {
    for (auto it_vec = itr->second.begin(); it_vec != itr->second.end();
         ++it_vec) {
      OAILOG_TRACE(
          LOG_MME_APP, "ue_ip: %s \t imsi:%lu \n", itr->first.c_str(), *it_vec);
    }
    OAILOG_TRACE(LOG_MME_APP, "\n");
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

int mme_app_insert_ue_ipv4_addr(uint32_t ipv4_addr, uint64_t imsi64) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  char ipv4[INET_ADDRSTRLEN] = {0};
  inet_ntop(AF_INET, (void*) &ipv4_addr, ipv4, INET_ADDRSTRLEN);
  auto itr = ipv4map.find(ipv4);
  if (itr == ipv4map.end()) {
    vector<uint64_t> vec = {};
    vec.insert(vec.begin(), imsi64);
    ipv4map[ipv4] = vec;
    OAILOG_DEBUG_UE(LOG_MME_APP, imsi64, "Inserting ue_ip:%x \n", ipv4_addr);
  } else {
    OAILOG_DEBUG_UE(
        LOG_MME_APP, imsi64, "Inserting imsi for existing ue_ip:%x \n",
        ipv4_addr);
    ipv4map[ipv4].push_back(imsi64);
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

int mme_app_get_imsi_from_ipv4(uint32_t ipv4_addr, uint64_t* imsi64) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  char ipv4[INET_ADDRSTRLEN] = {0};
  inet_ntop(AF_INET, (void*) &ipv4_addr, ipv4, INET_ADDRSTRLEN);
  auto itr = ipv4map.find(ipv4);
  if (itr == ipv4map.end()) {
    OAILOG_ERROR(LOG_MME_APP, " No imsi found for ip:%x \n", ipv4_addr);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  } else {
    uint8_t idx = 0;
    for (auto vect_itr = itr->second.begin(); vect_itr != itr->second.end();
         vect_itr++) {
      imsi64[idx++] = *vect_itr;
    }
  }
  OAILOG_DEBUG_UE(LOG_MME_APP, *imsi64, " Found imsi for ip:%x \n", ipv4_addr);
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

void mme_app_remove_ue_ipv4_addr(uint32_t ipv4_addr, uint64_t imsi64) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  char ipv4[INET_ADDRSTRLEN] = {0};
  inet_ntop(AF_INET, (void*) &ipv4_addr, ipv4, INET_ADDRSTRLEN);
  auto itr = ipv4map.find(ipv4);
  if (itr == ipv4map.end()) {
    OAILOG_ERROR_UE(
        LOG_MME_APP, imsi64, "No imsi found for ip:%x \n", ipv4_addr);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  } else {
    auto vec_it = itr->second.begin();
    for (; vec_it != itr->second.end(); ++vec_it) {
      if (*vec_it == imsi64) {
        OAILOG_DEBUG_UE(
            LOG_MME_APP, imsi64, "Deleted ue ipv4:%x from ipv4_imsi map \n",
            ipv4_addr);
        itr->second.erase(vec_it);
        vec_it--;
        break;
      }
    }
    if (vec_it == itr->second.end()) {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Failed to remove an entry for IP:%x from ipv4_imsi map \n",
          ipv4_addr);
    }
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}
