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
import Tokenizer from '../../components/design-system/Token/Tokenizer';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(_theme => ({
  root: {
    height: '100vh',
    backgroundColor: 'white',
    padding: '16px',
  },
  tokenizer: {
    margin: '16px',
  },
}));

const ALL_TOKENS = [
  {key: '0', label: 'Tel Aviv'},
  {key: '1', label: 'San Francisco'},
  {key: '2', label: 'A long name of a famous city'},
  {key: '3', label: 'Budapest'},
  {key: '4', label: 'Los Angeles'},
  {key: '5', label: 'Some other city with a long name'},
  ,
];

const TokenizerRoot = () => {
  const classes = useStyles();
  const [queryString, setQueryString] = useState('');
  const [tokens, setTokens] = useState([
    {key: '0', label: 'Tel Aviv'},
    {key: '1', label: 'San Francisco'},
    {key: '2', label: 'A long name of a famous city'},
  ]);

  return (
    <div className={classes.root}>
      <div className={classes.tokenizer}>
        <Tokenizer
          tokens={tokens}
          onTokensChange={setTokens}
          queryString={queryString}
          onQueryStringChange={setQueryString}
          dataSource={{
            fetchNetwork: _searchTerm => Promise.resolve(ALL_TOKENS),
          }}
        />
      </div>
      <div className={classes.tokenizer}>
        <Tokenizer
          disabled={true}
          tokens={tokens}
          onTokensChange={setTokens}
          queryString={queryString}
          onQueryStringChange={setQueryString}
          dataSource={{
            fetchNetwork: _searchTerm => Promise.resolve(ALL_TOKENS),
          }}
        />
      </div>
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('Tokenizer', () => (
  <TokenizerRoot />
));
