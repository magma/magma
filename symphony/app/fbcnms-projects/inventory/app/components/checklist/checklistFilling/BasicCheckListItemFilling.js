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

import Button from '@fbcnms/ui/components/design-system/Button';
import FormValidationContext from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import React, {useContext} from 'react';
import fbt from 'fbt';
import {makeStyles} from '@material-ui/styles';

type Props = {
  item: CheckListItem,
  onChange?: (updatedChecklistItem: CheckListItem) => void,
};

const useStyles = makeStyles(() => ({
  container: {
    display: 'flex',
    flexDirection: 'row',
  },
}));

const BasicCheckListItemFilling = ({item, onChange}: Props) => {
  const classes = useStyles();

  const _updateOnChange = () => {
    if (!onChange) {
      return;
    }
    const modifiedItem = {
      ...item,
      checked: !item.checked,
    };
    onChange(modifiedItem);
  };

  const validationContext = useContext(FormValidationContext);

  return (
    <div className={classes.container}>
      {!validationContext.editLock.detected && (
        <Button onClick={_updateOnChange} variant="text">
          {item.checked
            ? fbt(
                'Mark as Undone',
                'Caption of the simple checkbox item Uncheck button',
              )
            : fbt(
                'Mark as done',
                'Caption of the simple checkbox item Check button',
              )}
        </Button>
      )}
    </div>
  );
};

export default BasicCheckListItemFilling;
