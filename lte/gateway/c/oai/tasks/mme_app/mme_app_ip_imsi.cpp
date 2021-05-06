#include <unordered_map>
#include "mme_app_ip_imsi.h"
#include "magma_logging.h"

typedef std::unordered_map<uint32_t, uint64_t> Ipv4Map;

Ipv4Map ipv4map;

void initialize_ipv4_map() {
  OAILOG_FUNC_IN(LOG_MME_APP);
  ipv4map = Ipv4Map{};
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

int mme_app_insert_ue_ipv4_addr(uint32_t ipv4_addr, uint64_t imsi64) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  auto rc = ipv4map.insert({ipv4_addr, imsi64});
  if (!rc.second) {
    OAILOG_ERROR_UE(
        LOG_MME_APP, imsi64, "Failed to insert ipv4_addr:%x \n", ipv4_addr);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

void mme_app_log_ipv4_imsi_map() {
  OAILOG_FUNC_IN(LOG_MME_APP);
  for (auto itr = ipv4map.begin(); itr != ipv4map.end(); ++itr) {
    MLOG(MDEBUG) << "Key: " << itr->first << "and value: " << itr->second;
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

int mme_app_get_imsi_from_ipv4(uint32_t ipv4_addr, uint64_t* imsi64) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  auto itr = ipv4map.find(ipv4_addr);
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
  if (!(ipv4map.erase(ipv4_addr))) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to remove an entry for IP:%x from ipv4_imsi map \n", ipv4_addr);
  }
  OAILOG_ERROR(
      LOG_MME_APP, "Deleted ue ipv4:%x from ipv4_imsi map \n", ipv4_addr);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}
