/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
'use strict';

import type {FileItem} from '@fbcnms/ui/components/design-system/Experimental/FilesUploadSnackbar';

import * as React from 'react';
import emptyFunction from '@fbcnms/util/emptyFunction';
import {FileUploadStatuses} from '@fbcnms/ui/components/design-system/Experimental/FileUploadStatus';
import {Map as immMap} from 'immutable';
import {useState} from 'react';

type FilesStore = immMap<string, FileItem>;

type Context = {
  addFile: (id: string, name: string) => void,
  setFileProgress: (id: string, progress: number) => void,
  setFileUploadError: (id: string, errorMessage: React.Node) => void,
  files: FilesStore,
};

export const FilesUploadContext = React.createContext<Context>({
  addFile: emptyFunction,
  setFileProgress: emptyFunction,
  setFileUploadError: emptyFunction,
  files: new immMap<string, FileItem>(),
});

type Props = {
  children: React.Node,
};

export default function FilesUploadContextProvider({children}: Props) {
  const [files, setFiles] = useState(new immMap<string, FileItem>());
  const getStatus = (progress: number) =>
    progress === 100 ? FileUploadStatuses.DONE : FileUploadStatuses.UPLOADING;

  const value = {
    files: files,
    addFile: (id, name) =>
      setFiles(prevFiles =>
        prevFiles.set(id, {
          id,
          name,
          status: FileUploadStatuses.UPLOADING,
        }),
      ),
    setFileProgress: (id, progress) =>
      setFiles(prevFiles =>
        prevFiles.set(id, {
          ...prevFiles.get(id),
          status: getStatus(progress),
        }),
      ),
    setFileUploadError: (id, errorMessage) =>
      setFiles(prevFiles =>
        prevFiles.set(id, {
          ...prevFiles.get(id),
          status: 'error',
          errorMessage,
        }),
      ),
  };

  return (
    <FilesUploadContext.Provider value={value}>
      {children}
    </FilesUploadContext.Provider>
  );
}
