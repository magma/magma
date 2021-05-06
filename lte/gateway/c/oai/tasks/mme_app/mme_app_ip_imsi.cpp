#include <unordered_map>
#include <iostream>
#include "mme_app_ip_imsi.h"
using namespace std;
typedef unordered_map<string, uint64_t> Ipv4Map;

Ipv4Map ipv4map;

void initialize_ipv4_map() {
  OAILOG_FUNC_IN(LOG_MME_APP);
  ipv4map = Ipv4Map{};
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

void mme_app_log_ipv4_imsi_map() {
  OAILOG_FUNC_IN(LOG_MME_APP);
  for (auto itr = ipv4map.begin(); itr != ipv4map.end(); ++itr) {
    OAILOG_TRACE(
        LOG_MME_APP, "key: %s \t value:%lu \n", itr->first.c_str(),
        itr->second);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

int mme_app_insert_ue_ipv4_addr(uint32_t ipv4_addr, uint64_t imsi64) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  char ipv4[INET_ADDRSTRLEN];
  inet_ntop(AF_INET, (void*) &ipv4_addr, ipv4, INET_ADDRSTRLEN);
  auto rc = ipv4map.insert({ipv4, imsi64});
  if (!rc.second) {
    OAILOG_ERROR_UE(LOG_MME_APP, imsi64, "Failed to insert ipv4:%x \n", ipv4);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  mme_app_log_ipv4_imsi_map();
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

int mme_app_get_imsi_from_ipv4(uint32_t ipv4_addr, uint64_t* imsi64) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  char ipv4[INET_ADDRSTRLEN];
  inet_ntop(AF_INET, (void*) &ipv4_addr, ipv4, INET_ADDRSTRLEN);
  auto itr = ipv4map.find(ipv4);
  if (itr == ipv4map.end()) {
    OAILOG_ERROR(LOG_MME_APP, " No imsi found for ip:%x \n", ipv4_addr);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  *imsi64 = itr->second;
  OAILOG_DEBUG_UE(LOG_MME_APP, *imsi64, " Found imsi for ip:%x \n", ipv4_addr);
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

void mme_app_remove_ue_ipv4_addr(uint32_t ipv4_addr) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  char ipv4[INET_ADDRSTRLEN];
  inet_ntop(AF_INET, (void*) &ipv4_addr, ipv4, INET_ADDRSTRLEN);
  if (!(ipv4map.erase(ipv4))) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to remove an entry for IP:%x from ipv4_imsi map \n", ipv4_addr);
  }
  OAILOG_ERROR(
      LOG_MME_APP, "Deleted ue ipv4:%x from ipv4_imsi map \n", ipv4_addr);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}
