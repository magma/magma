import { Menu, Item, MenuProvider, Separator, IconFont } from "react-contexify";
import "react-contexify/dist/ReactContexify.min.css";
import * as React from "react";

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
    let taskRefName = "<no ref name>";
    if (this.props.node?.extras?.inputs?.taskReferenceName) {
      taskRefName = this.props.node?.extras?.inputs?.taskReferenceName;
    }

    return (
      <Menu id={this.props.node.id}>
        <Item disabled={true}>{taskRefName}</Item>
        <Separator></Separator>
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
      storeRef={false}
    >
      {props.children}
    </MenuProvider>
  );
}
