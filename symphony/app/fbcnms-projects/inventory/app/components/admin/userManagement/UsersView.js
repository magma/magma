/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {EditUserMutationResponse} from '../../../mutations/__generated__/EditUserMutation.graphql';
import type {MutationCallbacks} from '../../../mutations/MutationCallbacks.js';
import type {User} from './TempTypes';

import * as React from 'react';
import EditUserMutation from '../../../mutations/EditUserMutation';
import InventoryQueryRenderer from '../../InventoryQueryRenderer';
import UsersTable from './UsersTable';
import {graphql} from 'relay-runtime';

const usersQuery = graphql`
  query UsersView_UsersQuery {
    users {
      edges {
        node {
          id
          authID
          firstName
          lastName
          email
          status
          role
          profilePhoto {
            id
            fileName
            storeKey
          }
        }
      }
    }
  }
`;
type UsersQueryResponse = {
  users: {
    edges: Array<{
      node: User,
    }>,
  },
};

const editUser = (newUserValue: User) => {
  return new Promise<User>((resolve, reject) => {
    const callbacks: MutationCallbacks<EditUserMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          reject(errors[0].message);
        }
        resolve({
          id: response.editUser.id,
          authID: response.editUser.authID,
          firstName: response.editUser.firstName,
          lastName: response.editUser.lastName,
          role: response.editUser.role,
          status: response.editUser.status,
        });
      },
      onError: () => {
        reject('Error saving service');
      },
    };
    EditUserMutation(
      {
        input: {
          id: newUserValue.id,
          firstName: newUserValue.firstName,
          lastName: newUserValue.lastName,
          role: newUserValue.role,
          status: newUserValue.status,
        },
      },
      callbacks,
    );
  });
};

export default function UsersView() {
  return (
    <InventoryQueryRenderer
      query={usersQuery}
      variables={{}}
      render={(respons: UsersQueryResponse) => {
        const users: Array<User> = respons.users.edges.map(user => ({
          id: user.node.id,
          authID: user.node.authID,
          firstName: user.node.firstName,
          lastName: user.node.lastName,
          role: user.node.role,
          status: user.node.status,
        }));
        return <UsersTable users={users} onUseredit={editUser} />;
      }}
    />
  );
}
