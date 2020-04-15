/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 */
// package service_manager defines and implements API for service management
package service_manager

import (
	"log"
	"strings"

	"magma/gateway/config"
)

var (
	registry = map[string]ServiceController{
		strings.ToLower(DockerController{}.Name()):  DockerController{},
		strings.ToLower(SystemdController{}.Name()): SystemdController{},
		strings.ToLower(RunitController{}.Name()):   RunitController{},
	}
	defaultController = DockerController{}
)

// Get returns Service Controller for configured init system or the default controller if match cannot be found
func Get() ServiceController {
	initSystem := strings.ToLower(config.GetMagmadConfigs().InitSystem)
	if contr, ok := registry[initSystem]; ok {
		return contr
	}
	log.Printf("process controller for '%s' cannot be found, using '%s' controller",
		initSystem, defaultController.Name())
	return defaultController
}
