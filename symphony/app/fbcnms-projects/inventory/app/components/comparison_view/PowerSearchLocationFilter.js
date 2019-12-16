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
import React, {useState} from 'react';
import RelayEnvironment from '../../common/RelayEnvironment';
import Tokenizer from '@fbcnms/ui/components/Tokenizer';
import nullthrows from '@fbcnms/util/nullthrows';
import {fetchQuery, graphql} from 'relay-runtime';

const locationTokenizerQuery = graphql`
  query PowerSearchLocationFilterQuery($name: String!, $types: [ID!]) {
    locations(name: $name, first: 10, types: $types) {
      edges {
        node {
          id
          name
        }
      }
    }
  }
`;

const PowerSearchLocationFilter = (props: FilterProps) => {
  const {
    config,
    value,
    onInputBlurred,
    onValueChanged,
    onRemoveFilter,
    editMode,
  } = props;

  const [selectedLocations, setSelectedLocations] = useState([]);
  const [searchEntries, setSearchEntries] = useState([]);

  const fetchLocations = searchTerm =>
    fetchQuery(RelayEnvironment, locationTokenizerQuery, {
      name: searchTerm,
      types: [nullthrows(config.extraData).locationTypeId],
    }).then(data => {
      setSearchEntries(
        (data.locations.edges ?? [])
          .map(edge => edge.node)
          .map(location => ({id: location.id, label: location.name})),
      );
    });

  return (
    <PowerSearchFilter
      name={config.label}
      operator={value.operator}
      editMode={editMode}
      value={selectedLocations.map(location => location.label).join(', ')}
      onRemoveFilter={onRemoveFilter}
      input={
        <Tokenizer
          searchSource="Options"
          tokens={selectedLocations}
          onEntriesRequested={searchTerm => fetchLocations(searchTerm)}
          searchEntries={searchEntries}
          onBlur={onInputBlurred}
          onChange={newEntries => {
            setSelectedLocations(newEntries);
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

export default PowerSearchLocationFilter;
