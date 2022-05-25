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
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {DataRows} from '../../components/DataGrid';
// $FlowFixMe migrated to typescript
import type {EnodebInfo} from '../../components/lte/EnodebUtils';
import type {network_ran_configs} from '../../../generated/MagmaAPIBindings';

import AddEditEnodeButton from './EnodebDetailConfigEdit';
import Button from '@material-ui/core/Button';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import CardTitleRow from '../../components/layout/CardTitleRow';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import DataGrid from '../../components/DataGrid';
// $FlowFixMe migrated to typescript
import EnodebContext from '../../components/context/EnodebContext';
import Grid from '@material-ui/core/Grid';
// $FlowFixMe migrated to typescript
import JsonEditor from '../../components/JsonEditor';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';

import {EnodeConfigFdd} from './EnodebDetailConfigFdd';
import {EnodeConfigTdd} from './EnodebDetailConfigTdd';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useNavigate, useParams} from 'react-router-dom';

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
  const params = useParams();
  const [error, setError] = useState('');
  const enodebSerial: string = nullthrows(params.enodebSerial);
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
  const navigate = useNavigate();
  const ctx = useContext(EnodebContext);
  const params = useParams();
  const enodebSerial: string = nullthrows(params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];
  const lteRanConfigs = ctx.lteRanConfigs;
  const enbManaged = enbInfo.enb.enodeb_config?.config_type === 'MANAGED';

  function editJSON() {
    return (
      <Button
        className={classes.appBarBtn}
        onClick={() => {
          navigate('json');
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
          <CardTitleRow label="Config" icon={SettingsIcon} filter={editJSON} />
        </Grid>

        <Grid item xs={12} md={6}>
          <CardTitleRow label="eNodeb" filter={editEnodeb} />
          <EnodebInfoConfig />
        </Grid>

        <Grid item xs={12} md={6}>
          <CardTitleRow label="RAN" filter={editRAN} />
          {enbManaged ? (
            <EnodebManagedRanConfig
              enbInfo={enbInfo}
              lteRanConfigs={lteRanConfigs}
            />
          ) : (
            <EnodebUnmanagedRanConfig enbInfo={enbInfo} />
          )}
        </Grid>
      </Grid>
    </div>
  );
}

function EnodebManagedRanConfig({
  enbInfo,
  lteRanConfigs,
}: {
  enbInfo: EnodebInfo,
  lteRanConfigs?: network_ran_configs,
}) {
  const managedConfig: DataRows[] = [
    [
      {
        category: 'eNodeB Externally Managed',
        value: 'False',
      },
    ],
    [
      {
        category: 'Bandwidth',
        value: enbInfo.enb.enodeb_config?.managed_config?.bandwidth_mhz ?? '-',
      },
    ],
    [
      {
        category: 'Cell ID',
        value: enbInfo.enb.enodeb_config?.managed_config?.cell_id ?? '-',
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
            earfcndl={enbInfo.enb.enodeb_config?.managed_config?.earfcndl ?? 0}
            specialSubframePattern={
              enbInfo.enb.enodeb_config?.managed_config
                ?.special_subframe_pattern ?? 0
            }
            subframeAssignment={
              enbInfo.enb.enodeb_config?.managed_config?.subframe_assignment ??
              0
            }
          />
        ) : lteRanConfigs?.fdd_config ? (
          <EnodeConfigFdd
            earfcndl={enbInfo.enb.enodeb_config?.managed_config?.earfcndl ?? 0}
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
        value: enbInfo.enb.enodeb_config?.managed_config?.pci ?? '-',
      },
    ],
    [
      {
        category: 'TAC',
        value: enbInfo.enb.enodeb_config?.managed_config?.tac ?? '-',
      },
    ],
    [
      {
        category: 'Transmit',
        value: enbInfo.enb.enodeb_config?.managed_config?.transmit_enabled
          ? 'Enabled'
          : 'Disabled',
      },
    ],
  ];
  return <DataGrid data={managedConfig} testID="ran" />;
}

function EnodebUnmanagedRanConfig({enbInfo}: {enbInfo: EnodebInfo}) {
  const unmanagedConfig: DataRows[] = [
    [
      {
        category: 'eNodeB Externally Managed',
        value: 'True',
      },
    ],
    [
      {
        category: 'Cell ID',
        value: enbInfo.enb.enodeb_config?.unmanaged_config?.cell_id ?? '-',
      },
    ],
    [
      {
        category: 'TAC',
        value: enbInfo.enb.enodeb_config?.unmanaged_config?.tac ?? '-',
      },
    ],
    [
      {
        category: 'IP Address',
        value: enbInfo.enb.enodeb_config?.unmanaged_config?.ip_address ?? '-',
      },
    ],
  ];
  return <DataGrid data={unmanagedConfig} testID="ran" />;
}

function EnodebInfoConfig() {
  const ctx = useContext(EnodebContext);
  const params = useParams();
  const enodebSerial: string = nullthrows(params.enodebSerial);
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
