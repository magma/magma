/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FilterProps} from './ComparisonViewTypes';

import PowerSearchFilter from './PowerSearchFilter';
import React, {useEffect, useState} from 'react';
import RelayEnvironment from '../../common/RelayEnvironment';
import Tokenizer from '@fbcnms/ui/components/Tokenizer';
import {fetchQuery, graphql} from 'relay-runtime';

const locationTypesQuery = graphql`
  query PowerSearchLocationTypeFilterQuery {
    locationTypes {
      edges {
        node {
          id
          name
        }
      }
    }
  }
`;

const PowerSearchLocationTypeFilter = (props: FilterProps) => {
  const [locationTypes, setLocationTypes] = useState([]);
  const [searchEntries, setSearchEntries] = useState([]);
  const [tokens, setTokens] = useState([]);

  useEffect(() => {
    // $FlowFixMe (T62907961) Relay flow types
    fetchQuery(RelayEnvironment, locationTypesQuery).then(data => {
      setLocationTypes(data.locationTypes.edges.map(edge => edge.node));
    });
  }, []);

  const {
    value,
    onInputBlurred,
    onValueChanged,
    onRemoveFilter,
    editMode,
  } = props;
  return (
    <PowerSearchFilter
      name="Location Type"
      operator={value.operator}
      editMode={editMode}
      value={(value.idSet ?? [])
        .map(id => locationTypes.find(type => type.id === id)?.name)
        .join(', ')}
      onRemoveFilter={onRemoveFilter}
      input={
        <Tokenizer
          searchSource="Options"
          tokens={tokens}
          onEntriesRequested={searchTerm =>
            setSearchEntries(
              locationTypes
                .filter(type =>
                  type.name.toLowerCase().includes(searchTerm.toLowerCase()),
                )
                .map(type => ({id: type.id, label: type.name})),
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
              idSet: newEntries.map(entry => entry.id),
            });
          }}
        />
      }
    />
  );
};

export default PowerSearchLocationTypeFilter;
