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
import ProjectTypesList from './ProjectTypesList';
import React, {useState} from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';

import nullthrows from '@fbcnms/util/nullthrows';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  root: {
    position: 'relative',
  },
  avatar: {
    backgroundColor: '#e4f2ff',
  },
  dialogTitle: {
    padding: '24px',
    paddingBottom: '16px',
  },
  dialogTitleText: {
    fontSize: '20px',
    lineHeight: '24px',
    color: theme.palette.blueGrayDark,
    fontWeight: 500,
  },
  dialogContent: {
    padding: 0,
    height: '400px',
    overflowY: 'scroll',
  },
  dialogActions: {
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
  onProjectTypeSelected: (id: string) => void,
};

const AddProjectDialog = (props: Props) => {
  const [selectedProjectTypeId, setSelectedProjectTypeId] = useState(null);
  const classes = useStyles();
  return (
    <Dialog
      maxWidth="sm"
      open={props.open}
      onClose={props.onClose}
      fullWidth={true}
      className={classes.root}>
      <DialogTitle className={classes.dialogTitle}>
        <Text className={classes.dialogTitleText}>
          Select a template for this project
        </Text>
      </DialogTitle>
      <DialogContent className={classes.dialogContent}>
        <ProjectTypesList onSelect={type => setSelectedProjectTypeId(type)} />
      </DialogContent>
      <DialogActions className={classes.dialogActions}>
        <Button onClick={props.onClose} skin="regular">
          Cancel
        </Button>
        <Button
          disabled={selectedProjectTypeId === null}
          onClick={() => {
            props.onProjectTypeSelected(nullthrows(selectedProjectTypeId));
          }}>
          OK
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default AddProjectDialog;
