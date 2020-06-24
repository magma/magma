/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import BaseDialog from '../../components/design-system/Dialog/BaseDialog';
import Button from '../../components/design-system/Button';
import MessageDialog from '../../components/design-system/Dialog/MessageDialog';
import React, {useState} from 'react';
import Text from '../../components/design-system/Text';
import {POSITION} from '../../components/design-system/Dialog/DialogFrame';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(() => ({
  root: {
    width: '100%',
  },
  menuContainer: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'flex-start',
    '& > *': {
      marginBottom: '16px',
    },
  },
  content: {
    display: 'flex',
    flexDirection: 'column',
  },
  pageContent: {
    height: 'calc(100vh - 400px)',
    background: 'white',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
  },
}));

const DialogsRoot = () => {
  const classes = useStyles();
  const [isMessageDialogShown, setIsMessageDialogShown] = useState(false);
  const [isDialogShown, setIsDialogShown] = useState(false);
  const [isRightDialogShown, setIsRightDialogShown] = useState(false);
  const [isNoMaskDialogShown, setIsNoMaskDialogShown] = useState(false);
  const closeDialog = () => {
    setIsMessageDialogShown(false);
    setIsDialogShown(false);
    setIsRightDialogShown(false);
    setIsNoMaskDialogShown(false);
  };

  return (
    <div className={classes.root}>
      <div className={classes.menuContainer}>
        <Button onClick={() => setIsDialogShown(true)}>Open Dialog</Button>
        <Button onClick={() => setIsRightDialogShown(true)}>
          Open Right Dialog
        </Button>
        <Button onClick={() => setIsNoMaskDialogShown(true)}>
          Open Dialog with No Mask
        </Button>
        <Button onClick={() => setIsMessageDialogShown(true)}>
          Open Message Dialog
        </Button>
      </div>
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
        hidden={!isMessageDialogShown}
      />
      <BaseDialog
        title="Base Dialog"
        onClose={closeDialog}
        hidden={!isDialogShown}>
        <Text>This is the dialog content.</Text>
      </BaseDialog>
      <BaseDialog
        title="Base Dialog - on side!"
        position={POSITION.right}
        onClose={closeDialog}
        hidden={!isRightDialogShown}>
        <Text>This is the dialog content.</Text>
      </BaseDialog>
      <BaseDialog
        title="Base Dialog - no masked background"
        position={POSITION.right}
        isModal={false}
        onClose={closeDialog}
        hidden={!isNoMaskDialogShown}>
        <Text>This is the dialog content.</Text>
        <Text>Clicking out side of panel will not close it.</Text>
      </BaseDialog>
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('Dialog', () => (
  <DialogsRoot />
));
