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

#include <string>
#include <vector>
#include <regex>
#include <algorithm>
#include <fstream>
#include <cstdio>
#include <stdlib.h>
#include <unistd.h>
#include "glog_logging.h"
#include <glog/logging.h>
#include <glog/stl_logging.h>

#include <sys/types.h>
#include <dirent.h>

std::vector<std::string> read_directory(const std::string& dir_path) {
  char* absolute_dir_path = realpath(dir_path.c_str(), nullptr);

  std::vector<std::string> res;
  DIR* dir;
  struct dirent* ent;
  dir = opendir(absolute_dir_path);
  char file_name[MAX_FILE_NAME_LENGTH];
  while ((ent = readdir(dir)) != nullptr) {
    sprintf(file_name, "%s/%s", absolute_dir_path, ent->d_name);
    res.emplace_back(file_name);
  }
  closedir(dir);
  free(absolute_dir_path);
  return res;
}

void init_logging(
    const char* log_dir, const char* app_name, uint32_t default_verbosity) {
  google::InitGoogleLogging(app_name);

  // disable log prefix (file, line, etc as it's manually added in oai/log.c)
  FLAGS_log_prefix = false;
  FLAGS_log_dir    = log_dir;
  FLAGS_v          = default_verbosity;
}

void init_logging(const char* app_name, uint32_t default_verbosity) {
  auto log_dir = "/var/log/";
  init_logging(log_dir, app_name, default_verbosity);

  // symlink glog output to /var/log/mme.log
  auto mme_log_path  = "/var/log/mme.log";
  auto glog_log_path = "/var/log/MME.INFO";

  std::remove(mme_log_path);
  symlink(glog_log_path, mme_log_path);

  // remove old glog files and leave only 10 most recent
  auto dir_files = read_directory(log_dir);
  std::vector<std::string> glog_files;
  std::regex match_mme_log(".*MME.*log.INFO.*");
  for (const auto& filename : dir_files) {
    if (regex_search(filename, match_mme_log)) glog_files.push_back(filename);
  }

  sort(glog_files.begin(), glog_files.end());
  if (glog_files.size() <= 10) return;
  for (unsigned long i = 0; i < glog_files.size() - 10; ++i) {
    std::remove(glog_files[i].c_str());
  }
}

void flush_log(int32_t log_level) {
  google::FlushLogFiles(log_level);
}

void log_string(int32_t log_level, const char* str) {
  VLOG(log_level) << str;
  flush_log(log_level);
}
