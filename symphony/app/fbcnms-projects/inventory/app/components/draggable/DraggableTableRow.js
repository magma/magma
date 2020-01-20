/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {WithStyles} from '@material-ui/core';

import * as React from 'react';
import DragIndicatorIcon from '@fbcnms/ui/icons/DragIndicatorIcon';
import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';
import {Draggable} from 'react-beautiful-dnd';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  id: string,
  index: number,
  children?: React.Node,
  className?: string,
  draggableCellClassName?: string,
} & WithStyles<typeof styles>;

const styles = {
  dragIndicatorIcon: {
    cursor: 'grab',
    '&&': {
      fontSize: '15px',
    },
  },
};

class DraggableTableRow extends React.Component<Props> {
  _rowComponent = DraggableComponent(
    this.props.id,
    this.props.index,
    this.props.classes,
    this.props.draggableCellClassName ?? '',
  );

  render() {
    return (
      <TableRow className={this.props.className} component={this._rowComponent}>
        {this.props.children}
      </TableRow>
    );
  }
}

const DraggableComponent = (
  id: string,
  index: number,
  // eslint-disable-next-line flowtype/no-weak-types
  classes: any,
  draggableCellClassName: string,
) => (props: Props) => {
  const {children} = props;
  return (
    <Draggable draggableId={id} index={index}>
      {provided => (
        <div
          key={`draggable_${id}`}
          ref={provided.innerRef}
          {...provided.draggableProps}
          {...props}>
          <TableCell
            className={draggableCellClassName}
            style={{minWidth: '35px', width: '35px'}}
            size="small"
            padding="none"
            component="div">
            <div {...provided.dragHandleProps}>
              <DragIndicatorIcon className={classes.dragIndicatorIcon} />
            </div>
          </TableCell>
          {children}
        </div>
      )}
    </Draggable>
  );
};
export default withStyles(styles)(DraggableTableRow);
