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
import RadioButtonCheckedIcon from '@material-ui/icons/RadioButtonChecked';
import RadioButtonUncheckedIcon from '@material-ui/icons/RadioButtonUnchecked';
import SymphonyTheme from '../../../theme/symphony';
import Text from '../Text';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';
import {useFormElementContext} from '../Form/FormElementContext';
import {useMemo} from 'react';

const useStyles = makeStyles(() => ({
  option: {
    display: 'flex',
    marginBottom: '8px',
    '& $checkIcon': {
      color: SymphonyTheme.palette.D200,
    },
    '&$disabled': {
      '& $checkDetails $label, $details': {
        color: SymphonyTheme.palette.D200,
      },
    },
    '&:not($disabled)': {
      cursor: 'pointer',
      '&:hover': {
        '& $checkIcon, $checkDetails $label': {
          color: SymphonyTheme.palette.B700,
        },
      },
      '&$selected': {
        '& $checkIcon, $checkDetails $label': {
          color: SymphonyTheme.palette.primary,
        },
      },
    },
  },
  disabled: {},
  selected: {},
  checkIcon: {
    flexGrow: 0,
    flexShrink: 0,
    width: '24px',
    height: '24px',
  },
  checkDetails: {
    flexGrow: 1,
    marginLeft: '8px',
  },
  label: {},
  details: {},
}));

export type RadioOption = {
  value: string,
  label: React.Node,
  details: React.Node,
  disabled?: ?boolean,
};

type Props = {
  className?: string,
  optionClassName?: string,
  selectedOptionClassName?: string,
  value: string,
  options: RadioOption[],
  disabled?: ?boolean,
  onChange?: (newValue: string) => void,
};

const RadioGroup = (props: Props) => {
  const {
    className,
    optionClassName,
    selectedOptionClassName,
    value,
    options,
    onChange,
    disabled: propDisabled = false,
  } = props;
  const classes = useStyles();
  const {disabled: contextDisabled} = useFormElementContext();
  const disabled = useMemo(
    () => (propDisabled ? propDisabled : contextDisabled),
    [contextDisabled, propDisabled],
  );

  return (
    <div className={className}>
      {options.map(option => {
        const isSelected = option.value === value;
        const isDisabled = disabled || option.disabled;
        return (
          <div
            key={`radio_option_${option.value}`}
            className={classNames(
              classes.option,
              optionClassName,
              isSelected ? selectedOptionClassName : null,
              {
                [classes.selected]: isSelected,
                [classes.disabled]: isDisabled,
              },
            )}
            onClick={() => !isDisabled && onChange && onChange(option.value)}>
            <div className={classes.checkIcon}>
              {isSelected ? (
                <RadioButtonCheckedIcon />
              ) : (
                <RadioButtonUncheckedIcon />
              )}
            </div>
            <div className={classes.checkDetails}>
              <div>
                <Text variant="body1" className={classes.label}>
                  {option.label}
                </Text>
              </div>
              <div>
                <Text variant="body2" color="gray" className={classes.details}>
                  {option.details}
                </Text>
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};

export default RadioGroup;
