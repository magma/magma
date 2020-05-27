// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestProvider(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	logger, restorer, err := ProvideLogger(Config{})
	require.NoError(t, err)
	defer restorer()
	assert.Equal(t, logger.Background(), ProvideZapLogger(logger))
	assert.Equal(t, logger.Background(), zap.L())
	log.Println("suppressed message")
	assert.Zero(t, buf.Len())
}
