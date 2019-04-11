/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package utils

import (
	"fmt"
	"strconv"
	"time"
)

const (
	ParamQuery      = "query"
	ParamRangeStart = "start"
	ParamRangeEnd   = "end"
	ParamStepWidth  = "step"
	ParamTime       = "time"

	StatusSuccess = "success"
)

func ParseTime(timeString string, defaultTime *time.Time) (time.Time, error) {
	if timeString == "" {
		if defaultTime != nil {
			return *defaultTime, nil
		}
		return time.Time{}, fmt.Errorf("time parameter not provided")
	}
	time, err := parseUnixTime(timeString)
	if err == nil {
		return time, nil
	}
	return parseRFCTime(timeString)
}

func parseUnixTime(timeString string) (time.Time, error) {
	timeNum, err := strconv.ParseFloat(timeString, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(int64(timeNum), 0), nil
}

func parseRFCTime(timeString string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeString)
}

func ParseDuration(durationString, defaultDuration string) (time.Duration, error) {
	if durationString == "" {
		durationString = defaultDuration
	}
	return time.ParseDuration(durationString)
}
