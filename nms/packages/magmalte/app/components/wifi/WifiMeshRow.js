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

import type {WifiGateway} from './WifiUtils';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';

import AddIcon from '@material-ui/icons/Add';
import Button from '@material-ui/core/Button';
import ChevronRight from '@material-ui/icons/ChevronRight';
import ClipboardLink from '../../components/ClipboardLink';
import DeleteIcon from '@material-ui/icons/Delete';
import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';
import EditIcon from '@material-ui/icons/Edit';
import ExpandMore from '@material-ui/icons/ExpandMore';
import IconButton from '@material-ui/core/IconButton';
import InfoIcon from '@material-ui/icons/Info';
import LinkIcon from '@material-ui/icons/Link';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import React from 'react';
import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';
import Tooltip from '@material-ui/core/Tooltip';
import url from 'url';
import {groupBy} from 'lodash';

import WifiDeviceDetails, {InfoRow} from './WifiDeviceDetails';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(_ => ({
  actionsCell: {
    textAlign: 'right',
  },
  gatewayCell: {
    paddingBottom: '15px',
    paddingLeft: '75px',
    paddingRight: '15px',
    paddingTop: '15px',
  },
  deviceWarning: {
    color: 'red',
    paddingLeft: 40,
  },
  iconButton: {
    padding: '5px',
  },
  meshButton: {
    margin: 0,
    textTransform: 'none',
  },
  meshCell: {
    padding: '5px',
  },
  meshID: {
    color: colors.primary.brightGray,
    fontWeight: 'bolder',
  },
  meshIconButton: {
    color: colors.primary.brightGray,
    padding: '5px',
  },
  tableCell: {
    padding: '15px',
  },
  tableRow: {
    height: 'auto',
    whiteSpace: 'nowrap',
    verticalAlign: 'top',
  },
}));

type Props = WithAlert & {
  enableDeviceEditing?: boolean,
  meshID: string,
  gateways: WifiGateway[],
  onDeleteMesh: string => void,
  onDeleteDevice: () => void,
};

const EXPANDED_STATE_TYPES = {
  none: 0,
  device: 1,
  neighbors: 2,
  fullDump: 3,
  configs: 4,
};

const MESH_ID_PARAM = 'meshID';
const DEVICE_ID_PARAM = 'deviceID';
const EXPANDED_STATE_PARAM = 'expandedState';

type ExpandedGateways = {[key: string]: 0 | 1 | 2 | 3 | 4};

function buildLinkURL(
  meshID: string,
  deviceIDAndState: ?[string, number] = null,
): string {
  const query: {[string]: string | number} = {[MESH_ID_PARAM]: meshID};
  if (deviceIDAndState) {
    query[DEVICE_ID_PARAM] = deviceIDAndState[0];
    query[EXPANDED_STATE_PARAM] = deviceIDAndState[1] ?? 1;
  }
  const {protocol, host, pathname} = window.location;
  return url.format({protocol, host, pathname, query});
}

