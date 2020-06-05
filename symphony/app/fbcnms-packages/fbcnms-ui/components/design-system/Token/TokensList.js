/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Entries} from './Tokenizer';

import React from 'react';
import Token from './Token';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  token: {
    margin: '2px',
  },
}));

type Props<TEntry> = $ReadOnly<{|
  tokens: Entries<TEntry>,
  onTokensChange: (Entries<TEntry>) => void,
  disabled?: boolean,
|}>;

const TokensList = <TEntry>(props: Props<TEntry>) => {
  const {tokens, onTokensChange, disabled = false} = props;
  const classes = useStyles();
  return (
    <>
      {tokens.map(token => (
        <Token
          key={token.key}
          className={classes.token}
          disabled={disabled}
          label={token.label}
          onRemove={() =>
            onTokensChange(tokens.slice().filter(t => t.key !== token.key))
          }
        />
      ))}
    </>
  );
};

export default TokensList;
