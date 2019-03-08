/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include "daemonize.h"

#include <errno.h>
#include <stdlib.h>
#include <string.h>
#include <syslog.h>
#include <sys/stat.h>
#include <unistd.h>

#include "assertions.h"

void daemon_start(void)
{
  pid_t pid, sid; // Our process ID and Session ID

  pid = fork();
  if (pid < 0) {
    Fatal("fork failed: %s", strerror(errno));
  }

  // If we got a good PID, then we can exit the parent process.
  if (pid > 0) {
    exit(EXIT_SUCCESS);
  }
  umask(0);

  // Create a new SID for the child process
  sid = setsid();
  if (sid < 0) {
    Fatal("setsid failed: %s", strerror(errno));
  }

  if ((chdir("/")) < 0) {
    Fatal("chdir failed: %s", strerror(errno));
  }

  close(STDIN_FILENO);
  close(STDOUT_FILENO);
  close(STDERR_FILENO);

  openlog(NULL, 0, LOG_DAEMON);
}

void daemon_stop(void)
{
  closelog();
}
