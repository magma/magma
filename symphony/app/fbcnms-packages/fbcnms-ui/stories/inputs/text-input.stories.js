/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import CancelIcon from '@material-ui/icons/Cancel';
import InputAffix from '../../components/design-system/Input/InputAffix';
import React, {useEffect, useRef, useState} from 'react';
import TextInput from '../../components/design-system/Input/TextInput';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  root: {
    backgroundColor: 'white',
    height: '100vh',
    padding: '16px',
  },
  input: {
    marginBottom: '20px',
  },
  suffix: {
    cursor: 'pointer',
  },
}));

const InputsRoot = () => {
  const classes = useStyles();
  const [value, setValue] = useState('');
  const [numberValue, setNumberValue] = useState(123.4);
  const [affixValue, setAffixValue] = useState('Default');
  const inputRef = useRef<HTMLInputElement | null>(null);
  useEffect(() => {
    inputRef.current && inputRef.current.focus();
  }, []);

  return (
    <div className={classes.root}>
      <TextInput
        className={classes.input}
        type="string"
        placeholder="Placeholder"
        onChange={({target}) => setValue(target.value)}
        value={value}
        ref={inputRef}
      />
      <TextInput
        className={classes.input}
        type="string"
        placeholder="Placeholder"
        value=""
        disabled
      />
      <TextInput
        className={classes.input}
        type="string"
        placeholder="Placeholder"
        disabled={true}
        value="Value"
      />
      <TextInput
        className={classes.input}
        type="string"
        placeholder="Placeholder"
        hasError={true}
        value="Bad Value"
      />
      <TextInput
        className={classes.input}
        type="string"
        value="Default value"
      />
      <TextInput
        className={classes.input}
        type="number"
        value={numberValue}
        onChange={({target}) => setNumberValue(target.value)}
        prefix={<InputAffix>$</InputAffix>}
      />
      <TextInput
        className={classes.input}
        type="string"
        placeholder="Search..."
        value={affixValue}
        onChange={({target}) => setAffixValue(target.value)}
        suffix={
          affixValue ? (
            <InputAffix
              className={classes.suffix}
              onClick={() => setAffixValue('')}>
              <CancelIcon />
            </InputAffix>
          ) : null
        }
      />
      <TextInput className={classes.input} type="string" hint="Hinting" />
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.INPUTS}`, module).add('2.1 TextInput', () => (
  <InputsRoot />
));
