/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package cmd

import (
	"context"
	"fmt"
	"log"

	"magma/orc8r/cloud/go/services/state/indexer/reindex"
	"magma/orc8r/cloud/go/services/state/protos"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all indexer IDs and versions",
	Run:   runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listID, "id", "i", "", "restrict to specific indexer ID")
	listCmd.Flags().BoolVarP(&listShort, "name", "n", false, "only print indexer IDs")
}

func runList(cmd *cobra.Command, args []string) {
	res, err := getClient().GetIndexers(context.Background(), &protos.GetIndexersRequest{})
	if err != nil {
		log.Fatal(err)
	}

	if listID != "" {
		v, ok := res.IndexersById[listID]
		if !ok {
			log.Fatalf("No indexer found for ID %s", listID)
		}
		printVersions(protos.MakeVersion(v))
		return
	}

	printVersions(protos.MakeVersions(res.IndexersById)...)
}

func printVersions(versions ...*reindex.Version) {
	for _, v := range versions {
		if listShort {
			fmt.Println(v.IndexerID)
			continue
		}
		fmt.Println(v)
	}
}
