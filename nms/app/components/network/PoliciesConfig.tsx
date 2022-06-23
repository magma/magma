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
import CircularProgress from '@material-ui/core/CircularProgress';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '../LoadingFiller';
import LoadingFillerBackdrop from '../LoadingFillerBackdrop';
import MagmaAPI from '../../../api/MagmaAPI';
import NestedRouteLink from '../NestedRouteLink';
import PolicyBaseNameDialog from './PolicyBaseNameDialog';
import PolicyRuleEditDialog from './PolicyRuleEditDialog';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '../../theme/design-system/Text';
import Toolbar from '@material-ui/core/Toolbar';
import nullthrows from '../../../shared/util/nullthrows';
import useMagmaAPI from '../../../api/useMagmaAPI';
import withAlert from '../Alert/withAlert';
import {Route, Routes, useNavigate, useParams} from 'react-router-dom';
import {findIndex} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';
import type {WithAlert} from '../Alert/withAlert';

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
  const [ruleIDs, setRuleIDs] = useState<Array<string>>();
  const [baseNames, setBaseNames] = useState<Array<string>>();
  const networkID = nullthrows(params.networkId);

  useMagmaAPI(
    MagmaAPI.policies.networksNetworkIdPoliciesRulesGet,
    {
      networkId: networkID,
    },
    setRuleIDs,
  );
  useMagmaAPI(
    MagmaAPI.policies.networksNetworkIdPoliciesBaseNamesGet,
    {
      networkId: networkID,
    },
    setBaseNames,
  );

  const {response: qosProfiles} = useMagmaAPI(
    MagmaAPI.policies.lteNetworkIdPolicyQosProfilesGet,
    {
      networkId: networkID,
    },
  );

  if (!ruleIDs || !baseNames) {
    return <LoadingFiller />;
  }

  const onDelete = (id: string) => {
    const newRuleIDs = [...nullthrows(ruleIDs)];
    newRuleIDs.splice(
      findIndex(newRuleIDs, id2 => id2 === id),
      1,
    );
    setRuleIDs(newRuleIDs);
  };

  const deleteBaseName = async (name: string) => {
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
        data.push({
          networkId: props.mirrorNetwork,
          baseName: name,
        });
      }

      await Promise.all(
        data.map(d =>
          MagmaAPI.policies.networksNetworkIdPoliciesBaseNamesBaseNameDelete(d),
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
                <IconButton onClick={() => void deleteBaseName(name)}>
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
  ruleID: string;
  onDelete: (id: string) => void;
  mirrorNetwork?: string;
};

const RuleRow = withAlert((props: Props) => {
  const [lastRefreshTime] = useState(new Date().getTime());
  const navigate = useNavigate();
  const params = useParams();
  const networkID = nullthrows(params.networkId);

  const {response: rule} = useMagmaAPI(
    MagmaAPI.policies.networksNetworkIdPoliciesRulesRuleIdGet,
    {
      networkId: networkID,
      ruleId: encodeURIComponent(props.ruleID),
    },
    undefined,
    lastRefreshTime,
  );
  const {response: qosProfiles} = useMagmaAPI(
    MagmaAPI.policies.lteNetworkIdPolicyQosProfilesGet,
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
        MagmaAPI.policies.networksNetworkIdPoliciesRulesRuleIdDelete(d),
      ),
    );

    props.onDelete(props.ruleID);
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
        <IconButton onClick={() => void onDeleteRule()}>
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
