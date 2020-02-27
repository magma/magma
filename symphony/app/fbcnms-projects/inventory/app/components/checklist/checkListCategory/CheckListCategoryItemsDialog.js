/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CheckListCategoryItemsDialog_items} from './__generated__/CheckListCategoryItemsDialog_items.graphql';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import CheckListTable from '../CheckListTable';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Strings from '../../../common/CommonStrings';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '@fbcnms/ui/components/design-system/Text';
import fbt from 'fbt';
import {createFragmentContainer, graphql} from 'react-relay';
import {useState} from 'react';

type Props = {
  isOpened?: boolean,
  onClose?: () => void,
  categoryTitle: string,
  items: ?CheckListCategoryItemsDialog_items,
  onChecklistChanged?: (
    updatedList: CheckListCategoryItemsDialog_items,
  ) => void,
};

type View = {
  label: string,
  labelSuffix: (?CheckListCategoryItemsDialog_items) => string,
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
          (responsesCount, clItem) =>
            clItem.checked ? responsesCount + 1 : responsesCount,
          0,
        )})`
      : '',
  value: 2,
};
const VIEWS = [DESIGN_VIEW, RESPONSE_VIEW];

const CheckListCategoryItemsDialog = (props: Props) => {
  const {items} = props;
  const [pickedView, setPickedView] = useState(DESIGN_VIEW.value);
  return (
    <Dialog fullWidth={true} maxWidth="md" open={true}>
      <DialogTitle disableTypography={true}>
        <Text variant="h6">
          <fbt desc="">Checklist</fbt>
          {` / ${props.categoryTitle}`}
        </Text>
      </DialogTitle>
      <DialogContent>
        <Tabs
          value={pickedView}
          onChange={(_e, newValue) => setPickedView(newValue)}
          indicatorColor="primary">
          {VIEWS.map(view => (
            <Tab
              value={view.value}
              label={`${view.label}${view.labelSuffix(items)}`}
            />
          ))}
        </Tabs>
        <CheckListTable
          list={items}
          onChecklistChanged={props.onChecklistChanged}
          onDesignMode={pickedView === DESIGN_VIEW.value}
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>{Strings.common.closeButton}</Button>
      </DialogActions>
    </Dialog>
  );
};

export default createFragmentContainer(CheckListCategoryItemsDialog, {
  items: graphql`
    fragment CheckListCategoryItemsDialog_items on CheckListItem
      @relay(plural: true) {
      ...CheckListTable_list
      checked
    }
  `,
});
