/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_utils

import (
	"math/rand"
	"strings"
	"time"

	"magma/orc8r/cloud/go/services/magmad"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Generate a random network ID containing only lowercase letters
// If the generated ID is already registered, retry
func GenerateNetworkId() (string, error) {
	const (
		IdRandLen       = 24
		MaxIdGenRetries = 10
	)

	for retryCount := 0; retryCount < MaxIdGenRetries; retryCount++ {
		generatedId := randomString(IdRandLen)
		networkRecord, err := magmad.GetNetwork(generatedId)

		if networkRecord == nil && err != nil && strings.Contains(err.Error(), "No record") {
			return generatedId, nil
		} else if networkRecord == nil {
			return "", err
		}
	}
	return "", status.Errorf(codes.Internal, "Was not able to generate a unique network ID in %d attempts.", MaxIdGenRetries)
}

func randomString(strlen int) string {
	const (
		Alphanum    = "abcdefghijklmnopqrstuvwxyz"
		AlphanumLen = len(Alphanum)
	)

	rand.Seed(time.Now().UTC().UnixNano())
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = Alphanum[rand.Intn(AlphanumLen)]
	}
	return string(result)
}
