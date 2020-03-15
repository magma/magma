/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {
  EntityConfig,
  FilterConfig,
  SavedSearchConfig,
} from './ComparisonViewTypes';

import * as React from 'react';
import ChevronRightIcon from '@material-ui/icons/ChevronRight';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import nullthrows from '@fbcnms/util/nullthrows';
import {makeStyles} from '@material-ui/styles';
import {
  useCallback,
  useImperativeHandle,
  useMemo,
  useRef,
  useState,
} from 'react';

const useStyles = makeStyles(theme => ({
  root: {
    position: 'relative',
  },
  rootInput: {
    width: '100%',
    margin: theme.spacing() + 2,
    outline: 0,
    padding: 0,
    fontSize: '14px',
    border: 0,
    flexGrow: 1,
  },
  filterMenuItem: {
    padding: '6px 20px',
    cursor: 'pointer',
  },
  filterMenuItemText: {
    fontSize: '12px',
    lineHeight: '16px',
  },
  expansionPanel: {
    margin: '0px 0px',
  },
  expansionDetails: {
    padding: '0px',
  },
  panelExpanded: {
    '& > $headerRoot > $arrowRightIcon': {
      transform: 'rotate(90deg)',
    },
  },
  expansionSummary: {
    minHeight: '28px',
    padding: '0px 12px',
    '&$panelExpanded': {
      minHeight: '28px',
    },
    '& $headerRoot': {
      paddingRight: 0,
    },
  },
  headerRoot: {
    display: 'flex',
    alignItems: 'center',
    width: '100%',
  },
  arrowRightIcon: {
    fontSize: '16px',
    color: 'rgba(0, 0, 0, 0.54)',
    transition: 'transform 150ms cubic-bezier(0.4, 0, 0.2, 1) 0ms',
    marginLeft: '4px',
  },
  entityHeader: {
    flexGrow: 1,
    fontSize: '12px',
    lineHeight: '16px',
    fontWeight: 'bold',
  },
  expansionSummaryContent: {
    margin: '0px',
    '&$panelExpanded': {
      margin: '0px',
    },
  },
  popover: {
    maxHeight: '600px',
    overflowY: 'auto',
    minWidth: '172px',
    paddingBottom: '4px',
    paddingTop: '4px',
    backgroundColor: theme.palette.common.white,
    boxShadow: theme.shadows[1],
    position: 'fixed',
    zIndex: 10,
  },
  expansionPanelRoot: {
    '&:before': {
      display: 'none',
    },
  },
  entityFiltersList: {
    width: '100%',
    paddingBottom: '0px',
    paddingTop: '0px',
  },
  noMatchesText: {
    fontSize: '12px',
    lineHeight: '16px',
    margin: '6px 12px',
  },
  selectedFilterItem: {
    backgroundColor: theme.palette.action.hover,
  },
}));

type SelectableOption =
  | {
      option: FilterConfig,
      optionType: 'FILTER',
    }
  | {
      option: SavedSearchConfig,
      optionType: 'SAVED_SEARCH',
    };

type Props = {|
  options: Array<FilterConfig>,
  searchConfig: Array<EntityConfig>,
  savedSearches: Array<SavedSearchConfig>,
  onInputBlurred: () => void,
  onInputFocused: () => void,
  onFilterSelected: SelectableOption => void,
|};

const KEYBOARD_ENTER_KEY_CODE = 13;
const KEYBOARD_UP_KEY_CODE = 38;
const KEYBOARD_DOWN_KEY_CODE = 40;

/* $FlowFixMe - Flow doesn't support typing when using forwardRef on a
 * funcional component
 */
