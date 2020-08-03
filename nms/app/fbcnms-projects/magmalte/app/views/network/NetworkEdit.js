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
import type {
  network,
  network_epc_configs,
  network_id,
  network_ran_configs,
} from '@fbcnms/magma-api';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '../../theme/design-system/DialogTitle';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import {NetworkEpcEdit} from './NetworkEpc';
import {NetworkInfoEdit} from './NetworkInfo';
import {NetworkRanEdit} from './NetworkRanConfig';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

const NETWORK_TITLE = 'Network';
const EPC_TITLE = 'Epc';
const RAN_TITLE = 'Ran';

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
}));

const EditTableType = {
  info: 0,
  epc: 1,
  ran: 2,
};

type EditProps = {
  editTable: $Keys<typeof EditTableType>,
  networkInfo: network,
  epcConfigs: network_epc_configs,
  lteRanConfigs: network_ran_configs,
  onSaveNetworkInfo: network => void,
  onSaveEpcConfigs: network_epc_configs => void,
  onSaveLteRanConfigs: network_ran_configs => void,
};

type DialogProps = {
  open: boolean,
  onClose: () => void,
  editProps?: EditProps,
};

type ButtonProps = {
  title: string,
  isLink: boolean,
  editProps?: EditProps,
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

function NetworkEditDialog({open, onClose, editProps}: DialogProps) {
  const classes = useStyles();
  const [networkId, setNetworkId] = useState<network_id>(
    editProps?.networkInfo?.id || '',
  );
  const [networkInfo, setNetworkInfo] = useState<network>({});
  const [epcConfigs, setEpcConfigs] = useState<network_epc_configs>({});

  const [tabPos, setTabPos] = React.useState(
    editProps ? EditTableType[editProps.editTable] : 0,
  );

  return (
    <Dialog data-testid="editDialog" open={open} fullWidth={true} maxWidth="sm">
      <DialogTitle
        label={editProps ? 'Edit Network Settings' : 'Add Network'}
        onClose={onClose}
      />
      <Tabs
        value={tabPos}
        onChange={(_, v) => setTabPos(v)}
        indicatorColor="primary"
        className={classes.tabBar}>
        <Tab key="network" data-testid="networkTab" label={NETWORK_TITLE} />;
        <Tab
          key="epc"
          data-testid="epcTab"
          disabled={networkId ? false : true}
          label={EPC_TITLE}
        />
        ;
        <Tab
          key="ran"
          data-testid="ranTab"
          disabled={networkId ? false : true}
          label={RAN_TITLE}
        />
        ;
      </Tabs>
      {tabPos === 0 && (
        <NetworkInfoEdit
          saveButtonTitle={editProps ? 'Save' : 'Save And Continue'}
          networkInfo={
            Object.keys(networkInfo).length != 0
              ? networkInfo
              : editProps?.networkInfo
          }
          onClose={onClose}
          onSave={(networkInfo: network) => {
            setNetworkInfo(networkInfo);
            if (editProps) {
              editProps.onSaveNetworkInfo(networkInfo);
              onClose();
            } else {
              setNetworkId(networkInfo.id);
              setTabPos(tabPos + 1);
            }
          }}
        />
      )}
      {tabPos === 1 && (
        <NetworkEpcEdit
          saveButtonTitle={editProps ? 'Save' : 'Save And Continue'}
          networkId={networkId}
          epcConfigs={
            Object.keys(epcConfigs).length != 0
              ? epcConfigs
              : editProps?.epcConfigs
          }
          onClose={onClose}
          onSave={epcConfigs => {
            setEpcConfigs(epcConfigs);
            if (editProps) {
              editProps.onSaveEpcConfigs(epcConfigs);
              onClose();
            } else {
              setTabPos(tabPos + 1);
            }
          }}
        />
      )}
      {tabPos === 2 && (
        <NetworkRanEdit
          saveButtonTitle={editProps ? 'Save' : 'Save And Add Network'}
          networkId={networkId}
          lteRanConfigs={editProps?.lteRanConfigs}
          onClose={onClose}
          onSave={lteRanConfigs => {
            editProps?.onSaveLteRanConfigs(lteRanConfigs);
            onClose();
          }}
        />
      )}
    </Dialog>
  );
}
