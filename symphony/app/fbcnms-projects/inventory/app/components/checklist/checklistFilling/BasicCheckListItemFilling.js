/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Button from '@fbcnms/ui/components/design-system/Button';
import FormValidationContext from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import React, {useContext} from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import fbt from 'fbt';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';
import type {BasicCheckListItemFilling_item} from './__generated__/BasicCheckListItemFilling_item.graphql';

type Props = {
  item: BasicCheckListItemFilling_item,
  onChange?: (updatedChecklistItem: BasicCheckListItemFilling_item) => void,
};

const useStyles = makeStyles(() => ({
  container: {
    display: 'flex',
    flexDirection: 'row',
  },
  expandindPart: {
    flexGrow: 1,
    flexBasis: 0,
  },
}));

const BasicCheckListItemFilling = (props: Props) => {
  const {item, onChange} = props;
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
      <Text className={classes.expandindPart} variant="body2" weight="regular">
        {item.title}
      </Text>
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

export default createFragmentContainer(BasicCheckListItemFilling, {
  item: graphql`
    fragment BasicCheckListItemFilling_item on CheckListItem {
      title
      checked
      ...CheckListItem_item
    }
  `,
});
