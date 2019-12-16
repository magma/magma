/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import CheckListTableDefinition from './checklistDefinition/CheckListTableDefinition';
import CheckListTableFilling from './checklistFilling/CheckListTableFilling';
import React from 'react';
import {CHECKLIST_ITEM_TYPES} from './CheckListItem';
import {createFragmentContainer, graphql} from 'react-relay';
import {sortByIndex} from '../draggable/DraggableUtils';
import type {CheckListTable_list} from './__generated__/CheckListTable_list.graphql';

type Props = {
  list: ?CheckListTable_list,
  onChecklistChanged?: (updatedList: CheckListTable_list) => void,
  onDesignMode?: boolean,
};

const CheckListTable = (props: Props) => {
  const {list = []} = props;

  const checkListTableItems = Array.prototype.filter
    .call(list, item => CHECKLIST_ITEM_TYPES.hasOwnProperty(item.type))
    .sort(sortByIndex);

  const CheckListTableComponent = props.onDesignMode
    ? CheckListTableDefinition
    : CheckListTableFilling;

  const checkListTableProps = {
    ...props,
    list: checkListTableItems,
  };

  return <CheckListTableComponent {...checkListTableProps} />;
};

export default createFragmentContainer(CheckListTable, {
  list: graphql`
    fragment CheckListTable_list on CheckListItem @relay(plural: true) {
      id
      index
      type
      title
      checked
      ...CheckListItem_item
    }
  `,
});
