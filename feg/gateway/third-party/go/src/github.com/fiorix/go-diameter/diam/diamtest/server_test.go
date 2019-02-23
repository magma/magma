// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diamtest

import "testing"

func TestNewServer(t *testing.T) {
	srv := NewServer(nil, nil)
	srv.Close()
}

func TestNewServerTLS(t *testing.T) {
	srv := NewUnstartedServer(nil, nil)
	srv.StartTLS()
	srv.Close()
}

func TestNewServerSCTP(t *testing.T) {
	srv := NewServerNetwork("sctp", nil, nil)
	srv.Close()
}
