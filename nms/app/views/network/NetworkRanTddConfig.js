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
 *
 * @flow strict-local
 * @format
 */
import type {network_ran_configs} from '../../../generated/MagmaAPIBindings';

import Grid from '@material-ui/core/Grid';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';

// $FlowFixMe migrated to typescript
import {AltFormField} from '../../components/FormField';

type Props = {
  lteRanConfigs: network_ran_configs,
  setLteRanConfigs: network_ran_configs => void,
};

export default function TddConfig(props: Props) {
  return (
    <>
      <AltFormField label={'EARFCNDL'}>
        <OutlinedInput
          fullWidth={true}
          placeholder="Enter EARFCNDL"
          type="number"
          data-testid="earfcndl"
          value={props.lteRanConfigs.tdd_config?.earfcndl}
          onChange={({target}) =>
            props.setLteRanConfigs({
              ...props.lteRanConfigs,
              fdd_config: undefined,
              tdd_config: {
                special_subframe_pattern:
                  props.lteRanConfigs.tdd_config?.special_subframe_pattern ?? 0,
                subframe_assignment:
                  props.lteRanConfigs.tdd_config?.subframe_assignment ?? 0,
                earfcndl: parseInt(target.value),
              },
            })
          }
        />
      </AltFormField>
      <Grid container>
        <Grid item xs={12} sm={6}>
          <AltFormField label={'Special Subframe Pattern'}>
            <OutlinedInput
              fullWidth={true}
              placeholder="Special Subframe Pattern"
              type="number"
              data-testid="specialSubframePattern"
              value={props.lteRanConfigs.tdd_config?.special_subframe_pattern}
              onChange={({target}) =>
                props.setLteRanConfigs({
                  ...props.lteRanConfigs,
                  fdd_config: undefined,
                  tdd_config: {
                    special_subframe_pattern: parseInt(target.value),
                    subframe_assignment:
                      props.lteRanConfigs.tdd_config?.subframe_assignment ?? 0,
                    earfcndl: props.lteRanConfigs.tdd_config?.earfcndl ?? 0,
                  },
                })
              }
            />
          </AltFormField>
        </Grid>
        <Grid item xs={12} sm={6}>
          <AltFormField label={'Subframe Assignment'}>
            <OutlinedInput
              fullWidth={true}
              placeholder="Subframe Assignment"
              type="number"
              data-testid="subframeAssignment"
              value={props.lteRanConfigs.tdd_config?.subframe_assignment}
              onChange={({target}) => {
                props.setLteRanConfigs({
                  ...props.lteRanConfigs,
                  fdd_config: undefined,
                  tdd_config: {
                    subframe_assignment: parseInt(target.value),
                    special_subframe_pattern:
                      props.lteRanConfigs.tdd_config
                        ?.special_subframe_pattern ?? 0,
                    earfcndl: props.lteRanConfigs.tdd_config?.earfcndl ?? 0,
                  },
                });
              }}
            />
          </AltFormField>
        </Grid>
      </Grid>
    </>
  );
}
