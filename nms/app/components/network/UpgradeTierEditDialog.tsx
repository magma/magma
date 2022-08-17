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

import type {Tier} from '../../../generated';

import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import React from 'react';

import MagmaAPI from '../../api/MagmaAPI';
import nullthrows from '../../../shared/util/nullthrows';
import {AltFormField} from '../FormField';
import {OutlinedInput} from '@mui/material';
import {getErrorMessage} from '../../util/ErrorUtils';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useParams} from 'react-router-dom';
import {useState} from 'react';

type Props = {
  onSave: (tier: Tier) => void;
  onCancel: () => void;
  tier?: Tier;
};

export default function UpgradeTierEditDialog(props: Props) {
  const params = useParams();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [tier, setTier] = useState(
    props.tier || {
      id: '',
      name: '',
      version: '',
      images: [],
      gateways: [],
    },
  );

  const onSave = () => {
    if (!props.tier) {
      void MagmaAPI.upgrades
        .networksNetworkIdTiersPost({
          networkId: nullthrows(params.networkId),
          tier,
        })
        .then(() => props.onSave(tier))
        .catch(e => enqueueSnackbar(getErrorMessage(e), {variant: 'error'}));
    } else {
      void MagmaAPI.upgrades
        .networksNetworkIdTiersTierIdPut({
          networkId: nullthrows(params.networkId),
          tierId: tier.id,
          tier,
        })
        .then(() => props.onSave(tier))
        .catch(e => enqueueSnackbar(getErrorMessage(e), {variant: 'error'}));
    }
  };

  return (
    <Dialog open={true} onClose={props.onCancel} scroll="body">
      <DialogTitle>
        {props.tier ? 'Edit Upgrade Tier' : 'Add Upgrade Tier'}
      </DialogTitle>
      <DialogContent>
        <AltFormField label="Tier ID">
          <OutlinedInput
            required
            fullWidth
            placeholder="E.g. t1"
            disabled={Boolean(props.tier)}
            value={tier.id}
            onChange={({target}) => setTier({...tier, id: target.value})}
          />
        </AltFormField>
        <AltFormField label="Tier Name">
          <OutlinedInput
            required
            fullWidth
            label="Tier Name"
            placeholder="E.g. Example Tier"
            value={tier.name}
            onChange={({target}) => setTier({...tier, name: target.value})}
          />
        </AltFormField>
        <AltFormField label="Tier Version">
          <OutlinedInput
            required
            fullWidth
            placeholder="E.g. 1.0.0-0"
            value={tier.version}
            onChange={({target}) => setTier({...tier, version: target.value})}
          />
        </AltFormField>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onCancel}>Cancel</Button>
        <Button onClick={onSave} variant="contained" color="primary">
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}
