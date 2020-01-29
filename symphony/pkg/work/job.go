// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package work

import (
	"encoding/json"

	"go.uber.org/zap/zapcore"
)

// Args are the arguments passed into a job.
type Args map[string]interface{}

// String implements fmt.Stringer interface.
func (a Args) String() string {
	data, _ := json.Marshal(a)
	return string(data)
}

// Job to be processed by a Worker.
type Job struct {
	// Handler name to be executed by the worker.
	Handler string `json:"handler"`
	// Args that will be passed to Handler.
	Args Args `json:"args"`
}

// String implements fmt.Stringer interface.
func (j Job) String() string {
	data, _ := json.Marshal(j)
	return string(data)
}

// MarshalLogObject implement zapcore.ObjectMarshaler interface.
func (j Job) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("handler", j.Handler)
	enc.AddString("args", j.Args.String())
	return nil
}
