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
import type {CheckListItemEnumSelectionMode} from '../../work_orders/__generated__/WorkOrderDetails_workOrder.graphql';

import CheckListItemDefinitionBase from './CheckListItemDefinitionBase';
import React from 'react';
import Select from '@fbcnms/ui/components/design-system/Select/Select';
import Tokenizer from '@fbcnms/ui/components/Tokenizer';
import emptyFunction from '@fbcnms/util/emptyFunction';
import fbt from 'fbt';
import {enumStringToArray} from '../ChecklistUtils';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    flexDirection: 'row',
    marginTop: '20px',
  },
  tokenizer: {
    backgroundColor: 'white',
    flexGrow: 2,
  },
  select: {
    flexGrow: 1,
    marginRight: '16px',
  },
}));

type Props = {
  item: CheckListItem,
  onChange?: (updatedItem: CheckListItem) => void,
};

const MultipleChoiceCheckListItemDefinition = ({item, onChange}: Props) => {
  const classes = useStyles();
  return (
    <CheckListItemDefinitionBase item={item} onChange={onChange}>
      <div className={classes.root}>
        <Select
          className={classes.select}
          options={[
            {
              key: 'single',
              label: <fbt desc="">Select one option</fbt>,
              value: 'single',
            },
            {
              key: 'multiple',
              label: <fbt desc="">Select multiple options</fbt>,
              value: 'multiple',
            },
            ,
          ]}
          selectedValue={item.enumSelectionMode}
          onChange={(value: CheckListItemEnumSelectionMode) => {
            if (value === item.enumSelectionMode) {
              return;
            }
            const modifiedItem: CheckListItem = {
              ...item,
              enumSelectionMode: value,
              selectedEnumValues: '',
            };
            onChange && onChange(modifiedItem);
          }}
        />
        <Tokenizer
          placeholder={`${fbt('Press Enter after each value', '')}`}
          className={classes.tokenizer}
          searchSource="UserInput"
          tokens={enumStringToArray(item.enumValues).map(value => ({
            label: value,
            id: value,
          }))}
          onChange={entries =>
            onChange &&
            onChange({...item, enumValues: entries.map(e => e.label).join(',')})
          }
          onEntriesRequested={emptyFunction}
        />
      </div>
    </CheckListItemDefinitionBase>
  );
};

export default MultipleChoiceCheckListItemDefinition;
