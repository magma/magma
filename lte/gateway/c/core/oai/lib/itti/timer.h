/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are
 * those of the authors and should not be interpreted as representing official
 * policies, either expressed or implied, of the FreeBSD Project.
 */

#ifndef TIMER_H_
#define TIMER_H_

#include <signal.h>
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <sys/types.h>

#include "intertask_interface_types.h"
#include "intertask_interface.h"

#define SIGTIMER SIGRTMIN

typedef enum timer_type_s {
  TIMER_PERIODIC,
  TIMER_ONE_SHOT,
  TIMER_TYPE_MAX,
} timer_type_t;

typedef enum {
  TIMER_OK        = 0,
  TIMER_NOT_FOUND = -1,
  TIMER_ERR       = -2,
} timer_result_t;

int timer_handle_signal(siginfo_t* info, task_zmq_ctx_t* task_ctx);

/** \brief Request a new timer
 *  \param interval_sec timer interval in seconds
 *  \param interval_us  timer interval in micro seconds
 *  \param task_id      task id of the task requesting the timer
 *  \param instance     instance of the task requesting the timer
 *  \param type         timer type
 *  \param timer_arg    extra data to save with the timer
 *  \param arg_size     size of the data you are saving
 *  \param timer_id     unique timer identifier
 *  @returns -1 on failure, 0 otherwise
 **/
int timer_setup(
    uint32_t interval_sec, uint32_t interval_us, task_id_t task_id,
    int32_t instance, timer_type_t type, void* timer_arg, size_t arg_size,
    long* timer_id);

int timer_handle_expired(long timer_id);

bool timer_exists(long timer_id);

/** \brief Remove the timer from list
 *  \param timer_id unique timer id
 *  @returns -1 on failure, 0 otherwise
 **/

int timer_remove(long timer_id, void** arg);

#define timer_stop timer_remove

/** \brief Initialize timer task and its API
 *  \param mme_config MME common configuration
 *  @returns -1 on failure, 0 otherwise
 **/
int timer_init(void);

#endif
