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

#if HAVE_CONFIG_H
#include "config.h"
#endif

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <signal.h>
#include <ctype.h>

#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/common/backtrace.h"
#include "lte/gateway/c/core/oai/common/assertions.h"
#include "lte/gateway/c/core/oai/lib/itti/signals.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"

#ifndef SIG_DEBUG
#define SIG_DEBUG(x, args...)                                                  \
  do {                                                                         \
    fprintf(stdout, "[SIG][D]" x, ##args);                                     \
  } while (0)
#endif
#ifndef SIG_ERROR
#define SIG_ERROR(x, args...)                                                  \
  do {                                                                         \
    fprintf(stdout, "[SIG][E]" x, ##args);                                     \
  } while (0)
#endif

static sigset_t set;

#if LINK_GCOV
void gcov_flush(void);
#endif

// We have had cases where threads have been created before the main thread
// masks the signals mostly as a side effect of GRPC class instantiation. Add
// a simple assert to make sure we prevent this regression by asserting that
// the main thread is the only thread during signal masking.
//
// For now parse the output of the status file
// $ cat /proc/$$/status
// Name:   bash
// State:  S (sleeping)
// Tgid:   3515
// Pid:    3515
// PPid:   3452
// ..
// Threads: 1
//
static const char THREADS_STR[] = "Threads:";
static const char PROC_PATH[]   = "/proc/%d/status";

static int get_thread_count(pid_t pid) {
  char path[40], line[100], *p;
  int num_threads = -1;
  FILE* statusf;

  snprintf(path, sizeof(path), PROC_PATH, pid);

  statusf = fopen(path, "r");
  if (!statusf) return -1;

  while (fgets(line, sizeof(line), statusf)) {
    if (strncmp(line, "Threads:", 8) != 0) continue;
    // Ignore "Threads:" and whitespace
    p = line + strlen(THREADS_STR);
    while (isspace(*p)) ++p;
    num_threads = atoi(p);
    break;
  }
  fclose(statusf);
  return num_threads;
}

int signal_mask(void) {
  /*
   * We set the signal mask to avoid threads other than the main thread
   * to receive the timer signal. Note that threads created will inherit this
   * configuration.
   */
#if !MME_UNIT_TEST
  SIG_DEBUG("MME_UNIT_TEST Flag is Disabled\n");
  DevAssert(get_thread_count(getpid()) == 1);
#endif
  sigemptyset(&set);
  sigaddset(&set, SIGABRT);
  sigaddset(&set, SIGINT);
  sigaddset(&set, SIGTERM);

  if (sigprocmask(SIG_BLOCK, &set, NULL) < 0) {
    perror("sigprocmask");
    return -1;
  }

  return 0;
}

int signal_handle(int* end, task_zmq_ctx_t* task_ctx) {
  int ret;
  siginfo_t info;

  sigemptyset(&set);
  sigaddset(&set, SIGABRT);
  sigaddset(&set, SIGINT);
  sigaddset(&set, SIGTERM);

  if (sigprocmask(SIG_BLOCK, &set, NULL) < 0) {
    perror("sigprocmask");
    return -1;
  }

  /*
   * Block till a signal is received.
   * * * NOTE: The signals defined by set are required to be blocked at the time
   * * * of the call to sigwait() otherwise sigwait() is not successful.
   */
  if ((ret = sigwaitinfo(&set, &info)) == -1) {
    perror("sigwait");
    return ret;
  }

  /*
   * Dispatch the signal to sub-handlers
   */
  switch (info.si_signo) {
    case SIGABRT:
      SIG_DEBUG("Received SIGABORT\n");
      backtrace_handle_signal(&info);
      break;

    case SIGINT:
    case SIGTERM:
      printf("Received SIGINT or SIGTERM\n");
      send_terminate_message_fatal(task_ctx);
      *end = 1;
      break;

    default:
      SIG_ERROR("Received unknown signal %d\n", info.si_signo);
      break;
  }

  return 0;
}
