/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are those
 * of the authors and should not be interpreted as representing official policies,
 * either expressed or implied, of the FreeBSD Project.
 */

#include <stdlib.h>
#include <unistd.h>
#include "glog_logging.h"
#include <glog/logging.h>
#include <glog/stl_logging.h>


void init_logging(const char *log_dir,
                  const char *app_name,
                  uint32_t default_verbosity) {

  google::InitGoogleLogging(app_name);

  // disable log prefix (file, line, etc as it's manually added in oai/log.c)
  FLAGS_log_prefix = false;
  FLAGS_log_dir = log_dir;
  FLAGS_v = default_verbosity;
}

void init_logging(const char *app_name, uint32_t default_verbosity) {
  init_logging("/var/log", app_name, default_verbosity);

  // symlink glog output to /var/log/mme.log
  char *glog_log_path = realpath("/var/log/MME.INFO", NULL);
  symlink(glog_log_path, "/var/log/mme.log");
}


void log_string(int32_t log_level, const char *str) {
  VLOG(log_level) << str;
}

void flush_log(int32_t log_level) {
  google::FlushLogFiles(log_level);
}
