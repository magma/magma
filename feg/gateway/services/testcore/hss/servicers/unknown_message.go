/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"log"

	"github.com/fiorix/go-diameter/v4/diam"
)

// handleUnknownMessage is called when a diameter message is received with a
// code that we have not registered a handler for.
func handleUnknownMessage(_ diam.Conn, msg *diam.Message) {
	log.Printf("Unhandled diameter message with command code: %v", msg.Header.CommandCode)
}
