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

import type {LineMapLayer} from './WifiMapLayers';
import type {MagmaConnectionFeature} from '@fbcnms/ui/insights/map/GeoJSON';
import type {WifiGateway} from './WifiUtils';

import Checkbox from '@material-ui/core/Checkbox';
import Drawer from '@material-ui/core/Drawer';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MapView from '@fbcnms/ui/insights/map/MapView';
import React, {useCallback, useEffect, useMemo, useState} from 'react';
import RefreshIcon from '@material-ui/icons/Refresh';
import TextField from '@material-ui/core/TextField';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import WifiDeviceDetails from './WifiDeviceDetails';
import WifiMapFeatureDetail from './WifiMapFeatureDetail';
import WifiMapMarker from './WifiMapMarker';
import WifiSelectConnType from './WifiSelectConnType';
import WifiSelectMesh from './WifiSelectMesh';

import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import {Route, Switch} from 'react-router-dom';
import {buildLayer, layerNameForFeature, searchDevices} from './WifiMapLayers';
import {
  buildWifiGatewayFromPayloadV1,
  wifiGeoJson,
  wifiGeoJsonConnections,
} from './WifiUtils';
import {groupBy, map} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useRouter, useSnackbar} from '@fbcnms/ui/hooks';

const DRAWER_WIDTH = 440;

const useStyles = makeStyles(() => ({
  mapContainer: {
    width: `calc(100% - ${DRAWER_WIDTH}px)`,
    height: 'calc(100% - 64px)',
  },
  input: {
    display: 'inline-flex',
    margin: '8px',
    width: 'calc(100% - 15px)',
  },
  missingGPS: {
    margin: '8px',
    whiteSpace: 'nowrap',
  },
  viewControls: {
    margin: '8px',
  },
}));

export default function WifiMap() {
  const {match, relativePath} = useRouter();
  return (
    <>
      <Switch>
        <Route path={relativePath('/mesh/:meshID')} render={() => <Map />} />
        <Route path={match.path} render={() => <Map />} />
      </Switch>
    </>
  );
}

