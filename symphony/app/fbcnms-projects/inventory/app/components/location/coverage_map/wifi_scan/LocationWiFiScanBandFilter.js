/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local strict-local
 * @format
 */

'use strict';

import type {FilterProps} from '../../../comparison_view/ComparisonViewTypes';

import PowerSearchFilter from '../../../comparison_view/PowerSearchFilter';
import React, {useState} from 'react';
import Tokenizer from '@fbcnms/ui/components/Tokenizer';

// https://en.wikipedia.org/wiki/List_of_WLAN_channels
const POSSIBLE_BANDS = ['2.4GHz', '3.65GHz', '4.9GHz', '5GHz'];

const LocationWiFiScanBandFilter = (props: FilterProps) => {
  const {
    config,
    value,
    onInputBlurred,
    onValueChanged,
    onRemoveFilter,
    editMode,
  } = props;

  const bands = config.extraData?.bands ?? POSSIBLE_BANDS;
  const [searchEntries, setSearchEntries] = useState([]);
  const [tokens, setTokens] = useState([]);

  return (
    <PowerSearchFilter
      name="Band"
      operator={value.operator}
      editMode={editMode}
      value={(value.idSet ?? [])
        .map(id => bands.find(type => type === id))
        .join(', ')}
      onRemoveFilter={onRemoveFilter}
      input={
        <Tokenizer
          searchSource="Options"
          tokens={tokens}
          onEntriesRequested={searchTerm =>
            setSearchEntries(
              bands
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

export default LocationWiFiScanBandFilter;
