/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React, {useState} from 'react';
import TextInput from '../../components/design-system/TextInput.react';
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
    marginTop: '8px',
  },
}));

const InputsRoot = () => {
  const classes = useStyles();
  const [value, setValue] = useState('');
  const [numberValue, setNumberValue] = useState(123.4);
  return (
    <div className={classes.root}>
      <TextInput
        type="string"
        placeholder="Placeholder"
        onChange={({target}) => setValue(target.value)}
        value={value}
      />
      <TextInput
        className={classes.input}
        type="string"
        placeholder="Placeholder"
        disabled={true}
        value=""
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
        value=""
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
      />
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.INPUTS}`, module).add('2.1 TextInput', () => (
  <InputsRoot />
));
