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
import type {DataRows} from '../../components/DataGrid';

import AddEditEnodeButton from './EnodebDetailConfigEdit';
import Button from '@material-ui/core/Button';
import DataGrid from '../../components/DataGrid';
import EnodebContext from '../../components/context/EnodebContext';
import Grid from '@material-ui/core/Grid';
import JsonEditor from '../../components/JsonEditor';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import nullthrows from '@fbcnms/util/nullthrows';

import {CardTitleFilterRow} from '../../components/layout/CardTitleRow';
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

  function editJSON() {
    return (
      <Button
        className={classes.appBarBtn}
        onClick={() => {
          history.push(relativeUrl('/json'));
        }}>
        Edit JSON
      </Button>
    );
  }

  function editEnodeb() {
    return (
      <AddEditEnodeButton
        title={'Edit'}
        isLink={true}
        editProps={{
          editTable: 'config',
        }}
      />
    );
  }

  function editRAN() {
    return (
      <AddEditEnodeButton
        title={'Edit'}
        isLink={true}
        editProps={{
          editTable: 'ran',
        }}
      />
    );
  }

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <CardTitleFilterRow
            label="Config"
            icon={SettingsIcon}
            filter={editJSON}
          />
        </Grid>

        <Grid item xs={12} md={6}>
          <CardTitleFilterRow label="eNodeb" filter={editEnodeb} />
          <EnodebInfoConfig />
        </Grid>

        <Grid item xs={12} md={6}>
          <CardTitleFilterRow label="RAN" filter={editRAN} />
          <EnodebRanConfig />
        </Grid>
      </Grid>
    </div>
  );
}

function EnodebRanConfig() {
  const ctx = useContext(EnodebContext);
  const {match} = useRouter();
  const enodebSerial: string = nullthrows(match.params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];
  const lteRanConfigs = ctx.state.lteRanConfigs;

  const data: DataRows[] = [
    [
      {
        category: 'Bandwidth',
        value: enbInfo.enb.config.bandwidth_mhz ?? '-',
      },
    ],
    [
      {
        category: 'Cell ID',
        value: enbInfo.enb.config.cell_id ?? '-',
      },
    ],
    [
      {
        category: 'Description',
        value: enbInfo.enb.description ?? '-',
      },
    ],
    [
      {
        category: 'RAN Config',
        value: lteRanConfigs?.tdd_config
          ? 'TDD'
          : lteRanConfigs?.fdd_config
          ? 'FDD'
          : '-',
        collapse: lteRanConfigs?.tdd_config ? (
          <EnodeConfigTdd
            earfcndl={enbInfo.enb.config.earfcndl ?? 0}
            specialSubframePattern={
              enbInfo.enb.config.special_subframe_pattern ?? 0
            }
            subframeAssignment={enbInfo.enb.config.subframe_assignment ?? 0}
          />
        ) : lteRanConfigs?.fdd_config ? (
          <EnodeConfigFdd
            earfcndl={enbInfo.enb.config.earfcndl ?? 0}
            earfcnul={lteRanConfigs.fdd_config.earfcnul}
          />
        ) : (
          false
        ),
      },
    ],
    [
      {
        category: 'PCI',
        value: enbInfo.enb.config.pci ?? '-',
      },
    ],
    [
      {
        category: 'TAC',
        value: enbInfo.enb.config.tac ?? '-',
      },
    ],
    [
      {
        category: 'Transmit',
        value: enbInfo.enb.config.transmit_enabled ? 'Enabled' : 'Disabled',
      },
    ],
  ];

  return <DataGrid data={data} testID="ran" />;
}

function EnodebInfoConfig() {
  const ctx = useContext(EnodebContext);
  const {match} = useRouter();
  const enodebSerial: string = nullthrows(match.params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];

  const data: DataRows[] = [
    [
      {
        category: 'Name',
        value: enbInfo.enb.name,
      },
    ],
    [
      {
        category: 'Serial Number',
        value: enbInfo.enb.serial,
      },
    ],
    [
      {
        category: 'Description',
        value: enbInfo.enb.description ?? '-',
      },
    ],
  ];

  return <DataGrid data={data} testID="config" />;
}
