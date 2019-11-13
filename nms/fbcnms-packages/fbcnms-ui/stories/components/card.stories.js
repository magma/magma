/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Button from '../../components/design-system/Button';
import Card from '../../components/design-system/Card/Card';
import CardHeader from '../../components/design-system/Card/CardHeader';
import React from 'react';
import Text from '../../components/design-system/Text';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '100%',
  },
  card: {
    marginBottom: '16px',
  },
}));

const CardsRoot = () => {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Card className={classes.card}>
        <Text>Content</Text>
      </Card>
      <Card className={classes.card}>
        <CardHeader>Title</CardHeader>
        <Text>Content</Text>
      </Card>
      <Card className={classes.card}>
        <CardHeader rightContent={<Button>Action</Button>}>Title</CardHeader>
        <Text>Content</Text>
      </Card>
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('Card', () => (
  <CardsRoot />
));
