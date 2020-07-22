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

import Button from '../design-system/Button';
import Checkbox from '../design-system/Checkbox/Checkbox';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import React, {useState} from 'react';
import Text from '../design-system/Text';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  paper: {
    minWidth: `${theme.breakpoints.values.sm / 2}px`,
  },
}));

export type AlertSkin = 'primary' | 'red';

type Props = {|
  cancelLabel?: Node,
  confirmLabel?: Node,
  message: Node,
  checkboxLabel?: Node,
  skin?: AlertSkin,
  onCancel?: () => void,
  onClose?: () => void,
  onConfirm?: () => void,
  title?: ?Node,
  open?: boolean,
|};

const Alert = ({
  cancelLabel,
  confirmLabel,
  message,
  checkboxLabel,
  onCancel,
  onClose,
  onConfirm,
  title,
  open,
  skin = 'primary',
}: Props) => {
  const classes = useStyles();
  const [checkboxChecked, setCheckboxChecked] = useState(false);
  const hasActions = cancelLabel != null || confirmLabel != null;

  return (
    <Dialog
      classes={{paper: classes.paper}}
      open={open}
      onClose={onCancel}
      onExited={onClose}
      maxWidth="sm">
      {title && <DialogTitle>{title}</DialogTitle>}
      <DialogContent>
        <Text>{message}</Text>
        {checkboxLabel && (
          <Checkbox
            checked={checkboxChecked}
            title={checkboxLabel}
            onChange={selection =>
              setCheckboxChecked(selection === 'checked' ? true : false)
            }
          />
        )}
      </DialogContent>
      {hasActions && (
        <DialogActions>
          {cancelLabel && (
            <Button skin="regular" onClick={onCancel}>
              {cancelLabel}
            </Button>
          )}
          {confirmLabel && (
            <Button
              onClick={onConfirm}
              skin={skin}
              disabled={checkboxLabel != null && !checkboxChecked}>
              {confirmLabel}
            </Button>
          )}
        </DialogActions>
      )}
    </Dialog>
  );
};

export default Alert;
