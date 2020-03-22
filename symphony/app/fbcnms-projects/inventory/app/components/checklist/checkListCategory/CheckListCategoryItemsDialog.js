/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  CheckListItem,
  ChecklistItemsDialogStateType,
} from './ChecklistItemsDialogMutateState';
import type {ChecklistItemsDialogMutateStateActionType} from './ChecklistItemsDialogMutateAction';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import CheckListTableFilling from '../checklistFilling/CheckListTableFilling';
import ChecklistDefinitionsList from '../checklistDefinition/ChecklistDefinitionsList';
import ChecklistItemsDialogMutateDispatchContext from './ChecklistItemsDialogMutateDispatchContext';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import Strings from '../../../common/CommonStrings';
import TabsBar from '@fbcnms/ui/components/design-system/Tabs/TabsBar';
import Text from '@fbcnms/ui/components/design-system/Text';
import fbt from 'fbt';
import {PlusIcon} from '@fbcnms/ui/components/design-system/Icons';
import {getInitialState, reducer} from './ChecklistItemsDialogMutateReducer';
import {makeStyles} from '@material-ui/styles';
import {useReducer, useState} from 'react';

const useStyles = makeStyles(() => ({
  dialogHeader: {
    display: 'flex',
    flexDirection: 1,
    alignItems: 'center',
  },
  tabs: {
    flexGrow: 1,
  },
  dialogActions: {
    padding: '24px',
  },
}));

type Props = {
  isOpened?: boolean,
  onCancel?: () => void,
  onSave?: (items: Array<CheckListItem>) => void,
  categoryTitle: string,
  initialItems: Array<CheckListItem>,
};

const TabViewValues = {
  items: 0,
  responses: 1,
};

type TabViewValue = $Values<typeof TabViewValues>;

type View = {
  label: string,
  labelSuffix: (?Array<CheckListItem>) => string,
  value: TabViewValue,
};

const DESIGN_VIEW: View = {
  label: `${fbt('items', 'Header for tab showing checklist items')}`,
  labelSuffix: itemsList => (itemsList ? ` (${itemsList.length})` : ''),
  value: 0,
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
  value: 1,
};
const VIEWS = [DESIGN_VIEW, RESPONSE_VIEW];

const CheckListCategoryItemsDialog = ({
  initialItems,
  onCancel,
  onSave,
  categoryTitle,
}: Props) => {
  const classes = useStyles();
  const [dialogState, dispatch] = useReducer<
    ChecklistItemsDialogStateType,
    ChecklistItemsDialogMutateStateActionType,
    Array<CheckListItem>,
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
        <div className={classes.dialogHeader}>
          <TabsBar
            className={classes.tabs}
            tabs={VIEWS.map(view => ({
              label: `${view.label}${view.labelSuffix(dialogState.items)}`,
            }))}
            activeTabIndex={pickedView}
            onChange={setPickedView}
            spread={true}
          />
          {pickedView === TabViewValues.items && (
            <FormAction>
              <Button
                onClick={() => dispatch({type: 'ADD_ITEM'})}
                leftIcon={PlusIcon}>
                <fbt desc="">Add Item</fbt>
              </Button>
            </FormAction>
          )}
        </div>
        <ChecklistItemsDialogMutateDispatchContext.Provider value={dispatch}>
          {pickedView === TabViewValues.items ? (
            <ChecklistDefinitionsList
              items={dialogState.items}
              editedDefinitionId={dialogState.editedDefinitionId}
            />
          ) : (
            <CheckListTableFilling items={dialogState.items} />
          )}
        </ChecklistItemsDialogMutateDispatchContext.Provider>
      </DialogContent>
      <DialogActions className={classes.dialogActions}>
        <Button skin="gray" onClick={onCancel}>
          {Strings.common.cancelButton}
        </Button>
        <Button onClick={() => onSave && onSave(dialogState.items)}>
          {Strings.common.saveButton}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default CheckListCategoryItemsDialog;
