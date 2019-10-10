/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import CircularProgress from '@material-ui/core/CircularProgress';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';
import generateQRCode from './generateQRCode';
import red from '@material-ui/core/colors/red';
import {makeStyles} from '@material-ui/styles';
import {useAxios} from '@fbcnms/ui/hooks';

export type Props = {|
  endpoint: string,
|};

const useStyles = makeStyles(theme => ({
  errorAlert: {
    marginTop: theme.spacing(),
    padding: theme.spacing(),
    backgroundColor: red[100],
  },
  errorMessage: {color: theme.palette.error.dark},
  loadingSpinner: {
    padding: theme.spacing(2),
    textAlign: 'center',
  },
}));

export default function FBCMobileAppQRCode(props: Props) {
  const classes = useStyles();
  const {isLoading, error, response} = useAxios({
    method: 'GET',
    url: props.endpoint,
  });
  const [code, setCode] = React.useState(null);
  React.useEffect(() => {
    if (response && !error) {
      generateQRCode(JSON.stringify(response.data)).then(qr => {
        setCode(qr);
      });
    }
  }, [response, error]);

  if (error) {
    return (
      <Paper className={classes.errorAlert} elevation={0}>
        <Typography
          className={classes.errorMessage}
          data-testid="error-message">
          Could not load QR Code data
        </Typography>
      </Paper>
    );
  }

  if (isLoading || !response || !code) {
    return <CircularProgress data-testid="loading" />;
  }
  return (
    <img
      data-testid="qrcode"
      src={code}
      style={{
        height: 250,
        width: 250,
      }}
    />
  );
}
