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
import * as React from 'react';
import FBCMobileAppQRCode from './FBCMobileAppQRCode';
import Grid from '@material-ui/core/Grid';
import Typography from '@material-ui/core/Typography';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  root: {
    padding: theme.spacing(2),
  },
  qrCodeWrapper: {
    maxHeight: 250,
  },
}));

export type Props = {
  endpoint: string,
};

export default function FBCMobileAppConfigView(props: Props) {
  const {endpoint} = props;
  const classes = useStyles();
  return (
    <Grid
      className={classes.root}
      container
      justify="center"
      direction="column"
      alignItems="center"
      spacing={2}>
      <Grid className={classes.qrCodeWrapper} item>
        <FBCMobileAppQRCode endpoint={endpoint} />
      </Grid>
      <Grid item>
        <Typography>Scan this QR code using the FBC Mobile App</Typography>
      </Grid>
    </Grid>
  );
}
