/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import featureConfigs, {isFeatureEnabled} from '../features';
import {FeatureFlag, sequelize} from '@fbcnms/sequelize-models';

const ORG1 = 'org1';
const ORG2 = 'org2';
const FEATURE1 = featureConfigs[Object.keys(featureConfigs)[0]];
const FEATURE2 = featureConfigs[Object.keys(featureConfigs)[1]];

describe('feature tests', () => {
  beforeAll(async () => {
    await sequelize.sync({force: true});
    await FeatureFlag.create({
      featureId: FEATURE1.id,
      organization: ORG1,
      enabled: true,
    });

    await FeatureFlag.create({
      featureId: FEATURE2.id,
      organization: ORG2,
      enabled: false,
    });
  });

  it('should be enabled', async () => {
    const enabled = await isFeatureEnabled(FEATURE1.id, ORG1);
    expect(enabled).toBe(true);
  });

  it('should be disabled by default', async () => {
    const enabled = await isFeatureEnabled(FEATURE1.id, ORG2);
    expect(enabled).toBe(false);
  });

  it('should be disabled', async () => {
    const enabled = await isFeatureEnabled(FEATURE2.id, ORG2);
    expect(enabled).toBe(false);
  });

  it('should be enabled', async () => {
    const enabled = await isFeatureEnabled(FEATURE2.id, ORG1);
    expect(enabled).toBe(true);
  });
});
