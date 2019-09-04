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
import PrettyJSON from '@fbcnms/ui/components/PrettyJSON';
import React from 'react';

type Props = {
  open: boolean,
  onClose: () => void,
  onDelete: () => void,
  alertConfig: AlertConfig,
  deletingAlert: boolean,
};

export default function ViewDeleteAlertDialog(props: Props) {
  const {open, onClose, onDelete, alertConfig, deletingAlert} = props;
  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>{deletingAlert ? 'Delete Alert' : 'View Alert'}</DialogTitle>
      <DialogContent>
        <PrettyJSON jsonObject={alertConfig} />
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} color="primary">
          {deletingAlert ? 'Cancel' : 'Close'}
        </Button>
        {deletingAlert ? (
          <Button onClick={onDelete} color="primary" variant="contained">
            Delete
          </Button>
        ) : (
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
      </DialogActions>
    </Dialog>
  );
}
