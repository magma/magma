/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import React, {useState} from 'react';
import ServiceTypesList from './ServiceTypesList';
import Text from '@fbcnms/ui/components/design-system/Text';

import nullthrows from '@fbcnms/util/nullthrows';
import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_ => ({
  root: {
    position: 'relative',
  },
  avatar: {
    backgroundColor: symphony.palette.B50,
  },
  dialogTitle: {
    padding: '24px',
    paddingBottom: '16px',
  },
  dialogTitleText: {
    fontSize: '20px',
    lineHeight: '24px',
    color: symphony.palette.D900,
    fontWeight: 500,
  },
  dialogContent: {
    padding: 0,
    height: '400px',
    overflowY: 'scroll',
  },
  dialogActions: {
    position: 'absolute',
    padding: '24px',
    bottom: 0,
    display: 'flex',
    justifyContent: 'flex-end',
    width: '100%',
    backgroundColor: 'rgba(255, 255, 255, 0.9)',
    zIndex: 2,
  },
}));

type Props = {
  open: boolean,
  onClose: () => void,
  onServiceTypeSelected: (id: string) => void,
};

const AddServiceDialog = (props: Props) => {
  const [selectedServiceTypeId, setSelectedServiceTypeId] = useState(null);
  const classes = useStyles();
  return (
    <Dialog
      maxWidth="sm"
      open={props.open}
      onClose={props.onClose}
      fullWidth={true}
      className={classes.root}>
      <DialogTitle className={classes.dialogTitle}>
        <Text variant="h6">Select a Service type</Text>
      </DialogTitle>
      <DialogContent className={classes.dialogContent}>
        <ServiceTypesList onSelect={type => setSelectedServiceTypeId(type)} />
      </DialogContent>
      <DialogActions className={classes.dialogActions}>
        <Button onClick={props.onClose} skin="regular">
          Cancel
        </Button>
        <Button
          disabled={selectedServiceTypeId === null}
          onClick={() => {
            props.onServiceTypeSelected(nullthrows(selectedServiceTypeId));
          }}>
          OK
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default AddServiceDialog;
