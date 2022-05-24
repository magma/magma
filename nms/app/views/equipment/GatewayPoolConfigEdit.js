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
/*[object Object]*/
import type {GatewayPoolEditProps} from './GatewayPoolEdit';
import type {mutable_cellular_gateway_pool} from '../../../generated/MagmaAPIBindings';

import Button from '@material-ui/core/Button';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import FormLabel from '@material-ui/core/FormLabel';
import GatewayPoolsContext from '../../components/context/GatewayPoolsContext';
import List from '@material-ui/core/List';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';

import {AltFormField} from '../../components/FormField';
import {DEFAULT_GW_POOL_CONFIG} from '../../components/GatewayUtils';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';

const useStyles = makeStyles(_ => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '50%',
    fullWidth: true,
  },
}));
export default function ConfigEdit(props: GatewayPoolEditProps) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const [error, setError] = useState('');
  const ctx = useContext(GatewayPoolsContext);
  const [gwPool, setGwPool] = useState<mutable_cellular_gateway_pool>(
    Object.keys(props.gwPool || {}).length > 0
      ? props.gwPool
      : DEFAULT_GW_POOL_CONFIG,
  );
  const handleGwPoolConfigChange = (value: number) => {
    const newConfig = {
      ...gwPool,
      config: {...gwPool.config, ['mme_group_id']: value},
    };
    setGwPool(newConfig);
  };
  const onSave = async () => {
    try {
      await ctx.setState(gwPool.gateway_pool_id, gwPool);
      enqueueSnackbar('Gateway Pool saved successfully', {
        variant: 'success',
      });
      props.onSave(gwPool);
    } catch (e) {
      setError(e.response?.data?.message ?? e.message);
    }
  };
  const classes = useStyles();

  return (
    <>
      <DialogContent data-testid="configEdit">
        <List>
          {error !== '' && (
            <AltFormField label={''}>
              <FormLabel data-testid="configEditError" error>
                {error}
              </FormLabel>
            </AltFormField>
          )}
          <AltFormField label={'Name'}>
            <OutlinedInput
              data-testid="name"
              className={classes.input}
              placeholder="Enter Name"
              fullWidth={true}
              value={gwPool.gateway_pool_name}
              onChange={({target}) =>
                setGwPool({...gwPool, gateway_pool_name: target.value})
              }
            />
          </AltFormField>
          <AltFormField label={'ID'}>
            <OutlinedInput
              data-testid="poolId"
              className={classes.input}
              placeholder="Ex: pool1"
              fullWidth={true}
              value={gwPool.gateway_pool_id}
              readOnly={Object.keys(props.gwPool).length > 0 ? false : true}
              onChange={({target}) =>
                setGwPool({...gwPool, gateway_pool_id: target.value})
              }
            />
          </AltFormField>
          <AltFormField label={'MME Group ID'}>
            <OutlinedInput
              data-testid="mmeGroupId"
              className={classes.input}
              placeholder="Ex: 1"
              fullWidth={true}
              type="number"
              value={gwPool.config.mme_group_id}
              onChange={({target}) => {
                handleGwPoolConfigChange(parseInt(target.value));
              }}
            />
          </AltFormField>
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>Cancel</Button>
        <Button onClick={onSave} variant="contained" color="primary">
          {'Save And Continue'}
        </Button>
      </DialogActions>
    </>
  );
}
