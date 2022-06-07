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
import DialogTitle from '../../theme/design-system/DialogTitle';
import FormLabel from '@material-ui/core/FormLabel';
import List from '@material-ui/core/List';
import LteNetworkContext from '../../components/context/LteNetworkContext';
import PolicyAppEdit from './PolicyApp';
import PolicyContext from '../../components/context/PolicyContext';
import PolicyFlowsEdit from './PolicyFlows';
import PolicyHeaderEnrichmentEdit from './PolicyHeaderEnrichment';
import PolicyInfoEdit from './PolicyInfo';
import PolicyRedirectEdit from './PolicyRedirect';
import PolicyTrackingEdit from './PolicyTracking';
import React, {SetStateAction} from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import {AltFormField} from '../../components/FormField';
import {PolicyRule} from '../../../generated-ts';
import {colors} from '../../theme/default';
import {getErrorMessage} from '../../util/ErrorUtils';
import {makeStyles} from '@material-ui/styles';
import {useContext, useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';

const DEFAULT_POLICY_RULE = {
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
  header_enrichment_targets: undefined,
};
const useStyles = makeStyles(() => ({
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    color: colors.primary.white,
  },
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '50%',
    fullWidth: true,
  },
}));

type Props = {
  open: boolean;
  onClose: () => void;
  rule?: PolicyRule;
};

export default function PolicyRuleEditDialog(props: Props) {
  const classes = useStyles();
  const [tabPos, setTabPos] = React.useState(0);
  const lteNetworkCtx = useContext(LteNetworkContext);
  const lteNetwork = lteNetworkCtx.state;
  const ctx = useContext(PolicyContext);
  const qosProfiles = ctx.qosProfiles;
  const [isNetworkWide, setIsNetworkWide] = useState<boolean>(false);
  const enqueueSnackbar = useEnqueueSnackbar();
  const [error, setError] = useState('');
  const [rule, setRule] = useState(props.rule || DEFAULT_POLICY_RULE);

  useEffect(() => {
    setRule(props.rule || DEFAULT_POLICY_RULE);
    setError('');
    setTabPos(0);
  }, [props.open, props.rule]);

  useEffect(() => {
    if (props.rule?.id) {
      setIsNetworkWide(
        lteNetwork?.subscriber_config?.network_wide_rule_names?.includes(
          props.rule?.id,
        ) || false,
      );
    }
  }, [props.rule, lteNetwork]);

  const tabList = [
    'Policy',
    'Flows',
    'Tracking',
    'Redirect',
    'Header Enrichment',
  ];

  if (lteNetwork?.cellular?.epc?.network_services?.includes('dpi')) {
    tabList.push('App');
  }

  const isAdd = !props.rule;

  const onSave = async () => {
    try {
      if (isAdd) {
        // if we are trying to save first tab(containing ID information)
        // check if the policy exists so that we don't end up modifying
        // existing policies
        if (rule.id === '') {
          setError('empty rule id');
          return;
        }

        if (tabPos === 0 && rule.id in ctx.state) {
          setError(`Policy ${rule.id} already exists`);
          return;
        }
      }

      await ctx.setState(rule.id, rule, isNetworkWide);
      enqueueSnackbar('Policy saved successfully', {
        variant: 'success',
      });

      if (props.rule) {
        props.onClose();
      } else {
        if (tabPos < tabList.length - 1) {
          setTabPos(tabPos + 1);
        } else {
          props.onClose();
        }
      }
    } catch (error) {
      setError(getErrorMessage(error));
    }
  };

  const onClose = () => {
    props.onClose();
  };

  const editProps = {
    policyRule: rule,
    onChange: (policyRule: PolicyRule) => {
      setRule(policyRule);
    },
    isNetworkWide,
    setIsNetworkWide,
    qosProfiles: qosProfiles,
    inputClass: classes.input,
  };

  return (
    <Dialog
      data-testid="editDialog"
      open={props.open}
      scroll="body"
      fullWidth={true}
      maxWidth={'md'}>
      <DialogTitle
        onClose={onClose}
        label={props.rule ? 'Edit Policy' : 'Add New Policy'}
      />
      <Tabs
        value={tabPos}
        onChange={(_, v: SetStateAction<number>) => setTabPos(v)}
        indicatorColor="primary"
        variant="scrollable"
        scrollButtons="auto"
        className={classes.tabBar}>
        {tabList.map(tabKey => (
          <Tab key={tabKey} data-testid={tabKey + 'Tab'} label={tabKey} />
        ))}
      </Tabs>
      <DialogContent>
        <List>
          {error !== '' && (
            <AltFormField disableGutters label={''}>
              <FormLabel data-testid="configEditError" error>
                {error}
              </FormLabel>
            </AltFormField>
          )}

          {tabPos === 0 && <PolicyInfoEdit {...editProps} />}
          {tabPos === 1 && <PolicyFlowsEdit {...editProps} />}
          {tabPos === 2 && <PolicyTrackingEdit {...editProps} />}
          {tabPos === 3 && <PolicyRedirectEdit {...editProps} />}
          {tabPos === 4 && <PolicyHeaderEnrichmentEdit {...editProps} />}
          {tabPos === 5 && <PolicyAppEdit {...editProps} />}
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>Close</Button>
        <Button
          variant="contained"
          color="primary"
          onClick={() => void onSave()}>
          {props.rule ? 'Save' : 'Save And Continue'}
        </Button>
      </DialogActions>
    </Dialog>
  );
}
