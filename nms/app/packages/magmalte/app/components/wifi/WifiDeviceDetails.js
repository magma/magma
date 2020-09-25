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

import type {Node} from 'react';
import type {WifiGateway} from './WifiUtils';
import type {WithStyles} from '@material-ui/core';
import type {gateway_status} from '@fbcnms/magma-api';

import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import moment from 'moment';

import {withStyles} from '@material-ui/core/styles';

import {macColonfy} from './WifiUtils';

const styles = {
  deviceDetails: {
    color: 'black',
    margin: '8px',
    whiteSpace: 'nowrap',
  },
  deviceDetailList: {
    margin: '0px',
  },
  tableCell: {
    padding: '2px 5px 2px 5px',
  },
  tableRow: {
    height: 'auto',
  },
  validationField: {
    color: 'red',
  },
};

type Props = WithStyles<typeof styles> & {
  device: ?WifiGateway,
  hideHeader?: boolean,
  showConfigs?: boolean,
  showDevice?: boolean,
  showFullDump?: boolean,
  showNeighbors?: boolean,
};

const SpanSplit = withStyles(styles)(
  (
    props: WithStyles<typeof styles> & {
      field: ?string,
      separator?: string | RegExp,
    },
  ) => {
    if (!props.field) {
      return null;
    }
    const fieldSplit = props.field.split(props.separator || ',');
    if (fieldSplit.length > 0) {
      return (
        <ul className={props.classes.deviceDetailList}>
          {fieldSplit.filter(Boolean).map((data, i) => (
            <li key={i}>{data}</li>
          ))}
        </ul>
      );
    }
    return null;
  },
);

export const InfoRow = (props: {
  label: string,
  data: ?Node,
  className?: string,
}) => {
  if (props.data !== null && props.data !== undefined) {
    return (
      <div className={props.className}>
        <b>{props.label}: </b>
        {props.data}
      </div>
    );
  }
  return null;
};

const NeighborsList = (props: {
  status: ?gateway_status,
  classes: {+[string]: string},
}) => {
  const {status} = props;
  const {tableCell, tableRow} = props.classes;

  if (!status) {
    return null;
  }

  const {meta} = status;
  if (!meta) {
    return null;
  }

  const openrPeers = String(meta.openr_neighbors || '')
    .split(',')
    .filter(Boolean);
  const l2Peers = (meta.mesh0_stations || '')
    .replace(/:/g, '')
    .split(',')
    .filter(Boolean);

  // combine all the things uniquely
  const allPeersSet = new Set([...openrPeers, ...l2Peers]);
  const allPeers = Array.from(allPeersSet);

  if (allPeers.length == 0) {
    return null;
  }

  // sort by mac address
  // TODO: sort by signal strength?
  allPeers.sort();

  const peerList = allPeers.map(mac => {
    const l2 = `mesh0_${macColonfy(mac)}`;
    const l3 = `openr_${mac}`;

    const inactiveTimeMS = meta[`${l2}_inactive_time`]
      ? parseInt(meta[`${l2}_inactive_time`])
      : 0;
    return (
      <TableRow key={mac} className={tableRow}>
        <TableCell
          className={tableCell}
          style={{color: inactiveTimeMS > 5000 ? 'red' : 'black'}}>
          <Tooltip
            title={'inactive time: ' + meta[`${l2}_inactive_time`]}
            enterDelay={100}
            placement={'left'}>
            <span>{mac}</span>
          </Tooltip>
        </TableCell>
        <TableCell className={tableCell}>
          <Tooltip
            title={
              meta[`${l3}_ipv6`] +
              ' ' +
              meta[`${l2}_mesh_plink`] +
              ' ' +
              meta[`${l2}_metric`]
            }
            enterDelay={100}
            placement={'left'}>
            <span>
              {meta[`${l3}_ip`] || '<NONE>'}{' '}
              {meta[`${l3}_metric`] && <> ({meta[`${l3}_metric`]})</>}
              {meta[`${l2}_mesh_plink`] &&
                meta[`${l2}_mesh_plink`] != 'ESTAB' &&
                meta[`${l2}_mesh_plink`]}
            </span>
          </Tooltip>
        </TableCell>
        <TableCell className={tableCell}>
          {(meta[`${l2}_signal`] || '').replace(/\[.+\] /g, '')}
        </TableCell>
        <TableCell className={tableCell}>
          <Tooltip
            title={'TX Expected Throughput / RX Last Frame Bitrate'}
            enterDelay={100}
            placement={'left'}>
            <span>
              {(meta[`${l2}_expected_throughput`] || '').replace('Mbps', '') ||
                'NA'}{' '}
              /{' '}
              {(meta[`${l2}_rx_bitrate`] || '').replace(/ MBit.+/g, '') || 'NA'}
            </span>
          </Tooltip>
        </TableCell>
      </TableRow>
    );
  });

  return (
    <Table className={props.classes.table}>
      <TableHead>
        <TableRow className={tableRow}>
          <TableCell component="th" className={tableCell}>
            mac
          </TableCell>
          <TableCell component="th" className={tableCell}>
            IP (metric)
          </TableCell>
          <TableCell component="th" className={tableCell}>
            RSSI
          </TableCell>
          <TableCell component="th" className={tableCell}>
            last frame TX/RX Mbps
          </TableCell>
        </TableRow>
      </TableHead>
      <TableBody>{peerList}</TableBody>
    </Table>
  );
};

