/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package radius

// NonAuthenticResponseError is returned when a client was expecting
// a valid response but did not receive one.
type NonAuthenticResponseError struct {
}

func (e *NonAuthenticResponseError) Error() string {
	return `radius: non-authentic response`
}
