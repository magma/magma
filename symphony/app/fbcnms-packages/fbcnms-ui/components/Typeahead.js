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
import CancelIcon from '@material-ui/icons/Cancel';
import InputAffix from './design-system/Input/InputAffix';
import React, {useState} from 'react';
import TextInput from './design-system/Input/TextInput';
import emptyFunction from '@fbcnms/util/emptyFunction';
import symphony from '../theme/symphony';
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
    top: '32px',
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
  smallSuggest: {
    paddingTop: '9px',
    paddingBottom: '9px',
    paddingLeft: '14px',
    paddingRight: '14px',
    height: '14px',
  },
  cancelIcon: {
    color: symphony.palette.D300,
  },
}));

type Props = {
  margin?: ?string,
  required: ?boolean,
  suggestions: Array<Suggestion>,
  onEntitySelected: Suggestion => void,
  onSuggestionsFetchRequested: (searchTerm: string) => void,
  onSuggestionsClearRequested?: () => void,
  placeholder?: ?string,
  value?: ?Suggestion,
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
    placeholder,
    required,
    value,
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
          <TextInput
            type="string"
            required={!!required}
            variant="outlined"
            placeholder={placeholder ?? ''}
            fullWidth={true}
            disabled={selectedSuggestion != null}
            value={selectedSuggestion ? selectedSuggestion.name : ''}
            onChange={emptyFunction}
            suffix={
              searchTerm === '' ? (
                <InputAffix
                  onClick={() => {
                    setSearchTerm('');
                    setSelectedSuggestion(null);
                    onSuggestionsClearRequested &&
                      onSuggestionsClearRequested();
                  }}>
                  <CancelIcon className={classes.cancelIcon} />
                </InputAffix>
              ) : null
            }
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
            <TextInput
              type="string"
              placeholder={placeholder ?? ''}
              {...inputProps}
            />
          )}
          inputProps={{
            style: {},
            required: !!required,
            placeholder: placeholder ?? 'Search...',
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

export default Typeahead;
