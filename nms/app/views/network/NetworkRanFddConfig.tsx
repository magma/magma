/*
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
import type {NetworkRanConfigs} from '../../../generated-ts';

import Grid from '@material-ui/core/Grid';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';

import {AltFormField} from '../../components/FormField';

type Props = {
  lteRanConfigs: NetworkRanConfigs;
  setLteRanConfigs: (c: NetworkRanConfigs) => void;
};

export default function FddConfig(props: Props) {
  return (
    <Grid container xs={12}>
      <Grid item xs={12} sm={6}>
        <AltFormField label={'EARFCNDL'}>
          <OutlinedInput
            fullWidth={true}
            placeholder="Enter EARFCNDL"
            data-testid="earfcndl"
            type="number"
            value={props.lteRanConfigs.fdd_config?.earfcndl}
            onChange={({target}) =>
              props.setLteRanConfigs({
                ...props.lteRanConfigs,
                tdd_config: undefined,
                fdd_config: {
                  earfcndl: parseInt(target.value),
                  earfcnul: props.lteRanConfigs.fdd_config?.earfcnul ?? 0,
                },
              })
            }
          />
        </AltFormField>
      </Grid>
      <Grid item xs={12} sm={6}>
        <AltFormField label={'EARFCNUL'}>
          <OutlinedInput
            fullWidth={true}
            placeholder="Enter EARFCNUL"
            type="number"
            data-testid="earfcnul"
            value={props.lteRanConfigs.fdd_config?.earfcnul}
            onChange={({target}) =>
              props.setLteRanConfigs({
                ...props.lteRanConfigs,
                tdd_config: undefined,
                fdd_config: {
                  earfcndl: props.lteRanConfigs.fdd_config?.earfcndl ?? 0,
                  earfcnul: parseInt(target.value),
                },
              })
            }
          />
        </AltFormField>
      </Grid>
    </Grid>
  );
}
