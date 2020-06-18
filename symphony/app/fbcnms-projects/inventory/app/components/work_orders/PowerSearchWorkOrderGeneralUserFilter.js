/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

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
    userSearch(limit: 10, filters: $filters) {
      users {
        id
        email
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

const PowerSearchWorkOrderGeneralUserFilter = (props: FilterProps) => {
  const {
    value,
    title,
    onInputBlurred,
    onValueChanged,
    onRemoveFilter,
    editMode,
  } = props;
  const [selectedUsers, setSelectedUsers] = useState([]);
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
        (data.userSearch.users ?? []).map(user => ({
          id: user.id,
          label: user.email,
        })),
      );
    });

  useEffect(() => {
    value.idSet
      ?.filter(id => !selectedUsers.find(l => l.id == id))
      .map(id =>
        fetchQuery(RelayEnvironment, userQuery, {id: id}).then(user =>
          setSelectedUsers([
            ...selectedUsers,
            {id: user.node.id, label: user.node.email},
          ]),
        ),
      );
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
