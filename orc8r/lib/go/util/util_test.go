package util_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/orc8r/lib/go/util"
)

func TestEnv(t *testing.T) {
	testEnv := "This_test_env_Should_Not_be_SET"
	os.Unsetenv(testEnv)

	assert.True(t, util.GetEnvBool(testEnv, true))
	assert.False(t, util.GetEnvBool(testEnv))
	assert.False(t, util.GetEnvBool(testEnv, false))
	os.Setenv(testEnv, "1")

	assert.True(t, util.GetEnvBool(testEnv), "Env value: '%s'", os.Getenv(testEnv))
	os.Setenv(testEnv, "0")
	assert.False(t, util.GetEnvBool(testEnv))
	os.Setenv(testEnv, "true")
	assert.True(t, util.GetEnvBool(testEnv), "Env value: '%s'", os.Getenv(testEnv))
	os.Setenv(testEnv, "false")
	assert.False(t, util.GetEnvBool(testEnv))
}
