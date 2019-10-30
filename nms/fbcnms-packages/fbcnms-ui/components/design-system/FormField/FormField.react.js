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
import Text from '../Text.react';
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
  labelContainer: {
    marginBottom: '6px',
  },
  bottomText: {
    marginTop: '4px',
  },
}));

type Props = {
  className?: string,
  label: string,
  helpText?: string,
  children: React.Node,
  disabled: boolean,
  hasError: boolean,
  errorText?: string,
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
        <Text className={classes.labelContainer} variant="body2">
          {label}
        </Text>
        {children}
        {(helpText || (hasError && errorText)) && (
          <Text className={classes.bottomText} variant="caption">
            {nullthrows(hasError ? errorText : helpText)}
          </Text>
        )}
      </div>
    </FormFieldContext.Provider>
  );
};

FormField.defaultProps = {
  disabled: false,
  hasError: false,
};

export default FormField;
