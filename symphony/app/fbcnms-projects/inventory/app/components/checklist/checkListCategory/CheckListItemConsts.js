/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CheckListItemDefinitionProps} from '../checklistDefinition/CheckListItemDefinition';
import type {CheckListItemFillingProps} from '../checklistFilling/CheckListItemFilling';
import type {CheckListItemType} from '../../work_orders/__generated__/WorkOrderDetails_workOrder.graphql';
import type {SvgIconStyleProps} from '@fbcnms/ui/components/design-system/Icons/SvgIcon';

import * as React from 'react';
import BasicCheckListItemDefinition from '../checklistDefinition/BasicCheckListItemDefinition';
import BasicCheckListItemFilling from '../checklistFilling/BasicCheckListItemFilling';
import CellScanCheckListItemFilling from '../checklistFilling/CellScanCheckListItemFilling';
import FilesCheckListItemDefinition from '../checklistDefinition/FilesCheckListItemDefinition';
import FilesCheckListItemFilling from '../checklistFilling/FilesCheckListItemFilling';
import FreeTextCheckListItemDefinition from '../checklistDefinition/FreeTextCheckListItemDefinition';
import FreeTextCheckListItemFilling from '../checklistFilling/FreeTextCheckListItemFilling';
import MultipleChoiceCheckListItemDefinition from '../checklistDefinition/MultipleChoiceCheckListItemDefinition';
import MultipleChoiceCheckListItemFilling from '../checklistFilling/MultipleChoiceCheckListItemFilling';
import WifiScanCheckListItemFilling from '../checklistFilling/WifiScanCheckListItemFilling';
import YesNoCheckListItemDefinition from '../checklistDefinition/YesNoCheckListItemDefinition';
import YesNoCheckListItemFilling from '../checklistFilling/YesNoCheckListItemFilling';
import fbt from 'fbt';
import {
  AttachmentIcon,
  CellularIcon,
  ChecklistCheckIcon,
  MultipleSelectionIcon,
  TextIcon,
  WifiIcon,
  YesNoIcon,
} from '@fbcnms/ui/components/design-system/Icons';

export type CheckListItemConfigsType = {
  [CheckListItemType]: {|
    icon: React.ComponentType<SvgIconStyleProps>,
    definitionComponent: React.ComponentType<CheckListItemDefinitionProps>,
    fillingComponent: React.ComponentType<CheckListItemFillingProps>,
    selectLabel: React.Node,
    titlePlaceholder: string,
  |},
};

export const CheckListItemConfigs: CheckListItemConfigsType = {
  simple: {
    icon: ChecklistCheckIcon,
    definitionComponent: BasicCheckListItemDefinition,
    fillingComponent: BasicCheckListItemFilling,
    selectLabel: <fbt desc="">Check when complete</fbt>,
    titlePlaceholder: `${fbt('What needs to be done?', '')}`,
  },
  string: {
    icon: TextIcon,
    definitionComponent: FreeTextCheckListItemDefinition,
    fillingComponent: FreeTextCheckListItemFilling,
    selectLabel: <fbt desc="">Free text</fbt>,
    titlePlaceholder: `${fbt('What needs to be written?', '')}`,
  },
  enum: {
    icon: MultipleSelectionIcon,
    definitionComponent: MultipleChoiceCheckListItemDefinition,
    fillingComponent: MultipleChoiceCheckListItemFilling,
    selectLabel: <fbt desc="">Multiple choice</fbt>,
    titlePlaceholder: `${fbt('What needs to be chosen?', '')}`,
  },
  files: {
    icon: AttachmentIcon,
    definitionComponent: FilesCheckListItemDefinition,
    fillingComponent: FilesCheckListItemFilling,
    selectLabel: <fbt desc="">Upload files</fbt>,
    titlePlaceholder: `${fbt('What needs to be uploaded?', '')}`,
  },
  yes_no: {
    icon: YesNoIcon,
    definitionComponent: YesNoCheckListItemDefinition,
    fillingComponent: YesNoCheckListItemFilling,
    selectLabel: <fbt desc="">Yes/No</fbt>,
    titlePlaceholder: `${fbt('Write your yes/no question', '')}`,
  },
  cell_scan: {
    icon: CellularIcon,
    definitionComponent: BasicCheckListItemDefinition,
    fillingComponent: CellScanCheckListItemFilling,
    selectLabel: <fbt desc="">Cellular Scan</fbt>,
    titlePlaceholder: `${fbt('Scan cellular signals', '')}`,
  },
  wifi_scan: {
    icon: WifiIcon,
    definitionComponent: BasicCheckListItemDefinition,
    fillingComponent: WifiScanCheckListItemFilling,
    selectLabel: <fbt desc="">Wi-Fi Scan</fbt>,
    titlePlaceholder: `${fbt('Scan Wi-Fi signals', '')}`,
  },
};
