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
  FegLteNetwork,
  LteNetwork,
  NetworkEpcConfigs,
} from '../../../generated';

import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogTitle from '../../theme/design-system/DialogTitle';
import LteNetworkContext from '../../context/LteNetworkContext';
import React from 'react';
import Tab from '@mui/material/Tab';
import Tabs from '@mui/material/Tabs';

import {NetworkEpcEdit} from './NetworkEpc';
import {NetworkFederationEdit} from './NetworkFederationConfig';
import {NetworkInfoEdit} from './NetworkInfo';
import {NetworkRanEdit} from './NetworkRanConfig';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@mui/styles';
import {useContext, useEffect, useState} from 'react';

const NETWORK_TITLE = 'Network';
const FEDERATION_TITLE = 'Federation';
const EPC_TITLE = 'Epc';
const RAN_TITLE = 'Ran';

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
  },
});

const LTE_TABS = {
  info: 0,
  feg: -1,
  epc: 1,
  ran: 2,
};

const FEG_TABS = {
  info: 0,
  feg: 1,
  epc: 2,
  ran: 3,
};

type EditProps = {
  editTable: keyof typeof LTE_TABS & keyof typeof FEG_TABS;
};

type DialogProps = {
  open: boolean;
  onClose: () => void;
  editProps?: EditProps;
};

type ButtonProps = {
  title: string;
  isLink: boolean;
  editProps?: EditProps;
};

export default function AddEditNetworkButton(props: ButtonProps) {
  const classes = useStyles();
  const [open, setOpen] = React.useState(false);

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <>
      <NetworkEditDialog
        open={open}
        onClose={handleClose}
        editProps={props.editProps}
      />
      {props.isLink ? (
        <Button
          data-testid={(props.editProps?.editTable ?? '') + 'EditButton'}
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

export function NetworkEditDialog(props: DialogProps) {
  const {open, editProps} = props;
  const classes = useStyles();
  const ctx = useContext(LteNetworkContext);

  const [lteNetwork, setLteNetwork] = useState<
    Partial<LteNetwork & FegLteNetwork>
  >({});
  //  eslint-disable-next-line @typescript-eslint/ban-types
  const [epcConfigs, setEpcConfigs] = useState<NetworkEpcConfigs | {}>({});

  const lteRanConfigs = editProps ? ctx.state.cellular?.ran : undefined;

  const [tabPos, setTabPos] = React.useState<number>(0);

  useEffect(() => {
    if (editProps) {
      const network = ctx.state;
      setLteNetwork(network);
      setEpcConfigs(network.cellular?.epc ?? {});
      setTabPos(
        network.federation
          ? FEG_TABS[editProps.editTable]
          : LTE_TABS[editProps.editTable],
      );
    } else {
      setLteNetwork({});
      setEpcConfigs({});
      setTabPos(0);
    }
  }, [open, editProps, ctx.state]);

  const onClose = () => {
    props.onClose();
  };
  const isFegLet = !!lteNetwork.federation;
  const tabs = isFegLet ? FEG_TABS : LTE_TABS;

  return (
    <Dialog data-testid="editDialog" open={open} fullWidth={true} maxWidth="md">
      <DialogTitle
        label={editProps ? 'Edit Network Settings' : 'Add Network'}
        onClose={onClose}
      />
      <Tabs
        value={tabPos}
        onChange={(_, v: number) => setTabPos(v)}
        indicatorColor="primary"
        className={classes.tabBar}>
        <Tab key="network" data-testid="networkTab" label={NETWORK_TITLE} />
        {isFegLet && (
          <Tab
            key="federation"
            data-testid="federationTab"
            disabled={editProps ? false : true}
            label={FEDERATION_TITLE}
          />
        )}
        <Tab
          key="epc"
          data-testid="epcTab"
          disabled={editProps ? false : true}
          label={EPC_TITLE}
        />
        <Tab
          key="ran"
          data-testid="ranTab"
          disabled={editProps ? false : true}
          label={RAN_TITLE}
        />
      </Tabs>
      {tabPos === tabs.info && (
        <NetworkInfoEdit
          saveButtonTitle={editProps ? 'Save' : 'Save And Continue'}
          lteNetwork={lteNetwork}
          onClose={onClose}
          onSave={lteNetwork => {
            setLteNetwork(lteNetwork);
            if (editProps) {
              onClose();
            } else {
              setTabPos(isFegLet ? tabPos + 2 : tabPos + 1);
            }
          }}
        />
      )}
      {tabPos === tabs.feg && (
        <NetworkFederationEdit
          saveButtonTitle={editProps ? 'Save' : 'Save And Continue'}
          networkId={lteNetwork.id!}
          config={lteNetwork.federation!}
          onClose={onClose}
          onSave={federationConfigs => {
            setLteNetwork({...lteNetwork, federation: federationConfigs});
            if (editProps) {
              onClose();
            } else {
              setTabPos(tabPos + 1);
            }
          }}
        />
      )}
      {tabPos === tabs.epc && (
        <NetworkEpcEdit
          saveButtonTitle={editProps ? 'Save' : 'Save And Continue'}
          networkId={(lteNetwork as LteNetwork).id}
          epcConfigs={epcConfigs as NetworkEpcConfigs}
          onClose={onClose}
          onSave={epcConfigs => {
            setEpcConfigs(epcConfigs);
            if (editProps) {
              onClose();
            } else {
              setTabPos(tabPos + 1);
            }
          }}
        />
      )}
      {tabPos === tabs.ran && (
        <NetworkRanEdit
          saveButtonTitle={editProps ? 'Save' : 'Save And Add Network'}
          networkId={(lteNetwork as LteNetwork).id}
          lteRanConfigs={lteRanConfigs}
          onClose={onClose}
          onSave={onClose}
        />
      )}
    </Dialog>
  );
}
