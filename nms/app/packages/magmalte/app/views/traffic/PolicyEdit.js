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

import type {NetworkType} from '@fbcnms/types/network';
import type {
  policy_qos_profile,
  policy_rule,
  subscriber,
} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogTitle from '@material-ui/core/DialogTitle';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import PolicyAppEdit from './PolicyApp';
import PolicyFlowsEdit from './PolicyFlows';
import PolicyInfoEdit from './PolicyInfo';
import PolicyQosEdit from './PolicyQos';
import PolicyRedirectEdit from './PolicyRedirect';
import PolicySubscribersEdit from './PolicySubscribers';
import PolicyTrackingEdit from './PolicyTracking';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import grey from '@material-ui/core/colors/grey';

import {CWF, FEG, LTE} from '@fbcnms/types/network';
import {coalesceNetworkType} from '@fbcnms/types/network';
import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(() => ({
  input: {width: '100%'},
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    color: colors.primary.white,
  },
  dialog: {
    height: '540px',
  },
  inputField: {
    width: '360px',
  },
  description: {
    color: grey.A700,
    fontWeight: '400',
  },
}));

type Props = {
  onCancel: () => void,
  onSave: (policy_rule, boolean) => Promise<void>,
  subscribers: {[string]: subscriber},
  qosProfiles: {[string]: policy_qos_profile},
  rule?: policy_rule,
};

export default function PolicyRuleEditDialog(props: Props) {
  const classes = useStyles();
  const {match} = useRouter();
  const {networkId} = match.params;
  const {qosProfiles} = props;
  const [networkType, setNetworkType] = useState<?NetworkType>(null);
  const [networkWideRuleIDs, setNetworkWideRuldIDs] = useState(null);
  const [isNetworkWide, setIsNetworkWide] = useState<boolean>(false);

  const [rule, setRule] = useState(
    props.rule || {
      qos_profile: undefined,
      id: '',
      priority: 1,
      flow_list: [],
      rating_group: 0,
      redirect_information: {},
      monitoring_key: '',
      app_name: undefined,
      app_service_type: undefined,
      assigned_subscribers: undefined,
    },
  );

  const [tabPos, setTabPos] = React.useState(0);

  // Grab the network type for the network, and the mirrorNetwork if it exists.
  useEffect(() => {
    getNetworkType(networkId, setNetworkType);
  }, [networkId]);

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

  const onSave = async () => {
    await props.onSave(rule, isNetworkWide);
  };

  return (
    <Dialog open={true} onClose={props.onCancel} scroll="body" maxWidth={'md'}>
      <DialogTitle>{props.rule ? 'Edit' : 'Add New'} Policy</DialogTitle>
      <Tabs
        value={tabPos}
        onChange={(_, v) => setTabPos(v)}
        indicatorColor="primary"
        variant="scrollable"
        scrollButtons="auto"
        className={classes.tabBar}>
        <Tab key="policy" data-testid="epcTab" label={'Policy'} />
        <Tab key="flows" data-testid="ranTab" label={'Flows'} />
        <Tab key="tracking" data-testid="trackingTab" label={'Tracking'} />
        <Tab key="redirect" data-testid="redirectTab" label={'Redirect'} />
        <Tab key="app" data-testid="appTab" label={'App'} />
        <Tab key="qos" data-testid="qosTab" label={'QoS'} />
        <Tab
          key="subscribers"
          data-testid="subsribersTab"
          label={'Subscribers'}
        />
      </Tabs>
      {tabPos === 0 && (
        <PolicyInfoEdit
          policyRule={rule}
          onChange={(policyRule: policy_rule) => setRule(policyRule)}
          isNetworkWide={isNetworkWide}
          setIsNetworkWide={setIsNetworkWide}
          descriptionClass={classes.description}
          dialogClass={classes.dialog}
          inputClass={classes.inputField}
        />
      )}
      {tabPos === 1 && (
        <PolicyFlowsEdit
          policyRule={rule}
          onChange={(policyRule: policy_rule) => setRule(policyRule)}
          descriptionClass={classes.description}
          dialogClass={classes.dialog}
          inputClass={classes.inputField}
        />
      )}
      {tabPos === 2 && (
        <PolicyTrackingEdit
          policyRule={rule}
          onChange={(policyRule: policy_rule) => setRule(policyRule)}
          descriptionClass={classes.description}
          dialogClass={classes.dialog}
          inputClass={classes.inputField}
        />
      )}
      {tabPos === 3 && (
        <PolicyRedirectEdit
          policyRule={rule}
          onChange={(policyRule: policy_rule) => setRule(policyRule)}
          descriptionClass={classes.description}
          dialogClass={classes.dialog}
          inputClass={classes.inputField}
        />
      )}
      {tabPos === 4 && (
        <PolicyAppEdit
          policyRule={rule}
          onChange={(policyRule: policy_rule) => setRule(policyRule)}
          descriptionClass={classes.description}
          dialogClass={classes.dialog}
          inputClass={classes.inputField}
        />
      )}
      {tabPos === 5 && (
        <PolicyQosEdit
          policyRule={rule}
          onChange={(policyRule: policy_rule) => setRule(policyRule)}
          qosProfiles={qosProfiles}
          descriptionClass={classes.description}
          dialogClass={classes.dialog}
          inputClass={classes.inputField}
        />
      )}
      {tabPos === 6 && (
        <PolicySubscribersEdit
          policyRule={rule}
          onChange={(policyRule: policy_rule) => setRule(policyRule)}
          subscribers={props.subscribers}
          descriptionClass={classes.description}
          dialogClass={classes.dialog}
          inputClass={classes.inputField}
        />
      )}
      <DialogActions>
        <Button onClick={props.onCancel} skin="regular">
          Cancel
        </Button>
        <Button onClick={onSave}>Save</Button>
      </DialogActions>
    </Dialog>
  );
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
