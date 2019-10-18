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

import AddEditAlertConfigurationStep from './AddEditAlertConfigurationStep';
import AddEditAlertInfoStep from './AddEditAlertInfoStep';
import AddEditAlertNotificationStep from './AddEditAlertNotificationStep';
import Button from '@material-ui/core/Button';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import PrettyJSON from '@fbcnms/ui/components/PrettyJSON';
import React from 'react';
import Typography from '@material-ui/core/Typography';
import grey from '@material-ui/core/colors/grey';

import nullthrows from '@fbcnms/util/nullthrows';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

import {useState} from 'react';

type Props = {
  onExit: () => void,
};

const useStyles = makeStyles(theme => ({
  header: {
    padding: theme.spacing(3),
    display: 'flex',
    justifyContent: 'space-between',
    backgroundColor: 'white',
    borderBottom: `1px solid ${theme.palette.divider}`,
  },
  editingSpace: {
    height: '100%',
    float: 'left',
    width: '70%',
    padding: '70px',
  },
  previewSpace: {
    height: '100%',
    float: 'right',
    width: '30%',
    backgroundColor: grey[100],
    padding: '40px',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
  },
  alertPreview: {
    fontStyle: 'italic',
    fontSize: 15,
    fontWeight: 500,
    marginBottom: '20px',
  },
}));

export default function AddEditAlert(props: Props) {
  const [alertConfig, setAlertConfig] = useState<AlertConfig>({
    alert: '',
    expr: '',
    labels: {team: 'operations'},
  });

  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();

  const saveAlert = () => {
    MagmaV1API.postNetworksByNetworkIdPrometheusAlertConfig({
      networkId: nullthrows(match.params.networkId),
      alertConfig,
    })
      .then(() => props.onExit())
      .catch(error =>
        enqueueSnackbar(
          `Unable to create alert: ${
            error.response ? error.response.data.message : error.message
          }.`,
          {
            variant: 'error',
          },
        ),
      );
  };

  const classes = useStyles();

  return (
    <>
      <div className={classes.header}>
        <Typography variant="h5">New Alert</Typography>
        <Button
          variant="contained"
          color="secondary"
          onClick={() => props.onExit()}>
          Cancel
        </Button>
      </div>
      <div className={classes.editingSpace}>
        <ConfigurationStep
          alertConfig={alertConfig}
          setAlertConfig={setAlertConfig}
          saveAlert={saveAlert}
        />
      </div>
      <div className={classes.previewSpace}>
        <div className={classes.alertPreview}>ALERT PREVIEW</div>
        <div>
          <PrettyJSON jsonObject={alertConfig} />
        </div>
      </div>
    </>
  );
}

type ConfigProps = {|
  alertConfig: AlertConfig,
  setAlertConfig: ((AlertConfig => AlertConfig) | AlertConfig) => void,
  saveAlert: () => void,
|};

function ConfigurationStep(props: ConfigProps) {
  const [step, setStep] = useState(0);
  const onNext = () => setStep(step + 1);
  const onPrevious = () => setStep(step - 1);
  const {saveAlert, ...alertProps} = props;

  switch (step) {
    case 0:
      return <AddEditAlertConfigurationStep onNext={onNext} {...alertProps} />;
    case 1:
      return (
        <AddEditAlertInfoStep
          onNext={onNext}
          onPrevious={onPrevious}
          {...alertProps}
        />
      );
    case 2:
      return (
        <AddEditAlertNotificationStep
          onSave={saveAlert}
          onPrevious={onPrevious}
          {...alertProps}
        />
      );
    default:
      return <AddEditAlertConfigurationStep onNext={onNext} {...alertProps} />;
  }
}
