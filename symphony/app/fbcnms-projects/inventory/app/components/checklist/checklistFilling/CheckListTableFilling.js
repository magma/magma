/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CheckListItem} from '../checkListCategory/ChecklistItemsDialogMutateState';
import type {TableRowDataType} from '@fbcnms/ui/components/design-system/Table/Table';

import CheckListItemFilling from './CheckListItemFilling';
import ChecklistItemsDialogMutateDispatchContext from '../checkListCategory/ChecklistItemsDialogMutateDispatchContext';
import React, {useContext} from 'react';
import Table from '@fbcnms/ui/components/design-system/Table/Table';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {CheckListItemIcons} from '../checkListCategory/CheckListItemConsts';
import {isChecklistItemDone} from '../ChecklistUtils';
import {makeStyles} from '@material-ui/styles';

type Props = {
  items: Array<CheckListItem>,
};

const useStyles = makeStyles(() => ({
  container: {
    maxWidth: '1366px',
    overflowX: 'auto',
  },
  root: {
    marginBottom: '12px',
    maxWidth: '100%',
  },
  titleContainer: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
  },
  itemTitle: {
    marginLeft: '12px',
  },
  icon: {
    height: '24px',
    width: '24px',
    fill: symphony.palette.D300,
  },
  doneIcon: {
    fill: symphony.palette.G600,
  },
}));

const CheckListTableFilling = ({items}: Props) => {
  const dispatch = useContext(ChecklistItemsDialogMutateDispatchContext);
  const classes = useStyles();

  const tableData: Array<TableRowDataType<{|item: CheckListItem|}>> = items.map(
    item => ({item, key: item.id}),
  );
  return (
    <div className={classes.container}>
      <Table
        variant="embedded"
        dataRowsSeparator="border"
        data={tableData}
        columns={[
          {
            key: 'item',
            title: <fbt desc="">Item</fbt>,
            render: row => {
              const Icon = CheckListItemIcons[row.item.type];
              return (
                <div className={classes.titleContainer}>
                  <Icon
                    className={classNames(classes.icon, {
                      [classes.doneIcon]: isChecklistItemDone(row.item),
                    })}
                  />
                  <Text
                    weight="medium"
                    variant="body2"
                    className={classes.itemTitle}>
                    {row.item.title.trim() !== '' ? (
                      row.item.title
                    ) : (
                      <fbt desc="">Item</fbt>
                    )}
                  </Text>
                </div>
              );
            },
          },
          {
            key: 'response',
            title: <fbt desc="">Response</fbt>,
            render: row => (
              <CheckListItemFilling
                item={row.item}
                onChange={updatedItem =>
                  dispatch({
                    type: 'EDIT_ITEM',
                    value: updatedItem,
                  })
                }
              />
            ),
          },
        ]}
      />
    </div>
  );
};

export default CheckListTableFilling;
