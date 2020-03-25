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
import Tokenizer from '@fbcnms/ui/components/Tokenizer';
import axios from 'axios';

const PowerSearchWorkOrderGeneralUserFilter = (props: FilterProps) => {
  const [users, setUsers] = useState([]);
  const [searchEntries, setSearchEntries] = useState([]);
  const [tokens, setTokens] = useState([]);

  useEffect(() => {
    axios.get('/user/list/').then(response => setUsers(response.data.users));
  }, []);

  const {
    value,
    title,
    onInputBlurred,
    onValueChanged,
    onRemoveFilter,
    editMode,
  } = props;

  return (
    <PowerSearchFilter
      name={title ?? ''}
      operator={'is_one_of'}
      editMode={editMode}
      value={(value.stringSet ?? [])
        .map(email => users.find(users => users.email === email)?.email)
        .join(', ')}
      onRemoveFilter={onRemoveFilter}
      input={
        <Tokenizer
          searchSource="Options"
          tokens={tokens}
          onEntriesRequested={searchTerm =>
            setSearchEntries(
              users
                .filter(user =>
                  user.email.toLowerCase().includes(searchTerm.toLowerCase()),
                )
                .map(user => ({id: user.email, label: user.email})),
            )
          }
          searchEntries={searchEntries}
          onBlur={onInputBlurred}
          onChange={newEntries => {
            setTokens(newEntries);
            onValueChanged({
              id: value.id,
              key: value.key,
              name: value.name,
              operator: value.operator,
              stringSet: newEntries.map(entry => entry.id),
            });
          }}
        />
      }
    />
  );
};

export default PowerSearchWorkOrderGeneralUserFilter;
