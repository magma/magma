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

import type {ButtonProps} from '../Button';
import type {SelectMenuProps} from './SelectMenu';

import * as React from 'react';
import ArrowDropDownIcon from '@material-ui/icons/ArrowDropDown';
import BasePopoverTrigger from '../ContexualLayer/BasePopoverTrigger';
import Button from '../Button';
import SelectMenu from './SelectMenu';
import Text from '../../../../../app/theme/design-system/Text';
import classNames from 'classnames';
import theme, {colors} from '../../../../../app/theme/default';
import {makeStyles} from '@material-ui/styles';
import {useFormElementContext} from '../Form/FormElementContext';
import {useMemo} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    justifyContent: 'flex-start',
    '&&': {
      padding: '4px',
    },
    border: `1px solid ${colors.primary.mercury}`,
    '&$disabled': {
      backgroundColor: colors.primary.selago,
      color: colors.primary.brightGray,
    },
  },
  disabled: {
    '&&': {
      '&&': {
        color: colors.primary.brightGray,
        fill: colors.primary.brightGray,
      },
    },
  },
  formValue: {
    ...theme.typography.body2,
  },
  menu: {
    margin: '8px 0px',
  },
  label: {
    fontWeight: theme.typography.body2.fontWeight,
  },
}));

type Props<TValue> = $ReadOnly<{|
  className?: string,
  label?: React.Node,
  ...ButtonProps,
  ...SelectMenuProps<TValue>,
|}>;

const Select = <TValue>(props: Props<TValue>) => {
  const {
    label,
    className,
    disabled: disabledProp,
    skin,
    tooltip,
    useEllipsis = true,
    variant,
    ...selectMenuProps
  } = props;
  const {selectedValue, options} = selectMenuProps;
  const classes = useStyles();
  const {disabled: contextDisabled} = useFormElementContext();
  const disabled = useMemo(
    () => (disabledProp ? disabledProp : contextDisabled),
    [disabledProp, contextDisabled],
  );
  return (
    <BasePopoverTrigger
      popover={<SelectMenu {...selectMenuProps} className={classes.menu} />}>
      {(onShow, _onHide, contextRef) => (
        <Button
          className={classNames(classes.root, className, {
            [classes.disabled]: disabled,
          })}
          ref={contextRef}
          onClick={onShow}
          skin={skin ?? 'regular'}
          variant={variant}
          disabled={disabled}
          rightIcon={ArrowDropDownIcon}
          rightIconClass={classNames({[classes.disabled]: disabled})}
          tooltip={tooltip}
          useEllipsis={useEllipsis}>
          <Text variant="body2">
            <span className={classes.label}>{label}</span>
            {selectedValue != null && !!label ? ': ' : null}
            {selectedValue != null ? (
              <span
                className={
                  classNames({
                    [classes.formValue]: !label,
                    [classes.disabled]: !label && disabled,
                  }) || null
                }>
                {options.find(option => option.value === selectedValue)
                  ?.label ?? ''}
              </span>
            ) : null}
          </Text>
        </Button>
      )}
    </BasePopoverTrigger>
  );
};

export default Select;
