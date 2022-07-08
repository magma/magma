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
 */
import type {
  GatewayPoolRecordsType,
  gatewayPoolsStateType,
} from '../../components/context/GatewayPoolsContext';
import type {MutableCellularGatewayPool} from '../../../generated-ts';

import Button from '@material-ui/core/Button';
import ConfigEdit from './GatewayPoolConfigEdit';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '../../theme/design-system/DialogTitle';
import GatewayEdit from './GatewayPoolGatewaysEdit';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';

import {
  DEFAULT_GW_POOL_CONFIG,
  DEFAULT_GW_PRIMARY_CONFIG,
  DEFAULT_GW_SECONDARY_CONFIG,
} from '../../components/GatewayUtils';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';

const CONFIG_TITLE = 'Config';
const GATEWAY_PRIMARY_TITLE = 'Gateway Primary';
const GATEWAY_SECONDARY_TITLE = 'Gateway Secondary';

const useStyles = makeStyles({
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
});
type DialogProps = {
  open: boolean;
  onClose: () => void;
  pool?: gatewayPoolsStateType;
  isAdd: boolean;
};

type ButtonProps = {
  title: string;
  isLink: boolean;
};

export type GatewayPoolEditProps = {
  gwPool: MutableCellularGatewayPool;
  isPrimary?: boolean;
  gatewayPrimary: Array<GatewayPoolRecordsType>;
  gatewaySecondary: Array<GatewayPoolRecordsType>;
  onRecordChange?: (gateways: Array<GatewayPoolRecordsType>) => void;
  onClose: () => void;
  onSave: (pool: MutableCellularGatewayPool) => void;
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
  const [gwPool, setGwPool] = useState<MutableCellularGatewayPool>(
    (props.pool
      ? {...props.pool?.gatewayPool, gateway_ids: undefined}
      : {}) as MutableCellularGatewayPool,
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
        ? ({
            ...props.pool?.gatewayPool,
            gateway_ids: undefined,
          } as MutableCellularGatewayPool)
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
        onChange={(_, v) => setTabPos(v as number)}
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
          onSave={(pool: MutableCellularGatewayPool) => {
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
