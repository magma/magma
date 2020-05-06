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
import {CheckListItemConfigs} from '../checkListCategory/CheckListItemConsts';
import {makeStyles} from '@material-ui/styles';
import {useFormContext} from '../../../common/FormContext';

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

type Props = {
  item: CheckListItem,
  children?: Node,
  onChange?: (newItem: CheckListItem) => void,
};

const CheckListItemDefinitionBase = ({children, item, onChange}: Props) => {
  const classes = useStyles();
  const form = useFormContext();
  const dispatch = useContext(ChecklistItemsDialogMutateDispatchContext);
  const config = CheckListItemConfigs[item.type];
  return (
    <div className={classes.root}>
      <div className={classes.editIndicator} />
      <Grid className={classes.mainDetails} container spacing={2}>
        <Grid item xs={6} l={5}>
          <TextInput
            type="string"
            disabled={form.alerts.editLock.detected}
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
            disabled={form.alerts.editLock.detected}
            options={Object.keys(CheckListItemConfigs).map(
              (itemType: CheckListItemType) => {
                const Icon = CheckListItemConfigs[itemType].icon;
                return {
                  key: `${itemType}`,
                  label: (
                    <div className={classes.label}>
                      <Icon className={classes.selectIcon} />
                      <Text variant="body2">
                        {CheckListItemConfigs[itemType].selectLabel}
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
        disabled={form.alerts.editLock.detected}
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
          disabled={form.alerts.editLock.detected}
          onClick={() => dispatch({type: 'REMOVE_ITEM', itemId: item.id})}>
          <DeleteIcon />
        </Button>
      </div>
    </div>
  );
};

export default CheckListItemDefinitionBase;
