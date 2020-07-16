/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FeatureID} from '../../fbcnms-types/features';

import AppContext from './AppContext';
import {useContext} from 'react';

const useFeatureFlag = (featureId: FeatureID) => {
  return useContext(AppContext).isFeatureEnabled(featureId);
};

export default useFeatureFlag;
