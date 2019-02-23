/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 * @flow
 * @format
 */

const _ = require('lodash');
const path = require('path');
const yarn = require('@fbcnms/util/yarn');

it('ensures no mismatched versions in workspaces', () => {
  const root = path.resolve(__dirname, '..');
  const rootManifest = yarn.readManifest(path.resolve(root, 'package.json'));
  const workspaces = yarn.resolveWorkspaces(root, rootManifest);

  const allManifests = [rootManifest, ...workspaces];

  const allDepsMap = _.merge(
    {},
    ...allManifests.map(manifest => ({
      ...manifest.dependencies,
      ...manifest.devDependencies,
      ...manifest.optionalDependencies,
      ...manifest.peerDependencies,
    })),
  );

  for (const manifest of workspaces) {
    expect(allDepsMap).toMatchObject(manifest.dependencies || {});
    expect(allDepsMap).toMatchObject(manifest.devDependencies || {});
    expect(allDepsMap).toMatchObject(manifest.peerDependencies || {});
    expect(allDepsMap).toMatchObject(manifest.optionalDependencies || {});
  }
});
