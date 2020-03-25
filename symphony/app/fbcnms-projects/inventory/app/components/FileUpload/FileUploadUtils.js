/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import axios from 'axios';

export async function uploadFile(
  id: string,
  file: File,
  onUpload: (File, string) => void,
  onProgress?: (fileId: string, progress: number) => void,
) {
  const signingResponse = await axios.get('/store/put', {
    params: {
      contentType: file.type,
    },
  });

  const config = {
    headers: {
      'Content-Type': file.type,
    },
    onUploadProgress: function (progressEvent) {
      const percentCompleted = Math.round(
        (progressEvent.loaded * 100) / progressEvent.total,
      );
      onProgress && onProgress(id, percentCompleted);
    },
  };
  await axios.put(signingResponse.data.URL, file, config);

  onUpload(file, signingResponse.data.key);
}
