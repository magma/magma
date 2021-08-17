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
import PolicyBaseNameDialog from './PolicyBaseNameDialog';
import PolicyRuleEditDialog from './PolicyRuleEditDialog';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';
import Toolbar from '@material-ui/core/Toolbar';

import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {Route} from 'react-router-dom';
import {findIndex} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  header: {
    flexGrow: 1,
  },
  actionsColumn: {
    width: '300px',
  },
}));

function PoliciesConfig(props: WithAlert & {mirrorNetwork?: string}) {
  const classes = useStyles();
  const {match, relativePath, relativeUrl, history} = useRouter();
  const [ruleIDs, setRuleIDs] = useState();
  const [baseNames, setBaseNames] = useState();

  const networkID = nullthrows(match.params.networkId);
  useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPoliciesRules,
    {networkId: networkID},
    setRuleIDs,
  );
  useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPoliciesBaseNames,
    {networkId: networkID},
    setBaseNames,
  );
  const {response: qosProfiles} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdPolicyQosProfiles,
    {
      networkId: networkID,
    },
  );

  if (!ruleIDs || !baseNames) {
    return <LoadingFiller />;
  }

  const onDelete = id => {
    const newRuleIDs = [...nullthrows(ruleIDs)];
    newRuleIDs.splice(
      findIndex(newRuleIDs, id2 => id2 === id),
      1,
    );
    setRuleIDs(newRuleIDs);
  };

  const deleteBaseName = async name => {
    const confirmed = await props.confirm(
      `Are you sure you want to remove the base name "${name}"?`,
    );

    if (confirmed) {
      const data = [
        {
          networkId: networkID,
          baseName: name,
        },
      ];

      if (props.mirrorNetwork) {
        data.push({networkId: props.mirrorNetwork, baseName: name});
      }
      await Promise.all(
        data.map(d =>
          MagmaV1API.deleteNetworksByNetworkIdPoliciesBaseNamesByBaseName(d),
        ),
      );

      const newBaseNames = [...nullthrows(baseNames)];
      newBaseNames.splice(
        findIndex(newBaseNames, name2 => name2 === name),
        1,
      );
      setBaseNames(newBaseNames);
    }
  };

  return (
    <>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>ID</TableCell>
            <TableCell>Precedence</TableCell>
            <TableCell className={classes.actionsColumn}>
              <NestedRouteLink to="/add/">
                <Button>Add Rule</Button>
              </NestedRouteLink>
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {ruleIDs.map(id => (
            <RuleRow
              mirrorNetwork={props.mirrorNetwork}
              key={id}
              ruleID={id}
              onDelete={onDelete}
            />
          ))}
        </TableBody>
      </Table>
      <Route
        path={relativePath('/add')}
        component={() => (
          <PolicyRuleEditDialog
            qosProfiles={qosProfiles ?? {}}
            mirrorNetwork={props.mirrorNetwork}
            onCancel={() => history.push(relativeUrl(''))}
            onSave={ruleID => {
              setRuleIDs([...nullthrows(ruleIDs), ruleID]);
              history.push(relativeUrl(''));
            }}
          />
        )}
      />
      <Toolbar>
        <Text className={classes.header} variant="h5">
          Base Names
        </Text>
      </Toolbar>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>Name</TableCell>
            <TableCell className={classes.actionsColumn}>
              <NestedRouteLink to="/add_base_name/">
                <Button>Add Base Name</Button>
              </NestedRouteLink>
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {baseNames.map(name => (
            <TableRow key={name}>
              <TableCell>{name}</TableCell>
              <TableCell>
                <NestedRouteLink to={`/edit_base_name/${name}`}>
                  <IconButton>
                    <EditIcon />
                  </IconButton>
                </NestedRouteLink>
                <IconButton onClick={() => deleteBaseName(name)}>
                  <DeleteIcon />
                </IconButton>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <Route
        path={relativePath('/add_base_name')}
        component={() => (
          <PolicyBaseNameDialog
            mirrorNetwork={props.mirrorNetwork}
            onCancel={() => history.push(relativeUrl(''))}
            onSave={baseName => {
              setBaseNames([...nullthrows(baseNames), baseName]);
              history.push(relativeUrl(''));
            }}
          />
        )}
      />
      <Route
        path={relativePath('/edit_base_name/:baseName')}
        component={() => (
          <PolicyBaseNameDialog
            onCancel={() => history.push(relativeUrl(''))}
            onSave={() => history.push(relativeUrl(''))}
          />
        )}
      />
    </>
  );
}

type Props = WithAlert & {
  ruleID: string,
  onDelete: () => void,
  mirrorNetwork?: string,
};

const RuleRow = withAlert((props: Props) => {
  const [lastRefreshTime, setLastRefreshTime] = useState(new Date().getTime());
  const {match, relativePath, relativeUrl, history} = useRouter();
  const networkID = nullthrows(match.params.networkId);

  const {response: rule} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPoliciesRulesByRuleId,
    {networkId: networkID, ruleId: encodeURIComponent(props.ruleID)},
    undefined,
    lastRefreshTime,
  );
  const {response: qosProfiles} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdPolicyQosProfiles,
    {
      networkId: networkID,
    },
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

    const data = [
      {
        networkId: networkID,
        ruleId: props.ruleID,
      },
    ];
    if (props.mirrorNetwork) {
      data.push({
        networkId: props.mirrorNetwork,
        ruleId: props.ruleID,
      });
    }
    await Promise.all(
      data.map(d =>
        MagmaV1API.deleteNetworksByNetworkIdPoliciesRulesByRuleId(d),
      ),
    );

    props.onDelete();
  };

  const editPath = `/edit/${props.ruleID}/`;
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
                qosProfiles={qosProfiles ?? {}}
                mirrorNetwork={props.mirrorNetwork}
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

export default withAlert(PoliciesConfig);
