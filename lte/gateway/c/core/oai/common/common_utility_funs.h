#define pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>                // for uint8_t
#include "TrackingAreaIdentity.h"  // for tac_t
#include "common_types.h"          // for regional_subscription_t
#include "log.h"                   // for log_proto_t
#include "obj_hashtable.h"         // IWYU pragma: keep

int match_fed_mode_map(const char* imsi, log_proto_t module);
int verify_service_area_restriction(
    tac_t tac, const regional_subscription_t* reg_sub, uint8_t num_reg_sub);

#ifdef __cplusplus
}
#endif
