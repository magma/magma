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
import Button from '@material-ui/core/Button';
import CardTitleRow from './layout/CardTitleRow';
import FormLabel from '@material-ui/core/FormLabel';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import ReactJson, {InteractionProps, ReactJsonViewProps} from 'react-json-view';
import SettingsIcon from '@material-ui/icons/Settings';
import {Theme} from '@material-ui/core/styles';
import {colors, typography} from '../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';

const useStyles = makeStyles<Theme>(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  itemTitle: {
    color: colors.primary.comet,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  itemValue: {
    color: colors.primary.brightGray,
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
  jsonTextarea: {
    fontFamily: 'monospace',
    height: '95%',
    border: 'none',
    margin: theme.spacing(2),
  },
  configBody: {
    display: 'flex',
    flexFlow: 'column',
    flexGrow: 1,
    overflowX: 'hidden',
  },
  appBarBtnSecondary: {
    color: colors.primary.brightGray,
  },
}));

type Props<T> = {
  content: T;
  error: string;
  customFilter?: React.ReactNode;
  onSave: (arg0: T) => Promise<void>;
};

export default function JsonEditor<T>(props: Props<T>) {
  const classes = useStyles();
  const [error, setError] = useState<string>(props.error);
  const [content, setContent] = useState<T>(props.content);

  useEffect(() => {
    setError(props.error);
  }, [props.error]);

  const handleChange = (data: InteractionProps) => {
    setContent((data.updated_src as unknown) as T);
  };

  const JsonFilter = () => {
    return (
      <Grid container alignItems="center">
        {props.customFilter}
        <Grid item>
          <Button
            className={classes.appBarBtnSecondary}
            onClick={() => {
              setContent(props.content);
              setError('');
            }}>
            Clear
          </Button>
        </Grid>
        <Grid item>
          <Button
            className={classes.appBarBtn}
            onClick={() => {
              try {
                void props.onSave(content);
              } catch (e) {
                setError((e as Error).message);
              }
            }}>
            Save
          </Button>
        </Grid>
      </Grid>
    );
  };

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={3}>
        <Grid item xs={12}>
          <CardTitleRow
            icon={SettingsIcon}
            label="JSON Config"
            filter={JsonFilter}
          />
        </Grid>

        <Grid
          container
          className={classes.configBody}
          alignItems="stretch"
          item
          xs={12}>
          {error !== '' && <FormLabel error>{error}</FormLabel>}
          <ReactJson
            src={(content as unknown) as ReactJsonViewProps['src']}
            enableClipboard={false}
            displayDataTypes={false}
            onAdd={handleChange}
            onEdit={handleChange}
            onDelete={handleChange}
          />
        </Grid>
      </Grid>
    </div>
  );
}
