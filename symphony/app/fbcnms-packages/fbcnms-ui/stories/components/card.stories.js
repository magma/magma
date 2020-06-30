/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Button from '../../components/design-system/Button';
import Card, {
  CARD_MARGINS,
  CARD_VARIANTS,
} from '../../components/design-system/Card/Card';
import CardHeader from '../../components/design-system/Card/CardHeader';
import React from 'react';
import Text from '../../components/design-system/Text';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '100%',
  },
  mapContent: {
    padding: '120px',
    background: 'gray',
    backgroundClip: 'padding-box',
  },
}));

export const CardsRoot = () => {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <Card>
        <Text>Content</Text>
      </Card>
      <Card margins={CARD_MARGINS.none}>
        <div className={classes.mapContent}>
          <Text variant="h5" color="light">
            my content here needs no margins
          </Text>
        </div>
      </Card>
      <Card>
        <CardHeader>Title</CardHeader>
        <Text>Content</Text>
      </Card>
      <Card variant={CARD_VARIANTS.message}>
        <CardHeader>This is a system message</CardHeader>
        <Text>System message content</Text>
      </Card>
      <Card>
        <CardHeader rightContent={<Button>Action</Button>}>Title</CardHeader>
        <Text>Content</Text>
      </Card>
    </div>
  );
};

CardsRoot.story = {
  name: 'Card',
};

export default {
  title: `${STORY_CATEGORIES.COMPONENTS}`,
};
