/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CheckindGateway} from '../../common/MagmaAPIType';
import type {MagmaFeatureCollection} from '../../common/GeoJSON';

import Alert from '@fbcnms/ui/components/Alert/Alert';
import GatewayMapMarker from './GatewayMapMarker';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MapView from '../MapView';
import Paper from '@material-ui/core/Paper';
import React from 'react';

import {MagmaAPIUrls} from '../../common/MagmaAPI';
import {get} from 'lodash';
import {useAxios, useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

function buildGeoJson(gateways): MagmaFeatureCollection {
  const features = (gateways || [])
    .map((gateway, i) => {
      let longitude =
        parseFloat(get(gateway, 'status.meta.gps_longitude', '0')) || 0;
      if (longitude > 1000000 || longitude < -1000000) {
        // There's a bug in the enodeb that doesn't include the decimal point.
        // This is the best fix we can do for now.
        longitude = longitude / 1000000;
      }
      const latitude =
        parseFloat(get(gateway, 'status.meta.gps_latitude', '0')) || 0;

      // exclude gateways without valid coordinates
      // TODO: enable this after development is done
      // if (longitude === 0 && latitude === 0) {
      //   return null;
      // }

      return {
        type: 'Feature',
        geometry: {
          type: 'Point',
          coordinates: [longitude, latitude],
        },
        properties: {
          id: i,
          iconSize: [60, 60],
          gateway,
        },
      };
    })
    .filter(gateway => gateway !== null);

  return {
    type: 'FeatureCollection',
    features,
  };
}

export default function Insights() {
  const {match} = useRouter();
  const networkId = match.params.networkId || '';

  const [showDialog, setShowDialog] = useState(false);
  const {error, isLoading, response} = useAxios<null, CheckindGateway[]>({
    method: 'get',
    url: MagmaAPIUrls.gateways(networkId, true),
  });

  if (error) {
    return (
      <Alert
        confirmLabel="Okay"
        open={error && showDialog}
        message={error.response?.data?.message || error}
        onConfirm={() => setShowDialog(false)}
      />
    );
  }
  if (isLoading || !response) {
    return (
      <Paper elevation={2}>
        <LoadingFiller />
      </Paper>
    );
  }
  const gateways = response.data.filter(state => state.record);
  return (
    <MapView geojson={buildGeoJson(gateways)} MapMarker={GatewayMapMarker} />
  );
}
