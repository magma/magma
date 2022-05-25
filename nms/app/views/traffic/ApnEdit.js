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
import type {
  aggregated_maximum_bitrate,
  apn,
  apn_configuration,
  qos_profile,
} from '../../../generated/MagmaAPIBindings';

// $FlowFixMe migrated to typescript
import ApnContext from '../../components/context/ApnContext';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
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
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../theme/design-system/Text';

import {AltFormField, AltFormFieldSubheading} from '../../components/FormField';
import {useContext, useEffect, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';

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
  open: boolean,
  onClose: () => void,
  apn?: apn,
};

export default function ApnEditDialog(props: DialogProps) {
  const [_, setError] = useState('');
  const [apn, setApn] = useState<apn>(props.apn || DEFAULT_APN_CONFIG);
  const isAdd = props.apn ? false : true;

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
        onSave={(apn: apn) => {
          setApn(apn);
          onClose();
        }}
      />
    </Dialog>
  );
}

type Props = {
  isAdd: boolean,
  apn: apn,
  apnConfig?: apn_configuration,
  onClose: () => void,
  onSave: apn => void,
};

const PDNTypeEnum = Object.freeze({
  IPv4: 0,
  IPv6: 1,
  IPv4v6: 2,
  IPv4orv6: 3,
});

export function ApnEdit(props: Props) {
  const [error, setError] = useState('');
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(ApnContext);
  const [apn, setApn] = useState<apn>(props.apn || DEFAULT_APN_CONFIG);
  const [maxBandwidth, setMaxBandwidth] = useState<aggregated_maximum_bitrate>(
    props.apn?.apn_configuration?.ambr ||
      DEFAULT_APN_CONFIG.apn_configuration.ambr,
  );
  const [qosProfile, setQosProfile] = useState<qos_profile>(
    props.apn?.apn_configuration?.qos_profile ||
      DEFAULT_APN_CONFIG.apn_configuration.qos_profile,
  );

  const [pdnType, setPdnType] = useState<$Values<typeof PDNTypeEnum>>(
    props.apn?.apn_configuration?.pdn_type ||
      DEFAULT_APN_CONFIG.apn_configuration.pdn_type,
  );
  const onSave = async () => {
    if (apn.apn_name === '') {
      throw Error('Invalid Name');
    }
    if (props.isAdd && apn.apn_name in ctx.state) {
      setError(`APN ${apn.apn_name} already exists`);
      return;
    }
    try {
      const newApn: apn = {
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
    } catch (e) {
      setError(e.response?.data?.message ?? e?.message);
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
                min={0}
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
                min={0}
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
                  min={0}
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
                  min={0}
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
                  // $FlowIgnore: value guaranteed to match the number literals
                  setPdnType(parseInt(target.value));
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
        <Button onClick={onSave} variant="contained" color="primary">
          Save
        </Button>
      </DialogActions>
    </>
  );
}
