/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  AddWorkOrderTypeMutation,
  AddWorkOrderTypeMutationResponse,
  AddWorkOrderTypeMutationVariables,
} from './__generated__/AddWorkOrderTypeMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';
import type {WorkOrderType} from '../common/WorkOrder';

import RelayEnvironment from '../common/RelayEnvironment.js';
import {ConnectionHandler} from 'relay-runtime';
import {commitMutation, graphql} from 'react-relay';
import {convertPropertyTypeToMutationInput} from '../common/PropertyType';

const mutation = graphql`
  mutation AddWorkOrderTypeMutation($input: AddWorkOrderTypeInput!) {
    addWorkOrderType(input: $input) {
      id
      name
      description
      propertyTypes {
        id
        name
        type
        nodeType
        index
        stringValue
        intValue
        booleanValue
        floatValue
        latitudeValue
        longitudeValue
        rangeFromValue
        rangeToValue
        isEditable
        isMandatory
        isInstanceProperty
        isDeleted
        category
      }
    }
  }
`;

export const addWorkOrderType = (
  workOrderType: WorkOrderType,
): Promise<AddWorkOrderTypeMutationResponse> => {
  const {name, description, propertyTypes} = workOrderType;
  const variables: AddWorkOrderTypeMutationVariables = {
    input: {
      name,
      description,
      properties: convertPropertyTypeToMutationInput(propertyTypes),
    },
  };

  return new Promise((resolve, reject) => {
    const callbacks: MutationCallbacks<AddWorkOrderTypeMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          return reject(errors[0]);
        } else {
          resolve(response);
        }
      },
      onError: reject,
    };
    const updater = store => {
      const rootQuery = store.getRoot();
      const newNode = store.getRootField('addWorkOrderType');
      if (!newNode) {
        return;
      }
      const types = ConnectionHandler.getConnection(
        rootQuery,
        'Configure_workOrderTypes',
      );
      if (types == null) {
        return;
      }
      const edge = ConnectionHandler.createEdge(
        store,
        types,
        newNode,
        'WorkOrderTypesEdge',
      );
      ConnectionHandler.insertEdgeBefore(types, edge);
    };
    CommitWorkOrderTypeMutation(variables, callbacks, updater);
  });
};

const CommitWorkOrderTypeMutation = (
  variables: AddWorkOrderTypeMutationVariables,
  callbacks?: MutationCallbacks<AddWorkOrderTypeMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<AddWorkOrderTypeMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};

export default CommitWorkOrderTypeMutation;
