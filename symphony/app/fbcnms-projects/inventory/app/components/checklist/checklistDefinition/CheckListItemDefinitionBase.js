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
import type {CheckListItemType} from '../../work_orders/__generated__/WorkOrderDetails_workOrder.graphql';
import type {Node} from 'react';

import Button from '@fbcnms/ui/components/design-system/Button';
import ChecklistItemsDialogMutateDispatchContext from '../checkListCategory/ChecklistItemsDialogMutateDispatchContext';
import DeleteIcon from '@fbcnms/ui/components/design-system/Icons/Actions/DeleteIcon';
import Grid from '@material-ui/core/Grid';
import React, {useContext} from 'react';
import Select from '@fbcnms/ui/components/design-system/Select/Select';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {CheckListItemIcons} from '../checkListCategory/CheckListItemConsts';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    padding: '16px',
    border: `1px solid ${symphony.palette.D100}`,
    backgroundColor: symphony.palette.D10,
    borderRadius: '4px',
    position: 'relative',
    width: '100%',
  },
  typeSelector: {
    width: '100%',
  },
  mainDetails: {
    marginBottom: '20px',
  },
  divider: {
    height: '1px',
    backgroundColor: symphony.palette.D100,
    marginTop: '20px',
    marginBottom: '16px',
  },
  actions: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'flex-end',
  },
  editIndicator: {
    backgroundColor: symphony.palette.primary,
    width: '3px',
    height: '100%',
    position: 'absolute',
    left: 0,
    top: 0,
    bottom: 0,
    borderRadius: '4px 0px 0px 4px',
  },
  label: {
    display: 'flex',
    alignItems: 'center',
  },
  selectIcon: {
    width: 24,
    height: 24,
    marginRight: 6,
  },
}));

const CHECKLIST_ITEM_CONFIGS: {
  [CheckListItemType]: {|
    selectLabel: Node,
    titlePlaceholder: string,
  |},
} = {
  simple: {
    selectLabel: <fbt desc="">Check when complete</fbt>,
    titlePlaceholder: `${fbt('What needs to be done?', '')}`,
  },
  string: {
    selectLabel: <fbt desc="">Free text</fbt>,
    titlePlaceholder: `${fbt('What needs to be written?', '')}`,
  },
  enum: {
    selectLabel: <fbt desc="">Multiple choice</fbt>,
    titlePlaceholder: `${fbt('What needs to be chosen?', '')}`,
  },
};

type Props = {
  item: CheckListItem,
  children?: Node,
  onChange?: (newItem: CheckListItem) => void,
};

const CheckListItemDefinitionBase = ({children, item, onChange}: Props) => {
  const classes = useStyles();
  const dispatch = useContext(ChecklistItemsDialogMutateDispatchContext);
  const config = CHECKLIST_ITEM_CONFIGS[item.type];
  return (
    <div className={classes.root}>
      <div className={classes.editIndicator} />
      <Grid className={classes.mainDetails} container spacing={2}>
        <Grid item xs={6} l={5}>
          <TextInput
            type="string"
            placeholder={config.titlePlaceholder}
            value={item.title}
            onChange={({target: {value}}) =>
              onChange &&
              onChange({
                ...item,
                title: value,
              })
            }
          />
        </Grid>
        <Grid item xs={1} l={3} />
        <Grid item xs={5} l={4}>
          <Select
            className={classes.typeSelector}
            options={Object.keys(CHECKLIST_ITEM_CONFIGS).map(
              (itemType: CheckListItemType) => {
                const Icon = CheckListItemIcons[itemType];
                return {
                  key: `${itemType}`,
                  label: (
                    <div className={classes.label}>
                      <Icon className={classes.selectIcon} />
                      <Text variant="body2">
                        {CHECKLIST_ITEM_CONFIGS[itemType].selectLabel}
                      </Text>
                    </div>
                  ),
                  value: itemType,
                };
              },
            )}
            selectedValue={item.type}
            onChange={type =>
              onChange &&
              onChange({
                ...item,
                type,
              })
            }
          />
        </Grid>
      </Grid>
      <TextInput
        type="string"
        placeholder={fbt('Additional instructions (optional)', '')}
        value={item.helpText ?? ''}
        onChange={({target: {value}}) =>
          onChange &&
          onChange({
            ...item,
            helpText: value,
          })
        }
      />
      {children}
      <div className={classes.divider} />
      <div className={classes.actions}>
        <Button
          variant="text"
          skin="gray"
          onClick={() => dispatch({type: 'REMOVE_ITEM', itemId: item.id})}>
          <DeleteIcon />
        </Button>
      </div>
    </div>
  );
};

export default CheckListItemDefinitionBase;
