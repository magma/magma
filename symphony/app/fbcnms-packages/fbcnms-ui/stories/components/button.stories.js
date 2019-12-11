/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AddIcon from '@material-ui/icons/Add';
import Button from '../../components/design-system/Button';
import React from 'react';
import classNames from 'classnames';
import {STORY_CATEGORIES} from '../storybookUtils';
import {fbt} from 'fbt';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '100%',
  },
  container: {
    margin: '16px 0px',
    display: 'flex',
    '& > button': {
      marginRight: '8px',
    },
  },
  grayButtonContainer: {
    backgroundColor: 'white',
    padding: '10px',
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
        <Button onClick={onButtonClicked} rightIcon={AddIcon}>
          Default
        </Button>
        <Button onClick={onButtonClicked} leftIcon={AddIcon} disabled>
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
        <Button onClick={onButtonClicked} skin="regular" rightIcon={AddIcon}>
          Default
        </Button>
        <Button
          onClick={onButtonClicked}
          skin="regular"
          leftIcon={AddIcon}
          disabled>
          Default
        </Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked}>OK</Button>
        <Button onClick={onButtonClicked} rightIcon={AddIcon}>
          OK
        </Button>
        <Button onClick={onButtonClicked} leftIcon={AddIcon}>
          OK
        </Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked} skin="red">
          Default
        </Button>
        <Button onClick={onButtonClicked} skin="red" disabled>
          Default
        </Button>
        <Button onClick={onButtonClicked} skin="red" rightIcon={AddIcon}>
          Default
        </Button>
        <Button
          onClick={onButtonClicked}
          skin="red"
          leftIcon={AddIcon}
          disabled>
          Default
        </Button>
      </div>
      <div
        className={classNames(classes.container, classes.grayButtonContainer)}>
        <Button onClick={onButtonClicked} skin="gray">
          Default
        </Button>
        <Button onClick={onButtonClicked} skin="gray" disabled>
          Default
        </Button>
        <Button onClick={onButtonClicked} skin="gray" rightIcon={AddIcon}>
          Default
        </Button>
        <Button
          onClick={onButtonClicked}
          skin="gray"
          leftIcon={AddIcon}
          disabled>
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
        <Button onClick={onButtonClicked} variant="text" rightIcon={AddIcon}>
          Default
        </Button>
        <Button
          onClick={onButtonClicked}
          variant="text"
          disabled
          leftIcon={AddIcon}>
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
        <Button
          onClick={onButtonClicked}
          variant="text"
          skin="regular"
          rightIcon={AddIcon}>
          Default
        </Button>
        <Button
          onClick={onButtonClicked}
          variant="text"
          skin="regular"
          disabled
          leftIcon={AddIcon}>
          Default
        </Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked} variant="text">
          OK
        </Button>
        <Button onClick={onButtonClicked} variant="text" rightIcon={AddIcon}>
          OK
        </Button>
        <Button onClick={onButtonClicked} variant="text" rightIcon={AddIcon}>
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
        <Button
          onClick={onButtonClicked}
          variant="text"
          skin="red"
          rightIcon={AddIcon}>
          Default
        </Button>
        <Button
          onClick={onButtonClicked}
          variant="text"
          skin="red"
          disabled
          leftIcon={AddIcon}>
          Default
        </Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked} variant="text" skin="gray">
          Default
        </Button>
        <Button onClick={onButtonClicked} variant="text" skin="gray" disabled>
          Default
        </Button>
        <Button
          onClick={onButtonClicked}
          variant="text"
          skin="gray"
          rightIcon={AddIcon}>
          Default
        </Button>
        <Button
          onClick={onButtonClicked}
          variant="text"
          skin="gray"
          disabled
          leftIcon={AddIcon}>
          Default
        </Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked} rightIcon={AddIcon}>
          Button with a long label
        </Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked} rightIcon={AddIcon}>
          <fbt desc="wow, much desc">Translated</fbt>{' '}
          {fbt('with a function', 'wow, much desc')}
        </Button>
      </div>
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('Button', () => (
  <ButtonsRoot />
));
