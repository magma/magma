/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ChecklistCategoriesStateType} from '../ChecklistCategoriesMutateState';

import Button from '@fbcnms/ui/components/design-system/Button';
import CheckListCategoryItemsDialog from './CheckListCategoryItemsDialog';
import ChecklistCategoriesMutateDispatchContext from '../ChecklistCategoriesMutateDispatchContext';
import DeleteIcon from '@fbcnms/ui/components/design-system/Icons/Actions/DeleteIcon';
import FormValidationContext from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import React, {useContext, useMemo, useState} from 'react';
import Table from '@fbcnms/ui/components/design-system/Table/Table';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import fbt from 'fbt';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  categoryRow: {
    '&:hover $deleteButton': {
      visibility: 'visible',
    },
  },
  itemsCell: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  addItemsButton: {
    width: '100%',
    maxWidth: '160px',
  },
  deleteButton: {
    float: 'right',
    visibility: 'hidden',
  },
}));

type CheckListCategoryTableProps = $ReadOnly<{
  categories: ChecklistCategoriesStateType,
}>;

const CheckListCategoryTable = ({categories}: CheckListCategoryTableProps) => {
  const classes = useStyles();
  const formValidationContext = useContext(FormValidationContext);
  const dispatch = useContext(ChecklistCategoriesMutateDispatchContext);
  const [
    browsedCheckListCategoryId,
    setBrowsedCheckListCategoryId,
  ] = useState<?string>(null);

  const browsedCategory = useMemo(
    () => categories.find(c => c.id === browsedCheckListCategoryId) ?? null,
    [categories, browsedCheckListCategoryId],
  );
  const list = useMemo(() => {
    return categories.map((item, index) => ({
      index,
      key: item.id || `@key${index}`,
      value: item,
      responsesCount: item.checkList.reduce(
        (responsesCount, clItem) =>
          clItem.checked ? responsesCount + 1 : responsesCount,
        0,
      ),
    }));
  }, [categories]);

  return categories.length === 0 ? null : (
    <>
      <Table
        variant="embedded"
        dataRowsSeparator="border"
        dataRowClassName={classes.categoryRow}
        data={list}
        columns={[
          {
            key: '0',
            title: (
              <fbt desc="Category Name column header @ Checklist categories table">
                Category Name
              </fbt>
            ),
            render: row => (
              <TextInput
                id="title"
                variant="outlined"
                disabled={formValidationContext.editLock.detected}
                value={row.value.title}
                autoFocus={true}
                placeholder={`${fbt(
                  'Name of the category',
                  'hint text for checklist category name field',
                )}`}
                onChange={e => {
                  dispatch({
                    type: 'UPDATE_CATEGORY_TITLE',
                    categoryId: row.value.id,
                    value: e.target.value,
                  });
                }}
              />
            ),
          },
          {
            key: '1',
            title: (
              <fbt desc="Category Description column header @ Checklist categories table">
                Category Description
              </fbt>
            ),
            render: row => (
              <TextInput
                id="description"
                variant="outlined"
                disabled={formValidationContext.editLock.detected}
                value={row.value.description || ''}
                placeholder={`${fbt(
                  'Short description of category (optional)',
                  'hint text for optional checklist category description field',
                )}`}
                onChange={e => {
                  dispatch({
                    type: 'UPDATE_CATEGORY_DESCRIPTION',
                    categoryId: row.value.id,
                    value: e.target.value,
                  });
                }}
              />
            ),
          },
          {
            key: '2',
            title: (
              <fbt desc="Items (number of questions in category) column header @ Checklist categories table">
                Items
              </fbt>
            ),
            render: row => (
              <Button
                skin="gray"
                disabled={formValidationContext.editLock.detected}
                className={classes.addItemsButton}
                onClick={() => setBrowsedCheckListCategoryId(row.value.id)}>
                {row.value.checkList.length > 0 ? (
                  `${row.responsesCount}/${row.value.checkList.length}`
                ) : (
                  <fbt desc="Add checklist items button caption">Add Items</fbt>
                )}
              </Button>
            ),
          },
          {
            key: '3',
            title: '',
            render: row =>
              formValidationContext.editLock.detected ? null : (
                <Button
                  variant="text"
                  className={classes.deleteButton}
                  onClick={() =>
                    dispatch({
                      type: 'REMOVE_CATEGORY',
                      categoryId: row.value.id,
                    })
                  }>
                  <DeleteIcon color="gray" />
                </Button>
              ),
          },
        ]}
      />
      {browsedCheckListCategoryId != null && (
        <CheckListCategoryItemsDialog
          initialItems={browsedCategory?.checkList ?? []}
          categoryTitle={browsedCategory?.title ?? ''}
          onCancel={() => setBrowsedCheckListCategoryId(null)}
          onSave={items => {
            setBrowsedCheckListCategoryId(null);
            dispatch({
              type: 'UPDATE_CATEGORY_CHECKLIST',
              categoryId: browsedCheckListCategoryId,
              value: items,
            });
          }}
        />
      )}
    </>
  );
};

export default CheckListCategoryTable;
