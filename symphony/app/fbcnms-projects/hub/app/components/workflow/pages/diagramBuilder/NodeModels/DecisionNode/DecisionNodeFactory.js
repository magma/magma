import * as SRD from "@projectstorm/react-diagrams";
import * as React from "react";
import { DecisionNode } from "./DecisionNode";
import { DecisionNodeModel } from "./DecisionNodeModel";

export class DecisionNodeFactory extends SRD.AbstractNodeFactory {
  constructor() {
    super("decision");
  }

  generateReactWidget(
    diagramEngine: SRD.DiagramEngine,
    node: SRD.NodeModel
  ): JSX.Element {
    return <DecisionNode node={node} diagramEngine={diagramEngine} />;
  }

  getNewInstance() {
    return new DecisionNodeModel();
  }
}
