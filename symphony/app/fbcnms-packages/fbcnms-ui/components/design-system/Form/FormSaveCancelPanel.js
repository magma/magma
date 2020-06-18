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
import FormAction from './FormAction';
import React from 'react';
import Strings from '../../../../fbcnms-strings/Strings';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  cancelButton: {
    marginRight: '8px',
  },
}));

type Props = {
  isDisabled?: boolean,
  disabledMessage?: string,
  onSave: () => void,
  onCancel: () => void,
  classes?: {
    cancelButton?: string,
    saveButton?: string,
  },
  captions?: {
    cancelButton?: string,
    saveButton?: string,
  },
};

const FormSaveCancelPanel = (props: Props) => {
  const {
    captions,
    classes: propsClasses,
    isDisabled,
    disabledMessage,
    onCancel,
    onSave,
  } = props;
  const classes = useStyles();
  return (
    <div title={isDisabled && disabledMessage}>
      <FormAction
        ignorePermissions={true}
        disabled={isDisabled}
        tooltip={isDisabled ? disabledMessage : undefined}>
        <Button
          className={classNames(
            classes.cancelButton,
            propsClasses?.cancelButton,
          )}
          onClick={onCancel}
          skin="regular">
          {captions?.cancelButton || Strings.common.cancelButton}
        </Button>
      </FormAction>
      <FormAction
        disableOnFromError={true}
        disabled={isDisabled}
        tooltip={isDisabled ? disabledMessage : undefined}>
        <Button className={propsClasses?.saveButton} onClick={onSave}>
          {captions?.saveButton || Strings.common.saveButton}
        </Button>
      </FormAction>
    </div>
  );
};

export default FormSaveCancelPanel;
