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
import React from 'react';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '100%',
  },
  container: {
    margin: '16px 0px',
    '& > :first-child': {
      marginRight: '8px',
    },
  },
}));

const ButtonsRoot = () => {
  const classes = useStyles();
  const onButtonClicked = () => window.alert('clicked!');
  return (
    <div className={classes.root}>
      <div className={classes.container}>
        <Button onClick={onButtonClicked}>Default</Button>
        <Button onClick={onButtonClicked} disabled>
          Default
        </Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked} skin="regular">
          Default
        </Button>
        <Button onClick={onButtonClicked} skin="regular" disabled>
          Default
        </Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked}>OK</Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked} skin="red">
          Default
        </Button>
        <Button onClick={onButtonClicked} skin="red" disabled>
          Default
        </Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked} variant="text">
          Default
        </Button>
        <Button onClick={onButtonClicked} variant="text" disabled>
          Default
        </Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked} variant="text" skin="regular">
          Default
        </Button>
        <Button
          onClick={onButtonClicked}
          variant="text"
          skin="regular"
          disabled>
          Default
        </Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked} variant="text">
          OK
        </Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked} variant="text" skin="red">
          Default
        </Button>
        <Button onClick={onButtonClicked} variant="text" skin="red" disabled>
          Default
        </Button>
      </div>
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('Button', () => (
  <ButtonsRoot />
));
