/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package testcontroller

import (
	"context"

	"magma/orc8r/cloud/go/serde"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/registry"
	"orc8r/fbinternal/cloud/go/services/testcontroller/protos"
	"orc8r/fbinternal/cloud/go/services/testcontroller/statemachines"
	"orc8r/fbinternal/cloud/go/services/testcontroller/storage"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// UnmarshalledTestCase encapsulates a TestCase and the unmarshaled
// representation of its TestConfig field.
type UnmarshalledTestCase struct {
	*storage.TestCase
	UnmarshaledConfig interface{}
}

func getE2EClient() (protos.TestControllerClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewTestControllerClient(conn), nil
}

func ExecuteNextTestCase(testMachines map[string]statemachines.TestMachine, store storage.TestControllerStorage) error {
	tc, err := store.GetNextTestForExecution()
	if err != nil {
		return err
	}
	if tc == nil {
		return nil
	}

	machine, ok := testMachines[tc.TestCaseType]
	if !ok {
		return errors.Errorf("no test state machine found matching %s", tc.TestCaseType)
	}
	unmarshalledConfig, err := serde.Deserialize(SerdeDomain, tc.TestCaseType, tc.TestConfig)
	if err != nil {
		return errors.Wrapf(err, "could not deserialize test %s config", tc.TestCaseType)
	}
	var prevErr error
	if tc.Error != "" {
		prevErr = errors.New(tc.Error)
	}

	nextState, nextDuration, err := machine.Run(tc.State, unmarshalledConfig, prevErr)
	var newErr *string
	if err != nil {
		newErr = strPtr(err.Error())
	}
	err = store.ReleaseTest(tc.Pk, nextState, newErr, nextDuration)
	if err != nil {
		return err
	}
	return nil
}

func strPtr(s string) *string {
	return &s
}

func GetTestCases(pks []int64) (map[int64]*UnmarshalledTestCase, error) {
	client, err := getE2EClient()
	if err != nil {
		return nil, err
	}
	res, err := client.GetTestCases(context.Background(), &protos.GetTestCasesRequest{Pks: pks})
	if err != nil {
		return nil, err
	}

	ret := map[int64]*UnmarshalledTestCase{}
	for pk, tc := range res.Tests {
		unmarshalledConfig, err := serde.Deserialize(SerdeDomain, tc.TestCaseType, tc.TestConfig)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to deserialize test case of type %s", tc.TestCaseType)
		}
		ret[pk] = &UnmarshalledTestCase{
			TestCase:          tc,
			UnmarshaledConfig: unmarshalledConfig,
		}
	}
	return ret, nil
}

func CreateOrUpdateTestCase(pk int64, testCaseType string, testCaseConfig interface{}) error {
	marshaledConfig, err := serde.Serialize(SerdeDomain, testCaseType, testCaseConfig)
	if err != nil {
		return errors.Wrap(err, "failed to serialize config")
	}

	client, err := getE2EClient()
	if err != nil {
		return err
	}
	_, err = client.CreateOrUpdateTestCase(
		context.Background(),
		&protos.CreateTestCaseRequest{
			Test: &storage.MutableTestCase{
				Pk:           pk,
				TestCaseType: testCaseType,
				TestConfig:   marshaledConfig,
			},
		},
	)
	return err
}

func DeleteTestCase(pk int64) error {
	client, err := getE2EClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteTestCase(context.Background(), &protos.DeleteTestCaseRequest{Pk: pk})
	return err
}
