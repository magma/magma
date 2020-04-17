import * as SRD from "storm-react-diagrams";
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
    return <DecisionNode node={node} />;
  }

  getNewInstance() {
    return new DecisionNodeModel();
  }
}
