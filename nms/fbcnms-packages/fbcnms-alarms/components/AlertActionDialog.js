/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {AlertConfig} from './AlarmAPIType';

import Button from '@material-ui/core/Button';
import ClipboardLink from '@fbcnms/ui/components/ClipboardLink';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import React from 'react';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  paper: {
    minWidth: 360,
  },
  pre: {
    whiteSpace: 'pre-wrap',
    wordBreak: 'break-all',
  },
});

type Props = {
  open: boolean,
  onClose: () => void,
  title: string,
  additionalContent?: any,
  alertConfig: AlertConfig | Object,
  showCopyButton?: boolean,
  showDeleteButton?: boolean,
  onDelete?: () => Promise<void>,
};

export default function AlertActionDialog(props: Props) {
  const {
    open,
    onClose,
    title,
    additionalContent,
    alertConfig,
    showCopyButton,
    showDeleteButton,
    onDelete,
  } = props;
  const classes = useStyles();

  return (
    <Dialog
      PaperProps={{classes: {root: classes.paper}}}
      open={open}
      onClose={onClose}>
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        <pre className={classes.pre}>
          {JSON.stringify(alertConfig, null, 2)}
        </pre>
        {additionalContent}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} color="primary">
          {showDeleteButton ? 'Cancel' : 'Close'}
        </Button>
        {showCopyButton && (
          <ClipboardLink>
            {({copyString}) => (
              <Button
                onClick={() => copyString(JSON.stringify(alertConfig))}
                color="primary"
                variant="contained">
                Copy
              </Button>
            )}
          </ClipboardLink>
        )}
        {showDeleteButton && (
          <Button onClick={onDelete} color="primary" variant="contained">
            Delete
          </Button>
        )}
      </DialogActions>
    </Dialog>
  );
}
