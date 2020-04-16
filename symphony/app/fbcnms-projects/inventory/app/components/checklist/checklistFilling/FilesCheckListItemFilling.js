/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CheckListItemFillingProps} from './CheckListItemFilling';
import type {FileAttachmentType} from '../../../common/FileAttachment';

import * as React from 'react';
import FilePreview from '../../FilePreview/FilePreview';
import FileUploadArea from '@fbcnms/ui/components/design-system/Experimental/FileUpload/FileUploadArea';
import FileUploadButton from '../../FileUpload/FileUploadButton';
import {generateTempId} from '../../../common/EntUtils';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useMemo} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    flexDirection: 'row',
    flexWrap: 'wrap',
  },
  filePreview: {
    margin: '4px',
  },
  fileUploadArea: {
    margin: '4px',
  },
}));

const FilesCheckListItemFilling = ({
  item,
  onChange,
}: CheckListItemFillingProps): React.Node => {
  const classes = useStyles();
  const tempId = useMemo(() => generateTempId(), []);

  const removeItemFile = useCallback(
    (removedFile: FileAttachmentType) =>
      onChange &&
      onChange({
        ...item,
        files: item.files?.filter(file => file.id !== removedFile.id),
      }),
    [item, onChange],
  );

  return (
    <div className={classes.root}>
      {item.files?.map(file => (
        <FilePreview
          className={classes.filePreview}
          file={{
            id: file.id ?? tempId,
            fileName: file.fileName,
            sizeInBytes: file.sizeInBytes,
            modified: file.modificationTime
              ? `${file.modificationTime}`
              : undefined,
            uploaded: file.uploadTime ? `${file.uploadTime}` : undefined,
            storeKey: file.storeKey,
          }}
          onFileDeleted={removeItemFile}
        />
      ))}
      <FileUploadButton
        uploadUsingSnackbar={false}
        onProgress={(_fileId, _progress) => {
          // TODO: implement progress once there's design
        }}
        onFileUploaded={(file, storeKey) =>
          onChange &&
          onChange({
            ...item,
            files: [
              ...(item.files ?? []),
              {
                id: generateTempId(),
                storeKey,
                fileName: file.name,
                sizeInBytes: file.size,
                modificationTime: new Date().getTime(),
                uploadTime: new Date().getTime(),
              },
            ],
          })
        }>
        {openFileUploadDialog => (
          <FileUploadArea
            className={classes.fileUploadArea}
            icon="plus"
            onClick={openFileUploadDialog}
            dimensions="wide"
          />
        )}
      </FileUploadButton>
    </div>
  );
};

export default FilesCheckListItemFilling;
