/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import ActionTable from '../../components/ActionTable';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import RssFeedIcon from '@material-ui/icons/RssFeed';
import Text from '@fbcnms/ui/components/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const APN_TITLE = 'APNs';
const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  topBar: {
    backgroundColor: colors.primary.mirage,
    padding: '20px 40px 20px 40px',
    color: colors.primary.white,
  },
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    padding: `0 ${theme.spacing(5)}px`,
  },
  tabs: {
    color: colors.primary.white,
  },
  tab: {
    fontSize: '18px',
    textTransform: 'none',
  },
  tabLabel: {
    padding: '16px 0 16px 0',
    display: 'flex',
    alignItems: 'center',
  },
  tabIconLabel: {
    marginRight: '8px',
  },
  appBarBtn: {
    color: colors.primary.white,
    background: colors.primary.comet,
    fontFamily: typography.button.fontFamily,
    fontWeight: typography.button.fontWeight,
    fontSize: typography.button.fontSize,
    lineHeight: typography.button.lineHeight,
    letterSpacing: typography.button.letterSpacing,

    '&:hover': {
      background: colors.primary.mirage,
    },
  },
  appBarBtnSecondary: {
    color: colors.primary.white,
  },
  // TODO: remove this when we actually fill out the grid sections
  contentPlaceholder: {
    padding: '50px 0',
  },
  paper: {
    height: 100,
    padding: theme.spacing(10),
    textAlign: 'center',
  },
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
}));

type ApnRowType = {
  apnID: string,
  description: string,
  qosProfile: number,
  added: Date,
};

export default function ApnOverview() {
  const classes = useStyles();
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const {response, isLoading} = useMagmaAPI(MagmaV1API.getLteByNetworkIdApns, {
    networkId: networkId,
  });

  if (isLoading) {
    return <LoadingFiller />;
  }
  const apnRows: Array<ApnRowType> = response
    ? Object.keys(response).map((apn: string) => {
        return {
          apnID: apn,
          description: 'Test APN description',
          qosProfile: 1,
          added: new Date(0),
        };
      })
    : [];
  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={3}>
        <Grid container>
          <Grid item xs={6}>
            <Text key="title">
              <RssFeedIcon /> {APN_TITLE}
            </Text>
          </Grid>
          <Grid
            container
            item
            xs={6}
            justify="flex-end"
            alignItems="center"
            spacing={2}>
            <Button className={classes.appBarBtn}>Add New APN</Button>
          </Grid>
        </Grid>

        <Grid item xs={12}>
          <ActionTable
            data={apnRows}
            columns={[
              {title: 'Apn ID', field: 'apnID'},
              {title: 'Description', field: 'description'},
              {title: 'Qos Profile', field: 'qosProfile', type: 'numeric'},
              {title: 'Added', field: 'added', type: 'datetime'},
            ]}
            menuItems={[{name: 'Edit'}, {name: 'Deactivate'}, {name: 'Remove'}]}
            options={{
              actionsColumnIndex: -1,
              pageSizeOptions: [5, 10],
            }}
          />
        </Grid>
      </Grid>
    </div>
  );
}
