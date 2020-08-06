/*
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
import AddEditEnodeButton from './EnodebDetailConfigEdit';
import Button from '@material-ui/core/Button';
import Divider from '@material-ui/core/Divider';
import EnodebContext from '../../components/context/EnodebContext';
import Grid from '@material-ui/core/Grid';
import JsonEditor from '../../components/JsonEditor';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import Text from '@fbcnms/ui/components/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';

import {EnodeConfigFdd} from './EnodebDetailConfigFdd';
import {EnodeConfigTdd} from './EnodebDetailConfigTdd';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  itemTitle: {
    color: colors.primary.comet,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  itemValue: {
    color: colors.primary.brightGray,
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
}));

export function EnodebJsonConfig() {
  const ctx = useContext(EnodebContext);
  const {match} = useRouter();
  const [error, setError] = useState('');
  const enodebSerial: string = nullthrows(match.params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];
  const enqueueSnackbar = useEnqueueSnackbar();

  return (
    <JsonEditor
      content={enbInfo.enb}
      error={error}
      onSave={async enb => {
        try {
          ctx.setState(enbInfo.enb.serial, {...enbInfo, enb: enb});
          enqueueSnackbar('eNodeb saved successfully', {
            variant: 'success',
          });
          setError('');
        } catch (e) {
          setError(e.response?.data?.message ?? e.message);
        }
      }}
    />
  );
}

export default function EnodebConfig() {
  const classes = useStyles();
  const {history, relativeUrl} = useRouter();

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={3} alignItems="stretch">
        <Grid container spacing={3} alignItems="stretch" item xs={12}>
          <Grid container item xs={12}>
            <Grid item xs={6}>
              <Text>
                <SettingsIcon /> Config
              </Text>
            </Grid>
            <Grid container item xs={6} justify="flex-end">
              <Text>
                <Button
                  className={classes.appBarBtn}
                  onClick={() => {
                    history.push(relativeUrl('/json'));
                  }}>
                  Edit JSON
                </Button>
              </Text>
            </Grid>
          </Grid>

          <Grid item xs={6}>
            <Grid container>
              <Grid item xs={6}>
                <Text>eNodeB</Text>
              </Grid>
              <Grid container item xs={6} justify="flex-end">
                <AddEditEnodeButton
                  title={'Edit'}
                  isLink={true}
                  editProps={{
                    editTable: 'config',
                  }}
                />
              </Grid>
            </Grid>
            <EnodebInfoConfig />
          </Grid>

          <Grid item xs={6}>
            <Grid container>
              <Grid item xs={6}>
                <Text>RAN</Text>
              </Grid>
              <Grid container item xs={6} justify="flex-end">
                <AddEditEnodeButton
                  title={'Edit'}
                  isLink={true}
                  editProps={{
                    editTable: 'ran',
                  }}
                />
              </Grid>
            </Grid>
            <EnodebRanConfig />
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

function EnodebRanConfig() {
  const classes = useStyles();
  const ctx = useContext(EnodebContext);
  const {match} = useRouter();
  const enodebSerial: string = nullthrows(match.params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];
  const lteRanConfigs = ctx.state.lteRanConfigs;
  const typographyProps = {
    primaryTypographyProps: {
      variant: 'caption',
      className: classes.itemTitle,
    },
    secondaryTypographyProps: {
      variant: 'h6',
      className: classes.itemValue,
    },
  };

  return (
    <List component={Paper} data-testid="ran">
      <ListItem>
        <ListItemText
          primary="Bandwidth"
          secondary={enbInfo.enb.config.bandwidth_mhz}
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          secondary={enbInfo.enb.config.cell_id}
          primary="Cell ID"
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      {lteRanConfigs?.tdd_config && (
        <EnodeConfigTdd
          earfcndl={enbInfo.enb.config.earfcndl ?? 0}
          specialSubframePattern={
            enbInfo.enb.config.special_subframe_pattern ?? 0
          }
          subframeAssignment={enbInfo.enb.config.subframe_assignment ?? 0}
        />
      )}
      {lteRanConfigs?.fdd_config && (
        <EnodeConfigFdd
          earfcndl={enbInfo.enb.config.earfcndl ?? 0}
          earfcnul={lteRanConfigs.fdd_config.earfcnul}
        />
      )}
      <Divider />
      <ListItem>
        <ListItemText
          secondary={enbInfo.enb.config.pci}
          primary="PCI"
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          secondary={enbInfo.enb.config.tac}
          primary="TAC"
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          secondary={
            enbInfo.enb.config.transmit_enabled ? 'Enabled' : 'Disabled'
          }
          primary="Transmit"
          {...typographyProps}
        />
      </ListItem>
    </List>
  );
}

function EnodebInfoConfig() {
  const classes = useStyles();
  const ctx = useContext(EnodebContext);
  const {match} = useRouter();
  const enodebSerial: string = nullthrows(match.params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];

  const typographyProps = {
    primaryTypographyProps: {
      variant: 'caption',
      className: classes.itemTitle,
    },
    secondaryTypographyProps: {
      variant: 'h6',
      className: classes.itemValue,
    },
  };
  return (
    <List component={Paper} data-testid="config">
      <ListItem>
        <ListItemText
          primary="Name"
          secondary={enbInfo.enb.name}
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          primary="Serial Number"
          secondary={enbInfo.enb.serial}
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          primary="Description"
          secondary={enbInfo.enb.description}
          {...typographyProps}
        />
      </ListItem>
    </List>
  );
}
