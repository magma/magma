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
import Token from '../../components/design-system/Token/Token';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  token: {
    margin: '16px',
  },
}));

export const TokensRoot = () => {
  const classes = useStyles();

  return (
    <div>
      <div className={classes.token}>
        <Token label="San Francisco" onRemove={() => alert('removed')} />
      </div>
      <div className={classes.token}>
        <Token
          label="Disabled"
          onRemove={() => alert('removed')}
          disabled={true}
        />
      </div>
    </div>
  );
};

TokensRoot.story = {
  name: 'Token',
};

export default {
  title: `${STORY_CATEGORIES.COMPONENTS}`,
};
