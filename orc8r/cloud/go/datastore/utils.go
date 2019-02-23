/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package datastore

// We create a table for each network for all serivces.
// This utility function provides the table name to use.
func GetTableName(networkId string, store string) string {
	return networkId + "_" + store
}
