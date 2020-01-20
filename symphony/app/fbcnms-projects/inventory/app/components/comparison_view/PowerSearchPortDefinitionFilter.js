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

const portDefinitionsQuery = graphql`
  query PowerSearchPortDefinitionFilterQuery {
    equipmentPortDefinitions {
      edges {
        node {
          id
          name
        }
      }
    }
  }
`;

const PowerSearchPortDefinitionFilter = (props: FilterProps) => {
  const [portDefinitions, setPortDefinitions] = useState([]);
  const [searchEntries, setSearchEntries] = useState([]);
  const [tokens, setTokens] = useState([]);

  useEffect(() => {
    fetchQuery(RelayEnvironment, portDefinitionsQuery).then(data => {
      setPortDefinitions(
        data.equipmentPortDefinitions.edges.map(edge => edge.node),
      );
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
      name="Port Names"
      operator={value.operator}
      editMode={editMode}
      value={(value.idSet ?? [])
        .map(id => portDefinitions.find(type => type.id === id)?.name)
        .join(', ')}
      onRemoveFilter={onRemoveFilter}
      input={
        <Tokenizer
          searchSource="Options"
          tokens={tokens}
          onEntriesRequested={searchTerm =>
            setSearchEntries(
              portDefinitions
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

export default PowerSearchPortDefinitionFilter;
