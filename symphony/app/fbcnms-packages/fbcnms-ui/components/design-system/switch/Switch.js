/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {TextPairingContainerProps} from '../helpers/TextPairingContainer';

import React from 'react';
import TextPairingContainer from '../helpers/TextPairingContainer';
import classNames from 'classnames';
import symphony from '../../../theme/symphony';
import {hexToRgb} from '../../../utils/displayUtils';
import {makeStyles} from '@material-ui/styles';
import {useFormElementContext} from '../Form/FormElementContext';
import {useMemo} from 'react';

const SIZE = 16;
const SIZE_ATTR = `${SIZE}px`;

const useStyles = makeStyles(() => ({
  toggleContainer: {
    display: 'block',
    minWidth: `${2 * SIZE}px`,
    maxWidth: `${2 * SIZE}px`,
    minHeight: SIZE_ATTR,
    maxHeight: SIZE_ATTR,
    padding: '2px',
    clear: 'both',
    background: symphony.palette.D100,
    '&$critical': {
      background: symphony.palette.R600,
    },
    borderRadius: SIZE_ATTR,
    '&:not($disabled)': {
      cursor: 'pointer',
      '&$checked:not($critical)': {
        background: symphony.palette.primary,
      },
      '&:hover': {
        background: symphony.palette.D200,
        '&$checked:not($critical)': {
          background: symphony.palette.B700,
        },
        '&$critical': {
          background: symphony.palette.R800,
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
    '&$critical': {
      background: `rgba(${hexToRgb(symphony.palette.R600)},0.38)`,
    },
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
  checked: boolean,
  skin?: ?Skin,
  onChange?: ?(checked: boolean) => void,
  onClick?: ?(SyntheticMouseEvent<Element>) => void,
  ...TextPairingContainerProps,
|}>;

const Switch = (props: Props) => {
  const {
    checked,
    onChange,
    onClick,
    skin = 'regular',
    disabled: propDisabled = false,
    ...textPairingContainerProps
  } = props;
  const classes = useStyles();
  const {disabled: contextDisabled} = useFormElementContext();
  const disabled = useMemo(
    () => (propDisabled ? propDisabled : contextDisabled),
    [contextDisabled, propDisabled],
  );

  return (
    <TextPairingContainer
      {...textPairingContainerProps}
      margin="wide"
      disabled={disabled}>
      <div
        className={classNames(classes.toggleContainer, {
          [classes.disabled]: disabled,
          [classes.checked]: checked,
          [classes.critical]: skin === 'critical',
        })}
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
    </TextPairingContainer>
  );
};

export default Switch;
