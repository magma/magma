/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow
 * @format
 */

import type {TextPairingContainerProps} from '../helpers/TextPairingContainer';

import CheckBoxIcon from '@material-ui/icons/CheckBox';
import CheckBoxOutlineBlankIcon from '@material-ui/icons/CheckBoxOutlineBlank';
import IndeterminateCheckBoxIcon from '@material-ui/icons/IndeterminateCheckBox';
import React from 'react';
import SymphonyTheme from '../../../theme/symphony';
import TextPairingContainer from '../helpers/TextPairingContainer';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';
import {useFormElementContext} from '../Form/FormElementContext';
import {useMemo} from 'react';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '24px',
    height: '24px',
    '&:not($disabled)': {
      cursor: 'pointer',
      '&:hover': {
        '& $selection, & $noSelection': {
          fill: SymphonyTheme.palette.B700,
        },
      },
      '& $noSelection': {
        fill: SymphonyTheme.palette.D200,
      },
      '& $selection': {
        fill: SymphonyTheme.palette.primary,
      },
    },
  },
  disabled: {
    '& $noSelection, & $selection': {
      fill: SymphonyTheme.palette.disabled,
    },
  },
  selection: {},
  noSelection: {},
}));

export type SelectionType = 'checked' | 'unchecked';

type Props = $ReadOnly<{|
  checked: boolean,
  indeterminate?: boolean,
  disabled?: ?boolean,
  onChange?: (selection: SelectionType) => void,
  onClick?: ?(SyntheticMouseEvent<Element>) => void,
  ...TextPairingContainerProps,
|}>;

const Checkbox = (props: Props) => {
  const {
    checked,
    indeterminate,
    onChange,
    onClick,
    disabled: propDisabled = false,
    ...TextPairingContainerProps
  } = props;
  const classes = useStyles();
  const CheckboxIcon = indeterminate
    ? IndeterminateCheckBoxIcon
    : checked
    ? CheckBoxIcon
    : CheckBoxOutlineBlankIcon;

  const {disabled: contextDisabled} = useFormElementContext();
  const disabled = useMemo(
    () => (propDisabled ? propDisabled : contextDisabled),
    [contextDisabled, propDisabled],
  );

  return (
    <TextPairingContainer {...TextPairingContainerProps} disabled={disabled}>
      <div
        className={classNames(classes.root, {
          [classes.disabled]: disabled,
        })}
        onClick={e => {
          if (disabled) {
            return;
          }
          if (onChange) {
            onChange(
              indeterminate ? 'unchecked' : checked ? 'unchecked' : 'checked',
            );
          }
          if (onClick) {
            onClick(e);
          }
        }}>
        <CheckboxIcon
          className={classNames({
            [classes.selection]: checked || indeterminate,
            [classes.noSelection]: !checked && !indeterminate,
          })}
        />
      </div>
    </TextPairingContainer>
  );
};

Checkbox.defaultProps = {
  checked: false,
  indeterminate: false,
};

export default Checkbox;
