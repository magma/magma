/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package storage

import (
	"sort"

	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type internalEntityGraph struct {
	entsByPk map[string]*NetworkEntity
	edges    []loadedAssoc
}

// loadGraphInternal will load all entities and assocs for a given graph ID.
// This function will NOT fill entities with associations.
func (store *sqlConfiguratorStorage) loadGraphInternal(networkID string, graphID string, criteria EntityLoadCriteria) (internalEntityGraph, error) {
	loadFilter := EntityLoadFilter{graphID: &graphID}
	entsByPk, err := store.loadFromEntitiesTable(networkID, loadFilter, criteria)
	if err != nil {
		return internalEntityGraph{}, errors.Wrap(err, "failed to load entities for graph")
	}

	// always load all edges for a graph load
	criteria.LoadAssocsFromThis, criteria.LoadAssocsToThis = true, true
	assocs, _, err := store.loadFromAssocsTable(loadFilter, criteria, entsByPk)
	if err != nil {
		return internalEntityGraph{}, errors.Wrap(err, "failed to load edges for graph")
	}

	return internalEntityGraph{entsByPk: entsByPk, edges: assocs}, nil
}

// fixGraph will load a a graph which may have been partitioned, do a connected
// component search on it, and relabel components if a partition is detected.
// entToUpdateOut is an output parameter
func (store *sqlConfiguratorStorage) fixGraph(graphID string, entToUpdateOut *entWithPk) error {
	// TODO: implement after LoadGraphForEntity

	return nil
}

func findRootNodes(graph internalEntityGraph) []string {
	// A root node is a node with no edges terminating at it.
	// To make computation easier, we'll make a reverse adjacency list first
	reverseEdges := map[string][]string{}
	for _, edge := range graph.edges {
		reverseEdges[edge.toPk] = append(reverseEdges[edge.toPk], edge.fromPk)
	}

	// Set of root nodes is the set difference between node IDs and reverse
	// edge keyset.
	rootNodes := funk.Filter(
		funk.Keys(graph.entsByPk),
		func(pk string) bool {
			_, hasEdgesTo := reverseEdges[pk]
			return !hasEdgesTo
		},
	).([]string)
	sort.Strings(rootNodes)
	return rootNodes
}
