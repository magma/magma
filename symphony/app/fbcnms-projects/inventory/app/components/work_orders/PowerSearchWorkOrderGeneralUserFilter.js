/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {
  PowerSearchWorkOrderGeneralUserFilterIDsQuery as FetchUserQuery,
  // eslint-disable-next-line max-len
  PowerSearchWorkOrderGeneralUserFilterIDsQueryResponse as FetchUserQueryResponse,
} from './__generated__/PowerSearchWorkOrderGeneralUserFilterIDsQuery.graphql';
import type {FilterProps} from '../comparison_view/ComparisonViewTypes';

import PowerSearchFilter from '../comparison_view/PowerSearchFilter';
import React, {useEffect, useState} from 'react';
import RelayEnvironment from '../../common/RelayEnvironment';
import Tokenizer from '@fbcnms/ui/components/Tokenizer';
import {fetchQuery, graphql} from 'relay-runtime';

const usersQuery = graphql`
  query PowerSearchWorkOrderGeneralUserFilter_userQuery(
    $filters: [UserFilterInput!]!
  ) {
    users(first: 10, filterBy: $filters) {
      edges {
        node {
          id
          email
        }
      }
    }
  }
`;

const userQuery = graphql`
  query PowerSearchWorkOrderGeneralUserFilterIDsQuery($id: ID!) {
    node(id: $id) {
      ... on User {
        id
        email
      }
    }
  }
`;

type SelectedUser = $ReadOnly<{
  id: string,
  label: string,
}>;

const PowerSearchWorkOrderGeneralUserFilter = (props: FilterProps) => {
  const {
    value,
    title,
    onInputBlurred,
    onValueChanged,
    onRemoveFilter,
    editMode,
  } = props;
  const [selectedUsers, setSelectedUsers] = useState<Array<SelectedUser>>([]);
  const [searchEntries, setSearchEntries] = useState([]);

  const fetchUsers = searchTerm =>
    fetchQuery(RelayEnvironment, usersQuery, {
      filters: [
        {
          filterType: 'USER_NAME',
          operator: 'CONTAINS',
          stringValue: searchTerm,
        },
      ],
    }).then(data => {
      setSearchEntries(
        (data.users.edges ?? [])
          .map(edge => edge.node)
          .map(user => ({
            id: user.id,
            label: user.email,
          })),
      );
    });

  useEffect(() => {
    if (value.idSet == null) {
      return;
    }
    const missingUsersPromises = value.idSet
      .filter(id => !selectedUsers.find(l => l.id == id))
      .map<Promise<FetchUserQueryResponse>>(id =>
        fetchQuery<FetchUserQuery>(RelayEnvironment, userQuery, {id: id}),
      );
    Promise.allSettled<Array<Promise<FetchUserQueryResponse>>>(
      missingUsersPromises,
    ).then(userPromises => {
      const fetchedUsers = userPromises.map<?FetchUserQueryResponse>(
        userPromise =>
          userPromise.status === 'fulfilled' ? userPromise.value : null,
      );
      if (fetchedUsers.length === 0) {
        return;
      }
      const missingUsers: $ReadOnlyArray<SelectedUser> = fetchedUsers
        .map(fetchedUser => {
          const user = fetchedUser?.node;
          return user != null && user.id != null && user.email != null
            ? {
                id: user.id,
                label: user.email,
              }
            : null;
        })
        .filter(Boolean);
      setSelectedUsers([...selectedUsers, ...missingUsers]);
    });
  }, [selectedUsers, value.idSet]);

  return (
    <PowerSearchFilter
      name={title ?? ''}
      operator={'is_one_of'}
      editMode={editMode}
      value={selectedUsers.map(user => user.label).join(', ')}
      onRemoveFilter={onRemoveFilter}
      input={
        <Tokenizer
          searchSource="Options"
          tokens={selectedUsers}
          onEntriesRequested={searchTerm => {
            fetchUsers(searchTerm);
          }}
          searchEntries={searchEntries}
          onBlur={onInputBlurred}
          onChange={newEntries => {
            setSelectedUsers(newEntries);
            onValueChanged({
              id: value.id,
              key: value.key,
              name: value.name,
              operator: value.operator,
              idSet: newEntries.map(entry => entry.id),
            });
          }}
        />
      }
    />
  );
};

export default PowerSearchWorkOrderGeneralUserFilter;
