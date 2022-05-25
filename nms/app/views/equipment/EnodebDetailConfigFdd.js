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

import type {DataRows} from '../../components/DataGrid';

import DataGrid from '../../components/DataGrid';
import Grid from '@material-ui/core/Grid';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_ => ({
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
}));

type Props = {
  earfcndl: number,
  earfcnul: number,
};

export function EnodeConfigFdd(props: Props) {
  const fddData: DataRows[] = [
    [
      {
        category: 'EARFCNDL',
        value: props.earfcndl,
      },
    ],
    [
      {
        category: 'EARFCNUL',
        value: props.earfcnul,
      },
    ],
  ];

  return <DataGrid data={fddData} />;
}

type EditProps = {
  earfcndl: string,
  earfcnul: string,
  setEarfcndl: string => void,
};
export default function EnodeConfigEditFdd(props: EditProps) {
  const classes = useStyles();

  return (
    <Grid container>
      <Grid item xs={6}>
        <Grid container>
          <Grid item xs={12}>
            EARFCNDL
          </Grid>
          <Grid item xs={12}>
            <OutlinedInput
              data-testid="earfcndl"
              placeholder="Enter EARFCNDL"
              className={classes.input}
              fullWidth={true}
              value={props.earfcndl}
              onChange={({target}) => props.setEarfcndl(target.value)}
            />
          </Grid>
        </Grid>
      </Grid>
      <Grid item xs={6}>
        <Grid container>
          <Grid item xs={12}>
            EARFCNUL
          </Grid>
          <Grid item xs={12}>
            <OutlinedInput
              className={classes.input}
              placeholder="Enter EARFCNUL"
              fullWidth={true}
              value={props.earfcnul}
              readOnly={true}
            />
          </Grid>
        </Grid>
      </Grid>
    </Grid>
  );
}
