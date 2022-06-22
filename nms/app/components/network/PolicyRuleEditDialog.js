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

// $FlowFixMe migrated to typescript
import type {NetworkType} from '../../../shared/types/network';
import type {
  policy_qos_profile,
  policy_rule,
} from '../../../generated/MagmaAPIBindings';

import AddCircleOutline from '@material-ui/icons/AddCircleOutline';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import IconButton from '@material-ui/core/IconButton';
import InputLabel from '@material-ui/core/InputLabel';
import MagmaV1API from '../../../generated/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import PolicyFlowFields from './PolicyFlowFields';
import React from 'react';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import TypedSelect from '../TypedSelect';
import Typography from '@material-ui/core/Typography';

// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';
// $FlowFixMe migrated to typescript
import {ACTION, DIRECTION, PROTOCOL} from './PolicyTypes';
// $FlowFixMe migrated to typescript
import {CWF, FEG, LTE} from '../../../shared/types/network';
// $FlowFixMe[cannot-resolve-module]
import {base64ToHex, decodeBase64} from '../../util/strings';
// $FlowFixMe migrated to typescript
import {coalesceNetworkType} from '../../../shared/types/network';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';
import {useParams} from 'react-router-dom';

const useStyles = makeStyles(() => ({
  input: {width: '100%'},
}));

type Props = {
  onCancel: () => void,
  onSave: string => void,
  qosProfiles: {[string]: policy_qos_profile},
  rule?: policy_rule,
  mirrorNetwork?: string,
};

