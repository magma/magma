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

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import LoadingFillerBackdrop from './LoadingFillerBackdrop';
import MagmaAPI from '../../api/MagmaAPI';
import MenuItem from '@material-ui/core/MenuItem';
import React, {useState} from 'react';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';
import nullthrows from '../../shared/util/nullthrows';
import useMagmaAPI from '../../api/useMagmaAPI';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '../hooks/useSnackbar';
import {useParams} from 'react-router-dom';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

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
      // eslint-disable-next-line @typescript-eslint/no-floating-promises
      props.onSave(values);
    } catch (e) {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-argument,@typescript-eslint/no-unsafe-member-access
      enqueueSnackbar(e?.response?.data?.message || e?.message || e, {
        variant: 'error',
      });
    }
  };

  return (
    <Dialog open={true} onClose={props.onClose}>
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
  const classes = useStyles();

  return (
    <>
      <TextField
        label="Gateway Name"
        className={classes.input}
        value={props.values.name}
        onChange={({target}) =>
          props.onChange({...props.values, name: target.value})
        }
        placeholder="Gateway 1"
      />
      <TextField
        label="Gateway Description"
        className={classes.input}
        value={props.values.description}
        onChange={({target}) =>
          props.onChange({...props.values, description: target.value})
        }
        placeholder="Sample Gateway description"
      />
      <TextField
        label="Hardware UUID"
        className={classes.input}
        value={props.values.hardwareID}
        onChange={({target}) =>
          props.onChange({...props.values, hardwareID: target.value})
        }
        placeholder="Eg. 4dfe212f-df33-4cd2-910c-41892a042fee"
      />
      <TextField
        label="Gateway ID"
        className={classes.input}
        value={props.values.gatewayID}
        onChange={({target}) =>
          props.onChange({...props.values, gatewayID: target.value})
        }
        placeholder="<country>_<org>_<location>_<sitenumber>"
      />
      <TextField
        label="Challenge Key"
        className={classes.input}
        value={props.values.challengeKey}
        onChange={({target}) =>
          props.onChange({...props.values, challengeKey: target.value})
        }
        placeholder="A base64 bytestring of the key in DER format"
      />
      <FormControl className={classes.input}>
        <InputLabel htmlFor="types">Upgrade Tier</InputLabel>
        <Select
          className={classes.input}
          value={props.values.tier}
          onChange={({target}) =>
            props.onChange({...props.values, tier: target.value as string})
          }>
          {props.tiers.map(tier => (
            <MenuItem key={tier} value={tier}>
              {tier}
            </MenuItem>
          ))}
        </Select>
      </FormControl>
    </>
  );
};
