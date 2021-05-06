#pragma once
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#include "common_defs.h"
#include "common_types.h"
void initialize_ipv4_map(void);
int mme_app_insert_ue_ipv4_addr(uint32_t ipv4_addr, uint64_t imsi64);
int mme_app_get_imsi_from_ipv4(uint32_t ipv4_addr, uint64_t* imsi64);
void mme_app_remove_ue_ipv4_addr(uint32_t ipv4_addr);
#ifdef __cplusplus
}
#endif
