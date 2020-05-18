/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CheckListItemFillingProps} from './CheckListItemFilling';
import type {YesNoResponse} from '../../work_orders/__generated__/WorkOrderDetails_workOrder.graphql';

import * as React from 'react';
import CommonStrings from '@fbcnms/strings/Strings';
import Select from '@fbcnms/ui/components/design-system/Select/Select';
import fbt from 'fbt';
import {makeStyles} from '@material-ui/styles';
import {useFormContext} from '../../../common/FormContext';
import {useMemo} from 'react';

const useStyles = makeStyles(() => ({
  select: {
    width: '100%',
  },
}));

const YesNoCheckListItemFilling = ({
  item,
  onChange,
}: CheckListItemFillingProps): React.Node => {
  const classes = useStyles();
  const options = useMemo(
    () => [
      {
        key: 'yes',
        label: CommonStrings.common.yesButton,
        value: ('YES': YesNoResponse),
      },
      {
        key: 'no',
        label: CommonStrings.common.noButton,
        value: ('NO': YesNoResponse),
      },
    ],
    [],
  );

  const updateOnChange = (value: YesNoResponse) => {
    if (!onChange) {
      return;
    }
    const modifiedItem = {
      ...item,
      yesNoResponse: value,
    };
    onChange(modifiedItem);
  };

  const form = useFormContext();

  return (
    <div>
      <Select
        className={classes.select}
        label={<fbt desc="">Select option</fbt>}
        options={options}
        selectedValue={item.yesNoResponse ?? null}
        onChange={value => updateOnChange(value)}
        disabled={form.alerts.missingPermissions.detected}
      />
    </div>
  );
};

export default YesNoCheckListItemFilling;
