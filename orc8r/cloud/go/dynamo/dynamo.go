/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dynamo

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// BatchWriteRequests batches the given collection of WriteRequests into a slice of WriteRequest
// slices of the requested batch size.
// IMPORTANT: dynamoDB caps batch sizes at 25 items
func BatchWriteRequests(requests []*dynamodb.WriteRequest, batchSize int) ([][]*dynamodb.WriteRequest, error) {
	if batchSize < 1 {
		return nil, fmt.Errorf("Batch size must be at least 1")
	}

	batchCount := int(math.Ceil(float64(len(requests)) / float64(batchSize)))
	ret := make([][]*dynamodb.WriteRequest, 0, batchCount)

	batchStartIdx := 0
	for batchStartIdx < len(requests) {
		batchEndIdx := int(math.Min(float64(len(requests)), float64(batchStartIdx+batchSize)))
		ret = append(ret, requests[batchStartIdx:batchEndIdx])
		batchStartIdx += batchSize
	}
	return ret, nil
}

// ShouldInitTables returns whether DynamoDB tables should be initialized
func ShouldInitTables() bool {
	// Will be renamed to SHOULD_INIT_DYNAMO_TABLES
	env := os.Getenv("SHOULD_INIT_METERING_TABLES")
	return strings.ToLower(env) == "true"
}

// GetAWSSession returns the AWS session for this machine
func GetAWSSession() (*session.Session, error) {
	endpoint, defined := os.LookupEnv("AWS_ENDPOINT")
	region := os.Getenv("DYNAMO_REGION")

	if defined {
		return session.NewSession(&aws.Config{
			Endpoint:                      aws.String(endpoint),
			Region:                        aws.String(region),
			CredentialsChainVerboseErrors: aws.Bool(true),
		})
	} else {
		return session.NewSession(&aws.Config{
			Region:                        aws.String(region),
			CredentialsChainVerboseErrors: aws.Bool(true),
		})
	}
}
