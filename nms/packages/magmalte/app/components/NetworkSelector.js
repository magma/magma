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
import type {NetworkType} from '../../fbc_js_core/types/network';

import AppContext from '../../fbc_js_core/ui/context/AppContext';
import Divider from '@material-ui/core/Divider';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import MagmaV1API from '../../generated/WebClient';
import NetworkContext from './context/NetworkContext';
import Popout from '../../fbc_js_core/ui/components/Popout';
import React from 'react';
import SettingsEthernetIcon from '@material-ui/icons/SettingsEthernet';
import Text from '../../fbc_js_core/ui/components/design-system/Text';
import Tooltip from '@material-ui/core/Tooltip';
import useMagmaAPI from '../../api/useMagmaAPI';

import {LTE, coalesceNetworkType} from '../../fbc_js_core/types/network';
import {NetworkEditDialog} from '../views/network/NetworkEdit';
import {colors} from '../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useContext, useEffect, useState} from 'react';

const useStyles = makeStyles(_ => ({
  button: {
    display: 'flex',
    justifyContent: 'center',
    width: '100%',
    padding: '15px 0px',
    cursor: 'pointer',
    '&:hover $icon': {
      color: colors.primary.white,
    },
  },
  icon: {
    color: colors.primary.gullGray,
  },
  itemGutters: {
    '&&': {
      minWidth: '200px',
      padding: '6px 17px',
      '&:hover': {
        backgroundColor: colors.primary.concrete,
      },
    },
  },
  divider: {
    margin: '6px 17px',
  },
  networksList: {
    '&&': {
      maxHeight: '400px',
      overflowY: 'auto',
      padding: '10px 0',
    },
  },
  networkItemText: {
    fontSize: '14px',
    lineHeight: '20px',
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
        className={classes.button}
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
                <Divider className={classes.divider} />
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
        <Tooltip
          title={networkId}
          placement="right"
          data-testid="networkSelector">
          <SettingsEthernetIcon className={classes.icon} />
        </Tooltip>
      </Popout>
    </>
  );
};

export default NetworkSelector;
