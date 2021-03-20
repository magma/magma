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
 * @flow strict-local
 * @format
 */
/*[object Object]*/
import type {
  GatewayPoolRecordsType,
  gatewayPoolsStateType,
} from '../../components/context/GatewayPoolsContext';
import type {mutable_cellular_gateway_pool} from '@fbcnms/magma-api';

import AddCircleOutline from '@material-ui/icons/AddCircleOutline';
import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '../../theme/design-system/DialogTitle';
import FormLabel from '@material-ui/core/FormLabel';
import GatewayContext from '../../components/context/GatewayContext';
import GatewayPoolsContext from '../../components/context/GatewayPoolsContext';
import IconButton from '@material-ui/core/IconButton';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';

import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Select from '@material-ui/core/Select';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';

import {AltFormField} from '../../components/FormField';
import {
  DEFAULT_GW_PRIMARY_CONFIG,
  DEFAULT_GW_SECONDARY_CONFIG,
} from '../../components/GatewayUtils';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useEffect, useState} from 'react';

const CONFIG_TITLE = 'Config';
const GATEWAY_PRIMARY_TITLE = 'Gateway Primary';
const GATEWAY_SECONDARY_TITLE = 'Gateway Secondary';

const DEFAULT_GW_POOL_CONFIG = {
  config: {mme_group_id: 1},
  gateway_pool_id: '',
  gateway_pool_name: '',
};

const useStyles = makeStyles(_ => ({
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
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    color: colors.primary.white,
  },
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '50%',
    fullWidth: true,
  },
}));
type DialogProps = {
  open: boolean,
  onClose: () => void,
  pool?: gatewayPoolsStateType,
  isAdd: boolean,
};

type ButtonProps = {
  title: string,
  isLink: boolean,
};

export default function AddEditGatewayPoolButton(props: ButtonProps) {
  const classes = useStyles();
  const [open, setOpen] = useState(false);

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <>
      <GatewayPoolEditDialog open={open} onClose={handleClose} isAdd={false} />
      {props.isLink ? (
        <Button
          data-testid={'EditButton'}
          component="button"
          variant="text"
          onClick={handleClickOpen}>
          {props.title}
        </Button>
      ) : (
        <Button
          variant="text"
          className={classes.appBarBtn}
          onClick={handleClickOpen}>
          {props.title}
        </Button>
      )}
    </>
  );
}

export function GatewayPoolEditDialog(props: DialogProps) {
  const {open} = props;
  const classes = useStyles();
  const [gwPool, setGwPool] = useState<mutable_cellular_gateway_pool>(
    props.pool ? {...props.pool?.gatewayPool, gateway_ids: undefined} : {},
  );
  const [gwPrimary, setGwPrimary] = useState<Array<GatewayPoolRecordsType>>(
    (props.pool?.gatewayPoolRecords || []).filter(
      gw => gw.mme_relative_capacity === 255,
    ).length > 0
      ? (props.pool?.gatewayPoolRecords || []).filter(
          gw => gw.mme_relative_capacity === 255,
        )
      : [DEFAULT_GW_PRIMARY_CONFIG],
  );
  const [gwSecondary, setGwSecondary] = useState<Array<GatewayPoolRecordsType>>(
    (props.pool?.gatewayPoolRecords || []).filter(
      gw => gw.mme_relative_capacity === 1,
    ).length > 0
      ? (props.pool?.gatewayPoolRecords || []).filter(
          gw => gw.mme_relative_capacity === 1,
        )
      : [DEFAULT_GW_SECONDARY_CONFIG],
  );

  const [tabPos, setTabPos] = useState(0);

  const onClose = () => {
    props.onClose();
  };

  useEffect(() => {
    setTabPos(0);
    setGwPool(
      props.pool
        ? {...props.pool?.gatewayPool, gateway_ids: undefined}
        : DEFAULT_GW_POOL_CONFIG,
    );
    const primary =
      props.pool?.gatewayPoolRecords?.filter(
        gw => gw.mme_relative_capacity === 255,
      ) || [];
    setGwPrimary(
      primary.length > 0 ? primary : [{...DEFAULT_GW_PRIMARY_CONFIG}],
    );
    const secondary =
      props.pool?.gatewayPoolRecords?.filter(
        gw => gw.mme_relative_capacity === 1,
      ) || [];
    setGwSecondary(
      secondary.length > 0 ? secondary : [{...DEFAULT_GW_SECONDARY_CONFIG}],
    );
  }, [props.pool, props.open]);

  return (
    <Dialog
      data-testid="gatewayPoolEditDialog"
      open={open}
      fullWidth={true}
      maxWidth="md">
      <DialogTitle
        label={props.isAdd ? 'Edit Gateway Pool' : 'Add New Gateway Pool'}
        onClose={onClose}
      />
      <Tabs
        value={tabPos}
        onChange={(_, v) => setTabPos(v)}
        indicatorColor="primary"
        className={classes.tabBar}>
        <Tab key="config" data-testid="configTab" label={CONFIG_TITLE} />; ;
        <Tab
          key="gwPrimary"
          data-testid="gwPrimaryTab"
          label={GATEWAY_PRIMARY_TITLE}
        />
        <Tab
          key="gwSecondary"
          data-testid="gwSecondaryTab"
          label={GATEWAY_SECONDARY_TITLE}
        />
      </Tabs>
      {tabPos === 0 && (
        <ConfigEdit
          gwPool={gwPool}
          gatewayPrimary={gwPrimary}
          gatewaySecondary={gwSecondary}
          onClose={onClose}
          onSave={(pool: mutable_cellular_gateway_pool) => {
            setGwPool({...pool});
            setTabPos(tabPos + 1);
          }}
        />
      )}
      {tabPos === 1 && (
        <GatewayEdit
          isPrimary={true}
          gwPool={gwPool}
          onClose={onClose}
          onRecordChange={gateways => {
            setGwPrimary([...gateways]);
          }}
          gatewayPrimary={gwPrimary}
          gatewaySecondary={gwSecondary}
          onSave={() => {
            setTabPos(tabPos + 1);
          }}
        />
      )}
      {tabPos === 2 && (
        <GatewayEdit
          isPrimary={false}
          gwPool={gwPool}
          onClose={onClose}
          onRecordChange={gateways => {
            setGwSecondary([...gateways]);
          }}
          gatewayPrimary={gwPrimary}
          gatewaySecondary={gwSecondary}
          onSave={() => {
            onClose();
          }}
        />
      )}
    </Dialog>
  );
}
type Props = {
  gwPool: mutable_cellular_gateway_pool,
  isPrimary?: boolean,
  gatewayPrimary: Array<GatewayPoolRecordsType>,
  gatewaySecondary: Array<GatewayPoolRecordsType>,
  onRecordChange?: (gateways: Array<GatewayPoolRecordsType>) => void,
  onClose: () => void,
  onSave: mutable_cellular_gateway_pool => void,
};

