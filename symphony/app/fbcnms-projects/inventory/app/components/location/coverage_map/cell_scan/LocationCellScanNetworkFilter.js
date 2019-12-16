/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

'use strict';

import type {FilterProps} from '../../../comparison_view/ComparisonViewTypes';

import PowerSearchFilter from '../../../comparison_view/PowerSearchFilter';
import React, {useState} from 'react';
import Tokenizer from '@fbcnms/ui/components/Tokenizer';

const COMMON_NETWORK_TYPES = ['LTE', 'WCDMA', 'GSM'];

const LocationCellScanNetworkFilter = (props: FilterProps) => {
  const {
    config,
    value,
    onInputBlurred,
    onValueChanged,
    onRemoveFilter,
    editMode,
  } = props;

  const networkTypes = config.extraData?.networkTypes ?? COMMON_NETWORK_TYPES;
  const [searchEntries, setSearchEntries] = useState([]);
  const [tokens, setTokens] = useState([]);

  return (
    <PowerSearchFilter
      name="Network Type"
      operator={value.operator}
      editMode={editMode}
      value={(value.idSet ?? [])
        .map(id => networkTypes.find(type => type === id))
        .join(', ')}
      onRemoveFilter={onRemoveFilter}
      input={
        <Tokenizer
          searchSource="Options"
          tokens={tokens}
          onEntriesRequested={searchTerm =>
            setSearchEntries(
              networkTypes
                .filter(type =>
                  type.toLowerCase().includes(searchTerm.toLowerCase()),
                )
                .map(type => ({id: type, label: type})),
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

export default LocationCellScanNetworkFilter;
