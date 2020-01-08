// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package relay

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/ugorji/go/codec"
)

// PageInfo of a connection type.
type PageInfo struct {
	HasNextPage     bool   `json:"hasNextPage"`
	HasPreviousPage bool   `json:"hasPreviousPage"`
	StartCursor     Cursor `json:"startCursor"`
	EndCursor       Cursor `json:"endCursor"`
}

// Cursor of a relay edge type.
type Cursor struct {
	ID     string
	Offset int
}

var quote = []byte(`"`)

// MarshalGQL implements graphql.Marshaler interface.
func (c Cursor) MarshalGQL(w io.Writer) {
	w.Write(quote)
	defer w.Write(quote)
	wc := base64.NewEncoder(base64.StdEncoding, w)
	defer wc.Close()
	_ = codec.NewEncoder(wc, &codec.MsgpackHandle{}).Encode(c)
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (c *Cursor) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("%T is not a string", v)
	}
	if err := codec.NewDecoder(
		base64.NewDecoder(
			base64.StdEncoding,
			strings.NewReader(s),
		),
		&codec.MsgpackHandle{},
	).Decode(c); err != nil {
		return fmt.Errorf("decode cursor: %w", err)
	}
	return nil
}
