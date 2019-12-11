/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {OptionProps} from './SelectMenu';

import * as React from 'react';
import ArrowDropDownIcon from '@material-ui/icons/ArrowDropDown';
import BasePopoverTrigger from './BasePopoverTrigger';
import Button from '../Button';
import SelectMenu from './SelectMenu';
import Text from '../Text';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  root: {
    justifyContent: 'flex-start',
    padding: '4px',
  },
  label: {
    textAlign: 'left',
    flexGrow: 1,
    padding: '2px 4px',
  },
  arrow: {
    marginLeft: '4px',
  },
  value: {
    fontWeight: 500,
  },
  menu: {
    margin: '8px 0px',
  },
});

type Props<TValue> = {
  className?: string,
  label: React.Node,
  options: Array<OptionProps<TValue>>,
  onChange: (value: TValue) => void | (() => void),
  selectedValue: ?TValue,
};

const Select = <TValue>({
  label,
  className,
  ...selectMenuProps
}: Props<TValue>) => {
  const classes = useStyles();
  const {options, selectedValue} = selectMenuProps;
  return (
    <BasePopoverTrigger
      popover={<SelectMenu {...selectMenuProps} className={classes.menu} />}>
      {(onShow, contextRef) => (
        <Button
          className={classNames(classes.root, className)}
          ref={contextRef}
          onClick={onShow}
          skin="regular">
          <Text className={classes.label} variant="body2">
            {label}
            {selectedValue ? ': ' : null}
            {selectedValue ? (
              <span className={classes.value}>
                {options.find(option => option.value === selectedValue)
                  ?.label ?? ''}
              </span>
            ) : null}
          </Text>
          <ArrowDropDownIcon className={classes.arrow} size="small" />
        </Button>
      )}
    </BasePopoverTrigger>
  );
};

export default Select;
