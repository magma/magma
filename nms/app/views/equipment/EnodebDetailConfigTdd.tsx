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
 */

import DataGrid from '../../components/DataGrid';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import {AltFormField} from '../../components/FormField';
import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import type {DataRows} from '../../components/DataGrid';

const useStyles = makeStyles({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '80%',
  },
  itemTitle: {
    color: colors.primary.comet,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  itemValue: {
    color: colors.primary.brightGray,
  },
});

type Props = {
  earfcndl: number;
  specialSubframePattern: number;
  subframeAssignment: number;
};

export function EnodeConfigTdd(props: Props) {
  const tddData: Array<DataRows> = [
    [
      {
        category: 'EARFCNDL',
        value: props.earfcndl,
      },
    ],
    [
      {
        category: 'Special Subframe Pattern',
        value: props.specialSubframePattern,
      },
    ],
    [
      {
        category: 'Subframe Assignment',
        value: props.subframeAssignment,
      },
    ],
  ];

  return <DataGrid data={tddData} />;
}

type EditProps = {
  earfcndl: string;
  specialSubframePattern: string;
  subframeAssignment: string;
  setEarfcndl: (earfcndl: string) => void;
  setSpecialSubframePattern: (pattern: string) => void;
  setSubframeAssignment: (assignment: string) => void;
};

export default function EnodeConfigEditTdd(props: EditProps) {
  const classes = useStyles();

  return (
    <>
      <AltFormField label={'EARFCNDL'}>
        <OutlinedInput
          data-testid="earfcndl"
          placeholder="Enter EARFCNDL"
          className={classes.input}
          fullWidth={true}
          value={props.earfcndl}
          onChange={({target}) => props.setEarfcndl(target.value)}
        />
      </AltFormField>
      <AltFormField label={'Special Subframe Pattern'}>
        <OutlinedInput
          className={classes.input}
          placeholder="Enter Special Subframe Pattern"
          fullWidth={true}
          value={props.specialSubframePattern}
          onChange={({target}) => props.setSpecialSubframePattern(target.value)}
        />
      </AltFormField>
      <AltFormField label={'Subframe Assignment'}>
        <OutlinedInput
          className={classes.input}
          placeholder="Enter Subframe Assignment"
          fullWidth={true}
          value={props.subframeAssignment}
          onChange={({target}) => props.setSubframeAssignment(target.value)}
        />
      </AltFormField>
    </>
  );
}
