/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Sequelize from 'sequelize';

import type {AssociateProp} from './AssociateTypes.flow';
import type {DataTypes, Model} from 'sequelize';

export type FeatureFlagRawType = {|
  featureId: string,
  organization: string,
  enabled: boolean,
|};

type FeatureFlagReadAttributes = {|
  ...FeatureFlagRawType,
  id: number,
|};

type FeatureFlagModel = Model<FeatureFlagReadAttributes, FeatureFlagRawType>;
export type StaticFeatureFlagModel = Class<FeatureFlagModel>;
export type FeatureFlagType = FeatureFlagModel & FeatureFlagRawType;

export default (
  sequelize: Sequelize,
  types: DataTypes,
): AssociateProp & StaticFeatureFlagModel => {
  return sequelize.define(
    'FeatureFlag',
    {
      featureId: {
        type: types.STRING,
        allowNull: false,
      },
      organization: {
        type: types.STRING,
        allowNull: false,
      },
      enabled: {
        type: types.BOOLEAN,
        allowNull: false,
      },
    },
    {},
  );
};
