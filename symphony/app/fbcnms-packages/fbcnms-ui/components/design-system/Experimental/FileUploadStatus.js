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
import CheckCircleIcon from '@material-ui/icons/CheckCircle';
import CircularProgress from '@material-ui/core/CircularProgress';
import ErrorIcon from '@material-ui/icons/Error';
import FileIcon from '../Icons/Indications/FileIcon';
import Text from '../Text';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  root: {
    height: '52px',
    boxSizing: 'border-box',
    maxHeight: '64px',
    padding: '0px 8px',
    display: 'flex',
    alignItems: 'center',
  },
  content: {
    display: 'flex',
    flexDirection: 'column',
    flexGrow: 1,
    overflow: 'hidden',
    marginRight: '16px',
  },
  name: {
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
    overflow: 'hidden',
  },
  icon: {
    width: '20px',
    height: '20px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
  fileIcon: {
    marginRight: '8px',
  },
  errorMessage: {
    marginTop: '4px',
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
});

export type UploadStatus = 'uploading' | 'error' | 'done';

type Props = {
  name: React.Node,
  status: UploadStatus,
  errorMessage?: React.Node,
};

const StatusIcon = ({status}) => {
  const classes = useStyles();
  if (status === 'uploading') {
    return (
      <div className={classes.icon}>
        <CircularProgress color="primary" size={16.67} />
      </div>
    );
  }

  return status === 'done' ? (
    <CheckCircleIcon fontSize="small" color="primary" />
  ) : (
    <ErrorIcon fontSize="small" color="error" />
  );
};

const FileUploadStatus = ({name, status, errorMessage}: Props) => {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <FileIcon color="primary" className={classes.fileIcon} />
      <div className={classes.content}>
        <Text variant="body2" className={classes.name}>
          {name}
        </Text>
        {errorMessage && (
          <Text
            color="error"
            variant="caption"
            className={classes.errorMessage}>
            {errorMessage}
          </Text>
        )}
      </div>
      <StatusIcon status={status} />
    </div>
  );
};

export default FileUploadStatus;
