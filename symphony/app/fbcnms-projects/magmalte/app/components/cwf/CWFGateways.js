/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {cwf_gateway} from '@fbcnms/magma-api';

import AddGatewayDialog from '../AddGatewayDialog';
import Button from '@fbcnms/ui/components/design-system/Button';
import CWFEditGatewayDialog from './CWFEditGatewayDialog';
import ChevronRight from '@material-ui/icons/ChevronRight';
import DeleteIcon from '@material-ui/icons/Delete';
import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';
import EditIcon from '@material-ui/icons/Edit';
import ExpandMore from '@material-ui/icons/ExpandMore';
import IconButton from '@material-ui/core/IconButton';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';

import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import Text from '../../theme/design-system/Text';
import {colors} from '../../theme/default';
import {MAGMAD_DEFAULT_CONFIGS} from '../AddGatewayDialog';
import {Route} from 'react-router-dom';
import {findIndex} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {map} from 'lodash';
import {useCallback, useState} from 'react';
import {useInterval, useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  greCell: {
    paddingBottom: '15px',
    paddingLeft: '75px',
    paddingRight: '15px',
    paddingTop: '15px',
  },
  gatewayCell: {
    padding: '5px',
  },
  paper: {
    margin: theme.spacing(3),
  },
  expandIconButton: {
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
  gatewayName: {
    color: colors.primary.brightGray,
    fontWeight: 'bolder',
    paddingRight: '10px',
  },
}));

const FIVE_MINS = 5 * 60 * 1000;
const REFRESH_INTERVAL = 2 * 60 * 1000;

function CWFGateways(props: WithAlert & {}) {
  const [gateways, setGateways] = useState<?(cwf_gateway[])>(null);
  const [expanded, setExpanded] = useState<Set<string>>(new Set());
  const {match, history, relativePath, relativeUrl} = useRouter();
  const [lastFetchTime, setLastFetchTime] = useState(Date.now());
  const networkId = nullthrows(match.params.networkId);
  const classes = useStyles();

  const {isLoading} = useMagmaAPI(
    MagmaV1API.getCwfByNetworkIdGateways,
    {networkId},
    useCallback(response => setGateways(map(response, g => g)), []),
    lastFetchTime,
  );

  useInterval(() => setLastFetchTime(Date.now()), REFRESH_INTERVAL);

  if (!gateways || isLoading) {
    return <LoadingFiller />;
  }

  const deleteGateway = (gateway: cwf_gateway) => {
    props
      .confirm(`Are you sure you want to delete ${gateway.name}?`)
      .then(confirmed => {
        if (confirmed) {
          MagmaV1API.deleteCwfByNetworkIdGatewaysByGatewayId({
            networkId,
            gatewayId: gateway.id,
          }).then(() =>
            setGateways(gateways.filter(gw => gw.id != gateway.id)),
          );
        }
      });
  };

  const addGateway = async ({
    gatewayID,
    name,
    description,
    hardwareID,
    challengeKey,
    tier,
  }) => {
    await MagmaV1API.postCwfByNetworkIdGateways({
      networkId,
      gateway: {
        carrier_wifi: {
          allowed_gre_peers: [],
        },
        description,
        device: {
          hardware_id: hardwareID,
          key: {
            key: challengeKey,
            key_type: 'SOFTWARE_ECDSA_SHA256', // default key type
          },
        },
        id: gatewayID,
        magmad: MAGMAD_DEFAULT_CONFIGS,
        name,
        tier,
      },
    });

    const gateway = await MagmaV1API.getCwfByNetworkIdGatewaysByGatewayId({
      networkId,
      gatewayId: gatewayID,
    });

    setGateways([...gateways, gateway]);
    history.push(relativeUrl(''));
  };

  const rows = gateways.map(gateway => (
    <>
      <TableRow key={gateway.id}>
        <TableCell className={classes.gatewayCell}>
          <IconButton
            className={classes.expandIconButton}
            onClick={() => {
              const newExpanded = new Set(expanded);
              expanded.has(gateway.id)
                ? newExpanded.delete(gateway.id)
                : newExpanded.add(gateway.id);
              setExpanded(newExpanded);
            }}>
            {expanded.has(gateway.id) ? <ExpandMore /> : <ChevronRight />}
          </IconButton>
          <span className={classes.gatewayName}>{gateway.name}</span>
          <DeviceStatusCircle
            isGrey={!gateway.status?.checkin_time}
            isActive={
              Math.max(0, Date.now() - (gateway.status?.checkin_time || 0)) <
              FIVE_MINS
            }
          />
        </TableCell>
        <TableCell>{gateway.device.hardware_id}</TableCell>
        <TableCell>
          <IconButton
            color="primary"
            onClick={() => history.push(relativeUrl(`/edit/${gateway.id}`))}>
            <EditIcon />
          </IconButton>
          <IconButton color="primary" onClick={() => deleteGateway(gateway)}>
            <DeleteIcon />
          </IconButton>
        </TableCell>
      </TableRow>
      {expanded.has(gateway.id) &&
        gateway.carrier_wifi.allowed_gre_peers.map((gre, i) => (
          <TableRow key={i} classeName={classes.tableRow}>
            <TableCell className={classes.greCell}>{gre.ip}</TableCell>
            <TableCell>{gre.key}</TableCell>
            <TableCell />
          </TableRow>
        ))}
    </>
  ));

  return (
    <div className={classes.paper}>
      <div className={classes.header}>
        <Text variant="h5">Configure Gateways</Text>
        <NestedRouteLink to="/new">
          <Button>Add Gateway</Button>
        </NestedRouteLink>
      </div>
      <Paper elevation={2}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Hardware UUID / GRE Key</TableCell>
              <TableCell />
            </TableRow>
          </TableHead>
          <TableBody>{rows}</TableBody>
        </Table>
      </Paper>
      <Route
        path={relativePath('/new')}
        render={() => (
          <AddGatewayDialog
            onClose={() => history.push(relativeUrl(''))}
            onSave={addGateway}
          />
        )}
      />
      <Route
        path={relativePath('/edit/:gatewayID')}
        render={({match}) => (
          <CWFEditGatewayDialog
            gateway={nullthrows(
              gateways.find(gw => gw.id === match.params.gatewayID),
            )}
            onCancel={() => history.push(relativeUrl(''))}
            onSave={gateway => {
              const newGateways = [...gateways];
              const i = findIndex(newGateways, g => g.id === gateway.id);
              newGateways[i] = gateway;
              setGateways(newGateways);
              history.push(relativeUrl(''));
            }}
          />
        )}
      />
    </div>
  );
}

export default withAlert(CWFGateways);
