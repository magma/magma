/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Autosuggest from 'react-autosuggest';
import ClearIcon from '@material-ui/icons/Clear';
import IconButton from '@material-ui/core/IconButton';
import InputAdornment from '@material-ui/core/InputAdornment';
import React, {useState} from 'react';
import TextField from '@material-ui/core/TextField';
import emptyFunction from '@fbcnms/util/emptyFunction';
import {blue05} from '../theme/colors';
import {makeStyles, useTheme} from '@material-ui/styles';

const autoSuggestStyles = theme => ({
  container: {
    position: 'relative',
  },
  input: {
    width: '100%',
    padding: '0px',
    fontFamily: theme.typography.subtitle1.fontFamily,
    fontWeight: theme.typography.subtitle1.fontWeight,
    fontSize: theme.typography.subtitle1.fontSize,
    color: theme.typography.subtitle1.color,
  },
  suggestionsContainer: {
    display: 'none',
  },
  suggestionsContainerOpen: {
    backgroundColor: theme.palette.common.white,
    borderRadius: '2px',
    display: 'block',
    fontFamily: theme.typography.subtitle1.fontFamily,
    fontSize: theme.typography.subtitle1.fontSize,
    fontWeight: theme.typography.subtitle1.fontWeight,
    position: 'absolute',
    top: '42px',
    minWidth: '100%',
    zIndex: '2',
    boxShadow: theme.shadows[2],
  },
  suggestionsList: {
    margin: 0,
    padding: 0,
    listStyleType: 'none',
  },
  suggestion: {
    color: theme.palette.grey[900],
    cursor: 'pointer',
    padding: '10px 20px',
  },
  suggestionHighlighted: {
    backgroundColor: theme.palette.grey[100],
  },
});

const useStyles = makeStyles(theme => ({
  container: {
    minWidth: '250px',
    width: '100%',
  },
  suggestionRoot: {
    display: 'flex',
  },
  suggestionType: {
    color: theme.palette.text.secondary,
    fontSize: theme.typography.pxToRem(13),
    lineHeight: '21px',
    marginLeft: theme.spacing(),
  },
  outlinedInput: {
    '&&': {
      backgroundColor: blue05,
      color: theme.palette.text.primary,
    },
  },
  clearIcon: {
    padding: 0,
    color: theme.palette.grey[300],
    '&:hover': {
      color: theme.palette.grey[600],
      background: 'none',
    },
  },
  shrinkedInputLabel: {
    '&&': {
      color: theme.palette.grey[600],
    },
  },
  smallSuggest: {
    paddingTop: '9px',
    paddingBottom: '9px',
    paddingLeft: '14px',
    paddingRight: '14px',
    height: '14px',
  },
}));

type Props = {
  margin?: ?string,
  required: ?boolean,
  suggestions: Array<Suggestion>,
  onEntitySelected: Suggestion => void,
  onSuggestionsFetchRequested: (searchTerm: string) => void,
  onSuggestionsClearRequested?: () => void,
  headline?: ?string,
  value?: ?Suggestion,
  variant?: 'default' | 'small',
  displayText?: string,
};

export type Suggestion = {
  entityId: string,
  entityType: string,
  name: string,
  type: string,
};

const Typeahead = (props: Props) => {
  const {
    onSuggestionsFetchRequested,
    onSuggestionsClearRequested,
    onEntitySelected,
    suggestions,
    headline,
    required,
    value,
    variant,
    displayText,
    margin,
  } = props;
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedSuggestion, setSelectedSuggestion] = useState(value);
  const classes = useStyles();
  const theme = useTheme();
  return (
    <div className={classes.container}>
      {selectedSuggestion && onSuggestionsClearRequested ? (
        <div>
          <TextField
            required={!!required}
            variant="outlined"
            label={headline ?? ''}
            fullWidth={true}
            disabled={selectedSuggestion != null}
            value={selectedSuggestion ? selectedSuggestion.name : ''}
            onChange={emptyFunction}
            InputLabelProps={{
              classes: {
                shrink: classes.shrinkedInputLabel,
              },
            }}
            InputProps={{
              margin,
              classes: {
                root: classes.outlinedInput,
                input: variant === 'small' ? classes.smallSuggest : '',
              },
              endAdornment: (
                <InputAdornment position="end">
                  <IconButton
                    className={classes.clearIcon}
                    onClick={() => {
                      setSearchTerm('');
                      setSelectedSuggestion(null);
                      onSuggestionsClearRequested &&
                        onSuggestionsClearRequested();
                    }}>
                    <ClearIcon />
                  </IconButton>
                </InputAdornment>
              ),
            }}
          />
        </div>
      ) : (
        <Autosuggest
          suggestions={suggestions}
          getSuggestionValue={suggestion => suggestion.name}
          onSuggestionsFetchRequested={({value}) => {
            onSuggestionsFetchRequested(value);
          }}
          onSuggestionsClearRequested={emptyFunction}
          renderSuggestion={suggestion => (
            <div className={classes.suggestionRoot}>
              <div>{suggestion.name}</div>
              <div className={classes.suggestionType}>{suggestion.type}</div>
            </div>
          )}
          onSuggestionSelected={(e, data) => {
            const suggestion: Suggestion = data.suggestion;
            setSearchTerm('');
            setSelectedSuggestion(suggestion);
            onEntitySelected(suggestion);
          }}
          theme={autoSuggestStyles(theme)}
          renderInputComponent={inputProps => (
            <TextField
              variant="outlined"
              label={headline ?? ''}
              {...inputProps}
              InputProps={{
                classes:
                  variant === 'small' ? {input: classes.smallSuggest} : {},
              }}
            />
          )}
          inputProps={{
            required: !!required,
            placeholder: displayText ?? 'Search...',
            value: searchTerm,
            margin,
            onChange: (_e, {newValue}) => setSearchTerm(newValue),
          }}
          highlightFirstSuggestion={true}
        />
      )}
    </div>
  );
};

Typeahead.defaultProps = {
  variant: 'default',
};

export default Typeahead;