export function ConfigEdit(props: Props) {
  const [error, setError] = useState('');
  const ctx = useContext(GatewayPoolsContext);
  const [gwPool, setGwPool] = useState<mutable_cellular_gateway_pool>(
    Object.keys(props.gwPool || {}).length > 0
      ? props.gwPool
      : DEFAULT_GW_POOL_CONFIG,
  );
  const handleGwPoolConfigChange = (value: number) => {
    const newConfig = {
      ...gwPool,
      config: {...gwPool.config, ['mme_group_id']: value},
    };
    setGwPool(newConfig);
  };
  const onSave = async () => {
    try {
      await ctx.setState(gwPool.gateway_pool_id, gwPool);
      props.onSave(gwPool);
    } catch (e) {
      setError(e.response?.data?.message ?? e.message);
    }
  };
  const classes = useStyles();

  return (
    <>
      <DialogContent data-testid="configEdit">
        <List>
          {error !== '' && (
            <AltFormField label={''}>
              <FormLabel data-testid="configEditError" error>
                {error}
              </FormLabel>
            </AltFormField>
          )}
          <AltFormField label={'Name'}>
            <OutlinedInput
              data-testid="name"
              className={classes.input}
              placeholder="Enter Name"
              fullWidth={true}
              value={gwPool.gateway_pool_name}
              onChange={({target}) =>
                setGwPool({...gwPool, gateway_pool_name: target.value})
              }
            />
          </AltFormField>
          <AltFormField label={'ID'}>
            <OutlinedInput
              data-testid="poolId"
              className={classes.input}
              placeholder="Ex: pool1"
              fullWidth={true}
              value={gwPool.gateway_pool_id}
              readOnly={Object.keys(props.gwPool).length > 0 ? false : true}
              onChange={({target}) =>
                setGwPool({...gwPool, gateway_pool_id: target.value})
              }
            />
          </AltFormField>
          <AltFormField label={'MME Group ID'}>
            <OutlinedInput
              data-testid="mmeGroupId"
              className={classes.input}
              placeholder="Ex: 1"
              fullWidth={true}
              type="number"
              value={gwPool.config.mme_group_id}
              onChange={({target}) => {
                handleGwPoolConfigChange(parseInt(target.value));
              }}
            />
          </AltFormField>
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          Cancel
        </Button>
        <Button onClick={onSave} variant="contained" color="primary">
          {'Save And Continue'}
        </Button>
      </DialogActions>
    </>
  );
}

