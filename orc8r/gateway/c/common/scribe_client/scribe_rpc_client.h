/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once
#ifdef __cplusplus
extern "C" {
#endif

// caller is responsible for string memory allocation and release.
typedef struct scribe_int_param {
  const char *key; //string
  const int val; // integer
} scribe_int_param_t;

typedef struct scribe_string_param {
  const char *key;
  const char *val;
} scribe_string_param_t;

/**
 * Log one scribe entry to the given category on scribe. Default current
 * timestamp and sampleRate 1 will be used.
 *
 * @param category: category name of the scribe category to log to.
 * @param int_params[]: an array of scribe_int_param_t, where each
 * scribe_int_param_t contains a str name, and a int value.
 * @param int_params_len: length of the above array.
 * @param str_params[]: an array of scribe_string_param_t, where each
   scribe_string_param_t contains a str name, and a str value.
 * @param str_params_len: length of the above array.
 */
int log_to_scribe(
  char const *category,
  scribe_int_param_t *int_params,
  int int_params_len,
  scribe_string_param_t *str_params,
  int str_params_len);

/**
 * Log one scribe entry to the given category on scribe. Default current
 * timestamp will be used.
 *
 * @param category: category name of the scribe category to log to.
 * @param int_params[]: an array of scribe_int_param_t, where each
 * scribe_int_param_t contains a str name, and a int value.
 * @param int_params_len: length of the above array.
 * @param str_params[]: an array of scribe_string_param_t, where each
   scribe_string_param_t contains a str name, and a str value.
 * @param str_params_len: length of the above array.
 * @param sampling_rate: a float between 0 and 1 indicating the desired
 * samplingRate of the log. The ScribeClient will throw a die with value in
 * [0, 1) and drop the attempt to log the entry if the result of the die is
 * larger than the samplingRate.
 */
int log_to_scribe_with_sampling_rate(
  char const *category,
  scribe_int_param_t *int_params,
  int int_params_len,
  scribe_string_param_t *str_params,
  int str_params_len,
  float sampling_rate);

/**
 * Log one scribe entry to the given category on scribe with a timestamp.
 *
 * @param category: category name of the scribe category to log to.
 * @param time: a timestamp to associate the log message with.
 * @param int_params[]: an array of scribe_int_param_t, where each
 * scribe_int_param_t contains a str name, and a int value.
 * @param int_params_len: length of the above array.
 * @param str_params[]: an array of scribe_string_param_t, where each
   scribe_string_param_t contains a str name, and a str value.
 * @param str_params_len: length of the above array.
 * @param sampling_rate: a float between 0 and 1 indicating the desired
 * samplingRate of the log. The ScribeClient will throw a die with value in
 * [0, 1) and drop the attempt to log the entry if the result of the die is
 * larger than the samplingRate.
 */
int log_to_scribe_with_time_and_sampling_rate(
  char const *category,
  int time,
  scribe_int_param_t *int_params,
  int int_params_len,
  scribe_string_param_t *str_params,
  int str_params_len,
  float sampling_rate);

#ifdef __cplusplus
}
#endif
