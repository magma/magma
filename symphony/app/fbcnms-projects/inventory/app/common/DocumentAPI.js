/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

export const DocumentAPIUrls = {
  get_url: (documentId: string) => `/store/get?key=${documentId}`,
  download_url: (documentId: string, fileName: string) =>
    `/store/download?key=${documentId}&fileName=${fileName}`,
  delete_url: (documentId: string) => `/store/delete?key=${documentId}`,
};
