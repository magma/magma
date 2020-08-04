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
import type {KPIRows} from '../../components/KPIGrid';
import type {network_ran_configs} from '@fbcnms/magma-api';

import Button from '@material-ui/core/Button';
import CardHeader from '@material-ui/core/CardHeader';
import Collapse from '@material-ui/core/Collapse';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import Divider from '@material-ui/core/Divider';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import FddConfig from './NetworkRanFddConfig';
import FormLabel from '@material-ui/core/FormLabel';
import Grid from '@material-ui/core/Grid';
import KPIGrid from '../../components/KPIGrid';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Select from '@material-ui/core/Select';
import TddConfig from './NetworkRanTddConfig';

import {AltFormField, FormDivider} from '../../components/FormField';
import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  list: {
    padding: 0,
  },
  kpiLabel: {
    color: colors.primary.comet,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  kpiValue: {
    color: colors.primary.brightGray,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    width: '100%',
  },
  kpiBox: {
    width: '100%',
    padding: 0,
    '& > div': {
      width: '100%',
    },
  },
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '50%',
    fullWidth: true,
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
}));

type Props = {
  lteRanConfigs: network_ran_configs,
};

export default function NetworkRan(props: Props) {
  const classes = useStyles();
  const [open, setOpen] = React.useState(true);

  const tdd: KPIRows[] = [
    [
      {
        category: 'EARFCNDL',
        value: props.lteRanConfigs?.tdd_config?.earfcndl || '-',
      },
    ],
    [
      {
        category: 'Special Subframe Pattern',
        value: props.lteRanConfigs.tdd_config?.special_subframe_pattern || '-',
      },
    ],
    [
      {
        category: 'Subframe Assignment',
        value: props.lteRanConfigs?.tdd_config?.subframe_assignment || '-',
      },
    ],
  ];

  const fdd: KPIRows[] = [
    [
      {
        category: 'EARFCNDL',
        value: props.lteRanConfigs?.fdd_config?.earfcndl || '-',
      },
    ],
    [
      {
        category: 'EARFCNUL',
        value: props.lteRanConfigs?.fdd_config?.earfcnul || '-',
      },
    ],
  ];

  return (
    <Grid item xs={12}>
      <List
        component={Paper}
        elevation={0}
        data-testid="ran"
        className={classes.list}>
        {/* TODO: Temporary fix until Data Grid is made */}
        <ListItem>
          <ListItemText
            primary={'Bandwidth'}
            secondary={props.lteRanConfigs?.bandwidth_mhz}
          />
        </ListItem>
        <Divider />
        {props.lteRanConfigs?.tdd_config && (
          <List key="tddConfigs" className={classes.list}>
            <ListItem button onClick={() => setOpen(!open)}>
              <CardHeader
                title="RAN Config"
                className={classes.kpiBox}
                subheader="TDD"
                titleTypographyProps={{
                  variant: 'body3',
                  className: classes.kpiLabel,
                  title: 'RAN Config',
                }}
                subheaderTypographyProps={{
                  variant: 'body1',
                  className: classes.kpiValue,
                  title: 'TDD',
                }}
              />
              {open ? <ExpandLess /> : <ExpandMore />}
            </ListItem>
            <Divider />
            <Collapse key="tdd" in={open} timeout="auto" unmountOnExit>
              <KPIGrid data={tdd} />
            </Collapse>
          </List>
        )}
        {props.lteRanConfigs?.fdd_config && (
          <List key="fddConfigs" className={classes.list}>
            <ListItem button onClick={() => setOpen(!open)}>
              <CardHeader
                title="RAN Config"
                className={classes.kpiBox}
                subheader="FDD"
                titleTypographyProps={{
                  variant: 'body3',
                  className: classes.kpiLabel,
                  title: 'RAN Config',
                }}
                subheaderTypographyProps={{
                  variant: 'body1',
                  className: classes.kpiValue,
                  title: 'FDD',
                }}
              />
              {open ? <ExpandLess /> : <ExpandMore />}
            </ListItem>
            <Divider />
            <Collapse key="fdd" in={open} timeout="auto" unmountOnExit>
              <KPIGrid data={fdd} />
            </Collapse>
          </List>
        )}
      </List>
    </Grid>
  );
}

