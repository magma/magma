/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

const fs = require('fs');
const path = require('path');
const paths = require('./paths');

const DEV_MODE = process.env.NODE_ENV !== 'production';
const MANIFEST_FILE = path.join(paths.appSrc, '../static/dist/manifest.json');

let manifest = null;
if (fs.existsSync(MANIFEST_FILE)) {
  const manifestRaw = fs
    .readFileSync(MANIFEST_FILE)
    .toString('utf8')
    .trim();
  manifest = JSON.parse(manifestRaw);
}
export default function staticDist(
  projectName: string,
  filename: string,
): ?string {
  if (DEV_MODE || !manifest) {
    return '/' + projectName + '/static/dist/' + filename;
  }
  return manifest[filename] || '/dev/null/' + filename;
}
