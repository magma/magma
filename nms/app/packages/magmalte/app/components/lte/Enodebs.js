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

import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {enodeb} from '@fbcnms/magma-api';

import AddEditEnodebDialog from './AddEditEnodebDialog';
import Button from '@fbcnms/ui/components/design-system/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';

import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {Route} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  paper: {
    margin: theme.spacing(3),
  },
}));

export default function Enodebs() {
  const {match, relativePath, relativeUrl, history} = useRouter();
  const classes = useStyles();
  const [enodebs, setEnodebs] = useState([]);
  const [lastFetchTime, setLastFetchTime] = useState(Date.now());
  const {isLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdEnodebs,
    {networkId: nullthrows(match.params.networkId)},
    useCallback(
      response => setEnodebs(Object.keys(response).map(key => response[key])),
      [],
    ),
    lastFetchTime,
  );

  if (isLoading) {
    return <LoadingFiller />;
  }

  const rows = enodebs.map(enodeb => (
    <EnodebRow
      key={enodeb.serial}
      enodeb={enodeb}
      onSave={() => setLastFetchTime(Date.now())}
    />
  ));

  return (
    <>
      <div className={classes.paper}>
        <div className={classes.header}>
          <Text variant="h5">Configure eNodeB Devices</Text>
          <NestedRouteLink to="/new">
            <Button>Add eNodeB</Button>
          </NestedRouteLink>
        </div>
        <Paper elevation={2}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Serial ID</TableCell>
                <TableCell>Device Class</TableCell>
                <TableCell />
              </TableRow>
            </TableHead>
            <TableBody>{rows}</TableBody>
          </Table>
        </Paper>
        <Route
          path={relativePath('/new')}
          render={() => (
            <AddEditEnodebDialog
              editingEnodeb={null}
              onClose={() => history.push(relativeUrl(''))}
              onSave={() => {
                history.push(relativeUrl(''));
                setLastFetchTime(Date.now());
              }}
            />
          )}
        />
      </div>
    </>
  );
}

type Props = WithAlert & {
  enodeb: enodeb,
  onSave: () => void,
};

function EnodebRowItem(props: Props) {
  const {enodeb} = props;
  const {match, relativePath, history, relativeUrl} = useRouter();
  const deleteEnodeb = () => {
    props
      .confirm(`Are you sure you want to delete ${enodeb.serial}?`)
      .then(confirmed => {
        if (!confirmed) {
          return;
        }
        MagmaV1API.deleteLteByNetworkIdEnodebsByEnodebSerial({
          networkId: nullthrows(match.params.networkId),
          enodebSerial: enodeb.serial,
        }).then(props.onSave);
      });
  };

  return (
    <TableRow key={enodeb.serial}>
      <TableCell>
        {status}
        {enodeb.serial}
      </TableCell>
      <TableCell>{enodeb.config.device_class}</TableCell>
      <TableCell>
        <NestedRouteLink to={`/edit/${enodeb.serial}`}>
          <IconButton>
            <EditIcon />
          </IconButton>
        </NestedRouteLink>
        <IconButton onClick={deleteEnodeb}>
          <DeleteIcon />
        </IconButton>
      </TableCell>
      <Route
        path={relativePath(`/edit/${enodeb.serial}`)}
        render={() => (
          <AddEditEnodebDialog
            editingEnodeb={enodeb}
            onClose={() => history.push(relativeUrl(''))}
            onSave={() => {
              props.onSave();
              history.push(relativeUrl(''));
            }}
          />
        )}
      />
    </TableRow>
  );
}

const EnodebRow = withAlert(EnodebRowItem);
