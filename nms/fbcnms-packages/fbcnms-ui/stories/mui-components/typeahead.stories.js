/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Suggestion} from '../../components/Typeahead';

import React, {useState} from 'react';
import Typeahead from '../../components/Typeahead';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

const MOCK_ENTRIES: Array<Suggestion> = [
  {name: 'Tel Aviv', entityId: '0', entityType: 'location', type: 'City'},
  {name: 'Haifa', entityId: '1', entityType: 'location', type: 'City'},
  {
    name: 'Sarona Building',
    entityId: '2',
    entityType: 'location',
    type: 'Building',
  },
];

const AddLocation = () => {
  const [entries, setEntries] = useState([]);
  const [selectedLocationName, setSelectLocationName] = useState('');
  return (
    <>
      <div style={{marginBottom: '20px'}}>{selectedLocationName}</div>
      <Typeahead
        required={true}
        suggestions={entries}
        onEntitySelected={suggestion => setSelectLocationName(suggestion.name)}
        onSuggestionsFetchRequested={value => {
          setEntries(
            MOCK_ENTRIES.filter(e =>
              e.name.toLowerCase().includes(value.toLowerCase()),
            ),
          );
        }}
        onEntriesRequested={() => {}}
        searchEntries={entries}
        onSuggestionsClearRequested={() => setSelectLocationName('')}
        placeholder="Location"
      />
    </>
  );
};

const SearchBar = () => {
  const [entries, setEntries] = useState([]);
  const [selectedLocationName, setSelectLocationName] = useState('');
  return (
    <>
      <div style={{marginBottom: '20px'}}>{selectedLocationName}</div>
      <Typeahead
        required={false}
        suggestions={entries}
        onEntitySelected={suggestion => {
          setSelectLocationName(suggestion.name);
        }}
        onSuggestionsFetchRequested={value => {
          setEntries(
            MOCK_ENTRIES.filter(e =>
              e.name.toLowerCase().includes(value.toLowerCase()),
            ),
          );
        }}
        onEntriesRequested={() => {}}
        searchEntries={entries}
      />
    </>
  );
};

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/Typehead`, module)
  .add('addLocation', () => {
    return (
      <div style={{width: '300px'}}>
        <AddLocation />
      </div>
    );
  })
  .add('search', () => {
    return (
      <div style={{width: '300px'}}>
        <SearchBar />
      </div>
    );
  });
