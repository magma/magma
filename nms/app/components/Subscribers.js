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

import type {WithAlert} from '../../fbc_js_core/ui/components/Alert/withAlert';
import type {subscriber} from '../../generated/MagmaAPIBindings';

import AddEditSubscriberDialog from './lte/AddEditSubscriberDialog';
import Button from '../../fbc_js_core/ui/components/design-system/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import ImportSubscribersDialog from './ImportSubscribersDialog';
import LoadingFiller from '../../fbc_js_core/ui/components/LoadingFiller';
import MagmaV1API from '../../generated/WebClient';
import NestedRouteLink from '../../fbc_js_core/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React, {useState} from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableFooter from '@material-ui/core/TableFooter';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '../theme/design-system/Text';

import nullthrows from '../../fbc_js_core/util/nullthrows';
import useMagmaAPI from '../../api/useMagmaAPI';
import withAlert from '../../fbc_js_core/ui/components/Alert/withAlert';
import {CoreNetworkTypes} from '../views/subscriber/SubscriberUtils';
import {Route, Routes, useNavigate, useParams} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {map} from 'lodash';
import {useEnqueueSnackbar} from '../../fbc_js_core/ui/hooks/useSnackbar';

const useStyles = makeStyles(theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  buttons: {
    display: 'flex',
    justifyContent: 'flex-end',
    flexDirection: 'row',
  },
  paper: {
    margin: theme.spacing(3),
  },
  importButton: {
    marginRight: '8px',
  },
}));

const forbiddenNetworkTypes = Object.keys(CoreNetworkTypes).map(
  key => CoreNetworkTypes[key],
);

function Subscribers() {
  const classes = useStyles();
  const params = useParams();
  const navigate = useNavigate();

  const enqueueSnackbar = useEnqueueSnackbar();
  const [lastRefreshTime, setLastRefreshTime] = useState(new Date().getTime());
  const {error, isLoading, response: subscribers} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdSubscribers,
    {networkId: nullthrows(params.networkId)},
    undefined,
    lastRefreshTime,
  );

  const {isLoading: subProfilesLoading, response: epcConfigs} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdCellularEpc,
    {
      networkId: nullthrows(params.networkId),
    },
  );

  const {isLoading: apnsLoading, response: networkAPNs} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdApns,
    {
      networkId: nullthrows(params.networkId),
    },
  );

  if (isLoading || subProfilesLoading || apnsLoading) {
    return <LoadingFiller />;
  }

  const subProfiles = new Set(Object.keys(epcConfigs?.sub_profiles || {})).add(
    'default',
  );

  const apns = new Set(Object.keys(networkAPNs || {}));

  const onSave = () => {
    navigate('');
    setLastRefreshTime(new Date().getTime());
  };

  const onError = reason => {
    enqueueSnackbar(reason, {variant: 'error'});
  };

  const rows = map(subscribers, row => (
    <SubscriberTableRow
      key={row.id}
      subscriber={row}
      onSave={onSave}
      subProfiles={subProfiles}
      apns={apns}
    />
  ));

  return (
    <div className={classes.paper}>
      <div className={classes.header}>
        <Text variant="h5">Subscribers</Text>
        <div className={classes.buttons}>
          <NestedRouteLink to="import">
            <Button className={classes.importButton}>Import</Button>
          </NestedRouteLink>
          <NestedRouteLink to="add">
            <Button>Add Subscriber</Button>
          </NestedRouteLink>
        </div>
      </div>
      <Paper elevation={2}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>IMSI</TableCell>
              <TableCell>LTE Subscription State</TableCell>
              <TableCell>Data Plan</TableCell>
              <TableCell>Active APNs</TableCell>
              <TableCell />
            </TableRow>
          </TableHead>
          <TableBody>{rows}</TableBody>
          <TableFooter
            style={
              Object.keys(subscribers || {}).length === 0 && error === null
                ? {}
                : {display: 'none'}
            }>
            <TableRow>
              <TableCell colSpan="3">No subscribers found</TableCell>
            </TableRow>
          </TableFooter>
        </Table>
      </Paper>
      <div style={error !== null ? {} : {display: 'none'}}>
        <Text color="error" variant="body2">
          {error ?? ''}
        </Text>
      </div>
      <Routes>
        <Route
          path="/import"
          element={
            <ImportSubscribersDialog
              open={true}
              onClose={() => navigate('')}
              onSave={onSave}
              onSaveError={failureIDs => {
                enqueueSnackbar(
                  'Error adding the following subscribers: ' +
                    failureIDs.join(', '),
                  {variant: 'error'},
                );
              }}
            />
          }
        />
        <Route
          path="/add"
          element={
            <AddEditSubscriberDialog
              onClose={() => navigate('')}
              onSave={onSave}
              onSaveError={onError}
              forbiddenNetworkTypes={forbiddenNetworkTypes}
              subProfiles={Array.from(subProfiles)}
              apns={Array.from(apns)}
            />
          }
        />
      </Routes>
    </div>
  );
}

type Props = WithAlert & {
  subscriber: subscriber,
  subProfiles: Set<string>,
  apns: Set<string>,
  onSave: () => void,
};

function SubscriberTableRowComponent(props: Props) {
  const params = useParams();
  const navigate = useNavigate();
  const enqueueSnackbar = useEnqueueSnackbar();
  const {subscriber, subProfiles} = props;
  const displayID = subscriber.id.replace(/^IMSI/, '');
  const onDelete = async () => {
    const confirmed = await props.confirm(
      `Are you sure you want to delete subscriber ${displayID}?`,
    );
    if (confirmed) {
      MagmaV1API.deleteLteByNetworkIdSubscribersBySubscriberId({
        networkId: params.networkId || '',
        subscriberId: subscriber.id,
      })
        .then(props.onSave)
        .catch(error =>
          enqueueSnackbar(error.response.data.message, {variant: 'error'}),
        );
    }
  };

  const subProfile = subProfiles.has(subscriber.lte.sub_profile)
    ? subscriber.lte.sub_profile
    : 'default';

  return (
    <>
      <TableRow>
        <TableCell>{displayID}</TableCell>
        <TableCell>{subscriber.lte.state}</TableCell>
        <TableCell>{subProfile}</TableCell>
        <TableCell>{subscriber.active_apns?.join(', ')}</TableCell>
        <TableCell>
          <NestedRouteLink to={`edit/${subscriber.id}`}>
            <IconButton>
              <EditIcon />
            </IconButton>
          </NestedRouteLink>
          <IconButton onClick={onDelete}>
            <DeleteIcon />
          </IconButton>
        </TableCell>
      </TableRow>
      <Routes>
        <Route
          path={`/edit/${subscriber.id}`}
          element={
            <AddEditSubscriberDialog
              editingSubscriber={subscriber}
              onClose={() => navigate('')}
              onSave={props.onSave}
              onSaveError={reason => {
                enqueueSnackbar(reason, {variant: 'error'});
              }}
              forbiddenNetworkTypes={forbiddenNetworkTypes}
              subProfiles={Array.from(props.subProfiles)}
              apns={Array.from(props.apns)}
            />
          }
        />
      </Routes>
    </>
  );
}

const SubscriberTableRow = withAlert(SubscriberTableRowComponent);

export default Subscribers;