export default function PolicyRuleEditDialog(props: Props) {
  const classes = useStyles();
  const {networkId} = useParams();
  const {qosProfiles, mirrorNetwork} = props;
  const [networkType, setNetworkType] = useState<?NetworkType>(null);
  const [mirrorNetworkType, setMirrorNetworkType] = useState<?NetworkType>(
    null,
  );
  const [networkWideRuleIDs, setNetworkWideRuldIDs] = useState(null);
  const [isNetworkWide, setIsNetworkWide] = useState<boolean>(false);

  const [rule, setRule] = useState(
    props.rule || {
      qos_profile: undefined,
      id: '',
      priority: 1,
      flow_list: [],
      rating_group: 0,
      monitoring_key: '',
      app_name: undefined,
      app_service_type: undefined,
      assigned_subscribers: undefined,
    },
  );

  // Grab the network type for the network, and the mirrorNetwork if it exists.
  useEffect(() => {
    getNetworkType(networkId, setNetworkType);
    getNetworkType(mirrorNetwork, setMirrorNetworkType);
  }, [networkId, mirrorNetwork]);

  // Grab the network wide rule IDs from the network's subscriber config
  // Case on the network type to determine which endpoint to call.
  useEffect(() => {
    switch (networkType) {
      case LTE:
        MagmaV1API.getLteByNetworkIdSubscriberConfigRuleNames({
          networkId: networkId,
        }).then(ruleIDs => setNetworkWideRuldIDs(ruleIDs));
        break;
      case CWF:
        MagmaV1API.getCwfByNetworkIdSubscriberConfigRuleNames({
          networkId: networkId,
        }).then(ruleIDs => {
          setNetworkWideRuldIDs(ruleIDs);
        });
        break;
      case FEG:
        MagmaV1API.getFegByNetworkIdSubscriberConfigRuleNames({
          networkId: networkId,
        }).then(ruleIDs => setNetworkWideRuldIDs(ruleIDs));
        break;
    }
  }, [networkId, networkType]);

  // The rule is network-wide if the rule's ID exists in network-wide rule IDs
  useEffect(() => {
    networkWideRuleIDs
      ? setIsNetworkWide(networkWideRuleIDs.includes(props.rule?.id))
      : false;
  }, [networkWideRuleIDs, props.rule]);

  const handleAddFlow = () => {
    const flowList = [
      ...(rule.flow_list || []),
      {
        action: ACTION.DENY,
        match: {
          direction: DIRECTION.UPLINK,
          ip_proto: PROTOCOL.IPPROTO_IP,
        },
      },
    ];

    setRule({...rule, flow_list: flowList});
  };

  const onFlowChange = (index, flow) => {
    const flowList = [...(rule.flow_list || [])];
    flowList[index] = flow;
    setRule({...rule, flow_list: flowList});
  };

  const handleDeleteFlow = (index: number) => {
    const flowList = [...(rule.flow_list || [])];
    flowList.splice(index, 1);
    setRule({...rule, flow_list: flowList});
  };

  const onSave = async () => {
    const ruleData = [
      {
        networkId: nullthrows(networkId),
        ruleId: rule.id,
        policyRule: rule,
      },
    ];

    if (mirrorNetwork != null) {
      ruleData.push({
        networkId: mirrorNetwork,
        ruleId: rule.id,
        policyRule: rule,
      });
    }

    if (props.rule) {
      await Promise.all(
        ruleData.map(d =>
          MagmaV1API.putNetworksByNetworkIdPoliciesRulesByRuleId(d),
        ),
      );
    } else {
      await Promise.all(
        ruleData.map(d => MagmaV1API.postNetworksByNetworkIdPoliciesRules(d)),
      );
    }

    const networkWideRuleData = {
      networkId: nullthrows(networkId),
      ruleId: rule.id,
    };
    if (isNetworkWide) {
      await postNetworkWideRuleID(networkWideRuleData, networkType);
      if (mirrorNetwork != null) {
        networkWideRuleData.networkId = mirrorNetwork;
        await postNetworkWideRuleID(networkWideRuleData, mirrorNetworkType);
      }
    } else {
      await deleteNetworkWideRuleID(networkWideRuleData, networkType);
      if (mirrorNetwork != null) {
        networkWideRuleData.networkId = mirrorNetwork;
        await deleteNetworkWideRuleID(networkWideRuleData, mirrorNetworkType);
      }
    }
  };

  return (
    <Dialog open={true} onClose={props.onCancel} scroll="body">
      <DialogTitle>{props.rule ? 'Edit' : 'Add'} Rule</DialogTitle>
      <DialogContent>
        <TextField
          required
          className={classes.input}
          label="ID"
          margin="normal"
          disabled={!!props.rule}
          value={rule.id}
          onChange={({target}) => setRule({...rule, id: target.value})}
        />
        <TextField
          required
          className={classes.input}
          label="Precendence"
          margin="normal"
          value={rule.priority}
          onChange={({target}) =>
            setRule({...rule, priority: parseInt(target.value)})
          }
        />
        <TextField
          required
          className={classes.input}
          label="Monitoring Key (base64)"
          margin="normal"
          value={rule.monitoring_key ?? ''}
          onChange={({target}) =>
            setRule({...rule, monitoring_key: target.value})
          }
        />
        <TextField
          required
          className={classes.input}
          label="Monitoring Key (hex)"
          margin="normal"
          disabled={true}
          value={base64ToHex(rule.monitoring_key ?? '')}
        />
        <TextField
          required
          className={classes.input}
          label="Monitoring Key (ascii)"
          margin="normal"
          disabled={true}
          value={decodeBase64(rule.monitoring_key ?? '')}
        />
        <TextField
          required
          className={classes.input}
          label="Rating Group"
          margin="normal"
          value={rule.rating_group}
          type="number"
          onChange={({target}) =>
            setRule({
              ...rule,
              rating_group: parseInt(target.value),
            })
          }
        />
        <FormControl className={classes.input}>
          <InputLabel htmlFor="trackingType">Tracking Type</InputLabel>
          <TypedSelect
            items={{
              ONLY_OCS: 'Only OCS',
              ONLY_PCRF: 'Only PCRF',
              OCS_AND_PCRF: 'OCS and PCRF',
              NO_TRACKING: 'No Tracking',
            }}
            inputProps={{id: 'trackingType'}}
            value={rule.tracking_type || 'NO_TRACKING'}
            onChange={trackingType =>
              setRule({...rule, tracking_type: trackingType})
            }
          />
        </FormControl>
        <FormControl className={classes.input}>
          <InputLabel htmlFor="appName">App Name</InputLabel>
          <TypedSelect
            items={{
              NO_APP_NAME: 'No App Name',
              FACEBOOK: 'Facebook',
              FACEBOOK_MESSENGER: 'Facebook Messenger',
              INSTAGRAM: 'Instagram',
              YOUTUBE: 'Youtube',
              GOOGLE: 'Google',
              GMAIL: 'Gmail',
              GOOGLE_DOCS: 'Google Docs',
              NETFLIX: 'Netflix',
              APPLE: 'Apple',
              MICROSOFT: 'Microsoft',
              REDDIT: 'Reddit',
              WHATSAPP: 'WhatsApp',
              GOOGLE_PLAY: 'Google Play',
              APPSTORE: 'App Store',
              AMAZON: 'Amazon',
              WECHAT: 'Wechat',
              TIKTOK: 'TikTok',
              TWITTER: 'Twitter',
              WIKIPEDIA: 'Wikipedia',
              GOOGLE_MAPS: 'Google Maps',
              YAHOO: 'Yahoo',
              IMO: 'IMO',
            }}
            inputProps={{id: 'appName'}}
            value={rule.app_name || 'NO_APP_NAME'}
            onChange={appName => setRule({...rule, app_name: appName})}
          />
        </FormControl>
        <FormControl className={classes.input}>
          <InputLabel htmlFor="appServiceType">App Service Type</InputLabel>
          <TypedSelect
            items={{
              NO_SERVICE_TYPE: 'No Service Type',
              CHAT: 'Chat',
              AUDIO: 'Audio',
              VIDEO: 'Video',
            }}
            inputProps={{id: 'appServiceType'}}
            value={rule.app_service_type || 'NO_SERVICE_TYPE'}
            onChange={appServiceType =>
              setRule({...rule, app_service_type: appServiceType})
            }
          />
        </FormControl>
        <FormControl className={classes.input}>
          <InputLabel htmlFor="target">Network Wide</InputLabel>
          <TypedSelect
            items={{
              true: 'true',
              false: 'false',
            }}
            inputProps={{id: 'target'}}
            value={isNetworkWide ? 'true' : 'false'}
            onChange={target => {
              setIsNetworkWide(target === 'true');
            }}
          />
        </FormControl>
        <FormControl className={classes.input}>
          <InputLabel htmlFor="target">Qos Profile</InputLabel>
          <Select
            className={classes.input}
            value={rule?.qos_profile ?? ''}
            onChange={({target}) =>
              setRule({...rule, qos_profile: target.value})
            }>
            {Object.keys(qosProfiles).map(profileID => (
              <MenuItem key={profileID} value={profileID}>
                {profileID}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        <Typography variant="h6">
          Flows
          <IconButton onClick={handleAddFlow}>
            <AddCircleOutline />
          </IconButton>
        </Typography>
        {(rule.flow_list || []).slice(0, 30).map((flow, i) => (
          <PolicyFlowFields
            key={i}
            index={i}
            flow={flow}
            handleDelete={handleDeleteFlow}
            onChange={onFlowChange}
          />
        ))}
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

async function postNetworkWideRuleID(networkWideRuleData, networkType) {
  switch (networkType) {
    case LTE:
      MagmaV1API.postLteByNetworkIdSubscriberConfigRuleNamesByRuleId(
        networkWideRuleData,
      );
      break;
    case CWF:
      MagmaV1API.postCwfByNetworkIdSubscriberConfigRuleNamesByRuleId(
        networkWideRuleData,
      );
      break;
    case FEG:
      MagmaV1API.postFegByNetworkIdSubscriberConfigRuleNamesByRuleId(
        networkWideRuleData,
      );
      break;
  }
}

async function deleteNetworkWideRuleID(networkWideRuleData, networkType) {
  switch (networkType) {
    case LTE:
      MagmaV1API.deleteLteByNetworkIdSubscriberConfigRuleNamesByRuleId(
        networkWideRuleData,
      );
      break;
    case CWF:
      MagmaV1API.deleteCwfByNetworkIdSubscriberConfigRuleNamesByRuleId(
        networkWideRuleData,
      );
      break;
    case FEG:
      MagmaV1API.deleteFegByNetworkIdSubscriberConfigRuleNamesByRuleId(
        networkWideRuleData,
      );
      break;
  }
}

function getNetworkType(
  networkId: ?string,
  setNetworkType: (?NetworkType) => void,
) {
  if (networkId != null) {
    MagmaV1API.getNetworksByNetworkIdType({networkId}).then(networkType =>
      setNetworkType(coalesceNetworkType(networkId, networkType)),
    );
  }
}
