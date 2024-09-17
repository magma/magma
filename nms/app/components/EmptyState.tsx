/**
 * Copyright 2022 The Magma Authors.
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
import * as React from 'react';
import Button from '@mui/material/Button';
import Card from '@mui/material/Card';
import CardActions from '@mui/material/CardActions';
import CardContent from '@mui/material/CardContent';
import CardHeader from '@mui/material/CardHeader';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import Text from './../theme/design-system/Text';
import {Theme} from '@mui/material';
import {makeStyles} from '@mui/styles';

const useStyles = makeStyles<Theme>(theme => ({
  card: {
    padding: '10px',
    paddingBottom: '40px',
    height: '100%',
  },
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  cardAction: {
    padding: '16px',
  },
  emptyState: {
    margin: '0',
  },
  emptyStateLink: {
    marginLeft: '0px',
    fontSize: '14px',
  },
  instructions: {
    fontSize: '14px',
  },
  overviewDescription: {
    fontSize: '14px',
  },
  overviewTitle: {
    fontSize: '18px',
  },
}));

type CardActionsType = {
  // Add Button text in empty state instructions
  buttonText: string;
  // Action executed on Add Button click
  onClick: () => void;
  // Link text to magma documentation
  linkText: string;
  // Link to magma documentation
  link: string;
};

export type emptyStateProps = {
  title: string;
  customIntructions?: React.ReactNode;
  instructions: string;
  overviewTitle: string;
  overviewDescription: string;
  cardActions?: CardActionsType;
};

/**
 * Empty state instructions. Display 2 card with instructions when there is no data to display.
 * The left card contains instructions on how to use nms on the different pages and an Add Button.
 * The right card contains an overview of the page.
 *
 * @param title Empty state instructions title
 * @param customIntructions Custom instruction layout component
 * @param instructions Instructions text of the empty state
 * @param overviewTitle Title of the empty state overview,
 * @param overviewDescription Description of the empty state overview,
 * @param cardActions Instruction actions (button and link)
 */
export default function (props: emptyStateProps) {
  const classes = useStyles();
  return (
    <Grid container className={classes.emptyState} spacing={6}>
      <Grid item xs={8}>
        <Card raised={false} className={classes.card}>
          <CardHeader title={props.title} />
          <CardContent>
            <Text className={classes.instructions}>{props.instructions}</Text>
            {props.customIntructions}
          </CardContent>
          {Object.keys(props.cardActions || {}).length > 0 && (
            <CardActions
              disableSpacing={true}
              classes={{root: classes.cardAction}}>
              <Grid container direction="column" spacing={3}>
                <Grid item xs={3}>
                  <Button
                    variant="contained"
                    color="primary"
                    onClick={props.cardActions?.onClick}>
                    {props.cardActions?.buttonText || ''}
                  </Button>
                </Grid>
                <Grid item>
                  <Link
                    target="_blank"
                    className={classes.emptyStateLink}
                    href={props.cardActions?.link}
                    underline="hover">
                    {props.cardActions?.linkText}
                  </Link>
                </Grid>
              </Grid>
            </CardActions>
          )}
        </Card>
      </Grid>
      <Grid item xs={4}>
        <Card raised={false} className={classes.card}>
          <CardHeader title={props.overviewTitle} />
          <CardContent classes={{root: classes.overviewDescription}}>
            {props.overviewDescription}
          </CardContent>
        </Card>
      </Grid>
    </Grid>
  );
}
