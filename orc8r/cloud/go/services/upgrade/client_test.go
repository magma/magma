/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package upgrade_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	platform_protos "magma/orc8r/cloud/go/protos"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"
	"magma/orc8r/cloud/go/services/upgrade"
	"magma/orc8r/cloud/go/services/upgrade/protos"
	upgrade_test_init "magma/orc8r/cloud/go/services/upgrade/test_init"
)

func TestUpgradeServiceClientMethods(t *testing.T) {
	magmad_test_init.StartTestService(t)
	upgrade_test_init.StartTestService(t)

	initialProdChannel := &protos.ReleaseChannel{
		SupportedVersions: []string{"1.0.0-0", "1.1.0-0"},
	}
	err := upgrade.CreateReleaseChannel("stable", initialProdChannel)
	assert.NoError(t, err)

	actualChannel, err := upgrade.GetReleaseChannel("stable")
	assert.NoError(t, err)
	assert.Equal(t, platform_protos.TestMarshal(initialProdChannel), platform_protos.TestMarshal(actualChannel))
	_, err = upgrade.GetReleaseChannel("nochannel")
	assert.Error(t, err)

	updatedProdChannel := &protos.ReleaseChannel{
		SupportedVersions: []string{"1.1.0-0", "1.2.0-0"},
	}
	err = upgrade.UpdateReleaseChannel("stable", updatedProdChannel)
	assert.NoError(t, err)

	actualChannel, err = upgrade.GetReleaseChannel("stable")
	assert.NoError(t, err)
	assert.Equal(t, platform_protos.TestMarshal(updatedProdChannel), platform_protos.TestMarshal(actualChannel))

	err = upgrade.UpdateReleaseChannel("nochannel", updatedProdChannel)
	assert.Error(t, err)

	err = upgrade.CreateReleaseChannel(
		"delete",
		&protos.ReleaseChannel{SupportedVersions: []string{"1.0.0-0"}})
	assert.NoError(t, err)
	err = upgrade.DeleteReleaseChannel("delete")
	assert.NoError(t, err)
	_, err = upgrade.GetReleaseChannel("delete")
	assert.Error(t, err)

	err = upgrade.CreateTier("network", "t1", &protos.TierInfo{Name: "t1", Version: "1.1.0-0"})
	assert.NoError(t, err)
	err = upgrade.CreateTier("network", "t2", &protos.TierInfo{Name: "t2", Version: "1.2.0-0"})
	assert.NoError(t, err)

	tiers, err := upgrade.GetTiers("network", []string{})
	assert.Equal(t, 2, len(tiers))
	assert.Equal(t, "1.1.0-0", tiers["t1"].Version)
	assert.Equal(t, "1.2.0-0", tiers["t2"].Version)
	tiers, err = upgrade.GetTiers("network", []string{"t1"})
	assert.Equal(t, 1, len(tiers))
	assert.Equal(t, "1.1.0-0", tiers["t1"].Version)

	updatedT1 := &protos.TierInfo{Name: "t1v2", Version: "1.3.0-0"}
	err = upgrade.UpdateTier("network", "t1", updatedT1)
	assert.NoError(t, err)
	tiers, err = upgrade.GetTiers("network", []string{"t1"})
	assert.Equal(t, 1, len(tiers))
	assert.Equal(t, "1.3.0-0", tiers["t1"].Version)
	assert.Equal(t, "t1v2", tiers["t1"].Name)

	err = upgrade.DeleteTier("network", "t2")
	assert.NoError(t, err)
	tiers, err = upgrade.GetTiers("network", []string{})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tiers))
	assert.Equal(t, "1.3.0-0", tiers["t1"].Version)
}
