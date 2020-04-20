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

/*! \file shared_ts_log.h
   \brief
   \author  Lionel GAUTHIER
   \date 2016
   \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_SHARED_TS_LOG_SEEN
#define FILE_SHARED_TS_LOG_SEEN

#include <sys/time.h>
#include <liblfds710.h>

#include "log.h"
#include "bstrlib.h"

struct timeval;

typedef enum {
  MIN_SH_TS_LOG_CLIENT = 0,
  SH_TS_LOG_TXT        = MIN_SH_TS_LOG_CLIENT,
  MAX_SH_TS_LOG_CLIENT,
} sh_ts_log_app_id_t;

/*! \struct  shared_log_queue_item_t
 * \brief Structure containing a string to be logged.
 * This structure is pushed in thread safe queues by thread producers of logs.
 * This structure is then popped by a dedicated thread that will send back this
 * message to the logger producer in a thread safe manner.
 */
typedef struct shared_log_queue_item_s {
  struct lfds710_stack_element se;
  sh_ts_log_app_id_t app_id; /*!< \brief application identifier. */
  bstring bstr;              /*!< \brief string containing the message. */
  log_private_t log;         /*!< \brief string containing the message. */
} shared_log_queue_item_t;

/*! \struct  log_config_t
 * \brief Structure containing the dynamically configurable parameters of the
 * Logging facilities. This structure is filled by configuration facilities when
 * parsing a configuration file.
 */

//------------------------------------------------------------------------------
int shared_log_get_start_time_sec(void);
void shared_log_reuse_item(shared_log_queue_item_t* item_p);
shared_log_queue_item_t* get_new_log_queue_item(sh_ts_log_app_id_t app_id);
void shared_log_get_elapsed_time_since_start(
    struct timeval* const elapsed_time);
int shared_log_init(const int max_threadsP);
void shared_log_itti_connect(void);
void shared_log_start_use(void);
void shared_log_flush_messages(void);
void shared_log_item(shared_log_queue_item_t* messageP);
#endif /* FILE_SHARED_TS_LOG_SEEN */
