/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import FilesUploadSnackbar from '@fbcnms/ui/components/design-system/Experimental/FilesUploadSnackbar';
import React, {useContext, useEffect, useState} from 'react';
import {FilesUploadContext} from './context/FilesUploadContextProvider';

const SymphonyFilesUploadSnackbar = () => {
  const [isShown, setIsShown] = useState(false);
  const filesUploadContext = useContext(FilesUploadContext);

  useEffect(() => {
    setIsShown(filesUploadContext.files.size > 0);
  }, [filesUploadContext.files.size]);

  if (!isShown) {
    return null;
  }

  return (
    <FilesUploadSnackbar
      onClose={() => setIsShown(false)}
      files={filesUploadContext.files.toArray()}
    />
  );
};

export default SymphonyFilesUploadSnackbar;
