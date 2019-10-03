/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {WithStyles} from '@material-ui/core';

import * as React from 'react';
import classNames from 'classnames';
import {withStyles} from '@material-ui/core/styles';

const styles = ({symphony}) => ({
  root: {
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
  disabled: {},
  hasError: {},
  input: {
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
});

type Props = {
  /** Input type. See: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input#Form_%3Cinput%3E_types */
  type: string,
  value: string | number,
  className?: string,
  placeholder?: string,
  autoFocus?: boolean,
  disabled?: boolean,
  hasError?: boolean,
  onChange?: (e: SyntheticInputEvent<HTMLInputElement>) => void,
  onFocus?: () => void,
  onBlur?: () => void,
} & WithStyles<typeof styles>;

type State = {
  hasFocus: boolean,
};

class TextInput extends React.Component<Props, State> {
  static defaultProps = {
    autoFocus: false,
    disabled: false,
    hasError: false,
  };

  constructor(props: Props) {
    super(props);
    this.state = {
      hasFocus: props.autoFocus === true,
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

  render() {
    const {classes, className, hasError, disabled, ...rest} = this.props;
    const {hasFocus} = this.state;
    return (
      <div
        className={classNames(
          classes.root,
          {
            [classes.hasFocus]: hasFocus,
            [classes.disabled]: disabled,
            [classes.hasError]: hasError,
          },
          className,
        )}>
        <input
          className={classes.input}
          disabled={disabled}
          onFocus={this._onInputFocused}
          onBlur={this._onInputBlurred}
          onChange={this._onChange}
          {...rest}
        />
      </div>
    );
  }
}

export default withStyles(styles)(TextInput);
