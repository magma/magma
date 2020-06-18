/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Button from '@material-ui/core/Button';
import React, {useState} from 'react';
import SideBar from '../../components/layout/SideBar';
import Text from '../../components/design-system/Text';
import TopPageBar from '../../components/layout/TopPageBar';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(_theme => ({
  root: {
    margin: '-8px',
  },
}));

const Container = () => {
  const classes = useStyles();
  const [isShown, setIsShown] = useState(false);
  return (
    <div className={classes.root}>
      <TopPageBar>
        <Text variant="body2">I'm a Header</Text>
        <Button onClick={() => setIsShown(true)}>Open Drawer</Button>
      </TopPageBar>
      <SideBar isShown={isShown} top={60} onClose={() => setIsShown(false)}>
        Content
      </SideBar>
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/SideBar`, module).add(
  'default',
  () => <Container />,
);
