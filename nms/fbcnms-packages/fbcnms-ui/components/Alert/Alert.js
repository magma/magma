/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Node} from 'react';
import type {WithStyles} from '@material-ui/core';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import React from 'react';
import {withStyles} from '@material-ui/core/styles';

type Props = WithStyles<typeof styles> & {|
  cancelLabel?: Node,
  confirmLabel?: Node,
  message: Node,
  onCancel?: (element: HTMLElement) => void,
  onClose?: (element: HTMLElement) => void,
  onConfirm?: (element: HTMLElement) => void,
  title?: ?Node,
  open?: boolean,
|};

const styles = theme => ({
  paper: {
    minWidth: `${theme.breakpoints.values.sm / 2}px`,
  },
});

class Alert extends React.Component<Props> {
  render() {
    const {
      classes,
      cancelLabel,
      confirmLabel,
      message,
      onCancel,
      onClose,
      onConfirm,
      title,
      open,
    } = this.props;
    const hasActions = cancelLabel != null || confirmLabel != null;

    return (
      <Dialog
        classes={classes}
        open={open}
        onClose={onCancel}
        onExited={onClose}
        maxWidth="sm">
        {title && <DialogTitle>{title}</DialogTitle>}
        <DialogContent>
          <DialogContentText>{message}</DialogContentText>
        </DialogContent>
        {hasActions && (
          <DialogActions>
            {cancelLabel && (
              <Button onClick={onCancel} color="primary">
                {cancelLabel}
              </Button>
            )}
            {confirmLabel && (
              <Button
                onClick={onConfirm}
                color="primary"
                variant="contained"
                autoFocus>
                {confirmLabel}
              </Button>
            )}
          </DialogActions>
        )}
      </Dialog>
    );
  }
}

export default withStyles(styles)(Alert);
