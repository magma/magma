/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {RadioOption} from '../../components/design-system/RadioGroup/RadioGroup';

import Checkbox from '../../components/design-system/Checkbox/Checkbox';
import RadioGroup from '../../components/design-system/RadioGroup/RadioGroup';
import React, {useState} from 'react';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(() => ({
  root: {
    width: '100%',
  },
  optionsContainer: {
    display: 'flex',
    flexDirection: 'column',
    marginTop: '32px',
  },
}));

const options: Array<RadioOption> = [
  {
    value: 'option_1',
    label: 'Option 1',
    details: 'wow1',
  },
  {
    value: 'option_2',
    label: 'Option 2',
    details: 'wow2',
  },
  {
    value: 'option_3',
    label: 'Option 3',
    details: 'wow3',
  },
  {
    value: 'option_4',
    label: 'Option 4',
    details: 'wow4',
  },
];

const RadioGroupRoot = () => {
  const classes = useStyles();
  const [selectedValue, setSelectedValue] = useState(options[0].value);
  const [isDisabled, setIsDisabled] = useState(false);
  const [isOptionDisabled, setIsOptionDisabled] = useState(false);

  const changeIsOptionDisabled = newValue => {
    options[0].disabled = newValue;
    setIsOptionDisabled(newValue);
  };

  return (
    <div className={classes.root}>
      <RadioGroup
        disabled={isDisabled}
        options={options}
        value={selectedValue}
        onChange={value => setSelectedValue(value)}
      />
      <div className={classes.optionsContainer}>
        <Checkbox
          checked={isDisabled}
          title="All Disabled"
          onChange={selection =>
            setIsDisabled(selection === 'checked' ? true : false)
          }
        />
        <Checkbox
          checked={isOptionDisabled}
          title="First Option Disabled"
          onChange={selection =>
            changeIsOptionDisabled(selection === 'checked' ? true : false)
          }
        />
      </div>
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('RadioGroup', () => (
  <RadioGroupRoot />
));
