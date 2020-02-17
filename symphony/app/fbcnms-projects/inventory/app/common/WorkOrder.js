/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {CheckListCategoryTable_list} from '../components/checklist/checkListCategory/__generated__/CheckListCategoryTable_list.graphql';
import type {CheckListTable_list} from '../components/checklist/__generated__/CheckListTable_list.graphql';
import type {Equipment, Link} from './Equipment';
import type {FileAttachmentType} from './FileAttachment.js';
import type {ImageAttachmentType} from './ImageAttachment.js';
import type {Location} from './Location';
import type {Property} from './Property';
import type {PropertyType} from './PropertyType';

export type WorkOrderStatus = 'PENDING' | 'PLANNED' | 'DONE';
export type WorkOrderPriority = 'URGENT' | 'HIGH' | 'LOW' | 'NONE';

export type WorkOrderType = {
  id: string,
  name: string,
  description: ?string,
  propertyTypes: Array<PropertyType>,
  numberOfWorkOrders: number,
  checkListDefinitions: CheckListTable_list,
};

export type WorkOrder = {
  id: string,
  workOrderType: ?WorkOrderType,
  workOrderTypeId: ?string,
  name: string,
  description: ?string,
  location: ?Location,
  locationId: ?string,
  ownerName: string,
  creationDate: string,
  installDate: ?string,
  status: WorkOrderStatus,
  priority: WorkOrderPriority,
  equipmentToAdd: Array<Equipment>,
  equipmentToRemove: Array<Equipment>,
  linksToAdd: Array<Link>,
  linksToRemove: Array<Link>,
  images: Array<ImageAttachmentType>,
  files: Array<FileAttachmentType>,
  assignee: ?string,
  properties: Array<Property>,
  projectId: ?string,
  checkList: ?CheckListTable_list,
  checkListCategories: ?CheckListCategoryTable_list,
};

export type WorkOrderIdentifier = {
  +id: string,
  +name: string,
};

export const priorityValues = [
  {
    value: 'URGENT',
    label: 'Urgent',
  },
  {
    value: 'HIGH',
    label: 'High',
  },
  {
    value: 'MEDIUM',
    label: 'Medium',
  },
  {
    value: 'LOW',
    label: 'Low',
  },
  {
    value: 'NONE',
    label: 'None',
  },
];

export const doneStatus = {
  value: 'DONE',
  label: 'Done',
};

export const statusValues = [
  {
    value: 'PLANNED',
    label: 'Planned',
  },
  {
    value: 'PENDING',
    label: 'Pending',
  },
  doneStatus,
];

export type FutureState = 'INSTALL' | 'REMOVE';

export const FutureStateValues = [
  {
    value: 'INSTALL',
    label: 'Install',
  },
  {
    value: 'REMOVE',
    label: 'Remove',
  },
];
