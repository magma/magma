// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sign

import (
	"context"
	"errors"
	"time"
)

// Operation defines a signing operation.
type Operation int

const (
	// Invalid represents an invalid operation.
	Invalid Operation = iota

	// GetObject defines an object read operation.
	GetObject

	// PutObject defines an object write operation.
	PutObject

	// DeleteObject defines an object delete operation.
	DeleteObject

	// DownloadObject defines an object download operation.
	DownloadObject
)

type (
	// Signer defines the interface which pre-signs object operations.
	Signer interface {
		Sign(context.Context, Operation, string, string) (string, error)
	}

	// NopSigner is a signer that always fails to sign.
	NopSigner struct{}
)

// Sign implements Signer interface.
func (NopSigner) Sign(context.Context, Operation, string, string, time.Duration) (string, error) {
	return "", errors.New("not implemented")
}
