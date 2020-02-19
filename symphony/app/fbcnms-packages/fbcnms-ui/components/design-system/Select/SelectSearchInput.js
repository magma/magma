/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import ClearIcon from '@material-ui/icons/Clear';
import InputAffix from '../Input/InputAffix';
import React, {useEffect, useRef} from 'react';
import TextInput from '../Input/TextInput';
import fbt from 'fbt';
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';
import {useMenuContext} from './MenuContext';

const useStyles = makeStyles(() => ({
  input: {
    padding: '16px',
  },
  clearIconContainer: {
    backgroundColor: symphony.palette.background,
    padding: '6px',
    borderRadius: '100%',
    width: '20px',
    height: '20px',
    boxSizing: 'border-box',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
  },
  clearIcon: {
    color: symphony.palette.D800,
    fontSize: '13.66px',
  },
}));

type Props = {
  searchTerm: string,
  onChange: (searchTerm: string) => void,
};

const SelectSearchInput = ({searchTerm, onChange}: Props) => {
  const classes = useStyles();
  const inputRef = useRef(null);
  const {shown} = useMenuContext();
  const setFocus = () => inputRef.current?.focus();

  useEffect(() => {
    setFocus();
  }, [shown]);

  return (
    <TextInput
      className={classes.input}
      type="string"
      placeholder={fbt(
        'Type to filter...',
        'Input placeholder where user filters options',
      )}
      onChange={({target}) => onChange(target.value)}
      value={searchTerm}
      suffix={
        searchTerm ? (
          <InputAffix
            onClick={() => {
              onChange('');
              setFocus();
            }}
            className={classes.clearIconContainer}>
            <ClearIcon className={classes.clearIcon} />
          </InputAffix>
        ) : null
      }
      ref={inputRef}
    />
  );
};

export default SelectSearchInput;
