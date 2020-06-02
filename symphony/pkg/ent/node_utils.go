// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.
package ent

import (
	"encoding/json"
	"time"
)

func (f Field) MustGetInt() int {
	var val int
	err := json.Unmarshal([]byte(f.Value), &val)
	if err != nil {
		panic(err)
	}
	return val
}

func (f Field) MustGetUint() uint {
	var val uint
	err := json.Unmarshal([]byte(f.Value), &val)
	if err != nil {
		panic(err)
	}
	return val
}

func (f Field) MustGetFloat32() float32 {
	var val float32
	err := json.Unmarshal([]byte(f.Value), &val)
	if err != nil {
		panic(err)
	}
	return val
}

func (f Field) MustGetFloat64() float64 {
	var val float64
	err := json.Unmarshal([]byte(f.Value), &val)
	if err != nil {
		panic(err)
	}
	return val
}

func (f Field) MustGetTime() time.Time {
	var val time.Time
	err := json.Unmarshal([]byte(f.Value), &val)
	if err != nil {
		panic(err)
	}
	return val
}

func (f Field) MustGetString() string {
	var val string
	err := json.Unmarshal([]byte(f.Value), &val)
	if err != nil {
		panic(err)
	}
	return val
}
