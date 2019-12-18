/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  TriggerFilterOperator_data,
  TriggerFilterOperator_data$key,
} from './__generated__/TriggerFilterOperator_data.graphql';

import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';

import {graphql, useFragment} from 'react-relay/hooks';

const query = graphql`
  fragment TriggerFilterOperator_data on ActionsFilter {
    supportedOperators {
      operatorID
      description
      dataType
    }
  }
`;

type Props = {|
  filter: TriggerFilterOperator_data$key,
  selectedOperatorID: string,
  onChange: string => void,
|};

export default function TriggerFilterOperator(props: Props) {
  const data = useFragment<TriggerFilterOperator_data>(query, props.filter);
  const supportedOperators = data.supportedOperators;
  return (
    <Select
      value={props.selectedOperatorID || supportedOperators[0].operatorID}
      onChange={({target}) => {
        props.onChange(target.value);
      }}>
      {supportedOperators.map((operator, i) => (
        <MenuItem key={i} value={operator.operatorID}>
          {operator.description}
        </MenuItem>
      ))}
    </Select>
  );
}