const FiltersTyepahead = React.forwardRef((props: Props, ref) => {
  const {
    onFilterSelected,
    onInputFocused,
    onInputBlurred,
    options,
    searchConfig,
    savedSearches,
  } = props;
  const inputRef = useRef();
  const menuRef = useRef();
  const [isMenuOpen, toggleMenu] = useState(false);
  const [inputValue, setInputValue] = useState('');
  const [hoveredFilterIndex, setHoveredFilterIndex] = useState(0);
  const classes = useStyles();

  useImperativeHandle(ref, () => ({
    focus: () => inputRef.current && inputRef.current.focus(),
  }));

  const filterMatchesInput = useCallback(
    (filter: FilterConfig) =>
      !inputValue ||
      filter.label.toLowerCase().includes(inputValue.toLowerCase()),
    [inputValue],
  );
  const entityResultExists = useCallback(
    entityType =>
      options
        .filter(filter => filter.entityType === entityType)
        .some(filterMatchesInput),
    [filterMatchesInput, options],
  );
  const anyResultExists = useMemo(
    () => searchConfig.map(entity => entity.type).some(entityResultExists),
    [entityResultExists, searchConfig],
  );

  const optionKeys = useMemo(
    () => [...options.map(f => f.key), ...savedSearches.map(s => s.key)],
    [options, savedSearches],
  );

  const hoveredFilterKey =
    !optionKeys || optionKeys.length === 0
      ? null
      : optionKeys[hoveredFilterIndex];

  const selectFilter = useCallback(
    optionObj => {
      setInputValue('');
      onFilterSelected(optionObj);
      toggleMenu(false);
    },
    [onFilterSelected],
  );

  const onFilterHovered = useCallback(
    filter =>
      setHoveredFilterIndex(optionKeys.findIndex(k => filter.key === k)),
    [optionKeys],
  );

  const findOptionConfigByKey = (key: string): SelectableOption => {
    let conf = savedSearches.find(f => f.key == key);
    if (conf != null) {
      return {
        option: nullthrows(conf),
        optionType: 'SAVED_SEARCH',
      };
    }
    conf = options.find(f => f.key == key);
    if (conf != null) {
      return {
        option: nullthrows(conf),
        optionType: 'FILTER',
      };
    }
    throw new Error(`key not found`);
  };

  const expansionRowForSavedSearch = (
    savedSearches: Array<SavedSearchConfig>,
  ): React.Node => {
    if (savedSearches.length === 0) {
      return null;
    }
    return (
      <ExpansionPanelDetails classes={{root: classes.expansionDetails}}>
        <div className={classes.entityFiltersList}>
          {savedSearches.map(savedSearch => (
            <div
              key={`${savedSearch.id}`}
              className={classNames({
                [classes.filterMenuItem]: true,
                [classes.selectedFilterItem]:
                  hoveredFilterKey === savedSearch.key,
              })}
              onMouseDown={e => e.preventDefault()}
              onMouseOver={() => onFilterHovered(savedSearch)}
              onClick={() =>
                selectFilter({option: savedSearch, optionType: 'SAVED_SEARCH'})
              }>
              <Text className={classes.filterMenuItemText}>
                {savedSearch.label}
              </Text>
            </div>
          ))}
        </div>
      </ExpansionPanelDetails>
    );
  };
  const expansionRowForFilterConfig = (entity: EntityConfig): React.Node => {
    if (!entityResultExists(entity.type)) {
      return null;
    }
    const entityOptions = options.filter(
      filter => filter.entityType === entity.type,
    );
    return (
      <ExpansionPanelDetails classes={{root: classes.expansionDetails}}>
        <div className={classes.entityFiltersList}>
          {entityOptions.map(filter =>
            filterMatchesInput(filter) ? (
              <div
                key={`${filter.key}-${filter.name}`}
                className={classNames({
                  [classes.filterMenuItem]: true,
                  [classes.selectedFilterItem]: hoveredFilterKey === filter.key,
                })}
                onMouseDown={e => e.preventDefault()}
                onMouseOver={() => onFilterHovered(filter)}
                onClick={() =>
                  selectFilter({option: filter, optionType: 'FILTER'})
                }>
                <Text className={classes.filterMenuItemText}>
                  {filter.label}
                </Text>
              </div>
            ) : null,
          )}
        </div>
      </ExpansionPanelDetails>
    );
  };

  const searchConfigExpansionPanels = searchConfig.map(entity => {
    return {
      label: entity.label,
      options: expansionRowForFilterConfig(entity),
    };
  });

  const labelToExpansionPanelDetails: Array<{
    label: string,
    options: React.Node,
  }> = [
    {
      label: 'Saved Searches',
      options: expansionRowForSavedSearch(savedSearches),
    },
    ...searchConfigExpansionPanels,
  ];

  return (
    <div
      className={classes.root}
      onClick={e => {
        if (!menuRef.current) {
          return;
        }

        if (!menuRef.current.contains(e.target)) {
          inputRef.current && inputRef.current.focus();
        }
      }}>
      <input
        className={classes.rootInput}
        autoFocus={false}
        type="text"
        value={inputValue}
        onKeyDown={e => {
          switch (e.keyCode) {
            case KEYBOARD_UP_KEY_CODE:
              setHoveredFilterIndex(
                hoveredFilterIndex === 0
                  ? optionKeys.length - 1
                  : hoveredFilterIndex - 1,
              );
              break;
            case KEYBOARD_DOWN_KEY_CODE:
              setHoveredFilterIndex(
                (hoveredFilterIndex + 1) % optionKeys.length,
              );
              break;
            case KEYBOARD_ENTER_KEY_CODE:
              selectFilter(
                findOptionConfigByKey(optionKeys[hoveredFilterIndex]),
              );
              break;
            default:
              return;
          }

          e.preventDefault();
        }}
        onChange={({target}) => {
          setInputValue(target.value);
          setHoveredFilterIndex(0);
        }}
        onFocus={() => {
          toggleMenu(true);
          onInputFocused();
        }}
        onBlur={e => {
          if (!menuRef.current) {
            return;
          }

          if (!menuRef.current.contains(e.relatedTarget)) {
            toggleMenu(false);
            setInputValue('');
            onInputBlurred();
          }
        }}
        ref={inputRef}
      />
      {isMenuOpen ? (
        <div className={classes.popover} ref={menuRef}>
          {!anyResultExists ? (
            <Text className={classes.noMatchesText}>No matches</Text>
          ) : null}
          {labelToExpansionPanelDetails.map(entity => {
            if (entity.options == null) {
              return null;
            }
            return (
              <ExpansionPanel
                key={entity.label}
                classes={{
                  root: classes.expansionPanelRoot,
                  expanded: classes.expansionPanel,
                }}
                defaultExpanded={true}
                elevation={0}>
                <ExpansionPanelSummary
                  classes={{
                    root: classes.expansionSummary,
                    expanded: classes.panelExpanded,
                    content: classes.expansionSummaryContent,
                  }}>
                  <div className={classes.headerRoot}>
                    <Text className={classes.entityHeader}>{entity.label}</Text>
                    <ChevronRightIcon className={classes.arrowRightIcon} />
                  </div>
                </ExpansionPanelSummary>
                {entity.options}
              </ExpansionPanel>
            );
          })}
        </div>
      ) : null}
    </div>
  );
});

export default FiltersTyepahead;
