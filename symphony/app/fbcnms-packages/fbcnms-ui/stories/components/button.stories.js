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
import PlusIcon from '../../components/design-system/Icons/Actions/PlusIcon';
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
        <Button onClick={onButtonClicked} rightIcon={PlusIcon}>
          Default
        </Button>
        <Button onClick={onButtonClicked} leftIcon={PlusIcon} disabled>
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
        <Button onClick={onButtonClicked} skin="regular" rightIcon={PlusIcon}>
          Default
        </Button>
        <Button
          onClick={onButtonClicked}
          skin="regular"
          leftIcon={PlusIcon}
          disabled>
          Default
        </Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked}>OK</Button>
        <Button onClick={onButtonClicked} rightIcon={PlusIcon}>
          OK
        </Button>
        <Button onClick={onButtonClicked} leftIcon={PlusIcon}>
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
        <Button onClick={onButtonClicked} skin="red" rightIcon={PlusIcon}>
          Default
        </Button>
        <Button
          onClick={onButtonClicked}
          skin="red"
          leftIcon={PlusIcon}
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
        <Button onClick={onButtonClicked} skin="gray" rightIcon={PlusIcon}>
          Default
        </Button>
        <Button
          onClick={onButtonClicked}
          skin="gray"
          leftIcon={PlusIcon}
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
        <Button onClick={onButtonClicked} variant="text" rightIcon={PlusIcon}>
          Default
        </Button>
        <Button
          onClick={onButtonClicked}
          variant="text"
          disabled
          leftIcon={PlusIcon}>
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
          rightIcon={PlusIcon}>
          Default
        </Button>
        <Button
          onClick={onButtonClicked}
          variant="text"
          skin="regular"
          disabled
          leftIcon={PlusIcon}>
          Default
        </Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked} variant="text">
          OK
        </Button>
        <Button onClick={onButtonClicked} variant="text" rightIcon={PlusIcon}>
          OK
        </Button>
        <Button onClick={onButtonClicked} variant="text" rightIcon={PlusIcon}>
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
          rightIcon={PlusIcon}>
          Default
        </Button>
        <Button
          onClick={onButtonClicked}
          variant="text"
          skin="red"
          disabled
          leftIcon={PlusIcon}>
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
          rightIcon={PlusIcon}>
          Default
        </Button>
        <Button
          onClick={onButtonClicked}
          variant="text"
          skin="gray"
          disabled
          leftIcon={PlusIcon}>
          Default
        </Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked}>Go!</Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked}>Button with a long label</Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked} rightIcon={PlusIcon}>
          Button with a long label and right icon
        </Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked} leftIcon={PlusIcon}>
          Button with a long label and left icon
        </Button>
      </div>
      <div className={classes.container}>
        <Button onClick={onButtonClicked} rightIcon={PlusIcon}>
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
