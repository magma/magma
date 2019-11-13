/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import FormField from '../../components/design-system/FormField/FormField';
import React, {useState} from 'react';
import TextInput from '../../components/design-system/Input/TextInput';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  root: {
    backgroundColor: 'white',
    height: '400px',
    padding: '16px',
  },
  field: {
    marginBottom: '20px',
  },
}));

const FormFieldsRoot = () => {
  const classes = useStyles();
  const [value, setValue] = useState('');
  return (
    <div className={classes.root}>
      <FormField label="Label" helpText="Help Text" className={classes.field}>
        <TextInput
          type="string"
          placeholder="Placeholder"
          onChange={({target}) => setValue(target.value)}
          value={value}
        />
      </FormField>
      <FormField
        label="Required"
        helpText="Help Text"
        className={classes.field}
        required>
        <TextInput type="string" placeholder="Placeholder" value={value} />
      </FormField>
      <FormField
        label="Label"
        helpText="Help Text"
        className={classes.field}
        disabled>
        <TextInput type="string" placeholder="Placeholder" value="" />
      </FormField>
      <FormField
        label="Label"
        helpText="Help Text"
        className={classes.field}
        disabled>
        <TextInput
          type="string"
          placeholder="Placeholder"
          value="Default Value"
        />
      </FormField>
      <FormField
        label="Label"
        helpText="Help Text"
        className={classes.field}
        hasError
        errorText="Error Message">
        <TextInput
          className={classes.input}
          type="string"
          placeholder="Placeholder"
          value="Bad Value"
        />
      </FormField>
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.INPUTS}`, module)
  .add('2.2 FormField', () => <FormFieldsRoot />)
  .addParameters({order: 1});
