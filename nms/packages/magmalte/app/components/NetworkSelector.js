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
 * @flow
 * @format
 */
import type {NetworkType} from '@fbcnms/types/network';

import AppContext from '@fbcnms/ui/context/AppContext';
import Divider from '@material-ui/core/Divider';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NetworkContext from './context/NetworkContext';
import Popout from '@fbcnms/ui/components/Popout';
import React from 'react';
import SettingsEthernetIcon from '@material-ui/icons/SettingsEthernet';
import Text from '@fbcnms/ui/components/design-system/Text';
import Tooltip from '@material-ui/core/Tooltip';
import classNames from 'classnames';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {LTE, coalesceNetworkType} from '@fbcnms/types/network';
import {NetworkEditDialog} from '../views/network/NetworkEdit';
import {colors} from '../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useContext, useEffect, useState} from 'react';

const useStyles = makeStyles(_ => ({
  button: {
    backgroundColor: colors.primary.white,
    width: '28px',
    height: '28px',
    fontSize: '28px',
    cursor: 'pointer',
    borderRadius: '100%',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: '20px',
    border: `1px solid ${colors.primary.white}`,
    '&:hover, &$openButton': {
      border: `1px solid ${colors.secondary.dodgerBlue}`,
    },
  },
  openButton: {},
  itemGutters: {
    '&&': {
      minWidth: '170px',
      borderRadius: '4px',
      padding: '8px 10px',
      '&:hover': {
        backgroundColor: 'rgba(145, 145, 145, 0.1)',
      },
    },
  },
  networksList: {
    '&&': {
      maxHeight: '400px',
      overflowY: 'auto',
      padding: '10px 5px',
    },
  },
  networkItemText: {
    fontSize: '12px',
    lineHeight: '16px',
  },
  selectedNetwork: {
    color: colors.secondary.dodgerBlue,
    fontSize: '20px',
  },
  selectedListItem: {
    '& $networkItemText': {
      color: colors.secondary.dodgerBlue,
    },
  },
  listItemRoot: {
    '&$selectedListItem': {
      backgroundColor: colors.primary.concrete,
    },
    '&:not(:last-child)': {
      marginBottom: '8px',
    },
  },
}));

const NetworkSelector = () => {
  const classes = useStyles();
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const appContext = useContext(AppContext);
  const [networkIds, setNetworkIds] = useState([]);
  const [networkType, setNetworkType] = useState<?NetworkType>(null);
  const [lastRefreshTime, setLastRefreshTime] = useState(new Date().getTime());
  const [isNetworkAddOpen, setNetworkAddOpen] = useState(false);
  const {networkId} = useContext(NetworkContext);

  useMagmaAPI(
    MagmaV1API.getNetworks,
    {},
    useCallback(resp => setNetworkIds(resp), []),
    lastRefreshTime,
  );

  useEffect(() => {
    const fetchNetworkType = async () => {
      if (networkId) {
        const networkType = await MagmaV1API.getNetworksByNetworkIdType({
          networkId,
        });
        setNetworkType(coalesceNetworkType(networkId, networkType));
      }
    };

    fetchNetworkType();
  }, [networkId]);

  if (!networkId) {
    return null;
  }

  return (
    <>
      <NetworkEditDialog
        open={isNetworkAddOpen}
        onClose={() => {
          setNetworkAddOpen(false);
          setLastRefreshTime(new Date().getTime());
        }}
      />
      <Popout
        open={isMenuOpen}
        content={
          <List component="nav" className={classes.networksList}>
            {networkIds.map(id => (
              <ListItem
                key={id}
                selected={id === networkId}
                classes={{
                  root: classes.listItemRoot,
                  gutters: classes.itemGutters,
                  selected: classes.selectedListItem,
                }}
                button
                component="a"
                href={`/nms/${id}`}>
                <Text className={classes.networkItemText}>{id}</Text>
              </ListItem>
            ))}
            {appContext.user.isSuperUser && networkType === LTE && (
              <>
                <Divider />
                <ListItem
                  key="create_network"
                  classes={{
                    root: classes.listItemRoot,
                    gutters: classes.itemGutters,
                  }}
                  button
                  component="a"
                  onClick={() => setNetworkAddOpen(true)}>
                  <Text className={classes.networkItemText}>
                    Create Network
                  </Text>
                </ListItem>
              </>
            )}
          </List>
        }
        onOpen={() => setIsMenuOpen(true)}
        onClose={() => setIsMenuOpen(false)}>
        <Tooltip title={networkId} placement="right">
          <div
            className={classNames({
              [classes.button]: true,
              [classes.openButton]: isMenuOpen,
            })}>
            <SettingsEthernetIcon className={classes.selectedNetwork} />
          </div>
        </Tooltip>
      </Popout>
    </>
  );
};

export default NetworkSelector;
