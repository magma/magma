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
import ChecklistItemsDialogMutateDispatchContext from '../checkListCategory/ChecklistItemsDialogMutateDispatchContext';
import FilePreview from '../../FilePreview/FilePreview';
import FileUploadArea from '@fbcnms/ui/components/design-system/Experimental/FileUpload/FileUploadArea';
import FileUploadButton from '../../FileUpload/FileUploadButton';
import PendingFilePreview from '../../FilePreview/PendingFilePreview';
import {generateTempId} from '../../../common/EntUtils';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useContext, useMemo} from 'react';

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
  const dispatch = useContext(ChecklistItemsDialogMutateDispatchContext);
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

  const updatePendingFile = useCallback(
    (fileId: string, name: string, progress: number) => {
      dispatch({
        type: 'EDIT_ITEM_PENDING_FILE',
        itemId: item.id,
        file: {
          id: fileId,
          name,
          progress,
        },
      });
    },
    [dispatch, item.id],
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
      {item.pendingFiles?.map(pendingFile => (
        <PendingFilePreview
          className={classes.filePreview}
          fileName={pendingFile.name}
          progress={pendingFile.progress}
        />
      ))}
      <FileUploadButton
        uploadUsingSnackbar={false}
        onProgress={(fileId, file, progress) => {
          updatePendingFile(fileId, file.name, progress);
        }}
        onFileUploaded={(file, storeKey) =>
          dispatch({
            type: 'ADD_ITEM_FILE',
            itemId: item.id,
            file: {
              id: generateTempId(),
              storeKey,
              fileName: file.name,
              sizeInBytes: file.size,
              modificationTime: new Date().getTime(),
              uploadTime: new Date().getTime(),
            },
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
