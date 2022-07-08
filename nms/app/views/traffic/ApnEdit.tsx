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

import ApnContext from '../../components/context/ApnContext';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '../../theme/design-system/DialogTitle';
import FormLabel from '@material-ui/core/FormLabel';
import InputAdornment from '@material-ui/core/InputAdornment';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Select from '@material-ui/core/Select';
import Switch from '@material-ui/core/Switch';
import Text from '../../theme/design-system/Text';
import {AltFormField, AltFormFieldSubheading} from '../../components/FormField';
import {getErrorMessage} from '../../util/ErrorUtils';
import {useContext, useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import type {Apn, ApnConfiguration} from '../../../generated';

const DEFAULT_APN_CONFIG = {
  apn_configuration: {
    ambr: {
      max_bandwidth_dl: 1000000,
      max_bandwidth_ul: 1000000,
    },
    qos_profile: {
      class_id: 9,
      preemption_capability: false,
      preemption_vulnerability: false,
      priority_level: 15,
    },
    pdn_type: 0,
  },
  apn_name: '',
};

type DialogProps = {
  open: boolean;
  onClose: () => void;
  apn?: Apn;
};

export default function ApnEditDialog(props: DialogProps) {
  const [, setError] = useState('');
  const [apn, setApn] = useState<Apn>(props.apn || DEFAULT_APN_CONFIG);
  const isAdd = !props.apn;

  useEffect(() => {
    setApn(props.apn || DEFAULT_APN_CONFIG);
    setError('');
  }, [props.open, props.apn]);

  const onClose = () => {
    props.onClose();
  };

  return (
    <Dialog
      data-testid="editDialog"
      open={props.open}
      fullWidth={true}
      maxWidth="sm">
      <DialogTitle
        label={isAdd ? 'Add New APN' : 'Edit APN'}
        onClose={onClose}
      />
      <ApnEdit
        isAdd={isAdd}
        apn={apn}
        onClose={onClose}
        onSave={(apn: Apn) => {
          setApn(apn);
          onClose();
        }}
      />
    </Dialog>
  );
}

type Props = {
  isAdd: boolean;
  apn: Apn;
  apnConfig?: ApnConfiguration;
  onClose: () => void;
  onSave: (apn: Apn) => void;
};

export function ApnEdit(props: Props) {
  const [error, setError] = useState('');
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(ApnContext);
  const [apn, setApn] = useState<Apn>(props.apn || DEFAULT_APN_CONFIG);
  const [maxBandwidth, setMaxBandwidth] = useState(
    props.apn?.apn_configuration?.ambr ||
      DEFAULT_APN_CONFIG.apn_configuration.ambr,
  );
  const [qosProfile, setQosProfile] = useState(
    props.apn?.apn_configuration?.qos_profile ||
      DEFAULT_APN_CONFIG.apn_configuration.qos_profile,
  );
  const [pdnType, setPdnType] = useState(
    props.apn?.apn_configuration?.pdn_type ||
      DEFAULT_APN_CONFIG.apn_configuration.pdn_type,
  );

  const onSave = async () => {
    if (apn.apn_name === '') {
      throw Error('Invalid Name');
    }

    if (props.isAdd && apn.apn_name in ctx.state) {
      setError(`APN ${apn.apn_name} already exists`);
    }

    try {
      const newApn: Apn = {
        ...apn,
        apn_configuration: {
          ambr: maxBandwidth,
          qos_profile: qosProfile,
          pdn_type: pdnType,
        },
      };
      await ctx.setState(newApn.apn_name, newApn);
      enqueueSnackbar('APN saved successfully', {
        variant: 'success',
      });
      props.onSave(newApn);
    } catch (error) {
      getErrorMessage(error);
    }
  };

  return (
    <>
      <DialogContent data-testid="apnEditDialog">
        <List>
          {error !== '' && (
            <AltFormField label={''}>
              <FormLabel data-testid="configEditError" error>
                {error}
              </FormLabel>
            </AltFormField>
          )}
          <div>
            <ListItem dense disableGutters />
            <AltFormField label={'APN ID'}>
              <OutlinedInput
                data-testid="apnID"
                placeholder="apn_id"
                fullWidth={true}
                value={apn.apn_name}
                onChange={({target}) =>
                  setApn({...apn, apn_name: target.value})
                }
              />
            </AltFormField>
            <AltFormField label={'Class ID'}>
              <OutlinedInput
                data-testid="classID"
                placeholder="9"
                type="number"
                inputProps={{min: 0}}
                fullWidth={true}
                value={qosProfile.class_id}
                onChange={({target}) =>
                  setQosProfile({
                    ...qosProfile,
                    class_id: parseInt(target.value),
                  })
                }
              />
            </AltFormField>
            <AltFormField label={'ARP Priority Level'}>
              <OutlinedInput
                data-testid="apnPriority"
                placeholder="Value between 1 and 15"
                type="number"
                fullWidth={true}
                value={qosProfile.priority_level}
                onChange={({target}) =>
                  setQosProfile({
                    ...qosProfile,
                    priority_level: parseInt(target.value),
                  })
                }
              />
            </AltFormField>
            <AltFormField label={'Max Required Bandwidth'}>
              <AltFormFieldSubheading label={'Upload'}>
                <OutlinedInput
                  data-testid="apnBandwidthUL"
                  placeholder="1000000"
                  type="number"
                  value={maxBandwidth.max_bandwidth_ul}
                  onChange={({target}) =>
                    setMaxBandwidth({
                      ...maxBandwidth,
                      max_bandwidth_ul: parseInt(target.value),
                    })
                  }
                  endAdornment={
                    <InputAdornment position="end">
                      <Text variant="subtitle3">bps</Text>
                    </InputAdornment>
                  }
                />
              </AltFormFieldSubheading>
              <AltFormFieldSubheading label={'Download'}>
                <OutlinedInput
                  data-testid="apnBandwidthDL"
                  placeholder="1000000"
                  type="number"
                  value={maxBandwidth.max_bandwidth_dl}
                  onChange={({target}) =>
                    setMaxBandwidth({
                      ...maxBandwidth,
                      max_bandwidth_dl: parseInt(target.value),
                    })
                  }
                  endAdornment={
                    <InputAdornment position="end">
                      <Text variant="subtitle3">bps</Text>
                    </InputAdornment>
                  }
                />
              </AltFormFieldSubheading>
            </AltFormField>
            <AltFormField label={'ARP Pre-emption-Capability'}>
              <Switch
                data-testid="preemptionCapability"
                onChange={() => {
                  setQosProfile({
                    ...qosProfile,
                    preemption_capability: !qosProfile.preemption_capability,
                  });
                }}
                checked={qosProfile.preemption_capability}
              />
            </AltFormField>
            <AltFormField label={'ARP Pre-emption-Vulnerability'}>
              <Switch
                data-testid="preemptionVulnerability"
                onChange={() => {
                  setQosProfile({
                    ...qosProfile,
                    preemption_vulnerability: !qosProfile.preemption_vulnerability,
                  });
                }}
                checked={qosProfile.preemption_vulnerability}
              />
            </AltFormField>
            <AltFormField label={'PDN Type'}>
              <Select
                fullWidth={true}
                variant={'outlined'}
                value={pdnType || 0}
                onChange={({target}) => {
                  setPdnType(parseInt(target.value as string));
                }}
                input={<OutlinedInput data-testId="pdnType" />}>
                <MenuItem value={0}>
                  <ListItemText primary={'IPv4'} />
                </MenuItem>
                <MenuItem value={1}>
                  <ListItemText primary={'IPv6'} />
                </MenuItem>
                <MenuItem value={2}>
                  <ListItemText primary={'IPv4v6'} />
                </MenuItem>
                <MenuItem value={3}>
                  <ListItemText primary={'IPv4 or v6'} />
                </MenuItem>
              </Select>
            </AltFormField>
          </div>
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>Cancel</Button>
        <Button
          onClick={() => void onSave()}
          variant="contained"
          color="primary">
          Save
        </Button>
      </DialogActions>
    </>
  );
}
