import * as SRD from "@projectstorm/react-diagrams";
import { CircleEndNodeModel } from "./CircleEndNodeModel";
import * as React from "react";
import { CircleNodeEnd } from "./CircleNodeEnd";

export class CircleEndNodeFactory extends SRD.AbstractNodeFactory {
  constructor() {
    super("end");
  }

  generateReactWidget(
    diagramEngine: SRD.DiagramEngine,
    node: SRD.NodeModel
  ): JSX.Element {
    return <CircleNodeEnd node={node} diagramEngine={diagramEngine} />;
  }

  getNewInstance() {
    return new CircleEndNodeModel();
  }
}
