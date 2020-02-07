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
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(_theme => ({
  root: {
    padding: '16px',
    backgroundColor: 'white',
    boxShadow: '0px 1px 4px 0px rgba(0,0,0,0.17)',
    display: 'flex',
    alignItems: 'center',
  },
  input: {
    width: '200px',
  },
  name: {
    marginRight: '32px',
    fontWeight: 500,
    fontSize: '16px',
  },
}));

const ENTER_KEY_CODE = 13;

type Props = {};

const EntSearchBar = (_props: Props) => {
  const {history} = useRouter();
  const classes = useStyles();
  const [searchText, setSearchText] = useState('');

  return (
    <div className={classes.root}>
      <Text className={classes.name}>ID Tool</Text>
      <TextField
        className={classes.input}
        variant="outlined"
        margin="dense"
        value={searchText}
        onChange={e => setSearchText(e.target.value)}
        placeholder="Enter ID..."
        onKeyDown={e => {
          if (e.keyCode === ENTER_KEY_CODE) {
            if (e.target.value) {
              history.push(`/id/${e.target.value}`);
            }

            setSearchText('');
          }
        }}
      />
    </div>
  );
};

export default EntSearchBar;
