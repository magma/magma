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

import type {AssociateProp} from './AssociateTypes.flow';
import type {DataTypes, Model} from 'sequelize';

export type AuditLogEntryRawType = {
  actingUserId: number,
  organization: string,
  mutationType: 'CREATE' | 'UPDATE' | 'DELETE',
  objectId: string,
  objectType: string,
  objectDisplayName: string,
  mutationData: {[string]: mixed},
  url: string,
  ipAddress: string,
  status: 'SUCCESS' | 'FAILURE',
  statusCode: string,
};

type AuditLogEntryReadAttributes = AuditLogEntryRawType & {
  id: number,
};

type AuditLogEntryModel = Model<
  AuditLogEntryReadAttributes,
  AuditLogEntryRawType,
>;
export type StaticAuditLogEntryModel = Class<AuditLogEntryModel>;
export type AuditLogEntryType = AuditLogEntryModel & AuditLogEntryRawType;

export default (
  sequelize: Sequelize,
  types: DataTypes,
): StaticAuditLogEntryModel & AssociateProp => {
  return sequelize.define(
    'AuditLogEntry',
    {
      actingUserId: {
        type: types.INTEGER,
        allowNull: false,
      },
      organization: {
        type: types.STRING,
        allowNull: false,
      },
      mutationType: {
        type: types.STRING,
        allowNull: false,
      },
      objectId: {
        type: types.STRING,
        allowNull: false,
      },
      objectDisplayName: {
        type: types.STRING,
        allowNull: false,
      },
      objectType: {
        type: types.STRING,
        allowNull: false,
      },
      mutationData: {
        type: types.JSON,
        allowNull: false,
      },
      url: {
        type: types.STRING,
        allowNull: false,
      },
      ipAddress: {
        type: types.STRING,
        allowNull: false,
      },
      status: {
        type: types.STRING,
        allowNull: false,
      },
      statusCode: {
        type: types.STRING,
        allowNull: false,
      },
    },
    {},
  );
};
