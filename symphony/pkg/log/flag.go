// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	// LevelFlagName is the canonical flag name to configure the allowed log level.
	LevelFlagName = "log.level"
	// LevelFlagEnvVar is the environment variable name to configure the allowed log level.
	LevelFlagEnvar = "LOG_LEVEL"
	// LevelFlagHelp is the help description for the log.level flag.
	LevelFlagHelp = "Only log messages with the given severity or above. One of: [debug, info, warn, error]"
	// FormatFlagName is the canonical flag name to configure the log format.
	FormatFlagName = "log.format"
	// FormatFlagEnvar is the environment variable name to configure the allowed log format.
	FormatFlagEnvar = "LOG_FORMAT"
	// FormatFlagHelp is the help description for the log.format flag.
	FormatFlagHelp = "Output format of log messages. One of: [console, json]"
)

// AddFlagsVar adds the flags used by this package to the Kingpin application.
func AddFlagsVar(a *kingpin.Application, config *Config) {
	a.Flag(LevelFlagName, LevelFlagHelp).
		Envar(LevelFlagEnvar).
		Default("info").
		SetValue(&config.Level)
	a.Flag(FormatFlagName, FormatFlagHelp).
		Envar(FormatFlagEnvar).
		Default("console").
		SetValue(&config.Format)
}

// AddFlags adds the flags used by this package to the Kingpin application.
func AddFlags(a *kingpin.Application) *Config {
	config := &Config{}
	AddFlagsVar(a, config)
	return config
}
