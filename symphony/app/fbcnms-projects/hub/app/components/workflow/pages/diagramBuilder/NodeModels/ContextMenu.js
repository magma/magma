/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import 'react-contexify/dist/ReactContexify.min.css';
import * as React from 'react';
import {IconFont, Item, Menu, MenuProvider, Separator} from 'react-contexify';

export class NodeContextMenu extends React.Component {
  deleteNode = (node, diagramEngine) => {
    node.remove();
    diagramEngine.getDiagramModel().removeNode(node);
    diagramEngine.repaintCanvas();
  };

  handleDelete = () => {
    this.deleteNode(this.props.node, this.props.diagramEngine);
  };

  render() {
    let taskRefName = '<no ref name>';
    if (this.props.node?.extras?.inputs?.taskReferenceName) {
      taskRefName = this.props.node?.extras?.inputs?.taskReferenceName;
    }

    return (
      <Menu id={this.props.node.id}>
        <Item disabled={true}>{taskRefName}</Item>
        <Separator />
        <Item onClick={this.handleDelete}>
          <IconFont className="fa fa-trash" />
          Delete
        </Item>
      </Menu>
    );
  }
}

export function NodeMenuProvider(props) {
  return <MenuProvider id={props.node.id}>{props.children}</MenuProvider>;
}

export class LinkContextMenu extends React.Component {
  deleteLink = (link, diagramEngine) => {
    link.remove();
    diagramEngine.getDiagramModel().removeLink(link);
    diagramEngine.repaintCanvas();
  };

  handleDelete = () => {
    this.deleteLink(this.props.link, this.props.diagramEngine);
  };

  render() {
    return (
      <Menu id={this.props.link.id} event="onContextMenu" storeRef={false}>
        <Item onClick={this.handleDelete}>
          <IconFont className="fa fa-trash" />
          Delete
        </Item>
      </Menu>
    );
  }
}

export function LinkMenuProvider(props) {
  return (
    <MenuProvider
      id={props.link.id}
      component="g"
      event="onContextMenu"
      storeRef={false}>
      {props.children}
    </MenuProvider>
  );
}
