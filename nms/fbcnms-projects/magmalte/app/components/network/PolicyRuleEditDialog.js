/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ContextRouter} from 'react-router-dom';
import type {PolicyRule} from './PolicyTypes';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import AddCircleOutline from '@material-ui/icons/AddCircleOutline';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import IconButton from '@material-ui/core/IconButton';
import PolicyFlowFields from './PolicyFlowFields';
import React from 'react';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';

import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {ACTION, DIRECTION, PROTOCOL} from './PolicyTypes';
import {MagmaAPIUrls} from '../../common/MagmaAPI';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  input: {width: '100%'},
};

type Props = ContextRouter &
  WithStyles<typeof styles> &
  WithAlert & {
    onCancel: () => void,
    onSave: PolicyRule => void,
    rule: PolicyRule,
  };

type State = {
  rule: PolicyRule,
};

class PolicyRuleEditDialog extends React.Component<Props, State> {
  state = {
    rule: this.props.rule || {
      id: '',
      priority: 1,
      flow_list: [],
    },
  };

  render() {
    const {rule} = this.state;
    return (
      <Dialog open={true} onClose={this.props.onCancel} scroll="body">
        <DialogTitle>{this.props.rule ? 'Edit' : 'Add'} Rule</DialogTitle>
        <DialogContent>
          <TextField
            required
            className={this.props.classes.input}
            label="ID"
            margin="normal"
            disabled={!!this.props.rule}
            value={rule.id}
            onChange={this.handleIDChange}
          />
          <TextField
            required
            className={this.props.classes.input}
            label="Precendence"
            margin="normal"
            value={rule.priority}
            onChange={this.handlePriorityChange}
          />
          <Typography variant="h6">
            Flows
            <IconButton onClick={this.handleAddFlow}>
              <AddCircleOutline />
            </IconButton>
          </Typography>
          {(rule.flow_list || []).slice(0, 30).map((flow, i) => (
            <PolicyFlowFields
              key={i}
              index={i}
              flow={flow}
              handleActionChange={this.handleActionChange}
              handleFieldChange={this.handleFieldChange}
              handleDelete={this.handleDeleteFlow}
            />
          ))}
        </DialogContent>
        <DialogActions>
          <Button onClick={this.props.onCancel} color="primary">
            Cancel
          </Button>
          <Button onClick={this.onSave} color="primary" variant="contained">
            Save
          </Button>
        </DialogActions>
      </Dialog>
    );
  }

  onSave = async () => {
    if (this.props.rule) {
      await axios.put(
        MagmaAPIUrls.networkPolicyRule(this.props.match, this.state.rule.id),
        this.state.rule,
      );
    } else {
      await axios.post(
        MagmaAPIUrls.networkPolicyRules(this.props.match),
        this.state.rule,
      );
    }

    this.props.onSave(this.state.rule);
  };

  handleIDChange = ({target}) =>
    this.setState({rule: {...this.state.rule, id: target.value}});

  handlePriorityChange = ({target}) =>
    this.setState({
      rule: {...this.state.rule, priority: parseInt(target.value)},
    });

  handleAddFlow = () => {
    const flowList = (this.state.rule.flow_list || []).slice();
    flowList.push({
      action: ACTION.DENY,
      match: {
        direction: DIRECTION.UPLINK,
        ip_proto: PROTOCOL.IPPROTO_IP,
      },
    });

    this.setState({
      rule: {
        ...this.state.rule,
        flow_list: flowList,
      },
    });
  };

  handleActionChange = (index: number, action: string) => {
    const flowList = [...nullthrows(this.state.rule.flow_list)];
    flowList[index] = {...flowList[index], action};

    this.setState({
      rule: {
        ...this.state.rule,
        flow_list: flowList,
      },
    });
  };

  handleFieldChange = (
    index: number,
    field: string,
    value: string | number,
  ) => {
    const flowList = nullthrows(this.state.rule.flow_list).slice();
    flowList[index] = {
      ...flowList[index],
      match: {...flowList[index].match, [field]: value},
    };

    this.setState({
      rule: {
        ...this.state.rule,
        flow_list: flowList,
      },
    });
  };

  handleDeleteFlow = (index: number) => {
    const flowList = nullthrows(this.state.rule.flow_list).slice();
    flowList.splice(index, 1);

    this.setState({
      rule: {
        ...this.state.rule,
        flow_list: flowList,
      },
    });
  };
}

export default withStyles(styles)(withRouter(withAlert(PolicyRuleEditDialog)));
