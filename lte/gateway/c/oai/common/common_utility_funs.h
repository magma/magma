#define pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include "mme_config.h"
#include "log.h"
#include "conversions.h"
#include "common_defs.h"

int match_fed_mode_map(const char* imsi, log_proto_t module);
#ifdef __cplusplus
}
#endif
