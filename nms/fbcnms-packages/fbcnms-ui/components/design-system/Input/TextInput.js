/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FormFieldContextValue} from '../FormField/FormFieldContext';
import type {WithStyles} from '@material-ui/core';

import * as React from 'react';
import FormFieldContext from '../FormField/FormFieldContext';
import InputContext from './InputContext';
import Text from '../Text';
import classNames from 'classnames';
import {withStyles} from '@material-ui/core/styles';

export const KEYBOARD_KEYS = {
  CODES: {
    ENTER: 13,
  },
  MODIFIERS: {
    SHIFT: 'shift',
  },
};

const styles = ({symphony}) => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
  },
  inputContainer: {
    padding: '0px 8px',
    border: `1px solid ${symphony.palette.D100}`,
    borderRadius: '4px',
    display: 'flex',
    height: '32px',
    boxSizing: 'border-box',
    backgroundColor: symphony.palette.white,
    '&$hasFocus': {
      borderColor: symphony.palette.D500,
    },
    '&:hover:not($disabled):not($hasError)': {
      borderColor: symphony.palette.D500,
    },
    '&$disabled': {
      backgroundColor: symphony.palette.D50,
    },
    '&$hasError': {
      borderColor: symphony.palette.R600,
    },
  },
  hasFocus: {},
  disabled: {
    '& $input': {
      '&::placeholder': {
        color: symphony.palette.disabled,
      },
    },
  },
  hasError: {},
  input: {
    color: symphony.palette.D900,
    margin: 0,
    border: 0,
    outline: 0,
    background: 'transparent',
    flexGrow: 1,
    ...symphony.typography.body2,
    '&::placeholder': {
      color: symphony.palette.D400,
    },
  },
  prefix: {
    display: 'flex',
    alignItems: 'center',
    marginRight: '8px',
    marginLeft: '4px',
  },
  hint: {
    paddingTop: '4px',
  },
  hintText: {
    color: symphony.palette.D200,
  },
  suffix: {
    display: 'flex',
    alignItems: 'center',
    marginRight: '4px',
    marginLeft: '8px',
  },
});

type Props = {
  /** Input type. See: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input#Form_%3Cinput%3E_types */
  type?: string,
  value?: string | number,
  className?: string,
  placeholder?: string,
  autoFocus?: boolean,
  disabled?: boolean,
  hasError?: boolean,
  prefix?: React.Node,
  hint?: string,
  suffix?: React.Node,
  onChange?: (e: SyntheticInputEvent<HTMLInputElement>) => void,
  onFocus?: () => void,
  onBlur?: () => void,
  onEnterPressed?: (e: KeyboardEvent) => void,
} & WithStyles<typeof styles>;

type State = {
  hasFocus: boolean,
};

class TextInput extends React.Component<Props, State> {
  static contextType = FormFieldContext;

  static defaultProps = {
    autoFocus: false,
    disabled: false,
    hasError: false,
    type: 'string',
  };

  constructor(props: Props, context: FormFieldContextValue) {
    super(props);

    const disabled = props.disabled ? props.disabled : context.disabled;
    this.state = {
      hasFocus: disabled ? false : props.autoFocus === true,
    };
  }

  _onInputFocused = () => {
    this.setState({hasFocus: true});
    const {onFocus} = this.props;
    onFocus && onFocus();
  };

  _onInputBlurred = () => {
    this.setState({hasFocus: false});
    const {onBlur} = this.props;
    onBlur && onBlur();
  };

  _onChange = (e: SyntheticInputEvent<HTMLInputElement>) => {
    const {onChange} = this.props;
    onChange && onChange(e);
  };

  _onKeyDown = (e: KeyboardEvent) => {
    if (e.keyCode !== KEYBOARD_KEYS.CODES.ENTER) {
      return;
    }

    const {onEnterPressed} = this.props;
    onEnterPressed && onEnterPressed(e);
  };

  render() {
    const {
      classes,
      className,
      hasError: hasErrorProp,
      disabled: disabledProp,
      prefix,
      suffix,
      value,
      hint,
      ...rest
    } = this.props;
    const {hasFocus} = this.state;
    const disabled = disabledProp ? disabledProp : this.context.disabled;
    const hasError = hasErrorProp ? hasErrorProp : this.context.hasError;
    return (
      <div className={classNames(classes.root, className)}>
        <div
          className={classNames(classes.inputContainer, {
            [classes.hasFocus]: hasFocus,
            [classes.disabled]: disabled,
            [classes.hasError]: hasError,
          })}>
          <InputContext.Provider value={{disabled, value: value ?? ''}}>
            {prefix && <div className={classes.prefix}>{prefix}</div>}
            <input
              className={classes.input}
              disabled={disabled}
              onFocus={this._onInputFocused}
              onBlur={this._onInputBlurred}
              onChange={this._onChange}
              onKeyDown={this._onKeyDown}
              value={value}
              {...rest}
            />
            {suffix && <div className={classes.suffix}>{suffix}</div>}
          </InputContext.Provider>
        </div>
        {hint && (
          <div className={classes.hint}>
            <Text variant="caption" className={classes.hintText}>
              {hint}
            </Text>
          </div>
        )}
      </div>
    );
  }
}

export default withStyles(styles)(TextInput);
