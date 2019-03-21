/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package crypto

import "fmt"

// MockRNG yields a constant byte sequence instead of generating a new random sequence each time.
type MockRNG struct {
	rand []byte
}

func (rng MockRNG) Read(b []byte) (int, error) {
	copy(b, rng.rand)

	if len(b) <= len(rng.rand) {
		return len(b), nil
	}
	return len(rng.rand), fmt.Errorf("not enough data to read")
}

// NewMockMilenageCipher instantiates the Milenage algo using MockRNG for rng.
func NewMockMilenageCipher(amf []byte, rand []byte) (*MilenageCipher, error) {
	milenage, err := NewMilenageCipher(amf)
	if err != nil {
		return nil, err
	}
	milenage.rng = MockRNG{rand: rand}
	return milenage, nil
}