function Map() {
  const {match, history} = useRouter();
  const classes = useStyles();
  const {meshID, networkId} = match.params;

  const [searchFilter, setSearchFilter] = useState<string>('');
  const [clickedDeviceId, setClickedDeviceId] = useState<?string>(null);
  const [
    clickedFeatures,
    setClickedFeatures,
  ] = useState<?(MagmaConnectionFeature[])>(null);
  const [devices, setDevices] = useState<Array<WifiGateway>>([]);
  const [connType, setConnType] = useState<?LineMapLayer>('defaultRoute');
  const [showMarkerLabels, setShowMarkerLabels] = useState<boolean>(false);
  const [lastRefreshTime, setLastRefreshTime] = useState<string>(
    new Date().toLocaleString(),
  );

  const {
    isLoading: gatewaysIsLoading,
    error: gatewaysError,
    response: gatewaysResponse,
  } = useMagmaAPI(
    MagmaV1API.getWifiByNetworkIdGateways,
    {networkId: nullthrows(match.params.networkId)},
    undefined, // onResponse
    lastRefreshTime, // cacheCounter
  );

  const {
    isLoading: meshesIsLoading,
    error: meshesError,
    response: meshesResponse,
  } = useMagmaAPI(
    MagmaV1API.getWifiByNetworkIdMeshes,
    {networkId: nullthrows(match.params.networkId)},
    undefined, // onResponse
    lastRefreshTime, // cacheCounter
  );

  useSnackbar(
    'Unable to load map data',
    {variant: 'error'},
    meshesError || gatewaysError,
  );

  // TODO: use onResponse() instead
  useEffect(() => {
    setDevices(
      map(gatewaysResponse || {}) // turn id->device map into device list
        .filter(device => device.device)
        .map(device => buildWifiGatewayFromPayloadV1(device)),
    );
    setClickedFeatures(null); // clear; no easy way to update clicked feature
  }, [gatewaysResponse]);
  const meshes: string[] = meshesResponse || [];

  // Callback when the marker is clicked
  const onMarkerClick = useCallback(
    deviceIDorNumber => {
      const deviceID = String(deviceIDorNumber);
      // if click on same device, "unselect" a device
      if (deviceID === clickedDeviceId) {
        setClickedDeviceId(null);
      } else {
        setClickedDeviceId(deviceID);
      }
    },
    [clickedDeviceId],
  );

  // Callback when the filter is changed
  const onFilterChange = useCallback(
    meshID => {
      // clear selected state
      setClickedDeviceId(null);
      setClickedFeatures(null);
      history.push(`/nms/${networkId}/map/mesh/${meshID}`);
    },
    [history, networkId],
  );

  const onConnChange = useCallback(conn => {
    setClickedFeatures(null);
    setConnType(conn || null);
  }, []);

  // meshID specification controls visible devices on the map
  const [visibleDevices] = searchDevices(devices, null, null, meshID);

  const clickedDevice = useMemo(() => {
    if (!clickedDeviceId) {
      return null;
    }
    return devices.find(d => d.id === clickedDeviceId);
  }, [clickedDeviceId, devices]);

  // remainder of filter controls size of device icons on the map
  const [filteredDevices, layerFilters] = searchDevices(
    visibleDevices,
    searchFilter,
    clickedDevice,
    null,
  );

  // use large icons for matchedDeviceIds
  const matchedDeviceIds = new Set(filteredDevices.map(d => d.id));
  const [geojson, invalidDevices] = wifiGeoJson(
    visibleDevices,
    matchedDeviceIds,
  );

  // group connections into defaultRoute, l3, l2, none
  const groupedConnections = groupBy<LineMapLayer, MagmaConnectionFeature>(
    wifiGeoJsonConnections(visibleDevices),
    f => layerNameForFeature(f),
  );

  // defaultRoute connections ought to be included in l3 connections
  if ('defaultRoute' in groupedConnections) {
    groupedConnections['l3'] = (groupedConnections['l3'] || []).concat(
      groupedConnections['defaultRoute'],
    );
  }

  // create layers based on grouped connections
  const mapLayers = Object.keys(groupedConnections)
    .filter(layerName => connType === null || connType === layerName)
    .map(layerName => buildLayer(layerName, groupedConnections[layerName]));

  return (
    <>
      <MapView
        id="mapView"
        geojson={geojson}
        // buildLayer method returns a source that should be a string instead
        // $FlowFixMe[incompatible-type]
        mapLayers={mapLayers}
        MapMarker={WifiMapMarker}
        onMarkerClick={onMarkerClick}
        mapLayerFilters={layerFilters}
        classes={{mapContainer: classes.mapContainer}}
        onClickFeatures={setClickedFeatures}
        showMarkerLabels={showMarkerLabels}
        zoomHash={`url:${match.params.networkId || ''} mesh:${meshID || ''}`}
      />
      <Drawer
        variant="permanent"
        PaperProps={{style: {width: DRAWER_WIDTH}}}
        anchor="right">
        <div className={classes.viewControls}>
          <FormControlLabel
            control={<Checkbox color="primary" />}
            checked={showMarkerLabels}
            onChange={evt => setShowMarkerLabels(evt.target.checked)}
            label="Show Labels"
          />
          <Tooltip
            title={'Last refreshed: ' + lastRefreshTime}
            placement={'bottom-start'}>
            <IconButton
              color="inherit"
              onClick={() => setLastRefreshTime(new Date().toLocaleString())}>
              <RefreshIcon />
            </IconButton>
          </Tooltip>
        </div>
        {((meshesIsLoading || gatewaysIsLoading) && <LoadingFiller />) || (
          <>
            <div>
              <WifiSelectMesh
                meshes={meshes}
                onChange={onFilterChange}
                selectedMeshID={meshID || ''}
                helperText={'Filter by Mesh ID'}
              />
              <TextField
                label="Search ID or Info"
                className={classes.input}
                value={searchFilter}
                onChange={({target}) => setSearchFilter(target.value)}
              />
              <WifiSelectConnType
                onChange={onConnChange}
                selectedConnType={connType || ''}
              />
            </div>
            {invalidDevices.length > 0 && (
              <Typography variant="caption" className={classes.missingGPS}>
                {invalidDevices.length}{' '}
                {invalidDevices.length == 1 ? 'device' : 'devices'} in filter
                with missing GPS
              </Typography>
            )}
            <WifiDeviceDetails device={clickedDevice} showDevice={true} />
            <WifiMapFeatureDetail features={clickedFeatures} />
          </>
        )}
      </Drawer>
    </>
  );
}