const WifiDeviceDetails = (props: Props) => {
  if (!props.device) {
    return null;
  }

  const {
    device,
    hideHeader,
    showConfigs,
    showDevice,
    showFullDump,
    showNeighbors,
  } = props;

  const {status, wifi_config} = device;

  return (
    <Typography
      component="div"
      className={props.classes.deviceDetails}
      variant="body2">
      {!hideHeader && (
        <>
          <InfoRow label="ID" data={device.id} />
          <InfoRow label="Info" data={device.info} />
          <InfoRow
            label="Status"
            data={
              <>
                <DeviceStatusCircle isGrey={!status} isActive={!!device.up} />
                {device.up ? 'UP' : 'DOWN'}
              </>
            }
          />
        </>
      )}

      {showDevice && (
        <>
          <InfoRow label="Last Check-in" data={device.lastCheckin} />
          <InfoRow
            label="Last Check-in Time"
            data={
              (device.checkinTime != null &&
                moment(device.checkinTime)
                  .utc()
                  .format('ddd YYYY-MM-DD HH:mm:ss UTC')) ||
              null
            }
          />
          {status && (
            <>
              <InfoRow label="OS Uptime" data={status.meta?.uptime} />

              {device.isGateway && (
                <div>
                  <span style={{backgroundColor: 'LightBlue'}}>
                    <b>(GATEWAY)</b>
                  </span>
                </div>
              )}

              <InfoRow
                label="Default Route"
                data={
                  status.meta?.default_route && (
                    <SpanSplit
                      field={status.meta?.default_route}
                      separator=";"
                    />
                  )
                }
              />

              <InfoRow
                label="Carrier Detected"
                data={
                  status.meta?.carrier_detected && (
                    <SpanSplit field={status.meta?.carrier_detected} />
                  )
                }
              />

              <InfoRow
                label="Public IP"
                data={
                  status.meta?.public_ip && (
                    <SpanSplit field={status.meta?.public_ip} />
                  )
                }
              />

              <InfoRow
                label="VPN IP"
                data={
                  status.meta?.tun_vpn_ip && (
                    <SpanSplit field={status.meta?.tun_vpn_ip} />
                  )
                }
              />

              <InfoRow
                label="WAN IP"
                data={
                  status.meta?.eth0_ip && (
                    <SpanSplit field={status.meta?.eth0_ip} />
                  )
                }
              />

              <InfoRow label="Mesh Mac" data={status.meta?.mesh0_hw_addr} />

              <InfoRow
                label="Mesh IP"
                data={
                  status.meta?.mesh0_ip && (
                    <SpanSplit field={status.meta?.mesh0_ip} />
                  )
                }
              />

              <InfoRow
                label="OpenR Neighbors"
                data={
                  status.meta?.openr_neighbors &&
                  String(status.meta?.openr_neighbors).split(',').length
                }
              />

              <InfoRow
                label="L2 Established Peers"
                data={status.meta?.mesh0_num_stations_estab}
              />
              <InfoRow
                label="L2 Visible Peers"
                data={status.meta?.mesh0_num_stations}
              />

              <InfoRow label="SSID" data={status.meta?.wlan_soma_ssid} />
              <InfoRow label="BSSID" data={status.meta?.wlan_soma_hw_addr} />
              <InfoRow
                label="AP Channel"
                data={status.meta?.wlan_soma_channel}
              />
              <InfoRow
                label="AP Clients"
                data={status.meta?.wlan_soma_num_stations}
              />
            </>
          )}

          {device.versionParsed ? (
            <>
              <InfoRow label="Version Hash" data={device.versionParsed.hash} />
              <InfoRow
                label="Version Info"
                data={`${device.versionParsed.buildtime} by ${device.versionParsed.user}`}
              />
              <InfoRow label="Package" data={device.versionParsed.fbpkg} />
              <InfoRow label="Config Overlay" data={device.versionParsed.cfg} />
            </>
          ) : (
            <InfoRow label="Version" data={device.version} />
          )}

          {status && status.meta?.validation_status !== 'passed' && (
            <>
              <InfoRow
                label="Validation Status"
                data={status.meta?.validation_status}
                className={props.classes.validationField}
              />
              <InfoRow
                label="Validation Action"
                data={status.meta?.validation_decision}
                className={props.classes.validationField}
              />
              <InfoRow
                label="Validation Run count"
                data={status.meta?.validation_run_count}
                className={props.classes.validationField}
              />
              <InfoRow
                label="Validation Failed tests"
                data={status.meta?.validation_tests_failed}
                className={props.classes.validationField}
              />
            </>
          )}
        </>
      )}

      {showNeighbors &&
        ((status?.meta && (
          <>
            <InfoRow
              label="OpenR Neighbors"
              data={
                status?.meta?.openr_neighbors &&
                String(status?.meta?.openr_neighbors).split(',').length
              }
            />
            <InfoRow
              label="L2 Established Peers"
              data={status?.meta?.mesh0_num_stations_estab}
            />
            <InfoRow
              label="L2 Visible Peers"
              data={status?.meta?.mesh0_num_stations}
            />
            <NeighborsList status={status} classes={props.classes} />
          </>
        )) || <>No Neighbor Information</>)}

      {showConfigs &&
        ((wifi_config && (
          <>
            {Object.keys(wifi_config)
              .sort((a, b) => a.localeCompare(b))
              .map((key: string) => {
                if (key !== 'additional_props') {
                  return (
                    <div key={key}>
                      <span style={{fontWeight: 'bold'}}>config.{key}</span>:{' '}
                      <span style={{whiteSpace: 'normal'}}>
                        {wifi_config[key]}
                      </span>
                    </div>
                  );
                } else if (wifi_config[key]) {
                  return Object.keys(wifi_config[key])
                    .filter(Boolean)
                    .sort((a, b) => a.localeCompare(b))
                    .map((propkey: string) => (
                      <div key={propkey}>
                        <span style={{fontWeight: 'bold'}}>
                          config.prop.{propkey}
                        </span>
                        :{' '}
                        <span style={{whiteSpace: 'normal'}}>
                          {wifi_config[key]?.[propkey]}
                        </span>
                      </div>
                    ));
                }
              })}
          </>
        )) || <div>No wifi config</div>)}

      {showFullDump &&
        ((status?.meta && (
          <>
            {Object.keys(status?.meta || {})
              .sort((a, b) => a.localeCompare(b))
              .map((key: string) => (
                <div key={key}>
                  <span style={{fontWeight: 'bold'}}>{key}</span>:{' '}
                  <span style={{whiteSpace: 'normal'}}>
                    {(status?.meta || {})[key]}
                  </span>
                </div>
              ))}
          </>
        )) || <div>No status</div>)}
    </Typography>
  );
};

export default withStyles(styles)(WifiDeviceDetails);
