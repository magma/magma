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
 * @flow strict-local
 * @format
 */

import typeof SvgIcon from '@material-ui/core/@@SvgIcon';
import type {
  ErrorHandlingProps,
  PermissionHandlingProps,
} from '../Form/FormAction';

import * as React from 'react';
import FormAction from '../../../../../fbc_js_core/ui/components/design-system/Form/FormAction';
import FormElementContext from '../../../../../fbc_js_core/ui/components/design-system/Form/FormElementContext';
import Text from '../Text';
import classNames from 'classnames';
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  option: {
    display: 'flex',
    alignItems: 'center',
    padding: '8px 16px',
    cursor: 'pointer',
    '&:not($disabled)&:hover': {
      backgroundColor: symphony.palette.background,
    },
    '&$optionWithLeftAux': {
      paddingLeft: '12px',
      paddingTop: '6px',
      paddingBottom: '6px',
    },
  },
  optionWithLeftAux: {},
  disabled: {
    opacity: 0.38,
    cursor: 'not-allowed',
  },
  label: {
    flexGrow: 1,
  },
  checkIcon: {
    marginLeft: '6px',
    color: symphony.palette.primary,
  },
  leftAux: {
    display: 'inline-flex',
    marginRight: '8px',
  },
  contentContainer: {
    display: 'flex',
    flexDirection: 'column',
  },
}));

export type MenuItemLeftAux = $ReadOnly<
  | {|
      type: 'icon',
      icon: SvgIcon,
    |}
  | {
      type: 'node',
      node: React.Node,
    },
>;

export type SelectMenuItemBaseProps<TValue> = $ReadOnly<{|
  label: React.Node,
  value: TValue,
  isSelected?: boolean,
  className?: ?string,
  leftAux?: MenuItemLeftAux,
  secondaryText?: React.Node,
  disabled?: boolean,
  skin?: 'inherit' | 'red',
  ...PermissionHandlingProps,
  ...ErrorHandlingProps,
|}>;

type Props<TValue> = $ReadOnly<{|
  ...SelectMenuItemBaseProps<TValue>,
  onClick: (value: TValue) => void,
|}>;

const SelectMenuItem = <TValue>({
  label,
  value,
  onClick,
  isSelected = false,
  hideOnMissingPermissions = false,
  className,
  leftAux,
  secondaryText,
  skin = 'inherit',
  disabled: disabledProp = false,
  ...actionProps
}: Props<TValue>) => {
  const classes = useStyles();
  const LeftIcon = leftAux?.type === 'icon' ? leftAux.icon : null;
  const coercedSkin = disabledProp
    ? 'inherit'
    : skin === 'red'
    ? 'error'
    : skin;
  return (
    <FormAction
      {...actionProps}
      disabled={disabledProp}
      hideOnMissingPermissions={hideOnMissingPermissions}>
      <FormElementContext.Consumer>
        {({disabled}) => {
          return (
            <div
              className={classNames(classes.option, className, {
                [classes.disabled]: disabled,
                [classes.optionWithLeftAux]: leftAux != null,
              })}
              onClick={disabled ? null : () => onClick(value)}>
              {leftAux != null && (
                <div className={classes.leftAux}>
                  {leftAux.type === 'icon'
                    ? LeftIcon != null && (
                        <LeftIcon
                          color={isSelected ? 'primary' : coercedSkin}
                          size="small"
                        />
                      )
                    : leftAux.node}
                </div>
              )}
              <div className={classes.contentContainer}>
                <Text
                  className={classes.label}
                  variant="body2"
                  color={isSelected ? 'primary' : coercedSkin}>
                  {label}
                </Text>
                {secondaryText != null && (
                  <Text color="gray" variant="caption">
                    {secondaryText}
                  </Text>
                )}
              </div>
            </div>
          );
        }}
      </FormElementContext.Consumer>
    </FormAction>
  );
};

export default SelectMenuItem;
