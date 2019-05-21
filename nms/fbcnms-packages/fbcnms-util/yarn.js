/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

const _ = require('lodash');
const fs = require('fs');
const glob = require('glob');
const path = require('path');

type Dependencies = {
  [key: string]: string,
};

type Manifest = {
  dependencies?: Dependencies,
  devDependencies?: Dependencies,
  peerDependencies?: Dependencies,
  optionalDependencies?: Dependencies,
  workspaces?: Array<string>,
};

export function resolveWorkspaces(
  root: string,
  rootManifest: Manifest,
): Manifest[] {
  if (!rootManifest.workspaces) {
    return [];
  }

  const files = rootManifest.workspaces.map(pattern =>
    glob.sync(pattern.replace(/\/?$/, '/+(package.json)'), {
      cwd: root,
      ignore: pattern.replace(/\/?$/, '/node_modules/**/+(package.json)'),
    }),
  );

  return _.flatten(files).map(file => readManifest(path.join(root, file)));
}

export function readManifest(file: string): Manifest {
  return JSON.parse(fs.readFileSync(file, 'utf8'));
}
