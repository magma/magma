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
import FormValidationContext from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import React, {useContext} from 'react';
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
  const classes = useStyles();
  const validationContext = useContext(FormValidationContext);
  return (
    <div title={props.isDisabled && props.disabledMessage}>
      <Button
        className={classNames(
          classes.cancelButton,
          props.classes?.cancelButton,
        )}
        disabled={props.isDisabled}
        onClick={props.onCancel}
        skin="regular">
        {props.captions?.cancelButton || 'Cancel'}
      </Button>
      <FormAction>
        <Button
          className={props.classes?.saveButton}
          onClick={props.onSave}
          tooltip={validationContext.error.message}
          disabled={props.isDisabled || validationContext.error.detected}>
          {props.captions?.saveButton || 'Save'}
        </Button>
      </FormAction>
    </div>
  );
};

export default FormSaveCancelPanel;
