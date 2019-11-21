/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import * as React from 'react';
import TableBody from '@material-ui/core/TableBody';
import {DragDropContext, Droppable} from 'react-beautiful-dnd';
import type {DropResult, ResponderProvided} from 'react-beautiful-dnd';

type Props = {
  className?: string,
  children?: React.Node,
  onDragEnd: (result: DropResult, provided: ResponderProvided) => void,
};

class DroppableTableBody extends React.Component<Props> {
  _bodyComponent = DroppableComponent((a1, a2) => this.props.onDragEnd(a1, a2));

  render() {
    return (
      <TableBody
        className={this.props.className}
        component={this._bodyComponent}>
        {this.props.children}
      </TableBody>
    );
  }
}

const DroppableComponent = (
  onDragEnd: (result: DropResult, provided: ResponderProvided) => void,
) => (props: any) => {
  return (
    <DragDropContext onDragEnd={onDragEnd}>
      <Droppable droppableId={'1'}>
        {provided => {
          return (
            <div
              ref={provided.innerRef}
              {...provided.droppableProps}
              {...props}>
              {props.children}
              {provided.placeholder}
            </div>
          );
        }}
      </Droppable>
    </DragDropContext>
  );
};

export default DroppableTableBody;
