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
 */

import AddCircleOutline from '@material-ui/icons/AddCircleOutline';
import DeleteIcon from '@material-ui/icons/Delete';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Text from '../../theme/design-system/Text';
import {makeStyles} from '@material-ui/styles';
import {policyStyles} from './PolicyStyles';
import {useEffect, useState} from 'react';
import type {PolicyRule} from '../../../generated-ts';

const useStyles = makeStyles(() => policyStyles);

type Props = {
  policyRule: PolicyRule;
  onChange: (policyRule: PolicyRule) => void;
};

export default function PolicyHeaderEnrichmentEdit(props: Props) {
  const classes = useStyles();
  const [newUrl, setNewUrl] = useState('');
  const [urlList, setUrlList] = useState(
    props.policyRule.header_enrichment_targets || [],
  );

  const addUrl = () => {
    if (newUrl !== '' || urlList.includes(newUrl)) {
      const newList = [...urlList, newUrl];
      props.onChange({...props.policyRule, header_enrichment_targets: newList});
      setNewUrl('');
    }
  };

  const deleteUrl = (url: string) => {
    const newList = urlList.filter(item => item !== url);
    props.onChange({...props.policyRule, header_enrichment_targets: newList});
  };

  useEffect(
    () => setUrlList(props.policyRule.header_enrichment_targets || []),
    [props.policyRule],
  );
  return (
    <div data-testid="headerEnrichmentEdit">
      <Text weight="medium" variant="subtitle2" className={classes.description}>
        {'List of URL targets for header enrichment'}
      </Text>
      <Grid container direction="column">
        <Grid item xs={12} md={6}>
          <ListItem dense disableGutters />
          <OutlinedInput
            data-testid="newUrl"
            placeholder="E.g. example.com/"
            value={newUrl}
            onChange={({target}) => {
              setNewUrl(target.value);
            }}
          />
          <IconButton data-testid="addUrlButton" onClick={addUrl}>
            <AddCircleOutline />
          </IconButton>
        </Grid>
        <Grid item xs={12} md={6}>
          {urlList.length > 0 && (
            <List component={Paper} dense>
              {urlList.map((url, i) => (
                <ListItem key={i} button>
                  <ListItemText primary={url} />
                  <ListItemSecondaryAction>
                    <IconButton
                      edge="end"
                      aria-label="delete"
                      data-testid="deleteUrlButton"
                      onClick={() => deleteUrl(url)}>
                      <DeleteIcon />
                    </IconButton>
                  </ListItemSecondaryAction>
                </ListItem>
              ))}
            </List>
          )}
        </Grid>
      </Grid>
    </div>
  );
}
