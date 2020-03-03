/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {policy_rule} from '@fbcnms/magma-api';

import AddCircleOutline from '@material-ui/icons/AddCircleOutline';
import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import IconButton from '@material-ui/core/IconButton';
import InputLabel from '@material-ui/core/InputLabel';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import PolicyFlowFields from './PolicyFlowFields';
import React from 'react';
import TextField from '@material-ui/core/TextField';
import TypedSelect from '@fbcnms/ui/components/TypedSelect';
import Typography from '@material-ui/core/Typography';

import nullthrows from '@fbcnms/util/nullthrows';
import {ACTION, DIRECTION, PROTOCOL} from './PolicyTypes';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  input: {width: '100%'},
}));

type Props = {
  onCancel: () => void,
  onSave: string => void,
  rule?: policy_rule,
  mirrorNetwork?: string,
};

export default function PolicyRuleEditDialog(props: Props) {
  const classes = useStyles();
  const {match} = useRouter();
  const [rule, setRule] = useState(
    props.rule || {
      id: '',
      priority: 1,
      flow_list: [],
      rating_group: 0,
      monitoring_key: '',
    },
  );

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
    const data = [
      {
        networkId: nullthrows(match.params.networkId),
        ruleId: rule.id,
        policyRule: rule,
      },
    ];

    if (props.mirrorNetwork != null) {
      data.push({
        networkId: props.mirrorNetwork,
        ruleId: rule.id,
        policyRule: rule,
      });
    }

    if (props.rule) {
      await Promise.all(
        data.map(d =>
          MagmaV1API.putNetworksByNetworkIdPoliciesRulesByRuleId(d),
        ),
      );
    } else {
      await Promise.all(
        data.map(d => MagmaV1API.postNetworksByNetworkIdPoliciesRules(d)),
      );
    }

    props.onSave(rule.id);
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
          label="Monitoring Key"
          margin="normal"
          value={rule.monitoring_key}
          onChange={({target}) =>
            setRule({...rule, monitoring_key: target.value})
          }
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
        <Button onClick={props.onCancel} skin="regular">
          Cancel
        </Button>
        <Button onClick={onSave}>Save</Button>
      </DialogActions>
    </Dialog>
  );
}
