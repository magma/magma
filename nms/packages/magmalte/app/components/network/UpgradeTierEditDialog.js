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

import type {tier} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormGroup from '@material-ui/core/FormGroup';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import TextField from '@material-ui/core/TextField';

import nullthrows from '@fbcnms/util/nullthrows';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

type Props = {
  onSave: tier => void,
  onCancel: () => void,
  tier?: tier,
};

export default function UpgradeTierEditDialog(props: Props) {
  const {match} = useRouter();
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
      MagmaV1API.postNetworksByNetworkIdTiers({
        networkId: nullthrows(match.params.networkId),
        tier,
      })
        .then(() => props.onSave(tier))
        .catch(e =>
          enqueueSnackbar(e.response.data.message, {variant: 'error'}),
        );
    } else {
      MagmaV1API.putNetworksByNetworkIdTiersByTierId({
        networkId: nullthrows(match.params.networkId),
        tierId: tier.id,
        tier,
      })
        .then(_resp => props.onSave(tier))
        .catch(e =>
          enqueueSnackbar(e.response.data.message, {variant: 'error'}),
        );
    }
  };

  return (
    <Dialog open={true} onClose={props.onCancel} scroll="body">
      <DialogTitle>
        {props.tier ? 'Edit Upgrade Tier' : 'Add Upgrade Tier'}
      </DialogTitle>
      <DialogContent>
        <FormGroup row>
          <TextField
            required
            label="Tier ID"
            placeholder="E.g. t1"
            margin="normal"
            disabled={Boolean(props.tier)}
            value={tier.id}
            onChange={({target}) => setTier({...tier, id: target.value})}
          />
        </FormGroup>
        <FormGroup row>
          <TextField
            required
            label="Tier Name"
            placeholder="E.g. Example Tier"
            margin="normal"
            value={tier.name}
            onChange={({target}) => setTier({...tier, name: target.value})}
          />
        </FormGroup>
        <FormGroup row>
          <TextField
            required
            label="Tier Version"
            placeholder="E.g. 1.0.0-0"
            margin="normal"
            value={tier.version}
            onChange={({target}) => setTier({...tier, version: target.value})}
          />
        </FormGroup>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onCancel} skin="regular">
          Cancel
        </Button>
        <Button onClick={onSave}>Save</Button>
      </DialogActions>
    </Dialog>
  );
}
