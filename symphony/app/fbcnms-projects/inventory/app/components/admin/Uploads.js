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
import CSVFileUpload from '../CSVFileUpload';
import LoadingFillerBackdrop from '@fbcnms/ui/components/LoadingFillerBackdrop';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import {UploadAPIUrls} from '../../common/UploadAPI';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useState} from 'react';

const useStyles = makeStyles(_ => ({
  button: {
    margin: '10px',
  },
  paper: {
    margin: '10px',
  },
}));

export default function Uploads() {
  const classes = useStyles();
  const [isLoading, setIsLoading] = useState(false);
  const enqueueSnackbar = useEnqueueSnackbar();

  const onFileUploaded = () => {
    enqueueSnackbar('Upload successful', {variant: 'success'});
    setIsLoading(false);
  };

  const onUploadFailed = error => {
    enqueueSnackbar(`Upload failed ${error}`, {variant: 'error'});
    setIsLoading(false);
  };

  return (
    <div className={classes.paper}>
      <Paper className={classes.tableRoot} elevation={2}>
        <CSVFileUpload
          button={<Button className={classes.button}>Rural RAN</Button>}
          onProgress={() => setIsLoading(true)}
          onFileUploaded={onFileUploaded}
          onUploadFailed={onUploadFailed}
          uploadPath={UploadAPIUrls.rural_ran()}
        />
        <CSVFileUpload
          button={<Button className={classes.button}>Rural Transport</Button>}
          onProgress={() => setIsLoading(true)}
          onFileUploaded={onFileUploaded}
          onUploadFailed={onUploadFailed}
          uploadPath={UploadAPIUrls.rural_transport()}
        />
        <CSVFileUpload
          button={<Button className={classes.button}>Upload FTTH</Button>}
          onProgress={() => setIsLoading(true)}
          onFileUploaded={onFileUploaded}
          onUploadFailed={onUploadFailed}
          uploadPath={UploadAPIUrls.ftth()}
        />
        <CSVFileUpload
          button={
            <Button className={classes.button}>Express Wi-Fi Rural</Button>
          }
          onProgress={() => setIsLoading(true)}
          onFileUploaded={onFileUploaded}
          onUploadFailed={onUploadFailed}
          uploadPath={UploadAPIUrls.xwf1()}
        />
        <CSVFileUpload
          button={
            <Button className={classes.button}>
              Express Wi-Fi XPP Access Points
            </Button>
          }
          onProgress={() => setIsLoading(true)}
          onFileUploaded={onFileUploaded}
          onUploadFailed={onUploadFailed}
          uploadPath={UploadAPIUrls.xwfAps()}
        />
      </Paper>
      {isLoading && <LoadingFillerBackdrop />}
    </div>
  );
}