function WifiMeshRow(props: Props) {
  const classes = useStyles();
  const {location, match} = useRouter();
  const {meshID, gateways} = props;
  const [expandedGateways, setExpandedGateways] = useState<ExpandedGateways>(
    {},
  );
  const [expanded, setExpanded] = useState(false);
  useEffect(() => {
    const queryParams = new URLSearchParams(location.search);
    if (queryParams.get(MESH_ID_PARAM) === meshID) {
      setExpanded(true);
      const deviceID = queryParams.get(DEVICE_ID_PARAM);
      if (deviceID != null) {
        const expandedState = parseInt(queryParams.get(EXPANDED_STATE_PARAM));
        setExpandedGateways({[deviceID]: expandedState});
      }
    }
  }, [location.search, meshID]);

  const handleToggleAllDevices = () => {
    const {gateways} = props;

    // determine old max state by getting max() of all the states
    const maxState = gateways
      .map(gateway => expandedGateways[gateway.id])
      .reduce((max, state) => (state ? Math.max(max, state) : max), 0);

    // calculate next state
    const nextState = (maxState + 1) % Object.keys(EXPANDED_STATE_TYPES).length;

    // assign same next state to all gateways
    if (nextState === 0) {
      // no need to set any gateways for 0/unexpanded state
      setExpandedGateways({});
    } else {
      const newExpandedGateways = gateways
        .map(gateway => gateway.id)
        .reduce((expandedGateways, id) => {
          expandedGateways[id] = nextState;
          return expandedGateways;
        }, {});
      setExpandedGateways(newExpandedGateways);
    }
  };

  const showMeshDeleteDialog = () => {
    props
      .confirm(
        `Are you sure you want to delete mesh "${props.meshID}" and all its devices (count: ${props.gateways.length})?`,
      )
      .then(async confirmed => {
        if (!confirmed) {
          return;
        }

        await Promise.all(
          props.gateways.map(device =>
            MagmaV1API.deleteWifiByNetworkIdGatewaysByGatewayId({
              networkId: nullthrows(match.params.networkId),
              gatewayId: device.id,
            }),
          ),
        );

        await MagmaV1API.deleteWifiByNetworkIdMeshesByMeshId({
          networkId: nullthrows(match.params.networkId),
          meshId: props.meshID,
        });

        props.onDeleteMesh(props.meshID);
      });
  };

  const onExpandGateway = id =>
    setExpandedGateways({
      ...expandedGateways,
      [id]:
        ((expandedGateways[id] | 0) + 1) %
        Object.keys(EXPANDED_STATE_TYPES).length,
    });

  // construct version list per mesh
  const versionGroups: {string: Array<WifiGateway>} = groupBy(
    gateways,
    device => {
      if (device.versionParsed) {
        if (device.versionParsed.fbpkg !== 'none') {
          return device.versionParsed.fbpkg;
        } else {
          return device.versionParsed.hash;
        }
      }
      return device.version || 'UNKNOWN';
    },
  );

  // sort by device count, then version string
  const sortedVersions: Array<string> = Object.keys(versionGroups);
  sortedVersions.sort((a, b) => {
    // keep "Not Reported at the bottom"
    if (a === 'Not Reported') {
      return 1;
    } else if (b === 'Not Reported') {
      return -1;
    } else if (versionGroups[a].length === versionGroups[b].length) {
      // if device counts are equal, then use version string
      return a.localeCompare(b);
    } else {
      // sort by device count
      return versionGroups[b].length - versionGroups[a].length;
    }
  });

  const gatewayVersions = sortedVersions.map(version => (
    <div key={version}>
      <Tooltip
        title={`${versionGroups[version].length} device(s) with ${versionGroups[version][0].version}`}
        enterDelay={100}
        key={version}>
        <span style={{fontFamily: 'monospace'}}>{version}</span>
      </Tooltip>
      :{' '}
      <span style={{fontSize: '88%', fontWeight: 'bold'}}>
        {versionGroups[version].length}
      </span>
    </div>
  ));

  return (
    <>
      <TableRow className={classes.tableRow}>
        <TableCell className={classes.meshCell}>
          <IconButton
            className={classes.meshIconButton}
            onClick={
              gateways.length == 0 ? null : () => setExpanded(!expanded)
            }>
            {expanded ? <ExpandMore /> : <ChevronRight />}
          </IconButton>
          <span className={classes.meshID}>{meshID}</span>
        </TableCell>
        <TableCell className={classes.meshCell}>
          {gateways.length > 0 && (
            <>
              <InfoRow
                label="Up"
                data={`${gateways.filter(gateway => gateway.up).length} of ${
                  gateways.length
                }`}
              />
              {gatewayVersions}
              {expanded && (
                <>
                  <Tooltip
                    title="Click to toggle device info"
                    enterDelay={400}
                    placement={'right'}>
                    <Button
                      size="small"
                      className={classes.meshButton}
                      onClick={handleToggleAllDevices}>
                      toggle info
                    </Button>
                  </Tooltip>
                </>
              )}
            </>
          )}
        </TableCell>
        <TableCell className={classes.actionsCell}>
          <ClipboardLink title="Copy link to this mesh">
            {({copyString}) => (
              <IconButton
                className={classes.iconButton}
                onClick={() => copyString(buildLinkURL(meshID))}>
                <LinkIcon />
              </IconButton>
            )}
          </ClipboardLink>
          <NestedRouteLink to={`/add_device/${meshID}`}>
            <IconButton className={classes.meshIconButton}>
              <AddIcon />
            </IconButton>
          </NestedRouteLink>
          <NestedRouteLink to={`/edit_mesh/${meshID}`}>
            <IconButton className={classes.meshIconButton}>
              <EditIcon />
            </IconButton>
          </NestedRouteLink>
          <IconButton
            className={classes.meshIconButton}
            onClick={showMeshDeleteDialog}>
            <DeleteIcon />
          </IconButton>
        </TableCell>
      </TableRow>
      {expanded &&
        gateways.map(gateway => (
          <GatewayRow
            key={gateway.id}
            meshID={meshID}
            gateway={gateway}
            enableDeviceEditing={props.enableDeviceEditing}
            onDelete={props.onDeleteDevice}
            onExpandGateway={onExpandGateway}
            expandState={expandedGateways[gateway.id]}
          />
        ))}
    </>
  );
}

