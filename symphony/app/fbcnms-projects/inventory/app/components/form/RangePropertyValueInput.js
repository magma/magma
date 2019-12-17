/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FocusEvent} from '@fbcnms/ui/components/design-system/Input/TextInput';

import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import FormValidationContext from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import InputAffix from '@fbcnms/ui/components/design-system/Input/InputAffix';
import React from 'react';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

type Range = {
  rangeFrom: ?number,
  rangeTo: ?number,
};

type Props = {
  label?: ?string,
  value: Range,
  className: string,
  required: boolean,
  disabled: boolean,
  onBlur: (e: FocusEvent<HTMLInputElement>) => ?void,
  onRangeToChange: (event: SyntheticInputEvent<>) => void,
  onRangeFromChange: (event: SyntheticInputEvent<>) => void,
  margin: 'none' | 'dense' | 'normal',
  autoFocus?: boolean,
};

const useStyles = makeStyles(theme => ({
  container: {
    display: 'flex',
    width: '280px',
  },
  input: {
    marginLeft: '0px',
    marginRight: theme.spacing(),
    width: '100%',
  },
  lngField: {
    marginLeft: '16px',
  },
  formField: {
    flexGrow: 1,
  },
}));

const ENTER_KEY_CODE = 13;

const RangePropertyValueInput = (props: Props) => {
  const {
    className,
    disabled,
    margin,
    required,
    value,
    autoFocus,
    label,
  } = props;
  const classes = useStyles();

  const {rangeFrom, rangeTo} = value;
  const fieldIdPrefix = `range-${label || 'field'}-`;
  return (
    <FormValidationContext.Consumer>
      {validationContext => {
        const errorFrom = validationContext.error.check({
          fieldId: `${fieldIdPrefix}From`,
          fieldDisplayName: 'From',
          value: parseFloat(rangeFrom),
          required: required,
        });
        const errorTo = validationContext.error.check({
          fieldId: `${fieldIdPrefix}To`,
          fieldDisplayName: 'To',
          value: parseFloat(rangeTo),
          required: required,
        });

        return (
          <FormField label={label || ''} required={required}>
            <div className={classNames(classes.container, className)}>
              <FormField
                className={classes.formField}
                required={required}
                errorText={errorFrom}
                hasError={!!errorFrom}>
                <TextInput
                  autoFocus={autoFocus}
                  required={required}
                  disabled={disabled}
                  prefix={<InputAffix>From</InputAffix>}
                  id="from-value"
                  variant="outlined"
                  className={classes.input}
                  margin={margin}
                  onKeyDown={e => {
                    if (e.keyCode === ENTER_KEY_CODE) {
                      props.onBlur(e);
                    }
                  }}
                  value={parseFloat(rangeFrom)}
                  type="number"
                  onChange={props.onRangeFromChange}
                />
              </FormField>
              <FormField
                required={required}
                errorText={errorTo}
                hasError={!!errorTo}
                className={classNames(classes.lngField, classes.formField)}>
                <TextInput
                  required={required}
                  disabled={disabled}
                  prefix={<InputAffix>To</InputAffix>}
                  id="to-value"
                  variant="outlined"
                  className={classes.input}
                  margin={margin}
                  onKeyDown={e => {
                    if (e.keyCode === ENTER_KEY_CODE) {
                      props.onBlur(e);
                    }
                  }}
                  type="number"
                  value={parseFloat(rangeTo)}
                  onChange={props.onRangeToChange}
                />
              </FormField>
            </div>
          </FormField>
        );
      }}
    </FormValidationContext.Consumer>
  );
};

export default RangePropertyValueInput;
