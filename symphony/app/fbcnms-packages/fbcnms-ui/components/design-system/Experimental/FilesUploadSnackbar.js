/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {UploadStatus} from './FileUploadStatus';

import * as React from 'react';
import CloseIcon from '../Icons/Navigation/CloseIcon';
import FileUploadStatus, {FileUploadStatuses} from './FileUploadStatus';
import Portal from '../Core/Portal';
import Text from '../Text';
import fbt from 'fbt';
import symphony from '../../../theme/symphony';
import useVerticalScrollingEffect from '../hooks/useVerticalScrollingEffect';
import {makeStyles} from '@material-ui/styles';
import {useRef} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    position: 'absolute',
    zIndex: 1200,
    left: '96px',
    bottom: '16px',
    backgroundColor: symphony.palette.white,
    width: '320px',
    boxShadow: symphony.shadows.DP1,
    borderRadius: '4px',
  },
  header: {
    display: 'flex',
    alignItems: 'center',
    borderRadius: '4px 4px 0px 0px',
    padding: '8px',
    backgroundColor: symphony.palette.D900,
  },
  headerText: {
    flexGrow: 1,
  },
  closeIcon: {
    width: '20px',
    height: '20px',
    cursor: 'pointer',
    '&:hover': {
      fill: symphony.palette.B300,
    },
  },
  content: {
    maxHeight: '270px',
    overflowY: 'auto',
  },
}));

export type FileItem = {
  id: string,
  name: React.Node,
  status: UploadStatus,
  errorMessage?: React.Node,
};

type Props = {
  files: Array<FileItem>,
  onClose: () => void,
};

const FilesUploadSnackbar = ({files, onClose}: Props) => {
  const classes = useStyles();
  const thisElement = useRef(null);
  useVerticalScrollingEffect(thisElement);

  return (
    <Portal target={document.body}>
      <div className={classes.root}>
        <div className={classes.header}>
          <Text variant="body2" color="light" className={classes.headerText}>
            {files.every(f => f.status === FileUploadStatuses.DONE) ? (
              <fbt desc="Amount of uploaded files">
                <fbt:param name="Total number of files" number={true}>
                  {files.length}
                </fbt:param>
                Uploads Complete
              </fbt>
            ) : files.every(
                f =>
                  f.status === FileUploadStatuses.DONE ||
                  f.status === FileUploadStatuses.ERROR,
              ) ? (
              <fbt desc="Amount of files uploading">
                <fbt:param name="Number of successfuly uploaded files">
                  {
                    files.filter(f => f.status === FileUploadStatuses.DONE)
                      .length
                  }
                </fbt:param>
                Files Uploaded (<fbt:param name="Total number of files">
                  {
                    files.filter(f => f.status === FileUploadStatuses.ERROR)
                      .length
                  }
                </fbt:param>{' '}
                Errors)
              </fbt>
            ) : (
              <fbt desc="Amount of files uploading">
                Uploading
                <fbt:param name="Number of successfuly uploaded files">
                  {
                    files.filter(f => f.status === FileUploadStatuses.UPLOADING)
                      .length
                  }
                </fbt:param>
                {' / '}
                <fbt:param name="Total number of files">
                  {files.length}
                </fbt:param>
              </fbt>
            )}
          </Text>
          <CloseIcon
            color="light"
            className={classes.closeIcon}
            onClick={onClose}
          />
        </div>
        <div className={classes.content} ref={thisElement}>
          {files.map(file => (
            <FileUploadStatus
              key={file.id}
              name={file.name}
              status={file.status}
              errorMessage={file.errorMessage}
            />
          ))}
        </div>
      </div>
    </Portal>
  );
};

export default FilesUploadSnackbar;
