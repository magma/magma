/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Context} from 'react';
import type {FeatureID} from 'inventory/server/features';

import {useContext} from 'react';

type ContextType = {enabledFeatures: FeatureID[]} & {[string]: any};

export default function(
  appContext: Context<ContextType>,
  featureId: FeatureID,
): boolean {
  const {enabledFeatures} = useContext(appContext);
  return enabledFeatures.indexOf(featureId) !== -1;
}
