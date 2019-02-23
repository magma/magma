// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package sm provides diameter state machines for clients and servers.
//
// It currently handles CER/CEA handshakes, and automatic DWR/DWA. Peers
// that pass the handshake get metadata associated to their connection.
// See the peer sub-package for details on the metadata.
package sm
