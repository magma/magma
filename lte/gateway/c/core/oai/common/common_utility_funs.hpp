#define pragma once

#include "lte/gateway/c/core/oai/include/mme_config.hpp"

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif

int match_fed_mode_map(const char* imsi, log_proto_t module);
int verify_service_area_restriction(tac_t tac,
                                    const regional_subscription_t* reg_sub,
                                    uint8_t num_reg_sub);

int mme_config_find_mnc_length(const char mcc_digit1P, const char mcc_digit2P,
                               const char mcc_digit3P, const char mnc_digit1P,
                               const char mnc_digit2P, const char mnc_digit3P);
