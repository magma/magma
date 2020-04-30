import * as SRD from "@projectstorm/react-diagrams";
import { JoinNode } from "./JoinNode";
import { JoinNodeModel } from "./JoinNodeModel";
import * as React from "react";

export class JoinNodeFactory extends SRD.AbstractNodeFactory {
  constructor() {
    super("join");
  }

  generateReactWidget(
    diagramEngine: SRD.DiagramEngine,
    node: SRD.NodeModel
  ): JSX.Element {
    return <JoinNode node={node} />;
  }

  getNewInstance() {
    return new JoinNodeModel();
  }
}
