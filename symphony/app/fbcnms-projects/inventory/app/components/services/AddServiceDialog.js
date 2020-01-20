/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AddServiceDetails from './AddServiceDetails';
import CircularProgress from '@material-ui/core/CircularProgress';
import Dialog from '@material-ui/core/Dialog';
import React, {Suspense, useState} from 'react';
import ServiceTypesList from './ServiceTypesList';
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
  serviceTypesDialogContent: {
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
    backgroundColor: symphony.palette.white,
  },
  loading: {
    display: 'flex',
    height: '100%',
    alignItems: 'center',
    justifyContent: 'center',
  },
}));

type Props = {
  open: boolean,
  onClose: () => void,
  onServiceCreated: (id: string) => void,
};

const AddServiceDialog = (props: Props) => {
  const {open, onClose, onServiceCreated} = props;
  const [selectedServiceTypeId, setSelectedServiceTypeId] = useState(null);
  const [activeStep, setActiveStep] = useState(0);
  const classes = useStyles();

  return (
    <Dialog
      maxWidth="sm"
      open={open}
      onClose={onClose}
      fullWidth={true}
      className={classes.root}>
      {activeStep == 0 ? (
        <ServiceTypesList
          onSelect={type => {
            setSelectedServiceTypeId(type);
            setActiveStep(1);
          }}
          onClose={onClose}
        />
      ) : (
        <Suspense
          fallback={
            <div className={classes.loading}>
              <CircularProgress />
            </div>
          }>
          <AddServiceDetails
            serviceTypeId={nullthrows(selectedServiceTypeId)}
            onServiceCreated={onServiceCreated}
            onBackClicked={() => {
              setSelectedServiceTypeId(null);
              setActiveStep(0);
            }}
          />
        </Suspense>
      )}
    </Dialog>
  );
};

export default AddServiceDialog;
