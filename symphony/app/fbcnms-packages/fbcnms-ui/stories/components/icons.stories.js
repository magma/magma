/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

/*eslint import/namespace: ['error', { allowComputed: true }]*/
import * as Icons from '../../components/design-system/Icons';
import * as React from 'react';
import AddIcon from '../../components/design-system/Icons/Actions/AddIcon';
import Text from '../../components/design-system/Text';

import Grid from '@material-ui/core/Grid';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '100%',
  },
  iconColors: {
    marginBottom: '16px',
  },
  iconRoot: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: '16px',
  },
  iconName: {
    marginLeft: '8px',
    color: '#374050',
  },
  iconsContainer: {
    display: 'flex',
  },
}));

type IconProps = {
  name: string,
  icon: React.Node,
};

const Icon = ({icon, name}: IconProps) => {
  const classes = useStyles();
  return (
    <Grid item className={classes.iconRoot} xs={3}>
      {icon}
      <Text className={classes.iconName} variant="body1">
        {name}
      </Text>
    </Grid>
  );
};

export const IconsRoot = () => {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <div className={classes.iconColors}>
        <AddIcon />
        <AddIcon color="light" />
        <AddIcon color="primary" />
        <AddIcon color="error" />
        <AddIcon color="gray" />
      </div>
      <Grid className={classes.iconsContainer} container>
        {Object.keys(Icons)
          .sort()
          .map((name: $Keys<typeof Icons>) => {
            const SvgIcon = Icons[name];
            return <Icon key={name} icon={<SvgIcon />} name={name} />;
          })}
      </Grid>
    </div>
  );
};

IconsRoot.story = {
  name: 'Icons',
};

export default {
  title: `${STORY_CATEGORIES.FOUNDATION}`,
};
