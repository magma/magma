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

import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import LoadingFillerBackdrop from './LoadingFillerBackdrop';
import MagmaAPI from '../api/MagmaAPI';
import MenuItem from '@mui/material/MenuItem';
import OutlinedInput from '@mui/material/OutlinedInput';
import React, {useState} from 'react';
import Select from '@mui/material/Select';
import nullthrows from '../../shared/util/nullthrows';
import useMagmaAPI from '../api/useMagmaAPI';
import {AltFormField} from './FormField';
import {getErrorMessage} from '../util/ErrorUtils';
import {useEnqueueSnackbar} from '../hooks/useSnackbar';
import {useParams} from 'react-router-dom';

type GatewayData = {
  gatewayID: string;
  name: string;
  description: string;
  hardwareID: string;
  challengeKey: string;
  tier: string;
};

export const MAGMAD_DEFAULT_CONFIGS = {
  autoupgrade_enabled: true,
  autoupgrade_poll_interval: 300,
  checkin_interval: 60,
  checkin_timeout: 10,
};

export const EMPTY_GATEWAY_FIELDS = {
  gatewayID: '',
  name: '',
  description: '',
  hardwareID: '',
  challengeKey: '',
  tier: '',
};

type Props = {
  onClose: () => void;
  onSave: (data: GatewayData) => Promise<void>;
};

export default function AddGatewayDialog(props: Props) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const [values, setValues] = useState(EMPTY_GATEWAY_FIELDS);

  const params = useParams();
  const networkID = nullthrows(params.networkId);
  const {response: tiers, isLoading} = useMagmaAPI(
    MagmaAPI.upgrades.networksNetworkIdTiersGet,
    {
      networkId: networkID,
    },
  );

  if (isLoading || !tiers) {
    return <LoadingFillerBackdrop />;
  }

  const onSave = () => {
    if (
      !values.name ||
      !values.description ||
      !values.hardwareID ||
      !values.gatewayID ||
      !values.challengeKey
    ) {
      enqueueSnackbar('Please complete all fields', {variant: 'error'});
      return;
    }

    try {
      void props.onSave(values);
    } catch (e) {
      enqueueSnackbar(getErrorMessage(e), {
        variant: 'error',
      });
    }
  };

  return (
    <Dialog open={true} onClose={props.onClose} maxWidth="md" scroll="body">
      <DialogTitle>Add Gateway</DialogTitle>
      <DialogContent>
        <AddGatewayFields onChange={setValues} values={values} tiers={tiers} />
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>Cancel</Button>
        <Button onClick={onSave} variant="contained" color="primary">
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}

export const AddGatewayFields = (props: {
  values: GatewayData;
  onChange: (data: GatewayData) => void;
  tiers: Array<string>;
}) => {
  return (
    <>
      <AltFormField label="Gateway Name">
        <OutlinedInput
          fullWidth={true}
          value={props.values.name}
          onChange={({target}) =>
            props.onChange({...props.values, name: target.value})
          }
          placeholder="Gateway 1"
        />
      </AltFormField>
      <AltFormField label="Gateway Description">
        <OutlinedInput
          fullWidth={true}
          value={props.values.description}
          onChange={({target}) =>
            props.onChange({...props.values, description: target.value})
          }
          placeholder="Sample Gateway description"
        />
      </AltFormField>
      <AltFormField label="Hardware UUID">
        <OutlinedInput
          fullWidth={true}
          value={props.values.hardwareID}
          onChange={({target}) =>
            props.onChange({...props.values, hardwareID: target.value})
          }
          placeholder="Eg. 4dfe212f-df33-4cd2-910c-41892a042fee"
        />
      </AltFormField>
      <AltFormField label="Gateway ID">
        <OutlinedInput
          fullWidth={true}
          value={props.values.gatewayID}
          onChange={({target}) =>
            props.onChange({...props.values, gatewayID: target.value})
          }
          placeholder="<country>_<org>_<location>_<sitenumber>"
        />
      </AltFormField>
      <AltFormField label="Challenge Key">
        <OutlinedInput
          fullWidth={true}
          value={props.values.challengeKey}
          onChange={({target}) =>
            props.onChange({...props.values, challengeKey: target.value})
          }
          placeholder="A base64 bytestring of the key in DER format"
        />
      </AltFormField>
      <AltFormField label="Upgrade Tier">
        <Select
          fullWidth={true}
          variant={'outlined'}
          inputProps={{'data-testid': 'upgradeTier'}}
          value={props.values.tier}
          onChange={({target}) =>
            props.onChange({...props.values, tier: target.value})
          }>
          {props.tiers.map(tier => (
            <MenuItem key={tier} value={tier}>
              {tier}
            </MenuItem>
          ))}
        </Select>
      </AltFormField>
    </>
  );
};
