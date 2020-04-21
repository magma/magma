import * as SRD from "@projectstorm/react-diagrams";
import { ForkNode } from "./ForkNode";
import { ForkNodeModel } from "./ForkNodeModel";
import * as React from "react";

export class ForkNodeFactory extends SRD.AbstractNodeFactory {
  constructor() {
    super("fork");
  }

  generateReactWidget(
    diagramEngine: SRD.DiagramEngine,
    node: SRD.NodeModel
  ): JSX.Element {
    return <ForkNode node={node} />;
  }

  getNewInstance() {
    return new ForkNodeModel();
  }
}
