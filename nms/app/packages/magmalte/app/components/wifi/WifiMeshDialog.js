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

import type {mesh_wifi_configs} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import Checkbox from '@material-ui/core/Checkbox';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormGroup from '@material-ui/core/FormGroup';
import KeyValueFields from '@fbcnms/magmalte/app/components/KeyValueFields';
import LoadingFillerBackdrop from '@fbcnms/ui/components/LoadingFillerBackdrop';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import TextField from '@material-ui/core/TextField';

import nullthrows from '@fbcnms/util/nullthrows';
import {additionalPropsToArray, additionalPropsToObject} from './WifiUtils';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
  backdrop: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    top: 0,
    bottom: 0,
    left: 0,
    right: 0,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    position: 'fixed',
    zIndex: '13000',
  },
}));

type Props = {
  onCancel: () => void,
  onSave: string => void,
};

export default function WifiMeshDialog(props: Props) {
  const classes = useStyles();
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();

  const [meshID, setMeshID] = useState(match.params.meshID || '');
  const [configs, setConfigs] = useState<mesh_wifi_configs>({});
  const [additionalProps, setAdditionalProps] = useState<?Array<
    [string, string],
  >>([]);

  const editingMeshID = match.params.meshID;
  useEffect(() => {
    if (editingMeshID) {
      MagmaV1API.getWifiByNetworkIdMeshesByMeshIdConfig({
        networkId: nullthrows(match.params.networkId),
        meshId: meshID,
      })
        .then(configs => {
          setConfigs({...configs});
          setAdditionalProps(additionalPropsToArray(configs.additional_props));
        })
        .catch(e => {
          enqueueSnackbar(e?.response?.data?.message || e.message, {
            variant: 'error',
          });
          props.onCancel();
        });
    }
  }, [editingMeshID, enqueueSnackbar, match.params.networkId, meshID, props]);

  if (editingMeshID && Object.keys(configs).length === 0) {
    return <LoadingFillerBackdrop />;
  }

  const onSave = async () => {
    const meshWifiConfigs = {
      ...configs,
      mesh_frequency: parseInt(configs.mesh_frequency),
      additional_props: additionalPropsToObject(additionalProps) || undefined,
    };

    try {
      if (editingMeshID) {
        await MagmaV1API.putWifiByNetworkIdMeshesByMeshIdConfig({
          networkId: nullthrows(match.params.networkId),
          meshId: editingMeshID,
          meshWifiConfigs,
        });
        props.onSave(editingMeshID);
        return;
      }

      // create a mesh
      await MagmaV1API.postWifiByNetworkIdMeshes({
        networkId: nullthrows(match.params.networkId),
        wifiMesh: {
          id: meshID,
          config: meshWifiConfigs,
          name: meshID,
          gateway_ids: [],
        },
      });
      props.onSave(meshID);
    } catch (e) {
      enqueueSnackbar(e.response.data.message || e.message, {variant: 'error'});
    }
  };

  const onConfigChangeHandler = (fieldName: string) => ({target}) =>
    setConfigs({...configs, [fieldName]: target.value});

  return (
    <Dialog open={true} onClose={props.onCancel}>
      <DialogTitle>{editingMeshID ? 'Edit Mesh' : 'New Mesh'}</DialogTitle>
      <DialogContent>
        <FormGroup row>
          <TextField
            required
            className={classes.input}
            label="Mesh Name"
            margin="normal"
            onChange={({target}) => setMeshID(target.value)}
            value={meshID}
            disabled={!!editingMeshID}
          />
          <TextField
            required
            className={classes.input}
            label="SSID"
            margin="normal"
            value={configs.ssid}
            onChange={onConfigChangeHandler('ssid')}
          />
          <TextField
            className={classes.input}
            label="Password"
            margin="normal"
            value={configs.password}
            onChange={onConfigChangeHandler('password')}
          />
          <FormControlLabel
            control={
              <Checkbox
                checked={configs.xwf_enabled}
                onChange={({target}) =>
                  setConfigs({...configs, xwf_enabled: target.checked})
                }
                color="primary"
              />
            }
            label="Enable XWF"
          />
          <KeyValueFields
            key_label="key"
            value_label="value"
            keyValuePairs={additionalProps || [['', '']]}
            onChange={setAdditionalProps}
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
