/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ContextRouter, Match} from 'react-router-dom';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import PolicyRuleEditDialog from './PolicyRuleEditDialog';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import axios from 'axios';

import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {MagmaAPIUrls} from '../../common/MagmaAPI';
import {Route, withRouter} from 'react-router-dom';
import {findIndex} from 'lodash';
import {withStyles} from '@material-ui/core/styles';

const styles = {};

type Rule = {[string]: any};
type State = {
  rules: ?Array<Rule>,
};

type Props = WithStyles<typeof styles> & ContextRouter & WithAlert & {};

class PoliciesConfig extends React.Component<Props, State> {
  state = {
    rules: null,
  };

  componentDidMount() {
    this.loadData().catch(error =>
      this.props.alert(error.response?.data?.message || error),
    );
  }

  async loadData() {
    const rules = await axios.get(
      MagmaAPIUrls.networkPolicyRules(this.props.match),
    );
    const results = await axios.all(
      rules.data.map(ruleId =>
        axios.get(MagmaAPIUrls.networkPolicyRule(this.props.match, ruleId)),
      ),
    );

    this.setState({rules: results.map(r => r.data)});
  }

  render() {
    const {rules} = this.state;
    if (!rules) {
      return <LoadingFiller />;
    }

    const rows = rules.map(rule => (
      <TableRow key={rule.id}>
        <TableCell>{rule.id}</TableCell>
        <TableCell>{rule.priority}</TableCell>
        <TableCell>
          <NestedRouteLink to={`/edit/${encodeURIComponent(rule.id)}/`}>
            <IconButton>
              <EditIcon />
            </IconButton>
          </NestedRouteLink>
          <IconButton onClick={() => this.onDeleteRule(rule)}>
            <DeleteIcon />
          </IconButton>
        </TableCell>
      </TableRow>
    ));

    return (
      <>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Precedence</TableCell>
              <TableCell>
                <NestedRouteLink to="/add/">
                  <Button color="primary" variant="contained">
                    Add Rule
                  </Button>
                </NestedRouteLink>
              </TableCell>
            </TableRow>
          </TableHead>
          {<TableBody>{rows}</TableBody>}
        </Table>
        <Route
          path={`${this.props.match.path}/add`}
          component={this.renderAddDialog}
        />
        <Route
          path={`${this.props.match.path}/edit/:ruleID`}
          component={this.renderEditDialog}
        />
      </>
    );
  }

  renderAddDialog = () => (
    // $FlowFixMe: missing prop. Please fix.
    <PolicyRuleEditDialog
      onCancel={this.handleCloseDialog}
      onSave={this.onAddRule}
    />
  );

  onAddRule = rule => {
    const rules = [...nullthrows(this.state.rules), rule];
    this.setState({rules});
    this.handleCloseDialog();
  };

  renderEditDialog = ({match}: {match: Match}) => {
    const rule = nullthrows(this.state.rules).find(
      r => r.id == match.params.ruleID,
    );
    return (
      <PolicyRuleEditDialog
        // $FlowFixMe: rule is nullable. Please fix.
        rule={rule}
        onCancel={this.handleCloseDialog}
        onSave={this.onEditRule}
      />
    );
  };

  onEditRule = rule => {
    const rules = [...nullthrows(this.state.rules)];
    rules[findIndex(rules, r => r.id === rule.id)] = rule;
    this.setState({rules});
    this.handleCloseDialog();
  };

  onDeleteRule = async (rule: Rule) => {
    const confirmed = await this.props.confirm(
      `Are you sure you want to remove the rule "${rule.id}"?`,
    );

    if (!confirmed) {
      return;
    }

    await axios.delete(
      MagmaAPIUrls.networkPolicyRule(this.props.match, rule.id),
    );
    const rules = nullthrows(this.state.rules).slice();
    rules.splice(findIndex(rules, r => r.id === rule.id), 1);
    this.setState({rules});
  };

  handleCloseDialog = () => {
    this.props.history.push(`${this.props.match.url}/`);
  };
}

export default withStyles(styles)(withRouter(withAlert(PoliciesConfig)));
