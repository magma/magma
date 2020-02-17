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

import AddIcon from '@fbcnms/ui/components/design-system/Icons/Actions/AddIcon';
import Button from '@fbcnms/ui/components/design-system/Button';
import DeleteIcon from '@fbcnms/ui/components/design-system/Icons/Actions/DeleteIcon';
import React, {useCallback, useMemo, useState} from 'react';
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
  deleteButton: {
    float: 'right',
    visibility: 'hidden',
  },
}));

type Props = {
  list: ?CheckListCategoryTable_list,
  onListChanged?: (updatedList: CheckListCategoryTable_list) => void,
};

const CheckListCategoryTable = (props: Props) => {
  const classes = useStyles();
  const {list: propsList = [], onListChanged} = props;
  const list = useMemo(() => {
    const itemsList = propsList ?? [];
    return [...itemsList].map((item, index) => ({
      index,
      key: item.id,
      value: item,
    }));
  }, [propsList]);
  const [nextNewItemTempId, setNextNewItemTempId] = useState(list.length + 1);
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
      id: `@tmp${newId}`,
      title: `this is my title - ${newId}`,
      description: 'very nice description',
      checkList: [],
    };
  }, [nextNewItemTempId]);

  const _updateCheckListCategory = useCallback(
    // eslint-disable-next-line flowtype/no-weak-types
    (updatedItem: any, index) =>
      (propsList &&
        _updateList([
          ...propsList.slice(0, index),
          updatedItem,
          ...propsList.slice(index + 1, propsList.length),
        ])) ??
      undefined,
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
  return (
    <>
      {list.length === 0 ? null : (
        <Table
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
                  onChange={e => {
                    _updateCheckListCategory(
                      Object.assign({}, row.value, {title: e.target.value}),
                      row.index,
                    );
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
              title: (
                <fbt desc="Items (number of questions in category) column header @ Checklist categories table">
                  Items
                </fbt>
              ),
              render: row => (
                <div className={classes.itemsCell}>
                  <Button disabled={true} skin="gray">
                    {`0/${row.value.checkList.length}`}
                  </Button>
                  <Button
                    variant="text"
                    className={classes.deleteButton}
                    onClick={() => _removeCheckListCategory(row.index)}>
                    <DeleteIcon color="gray" />
                  </Button>
                </div>
              ),
            },
          ]}
        />
      )}
      <Button variant="text" onClick={_addCheckListCategory}>
        <AddIcon color="primary" />
      </Button>
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
        id
      }
    }
  `,
});
