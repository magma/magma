// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package relay

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCursorEncoding(t *testing.T) {
	t.Run("EncodeDecode", func(t *testing.T) {
		var buf bytes.Buffer
		c := Cursor{ID: "42", Offset: 15}
		c.MarshalGQL(&buf)
		var uc Cursor
		s, err := strconv.Unquote(buf.String())
		assert.NoError(t, err)
		err = uc.UnmarshalGQL(s)
		assert.NoError(t, err)
		assert.Equal(t, c, uc)
	})
	t.Run("DecodeBadInput", func(t *testing.T) {
		inputs := []interface{}{
			0xbadbeef,
			"cursor@bad123",
			"Y3Vyc29yQGJhZDEyMw==",
		}
		for _, input := range inputs {
			var c Cursor
			err := c.UnmarshalGQL(input)
			assert.Error(t, err)
		}
	})
}
