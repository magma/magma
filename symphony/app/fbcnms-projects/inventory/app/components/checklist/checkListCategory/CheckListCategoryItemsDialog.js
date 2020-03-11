/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ChecklistItemsDialogMutateStateActionType} from './ChecklistItemsDialogMutateAction';
import type {ChecklistItemsDialogStateType} from './ChecklistItemsDialogMutateState';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import CheckListTable from '../CheckListTable';
import ChecklistItemsDialogMutateDispatchContext from './ChecklistItemsDialogMutateDispatchContext';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Strings from '../../../common/CommonStrings';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '@fbcnms/ui/components/design-system/Text';
import fbt from 'fbt';
import {getInitialState, reducer} from './ChecklistItemsDialogMutateReducer';
import {useReducer, useState} from 'react';

type Props = {
  isOpened?: boolean,
  onCancel?: () => void,
  onSave?: (items: ChecklistItemsDialogStateType) => void,
  categoryTitle: string,
  initialItems: ChecklistItemsDialogStateType,
};

type View = {
  label: string,
  labelSuffix: (?ChecklistItemsDialogStateType) => string,
  value: number,
};

const DESIGN_VIEW: View = {
  label: `${fbt('items', 'Header for tab showing checklist items')}`,
  labelSuffix: itemsList => (itemsList ? ` (${itemsList.length})` : ''),
  value: 1,
};
const RESPONSE_VIEW: View = {
  label: `${fbt(
    'Responses',
    'Header for tab showing checklist response items',
  )}`,
  labelSuffix: itemsList =>
    itemsList
      ? ` (${itemsList.reduce(
          (responsesCount: number, clItem) =>
            clItem.checked ? responsesCount + 1 : responsesCount,
          0,
        )})`
      : '',
  value: 2,
};
const VIEWS = [DESIGN_VIEW, RESPONSE_VIEW];

const CheckListCategoryItemsDialog = ({
  initialItems,
  onCancel,
  onSave,
  categoryTitle,
}: Props) => {
  const [editingItems, dispatch] = useReducer<
    ChecklistItemsDialogStateType,
    ChecklistItemsDialogMutateStateActionType,
    ChecklistItemsDialogStateType,
  >(reducer, initialItems, getInitialState);

  const [pickedView, setPickedView] = useState<number>(DESIGN_VIEW.value);
  return (
    <Dialog fullWidth={true} maxWidth="md" open={true}>
      <DialogTitle disableTypography={true}>
        <Text variant="h6">
          <fbt desc="">Checklist</fbt>
          {` / ${categoryTitle}`}
        </Text>
      </DialogTitle>
      <DialogContent>
        <Tabs
          value={pickedView}
          onChange={(_e, newValue: number) => setPickedView(newValue)}
          indicatorColor="primary">
          {VIEWS.map(view => (
            <Tab
              value={view.value}
              label={`${view.label}${view.labelSuffix(editingItems)}`}
            />
          ))}
        </Tabs>
        <ChecklistItemsDialogMutateDispatchContext.Provider value={dispatch}>
          <CheckListTable
            items={editingItems}
            onDesignMode={pickedView === DESIGN_VIEW.value}
          />
        </ChecklistItemsDialogMutateDispatchContext.Provider>
      </DialogContent>
      <DialogActions>
        <Button skin="gray" onClick={onCancel}>
          {Strings.common.cancelButton}
        </Button>
        <Button onClick={() => onSave && onSave(editingItems)}>
          {Strings.common.saveButton}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default CheckListCategoryItemsDialog;
