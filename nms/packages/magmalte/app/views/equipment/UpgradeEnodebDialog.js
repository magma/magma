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

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '../../theme/design-system/DialogTitle';
import EnodebContext from '../../components/context/EnodebContext';
import FormLabel from '@material-ui/core/FormLabel';
import List from '@material-ui/core/List';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';

import {AltFormField} from '../../components/FormField';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {RunGatewayCommands} from '../../state/lte/EquipmentState';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';

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
    color: colors.primary.white,
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
}));

type DialogProps = {
  open: boolean,
  onClose: () => void,
};

export default function UpgradeEnodebButton() {
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
      <UpgradeEnodebDialog open={open} onClose={handleClose} />
      <Button
        variant="text"
        className={classes.appBarBtnSecondary}
        onClick={handleClickOpen}>
        {'Upgrade'}
      </Button>
    </>
  );
}

function UpgradeEnodebDialog(props: DialogProps) {
  const {match} = useRouter();
  const enodebSerial: string = nullthrows(match.params.enodebSerial);

  return (
    <Dialog data-testid="editDialog" open={props.open} fullWidth={true} maxWidth="sm">
      <DialogTitle
        label={`Upgrade eNodeB/${enodebSerial}`}
        onClose={props.onClose}
      />
      <ConfigUpgrade onClose={props.onClose} />
    </Dialog>
  );
}

type UpgradeData = {
  name: string,
  url: string,
  md5: string,
  size: string,
};

const ConfigUpgrade = withAlert(ConfigUpgradeInternal);

function ConfigUpgradeInternal(props: WithAlert & {onClose: () => void}) {
  const [error, setError] = useState('');

  const {match} = useRouter();
  const ctx = useContext(EnodebContext);
  const networkId: string = nullthrows(match.params.networkId);
  const enodebSerial: string = nullthrows(match.params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];
  const gatewayId = enbInfo?.enb_state?.reporting_gateway_id;
  const enqueueSnackbar = useEnqueueSnackbar();

  const [upgradeData, setUpgradeData] = useState<UpgradeData>({});
  const handleUpgradeDataChange = (key: string, val) =>
    setUpgradeData({...upgradeData, [key]: val});

  const upgrade = (name, url, md5, size) => {
    if (gatewayId == null) {
      enqueueSnackbar(
        'Unable to trigger Upgrade, reporting gateway not found',
        {variant: 'error'},
      );
      return;
    }
    props
      .confirm(`Are you sure you want to upgrade the ${enodebSerial}?`)
      .then(async confirmed => {
        if (!confirmed) {
          return;
        }
        const params = {
          command: 'download_enodeb',
          params: {
            shell_params: {
              [enodebSerial]: {
                url: url,
                user_name: 'admin',
                password: 'admin',
                file_size: size,
                target_file_name: name,
                md5: md5,
              },
            },
          },
        };
        try {
          await RunGatewayCommands({
            networkId,
            gatewayId,
            command: 'generic',
            params,
          });
          enqueueSnackbar('eNodeb upgrade command is triggered successfully', {
            variant: 'success',
          });
        } catch (e) {
          enqueueSnackbar(e.response?.data?.message ?? e.message, {
            variant: 'error',
          });
        }
      });
  };

  const onUpgrade = async () => {
    setError('');
    if (upgradeData.name == '') {
      setError('Filename is empty!');
      return;
    }
    if (upgradeData.url == '') {
      setError('Url is empty!');
      return;
    }
    if (upgradeData.md5 == '') {
      setError('MD5 is empty!');
      return;
    }
    const val = parseInt(upgradeData.size, 10);
    if (upgradeData.size == '' || isNaN(val)) {
      setError('Size is empty!');
      return;
    }

    upgrade(
      upgradeData.name,
      upgradeData.url,
      upgradeData.md5,
      upgradeData.size,
    );
  };

  return (
    <>
      <DialogContent data-testid="configUpgrade">
        <List>
          {error !== '' && (
            <AltFormField label={''}>
              <FormLabel data-testid="configUpgradeError" error>
                {error}
              </FormLabel>
            </AltFormField>
          )}
          <AltFormField label={'Filename'}>
            <OutlinedInput
              data-testid="filename"
              placeholder="Enter Filename Ex: xxx.img"
              fullWidth={true}
              value={upgradeData.name}
              onChange={({target}) =>
                handleUpgradeDataChange('name', target.value)
              }
            />
          </AltFormField>
          <AltFormField label={'Url'}>
            <OutlinedInput
              data-testid="url"
              placeholder="Enter Url Ex: http://xxx/xxx.img"
              fullWidth={true}
              value={upgradeData.url}
              onChange={({target}) =>
                handleUpgradeDataChange('url', target.value)
              }
            />
          </AltFormField>
          <AltFormField label={'MD5'}>
            <OutlinedInput
              data-testid="md5"
              placeholder="Enter MD5 Ex: 7c65f88108cc554593a163043b845805"
              fullWidth={true}
              value={upgradeData.md5}
              onChange={({target}) =>
                handleUpgradeDataChange('md5', target.value)
              }
            />
          </AltFormField>
          <AltFormField label={'Size'}>
            <OutlinedInput
              data-testid="size"
              placeholder="Enter Size(Bytes) Ex: 9574368"
              fullWidth={true}
              value={upgradeData.size}
              onChange={({target}) =>
                handleUpgradeDataChange('size', target.value)
              }
            />
          </AltFormField>
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          Cancel
        </Button>
        <Button onClick={onUpgrade} variant="contained" color="primary">
          {'Upgrade'}
        </Button>
      </DialogActions>
    </>
  );
}
