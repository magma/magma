/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef C_SCRIBE_RPC_CLIENT_FOR_CPP_H
#define C_SCRIBE_RPC_CLIENT_FOR_CPP_H

/**
 *  C++ API for logging one scribe entry to the given category with a timestamp and a sampling rate .
 * @param category category name of the scribe category to log to.
 * @param time  a timestamp to associate the log message with.
 * @param int_params a map of string keys to int values to log
 * @param str_params a map of string keys to string values to log
 * @param sampling_rate a float between 0 and 1 indicating the desired
 * sampling_rate of the log. The ScribeClient will throw a die with value in
 * [0, 1) and drop the attempt to log the entry if the result of the die is
 * larger than the sampling_rate.
 */
void log_to_scribe_with_time_and_sampling_rate(
  std::string category,
  time_t time,
  std::map<std::string, int> int_params,
  std::map<std::string, std::string> str_params,
  float sampling_rate);

/**
 * C++ API for logging one scribe entry with default timestamp(current time), and sampling rate 1.
 * @param category category name of the scribe category to log to.
 * @param int_params a map of string keys to int values to log
 * @param str_params a map of string keys to string values to log
 */
void log_to_scribe(
  std::string category,
  std::map<std::string, int> int_params,
  std::map<std::string, std::string> str_params);

#endif //C_SCRIBE_RPC_CLIENT_FOR_CPP_H
