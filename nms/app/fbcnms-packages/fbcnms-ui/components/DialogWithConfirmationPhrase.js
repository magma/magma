/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Button from './design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import React from 'react';
import TextField from '@material-ui/core/TextField';

import {useState} from 'react';

type Props = {
  title: string,
  message: string,
  confirmationPhrase: string,
  label: string,
  onClose: () => void,
  onConfirm: () => void | Promise<void>,
};

export default function DialogWithConfirmationPhrase(props: Props) {
  const [confirmationPhrase, setConfirmationPhrase] = useState('');
  const {title, message, label, onClose, onConfirm} = props;

  return (
    <Dialog open={true} onClose={onClose} onExited={onClose} maxWidth="sm">
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        <DialogContentText>
          {message}
          <TextField
            label={label}
            value={confirmationPhrase}
            onChange={({target}) => setConfirmationPhrase(target.value)}
          />
        </DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button skin="regular" onClick={onClose}>
          Cancel
        </Button>
        <Button
          skin="red"
          onClick={onConfirm}
          disabled={confirmationPhrase !== props.confirmationPhrase}>
          Confirm
        </Button>
      </DialogActions>
    </Dialog>
  );
}
