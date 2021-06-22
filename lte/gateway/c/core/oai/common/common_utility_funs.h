#define pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include "mme_config.h"
#include "log.h"
#include "conversions.h"
#include "common_defs.h"

int match_fed_mode_map(const char* imsi, log_proto_t module);
int verify_service_area_restriction(
    tac_t tac, const regional_subscription_t* reg_sub, uint8_t num_reg_sub);

#ifdef __cplusplus
}
#endif
