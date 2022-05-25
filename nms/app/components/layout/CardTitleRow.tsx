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

import Grid from '@material-ui/core/Grid';
import React from 'react';
import SvgIcon from '@material-ui/core/SvgIcon';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../theme/design-system/Text';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import type {Theme} from '@material-ui/core/styles';

const useStyles = makeStyles<Theme>(theme => ({
  cardTitleRow: {
    marginBottom: theme.spacing(1),
    minHeight: '36px',
  },
  cardTitleIcon: {
    fill: colors.primary.comet,
    marginRight: theme.spacing(1),
  },
}));

type CardTitleRowProps = {
  icon?: typeof SvgIcon;
  label: string;
  filter?: () => React.ReactElement;
};

export default function CardTitleRow(props: CardTitleRowProps) {
  const classes = useStyles();
  const Filters = props.filter;
  const Icon = props.icon;

  return (
    <Grid container alignItems="center" className={classes.cardTitleRow}>
      <Grid item xs>
        <Grid container alignItems="center">
          {Icon ? <Icon className={classes.cardTitleIcon} /> : null}
          <Text variant="body1">{props.label}</Text>
        </Grid>
      </Grid>
      {Filters ? (
        <Grid item>
          <Filters />
        </Grid>
      ) : null}
    </Grid>
  );
}
