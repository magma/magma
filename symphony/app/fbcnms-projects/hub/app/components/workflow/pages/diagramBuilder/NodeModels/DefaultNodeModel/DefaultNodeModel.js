import { NodeModel, DiagramEngine, Toolkit } from "storm-react-diagrams";
import * as _ from "lodash";
import { DefaultPortModel } from "./DefaultPortModel";

export class DefaultNodeModel extends NodeModel {
  name: string;
  color: string;
  inputs: {};

  constructor(
    name: string = "Untitled",
    color: string = "rgb(0,192,255)",
    inputs: {}
  ) {
    super("default");
    this.name = name;
    this.color = color;
    super.extras = { inputs: inputs };

    this.addInPort("In");
    this.addOutPort("Out");
  }

  addInPort(label: string): DefaultPortModel {
    return this.addPort(new DefaultPortModel(true, Toolkit.UID(), label));
  }

  addOutPort(label: string): DefaultPortModel {
    return this.addPort(new DefaultPortModel(false, Toolkit.UID(), label));
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

  getInPorts(): DefaultPortModel[] {
    return _.filter(this.ports, portModel => {
      return portModel.in;
    });
  }

  getOutPorts(): DefaultPortModel[] {
    return _.filter(this.ports, portModel => {
      return !portModel.in;
    });
  }

  getInputs() {
    return this.inputs;
  }
}
