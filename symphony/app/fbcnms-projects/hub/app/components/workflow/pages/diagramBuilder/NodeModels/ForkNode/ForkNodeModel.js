import { NodeModel } from "@projectstorm/react-diagrams";
import { ForkNodePortModel } from "./ForkNodePortModel";
import * as _ from "lodash";
import { DefaultPortModel } from "@projectstorm/react-diagrams";
import { DiagramEngine } from "@projectstorm/react-diagrams";

export class ForkNodeModel extends NodeModel {
  name: string;
  color: string;
  inputs: {};

  constructor(
    name: string = "Untitled",
    color: string = "rgb(0,192,255)",
    inputs: {}
  ) {
    super("fork");
    this.name = name;
    this.color = color;
    super.extras = { inputs: inputs };

    this.addPort(new ForkNodePortModel(true, "left"));
    this.addPort(new ForkNodePortModel(false, "right"));
  }

  deSerialize(object, engine: DiagramEngine) {
    super.deSerialize(object, engine);
    this.name = object.name;
    this.color = object.color;
  }

  serialize() {
    return _.merge(super.serialize(), {
      name: this.name,
      color: this.color
    });
  }

  getInputs() {
    return this.inputs;
  }
}
