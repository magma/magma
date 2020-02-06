/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import InfoTooltip from '../../components/design-system/Tooltip/InfoTooltip';
import React from 'react';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(_theme => ({
  root: {
    paddingTop: '100px',
    width: '100%',
  },
}));

const InfoTooltipRoot = () => {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <InfoTooltip description="This is a tooltip with extra information about the section" />
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('InfoTooltip', () => (
  <InfoTooltipRoot />
));
