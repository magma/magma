import Grid from '@material-ui/core/Grid';
import React from 'react';
import Text from '../../theme/design-system/Text';

import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  cardTitleRow: {
    marginBottom: theme.spacing(1),
    minHeight: '36px',
  },
  cardTitleIcon: {
    fill: colors.primary.comet,
    marginRight: theme.spacing(1),
  },
}));

export const CardTitleRow = props => {
  const classes = useStyles();
  const Icon = props.icon;

  return (
    <Grid container alignItems="center" className={classes.cardTitleRow}>
      <Icon className={classes.cardTitleIcon} />
      <Text variant="body1">{props.label}</Text>
    </Grid>
  );
};

export const CardTitleFilterRow = props => {
  const classes = useStyles();
  const Icon = props.icon;
  const Filters = props.filter;

  return (
    <Grid container alignItems="center" className={classes.cardTitleRow}>
      <Grid container xs>
        <Icon className={classes.cardTitleIcon} />
        <Text variant="body1">{props.label}</Text>
      </Grid>
      <Grid item>
        <Filters />
      </Grid>
    </Grid>
  );
};
