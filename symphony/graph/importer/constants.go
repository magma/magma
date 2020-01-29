// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

// ReturnMessageCode is a type for codes to be returned to client
type ReturnMessageCode int

const (
	// SuccessfullyUploaded code for successful upload
	SuccessfullyUploaded ReturnMessageCode = iota
	// FailedToUpload code for fail in upload
	FailedToUpload
)
