/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

// import type {CheckListCategoryInput} from '../../../mutations/__generated__/AddWorkOrderMutation.graphql';
import type {CheckListCategoryTable_list} from './__generated__/CheckListCategoryTable_list.graphql';

import Button from '@fbcnms/ui/components/design-system/Button';
import CheckListCategoryContext from './CheckListCategoryContext';
import CheckListCategoryItemsDialog from './CheckListCategoryItemsDialog';
import DeleteIcon from '@fbcnms/ui/components/design-system/Icons/Actions/DeleteIcon';
import React, {
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from 'react';
import Table from '@fbcnms/ui/components/design-system/Table/Table';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import fbt from 'fbt';
import {createFragmentContainer, graphql} from 'react-relay';
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

type Props = {
  list: CheckListCategoryTable_list,
  onListChanged?: (updatedList: CheckListCategoryTable_list) => void,
};

const CheckListCategoryTable = (props: Props) => {
  const classes = useStyles();
  const {list: propsList, onListChanged} = props;
  const list = useMemo(() => {
    return propsList.map((item, index) => ({
      index,
      key: item.id || `@key${index}`,
      value: item,
      responsesCount: item.checkList.reduce(
        (responsesCount, clItem) =>
          clItem.checked ? responsesCount + 1 : responsesCount,
        0,
      ),
    }));
  }, [propsList]);
  const [
    browsedCheckListCategory,
    setBrowsedCheckListCategory,
  ] = useState<?number>(null);
  const [nextNewItemTempId, setNextNewItemTempId] = useState(list.length + 1);
  const defaultCategoryName = `${fbt(
    'New Category',
    'Default name for checklist category',
  )}`;
  const _updateList = useCallback(
    (updatedList: CheckListCategoryTable_list) => {
      if (!onListChanged) {
        return;
      }
      onListChanged(updatedList);
    },
    [onListChanged],
  );
  // eslint-disable-next-line flowtype/no-weak-types
  const _createNewItem: () => any = useCallback(() => {
    const newId = nextNewItemTempId;
    setNextNewItemTempId(newId + 1);
    return {
      title: defaultCategoryName,
      description: '',
      checkList: [],
    };
  }, [defaultCategoryName, nextNewItemTempId]);
  const _updateCheckListCategory = useCallback(
    // eslint-disable-next-line flowtype/no-weak-types
    (updatedItem: any, index) =>
      _updateList([
        ...propsList.slice(0, index),
        updatedItem,
        ...propsList.slice(index + 1, propsList.length),
      ]),
    [_updateList, propsList],
  );

  const _removeCheckListCategory = useCallback(
    index =>
      (propsList &&
        _updateList([
          ...propsList.slice(0, index),
          ...propsList.slice(index + 1, propsList.length),
        ])) ??
      undefined,
    [_updateList, propsList],
  );

  const _addCheckListCategory = useCallback(
    () =>
      (propsList && _updateList([...propsList, _createNewItem()])) ?? undefined,
    [_createNewItem, _updateList, propsList],
  );
  const context = useContext(CheckListCategoryContext);
  useEffect(() => {
    context.override.addNewCategory(_addCheckListCategory);
  }, [_addCheckListCategory, context.override]);
  return list.length === 0 ? null : (
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
                value={row.value.title}
                autoFocus={true}
                placeholder={`${fbt(
                  'Name of the category',
                  'hint text for checklist category name field',
                )}`}
                onChange={e => {
                  _updateCheckListCategory(
                    Object.assign({}, row.value, {
                      title: e.target.value,
                    }),
                    row.index,
                  );
                }}
                onBlur={() => {
                  if (!row.value.title) {
                    _updateCheckListCategory(
                      Object.assign({}, row.value, {
                        title: defaultCategoryName,
                      }),
                      row.index,
                    );
                  }
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
                value={row.value.description || ''}
                placeholder={`${fbt(
                  'Short description of category (optional)',
                  'hint text for optional checklist category description field',
                )}`}
                onChange={e => {
                  _updateCheckListCategory(
                    Object.assign({}, row.value, {
                      description: e.target.value,
                    }),
                    row.index,
                  );
                }}
              />
            ),
          },
          {
            key: '2',
            title: !!list.find(row => row.value.checkList.length > 0) ? (
              <fbt desc="Completed Items (number of filled questions in category) column header @ Checklist categories table">
                Completed Items
              </fbt>
            ) : (
              <fbt desc="Items (number of questions in category) column header @ Checklist categories table">
                Items
              </fbt>
            ),
            render: row => (
              <Button
                skin="gray"
                className={classes.addItemsButton}
                onClick={() => setBrowsedCheckListCategory(row.index)}>
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
            render: row => (
              <Button
                variant="text"
                className={classes.deleteButton}
                onClick={() => _removeCheckListCategory(row.index)}>
                <DeleteIcon color="gray" />
              </Button>
            ),
          },
        ]}
      />

      {browsedCheckListCategory != null && (
        <CheckListCategoryItemsDialog
          items={list[browsedCheckListCategory]?.value.checkList}
          categoryTitle={list[browsedCheckListCategory]?.value.title}
          onChecklistChanged={updatedList =>
            _updateCheckListCategory(
              Object.assign({}, list[browsedCheckListCategory].value, {
                checkList: updatedList,
              }),
              browsedCheckListCategory,
            )
          }
          onClose={() => setBrowsedCheckListCategory(null)}
        />
      )}
    </>
  );
};

export default createFragmentContainer(CheckListCategoryTable, {
  list: graphql`
    fragment CheckListCategoryTable_list on CheckListCategory
      @relay(plural: true) {
      id
      title
      description
      checkList {
        ...CheckListCategoryItemsDialog_items
        id
        title
        type
        index
        helpText
        enumValues
        stringValue
        checked
      }
    }
  `,
});
