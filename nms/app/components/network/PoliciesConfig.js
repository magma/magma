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

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Button from '@material-ui/core/Button';
import CircularProgress from '@material-ui/core/CircularProgress';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {WithAlert} from '../Alert/withAlert';
// $FlowFixMe migrated to typescript
import LoadingFiller from '../LoadingFiller';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import LoadingFillerBackdrop from '../LoadingFillerBackdrop';
import MagmaV1API from '../../../generated/WebClient';
// $FlowFixMe migrated to typescript
import NestedRouteLink from '../NestedRouteLink';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import PolicyBaseNameDialog from './PolicyBaseNameDialog';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import PolicyRuleEditDialog from './PolicyRuleEditDialog';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../theme/design-system/Text';
import Toolbar from '@material-ui/core/Toolbar';

// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';
import useMagmaAPI from '../../../api/useMagmaAPIFlow';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import withAlert from '../Alert/withAlert';
import {Route, Routes, useNavigate, useParams} from 'react-router-dom';
import {findIndex} from 'lodash';
import {makeStyles} from '@material-ui/styles';
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
  const navigate = useNavigate();
  const params = useParams();
  const [ruleIDs, setRuleIDs] = useState();
  const [baseNames, setBaseNames] = useState();

  const networkID = nullthrows(params.networkId);
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
              <NestedRouteLink to="add/">
                <Button variant="contained" color="primary">
                  Add Rule
                </Button>
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
      <Routes>
        <Route
          path="/add"
          element={
            <PolicyRuleEditDialog
              qosProfiles={qosProfiles ?? {}}
              mirrorNetwork={props.mirrorNetwork}
              onCancel={() => navigate('')}
              onSave={ruleID => {
                setRuleIDs([...nullthrows(ruleIDs), ruleID]);
                navigate('');
              }}
            />
          }
        />
      </Routes>
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
              <NestedRouteLink to="add_base_name/">
                <Button variant="contained" color="primary">
                  Add Base Name
                </Button>
              </NestedRouteLink>
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {baseNames.map(name => (
            <TableRow key={name}>
              <TableCell>{name}</TableCell>
              <TableCell>
                <NestedRouteLink to={`edit_base_name/${name}`}>
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
      <Routes>
        <Route
          path="/add_base_name"
          element={
            <PolicyBaseNameDialog
              mirrorNetwork={props.mirrorNetwork}
              onCancel={() => navigate('')}
              onSave={baseName => {
                setBaseNames([...nullthrows(baseNames), baseName]);
                navigate('');
              }}
            />
          }
        />
        <Route
          path="/edit_base_name/:baseName"
          element={
            <PolicyBaseNameDialog
              onCancel={() => navigate('')}
              onSave={() => navigate('')}
            />
          }
        />
      </Routes>
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
  const navigate = useNavigate();
  const params = useParams();
  const networkID = nullthrows(params.networkId);

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

  const editPath = `edit/${props.ruleID}/`;
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
        <Routes>
          <Route
            path={editPath}
            element={
              rule ? (
                <PolicyRuleEditDialog
                  qosProfiles={qosProfiles ?? {}}
                  mirrorNetwork={props.mirrorNetwork}
                  rule={rule}
                  onCancel={() => navigate('')}
                  onSave={() => {
                    setLastRefreshTime(new Date().getTime());
                    navigate('');
                  }}
                />
              ) : (
                <LoadingFillerBackdrop />
              )
            }
          />
        </Routes>
      </TableCell>
    </TableRow>
  );
});

export default withAlert(PoliciesConfig);
