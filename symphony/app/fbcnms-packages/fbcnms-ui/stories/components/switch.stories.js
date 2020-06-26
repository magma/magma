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
  optionCheckbox: {
    marginTop: '4px',
  },
}));

export const SwitchRoot = () => {
  const classes = useStyles();
  const [checked, setChecked] = useState(true);
  const [isCritical, setIsCritical] = useState(false);
  const [isDisabled, setIsDisabled] = useState(false);
  const [isBold, setIsBold] = useState(false);
  return (
    <div className={classes.root}>
      <div className={classes.sample}>
        <Switch
          checked={checked}
          title={checked ? 'On' : 'Off'}
          variant={isBold ? 'subtitle2' : 'body2'}
          onChange={setChecked}
          disabled={isDisabled}
          skin={isCritical ? 'critical' : undefined}
        />
      </div>
      <div className={classes.optionsContainer}>
        <Text variant="h6">Variants:</Text>
        <Checkbox
          className={classes.optionCheckbox}
          checked={isDisabled}
          title="Show disabled"
          onChange={selection =>
            setIsDisabled(selection === 'checked' ? true : false)
          }
        />
        <Checkbox
          className={classes.optionCheckbox}
          checked={isCritical}
          title="Show critical"
          onChange={selection =>
            setIsCritical(selection === 'checked' ? true : false)
          }
        />
        <Checkbox
          className={classes.optionCheckbox}
          checked={isBold}
          title="Show bold"
          onChange={selection =>
            setIsBold(selection === 'checked' ? true : false)
          }
        />
      </div>
    </div>
  );
};

SwitchRoot.story = {
  name: 'Switch',
};

export default {
  title: `${STORY_CATEGORIES.COMPONENTS}`,
};
