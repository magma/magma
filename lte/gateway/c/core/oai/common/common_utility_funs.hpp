#define pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/include/mme_config.h"

int match_fed_mode_map(const char* imsi, log_proto_t module);
int verify_service_area_restriction(tac_t tac,
                                    const regional_subscription_t* reg_sub,
                                    uint8_t num_reg_sub);

#ifdef __cplusplus
}

// C includes --------------------------------------------------------------
extern "C" {
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/esm_data.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_timer.h"
}
// C++ includes ------------------------------------------------------------
#include <czmq.h>
#include <map>
#include <utility>
#include <stddef.h>
#include <stdint.h>

namespace magma {
namespace utils {

template <typename TimerArgType>
class AppTimerUeContext {
 private:
  std::map<int, TimerArgType> app_task_timers_;
  task_zmq_ctx_s* app_task_zmq_ctx_;

 public:
  AppTimerUeContext(AppTimerUeContext const&) = delete;
  void operator=(AppTimerUeContext const&) = delete;

  explicit AppTimerUeContext(task_zmq_ctx_s* zctx) : app_task_zmq_ctx_{zctx} {}

  int StartTimer(size_t msec, timer_repeat_t repeat, zloop_timer_fn handler,
                 const TimerArgType& arg) {
    int timer_id = -1;
    if ((timer_id = start_timer(app_task_zmq_ctx_, msec, repeat, handler,
                                nullptr)) != -1) {
      app_task_timers_.insert(std::pair<int, TimerArgType>(timer_id, arg));
    }
    return timer_id;
  }

  void StopTimer(int timer_id) {
    stop_timer(app_task_zmq_ctx_, timer_id);
    app_task_timers_.erase(timer_id);
  }

  bool PopTimerById(const int timer_id, TimerArgType* arg) {
    try {
      *arg = app_task_timers_.at(timer_id);
      app_task_timers_.erase(timer_id);
      return true;
    } catch (std::out_of_range& e) {
      return false;
    }
  }
};
}  // namespace utils
}  // namespace magma

#endif  /* __cplusplus */
