/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <string>

namespace magma {
namespace service303 {

/**
 * ProcFileUtils is a helper class to parse proc files for process information
 */
class ProcFileUtils final {
  public:
    /*
     * memory_info_t wraps the needed information from the /status file
     */
    typedef struct memory_info_s {
      double physical_mem = -1;
      double virtual_mem = -1;
    } memory_info_t;

  public:
    /*
     * Parses the /proc/self/status file for information on memory usage
     *
     * @return memory_info_t containing virtual and physical memory usage
     */
    static const memory_info_t getMemoryInfo();

  private:
    /*
     * Helper function to read from the proc file stream and output the value if
     * the prefix is found
     *
     * @return -1 if the string isn't the prefix we're looking for, otherwise
     *    return the actual value
     */
    static double parseForPrefix(
      std::ifstream& infile,
      const std::string& to_compare,
      const std::string& prefix_name);

  private:
    static const std::string STATUS_FILE;
    // status file labels
    static const std::string VIRTUAL_MEM_PREFIX;
    static const std::string PHYSICAL_MEM_PREFIX;
};

}
}
