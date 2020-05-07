import * as React from "react";
import { DefaultNodeWidget } from "@projectstorm/react-diagrams";
import { NodeContextMenu, NodeMenuProvider } from "./ContextMenu";

export class NodeWithContextWidget extends DefaultNodeWidget {
  render() {
    return (
      <div {...this.getProps()}>
        <NodeMenuProvider node={this.props.node}>
          {super.render()}
        </NodeMenuProvider>
        <NodeContextMenu
          node={this.props.node}
          diagramEngine={this.props.diagramEngine}
        />
      </div>
    );
  }
}
