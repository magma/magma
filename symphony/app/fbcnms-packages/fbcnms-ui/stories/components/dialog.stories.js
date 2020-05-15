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
import MessageDialog from '../../components/design-system/Dialog/MessageDialog';
import React, {useState} from 'react';
import Text from '../../components/design-system/Text';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(() => ({
  root: {
    width: '100%',
  },
  content: {
    display: 'flex',
    flexDirection: 'column',
  },
  pageContent: {
    height: 'calc(100vh + 50px)',
    background: 'white',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
  },
}));

const DialogsRoot = () => {
  const classes = useStyles();
  const [isDialogShown, setIsDialogShown] = useState(false);
  const closeDialog = () => {
    setIsDialogShown(false);
  };

  return (
    <div className={classes.root}>
      <Button onClick={() => setIsDialogShown(true)}>Open Dialog</Button>
      <div className={classes.pageContent}>
        <Text>Page Content</Text>
      </div>
      <MessageDialog
        title="Message Dialog"
        message={
          <div className={classes.content}>
            <Text>This is the message of the popup.</Text>
            <Text>Click Save to approve it or cancel to cancel.</Text>
          </div>
        }
        onClose={closeDialog}
        verificationCheckbox={{label: 'I understand', isMandatory: true}}
        cancelLabel="Cancel"
        confirmLabel="Save"
        onCancel={closeDialog}
        onConfirm={closeDialog}
        hidden={!isDialogShown}
      />
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('MessageDialog', () => (
  <DialogsRoot />
));
