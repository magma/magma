// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package models

import (
	"bytes"
	"encoding/base64"
	"io"
	"strconv"

	"golang.org/x/xerrors"
)

// Cursor is a graphql pagination cursor.
type Cursor int

// cursor encoding prefix.
const cursorPrefix = "cursor@"

// MarshalGQL implements graphql.Marshaler interface.
func (c Cursor) MarshalGQL(w io.Writer) {
	s := base64.StdEncoding.EncodeToString(
		[]byte(cursorPrefix + strconv.Itoa(c.Int())),
	)
	_, _ = io.WriteString(w, strconv.Quote(s))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (c *Cursor) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return xerrors.Errorf("%T is not a string", v)
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return xerrors.Errorf("decoding cursor: %w", err)
	}
	b = bytes.TrimPrefix(b, []byte(cursorPrefix))
	i, err := strconv.Atoi(string(b))
	if err != nil {
		return xerrors.Errorf("extracting cursor value: %w", err)
	}
	*c = Cursor(i)
	return nil
}

// Int returns the int value stored in cursor.
func (c Cursor) Int() int {
	return int(c)
}
