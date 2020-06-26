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
import symphony from '../../theme/symphony';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '100%',
  },
  lightText: {
    display: 'inline-block',
    backgroundColor: symphony.palette.D900,
  },
  text: {
    margin: '16px 0px',
  },
  mediumWeight: {
    fontWeight: 500,
  },
}));

export const TextRoot = () => {
  const classes = useStyles();
  const param = 'much wow';
  return (
    <div className={classes.root}>
      <div className={classes.text}>
        <Text variant="body2" color="regular">
          Regular Text
        </Text>
      </div>
      <div className={classes.text}>
        <div className={classes.lightText}>
          <Text variant="body2" color="light">
            Light Text
          </Text>
        </div>
      </div>
      <div className={classes.text}>
        <Text variant="body2" color="primary">
          Primary Text
        </Text>
      </div>
      <div className={classes.text}>
        <Text variant="body2" color="error">
          Error Text
        </Text>
      </div>
      <div className={classes.text}>
        <Text variant="body2" weight="light">
          Light Weight Text
        </Text>
      </div>
      <div className={classes.text}>
        <Text variant="body2" weight="regular">
          Regular Weight Text
        </Text>
      </div>
      <div className={classes.text}>
        <Text variant="body2" weight="medium">
          Medium Weight Text
        </Text>
      </div>
      <div className={classes.text}>
        <Text variant="body2" weight="bold">
          Bold Weight Text
        </Text>
      </div>
      <div className={classes.text}>
        <Text variant="body2" weight="bold">
          Text with {param} parameter
        </Text>
      </div>
      <div className={classes.text}>
        <Text variant="body2">
          <span className={classes.mediumWeight}>
            Text with multiple children:
          </span>{' '}
          yay
        </Text>
      </div>
    </div>
  );
};

TextRoot.story = {
  name: 'Text',
};

export default {
  title: `${STORY_CATEGORIES.COMPONENTS}`,
};