type EditProps = {
  saveButtonTitle: string,
  networkId: string,
  lteRanConfigs: ?network_ran_configs,
  onClose: () => void,
  onSave: network_ran_configs => void,
};
type BandType = 'tdd' | 'fdd';
const ValidBandwidths = [3, 5, 10, 15, 20];

export function NetworkRanEdit(props: EditProps) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const [error, setError] = useState('');
  const [bandType, setBandType] = useState<BandType>('tdd');
  const defaultTddConfig = {
    earfcndl: 0,
    special_subframe_pattern: 0,
    subframe_assignment: 0,
  };
  const defaulFddConfig = {
    earfcndl: 0,
    earfcnul: 0,
  };
  const [lteRanConfigs, setLteRanConfigs] = useState(
    props?.lteRanConfigs || {
      bandwidth_mhz: 20,
      fdd_config: undefined,
      tdd_config: defaultTddConfig,
    },
  );

  const onSave = async () => {
    const config: network_ran_configs = {
      ...lteRanConfigs,
    };
    if (bandType === 'tdd') {
      config.fdd_config = undefined;
    } else {
      config.tdd_config = undefined;
    }
    try {
      await MagmaV1API.putLteByNetworkIdCellularRan({
        networkId: props.networkId,
        config: config,
      });
      enqueueSnackbar('RAN configs saved successfully', {variant: 'success'});
      props.onSave(config);
    } catch (e) {
      setError(e.response.data?.message ?? e.message);
    }
  };

  return (
    <>
      <DialogContent data-testid="networkRanEdit">
        {error !== '' && (
          <AltFormField label={''}>
            <FormLabel error>{error}</FormLabel>
          </AltFormField>
        )}
        <List>
          <AltFormField label={'Bandwidth'}>
            <Select
              variant={'outlined'}
              fullWidth={true}
              value={lteRanConfigs.bandwidth_mhz}
              onChange={({target}) => {
                if (
                  target.value === 3 ||
                  target.value === 5 ||
                  target.value === 10 ||
                  target.value === 15 ||
                  target.value === 20
                ) {
                  setLteRanConfigs({
                    ...lteRanConfigs,
                    bandwidth_mhz: target.value,
                  });
                }
              }}
              input={<OutlinedInput fullWidth={true} id="bandwidth" />}>
              {ValidBandwidths.map((k: number, idx: number) => (
                <MenuItem key={idx} value={k}>
                  <ListItemText primary={k} />
                </MenuItem>
              ))}
            </Select>
          </AltFormField>
          <AltFormField label={'Band Type'}>
            <Select
              variant={'outlined'}
              fullWidth={true}
              value={bandType}
              onChange={({target}) => {
                if (target.value === 'fdd') {
                  setLteRanConfigs({
                    fdd_config: defaulFddConfig,
                    ...lteRanConfigs,
                  });
                  setBandType('fdd');
                } else {
                  setLteRanConfigs({
                    tdd_config: defaultTddConfig,
                    ...lteRanConfigs,
                  });
                  setBandType(target.value === 'tdd' ? 'tdd' : 'fdd');
                }
              }}
              input={<OutlinedInput fullWidth={true} id="bandType" />}>
              <MenuItem value={'tdd'}>
                <ListItemText primary={'TDD'} />
              </MenuItem>
              <MenuItem value={'fdd'}>
                <ListItemText primary={'FDD'} />
              </MenuItem>
            </Select>
          </AltFormField>
          <FormDivider />
          {bandType === 'tdd' && (
            <TddConfig
              lteRanConfigs={lteRanConfigs}
              setLteRanConfigs={setLteRanConfigs}
            />
          )}
          {bandType === 'fdd' && (
            <FddConfig
              lteRanConfigs={lteRanConfigs}
              setLteRanConfigs={setLteRanConfigs}
            />
          )}
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          Cancel
        </Button>
        <Button onClick={onSave} variant="contained" color="primary">
          {props.saveButtonTitle}
        </Button>
      </DialogActions>
    </>
  );
}
