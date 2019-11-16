/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import LoadingFillerBackdrop from '@fbcnms/ui/components/LoadingFillerBackdrop';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import TextField from '@material-ui/core/TextField';

import nullthrows from '@fbcnms/util/nullthrows';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
});

type Props = {
  onCancel: () => void,
  onSave: string => void,
};

export default function(props: Props) {
  const classes = useStyles();
  const {match} = useRouter();
  const [baseName, setBaseName] = useState();
  const enqueueSnackbar = useEnqueueSnackbar();
  const editingBaseName = match.params.baseName;

  useEffect(() => {
    if (editingBaseName) {
      MagmaV1API.getNetworksByNetworkIdPoliciesBaseNamesByBaseName({
        networkId: nullthrows(match.params.networkId),
        baseName: editingBaseName,
      }).then(ruleNames =>
        setBaseName({
          name: editingBaseName,
          ruleNames: ruleNames.join(','),
        }),
      );
    } else {
      setBaseName({name: '', ruleNames: ''});
    }
  }, [editingBaseName, match.params.networkId, setBaseName]);

  if (!baseName) {
    return <LoadingFillerBackdrop />;
  }

  const onSave = async () => {
    try {
      if (editingBaseName) {
        await MagmaV1API.putNetworksByNetworkIdPoliciesBaseNamesByBaseName({
          networkId: nullthrows(match.params.networkId),
          baseName: editingBaseName,
          chargingRuleBaseName: baseName.ruleNames.split(','),
        });
        props.onSave(editingBaseName);
      } else {
        await MagmaV1API.postNetworksByNetworkIdPoliciesBaseNames({
          networkId: nullthrows(match.params.networkId),
          chargingRuleBaseName: {
            name: baseName.name,
            rule_names: baseName.ruleNames.split(','),
          },
        });
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
        <Button onClick={props.onCancel} skin="regular">
          Cancel
        </Button>
        <Button onClick={onSave}>Save</Button>
      </DialogActions>
    </Dialog>
  );
}
