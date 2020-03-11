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

import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import React, {useCallback} from 'react';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import fbt from 'fbt';
import {makeStyles} from '@material-ui/styles';

type Props = {
  item: CheckListItem,
  onChange?: (updatedItem: CheckListItem) => void,
};

const useStyles = makeStyles(() => ({
  container: {
    display: 'flex',
    flexDirection: 'row',
  },
  expandingPart: {
    flexGrow: 1,
    flexBasis: 0,
  },
}));

const BasicCheckListItemDefinition = ({item, onChange}: Props) => {
  const classes = useStyles();

  const _updateOnChange = useCallback(
    newTitle => {
      const newItem = {
        ...item,
        title: newTitle,
      };
      onChange && onChange(newItem);
    },
    [item, onChange],
  );

  return (
    <div className={classes.container}>
      <FormField className={classes.expandingPart}>
        <TextInput
          type="string"
          placeholder={fbt(
            'What needs to be done?',
            'Placeholder for checkbox field title (user needs to type the title of the checkbox in this field).',
          )}
          value={item.title || ''}
          onChange={event => _updateOnChange(event.target.value)}
        />
      </FormField>
    </div>
  );
};

export default BasicCheckListItemDefinition;
