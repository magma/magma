#define pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/common/log.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/include/mme_config.h"

int match_fed_mode_map(const char* imsi, log_proto_t module);
int verify_service_area_restriction(tac_t tac,
                                    const regional_subscription_t* reg_sub,
                                    uint8_t num_reg_sub);

#ifdef __cplusplus
}
#endif
