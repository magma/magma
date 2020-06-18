/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import React from 'react';
import Text from '../../components/design-system/Text';
import classNames from 'classnames';
import symphony from '@fbcnms/ui/theme/symphony';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(() => ({
  root: {
    padding: '52px',
    display: 'flex',
    flexDirection: 'row',
    flexWrap: 'wrap',
  },
  blockRoot: {
    marginBottom: '40px',
    marginRight: '40px',
  },
  shadowBlock: {
    width: '200px',
    height: '200px',
    backgroundColor: symphony.palette.white,
    marginBottom: '12px',
  },
  dp1: {boxShadow: symphony.shadows.DP1},
  dp2: {boxShadow: symphony.shadows.DP2},
  dp3: {boxShadow: symphony.shadows.DP3},
  dp4: {boxShadow: symphony.shadows.DP4},
}));

const ShadowBlock = (props: {shadow: string}) => {
  const {shadow} = props;
  const classes = useStyles();

  return (
    <div className={classes.blockRoot}>
      <div
        className={classNames(
          classes.shadowBlock,
          classes[shadow.toLowerCase()],
        )}
      />
      <Text weight="medium">{shadow.toUpperCase()}</Text>
    </div>
  );
};

const ShadowBlocksRoot = () => {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      {Object.keys(symphony.shadows).map(shadow => (
        <ShadowBlock key={shadow} shadow={shadow} />
      ))}
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.FOUNDATION}`, module).add('1.2 Elevation', () => (
  <ShadowBlocksRoot />
));
