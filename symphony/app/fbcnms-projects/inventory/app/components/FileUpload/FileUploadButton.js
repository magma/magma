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
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import fbt from 'fbt';
import shortid from 'shortid';
import {FilesUploadContext} from '../context/FilesUploadContextProvider';
import {makeStyles} from '@material-ui/styles';
import {uploadFile} from './FileUploadUtils';
import {useCallback, useContext, useRef} from 'react';

const MAX_FILES_LIMIT = 15;

const useStyles = makeStyles(() => ({
  hiddenInput: {
    width: '0px',
    height: '0px',
    opacity: 0,
    overflow: 'hidden',
    position: 'absolute',
    zIndex: -1,
  },
}));

type Props = {
  children: (openFileUploadDialog: () => void) => React.Node,
  className?: ?string,
  onProgress?: (fileId: string, progress: number) => void,
  onFileUploaded: (file: File, key: string) => void,
  multiple?: boolean,
  fileTypes?: string,
  uploadUsingSnackbar?: boolean,
  uploadType?: 'locally' | 'upload',
};

const FileUploadButton = ({
  children,
  onProgress,
  onFileUploaded,
  multiple = true,
  fileTypes = 'file',
  uploadUsingSnackbar = true,
  uploadType = 'upload',
}: Props) => {
  const classes = useStyles();
  const uploadContext = useContext(FilesUploadContext);
  const inputRef = useRef();
  const buttonClick = useCallback(() => inputRef?.current?.click(), [inputRef]);

  const onFileProgress = (fileId, progress) => {
    uploadUsingSnackbar && uploadContext.setFileProgress(fileId, progress);
    onProgress && onProgress(fileId, progress);
  };

  const onFilesChanged = async (e: SyntheticEvent<HTMLInputElement>) => {
    const eventFiles = Array.from(e.currentTarget.files).slice(
      0,
      MAX_FILES_LIMIT,
    );
    if (!eventFiles || eventFiles.length === 0) {
      return;
    }

    await Promise.all(
      eventFiles.map(async file => {
        const fileId = shortid.generate();
        uploadUsingSnackbar && uploadContext.addFile(fileId, file.name);
        try {
          if (uploadType === 'upload') {
            await uploadFile(fileId, file, onFileUploaded, onFileProgress);
          } else {
            onFileUploaded(file, fileId);
          }
        } catch (e) {
          uploadUsingSnackbar &&
            uploadContext.setFileUploadError(
              fileId,
              fbt(
                'We had a problem uploading this file',
                'Error message describing that we had an error while uploading the file',
              ),
            );
        }
      }),
    );
  };
  return (
    <FormAction>
      <input
        className={classes.hiddenInput}
        type={fileTypes}
        onChange={async e => await onFilesChanged(e)}
        ref={inputRef}
        multiple={multiple}
        fileTypes={fileTypes}
      />
      {children(buttonClick)}
    </FormAction>
  );
};

export default FileUploadButton;
