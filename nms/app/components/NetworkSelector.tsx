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
import AppContext from '../components/context/AppContext';
import Divider from '@material-ui/core/Divider';
import MagmaAPI from '../../api/MagmaAPI';
import MenuButton from './MenuButton';
import MenuItem from '@material-ui/core/MenuItem';
import NetworkContext from './context/NetworkContext';
import React, {useCallback, useContext, useEffect, useState} from 'react';
import Text from '../theme/design-system/Text';
import useMagmaAPI from '../../api/useMagmaAPI';
import {LTE, coalesceNetworkType} from '../../shared/types/network';
import {NetworkEditDialog} from '../views/network/NetworkEdit';
import {makeStyles} from '@material-ui/styles';
import {useNavigate} from 'react-router-dom';
import type {NetworkType} from '../../shared/types/network';

const useStyles = makeStyles({
  button: {
    '&&': {
      textTransform: 'none',
    },
  },
});

const NetworkSelector = () => {
  const classes = useStyles();
  const navigate = useNavigate();
  const appContext = useContext(AppContext);
  const [networkIds, setNetworkIds] = useState<Array<string>>([]);
  const [networkType, setNetworkType] = useState<NetworkType | null>(null);
  const [lastRefreshTime, setLastRefreshTime] = useState(new Date().getTime());
  const [isNetworkAddOpen, setNetworkAddOpen] = useState(false);
  const {networkId: selectedNetworkId} = useContext(NetworkContext);

  useMagmaAPI(
    MagmaAPI.networks.networksGet,
    {},
    useCallback((resp: Array<string>) => setNetworkIds(resp), []),
    lastRefreshTime,
  );

  useEffect(() => {
    const fetchNetworkType = async () => {
      if (selectedNetworkId) {
        const networkType = (
          await MagmaAPI.networks.networksNetworkIdTypeGet({
            networkId: selectedNetworkId,
          })
        ).data;
        setNetworkType(coalesceNetworkType(selectedNetworkId, networkType));
      }
    };

    void fetchNetworkType();
  }, [selectedNetworkId]);

  if (!selectedNetworkId) {
    return null;
  }
  const {isSuperUser} = appContext.user;

  if (!isSuperUser && networkIds.length < 2) {
    return <Text variant="body2">{`Network: ${selectedNetworkId}`}</Text>;
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
      <MenuButton
        data-testid="networkSelector"
        className={classes.button}
        label={`Network: ${selectedNetworkId}`}>
        {isSuperUser && networkType === LTE && (
          <MenuItem onClick={() => setNetworkAddOpen(true)}>
            <Text variant="body2">Create Network</Text>
          </MenuItem>
        )}
        {isSuperUser && (
          <MenuItem
            component="a"
            onClick={() => {
              navigate(`/nms/${selectedNetworkId}/admin/networks`);
            }}>
            <Text variant="body2">Manage Networks</Text>
          </MenuItem>
        )}
        {isSuperUser && networkIds.length > 0 && <Divider variant="middle" />}
        {networkIds.map(id => (
          <MenuItem key={id} component="a" href={`/nms/${id}`}>
            <Text variant="body2">{id}</Text>
          </MenuItem>
        ))}
      </MenuButton>
    </>
  );
};

export default NetworkSelector;
