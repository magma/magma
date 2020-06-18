/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {InventoryTreeNodeWithChildrenQuery} from './__generated__/InventoryTreeNodeWithChildrenQuery.graphql';
import type {Location} from '../common/Location';

import * as React from 'react';
import ArrowRightIcon from '@material-ui/icons/ArrowRight';
import CircularProgress from '@material-ui/core/CircularProgress';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import {Suspense, useEffect, useState} from 'react';
import {extractEntityIdFromUrl} from '../common/RouterUtils';
import {graphql, useLazyLoadQuery} from 'react-relay/hooks';
import {makeStyles} from '@material-ui/styles';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';
import {useHistory} from 'react-router';

import 'react-perfect-scrollbar/dist/css/styles.css';

const useStyles = makeStyles(theme => ({
  panelDetails: {
    padding: 0,
    display: 'block',
  },
  panelExpanded: {
    '& > $headerRoot > $arrowRightIcon': {
      transform: 'rotate(90deg)',
    },
    margin: 0,
    '&:before': {
      opacity: 0,
    },
  },
  childPanel: {
    '&:before': {
      opacity: 0,
    },
  },
  headerRoot: {
    display: 'flex',
    alignItems: 'center',
    width: '100%',
    cursor: 'pointer',
    '&:hover': {
      backgroundColor: 'rgba(0, 0, 0, 0.04)',
    },
    '&:hover::before, &$selected::before': {
      borderLeft: `2px solid ${theme.palette.primary.main}`,
    },
    '&:hover $hoverRightContent': {
      display: 'flex',
    },
    '&:hover $childrenLabel': {
      display: 'block',
      flexGrow: 1,
      alignItems: 'right',
    },
    '&:before': {
      content: '""',
      paddingTop: '16px',
      paddingBottom: '16px',
      borderLeft: `1px solid rgba(0, 0, 0, 0.086)`,
    },
    '&$selected': {
      backgroundColor: theme.palette.fadedBlue,
    },
  },
  selected: {},
  heading: {
    fontSize: theme.typography.pxToRem(13),
    marginRight: '4px',
    lineHeight: '100%',
  },
  childrenLabelContainer: {
    flexGrow: 1,
  },
  childrenLabel: {
    display: 'none',
    textAlign: 'right',
    marginRight: '8px',
    fontSize: theme.typography.pxToRem(11),
    color: theme.palette.text.secondary,
    lineHeight: '13px',
    marginLeft: '4px',
  },
  secondaryHeading: {
    fontSize: theme.typography.pxToRem(11),
    color: theme.palette.text.secondary,
    lineHeight: '13px',
    marginLeft: '4px',
  },
  headerContainer: {
    display: 'flex',
    flexGrow: 1,
  },
  headerLeftContent: {
    display: 'flex',
    flexGrow: 1,
    paddingBottom: '9px',
    paddingTop: '9px',
  },
  hoverRightContent: {
    display: 'none',
    alignItems: 'center',
    marginRight: '24px',
  },
  arrowRightIcon: {
    color: 'rgba(0, 0, 0, 0.54)',
    transition: 'transform 150ms cubic-bezier(0.4, 0, 0.2, 1) 0ms',
    marginLeft: '4px',
    cursor: 'pointer',
  },
  progress: {
    margin: '0px 8px',
  },
  leaf: {
    marginLeft: '12px',
  },
}));

const locationsWithChildrenQuery = graphql`
  query InventoryTreeNodeWithChildrenQuery($id: ID!) {
    location: node(id: $id) {
      ... on Location {
        ...LocationsTree_location @relay(mask: false)
        children {
          ...LocationsTree_location @relay(mask: false)
        }
      }
    }
  }
`;

type Props = {
  element: Location,
  parent: ?Location,
  depth: number,
  onClick: ?(string) => void,
  getHoverRightContent: ?(Object) => ?React.Node,
  selectedHierarchy: Array<string>,
};

function InventoryTreeNodeChildren(props: Props) {
  const {
    element,
    onClick,
    getHoverRightContent,
    selectedHierarchy,
    depth,
  } = props;

  const classes = useStyles();

  const data = useLazyLoadQuery<InventoryTreeNodeWithChildrenQuery>(
    locationsWithChildrenQuery,
    {id: element.id},
  );

  return (
    <div className={classes.panelDetails}>
      {data.location.children
        .slice()
        .filter(Boolean)
        .sort((x, y) => sortLexicographically(x.name ?? '', y.name ?? ''))
        .map(childLocation => (
          <InventoryTreeNode
            key={childLocation.id}
            onClick={onClick}
            selectedHierarchy={selectedHierarchy}
            element={childLocation}
            parent={element}
            depth={depth}
            getHoverRightContent={getHoverRightContent}
          />
        ))}
    </div>
  );
}

export default function InventoryTreeNode(props: Props) {
  const {
    element,
    parent,
    depth,
    getHoverRightContent,
    onClick,
    selectedHierarchy,
  } = props;
  const defaultIsSelected =
    extractEntityIdFromUrl('location', location.search) === element.id;

  const classes = useStyles();
  const history = useHistory();

  const [isExpanded, setIsExpanded] = useState<?boolean>(null);
  const [selected, setSelected] = useState(defaultIsSelected);

  useEffect(() => {
    if (selectedHierarchy.includes(element.id)) {
      setIsExpanded(true);
    }
  }, [selectedHierarchy, element.id]);

  useEffect(() => {
    const unlistener = history.listen(location => {
      const locationId = extractEntityIdFromUrl('location', location.search);
      setSelected(locationId === element.id);
    });
    return () => unlistener();
  }, [element.id, history]);

  const {numChildren} = element;
  const hasChildren = numChildren > 0;
  const unit = 8;
  const spacingPx = unit * 2;
  const key = element.id ?? element.name;
  const hoverRightContent =
    getHoverRightContent && getHoverRightContent(element);
  const paddingLeft = depth * spacingPx + spacingPx + unit;

  return (
    <div
      className={classNames({
        [classes.panelExpanded]: isExpanded,
        [classes.childPanel]: !!parent,
      })}>
      <div
        style={{paddingLeft}}
        className={classNames({
          [classes.headerRoot]: true,
          [classes.selected]: selected,
        })}>
        {hasChildren ? (
          <ArrowRightIcon
            data-testid={'inventory-expand-' + element.id}
            className={classes.arrowRightIcon}
            onClick={() => setIsExpanded(!isExpanded)}
          />
        ) : null}
        <div
          className={classNames({
            [classes.headerContainer]: true,
            [classes.leaf]: !hasChildren,
          })}
          onClick={() => {
            onClick && onClick(element.id);
            setSelected(true);
          }}>
          <div className={classes.headerLeftContent}>
            <Text className={classes.heading}>{element.name}</Text>
            <Text className={classes.secondaryHeading}>
              {element.locationType.name +
                (element.externalId ? ` - ${element.externalId}` : '')}
            </Text>
            <div className={classes.childrenLabelContainer}>
              <Text className={classes.childrenLabel}>({numChildren})</Text>
            </div>
          </div>
          <div className={classes.hoverRightContent}>{hoverRightContent}</div>
        </div>
      </div>
      {key && isExpanded ? (
        <Suspense
          fallback={
            <CircularProgress
              style={{marginLeft: paddingLeft}}
              className={classes.progress}
              size={16}
            />
          }>
          <InventoryTreeNodeChildren {...props} depth={depth + 1} />
        </Suspense>
      ) : null}
    </div>
  );
}
