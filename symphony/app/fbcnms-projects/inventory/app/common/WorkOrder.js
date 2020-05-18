/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {AddWorkOrderTypeInput} from '../mutations/__generated__/AddWorkOrderTypeMutation.graphql';
import type {CheckListCategoryExpandingPanel_list} from '../components/checklist/checkListCategory/__generated__/CheckListCategoryExpandingPanel_list.graphql';
import type {ChecklistCategoryDefinition} from '../components/checklist/ChecklistCategoriesMutateState';
import type {Equipment, Link} from './Equipment';
import type {FileAttachmentType} from './FileAttachment.js';
import type {ImageAttachmentType} from './ImageAttachment.js';
import type {Location} from './Location';
import type {Property} from './Property';
import type {PropertyType} from './PropertyType';
import type {ShortUser} from './EntUtils';

import {convertPropertyTypeToMutationInput} from './PropertyType';
import {isTempId} from './EntUtils';

export type WorkOrderStatus = 'PENDING' | 'PLANNED' | 'DONE';
export type WorkOrderPriority = 'URGENT' | 'HIGH' | 'LOW' | 'NONE';

export type WorkOrderType = {
  id: string,
  name: string,
  description: ?string,
  propertyTypes: Array<PropertyType>,
  numberOfWorkOrders: number,
  checklistCategoryDefinitions: Array<ChecklistCategoryDefinition>,
};

export type WorkOrder = {
  id: string,
  workOrderType: ?WorkOrderType,
  workOrderTypeId: ?string,
  name: string,
  description: ?string,
  location: ?Location,
  locationId: ?string,
  owner: ShortUser,
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
  assignedTo: ?ShortUser,
  properties: Array<Property>,
  projectId: ?string,
  checkListCategories: ?CheckListCategoryExpandingPanel_list,
};

export type WorkOrderIdentifier = {
  +id: string,
  +name: string,
};

export const priorityValues = [
  {
    key: 'urgent',
    value: 'URGENT',
    label: 'Urgent',
  },
  {
    key: 'high',
    value: 'HIGH',
    label: 'High',
  },
  {
    key: 'medium',
    value: 'MEDIUM',
    label: 'Medium',
  },
  {
    key: 'low',
    value: 'LOW',
    label: 'Low',
  },
  {
    key: 'none',
    value: 'NONE',
    label: 'None',
  },
];

export const doneStatus = {
  key: 'done',
  value: 'DONE',
  label: 'Done',
};

export const statusValues = [
  {
    key: 'planned',
    value: 'PLANNED',
    label: 'Planned',
  },
  {
    key: 'pending',
    value: 'PENDING',
    label: 'Pending',
  },
  doneStatus,
];

export type FutureState = 'INSTALL' | 'REMOVE';

export const FutureStateValues = [
  {
    key: 'install',
    value: 'INSTALL',
    label: 'Install',
  },
  {
    key: 'remove',
    value: 'REMOVE',
    label: 'Remove',
  },
];

export const convertWorkOrderTypeToMutationInput = (
  workOrderType: WorkOrderType,
): AddWorkOrderTypeInput => {
  return {
    name: workOrderType.name,
    description: workOrderType.description,
    properties: convertPropertyTypeToMutationInput(workOrderType.propertyTypes),
    checkListCategories: workOrderType.checklistCategoryDefinitions
      .slice()
      .map(categoryDef => ({
        ...categoryDef,
        id: isTempId(categoryDef.id) ? undefined : categoryDef.id,
        checkList: categoryDef.checkList.slice().map(item => ({
          ...item,
          id: isTempId(item.id) ? undefined : item.id,
        })),
      })),
  };
};
