/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package migration

import (
	"io/ioutil"
	"os"
	"plugin"
	"strings"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// ConfiguratorMigrationPlugin defines the surface area for modules to
// integrate into the configurator data migration.
type ConfiguratorMigrationPlugin interface {
	GetConfigMigrators() []ConfigMigrator
}

// ConfigMigrator is responsible for deserializing a legacy config and
// reserializing it into a configurator config
type ConfigMigrator interface {
	GetType() string

	ToNewConfig(oldConfig []byte) ([]byte, error)
}

const defaultPluginDir = "/var/opt/magma/plugins/migrations/m003_configurator"
const factoryFunction = "GetPlugin"

// allow override for local testing
func getPluginDir() string {
	override, found := os.LookupEnv("PLUGIN_DIR_OVERRIDE")
	if found {
		return override
	}
	return defaultPluginDir
}

func LoadPlugins() error {
	_, err := os.Stat(getPluginDir())
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
	}

	files, err := ioutil.ReadDir(getPluginDir())
	if err != nil {
		return errors.Wrap(err, "could not read plugin dir")
	}

	for _, file := range files {
		isPlugin := strings.HasSuffix(file.Name(), ".so") && !file.IsDir()
		if !isPlugin {
			glog.Infof("not loading file %s because it does not appear to be a valid plugin", file.Name())
			continue
		}

		p, err := plugin.Open(getPluginDir() + file.Name())
		if err != nil {
			return errors.Wrapf(err, "could not open plugin %s", file.Name())
		}
		factory, err := p.Lookup(factoryFunction)
		if err != nil {
			return errors.Errorf("no factory function %s found for %s", factoryFunction, file.Name())
		}
		castedFactory, ok := factory.(func() ConfiguratorMigrationPlugin)
		if !ok {
			return errors.Errorf("expected func() ConfiguratorMigrationPlugin, got %T", factory)
		}

		plug := castedFactory()
		for _, migrator := range plug.GetConfigMigrators() {
			if _, exists := migratorRegistry[migrator.GetType()]; exists {
				return errors.Errorf("plugin %s: migrator with type %s already exists", file.Name(), migrator.GetType())
			}

			migratorRegistry[migrator.GetType()] = migrator
		}
	}
	return nil
}
