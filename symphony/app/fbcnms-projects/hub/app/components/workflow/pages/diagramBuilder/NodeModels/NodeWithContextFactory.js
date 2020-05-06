import * as React from "react";
import { DefaultNodeFactory } from "@projectstorm/react-diagrams";
import { NodeWithContextWidget } from "./NodeWithContextWidget";

export class NodeWithContextFactory extends DefaultNodeFactory {
  generateReactWidget(diagramEngine, node) {
		return React.createElement(NodeWithContextWidget, {
			node: node,
			diagramEngine: diagramEngine
		});
  }
}