function GatewayRowElement(
  props: WithAlert & {
    meshID: string,
    gateway: WifiGateway,
    onExpandGateway: string => void,
    expandState: number,
    enableDeviceEditing?: boolean,
    onDelete: () => void,
  },
) {
  const classes = useStyles();
  const match = useRouter();
  const {
    meshID,
    gateway,
    onExpandGateway,
    expandState,
    enableDeviceEditing,
  } = props;

  const showDeviceDeleteDialog = () => {
    props
      .confirm(`Are you sure you want to delete "${gateway.id}"?`)
      .then(async confirmed => {
        if (!confirmed) {
          return;
        }

        // V1 API call will delete all parts of the device
        await MagmaV1API.deleteWifiByNetworkIdGatewaysByGatewayId({
          networkId: nullthrows(match.params.networkId),
          gatewayId: gateway.id,
        });

        props.onDelete();
      });
  };

  return (
    <TableRow className={classes.tableRow} key={gateway.id}>
      <TableCell className={classes.gatewayCell}>{gateway.info}</TableCell>
      <TableCell className={classes.tableCell}>
        {status}
        <DeviceStatusCircle isGrey={!gateway.status} isActive={!!gateway.up} />
        <Tooltip
          title="Click to toggle device info"
          enterDelay={400}
          placement={'right'}>
          <span onClick={() => onExpandGateway(gateway.id)}>{gateway.id}</span>
        </Tooltip>
        {gateway.coordinates.includes(NaN) && (
          <span className={classes.deviceWarning}>
            {' '}
            Please configure Lat/Lng
          </span>
        )}
        {gateway.status &&
          gateway.status.meta &&
          gateway.status.meta['validation_status'] !== 'passed' && (
            <span className={classes.deviceWarning}>
              {' '}
              Please check image validation status
            </span>
          )}

        {!!expandState && gateway.status && (
          <WifiDeviceDetails
            device={gateway}
            hideHeader={true}
            showConfigs={expandState === EXPANDED_STATE_TYPES.configs}
            showDevice={expandState === EXPANDED_STATE_TYPES.device}
            showNeighbors={expandState === EXPANDED_STATE_TYPES.neighbors}
            showFullDump={expandState === EXPANDED_STATE_TYPES.fullDump}
          />
        )}
      </TableCell>
      <TableCell className={classes.actionsCell}>
        <ClipboardLink title="Copy link to this device">
          {({copyString}) => (
            <IconButton
              className={classes.iconButton}
              onClick={() =>
                copyString(buildLinkURL(meshID, [gateway.id, expandState]))
              }>
              <LinkIcon />
            </IconButton>
          )}
        </ClipboardLink>
        <Tooltip title="Click to toggle device info" enterDelay={400}>
          <IconButton
            className={classes.iconButton}
            onClick={() => onExpandGateway(gateway.id)}>
            <InfoIcon />
          </IconButton>
        </Tooltip>
        {enableDeviceEditing && (
          <NestedRouteLink to={`/${meshID}/edit_device/${gateway.id}`}>
            <IconButton className={classes.iconButton}>
              <EditIcon />
            </IconButton>
          </NestedRouteLink>
        )}
        <IconButton
          className={classes.iconButton}
          onClick={showDeviceDeleteDialog}>
          <DeleteIcon />
        </IconButton>
      </TableCell>
    </TableRow>
  );
}

const GatewayRow = withAlert(GatewayRowElement);

export default withAlert(WifiMeshRow);
