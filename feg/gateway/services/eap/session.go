/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSDstyle license found in the
LICENSE file in the root directory of this source tree.
*/

package eap

import (
	"fmt"
	"math/rand"
	"time"
)

// CreateSessionId creates & returns unique session ID string
func CreateSessionId() string {
	return fmt.Sprintf("%X-%X", time.Now().UnixNano()>>16, rand.Uint32())
}
