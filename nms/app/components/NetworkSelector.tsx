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
import AppContext from '../context/AppContext';
import Divider from '@mui/material/Divider';
import MenuButton from './MenuButton';
import MenuItem from '@mui/material/MenuItem';
import NetworkContext from '../context/NetworkContext';
import React, {useContext, useState} from 'react';
import Text from '../theme/design-system/Text';
import {LTE} from '../../shared/types/network';
import {NetworkEditDialog} from '../views/network/NetworkEdit';
import {makeStyles} from '@mui/styles';
import {useNavigate} from 'react-router-dom';

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
  const {networkIds, user} = useContext(AppContext);
  const [isNetworkAddOpen, setNetworkAddOpen] = useState(false);
  const {networkId: selectedNetworkId, networkType} = useContext(
    NetworkContext,
  );
  if (!selectedNetworkId) {
    return null;
  }

  if (!user.isSuperUser && networkIds.length < 2) {
    return <Text variant="body2">{`Network: ${selectedNetworkId}`}</Text>;
  }

  return (
    <>
      <NetworkEditDialog
        open={isNetworkAddOpen}
        onClose={() => {
          setNetworkAddOpen(false);
        }}
      />
      <MenuButton
        data-testid="networkSelector"
        className={classes.button}
        label={`Network: ${selectedNetworkId}`}>
        {user.isSuperUser && networkType === LTE && (
          <MenuItem onClick={() => setNetworkAddOpen(true)}>
            <Text variant="body2">Create Network</Text>
          </MenuItem>
        )}
        {user.isSuperUser && (
          <MenuItem
            component="a"
            onClick={() => {
              navigate(`/nms/${selectedNetworkId}/admin/networks`);
            }}>
            <Text variant="body2">Manage Networks</Text>
          </MenuItem>
        )}
        {user.isSuperUser && networkIds.length > 0 && (
          <Divider variant="middle" />
        )}
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
