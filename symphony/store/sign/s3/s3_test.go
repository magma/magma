// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package s3

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
)

func TestAddFlags(t *testing.T) {
	a := kingpin.New(t.Name(), "")
	config := AddFlags(a)
	_, err := a.Parse([]string{
		"--s3.bucket", t.Name(),
		"--s3.region", "us-east-1",
		"--s3.endpoint", "localtest.me",
	})
	assert.NoError(t, err)
	assert.Equal(t, t.Name(), config.Bucket)
	assert.Equal(t, "us-east-1", config.Region)
	assert.Equal(t, "localtest.me", config.Endpoint)
	assert.Equal(t, 24*time.Hour, config.Expire)
}
