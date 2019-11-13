/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';

import Button from '@fbcnms/ui/components/design-system/Button';
import CircularProgress from '@material-ui/core/CircularProgress';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import LoadingFillerBackdrop from '@fbcnms/ui/components/LoadingFillerBackdrop';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import PolicyRuleEditDialog from './PolicyRuleEditDialog';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';

import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '../../common/useMagmaAPI';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {Route} from 'react-router-dom';
import {findIndex} from 'lodash';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

function PoliciesConfig() {
  const {match, relativePath, relativeUrl, history} = useRouter();
  const [ruleIDs, setRuleIDs] = useState();

  const networkID = nullthrows(match.params.networkId);
  useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPoliciesRules,
    {networkId: networkID},
    setRuleIDs,
  );

  if (!ruleIDs) {
    return <LoadingFiller />;
  }

  const onDelete = id => {
    const newRuleIDs = [...nullthrows(ruleIDs)];
    newRuleIDs.splice(findIndex(newRuleIDs, id2 => id2 === id), 1);
    setRuleIDs(newRuleIDs);
  };

  return (
    <>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>ID</TableCell>
            <TableCell>Precedence</TableCell>
            <TableCell>
              <NestedRouteLink to="/add/">
                <Button>Add Rule</Button>
              </NestedRouteLink>
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {ruleIDs.map(id => (
            <RuleRow key={id} ruleID={id} onDelete={onDelete} />
          ))}
        </TableBody>
      </Table>
      <Route
        path={relativePath('/add')}
        component={() => (
          <PolicyRuleEditDialog
            onCancel={() => history.push(relativeUrl(''))}
            onSave={ruleID => {
              setRuleIDs([...nullthrows(ruleIDs), ruleID]);
              history.push(relativeUrl(''));
            }}
          />
        )}
      />
    </>
  );
}

type Props = WithAlert & {ruleID: string, onDelete: () => void};

const RuleRow = withAlert((props: Props) => {
  const [lastRefreshTime, setLastRefreshTime] = useState(new Date().getTime());
  const {match, relativePath, relativeUrl, history} = useRouter();
  const networkID = nullthrows(match.params.networkId);

  const {response: rule} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPoliciesRulesByRuleId,
    {networkId: networkID, ruleId: props.ruleID},
    undefined,
    lastRefreshTime,
  );

  const onDeleteRule = async () => {
    const confirmed = await props.confirm(
      `Are you sure you want to remove the rule "${props.ruleID}"?`,
    );

    if (!confirmed) {
      return;
    }

    await MagmaV1API.deleteNetworksByNetworkIdPoliciesRulesByRuleId({
      networkId: networkID,
      ruleId: props.ruleID,
    });

    props.onDelete();
  };

  const editPath = `/edit/${encodeURIComponent(props.ruleID)}/`;
  return (
    <TableRow>
      <TableCell>{props.ruleID}</TableCell>
      <TableCell>
        {rule ? rule.priority : <CircularProgress size={20} />}
      </TableCell>
      <TableCell>
        <NestedRouteLink to={editPath}>
          <IconButton>
            <EditIcon />
          </IconButton>
        </NestedRouteLink>
        <IconButton onClick={onDeleteRule}>
          <DeleteIcon />
        </IconButton>
        <Route
          path={relativePath(editPath)}
          component={() =>
            rule ? (
              <PolicyRuleEditDialog
                rule={rule}
                onCancel={() => history.push(relativeUrl(''))}
                onSave={() => {
                  setLastRefreshTime(new Date().getTime());
                  history.push(relativeUrl(''));
                }}
              />
            ) : (
              <LoadingFillerBackdrop />
            )
          }
        />
      </TableCell>
    </TableRow>
  );
});

export default PoliciesConfig;
