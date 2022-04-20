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

import type {MultiSelectMenuProps} from './MultiSelectMenu';

import * as React from 'react';
import ArrowDropDownIcon from '@material-ui/icons/ArrowDropDown';
import BasePopoverTrigger from '../ContexualLayer/BasePopoverTrigger';
import Button from '../Button';
import MultiSelectMenu from './MultiSelectMenu';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    justifyContent: 'flex-start',
    padding: '4px',
    border: `1px solid ${symphony.palette.D100}`,
  },
  value: {
    fontWeight: 500,
  },
  menu: {
    margin: '8px 0px',
  },
  label: {
    fontWeight: 400,
  },
}));

type Props<TValue> = $ReadOnly<{|
  ...MultiSelectMenuProps<TValue>,
  label: React.Node,
  disabled?: boolean,
|}>;

const MultiSelect = <TValue>({
  label,
  className,
  disabled,
  ...selectMenuProps
}: Props<TValue>) => {
  const classes = useStyles();
  const {selectedValues} = selectMenuProps;
  return (
    <BasePopoverTrigger
      popover={
        <MultiSelectMenu
          {...selectMenuProps}
          className={classes.menu}
          size="normal"
        />
      }>
      {(onShow, _onHide, contextRef) => (
        <Button
          className={classNames(classes.root, className)}
          ref={contextRef}
          onClick={onShow}
          skin="regular"
          rightIcon={ArrowDropDownIcon}
          disabled={disabled}>
          <span className={classes.label}>{label}</span>
          {selectedValues.length > 0 ? ': ' : null}
          {selectedValues.length === 1 ? (
            <span className={classes.value} key={String(selectedValues[0])}>
              {selectedValues[0].label ?? ''}
            </span>
          ) : null}
          {selectedValues.length > 1 ? (
            <span className={classes.value}>
              <fbt desc="Amount of selected items">
                <fbt:param name="num_selected" number={true}>
                  {selectedValues.length}
                </fbt:param>
                Selected
              </fbt>
            </span>
          ) : null}
        </Button>
      )}
    </BasePopoverTrigger>
  );
};

export default MultiSelect;
