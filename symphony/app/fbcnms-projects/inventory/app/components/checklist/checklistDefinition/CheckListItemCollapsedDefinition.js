/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {CheckListItem} from '../checkListCategory/ChecklistItemsDialogMutateState';

import ChecklistItemsDialogMutateDispatchContext from '../checkListCategory/ChecklistItemsDialogMutateDispatchContext';
import React, {useContext} from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {CheckListItemConfigs} from '../checkListCategory/CheckListItemConsts';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    borderRadius: '4px',
    padding: '12px',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    cursor: 'pointer',
    backgroundColor: symphony.palette.white,
    border: `1px solid ${symphony.palette.D100}`,
    '&:hover': {
      backgroundColor: symphony.palette.background,
    },
    width: '100%',
  },
  title: {
    display: 'flex',
    flexDirection: 'column',
  },
  icon: {
    marginRight: '12px',
    fill: symphony.palette.D300,
  },
}));

type Props = {
  className?: string,
  item: CheckListItem,
};

const CheckListItemCollapsedDefinition = ({item, className}: Props) => {
  const classes = useStyles();
  const dispatch = useContext(ChecklistItemsDialogMutateDispatchContext);
  const itemConfig = CheckListItemConfigs[item.type];
  if (itemConfig == null) {
    return null;
  }
  const Icon = itemConfig.icon;
  return (
    <div
      className={classNames(classes.root, className)}
      onClick={() =>
        dispatch({type: 'SET_EDITED_DEFINITION_ID', itemId: item.id})
      }>
      <Icon className={classes.icon} />
      <div className={classes.title}>
        <Text variant="body2">
          {item.title.trim() !== '' ? item.title : <fbt desc="">Item</fbt>}
        </Text>
        {item.helpText != null && (
          <Text variant="body2" color="gray">
            {item.helpText}
          </Text>
        )}
      </div>
    </div>
  );
};

export default CheckListItemCollapsedDefinition;
