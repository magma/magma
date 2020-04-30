/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React from 'react';
import classNames from 'classnames';
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';
import {useFormElementContext} from '../Form/FormElementContext';
import {useMemo} from 'react';

const SIZE = 16;
const SIZE_ATTR = `${SIZE}px`;

const useStyles = makeStyles(() => ({
  root: {
    display: 'inline-block',
    minWidth: `${2 * SIZE}px`,
    maxWidth: `${2 * SIZE}px`,
    minHeight: SIZE_ATTR,
    maxHeight: SIZE_ATTR,
    padding: '2px',
    clear: 'both',
    background: symphony.palette.D100,
    borderRadius: SIZE_ATTR,
    '&:not($disabled)': {
      cursor: 'pointer',
      '&$checked': {
        background: symphony.palette.primary,
        '&$critical': {
          background: symphony.palette.R600,
        },
      },
      '&:hover': {
        background: symphony.palette.D200,
        '&$checked': {
          background: symphony.palette.B700,
          '&$critical': {
            background: symphony.palette.R800,
          },
        },
      },
    },
  },
  checked: {
    '& $toggle': {
      float: 'right',
    },
  },
  critical: {},
  disabled: {
    background: symphony.palette.disabled,
  },
  toggle: {
    height: SIZE_ATTR,
    width: SIZE_ATTR,
    background: symphony.palette.white,
    borderRadius: '50%',
  },
}));

type Skin = 'regular' | 'critical';

type Props = $ReadOnly<{|
  className?: ?string,
  checked: boolean,
  skin?: ?Skin,
  disabled?: ?boolean,
  onChange?: ?(checked: boolean) => void,
  onClick?: ?(SyntheticMouseEvent<Element>) => void,
|}>;

const Switch = (props: Props) => {
  const {
    className,
    checked,
    onChange,
    onClick,
    skin = 'regular',
    disabled: propDisabled = false,
  } = props;
  const classes = useStyles();
  const {disabled: contextDisabled} = useFormElementContext();
  const disabled = useMemo(
    () => (propDisabled ? propDisabled : contextDisabled),
    [contextDisabled, propDisabled],
  );

  return (
    <div
      className={classNames(
        classes.root,
        {
          [classes.disabled]: disabled,
          [classes.checked]: checked,
          [classes.critical]: skin === 'critical',
        },
        className,
      )}
      onClick={e => {
        if (disabled) {
          return;
        }
        if (onChange) {
          onChange(!checked);
        }
        if (onClick) {
          onClick(e);
        }
      }}>
      <div className={classes.toggle} />
    </div>
  );
};

export default Switch;
