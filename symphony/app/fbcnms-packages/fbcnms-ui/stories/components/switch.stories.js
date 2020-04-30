/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Checkbox from '../../components/design-system/Checkbox/Checkbox';
import React, {useState} from 'react';
import Switch from '../../components/design-system/switch/Switch';
import Text from '../../components/design-system/Text';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(_theme => ({
  root: {},
  sample: {
    display: 'flex',
    margin: '8px',
    '& > *': {
      marginRight: '4px',
    },
  },
  optionsContainer: {
    display: 'flex',
    flexDirection: 'column',
    marginTop: '32px',
  },
  displayOption: {
    marginTop: '4px',
    display: 'flex',
    alignItems: 'center',
  },
  optionCheckbox: {
    marginRight: '8px',
  },
}));

const SwitchRoot = () => {
  const classes = useStyles();
  const [checked, setChecked] = useState(true);
  const [isCritical, setIsCritical] = useState(false);
  const [isDisabled, setIsDisabled] = useState(false);
  return (
    <div className={classes.root}>
      <div className={classes.sample}>
        <Switch
          checked={checked}
          onChange={setChecked}
          disabled={isDisabled}
          skin={isCritical ? 'critical' : undefined}
        />
        <Text variant="body1">{checked ? 'On' : 'Off'}</Text>
      </div>
      <div className={classes.optionsContainer}>
        <Text variant="h6">Variants:</Text>
        <div className={classes.displayOption}>
          <Checkbox
            className={classes.optionCheckbox}
            checked={isDisabled}
            onChange={selection =>
              setIsDisabled(selection === 'checked' ? true : false)
            }
          />
          <Text>Show disabled</Text>
        </div>
        <div className={classes.displayOption}>
          <Checkbox
            className={classes.optionCheckbox}
            checked={isCritical}
            onChange={selection =>
              setIsCritical(selection === 'checked' ? true : false)
            }
          />
          <Text>Show critical</Text>
        </div>
      </div>
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('Switch', () => (
  <SwitchRoot />
));
