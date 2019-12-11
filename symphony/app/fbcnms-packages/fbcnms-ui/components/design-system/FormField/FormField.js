/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import FormFieldContext from './FormFieldContext';
import Text from '../Text';
import classNames from 'classnames';
import nullthrows from 'nullthrows';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(({symphony}) => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
  },
  disabled: {
    '& $bottomText': {
      color: symphony.palette.disabled,
    },
  },
  hasError: {
    '& $bottomText': {
      color: symphony.palette.R600,
    },
  },
  labelPlaceholder: {
    opacity: 0,
  },
  labelContainer: {
    marginBottom: '6px',
  },
  bottomText: {
    marginTop: '4px',
    lineHeight: '16px',
  },
  spacer: {
    marginTop: '4px',
    height: '16px',
  },
}));

type Props = {
  className?: string,
  label?: string,
  helpText?: string,
  children: React.Node,
  disabled: boolean,
  hasError: boolean,
  required: boolean,
  errorText?: ?string,
  hasSpacer?: boolean,
  useLabelPlaceholder?: boolean,
};

const FormField = (props: Props) => {
  const {
    children,
    label,
    helpText,
    disabled,
    className,
    hasError,
    errorText,
    hasSpacer,
    required,
    useLabelPlaceholder,
  } = props;
  const classes = useStyles();
  return (
    <FormFieldContext.Provider value={{disabled, hasError}}>
      <div
        className={classNames(
          classes.root,
          {[classes.disabled]: disabled},
          {[classes.hasError]: hasError},
          className,
        )}>
        {(label || useLabelPlaceholder) && (
          <Text
            variant="body2"
            className={classNames(classes.labelContainer, {
              [classes.labelPlaceholder]: !label && useLabelPlaceholder,
            })}>
            {label || (useLabelPlaceholder && '.') || ''}
            {required && ' *'}
          </Text>
        )}
        {children}
        {(helpText || (hasError && errorText)) && (
          <Text className={classes.bottomText} variant="caption">
            {nullthrows(hasError ? errorText : helpText)}
          </Text>
        )}
        {!helpText && !hasError && hasSpacer && (
          <div className={classes.spacer} />
        )}
      </div>
    </FormFieldContext.Provider>
  );
};

FormField.defaultProps = {
  disabled: false,
  hasError: false,
  required: false,
};

export default FormField;
