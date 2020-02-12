/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import React, {useState} from 'react';
import Tokenizer from '../../components/Tokenizer';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';
import type {Entry} from '../../components/Tokenizer';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '300px',
  },
}));

const entries = [
  {label: 'Chassis', id: '0'},
  {label: 'Rack', id: '1'},
  {label: 'Card', id: '2'},
  {label: 'AP', id: '3'},
];

function TestTokenizer(props: {
  searchSource: 'Options' | 'UserInput',
  searchEntries?: Array<Entry>,
  defaultTokens?: Array<Entry>,
}) {
  const {searchSource, searchEntries, defaultTokens = []} = props;
  const classes = useStyles();
  const [tokens, setTokens] = useState(defaultTokens);
  return (
    <div className={classes.root}>
      <Tokenizer
        searchSource={searchSource}
        tokens={tokens}
        searchEntries={searchEntries}
        onEntriesRequested={() => {}}
        onChange={entires => setTokens(entires)}
      />
    </div>
  );
}

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/Tokenizer`, module)
  .add('options', () => (
    <TestTokenizer searchSource="Options" searchEntries={entries} />
  ))
  .add('userInput', () => <TestTokenizer searchSource="UserInput" />)
  .add('userInputDefaultToken', () => (
    <TestTokenizer
      searchSource="UserInput"
      defaultTokens={[{id: 'Default', label: 'Default'}]}
    />
  ));
