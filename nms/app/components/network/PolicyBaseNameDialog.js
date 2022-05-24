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

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import LoadingFillerBackdrop from '../LoadingFillerBackdrop';
import MagmaV1API from '../../../generated/WebClient';
import React from 'react';
import TextField from '@material-ui/core/TextField';

// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useParams} from 'react-router-dom';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

type Props = {
  mirrorNetwork?: string,
  onCancel: () => void,
  onSave: string => void,
};

export default function (props: Props) {
  const classes = useStyles();
  const params = useParams();
  const [baseName, setBaseName] = useState();
  const enqueueSnackbar = useEnqueueSnackbar();
  const editingBaseName = params.baseName;

  useEffect(() => {
    if (editingBaseName) {
      MagmaV1API.getNetworksByNetworkIdPoliciesBaseNamesByBaseName({
        networkId: nullthrows(params.networkId),
        baseName: editingBaseName,
      }).then(response =>
        setBaseName({
          name: editingBaseName,
          ruleNames: response.rule_names.join(','),
        }),
      );
    } else {
      setBaseName({name: '', ruleNames: ''});
    }
  }, [editingBaseName, params.networkId, setBaseName]);

  if (!baseName) {
    return <LoadingFillerBackdrop />;
  }

  const onSave = async () => {
    const baseNameRecord = {
      name: baseName.name,
      rule_names: baseName.ruleNames.split(','),
    };

    const data = [
      {
        networkId: nullthrows(params.networkId),
        baseNameRecord,
      },
    ];

    if (props.mirrorNetwork) {
      data.push({
        networkId: props.mirrorNetwork,
        baseNameRecord,
      });
    }

    try {
      if (editingBaseName) {
        await Promise.all(
          data.map(d =>
            MagmaV1API.putNetworksByNetworkIdPoliciesBaseNamesByBaseName({
              ...d,
              baseName: editingBaseName,
            }),
          ),
        );
      } else {
        await Promise.all(
          data.map(d => MagmaV1API.postNetworksByNetworkIdPoliciesBaseNames(d)),
        );
      }

      props.onSave(baseName.name);
    } catch (e) {
      enqueueSnackbar(e?.response?.data?.message || e, {variant: 'error'});
    }
  };

  return (
    <Dialog open={true} onClose={props.onCancel} scroll="body">
      <DialogTitle>{editingBaseName ? 'Edit' : 'Add'} Base Name</DialogTitle>
      <DialogContent>
        <TextField
          required
          className={classes.input}
          label="Base Name"
          margin="normal"
          disabled={!!editingBaseName}
          value={baseName.name}
          onChange={({target}) =>
            setBaseName({...baseName, name: target.value})
          }
        />
        <TextField
          required
          className={classes.input}
          label="Rule Names (CSV)"
          margin="normal"
          value={baseName.ruleNames}
          onChange={({target}) =>
            setBaseName({...baseName, ruleNames: target.value})
          }
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onCancel}>Cancel</Button>
        <Button onClick={onSave} variant="contained" color="primary">
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}
