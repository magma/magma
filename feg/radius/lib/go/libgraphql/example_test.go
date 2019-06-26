/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package libgraphql

import (
	"log"
	"os"
)

// get ACCESS_TOKEN from "ixpp <partner short name>"
func ExampleClient_Do() {
	c := NewClient(ClientConfig{
		Token:    os.Getenv("ACCESS_TOKEN"),
		Endpoint: "https://graph.expresswifi.com/graphql",
	})
	op := NewUpsertCustomer(&AppCustomer{
		MobileNumber: "12311728371117",
	})
	if err := c.Do(op); err != nil {
		log.Fatalf("failed executing graphql request: %v", err)
	}
	log.Printf("graphql response: %v", op.Response())
}
