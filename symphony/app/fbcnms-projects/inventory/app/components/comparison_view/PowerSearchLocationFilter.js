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
import React, {useContext, useEffect, useState} from 'react';
import RelayEnvironment from '../../common/RelayEnvironment';
import Tokenizer from '@fbcnms/ui/components/Tokenizer';
import WizardContext from '@fbcnms/ui/components/design-system/Wizard/WizardContext';
import nullthrows from '@fbcnms/util/nullthrows';
import {fetchQuery, graphql} from 'relay-runtime';
import {useTokens} from './tokensHook';

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

const locationQuery = graphql`
  query PowerSearchLocationFilterIDsQuery($id: ID!) {
    node(id: $id) {
      ... on Location {
        id
        name
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
  const wizardContext = useContext(WizardContext);
  const [searchEntries, setSearchEntries] = useState([]);
  const tokens = useTokens(value);
  const [selectedLocations, setSelectedLocations] = useState(tokens);

  const fetchLocations = searchTerm =>
    fetchQuery(RelayEnvironment, locationTokenizerQuery, {
      name: searchTerm,
      types: [nullthrows(config.extraData).locationTypeId],
    }).then(data => {
      setSearchEntries(
        (data.locations.edges ?? [])
          .map(edge => edge.node)
          .filter(node => !selectedLocations.find(loc => loc.id == node.id))
          .map(location => ({id: location.id, label: location.name})),
      );
    });

  useEffect(() => {
    value.idSet
      ?.filter(id => !selectedLocations.find(l => l.id == id))
      .map(id =>
        fetchQuery(RelayEnvironment, locationQuery, {id: id}).then(location =>
          setSelectedLocations(locations => [
            ...locations,
            {id: location.node.id, label: location.node.name},
          ]),
        ),
      );
  }, [selectedLocations, value.idSet]);

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
            newEntries.map(entry =>
              wizardContext.set(entry.id, {id: entry.id, label: entry.label}),
            );
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
