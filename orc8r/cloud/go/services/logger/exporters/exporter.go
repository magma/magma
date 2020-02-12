/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package exporters

import "magma/orc8r/lib/go/protos"

type Exporter interface {
	// export logEntries
	Submit(logEntries []*protos.LogEntry) error
}
