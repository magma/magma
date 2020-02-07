/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

export type FileAttachmentType = {
  id: string,
  fileName: string,
  sizeInBytes: number,
  modified: string,
  uploaded: string,
  fileType: string,
  storeKey: string,
  category: string,
};
