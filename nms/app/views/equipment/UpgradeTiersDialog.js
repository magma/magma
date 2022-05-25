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
import ActionTable from '../../components/ActionTable';
import Button from '@material-ui/core/Button';
import Collapse from '@material-ui/core/Collapse';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import DialogTitle from '../../theme/design-system/DialogTitle';
import Divider from '@material-ui/core/Divider';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import FormLabel from '@material-ui/core/FormLabel';
import GatewayTierContext from '../../components/context/GatewayTierContext';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '../../theme/design-system/Text';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {AutoCompleteEditComponent} from '../../components/ActionTable';

import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';

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
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '50%',
    fullWidth: true,
  },
}));

export default function UpgradeButton() {
  const classes = useStyles();
  const [open, setOpen] = useState(false);

  return (
    <>
      <UpgradeDialog open={open} onClose={() => setOpen(false)} />
      <Button
        variant="text"
        onClick={() => setOpen(true)}
        className={classes.appBarBtnSecondary}>
        {'Upgrade'}
      </Button>
    </>
  );
}

type DialogProps = {
  open: boolean,
  onClose: () => void,
};

function UpgradeDialog(props: DialogProps) {
  const classes = useStyles();
  const [tabPos, setTabPos] = useState(0);

  return (
    <Dialog
      data-testid="editDialog"
      open={props.open}
      fullWidth={true}
      maxWidth="lg">
      <DialogTitle label={'Upgrade Tiers'} onClose={props.onClose} />
      <Tabs
        value={tabPos}
        onChange={(_, v) => setTabPos(v)}
        indicatorColor="primary"
        className={classes.tabBar}>
        <Tab key="upgradeTiers" label={'Tiers'} />; ;
        <Tab key="supportedVersions" label={'Supported Versions'} />
      </Tabs>
      {tabPos === 0 && <UpgradeDetails {...props} />}
      {tabPos === 1 && <SupportedVersions {...props} />}
    </Dialog>
  );
}

function UpgradeDetails(props: DialogProps) {
  const ctx = useContext(GatewayTierContext);
  const [error, setError] = useState('');
  const [updatedTierEntries, setUpdatedTierEntries] = useState(new Set());
  const [removedTierEntries, setRemovedTierEntries] = useState(new Set());
  const [tierEntries, setTierEntries] = useState(
    Object.keys(ctx.state.tiers).map(tierId => ({
      id: tierId,
      name: ctx.state.tiers[tierId].name,
      version: ctx.state.tiers[tierId].version,
    })),
  );
  const saveTier = async () => {
    for (const tier of tierEntries) {
      if (!updatedTierEntries.has(tier.id)) {
        continue;
      }
      try {
        await ctx.setState(tier.id, {
          ...(ctx.state.tiers[tier.id] ?? {gateways: [], images: []}),
          ...tier,
        });
      } catch (e) {
        const errMsg = e.response?.data?.message ?? e.message ?? e;
        setError('error saving ' + tier.id + ' : ' + errMsg);
        return;
      }
    }

    for (const tierId of removedTierEntries) {
      try {
        await ctx.setState(tierId);
      } catch (e) {
        const errMsg = e.response?.data?.message ?? e.message ?? e;
        setError('error removing ' + tierId + ' : ' + errMsg);
        return;
      }
    }
    props.onClose();
  };

  return (
    <>
      <DialogContent>
        {error !== '' && <FormLabel error>{error}</FormLabel>}

        <ActionTable
          data={tierEntries}
          columns={[
            {
              title: 'Tier ID',
              field: 'id',
              editable: 'onAdd',
              editComponent: props => (
                <OutlinedInput
                  variant="outlined"
                  type="text"
                  value={props.value}
                  onChange={e => props.onChange(e.target.value)}
                />
              ),
            },
            {
              title: 'Tier Name',
              field: 'name',
              editComponent: props => (
                <OutlinedInput
                  variant="outlined"
                  type="text"
                  value={props.value}
                  onChange={e => props.onChange(e.target.value)}
                />
              ),
            },
            {
              title: 'Software Version',
              field: 'version',
              editComponent: props => (
                <AutoCompleteEditComponent
                  {...props}
                  value={props.value}
                  content={ctx.state.supportedVersions}
                  onChange={value => props.onChange(value)}
                />
              ),
            },
          ]}
          options={{
            actionsColumnIndex: -1,
            pageSizeOptions: [10, 15],
          }}
          editable={{
            onRowAdd: newData =>
              new Promise((resolve, _) => {
                setTierEntries([...tierEntries, newData]);
                setUpdatedTierEntries(
                  new Set([...updatedTierEntries, newData.id]),
                );
                resolve();
              }),
            onRowUpdate: (newData, oldData) =>
              new Promise((resolve, _) => {
                const dataUpdate = [...tierEntries];
                const index = oldData.tableData.id;
                dataUpdate[index] = newData;
                setTierEntries([...dataUpdate]);
                setUpdatedTierEntries(
                  new Set([...updatedTierEntries, newData.id]),
                );
                resolve();
              }),
            onRowDelete: oldData =>
              new Promise(resolve => {
                const dataDelete = [...tierEntries];
                const index = oldData.tableData.id;
                dataDelete.splice(index, 1);
                setTierEntries([...dataDelete]);
                setRemovedTierEntries(
                  new Set([...removedTierEntries, oldData.id]),
                );
                resolve();
              }),
          }}
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}> Cancel </Button>
        <Button onClick={saveTier}> Save </Button>
      </DialogActions>
    </>
  );
}

function SupportedVersions() {
  const [open, setOpen] = useState(false);
  const ctx = useContext(GatewayTierContext);
  const newVersions = ctx.state.supportedVersions.slice(0, 10);
  const olderVersion = ctx.state.supportedVersions.slice(10);
  return newVersions.length > 0 ? (
    <DialogContent>
      <List component={Paper}>
        {newVersions.map((version, i) => (
          <>
            <ListItem key={version}>
              {version}
              {i === 0 && <b> (Newest Version)</b>}
            </ListItem>
            {i < newVersions.length - 1 && <Divider />}
          </>
        ))}
      </List>
      {olderVersion.length > 0 && (
        <List component={Paper}>
          <ListItem button onClick={() => setOpen(!open)}>
            {open ? <ExpandLess /> : <ExpandMore />}
            Older Versions
          </ListItem>
          <Collapse key="olderVersions" in={open} timeout="auto" unmountOnExit>
            {olderVersion.map((version, i) => (
              <>
                <ListItem key={version}>{version}</ListItem>
                {i < olderVersion.length - 1 && <Divider />}
              </>
            ))}
          </Collapse>
        </List>
      )}
    </DialogContent>
  ) : (
    <Text>no supported versions found</Text>
  );
}
