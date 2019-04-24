/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Sequelize from 'sequelize';

import type {DataTypes, Model} from 'sequelize';
import type {AssociateProp} from './AssociateTypes.flow';

export type FeatureFlagRawType = {
  featureId: string,
  organizationId: number,
  enabled: boolean,
};

type FeatureFlagModel = Model<FeatureFlagRawType>;
export type StaticFeatureFlagModel = Class<FeatureFlagModel>;
export type FeatureFlagType = StaticFeatureFlagModel & FeatureFlagRawType;

export default (
  sequelize: Sequelize,
  types: DataTypes,
): StaticFeatureFlagModel & AssociateProp => {
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