export function GatewayEdit(props: Props) {
  const [error, setError] = useState('');
  const ctx = useContext(GatewayPoolsContext);
  const gwCtx = useContext(GatewayContext);
  const [gwIds, _setGwIds] = useState(Object.keys(gwCtx.state));
  const [isPrimary, setIsPrimary] = useState(props.isPrimary || false);
  const [gwPool, setGwPool] = useState(
    props.gwPool || {...DEFAULT_GW_POOL_CONFIG, gateway_ids: []},
  );
  const [gateways, setGateways] = useState(
    isPrimary ? props.gatewayPrimary : props.gatewaySecondary,
  );

  useEffect(() => {
    setGwPool(props.gwPool || DEFAULT_GW_POOL_CONFIG);
  }, [props.gwPool]);

  const handleGwIdChange = (id: string, index: number) => {
    const newGwList = gateways;
    newGwList[index].gateway_id = id;
    props.onRecordChange?.([...newGwList]);
    setGateways([...newGwList]);
  };
  const handlePrimaryChange = (index: number, value: number, key) => {
    const newGwList = gateways;
    newGwList[index][key] = value;
    props.onRecordChange?.([...newGwList]);
    setGateways([...newGwList]);
  };

  const handleAddGw = () => {
    const newGwList = [
      ...gateways,
      isPrimary
        ? {
            gateway_id: '',
            gateway_pool_id: '',
            mme_code: 1,
            mme_relative_capacity: 255,
          }
        : {
            gateway_id: '',
            gateway_pool_id: '',
            mme_code: 1,
            mme_relative_capacity: 1,
          },
    ];

    setGateways([...newGwList]);
  };
  const deleteGateway = (gatewayId: string) => {
    const newGwList = isPrimary
      ? props.gatewayPrimary.filter(gw => gw.gateway_id !== gatewayId)
      : props.gatewaySecondary.filter(gw => gw.gateway_id !== gatewayId);
    if (newGwList) {
      props.onRecordChange?.([...newGwList]);
      setGateways([...newGwList]);
    }
  };

  useEffect(() => {
    setIsPrimary(props.isPrimary || false);
  }, [props.isPrimary]);

  const onSave = async () => {
    try {
      await ctx.setState(gwPool.gateway_pool_id, gwPool, [
        ...props.gatewayPrimary,
        ...props.gatewaySecondary,
      ]);
      props.onSave(gwPool);
    } catch (e) {
      setError(e.response?.data?.message ?? e.message);
    }
  };
  return (
    <>
      <DialogContent data-testid={`${isPrimary ? 'Primary' : 'Secondary'}Edit`}>
        <List>
          {error !== '' && (
            <AltFormField label={''}>
              <FormLabel data-testid="configEditError" error>
                {error}
              </FormLabel>
            </AltFormField>
          )}

          {gateways.length > 0 &&
            gateways.map((gw, index) => (
              <ListItem component={Paper}>
                <AltFormField
                  label={`${isPrimary ? 'Primary' : 'Secondary'} Gateway ID`}>
                  <Select
                    variant={'outlined'}
                    displayEmpty={true}
                    value={gw.gateway_id}
                    onChange={({target}) =>
                      handleGwIdChange(target.value, index)
                    }
                    input={
                      <OutlinedInput
                        data-testid={`gwId${
                          isPrimary ? 'Primary' : 'Secondary'
                        }`}
                        fullWidth={true}
                      />
                    }>
                    {gwIds.map(id => (
                      <MenuItem key={id} value={id}>
                        <ListItemText primary={id} />
                      </MenuItem>
                    ))}
                  </Select>
                </AltFormField>
                <AltFormField label={'MME Code'}>
                  <OutlinedInput
                    data-testid="mmeCode"
                    placeholder="Ex: 12020000261814C0021"
                    fullWidth={true}
                    type="number"
                    value={gw.mme_code}
                    onChange={({target}) => {
                      handlePrimaryChange(
                        index,
                        parseInt(target.value),
                        'mme_code',
                      );
                    }}
                  />
                </AltFormField>
                <AltFormField label={'MME Relative Capacity'}>
                  <OutlinedInput
                    data-testid="mmeCapacity"
                    placeholder="Enter Description"
                    fullWidth={true}
                    type="number"
                    value={gw.mme_relative_capacity}
                    onChange={({target}) => {
                      handlePrimaryChange(
                        index,
                        parseInt(target.value),
                        'mme_relative_capacity',
                      );
                    }}
                  />
                </AltFormField>
                <ListItemSecondaryAction>
                  <IconButton
                    edge="end"
                    aria-label="delete"
                    onClick={() => deleteGateway(gw.gateway_id)}>
                    <DeleteIcon />
                  </IconButton>
                </ListItemSecondaryAction>
              </ListItem>
            ))}
          <>
            Add New Gateway
            <IconButton
              data-testid="addGwButton"
              onClick={handleAddGw}
              disabled={isPrimary ? false : gateways.length > 0}>
              <AddCircleOutline />
            </IconButton>
          </>
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          Cancel
        </Button>
        <Button onClick={onSave} variant="contained" color="primary">
          {props.isPrimary ?? false ? 'Save And Continue' : 'Save'}
        </Button>
      </DialogActions>
    </>
  );
}
