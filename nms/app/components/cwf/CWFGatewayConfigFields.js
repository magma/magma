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
  allowed_gre_peers,
  ipdr_export_dst,
} from '../../../generated/MagmaAPIBindings';

import AddCircleOutline from '@material-ui/icons/AddCircleOutline';
import Button from '@material-ui/core/Button';
import Divider from '@material-ui/core/Divider';
import IconButton from '@material-ui/core/IconButton';
import React from 'react';
import RemoveCircleOutline from '@material-ui/icons/RemoveCircleOutline';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../theme/design-system/Text';
import TextField from '@material-ui/core/TextField';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  container: {
    display: 'block',
    margin: '5px 0',
    whiteSpace: 'nowrap',
    width: '100%',
  },
  divider: {
    margin: '10px 0',
  },
  inputKey: {
    width: '245px',
    paddingRight: '10px',
  },
  inputValue: {
    width: '240px',
  },
  icon: {
    width: '30px',
    height: '30px',
    verticalAlign: 'bottom',
  },
}));

type Props = {
  allowedGREPeers: allowed_gre_peers,
  onChange: allowed_gre_peers => void,
  ipdrExportDst: ?ipdr_export_dst,
  onIPDRChanged: ipdr_export_dst => void,
};

export default function (props: Props) {
  const classes = useStyles();
  const {allowedGREPeers, onIPDRChanged} = props;
  const ipdrExportDst = props.ipdrExportDst || {ip: '', port: 0};

  const onChange = (index, field: 'ip' | 'key', value) => {
    const newValue = [...allowedGREPeers];
    newValue[index] = {...allowedGREPeers[index]};
    if (field === 'key') {
      newValue[index].key = parseInt(value.replace(/[^0-9]*/g, '') || 0);
    } else {
      newValue[index][field] = value;
    }
    props.onChange(newValue);
  };

  const removePeer = index => {
    const newValue = [...allowedGREPeers];
    newValue.splice(index, 1);
    props.onChange(newValue);
  };

  const addPeer = () => {
    const newValue = [...allowedGREPeers];
    newValue.push({ip: '', key: 0});
    props.onChange(newValue);
  };

  let grePeersContent;
  if (allowedGREPeers.length === 0) {
    grePeersContent = (
      <Button color="primary" variant="text" onClick={addPeer}>
        Add GRE Peer
      </Button>
    );
  } else {
    grePeersContent = allowedGREPeers.map((peer, index) => (
      <div className={classes.container} key={index}>
        <TextField
          label="IP"
          margin="none"
          value={peer.ip}
          onChange={({target}) => onChange(index, 'ip', target.value)}
          className={classes.inputKey}
        />
        <TextField
          label="Key"
          margin="none"
          value={peer.key}
          onChange={({target}) => onChange(index, 'key', target.value)}
          className={classes.inputValue}
        />
        <IconButton onClick={() => removePeer(index)} className={classes.icon}>
          <RemoveCircleOutline />
        </IconButton>
        {index === allowedGREPeers.length - 1 && (
          <IconButton onClick={addPeer} className={classes.icon}>
            <AddCircleOutline />
          </IconButton>
        )}
      </div>
    ));
  }

  return (
    <div>
      <Text variant="subtitle1">IPDR Connections</Text>
      <div className={classes.container}>
        <TextField
          label="IP"
          margin="none"
          value={ipdrExportDst.ip}
          onChange={({target}) =>
            onIPDRChanged({...ipdrExportDst, ip: target.value})
          }
          className={classes.inputKey}
        />
        <TextField
          type="number"
          label="Port"
          margin="none"
          value={ipdrExportDst.port}
          onChange={({target}) =>
            onIPDRChanged({...ipdrExportDst, port: parseInt(target.value)})
          }
          className={classes.inputValue}
        />
      </div>
      <Divider className={classes.divider} />
      <Text variant="subtitle1">GRE Peers</Text>
      <div>{grePeersContent}</div>
    </div>
  );
}
