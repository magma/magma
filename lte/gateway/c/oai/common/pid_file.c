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

/*! \file pid_file.c
   \brief
   \author  Lionel GAUTHIER
   \date 2016
   \email: lionel.gauthier@eurecom.fr
*/
#include <fcntl.h>
#include <unistd.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <libgen.h>

#include "bstrlib.h"
#include "pid_file.h"

#define PID_DEC_BUF_SIZE 64 /* should be big enough */

int g_fd_pid_file = -1;
__pid_t g_pid     = -1;

//------------------------------------------------------------------------------
char* get_pid_file_name(bstring pid_dir) {
  char* pid_file_name;

  pid_file_name = get_exe_absolute_path(bdatae(pid_dir, "var/run"));

  return pid_file_name;
}

//------------------------------------------------------------------------------
char* get_exe_absolute_path(char const* basepath) {
  char pid_file_name[256] = {0};
  char* exe_basename      = NULL;
  int len                 = 0;

  // get executable name
  len = readlink("/proc/self/exe", pid_file_name, sizeof(pid_file_name));
  if (len == -1) {
    return NULL;
  }
  pid_file_name[len] = '\0';
  exe_basename       = basename(pid_file_name);

  snprintf(
      pid_file_name, sizeof(pid_file_name), "%s/%s.pid", basepath,
      exe_basename);

  return strdup(pid_file_name);
}

//------------------------------------------------------------------------------
int lockfile(int fd, int lock_type) {
  // lock on fd only, not on file on disk (do not prevent another process from
  // modifying the file)
  return lockf(fd, F_TLOCK, 0);
}

//------------------------------------------------------------------------------
bool pid_file_lock(char const* pid_file_name) {
  char pid_dec[PID_DEC_BUF_SIZE] = {0};

  g_fd_pid_file = open(
      pid_file_name, O_RDWR | O_CREAT,
      S_IRUSR | S_IWUSR | S_IRGRP |
          S_IROTH); /* Read/write by owner, read by grp, others */
  if (0 > g_fd_pid_file) {
    printf("filename %s failed %d:%s\n", pid_file_name, errno, strerror(errno));
    return false;
  }

  if (0 > lockfile(g_fd_pid_file, F_TLOCK)) {
    if (EACCES == errno || EAGAIN == errno) {
      printf(
          "filename %s failed %d:%s\n", pid_file_name, errno, strerror(errno));
      close(g_fd_pid_file);
    }
    printf("filename %s failed %d:%s\n", pid_file_name, errno, strerror(errno));
    return false;
  }
  // fruncate file content
  if (ftruncate(g_fd_pid_file, 0) != 0) {
    printf("filename %s failed to truncate\n", pid_file_name);
  }
  // write PID in file
  g_pid = getpid();
  snprintf(pid_dec, sizeof(pid_dec), "%ld", (long) g_pid);
  if (write(g_fd_pid_file, pid_dec, strlen(pid_dec)) != 0) {
    printf("filename %s failed to be written\n", pid_file_name);
  }
  return true;
}

//------------------------------------------------------------------------------
void pid_file_unlock(void) {
  lockfile(g_fd_pid_file, F_ULOCK);
  close(g_fd_pid_file);
  g_fd_pid_file = -1;
}
